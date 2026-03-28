package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIError struct {
	Message string `json:"error"`
	Code    int    `json:"code"`
}

func (e *APIError) Error() string { return e.Message }

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("writeJSON encode error: %v", err)
	}
}

func WriteError(w http.ResponseWriter, apiErr *APIError) {
	WriteJSON(w, apiErr.Code, apiErr)
}
