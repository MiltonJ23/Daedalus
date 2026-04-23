package storage

import (
	"context"
	"testing"
)

func TestNewMinIOClient(t *testing.T) {
	cfg := MinIOConfig{
		Endpoint:      "localhost:9000",
		AccessKey:     "minioadmin",
		SecretKey:     "minioadmin",
		UseSSL:        false,
		Bucket3D:      "daedalus-3d-assets",
		BucketExports: "daedalus-exports",
	}

	client, err := NewMinIOClient(cfg)
	if err != nil {
		t.Skipf("MinIO not available: %v", err)
	}

	if client.client == nil {
		t.Error("minio client should not be nil")
	}
}

func TestHealthCheck(t *testing.T) {
	cfg := MinIOConfig{
		Endpoint:      "localhost:9000",
		AccessKey:     "minioadmin",
		SecretKey:     "minioadmin",
		UseSSL:        false,
		Bucket3D:      "daedalus-3d-assets",
		BucketExports: "daedalus-exports",
	}

	client, err := NewMinIOClient(cfg)
	if err != nil {
		t.Skipf("MinIO not available: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = client.HealthCheck(ctx)
	if err != nil {
		t.Errorf("HealthCheck failed: %v", err)
	}
}

func TestGetContentType(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"model.glb", "model/gltf-binary"},
		{"model.gltf", "model/gltf+json"},
		{"export.pdf", "application/pdf"},
		{"data.csv", "text/csv"},
		{"config.json", "application/json"},
		{"unknown.bin", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := getContentType(tt.filename)
			if got != tt.expected {
				t.Errorf("getContentType = %s, want %s", got, tt.expected)
			}
		})
	}
}
