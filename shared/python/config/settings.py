import os
from pydantic_settings import BaseSettings
from pydantic import Field, validator
from typing import Optional


class InfrastructureSettings(BaseSettings):
    redis_url: str = Field(..., alias="REDIS_URL")
    rabbitmq_url: str = Field(..., alias="RABBITMQ_URL")
    minio_endpoint: str = Field(..., alias="MINIO_ENDPOINT")
    minio_access_key: str = Field(..., alias="MINIO_ACCESS_KEY")
    minio_secret_key: str = Field(..., alias="MINIO_SECRET_KEY")
    minio_bucket_3d: str = Field(default="daedalus-3d-assets", alias="MINIO_BUCKET_3D")
    minio_bucket_exports: str = Field(default="daedalus-exports", alias="MINIO_BUCKET_EXPORTS")
    minio_use_ssl: bool = Field(default=False, alias="MINIO_USE_SSL")
    database_url: str = Field(..., alias="DATABASE_URL")
    jaeger_host: str = Field(default="jaeger:4317", alias="JAEGER_HOST")

    class Config:
        env_file = ".env"
        case_sensitive = False
        populate_by_name = True


class AuthSettings(BaseSettings):
    jwt_secret: str = Field(..., alias="JWT_SECRET")
    jwt_algorithm: str = Field(default="HS256", alias="JWT_ALGORITHM")
    jwt_access_expire_minutes: int = Field(default=30, alias="JWT_ACCESS_EXPIRE_MINUTES")
    jwt_refresh_expire_days: int = Field(default=7, alias="JWT_REFRESH_EXPIRE_DAYS")
    oauth2_google_client_id: Optional[str] = Field(None, alias="OAUTH2_GOOGLE_CLIENT_ID")
    oauth2_google_client_secret: Optional[str] = Field(None, alias="OAUTH2_GOOGLE_CLIENT_SECRET")
    oauth2_github_client_id: Optional[str] = Field(None, alias="OAUTH2_GITHUB_CLIENT_ID")
    oauth2_github_client_secret: Optional[str] = Field(None, alias="OAUTH2_GITHUB_CLIENT_SECRET")
    traefik_forward_auth_secret: Optional[str] = Field(None, alias="TRAEFIK_FORWARD_AUTH_SECRET")

    class Config:
        env_file = ".env"
        case_sensitive = False
        populate_by_name = True


class PaymentSettings(BaseSettings):
    paystack_secret_key: Optional[str] = Field(None, alias="PAYSTACK_SECRET_KEY")
    paystack_public_key: Optional[str] = Field(None, alias="PAYSTACK_PUBLIC_KEY")
    stripe_secret_key: Optional[str] = Field(None, alias="STRIPE_SECRET_KEY")
    stripe_webhook_secret: Optional[str] = Field(None, alias="STRIPE_WEBHOOK_SECRET")
    payment_currency_default: str = Field(default="XAF", alias="PAYMENT_CURRENCY_DEFAULT")

    class Config:
        env_file = ".env"
        case_sensitive = False
        populate_by_name = True


class NotificationSettings(BaseSettings):
    smtp_host: str = Field(default="smtp.mailgun.org", alias="SMTP_HOST")
    smtp_port: int = Field(default=587, alias="SMTP_PORT")
    smtp_user: str = Field(..., alias="SMTP_USER")
    smtp_password: str = Field(..., alias="SMTP_PASSWORD")
    smtp_from: str = Field(default="no-reply@daedalus.io", alias="SMTP_FROM")
    notification_queue: str = Field(default="daedalus.notifications", alias="NOTIFICATION_QUEUE")

    class Config:
        env_file = ".env"
        case_sensitive = False
        populate_by_name = True


class AISettings(BaseSettings):
    openai_api_key: Optional[str] = Field(None, alias="OPENAI_API_KEY")
    anthropic_api_key: Optional[str] = Field(None, alias="ANTHROPIC_API_KEY")
    google_genai_api_key: Optional[str] = Field(None, alias="GOOGLE_GENAI_API_KEY")
    meshy_api_key: Optional[str] = Field(None, alias="MESHY_API_KEY")
    agent_model: str = Field(default="gpt-4o", alias="AGENT_MODEL")

    class Config:
        env_file = ".env"
        case_sensitive = False
        populate_by_name = True


class ServiceSettings(BaseSettings):
    service_name: str = Field(..., alias="SERVICE_NAME")
    environment: str = Field(default="development", alias="ENVIRONMENT")
    log_level: str = Field(default="INFO", alias="LOG_LEVEL")

    class Config:
        env_file = ".env"
        case_sensitive = False
        populate_by_name = True


class Settings(BaseSettings):
    infrastructure: InfrastructureSettings
    auth: AuthSettings
    payment: PaymentSettings
    notification: NotificationSettings
    ai: AISettings
    service: ServiceSettings

    class Config:
        env_file = ".env"
        case_sensitive = False
        populate_by_name = True

    def __init__(self, **data):
        super().__init__(**data)
        self.infrastructure = InfrastructureSettings()
        self.auth = AuthSettings()
        self.payment = PaymentSettings()
        self.notification = NotificationSettings()
        self.ai = AISettings()
        self.service = ServiceSettings()
