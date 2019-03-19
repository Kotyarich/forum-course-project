package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func SetPostRouter(router *mux.Router) {
	router.HandleFunc("/api/post/{id:[0-9]+}/details", postHandler)
}

func postHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /post")
}