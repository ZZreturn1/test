package global

import (
	"context"
	"github.com/robfig/cron/v3"
	_ "unsafe"
)

var webServer WebServer	// 全局Web服务器实例

// WebServer是Web服务器接口
type WebServer interface {
	GetCron() *cron.Cron	// 获取定时任务调度器
	GetCtx() context.Context	// 获取上下文
}

// SetWebServer设置全局Web服务器实例
func SetWebServer(s WebServer) {
	webServer = s	// 设置全局Web服务器实例
}

// GetWebServer获取全局Web服务器实例
func GetWebServer() WebServer {
	return webServer	// 获取全局Web服务器实例
}