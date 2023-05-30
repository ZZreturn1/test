package common

import (
	"fmt"
)

func FormatTraffic(trafficBytes int64) (size string) {
	if trafficBytes < 1024 {
                                // 小于1KB，以字节单位（B）格式化
		return fmt.Sprintf("%.2fB", float64(trafficBytes)/float64(1)) 
	} else if trafficBytes < (1024 * 1024) {
                                // 小于1MB，以千字节单位（KB）格式化
		return fmt.Sprintf("%.2fKB", float64(trafficBytes)/float64(1024)) 
	} else if trafficBytes < (1024 * 1024 * 1024) {
                                // 小于1GB，以兆字节单位（MB）格式化
		return fmt.Sprintf("%.2fMB", float64(trafficBytes)/float64(1024*1024)) 
	} else if trafficBytes < (1024 * 1024 * 1024 * 1024) {
                                // 小于1TB，以吉字节单位（GB）格式化
		return fmt.Sprintf("%.2fGB", float64(trafficBytes)/float64(1024*1024*1024)) 
	} else if trafficBytes < (1024 * 1024 * 1024 * 1024 * 1024) {
                                // 小于1PB，以太字节单位（TB）格式化
		return fmt.Sprintf("%.2fTB", float64(trafficBytes)/float64(1024*1024*1024*1024)) 
	} else {
                                // 大于等于1PB，以艾字节单位（EB）格式化
		return fmt.Sprintf("%.2fEB", float64(trafficBytes)/float64(1024*1024*1024*1024*1024)) 
	}
}