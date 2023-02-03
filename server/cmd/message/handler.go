package main

import (
	"context"
	message "douyin_rpc/server/cmd/message/kitex_gen/message"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct{}

// MessageAction implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageAction(ctx context.Context, req *message.MessageActionRequest) (resp *message.MessageActionResponse, err error) {
	// TODO: Your code here...
	return
}

// MessageList implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) MessageList(ctx context.Context, req *message.MessageListRequest) (resp *message.MessageListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetMessage implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) GetMessage(ctx context.Context, req *message.GetMessageRequest) (resp *message.GetMessageResponse, err error) {
	// TODO: Your code here...
	return
}
