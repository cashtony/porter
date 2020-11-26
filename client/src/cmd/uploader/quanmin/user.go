package quanmin

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"porter/requester"
	"porter/util"
	"porter/wlog"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var (
	baiduComURL = &url.URL{
		Scheme: "http",
		Host:   "baidu.com",
	}
	cookieDomain = ".baidu.com"
	uploadURL    = "https://quanmin.baidu.com/web/publish/upload"
)

type baiduAccount struct {
	uid    uint64 // 百度ID对应的uid
	name   string
	bduss  string // 百度BDUSS
	PTOKEN string
	STOKEN string

	client *requester.HTTPClient
}

// NewUserInfo 检测BDUSS有效性, 同时获取百度详细信息
func (b *baiduAccount) Login() error {
	if b.uid != 0 {
		return errors.New("当前已经登录")
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	post := map[string]string{
		"bdusstoken":  b.bduss + "|null",
		"channel_id":  "",
		"channel_uid": "",
		"stErrorNums": "0",
		"subapp_type": "mini",
		"timestamp":   timestamp + "922",
	}
	ClientSignature(post)

	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"Cookie":       "ka=open",
		"net":          "1",
		"User-Agent":   "bdtb for Android 6.9.2.1",
		"client_logid": timestamp + "416",
		"Connection":   "Keep-Alive",
	}

	resp, err := b.client.Req("POST", "http://tieba.baidu.com/c/s/login", post, header) // 获取百度ID的UID，BDUSS等
	if err != nil {
		return fmt.Errorf("检测BDUSS有效性网络错误, %s", err)
	}
	defer resp.Body.Close()

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("检测BDUSS有效性json解析出错: %s", err)
	}

	errCode := json.Get("error_code").MustString()
	errMsg := json.Get("error_msg").MustString()
	if errCode != "0" {
		return fmt.Errorf("检测BDUSS有效性错误代码: %s, 消息: %s", errCode, errMsg)
	}

	userJSON := json.Get("user")
	uidStr := userJSON.Get("id").MustString()
	b.uid, _ = strconv.ParseUint(uidStr, 10, 64)
	b.name = userJSON.Get("name").MustString()

	return nil
}

func (b *baiduAccount) Test() {
	resp, err := b.client.Req("GET", "https://quanmin.baidu.com/web/publish/upload", nil, nil) // 获取百度ID的UID，BDUSS等
	if err != nil {
		wlog.Info("请求失败")
		return
	}
	defer resp.Body.Close()

	util.PrintBody(resp.Body)
}

func (b *baiduAccount) Upload(filePath, desc string) error {
	var ids []cdp.NodeID
	cookies := map[string]string{
		"BDUSS": b.bduss,
	}
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.UserAgent(requester.UserAgent),
	}

	allocatorCtx, allocatorCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocatorCancel()

	optsctx, optCancel := chromedp.NewContext(
		allocatorCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer optCancel()

	ctx, cancel := context.WithTimeout(optsctx, 600*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// create cookie expiration
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			// add cookies to chrome
			for k, v := range cookies {
				success, err := network.SetCookie(k, v).
					WithExpires(&expr).
					WithDomain(cookieDomain).
					WithHTTPOnly(false).
					Do(ctx)
				if err != nil {
					return err
				}
				if !success {
					return fmt.Errorf("could not set cookie %s to %s", k, v)
				}
			}
			return nil
		}),
		// navigate to site
		chromedp.Navigate(uploadURL),
		chromedp.WaitReady(`div > span > div > span > input[type=file]`, chromedp.ByQuery),
		chromedp.NodeIDs(`div > span > div > span > input[type=file]`, &ids, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			dom.SetFileInputFiles([]string{filePath}).WithNodeID(ids[0]).Do(ctx)
			return nil
		}),
		chromedp.WaitReady(`div[class^="success"]`, chromedp.ByQuery),
		chromedp.Click(`div > button.ant-btn.ant-btn-primary`, chromedp.ByQuery),
		chromedp.SendKeys(`span.ant-form-item-children > textarea`, desc, chromedp.ByQuery),
		chromedp.Click(`button[class^="ant-btn btn-publish"]`, chromedp.ByQuery),
		chromedp.Sleep(5*time.Second),
	)
	if err != nil {
		return fmt.Errorf("自动上传失败:%s", err)
	}

	return nil
}

func (b *baiduAccount) streamUpload(data []byte) error {
	//https://quanmin.baidu.com/wise/video/pcpub/getuploadid?video_num=1
	resp, err := b.client.Req("POST", "https://quanmin.baidu.com/wise/video/pcpub/getuploadid?video_num=1", nil, nil) // 获取百度ID的UID，BDUSS等
	if err != nil {
		return fmt.Errorf("检测BDUSS有效性网络错误, %s", err)
	}
	defer resp.Body.Close()

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("请求上传的json解析出错: %s", err)
	}

	errCode := json.Get("errno").MustInt()
	errMsg := json.Get("errmsg").MustString()
	if errCode != 0 {
		return fmt.Errorf("请求获取上传视频信息失败: %d, 消息: %s", errCode, errMsg)
	}

	var bucket, upload_id, key, media_id string

	for _, single := range json.Get("data").Get("upload").MustArray() {
		item := single.(map[string]interface{})
		bucket = item["bucket"].(string)
		upload_id = item["upload_id"].(string)
		media_id = item["media_id"].(string)
		key = item["key"].(string)
	}
	wlog.Infof("bucket: %s upload_id:%s media_id: %s, key:%s \n", bucket, upload_id, media_id, key)
	// 开始发送
	body_buf := bytes.NewBuffer(make([]byte, 0))
	body_writer := multipart.NewWriter(body_buf)

	fileWriter, err := body_writer.CreateFormFile("file", "file_1")
	if err != nil {
		return fmt.Errorf("创建multipart失败:%s", err)
	}

	f, err := os.Open("D:/喜欢的视频.mp4")
	if err != nil {
		return fmt.Errorf("打开文件失败:%s", err)
	}

	io.Copy(fileWriter, f)
	body_writer.WriteField("part_num", "1")
	body_writer.WriteField("media_id", media_id)
	body_writer.WriteField("upload_id", upload_id)
	body_writer.Close()

	header := map[string]string{
		"Content-Type":   body_writer.FormDataContentType(),
		"content-length": strconv.Itoa(body_buf.Len()),
	}

	uploadResp, uploadErr := b.client.Req("POST", "https://quanmin.baidu.com/wise/video/pcpub/uploadvideopart", body_buf, header)
	if uploadErr != nil {
		return fmt.Errorf("上传请求出现错误 %s", err)
	}
	defer uploadResp.Body.Close()

	upJSON, upErr := simplejson.NewFromReader(uploadResp.Body)
	if upErr != nil {
		return fmt.Errorf("json解析失败: %s", err)
	}
	errCode = upJSON.Get("errno").MustInt()
	errMsg = upJSON.Get("errmsg").MustString()
	if errCode != 0 {
		return fmt.Errorf("上传视频信息失败: %d, 消息: %s", errCode, errMsg)
	}

	wlog.Info("etag:", upJSON.Get("data").Get("eTag").MustString())

	return nil
}

func (b *baiduAccount) Name() string {
	return b.name
}

func NewUser(bduss string) *baiduAccount {
	b := &baiduAccount{
		bduss:  bduss,
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
