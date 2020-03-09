package user

import (
	"context"
	"dbProject/models"
)

type UseCase interface {
	SignUp(ctx context.Context, user *models.User) ([]*models.User, string, error)
	SignIn(ctx context.Context, username, password string) (*models.User, string, error)
	SignOut(ctx context.Context, token string) error
	GetProfile(ctx context.Context, username string) (*models.User, error)
	ChangeProfile(ctx context.Context, user *models.User) (*models.User, error)
	CheckAuth(ctx context.Context, token string) (*models.User, error)
}
