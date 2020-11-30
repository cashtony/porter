package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"porter/requester"
	"porter/wlog"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
)

type DYVideo struct {
	AuthorUID     string
	AuthorName    string
	Awemeid       string
	Desc          string // 视频描述
	DownloadURL   string
	LocalFilePath string
}

func DouyinVideoList(userid string) []string {
	result := make([]string, 0)

	return result
}

func Mock() []*DYVideo {
	videoList := make([]*DYVideo, 0)
	j, err := simplejson.NewFromReader(strings.NewReader(mock))
	if err != nil {
		wlog.Error("数据解析失败:", err)
		return videoList
	}

	for _, item := range j.Get("aweme_list").MustArray() {
		itemJSON := item.(map[string]interface{})
		downloadURL := itemJSON["video"].(map[string]interface{})["play_addr"].(map[string]interface{})["url_list"].([]interface{})[0].(string)
		authorInfo := itemJSON["author"].(map[string]interface{})
		video := &DYVideo{
			AuthorUID:   authorInfo["uid"].(string),
			AuthorName:  authorInfo["nickname"].(string),
			Awemeid:     itemJSON["aweme_id"].(string),
			Desc:        itemJSON["desc"].(string),
			DownloadURL: downloadURL,
		}

		videoList = append(videoList, video)
	}

	return videoList
}

func UserVideoList(uid string) []*DYVideo {
	url := "https://www.iesdouyin.com/web/api/v2/aweme/post/?user_id=%s&sec_uid=&count=20&max_cursor=0&aid=1128&_signature=&dytk="
	client := requester.NewHTTPClient()
	client.SetUserAgent(requester.MobileUserAgent)

	videoList := make([]*DYVideo, 0)
	for {
		time.Sleep(5 * time.Millisecond)
		resp, err := client.Req("GET", fmt.Sprintf(url, uid), nil, nil)
		if err != nil {
			wlog.Error("访问抖音服务器报错", err)
			continue
		}
		if resp.Header.Get("status_code") == "" {
			continue
		}

		j, err := simplejson.NewFromReader(resp.Body)
		if err != nil {
			wlog.Error("数据解析失败:", err)
			return videoList
		}

		for _, item := range j.Get("aweme_list").MustArray() {
			itemJSON := item.(map[string]interface{})
			downloadURL := itemJSON["video"].(map[string]interface{})["play_addr"].(map[string]interface{})["url_list"].([]interface{})[0].(string)
			authorInfo := itemJSON["author"].(map[string]interface{})
			video := &DYVideo{
				AuthorUID:   authorInfo["uid"].(string),
				AuthorName:  authorInfo["nickname"].(string),
				Awemeid:     itemJSON["aweme_id"].(string),
				Desc:        itemJSON["desc"].(string),
				DownloadURL: downloadURL,
			}

			videoList = append(videoList, video)
		}

		break
	}

	return videoList
}

func download(video *DYVideo) {
	client := requester.NewHTTPClient()
	client.SetUserAgent(requester.MobileUserAgent)

	resp, err := client.Req("GET", video.DownloadURL, nil, nil)
	if err != nil {
		wlog.Error("下载的请求发生错误", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		wlog.Warn("访问被拒绝了")
		return
	}

	dataResp := resp
	if resp.StatusCode == http.StatusFound {
		location := resp.Header.Get("location")
		jumpURLResp, err := client.Req("GET", location, nil, nil)
		if err != nil {
			wlog.Error("下载的请求发生错误", err)
			return
		}
		defer jumpURLResp.Body.Close()

		dataResp = jumpURLResp
	}

	contentType := dataResp.Header.Get("Content-Type")
	if contentType != "video/mp4" {
		wlog.Warn("获取的类型不正确")
		return
	}
	//创建文件夹
	cmd, _ := os.Getwd()
	dirPath := fmt.Sprintf("temp/%s", video.AuthorName)
	os.MkdirAll(dirPath, os.ModePerm)
	video.LocalFilePath = fmt.Sprintf("%s/%s/%s.mp4", cmd, dirPath, video.Desc)
	f, err := os.Create(video.LocalFilePath)
	if err != nil {
		wlog.Error("下载时申请空间失败:", err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, dataResp.Body)
	if err != nil {
		wlog.Error("下载时出现错误:", err)
		return
	}
}
