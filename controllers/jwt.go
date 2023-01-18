package controllers

import (
	auth2 "douyin/auth"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var MySecret = []byte("手写的从前")

func MakeToken(auth *auth2.Auth) (tokenString string, err error) {
	claim := auth2.Auth{
		UserName:      auth.UserName,
		UserID:        auth.UserID,
		FollowCount:   auth.FollowCount,
		FollowerCount: auth.FollowerCount,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))), // 过期时间3小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                       // 生效时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
	tokenString, err = token.SignedString(MySecret)
	return tokenString, err
}

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte("手写的从前"), nil // 这是我的secret
	}
}

func ParseToken(originToken string) (*auth2.Auth, error) {
	token, err := jwt.ParseWithClaims(originToken, &auth2.Auth{}, Secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
	}
	if claims, ok := token.Claims.(*auth2.Auth); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}
