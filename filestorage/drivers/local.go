package drivers

import (
	"bytes"
	"fmt"
	"github.com/rayyone/go-core/filestorage/config"
	"github.com/rayyone/go-core/filestorage/contracts"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type LocalDisk struct {
	root string
	url  string
}

// NewLocalDisk creates a new local disk storage
func NewLocalDisk(config config.LocalConfig) *LocalDisk {
	return &LocalDisk{
		root: config.RootPath,
		url:  strings.TrimSuffix(config.BaseURL, "/"),
	}
}
func (ld *LocalDisk) Put(path string, content io.Reader, options ...contracts.Option) (*string, error) {
	fullPath := filepath.Join(ld.root, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	return &fullPath, err
}

func (ld *LocalDisk) PutFile(path string, filePath string, options ...contracts.Option) (*string, error) {
	sourceFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer sourceFile.Close()

	return ld.Put(path, sourceFile, options...)
}

func (ld *LocalDisk) Get(path string) ([]byte, error) {
	fullPath := filepath.Join(ld.root, path)
	return os.ReadFile(fullPath)
}

func (ld *LocalDisk) GetStream(path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(ld.root, path)
	return os.Open(fullPath)
}

func (ld *LocalDisk) URL(path string) string {
	return fmt.Sprintf("%s/%s", ld.url, strings.TrimPrefix(path, "/"))
}

func (ld *LocalDisk) SignedURL(path string, expiration time.Duration) (string, error) {
	return ld.URL(path), nil
}

func (ld *LocalDisk) Download(path string, localPath string) error {
	content, err := ld.Get(path)
	if err != nil {
		return err
	}
	return os.WriteFile(localPath, content, 0644)
}

func (ld *LocalDisk) Delete(path string) error {
	fullPath := filepath.Join(ld.root, path)
	return os.Remove(fullPath)
}

func (ld *LocalDisk) Exists(path string) bool {
	fullPath := filepath.Join(ld.root, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

func (ld *LocalDisk) Copy(from, to string) error {
	content, err := ld.Get(from)
	if err != nil {
		return err
	}
	_, err = ld.Put(to, bytes.NewReader(content))
	return err
}

func (ld *LocalDisk) Move(from, to string) error {
	if err := ld.Copy(from, to); err != nil {
		return err
	}
	return ld.Delete(from)
}

func (ld *LocalDisk) Size(path string) (int64, error) {
	fullPath := filepath.Join(ld.root, path)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}
