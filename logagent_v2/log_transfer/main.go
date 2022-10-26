package main

import (
	"github.com/sirupsen/logrus"
	"logagent_v2/log_transfer/es"
	"logagent_v2/log_transfer/kafka"
	"logagent_v2/log_transfer/model"
)

// 从kafka中读取日志发送es中

func main() {
	// 初始化配置文件
	config, err := model.InitModel()
	if err != nil {
		logrus.Error("init model failed, err", err)
		return
	}
	logrus.Info("init model success!")

	// 初始化kafka  -- 从kafka中读取数据
	err = kafka.InitKafka([]string{config.ConKafka.Address}, config.Topic)
	if err != nil {
		logrus.Error("init kafka failed, err", err)
		return
	}

	// 初始化es
	err = es.InitElastic(config.ConfEs.ChanSize, config.ConfEs.Index, config.ConfEs.Address)
	if err != nil {
		logrus.Error("init es failed,err:", err)
		return
	}
	select {}
}
