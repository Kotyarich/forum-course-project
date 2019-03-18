package routes

import (
	"fmt"
	"net/http"
)

func threadHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /thread")
}