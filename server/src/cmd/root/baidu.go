package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"porter/api"
	"porter/define"
	"porter/requester"
	"porter/task"
	"porter/wlog"
	"time"

	"github.com/bitly/go-simplejson"
	"gorm.io/gorm/clause"
)

type TableBaiduUser struct {
	UID            string          `gorm:"primaryKey" json:"uid"`
	Username       string          `json:"username"` // 账户名称
	Nickname       string          `json:"nickname"`
	Bduss          string          `gorm:"primaryKey" json:"dbuss"`
	FansNum        int             `json:"fansNum"`
	Diamond        int             `json:"diamond"`
	Sex            string          `json:"sex"`
	Age            string          `json:"age"`
	Area           string          `json:"area"`
	Autograph      string          `json:"autograph"`
	VideoCount     int             `json:"videoCount"`
	DouyinUID      string          `json:"douyinUID"` // 绑定的抖音uid
	CreateTime     define.JsonTime `gorm:"default:now()" json:"createTime"`
	LastUploadTime time.Time       `json:"lastUploadTime"`
	Status         int             `json:"status" gorm:"default:1"` // 1:正常, 0:不搬运视频
}

func (TableBaiduUser) TableName() string {
	return "baidu_users"
}

func (b *TableBaiduUser) fetcBaseInfo() error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, define.GetBaiduBaseInfo, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %s", err)
	}
	req.Header.Add("User-Agent", requester.UserAgent)
	cookie := http.Cookie{Name: "BDUSS", Value: b.Bduss, Expires: time.Now().Add(180 * 24 * time.Hour)}
	req.AddCookie(&cookie)

	// 先获取基本数据
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("获取百度账号的基本信息出错: %s", err)
	}
	defer resp.Body.Close()

	j, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("解析百度账号数据失败: %s", err)
	}
	errCode, err := j.Get("errno").Int()
	if err != nil {
		return fmt.Errorf("接口返回值错误: %s", err)
	}
	errMsg, err := j.Get("show_msg").String()
	if err != nil {
		return fmt.Errorf("接口返回消息错误: %s", err)
	}

	if errCode == -6 {
		return fmt.Errorf("bduss错误: %d, 消息: %s", errCode, errMsg)
	}

	if errCode != 0 {
		return fmt.Errorf("请求获取百度账号数据失败: %d, 消息: %s", errCode, errMsg)
	}

	username, err := j.Get("login_info").Get("username").String()
	if err != nil {
		return fmt.Errorf("解析百度username失败: %s", err)
	}
	uid, err := j.Get("login_info").Get("uk_str").String()
	if err != nil {
		return fmt.Errorf("解析百度username失败: %s", err)
	}

	b.UID = uid
	b.Username = username

	return nil
}

func (b *TableBaiduUser) getVideoList(num int, mode GetMode) ([]*define.TableDouyinVideo, error) {
	page, limit := 1, 20

	videoList := make([]*define.TableDouyinVideo, 0)
	for {
		tmpList := make([]*define.TableDouyinVideo, 0)
		subDB := DB.Model(&define.TableDouyinVideo{}).Where("douyin_uid = ? and state = ?", b.DouyinUID, WaitUpload).Order("create_time desc")
		if mode == GetModeNewly {
			subDB.Where("date(create_time) >= current_date - 1")
		}
		subDB.Offset((page - 1) * limit).Limit(limit).Find(&tmpList)
		if subDB.Error != nil {
			return nil, subDB.Error
		}

		page++

		if len(tmpList) == 0 {
			return videoList, nil
		}

		for _, v := range tmpList {
			// 检测视频是否还有效
			info, err := api.GetVideoExtraInfo(v.AwemeID)
			if err != nil {
				continue
			}
			if len(info.ItemList) == 0 {
				// 数据库中记录
				DB.Model(&define.TableDouyinVideo{}).Where("aweme_id = ?", v.AwemeID).Update("state", 2)
				continue
			}
			// 之前有些视频vid是没有储存的,所以需要在这里获取一次
			v.Vid = info.ItemList[0].Video.Vid

			videoList = append(videoList, v)
			if len(videoList) >= num {
				return videoList, nil
			}
		}
	}
}

func (b *TableBaiduUser) UploadVideo(utype UploadType) {
	uploadVideoList := make([]*define.TableDouyinVideo, 0)
	num := rand.Intn(MaxUploadNum-MinUploadNum) + MinUploadNum
	uploadVideoList, err := b.getVideoList(num, GetModeNewly)
	if err != nil {
		wlog.Errorf("从数据库中获取用户[%s][%s]最新视频列表信息失败:%s \n", b.UID, b.Nickname, err)
		return
	}

	// 没有新视频并且今天没有更新过视频才从以前的视频中提取
	if utype == UploadTypeDaily && len(uploadVideoList) == 0 && b.LastUploadTime.Day() != time.Now().Day() {
		wlog.Infof("用户[%s][%s]绑定的抖音号昨天没有更新,将获取以前的视频 \n", b.UID, b.Nickname)
		uploadVideoList, err = b.getVideoList(num, GetModeOlder)
		if err != nil {
			wlog.Errorf("从数据库中获取用户[%s][%s]视频列表信息失败:%s \n", b.UID, b.Nickname, err)
			return
		}
	}

	if len(uploadVideoList) == 0 {
		wlog.Infof("用户[%s][%s]没有可更新内容,退出 \n", b.UID, b.Nickname)
		return
	}

	b.pulishTask(uploadVideoList)
}

func (b *TableBaiduUser) pulishTask(uploadVideoList []*define.TableDouyinVideo) {
	// 查找视频下载url
	taskVideoList := make([]*task.TaskUploadVideo, 0)
	statisticList := make([]*Statistic, 0)

	for _, v := range uploadVideoList {
		taskVideoList = append(taskVideoList, &task.TaskUploadVideo{
			AwemeID:     v.AwemeID,
			Desc:        v.Desc,
			DownloadURL: fmt.Sprintf("%s/?video_id=%s&ratio=720p&line=0", define.GetVideoDownload, v.Vid),
		})

		statisticList = append(statisticList, &Statistic{
			BaiduUID:  b.UID,
			DouyinUID: b.DouyinUID,
			AwemeID:   v.AwemeID,
			State:     WaitUpload,
		})
	}

	if len(uploadVideoList) == 0 {
		wlog.Infof("用户[%s][%s]没有可更新内容,退出 \n", b.UID, b.Nickname)
		return
	}

	// 封装成task投递到任务队列中
	wlog.Debugf("开始投放用户[%s][%s]任务, 数量:[%d] \n", b.UID, b.Nickname, len(taskVideoList))
	t := &task.TaskUpload{
		Bduss:    b.Bduss,
		Videos:   taskVideoList,
		Nickname: b.Nickname,
	}

	data, err := json.Marshal(t)
	if err != nil {
		wlog.Error("task解析成json错误", err)
		return
	}

	// 增加数据统计
	DB.Create(&statisticList)

	err = Q.Publish(define.TaskPushTopic, data)
	if err != nil {
		wlog.Error("任务发布失败:", err)
	}

	//更新用户的最后上传字段
	result := DB.Model(b).Update("last_upload_time", time.Now())
	if result.Error != nil {
		wlog.Errorf("从数据库中更新用户[%s][%s]last_upload_time字段失败: %s \n", b.UID, b.Nickname, DB.Error)
		return
	}
}

func (b *TableBaiduUser) Store() {
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}},
		UpdateAll: true,
	}).Create(b)
	if DB.Error != nil {
		wlog.Errorf("用户[%s][%s]存入数据库失败:%s \n", b.UID, b.Nickname, DB.Error)
		return
	}
}

func NewBaiduUser(bduss string) (*TableBaiduUser, error) {
	b := &TableBaiduUser{Bduss: bduss}
	err := b.fetcBaseInfo()
	if err != nil {
		return nil, fmt.Errorf("获取基本百度用户基本数据失败:%s", err)
	}

	qmInfo, err := api.GetQuanminInfo(bduss)
	if err != nil {
		return nil, fmt.Errorf("获取全民用户基本数据失败:%s", err)
	}
	userInfo := qmInfo.Mine.Data.User
	charmInfo := qmInfo.Mine.Data.Charm

	b.Nickname = userInfo.UserName
	b.FansNum = userInfo.FansNum
	b.Sex = userInfo.Sex
	b.Area = userInfo.Area
	b.Autograph = userInfo.Autograph
	b.Age = userInfo.Age

	b.Diamond = charmInfo.CharmpointsNumber

	return b, nil
}

func UpdateBaiduUser(users []*TableBaiduUser) {
	for _, u := range users {
		qmInfo, err := api.GetQuanminInfo(u.Bduss)
		if err != nil {
			wlog.Errorf("获取[%s][%s]全民视频用户数据时错误:%s", u.UID, u.Nickname, err)
			continue
		}
		nickname := qmInfo.Mine.Data.User.UserName
		fansNum := qmInfo.Mine.Data.User.FansNum
		diamond := qmInfo.Mine.Data.Charm.CharmpointsNumber

		DB.Model(&TableBaiduUser{}).Where("uid = ?", u.UID).Updates(&TableBaiduUser{Diamond: diamond, Nickname: nickname, FansNum: fansNum})
	}
}
