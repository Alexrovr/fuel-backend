package storage

import (
	"fmt"
	"os"
	"strings"
)

type MinioMedia struct {
	Endpoint string
	Bucket   string
}

func NewMinioMedia() *MinioMedia {
	endpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:9000"
	}
	bucket := os.Getenv("MINIO_BUCKET")
	if bucket == "" {
		bucket = "heat-fuel-media"
	}
	return &MinioMedia{
		Endpoint: strings.TrimRight(endpoint, "/"),
		Bucket:   bucket,
	}
}

// ObjectURL возвращает абсолютный адрес объекта в бакете Minio.
func (m *MinioMedia) ObjectURL(objectKey string) string {
	if objectKey == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s", m.Endpoint, m.Bucket, objectKey)
}
