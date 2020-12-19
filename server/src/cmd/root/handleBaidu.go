package main

import (
	"encoding/json"
	"net/http"
	"porter/api"
	"porter/define"
	"porter/wlog"
	"strings"
	"time"

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
	sucNum, total := 0, 0
	for _, line := range lines {
		if line == "" {
			continue
		}
		time.Sleep(300 * time.Millisecond)
		total++
		accounts := strings.Split(line, "\t")
		bduss, douyinShare := accounts[0], accounts[1]

		num := int64(0)
		result := DB.Model(&BaiduUser{}).Where("douyin_url = ?", douyinShare).Count(&num)
		if result.Error != nil {
			wlog.Errorf("新增账号时检测重复抖音号[%s]时出现问题:%s", douyinShare, result.Error)
			continue
		}
		if num != 0 {
			wlog.Info("抖音号[%s]已经被绑定,将跳过此绑定", douyinShare)
			continue
		}
		var baiduErr, douyinErr bool

		bd, err := NewBaiduUser(bduss)
		if err != nil {
			wlog.Error("百度bduss解析失败", bduss, err)
			baiduErr = true
			failRecords = append(failRecords, bduss)
		}
		apiDouyinUser, err := api.NewAPIDouyinUser(douyinShare)
		if err != nil {
			wlog.Error("抖音用户解析失败", douyinShare, err)
			douyinErr = true
			failRecords = append(failRecords, douyinShare)
		}

		// 任意一个解析有问题就不进行绑定
		if baiduErr || douyinErr {
			storeFaild(bduss, douyinShare)
			continue
		}
		wlog.Debugf("开始绑定%s %s \n", bd.Nickname, apiDouyinUser.Nickname)
		bd.DouyinURL = douyinShare

		tableDouyinUser := &DouyinUser{
			UID:        apiDouyinUser.UID,
			UniqueUID:  apiDouyinUser.UniqueUID,
			Nickname:   apiDouyinUser.Nickname,
			ShareURL:   douyinShare,
			VideoCount: apiDouyinUser.AwemeCount,
			FansCount:  apiDouyinUser.FollowerCount,
		}
		bd.Store()
		tableDouyinUser.Store()

		ThreadTraffic <- 1
		go func() {
			tableDouyinUser.initVideoList()
			// bd.UploadVideo(UploadTypeDaily)
			<-ThreadTraffic
		}()

		sucNum++
	}

	c.JSON(http.StatusOK, gin.H{"code": define.Success, "total": total, "sucNum": sucNum})
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
		c.JSON(http.StatusOK, gin.H{"code": define.QueryDataErr})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     define.Success,
		"users":    users,
		"totalNum": totalNum,
	})
}

func BaiduUserEdit(c *gin.Context) {
	param := &struct {
		UID       string `json:"uid"`
		DouyinURL string `json:"douyinURL"`
	}{}

	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	result := DB.Model(&BaiduUser{}).Where("uid = ?", param.UID).Update("douyin_url", param.DouyinURL)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.CannotBind})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": define.Success,
	})
}

func BaiduUserUpdate(c *gin.Context) {
	param := &struct {
		UID string `json:"uid,omitempty"`
	}{}

	err := c.BindJSON(param)
	if err != nil {
		wlog.Error("参数解析错误", err)
		return
	}

	usersList := make([]*BaiduUser, 0)
	subDB := DB.Model(&BaiduUser{})
	if param.UID != "" {
		subDB.Where("uid = ?", param.UID)
	}

	subDB.Find(&usersList)

	for _, u := range usersList {
		err := u.fetchQuanminInfo()
		if err != nil {
			wlog.Errorf("获取[%s][%s]全民视频用户数据时错误:%s", u.UID, u.Nickname, err)
			continue
		}
		DB.Model(&BaiduUser{}).Where("uid = ?", u.UID).Updates(&BaiduUser{Diamond: u.Diamond})
		time.Sleep(300 * time.Millisecond)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": define.Success,
	})
}

func SyncBaiduUser(c *gin.Context) {
	param := &struct {
		Content string `json:"content"`
	}{}

	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	if param.Content == "" {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
	}

	items := strings.Split(param.Content, "\n")
	list := make([]define.TaskChangeInfoItem, 0)

	for _, value := range items {
		if value == "" {
			continue
		}
		value = strings.TrimSpace(value)
		item := strings.Split(value, "\t")
		if len(item) == 2 {
			list = append(list, define.TaskChangeInfoItem{
				Bduss:     item[0],
				DouyinURL: item[1],
			})

		}
	}

	data, err := json.Marshal(&define.TaskChangeInfo{
		List: list,
	})
	if err != nil {
		wlog.Error("task解析成json错误", err)
		return
	}

	err = Q.Publish(define.TaskChangeInfoTopic, data)
	if err != nil {
		wlog.Error("任务发布失败:", err)
	}

	c.JSON(http.StatusOK, gin.H{"code": define.Success})
}

func ChangeBaiduUserStatus(c *gin.Context) {
	param := &struct {
		UID    string `json:"uid"`
		Status int    `json:"status"`
	}{}

	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	result := DB.Model(&BaiduUser{}).Where("uid = ?", param.UID).Update("status", param.Status)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.CannotBind})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": define.Success,
	})
}
