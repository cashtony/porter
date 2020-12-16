package main

import (
	"testing"
)

func TestSetParait(t *testing.T) {
	c := NewBaiduClient("NDOGFGQVJ6QVMxcVozNkNHZnA0UFFJcVREb1p2WW1tSWN-MEI0SWU2TjloRVJmSVFBQUFBJCQAAAAAAAAAAAEAAAA48Ihab3RzdW1heGxsZm5pAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAH33HF999xxfaG")
	err := c.SyncFromDouyin("https://v.douyin.com/JCwHG3w/")
	if err != nil {
		t.Error(err)
	}
}
