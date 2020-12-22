package main

import (
	"encoding/json"
	"porter/task"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
)

type TaskChangeInfoHandler struct{}

func (*TaskChangeInfoHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	changeInfo := &task.TaskChangeInfo{}
	err := json.Unmarshal(m.Body, changeInfo)
	if err != nil {
		wlog.Error("任务解析失败:", err)
		return nil
	}

	wlog.Infof("接收到复制信息任务, 数量:%d", len(changeInfo.List))

	for _, item := range changeInfo.List {
		ThreadTraffic <- 1
		go excuteChangeInfo(item)
	}

	return nil
}

func excuteChangeInfo(item task.TaskChangeInfoItem) {
	defer func() {
		<-ThreadTraffic
	}()
	client := NewBaiduClient(item.Bduss)
	err := client.SyncFromDouyin(item.DouyinURL)
	if err != nil {
		wlog.Error("从抖音复制用户数据到全民失败:", err)
	}
}
