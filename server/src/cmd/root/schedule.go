package main

import (
	"encoding/json"
	"porter/define"
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

	wlog.Info("开始更新百度用户数据")
	DailyUpdateBaiduUsers()
	wlog.Info("百度用户数据更新完毕")

	// 更新百度视频数据
	wlog.Info("开始发布百度用户视频任务")
	ScheduleUpload(UploadTypeDaily)
	// BaiduUsersUpload(UploadTypeDaily)
	wlog.Info("百度用户视频任务发布完毕")
	isUpdading = false
	wlog.Info("更新完毕")
}

func ScheduleUpdate() {
	// todo 查看topic长度, 如果不为0说明之前的更新任务没处理完就不进行更新

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

}
