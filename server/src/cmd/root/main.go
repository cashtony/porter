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

var Mode = flag.String("mode", "release", "运行模式 debug:开发模式, release:产品模式")
var Host = flag.String("host", ":5000", "指定的地址")
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
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// db初始化
	DB = db.NewPG()
	DB.AutoMigrate(&DouyinVideo{})
	DB.AutoMigrate(&DouyinUser{})
	DB.AutoMigrate(&BaiduUser{})
	DB.AutoMigrate(&Account{})
	DB.AutoMigrate(&FaildRecords{})
	DB.AutoMigrate(&Statistic{})
	// 队列初始化
	comsumer := queue.InitComsumer(define.TaskFinishedTopic, &queueMsgHandler{})
	Q = queue.InitProducer()

	// 定时器初始化, 每天固定时间开始进行用户视频的检测
	c := cron.New()
	c.AddFunc("0 22 * * *", DailyUpdate)
	go c.Run()

	// gin初始化
	g := gin.Default()
	g.Static("/bg", "./www/")
	g.Static("/static", "./www/static")

	g.POST("/account/login", Login)
	g.POST("/account/logout", Logout)
	g.GET("/account/info", AccountInfo)

	g.POST("/douyin/user/list", DouyinUserList)
	g.POST("/baidu/user/list", BaiduUserList)
	g.POST("/baidu/user/edit", BaiduUserEdit)
	g.POST("/baidu/user/update", BaiduUserUpdate)
	g.POST("/bind/add", BindAdd)
	g.POST("/baidu/user/sync", SyncBaiduUser)
	g.POST("/statistic", GetStatistic)

	g.POST("/manage/manuallyDailyUpdate", ManuallyDailyUpdate)
	g.POST("/manage/manuallyNewlyUpdate", ManuallyNewVideoUpdate)

	g.Run(*Host)

	comsumer.Stop()
	c.Stop()
}
