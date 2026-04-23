import asyncio
import json
import logging
from typing import Any, Awaitable, Callable, Optional
import aio_pika
from aio_pika import IncomingMessage

logger = logging.getLogger(__name__)

MessageHandler = Callable[[bytes], Awaitable[Any]]


class RabbitMQClient:
    def __init__(self, rabbitmq_url: str, worker_count: int = 5):
        self.rabbitmq_url = rabbitmq_url
        self.worker_count = worker_count
        self.connection: Optional[aio_pika.Connection] = None
        self.channel: Optional[aio_pika.Channel] = None
        self.exchange: Optional[aio_pika.Exchange] = None
        self.semaphore = asyncio.Semaphore(worker_count)

    async def connect(self) -> None:
        try:
            self.connection = await aio_pika.connect_robust(self.rabbitmq_url)
            self.channel = await self.connection.channel()
            self.exchange = await self.channel.declare_exchange(
                "daedalus.events",
                aio_pika.ExchangeType.TOPIC,
                durable=True,
            )
            logger.info("Connected to RabbitMQ")
        except Exception as e:
            logger.error(f"Failed to connect to RabbitMQ: {e}")
            raise

    async def publish(self, routing_key: str, message: dict) -> None:
        if not self.channel or not self.exchange:
            raise RuntimeError("Not connected to RabbitMQ")

        try:
            msg = aio_pika.Message(
                body=json.dumps(message).encode("utf-8"),
                content_type="application/json",
                delivery_mode=aio_pika.DeliveryMode.PERSISTENT,
            )
            await self.exchange.publish(msg, routing_key=routing_key)
            logger.debug(f"Published message to {routing_key}")
        except Exception as e:
            logger.error(f"Failed to publish message: {e}")
            raise

    async def subscribe(
        self, queue_name: str, routing_key: str, handler: MessageHandler
    ) -> None:
        if not self.channel:
            raise RuntimeError("Not connected to RabbitMQ")

        try:
            queue = await self.channel.declare_queue(queue_name, durable=True)
            await queue.bind(self.exchange, routing_key=routing_key)

            async with queue.iterator() as queue_iter:
                async for message in queue_iter:
                    async with self.semaphore:
                        try:
                            await handler(message.body)
                            await message.ack()
                        except Exception as e:
                            logger.error(f"Handler failed: {e}")
                            await message.nack(requeue=True)
        except Exception as e:
            logger.error(f"Subscription failed: {e}")
            raise

    async def health_check(self) -> bool:
        try:
            if not self.connection or self.connection.is_closed:
                return False
            return True
        except Exception as e:
            logger.error(f"Health check failed: {e}")
            return False

    async def close(self) -> None:
        if self.connection:
            await self.connection.close()
            logger.info("Closed RabbitMQ connection")
