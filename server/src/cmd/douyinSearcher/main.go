package main

import (
	"flag"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"porter/define"
	"porter/queue"
	"porter/wlog"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

var proxy = func(_ *http.Request) (*url.URL, error) {
	return url.Parse("http://127.0.0.1:8888")
}
var Mode = flag.String("mode", "debug", "运行模式 debug:开发模式, release:产品模式")
var Q *nsq.Producer

func main() {
	rand.Seed(time.Now().Unix())
	flag.Parse()

	if *Mode == "debug" {
		wlog.DevelopMode()
	}

	searchKeyword := queue.InitComsumer(define.TaskSearchKeyword, &TaskSearchKeyword{})
	parseDouyinURL := queue.InitComsumer(define.TaskParseDouyinURL, &TaskParseDouyinURL{})

	Q = queue.InitProducer()

	wlog.Info("等待新任务中...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	searchKeyword.Stop()
	parseDouyinURL.Stop()
}
