package server

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string
	Code    int
}

func WriteError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Message: message,
		Code:    code,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

var (
	InternalServerError = func(w http.ResponseWriter) {
		WriteError(w, "Unexpected error occurred", http.StatusInternalServerError)
	}
	BadRequestError = func(w http.ResponseWriter, err error) {
		WriteError(w, err.Error(), http.StatusBadRequest)
	}
	BadGatewayError = func(w http.ResponseWriter, err error) {
		WriteError(w, err.Error(), http.StatusBadGateway)
	}
)
