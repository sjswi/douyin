package controllers

import (
	"douyin/bootstrap/driver"
	"douyin/models"
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	ID         uint
	Content    string
	CreateTime string
}
type MessageChatResponse struct {
	utils.Response
	MessageList []Message `json:"message_list"`
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
	auth := Auth{}.GetAuth(c)

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
	userID := c.Query("user_id")
	// 2、验证参数
	id, err := strconv.Atoi(userID)
	if err != nil || id <= 0 {
		failureResponse.StatusMsg = "user_id必须为正整数"
		c.JSON(409, failureResponse)
		return
	}
	if id == 0 {
		id = int(auth.UserID)
	}
	// 3、获取消息
	var messageList []models.Message
	if err := driver.Db.Debug().Model(models.Message{}).Where("user_id=? and target_id=?", auth.UserID, id).Or("user_id=? and target_id=?", id, auth.UserID).Order("create_time DESC").Find(&messageList).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库错误" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	// 4、构建返回值
	returnMessage := make([]Message, len(messageList))
	for i := 0; i < len(messageList); i++ {
		returnMessage[i].ID = messageList[i].ID
		returnMessage[i].Content = messageList[i].Content
		returnMessage[i].CreateTime = messageList[i].CreateTime.Format("2006-01-02 15:04:05")
	}
	successResponse.MessageList = returnMessage
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
// @Router /douyin/message/action [post]
func MessageAction(c *gin.Context) {
	auth := Auth{}.GetAuth(c)
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
	// 3、构建消息
	message := models.Message{
		Model:      gorm.Model{},
		Content:    content,
		UserID:     auth.UserID,
		TargetID:   uint(id),
		CreateTime: time.Now().UTC(),
	}
	//TODO
	// 布隆过滤器过滤ID是否存在
	// 现在查询数据库
	var user models.User
	if err := driver.Db.Debug().Model(user).Where("id = ?", id).Find(&user).Error; err != nil {
		failureResponse.StatusMsg = "查询用户信息错误" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	if user.ID == 0 {
		failureResponse.StatusMsg = "用户不存在"
		c.JSON(409, failureResponse)
		return
	}
	// 4、创建消息
	tx := driver.Db.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Model(message).Create(&message).Error; err != nil {
		failureResponse.StatusMsg = "创建消息失败"
		c.JSON(409, failureResponse)
		return
	}
	if err := tx.Commit().Error; err != nil {
		failureResponse.StatusMsg = "提交失败"
		c.JSON(409, failureResponse)
		return
	}
	c.JSON(200, successResponse)
	return
}
