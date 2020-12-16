package main

import (
	"encoding/json"
	"porter/define"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
)

var changeInfoTraffic = make(chan int, define.ParallelNum)

type TaskChangeInfoHandler struct{}

func (t *TaskChangeInfoHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	task := &define.TaskChangeInfo{}
	err := json.Unmarshal(m.Body, task)
	if err != nil {
		wlog.Error("任务解析失败:", err)
		return nil
	}

	wlog.Infof("接收到复制信息任务, 数量:%d", len(task.List))

	for _, item := range task.List {
		client := NewBaiduClient(item.Bduss)
		err := client.SyncFromDouyin(item.DouyinURL)
		if err != nil {
			wlog.Error("从抖音复制到全民全台失败:", err)
			continue
		}
	}

	return nil
}
