package usecase

import (
	"context"
	"time"
	"user-service/models"
	userPkg "user-service/user"
)

type UserUseCase struct {
	userRepo       userPkg.Repository
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
	producer       userPkg.Producer
}

func NewUserUseCase(
	userRepo userPkg.Repository,
	producer userPkg.Producer,
	hashSalt string,
	signingKey []byte,
	sessionDurationSeconds time.Duration) *UserUseCase {
	return &UserUseCase{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * sessionDurationSeconds,
		producer:       producer,
	}
}

func (u *UserUseCase) CreateUser(ctx context.Context, user *models.User) ([]*models.User, error) {
	conflicts, _, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return conflicts, err
	}

	u.producer.Produce(user.Nickname)

	return nil, nil
}

func (u *UserUseCase) CheckAuth(ctx context.Context, username string, password string) (*models.User, error) {
	user, _, err := u.userRepo.AuthUser(ctx, username, password)
	return user, err
}

func (u *UserUseCase) GetProfile(ctx context.Context, username string) (*models.User, error) {
	user, err := u.userRepo.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) ChangeProfile(ctx context.Context, user *models.User) (*models.User, error) {
	oldUser, err := u.userRepo.GetUser(ctx, user.Nickname)
	if err != nil {
		return nil, err
	}

	// do nothing if we dont need to change anything
	if user.Email == "" && user.Fullname == "" && user.About == "" {
		return oldUser, nil
	}
	// check empty fields
	if user.Fullname == "" {
		user.Fullname = oldUser.Fullname
	}
	if user.Email == "" {
		user.Email = oldUser.Email
	}
	if user.About == "" {
		user.About = oldUser.About
	}

	_, err = u.userRepo.ChangeUser(ctx, user)
	return user, err
}
