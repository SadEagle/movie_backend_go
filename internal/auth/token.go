package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

var (
	JWTSignKey  = []byte("RandomKeyNeedToChangeLater")
	EXPIRE_TIME = 24 * time.Hour
)

func OauthTokenGenerate(userTokenData UserTokenData) (oauth2.Token, error) {
	experify := time.Now().Add(EXPIRE_TIME)
	expiresIn := int64(EXPIRE_TIME.Seconds())

	claims := UserClaims{
		UserTokenData: userTokenData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(experify),
		},
	}

	tokenData := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := tokenData.SignedString(JWTSignKey)
	if err != nil {
		return oauth2.Token{}, fmt.Errorf("sign jwt key with token data: %w", err)
	}
	oauthToken := oauth2.Token{AccessToken: accessToken, TokenType: "Bearer", Expiry: experify, ExpiresIn: expiresIn}

	return oauthToken, nil
}

func BearerTokenExtract(tokenStr string) (UserTokenData, error) {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(tokenStr, bearerPrefix) {
		return UserTokenData{}, ErrWrongTokenExtractor
	}
	tokenStr = strings.TrimPrefix(tokenStr, bearerPrefix)

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
