package http

import (
	"dbProject/user"
	"github.com/dimfeld/httptreemux"
)

func RegisterHTTPEndpoints(router *httptreemux.TreeMux, uc user.UseCase) {
	handler := NewHandler(uc)

	router.POST("/api/user/:nickname/create", handler.UserCreateHandler)
	router.GET("/api/user/:nickname/profile", handler.UserGetHandler)
	router.POST("/api/user/:nickname/profile", handler.UserPostHandler)
	router.GET("/api/user/auth", handler.UserAuthHandler)
}
