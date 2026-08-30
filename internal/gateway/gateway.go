// Package gateway handles user requests
package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"imago/pkg/broker"
	"imago/pkg/storage"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func acceptFileTransformRequest(
	b *broker.Broker,
	repo storage.FileRepo,
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Header().Set("Content-Type", "application/json")
	email := r.FormValue("email")
	encoder := json.NewEncoder(w)
	if email == "" {
		w.WriteHeader(400)
		encoder.Encode(MessageResponse{Message: "Missing field 'email'"})
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(400)
		encoder.Encode(MessageResponse{Message: "Failed to save file"})
		return
	}
	fileID, err := repo.Save(r.Context(), file)
	if err != nil {
		w.WriteHeader(400)
		encoder.Encode(MessageResponse{Message: "Failed to save file"})
		return
	}
	log.Printf("Uploaded file ID: %s\n", fileID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(MessageResponse{Message: "Uploaded"})
	b.SendJSON(
		MessageResponse{
			Message: fmt.Sprintf("uploaded file: %s", fileID),
		},
		"file.uploaded",
	)
}

func createHandler(
	b *broker.Broker,
	repo storage.FileRepo,
	handler func(
		b *broker.Broker,
		repo storage.FileRepo,
		w http.ResponseWriter,
		r *http.Request,
	),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(b, repo, w, r)
	}
}

func Run() error {
	config, err := NewConfig()
	if err != nil {
		log.Fatal("failed to initialize config")
	}

	b, err := broker.NewBroker(
		broker.BrokerConfig{
			User:  "imago",
			Pass:  "strongpassword",
			Host:  "localhost",
			Port:  "5672",
			Vhost: "imago",
		},
	)
	if err != nil {
		log.Fatal("failed to initialize config")
	}
	repo := storage.NewRemoteFileRepo(config.FileRepoURL)
	http.HandleFunc("POST /transform", createHandler(b, repo, acceptFileTransformRequest))

	log.Printf("starting gateway server. listening at port %v", config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
}
