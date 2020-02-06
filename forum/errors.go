package forum

import "errors"

var (
	ErrForumNotFound       = errors.New("forum not found")
	ErrThreadNotFound      = errors.New("thread not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrForumAlreadyExists  = errors.New("forum already exists")
	ErrThreadAlreadyExists = errors.New("thread already exists")
	ErrWrongParentsThread  = errors.New("post parent in another thread")
	ErrPostPatentNotFound  = errors.New("post parent not found")
)
