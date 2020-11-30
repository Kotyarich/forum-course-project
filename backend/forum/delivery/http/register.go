package http

import (
	"dbProject/common"
	"dbProject/forum"
	"github.com/labstack/echo/v4"
)

func registerHTTPForumEndpoints(router *echo.Echo, uc forum.UseCaseForum) {
	handler := NewForumHandler(uc)

	router.POST("/api/forum/create",
		common.CORSMiddlware(handler.ForumCreateHandler))
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
}

func registerHTTPThreadEndpoints(router *echo.Echo, uc forum.UseCaseThread) {
	handler := NewThreadHandler(uc)

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
}

func registerHTTPPostEndpoints(router *echo.Echo, uc forum.UseCasePost) {
	handler := NewPostHandler(uc)

	router.GET("/api/post/:id/details",
		common.CORSMiddlware(handler.GetPostHandler))
	router.POST("/api/post/:id/details",
		common.CORSMiddlware(handler.ChangePostHandler))
}

func registerHTTPServiceEndpoints(router *echo.Echo, uc forum.UseCaseService) {
	handler := NewServiceHandler(uc)

	router.POST("/api/service/clear",
		common.CORSMiddlware(handler.ClearHandler))
	router.GET("/api/service/status",
		common.CORSMiddlware(handler.StatusHandler))
}

func RegisterHTTPEndpoints(router *echo.Echo, uc forum.UseCase) {
	registerHTTPForumEndpoints(router, uc.ForumUseCase)
	registerHTTPThreadEndpoints(router, uc.ThreadUseCase)
	registerHTTPPostEndpoints(router, uc.PostUseCase)
	registerHTTPServiceEndpoints(router, uc.ServiceUseCase)
}
