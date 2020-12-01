package usecase

import (
	"context"
	"crypto/sha1"
	"dbProject/models"
	userPkg "dbProject/user"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
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

type Claims struct {
	jwt.StandardClaims
	Nickname string `json:"nickname"`
}

func (u *UserUseCase) SignUp(ctx context.Context, user *models.User) ([]*models.User, string, error) {
	user.Password = passwordHash(user.Password, u.hashSalt)
	conflicts, _, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return conflicts, "", err
	}

	tk := &Claims{Nickname: user.Nickname}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("some_key"))

	return nil, tokenString, nil
}

func (u *UserUseCase) SignIn(ctx context.Context, username, password string) (*models.User, string, error) {
	password = passwordHash(password, u.hashSalt)
	user, _, err := u.userRepo.AuthUser(ctx, username, password)
	if err != nil {
		return nil, "", userPkg.ErrUserNotFound
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Nickname: user.Nickname,
	})

	tokenStr, err := token.SignedString([]byte("some_key"))

	return user, tokenStr, err
}

func ParseToken(token string, signKey []byte) (string, error) {
	parsedToken, err := jwt.ParseWithClaims(
		token, &Claims{},
		func(token *jwt.Token) (i interface{}, e error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return signKey, nil
		})

	if err != nil {
		return "", err
	}

	if claims, ok := parsedToken.Claims.(*Claims); ok && parsedToken.Valid {
		return claims.Nickname, nil
	}

	return "", fmt.Errorf("invalid access token")
}

func (u *UserUseCase) CheckAuth(ctx context.Context, token string) (*models.User, error) {
	nickname, err := ParseToken(token, []byte("some_key"))
	if err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetUser(ctx, nickname)

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

func passwordHash(password, hashSalt string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(hashSalt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
