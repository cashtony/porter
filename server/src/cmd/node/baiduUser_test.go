package main

import (
	"testing"
)

// func TestSetParait(t *testing.T) {
// 	c := NewBaiduClient("NDOGFGQVJ6QVMxcVozNkNHZnA0UFFJcVREb1p2WW1tSWN-MEI0SWU2TjloRVJmSVFBQUFBJCQAAAAAAAAAAAEAAAA48Ihab3RzdW1heGxsZm5pAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAH33HF999xxfaG")
// 	err := c.SyncFromDouyin("https://v.douyin.com/JCwHG3w/")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestFilterCharacter(t *testing.T) {
	str := "🔥拍车纯属娱乐↵🔥一直被模仿，从未被超越😘😘😘"
	want := "拍车纯属娱乐一直被模仿，从未被超越"
	result := filterSpecial(str)
	if want != result {
		t.Error("过滤不一致", result)
	}
}

func TestFilterKeywords(t *testing.T) {
	str := "第一个龙珠#不火系列 @DOU+小助手 @抖音小助手"
	want := "第一个龙珠#不火系列 @小助手 "

	result := filterKeyword(str)
	if want != result {
		t.Error("过滤不一致", result)
	}
}
