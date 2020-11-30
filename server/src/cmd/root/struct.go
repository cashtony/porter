package main

import "time"

type DouyinUser struct {
	UID             string `gorm:"primaryKey"`
	UniqueUID       string `gorm:"primaryKey"` // 抖音号
	Nickname        string
	ShareURL        string    `gorm:"primaryKey"`
	BindUID         string    // 绑定的百度uid
	LastCollectTime time.Time // 最后一次采集时间

	sec_uid   string         // 用于填充获取用户数据接口
	videoList []*DouyinVideo // 此用户的视频信息
}

type BaiduUser struct {
	UID      string `gorm:"primaryKey"`
	Nickname string
	Bduss    string `gorm:"primaryKey"`
	BindUID  string // 绑定的抖音uid
}

type DouyinVideo struct {
	AwemeID    string `gorm:"primaryKey"`
	AuthorUID  string // 抖音uid
	Desc       string // 视频描述
	CreateTime time.Time
	State      int // 0未搬运 1:已搬运
}
