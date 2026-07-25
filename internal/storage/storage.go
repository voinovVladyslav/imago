// Package storage handles file uploads/retrieval
package storage

import (
	"database/sql"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type FileRecord struct {
	ID        uuid.UUID
	FilePath  string
	CreatedAt time.Time
}

type FileStorage interface {
	Save(f io.Reader) (id uuid.UUID)
	Get(id uuid.UUID) (f io.Reader)
	Delete(id uuid.UUID)
}

type LocalFileStorage struct {
	baseDir string
}

func (s *LocalFileStorage) Save(f io.Reader) uuid.UUID {
	id := uuid.Must(uuid.NewV7())
	saveDir := filepath.Join(s.baseDir, id.String())
	fmt.Println("created file", saveDir)
	return id
}

func (s *LocalFileStorage) Get(id uuid.UUID) io.Reader {
	return strings.NewReader("test")
}

func (s *LocalFileStorage) Delete(id uuid.UUID) {
	fmt.Println("deleting", id)
}

func readFromDB() {
	db, err := sql.Open("sqlite", "storage.sqlite3")
	if err != nil {
		panic("Failed to open db connection")
	}
	rows, err := db.Query("SELECT id, filepath, created_at FROM file_registry;")
	if err != nil {
		panic("failed to fetch rows")
	}

	var files []FileRecord
	defer rows.Close()
	for rows.Next() {
		var f FileRecord
		var timestamp int64
		err := rows.Scan(&f.ID, &f.FilePath, &timestamp)
		if err != nil {
			panic("no scan")
		}
		f.CreatedAt = time.Unix(timestamp, 0).UTC()
		files = append(files, f)
	}
	fmt.Println(files)
}

func writeToDB() {
	db, err := sql.Open("sqlite", "storage.sqlite3")
	if err != nil {
		panic("Failed to open db connection")
	}
	record := FileRecord{uuid.Must(uuid.NewV7()), "/main/file/path/photo.jpg", time.Now()}

	query := "INSERT INTO file_registry (id, filepath, created_at) values (?, ?, ?)"
	_, err = db.Exec(query, record.ID, record.FilePath, record.CreatedAt.UTC().Unix())
	if err != nil {
		fmt.Println("Error while inserting", err)
		return
	}
	fmt.Println("done")
}

func Run() {
	readFromDB()
}
