package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		logrus.Errorf("connect etcd failed,err:%v\n", err)
		return
	}
	logrus.Infoln("connect etcd success!")
	defer cli.Close()

	// 使用watch 监视key的变化  返回值是一个chan类型
	wactChan := cli.Watch(context.Background(), "name")
	// 遍历管道中的类容获name 的变化情况  -- name发生改变就会给chan里面发送数据
	for wresp := range wactChan {
		for _, ev := range wresp.Events {
			// 对key的操作类型（增删改）已经对应的值
			fmt.Printf("Type : %s key:%s value:%s", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
