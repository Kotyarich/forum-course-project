package http

import (
	"context"
	"forum-service/forum"
	"github.com/labstack/echo/v4"
)

func AuthMiddlware(as forum.AuthService, f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Auth")
		ctx := c.Request().Context()

		if err != nil {
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "auth", nil)))
			return f(c)
		}

		u, err := as.CheckAuth(cookie.Value)
		if err != nil {
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "auth", nil)))
		} else {
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "auth", u)))
		}

		return f(c)
	}
}
