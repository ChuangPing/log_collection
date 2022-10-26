package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	// kafka 操作对象
	client sarama.SyncProducer

	// MsgChan 接收tail读取数据的chan
	msgChan chan *sarama.ProducerMessage
)

// Init 初始化全局的kafka client  -- 连接kafka
func Init(address []string, chanSize int64) (err error) {
	// 1. 生产者者配置信息
	config := sarama.NewConfig()
	// 等待所有的确认
	config.Producer.RequiredAcks = sarama.WaitForAll // ACK
	//随机选取分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	// 2. 连接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("kafka: producer closed, err:", err)
		return err
	}
	// 初始化 MsgChan：chan类型只要有初始化才能使用，前面只是定义
	msgChan = make(chan *sarama.ProducerMessage, chanSize)
	// 起一个携程从msgChan中取数据并进行发送
	go sendMsg()
	return nil
}

// sendMsg:将MsgChan 中读取的msg发送到kafka
func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Warningf("send msg failed,err:%v\n", err)
				return
			}
			logrus.Infof("send msg to kafka success, pid:%v, offset:%v", pid, offset)
		}
	}
}

// SendMsgChan 向外提供获取可以将消息发送到msgChan管道的接口(函数）
func SendMsgChan(message *sarama.ProducerMessage) {
	msgChan <- message
}
