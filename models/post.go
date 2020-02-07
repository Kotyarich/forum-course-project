package models

import "time"

type Post struct {
	Author    string
	Created   time.Time
	ForumName string
	Id        int
	IsEdited  bool
	Message   string
	Parent    int
	Tid       int
}

type DetailedInfo struct {
	PostInfo   Post
	AuthorInfo *User
	ThreadInfo *Thread
	ForumInfo  *Forum
}
