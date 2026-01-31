package secret

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var JWTSignKey = []byte("RandomKeyNeedToChangeLater")
var ErrExpiredToken = errors.New("Token expired")

type UserTokenData struct {
	UserID  pgtype.UUID `json:"user_id"`
	IsAdmin bool        `json:"-"`
}

type UserClaims struct {
	UserTokenData
	jwt.RegisteredClaims
}

func TokenGenerate(userTokenData UserTokenData, expiresAt time.Time) (string, error) {
	claims := UserClaims{
		UserTokenData: userTokenData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(JWTSignKey)
	if err != nil {
		return "", fmt.Errorf("sign jwt key with token data - %+v: %w", userTokenData, err)
	}
	return ss, nil
}

func TokenExtract(tokenStr string) (UserTokenData, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (any, error) {
		return JWTSignKey, nil
	})
	if err != nil {
		return UserTokenData{}, fmt.Errorf("parse JWT: %w", err)
	}
	claims := token.Claims.(*UserClaims)

	exp_time, err := token.Claims.GetExpirationTime()
	if err != nil {
		return UserTokenData{}, fmt.Errorf("get expiration time: %w", err)
	}
	if time.Now().After(exp_time.Time) {
		return UserTokenData{}, ErrExpiredToken
	}

	return claims.UserTokenData, nil
}
