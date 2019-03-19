package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func SetForumRouter(router *mux.Router) {
	router.HandleFunc("/api/forum/create", createHandler)
	router.HandleFunc("/api/forum/{slug}/create", slugCreateHandler)
	router.HandleFunc("/api/forum/{slug}/{type}", slugHandler)
}

func createHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /forum")
	if request.Method == "POST" {
		//
	}
}

func slugHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /forum")
	if request.Method == "POST" {
		//
	}
}

func slugCreateHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /forum")
	if request.Method == "POST" {
		//
	}
}