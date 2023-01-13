package controllers

import (
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RelationList struct {
	utils.Response
	UserList []*User `json:"user_list"`
}

// RelationAction
// @Summary 关注和取消关注操作
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param to_user_id query int true "用户id"
// @Param token query string true "token"
// @Param action_type query int true "操作类型"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/relation/action [post]
func RelationAction(c *gin.Context) {
	//TODO
	// 业务代码
	response := utils.Response{
		StatusCode: 0,
		StatusMsg:  "",
	}
	c.JSON(http.StatusOK, response)
}

// RelationFollowList
// func RelationList(c *gin.Context) {
//
// }
// RelationFollowList
// @Summary 获取关注列表
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param to_user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/relation/follow/list [get]
func RelationFollowList(c *gin.Context) {

	//TODO
	// 业务代码
	response := RelationList{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	c.JSON(http.StatusOK, response)
}

// RelationFollowerList
// @Summary 获取关注者列表
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param to_user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/relation/follower/list [get]
func RelationFollowerList(c *gin.Context) {

	//TODO
	// 业务代码
	response := RelationList{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	c.JSON(http.StatusOK, response)
}

// RelationFriendList
// @Summary 获取聊天记录
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param to_user_id query int true "用户id"
// @Param token query string true "token"
// @Param action_type query string true "操作类型"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/relation/friend/list [get]
func RelationFriendList(c *gin.Context) {

	//TODO
	// 业务代码
	response := RelationList{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	c.JSON(http.StatusOK, response)
}
