package forum

import (
	"context"
	"post-service/models"
)

type UseCasePost interface {
	CreateThreadPost(ctx context.Context, threadId int, posts []*models.Post) ([]*models.Post, error)
	GetThreadPosts(ctx context.Context, threadId int, limit, offset, since int, desc bool, sort models.PostSortType) ([]*models.Post, error)
	ChangePost(ctx context.Context, id int, message string) (*models.Post, error)
	// TODO add removing
}

type UseCase struct {
	PostUseCase    UseCasePost
}
