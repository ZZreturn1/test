package common

import (
	"strings"
)

type multiError []error

func (e multiError) Error() string {
	var r strings.Builder
                // 拼接错误字符串的前缀
	r.WriteString("multierr: ") 
	for _, err := range e {
                                // 将每个错误的字符串表示追加到结果字符串中
		r.WriteString(err.Error()) 
                                // 在错误之间添加分隔符
		r.WriteString(" | ") 
	}
                // 返回拼接后的错误字符串
	return r.String() 
}

func Combine(maybeError ...error) error {
                // 定义一个多错误类型的切片
	var errs multiError 
	for _, err := range maybeError {
		if err != nil {
                                                // 将非空错误追加到切片中
			errs = append(errs, err) 
		}
	}
	if len(errs) == 0 {
                                // 如果切片为空，则返回nil表示没有错误
		return nil 
	}
                // 返回多错误切片作为组合后的错误
	return errs 
}
