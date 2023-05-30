package network

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"sync"
)

// AutoHttpsConn是自动HTTPS连接
type AutoHttpsConn struct {
	net.Conn		// 网络连接

	firstBuf []byte	// 第一次读取的数据缓冲区
	bufStart int	// 缓冲区起始位置

	readRequestOnce sync.Once		// 读取请求的同步锁
}

// NewAutoHttpsConn创建一个新的自动HTTPS连接
func NewAutoHttpsConn(conn net.Conn) net.Conn {
	return &AutoHttpsConn{
		Conn: conn,	// 使用传入的连接初始化AutoHttpsConn
	}
}

// readRequest读取HTTP请求并进行重定向
func (c *AutoHttpsConn) readRequest() bool {
	c.firstBuf = make([]byte, 2048)	// 创建数据缓冲区
	n, err := c.Conn.Read(c.firstBuf)	// 读取数据
	c.firstBuf = c.firstBuf[:n]	// 截取实际读取的数据
	if err != nil {
		return false	// 返回读取失败
	}
	reader := bytes.NewReader(c.firstBuf)		// 创建读取器
	bufReader := bufio.NewReader(reader)	// 创建缓冲读取器
	request, err := http.ReadRequest(bufReader)	// 读取HTTP请求
	if err != nil {
		return false	// 返回读取失败
	}
	resp := http.Response{
		Header: http.Header{},	// 创建响应头
	}
	resp.StatusCode = http.StatusTemporaryRedirect	// 设置重定向状态码
	location := fmt.Sprintf("https://%v%v", request.Host, request.RequestURI)	// 构建重定向地址
	resp.Header.Set("Location", location)	// 设置重定向头
	resp.Write(c.Conn)		// 发送重定向响应
	c.Close()		// 关闭连接
	c.firstBuf = nil	// 清空缓冲区
	return true	// 返回读取成功
}

// Read读取数据
func (c *AutoHttpsConn) Read(buf []byte) (int, error) {
	c.readRequestOnce.Do(func() {
		c.readRequest()		// 读取请求并进行重定向
	})

	if c.firstBuf != nil {
		n := copy(buf, c.firstBuf[c.bufStart:])		// 将缓冲区的数据拷贝到buf中
		c.bufStart += n		// 更新缓冲区的起始位置
		if c.bufStart >= len(c.firstBuf) {
			c.firstBuf = nil	// 清空缓冲区
		}
		return n, nil	// 返回读取的字节数和nil错误
	}

	return c.Conn.Read(buf)	// 从底层连接读取数据
}
