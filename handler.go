package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/shouxian92/email/service"
)

// Type represents the kind of email template that will be used for sending
type Type string

type sendResponse struct {
	Error string `json:"error,omitempty"`
}

// SendHandler handles http requests for the /send resource
func SendHandler(m map[string]service.Client) func(http.ResponseWriter, *http.Request) {
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

		c, ok := m[req.Type]
		if !ok {
			res := &sendResponse{"unrecognized email request type"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		req.Subject = "New driving lessons available!"
		err = c.Send(*req)

		if err != nil {
			switch err.Error() {
			case "too many requests":
				w.WriteHeader(http.StatusTooManyRequests)
				break
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				break
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func decodeRequestBody(b []byte) (*service.RequestContext, error) {
	var req service.RequestContext
	err := json.Unmarshal(b, &req)

	if err != nil {
		log.Printf("parsing of request body failed: %v", err)
		return nil, errors.New("failed to parse request body")
	}

	if len(req.To) == 0 {
		return &req, errors.New("'to' is required")
	}

	if len(req.Type) == 0 {
		return &req, errors.New("'type' is required")
	}

	return &req, nil
}
