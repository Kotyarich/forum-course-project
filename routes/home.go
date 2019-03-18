package routes

import (
	"fmt"
	"net/http"
)

func homeHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /")
}