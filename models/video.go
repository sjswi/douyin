package models

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	AuthorID      uint   //作者id
	Title         string //视频标题
	CommentCount  int    //评论数
	FavoriteCount int    //点赞数
	PlayURL       string //播放地址
	CoverURL      string //封面地址
}
