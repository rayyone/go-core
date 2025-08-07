package option

import (
	"github.com/rayyone/go-core/filestorage/contracts"
	"time"
)

const (
	ACLPrivate    contracts.ACLType = "private"
	ACLPublicRead contracts.ACLType = "public-read"
)

func WithACL(acl contracts.ACLType) contracts.Option {
	return func(o *contracts.PutOptions) {
		o.ACL = acl
	}
}

func WithMimeType(mimeType string) contracts.Option {
	return func(o *contracts.PutOptions) {
		o.MimeType = mimeType
	}
}

func WithContentDisposition(disposition string) contracts.Option {
	return func(o *contracts.PutOptions) {
		o.ContentDisposition = disposition
	}
}

func WithBucket(bucket string) contracts.Option {
	return func(o *contracts.PutOptions) {
		o.Bucket = bucket
	}
}
func WithCacheControl(cacheControl string) contracts.Option {
	return func(o *contracts.PutOptions) {
		o.CacheControl = cacheControl
	}
}

func WithContentEncoding(encoding string) contracts.Option {
	return func(o *contracts.PutOptions) {
		o.ContentEncoding = encoding
	}
}

func WithContentLanguage(language string) contracts.Option {
	return func(o *contracts.PutOptions) {
		o.ContentLanguage = language
	}
}

func WithExpires(expires time.Time) contracts.Option {
	return func(o *contracts.PutOptions) {
		o.Expires = &expires
	}
}

func WithMetadata(metadata map[string]string) contracts.Option {
	return func(o *contracts.PutOptions) {
		o.Metadata = metadata
	}
}

func WithServerSideEncryption(encryption string) contracts.Option {
	return func(o *contracts.PutOptions) {
		o.ServerSideEncryption = encryption
	}
}

func WithStorageClass(class string) contracts.Option {
	return func(o *contracts.PutOptions) {
		o.StorageClass = class
	}
}
