package auth

// TODO: Implement a custom error type for easier and more precise error handling

import "errors"

var (
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionUnverified  = errors.New("session not verified")
	ErrSessionExpired     = errors.New("session is expired")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
