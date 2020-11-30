package common

import (
	"github.com/labstack/echo/v4"
)

func setCORSHeaders(c echo.Context) {
	c.Response().Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD, PUT")
	c.Response().Header().Set("Access-Control-Max-Age", "600")
}

func CORSMiddlware(f echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		setCORSHeaders(c)
		return f(c)
	}
}

func CORSHandler(c echo.Context) {
	setCORSHeaders(c)
}

