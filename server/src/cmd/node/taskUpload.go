package main

import (
	"encoding/json"
	"porter/define"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
)

type TaskUploadHandler struct{}

func (q *TaskUploadHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	ThreadTraffic <- 1
	// 接收root发布的任务
	task := &define.TaskUpload{}
	err := json.Unmarshal(m.Body, task)
	if err != nil {
		wlog.Error("任务解析失败:", err)
		return nil
	}
	go excuteTask(task)

	return nil
}

func excuteTask(task *define.TaskUpload) {
	wlog.Infof("用户[%s]接收到任务 数量:%d \n", task.Nickname, len(task.Videos))
	if task.Bduss == "" {
		wlog.Errorf("[%s] bduss为空 任务错误: \n", task.Nickname)
		return
	}

	client := NewBaiduClient(task.Bduss)
	if err := client.FetchSecretInfo(); err != nil {
		wlog.Errorf("用户解密数据获取失败: %s", err)
		return
	}
	finishedList := make([]string, 0)
	sucNum := 0

	for i, v := range task.Videos {
		vid := i + 1
		filterDesc := filterKeyword(v.Desc)
		wlog.Infof("[%s][%d][%s]开始下载 \n", task.Nickname, i+1, filterDesc)
		// 下载视频
		filePath, err := download(task.Nickname, filterDesc, v.DownloadURL)
		if err != nil {
			wlog.Errorf("[%s][%d][%s][%s]下载发生错误:%s \n", task.Nickname, vid, filterDesc, v.DownloadURL, err)
			continue
		}
		wlog.Infof("[%s][%d][%s]下载结束,开始上传 \n", task.Nickname, i+1, filterDesc)

		err = client.Upload(filePath, filterDesc)
		if err != nil {
			wlog.Errorf("[%s][%d][%s][%s]上传发生错误:%s \n", task.Nickname, vid, filterDesc, v.DownloadURL, err)
			continue
		}
		wlog.Infof("[%s][%d][%s]上传完毕 \n", task.Nickname, vid, filterDesc)

		finishedList = append(finishedList, v.AwemeID)

		// 发送任务完成消息
		data, err := json.Marshal(&define.TaskUploadFinished{
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

		sucNum++
	}

	deleteUserDir(task.Nickname)

	wlog.Infof("用户[%s]任务完成, 成功[%d]条, 失败[%d]条 \n", task.Nickname, sucNum, len(task.Videos)-sucNum)

	<-ThreadTraffic
}
