package models

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID  uint //用户id
	VideoID uint //视频id
	Exist   bool //是否存在，避免重复的点赞取消点赞使得该表内容变得很大
}
