package forum

import (
	"context"
	"forum-service/models"
)

type UseCaseForum interface {
	CreateForum(ctx context.Context, forum *models.Forum) (*models.Forum, error)
	GetForums(ctx context.Context) ([]*models.Forum, error)
	GetForumDetails(ctx context.Context, slug string) (*models.Forum, error)
	GetForumUsers(ctx context.Context, slug, since string, limit int, sort bool) ([]*models.User, error)
	// TODO add removing
	// TODO add changing
}

type UseCase struct {
	ForumUseCase   UseCaseForum
}
