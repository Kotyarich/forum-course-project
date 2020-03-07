package http

import (
	"dbProject/common"
	"dbProject/forum"
	"github.com/dimfeld/httptreemux"
)

func RegisterHTTPEndpoints(router *httptreemux.TreeMux, uc forum.UseCase) {
	handler := NewHandler(uc)

	router.POST("/api/forum/:slug/create",
		common.CORSMiddlware(handler.ThreadCreateHandler))
	router.GET("/api/forums",
		common.CORSMiddlware(handler.ForumsHandler))
	router.GET("/api/forum/:slug/details",
		common.CORSMiddlware(handler.ForumDetailsHandler))
	router.GET("/api/forum/:slug/threads",
		common.CORSMiddlware(handler.ForumThreadsHandler))
	router.GET("/api/forum/:slug/users",
		common.CORSMiddlware(handler.ForumUsersHandler))
	router.POST("/api/forum/create",
		common.CORSMiddlware(handler.ForumCreateHandler))

	router.POST("/api/thread/:slug/create",
		common.CORSMiddlware(handler.ThreadPostCreateHandler))
	router.GET("/api/thread/:slug/details",
		common.CORSMiddlware(handler.GetThreadHandler))
	router.POST("/api/thread/:slug/details",
		common.CORSMiddlware(handler.PostThreadHandler))
	router.GET("/api/thread/:slug/posts",
		common.CORSMiddlware(handler.GetThreadPosts))
	router.POST("/api/thread/:slug/vote",
		common.CORSMiddlware(handler.ThreadVoteHandler))

	router.GET("/api/post/:id/details",
		common.CORSMiddlware(handler.GetPostHandler))
	router.POST("/api/post/:id/details",
		common.CORSMiddlware(handler.ChangePostHandler))

	router.POST("/api/service/clear",
		common.CORSMiddlware(handler.ClearHandler))
	router.GET("/api/service/status",
		common.CORSMiddlware(handler.StatusHandler))
}
