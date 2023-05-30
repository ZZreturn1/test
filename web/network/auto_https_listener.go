package network

import "net"

// AutoHttpsListener是自动HTTPS监听器
type AutoHttpsListener struct {
	net.Listener	 // 网络监听器
}

// NewAutoHttpsListener创建一个新的自动HTTPS监听器
func NewAutoHttpsListener(listener net.Listener) net.Listener {
	return &AutoHttpsListener{
		Listener: listener,	// 使用传入的监听器初始化AutoHttpsListener
	}
}

// Accept接受新的连接
func (l *AutoHttpsListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()	// 接受新的连接
	if err != nil {
		return nil, err	// 返回错误
	}
	// 返回新的自动HTTPS连接
	return NewAutoHttpsConn(conn), nil
}