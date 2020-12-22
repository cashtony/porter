package define

import (
	"fmt"
	"strconv"
	"time"
)

type JsonTime time.Time

// MarshalJSON 实现它的json序列化方法
func (this *JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", time.Time(*this).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

type StampTime time.Time

func (this *StampTime) UnmarshalJSON(b []byte) error {
	stamp, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(this) = time.Unix(stamp, 0)

	return nil
}

func (this *StampTime) IsZero() bool {
	return (*time.Time)(this).IsZero()
}

type TableDouyinVideo struct {
	AwemeID    string `gorm:"primaryKey"`
	DouyinUID  string // 抖音UID
	Desc       string // 视频描述
	Vid        string // 用于下载时填充链接
	CreateTime time.Time
	Duration   int
	State      int // 0未搬运 1:已搬运
}

func (*TableDouyinVideo) TableName() string {
	return "douyin_videos"
}

type DouyinVideoExtraInfo struct {
	CreateTime time.Time
	VID        string
}
