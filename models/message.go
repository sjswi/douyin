package models

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	gorm.Model
	Content    string
	UserID     uint
	TargetID   uint
	CreateTime time.Time
}
