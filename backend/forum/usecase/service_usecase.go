package usecase

import (
	"context"
	"dbProject/forum"
	"dbProject/models"
)

type ServiceUseCase struct {
	serviceRepo forum.RepositoryService
}

func (u *ServiceUseCase) Clear(ctx context.Context) error {
	return u.serviceRepo.Clear(ctx)
}

func (u *ServiceUseCase) Status(ctx context.Context) (*models.Status, error) {
	status, err := u.serviceRepo.Status(ctx)
	if err != nil {
		return nil, err
	}

	return status, nil
}
