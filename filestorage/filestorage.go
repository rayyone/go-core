package filestorage

import (
	"github.com/rayyone/go-core/filestorage/config"
	"github.com/rayyone/go-core/filestorage/contracts"
	"github.com/rayyone/go-core/filestorage/drivers"
	"github.com/rayyone/go-core/filestorage/manager"
	"github.com/rayyone/go-core/filestorage/option"
	"io"
)

var defaultManager = manager.NewFileManager()

type (
	Storage     = contracts.Storage
	Option      = contracts.Option
	PutOptions  = contracts.PutOptions
	S3Config    = config.S3Config
	LocalConfig = config.LocalConfig
)

var (
	WithACL                  = option.WithACL
	WithMimeType             = option.WithMimeType
	WithContentDisposition   = option.WithContentDisposition
	WithBucket               = option.WithBucket
	WithCacheControl         = option.WithCacheControl
	WithContentEncoding      = option.WithContentEncoding
	WithContentLanguage      = option.WithContentLanguage
	WithExpires              = option.WithExpires
	WithMetadata             = option.WithMetadata
	WithServerSideEncryption = option.WithServerSideEncryption
	WithStorageClass         = option.WithStorageClass
)

func NewLocalDisk(cfg LocalConfig) (Storage, error) {
	return drivers.NewLocalDisk(cfg), nil
}

func NewS3Disk(cfg S3Config) (Storage, error) {
	return drivers.NewS3Disk(cfg)
}

func AddDisk(name string, storage Storage) {
	defaultManager.AddDisk(name, storage)
}

func SetDefaultDisk(name string) {
	defaultManager.SetDefaultDisk(name)
}

func Disk(name string) Storage {
	return defaultManager.Disk(name)
}

func Default() Storage {
	return defaultManager.Default()
}

func Put(path string, content io.Reader, options ...Option) (*string, error) {
	return defaultManager.Put(path, content, options...)
}

func PutFile(path string, filePath string, options ...Option) (*string, error) {
	return defaultManager.PutFile(path, filePath, options...)
}

func Get(path string) ([]byte, error) {
	return defaultManager.Get(path)
}

func URL(path string) string {
	return defaultManager.URL(path)
}

func Delete(path string) error {
	return defaultManager.Delete(path)
}

func Exists(path string) bool {
	return defaultManager.Exists(path)
}

func Initialize(confs []config.FileStorageConfig, defaultDisk string) *manager.FileManager {
	for _, conf := range confs {
		var (
			storage Storage
			err     error
		)
		switch conf.Driver {
		case manager.S3Driver:
			storage, err = NewS3Disk(conf.S3)
			if err != nil {
				panic(err)
			}
			break
		case manager.LocalDriver:
			storage, err = NewLocalDisk(conf.Local)
			if err != nil {
				panic(err)
			}
			break
		default:
			panic("Driver is not support")
		}
		AddDisk(conf.Name, storage)
	}
	if defaultDisk != "" {
		SetDefaultDisk(defaultDisk)
	}
	return defaultManager
}
