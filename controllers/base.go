package controllers

import (
	"github.com/golang-jwt/jwt/v4"
)

type Auth struct {
	UserID        uint
	UserName      string
	FollowCount   int
	FollowerCount int
	jwt.RegisteredClaims
}
