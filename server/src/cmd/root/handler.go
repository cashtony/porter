package main

import (
	"net/http"
	"porter/define"
	"porter/wlog"
	"strings"

	"github.com/gin-gonic/gin"
)

func BindAdd(c *gin.Context) {
	param := &struct {
		Content string `json:"content"`
	}{}
	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	// 冒号和换行作为分割符
	// 单条格式为 抖音分享号:百度bduss
	//首先清理无用的空格之类符号
	cleanContent := strings.Replace(param.Content, " ", "", -1)
	if len(cleanContent) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	failRecords := make([]string, 0)
	// 解析成单行
	lines := strings.Split(cleanContent, "\n")
	sucNum, faildNum := 0, 0
	for _, line := range lines {
		if line == "" {
			continue
		}
		accounts := strings.Split(line, "\t")
		bduss, douyinShare := accounts[0], accounts[1]
		var baiduErr, douyinErr bool

		bd, err := NewBaiduUser(bduss)
		if err != nil {
			wlog.Error("百度bduss解析失败", bduss, err)
			baiduErr = true
			failRecords = append(failRecords, bduss)
		}
		dy, err := NewDouYinUser(douyinShare)
		if err != nil {
			wlog.Error("抖音用户解析失败", douyinShare, err)
			douyinErr = true
			failRecords = append(failRecords, douyinShare)
		}

		// 任意一个解析有问题就不进行绑定
		if baiduErr || douyinErr {
			storeFaild(bduss, douyinShare)
			faildNum++
			continue
		}
		wlog.Debugf("开始绑定%s %s \n", bd.Nickname, dy.Nickname)
		bd.DouyinUID = dy.UID
		dy.BaiduUID = bd.UID

		bd.Store()
		dy.Store()

		go func() {
			dy.initVideoList()
			publishTask(dy)
		}()

		sucNum++
	}

	c.JSON(http.StatusOK, gin.H{"code": define.Success, "sucNum": sucNum, "faildNum": faildNum})
}

func storeFaild(bduss, douyin string) {
	r := &FaildRecords{
		Bduss:  bduss,
		Douyin: douyin,
	}

	DB.Create(r)
}

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
		c.JSON(http.StatusOK, gin.H{"code": define.QueryErr})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     define.Success,
		"users":    users,
		"totalNum": totalNum,
	})
}

func BaiduUserList(c *gin.Context) {
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

	users := make([]*BaiduUser, 0)
	totalNum := int64(0)
	result := DB.Model(&BaiduUser{}).Count(&totalNum).Offset((param.Page - 1) * param.Limit).Limit(param.Limit).Order("create_time desc").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.QueryErr})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     define.Success,
		"users":    users,
		"totalNum": totalNum,
	})
}

func ImmediatelyUpdate(c *gin.Context) {
	UpdateAndUpload()
}

func ReloadUserVideoList(c *gin.Context) {
	param := &struct {
		ShareURL string `json:"shareURL"`
	}{}

	err := c.BindJSON(param)
	if err != nil {
		wlog.Error("参数解析错误", err)
		return
	}

	dy, _ := NewDouYinUser(param.ShareURL)
	go dy.initVideoList()

	c.JSON(http.StatusOK, gin.H{
		"code": define.Success,
	})
}
