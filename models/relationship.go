package models

import "gorm.io/gorm"

type Relation struct {
	gorm.Model
	UserID   uint
	TargetID uint
	Type     int // 1:UserID关注TargetID， 2:TargetID关注UserID  3:互相关注

}
