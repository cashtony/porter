package main

import (
	"porter/define"
	"porter/wlog"

	"gorm.io/gorm/clause"
)

type TableDouyinUser struct {
	UID             string          `json:"uid" gorm:"primaryKey"`
	UniqueUID       string          `json:"uniqueUID" gorm:"primaryKey"` // 抖音号
	Nickname        string          `json:"nickName"`
	ShareURL        string          `json:"shareURL" gorm:"primaryKey"`
	VideoCount      int             `json:"videoCount"`
	FansCount       int             `json:"fansNum"`
	LastCollectTime define.JsonTime `json:"lastCollectTime"` // 最后一次采集时间
	CreateTime      define.JsonTime `json:"createTime" gorm:"default:now()"`

	secUID string // 用于填充获取用户数据接口
}

func (TableDouyinUser) TableName() string {
	return "douyin_users"
}

func (d *TableDouyinUser) Store() {
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}},
		UpdateAll: true,
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
