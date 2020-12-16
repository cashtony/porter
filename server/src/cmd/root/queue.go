package main

import (
	"encoding/json"
	"porter/define"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
)

type queueMsgHandler struct{}

func (q *queueMsgHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}

	// do whatever actual message processing is desired
	// err := processMessage(m.Body)
	// 节点完成了某个视频的上传
	finishedVideoID := &define.TaskUploadFinished{}
	err := json.Unmarshal(m.Body, finishedVideoID)
	if err != nil {
		wlog.Errorf("队列事件解析失败:%s \n", err)
		return nil
	}

	DB.Model(&DouyinVideo{}).Where("aweme_id = ?", finishedVideoID.AwemeID).Update("state", FinishedUpload)
	// 加入统计
	DB.Model(&Statistic{}).Where("aweme_id = ?", finishedVideoID.AwemeID).Update("state", FinishedUpload)

	return nil
}
