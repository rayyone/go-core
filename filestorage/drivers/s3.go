package drivers

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	appconfig "github.com/rayyone/go-core/filestorage/config"
	"github.com/rayyone/go-core/filestorage/contracts"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type S3Disk struct {
	client        *s3.Client
	uploader      *manager.Uploader
	downloader    *manager.Downloader
	bucket        string
	region        string
	baseURL       string
	cloudFrontURL string
}

// NewS3Disk creates a new S3 disk storage
func NewS3Disk(cfg appconfig.S3Config) (*S3Disk, error) {
	var optFn config.LoadOptionsFunc
	if cfg.Profile != "" {
		optFn = config.WithSharedConfigProfile(cfg.Profile)
	}
	if cfg.AccessKey != "" && cfg.SecretKey != "" {
		creds := credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, cfg.SessionToken)
		optFn = config.WithCredentialsProvider(creds)
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		optFn,
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	uploader := manager.NewUploader(client)
	downloader := manager.NewDownloader(client)

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", cfg.Bucket, cfg.Region)
	}
	return &S3Disk{
		client:        client,
		uploader:      uploader,
		downloader:    downloader,
		bucket:        cfg.Bucket,
		region:        cfg.Region,
		baseURL:       strings.TrimSuffix(baseURL, "/"),
		cloudFrontURL: strings.TrimSuffix(cfg.CloudFrontURL, "/"),
	}, nil
}
func (s3d *S3Disk) Put(path string, content io.Reader, options ...contracts.Option) (*string, error) {
	opts := &contracts.PutOptions{}
	for _, option := range options {
		option(opts)
	}

	bucket := s3d.bucket
	if opts.Bucket != "" {
		bucket = opts.Bucket
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
		Body:   content,
	}

	// Set options
	if opts.ACL != "" {
		input.ACL = types.ObjectCannedACL(opts.ACL)
	}
	if opts.MimeType != "" {
		input.ContentType = aws.String(opts.MimeType)
	} else {
		if contentType := mime.TypeByExtension(filepath.Ext(path)); contentType != "" {
			input.ContentType = aws.String(contentType)
		}
	}
	if opts.ContentDisposition != "" {
		input.ContentDisposition = aws.String(opts.ContentDisposition)
	}
	if opts.CacheControl != "" {
		input.CacheControl = aws.String(opts.CacheControl)
	}
	if opts.ContentEncoding != "" {
		input.ContentEncoding = aws.String(opts.ContentEncoding)
	}
	if opts.ContentLanguage != "" {
		input.ContentLanguage = aws.String(opts.ContentLanguage)
	}
	if opts.Expires != nil {
		input.Expires = opts.Expires
	}
	if opts.Metadata != nil {
		input.Metadata = opts.Metadata
	}
	if opts.ServerSideEncryption != "" {
		input.ServerSideEncryption = types.ServerSideEncryption(opts.ServerSideEncryption)
	}
	if opts.StorageClass != "" {
		input.StorageClass = types.StorageClass(opts.StorageClass)
	}

	_, err := s3d.uploader.Upload(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	return &path, err
}
func (s3d *S3Disk) PutFile(path string, filePath string, options ...contracts.Option) (*string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return s3d.Put(path, file, options...)
}

func (s3d *S3Disk) Get(path string) ([]byte, error) {
	result, err := s3d.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s3d.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (s3d *S3Disk) GetStream(path string) (io.ReadCloser, error) {
	result, err := s3d.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s3d.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (s3d *S3Disk) URL(path string) string {
	if s3d.cloudFrontURL != "" {
		return fmt.Sprintf("%s/%s", s3d.cloudFrontURL, strings.TrimPrefix(path, "/"))
	}
	return fmt.Sprintf("%s/%s", s3d.baseURL, strings.TrimPrefix(path, "/"))
}

func (s3d *S3Disk) SignedURL(path string, expiration time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s3d.client)

	request, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s3d.bucket),
		Key:    aws.String(path),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func (s3d *S3Disk) Download(path string, localPath string) error {
	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s3d.downloader.Download(context.TODO(), file, &s3.GetObjectInput{
		Bucket: aws.String(s3d.bucket),
		Key:    aws.String(path),
	})

	return err
}

func (s3d *S3Disk) Delete(path string) error {
	_, err := s3d.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s3d.bucket),
		Key:    aws.String(path),
	})
	return err
}

func (s3d *S3Disk) Exists(path string) bool {
	_, err := s3d.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s3d.bucket),
		Key:    aws.String(path),
	})
	return err == nil
}

func (s3d *S3Disk) Copy(from, to string) error {
	copySource := fmt.Sprintf("%s/%s", s3d.bucket, from)
	_, err := s3d.client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(s3d.bucket),
		Key:        aws.String(to),
		CopySource: aws.String(url.PathEscape(copySource)),
	})
	return err
}

func (s3d *S3Disk) Move(from, to string) error {
	if err := s3d.Copy(from, to); err != nil {
		return err
	}
	return s3d.Delete(from)
}

func (s3d *S3Disk) Size(path string) (int64, error) {
	result, err := s3d.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s3d.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return 0, err
	}
	return aws.ToInt64(result.ContentLength), nil
}
