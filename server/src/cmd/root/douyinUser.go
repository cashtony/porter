package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"porter/define"
	"porter/requester"
	"porter/wlog"
	"time"

	"github.com/bitly/go-simplejson"
	"gorm.io/gorm/clause"
)

type DouyinUser struct {
	UID             string          `json:"uid" gorm:"primaryKey"`
	UniqueUID       string          `json:"uniqueUID" gorm:"primaryKey"` // 抖音号
	Nickname        string          `json:"nickName"`
	ShareURL        string          `json:"shareURL" gorm:"primaryKey"`
	VideoCount      int             `json:"videoCount"`
	FansCount       int             `json:"fansNum"`
	LastCollectTime define.JsonTime `json:"lastCollectTime"` // 最后一次采集时间
	CreateTime      define.JsonTime `json:"createTime" gorm:"default:now()"`

	secUID string // 用于填充获取用户数据接口
}

// NewDouYinUser 传入一个分享的url,类似https://v.douyin.com/qKDMXG/
func NewDouYinUser(shareURL string) (*DouyinUser, error) {
	// 获取个人信息
	user := &DouyinUser{
		ShareURL: shareURL,
	}
	err := user.initUserInfo()
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *DouyinUser) initSecID() {
	resp, err := requester.DefaultClient.Req("GET", u.ShareURL, nil, nil)
	if err != nil {
		wlog.Error("获取secid失败: %s %s", u.ShareURL, err)
		return
	}
	defer resp.Body.Close()

	u.secUID = resp.Request.URL.Query().Get("sec_uid")
}

func (u *DouyinUser) initUserInfo() error {
	u.initSecID()

	infoReq := fmt.Sprintf("%s?sec_uid=%s", define.GetUserInfo, u.secUID)
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
	u.Nickname, _ = jUser.Get("nickname").String()
	u.UniqueUID, _ = jUser.Get("unique_id").String()
	u.UID, _ = jUser.Get("uid").String()
	u.VideoCount, _ = jUser.Get("aweme_count").Int()
	u.FansCount, _ = jUser.Get("follower_count").Int()
	return nil
}

func (u *DouyinUser) initVideoList() {

	var (
		onePageList []*DouyinVideo
		nextCursor  int64 = 0
		hasMore           = true
		page              = 1
		err         error
	)

	for hasMore {
		wlog.Debugf("开始解析[%s]第[%d]页视频 \n", u.Nickname, page)
		onePageList, hasMore, nextCursor, err = u.OnePageVideo(nextCursor)
		if err != nil {
			wlog.Error("获取单页视频发生错误:", err)
			return
		}

		u.StoreVideo(onePageList)
		wlog.Debugf("[%s]第[%d]页视频解析完毕 hasmore:%b nextCursor:%d videoLen:%d \n", u.Nickname, page, hasMore, nextCursor, len(onePageList))
		page++
	}

	wlog.Debugf("用户[%s]视频解析完毕 \n", u.Nickname)
}

// OnePageVideo 接收一个游标,返回是否有下一页以及相关游标
func (u *DouyinUser) OnePageVideo(cursor int64) ([]*DouyinVideo, bool, int64, error) {
	videoList := make([]*DouyinVideo, 0)
	var (
		nextCursor int64 = 0
		hasMore          = false
		tryTimes         = 0
	)

	for {
		time.Sleep(100 * time.Millisecond)
		if tryTimes > 500 {
			wlog.Infof("[警告]获取视频列表尝试超过%d次仍然没有获得数据 \n", tryTimes)
			tryTimes = 0
		}

		url := fmt.Sprintf("%s?user_id=%s&sec_uid=%s&count=20&max_cursor=%d&aid=1128&_signature=&dytk=", define.GetVideoList, u.UID, u.secUID, cursor)
		resp, err := requester.DefaultClient.Req("GET", url, nil, nil)
		if err != nil {
			tryTimes++
			continue
		}

		if resp.Header.Get("status_code") == "" {
			resp.Body.Close()
			tryTimes++
			continue
		}

		defer resp.Body.Close()

		j, err := simplejson.NewFromReader(resp.Body)
		if err != nil {
			return videoList, false, 0, fmt.Errorf("json数据解析失败:%s", err)
		}

		for _, item := range j.Get("aweme_list").MustArray() {
			itemJSON := item.(map[string]interface{})
			authorInfo := itemJSON["author"].(map[string]interface{})
			video := &DouyinVideo{
				AuthorUID: authorInfo["uid"].(string),
				AwemeID:   itemJSON["aweme_id"].(string),
				Desc:      itemJSON["desc"].(string),
			}

			//获取视频上传时间
			videoExtraInfo, _ := getVideoCreateTime(video.AwemeID)

			video.CreateTime = videoExtraInfo.CreateTime
			videoList = append(videoList, video)
		}

		hasMore, err = j.Get("has_more").Bool()
		if err != nil {
			wlog.Error("获取has_more字段错误", err)
		}
		if hasMore {
			nextCursor, err = j.Get("max_cursor").Int64()
			if err != nil {
				wlog.Error("获取max_cursor字段错误", err)
			}
		}

		if hasMore && nextCursor == 0 {
			hasMore = false
		}

		break
	}

	return videoList, hasMore, nextCursor, nil
}

func (d *DouyinUser) Store() {
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}},
		UpdateAll: true,
	}).Create(d)
	if DB.Error != nil {
		wlog.Errorf("抖音用户[%s][%s]存入数据库失败:%s \n", d.UID, d.Nickname, DB.Error)
		return
	}
}
func (d *DouyinUser) StoreVideo(list []*DouyinVideo) {
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "aweme_id"}},
		DoNothing: true,
	}).Create(list)
	if DB.Error != nil {
		wlog.Errorf("抖音用户[%s][%s]存入数据库失败:%s \n", d.UID, d.Nickname, DB.Error)
		return
	}
}

func (d *DouyinUser) Update() {
	onePageList, _, _, err := d.OnePageVideo(0)
	if err != nil {
		wlog.Errorf("用户[%s][%s]获取视频列表失败:%s \n", d.UID, d.Nickname, err)
		return
	}

	// 将没有的视频传入到数据库中
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "aweme_id"}},
		DoNothing: true,
	}).Create(onePageList)
	if DB.Error != nil {
		wlog.Errorf("用户[%s][%s]新视频信息存入数据库失败:%s \n", d.UID, d.Nickname, DB.Error)
		return
	}

	//更新用户的last_collect_time字段
	DB.Model(d).Update("last_collect_time", time.Now())
	if DB.Error != nil {
		wlog.Errorf("从数据库中更新用户[%s][%s]last_collect_time字段失败: %s \n", d.UID, d.Nickname, DB.Error)
		return
	}
}

func getVideoCreateTime(awemeid string) (*DouyinVideoExtraInfo, error) {
	info := &DouyinVideoExtraInfo{}
	url := fmt.Sprintf("%s?item_ids=%s", define.GetVideoURI, awemeid)

	tryTimes := 10
	var resp *http.Response
	var err error
	for tryTimes > 0 {
		// 设置间隔为了防止两次调用时间间隔过短导致握手失败
		time.Sleep(500 * time.Millisecond)
		resp, err = requester.DefaultClient.Req("GET", url, nil, nil)
		if err != nil {
			tryTimes--
			continue
		}
		break
	}
	if resp == nil {
		return info, fmt.Errorf("[%s]请求视频信息失败: 超过重试次数", awemeid)
	}
	defer resp.Body.Close()

	j, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return info, fmt.Errorf("[%s]数据解析失败:%s", awemeid, err)
	}

	list, err := j.Get("item_list").Array()
	if err != nil {
		return info, fmt.Errorf("[%s]解析item_list字段失败:%s", awemeid, err)
	}

	if len(list) == 0 {
		// 为0说明已经失效,如果视频主删掉了视频就会导致此结果
		DB.Model(&DouyinVideo{}).Where("aweme_id = ?", awemeid).Update("state", 2)
		return info, fmt.Errorf("[%s]此视频可能已经被清理,将在数据库中进行标记", awemeid)
	}
	videoJSONInfo := list[0].(map[string]interface{})
	t := videoJSONInfo["create_time"].(json.Number)
	timeStamp, err := t.Int64()
	if err != nil {
		return info, fmt.Errorf("[%s]视频的createTime字段获取失败:%s", awemeid, err)
	}

	info.CreateTime = time.Unix(timeStamp, 0)
	info.VID = videoJSONInfo["video"].(map[string]interface{})["vid"].(string)

	return info, nil

}
