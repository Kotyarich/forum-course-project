package http

import (
	"dbProject/common"
	"dbProject/user"
	"github.com/labstack/echo/v4"
)

func RegisterHTTPEndpoints(router *echo.Echo, uc user.UseCase) {
	handler := NewHandler(uc)


	router.GET("/api/user/signout",
		common.CORSMiddlware(
			AuthMiddleware(handler.SignOutHandler, uc)))
	router.GET("/api/user/check",
		common.CORSMiddlware(
			AuthMiddleware(handler.UserCheckAuthHandler, uc)))
	router.POST("/api/user/:nickname/create",
		common.CORSMiddlware(handler.UserCreateHandler))
	router.GET("/api/user/:nickname/profile",
		common.CORSMiddlware(handler.UserGetHandler))
	router.POST("/api/user/:nickname/profile",
		common.CORSMiddlware(handler.UserPostHandler))

	router.POST("/api/user/auth", common.CORSMiddlware(handler.UserAuthHandler))
}
