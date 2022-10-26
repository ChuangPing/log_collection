package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
)

// cpu info
func getCpuInfo() {
	cpuInfos, err := cpu.Info()
	if err != nil {
		fmt.Println("get cpu info failed, err:", err)
		return
	}

	for _, cpuInfo := range cpuInfos {
		fmt.Printf("cpuInfo:%v\n", cpuInfo)
	}

	// 获取cpu的使用频率  (每秒获取一次）
	for {
		precent, _ := cpu.Percent(time.Second, false)
		fmt.Printf("cpu percent:%v\n", precent) // cpu percent:[8.3984375]

	}
}

// 获取cpu负载 (windows 还没实现）
func getCpuLoad() {
	info, _ := load.Avg()
	fmt.Printf("avg:%v\n", info)
}

// 获取内存信息
func getMemInfo() {
	info, err := mem.VirtualMemory()
	if err != nil {
		fmt.Printf("get mem info failed, err:%v\n", err)
		return
	}
	fmt.Println("mem:info", info)
}

func main() {
	//getCpuInfo()
	getCpuLoad()
	getMemInfo()
}
