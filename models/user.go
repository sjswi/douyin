package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name          string
	FollowCount   int
	FollowerCount int
	Password      string
	Salt          string
}
