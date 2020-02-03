package postgres

import (
	"dbProject/db"
	"dbProject/models"
	"github.com/jackc/pgx"
	"time"
)

type ThreadRepository struct {
	db *pgx.ConnPool
}

func NewThreadRepository() *ForumRepository {
	return &ForumRepository{db.GetDB()}
}

type Thread struct {
	Author    string
	Slug      *string
	Votes     int
	Title     string
	Created   time.Time
	ForumName string
	Id        int
	Message   string
}

func toPostgresThread(thread *models.Thread) *Thread {
	return &Thread{
		Id:        thread.Id,
		Author:    thread.Author,
		Slug:      thread.Slug,
		Votes:     thread.Votes,
		Title:     thread.Title,
		Created:   thread.Created,
		ForumName: thread.ForumName,
		Message:   thread.Message,
	}
}

func toModelThread(thread *Thread) *models.Thread {
	return &models.Thread{
		Id:        thread.Id,
		Author:    thread.Author,
		Slug:      thread.Slug,
		Votes:     thread.Votes,
		Title:     thread.Title,
		Created:   thread.Created,
		ForumName: thread.ForumName,
		Message:   thread.Message,
	}
}
