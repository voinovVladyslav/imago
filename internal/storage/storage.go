// Package storage handles file uploads/retrieval
package storage

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/uuid"
)

type FileStorage interface {
	Save(f io.Reader) (id uuid.UUID)
	Get(id uuid.UUID) (f io.Reader)
	Delete(id uuid.UUID)
}

type LocalFileStorage struct{}

func (s *LocalFileStorage) Save(f io.Reader) uuid.UUID {
	id := uuid.Must(uuid.NewV7())
	return id
}

func Run() {
	fmt.Println("Hello from storage")
	s := LocalFileStorage{}
	reader := strings.NewReader("test")
	fmt.Println(s.Save(reader))
}
