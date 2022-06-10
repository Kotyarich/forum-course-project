package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"statistic-service/models"
	"statistic-service/statistic"
	"time"
)

const authUrl = "http://auths:5002"

// Информация о пользователе
// swagger:model User
type userInput struct {
	// Описание пользователя
	//
	// example: This is the day you will always remember as the day that you almost caught Captain Jack Sparrow!
	About string `json:"about"`

	// Почтовый адрес пользователя
	//
	// format: email
	// example: captaina@blackpearl.sea
	Email string `json:"email"`

	// Полное имя пользователя
	//
	// example: Captain Jack Sparrow
	Fullname string `json:"fullname"`

	// Имя пользователя (уникальное поле)
	//
	// format: identity
	// read only: true
	// example: j.sparrow
	Nickname string `json:"nickname"`

	// Пароль пользователя
	//
	// example: 123456
	Password string `json:"password,omitempty"`

	IsAdmin bool `json:"isAdmin"`
}

func userInputToModel(user userInput) *models.User {
	return &models.User{
		Nickname: user.Nickname,
		Fullname: user.Fullname,
		Password: user.Password,
		Email:    user.Email,
		About:    user.About,
		IsAdmin:  user.IsAdmin,
	}
}

func checkAuth(_ context.Context, token string) (*models.User, error) {
	url := fmt.Sprintf("%s/user/check", authUrl)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("content-type", "application/json")
	req.AddCookie(&http.Cookie{Name: "Auth", Value: token, Expires: time.Now().Add(time.Hour)})

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	//defer func() { _ = resp.Body.Close() }()

	var user userInput
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		_ = resp.Body.Close()
		return nil, err
	}
	_ = resp.Body.Close()
	return userInputToModel(user), nil
}

func AuthMiddleware(f echo.HandlerFunc, _ statistic.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Auth")
		ctx := c.Request().Context()

		if err != nil {
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "user", nil)))
			log.Println("user is nil")
			return f(c)
		}

		u, err := checkAuth(c.Request().Context(), cookie.Value)
		if err != nil || !u.IsAdmin {
			log.Println("user is nil or not admin")
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "user", nil)))
		} else {
			log.Println("user is admin")
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "user", u)))
		}

		return f(c)
	}
}
