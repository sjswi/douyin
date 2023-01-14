package models

import "gorm.io/gorm"

type Relation struct {
	gorm.Model
	UserID   uint
	TargetID uint
	Type     int  // 1:UserID关注TargetID   2:互相关注
	Exist    bool // 判断是否存在，避免一个用户频繁关注取消关注造成数据膨胀，并且这样可以使用唯一索引
}
