package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return db, nil
}

type FileRecord struct {
	ID        uuid.UUID
	FilePath  string
	CreatedAt time.Time
}

type FileRepo interface {
	Save(ctx context.Context, f io.Reader) (FileRecord, error)
	Get(ctx context.Context, id uuid.UUID) (io.ReadCloser, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type LocalFileRepo struct {
	baseDir string
	db      *sql.DB
}

func (r *LocalFileRepo) Save(ctx context.Context, f io.Reader) (record *FileRecord, err error) {
	id := uuid.Must(uuid.NewV7())
	record = &FileRecord{id, path.Join(r.baseDir, id.String()), time.Now()}

	dest, err := os.Create(record.FilePath)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(dest, f)

	defer func() {
		dest.Close()
		if err != nil {
			os.Remove(record.FilePath)
		}
	}()

	if err != nil {
		return nil, err
	}

	query := "INSERT INTO file_registry (id, filepath, created_at) values (?, ?, ?)"
	_, err = r.db.ExecContext(ctx, query, record.ID, record.FilePath, record.CreatedAt.UTC().Unix())
	if err != nil {
		os.Remove(record.FilePath)
		return nil, err
	}
	return record, nil
}

func (r *LocalFileRepo) Get(ctx context.Context, id uuid.UUID) (f io.ReadCloser, err error) {
	query := "SELECT filepath FROM file_registry WHERE id = ? LIMIT 1"
	rows, err := r.db.QueryContext(ctx, query, id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	hasResult := rows.Next()
	if !hasResult {
		return nil, errors.New("not found in the database")
	}
	var fp string
	err = rows.Scan(&fp)
	if err != nil {
		return nil, err
	}
	return os.Open(fp)
}

func (r *LocalFileRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM file_registry WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func NewLocalFileRepo(db *sql.DB) *LocalFileRepo {
	return &LocalFileRepo{baseDir: ".storage",db: db}
}
