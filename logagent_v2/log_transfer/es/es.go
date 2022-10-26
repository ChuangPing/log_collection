package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

//	ESClient将日志写入ElasticSearch
type ESClient struct {
	client      *elastic.Client
	index       string
	logDataChan chan interface{}
}
type msgType struct {
	key   string
	value string
}

var (
	esClient *ESClient
)

func InitElastic(chanSize int, index string, address string) (err error) {
	client, err := elastic.NewClient(elastic.SetURL("http://" + address))
	if err != nil {
		logrus.Error("connect elasticSearch failed,err:", err)
		return
	}
	logrus.Infof("connect elastic success!\n")
	// 初始化client和存放日志的chan
	esClient = &ESClient{
		client:      client,
		index:       index,
		logDataChan: make(chan interface{}, chanSize),
	}

	// 从logChan 中取数据放到es中
	go SendToEs()
	return
}

// SendLogData 向外暴露向logData chan中写入数据的方法
func SendLogData(msg interface{}) {
	esClient.logDataChan <- msg
	logrus.Infof("get logData from kafka:%v\n", msg)
}

// SendToEs	将logDataChan 中的数据写入到es中
func SendToEs() {
	select {
	case msg := <-esClient.logDataChan:
		msgString := msgType{
			key:   "web_log",
			value: msg.(string),
		}
		//msgByte, err := json.Marshal(&msgString)
		//if err != nil {
		//	fmt.Println("json marsh failed,err:", err)
		//}
		put1, err := esClient.client.Index().
			Index(esClient.index).
			BodyJson(msgString).
			Do(context.Background())
		if err != nil {
			logrus.Error("send to es failed,err:", err)
		}
		//fmt.Println("pu1", *put1)
		fmt.Printf("indexed user %s to index %s,tpye %s\n", put1.Id, put1.Index, put1.Type)
	}

}
