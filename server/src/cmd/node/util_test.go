package main

import (
	"fmt"
	"testing"
)

func TestGetSecSig(t *testing.T) {
	url := "https://v.douyin.com/JqQQY4p/"
	sig := GetSecSig(url)
	if sig == "" {
		t.Error("获取加密signature失败")
	}

	fmt.Println("signature:", sig)
}
