// Package storage handles file uploads/retrieval
package storage

import (
	"context"

	"github.com/google/uuid"
)

func Run() error {
	db, err := NewDB("storage.sqlite3")
	if err != nil {
		return err
	}
	repo := NewLocalFileRepo(db)

	ctx := context.Background()
	_, err = repo.Get(ctx, uuid.Must(uuid.NewV7()))
	return err
}
