import logging
from typing import Optional
import redis.asyncio as redis

logger = logging.getLogger(__name__)


class RedisClient:
    def __init__(self, redis_url: str):
        self.redis_url = redis_url
        self.client: Optional[redis.Redis] = None

    async def connect(self) -> None:
        try:
            self.client = redis.from_url(
                self.redis_url,
                encoding="utf8",
                decode_responses=True,
                socket_keepalive=True,
                socket_keepalive_options={1: 1, 2: 3, 3: 3},
            )
            await self.client.ping()
            logger.info("Connected to Redis")
        except Exception as e:
            logger.error(f"Failed to connect to Redis: {e}")
            raise

    async def get(self, key: str) -> Optional[str]:
        if not self.client:
            raise RuntimeError("Not connected to Redis")
        return await self.client.get(key)

    async def set(self, key: str, value: str, expiration: int = None) -> None:
        if not self.client:
            raise RuntimeError("Not connected to Redis")
        if expiration:
            await self.client.setex(key, expiration, value)
        else:
            await self.client.set(key, value)

    async def delete(self, *keys: str) -> None:
        if not self.client:
            raise RuntimeError("Not connected to Redis")
        if keys:
            await self.client.delete(*keys)

    async def exists(self, key: str) -> bool:
        if not self.client:
            raise RuntimeError("Not connected to Redis")
        return bool(await self.client.exists(key))

    async def incr(self, key: str) -> int:
        if not self.client:
            raise RuntimeError("Not connected to Redis")
        return await self.client.incr(key)

    async def decr(self, key: str) -> int:
        if not self.client:
            raise RuntimeError("Not connected to Redis")
        return await self.client.decr(key)

    async def expire(self, key: str, seconds: int) -> None:
        if not self.client:
            raise RuntimeError("Not connected to Redis")
        await self.client.expire(key, seconds)

    async def keys(self, pattern: str) -> list[str]:
        if not self.client:
            raise RuntimeError("Not connected to Redis")
        return await self.client.keys(pattern)

    async def flush(self) -> None:
        if not self.client:
            raise RuntimeError("Not connected to Redis")
        await self.client.flushdb()

    async def health_check(self) -> bool:
        try:
            if not self.client:
                return False
            await self.client.ping()
            return True
        except Exception as e:
            logger.error(f"Health check failed: {e}")
            return False

    async def close(self) -> None:
        if self.client:
            await self.client.close()
            logger.info("Closed Redis connection")


def cache_key_pattern(service: str, key: str) -> str:
    return f"cache:{service}:{key}"
