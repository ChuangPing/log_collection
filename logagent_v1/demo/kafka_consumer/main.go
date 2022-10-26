package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"sync"
)

// kafka consumer demo:获取kafka中存储的数据进行消费

// 利用sarama:sarama包不仅可以发送数据到kafka，还可以从kafka读取数据

func main() {
	var wg sync.WaitGroup
	// 创建新的消费者
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		logrus.Errorf("start consumer failed,err:%v", err)
		return
	}
	// 拿到指定的topic下面的分区列表
	//topics, err := consumer.Topics()
	//fmt.Println("topics:", topics)
	partitionList, err := consumer.Partitions("app_log ") // 根据topic获取到所有的分区
	if err != nil {
		logrus.Errorf("get partitions list failed, err:%v", err)
		return
	}
	fmt.Println("list partition:", partitionList)
	// 遍历所有分区
	for partition := range partitionList {
		// 针对每一个分区建立一个对应的分区消费者
		pc, err := consumer.ConsumePartition("app_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			logrus.Errorf("failed to start consumer for partition %d, err %v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费者中获取信息
		wg.Add(1)
		go func(partitionConsumer sarama.PartitionConsumer) {
			wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("parition:%d offset:%d Key:%v value:%v", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}
	wg.Wait()
}

// 命令行消费者：bin\windows\kafka-console-consumer.bat --bootstrap-server 127.0.0.1:9092 --topic app_log --from-beginning
