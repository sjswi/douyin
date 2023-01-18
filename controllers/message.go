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

type MessageChatResponse struct {
	utils.Response
	MessageList []vo.Message `json:"message_list"`
}

// MessageChat
// @Summary 获取聊天记录
// @Tags 消息
// @version 1.0
// @Accept application/x-json-stream
// @Param to_user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object FavoriteListResponse 成功后返回值
// @Failure 409 object FavoriteListResponse 失败后返回值
// @Router /douyin/message/chat/ [get]
func MessageChat(c *gin.Context) {
	auth := auth2.Auth{}.GetAuth(c)

	successResponse := MessageChatResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		MessageList: nil,
	}
	failureResponse := MessageChatResponse{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		MessageList: nil,
	}
	// 1、解析参数
	toUserID := c.Query("to_user_id")
	// 2、验证参数
	toUserId, err := strconv.Atoi(toUserID)
	if err != nil || toUserId <= 0 {
		failureResponse.StatusMsg = "user_id必须为正整数"
		c.JSON(409, failureResponse)
		return
	}
	// 3、获取消息
	messages, err := service.MessageChatGet(uint(toUserId), &auth)
	if err != nil {
		return
	}
	successResponse.MessageList = messages
	// 返回
	c.JSON(http.StatusOK, successResponse)
	return
}

// MessageAction
// @Summary 发送消息
// @Tags 消息
// @version 1.0
// @Accept application/x-json-stream
// @Param token query string true "token"
// @Param to_user_id query int true "对方用户id"
// @Param action_type query int true "类型"
// @Param content query string true "消息内容"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/message/action/ [post]
func MessageAction(c *gin.Context) {
	auth := auth2.Auth{}.GetAuth(c)
	successResponse := utils.Response{
		StatusCode: 0,
		StatusMsg:  "",
	}
	failureResponse := utils.Response{
		StatusCode: 1,
		StatusMsg:  "",
	}
	// 1、解析参数
	toUserID := c.Query("to_user_id")
	actionType := c.Query("action_type")
	content := c.Query("content")
	// 2、验证参数
	// toUserID必须为正整数
	id, err := strconv.Atoi(toUserID)
	if err != nil || id <= 0 {
		failureResponse.StatusMsg = "toUserID必须为正整数"
		c.JSON(409, failureResponse)
		return
	}
	action, err := strconv.Atoi(actionType)
	if err != nil || action != 1 {
		failureResponse.StatusMsg = "action必须为1"
		c.JSON(409, failureResponse)
		return
	}

	//TODO
	// 布隆过滤器过滤ID是否存在
	// 现在查询数据库
	err = service.MessageActionPost(uint(id), action, content, &auth)
	if err != nil {
		c.JSON(409, failureResponse)
		return
	}
	c.JSON(200, successResponse)
	return
}
