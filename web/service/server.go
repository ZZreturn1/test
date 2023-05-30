package service

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"io"
	"io/fs"
	"net/http"
	"os"
	"runtime"
	"time"
	"x-ui/logger"
	"x-ui/util/sys"
	"x-ui/xray"
)

type ProcessState string	// 进程状态类型

// 定义进程状态常量
const (
	Running ProcessState = "running"	// 运行中
	Stop    ProcessState = "stop"		// 停止
	Error   ProcessState = "error"		// 错误
)

// 定义服务器状态结构体
type Status struct {
	T   time.Time `json:"-"`	// 时间戳
	Cpu float64   `json:"cpu"`	// CPU使用率
	Mem struct {
		Current uint64 `json:"current"`	// 当前内存使用量
		Total   uint64 `json:"total"`	// 总内存量
	} `json:"mem"`	// 内存使用情况
	Swap struct {
		Current uint64 `json:"current"`	// 当前交换空间使用量
		Total   uint64 `json:"total"`		// 总交换空间量
	} `json:"swap"`
	Disk struct {
		Current uint64 `json:"current"`	// 当前磁盘使用量
		Total   uint64 `json:"total"`	// 总磁盘空间量
	} `json:"disk"`
	Xray struct {
		State    ProcessState `json:"state"`	// Xray进程状态
		ErrorMsg string       `json:"errorMsg"`	// Xray错误信息
		Version  string       `json:"version"`	// Xray版本号
	} `json:"xray"`
	Uptime   uint64    `json:"uptime"`	// 服务器运行时间
	Loads    []float64 `json:"loads"`	// 系统负载
	TcpCount int       `json:"tcpCount"`	// TCP连接数
	UdpCount int       `json:"udpCount"`	// UDP连接数
	NetIO    struct {
		Up   uint64 `json:"up"`	// 网络上传速率
		Down uint64 `json:"down"`	// 网络下载速率
	} `json:"netIO"`
	NetTraffic struct {
		Sent uint64 `json:"sent"`	// 发送的网络流量
		Recv uint64 `json:"recv"`	// 接收的网络流量
	} `json:"netTraffic"`
}

// 定义Release结构体，用于获取Xray的版本信息
type Release struct {
	TagName string `json:"tag_name"`	// 版本标签名
}

// ServerService提供服务器相关的服务功能
type ServerService struct {
	xrayService XrayService	// Xray服务
}

// GetStatus获取服务器状态信息
func (s *ServerService) GetStatus(lastStatus *Status) *Status {
	now := time.Now()		// 当前时间
	// 创建状态对象并设置时间戳
	status := &Status{
		T: now,
	}

	// 获取CPU使用率
	percents, err := cpu.Percent(0, false)
	// 
	if err != nil {
		logger.Warning("获取CPU使用率失败:", err)	// 
	} else {
		status.Cpu = percents[0]	// 
	}

	// 获取服务器运行时间
	upTime, err := host.Uptime()
	if err != nil {
		logger.Warning("获取服务器运行时间失败:", err)
	} else {
		status.Uptime = upTime
	}

	// 获取内存使用情况
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		logger.Warning("获取内存使用情况失败:", err)
	} else {
		status.Mem.Current = memInfo.Used
		status.Mem.Total = memInfo.Total
	}

	// 获取交换空间使用情况
	swapInfo, err := mem.SwapMemory()
	if err != nil {
		logger.Warning("获取交换空间使用情况失败:", err)
	} else {
		status.Swap.Current = swapInfo.Used
		status.Swap.Total = swapInfo.Total
	}

	// 获取磁盘使用情况
	distInfo, err := disk.Usage("/")
	if err != nil {
		logger.Warning("获取磁盘使用情况失败:", err)
	} else {
		status.Disk.Current = distInfo.Used
		status.Disk.Total = distInfo.Total
	}

	// 获取系统负载情况
	avgState, err := load.Avg()
	if err != nil {
		logger.Warning("获取系统负载情况失败:", err)
	} else {
		status.Loads = []float64{avgState.Load1, avgState.Load5, avgState.Load15}
	}

	// 获取网络IO统计信息
	ioStats, err := net.IOCounters(false)
	if err != nil {
		logger.Warning("获取网络IO统计信息失败:", err)
	} else if len(ioStats) > 0 {
		ioStat := ioStats[0]		// 获取第一个网络接口的统计信息
		status.NetTraffic.Sent = ioStat.BytesSent // 设置发送流量
		status.NetTraffic.Recv = ioStat.BytesRecv // 设置接收流量

		// 计算网络上传和下载速率
		if lastStatus != nil {
			duration := now.Sub(lastStatus.T) 	// 计算时间间隔
			seconds := float64(duration) / float64(time.Second) 	// 转换为秒数
			up := uint64(float64(status.NetTraffic.Sent-lastStatus.NetTraffic.Sent) / seconds) 	// 计算上行速率
			down := uint64(float64(status.NetTraffic.Recv-lastStatus.NetTraffic.Recv) / seconds) 	// 计算下行速率
			status.NetIO.Up = up 	// 设置上行速率
			status.NetIO.Down = down 	// 设置下行速率
		}
	} else {
		logger.Warning("无法找到网络IO统计信息")
	}

	// 获取TCP连接数
	status.TcpCount, err = sys.GetTCPCount()	// 获取TCP连接数
	if err != nil {
		logger.Warning("获取TCP连接数失败:", err)
	}

	// 获取UDP连接数
	status.UdpCount, err = sys.GetUDPCount()	// 获取UDP连接数
	if err != nil {
		logger.Warning("获取UDP连接数失败:", err)
	}

	// 检查Xray进程状态和错误信息
	if s.xrayService.IsXrayRunning() {
		status.Xray.State = Running		// 设置Xray状态为运行中
		status.Xray.ErrorMsg = ""		// 清空错误消息
	} else {
		err := s.xrayService.GetXrayErr()	// 获取Xray错误信息
		if err != nil {
			status.Xray.State = Error	// 设置Xray状态为错误
		} else {
			status.Xray.State = Stop	// 设置Xray状态为停止
		}
		status.Xray.ErrorMsg = s.xrayService.GetXrayResult()	// 设置Xray错误消息
	}
	status.Xray.Version = s.xrayService.GetXrayVersion()	// 获取Xray版本号

	return status	// 返回状态信息
}

// GetXrayVersions 返回可用的Xray版本列表
func (s *ServerService) GetXrayVersions() ([]string, error) {
	url := "https://api.github.com/repos/XTLS/Xray-core/releases"	// Xray版本列表API地址
	resp, err := http.Get(url)	// 发起HTTP GET请求
	if err != nil {
		return nil, err	// 返回错误信息
	}

	defer resp.Body.Close() 	// 延迟关闭响应体
	buffer := bytes.NewBuffer(make([]byte, 8192)) 	// 创建字节缓冲区
	buffer.Reset()	 // 重置缓冲区
	_, err = buffer.ReadFrom(resp.Body)	 // 从响应体读取数据到缓冲区
	if err != nil {
		return nil, err	 // 返回错误信息
	}

	releases := make([]Release, 0) 	// 创建Release切片
	err = json.Unmarshal(buffer.Bytes(), &releases)	 // 解析JSON数据到Release切片
	if err != nil {
		return nil, err	 // 返回错误信息
	}
	versions := make([]string, 0, len(releases)) 	// 创建版本号切片
	for _, release := range releases {
		versions = append(versions, release.TagName) 	// 将版本号添加到切片中
	}
	return versions, nil		// 返回版本号切片
}

// downloadXRay 下载指定版本的Xray二进制文件并返回文件名
func (s *ServerService) downloadXRay(version string) (string, error) {
	osName := runtime.GOOS	// 获取操作系统名称
	arch := runtime.GOARCH	// 获取CPU架构

	// 调整操作系统名称和架构名称
	switch osName {
	case "darwin":
		osName = "macos"	// 将"darwin"替换为"macos"
	}

	switch arch {
	case "amd64":
		arch = "64"	// 将"amd64"替换为"64"
	case "arm64":
		arch = "arm64-v8a"	// 将"arm64"替换为"arm64-v8a"
	}

	// 构造文件名和下载链接
	fileName := fmt.Sprintf("Xray-%s-%s.zip", osName, arch)
	url := fmt.Sprintf("https://github.com/XTLS/Xray-core/releases/download/%s/%s", version, fileName)

	// 发起HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		return "", err	// 返回错误信息
	}
	defer resp.Body.Close()	// 延迟关闭响应体

	os.Remove(fileName)	// 删除已存在的文件
	file, err := os.Create(fileName)	// 创建文件
	if err != nil {
		return "", err	// 返回错误信息
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// UpdateXray 更新Xray到指定版本
func (s *ServerService) UpdateXray(version string) error {
	zipFileName, err := s.downloadXRay(version)	// 下载Xray文件
	if err != nil {
		return err		// 返回错误信息
	}

	zipFile, err := os.Open(zipFileName)	// 打开下载的压缩文件
	if err != nil {
		return err		// 返回错误信息
	}
	// 
	defer func() {
		zipFile.Close()	// 关闭压缩文件
		os.Remove(zipFileName)	// 删除临时压缩文件
	}()

	stat, err := zipFile.Stat()
	if err != nil {
		return err		// 返回错误信息
	}
	reader, err := zip.NewReader(zipFile, stat.Size())	// 创建zip.Reader读取压缩文件内容
	if err != nil {
		return err		// 返回错误信息
	}

	s.xrayService.StopXray()	// 停止Xray服务
	defer func() {
		err := s.xrayService.RestartXray(true)	// 重启Xray服务
		if err != nil {
			logger.Error("启动Xray失败:", err)	// 输出错误日志
		}
	}()

	// 复制压缩文件内容到指定目录
	copyZipFile := func(zipName string, fileName string) error {
		zipFile, err := reader.Open(zipName)	// 打开压缩文件中的文件
		if err != nil {
			return err		// 返回错误信息
		}

		os.Remove(fileName)	// 删除已存在的文件

		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, fs.ModePerm)
		if err != nil {
			return err		// 返回错误信息
		}
		defer file.Close()	// 关闭文件
		_, err = io.Copy(file, zipFile)	// 复制文件内容
		return err
	}

	err = copyZipFile("xray", xray.GetBinaryPath())	// 复制xray文件
	if err != nil {
		return err		// 返回错误信息
	}
	err = copyZipFile("geosite.dat", xray.GetGeositePath())	// 复制geosite.dat文件
	if err != nil {
		return err		// 返回错误信息
	}
	err = copyZipFile("geoip.dat", xray.GetGeoipPath())	// 复制geoip.dat文件
	if err != nil {
		return err		// 返回错误信息
	}

	return nil		// 返回nil表示更新成功

}
