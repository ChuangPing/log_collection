package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

// 测是利用tail读取日志
func main() {
	// 读取日志文件路径
	fileName := "demo/tailfdemo/my.log"

	// 配置tail
	config := tail.Config{
		// 重复读  -日志文件切割可以跟随文件名继续读
		ReOpen:   true,
		Follow:   true,
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		// 允许读取文件不存在（不会报错）
		MustExist: false,
		Poll:      true,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tails err:", err)
		return
	}
	// 循环读取数据
	var (
		msg *tail.Line
		ok  bool
	)
	for {
		msg, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen,filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("msg:", msg.Text)
	}
}
