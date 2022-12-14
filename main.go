package main

import (
	"net/http"
	"scootin/db"
	"scootin/logger"
	"scootin/service"
)

const appName = "Scootin"

func main() {

	log := logger.NewLogger()
	logger.InitLogger(log)
	defer logger.Sync()

	if err := db.InitiatePostgre(); err != nil {
		panic(err)
	}
	//  create a new *router instance
	router := service.NewRouter()
	logger.Fatal(http.ListenAndServe(":8080", router))
}
