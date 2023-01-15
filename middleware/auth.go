package middleware

import (
	"douyin/controllers"
	"douyin/utils"
	"github.com/gin-gonic/gin"
)

// GetAuthorization
// token存在query
// 解析token，由于get方法和post方法token的存在位置不同因此需要分别解析
// 验证失败全部redirect到登录页面
func GetAuthorization(c *gin.Context) {
	token := c.Query("token")
	auth, err := controllers.ParseToken(token)

	// 验证失败，返回登录页面登录
	if err != nil {
		utils.Redirect(c, "/douyin/user/login")

	}
	c.Set("auth", *auth)
	c.Next()

}

// PostAuthorization
// token存在form表单中
// 解析token，由于get方法和post方法token的存在位置不同因此需要分别解析
// 验证失败全部redirect到登录页面
func PostAuthorization(c *gin.Context) {
	token := c.PostForm("token")
	auth, err := controllers.ParseToken(token)

	// 验证失败，返回登录页面登录
	if err != nil {
		utils.Redirect(c, "/douyin/user/login")

	}
	c.Set("auth", *auth)
	c.Next()
}
