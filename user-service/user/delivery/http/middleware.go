package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
	"user-service/models"
	"user-service/user"
)

const authUrl = "http://auths:5002"

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

func AuthMiddleware(f echo.HandlerFunc, _ user.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Auth")
		ctx := c.Request().Context()

		if err != nil {
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "user", nil)))
			return f(c)
		}

		u, err := checkAuth(c.Request().Context(), cookie.Value)
		log.Println("checked")
		if err != nil {
			log.Println(err)
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "user", nil)))
		} else {
			log.Println(u.Nickname)
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "user", u)))
		}

		return f(c)
	}
}
