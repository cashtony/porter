package main

import (
	"porter/define"
	"porter/wlog"
	"sync"
)

var isUpdading = false

const (
	WaitUpload     = 0
	FinishedUpload = 1
)

// 每天固定时间更新账号
func UpdateAndUpload() {
	if isUpdading {
		wlog.Warn("当前正在运行中,不能重复操作")
		return
	}

	isUpdading = true

	// 更新抖音用户视频数据
	UpdateDouyinUsers()
	// 更新百度用户数据(主要是钻石)
	UpdateBaiduUsers()
	// 更新百度视频数据
	BaiduUsersUpload()

	isUpdading = false
	wlog.Info("更新完毕")

}

func UpdateDouyinUsers() {
	defer func() {
		if r := recover(); r != nil {
			wlog.Error("任务panic", r)
		}
	}()
	traffic := make(chan int, define.ParallelNum)
	wg := sync.WaitGroup{}

	users := make([]*DouyinUser, 0)
	result := DB.Where("last_collect_time < current_date").Find(&users)
	if result.Error != nil {
		wlog.Error("定时任务获取数据库数据时发生错误:", result.Error)
		return
	}
	wlog.Info("本次更新的用户数量为:", len(users))
	wg.Add(len(users))

	for _, user := range users {
		traffic <- 1
		go func(u *DouyinUser) {
			// 获取最新一页视频
			wlog.Debugf("开始更新用户[%s][%s]数据: \n", u.UID, u.Nickname)
			u.Update()

			wg.Done()
			<-traffic
		}(user)

	}

	wg.Wait()
}

func UpdateBaiduUsers() {

}

func BaiduUsersUpload() {
	defer func() {
		if r := recover(); r != nil {
			wlog.Error("任务panic", r)
		}
	}()
	wg := sync.WaitGroup{}
	traffic := make(chan int, define.ParallelNum)
	// 开始上传视频
	bdUsers := make([]*BaiduUser, 0)
	result := DB.Model(&BaiduUser{}).Where("last_upload_time < current_date and douyin_uid != ''").Find(bdUsers)
	if result.Error != nil {
		wlog.Error("定时任务获取百度用户时发生错误:", result.Error)
		return
	}
	wlog.Info("开始上传: 本次要上传的用户数量为:", len(bdUsers))
	wg.Add(len(bdUsers))

	for _, bduser := range bdUsers {
		traffic <- 1

		go func(u *BaiduUser) {
			u.Upload()

			wg.Done()
			<-traffic
		}(bduser)
	}

	wg.Wait()
}
