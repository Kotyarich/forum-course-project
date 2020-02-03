package user

import (
	"context"
	"dbProject/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) ([]*models.User, error)
	AuthUser(ctx context.Context, username, password string) (*models.User, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
	ChangeUser(ctx context.Context, user *models.User) (*models.User, error)
}
