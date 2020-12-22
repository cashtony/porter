package main

import (
	"encoding/json"
	"porter/task"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
	"gorm.io/gorm/clause"
)

type taskParseVideoResult struct{}

func (t *taskParseVideoResult) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}
	result := &task.TaskParseVideoResult{}
	err := json.Unmarshal(m.Body, result)
	if err != nil {
		wlog.Errorf("视频解析事件解析失败:%s \n", err)
		return nil
	}

	if len(result.List) != 0 {
		// 将没有的视频传入到数据库中
		DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "aweme_id"}},
			DoNothing: true,
		}).Create(result.List)
		if DB.Error != nil {
			wlog.Errorf("用户[%s][%s]新视频信息存入数据库失败:%s \n", result.DouyinUID, result.DouyinNickname, DB.Error)
			return nil
		}
	}

	return nil
}
