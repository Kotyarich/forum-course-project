package statistic

import (
	"context"
	"statistic-service/models"
)

type Repository interface {
	CreateUserRecord(ctx context.Context, id string) error
	CreatePostRecord(ctx context.Context, id int) error
	CreateVoteRecord(ctx context.Context, id string) error
	CreateThreadRecord(ctx context.Context, id int) error
	CreateForumRecord(ctx context.Context, id int) error
	GetStatistic(ctx context.Context) (*models.Status, error)
}
