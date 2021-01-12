package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"porter/define"
	"porter/task"
	"porter/wlog"

	"github.com/gin-gonic/gin"
)

func BindUser(c *gin.Context) {
	param := &struct {
		Content []string `json:"content"`
	}{}

	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr, "message": "参数错误"})
		return
	}

	resultMsg := make(map[string]string, 0)
	for _, uid := range param.Content {
		douyinUser := &TableDouyinUser{}
		if result := DB.Model(&TableDouyinUser{}).Where("uid = ?", uid).First(douyinUser); result.Error != nil {
			wlog.Error(result.Error)
			resultMsg[uid] = "获取抖音用户数据时操作数据库失败"
			continue
		}
		if douyinUser.BaiduUID != "" {
			resultMsg[uid] = fmt.Sprintf("抖音号[%s]已经被绑定", douyinUser.Nickname)
			continue
		}

		baiduUser := &TableBaiduUser{}
		baiduAvalibleNum := int64(0)
		subDB := DB.Model(&TableBaiduUser{}).Where("douyin_uid = ''").Count(&baiduAvalibleNum)

		if baiduAvalibleNum <= 0 {
			resultMsg[uid] = "已经没有可用的百度账户了"
			continue
		}

		if result := subDB.First(baiduUser); result.Error != nil {
			wlog.Error(result.Error)
			resultMsg[uid] = "获取百度用户数据时操作数据库失败"
			continue
		}

		t := &task.TaskParseVideo{
			Type:   define.ParseVideoTypeAll,
			SecUID: douyinUser.SecUID,
		}
		data, err := json.Marshal(t)
		if err != nil {
			wlog.Error("任务数据解析失败:", err)
			resultMsg[uid] = "数据解析失败"
			continue
		}

		if err := Q.Publish(define.TaskParseVideoTopic, data); err != nil {
			wlog.Error("视频解析任务发布失败:", err)
			resultMsg[uid] = "发布解析视频任务失败"
			continue
		}

		// 绑定用户数据
		if result := DB.Model(&TableBaiduUser{}).Where("uid = ?", baiduUser.UID).Updates(TableBaiduUser{DouyinUID: douyinUser.UID}); result.Error != nil {
			wlog.Error("更新数据库中百度用户所绑定的抖音号时失败:", result.Error)
			resultMsg[uid] = "绑定百度用户数据时发生错误"
			continue
		}
		if result := DB.Model(&TableDouyinUser{}).Where("uid = ?", douyinUser.UID).Updates(TableDouyinUser{BaiduUID: baiduUser.UID}); result.Error != nil {
			wlog.Error("更新数据库中抖音用户所绑定的百度号时失败:", result.Error)
			resultMsg[uid] = "绑定抖音用户数据时发生错误"
			continue
		}

		// 同步全民账号的头像,昵称,地区等信息
		taskChangeInfo := &task.TaskChangeInfoItem{
			Bduss:     baiduUser.Bduss,
			Avatar:    douyinUser.Avatar,
			Nickname:  douyinUser.Nickname,
			Gender:    douyinUser.Gender,
			Signature: douyinUser.Signature,
			Birthday:  douyinUser.Birthday,
			Location:  douyinUser.Location,
			Province:  douyinUser.Province,
			City:      douyinUser.City,
		}

		data, err = json.Marshal(taskChangeInfo)
		if err != nil {
			wlog.Error("task解析成json错误", err)
			continue
		}

		err = Q.Publish(define.TaskChangeInfoTopic, data)
		if err != nil {
			wlog.Error("任务发布失败:", err)
			continue
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":      define.Success,
		"resultMsg": resultMsg,
	})
}

func ManuallyDailyUpdate(c *gin.Context) {
	if isUpdading {
		c.JSON(http.StatusOK, gin.H{
			"code": define.AlreadyUpdating,
		})

		return
	}

	go DailyUpload()

	c.JSON(http.StatusOK, gin.H{
		"code": define.Success,
	})
}

func ManuallyNewVideoUpdate(c *gin.Context) {
	if isUpdading {
		c.JSON(http.StatusOK, gin.H{
			"code": define.AlreadyUpdating,
		})

		return
	}

	go NewlyUpload()

	c.JSON(http.StatusOK, gin.H{
		"code": define.Success,
	})
}

func ReloadUserVideoList(c *gin.Context) {
	param := &struct {
		ShareURL string `json:"shareURL"`
	}{}

	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": define.Success,
	})
}
