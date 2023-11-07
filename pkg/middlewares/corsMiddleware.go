package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CorsMiddleware(c *gin.Context) {
	//允许访问所有域
	method := c.Request.Method
	// 必须，接受指定域的请求，可以使用*不加以限制。* 常常使用在公开的无敏感信息的接口
	c.Header("Access-Control-Allow-Origin", c.GetHeader("*"))
	fmt.Println(c.GetHeader("Origin"))
	// 必须，设置服务器支持的所有跨域请求的方法
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
	// 服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Token")
	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()

}
