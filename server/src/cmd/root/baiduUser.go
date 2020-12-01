package main

import (
	"fmt"
	"net/http"
	"porter/define"
	"porter/requester"
	"porter/wlog"
	"time"

	"github.com/bitly/go-simplejson"
	"gorm.io/gorm/clause"
)

type BaiduUser struct {
	UID        string `gorm:"primaryKey"`
	Username   string // 账户名称
	Nickname   string // 昵称
	Bduss      string `gorm:"primaryKey"`
	FansNum    int
	Diamond    int
	videoCount int
	DouyinUID  string // 绑定的抖音uid
}

func NewBaiduUser(bduss string) (*BaiduUser, error) {
	b := &BaiduUser{Bduss: bduss}
	err := b.fetchUsernInfo()

	return b, err
}
func (b *BaiduUser) fetchUsernInfo() error {
	client := requester.NewHTTPClient()
	cookie := http.Cookie{Name: "BDUSS", Value: b.Bduss, Expires: time.Now().Add(180 * 24 * time.Hour)}
	req, err := http.NewRequest("GET", define.GetBaiduBaseInfo, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %s", err)
	}
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

	// 再获取全民视频相关数据
	quanminReq, err := http.NewRequest("GET", define.GetQuanminInfo, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %s", err)
	}

	quanminReq.AddCookie(&cookie)

	quanminResp, err := client.Do(quanminReq)
	if err != nil {
		return fmt.Errorf("获取百度账号的基本信息出错: %s", err)
	}
	defer quanminResp.Body.Close()

	quanminJ, err := simplejson.NewFromReader(quanminResp.Body)
	if err != nil {
		return fmt.Errorf("解析全民视频数据失败: %s", err)
	}

	errCode, err = quanminJ.Get("errno").Int()
	if err != nil {
		return fmt.Errorf("接口返回值错误: %s", err)
	}
	errMsg, err = quanminJ.Get("errmsg").String()
	if err != nil {
		return fmt.Errorf("接口返回消息错误: %s", err)
	}
	if errCode != 0 {
		return fmt.Errorf("请求获取全民账号数据失败: %d, 消息: %s", errCode, errMsg)
	}

	dataJ := quanminJ.Get("data")
	nickname, err := dataJ.Get("name").String()
	if err != nil {
		return fmt.Errorf("解析全民视频nickname失败: %s", err)
	}
	fansnum, err := dataJ.Get("fans_num").Int()
	if err != nil {
		return fmt.Errorf("解析全民视频粉丝数量失败: %s", err)
	}
	diamond, err := dataJ.Get("points").Int()
	if err != nil {
		return fmt.Errorf("解析全民视频钻石数量失败: %s", err)
	}
	videoCount, err := dataJ.Get("video_count").Int()
	if err != nil {
		return fmt.Errorf("解析全民视频钻石数量失败: %s", err)
	}
	b.Nickname = nickname
	b.FansNum = fansnum
	b.Diamond = diamond
	b.videoCount = videoCount

	return nil
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
