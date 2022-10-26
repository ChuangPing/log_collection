package main

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"logagent/kafka"
	"logagent/tailfile"
	"strings"
	"time"
)

// Config 使用struct对配置文件进行封装  -- 将配置信息映射到结构体
type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	ChanSize int64  `ini:"chan_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

// run：将tail读取的日志发送到kafka
func run() error {
	for {
		// 循环读数据
		line, ok := <-tailfile.TailObj.Lines // tailObj.lines 将读取到的一行数据放在chan中
		if !ok {
			logrus.Warn("tail file close reopen filenam:%s\n", tailfile.TailObj.Filename)
			time.Sleep(time.Second)
			continue
		}
		//fmt.Printf("line:%#v", line.Text)  // line:"\r"  只是换行会右 \r -- window平台有 \n\r作为换行
		//	判断是否是空行，空行就不写如
		if len(strings.Trim(line.Text, "\r")) == 0 {
			continue
		}

		// 利用管道将同步的代码改为异步：这边tail读取到数据后就仍进管道
		// 把读出来的一行日志包装成kafka里面的msg类型
		msg := &sarama.ProducerMessage{}
		msg.Topic = "web_log"
		msg.Value = sarama.StringEncoder(line.Text)

		// 将消息丢到通道
		kafka.SendMsgChan(msg)
	}
}

func main() {
	// 0.读取配置文件 ‘go-ini'

	// 测试使用直接读取配置文件
	//cfg, err := ini.Load("./conf/config.ini")
	//if err != nil {
	//	logrus.Errorf("load config failed, err:%v", err)
	//}

	//kafkaAddr := cfg.Section("kafka").Key("address").String()
	//fmt.Println("kafka address", kafkaAddr)// kafka address 127.0.0.1:9092

	// 将配置文件中的信息映射到config结构体
	var configObj = new(Config) // 由于需要将config传入函数参数,因此使用new的方式初始化Config结构体指针,这样才函数体内部赋值才会影响外面
	err := ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Error("load config failed, err:%v", err)
		return
	}

	//1. 初始化连接kafka
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Error("init kafka failed, err:", err)
		return
	}
	logrus.Info("init kafka success!")

	// 2.根据配置中的日志路径初始化tail  -- 读取文件
	err = tailfile.Init(configObj.CollectConfig.LogFilePath)
	if err != nil {
		logrus.Error("init tail failed, err:", err)
		return
	}
	logrus.Info("init tailfile success!")

	// 3.把利用tail读取的文件发送到kafka
	err = run()
	if err != nil {
		logrus.Error("run failed, err:%v", err)
	}
}

/*
	当前版本日志收集项目存在的问题
1、只能读取一个日志文件，不支持多个文件
2、无法管理日志的topic
3、无法根据实际需要动态扩展：
想要达到的效果：当加入新的配置可以自动完成不同种类的日志收集

*/
