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
		"/scooters",
		ListAvailableScooter,
	},
	Route{
		"POST",
		"/scooter",
		CreateScooter,
	},
	Route{
		"PUT",
		"/scooter/book/:id",
		BookScooter,
	},
	Route{
		"POST",
		"/user",
		CreateUser,
	},
}
