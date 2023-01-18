package service

import (
	auth2 "douyin/auth"
	"douyin/bootstrap/driver"
	"douyin/models"
	"douyin/vo"
	"gorm.io/gorm"
	"time"
)

type MessageChatFlow struct {
	ToUserId    uint
	AuthUser    *auth2.Auth
	MessageList []vo.Message `json:"message_list"`
}

func MessageChatGet(toUserId uint, auth *auth2.Auth) ([]vo.Message, error) {
	return (&MessageChatFlow{
		MessageList: nil,
		ToUserId:    toUserId,
		AuthUser:    auth,
	}).Do()
}
func (c *MessageChatFlow) Do() ([]vo.Message, error) {
	if err := c.checkParam(); err != nil {
		return nil, err
	}
	if err := c.messageChat(); err != nil {
		return nil, err
	}
	return c.MessageList, nil
}
func (c *MessageChatFlow) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}

func (c *MessageChatFlow) messageChat() error {
	tx := driver.Db.Debug()
	messages1, err := models.QueryMessageByUserIDAndTargetIDWithCache(tx, c.AuthUser.UserID, c.ToUserId)
	if err != nil {
		return err
	}
	messages2, err := models.QueryMessageByUserIDAndTargetIDWithCache(tx, c.ToUserId, c.AuthUser.UserID)
	if err != nil {
		return err
	}
	messages := append(messages1, messages2...)
	// 4、构建返回值
	c.MessageList = make([]vo.Message, len(messages))
	for j := 0; j < len(messages); j++ {
		i := j
		go func() {
			c.MessageList[i].ID = messages[i].ID
			c.MessageList[i].Content = messages[i].Content
			c.MessageList[i].CreateTime = messages[i].CreateTime.Format("2006-01-02 15:04:05")
		}()

	}
	return nil
}

type MessageActionFlow struct {
	ToUserId   uint
	ActionType int
	Content    string
	AuthUser   *auth2.Auth
}

func MessageActionPost(toUserId uint, actionType int, content string, auth *auth2.Auth) error {
	return (&MessageActionFlow{
		ToUserId:   toUserId,
		ActionType: actionType,
		Content:    content,
		AuthUser:   auth,
	}).Do()
}
func (c *MessageActionFlow) Do() error {
	if err := c.checkParam(); err != nil {
		return err
	}
	if err := c.message(); err != nil {
		return err
	}
	return nil
}
func (c *MessageActionFlow) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}

func (c *MessageActionFlow) message() error {
	message := models.Message{
		Model:      gorm.Model{},
		Content:    c.Content,
		UserID:     c.AuthUser.UserID,
		TargetID:   c.ToUserId,
		CreateTime: time.Now().UTC(),
	}
	//TODO
	// 布隆过滤器过滤ID是否存在
	// 现在查询数据库
	tx := driver.Db.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if c.ActionType == 1 {
		err := models.CreateMessage(tx, message)
		if err != nil {
			return err
		}
	} else {
		// 删除消息未实现
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
