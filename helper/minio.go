package helper

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Minio(endpoint string, accessKey string, secretAccessKey string, secure bool) (minioClient *minio.Client, err error) {
	minioClient, err = minio.New(
		endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
			Secure: secure,
		},
	)

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
