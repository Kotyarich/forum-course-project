package auth

import (
	"context"
	"user-service/models"
)

type UseCase interface {
	SignUp(ctx context.Context, user *models.User) ([]*models.User, string, error)
	SignIn(ctx context.Context, username, password string) (*models.User, string, error)
	CheckAuth(ctx context.Context, token string) (*models.User, error)
}
