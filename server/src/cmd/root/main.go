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
var Host = flag.String("host", ":5000", "指定的地址")
var DB *gorm.DB
var Q *nsq.Producer

func main() {
	rand.Seed(time.Now().Unix())
	flag.Parse()

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

	// 定时器初始化, 每天早上8点开始进行用户视频的检测
	c := cron.New()
	c.AddFunc("0 8 * * *", DailyUpdate)

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

	g.POST("/statistic", GetStatistic)

	g.POST("/manage/manuallyDailyUpdate", ManuallyDailyUpdate)
	g.POST("/manage/manuallyNewlyUpdate", ManuallyNewVideoUpdate)

	g.Run(*Host)

	comsumer.Stop()

}
