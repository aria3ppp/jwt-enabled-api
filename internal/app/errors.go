package app

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrUsedEmail         = errors.New("email used")
	ErrIncorrectPassword = errors.New("incorrect password")
	// ErrExpiredToken      = errors.New("expired token")
	// ErrInvalidToken      = errors.New("invalid token")
	ErrSamePassword = errors.New("same password")
)
