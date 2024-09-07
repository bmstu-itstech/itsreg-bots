package jwtauth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type accessTokenClaims struct {
	jwt.RegisteredClaims
	UserUUID string
}

type AccessTokenPayload struct {
	UserUUID string
}

func NewAccessToken(
	userUUID string,
	ttl time.Duration,
) (string, error) {
	claims := accessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
		UserUUID: userUUID,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func ParseAccessToken(token string) (AccessTokenPayload, error) {
	parsed, err := jwt.ParseWithClaims(token, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return AccessTokenPayload{}, nil
	}

	claims, ok := parsed.Claims.(*accessTokenClaims)
	if !ok {
		return AccessTokenPayload{}, errors.New("can not cast to access token claims")
	}

	return AccessTokenPayload{
		UserUUID: claims.UserUUID,
	}, nil
}
