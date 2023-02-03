package main

import (
	"context"
	video "douyin_rpc/server/cmd/video/kitex_gen/video"
)

// FeedServiceImpl implements the last service interface defined in the IDL.
type FeedServiceImpl struct{}

// Feed implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	// TODO: Your code here...
	return
}

// PublishAction implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {
	// TODO: Your code here...
	return
}

// PublishList implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetVideo implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) GetVideo(ctx context.Context, req *video.GetVideoRequest) (resp *video.GetVideoResponse, err error) {
	// TODO: Your code here...
	return
}
