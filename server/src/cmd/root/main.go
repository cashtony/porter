package main

import (
	"flag"
	"math/rand"
	"porter/db"
	"porter/define"
	"porter/queue"
	"porter/wlog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var Mode = flag.String("mode", "debug", "运行模式 debug:开发模式, release:产品模式")
var DB *gorm.DB
var Q *nsq.Producer

func main() {
	rand.Seed(time.Now().Unix())
	flag.Parse()

	if *Mode == "debug" {
		wlog.DevelopMode()
	}
	// db初始化
	DB = db.NewPG()
	DB.AutoMigrate(&DouyinVideo{})
	DB.AutoMigrate(&DouyinUser{})
	DB.AutoMigrate(&BaiduUser{})

	// 队列初始化
	comsumer := queue.InitComsumer(define.TaskFinishedTopic, &queueMsgHandler{})
	Q = queue.InitProducer()

	// 定时器初始化, 每天早上8点开始进行用户视频的检测
	c := cron.New()
	c.AddFunc("0 8 * * *", ScheduleUpdate)
	// gin初始化

	// user, err := NewDouYinUser("https://v.douyin.com/qKDMXG/")
	// if err != nil {
	// 	wlog.Error("新建抖音用户失败:", err)
	// }
	// user.Store()1460800875 380917571

	// user := &BaiduUser{
	// 	UID: "1",
	// }

	// ScheduleUpdate()

	g := gin.Default()
	g.Run()

	comsumer.Stop()

}
