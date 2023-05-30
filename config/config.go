package config

import (
	_ "embed" // 导入空白标识符以使用 embed 包
	"fmt"     // 导入 fmt 包提供格式化和打印功能
	"os"      // 导入 os 包提供与操作系统交互的函数
	"strings" // 导入 strings 包提供字符串相关的函数
)

//go:embed version // 使用 embed 标记嵌入 version 文件内容
var version string

//go:embed name // 使用 embed 标记嵌入 name 文件内容
var name string

// LogLevel 日志级别类型声明
type LogLevel string

const (
	Debug LogLevel = "debug" // 调试日志级别
	Info  LogLevel = "info"  // 信息日志级别
	Warn  LogLevel = "warn"  // 警告日志级别
	Error LogLevel = "error" // 错误日志级别
)

// GetVersion 返回版本号
func GetVersion() string {
	return strings.TrimSpace(version) // 返回修剪后的版本号字符串
}

// GetName 返回名称
func GetName() string {
	return strings.TrimSpace(name) // 返回修剪后的名称字符串
}

// GetLogLevel 返回日志级别
func GetLogLevel() LogLevel {
	if IsDebug() { // 如果启用调试模式
		return Debug // 返回调试日志级别
	}
	logLevel := os.Getenv("XUI_LOG_LEVEL") // 获取环境变量 XUI_LOG_LEVEL 的值
	if logLevel == "" { // 如果环境变量未设置
		return Info // 返回信息日志级别
	}
	return LogLevel(logLevel) // 返回根据环境变量设置的日志级别
}

// IsDebug 检查是否启用调试模式
func IsDebug() bool {
	return os.Getenv("XUI_DEBUG") == "true" // 如果环境变量 XUI_DEBUG 的值为 "true"，则返回 true，否则返回 false
}

// GetDBPath 返回数据库路径
func GetDBPath() string {
	return fmt.Sprintf("/etc/%s/%s.db", GetName(), GetName()) // 返回格式化后的数据库路径字符串
}