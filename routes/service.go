package routes

import (
	"fmt"
	"net/http"
)

func serviceHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("hi /serve")
}