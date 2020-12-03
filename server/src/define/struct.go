package define

import (
	"fmt"
	"time"
)

type Task struct {
	Bduss    string
	Nickname string
	Videos   []*TaskVideo
}
type TaskVideo struct {
	AwemeID     string
	Desc        string
	DownloadURL string
}

type TaskFinished struct {
	AwemeID string
}

type JsonTime time.Time

// MarshalJSON 实现它的json序列化方法
func (this *JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(*this).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}
