package main

import "testing"

func TestIsSimilarInQuanmin(t *testing.T) {
	has, err := IsSimilarInQuanmin("小鹏动漫｀")
	if err != nil {
		t.Error(err)
	}

	if !has {
		t.Error("结果不正确")
	}
}
