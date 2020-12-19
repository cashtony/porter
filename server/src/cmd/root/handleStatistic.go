package main

import (
	"net/http"
	"porter/define"
	"porter/wlog"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetStatistic(c *gin.Context) {
	param := &struct {
		UID      string `json:"uid"`
		PickDate string `json:"pickDate"`
		Page     int    `json:"page"`
		Limit    int    `json:"limit"`
	}{Page: 1, Limit: 10}

	err := c.BindJSON(param)
	if err != nil {
		wlog.Error("参数解析错误", err)
		return
	}

	totalNum := int64(0)
	list := make([]*StatisticRough, 0)
	subDB := DB.Model(&Statistic{}).
		Select("date(upload_time) as date, baidu_users.nickname as nickname, baidu_uid as uid, count(id) as num").
		Group("baidu_uid, nickname").
		Joins("left join baidu_users on statistics.baidu_uid = baidu_users.uid").Where("state = ?", 1)
	if param.UID != "" {
		subDB = subDB.Where("baidu_uid = ?", param.UID)
	}
	if param.PickDate != "" {
		subDB = subDB.Where("date(upload_time) = ?", param.PickDate)
	} else {
		subDB = subDB.Where("date(upload_time) = ?", time.Now().Format("2006-01-02"))
	}

	// 这里gorm有个bug, 如果直接使用一个session进行查询数量的时候生成的语句是有问题的
	subDB.Session(&gorm.Session{}).Count(&totalNum)
	subDB.Order("num desc").Group("date").Find(&list)

	c.JSON(http.StatusOK, gin.H{
		"code":     define.Success,
		"list":     list,
		"totalNum": totalNum,
	})
}
func storeFaild(bduss, douyin string) {
	r := &FaildRecords{
		Bduss:  bduss,
		Douyin: douyin,
	}

	DB.Create(r)
}
