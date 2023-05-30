package controller

// 导入所需的包
import (
	"github.com/gin-gonic/gin"		// 导入gin包
	"net/http"			// 导入http包
	"x-ui/web/session"			// 导入session包
)

// BaseController是一个基础控制器结构体
type BaseController struct {
}

// checkLogin用于检查用户登录状态的中间件
func (a *BaseController) checkLogin(c *gin.Context) {
	// 如果用户未登录
	if !session.IsLogin(c) {
		// 如果是Ajax请求
		if isAjax(c) {
			pureJsonMsg(c, false, "登录时效已过，请重新登录")	// 返回纯JSON格式的错误信息
		} else {
			c.Redirect(http.StatusTemporaryRedirect, c.GetString("base_path"))		// 重定向到基础路径
		}
		c.Abort()	// 终止请求处理
	} else {
		c.Next()	// 用户已登录，继续处理请求
	}
}