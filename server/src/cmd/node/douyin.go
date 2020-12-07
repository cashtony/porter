package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"porter/requester"
	"porter/wlog"
	"regexp"
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

var re, _ = regexp.Compile(`[\*\\:?"/<>|]`)

func download(dirName, fileName, downloadURL string) (string, error) {
	client := requester.NewHTTPClient()
	client.SetUserAgent(requester.MobileUserAgent)

	resp, err := client.Req("GET", downloadURL, nil, nil)
	if err != nil {
		return "", fmt.Errorf("下载的请求发生错误:%s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return "", errors.New("访问被拒绝了")
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "video/mp4" {
		return "", fmt.Errorf("获取的类型不正确:%s", contentType)
	}
	//创建文件夹
	cmd, _ := os.Getwd()
	dirPath := fmt.Sprintf("temp/%s", dirName)
	os.MkdirAll(dirPath, os.ModePerm)
	// 去掉特殊字符,否则windows下会报错The filename, directory name, or volume label syntax is incorrect.
	fileName = re.ReplaceAllString(fileName, "")
	filePath := fmt.Sprintf("%s/%s/%s.mp4", cmd, dirPath, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("下载时申请空间失败:%s", err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", fmt.Errorf("交换空间时发生错误:%s", err)
	}

	return filePath, nil
}

func deleteUserDir(nickname string) {
	dirPath := fmt.Sprintf("temp/%s", nickname)
	os.RemoveAll(dirPath)
}
