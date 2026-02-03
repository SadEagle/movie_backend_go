package auth

import "errors"

var (
	ErrNoContextValue      = errors.New("CRITICAL: Expected context value wasn't found")
	ErrWrongContextType    = errors.New("CRITICAL: Context value expected different type")
	ErrExpiredToken        = errors.New("Token expired")
	ErrWrongTokenExtractor = errors.New("Wrong token extractor")
)
