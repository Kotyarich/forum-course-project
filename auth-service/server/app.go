package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"user-service/auth"
	userHttp "user-service/auth/delivery/http"
	userUceCase "user-service/auth/usecase"
	"user-service/user-service"
)

type App struct {
	httpServer *http.Server

	userUC auth.UseCase
}

const userServiceUrl = "http://users:5001/"

func NewApp() *App {
	userService := user_service.NewUserService(userServiceUrl)

	return &App{
		userUC: userUceCase.NewUserUseCase(
			userService,
			"hash_salt",
			[]byte("signing_key"),
			time.Hour*24*7),
	}
}

func (a *App) Run(port string) error {
	router := echo.New()

	userHttp.RegisterHTTPEndpoints(router, a.userUC)
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
