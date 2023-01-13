package models

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	gorm.Model
	Content    string
	UserID     uint
	TargetId   uint
	CreateTime time.Time
}
