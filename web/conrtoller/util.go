package controller

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"strings"
	"x-ui/config"
	"x-ui/logger"
	"x-ui/web/entity"
)

// getUriId从上下文中获取URI中的ID参数
func getUriId(c *gin.Context) int64 {
	// 定义结构体s来接收URI参数中的ID值
	s := struct {
		Id int64 `uri:"id"`	// 使用uri标签指定字段名为id
	}{}

	_ = c.BindUri(&s)	// 将URI参数绑定到结构体s中
	return s.Id
}

// getRemoteIp获取请求的远程IP地址
func getRemoteIp(c *gin.Context) string {
	value := c.GetHeader("X-Forwarded-For")	// 从请求头中获取X-Forwarded-For字段的值
	if value != "" {
		ips := strings.Split(value, ",")	// 根据逗号分隔IP地址列表
		return ips[0]	// 返回第一个IP地址
	} else {
		addr := c.Request.RemoteAddr	// 获取请求的远程地址
		ip, _, _ := net.SplitHostPort(addr)	// 将远程地址拆分为IP地址和端口号
		return ip	// 返回IP地址
	}
}

// jsonMsg将消息以JSON格式返回给客户端
func jsonMsg(c *gin.Context, msg string, err error) {
	jsonMsgObj(c, msg, nil, err)
}

// jsonObj将对象以JSON格式返回给客户端
func jsonObj(c *gin.Context, obj interface{}, err error) {
	jsonMsgObj(c, "", obj, err)	// 
}

// jsonMsgObj将带有消息和对象的JSON响应返回给客户端
func jsonMsgObj(c *gin.Context, msg string, obj interface{}, err error) {
	m := entity.Msg{
		Obj: obj,	// 将对象作为响应中的Obj字段
	}
	if err == nil {
		m.Success = true
		if msg != "" {
			m.Msg = msg + "成功"
		}
	} else {
		m.Success = false
		m.Msg = msg + "失败: " + err.Error()
		logger.Warning(msg+"失败: ", err)
	}
	c.JSON(http.StatusOK, m)
}

// pureJsonMsg将纯粹的JSON消息返回给客户端
func pureJsonMsg(c *gin.Context, success bool, msg string) {
	if success {
		c.JSON(http.StatusOK, entity.Msg{
			Success: true,
			Msg:     msg,
		})
	} else {
		c.JSON(http.StatusOK, entity.Msg{
			Success: false,
			Msg:     msg,
		})
	}
}

// html将HTML模板渲染并返回给客户端
func html(c *gin.Context, name string, title string, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	data["title"] = title		// 设置模板中的title变量
	data["request_uri"] = c.Request.RequestURI 	// 设置模板中的request_uri变量
	data["base_path"] = c.GetString("base_path") 	// 设置模板中的base_path变量
	c.HTML(http.StatusOK, name, getContext(data)) // 渲染HTML模板并传入数据
}

// getContext获取模板渲染的上下文数据
func getContext(h gin.H) gin.H {
	a := gin.H{
		"cur_ver": config.GetVersion(),	// 设置了当前版本的上下文数据变量cur_ver
	}
	if h != nil {
		for key, value := range h {
			a[key] = value	// 将传入的额外数据添加到上下文数据中
		}
	}
	return a
}

// isAjax判断请求是否为Ajax请求
func isAjax(c *gin.Context) bool {
	// 检查请求头中的X-Requested-With字段是否为XMLHttpRequest，判断是否为Ajax请求
	return c.GetHeader("X-Requested-With") == "XMLHttpRequest"
}