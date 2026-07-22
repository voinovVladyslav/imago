// Package storage handles file uploads/retrieval
package storage

import (
	"fmt"
	"io"
)

type FileStorage interface {
	Save(f io.Reader) (id string)
	Get()
	Delete()
}

func Run() {
	fmt.Println("Hello from storage")
}
