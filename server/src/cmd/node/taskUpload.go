package main

import (
	"encoding/json"
	"porter/define"
	"porter/task"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
)

type TaskUploadHandler struct{}

func (q *TaskUploadHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	ThreadTraffic <- 1
	// 接收root发布的任务t
	t := &task.TaskUpload{}
	err := json.Unmarshal(m.Body, t)
	if err != nil {
		wlog.Error("任务解析失败:", err)
		return nil
	}
	go excuteTask(t)

	return nil
}

func excuteTask(t *task.TaskUpload) {
	defer func() {
		<-ThreadTraffic
	}()

	wlog.Infof("用户[%s]接收到任务 数量:%d \n", t.Nickname, len(t.Videos))
	if t.Bduss == "" {
		wlog.Errorf("[%s] bduss为空 任务错误: \n", t.Nickname)
		return
	}

	client := NewBaiduClient(t.Bduss)
	if err := client.FetchSecretInfo(); err != nil {
		wlog.Errorf("用户解密数据获取失败: %s", err)
		return
	}
	finishedList := make([]string, 0)
	sucNum := 0

	for i, v := range t.Videos {
		vid := i + 1
		filterDesc := filterKeyword(v.Desc)
		wlog.Infof("[%s][%d][%s]开始下载 \n", t.Nickname, i+1, filterDesc)
		// 下载视频
		filePath, err := download(t.Nickname, filterDesc, v.DownloadURL)
		if err != nil {
			wlog.Errorf("[%s][%d][%s][%s]下载发生错误:%s \n", t.Nickname, vid, filterDesc, v.DownloadURL, err)
			continue
		}
		wlog.Infof("[%s][%d][%s]下载结束,开始上传 \n", t.Nickname, i+1, filterDesc)

		err = client.Upload(filePath, filterDesc)
		if err != nil {
			wlog.Errorf("[%s][%d][%s][%s]上传发生错误:%s \n", t.Nickname, vid, filterDesc, v.DownloadURL, err)
			continue
		}
		wlog.Infof("[%s][%d][%s]上传完毕 \n", t.Nickname, vid, filterDesc)

		finishedList = append(finishedList, v.AwemeID)

		// 发送任务完成消息
		data, err := json.Marshal(&task.TaskUploadFinished{
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

	deleteUserDir(t.Nickname)

	wlog.Infof("用户[%s]任务完成, 成功[%d]条, 失败[%d]条 \n", t.Nickname, sucNum, len(t.Videos)-sucNum)
}
