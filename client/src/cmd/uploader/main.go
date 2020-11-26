package main

import (
	"flag"
	"porter/cmd/uploader/quanmin"
	"porter/wlog"

	"gorm.io/gorm"
)

var Mode = flag.String("mode", "debug", "运行模式 debug:开发模式, release:产品模式")
var Host = flag.String("host", ":1213", "指定主机地址")
var bduss = flag.String("bduss", "", "百度BDUSS")
var dyuid = flag.String("uid", "", "抖音uid")
var DB *gorm.DB

func main() {
	flag.Parse()

	if *Mode == "debug" {
		wlog.DevelopMode()
	}
	// DB = db.NewPG()

	// r := gin.Default()
	// r.GET("/bind", Bind)

	wlog.Infof("开始解析用户[%s]视频列表 \n", *dyuid)
	vlist := UserVideoList(*dyuid)
	wlog.Infof("用户[%s]视频列表解析成功 \n", *dyuid)
	// vlist := Mock()
	q := quanmin.NewUser(*bduss)
	for i, v := range vlist[:5] {
		wlog.Infof("[%d][%s]开始下载 \n", i+1, v.Desc)
		// 下载视频
		download(v)
		wlog.Infof("[%d][%s]下载结束,开始上传 \n", i+1, v.Desc)
		q.Upload(v.LocalFilePath, v.Desc)

		wlog.Infof("[%s]上传完毕 \n", v.Desc)
	}
	// err := q.Login()
	// if err != nil {
	// 	wlog.Error("登录失败", err)
	// 	return
	// }

	// wlog.Info("登录成功:", q.Name())
	// file := make([]byte, 0)
	// err := q.Upload(file)
	// if err != nil {
	// 	wlog.Error("文件上传失败:", err)
	// 	return
	// }

	// r.Run(*Host)

}
