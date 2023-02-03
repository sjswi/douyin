package main

import (
	"context"
	comment "douyin_rpc/server/cmd/comment/kitex_gen/comment"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CommentAction implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentAction(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	// TODO: Your code here...
	return
}

// CommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetComment(ctx context.Context, req *comment.GetCommentRequest) (resp *comment.GetCommentResponse, err error) {
	// TODO: Your code here...
	return
}
