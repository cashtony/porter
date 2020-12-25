package main

import (
	"encoding/json"
	"porter/define"
	"porter/queue"
	"porter/task"
	"porter/wlog"
	"sync"
)

var isUpdading = false

const (
	WaitUpload     = 0
	FinishedUpload = 1
)

// 只更新当天最新的
func NewlyUpload() {
	if isUpdading {
		wlog.Warn("当前正在运行中,不能重复操作")
		return
	}

	isUpdading = true
	ScheduleUpload(UploadTypeNewly)
	isUpdading = false

	wlog.Info("当天最新视频更新完毕")
}

// 每天固定时间更新账号
func DailyUpload() {
	if isUpdading {
		wlog.Warn("当前正在运行中,不能重复操作")
		return
	}

	isUpdading = true
	// 更新百度视频数据
	wlog.Info("开始发布百度用户视频任务")
	ScheduleUpload(UploadTypeDaily)
	// BaiduUsersUpload(UploadTypeDaily)
	wlog.Info("百度用户视频任务发布完毕")
	isUpdading = false
	wlog.Info("更新完毕")
}

func ScheduleFetchNewVideo() {
	api, err := queue.GetTopicStat("TaskParseVideo")
	if err != nil {
		wlog.Info("获取队列数据错误,将跳过此次检测")
		return
	}

	if len(api.Topics) == 0 || len(api.Topics[0].Channels) == 0 {
		wlog.Info("获取队列长度错误,将跳过此次检测")
		return
	}

	if api.Topics[0].Channels[0].Depth != 0 {
		wlog.Info("队列中还有未处理的消息,将跳过此次检测")
		return
	}

	bdUsers := make([]*TableBaiduUser, 0)
	result := DB.Model(&TableBaiduUser{}).Where("douyin_url != ''").Find(&bdUsers)
	if result.Error != nil {
		wlog.Error("获取百度用户时发生错误:", result.Error)
		return
	}

	for _, u := range bdUsers {
		// 发布任务
		t := &task.TaskParseVideo{
			Type:     define.ParseVideoTypeOnePage,
			ShareURL: u.DouyinURL,
		}
		data, err := json.Marshal(t)
		if err != nil {
			wlog.Error("任务解析失败:", err)
			continue
		}

		if err := Q.Publish(define.TaskParseVideoTopic, data); err != nil {
			wlog.Error("视频解析任务发布失败:", err)
		}
	}
}

func ScheduleUpload(utype UploadType) {
	// 取出适合数据的百度用户
	defer func() {
		if r := recover(); r != nil {
			wlog.Error("任务panic", r)
		}
	}()
	wg := sync.WaitGroup{}

	bdUsers := make([]*TableBaiduUser, 0)
	result := DB.Model(&TableBaiduUser{}).Where("douyin_url != ''").Find(&bdUsers)
	if result.Error != nil {
		wlog.Error("获取百度用户时发生错误:", result.Error)
		return
	}

	wg.Add(len(bdUsers))

	for _, bduser := range bdUsers {
		ThreadTraffic <- 1

		go func(u *TableBaiduUser) {
			if u.Status == int(BaiduUserStatusNormal) {
				u.UploadVideo(utype)
			}

			wg.Done()
			<-ThreadTraffic

		}(bduser)
	}

	wg.Wait()
}

func DailyUpdateBaiduUsers() {
	usersList := make([]*TableBaiduUser, 0)
	DB.Model(&TableBaiduUser{}).Where("status = ?", BaiduUserStatusNormal).Find(&usersList)
	wlog.Infof("本次更新基础数据的用户数量为:%d", len(usersList))
	UpdateBaiduUser(usersList)
	wlog.Infof("更新百度用户基础数据完毕")
}
