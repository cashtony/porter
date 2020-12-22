package main

import (
	"encoding/json"
	"porter/api"
	"porter/define"
	"porter/task"
	"porter/wlog"
	"time"

	"github.com/nsqio/go-nsq"
)

type TaskParseVideoHandler struct{}

func (*TaskParseVideoHandler) HandleMessage(m *nsq.Message) error {
	if len(m.Body) == 0 {
		return nil
	}

	parseVideo := &task.TaskParseVideo{}
	err := json.Unmarshal(m.Body, parseVideo)
	if err != nil {
		wlog.Error("任务解析失败:", err)
		return nil
	}

	ThreadTraffic <- 1
	go excuteParseVideo(parseVideo.Type, parseVideo.ShareURL)

	return nil
}

func excuteParseVideo(t define.ParseVideoType, shareURL string) {
	defer func() {
		<-ThreadTraffic
	}()

	var (
		nextCursor int64 = 0
		page             = 1
	)
	apiDouyinUser, err := api.NewAPIDouyinUser(shareURL)
	if err != nil {
		wlog.Error("解析抖音用户数据失败:", shareURL, err)
		return
	}

	secUID := api.GetSecID(shareURL)
	secSig := GetSecSig(shareURL)

	for {
		wlog.Debugf("开始解析[%s]第[%d]页视频 \n", filterSpecial(apiDouyinUser.Nickname), page)
		apiVideoList, newCursor, err := api.GetDouyinVideo(secUID, secSig, nextCursor)
		if err != nil {
			wlog.Error("获取单页视频发生错误:", err)
			return
		}

		wlog.Debugf("[%s]第[%d]页视频解析完毕 newCursor:%d videoLen:%d \n", filterSpecial(apiDouyinUser.Nickname), page, newCursor, len(apiVideoList))

		tableVideoList := make([]*define.TableDouyinVideo, 0)
		for _, v := range apiVideoList {
			videoExtraInfo, _ := api.GetVideoExtraInfo(v.AwemeID)
			tableVideoList = append(tableVideoList, &define.TableDouyinVideo{
				AwemeID:    v.AwemeID,
				Desc:       v.Desc,
				DouyinUID:  apiDouyinUser.UID,
				Vid:        v.Video.Vid,
				Duration:   v.Video.Duration,
				CreateTime: time.Unix(videoExtraInfo.ItemList[0].CreateTime, 0),
			})
		}

		result := task.TaskParseVideoResult{
			DouyinNickname: apiDouyinUser.Nickname,
			DouyinUID:      apiDouyinUser.UID,
			List:           tableVideoList,
		}

		data, err := json.Marshal(result)
		if err != nil {
			wlog.Error("传递视频解析结果失败:", err)
		}

		if err := Q.Publish(define.TaskParseVideoResultTopic, data); err != nil {
			wlog.Error("传递视频解析结果失败:", err)
		}

		if t == define.ParseVideoTypeOnePage {
			break
		}

		if newCursor == 0 {
			break
		}

		nextCursor = newCursor
		page++
	}

	wlog.Debugf("用户[%s]视频解析完毕 \n", apiDouyinUser.Nickname)
}
