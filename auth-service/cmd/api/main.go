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
	"log"
	"os"
	"time"
	"user-service/server"
)

func main() {
	_ = os.Setenv("TZ", "Europe/Moscow")
	app := server.NewApp()
	time.Sleep(time.Minute)
	if err := app.Run(":5002"); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
