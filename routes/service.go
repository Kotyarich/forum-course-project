package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func SetServiceRouter(router *mux.Router) {
	router.HandleFunc("/api/service/clear", clearHandler)
	router.HandleFunc("/api/service/status", statusHandler)
}

func clearHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /serve")
}

func statusHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /serve")
}