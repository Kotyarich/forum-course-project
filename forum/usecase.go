package forum

import (
	"context"
	"dbProject/models"
)

type UseCase interface {
	CreateForum(ctx context.Context, forum *models.Forum) (*models.Forum, error)
	CreateForumThread(ctx context.Context, slug string, thread *models.Thread) (*models.Thread, error)
	GetForumDetails(ctx context.Context, slug string) (*models.Forum, error)
	GetForumThreads(ctx context.Context, slug, since string, limit int, sort bool) ([]*models.Thread, error)
	GetForumUsers(ctx context.Context, slug, since string, limit int, sort bool) ([]*models.User, error)
}