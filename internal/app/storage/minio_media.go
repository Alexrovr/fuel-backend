package storage

import (
	"fmt"
	"os"
	"strings"
)

// MinioMedia собирает публичные url объектов Minio по ключам,
// которые хранятся в коллекции топлив (поля ImageKey и VideoKey).
//
// В модели лежит только ключ на латинице (methane.jpg), адрес хранилища
// задаётся окружением, поэтому при переезде Minio коллекция не меняется.
type MinioMedia struct {
	Endpoint string // http://localhost:9000
	Bucket   string // heat-fuel-media
}

// NewMinioMedia читает настройки Minio из переменных окружения,
// подставляя значения docker-compose по умолчанию.
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
