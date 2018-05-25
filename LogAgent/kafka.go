package main

import (
	"strings"

	"github.com/astaxie/beego/logs"

	"github.com/Shopify/sarama"
)

// KafkaSender 用于管理kafka连接
type KafkaSender struct {
	kafkaClient sarama.SyncProducer
	msgChan     chan string
}

// NewKafkaSender 生成kafkaMgr的实例
func NewKafkaSender() *KafkaSender {
	return &KafkaSender{
		msgChan: make(chan string, 10000),
	}
}

func initKafka(kafkaAddr string) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		return
	}

	kafkaSend.kafkaClient = client
	kafkaSend.process()
	return
}

func (k *KafkaSender) process() {
	for i := 0; i < appConf.kafkaThread; i++ {
		wg.Add(1)
		go kafkaSend.sendMsgToKafka()
	}
}

// AddMsg add message to msgChan
func (k *KafkaSender) AddMsg(msg string) {
	msg = strings.TrimSpace(msg)
	if len(msg) == 0 {
		return
	}
	k.msgChan <- msg
}

func (k *KafkaSender) sendMsgToKafka() {
	defer wg.Done()
	for line := range k.msgChan {
		msg := &sarama.ProducerMessage{}
		msg.Topic = "nginx_log"
		msg.Value = sarama.StringEncoder(line)
		_, _, err := k.kafkaClient.SendMessage(msg)
		if err != nil {
			logs.Warn("send message to kafka failed. msg:[%s]  err:%v", line, err)
			continue
		}
		logs.Debug("send message to kafka success. msg:", line)
	}
}
