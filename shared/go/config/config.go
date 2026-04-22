package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Infrastructure InfrastructureConfig
	Auth           AuthConfig
	Payment        PaymentConfig
	Notification   NotificationConfig
	AI             AIConfig
	Service        ServiceConfig
}

type InfrastructureConfig struct {
	RedisURL           string
	RabbitMQURL        string
	MinIOEndpoint      string
	MinIOAccessKey     string
	MinIOSecretKey     string
	MinIOBucket3D      string
	MinIOBucketExports string
	MinIOUseSSL        bool
	DatabaseURL        string
	JaegerHost         string
}

type AuthConfig struct {
	JWTSecret                string
	JWTAlgorithm             string
	JWTAccessExpireMinutes   int
	JWTRefreshExpireDays     int
	OAuth2GoogleClientID     string
	OAuth2GoogleClientSecret string
	OAuth2GitHubClientID     string
	OAuth2GitHubClientSecret string
	TraefikForwardAuthSecret string
}

type PaymentConfig struct {
	PaystackSecretKey      string
	PaystackPublicKey      string
	StripeSecretKey        string
	StripeWebhookSecret    string
	PaymentCurrencyDefault string
}

type NotificationConfig struct {
	SMTPHost          string
	SMTPPort          int
	SMTPUser          string
	SMTPPassword      string
	SMTPFrom          string
	NotificationQueue string
}

type AIConfig struct {
	OpenAIAPIKey      string
	AnthropicAPIKey   string
	GoogleGenAIAPIKey string
	MeshyAPIKey       string
	AgentModel        string
}

type ServiceConfig struct {
	ServiceName string
	Environment string
	LogLevel    string
}

func Load() (*Config, error) {
	viper.SetEnvPrefix("DAEDALUS")
	viper.AutomaticEnv()

	viper.SetDefault("REDIS_URL", os.Getenv("REDIS_URL"))
	viper.SetDefault("RABBITMQ_URL", os.Getenv("RABBITMQ_URL"))
	viper.SetDefault("MINIO_ENDPOINT", os.Getenv("MINIO_ENDPOINT"))
	viper.SetDefault("MINIO_ACCESS_KEY", os.Getenv("MINIO_ACCESS_KEY"))
	viper.SetDefault("MINIO_SECRET_KEY", os.Getenv("MINIO_SECRET_KEY"))
	viper.SetDefault("MINIO_BUCKET_3D", os.Getenv("MINIO_BUCKET_3D"))
	viper.SetDefault("MINIO_BUCKET_EXPORTS", os.Getenv("MINIO_BUCKET_EXPORTS"))
	viper.SetDefault("MINIO_USE_SSL", os.Getenv("MINIO_USE_SSL") == "true")
	viper.SetDefault("DATABASE_URL", os.Getenv("DATABASE_URL"))
	viper.SetDefault("JAEGER_HOST", os.Getenv("JAEGER_HOST"))

	viper.SetDefault("JWT_SECRET", os.Getenv("JWT_SECRET"))
	viper.SetDefault("JWT_ALGORITHM", os.Getenv("JWT_ALGORITHM"))
	viper.SetDefault("JWT_ACCESS_EXPIRE_MINUTES", getenvInt("JWT_ACCESS_EXPIRE_MINUTES", 30))
	viper.SetDefault("JWT_REFRESH_EXPIRE_DAYS", getenvInt("JWT_REFRESH_EXPIRE_DAYS", 7))
	viper.SetDefault("OAUTH2_GOOGLE_CLIENT_ID", os.Getenv("OAUTH2_GOOGLE_CLIENT_ID"))
	viper.SetDefault("OAUTH2_GOOGLE_CLIENT_SECRET", os.Getenv("OAUTH2_GOOGLE_CLIENT_SECRET"))
	viper.SetDefault("OAUTH2_GITHUB_CLIENT_ID", os.Getenv("OAUTH2_GITHUB_CLIENT_ID"))
	viper.SetDefault("OAUTH2_GITHUB_CLIENT_SECRET", os.Getenv("OAUTH2_GITHUB_CLIENT_SECRET"))
	viper.SetDefault("TRAEFIK_FORWARD_AUTH_SECRET", os.Getenv("TRAEFIK_FORWARD_AUTH_SECRET"))

	viper.SetDefault("PAYSTACK_SECRET_KEY", os.Getenv("PAYSTACK_SECRET_KEY"))
	viper.SetDefault("PAYSTACK_PUBLIC_KEY", os.Getenv("PAYSTACK_PUBLIC_KEY"))
	viper.SetDefault("STRIPE_SECRET_KEY", os.Getenv("STRIPE_SECRET_KEY"))
	viper.SetDefault("STRIPE_WEBHOOK_SECRET", os.Getenv("STRIPE_WEBHOOK_SECRET"))
	viper.SetDefault("PAYMENT_CURRENCY_DEFAULT", os.Getenv("PAYMENT_CURRENCY_DEFAULT"))

	viper.SetDefault("SMTP_HOST", os.Getenv("SMTP_HOST"))
	viper.SetDefault("SMTP_PORT", getenvInt("SMTP_PORT", 587))
	viper.SetDefault("SMTP_USER", os.Getenv("SMTP_USER"))
	viper.SetDefault("SMTP_PASSWORD", os.Getenv("SMTP_PASSWORD"))
	viper.SetDefault("SMTP_FROM", os.Getenv("SMTP_FROM"))
	viper.SetDefault("NOTIFICATION_QUEUE", os.Getenv("NOTIFICATION_QUEUE"))

	viper.SetDefault("OPENAI_API_KEY", os.Getenv("OPENAI_API_KEY"))
	viper.SetDefault("ANTHROPIC_API_KEY", os.Getenv("ANTHROPIC_API_KEY"))
	viper.SetDefault("GOOGLE_GENAI_API_KEY", os.Getenv("GOOGLE_GENAI_API_KEY"))
	viper.SetDefault("MESHY_API_KEY", os.Getenv("MESHY_API_KEY"))
	viper.SetDefault("AGENT_MODEL", os.Getenv("AGENT_MODEL"))

	viper.SetDefault("SERVICE_NAME", os.Getenv("SERVICE_NAME"))
	viper.SetDefault("ENVIRONMENT", os.Getenv("ENVIRONMENT"))
	viper.SetDefault("LOG_LEVEL", os.Getenv("LOG_LEVEL"))

	cfg := &Config{
		Infrastructure: InfrastructureConfig{
			RedisURL:           viper.GetString("REDIS_URL"),
			RabbitMQURL:        viper.GetString("RABBITMQ_URL"),
			MinIOEndpoint:      viper.GetString("MINIO_ENDPOINT"),
			MinIOAccessKey:     viper.GetString("MINIO_ACCESS_KEY"),
			MinIOSecretKey:     viper.GetString("MINIO_SECRET_KEY"),
			MinIOBucket3D:      viper.GetString("MINIO_BUCKET_3D"),
			MinIOBucketExports: viper.GetString("MINIO_BUCKET_EXPORTS"),
			MinIOUseSSL:        viper.GetBool("MINIO_USE_SSL"),
			DatabaseURL:        viper.GetString("DATABASE_URL"),
			JaegerHost:         viper.GetString("JAEGER_HOST"),
		},
		Auth: AuthConfig{
			JWTSecret:                viper.GetString("JWT_SECRET"),
			JWTAlgorithm:             viper.GetString("JWT_ALGORITHM"),
			JWTAccessExpireMinutes:   viper.GetInt("JWT_ACCESS_EXPIRE_MINUTES"),
			JWTRefreshExpireDays:     viper.GetInt("JWT_REFRESH_EXPIRE_DAYS"),
			OAuth2GoogleClientID:     viper.GetString("OAUTH2_GOOGLE_CLIENT_ID"),
			OAuth2GoogleClientSecret: viper.GetString("OAUTH2_GOOGLE_CLIENT_SECRET"),
			OAuth2GitHubClientID:     viper.GetString("OAUTH2_GITHUB_CLIENT_ID"),
			OAuth2GitHubClientSecret: viper.GetString("OAUTH2_GITHUB_CLIENT_SECRET"),
			TraefikForwardAuthSecret: viper.GetString("TRAEFIK_FORWARD_AUTH_SECRET"),
		},
		Payment: PaymentConfig{
			PaystackSecretKey:      viper.GetString("PAYSTACK_SECRET_KEY"),
			PaystackPublicKey:      viper.GetString("PAYSTACK_PUBLIC_KEY"),
			StripeSecretKey:        viper.GetString("STRIPE_SECRET_KEY"),
			StripeWebhookSecret:    viper.GetString("STRIPE_WEBHOOK_SECRET"),
			PaymentCurrencyDefault: viper.GetString("PAYMENT_CURRENCY_DEFAULT"),
		},
		Notification: NotificationConfig{
			SMTPHost:          viper.GetString("SMTP_HOST"),
			SMTPPort:          viper.GetInt("SMTP_PORT"),
			SMTPUser:          viper.GetString("SMTP_USER"),
			SMTPPassword:      viper.GetString("SMTP_PASSWORD"),
			SMTPFrom:          viper.GetString("SMTP_FROM"),
			NotificationQueue: viper.GetString("NOTIFICATION_QUEUE"),
		},
		AI: AIConfig{
			OpenAIAPIKey:      viper.GetString("OPENAI_API_KEY"),
			AnthropicAPIKey:   viper.GetString("ANTHROPIC_API_KEY"),
			GoogleGenAIAPIKey: viper.GetString("GOOGLE_GENAI_API_KEY"),
			MeshyAPIKey:       viper.GetString("MESHY_API_KEY"),
			AgentModel:        viper.GetString("AGENT_MODEL"),
		},
		Service: ServiceConfig{
			ServiceName: viper.GetString("SERVICE_NAME"),
			Environment: viper.GetString("ENVIRONMENT"),
			LogLevel:    viper.GetString("LOG_LEVEL"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Infrastructure.RedisURL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	if c.Infrastructure.RabbitMQURL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}
	if c.Service.ServiceName == "" {
		return fmt.Errorf("SERVICE_NAME is required")
	}
	if c.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

func getenvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func (c *InfrastructureConfig) RedisConnectTimeout() time.Duration {
	return 5 * time.Second
}

func (c *InfrastructureConfig) RabbitMQConnectTimeout() time.Duration {
	return 5 * time.Second
}
