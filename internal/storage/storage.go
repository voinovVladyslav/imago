// Package storage handles file uploads/retrieval
package storage

import (
	"context"
	"fmt"
	"io"
	"strings"
)

func _Run() error {
	db, err := NewDB("storage.sqlite3")
	if err != nil {
		return err
	}
	repo := NewLocalFileRepo(db)

	ctx := context.Background()
	record, err := repo.Save(ctx, strings.NewReader("test file contents"))
	if err != nil {
		return err
	}
	fmt.Println("Saved to db:", record)
	f, err := repo.Get(ctx, record.ID)
	if err != nil {
		return err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	fmt.Println("Got from file itself:", string(data))
	err = repo.Delete(ctx, record.ID)
	return err
}

func Run() error {
	return InitServer()
}
