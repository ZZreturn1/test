package common

// 导入了 errors 包用于创建错误，fmt 包用于格式化输出，x-ui/logger 包用于日志记录。
import (
	"errors"
	"fmt"
	"x-ui/logger"
)

// CtxDone 表示上下文已完成的错误，定义了一个错误变量 CtxDone，表示上下文已完成。
var CtxDone = errors.New("context done")

// NewErrorf 根据格式化字符串创建新的错误
func NewErrorf(format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	return errors.New(msg)
}

// NewErrorf 函数根据格式化字符串和参数创建新的错误。
func NewError(a ...interface{}) error {
	msg := fmt.Sprintln(a...)
	return errors.New(msg)
}

// Recover 恢复 panic 并进行日志记录
// Recover 函数用于捕获 panic，并进行日志记录。如果发生 panic，将打印日志信息并返回 panic 的值。
/* 
    在Go语言中，panic是一种错误处理机制，用于表示程序发生了无法恢复的错误。
    当程序遇到无法处理的异常情况时，会触发panic，并立即停止当前函数的执行，然后逐层向上回溯，
    执行每层调用函数的defer语句，最终导致程序终止并输出panic信息。
*/
func Recover(msg string) interface{} {
	panicErr := recover()
	if panicErr != nil {
		if msg != "" {
			logger.Error(msg, "panic:", panicErr)
		}
	}
	return panicErr
}