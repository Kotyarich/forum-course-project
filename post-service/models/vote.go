package models

type Vote struct {
	Voice    Voice
	Nickname string
}

type Voice int

const (
	Up   Voice = 1
	Down Voice = -1
)
