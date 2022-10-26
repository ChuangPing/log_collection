package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	client  sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)

// Init 始化全局的kafka client  -- 连接kafka
func Init(address []string, chanSize int64) (err error) {
	// 1.生产者配置信息
	config := sarama.NewConfig()
	// 等待所有的确认
	config.Producer.RequiredAcks = sarama.WaitForAll // ACK
	//随机选取分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	// 2. 连接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Error("connect kafka failed,err:", err)
		return
	}
	// 初始化 msgChan
	msgChan = make(chan *sarama.ProducerMessage, chanSize)
	// 使用携程等待msgChain中的数据并发送到kafka
	go sendMsg()
	return err
}

// sendMsg:将msgChan 中msg发送到kafka
func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Error("send msg to kafka failed, err:", err)
				return
			}
			logrus.Infof("send msg to kafka success, pid:%v, offset:%v", pid, offset)
		}
	}
}

// SendMsgChan 向外提供向外提供获取可以将消息发送到msgChan管道的接口(函数）
func SendMsgChan(message *sarama.ProducerMessage) {
	msgChan <- message
}
