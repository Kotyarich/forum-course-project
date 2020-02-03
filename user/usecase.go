package user

import (
	"context"
	"dbProject/models"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUp(ctx context.Context, user *models.User) ([]*models.User, error)
	SignIn(ctx context.Context, username, password string) (string, error)
	GetProfile(ctx context.Context, username string) (*models.User, error)
	ChangeProfile(ctx context.Context, user *models.User) (*models.User, error)
	ParseToken(ctx context.Context, accessToken string) (*models.User, error)
}
