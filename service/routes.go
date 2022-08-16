package service

import (
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Method string            //HTTP method
	Path   string            //url endpoint
	Handle httprouter.Handle //Controller function which dispatches the right HTML page and/or data for each route
}

type Routes []Route

var routes = Routes{
	Route{
		"GET",
		"/",
		Index,
	},
	Route{
		"GET",
		"/v0.1/scooters",
		ListAvailableScooter,
	},
	Route{
		"POST",
		"/v0.1/scooter",
		CreateScooter,
	},
	Route{
		"POST",
		"/v0.1/user",
		CreateUser,
	},
	Route{
		"PUT",
		"/v0.1/scooter/book/:id",
		BookScooter,
	},
	Route{
		"PUT",
		"/v0.1/scooter/release/:id",
		ReleaseScooter,
	},
}
