package controller

// 导入所需的包
import (
	"github.com/gin-gonic/gin"
	"time"
	"x-ui/web/global"
	"x-ui/web/service"
)

// ServerController处理服务器相关请求
type ServerController struct {
	BaseController	// 继承基础控制器

	serverService service.ServerService // 服务器服务

	lastStatus        *service.Status // 最后的状态
	lastGetStatusTime time.Time       // 最后获取状态的时间

	lastVersions        []string  // 最后的版本
	lastGetVersionsTime time.Time // 最后获取版本的时间
}

// NewServerController创建一个新的ServerController实例
func NewServerController(g *gin.RouterGroup) *ServerController {
	a := &ServerController{
		lastGetStatusTime: time.Now(),	// 初始化最后获取状态的时间
	}
	a.initRouter(g) 	// 初始化路由
	a.startTask() 	// 启动任务
	return a
}

// initRouter初始化服务器路由
func (a *ServerController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/server")	// 

	g.Use(a.checkLogin)                   // 添加登录检查中间件
	g.POST("/status", a.status)           // 处理获取服务器状态的请求
	g.POST("/getXrayVersion", a.getXrayVersion)         // 处理获取Xray版本的请求
	g.POST("/installXray/:version", a.installXray)       // 处理安装Xray的请求
}

// refreshStatus刷新服务器状态
func (a *ServerController) refreshStatus() {
	a.lastStatus = a.serverService.GetStatus(a.lastStatus)	// 获取最新的服务器状态
}

// startTask启动定时任务
func (a *ServerController) startTask() {
	webServer := global.GetWebServer()	// 获取全局Web服务器实例
	c := webServer.GetCron()	// 获取定时任务管理器
	// 每2秒执行一次任务
	c.AddFunc("@every 2s", func() {
		now := time.Now()	// 当前时间
		// 如果超过3分钟未获取状态，则不执行任务
		if now.Sub(a.lastGetStatusTime) > time.Minute*3 {
			return
		}
		a.refreshStatus()	// 刷新服务器状态
	})
}

// status处理获取服务器状态的请求
func (a *ServerController) status(c *gin.Context) {
	a.lastGetStatusTime = time.Now()	// 更新最后获取状态的时间

	jsonObj(c, a.lastStatus, nil)	// 返回服务器状态的JSON对象
}

// getXrayVersion处理获取Xray版本的请求
func (a *ServerController) getXrayVersion(c *gin.Context) {
	now := time.Now()		// 当前时间
	if now.Sub(a.lastGetVersionsTime) <= time.Minute {
		jsonObj(c, a.lastVersions, nil)		// 如果在1分钟内已获取过版本，则直接返回缓存的版本列表
		return
	}
	
	versions, err := a.serverService.GetXrayVersions() // 获取Xray版本列表
	if err != nil {
		jsonMsg(c, "获取版本", err) // 返回获取版本失败的JSON消息
		return
	}

	a.lastVersions = versions       // 更新最后的版本列表
	a.lastGetVersionsTime = time.Now() // 更新最后获取版本的时间

	jsonObj(c, versions, nil) // 返回版本列表的JSON对象
}

// installXray处理安装Xray的请求
func (a *ServerController) installXray(c *gin.Context) {
	version := c.Param("version")	// 获取URL参数中的版本号
	err := a.serverService.UpdateXray(version)	// 安装指定版本的Xray
	jsonMsg(c, "安装 xray", err)	// 返回安装Xray的JSON消息
}