package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func SetHomeRouter(router *mux.Router) {
	router.HandleFunc("/api", homeHandler)
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /home")
}