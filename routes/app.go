package routes

import (
	"douyin/controllers"
	"douyin/docs"
	"douyin/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Home(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	home := r.Group("/douyin")
	{
		// 首页
		home.GET("/feed", controllers.Feed)

		// 个人中心
		user := home.Group("/user", middleware.Authorization)
		{
			user.GET("/", controllers.UserInfo)
			user.POST("/register", controllers.UserRegister)
			user.POST("/login", controllers.UserLogin)

		}
		publish := home.Group("/publish", middleware.Authorization)
		{
			publish.POST("/action", controllers.PublishAction)
			publish.GET("/list", controllers.PublishList)
		}
		favorite := home.Group("/favorite", middleware.Authorization)
		{
			favorite.POST("/action", controllers.FavoriteAction)
			favorite.GET("/list", controllers.FavoriteList)

		}
		comment := home.Group("/comment", middleware.Authorization)
		{
			comment.POST("/action", controllers.CommentAction)
			comment.GET("/list", controllers.CommentList)
		}
		relation := home.Group("/relation", middleware.Authorization)
		{
			relation.POST("/action", controllers.RelationAction)
			//relation.GET("/list", controllers.RelationList)
			relation.GET("/follow/list", controllers.RelationFollowList)
			relation.GET("/follower/list", controllers.RelationFollowerList)
			relation.GET("/friend/list", controllers.RelationFriendList)
		}
		message := home.Group("/message", middleware.Authorization)
		{
			message.GET("/chat", controllers.MessageChat)
			message.POST("/action", controllers.MessageAction)
		}
	}
}
