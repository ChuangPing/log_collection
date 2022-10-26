package etcd

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"logagent_v2/comm"
	"logagent_v2/tailfile"
	"time"
)

var (
	client *clientv3.Client
)

// Init 连接etcd客户端初始化client对象
func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		logrus.Error("connect etcd failed, err:", err)
		return
	}
	return
}

// PutConf 将文件路径和topic等信息保存在etcd
func PutConf(key string, jsonStr string) (err error) {
	//valByte, err := json.Marshal(collectEntryList)
	//if err != nil {
	//	logrus.Errorf("json marshal collectEntryList:%v failed, err", collectEntryList, err)
	//	return
	//}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err = client.Put(ctx, key, jsonStr)
	if err != nil {
		logrus.Errorf("put key:%s to etcd failed,err:%v\n", key, err)
		return
	}
	return
}

// GetConf 根据配置函数中的key，获取存储在etcd中value（收集日志的位置和topic）
func GetConf(key string) (collectEntryList []comm.CollectEntry, err error) {
	// 准备具有操作时自动停止的context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	// 获取etcd中的值  key -> value
	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get key:%s from etcd failed, err:%v\n", key, err)
		return
	}
	// 判断是否读取到数据  -- 防止空值
	if len(resp.Kvs) == 0 {
		logrus.Warningf("get len:0 conf from etcd by key:%s", key)
		return
	}
	ret := resp.Kvs[0]
	// ret.Value  中存储的是json格式的字符串
	// 将对应的json数据反序列化
	err = json.Unmarshal(ret.Value, &collectEntryList)
	if err != nil {
		logrus.Error("json unmarshal failed,err", err)
		return
	}
	return
}

// WatchConf watch监控etcd日志收集项配置的变化函数
func WatchConf(key string) {
	// 不停的监控
	for {
		watchChan := client.Watch(context.Background(), key)
		for wresp := range watchChan {
			// 提升监控到配置项信息出现变化
			logrus.Infof("get new conf from etcd")
			for _, evt := range wresp.Events {
				// 接收变化后的配置项
				var newConf []comm.CollectEntry
				if evt.Type == clientv3.EventTypeDelete {
					// 如果是删除操作
					tailfile.SendNewConf(newConf)
					continue
				}
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logrus.Errorf("json unmarshal failed, conf:%v err:%v\n", evt.Kv.Value, err)
					// 只是这一次监控出现问题，不因该退出任然还要继续监控
					continue
				}
				//新配置出现告诉tailfile这个模块根据最新的配置进行读取文件（将配置信息传递到tailfile模块中的newConf管道
				tailfile.SendNewConf(newConf) // 没有监听到新配置，chan就是阻塞状态
			}
		}
	}
}
