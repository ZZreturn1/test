package job

import (
	"x-ui/logger"
	"x-ui/web/service"
)

// XrayTrafficJob是Xray流量统计任务
type XrayTrafficJob struct {
	xrayService    service.XrayService	// Xray服务
	inboundService service.InboundService	// 入站服务
}

// NewXrayTrafficJob创建一个新的Xray流量统计任务
func NewXrayTrafficJob() *XrayTrafficJob {
	return new(XrayTrafficJob)	// 创建并返回一个XrayTrafficJob实例
}

// Run运行Xray流量统计任务
func (j *XrayTrafficJob) Run() {
	if !j.xrayService.IsXrayRunning() {
		return
	}
	traffics, err := j.xrayService.GetXrayTraffic()
	if err != nil {
		logger.Warning("get xray traffic failed:", err)	// 获取Xray流量统计失败
		return
	}
	err = j.inboundService.AddTraffic(traffics)
	if err != nil {
		logger.Warning("add traffic failed:", err)	// 添加流量统计失败
	}
}