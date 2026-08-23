// Package gateway handles user requests
package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"imago/pkg/storage"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func acceptFileTransformRequest(repo storage.FileRepo, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	w.Header().Set("Content-Type", "application/json")
	email := r.FormValue("email")
	encoder := json.NewEncoder(w)
	if email == "" {
		encoder.Encode(MessageResponse{Message: "Missing field 'email'"})
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		encoder.Encode(MessageResponse{Message: "Failed to save file"})
		return
	}
	fileID, err := repo.Save(r.Context(), file)
	if err != nil {
		encoder.Encode(MessageResponse{Message: "Failed to save file"})
		return
	}
	log.Printf("Uploaded file ID: %s\n", fileID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(MessageResponse{Message: "Uploaded"})
}

func createHandler(
	repo storage.FileRepo,
	handler func(repo storage.FileRepo, w http.ResponseWriter, r *http.Request),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(repo, w, r)
	}
}

func Run() error {
	config, err := NewConfig()
	if err != nil {
		log.Fatal("failed to initialize config")
	}
	repo := storage.NewRemoteFileRepo(config.FileRepoURL)
	http.HandleFunc("POST /transform", createHandler(repo, acceptFileTransformRequest))

	log.Printf("starting gateway server. listening at port %v", config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
}
