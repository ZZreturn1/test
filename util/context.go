package util

import "context"

func IsDone(ctx context.Context) bool {
	// 检查上下文是否已完成
	select {

	// 如果上下文已完成
	case <-ctx.Done(): 
		return true
	// 如果上下文未完成
	default: 
		return false
	}
}
