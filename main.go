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

/*
1. The scooters report an event when a trip begins, report an event when the
trip ends, and send in periodic updates on their location. After beginning a
trip, the scooter is considered occupied. After a trip ends the scooter
becomes free for use. A location update must contain the time, and
geographical coordinates.

2. Mobile clients can query scooter locations and statuses in any rectangular
location (e.g. two pair of coordinates), and filter them by status. While there
will be no actual mobile clients, implement child process that would start
with main process and spawn three fake clients using API randomly (finding
scooters, travelling for 10-15 seconds whilst updating location every 3
seconds, and resting for 2-5 seconds before starting next trip). Client
movement in straight line will be good enough.

*/
