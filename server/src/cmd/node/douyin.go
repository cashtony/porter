package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"porter/requester"
	"regexp"
)

var re, _ = regexp.Compile(`[\*\\:?"/<>|]`)

func download(dirName, fileName, downloadURL string) (string, error) {
	client := requester.NewHTTPClient()
	client.SetUserAgent(requester.MobileUserAgent)

	resp, err := client.Req("GET", downloadURL, nil, nil)
	if err != nil {
		return "", fmt.Errorf("下载的请求发生错误:%s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载时访问发生错误: %d", resp.StatusCode)
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
		return "", fmt.Errorf("复制文件时发生错误:%s", err)
	}

	return filePath, nil
}

func deleteUserDir(nickname string) {
	dirPath := fmt.Sprintf("temp/%s", nickname)
	os.RemoveAll(dirPath)
	dirPath = fmt.Sprintf("thumbsnails/%s", nickname)
	os.RemoveAll(dirPath)
}
