package auth

import (
	"context"
	"user-service/models"
)

type UserService interface {
	CheckAuth(ctx context.Context, username string, password string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, []*models.User, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
}
