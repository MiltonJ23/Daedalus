package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name     string
		setup    func()
		teardown func()
		wantErr  bool
	}{
		{
			name: "valid config loads successfully",
			setup: func() {
				os.Setenv("REDIS_URL", "redis://localhost:6379")
				os.Setenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672")
				os.Setenv("SERVICE_NAME", "test-service")
				os.Setenv("JWT_SECRET", "test-secret")
				os.Setenv("JWT_ALGORITHM", "HS256")
			},
			teardown: func() {
				os.Unsetenv("REDIS_URL")
				os.Unsetenv("RABBITMQ_URL")
				os.Unsetenv("SERVICE_NAME")
				os.Unsetenv("JWT_SECRET")
				os.Unsetenv("JWT_ALGORITHM")
			},
			wantErr: false,
		},
		{
			name: "missing required REDIS_URL fails",
			setup: func() {
				os.Unsetenv("REDIS_URL")
				os.Setenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672")
				os.Setenv("SERVICE_NAME", "test-service")
				os.Setenv("JWT_SECRET", "test-secret")
			},
			teardown: func() {
				os.Unsetenv("RABBITMQ_URL")
				os.Unsetenv("SERVICE_NAME")
				os.Unsetenv("JWT_SECRET")
			},
			wantErr: true,
		},
		{
			name: "missing JWT_SECRET fails",
			setup: func() {
				os.Setenv("REDIS_URL", "redis://localhost:6379")
				os.Setenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672")
				os.Setenv("SERVICE_NAME", "test-service")
				os.Unsetenv("JWT_SECRET")
			},
			teardown: func() {
				os.Unsetenv("REDIS_URL")
				os.Unsetenv("RABBITMQ_URL")
				os.Unsetenv("SERVICE_NAME")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			defer tt.teardown()

			cfg, err := Load()

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && cfg == nil {
				t.Error("Load() returned nil config")
			}
		})
	}
}

func TestConfigDefaults(t *testing.T) {
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672")
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("JWT_SECRET", "test-secret")
	defer func() {
		os.Unsetenv("REDIS_URL")
		os.Unsetenv("RABBITMQ_URL")
		os.Unsetenv("SERVICE_NAME")
		os.Unsetenv("JWT_SECRET")
	}()

	cfg, _ := Load()

	if cfg.Auth.JWTAccessExpireMinutes != 30 {
		t.Errorf("default JWT_ACCESS_EXPIRE_MINUTES = %d, want 30", cfg.Auth.JWTAccessExpireMinutes)
	}

	if cfg.Auth.JWTRefreshExpireDays != 7 {
		t.Errorf("default JWT_REFRESH_EXPIRE_DAYS = %d, want 7", cfg.Auth.JWTRefreshExpireDays)
	}

	if cfg.Notification.SMTPPort != 587 {
		t.Errorf("default SMTP_PORT = %d, want 587", cfg.Notification.SMTPPort)
	}
}

func TestConfigTimeouts(t *testing.T) {
	cfg := &InfrastructureConfig{}

	if cfg.RedisConnectTimeout() == 0 {
		t.Error("RedisConnectTimeout should return non-zero duration")
	}

	if cfg.RabbitMQConnectTimeout() == 0 {
		t.Error("RabbitMQConnectTimeout should return non-zero duration")
	}
}
