package job

import (
	"x-ui/logger"
	"x-ui/web/service"
)

// CheckInboundJob是检查入站任务
type CheckInboundJob struct {
	xrayService    service.XrayService	// Xray服务
	inboundService service.InboundService	// 入站服务
}

// NewCheckInboundJob创建一个新的检查入站任务
func NewCheckInboundJob() *CheckInboundJob {
	return new(CheckInboundJob)	// 创建并返回一个CheckInboundJob实例
}

// Run运行检查入站任务
func (j *CheckInboundJob) Run() {
	count, err := j.inboundService.DisableInvalidInbounds()	// 禁用无效的入站配置
	if err != nil {
		logger.Warning("disable invalid inbounds err:", err)	// 记录禁用无效入站配置时的错误日志
	} else if count > 0 {
		logger.Debugf("disabled %v inbounds", count)	// 记录禁用的入站配置数量
		j.xrayService.SetToNeedRestart()	// 设置Xray服务需要重启
	}
}