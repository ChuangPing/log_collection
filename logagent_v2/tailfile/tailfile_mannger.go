package tailfile

import (
	"github.com/sirupsen/logrus"
	"logagent_v2/comm"
)

// tailTask 日志收集任务管理
type tailTaskMg struct {
	tailTaskMap      map[string]*TailTask
	collectEntryList []comm.CollectEntry
	newConfChan      chan []comm.CollectEntry
}

var (
	ttMg *tailTaskMg
)

func Init(allConf []comm.CollectEntry) (err error) {
	// allConf 里面存了若干日志的收集
	// 正对每一个日志收集项创建一个一个对应的tailObj
	ttMg = &tailTaskMg{
		tailTaskMap:      make(map[string]*TailTask),
		collectEntryList: allConf,
		newConfChan:      make(chan []comm.CollectEntry, 1),
	}
	for _, conf := range allConf {
		tt := newTailTask(conf.Path, conf.Topic) // 创建一个日志收集任务
		err := tt.InitTail()                     //打开文件准备读
		if err != nil {
			if err != nil {
				logrus.Errorf("create tailObj for path:%s,topic:%s is failed.err:%v", tt.Path, tt.Topic, err)
				// 只是一个tail任务失败，因该继续下一次循环而不是终止任务
				continue
			}
		}
		logrus.Infof("create a tail task for pathL%v success", tt.Path)
		// 将创建成功的日志收集任务添加ttMg，方便后续新配置来可以进行管理
		ttMg.tailTaskMap[tt.Path] = tt
		// 开启后台的goroutine去收集日志
		go tt.run()
	}
	// 在后台等待新配置,并进行管理
	go ttMg.watch()
	return
}

func (tm *tailTaskMg) watch() {
	// 不停的监视，不断的从newConfChan中读取数据
	for {
		// 从newConfChan中取值，等待新的配置到来   -- etcd监孔到 key对应的值发生变化就会发送数据到管道
		newConf := <-tm.newConfChan
		// 新配置来了之后应该管理一下我之前启动的那些tailTask
		logrus.Infof("get new conf from etcd, conf:%v, start manage tailTask...", newConf)

		// 当有新的配置来，此时就要进行管理taiTaskMg中的任务
		for _, conf := range newConf {
			// 情况1. 原来已经存在的配置项（日志任务）不用处理
			if tm.isExsit(conf) {
				continue
			}
			// 情况2. 出现新的配置项，需要创建新的日志任务
			tt := newTailTask(conf.Path, conf.Topic) // 初始化日志任务
			err := tt.InitTail()                     // 打开文件准备读
			if err != nil {
				logrus.Errorf("create tailObj for path:%s,topic:%s is failed.err:%v", tt.Path, tt.Topic, err)
				// 只是一个tail任务失败，因该继续下一次循环而不是终止任务
				continue
			}
			logrus.Infof("create a tail for path:%s success", conf.Path)
			// 将新创建的日志任务加入tailMg方便管理
			tm.tailTaskMap[tt.Path] = tt
			// 新创建的日志任务开始读日志
			go tt.run()
		}
		// 情况3：原来存在的配置现在没有，则需要中断这些任务
		for key, task := range tm.tailTaskMap {
			var found bool
			for _, conf := range newConf {
				if key == conf.Path {
					found = true
					break
				}
			}
			if !found {
				// task 对应的任务在newconf中没有，需要关闭
				logrus.Infof("the task collect path:%s need to stop", task.Path)
				// 从tailMg中删除
				delete(tm.tailTaskMap, key)
				task.cancel() // 发送关闭信号，并不能立即停止，需要在运行日志任务的goroutine中接收这个信号并退出

			}
		}
	}

}

// SendNewConf 将新配置信息发送到chan
func SendNewConf(newConf []comm.CollectEntry) {
	ttMg.newConfChan <- newConf
}

func (t *tailTaskMg) isExsit(conf comm.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok
}
