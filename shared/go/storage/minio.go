package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client *minio.Client
	config MinIOConfig
}

type MinIOConfig struct {
	Endpoint      string
	AccessKey     string
	SecretKey     string
	UseSSL        bool
	Bucket3D      string
	BucketExports string
}

func NewMinIOClient(cfg MinIOConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.BucketExists(ctx, cfg.Bucket3D); err != nil {
		return nil, fmt.Errorf("failed to verify bucket: %w", err)
	}

	return &MinIOClient{client: client, config: cfg}, nil
}

func (mc *MinIOClient) UploadFile(ctx context.Context, bucket string, objectName string, filePath string) error {
	_, err := mc.client.FPutObject(ctx, bucket, objectName, filePath, minio.PutObjectOptions{
		ContentType: getContentType(objectName),
	})
	return err
}

func (mc *MinIOClient) DownloadFile(ctx context.Context, bucket string, objectName string, filePath string) error {
	return mc.client.FGetObject(ctx, bucket, objectName, filePath, minio.GetObjectOptions{})
}

func (mc *MinIOClient) GetPresignedURL(ctx context.Context, bucket string, objectName string, expiration time.Duration) (string, error) {
	url, err := mc.client.PresignedGetObject(ctx, bucket, objectName, expiration, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

func (mc *MinIOClient) DeleteObject(ctx context.Context, bucket string, objectName string) error {
	return mc.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
}

func (mc *MinIOClient) ObjectExists(ctx context.Context, bucket string, objectName string) (bool, error) {
	_, err := mc.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (mc *MinIOClient) ListObjects(ctx context.Context, bucket string, prefix string) ([]string, error) {
	var objects []string
	for obj := range mc.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Prefix: prefix}) {
		if obj.Err != nil {
			return nil, obj.Err
		}
		objects = append(objects, obj.Key)
	}
	return objects, nil
}

func (mc *MinIOClient) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := mc.client.BucketExists(ctx, mc.config.Bucket3D)
	return err
}

func (mc *MinIOClient) Close() error {
	return nil
}

func getContentType(filename string) string {
	if len(filename) < 4 {
		return "application/octet-stream"
	}

	if strings.HasSuffix(filename, ".gltf") {
		return "model/gltf+json"
	}
	if strings.HasSuffix(filename, ".glb") {
		return "model/gltf-binary"
	}
	if strings.HasSuffix(filename, ".pdf") {
		return "application/pdf"
	}
	if strings.HasSuffix(filename, ".csv") {
		return "text/csv"
	}
	if strings.HasSuffix(filename, ".json") {
		return "application/json"
	}
	return "application/octet-stream"
}
