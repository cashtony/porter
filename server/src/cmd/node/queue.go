package main

import (
	"encoding/json"
	"porter/cmd/node/quanmin"
	"porter/define"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
)

var traffic = make(chan int, define.ParallelNum)

type queueMsgHandler struct{}

func (q *queueMsgHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
		// In this case, a message with an empty body is simply ignored/discarded.
		return nil
	}
	traffic <- 1
	// do whatever actual message processing is desired
	// err := processMessage(m.Body)
	// 接收root发布的任务
	task := &define.Task{}
	err := json.Unmarshal(m.Body, task)
	if err != nil {
		wlog.Error("任务解析失败:", err)
		return nil
	}
	go excuteTask(task)

	// Returning a non-nil error will automatically send a REQ command to NSQ to re-queue the message.
	return nil
}

func excuteTask(task *define.Task) {
	wlog.Infof("接收到任务:%+v \n", task)

	quanmin := quanmin.NewUser(task.Bduss)
	finishedList := make([]string, 0)
	for i, v := range task.Videos {
		wlog.Infof("[%d][%s]开始下载 \n", i+1, v.Desc)
		// 下载视频
		filePath, err := download(task.Nickname, v.Desc, v.DownloadURL)
		if err != nil {
			wlog.Errorf("[%s][%s]下载发生错误:%s \n", v.Desc, v.DownloadURL, err)
			continue
		}
		wlog.Infof("[%d][%s]下载结束,开始上传 \n", i+1, v.Desc)
		err = quanmin.Upload(filePath, v.Desc)
		if err != nil {
			wlog.Errorf("[%s][%s]上传发生错误:%s \n", v.Desc, v.DownloadURL, err)
			continue
		}

		wlog.Infof("[%s]上传完毕 \n", v.Desc)
		finishedList = append(finishedList, v.AwemeID)

		// 发送任务完成消息
		data, err := json.Marshal(&define.TaskFinished{
			AwemeID: v.AwemeID,
		})
		if err != nil {
			wlog.Errorf("json解析错误:%s", err)
			continue
		}

		err = Q.Publish(define.TaskFinishedTopic, data)
		if err != nil {
			wlog.Errorf("上传完成事件发布失败:%s", err)
			continue
		}
	}
	deleteUserDir(task.Nickname)

	wlog.Infof("用户[%s]任务完成 \n", task.Nickname)

	<-traffic
}
