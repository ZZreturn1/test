package sys

import (
	_ "unsafe"
)

//go:linkname HostProc github.com/shirou/gopsutil/internal/common.HostProc
// HostProc 是一个外部函数的链接名称，位于 github.com/shirou/gopsutil/internal/common 包中。
// 它可能用于获取与主机进程相关的信息。
func HostProc(combineWith ...string) string