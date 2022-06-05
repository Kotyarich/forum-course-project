package server

import (
	"context"
	"forum-service/auth-service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"forum-service/forum"
	forumHttp "forum-service/forum/delivery/http"
	forumPostgres "forum-service/forum/repository/postgres"
	forumUseCase "forum-service/forum/usecase"
)

type App struct {
	httpServer *http.Server

	forumUC forum.UseCase
}

func NewApp() *App {
	forumRepo := forumPostgres.NewForumRepository()

	return &App{
		forumUC: forumUseCase.NewForumUseCase(forumRepo),
	}
}

func (a *App) Run(port string) error {
	router := echo.New()

	auth := auth_service.NewAuthService("http://localhost:5002/")
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
