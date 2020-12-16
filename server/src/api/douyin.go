package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"porter/define"
	"porter/requester"
	"porter/wlog"
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

	secUID string
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
	resp, err := requester.DefaultClient.Req("GET", infoReq, nil, nil)
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

	data.UserInfo.secUID = secUID

	return data.UserInfo, nil
}

func GetSecID(shareURL string) string {
	resp, err := requester.DefaultClient.Req("GET", shareURL, nil, nil)
	if err != nil {
		wlog.Error("获取secid失败: %s %s", shareURL, err)
		return ""
	}
	defer resp.Body.Close()

	return resp.Request.URL.Query().Get("sec_uid")
}
