package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	c, err := NewInstance()

	if err != nil {
		panic("unable to init email client instance")
	}

	r := mux.NewRouter()
	r.HandleFunc("/send", SendHandler(c)).Methods(http.MethodPost)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", r))
}
