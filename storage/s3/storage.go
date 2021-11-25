package s3

import (
	"fmt"
	"io"
	"path"
	"time"

	"github.com/rayyone/go-core/errors"
	fileHelper "github.com/rayyone/go-core/helpers/file"
	stgoption "github.com/rayyone/go-core/storage/option"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3 S3 storage
type S3 struct {
	config   Configuration
	uploader *s3manager.Uploader
	service  *s3.S3
}

type Configuration struct {
	Bucket string
	Region string
	ACL    string
}

// Store Store the file
func (s *S3) Store(file io.Reader, filename string, filePath string, opts ...stgoption.OptionFunc) (location *string, err error) {
	options := stgoption.GetDefaultOptions()
	for _, o := range opts {
		o(&options)
	}

	mimeType := fileHelper.GetMimeType(filename)

	upParams := &s3manager.UploadInput{
		Bucket:             aws.String(s.config.Bucket),
		Key:                aws.String(path.Join(filePath, filename)),
		Body:               file,
		ACL:                aws.String(s.config.ACL),
		ContentDisposition: aws.String(options.ContentDisposition),
		ContentType:        aws.String(mimeType),
	}

	result, err := s.uploader.Upload(upParams)
	if err != nil {
		errMsg := fmt.Sprintf("Error: Cannot store file to S3. Error: %v", err)
		return nil, errors.BadRequest.New(errMsg)
	}

	return &result.Location, nil
}

// GetSignedUrl get signed url
func (s *S3) GetSignedUrl(key string, expireIn time.Duration, opts ...stgoption.OptionFunc) (url string, err error) {
	options := stgoption.GetDefaultOptions()
	for _, o := range opts {
		o(&options)
	}

	req, _ := s.service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})
	url, err = req.Presign(expireIn)
	if err != nil {
		return "", errors.BadRequest.Newf("Cannot get s3 signed url for key '%s'", key)
	}

	return url, nil
}

func (s *S3) GetPutSignedURL(key string, expireIn time.Duration, opts ...stgoption.OptionFunc) (url string, err error) {
	options := stgoption.GetDefaultOptions()
	for _, o := range opts {
		o(&options)
	}

	req, _ := s.service.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
		ACL:    aws.String(s.config.ACL),
	})
	url, err = req.Presign(expireIn)
	if err != nil {
		return "", errors.BadRequest.Newf("Cannot get s3 put signed url for key '%s'", key)
	}

	return url, nil
}

// NewStorage New S3 storage
func NewStorage(conf Configuration) *S3 {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(conf.Region),
	})
	if err != nil {
		panic("Cannot create new S3 session: " + err.Error())
	}

	service := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	return &S3{
		config:   conf,
		service:  service,
		uploader: uploader,
	}
}
