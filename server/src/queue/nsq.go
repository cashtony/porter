package queue

import (
	"flag"
	"porter/wlog"

	"github.com/nsqio/go-nsq"
)

var NSQLookupd = flag.String("NSQLookupd", "localhost:4161", "NSQ lookup 服务")
var NSQD = flag.String("NSQD", "localhost:4150", "NSQ 监控")

func InitComsumer(topic string, f nsq.Handler) *nsq.Consumer {
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, "channel", config)
	if err != nil {
		wlog.Fatal("Comsumer消息队列初始化失败", err)
	}
	consumer.AddHandler(f)

	// Use nsqlookupd to discover nsqd instances.
	// See also ConnectToNSQD, ConnectToNSQDs, ConnectToNSQLookupds.
	err = consumer.ConnectToNSQLookupd(*NSQLookupd)
	if err != nil {
		wlog.Fatal("Comsumer消息队列连接lookup服务失败", err)
	}

	return consumer
}

func InitProducer() *nsq.Producer {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(*NSQD, config)
	if err != nil {
		wlog.Fatal("Producer消息队列初始化失败", err)
	}

	return producer
}
