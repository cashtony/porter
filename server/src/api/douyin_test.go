package api

import (
	"fmt"
	"testing"
)

func TestNewApiDouyinUser(t *testing.T) {
	shareURL := "https://v.douyin.com/qKDMXG/"

	d, err := NewAPIDouyinUser(shareURL)
	if err != nil {
		t.Fatal(err)
	}
	if d.Nickname == "" {
		t.Error("昵称为空")
	}
	if d.AvatarMedium.URLList[0] == "" {
		t.Error("头像为空")
	}

	fmt.Printf("抖音用户数据为: %+v \n", d)
}
