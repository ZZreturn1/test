package logger

import (
	"github.com/op/go-logging"
	"os"
)

var logger *logging.Logger

func init() {
	InitLogger(logging.INFO)
}

// InitLogger 初始化日志记录器
func InitLogger(level logging.Level) {
                // 定义日志输出格式
	format := logging.MustStringFormatter(
		`%{time:2006/01/02 15:04:05} %{level} - %{message}`,
	)

                // 创建新的日志记录器实例
	newLogger := logging.MustGetLogger("x-ui")

                // 创建日志输出的后端
	backend := logging.NewLogBackend(os.Stderr, "", 0)

                // 创建格式化后的后端
	backendFormatter := logging.NewBackendFormatter(backend, format)

                // 创建可设置级别的后端
	backendLeveled := logging.AddModuleLevel(backendFormatter)

                // 设置日志输出级别
	backendLeveled.SetLevel(level, "")

                // 将后端应用于日志记录器
	newLogger.SetBackend(backendLeveled)

                // 实例化，并保存到logger
	logger = newLogger
}


// Debug 输出调试级别的日志
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugf 格式化输出调试级别的日志
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Info 输出信息级别的日志
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infof 格式化输出信息级别的日志
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warning 输出警告级别的日志
func Warning(args ...interface{}) {
	logger.Warning(args...)
}

// Warningf 格式化输出警告级别的日志
func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

// Error 输出错误级别的日志
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorf 格式化输出错误级别的日志
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
