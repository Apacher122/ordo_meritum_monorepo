package middleware

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON is a helper function that writes a JSON response to the http.ResponseWriter.
// It sets the "Content-Type" header to "application/json", writes the provided
// statusCode, and then marshals the payload into the response body.
// If the payload is nil, no body is written. It logs an error if the JSON
// encoding fails but does not panic.
func JSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}
