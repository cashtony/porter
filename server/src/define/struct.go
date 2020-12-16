package define

import (
	"fmt"
	"time"
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
type JsonTime time.Time

// MarshalJSON 实现它的json序列化方法
func (this *JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(*this).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}
