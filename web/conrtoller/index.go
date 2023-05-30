package controller

// 导入所需的包
import (
	"net/http" 	// 导入http包
	"time" 		// 导入time包
	"x-ui/logger" 	// 导入logger包
	"x-ui/web/job" 		// 导入job包
	"x-ui/web/service" 		// 导入service包
	"x-ui/web/session" 	// 导入session包

	"github.com/gin-gonic/gin"		// 导入gin包
)

// LoginForm是登录表单结构体
type LoginForm struct {
	Username string json:"username" form:"username" 	// 用户名字段
	Password string json:"password" form:"password" 	// 密码字段
}

// IndexController是首页控制器结构体
type IndexController struct {
	BaseController	// 继承BaseController结构体

	userService service.UserService	// 用户服务
}

// NewIndexController创建一个首页控制器实例
func NewIndexController(g *gin.RouterGroup) *IndexController {
	a := &IndexController{}	// 创建首页控制器实例
	a.initRouter(g)	// 初始化路由
	return a
}

// initRouter用于初始化首页控制器的路由
func (a *IndexController) initRouter(g *gin.RouterGroup) {
	g.GET("/", a.index) 		// 路由：首页
	g.POST("/login", a.login) 	// 路由：登录
	g.GET("/logout", a.logout) 	// 路由：登出
}

// index处理首页请求
func (a *IndexController) index(c *gin.Context) {
	if session.IsLogin(c) {
		c.Redirect(http.StatusTemporaryRedirect, "xui/")		// 如果已登录，重定向到xui页面
		return
	}
	html(c, "login.html", "登录", nil)	// 返回登录页面HTML
}

// login处理登录请求
func (a *IndexController) login(c *gin.Context) {
	var form LoginForm
	err := c.ShouldBind(&form)	// 绑定请求数据到登录表单
	if err != nil {
	pureJsonMsg(c, false, "数据格式错误")   // 返回错误JSON消息
		return
	}

	if form.Username == "" {
		pureJsonMsg(c, false, "请输入用户名")   // 返回错误JSON消息
		return
	}

	if form.Password == "" {
		pureJsonMsg(c, false, "请输入密码")   // 返回错误JSON消息
		return
	}

	user := a.userService.CheckUser(form.Username, form.Password)   // 验证用户登录

	timeStr := time.Now().Format("2006-01-02 15:04:05")   // 获取当前时间字符串

	if user == nil {
		// 用户登录通知任务
		job.NewStatsNotifyJob().UserLoginNotify(form.Username, getRemoteIp(c), timeStr, 0)
		// 输出错误日志
		logger.Infof("wrong username or password: \"%s\" \"%s\"", form.Username, form.Password)
		// 返回错误JSON消息
		pureJsonMsg(c, false, "用户名或密码错误")
		return
	} else {
		logger.Infof("%s 登录成功，IP地址：%s\n", form.Username, getRemoteIp(c))   	// 输出登录成功日志
		job.NewStatsNotifyJob().UserLoginNotify(form.Username, getRemoteIp(c), timeStr, 1)  	 // 用户登录通知任务
	}

	err = session.SetLoginUser(c, user)	// 将登录用户信息设置到会话中
	logger.Info("user", user.Id, "login success")	// 输出登录成功日志
	jsonMsg(c, "登录", err)	// 返回登录结果的JSON消息
}

// logout处理登出请求
func (a *IndexController) logout(c *gin.Context) {
	user := session.GetLoginUser(c)	// 获取登录用户信息
	if user != nil {
		logger.Info("用户", user.Id, "登出")   // 输出登出日志
	}
	session.ClearSession(c)	// 清除会话信息
	c.Redirect(http.StatusTemporaryRedirect, c.GetString("base_path"))		// 重定向到基础路径
}