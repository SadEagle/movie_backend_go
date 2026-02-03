package auth

import (
	"context"
	"net/http"
)

func TokenExtractionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(rw, "No Authorization header", http.StatusBadRequest)
			return
		}

		userTokenData, err := TokenExtract(tokenStr)
		if err != nil {
			http.Error(rw, "Can't extract token data", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "tokenData", userTokenData)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

func GetTokenDataContext(ctx context.Context) (UserTokenData, error) {
	// Extract token data from request context after middleware call
	userTokenDataAny := ctx.Value("userTokenData")
	if userTokenDataAny == nil {
		return UserTokenData{}, ErrNoContextValue
	}
	userTokenData, ok := userTokenDataAny.(*UserTokenData)
	if !ok {
		return UserTokenData{}, ErrWrongContextType
	}
	return *userTokenData, nil
}
