package main

import (
	"net/http"
	"porter/define"

	"github.com/gin-gonic/gin"
)

func DouyinUserList(c *gin.Context) {
	param := &struct {
		UID   string `json:"uid,omitempty"`
		Page  int    `json:"page"`
		Limit int    `json:"limit"`
	}{Page: 1, Limit: 10}

	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	users := make([]*DouyinUser, 0)
	totalNum := int64(0)
	result := DB.Model(&DouyinUser{}).Count(&totalNum).Offset((param.Page - 1) * param.Limit).Limit(param.Limit).Order("create_time desc").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.QueryDataErr})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     define.Success,
		"users":    users,
		"totalNum": totalNum,
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
