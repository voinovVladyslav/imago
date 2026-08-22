// Package shared contains reusable structucts
package shared

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
)

type FileRepo interface {
	Save(ctx context.Context, f io.Reader) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (io.ReadCloser, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type RemoteFileRepo struct {
	endpoint string
}

type FileUploadResponse struct {
	ID uuid.UUID `json:"id"`
}

func (r *RemoteFileRepo) Save(ctx context.Context, f io.Reader) (fileID uuid.UUID, err error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile("file", "image.jpg")
	io.Copy(fw, f)
	w.Close()
	resp, err := http.Post(fmt.Sprintf("%s/file", r.endpoint), w.FormDataContentType(), &buf)
	if err != nil {
		return uuid.Nil, err
	}
	var msg FileUploadResponse
	err = json.NewDecoder(resp.Body).Decode(&msg)
	if err != nil {
		return uuid.Nil, err
	}
	return msg.ID, nil
}

func (r *RemoteFileRepo) Get(ctx context.Context, id uuid.UUID) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/file/%s/", r.endpoint, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("failed to get file")
	}
	return resp.Body, nil
}

func (r *RemoteFileRepo) Delete(ctx context.Context, id uuid.UUID) error {
	url := fmt.Sprintf("%s/file/%s/", r.endpoint, id)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return errors.New("failed to delete file")
	}
	return nil
}
