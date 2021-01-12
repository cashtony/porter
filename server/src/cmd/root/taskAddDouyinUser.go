package main

import (
	"encoding/json"
	"porter/task"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
)

type taskAddDouyinUser struct{}

func (t *taskAddDouyinUser) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}
	taskAddDouyinUser := &task.TaskAddDouyinUser{}
	err := json.Unmarshal(m.Body, taskAddDouyinUser)
	if err != nil {
		wlog.Errorf("上传事件解析失败:%s \n", err)
		return nil
	}

	similarSign := 0
	if taskAddDouyinUser.HasSimilar {
		similarSign = 1
	}
	douyinUserInfo := taskAddDouyinUser.APIPhoneDouyinUser.User
	tableDouyinUser := &TableDouyinUser{
		UID:            douyinUserInfo.UID,
		UniqueUID:      douyinUserInfo.UniqueID,
		Nickname:       douyinUserInfo.Nickname,
		AwemeCount:     douyinUserInfo.AwemeCount,
		FollowerCount:  douyinUserInfo.FollowerCount,
		Gender:         douyinUserInfo.Gender,
		Signature:      douyinUserInfo.Signature,
		Birthday:       douyinUserInfo.Birthday,
		TotalFavorited: douyinUserInfo.TotalFavorited,
		Location:       douyinUserInfo.Location,
		Province:       douyinUserInfo.Province,
		City:           douyinUserInfo.City,
		Avatar:         douyinUserInfo.Avatar.URLList[0],
		SecUID:         douyinUserInfo.SecUID,
		SimilarSign:    similarSign,
	}

	tableDouyinUser.Store()

	return nil
}
