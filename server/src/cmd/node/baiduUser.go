package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"porter/define"
	"porter/requester"
	"porter/task"
	"porter/util"
	"porter/wlog"
	"strconv"
	"strings"
	"time"

	"github.com/go-errors/errors"
)

var (
	baiduComURL = &url.URL{
		Scheme: "http",
		Host:   "baidu.com",
	}
	cookieDomain = ".baidu.com"
	uploadURL    = "https://quanmin.baidu.com/web/publish/upload"
)

type BaiduClient struct {
	Nickname string
	BDUSS    string // 百度BDUSS

	k  string
	v  string
	sv string

	client *requester.HTTPClient
}

// 将抖音那边的头像,昵称,签名等数据同步到百度
func (b *BaiduClient) SyncFromDouyin(item *task.TaskChangeInfoItem) error {
	wlog.Infof("开始从抖音用户[%s]复制用户信息", item.Nickname)
	// 更换头像
	if err := b.Setportrait(item.Avatar); err != nil {
		wlog.Infof("[%s]设置头像失败: %s", item.Nickname, err)
	} else {
		wlog.Infof("[%s]设置头像成功", item.Nickname)
	}

	// 昵称和签名关键字过滤后进行更换
	newName := filterKeyword(addExtraName(util.FilterSpecial(item.Nickname)))
	err := b.SetProfile(map[string]string{
		"nickname": newName,
	})
	if err != nil {
		wlog.Infof("[%s]设置昵称[%s]失败: %s", item.Nickname, newName, err)
	} else {
		wlog.Infof("[%s]设置昵称为[%s]成功", item.Nickname, newName)
	}

	autograph := filterKeyword(util.FilterSpecial(item.Signature))
	err = b.SetProfile(map[string]string{
		"autograph": autograph,
	})
	if err != nil {
		wlog.Infof("[%s]设置签名[%s]失败: %s", item.Nickname, autograph, err)
	} else {
		wlog.Infof("[%s]设置签名成功", item.Nickname)
	}

	// 性别 sex=1&user_type=ugc 1是女 2是男 0是未知
	if item.Gender != 0 {
		// 抖音和百度的性别设置是反的
		sex := "1"
		if item.Gender == 1 {
			sex = "2"
		}
		err = b.SetProfile(map[string]string{
			"sex": sex,
		})
		if err != nil {
			wlog.Infof("[%s]设置性别[%s]失败: %s", item.Nickname, sex, err)
		} else {
			wlog.Infof("[%s]设置性别成功", item.Nickname)
		}
	}

	// 生日birthday=19830104&user_type=ugc
	if item.Birthday != "" {
		birthday := strings.ReplaceAll(item.Birthday, "-", "")
		err = b.SetProfile(map[string]string{
			"birthday": birthday,
		})
		if err != nil {
			wlog.Infof("[%s]设置生日[%s]失败: %s", item.Nickname, birthday, err)
		} else {
			wlog.Infof("[%s]设置生日成功", item.Nickname)
		}
	}

	if item.Province != "" || item.City != "" {
		var cityCode int
		if item.City != "" {
			cityName := strings.ReplaceAll(item.City, "市", "")
			cityCode = Area[cityName]
		} else {
			provinceName := strings.ReplaceAll(item.Province, "省", "")
			cityCode = Area[provinceName]
		}

		err = b.SetProfile(map[string]string{
			"city": strconv.Itoa(cityCode),
		})
		if err != nil {
			wlog.Infof("[%s]设置城市[%d]失败: %s", item.Nickname, item.Location, err)
		} else {
			wlog.Infof("[%s]设置城市成功", item.Nickname)
		}
	}

	wlog.Infof("[%s]信息复制完毕", item.Nickname)

	return nil
}

func (b *BaiduClient) Setportrait(imageURL string) error {
	if imageURL == "" {
		return errors.New("头像链接不能为空")
	}
	resp, err := b.client.Req(http.MethodGet, imageURL, nil, nil)
	if err != nil {
		return fmt.Errorf("获取头像数据失败:%s", err)
	}
	headData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取头像数据失败:%s", err)
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	m := make(map[string]string)
	m["client"] = "android"
	m["cuid"] = "133B4B0017CB08044A8983B54D34C3F5"
	m["clientid"] = "133B4B0017CB08044A8983B54D34C3F5"
	m["zid"] = "VhWdd6aWoVz4pg6CYjp1yysdG8r69POjHcxLHy1pyQ7Cwbj8kVb_mEkhPqKH3ASpKY_6e9v6QxAEYmDzuhNA5FA"
	m["clientip"] = "10.0.2.15"
	m["appid"] = "1"
	m["tpl"] = "bdmv"
	m["app_version"] = "2.3.2.10"
	m["sdk_version"] = "8.9.3"
	m["sdkversion"] = "8.9.3"
	m["bduss"] = b.BDUSS
	m["portrait_type"] = "0"
	for k, v := range m {
		writer.WriteField(k, v)
	}
	sig := calculateSig(m, "0c7c877d1c7825fa4438c44dbb645d1b")
	writer.WriteField("sig", sig)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("file"), escapeQuotes("portrait.jpg")))
	h.Set("Content-Type", "image/jpeg")

	part1, _ := writer.CreatePart(h)
	part1.Write(headData)
	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, define.SetPortrait, payload)

	if err != nil {
		return err
	}

	req.Header.Add("XRAY-TRACEID", "ec8f1b76-3f3c-4e65-936b-a0a26ffaf0c8")
	req.Header.Add("XRAY-REQ-FUNC-ST-DNS", "httpsUrlConn;"+strconv.FormatInt(time.Now().UnixNano()/10e5, 10)+";5")
	req.Header.Add("Content-Type", "multipart/form-data;boundary="+writer.Boundary())
	req.Header.Add("User-Agent", "tpl:bdmv;android_sapi_v8.9.3")
	req.Header.Add("Host", "passport.baidu.com")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Content-Length", strconv.Itoa(len(payload.String())))
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	result := &struct {
		Errno  int    `json:"errno"`
		Errmsg string `json:"errmsg"`
	}{}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	if result.Errno != 0 {
		return errors.New(result.Errmsg)
	}

	return nil
}

// 昵称, 签名, 生日, 性别等等都是由这个接口来完成
func (b *BaiduClient) SetProfile(args map[string]string) error {
	// https://quanmin.baidu.com/mvideo/api?api_name=userprofilesubmit&nickname=超级码力
	url := fmt.Sprintf("%s?api_name=userprofilesubmit", define.QuanminAPI)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", requester.UserAgent)

	q := req.URL.Query()

	for k, v := range args {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	result := &struct {
		Timestamp         int    `json:"timestamp"`
		Logid             string `json:"logid"`
		ServLogin         bool   `json:"servLogin"`
		Userprofilesubmit struct {
			Status int    `json:"status"`
			Msg    string `json:"msg"`
		} `json:"userprofilesubmit"`
	}{}

	err = json.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("解析错误:%s", err)
	}

	if result.Userprofilesubmit.Status != 0 {
		return errors.New(result.Userprofilesubmit.Msg)
	}

	return nil
}

func NewBaiduClient(bduss string) *BaiduClient {
	b := &BaiduClient{
		BDUSS:  bduss,
		client: requester.NewHTTPClient(),
	}
	b.client.Jar.SetCookies(baiduComURL, []*http.Cookie{
		{
			Name:   "BDUSS",
			Value:  bduss,
			Domain: cookieDomain,
		},
	})

	return b
}
