package model

import "github.com/golang-jwt/jwt/v4"

type UserClaims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}
