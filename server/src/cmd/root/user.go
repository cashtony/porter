package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"porter/define"
	"porter/requester"
	"porter/wlog"
	"time"

	"github.com/bitly/go-simplejson"
)

// NewDouYinUser 传入一个分享的url,类似https://v.douyin.com/qKDMXG/
func NewDouYinUser(shareURL string) (*DouyinUser, error) {
	// 获取个人信息
	user := &DouyinUser{
		ShareURL: shareURL,
	}
	err := user.getDouyinUserInfo()
	if err != nil {
		return user, err
	}

	// 获取所有视频信息
	user.getDouYinVideoList()

	wlog.Infof("user %+v", user)
	return user, nil
}

func NewBaiduUser(bduss string) (user *BaiduUser) {
	// 获取个人信息
	return
}

func (u *DouyinUser) getDouyinUserInfo() error {
	resp, err := requester.DefaultClient.Req("GET", u.ShareURL, nil, nil)
	if err != nil {
		return fmt.Errorf("访问失败: %s %s", u.ShareURL, err)
	}
	defer resp.Body.Close()

	u.sec_uid = resp.Request.URL.Query().Get("sec_uid")

	infoReq := fmt.Sprintf("%s?sec_uid=%s", define.GetUserInfo, u.sec_uid)
	infoResp, err := requester.DefaultClient.Req("GET", infoReq, nil, nil)
	if err != nil {
		return fmt.Errorf("获取用户数据时失败:%s %s", infoReq, err)
	}
	defer infoResp.Body.Close()

	j, err := simplejson.NewFromReader(infoResp.Body)
	if err != nil {
		return fmt.Errorf("解析用户json数据失败: %s", err)
	}

	jUser := j.Get("user_info")
	u.Nickname = jUser.Get("nickname").MustString()
	u.UniqueUID = jUser.Get("unique_id").MustString()
	u.UID = jUser.Get("uid").MustString()

	return nil
}

func (u *DouyinUser) getDouYinVideoList() {
	u.videoList = make([]*DouyinVideo, 0)
	var (
		maxCursor int64 = 0
		hasMore         = true
		page            = 1
	)
	for hasMore {
		time.Sleep(5 * time.Millisecond)

		url := fmt.Sprintf("%s?user_id=%s&sec_uid=%s&count=35&max_cursor=%d&aid=1128&_signature=&dytk=", define.GetVideoList, u.UID, u.sec_uid, maxCursor)
		resp, err := requester.DefaultClient.Req("GET", url, nil, nil)
		if err != nil {
			continue
		}

		if resp.Header.Get("status_code") == "" {
			resp.Body.Close()
			continue
		}

		j, err := simplejson.NewFromReader(resp.Body)
		if err != nil {
			wlog.Error("数据解析失败:", err)
			resp.Body.Close()
			return
		}

		wlog.Debugf("开始解析[%s]第[%d]页视频 \n", u.Nickname, page)
		for _, item := range j.Get("aweme_list").MustArray() {
			itemJSON := item.(map[string]interface{})
			// downloadURL := itemJSON["video"].(map[string]interface{})["play_addr"].(map[string]interface{})["url_list"].([]interface{})[0].(string)
			authorInfo := itemJSON["author"].(map[string]interface{})
			video := &DouyinVideo{
				AuthorUID: authorInfo["uid"].(string),
				AwemeID:   itemJSON["aweme_id"].(string),
				Desc:      itemJSON["desc"].(string),
			}

			//获取视频上传时间
			video.CreateTime, err = getVideoCreateTime(video.AwemeID)
			if err != nil {
				wlog.Error("视频的创建时间获取错误", err)
				continue
			}
			u.videoList = append(u.videoList, video)
			wlog.Debugf("解析视频成功: %+v \n", video)
		}

		hasMore, err = j.Get("has_more").Bool()
		if err != nil {
			wlog.Error("获取has_more字段错误", err)
		}
		if hasMore {
			maxCursor, err = j.Get("max_cursor").Int64()
			if err != nil {
				wlog.Error("获取max_cursor字段错误", err)
			}
		}

		resp.Body.Close()
		page++
	}
}

func getVideoCreateTime(awemeid string) (time.Time, error) {
	url := fmt.Sprintf("%s?item_ids=%s", define.GetVideoURI, awemeid)

	tryTimes := 10
	var resp *http.Response
	var err error
	for tryTimes > 0 {
		resp, err = requester.DefaultClient.Req("GET", url, nil, nil)
		if err != nil {
			if err.Error() == "TLS handshake timeout" {
				tryTimes--
				continue
			}

			return time.Now(), fmt.Errorf("请求视频信息失败:%s", err)
		}
	}
	if resp == nil {
		return time.Now(), errors.New("请求视频信息失败: 超过重试次数")
	}
	defer resp.Body.Close()

	j, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return time.Now(), fmt.Errorf("数据解析失败:%s", err)
	}

	list, err := j.Get("item_list").Array()
	if err != nil {
		return time.Now(), fmt.Errorf("解析item_list字段失败:%s", err)
	}

	if len(list) == 0 {
		return time.Now(), fmt.Errorf("获取item_list数据长度为0:%s", err)
	}

	t := list[0].(map[string]interface{})["create_time"].(json.Number)
	timeStamp, err := t.Int64()
	if err != nil {
		return time.Now(), fmt.Errorf("视频的createTime字段获取失败:%s", err)
	}
	return time.Unix(timeStamp, 0), nil

}
func (v *DouyinVideo) getVideoInfo() {

}
func (v *DouyinVideo) getVideoCreateTime() {

}

// 查看是否有新的视频数据,只需查看第一页即可
func UpdateDouYinUser(uid string) {

}

func StoreDouYinUser(user *DouyinUser) {
	// 储存到数据库
	DB.Create(&user)
	for _, v := range user.videoList {
		DB.Create(v)
	}
}
