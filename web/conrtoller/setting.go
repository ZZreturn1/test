package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"time"
	"x-ui/web/entity"
	"x-ui/web/service"
	"x-ui/web/session"
)

// updateUserForm定义了更新用户的表单结构
type updateUserForm struct {
	OldUsername string `json:"oldUsername" form:"oldUsername"`	// 原用户名
	OldPassword string `json:"oldPassword" form:"oldPassword"`	// 原密码
	NewUsername string `json:"newUsername" form:"newUsername"`	// 新用户名
	NewPassword string `json:"newPassword" form:"newPassword"`	// 新密码
}

// SettingController定义了设置页面的控制器 
type SettingController struct {
	settingService service.SettingService	// 设置服务
	userService    service.UserService	// 用户服务
	panelService   service.PanelService	// 面板服务
}

// NewSettingController创建SettingController的实例
func NewSettingController(g *gin.RouterGroup) *SettingController {
	a := &SettingController{}	// 创建SettingController对象
	a.initRouter(g)	// 初始化路由
	return a
}

// initRouter初始化路由配置
func (a *SettingController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/setting")	// 创建setting路由组

	g.POST("/all", a.getAllSetting)         // 获取所有设置的路由
	g.POST("/update", a.updateSetting)      // 更新设置的路由
	g.POST("/updateUser", a.updateUser)     // 更新用户的路由
	g.POST("/restartPanel", a.restartPanel) // 重启面板的路由
}

// getAllSetting获取所有设置
func (a *SettingController) getAllSetting(c *gin.Context) {
	allSetting, err := a.settingService.GetAllSetting()	// 获取所有设置
	if err != nil {
		jsonMsg(c, "获取设置", err)	// 返回获取设置的JSON消息 
		return
	}
	jsonObj(c, allSetting, nil)	// 返回所有设置的JSON对象
}

// updateSetting更新设置
func (a *SettingController) updateSetting(c *gin.Context) {
	allSetting := &entity.AllSetting{}	// 创建AllSetting对象
	err := c.ShouldBind(allSetting)	// 绑定请求参数到AllSetting对象
	if err != nil {
		jsonMsg(c, "修改设置", err)	// 返回修改设置的JSON消息
		return
	}
	err = a.settingService.UpdateAllSetting(allSetting)	// 更新所有设置
	jsonMsg(c, "修改设置", err)	// 返回修改设置的JSON消息
}

// updateUser更新用户
func (a *SettingController) updateUser(c *gin.Context) {
	form := &updateUserForm{}	// 创建updateUserForm对象
	err := c.ShouldBind(form)	// 绑定请求参数到updateUserForm对象
	if err != nil {
		jsonMsg(c, "修改用户", err)	// 返回修改用户的JSON消息
		return
	}

	user := session.GetLoginUser(c)	// 获取登录用户
	if user.Username != form.OldUsername || user.Password != form.OldPassword {
		jsonMsg(c, "修改用户", errors.New("原用户名或原密码错误"))	// 返回错误消息
		return
	}

	if form.NewUsername == "" || form.NewPassword == "" {
		jsonMsg(c, "修改用户", errors.New("新用户名和新密码不能为空"))	// 返回错误消息
		return
	}

	// 更新用户
	err = a.userService.UpdateUser(user.Id, form.NewUsername, form.NewPassword)
	if err == nil {
		user.Username = form.NewUsername
		user.Password = form.NewPassword
		session.SetLoginUser(c, user)
	}
	jsonMsg(c, "修改用户", err)	// 返回修改用户的JSON消息
}

// restartPanel重启面板
func (a *SettingController) restartPanel(c *gin.Context) {
	err := a.panelService.RestartPanel(time.Second * 3)	// 重启面板
	jsonMsg(c, "重启面板", err)	// 返回重启面板的JSON消息
}