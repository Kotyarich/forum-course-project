package usecase

import (
	"context"
	"statistic-service/models"
	statisticPkg "statistic-service/statistic"
)

type StatisticUseCase struct {
	statisticRepo  statisticPkg.Repository
}

func NewStatisticUseCase(userRepo statisticPkg.Repository) *StatisticUseCase {
	return &StatisticUseCase{
		statisticRepo:  userRepo,
	}
}

func (u *StatisticUseCase) GetStatistic(ctx context.Context) (*models.Status, error) {
	statistic, err := u.statisticRepo.GetStatistic(ctx)
	if err != nil {
		return nil, err
	}

	return statistic, nil
}
