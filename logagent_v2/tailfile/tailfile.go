package tailfile

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"logagent_v2/kafka"
	"strings"
	"time"
)

type TailTask struct {
	Path    string
	Topic   string
	TailObj *tail.Tail
	// 方便进行任务管理，在合适的时候杀死任务
	ctx    context.Context
	cancel context.CancelFunc
}

// 初始化日志收集任务
func newTailTask(path, topic string) *TailTask {
	// 带有cancel功能的上下文
	ctx, cancel := context.WithCancel(context.Background())
	return &TailTask{
		Path:   path,
		Topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}
}

// InitTail 日志任务tail对象初始化
func (tt *TailTask) InitTail() (err error) {
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	tt.TailObj, err = tail.TailFile(tt.Path, config)
	return
}

// run 日志收集任务读取日志的方法
func (tt *TailTask) run() {
	logrus.Infof("run log task,file:%s,topic:%s\n", tt.Path, tt.Topic)
	// 读取日志，发往kafka存储日志
	for {
		select {
		case <-tt.ctx.Done(): //接收到停止信号
			logrus.Infof("the task get Done sign is Doned, task path:%s", tt.Path)
			return
		case line, ok := <-tt.TailObj.Lines:
			// 循环读数据
			// tail.Lines 读取结果存放在管道
			if !ok {
				logrus.Warn("tail file failed，file:", tt.TailObj.Filename)
				time.Sleep(time.Second)
				continue
			}
			// 处理空行情况
			if len(strings.Trim(line.Text, "\r")) == 0 {
				logrus.Infof("日志文件:%s 有空行\n", tt.TailObj.Filename)
				continue
			}
			// 将读取出来的数据发送给kafka中存数据的管道，kafka再从管道中读取数据  --将代码改为异步，提高效率
			// 封装成kafka消息类型
			msg := &sarama.ProducerMessage{}
			msg.Topic = tt.Topic
			msg.Value = sarama.StringEncoder(line.Text)
			// 发送至通道
			kafka.SendMsgChan(msg)
		}
	}
}

//// Init 初始化日志收集任务
//func Init(allconf []comm.CollectEntry) (err error) {
//	// allconf 里面存了若干个日志收集任务
//	// 针对每一个日志收集项创建一个对应的tailObj
//	for _, conf := range allconf {
//		// 初始化一个日志收集任务
//		tt := newTailTask(conf.Path, conf.Topic)
//		// 创建一个日志收集任务
//		err = tt.InitTail()
//		if err != nil {
//			logrus.Errorf("create tailObj for path:%s,topic:%s is failed.err:%v", tt.Path, tt.Topic, err)
//			// 只是一个tail任务失败，因该继续下一次循环而不是终止任务
//			continue
//		}
//		// 初始化一个日志收集任务就去读数据
//		go tt.run()
//	}
//	// 初始化全局newConfChan -- 全局只是声明并未初始化
//	newConfChan = make(chan []comm.CollectEntry) // 不带缓冲区阻塞式的管道
//	newConf := <-newConfChan
//	logrus.Infof("get new conf:%v\n", newConf)
//	return
//}

// SendNewConf 将新配置信息发送到chan
//func SendNewConf(newConf []comm.CollectEntry) {
//	//newConfChan = make(chan []comm.CollectEntry)
//	newConfChan <- newConf
//	fmt.Println("newConfChan", <-newConfChan)
//}
