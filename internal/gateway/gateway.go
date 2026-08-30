// Package gateway handles user requests
package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"imago/pkg/broker"
	"imago/pkg/cache"
	"imago/pkg/storage"

	"github.com/redis/go-redis/v9"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func acceptFileTransformRequest(
	session *Session,
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
	fileID, err := session.files.Save(r.Context(), file)
	if err != nil {
		w.WriteHeader(400)
		encoder.Encode(MessageResponse{Message: "Failed to save file"})
		return
	}
	log.Printf("Uploaded file ID: %s\n", fileID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(MessageResponse{Message: "Uploaded"})
	key := fmt.Sprintf("request:%s", fileID)
	err = session.cache.Set(r.Context(), key, email, time.Hour*24).Err()
	if err != nil {
		log.Printf("failed to set cache key: %s", err)
	}
	session.broker.SendJSON(
		MessageResponse{
			Message: fmt.Sprintf("uploaded file: %s", fileID),
		},
		"file.uploaded",
	)
}

func createHandler(
	session *Session,
	handler func(
		session *Session,
		w http.ResponseWriter,
		r *http.Request,
	),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(session, w, r)
	}
}

type Session struct {
	cache  *redis.Client
	broker *broker.Broker
	files  storage.FileRepo
}

func Run() error {
	config, err := NewConfig()
	if err != nil {
		log.Fatal("failed to initialize config")
	}
	repo := storage.NewRemoteFileRepo(config.FileRepoURL)

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
		log.Fatal("failed to initialize broker")
	}

	c := cache.NewCache()
	defer c.Close()
	s := &Session{cache: c, broker: b, files: repo}

	http.HandleFunc("POST /transform", createHandler(s, acceptFileTransformRequest))

	log.Printf("starting gateway server. listening at port %v", config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
}
