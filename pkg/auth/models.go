package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserTokenData struct {
	UserID  pgtype.UUID `json:"user_id"`
	IsAdmin bool        `json:"is_admin"`
}

type UserClaims struct {
	UserTokenData
	jwt.RegisteredClaims
}
