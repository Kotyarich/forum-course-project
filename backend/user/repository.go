package user

import (
	"context"
	"dbProject/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) ([]*models.User, int, error)
	AuthUser(ctx context.Context, username, password string) (*models.User, int, error)
	CreateSession(ctx context.Context, userId int) (string, error)
	CheckSession(ctx context.Context, token string) (*models.User, error)
	DeleteSession(ctx context.Context, token string) error
	GetUser(ctx context.Context, username string) (*models.User, error)
	ChangeUser(ctx context.Context, user *models.User) (*models.User, error)
}
