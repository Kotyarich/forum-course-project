package forum

import "errors"

var (
	ErrForumNotFound       = errors.New("forum not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrForumAlreadyExists  = errors.New("forum already exists")
	ErrThreadAlreadyExists = errors.New("thread already exists")
)
