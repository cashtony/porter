package task

import (
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

type TaskChangeInfo struct {
	List []TaskChangeInfoItem
}

type TaskChangeInfoItem struct {
	Bduss     string
	DouyinURL string
}

type TaskParseVideo struct {
	Type     define.ParseVideoType
	ShareURL string
}

type TaskParseVideoResult struct {
	DouyinNickname string
	DouyinUID      string
	List           []*define.TableDouyinVideo
}
