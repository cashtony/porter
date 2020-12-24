package queue

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"porter/define"
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
	consumer.SetLoggerLevel(nsq.LogLevelWarning)
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

type APITopicStat struct {
	Topics []struct {
		Channels []struct {
			ChannelName string `json:"channel_name"`
			Depth       int    `json:"depth"`
		} `json:"channels"`
		TopicName string `json:"topic_name"`
	}
}

func GetTopicStat(topic string) (*APITopicStat, error) {
	u := fmt.Sprintf("%s&topic=%s", define.TopicStats, topic)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &APITopicStat{}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}

	return result, nil
}
