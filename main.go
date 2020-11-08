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
	log.Printf("SENDGRID_API_KEY: %v", os.Getenv("SENDGRID_API_KEY"))

	r := mux.NewRouter()
	r.HandleFunc("/send", SendHandler(m)).Methods(http.MethodPost)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":80", r))
}
