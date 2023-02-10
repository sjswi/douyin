package main

import (
	"context"
	"douyin_rpc/server/cmd/message/global"
	message "douyin_rpc/server/cmd/message/kitex_gen/message"
	"douyin_rpc/server/cmd/message/model"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct{}

// MessageAction implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageAction(ctx context.Context, req *message.MessageActionRequest) (resp *message.MessageActionResponse, err error) {
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

// MessageList implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageList(ctx context.Context, req *message.MessageListRequest) (resp *message.MessageListResponse, err error) {
	tx := global.DB.Debug()
	messages1, err := model.QueryMessageByUserIDAndTargetIDWithCache(tx, req.AuthId, req.ToUserId)
	if err != nil {
		return
	}
	messages2, err := model.QueryMessageByUserIDAndTargetIDWithCache(tx, req.ToUserId, req.AuthId)
	if err != nil {
		return
	}
	messages := append(messages1, messages2...)
	// 4、构建返回值
	resp.MessageList = make([]*message.Message, len(messages))
	for j := 0; j < len(messages); j++ {
		i := j
		go func() {
			resp.MessageList[i].Id = messages[i].ID

			resp.MessageList[i].Content = messages[i].Content
			resp.MessageList[i].CreateTime = messages[i].CreatedAt.Format("2006-01-02 15:04:05")
		}()

	}
	return
}

// GetMessage implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) GetMessage(ctx context.Context, req *message.GetMessageRequest) (resp *message.GetMessageResponse, err error) {
	if req.QueryType == 1 {

	} else if req.QueryType == 2 {

	} else if req.QueryType == 3 {

	}
}
