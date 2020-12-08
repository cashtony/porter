package main

import "time"

type DouyinVideo struct {
	AwemeID    string `gorm:"primaryKey"`
	AuthorUID  string // 抖音uid
	Desc       string // 视频描述
	CreateTime time.Time
	State      int // 0未搬运 1:已搬运
}

type DouyinVideoExtraInfo struct {
	CreateTime time.Time
	VID        string
}

type Account struct {
	UID        int    `json:"uid" gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Name       string `json:"name"`
	Password   string `json:"-"`
	Token      string `json:"-"`
	Rule       int    `json:"role"` // 1:管理员 50:文员
	CreateTime time.Time
}

type FaildRecords struct {
	ID         int `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Bduss      string
	Douyin     string
	CreateTime time.Time `gorm:"default:now()"`
}

type Statistic struct {
	ID         int       `gorm:"AUTO_INCREMENT"`
	BaiduUID   string    // 传到哪个百度uid
	DouyinUID  string    // 从哪个抖音号中搬运
	AwemeID    string    // 视频id
	UploadTime time.Time `gorm:"default:now()"` // 上传时间
	State      int       // 上传状态 0:上传成功 1:上传中
}

// 一级页面显示 XX用户多少条视频
type StatisticRough struct {
	UID      string `json:"uid"`
	Nickname string `json:"nickname"`
	Date     string `json:"date"`
	Num      int    `json:"num"`
}

// 二级页面,选择了某个用户查看详情时显示
type StatisticDetail struct {
	UID        string    `json:"uid"`
	VideoDesc  string    `json:"videoDesc"`
	UploadTime time.Time `json:"uploadTime"`
}
