package main

import (
	"context"
	"douyin_rpc/server/cmd/message/global"
	message "douyin_rpc/server/cmd/message/kitex_gen/message"
	"douyin_rpc/server/cmd/message/model"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct{}

// MessageAction implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageAction(ctx context.Context, req *message.MessageActionRequest) (resp *message.MessageActionResponse, err error) {
	message1 := model.Message{
		Model:      gorm.Model{},
		Content:    req.Content,
		UserID:     req.AuthId,
		TargetID:   req.ToUserId,
		CreateTime: time.Now().Unix(),
	}
	resp = new(message.MessageActionResponse)
	//TODO
	// 布隆过滤器过滤ID是否存在
	// 现在查询数据库
	tx := global.DB.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if req.ActionType == 1 {
		err = model.CreateMessage(tx, &message1)
		if err != nil {
			return
		}
	} else {
		// 删除消息未实现
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return
	}
	return
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
	resp = new(message.MessageListResponse)
	// 4、构建返回值
	var wg sync.WaitGroup
	wg.Add(len(messages))
	resp.MessageList = make([]*message.Message, len(messages))
	for j := 0; j < len(messages); j++ {
		i := j
		go func() {
			defer wg.Done()
			resp.MessageList[i] = &message.Message{
				Id:         strconv.FormatInt(messages[i].ID, 10),
				Content:    messages[i].Content,
				CreateTime: messages[i].CreatedAt.Format("2006-01-02 15:04:05"),
				FromUserId: strconv.FormatInt(messages[i].UserID, 10),
				ToUserId:   strconv.FormatInt(messages[i].TargetID, 10),
			}
		}()

	}
	wg.Wait()
	return
}

// GetMessage implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) GetMessage(ctx context.Context, req *message.GetMessageRequest) (resp *message.GetMessageResponse, err error) {
	/*
	    query_type=1  根据id查询
	    query_type=2  根据user_id查询
	    query_type=3  根据target_id查询
	    query_type=4  根据user_id和target_id查询
	   如果给定了根据create_time也需要根据这个查询，默认-1不给定
	*/
	tx := global.DB.Debug()
	//TODO
	// 有关消息时间设置
	//if req.CreateTime!=-1{
	//	after := time.Unix(req.CreateTime, 0)
	//}
	var cache []model.Message
	if req.QueryType == 1 {
		cache, err = model.QueryMessageByUserIDAndTargetIDWithCache(tx, req.UserId, req.TargetId)
		if err != nil {
			return
		}
		resp.Messages = make([]*message.Message1, len(cache))
		for j := 0; j < len(cache); j++ {
			i := j
			go func() {
				resp.Messages[i] = &message.Message1{
					Id:         cache[i].ID,
					UserId:     cache[i].UserID,
					TargetId:   cache[i].TargetID,
					Content:    cache[i].Content,
					CreateTime: cache[i].CreateTime,
					CreatedAt:  cache[i].CreatedAt.Unix(),
					UpdatedAt:  cache[i].UpdatedAt.Unix(),
				}
			}()
		}
	} else if req.QueryType == 2 {
		cache, err = model.QueryMessageByUserIDAndTargetIDWithCache(tx, req.UserId, req.TargetId)
		if err != nil {
			return
		}
	} else if req.QueryType == 3 {
		cache, err = model.QueryMessageByUserIDAndTargetIDWithCache(tx, req.UserId, req.TargetId)
		if err != nil {
			return
		}
	} else if req.QueryType == 4 {
		cache, err = model.QueryMessageByUserIDAndTargetIDWithCache(tx, req.UserId, req.TargetId)
		if err != nil {
			return
		}
	}
	resp.Messages = make([]*message.Message1, len(cache))
	for j := 0; j < len(cache); j++ {
		i := j
		go func() {
			resp.Messages[i] = &message.Message1{
				Id:         cache[i].ID,
				UserId:     cache[i].UserID,
				TargetId:   cache[i].TargetID,
				Content:    cache[i].Content,
				CreateTime: cache[i].CreateTime,
				CreatedAt:  cache[i].CreatedAt.Unix(),
				UpdatedAt:  cache[i].UpdatedAt.Unix(),
			}
		}()
	}
	return
}
