package main

import (
	"porter/api"
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

	// 更新抖音用户视频数据
	// wlog.Info("开始更新抖音用户数据")
	// UpdateDouyinUsers()
	// wlog.Info("抖音用户数据更新完毕")
	// 更新百度用户数据(主要是钻石)
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

func ScheduleUpload(utype UploadType) {
	// 取出适合数据的百度用户
	defer func() {
		if r := recover(); r != nil {
			wlog.Error("任务panic", r)
		}
	}()
	wg := sync.WaitGroup{}

	// 开始上传视频
	bdUsers := make([]*BaiduUser, 0)
	result := DB.Model(&BaiduUser{}).Where("douyin_url != ''").Find(&bdUsers)
	if result.Error != nil {
		wlog.Error("获取百度用户时发生错误:", result.Error)
		return
	}

	wg.Add(len(bdUsers))

	for _, bduser := range bdUsers {
		ThreadTraffic <- 1

		go func(u *BaiduUser) {
			// 更新绑定的抖音用户最新视频
			apiDouyinUser, err := api.NewAPIDouyinUser(u.DouyinURL)
			if err != nil {
				wlog.Error("获取抖音用户数据失败:", err)
				<-ThreadTraffic
				return
			}
			duser := &DouyinUser{
				UID:      apiDouyinUser.UID,
				ShareURL: u.DouyinURL,
				Nickname: apiDouyinUser.Nickname,
				secUID:   apiDouyinUser.SecUID,
			}
			duser.Update()
			if u.Status == int(BaiduUserStatusNormal) {
				u.UploadVideo(utype)
			}

			wg.Done()
			<-ThreadTraffic

		}(bduser)
	}

	wg.Wait()

}

// func UpdateDouyinUsers() {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			wlog.Error("任务panic", r)
// 		}
// 	}()
// 	wg := sync.WaitGroup{}

// 	users := make([]*DouyinUser, 0)
// 	result := DB.Find(&users)
// 	if result.Error != nil {
// 		wlog.Error("定时任务获取数据库数据时发生错误:", result.Error)
// 		return
// 	}

// 	wlog.Info("本次更新的用户数量为:", len(users))
// 	wg.Add(len(users))

// 	for _, user := range users {
// 		ThreadTraffic <- 1
// 		go func(u *DouyinUser) {
// 			// 获取最新一页视频
// 			u.secUID = api.GetSecID(u.ShareURL)
// 			wlog.Debugf("开始更新用户[%s][%s]数据:", u.UID, u.Nickname)
// 			u.Update()

// 			wg.Done()
// 			<-ThreadTraffic
// 		}(user)

// 	}

// 	wg.Wait()
// }

func DailyUpdateBaiduUsers() {

}

// func BaiduUsersUpload(uType UploadType) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			wlog.Error("任务panic", r)
// 		}
// 	}()
// 	wg := sync.WaitGroup{}

// 	// 开始上传视频
// 	bdUsers := make([]*BaiduUser, 0)
// 	subDB := DB.Model(&BaiduUser{}).Where("douyin_url != ''")
// 	if uType == UploadTypeDaily {
// 		subDB = subDB.Where("last_upload_time < current_date")
// 	}
// 	subDB.Find(&bdUsers)
// 	if subDB.Error != nil {
// 		wlog.Error("定时任务获取百度用户时发生错误:", subDB.Error)
// 		return
// 	}
// 	wlog.Info("开始上传: 本次要上传的用户数量为:", len(bdUsers))
// 	wg.Add(len(bdUsers))

// 	for _, bduser := range bdUsers {
// 		ThreadTraffic <- 1

// 		go func(u *BaiduUser) {
// 			u.UploadVideo(uType)

// 			wg.Done()
// 			<-ThreadTraffic
// 		}(bduser)
// 	}

// 	wg.Wait()
// }
