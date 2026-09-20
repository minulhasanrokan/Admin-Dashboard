package storage

import (
	appconfig "admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/storage/local"
	s3storage "admin-dashboard/cmd/internal/storage/s3"

	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func New(cfg *appconfig.Config) Storage {

	switch cfg.FileSystem {

	case "local":
		return local.New(
			cfg.FileSystemLocation,
			cfg.AppBaseUrl+"/file",
			cfg.FileBackupLocation,
		)

	case "aws-s3":
		client := s3.New(s3.Options{
			Region: cfg.FileSystemAWSRegion,
			Credentials: credentials.NewStaticCredentialsProvider(
				cfg.AWSAccessKey,
				cfg.AWSSecretKey,
				"",
			),
		})

		return s3storage.New(
			client,
			cfg.FileSystemS3Bucket,
			cfg.FileSystemAWSRegion,
		)

	default:
		return local.New(
			cfg.FileSystemLocation,
			cfg.AppBaseUrl+"/file",
			cfg.FileBackupLocation,
		)
	}
}
