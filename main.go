package main

import (
	"douyin/bootstrap/driver"
	"douyin/config"
	"douyin/routes"
	"github.com/gin-gonic/gin"
)

// @title douyin
// @version 1.0
// @description 青训营抖音大项目
// @termsOfService http://swagger.io/terms/

// @contact.name 余晓兵
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath /douyin
func main() {

	App := gin.Default()
	config.InitConfig()
	driver.InitConn("mysql")
	driver.InitOSS()
	driver.InitRedis()

	routes.Home(App)
	//App.NoRoute(func(c *gin.Context) {
	//	c.HTML(http.StatusOK, "404", nil)
	//})
	App.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
