package storage

import (
	"context"
	"errors"
	"io"
	"time"
)

var (
	ErrNotFound   = errors.New("file not found")
	ErrPermission = errors.New("permission denied")
)

type FileInfo struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	IsDir   bool      `json:"isDir"`
	ModTime time.Time `json:"modTime"`
}

type Storage interface {
	List(ctx context.Context, path string) ([]FileInfo, error)
	Read(ctx context.Context, path string) (io.ReadCloser, error)
	Write(ctx context.Context, path string, r io.Reader) error
	Delete(ctx context.Context, path string) error
	Stat(ctx context.Context, path string) (*FileInfo, error)
}
