package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"porter/define"
	"porter/queue"
	"porter/wlog"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
	"gorm.io/gorm"
)

var Mode = flag.String("mode", "debug", "运行模式 debug:开发模式, release:产品模式")
var Host = flag.String("host", ":1213", "指定主机地址")
var Thread = flag.Int("thread", 16, "同时运行任务数量")
var ThreadTraffic chan int

var DB *gorm.DB
var Q *nsq.Producer

func main() {
	rand.Seed(time.Now().Unix())
	flag.Parse()

	ThreadTraffic = make(chan int, *Thread)

	if *Mode == "debug" {
		wlog.DevelopMode()
	}
	// 读取一些配置文件
	InitConfigFile()
	// 从消息队列中获取任务
	uploadComsumer := queue.InitComsumer(define.TaskPushTopic, &TaskUploadHandler{})
	changeInfoComsumer := queue.InitComsumer(define.TaskChangeInfoTopic, &TaskChangeInfoHandler{})
	parseVideoComsumer := queue.InitComsumer(define.TaskParseVideoTopic, &TaskParseVideoHandler{})

	Q = queue.InitProducer()
	wlog.Info("当前设定Thread数量为:", *Thread)
	wlog.Info("等待新任务中...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	uploadComsumer.Stop()
	changeInfoComsumer.Stop()
	parseVideoComsumer.Stop()
}
