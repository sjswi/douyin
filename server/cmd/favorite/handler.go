package main

import (
	"context"
	favorite "douyin_rpc/server/cmd/favorite/kitex_gen/favorite"
)

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	// TODO: Your code here...
	return
}

// FavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavorite(ctx context.Context, req *favorite.GetFavoriteRequest) (resp *favorite.GetFavoriteResponse, err error) {
	// TODO: Your code here...
	return
}

// GetFavoriteCount implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteCount(ctx context.Context, req *favorite.GetFavoriteCountRequest) (resp *favorite.GetFavoriteCountResponse, err error) {
	// TODO: Your code here...
	return
}
