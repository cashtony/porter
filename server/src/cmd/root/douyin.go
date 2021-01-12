package main

import (
	"porter/define"
	"porter/wlog"

	"gorm.io/gorm/clause"
)

type TableDouyinUser struct {
	UID            string `json:"uid" gorm:"primaryKey"`
	UniqueUID      string `json:"uniqueUID" gorm:"primaryKey"` // 抖音号
	Nickname       string `json:"nickName"`
	Gender         int    `json:"gender" gorm:"default:0"`
	Signature      string `json:"signature" gorm:"default:'"`
	Birthday       string `json:"birthday" gorm:"default:''"`
	TotalFavorited int    `json:"total_favorited" gorm:"default:0"`
	AwemeCount     int    `json:"aweme_count" gorm:"default:0"`
	FollowerCount  int    `json:"follower_count" gorm:"default:0"`
	Location       string `json:"location" gorm:"default:''"`
	Province       string `json:"province" gorm:"default:''"`
	City           string `json:"city" gorm:"default:''"`
	Avatar         string `json:"-" gorm:"default:''"`
	SecUID         string `json:"-" gorm:"default:''"` // 用于填充获取用户数据接口

	BaiduUID        string          `json:"baiduUID" gorm:"default:'"`
	SimilarSign     int             `json:"similarSign" gorm:"default:0"` // 有效位 1:全民小视频
	LastCollectTime define.JsonTime `json:"lastCollectTime"`              // 最后一次采集时间
	CreateTime      define.JsonTime `json:"createTime" gorm:"default:now()"`
}

func (TableDouyinUser) TableName() string {
	return "douyin_users"
}

func (d *TableDouyinUser) Store() {
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}, {Name: "unique_uid"}},
		DoNothing: true,
	}).Create(d)
	if DB.Error != nil {
		wlog.Errorf("抖音用户[%s][%s]存入数据库失败:%s \n", d.UID, d.Nickname, DB.Error)
		return
	}
}
func (d *TableDouyinUser) StoreVideo(list []*define.TableDouyinVideo) {
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "aweme_id"}},
		DoNothing: true,
	}).Create(list)
	if DB.Error != nil {
		wlog.Errorf("抖音用户[%s][%s]存入数据库失败:%s \n", d.UID, d.Nickname, DB.Error)
		return
	}
}
