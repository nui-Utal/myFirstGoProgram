package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"xinyeOfficalWebsite/pkg/middlewares"
	"xinyeOfficalWebsite/pkg/routers"
	"xinyeOfficalWebsite/pkg/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.InitRedis()

	r := gin.Default()
	// 应用全局中间件
	r.Static("/web/", "./web/")       // 静态资源访问
	r.Use(middlewares.CorsMiddleware) // 跨域请求

	// 加载路由
	routers.LoadRoutes(r)

	port := viper.GetString("port")
	err := r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("监听端口", "http://127.0.0.1:"+port)
}
