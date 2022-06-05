package forum

import (
	"context"
	"thread-service/models"
)

type RepositoryThread interface {
	GetThreadBySlug(ctx context.Context, slug string) (*models.Thread, error)
	GetThreadById(ctx context.Context, id int) (*models.Thread, error)
	GetForumThreads(ctx context.Context, slug, since string, limit, offset int, sort bool) ([]*models.Thread, error)
	CreateThread(ctx context.Context, slug string, thread *models.Thread) (*models.Thread, error)
	DeleteThread(ctx context.Context, id int) error
	ChangeThread(ctx context.Context, slug, title, message string) (*models.Thread, error)
	VoteForThread(ctx context.Context, slug string, vote *models.Vote) (*models.Thread, error)
}
