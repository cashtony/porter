package main

import (
	"encoding/json"
	"porter/api"
	"porter/define"
	"porter/task"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
)

type TaskParseDouyinURL struct{}

func (*TaskParseDouyinURL) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	t := &task.TaskParseDouyinURL{}
	if err := json.Unmarshal(m.Body, t); err != nil {
		wlog.Error("任务解析失败:", err)
		return nil
	}
	wlog.Info("开始解析:", t.DouyinURL)
	// 解析url获取secuid, 然后通过secuid来获取手机版用户数据
	secUID := api.GetSecID(t.DouyinURL)
	if secUID == "" {
		wlog.Error("secuid 获取失败", t.DouyinURL)
		return nil
	}

	phoneUserInfo, err := api.NewPhoneDouyinUser(secUID)
	if err != nil {
		wlog.Error("获取用户信息失败:", err)
		return nil
	}
	var hasSimilar bool
	// 检查全民那边是否有相同的号
	hasSimilar, err = IsSimilarInQuanmin(phoneUserInfo.User.Nickname)
	if err != nil {
		wlog.Error("查找相似名字时出现错误:", err)
	}

	// 发布给服务端
	fin := &task.TaskAddDouyinUser{
		HasSimilar:         hasSimilar,
		APIPhoneDouyinUser: phoneUserInfo,
	}

	data, err := json.Marshal(fin)
	if err != nil {
		wlog.Error("解析数据时错误:", err)
		return nil
	}

	Q.Publish(define.TaskAddDouyinUser, data)

	return nil
}
