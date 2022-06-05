package user

import (
	"context"
	"user-service/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) ([]*models.User, int, error)
	AuthUser(ctx context.Context, username, password string) (*models.User, int, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
	ChangeUser(ctx context.Context, user *models.User) (*models.User, error)
}
