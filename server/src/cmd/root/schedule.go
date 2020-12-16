package main

import (
	"porter/api"
	"porter/define"
	"porter/wlog"
	"sync"
)

var isUpdading = false

const (
	WaitUpload     = 0
	FinishedUpload = 1
)

// 只更新当天最新的
func NewlyUpdate() {
	if isUpdading {
		wlog.Warn("当前正在运行中,不能重复操作")
		return
	}

	isUpdading = true

	wlog.Info("开始更新抖音用户数据")
	UpdateDouyinUsers(UpdateTypeNewly)
	wlog.Info("开始更新最新视频")
	BaiduUsersUpload(UpdateTypeNewly)

	isUpdading = false

	wlog.Info("当天最新视频更新完毕")
}

// 每天固定时间更新账号
func DailyUpdate() {
	if isUpdading {
		wlog.Warn("当前正在运行中,不能重复操作")
		return
	}

	isUpdading = true

	// 更新抖音用户视频数据
	wlog.Info("开始更新抖音用户数据")
	UpdateDouyinUsers(UpdateTypeDaily)
	wlog.Info("抖音用户数据更新完毕")
	// 更新百度用户数据(主要是钻石)
	wlog.Info("开始更新百度用户数据")
	DailyUpdateBaiduUsers()
	wlog.Info("百度用户数据更新完毕")

	// 更新百度视频数据
	wlog.Info("开始发布百度用户视频任务")
	BaiduUsersUpload(UpdateTypeDaily)
	wlog.Info("百度用户视频任务发布完毕")
	isUpdading = false
	wlog.Info("更新完毕")
}

func UpdateDouyinUsers(utype UpdateType) {
	defer func() {
		if r := recover(); r != nil {
			wlog.Error("任务panic", r)
		}
	}()
	traffic := make(chan int, define.ParallelNum)
	wg := sync.WaitGroup{}

	users := make([]*DouyinUser, 0)
	if utype == UpdateTypeDaily {
		result := DB.Where("last_collect_time < current_date").Find(&users)
		if result.Error != nil {
			wlog.Error("定时任务获取数据库数据时发生错误:", result.Error)
			return
		}
	} else {
		result := DB.Find(&users)
		if result.Error != nil {
			wlog.Error("定时任务获取数据库数据时发生错误:", result.Error)
			return
		}
	}

	wlog.Info("本次更新的用户数量为:", len(users))
	wg.Add(len(users))

	for _, user := range users {
		traffic <- 1
		go func(u *DouyinUser) {
			// 获取最新一页视频
			u.secUID = api.GetSecID(u.ShareURL)
			wlog.Debugf("开始更新用户[%s][%s]数据: \n", u.UID, u.Nickname)
			u.Update()

			wg.Done()
			<-traffic
		}(user)

	}

	wg.Wait()
}

func DailyUpdateBaiduUsers() {

}

func BaiduUsersUpload(uType UpdateType) {
	defer func() {
		if r := recover(); r != nil {
			wlog.Error("任务panic", r)
		}
	}()
	wg := sync.WaitGroup{}
	traffic := make(chan int, define.ParallelNum)
	// 开始上传视频
	bdUsers := make([]*BaiduUser, 0)
	subDB := DB.Model(&BaiduUser{}).Where("douyin_uid != ''")
	if uType == UpdateTypeDaily {
		subDB = subDB.Where("last_upload_time < current_date")
	}
	subDB.Find(&bdUsers)
	if subDB.Error != nil {
		wlog.Error("定时任务获取百度用户时发生错误:", subDB.Error)
		return
	}
	wlog.Info("开始上传: 本次要上传的用户数量为:", len(bdUsers))
	wg.Add(len(bdUsers))

	for _, bduser := range bdUsers {
		traffic <- 1

		go func(u *BaiduUser) {
			u.UploadVideo(uType)

			wg.Done()
			<-traffic
		}(bduser)
	}

	wg.Wait()
}

// func NewlyBaiduUsersUpload() {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			wlog.Error("任务panic", r)
// 		}
// 	}()
// 	wg := sync.WaitGroup{}
// 	traffic := make(chan int, define.ParallelNum)
// 	// 开始上传视频
// 	bdUsers := make([]*BaiduUser, 0)
// 	result := DB.Model(&BaiduUser{}).Where("douyin_uid != ''").Find(&bdUsers)
// 	if result.Error != nil {
// 		wlog.Error("定时任务获取百度用户时发生错误:", result.Error)
// 		return
// 	}
// 	wlog.Info("开始上传: 本次要上传的用户数量为:", len(bdUsers))
// 	wg.Add(len(bdUsers))

// 	for _, bduser := range bdUsers {
// 		traffic <- 1

// 		go func(u *BaiduUser) {
// 			u.UploadVideo(UpdateTypeNewly)

// 			wg.Done()
// 			<-traffic
// 		}(bduser)
// 	}

// 	wg.Wait()

// 	wlog.Info("本次只更新最新视频完成")
// }
