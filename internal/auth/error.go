package auth

import "errors"

var (
	ErrNoContextValue      = errors.New("Expected context value wasn't found, user wasn't authorized or critical: middleware wasn't set")
	ErrWrongContextType    = errors.New("CRITICAL: Token extractor middleware use different token type")
	ErrExpiredToken        = errors.New("Token expired")
	ErrWrongTokenExtractor = errors.New("CRITICAL: generated token type and expected one are different")
)
