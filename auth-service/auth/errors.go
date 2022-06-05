package auth

import "errors"

var (
	ErrUserAlreadyExists = errors.New("auth with this email or nickname exists")
	ErrUserNotFound      = errors.New("auth not found")
	ErrInvalidSession    = errors.New("invalid session")
)
