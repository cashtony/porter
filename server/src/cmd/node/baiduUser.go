package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"porter/api"
	"porter/define"
	"porter/requester"
	"porter/wlog"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

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

var keywordList, nameList []string

func init() {
	nameData := ReadFileData("./names.txt")
	if nameData == "" {
		wlog.Warn("名字前后缀词表为空")
	}

	nameList = strings.Split(nameData, "\r\n")

	wordData := ReadFileData("./keywords.txt")
	if wordData == "" {
		wlog.Warn("过滤关键词表为空")
	}

	keywordList = strings.Split(nameData, "\r\n")
}

type BaiduClient struct {
	Nickname string
	BDUSS    string // 百度BDUSS

	k  string
	v  string
	sv string

	client *requester.HTTPClient
}

// 将抖音那边的头像,昵称,签名等数据同步到百度
func (b *BaiduClient) SyncFromDouyin(shareURL string) error {
	if shareURL == "" {
		return errors.Errorf("用户[%s]绑定的抖音为空,不同步数据 \n", b.Nickname)
	}

	apiDouyinUser, err := api.NewAPIDouyinUser(shareURL)
	if err != nil {
		return fmt.Errorf("获取抖音用户数据失败:%s", err)
	}

	wlog.Infof("开始从抖音用户[%s]复制用户信息 \n", apiDouyinUser.Nickname)

	headImgURL := apiDouyinUser.AvatarMedium.URLList[0]
	// 更换头像
	if err := b.Setportrait(headImgURL); err != nil {
		wlog.Infof("[%s]更换头像失败: %s \n", apiDouyinUser.Nickname, err)
	} else {
		wlog.Infof("[%s]头像更换成功 \n", apiDouyinUser.Nickname)
	}

	// 昵称和签名关键字过滤后进行更换
	newName := fileterKeyWord(addExtraName(filterSpecial(apiDouyinUser.Nickname)))
	err = b.SetProfile(map[string]string{
		"nickname": newName,
	})
	if err != nil {
		wlog.Infof("[%s]更换昵称失败: %s \n", apiDouyinUser.Nickname, err)
	} else {
		wlog.Infof("[%s]昵称更换为[%s]成功 \n", apiDouyinUser.Nickname, newName)
	}

	err = b.SetProfile(map[string]string{
		"autograph": fileterKeyWord(filterSpecial(apiDouyinUser.Signature)),
	})
	if err != nil {
		wlog.Infof("[%s]更换签名失败:%s \n", apiDouyinUser.Nickname, err)
	} else {
		wlog.Infof("[%s]签名更换成功 \n", apiDouyinUser.Nickname)
	}

	wlog.Infof("[%s]信息复制完毕 \n", apiDouyinUser.Nickname)

	return nil
}

func fileterKeyWord(content string) string {
	newstr := ""
	for _, value := range keywordList {
		newstr = strings.ReplaceAll(content, value, "")
	}

	return newstr
}

// 过滤掉特殊字符
func filterSpecial(content string) string {
	var buffer bytes.Buffer
	for _, v := range content {
		if unicode.Is(unicode.Han, v) || unicode.IsLetter(v) || unicode.IsDigit(v) || unicode.IsPunct(v) {
			buffer.WriteRune(v)
			continue
		}
	}

	return buffer.String()
}

func addExtraName(name string) string {
	newName := ""

	randIndex := rand.Intn(len(nameList) - 1)

	// begin := rand.Intn(100) % 2
	// if begin == 0 {
	// newName = nameList[randIndex] + name
	// } else {
	newName = name + nameList[randIndex]
	// }

	return newName
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func calculateSig(m map[string]string, signKey string) string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sb := ""
	for k := range keys {
		sb += keys[k]
		sb += "="
		sb += m[keys[k]]
		sb += "&"
	}
	sb += "sign_key="
	sb += signKey

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(sb))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
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
