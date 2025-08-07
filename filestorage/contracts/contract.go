package contracts

import (
	"io"
	"time"
)

type Option func(*PutOptions)

type Storage interface {
	Put(path string, content io.Reader, options ...Option) (*string, error)
	PutFile(path string, filePath string, options ...Option) (*string, error)
	Get(path string) ([]byte, error)
	GetStream(path string) (io.ReadCloser, error)
	URL(path string) string
	SignedURL(path string, expiration time.Duration) (string, error)
	Download(path string, localPath string) error
	Delete(path string) error
	Exists(path string) bool
	Copy(from, to string) error
	Move(from, to string) error
	Size(path string) (int64, error)
}
type ACLType string
type DiskType string
type DriverType string
type PutOptions struct {
	ACL                  ACLType
	MimeType             string
	ContentDisposition   string
	Bucket               string
	CacheControl         string
	ContentEncoding      string
	ContentLanguage      string
	Expires              *time.Time
	Metadata             map[string]string
	ServerSideEncryption string
	StorageClass         string
}
