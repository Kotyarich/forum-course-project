package usecase

import (
	"context"
	"crypto/sha1"
	"dbProject/models"
	userPkg "dbProject/user"
	"fmt"
	"time"
)

type UserUseCase struct {
	userRepo       userPkg.Repository
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewUserUseCase(
	userRepo userPkg.Repository,
	hashSalt string,
	signingKey []byte,
	sessionDurationSeconds time.Duration) *UserUseCase {
	return &UserUseCase{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * sessionDurationSeconds,
	}
}

func (u *UserUseCase) SignUp(ctx context.Context, user *models.User) ([]*models.User, string, error) {
	user.Password = passwordHash(user.Password, u.hashSalt)
	conflicts, userId, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return conflicts, "", err
	}

	token, err := u.userRepo.CreateSession(ctx, userId)
	if err != nil {
		return nil, "", err
	}

	return nil, token, nil
}

func (u *UserUseCase) SignIn(ctx context.Context, username, password string) (*models.User, string, error) {
	password = passwordHash(password, u.hashSalt)
	user, userId, err := u.userRepo.AuthUser(ctx, username, password)
	if err != nil {
		return nil, "", userPkg.ErrUserNotFound
	}

	token, err := u.userRepo.CreateSession(ctx, userId)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (u *UserUseCase) SignOut(ctx context.Context, token string) error {
	err := u.userRepo.DeleteSession(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) CheckAuth(ctx context.Context, token string) (*models.User, error) {
	user, err := u.userRepo.CheckSession(ctx, token)
	if err != nil {
		return nil, err
	}

	return user, nil
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

func passwordHash(password, hashSalt string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(hashSalt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
