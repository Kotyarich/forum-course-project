package main

import (
	"dbProject/routes"
	"fmt"
	"github.com/dimfeld/httptreemux"
	"net/http"
)

func main() {
	router := httptreemux.New()

	routes.SetHomeRouter(router)
	routes.SetForumRouter(router)
	routes.SetServiceRouter(router)
	routes.SetPostRouter(router)
	routes.SetThreadRouter(router)
	routes.SetUserRouter(router)

	server := http.Server{
		Addr:    ":5000",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
