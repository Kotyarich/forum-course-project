package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"user-service/auth"
)

func AuthMiddleware(f echo.HandlerFunc, uc auth.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Auth")
		ctx := c.Request().Context()

		if err != nil {
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "auth", nil)))
			return f(c)
		}

		u, err := uc.CheckAuth(c.Request().Context(), cookie.Value)
		if err != nil {
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "auth", nil)))
		} else {
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "auth", u)))
		}

		return f(c)
	}
}