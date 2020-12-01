package main

import (
	"flag"
	"os"
	"os/signal"
	"porter/define"
	"porter/queue"
	"porter/wlog"
	"syscall"

	"github.com/nsqio/go-nsq"
	"gorm.io/gorm"
)

var Mode = flag.String("mode", "debug", "运行模式 debug:开发模式, release:产品模式")
var Host = flag.String("host", ":1213", "指定主机地址")
var bduss = flag.String("bduss", "", "百度BDUSS")
var dyuid = flag.String("uid", "", "抖音uid")
var DB *gorm.DB
var Q *nsq.Producer

func main() {
	flag.Parse()

	if *Mode == "debug" {
		wlog.DevelopMode()
	}

	// 从消息队列中获取任务
	comsumer := queue.InitComsumer(define.TaskPushTopic, &queueMsgHandler{})
	Q = queue.InitProducer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	comsumer.Stop()
}
