package main

import (
	"flag"
	"porter/db"
	"porter/wlog"

	"gorm.io/gorm"
)

var Mode = flag.String("mode", "debug", "运行模式 debug:开发模式, release:产品模式")
var DB *gorm.DB

// center负责采集,分类,分发任务
func main() {
	flag.Parse()

	if *Mode == "debug" {
		wlog.DevelopMode()
	}
	DB = db.NewPG()
	DB.AutoMigrate(&DouyinVideo{})
	DB.AutoMigrate(&DouyinUser{})
	DB.AutoMigrate(&BaiduUser{})

	user, err := NewDouYinUser("https://v.douyin.com/qKDMXG/")
	if err != nil {
		wlog.Error("新建抖音用户失败:", err)
	}
	StoreDouYinUser(user)

}
