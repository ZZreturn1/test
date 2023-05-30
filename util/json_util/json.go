package json_util

import (
	"errors"
)

// RawMessage 是一个类型别名，表示 JSON 原始消息的字节切片
type RawMessage []byte

// MarshalJSON 自定义 json.RawMessage 的默认行为
func (m RawMessage) MarshalJSON() ([]byte, error) {
	if len(m) == 0 {
                                // 如果 RawMessage 为空，则返回 "null"
		return []byte("null"), nil 
	}
                // 返回 RawMessage 本身
	return m, nil 
}

// UnmarshalJSON 将 *m 设置为 data 的副本
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
                                // 如果 m 为空指针，则返回错误：在空指针上调用 UnmarshalJSON
		return errors.New("json.RawMessage: 在空指针上调用 UnmarshalJSON") 
	}
                // 将 data 的副本追加到 *m 中
	*m = append((*m)[0:0], data...) 

                // 返回 nil 表示无错误
	return nil 
}