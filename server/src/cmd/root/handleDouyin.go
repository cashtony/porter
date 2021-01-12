package main

import (
	"encoding/json"
	"net/http"
	"porter/define"
	"porter/task"
	"porter/wlog"
	"strings"

	"github.com/gin-gonic/gin"
)

func DouyinUserList(c *gin.Context) {
	param := &struct {
		UID         string `json:"uid,omitempty"`
		Page        int    `json:"page"`
		Limit       int    `json:"limit"`
		DouyinUID   string `json:"douyinUID"`
		Nickname    string `json:"nickname"`
		HideSimilar bool   `json:"hideSimilar"`
		HideBinded  bool   `json:"hideBinded"`
	}{Page: 1, Limit: 10}

	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	users := make([]*TableDouyinUser, 0)
	totalNum := int64(0)

	query := DB.Model(&TableDouyinUser{})
	if param.DouyinUID != "" {
		query.Where("uid = ?", param.DouyinUID)
	}
	if param.Nickname != "" {
		query.Where("nickname = ?", param.Nickname)
	}
	if param.HideBinded {
		query.Where("baidu_uid = ''")
	}

	if param.HideSimilar {
		query.Where("similar_sign = 0")
	}

	result := query.Count(&totalNum).Offset((param.Page - 1) * param.Limit).Limit(param.Limit).Order("create_time desc").Find(&users)
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

func AddDouyinUser(c *gin.Context) {
	param := &struct {
		Content string `json:"content"`
	}{}
	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr, "message": "参数错误"})
		return
	}

	// 	//首先清理无用的空格之类符号
	cleanContent := strings.Replace(param.Content, " ", "", -1)
	if len(cleanContent) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr, "message": "参数错误"})
		return
	}

	// 解析成单行
	lines := strings.Split(cleanContent, "\n")

	for _, shareURL := range lines {
		if shareURL == "" {
			continue
		}

		t := &task.TaskParseDouyinURL{
			DouyinURL: shareURL,
		}
		data, err := json.Marshal(t)
		if err != nil {
			wlog.Error("json解析失败:", err)
			continue
		}

		Q.Publish(define.TaskParseDouyinURL, data)
	}

	c.JSON(http.StatusOK, gin.H{"code": define.Success})
}

func DeleteDouyinUser(c *gin.Context) {
	param := &struct {
		UID string `json:"uid"`
	}{}
	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	if result := DB.Model(&TableBaiduUser{}).Where("douyin_uid = ?", param.UID).Updates(TableBaiduUser{DouyinUID: ""}); result.Error != nil {
		wlog.Error("删除抖音用户时更新相应的百度账号失败:", result.Error)
	}

	result := DB.Where("uid = ?", param.UID).Delete(&TableDouyinUser{}).Limit(1)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.CannotDelete, "message": "删除时发生错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": define.Success,
	})
}

func DouyinUserSearch(c *gin.Context) {
	param := &struct {
		Content string `json:"content"`
		Total   int    `json:"total"`
	}{}
	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	// 	//首先清理无用的空格之类符号
	cleanContent := strings.Replace(param.Content, " ", "", -1)
	if len(cleanContent) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	// 解析成单行
	keywords := strings.Split(cleanContent, "\n")

	for _, keyword := range keywords {
		// 发布给服务端
		fin := &task.TaskSearchKeyword{
			Keyword: keyword,
			Total:   param.Total,
		}
		data, err := json.Marshal(fin)
		if err != nil {
			wlog.Error("解析数据时错误:", err)
			return
		}

		Q.Publish(define.TaskSearchKeyword, data)
	}

	c.JSON(http.StatusOK, gin.H{"code": define.Success})
}
