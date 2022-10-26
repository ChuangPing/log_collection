package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"strings"
)

// GetLocalIP 使用内置的net包获取IP
func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logrus.Error("net.InterfaceAddrs failed, err:", err)
		return
	}

	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return
}

// 获取本机IP
func GetOutbountIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	return strings.Split(localAddr.IP.String(), ":")[0]
}

func main() {
	ipStr, _ := GetLocalIP()
	fmt.Println("ip:", ipStr)
	ip := GetOutbountIP()
	fmt.Println("ip:", ip)
}
