package controller

import (
	"github.com/gin-gonic/gin"
)

// XUIController是XUI控制器的结构体
type XUIController struct {
	BaseController	// 继承基础控制器

	inboundController *InboundController    // 入站控制器
	settingController *SettingController    // 设置控制器
}

// NewXUIController创建一个新的XUI控制器
func NewXUIController(g *gin.RouterGroup) *XUIController {
	a := &XUIController{}	// 实例化XUI控制器对象
	a.initRouter(g)	// 初始化路由
	return a
}

// initRouter初始化XUI控制器的路由
func (a *XUIController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/xui")	// 创建/xui路由分组
	g.Use(a.checkLogin)	// 应用登录检查中间件 

	g.GET("/", a.index)    // 处理根路径请求
	g.GET("/inbounds", a.inbounds)    // 处理入站列表请求
	g.GET("/setting", a.setting)    // 处理设置请求

	a.inboundController = NewInboundController(g)    // 创建入站控制器
	a.settingController = NewSettingController(g)    // 创建设置控制器
}

// index处理根路径请求
func (a *XUIController) index(c *gin.Context) {
	html(c, "index.html", "系统状态", nil)	// 返回系统状态页面
}

// inbounds处理入站列表请求
func (a *XUIController) inbounds(c *gin.Context) {
	html(c, "inbounds.html", "入站列表", nil)	// 返回入站列表页面
}

// setting处理设置请求
func (a *XUIController) setting(c *gin.Context) {
	html(c, "setting.html", "设置", nil)	// 返回设置页面
}
