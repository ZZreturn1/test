package service

import (
	"os"
	"syscall"
	"time"
	"x-ui/logger"
)

// PanelService提供面板服务功能
type PanelService struct {
}

// RestartPanel重新启动面板
func (s *PanelService) RestartPanel(delay time.Duration) error {
	p, err := os.FindProcess(syscall.Getpid())	// 查找当前进程
	// 如果查找进程出现错误，则返回错误
	if err != nil {
		return err
	}
	// 启动一个goroutine，在延迟时间后发送SIGHUP信号给进程
	go func() {
		time.Sleep(delay)		// 延迟一段时间
		err := p.Signal(syscall.SIGHUP)	// 发送SIGHUP信号给进程
		if err != nil {
			logger.Error("send signal SIGHUP failed:", err)	// 如果发送信号失败，则记录错误日志
		}
	}()
	return nil
}