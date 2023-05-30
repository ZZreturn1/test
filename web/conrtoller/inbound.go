package controller

// 导入所需的包
import (
	"fmt"
	"github.com/gin-gonic/gin"		// 导入gin包
	"strconv" 	// 导入strconv包
	"x-ui/database/model" 	// 导入model包
	"x-ui/logger" 		// 导入logger包
	"x-ui/web/global" 		// 导入global包
	"x-ui/web/service" 		// 导入service包
	"x-ui/web/session" 	// 导入session包
)

// InboundController是一个入站控制器结构体
type InboundController struct {
	inboundService service.InboundService	// 入站服务
	xrayService    service.XrayService	// Xray服务
}

// NewInboundController创建一个入站控制器实例
func NewInboundController(g *gin.RouterGroup) *InboundController {
	a := &InboundController{}	// 创建入站控制器实例
	a.initRouter(g)	// 初始化路由
	a.startTask()	// 启动任务
	return a
}

// initRouter用于初始化入站控制器的路由
func (a *InboundController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/inbound")	// 设置路由组的基础路径为"/inbound"

	g.POST("/list", a.getInbounds)             // 路由：获取入站列表
	g.POST("/add", a.addInbound)                // 路由：添加入站
	g.POST("/del/:id", a.delInbound)            // 路由：删除入站
	g.POST("/update/:id", a.updateInbound)      // 路由：更新入站
}

// startTask用于启动定时任务
func (a *InboundController) startTask() {
	webServer := global.GetWebServer()	// 获取全局Web服务器实例
	c := webServer.GetCron()	// 获取Cron实例
	//
	// 定时执行任务，每10秒执行一次
	c.AddFunc("@every 10s", func() {
		if a.xrayService.IsNeedRestartAndSetFalse() {
			err := a.xrayService.RestartXray(false)   // 重新启动Xray服务，关闭日志输出

			if err != nil {
				logger.Error("重新启动Xray失败：", err)   // 输出错误日志
			}
		}
	})
}

// getInbounds用于获取入站列表
func (a *InboundController) getInbounds(c *gin.Context) {
	// 获取登录用户信息
	user := session.GetLoginUser(c)
	inbounds, err := a.inboundService.GetInbounds(user.Id)   // 调用入站服务获取入站列表

	if err != nil {
		jsonMsg(c, "获取", err)   // 返回错误信息
		return
	}

	jsonObj(c, inbounds, nil)   // 返回JSON格式的入站列表
}

// addInbound用于添加入站
func (a *InboundController) addInbound(c *gin.Context) {
	// 创建一个入站实例
	inbound := &model.Inbound{}
	//
	err := c.ShouldBind(inbound)   // 绑定请求数据到入站实例

	if err != nil {
		jsonMsg(c, "添加", err)   // 返回错误信息
		return
	}
	user := session.GetLoginUser(c)	// 获取登录用户信息
	inbound.UserId = user.Id	// 设置入站的用户ID
	inbound.Enable = true	// 设置入站为启用状态
	inbound.Tag = fmt.Sprintf("inbound-%v", inbound.Port)	// 设置入站标签，格式为"inbound-端口号"
	err = a.inboundService.AddInbound(inbound)	// 调用入站服务添加入站
	jsonMsg(c, "添加", err)	// 返回添加结果的JSON消息

	if err == nil {
		a.xrayService.SetToNeedRestart()	// 设置Xray服务需要重新启动
	}
}

// delInbound用于删除入站
func (a *InboundController) delInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))	// 获取要删除的入站ID
	//
	if err != nil {
		jsonMsg(c, "删除", err)   // 返回错误信息
		return
	}

	err = a.inboundService.DelInbound(id)   // 调用入站服务删除入站

	jsonMsg(c, "删除", err)   // 返回删除结果的JSON消息

	if err == nil {
		a.xrayService.SetToNeedRestart()   // 设置Xray服务需要重新启动
	}
}

// updateInbound用于更新入站
func (a *InboundController) updateInbound(c *gin.Context) {
	// 获取要更新的入站ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "修改", err)   // 返回错误信息
		return
	}

	inbound := &model.Inbound{
		Id: id,   // 设置入站的ID
	}

	err = c.ShouldBind(inbound)   // 绑定请求数据到入站实例

	if err != nil {
		jsonMsg(c, "修改", err)   // 返回错误信息
		return
	}

	err = a.inboundService.UpdateInbound(inbound)   // 调用入站服务更新入站

	jsonMsg(c, "修改", err)   // 返回更新结果的JSON消息

	if err == nil {
		a.xrayService.SetToNeedRestart()   // 设置Xray服务需要重新启动
	}
}