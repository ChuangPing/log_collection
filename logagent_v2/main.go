package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"logagent_v2/comm"
	"logagent_v2/etcd"
	"logagent_v2/kafka"
	"logagent_v2/tailfile"
)

// kafkaConf kafka配置
type KafkaConf struct {
	Address  string `ini:"address"`
	ChanSize int64  `ini:"chan_size"`
}

//// CollectConf 日志文件存储位置
//type CollectConf struct {
//	LogFilePath string `ini:"logfile_path"`
//}

// EtcdConf etcd配置信息
type EtcdConf struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

// Conf 系统配置信息
type Conf struct {
	KafkaConf `ini:"kafka"`
	//CollectConf `ini:"collect"`
	EtcdConf `ini:"etcd"`
}

func main() {
	// 获取本机IP，为后续去etcd获取配置信息做准备（etcd中的key更具带有ip）
	ip, err := comm.GetOutbountIP()
	fmt.Println("ip:", ip)
	if err != nil {
		logrus.Error("get ip failed,err:", err)
		return
	}

	// 0、读取系统配置  -- 注意这里为什么需要指针
	var conf = new(Conf)
	err = ini.MapTo(conf, "./conf/config.ini")
	if err != nil {
		logrus.Errorf("load config failed,err:%v\n", err)
		return
	}
	// 1.初始化 连接kafka并向kafka发送数据
	err = kafka.Init([]string{conf.KafkaConf.Address}, conf.KafkaConf.ChanSize)
	if err != nil {
		logrus.Error("init kafka failed,err:", err)
		return
	}
	logrus.Info("init kafka success!")

	// 2.将日志配置文件存储到etcd
	err = etcd.Init([]string{conf.EtcdConf.Address})
	if err != nil {
		logrus.Error("init etcd failed,err:", err)
		return
	}
	//var jsonStr string = `[
	//{
	//  "path":"F:/Microservice/logagentProject/logagent_v2/logs/my.log",
	//  "topic":"app_log"
	//},
	//{
	//  "path":"F:/Microservice/logagentProject/logagent_v2/logs/my1.log",
	//  "topic":"admin_log"
	//}
	//]`

	//err = etcd.PutConf(conf.EtcdConf.CollectKey, jsonStr)
	//if err != nil {
	//	logrus.Error("put collectkey failed,err:", err)
	//	return
	//}

	//从etcd中读取配置信息
	// 根据IP拼接 key
	collectKey := fmt.Sprintf(conf.EtcdConf.CollectKey, ip)
	allConf, err := etcd.GetConf(collectKey)
	if err != nil {
		logrus.Error("get collectkey failed,err:", err)
		return
	}
	logrus.Info("init etcd success!")
	// 监控etcd存储的key vaue数据的变化，监控conf.EtcdConf.CollectKey 对应的值的变化，从而动态响应日志任务无需重启业务（使用etcd的watch功能）
	go etcd.WatchConf(conf.EtcdConf.CollectKey)
	// 3.根据配置信息初始化tail，并将对应的日志发送至kafka
	err = tailfile.Init(allConf)
	if err != nil {
		logrus.Error("init tailfile failed,err:", err)
		return
	}
	logrus.Info("init tailfile success!")
	// 阻塞main函数携程，防止主协程直接退出
	run()
}

func run() {
	select {}
}
