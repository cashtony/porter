package main

import (
	"encoding/json"
	"porter/define"
	"porter/task"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
)

type taskUploadHandler struct{}

func (q *taskUploadHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	// 节点完成了某个视频的上传
	finishedVideoID := &task.TaskUploadFinished{}
	err := json.Unmarshal(m.Body, finishedVideoID)
	if err != nil {
		wlog.Errorf("上传事件解析失败:%s \n", err)
		return nil
	}

	DB.Model(&define.TableDouyinVideo{}).Where("aweme_id = ?", finishedVideoID.AwemeID).Update("state", FinishedUpload)
	// 加入统计
	DB.Model(&Statistic{}).Where("aweme_id = ?", finishedVideoID.AwemeID).Update("state", FinishedUpload)

	return nil
}
