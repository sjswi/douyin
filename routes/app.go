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
		home.GET("/user", controllers.UserInfo)

		user := home.Group("/user")
		{
			user.POST("/register", controllers.UserRegister)
			user.POST("/login", controllers.UserLogin)
		}

		publish := home.Group("/publish")
		{
			publish.POST("/action", controllers.PublishAction).Use(middleware.PostAuthorization)
			publish.GET("/list", controllers.PublishList).Use(middleware.GetAuthorization)
		}
		favorite := home.Group("/favorite", middleware.GetAuthorization)
		{
			favorite.POST("/action", controllers.FavoriteAction)
			favorite.GET("/list", controllers.FavoriteList)

		}
		comment := home.Group("/comment", middleware.GetAuthorization)
		{
			comment.POST("/action", controllers.CommentAction)
			comment.GET("/list", controllers.CommentList)
		}
		relation := home.Group("/relation", middleware.GetAuthorization)
		{
			relation.POST("/action", controllers.RelationAction)
			//relation.GET("/list", controllers.RelationList)
			relation.GET("/follow/list", controllers.RelationFollowList)
			relation.GET("/follower/list", controllers.RelationFollowerList)
			relation.GET("/friend/list", controllers.RelationFriendList)
		}
		message := home.Group("/message", middleware.GetAuthorization)
		{
			message.GET("/chat", controllers.MessageChat)
			message.POST("/action", controllers.MessageAction)
		}
	}
}
