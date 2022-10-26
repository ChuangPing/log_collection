package main

import (
	"fmt"
	"github.com/Shopify/sarama"
)

// 利用go第三方包实现kafka发送数据  -- go get github.com/Shopify/sarama

func main() {
	// 1、生产者配置信息
	config := sarama.NewConfig()
	// 发送完数据需要leader 和 follow 都确认
	config.Producer.RequiredAcks = sarama.WaitForAll
	// 新选出一个partition
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	// 成功交付的消息将在success channel中返回
	config.Producer.Return.Successes = true

	// 2、连接kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("producer closed, err", err)
		return
	}
	defer client.Close()

	//3、封装消息
	msg := &sarama.ProducerMessage{
		// 主题
		Topic: "shopping",
		Value: sarama.StringEncoder("2022-10-15, "),
	}

	//4、发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Printf("pid:%v，offest:%v\n", pid, offset)
}

// 启动zookeeper ：bin\windows\zookeeper-server-start.bat config\zookeeper.properties
// 启动kafka: bin\windows\kafka-server-start.bat config\server.properties
// 查看： bin\window\kafka-console-consumer.bat --bootstrap-server 127.0.0.1:9092 --topic web_log --from-beginning
