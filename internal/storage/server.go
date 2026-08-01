package storage

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type MessageResponse struct {
	Message string `json:"message"`
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{Message: "Not Found"})
}

func getFile(repo FileRepo, w http.ResponseWriter, r *http.Request) {
	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte("File not found"))
		return
	}
	f, err := repo.Get(r.Context(), fileID)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("File not found"))
		return
	}
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
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageResponse{Message: "Failed to save file"})
		return
	}
	fileRecord, err := repo.Save(r.Context(), uploadedFile)
	if err != nil {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageResponse{Message: "Failed to save file"})
		return
	}
	w.WriteHeader(201)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FileUploadResponse{ID: fileRecord.ID})
}

func deleteFile(repo FileRepo, w http.ResponseWriter, r *http.Request) {
	fileID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte("Invalid file ID"))
		return
	}
	err = repo.Delete(r.Context(), fileID)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Failed to delete file"))
		return
	}
	w.WriteHeader(204)
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
	db, err := NewDB("storage.sqlite3")
	if err != nil {
		return err
	}
	repo := NewLocalFileRepo(db)

	http.HandleFunc("GET /", statusHandler)
	http.HandleFunc("GET /file/{id}/", createHandler(repo, getFile))
	http.HandleFunc("POST /file/", createHandler(repo, saveFile))
	http.HandleFunc("DELETE /file/{id}/", createHandler(repo, deleteFile))

	return http.ListenAndServe(":8000", nil)
}
