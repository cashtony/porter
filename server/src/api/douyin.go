package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"porter/define"
	"porter/requester"
	"porter/wlog"
	"time"
)

type APIDouyinUser struct {
	UID            string `json:"uid" gorm:"primaryKey"`
	UniqueUID      string `json:"unique_id" gorm:"primaryKey"` // 抖音号
	Nickname       string `json:"nickname"`
	AwemeCount     int    `json:"aweme_count"`
	FollowerCount  int    `json:"follower_count"`
	Signature      string `json:"signature" gorm:"-"` // 个人签名
	TotalFavorited string `json:"total_favorited"`
	AvatarMedium   struct {
		URLList []string `json:"url_list"`
	} `json:"avatar_medium"`

	SecUID string `json:"-"`
}

func NewAPIDouyinUser(shareURL string) (*APIDouyinUser, error) {
	if shareURL == "" {
		return nil, errors.New("shareURL不能为空")
	}
	secUID := GetSecID(shareURL)
	// 基本数据
	if secUID == "" {
		return nil, errors.New("secUID为空,不能获取抖音用户数据")
	}

	data := struct {
		StatusCode int            `json:"status_code"`
		UserInfo   *APIDouyinUser `json:"user_info"`
	}{}

	infoReq := fmt.Sprintf("%s?sec_uid=%s", define.GetUserInfo, secUID)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, infoReq, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %s", err)
	}
	req.Header.Add("User-Agent", requester.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取用户数据时失败:%s %s", infoReq, err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取response内容失败: %s", err)
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("解析用户json数据失败: %s", err)
	}

	if data.StatusCode != 0 {
		return nil, fmt.Errorf("返回的状态代码未成功: %d", data.StatusCode)
	}

	data.UserInfo.SecUID = secUID

	return data.UserInfo, nil
}

func GetSecID(shareURL string) string {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, shareURL, nil)
	if err != nil {
		wlog.Error("获取secid时创建请求失败:", err)
		return ""
	}
	req.Header.Add("User-Agent", requester.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		wlog.Error("获取secid失败: %s %s", shareURL, err)
		return ""
	}
	defer resp.Body.Close()

	return resp.Request.URL.Query().Get("sec_uid")
}

type APIDouyinVideo struct {
	AwemeID string `json:"aweme_id"`
	Desc    string `json:"desc"`
	Video   struct {
		Cover struct {
			URLList []string `json:"url_list"`
		} `json:"cover"`
		Duration int    `json:"duration"`
		Height   int    `json:"height"`
		Width    int    `json:"width"`
		VID      string `json:"vid"`
		PlayAddr struct {
			URLList []string `json:"url_list"`
		} `json:"play_addr"`
	} `json:"video"`
}

func GetDouyinVideo(secUID string, cursor int64) ([]*APIDouyinVideo, int64, error) {
	if secUID == "" {
		return nil, 0, errors.New("secUID不能为空")
	}

	data := struct {
		AwemeList  []*APIDouyinVideo `json:"aweme_list"`
		HasMore    interface{}       `json:"has_more"`
		MaxCursor  int64             `json:"max_cursor"`
		MinCursor  int64             `json:"min_cursor"`
		StatusCode int               `json:"status_code"`
	}{}

	tryTimes := 0
	url := fmt.Sprintf("%s?sec_uid=%s&count=20&max_cursor=%d&aid=1128&_signature=&dytk=", define.GetVideoList, secUID, cursor)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", requester.UserAgent)

	for {
		time.Sleep(500 * time.Millisecond)

		if tryTimes > 500 {
			wlog.Infof("[警告]获取视频列表尝试超过%d次仍然没有获得数据", tryTimes)
			tryTimes = 0
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			tryTimes++
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, err
		}

		if err := json.Unmarshal(body, &data); err != nil {
			return nil, 0, err
		}

		if len(data.AwemeList) == 0 && data.StatusCode == 0 && data.MaxCursor == 0 {
			tryTimes++
			continue
		}

		break
	}

	maxCursor := int64(0)
	hasmore := false
	switch data.HasMore.(type) {
	case bool:
		hasmore = data.HasMore.(bool)
	}

	if hasmore && data.MaxCursor != 0 {
		maxCursor = data.MaxCursor
	}

	return data.AwemeList, maxCursor, nil

}
