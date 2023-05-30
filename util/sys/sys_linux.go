// +build linux

package sys

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func getLinesNum(filename string) (int, error) {
	// 打开指定文件
	file, err := os.Open(filename) 
	if err != nil {
		// 如果打开文件失败，则返回错误
		return 0, err 
	}

	// 延迟关闭文件，确保在函数结束时关闭文件
	defer file.Close() 

	// 行数计数器
	sum := 0 

	// 缓冲区大小为 8192 字节
	buf := make([]byte, 8192) 
	for {
		// 从文件中读取内容到缓冲区
		n, err := file.Read(buf) 

		// 缓冲区位置
		var buffPosition int 
		for {
			// 在缓冲区中查找换行符的位置
			i := bytes.IndexByte(buf[buffPosition:], '\n') 
			if i < 0 || n == buffPosition {
				// 如果未找到换行符或已读取到缓冲区末尾，则退出循环
				break 
			}
			// 更新缓冲区位置
			buffPosition += i + 1 
			// 增加行数计数
			sum++ 
		}

		if err == io.EOF {
			// 如果已到达文件末尾，则返回行数和 nil
			return sum, nil 
		} else if err != nil {
			// 如果出现其他错误，则返回行数和错误信息
			return sum, err 
		}
	}
}

func GetTCPCount() (int, error) {
	// 获取主机的 proc 目录路径
	root := HostProc() 

	// 获取 tcp4 连接数
	tcp4, err := getLinesNum(fmt.Sprintf("%v/net/tcp", root)) 
	if err != nil {
		return tcp4, err
	}

	// 获取 tcp6 连接数
	tcp6, err := getLinesNum(fmt.Sprintf("%v/net/tcp6", root)) 
	if err != nil {
		return tcp4 + tcp6, nil
	}

	// 返回 tcp4 和 tcp6 连接数之和
	return tcp4 + tcp6, nil 
}

func GetUDPCount() (int, error) {
	// 获取主机的 proc 目录路径
	root := HostProc() 

	// 获取 udp4 连接数
	udp4, err := getLinesNum(fmt.Sprintf("%v/net/udp", root)) 
	if err != nil {
		return udp4, err
	}

	// 获取 udp6 连接数
	udp6, err := getLinesNum(fmt.Sprintf("%v/net/udp6", root)) 
	if err != nil {
		return udp4 + udp6, nil
	}

	// 返回 udp4 和 udp6 连接数之和
	return udp4 + udp6, nil 
}
