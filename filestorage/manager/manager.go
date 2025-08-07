package manager

import (
	"fmt"
	"github.com/rayyone/go-core/filestorage/contracts"
	"io"
)

const (
	S3Driver    contracts.DriverType = "s3"
	LocalDriver contracts.DriverType = "local"
)

type FileManager struct {
	disks       map[string]contracts.Storage
	defaultDisk string
}

func NewFileManager() *FileManager {
	return &FileManager{
		disks: make(map[string]contracts.Storage),
	}
}

func (fm *FileManager) AddDisk(name string, storage contracts.Storage) {
	fm.disks[name] = storage
}

func (fm *FileManager) SetDefaultDisk(name string) {
	fm.defaultDisk = name
}

func (fm *FileManager) Disk(name string) contracts.Storage {
	if storage, exists := fm.disks[name]; exists {
		return storage
	}
	panic(fmt.Sprintf("Disk '%s' not found", name))
}

func (fm *FileManager) Default() contracts.Storage {
	if fm.defaultDisk == "" {
		panic("No default disk set")
	}
	return fm.Disk(fm.defaultDisk)
}

func (fm *FileManager) Put(path string, content io.Reader, options ...contracts.Option) (*string, error) {
	return fm.Default().Put(path, content, options...)
}
func (fm *FileManager) PutFile(path string, filePath string, options ...contracts.Option) (*string, error) {
	return fm.Default().PutFile(path, filePath, options...)
}

func (fm *FileManager) Get(path string) ([]byte, error) {
	return fm.Default().Get(path)
}

func (fm *FileManager) URL(path string) string {
	return fm.Default().URL(path)
}

func (fm *FileManager) Delete(path string) error {
	return fm.Default().Delete(path)
}

func (fm *FileManager) Exists(path string) bool {
	return fm.Default().Exists(path)
}
