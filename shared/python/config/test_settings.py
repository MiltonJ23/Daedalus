import pytest
import os
from config.settings import Settings, InfrastructureSettings


@pytest.fixture
def env_setup(monkeypatch):
    monkeypatch.setenv("REDIS_URL", "redis://localhost:6379/0")
    monkeypatch.setenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
    monkeypatch.setenv("MINIO_ENDPOINT", "localhost:9000")
    monkeypatch.setenv("MINIO_ACCESS_KEY", "minioadmin")
    monkeypatch.setenv("MINIO_SECRET_KEY", "minioadmin")
    monkeypatch.setenv("DATABASE_URL", "postgresql://localhost/daedalus")
    monkeypatch.setenv("JWT_SECRET", "test-secret-key")
    monkeypatch.setenv("SERVICE_NAME", "test-service")
    monkeypatch.setenv("SMTP_USER", "test@example.com")
    monkeypatch.setenv("SMTP_PASSWORD", "password")


def test_infrastructure_settings_valid(env_setup):
    settings = InfrastructureSettings()
    assert settings.redis_url == "redis://localhost:6379/0"
    assert settings.minio_endpoint == "localhost:9000"


def test_infrastructure_settings_defaults(env_setup):
    settings = InfrastructureSettings()
    assert settings.minio_bucket_3d == "daedalus-3d-assets"
    assert settings.minio_use_ssl is False
    assert settings.jaeger_host == "jaeger:6831"


def test_settings_all_subsettings(env_setup):
    settings = Settings()
    assert settings.infrastructure is not None
    assert settings.auth is not None
    assert settings.payment is not None
    assert settings.notification is not None
    assert settings.ai is not None
    assert settings.service is not None


def test_auth_settings_defaults(env_setup):
    from config.settings import AuthSettings
    settings = AuthSettings()
    assert settings.jwt_algorithm == "HS256"
    assert settings.jwt_access_expire_minutes == 30
    assert settings.jwt_refresh_expire_days == 7


def test_payment_settings_defaults(env_setup):
    from config.settings import PaymentSettings
    settings = PaymentSettings()
    assert settings.payment_currency_default == "XAF"


def test_service_settings_defaults(env_setup):
    from config.settings import ServiceSettings
    settings = ServiceSettings()
    assert settings.environment == "development"
    assert settings.log_level == "INFO"
