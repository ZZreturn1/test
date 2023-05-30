package model

// 导入所需的包
import (
	"fmt"    // 导入 fmt 包提供格式化和打印功能
	"x-ui/util/json_util"
	"x-ui/xray"
)

// Protocol 日志级别类型声明
type Protocol string

const (
	VMess       Protocol = "vmess"           // VMess 协议
	VLESS       Protocol = "vless"           // VLESS 协议
	Dokodemo    Protocol = "Dokodemo-door"   // Dokodemo 协议
	Http        Protocol = "http"            // HTTP 协议
	Trojan      Protocol = "trojan"          // Trojan 协议
	Shadowsocks Protocol = "shadowsocks"     // Shadowsocks 协议
)

type User struct {
	Id       int    `json:"id" gorm:"primaryKey;autoIncrement"`  // 用户ID
	Username string `json:"username"`                             // 用户名
	Password string `json:"password"`                             // 密码
}

type Inbound struct {
	Id         int    `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`  // 入站ID
	UserId     int    `json:"-"`                                             // 用户ID
	Up         int64  `json:"up" form:"up"`                                  // 上行流量
	Down       int64  `json:"down" form:"down"`                              // 下行流量
	Total      int64  `json:"total" form:"total"`                            // 总流量
	Remark     string `json:"remark" form:"remark"`                          // 备注
	Enable     bool   `json:"enable" form:"enable"`                          // 是否启用
	ExpiryTime int64  `json:"expiryTime" form:"expiryTime"`                  // 过期时间

	// 配置部分
	Listen         string   `json:"listen" form:"listen"`                 // 监听地址
	Port           int      `json:"port" form:"port" gorm:"unique"`       // 端口号
	Protocol       Protocol `json:"protocol" form:"protocol"`             // 协议类型
	Settings       string   `json:"settings" form:"settings"`             // 协议设置
	StreamSettings string   `json:"streamSettings" form:"streamSettings"` // 传输配置
	Tag            string   `json:"tag" form:"tag" gorm:"unique"`         // 标签
	Sniffing       string   `json:"sniffing" form:"sniffing"`             // 流量识别设置
}

// GenXrayInboundConfig 生成 Xray 的入站配置
func (i *Inbound) GenXrayInboundConfig() *xray.InboundConfig {
                // 获取监听地址
	listen := i.Listen
	if listen != "" {
		listen = fmt.Sprintf("\"%v\"", listen)
	}

                // 返回生成的 Xray 入站配置
	return &xray.InboundConfig{
		Listen: json_util.RawMessage(listen),     // 监听地址
		Port: i.Port,                            // 端口号
		Protocol: string(i.Protocol),                // 协议类型
		Settings:  json_util.RawMessage(i.Settings),  // 设置
		StreamSettings: json_util.RawMessage(i.StreamSettings),  // 传输设置
		Tag: i.Tag,                             // 标签
		Sniffing: json_util.RawMessage(i.Sniffing),  // 数据嗅探
	}
}

type Setting struct {
	Id    int    `json:"id" form:"id" gorm:"primaryKey;autoIncrement"`  // 设置ID
	Key   string `json:"key" form:"key"`                                // 键名
	Value string `json:"value" form:"value"`                            // 键值
}