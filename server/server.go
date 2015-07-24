package main

import (
	"net/http"

	"github.com/cyarie/suchgopcandidatewow.com/server/router"
)

func main() {
	router := router.Router()

	http.ListenAndServe(":8080", router)
}
