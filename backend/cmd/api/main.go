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
