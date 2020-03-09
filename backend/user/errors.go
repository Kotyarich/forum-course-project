package user

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user with this email or nickname exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidSession    = errors.New("invalid session")
)
