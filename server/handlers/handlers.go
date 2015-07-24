package handlers

import (
	"log"
	"net/http"
)

// Rolling up a custom handler to better deal with HTTP errors and the such. Also lets us pass through an environment
// struct/variables without having to use globals.
type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

type WebHandler struct {
	H func(w http.ResponseWriter, r *http.Request) error
}

// Write a function receiver to integrate the structs we made up above
func (h WebHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(w, r)
	if err != nil {
		switch e := err.(type) {
			case Error:
			log.Printf("HTTP %d - %s", e.Status(), e)
			default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) error {
	var err error

	http.ServeFile(w, r, "./html/index.html")

	return err
}
