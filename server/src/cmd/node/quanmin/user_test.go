package quanmin

import (
	"fmt"
	"os"
	"porter/define"
	"strconv"
	"testing"
)

// func TestGetVideoDuration(t *testing.T) {
// 	want := "57.946"

// 	result, err := getVideoDuration("D:/#龙城传奇 #是兄弟就一起来战 #我的沙城我做主 @+小助手.mp4")
// 	if err != nil {
// 		t.Errorf("获取视频长度错误:%s", err)
// 	}

// 	if result != want {
// 		t.Errorf("获取视频长度结果不一致: wangt[%s] result[%s]", want, result)
// 	}
// }

// func TestGenThunbsnails(t *testing.T) {
// 	err := genThumbnails("D:/#龙城传奇 #是兄弟就一起来战 #我的沙城我做主 @+小助手.mp4", "./thumbsnails/#龙城传奇 #是兄弟就一起来战 #我的沙城我做主 @+小助手.jpg")
// 	if err != nil {
// 		t.Errorf("生成缩略图失败: %s", err)
// 	}
// }

func TestCut(t *testing.T) {
	f, err := os.Open("D:/喜欢的视频.mp4")
	if err != nil {
		t.Error("读取错误:", err)
	}
	partNum := 1
	s := make([]byte, 4*define.MB)
	for {
		switch nr, err := f.Read(s); true {

		case nr < 0:
			t.Error("从视频文件中读取数据出错", err)
			return
		case nr == 0: // EOF
			return
		case nr > 0:
			nf, err := os.Create("./" + strconv.Itoa(partNum))
			defer nf.Close()

			if err != nil {
				t.Error("创建文件失败")
			}

			l, err := nf.Write(s[:nr])
			if err != nil {
				t.Error("写入失败")
			}
			fmt.Println("l:", l)
			if l != nr {
				t.Error("长度错误", l, nr)
			}
			partNum++
		}
	}
}
