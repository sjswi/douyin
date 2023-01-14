package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	UserID        uint
	UserName      string
	FollowCount   int
	FollowerCount int
	jwt.RegisteredClaims
}

// GetAuth
// 从context获取auth信息
func (a Auth) GetAuth(c *gin.Context) Auth {
	auth, exists := c.Get("auth")
	if !exists {
		auth = Auth{
			UserID:           0,
			UserName:         "",
			FollowCount:      0,
			FollowerCount:    0,
			RegisteredClaims: jwt.RegisteredClaims{},
		}
	}
	return auth.(Auth)
}
