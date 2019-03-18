package routes

import (
	"fmt"
	"net/http"
)

func userHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /")
}