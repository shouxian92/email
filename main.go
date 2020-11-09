package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/shouxian92/email/driving"
	"github.com/shouxian92/email/service"
)

func main() {
	m := map[string]service.Client{
		"driving": driving.GetEmailInstance(),
	}

	log.Println("Starting up..")
	log.Printf("EMAIL_FROM: %v", os.Getenv("EMAIL_FROM"))
	log.Printf("EMAIL_FROM_ADDRESS: %v", os.Getenv("EMAIL_FROM_ADDRESS"))

	r := mux.NewRouter()
	r.HandleFunc("/send", SendHandler(m)).Methods(http.MethodPost)
	http.Handle("/", r)

	port := os.Getenv("PORT")

	if len(port) == 0 {
		port = "80"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
