package router

import (
	"github.com/gorilla/mux"
	"github.com/cyarie/suchgopcandidatewow.com/server/handlers"
)

func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler handlers.WebHandler

		handler = handlers.WebHandler{route.HandlerFunc.H}
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(&handler)
	}

	return router
}