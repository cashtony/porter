package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"porter/define"
	"porter/requester"
	"time"
)

type ApiQuanminUser struct {
	Mine struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
		Data   struct {
			User struct {
				ID       string `json:"id"`
				UserName string `json:"userName"`
				HeadImg  string `json:"headImg"`
				FansNum  int    `json:"fansNum"`
			} `json:"user"`
			Charm struct {
				CharmpointsNumber     int `json:"charmpointsNumber"`
				AvailablePointsNumber int `json:"availablePointsNumber"`
			}
		} `json:"data"`
	} `json:"mine"`
}

func GetQuanminInfo(bduss string) (*ApiQuanminUser, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, define.GetQuanminInfoV2, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %s", err)
	}
	req.Header.Add("User-Agent", requester.UserAgent)
	cookie := http.Cookie{Name: "BDUSS", Value: bduss, Expires: time.Now().Add(180 * 24 * time.Hour)}
	req.AddCookie(&cookie)

	resq, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取全民账号的基本信息出错: %s", err)
	}
	defer resq.Body.Close()

	data, err := ioutil.ReadAll(resq.Body)
	if err != nil {
		return nil, err
	}

	user := &ApiQuanminUser{}
	if err := json.Unmarshal(data, user); err != nil {
		return nil, err
	}

	return user, nil
}
