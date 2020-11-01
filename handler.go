package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

//`{
//	"type"
//}`

// Type represents the kind of email template that will be used for sending
type Type string

const (
	Driving Type = "Driving"
)

// sendRequest holds the request structure for SendHandler
type sendRequest struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

type sendResponse struct {
	Error string `json:"error,omitempty"`
}

// SendHandler handles http requests for the /send resource
func SendHandler(c *Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("error reading request body: %v", err)
		}

		req, err := decodeRequestBody(b)

		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			res := &sendResponse{err.Error()}
			json.NewEncoder(w).Encode(res)
			return
		}

		switch req.Type {
		case Driving:
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNoContent)
			break
		default:
			http.Error(w, errors.New("unrecognized email request type").Error(), http.StatusBadRequest)
			break
		}
	}
}

func decodeRequestBody(b []byte) (*sendRequest, error) {
	var req sendRequest
	err := json.Unmarshal(b, &req)

	if err != nil {
		log.Printf("parsing of request body failed: %v", err)
		return nil, errors.New("failed to parse request body")
	}

	return &req, nil
}
