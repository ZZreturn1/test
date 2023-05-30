// +build darwin

package sys

import (
	"github.com/shirou/gopsutil/net"
)

// GetTCPCount 获取当前系统上的 TCP 连接数。
// 它使用 gopsutil/net 包来获取 TCP 连接统计信息。
func GetTCPCount() (int, error) {
	stats, err := net.Connections("tcp")
	if err != nil {
		return 0, err
	}
	return len(stats), nil
}

// GetUDPCount 获取当前系统上的 UDP 连接数。
// 它使用 gopsutil/net 包来获取 UDP 连接统计信息。
func GetUDPCount() (int, error) {
	stats, err := net.Connections("udp")
	if err != nil {
		return 0, err
	}
	return len(stats), nil
}
