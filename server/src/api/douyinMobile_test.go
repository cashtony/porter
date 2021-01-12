package api

import (
	"testing"
)

func TestGetPhoneUser(t *testing.T) {
	secUID := GetSecID("https://v.douyin.com/JCTxF7S/")
	if secUID == "" {
		t.Error("secuid 获取失败")
		return
	}

	apiUser, err := NewPhoneDouyinUser(secUID)
	if err != nil {
		t.Error(err)
		return
	}
	if apiUser.User.Nickname != "小鹏动漫｀" {
		t.Error("获取到的用户数据不一致")
	}
	if len(apiUser.User.Avatar.URLList) == 0 {
		t.Error("没有获取到用户头像")
	}
}
