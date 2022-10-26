package comm

import (
	"github.com/sirupsen/logrus"
	"net"
	"strings"
)

// CollectEntry :etcd中读取数据后的结构体
type CollectEntry struct {
	Path  string `json:"path"`  // 读取日志文件路径
	Topic string `json:"topic"` // 日志文件发往kafka中的那个topic
}

// GetOutbountIP 获取本机IP
func GetOutbountIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logrus.Error("etOutbountIP:net.Dial failed,err", err)
		return "get IP failed", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return strings.Split(localAddr.IP.String(), ":")[0], nil
}
