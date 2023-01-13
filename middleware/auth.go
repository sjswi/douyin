package middleware

import (
	"douyin/controllers"
	"douyin/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
)

func Authorization(c *gin.Context) {
	if c.Request.Method == "GET" {
		token := c.Query("token")
		auth, err := controllers.ParseToken(token)

		// 验证失败，返回登录页面登录
		if err != nil {
			utils.Redirect(c, "/douyin/user/login")

		}
		c.Set("auth", auth)
		c.Next()
	} else if c.Request.Method == "POST" {
		all, err := io.ReadAll(c.Request.Body)
		if err != nil {
			utils.Redirect(c, "/douyin/user/login")
		}
		var temp map[string]interface{}
		err = json.Unmarshal(all, &temp)
		if err != nil {
			utils.Redirect(c, "/douyin/user/login")
		}
		if value, ok := temp["token"]; ok {
			token := value.(string)
			auth, err := controllers.ParseToken(token)

			// 验证失败，返回登录页面登录
			if err != nil {
				utils.Redirect(c, "/douyin/user/login")

			}
			c.Set("auth", auth)
			c.Next()
		}
		utils.Redirect(c, "/douyin/user/login")

	} else {
		utils.Redirect(c, "/douyin/user/login")
	}

}
