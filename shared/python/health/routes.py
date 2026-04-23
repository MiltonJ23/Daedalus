import logging
from typing import Callable, Optional
from datetime import datetime
from fastapi import APIRouter

logger = logging.getLogger(__name__)

router = APIRouter()


class HealthStatus:
    UP = "UP"
    DOWN = "DOWN"


class HealthCheckResult:
    def __init__(self):
        self.status = HealthStatus.UP
        self.checks: dict[str, str] = {}
        self.timestamp = datetime.utcnow().isoformat()
        self.uptime_seconds = 0

    def to_dict(self) -> dict:
        return {
            "status": self.status,
            "checks": self.checks,
            "timestamp": self.timestamp,
            "uptime_seconds": self.uptime_seconds,
        }


class HealthChecker:
    def __init__(self):
        self.checkers: dict[str, Callable] = {}

    def register(self, name: str, checker: Callable) -> None:
        self.checkers[name] = checker

    async def check_all(self) -> HealthCheckResult:
        result = HealthCheckResult()

        for name, checker in self.checkers.items():
            try:
                is_healthy = await checker()
                result.checks[name] = HealthStatus.UP if is_healthy else HealthStatus.DOWN
                if not is_healthy:
                    result.status = HealthStatus.DOWN
            except Exception as e:
                logger.error(f"Health check failed for {name}: {e}")
                result.checks[name] = HealthStatus.DOWN
                result.status = HealthStatus.DOWN

        return result


health_checker = HealthChecker()


@router.get("/health", tags=["health"])
async def health():
    result = await health_checker.check_all()
    return result.to_dict()


@router.get("/health/live", tags=["health"])
async def liveness():
    return {"status": "alive"}


@router.get("/health/ready", tags=["health"])
async def readiness():
    result = await health_checker.check_all()
    if result.status == HealthStatus.UP:
        return {"status": "ready"}
    return {"status": "not ready", "checks": result.checks}
