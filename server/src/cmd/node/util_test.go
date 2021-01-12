package main

import (
	"fmt"
	"porter/api"
	"testing"
)

func TestGetSecSig(t *testing.T) {
	url := "https://v.douyin.com/JqQQY4p/"
	sig := api.GetSecSignature(url)
	if sig == "" {
		t.Error("获取加密signature失败")
	}

	fmt.Println("signature:", sig)
}
