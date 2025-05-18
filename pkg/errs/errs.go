package errs

import "errors"

var (
	ErrTokenExpired          = errors.New("token expired")
	ErrInvalidToken          = errors.New("invalid token")
	ErrFailedToParseClaims   = errors.New("failed to parse claims")
	ErrInvalidUserID         = errors.New("invalid user ID")
	ErrFailedToParseUserID   = errors.New("failed to parse user ID")
	ErrFailedToParseUserRole = errors.New("failed to parse user role")
	ErrInvalidUserRole       = errors.New("invalid user role")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrInvalidLogPathFile    = errors.New("invalid log file path")
)
