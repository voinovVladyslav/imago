package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("request path /")
	w.WriteHeader(404)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{Message: "Not Found"})
}

func getFile(repo FileRepo, w http.ResponseWriter, r *http.Request) {
	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte("File not found"))
		log.Println("unable to parse uuid:", err)
		return
	}
	f, err := repo.Get(r.Context(), fileID)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("File not found"))
		log.Println("file not found:", fileID, err)
		return
	}
	log.Println("return requested file:", fileID)
	defer f.Close()
	w.Header().Set("Content-Disposition", "attachment; filename=file")
	io.Copy(w, f)
}

type FileUploadResponse struct {
	ID uuid.UUID `json:"id"`
}

// User uploaded file using form and sending file as "file" field
func saveFile(repo FileRepo, w http.ResponseWriter, r *http.Request) {
	uploadedFile, _, err := r.FormFile("file")
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err != nil {
		w.WriteHeader(400)
		enc.Encode(MessageResponse{Message: "Failed to save file"})
		log.Println("failed to save file:", err)
		return
	}
	fileRecord, err := repo.Save(r.Context(), uploadedFile)
	if err != nil {
		w.WriteHeader(400)
		enc.Encode(MessageResponse{Message: "Failed to save file"})
		log.Println("failed to save file:", err)
		return
	}
	w.WriteHeader(201)
	enc.Encode(FileUploadResponse{ID: fileRecord.ID})
	log.Println("saved file:", fileRecord.ID)
}

func deleteFile(repo FileRepo, w http.ResponseWriter, r *http.Request) {
	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte("Invalid file ID"))
		log.Println("invalid file id:", err)
		return
	}
	err = repo.Delete(r.Context(), fileID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Failed to delete file"))
		log.Println("failed to delete file:", err)
		return
	}
	w.WriteHeader(204)
	log.Println("deleted file")
}

func createHandler(
	repo FileRepo,
	handler func(repo FileRepo, w http.ResponseWriter, r *http.Request),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(repo, w, r)
	}
}

func InitServer() error {
	config, err := NewConfig()
	if err != nil {
		return err
	}
	db, err := NewDB(config.SqliteDSN)
	if err != nil {
		return err
	}
	repo := NewLocalFileRepo(db, config.FileStorageDir)

	http.HandleFunc("GET /", statusHandler)
	http.HandleFunc("GET /file/{id}/", createHandler(repo, getFile))
	http.HandleFunc("POST /file/", createHandler(repo, saveFile))
	http.HandleFunc("DELETE /file/{id}/", createHandler(repo, deleteFile))

	log.Printf("starting storage server. listening at port %v", config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
}
