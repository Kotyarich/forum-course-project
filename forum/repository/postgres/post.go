package postgres

import (
	"dbProject/models"
	"time"
)

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

func toModelPost(p *Post) *models.Post {
	return &models.Post{
		Author:    p.Author,
		Created:   p.Created,
		ForumName: p.ForumName,
		Id:        p.Id,
		IsEdited:  p.IsEdited,
		Message:   p.Message,
		Parent:    p.Parent,
		Tid:       p.Tid,
	}
}

func modelToPost(p *models.Post) *Post {
	return &Post{
		Author:    p.Author,
		Created:   p.Created,
		ForumName: p.ForumName,
		Id:        p.Id,
		IsEdited:  p.IsEdited,
		Message:   p.Message,
		Parent:    p.Parent,
		Tid:       p.Tid,
	}
}
