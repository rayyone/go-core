package config

import "github.com/rayyone/go-core/filestorage/contracts"

type S3Config struct {
	Region        string
	Bucket        string
	BaseURL       string
	CloudFrontURL string
	Profile       string
	AccessKey     string
	SecretKey     string
	SessionToken  string
	Endpoint      string
}

type LocalConfig struct {
	RootPath string
	BaseURL  string
}

type FileStorageConfig struct {
	Name   string
	Driver contracts.DriverType
	S3     S3Config
	Local  LocalConfig
}
