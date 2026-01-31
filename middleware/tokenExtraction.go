package middleware

// import (
// 	"errors"
// 	"movie_backend_go/secret"
// 	"net/http"
// 	"strings"
// )
//
// var (
// 	ErrNoAuthHeader      = errors.New("no authorization header")
// 	ErrInvalidAuthFormat = errors.New("invalid authorization format")
// )
//
// const bearerAuthPrefix = "Bearer "
//
// // JWT manager
// func extractBearerTokenHTTP(r *http.Request) (secret.UserTokenData, error) {
// 	authTokenStr := r.Header.Get("Authorization")
// 	if authTokenStr == "" {
// 		return secret.UserTokenData{}, ErrNoAuthHeader
// 	}
//
// 	if !strings.HasPrefix(authTokenStr, bearerAuthPrefix) {
// 		return secret.UserTokenData{}, ErrInvalidAuthFormat
// 	}
// 	tokenStr := strings.TrimPrefix(authTokenStr, bearerAuthPrefix)
// 	if tokenStr == "" {
// 		return secret.UserTokenData{}, ErrInvalidAuthFormat
// 	}
// 	userTokenData, err := secret.TokenExtract(tokenStr)
// 	return userTokenData, err
// }
//
// func MiddlewareAccessToken(next http.Handler) http.Handler {
// 	return
// }
//
// func MiddlewareUserByTokenData(next http.Handler) http.Handler {}
