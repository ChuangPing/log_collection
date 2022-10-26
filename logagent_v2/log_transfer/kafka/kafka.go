package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"logagent_v2/log_transfer/es"
)

func InitKafka(address []string, topic string) error {
	//var wg sync.WaitGroup
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	client, _ := sarama.NewClient(address, config)
	//consumer, err := sarama.NewConsumer(address, nil)
	consumer, _ := sarama.NewConsumerFromClient(client)
	// 获取topic下面的分区列表
	fmt.Println(consumer.Topics())
	fmt.Println(consumer.Partitions("qq_log"))
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		logrus.Error("get consumer partitions failed, err:", err)
		return err
	}
	fmt.Println("partitionList:", partitionList)
	// 遍历分区
	for partition, _ := range partitionList {
		// 针对每个分区建立一个消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			logrus.Errorf("failed start consumer for partition:%d err:%v\n", partition, err)
			return err
		}
		// 异步从每个分区读取消费者信息
		go ReadMsg(pc)
		//go func() {}()
		//msg := pc.Messages()
		//fmt.Println("1111111111",)
		//for msg := range pc.Messages() {
		//	fmt.Printf("mag%s\n", msg.Value)
		//}

		//wg.Add(1)
		//go func(partitionConsumer sarama.PartitionConsumer) {
		//	for msg := range partitionConsumer.Messages() {
		//		fmt.Printf("read msg from kafka key:%v, value=%v", msg.Key, msg.Value)
		//	}
		//	//wg.Done()
		//	//defer pc.AsyncClose()
		//}(pc)
	}
	//wg.Wait()
	return nil
}

func ReadMsg(pc sarama.PartitionConsumer) {
	for msg := range pc.Messages() {
		fmt.Printf("read msg from kafka key:%s, value=%s\n", msg.Key, msg.Value)
		// 将日志进行序列化
		//var m map[interface{}]interface{}
		//err := json.Unmarshal(msg, &m)
		//if err != nil {
		//	logrus.Error("json unmarshal msg failed,err:", err)
		//	return
		//}
		// 将日志内容发送的es的chan中
		es.SendLogData(string(msg.Value))
	}
}
