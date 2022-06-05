package forum

import (
	"context"
	"post-service/models"
)

type RepositoryPost interface {
	GetThreadPostsFlat(ctx context.Context, threadId int, limit, offset, since int, desc bool) ([]*models.Post, error)
	GetThreadPostsTree(ctx context.Context, threadId int, limit, offset, since int, desc bool) ([]*models.Post, error)
	GetThreadPostsParentTree(ctx context.Context, threadId int, limit, offset, since int, desc bool) ([]*models.Post, error)
	ThreadPostCreate(ctx context.Context, threadId int, posts []*models.Post) ([]*models.Post, error)
	GetPost(ctx context.Context, id int) (*models.Post, error)
	DeletePost(ctx context.Context, id int) error
	ChangePost(ctx context.Context, newMessage string, post *models.Post) error
}
