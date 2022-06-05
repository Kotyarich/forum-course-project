package forum

import (
	"context"
	"thread-service/models"
)

type UseCaseThread interface {
	GetThread(ctx context.Context, slug string) (*models.Thread, error)
	ChangeThread(ctx context.Context, slug, title, message string) (*models.Thread, error)
	CreateForumThread(ctx context.Context, slug string, thread *models.Thread) (*models.Thread, error)
	VoteForThread(ctx context.Context, slug string, vote models.Vote) (*models.Thread, error)
	GetForumThreads(ctx context.Context, slug, since string, limit, offset int, sort bool) ([]*models.Thread, error)
}

type Producer interface {
	ProduceNewThread(string)
	ProduceNewVote(string)
}

type UseCase struct {
	ThreadUseCase  UseCaseThread
}
