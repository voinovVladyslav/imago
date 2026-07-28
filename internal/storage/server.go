package storage

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(MessageResponse{Message: "Not Found"})
	if err != nil {
		return
	}
	w.Write(data)
}

func getFile(w http.ResponseWriter, r *http.Request) {
	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(422)
		data, err := json.Marshal(MessageResponse{Message: "file ID should be valid uuid"})
		if err != nil {
			return
		}
		w.Write(data)
		return
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, "Valid uuid: %s", fileID)
}

func Serve() {
	http.HandleFunc("GET /", statusHandler)
	http.HandleFunc("GET /file/{id}/", getFile)

	http.ListenAndServe(":8000", nil)
}
