package task

import (
	"porter/api"
	"porter/define"
)

type TaskUpload struct {
	Bduss    string
	Nickname string
	Videos   []*TaskUploadVideo
}
type TaskUploadVideo struct {
	AwemeID     string
	Desc        string
	DownloadURL string
}

type TaskUploadFinished struct {
	AwemeID string
}

type TaskChangeInfoItem struct {
	Bduss     string
	Avatar    string
	Nickname  string
	Gender    int
	Signature string
	Birthday  string
	Location  string
	Province  string
	City      string
}

type TaskParseVideo struct {
	Type   define.ParseVideoType
	SecUID string
}

type TaskParseVideoResult struct {
	DouyinNickname string
	DouyinUID      string
	List           []*define.TableDouyinVideo
}

type TaskSearchKeyword struct {
	Keyword string
	Total   int // 累计搜索多少条
}

type TaskParseDouyinURL struct {
	DouyinURL string
}

type TaskAddDouyinUser struct {
	*api.APIPhoneDouyinUser
	HasSimilar bool
}
