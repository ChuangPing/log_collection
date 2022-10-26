package model

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// 处理配置文件

type Config struct {
	ConKafka `ini:"kafka"`
	ConfEs   `ini:"es"`
}

//	ConKafka kafka配置选项
type ConKafka struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

//	ConfEs es配置选项
type ConfEs struct {
	Address  string `ini:"address"`
	Index    string `ini:"index"`
	ChanSize int `ini:"chanSize"`
}

func InitModel() (*Config, error) {
	var config = new(Config)
	err := ini.MapTo(config, "F:/Microservice/logagentProject/logagent_v2/log_transfer/conf/config.ini")
	if err != nil {
		logrus.Errorf("ini mapTo failed,err:%v\n", err)
		return config, err
	}
	return config, nil
}
