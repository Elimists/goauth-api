package controller

import (
	"github.com/Elimists/go-app/models"
	"github.com/golang-jwt/jwt/v4"
)

type Oauth struct {
	Error   bool
	Code    string
	Message string
	Claim   *models.CustomClaims `json:"-"`
}

func Checkauth(cookie string) Oauth {

	token, err := jwt.ParseWithClaims(cookie, &models.CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return SECRET_KEY, nil
		})

	if err != nil {
		return Oauth{Error: true, Code: "parsing_error", Message: "Could not parse token", Claim: nil}
	}

	if !token.Valid {
		return Oauth{Error: true, Code: "invalid_token", Message: "The token is not valid", Claim: nil}
	}

	claims, ok := token.Claims.(*models.CustomClaims)

	if !ok {
		return Oauth{Error: true, Code: "invalid_claims", Message: "Un..claimable", Claim: nil}
	}

	return Oauth{Error: false, Code: "validated", Message: "Success", Claim: claims}
}
