package api

import (
	"testing"
)

func TestNewApiDouyinUser(t *testing.T) {
	shareURL := "https://v.douyin.com/qKDMXG/"

	d, err := NewAPIWebDouyinUser(shareURL)
	if err != nil {
		t.Fatal(err)
	}
	if d.Nickname == "" {
		t.Error("昵称为空")
	}
	if d.AvatarMedium.URLList[0] == "" {
		t.Error("头像为空")
	}
}

func TestGetVideoExtraInfo(t *testing.T) {
	v, err := GetVideoExtraInfo("6834090710124236043")
	if err != nil {
		t.Error(err)
	}

	if len(v.ItemList) == 0 {
		t.Error("获取的长度错误")
	}

	if v.ItemList[0].CreateTime == 0 {
		t.Error("时间获取失败")
	}

	if v.ItemList[0].Video.Vid == "" {
		t.Error("vid获取失败")
	}
}
