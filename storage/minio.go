package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

var (
	minioClient *minio.Client
)

func InitMinio(endpoint string) error {
	accessKeyID := "minioadmin"
	secretAccessKey := "xx123456"

	var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return err
	}
	return nil
}

func UploadFile(bucket string, filename string, ioReader io.Reader, size int64) error {
	_, err := minioClient.PutObject(context.Background(), bucket, filename, ioReader, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}
	return nil
}
