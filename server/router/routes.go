package router

import (
	"github.com/cyarie/suchgopcandidatewow.com/server/handlers"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	HandlerFunc handlers.WebHandler
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handlers.WebHandler{handlers.IndexHandler},
	},
}