package controllers

import (
	auth2 "douyin/auth"
	"douyin/service"
	"douyin/utils"
	"douyin/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//TODO
// 关注，粉丝，好友列表机器相似，可以整合为一个函数

type RelationList struct {
	utils.Response
	UserList []vo.User `json:"user_list"`
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
// @Router /douyin/relation/action/ [post]
func RelationAction(c *gin.Context) {
	successResponse := utils.Response{
		StatusCode: 0,
		StatusMsg:  "",
	}
	failureResponse := utils.Response{
		StatusCode: 1,
		StatusMsg:  "",
	}
	auth := auth2.Auth{}.GetAuth(c)
	//TODO
	// 业务代码
	// 1、解析参数
	toUserID := c.Query("to_user_id")
	actionType := c.Query("action_type")
	// 2、验证参数
	// 2.1 toUserID必须是整数
	id, err := strconv.Atoi(toUserID)
	if err != nil || id <= 0 {
		failureResponse.StatusMsg = "to_user_id必须为非0正整数"
		c.JSON(409, failureResponse)
		return
	}
	// 2.2 actionType 必须为1或2
	action, err := strconv.Atoi(actionType)
	if err != nil || (action != 1 && action != 2) {
		failureResponse.StatusMsg = "action必须为1或2"
		c.JSON(409, failureResponse)
		return
	}
	if auth.UserID == uint(id) {
		failureResponse.StatusMsg = "没办法关注自己或取消关注自己"
		c.JSON(409, failureResponse)
		return
	}
	// 3、查询数据库获取两个用户信息，使用for update加锁（用户一般都存在）
	err = service.RelationActionPost(uint(id), action, &auth)
	if err != nil {
		c.JSON(409, failureResponse)
		return
	}
	c.JSON(http.StatusOK, successResponse)
}

// RelationFollowList
// @Summary 获取关注列表
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object RelationList 成功后返回值
// @Failure 409 object RelationList 失败后返回值
// @Router /douyin/relation/follow/list/ [get]
func RelationFollowList(c *gin.Context) {

	auth := auth2.Auth{}.GetAuth(c)
	successResponse := RelationList{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	failureResponse := RelationList{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	//1、解析参数
	userID := c.Query("user_id")
	// 2、验证参数（确保user_id为正整数）
	id, err := strconv.Atoi(userID)
	if err != nil || id < 0 {
		failureResponse.StatusMsg = "user_id必须为正整数"
		c.JSON(http.StatusOK, failureResponse)
		return
	}
	if id == 0 {
		id = int(auth.UserID)
	}
	// 3、查询数据库获取该用的所有关注
	userList, err := service.RelationFollowListGet(uint(id), &auth)
	if err != nil {
		c.JSON(409, failureResponse)
		return
	}
	// 返回
	successResponse.UserList = userList
	c.JSON(http.StatusOK, successResponse)
}

// RelationFollowerList
// @Summary 获取关注者列表
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object RelationList 成功后返回值
// @Failure 409 object RelationList 失败后返回值
// @Router /douyin/relation/follower/list/ [get]
func RelationFollowerList(c *gin.Context) {

	auth := auth2.Auth{}.GetAuth(c)
	successResponse := RelationList{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	failureResponse := RelationList{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	//1、解析参数
	userID := c.Query("user_id")
	// 2、验证参数（确保user_id为正整数）
	id, err := strconv.Atoi(userID)
	if err != nil || id < 0 {
		failureResponse.StatusMsg = "user_id必须为正整数"
		c.JSON(http.StatusOK, failureResponse)
		return
	}
	if id == 0 {
		id = int(auth.UserID)
	}
	userList, err := service.RelationFollowerListGet(uint(id), &auth)
	if err != nil {
		c.JSON(409, failureResponse)
		return
	}
	// 返回
	successResponse.UserList = userList
	c.JSON(http.StatusOK, successResponse)
}

// RelationFriendList
// @Summary 获取好友列表
// @Tags 关系
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object RelationList 成功后返回值
// @Failure 409 object RelationList 失败后返回值
// @Router /douyin/relation/friend/list/ [get]
func RelationFriendList(c *gin.Context) {

	auth := auth2.Auth{}.GetAuth(c)
	successResponse := RelationList{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	failureResponse := RelationList{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		UserList: nil,
	}
	//1、解析参数
	userID := c.Query("user_id")
	// 2、验证参数（确保user_id为正整数）
	id, err := strconv.Atoi(userID)
	if err != nil || id <= 0 {
		failureResponse.StatusMsg = "user_id必须为正整数"
		c.JSON(http.StatusOK, failureResponse)
		return
	}
	// 3、查询数据库获取该用的所有关注
	userList, err := service.RelationFriendListGet(uint(id), &auth)
	if err != nil {
		c.JSON(409, failureResponse)
		return
	}
	// 返回
	successResponse.UserList = userList
	c.JSON(http.StatusOK, successResponse)
}
