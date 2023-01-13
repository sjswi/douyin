package models

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	gorm.Model
	Content    string
	CreateTime time.Time
	VideoID    uint
	UserID     uint
}
