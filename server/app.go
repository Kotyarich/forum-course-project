package server

import (
	"context"
	"dbProject/forum"
	forumHttp "dbProject/forum/delivery/http"
	forumPostgres "dbProject/forum/repository/postgres"
	forumUseCase "dbProject/forum/usecase"
	"dbProject/routes"
	"dbProject/user"
	userHttp "dbProject/user/delivery/http"
	userPostgres "dbProject/user/repository/postgres"
	userUceCase "dbProject/user/usecase"
	"github.com/dimfeld/httptreemux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer *http.Server

	userUC  user.UseCase
	forumUC forum.UseCase
}

func NewApp() *App {
	userRepo := userPostgres.NewUserRepository()
	forumRepo := forumPostgres.NewForumRepository()
	threadRepo := forumPostgres.NewThreadRepository()

	return &App{
		userUC: userUceCase.NewUserUseCase(
			userRepo,
			"hash_salt",
			[]byte("signing_key"),
			time.Hour*24*7),
		forumUC: forumUseCase.NewForumUseCase(forumRepo, threadRepo),
	}
}

func (a *App) Run(port string) error {
	router := httptreemux.New()

	routes.SetHomeRouter(router)
	routes.SetServiceRouter(router)
	routes.SetPostRouter(router)

	forumHttp.RegisterHTTPEndpoints(router, a.forumUC)
	userHttp.RegisterHTTPEndpoints(router, a.userUC)

	a.httpServer = &http.Server{
		Addr:           port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//go func() {
	if err := a.httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to listen and serve: %+v", err)
	}
	//}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}
