package main

import (
	"encoding/json"
	"porter/api"
	"porter/define"
	"porter/task"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
)

type TaskSearchKeyword struct{}

func (*TaskSearchKeyword) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	t := &task.TaskSearchKeyword{}
	if err := json.Unmarshal(m.Body, t); err != nil {
		wlog.Error("任务解析失败:", err)
		return nil
	}
	wlog.Info("开始处理关键字:", t.Keyword)
	cursor := 0

	for {
		if cursor >= t.Total {
			break
		}

		result, err := api.DouyinSearchKeyword(t.Keyword, cursor)
		if err != nil {
			wlog.Error("搜索发生错误:", err)
			break
		}

		// 搜索用户详情
		for _, item := range result.UserList {
			phoneUserInfo, err := api.NewPhoneDouyinUser(item.UserInfo.SecUID)
			if err != nil {
				wlog.Error("获取用户信息失败:", err)
				continue
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
				continue
			}

			Q.Publish(define.TaskAddDouyinUser, data)

			wlog.Infof("增加账号[%s]", phoneUserInfo.User.Nickname)
		}

		if result.HasMore == 0 && len(result.UserList) == 0 {
			break
		}

		cursor += 20
	}

	return nil
}
