package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"time"
	userPkg "user-service/auth"
	"user-service/models"
)

type UserUseCase struct {
	userService    userPkg.UserService
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewUserUseCase(
	userService userPkg.UserService,
	hashSalt string,
	signingKey []byte,
	sessionDurationSeconds time.Duration) *UserUseCase {
	return &UserUseCase{
		userService:    userService,
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
	_, conflicts, err := u.userService.CreateUser(ctx, user)
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
	user, err := u.userService.CheckAuth(ctx, username, password)
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

	user, err := u.userService.GetUser(ctx, nickname)
	return user, err
}

func passwordHash(password, hashSalt string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(hashSalt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}
