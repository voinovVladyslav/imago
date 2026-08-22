// Package gateway handles user requests
package gateway

import (
	"encoding/json"
	"net/http"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func acceptFileTransformRequest(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageResponse{Message: "Missing field 'email'"})
		return
	}
	_, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageResponse{Message: "Failed to save file"})
		return
	}
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{Message: "Uploaded"})
}

func Run() error {
	http.HandleFunc("POST /transform", acceptFileTransformRequest)
	return http.ListenAndServe(":8000", nil)
}
