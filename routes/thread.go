package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func SetThreadRouter(router *mux.Router) {
	router.HandleFunc("/api/thread/{slug}/create", threadHandler)
	router.HandleFunc("/api/thread/{slug}/details", threadHandler)
	router.HandleFunc("/api/thread/{slug}/posts", threadHandler)
	router.HandleFunc("/api/thread/{slug}/vote", threadHandler)

}

func threadHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /thread")
}