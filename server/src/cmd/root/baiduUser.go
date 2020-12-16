package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"porter/define"
	"porter/requester"
	"porter/wlog"
	"regexp"
	"time"

	"github.com/bitly/go-simplejson"
	"gorm.io/gorm/clause"
)

type BaiduUser struct {
	UID            string          `gorm:"primaryKey" json:"uid"`
	Username       string          `json:"userName"` // 账户名称
	Nickname       string          `json:"nickName"`
	Bduss          string          `gorm:"primaryKey" json:"dbuss"`
	FansNum        int             `json:"fansNum"`
	Diamond        int             `json:"diamond"`
	Tbean          int             `json:"tbean"` // T豆
	VideoCount     int             `json:"videoCount"`
	DouyinURL      string          `json:"douyinURL"` // 绑定的抖音uid
	CreateTime     define.JsonTime `gorm:"default:now()" json:"createTime"`
	LastUploadTime time.Time       `json:"lastUploadTime"`
}

func NewBaiduUser(bduss string) (*BaiduUser, error) {
	b := &BaiduUser{Bduss: bduss}
	err := b.fetcBaseInfo()
	if err != nil {
		return nil, fmt.Errorf("获取基本百度用户基本数据失败:%s", err)
	}

	err = b.fetchQuanminInfo()
	if err != nil {
		return nil, fmt.Errorf("获取全民用户基本数据失败:%s", err)
	}

	return b, nil
}

func (b *BaiduUser) fetchQuanminInfo() error {
	cookie := http.Cookie{Name: "BDUSS", Value: b.Bduss, Expires: time.Now().Add(180 * 24 * time.Hour)}
	// 再获取全民视频相关数据
	quanminReq, err := http.NewRequest("POST", define.GetQuanminInfoV2, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %s", err)
	}

	quanminReq.AddCookie(&cookie)

	quanminResp, err := requester.DefaultClient.Do(quanminReq)
	if err != nil {
		return fmt.Errorf("获取百度账号的基本信息出错: %s", err)
	}
	defer quanminResp.Body.Close()

	quanminJ, err := simplejson.NewFromReader(quanminResp.Body)
	if err != nil {
		return fmt.Errorf("解析全民视频数据失败: %s", err)
	}

	mineJ := quanminJ.Get("mine")
	errCode, err := mineJ.Get("status").Int()
	if err != nil {
		return fmt.Errorf("接口返回值错误: %s", err)
	}

	if errCode != 0 {
		errMsg, _ := quanminJ.Get("msg").String()
		return fmt.Errorf("请求获取全民账号数据失败: %d, 消息: %s", errCode, errMsg)
	}

	data := mineJ.Get("data")
	diamond, err := data.Get("charm").Get("availablePointsNumber").Int()
	if err != nil {
		return fmt.Errorf("解析全民视频钻石数量获取失败: %s", err)
	}
	userJ := data.Get("user")
	nickname, err := userJ.Get("userName").String()
	if err != nil {
		return fmt.Errorf("解析全民视频nickname失败: %s", err)
	}
	fansnum, err := userJ.Get("fansNum").Int()
	if err != nil {
		return fmt.Errorf("解析全民视频粉丝数量失败: %s", err)
	}

	b.Nickname = nickname
	b.FansNum = fansnum
	b.Diamond = diamond

	return nil
}

func (b *BaiduUser) fetcBaseInfo() error {
	cookie := http.Cookie{Name: "BDUSS", Value: b.Bduss, Expires: time.Now().Add(180 * 24 * time.Hour)}
	req, err := http.NewRequest("GET", define.GetBaiduBaseInfo, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %s", err)
	}
	req.AddCookie(&cookie)

	// 先获取基本数据
	resp, err := requester.DefaultClient.Do(req)
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

var uploadDescReg, _ = regexp.Compile(`(抖音|dou|DOU)`)

func (b *BaiduUser) getVideoList(num int, mode GetMode) ([]*DouyinVideo, error) {
	page, limit := 1, 20

	videoList := make([]*DouyinVideo, 0)
	for {
		tmpList := make([]*DouyinVideo, 0)
		subDB := DB.Model(&DouyinVideo{}).Where("share_url = ? and state = ?", b.DouyinURL, WaitUpload).Order("create_time desc")
		if mode == GetModeNewly {
			subDB.Where(" date(create_time) >= current_date - 1")
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
			info, err := getVideoExtraInfo(v.AwemeID)
			if err != nil {
				continue
			}
			v.VID = info.VID

			videoList = append(videoList, v)
			if len(videoList) >= num {
				return videoList, nil
			}
		}
	}
}

// UploadNewOrOlderVideo 上传之前的老视频
func (b *BaiduUser) UploadVideo(utype UpdateType) {
	// 上传视频(从未上传的视频中挑选8-12条)
	uploadVideoList := make([]*DouyinVideo, 0)
	num := rand.Intn(MaxUploadNum-MinUploadNum) + MinUploadNum
	uploadVideoList, err := b.getVideoList(num, GetModeNewly)
	if err != nil {
		wlog.Errorf("从数据库中获取用户[%s][%s]最新视频列表信息失败:%s \n", b.UID, b.Nickname, DB.Error)
		return
	}

	if len(uploadVideoList) == 0 && utype == UpdateTypeDaily {
		wlog.Infof("用户[%s][%s]绑定的抖音号昨天没有更新,将获取以前的视频 \n", b.UID, b.Nickname)
		uploadVideoList, err = b.getVideoList(num, GetModeOlder)
		if err != nil {
			wlog.Errorf("从数据库中获取用户[%s][%s]视频列表信息失败:%s \n", b.UID, b.Nickname, DB.Error)
			return
		}
	}

	if len(uploadVideoList) == 0 {
		wlog.Infof("用户[%s][%s]绑定的抖音号没有可更新内容,退出 \n", b.UID, b.Nickname)
		return
	}

	b.pulishTask(uploadVideoList)
}

func (b *BaiduUser) pulishTask(uploadVideoList []*DouyinVideo) {
	// 查找视频下载url
	taskVideoList := make([]*define.TaskUploadVideo, 0)
	statisticList := make([]*Statistic, 0)

	for _, v := range uploadVideoList {
		taskVideoList = append(taskVideoList, &define.TaskUploadVideo{
			AwemeID:     v.AwemeID,
			Desc:        uploadDescReg.ReplaceAllString(v.Desc, ""),
			DownloadURL: fmt.Sprintf("%s/?video_id=%s&ratio=720p&line=0", define.GetVideoDownload, v.VID),
		})

		statisticList = append(statisticList, &Statistic{
			BaiduUID:  b.UID,
			DouyinURL: b.DouyinURL,
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
	t := &define.TaskUpload{
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

func (b *BaiduUser) Store() {
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}},
		UpdateAll: true,
	}).Create(b)
	if DB.Error != nil {
		wlog.Errorf("用户[%s][%s]存入数据库失败:%s \n", b.UID, b.Nickname, DB.Error)
		return
	}
}
