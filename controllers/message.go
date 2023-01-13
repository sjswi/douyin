package controllers

import (
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Message struct {
	ID         uint
	Content    string
	CreateTime time.Time
}
type MessageChatResponse struct {
	utils.Response
	MessageList []*Message `json:"message_list"`
}

// MessageChat
// @Summary 获取聊天记录
// @Tags 消息
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object FavoriteListResponse 成功后返回值
// @Failure 409 object FavoriteListResponse 失败后返回值
// @Router /douyin/message/chat [get]
func MessageChat(c *gin.Context) {
	c.Query("user_id")
	//TODO
	// 业务代码

	response := MessageChatResponse{
		Response:    utils.Response{},
		MessageList: nil,
	}

	c.JSON(http.StatusOK, response)
	return
}

type PostMessageActionForm struct {
	Token      string `json:"token"`
	ToUserID   uint   `json:"to_user_id"`
	ActionType int    `json:"action_type"`
	Content    string `json:"content"`
}

// MessageAction
// @Summary 发送消息
// @Tags 消息
// @version 1.0
// @Accept application/x-json-stream
// @Param messageAction body PostMessageActionForm true "消息"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/message/action [post]
func MessageAction(c *gin.Context) {
	var messageActionForm *PostMessageActionForm
	if err := c.ShouldBind(&messageActionForm); err != nil {
		response := utils.Response{
			StatusCode: 1,
			StatusMsg:  "data数据有误",
		}

		c.JSON(409, response)
		return
	}
	///TODO
	// 业务代码
	response := utils.Response{
		StatusCode: 0,
		StatusMsg:  "data数据有误",
	}

	c.JSON(409, response)
	return
}
