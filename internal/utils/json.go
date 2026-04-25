package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJSON sets the Content-Type to application/json, writes the status code,
// and encodes the payload as JSON to the response writer.
func WriteJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

// WriteError sends a JSON-formatted error response.
func WriteError(w http.ResponseWriter, status int, message string) error {
	return WriteJSON(w, status, map[string]string{"error": message})
}
