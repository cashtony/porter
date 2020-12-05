package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"porter/define"
	"porter/wlog"
	"time"

	"gorm.io/gorm/clause"
)

var traffic = make(chan int, define.ParallelNum)

const (
	WaitUpload     = 0
	FinishedUpload = 1
)

// 每天固定时间更新账号
func UpdateAndUpload() {
	users := make([]*DouyinUser, 0)
	DB.Where("last_collect_time < current_date").Find(&users)
	if DB.Error != nil {
		wlog.Error("定时任务获取数据库数据时发生错误:", DB.Error)
		return
	}
	wlog.Info("本次更新的用户数量为:", len(users))
	for _, user := range users {
		traffic <- 1
		go updateOneUser(user)
	}

}

func updateOneUser(user *DouyinUser) {
	defer func() {
		<-traffic

		if r := recover(); r != nil {
			wlog.Error("任务panic", user.UID, r)
		}
	}()
	// 获取最新一页视频
	wlog.Debugf("开始更新用户[%s][%s]数据: \n", user.UID, user.Nickname)
	onePageList, _, _, err := user.OnePageVideo(0)
	if err != nil {
		wlog.Errorf("用户[%s][%s]获取视频列表失败:%s \n", user.UID, user.Nickname, err)
		return
	}

	// 将没有的视频传入到数据库中
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "aweme_id"}},
		DoNothing: true,
	}).Create(onePageList)
	if DB.Error != nil {
		wlog.Errorf("用户[%s][%s]新视频信息存入数据库失败:%s \n", user.UID, user.Nickname, DB.Error)
		return
	}

	//更新用户的last_collect_time字段
	DB.Model(user).Update("last_collect_time", time.Now())
	if DB.Error != nil {
		wlog.Errorf("从数据库中更新用户[%s][%s]last_collect_time字段失败: %s \n", user.UID, user.Nickname, DB.Error)
		return
	}

	if user.BaiduUID != "" {
		publishTask(user)
	}
}

func publishTask(user *DouyinUser) {
	// 上传视频(从未上传的视频中挑选8-12条)
	randomNum := rand.Intn(MaxUploadNum-MinUploadNum) + MinUploadNum
	uploadVideoList := make([]*DouyinVideo, 0)
	videoModel := DB.Model(&DouyinVideo{}).Where("author_uid = ? and state = ?", user.UID, WaitUpload).Order("create_time desc").Limit(randomNum)

	videoModel.Debug().Where("date(create_time) = current_date - 1").Find(&uploadVideoList)
	if DB.Error != nil {
		wlog.Errorf("从数据库中获取用户[%s][%s]视频列表信息失败:%s \n", user.UID, user.Nickname, DB.Error)
		return
	}
	if len(uploadVideoList) == 0 {
		wlog.Infof("用户[%s][%s]昨天没有更新,将获取以前的视频 \n", user.UID, user.Nickname)
		videoModel.Debug().Find(&uploadVideoList)
		if DB.Error != nil {
			wlog.Errorf("从数据库中获取用户[%s][%s]视频列表信息失败:%s \n", user.UID, user.Nickname, DB.Error)
			return
		}
	}

	if len(uploadVideoList) == 0 {
		wlog.Infof("用户[%s][%s]没有可更新内容,退出 \n", user.UID, user.Nickname)
		return
	}

	// 查找视频下载url
	taskVideoList := make([]*define.TaskVideo, 0)
	statisticList := make([]*Statistic, 0)

	for _, v := range uploadVideoList {
		videoExtranInfo, err := getVideoCreateTime(v.AwemeID)
		if err != nil {
			wlog.Error("获取视频额外数据发生错误:", err)
			continue
		}
		taskVideoList = append(taskVideoList, &define.TaskVideo{
			AwemeID:     v.AwemeID,
			Desc:        v.Desc,
			DownloadURL: fmt.Sprintf("%s/?video_id=%s&ratio=720p&line=0", define.GetVideoDownload, videoExtranInfo.VID),
		})

		statisticList = append(statisticList, &Statistic{
			BaiduUID:  user.BaiduUID,
			DouyinUID: user.UID,
			AwemeID:   v.AwemeID,
			State:     WaitUpload,
		})
	}

	if len(uploadVideoList) == 0 {
		wlog.Infof("用户[%s][%s]没有可更新内容,退出 \n", user.UID, user.Nickname)
		return
	}

	bduss := ""
	DB.Model(&BaiduUser{}).Select("bduss").Where("uid = ?", user.BaiduUID).First(&bduss)
	if DB.Error != nil {
		wlog.Errorf("从数据库中获取用[%s][%s]绑定的bduss字段失败: %s \n", user.UID, user.Nickname, DB.Error)
		return
	}

	// 封装成task投递到任务队列中
	wlog.Debugf("开始投放用户[%s][%s]任务: \n", user.UID, user.Nickname)
	t := &define.Task{
		Bduss:    bduss,
		Videos:   taskVideoList,
		Nickname: user.Nickname,
	}

	data, err := json.Marshal(t)
	if err != nil {
		wlog.Error("task解析成json错误", err)
		return
	}

	// 增加数据统计
	DB.Create(&statisticList)

	err = Q.Publish(define.TaskPushTopic, data)
	if err != nil {
		wlog.Error("任务发布失败:", err)
	}
}
