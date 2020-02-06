package models

type Vote struct {
	Voice    Voice  `json:"voice"`
	Nickname string `json:"nickname"`
}

type Voice int

const (
	Up   Voice = 1
	Down Voice = -1
)
