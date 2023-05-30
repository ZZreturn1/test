package job

import (
	"fmt"
	"net"
	"os"

	"time"

	"x-ui/logger"
	"x-ui/util/common"
	"x-ui/web/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// LoginStatus表示登录状态
type LoginStatus byte

const (
	LoginSuccess LoginStatus = 1	// 登录成功
	LoginFail    LoginStatus = 0	// 登录失败
)

// StatsNotifyJob是统计通知任务
type StatsNotifyJob struct {
	enable         bool	// 是否启用
	xrayService    service.XrayService	// Xray服务
	inboundService service.InboundService	// 入站服务
	settingService service.SettingService	// 设置服务
}

// NewStatsNotifyJob创建一个新的统计通知任务
func NewStatsNotifyJob() *StatsNotifyJob {
	return new(StatsNotifyJob)	// 创建并返回一个StatsNotifyJob实例
}

// SendMsgToTgbot向Telegram Bot发送消息
func (j *StatsNotifyJob) SendMsgToTgbot(msg string) {
	tgBottoken, err := j.settingService.GetTgBotToken()
	if err != nil {
		// 获取Telegram Bot Token失败
		logger.Warning("sendMsgToTgbot failed,GetTgBotToken fail:", err)
		return
	}

	tgBotid, err := j.settingService.GetTgBotChatId()
	if err != nil {
		// 获取Telegram Bot Chat ID失败
		logger.Warning("sendMsgToTgbot failed,GetTgBotChatId fail:", err)
		return
	}

	bot, err := tgbotapi.NewBotAPI(tgBottoken)
	if err != nil {
		fmt.Println("get tgbot error:", err)	// 获取Telegram Bot实例失败
		return
	}
	bot.Debug = true	// 调试模式开启
	fmt.Printf("Authorized on account %s", bot.Self.UserName)	// 打印Bot用户名
	info := tgbotapi.NewMessage(int64(tgBotid), msg)	// 创建新的消息实例
	bot.Send(info)
}

// Run运行统计通知任务
func (j *StatsNotifyJob) Run() {
	if !j.xrayService.IsXrayRunning() {
		return
	}
	var info string	// 统计信息
	//
	name, err := os.Hostname()
	if err != nil {
		fmt.Println("get hostname error:", err)	// 获取主机名失败
		return
	}
	info = fmt.Sprintf("主机名称:%s\r\n", name)

	var ip string	// IP地址
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())	// 获取网络接口失败
		return
	}

	// 
	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()

			// 
			for _, address := range addrs {
				// 
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					// 
					if ipnet.IP.To4() != nil {
						ip = ipnet.IP.String()	// 获取IPv4地址
						break
					} else {
						ip = ipnet.IP.String()	// 获取IPv6地址
						break
					}
				}
			}
		}
	}

	info += fmt.Sprintf("IP地址:%s\r\n \r\n", ip)

	inbouds, err := j.inboundService.GetAllInbounds()
	if err != nil {
		logger.Warning("StatsNotifyJob run failed:", err)	// 获取所有入站信息失败
		return
	}

	// NOTE:如果这里没有任何会话，需要在这里通知
	// TODO:分节点推送,自动转化格式
	for _, inbound := range inbouds {
		info += fmt.Sprintf("节点名称:%s\r\n端口:%d\r\n上行流量↑:%s\r\n下行流量↓:%s\r\n总流量:%s\r\n", inbound.Remark, inbound.Port, common.FormatTraffic(inbound.Up), common.FormatTraffic(inbound.Down), common.FormatTraffic((inbound.Up + inbound.Down)))
		if inbound.ExpiryTime == 0 {
			info += fmt.Sprintf("到期时间:无限期\r\n \r\n")	// 到期时间为无限期
		} else {
			// 
			info += fmt.Sprintf("到期时间:%s\r\n \r\n", time.Unix((inbound.ExpiryTime/1000), 0).Format("2006-01-02 15:04:05"))
		}
	}
	j.SendMsgToTgbot(info)	// 发送统计信息至Telegram Bot
}

// UserLoginNotify向用户发送登录通知
func (j *StatsNotifyJob) UserLoginNotify(username string, ip string, time string, status LoginStatus) {
	// 
	if username == "" || ip == "" || time == "" {
		logger.Warning("UserLoginNotify failed,invalid info")	// 无效的用户登录信息
		return
	}
	var msg string	// 通知消息
	//
	name, err := os.Hostname()
	// 
	if err != nil {
		fmt.Println("get hostname error:", err)	// 获取主机名失败
		return
	}
	// 
	if status == LoginSuccess {
		msg = fmt.Sprintf("面板登录成功提醒\r\n主机名称:%s\r\n", name)	// 面板登录成功提醒
	} else if status == LoginFail {
		msg = fmt.Sprintf("面板登录失败提醒\r\n主机名称:%s\r\n", name)	// 面板登录失败提醒
	}
	msg += fmt.Sprintf("时间:%s\r\n", time)	// 登录时间
	msg += fmt.Sprintf("用户:%s\r\n", username)	// 用户名
	msg += fmt.Sprintf("IP:%s\r\n", ip)	// IP地址
	j.SendMsgToTgbot(msg)	// 发送登录通知至Telegram Bot
}