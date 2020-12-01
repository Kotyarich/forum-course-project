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
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: KEY
//          in: header
//
// swagger:meta
package main

import (
	"dbProject/server"
	"log"
	"os"
)

func main() {
	os.Setenv("TZ", "Europe/Moscow")
	app := server.NewApp()

	if err := app.Run(":5000"); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
