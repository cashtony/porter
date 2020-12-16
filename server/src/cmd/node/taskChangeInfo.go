package main

import (
	"porter/define"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
)

var changeInfoTraffic = make(chan int, define.ParallelNum)

type TaskChangeInfoHandler struct{}

func (t *TaskChangeInfoHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}
	// do whatever actual message processing is desired
	// err := processMessage(m.Body)
	// 接收root发布的任务
	task := &define.TaskChangeInfo{}
	for _, item := range task.List {
		client := NewBaiduClient(item.Bduss)
		err := client.SyncFromDouyin(item.DouyinURL)
		if err != nil {
			wlog.Error("从抖音复制到全民全台失败:", err)
			continue
		}
	}
	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}
