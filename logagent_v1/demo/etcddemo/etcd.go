package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
	"time"
)

// 解压运行etcd.exe 后，测试：存储值：etcdctl.exe --endpoints=http://127.0.0.1:2379 put name "tome
// 获取值：etcdctl.exe --endpoints=http://127.0.0.1:2379 get name
//etcdctl.exe --endpoints=http://127.0.0.1:2379 get collect_log_222.198.39.41_conf

// go操作etcd  -- 使用etcd/clientv3:go get go.etcd.io/etcd/client/v3

// 测试put和get操作

func main() {
	client, err := clientv3.New(clientv3.Config{
		// etcd服务地址
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logrus.Errorf("connect to etcd failed,err:%v\n", err)
		return
	}
	logrus.Infof("connect to etcd success")
	defer client.Close()

	// put操作，向etcd发送数据
	// 传入一个具有超时机制的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	jsonStr := `[{"path":"F:/Microservice/logagentProject/logagent_v2/logs/my.log","topic":"app_log"},{"path":"F:/Microservice/logagentProject/logagent_v2/logs/my1.log","topic":"admin_log"},{"path":"F:/Microservice/logagentProject/logagent_v2/logs/my2.log","topic":"web_log"}]`
	//jsonStr := `[{"path":"F:/Microservice/logagentProject/logagent_v2/logs/my.log","topic":"web_log"},{"path":"F:/Microservice/logagentProject/logagent_v2/logs/my1.log","topic":"admin_log"}]`
	//jsonStr := `[{"path":"F:/Microservice/logagentProject/logagent_v2/logs/my.log","topic":"app_log"}]`
	//_, err = client.Put(ctx, "collect_log_conf", jsonStr)
	// 加入 ip
	_, err = client.Put(ctx, "collect_log_222.198.39.41_conf", jsonStr)
	//_, err = client.Put(ctx, "collect_log_222.198.39.41_")
	if err != nil {
		logrus.Errorf("put keyvalue to etcd failed,err:%v\n", err)
		return
	}
	cancel()

	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	getResponse, err := client.Get(ctx, "collect_log_conf")
	if err != nil {
		logrus.Errorf("get from etcd failed, err:%v\n", err)
		return
	}
	// 遍历获取：key value
	for _, ev := range getResponse.Kvs {
		fmt.Printf("key:%s,value:%s\n", ev.Key, ev.Value) // Type : PUT key:name value:tome

	}
	cancel()
}
