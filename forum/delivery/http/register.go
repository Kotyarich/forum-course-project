package http

import (
	"dbProject/forum"
	"github.com/dimfeld/httptreemux"
)

func RegisterHTTPEndpoints(router *httptreemux.TreeMux, uc forum.UseCase) {
	handler := NewHandler(uc)

	router.POST("/api/forum/:slug/create", handler.ThreadCreateHandler)
	router.GET("/api/forum/:slug/details", handler.ForumDetailsHandler)
	router.GET("/api/forum/:slug/threads", handler.ForumThreadsHandler)
	router.GET("/api/forum/:slug/users", handler.ForumUsersHandler)
	router.POST("/api/forum/create", handler.ForumCreateHandler)
}
