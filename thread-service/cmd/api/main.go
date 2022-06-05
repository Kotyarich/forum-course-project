// Forums API
//
// BMSTU Web Course 2020 project
//
// Terms Of Service:
//
//     Schemes: http
//     Host: localhost:5000
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"thread-service/server"
	"log"
	"os"
)

func main() {
	_ = os.Setenv("TZ", "Europe/Moscow")
	app := server.NewApp()

	if err := app.Run(":5005"); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
