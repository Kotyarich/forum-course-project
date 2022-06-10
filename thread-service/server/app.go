package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"thread-service/auth-service"
	"thread-service/forum"
	"thread-service/forum-service"
	forumHttp "thread-service/forum/delivery/http"
	forumPostgres "thread-service/forum/repository/postgres"
	forumUseCase "thread-service/forum/usecase"
	"thread-service/queue"
	"time"
)

type App struct {
	httpServer *http.Server

	forumUC forum.UseCase
}

func NewApp() *App {
	threadRepo := forumPostgres.NewThreadRepository()
	forumService := forum_service.NewForumService("http://forums:5000/")
	producer := queue.NewProducer()

	return &App{
		forumUC: forumUseCase.NewForumUseCase(threadRepo, *forumService, producer),
	}
}

func (a *App) Run(port string) error {
	router := echo.New()

	auth := auth_service.NewAuthService("http://auths:5002/")
	forumHttp.RegisterHTTPEndpoints(router, a.forumUC, auth)
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	router.Use(middleware.Logger())

	router.Static("/", "swaggerui")

	a.httpServer = &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
