package models

type User struct {
	Nickname string
	Password string
	Fullname string
	Email    string
	About    string
	IsAdmin  bool
}
