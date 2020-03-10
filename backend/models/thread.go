package models

import "time"

type Thread struct {
	Author     string
	Slug       *string
	Votes      int
	Title      string
	Created    time.Time
	ForumName  string
	Id         int
	Message    string
	PostsCount int
}

type PostSortType int

const (
	Flat       PostSortType = 0
	Tree       PostSortType = 1
	ParentTree PostSortType = 2
)
