package job

import "x-ui/web/service"

// CheckXrayRunningJob是检查Xray是否运行的任务
type CheckXrayRunningJob struct {
	xrayService service.XrayService	// Xray服务

	checkTime int	// 检查次数
}

// NewCheckXrayRunningJob创建一个新的检查Xray是否运行的任务
func NewCheckXrayRunningJob() *CheckXrayRunningJob {
	return new(CheckXrayRunningJob)	// 创建并返回一个CheckXrayRunningJob实例
}

// Run运行检查Xray是否运行的任务
func (j *CheckXrayRunningJob) Run() {
	// 如果Xray正在运行，重置检查次数并返回
	if j.xrayService.IsXrayRunning() {
		j.checkTime = 0	// 重置检查次数
		return
	}
	j.checkTime++	// 增加检查次数
	// 如果检查次数小于2，直接返回
	if j.checkTime < 2 {
		return
	}
	j.xrayService.SetToNeedRestart()	// 设置Xray服务需要重启
}