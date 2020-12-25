package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"porter/api"
	"porter/define"
	"porter/task"
	"porter/wlog"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
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

	// 解析成单行
	lines := strings.Split(cleanContent, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		accounts := strings.Split(line, "\t")
		bduss, douyinShare := accounts[0], accounts[1]
		errcode := 0
		apiDouyinUser, err := api.NewAPIDouyinUser(douyinShare)
		if err != nil {
			wlog.Info("抖音用户解析失败", douyinShare, err)
			errcode = BindErrDouyinUser
		}

		num := int64(0)
		result := DB.Model(&TableBaiduUser{}).Where("douyin_uid = ?", apiDouyinUser.UID).Count(&num)
		if result.Error != nil {
			wlog.Errorf("新增账号时检测重复抖音号[%s]时出现问题:%s", douyinShare, result.Error)
			errcode = BindErrSqlQuery
		}
		if num != 0 {
			wlog.Infof("抖音号[%s]已经被绑定,将跳过此绑定", apiDouyinUser.UID)
			errcode = BindErrAlreadyBind
		}

		bd, err := NewBaiduUser(bduss)
		if err != nil {
			wlog.Info("百度bduss解析失败:", bduss, err)
			errcode = BindErrBdussWrong
		}

		// 任意一个解析有问题就不进行绑定
		if errcode != 0 {
			storeFaild(bduss, douyinShare, errcode)
			continue
		}
		wlog.Debugf("开始绑定[%s]和[%s] \n", bd.Nickname, apiDouyinUser.Nickname)
		bd.DouyinUID = apiDouyinUser.UID
		bd.DouyinURL = douyinShare
		bd.Status = int(BaiduUserStatusNormal)

		tableDouyinUser := &TableDouyinUser{
			UID:        apiDouyinUser.UID,
			UniqueUID:  apiDouyinUser.UniqueUID,
			Nickname:   apiDouyinUser.Nickname,
			ShareURL:   douyinShare,
			VideoCount: apiDouyinUser.AwemeCount,
			FansCount:  apiDouyinUser.FollowerCount,
		}
		bd.Store()
		tableDouyinUser.Store()

		// 发布任务
		t := &task.TaskParseVideo{
			Type:     define.ParseVideoTypeAll,
			ShareURL: douyinShare,
		}
		data, err := json.Marshal(t)
		if err != nil {
			wlog.Error("任务解析失败:", err)
			continue
		}

		if err := Q.Publish(define.TaskParseVideoTopic, data); err != nil {
			wlog.Error("视频解析任务发布失败:", err)
		}

	}

	c.JSON(http.StatusOK, gin.H{"code": define.Success})
}

func BaiduUserList(c *gin.Context) {
	param := &struct {
		UID       string `json:"uid,omitempty"`
		DouyinUID string `json:"douyinUID,omitempty"`
		DouyinURL string `json:"douyinURL,omitempty"`
		BDUSS     string `json:"bduss,omitempty"`
		Nickname  string `json:"nickname,omitempty"`
		Page      int    `json:"page"`
		Limit     int    `json:"limit"`
	}{Page: 1, Limit: 10}

	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	users := make([]*TableBaiduUser, 0)
	totalNum := int64(0)
	subDB := DB.Model(&TableBaiduUser{})
	if param.DouyinUID != "" {
		subDB.Where("douyin_uid = ?", param.DouyinUID)
	}
	if param.DouyinURL != "" {
		subDB.Where("douyin_url = ?", param.DouyinURL)
	}
	if param.Nickname != "" {
		subDB.Where("nickname like ?", "%"+param.Nickname+"%")
	}

	result := subDB.
		Count(&totalNum).
		Offset((param.Page - 1) * param.Limit).
		Limit(param.Limit).
		Order("create_time desc").Find(&users)
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
	douyinUser, err := api.NewAPIDouyinUser(param.DouyinURL)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.WrongDouyinShareURL})
		return
	}
	num := int64(0)
	if result := DB.Model(&TableBaiduUser{}).Where("douyin_uid = ?", douyinUser.UID).Count(&num); result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.QueryDataErr, "message": "数据查找错误,请稍后再试"})
		return
	}
	if num > 0 {
		c.JSON(http.StatusOK, gin.H{"code": define.AlreadyBind, "message": "该抖音账号已经绑定其他百度账号"})
		return
	}
	if result := DB.Model(&TableBaiduUser{}).Where("uid = ?", param.UID).Updates(TableBaiduUser{DouyinUID: douyinUser.UID, DouyinURL: param.DouyinURL}); result.Error != nil {
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

	DailyUpdateBaiduUsers()

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
	list := make([]task.TaskChangeInfoItem, 0)

	for _, value := range items {
		if value == "" {
			continue
		}
		value = strings.TrimSpace(value)
		item := strings.Split(value, "\t")
		if len(item) == 2 {
			list = append(list, task.TaskChangeInfoItem{
				Bduss:     item[0],
				DouyinURL: item[1],
			})

		}
	}

	data, err := json.Marshal(&task.TaskChangeInfo{
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

	result := DB.Model(&TableBaiduUser{}).Where("uid = ?", param.UID).Update("status", param.Status)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.CannotBind})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": define.Success,
	})
}

func BaiduUserDelete(c *gin.Context) {
	param := &struct {
		UID string `json:"uid"`
	}{}
	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	if param.UID == "" {
		c.JSON(http.StatusOK, gin.H{"code": define.ParamErr})
		return
	}

	result := DB.Debug().Where("uid = ?", param.UID).Delete(&TableBaiduUser{}).Limit(1)
	if result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"code": define.CannotDelete, "message": "删除时发生错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": define.Success,
	})
}

func ExcelBaiduUsers(c *gin.Context) {
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("百度用户")

	row := sheet.AddRow()

	cell := row.AddCell()
	cell.Value = "UID"

	cell = row.AddCell()
	cell.Value = "昵称"

	cell = row.AddCell()
	cell.Value = "BDUSS"

	cell = row.AddCell()
	cell.Value = "抖音URL"

	cell = row.AddCell()
	cell.Value = "抖音UID"

	cell = row.AddCell()
	cell.Value = "加入时间"

	cell = row.AddCell()
	cell.Value = "钻石数量"

	cell = row.AddCell()
	cell.Value = "粉丝数量"

	usersList := make([]*TableBaiduUser, 0)
	DB.Model(&TableBaiduUser{}).Find(&usersList)

	for _, u := range usersList {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.Value = u.UID

		cell = row.AddCell()
		cell.Value = u.Nickname

		cell = row.AddCell()
		cell.Value = u.Bduss

		cell = row.AddCell()
		cell.Value = u.DouyinURL

		cell = row.AddCell()
		cell.Value = u.DouyinUID

		cell = row.AddCell()
		cell.Value = u.CreateTime.String()

		cell = row.AddCell()
		cell.Value = strconv.Itoa(u.Diamond)

		cell = row.AddCell()
		cell.Value = strconv.Itoa(u.FansNum)
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	file.Write(writer)
	theBytes := b.Bytes()

	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Content-Disposition", "attachment; filename=百度用户.xlsx")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(theBytes)))
	c.Writer.Write([]byte(theBytes))
}
