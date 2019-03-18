package routes

import (
	"fmt"
	"net/http"
)

func forumHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /forum")
}