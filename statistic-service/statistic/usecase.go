package statistic

import (
	"context"
	"statistic-service/models"
)

type UseCase interface {
	GetStatistic(ctx context.Context) (*models.Status, error)
}
