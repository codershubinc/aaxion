package helpers

import (
	"encoding/json"
	"net/http"
)

func SetJSONResponce(w http.ResponseWriter, data any) {
	defer w.(http.Flusher).Flush()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
