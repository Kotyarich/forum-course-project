package routes

import (
	"fmt"
	"net/http"
)

func postHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /post")
}