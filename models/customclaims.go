package models

import "github.com/golang-jwt/jwt/v4"

type CustomClaims struct {
	Email     string `json:"email"`
	Verified  bool   `json:"verfied"`
	Privilege int8   `json:"privilege"`
	jwt.RegisteredClaims
}
