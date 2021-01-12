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

	changeInfo := &task.TaskChangeInfoItem{}
	err := json.Unmarshal(m.Body, changeInfo)
	if err != nil {
		wlog.Error("任务解析失败:", err)
		return nil
	}

	ThreadTraffic <- 1
	go excuteChangeInfo(changeInfo)

	return nil
}

func excuteChangeInfo(item *task.TaskChangeInfoItem) {
	defer func() {
		<-ThreadTraffic
	}()

	client := NewBaiduClient(item.Bduss)
	err := client.SyncFromDouyin(item)
	if err != nil {
		wlog.Error("从抖音复制用户数据到全民失败:", err)
	}
}
