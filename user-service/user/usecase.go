package user

import (
	"context"
	"user-service/models"
)


type Producer interface {
	Produce(string)
}

type UseCase interface {
	CreateUser(ctx context.Context, user *models.User) ([]*models.User, error)
	GetProfile(ctx context.Context, username string) (*models.User, error)
	ChangeProfile(ctx context.Context, user *models.User) (*models.User, error)
	CheckAuth(ctx context.Context, username string, password string) (*models.User, error)
}
