package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xinyeOfficalWebsite/pkg/common"
)

func CheckUserLogin(c *gin.Context) {
	AuthenticateToken(c, "user")
	c.Next() // 继续处理下一个中间件或路由处理程序
}

func CheckAdminLogin(c *gin.Context) {
	AuthenticateToken(c, "admin")
	c.Next() // 继续处理下一个中间件或路由处理程序
}

func CheckNotAsk(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString != "" {
		claims, _ := common.ParseToken(tokenString, "user")
		c.Set("user", claims.Userid)
	}
	c.Next()
}

func AuthenticateToken(c *gin.Context, end string) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.AbortWithStatus(http.StatusUnauthorized) // 如果没有提供令牌，则返回未授权状态码
	}
	claims, err := common.ParseToken(tokenString, end)
	if err != nil {
		common.HandleError(c, http.StatusUnauthorized, "请登录后重试")
		c.Abort()
	}
	c.Set("user", claims.Userid)
}
