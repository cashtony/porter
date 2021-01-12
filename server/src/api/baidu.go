package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"porter/define"
	"porter/requester"
	"porter/wlog"
	"time"
)

type ApiQuanminUser struct {
	Mine struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
		Data   struct {
			User struct {
				ID        string `json:"id"`
				UserName  string `json:"userName"`
				Sex       string `json:"sex"`
				Age       string `json:"age"`
				Area      string `json:"area"`
				Autograph string `json:"autograph"`
				HeadImg   string `json:"headImg"`
				FansNum   int    `json:"fansNum"`
			} `json:"user"`
			Charm struct {
				CharmpointsNumber     int `json:"charmpointsNumber"`
				AvailablePointsNumber int `json:"availablePointsNumber"`
			}
		} `json:"data,omitempty"`
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
		wlog.Info("无法解析此数据:", string(data))
		return nil, err
	}

	if user.Mine.Status != 0 {
		return nil, fmt.Errorf("没有获取到全民账户数据:%s", user.Mine.Msg)
	}

	return user, nil
}

type APIQuanminSearch struct {
	NewTabSearch struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
		Data   struct {
			AuthorList []struct {
				Content struct {
					IsUserSelf int    `json:"isUserSelf"`
					Nickname   string `json:"nickname"`
					AuthorID   string `json:"authorId"`
					Daren      int    `json:"daren"` // 0是未认证
				} `json:"content"`
			} `json:"author_list"`
		} `json:"data"`
	} `json:"newTabSearch"`
}
