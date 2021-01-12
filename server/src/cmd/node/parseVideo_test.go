package main

import (
	"fmt"
	"porter/api"
	"testing"
)

func TestParseVideo(t *testing.T) {
	secUID := "MS4wLjABAAAAr9bIGGSMKw79HUUa8t_bsYBGpZ2O1NW_mpSeRi7NqDcWL-K1ddZCA-bO38ietqed"
	apiVideoList, _, err := api.MobileDouyinVideo(secUID, 0)
	if err != nil {
		t.Error("获取单页视频发生错误:", err)
		return
	}

	for _, item := range apiVideoList {
		fmt.Println(item.Desc)
	}

	if len(apiVideoList) == 0 {
		t.Error("获取的视频长度错误")
	}
}
