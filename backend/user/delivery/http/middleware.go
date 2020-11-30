package http

import (
	"context"
	"dbProject/user"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(f echo.HandlerFunc, uc user.UseCase) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Auth")
		ctx := c.Request().Context()

		if err != nil {
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "user", nil)))
			return f(c)
		}

		u, err := uc.CheckAuth(c.Request().Context(), cookie.Value)
		if err != nil {
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "user", nil)))
		} else {
			c.SetRequest(c.Request().WithContext(context.WithValue(ctx, "user", u)))
		}

		return f(c)
	}
}