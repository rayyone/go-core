package storage

import (
	"io"

	stgoption "github.com/rayyone/go-core/storage/option"
	"github.com/rayyone/go-core/storage/s3"
)

type Configuration struct {
	Default string
	S3      s3.Configuration
	Folders map[string]string
}

// Driver is an storage driver interface
type Driver interface {
	Store(file io.Reader, filename string, filePath string, opts ...stgoption.OptionFunc) (location *string, err error)
	Delete(fullPath string) error
}

// Storage Storage
type Storage struct {
	driver Driver
}

// Store Store file based on driver
func (s *Storage) Store(file io.Reader, filename string, filePath string, opts ...stgoption.OptionFunc) (location *string, err error) {
	return s.driver.Store(file, filename, filePath, opts...)
}

// Delete Delete file based on driver
func (s *Storage) Delete(fullPath string) error {
	return s.driver.Delete(fullPath)
}

// Driver Set driver
func (s *Storage) Driver(driver Driver) *Storage {
	return &Storage{driver: driver}
}

// NewStorage creates new storage
func NewStorage(driver Driver) *Storage {
	return &Storage{
		driver: driver,
	}
}

// NewDefaultDriver create default driver
func NewDefaultDriver(config Configuration) Driver {
	switch config.Default {
	case "s3":
		return s3.NewStorage(config.S3)
	default:
		panic("Default Storage Driver is not found")
	}
}
