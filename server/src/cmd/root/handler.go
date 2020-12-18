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
	"gorm.io/gorm"
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
		apiDouyinUser, err := api.NewAPIDouyinUser(douyinShare)
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

		go func() {
			tableDouyinUser.initVideoList()
			bd.UploadVideo(UploadTypeDaily)
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
		c.JSON(http.StatusOK, gin.H{"code": define.QueryDataErr})
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
		time.Sleep(200 * time.Millisecond)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": define.Success,
	})
}

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

	items := strings.Split(param.Content, "\r\n")
	list := make([]define.TaskChangeInfoItem, 0)

	for _, value := range items {
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
