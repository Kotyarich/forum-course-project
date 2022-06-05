package forum

import (
	"context"
	"forum-service/models"
)

type RepositoryForum interface {
	CreateForum(ctx context.Context, forum *models.Forum) (*models.Forum, error)
	GetForum(ctx context.Context, slug string) (*models.Forum, error)
	DeleteForum(ctx context.Context, slug string) error
	GetForums(ctx context.Context) ([]*models.Forum, error)
	GetForumUsers(ctx context.Context, slug, since string, limit int, sort bool) ([]*models.User, error)
}
