import asyncio
import json
import pytest
from unittest.mock import AsyncMock, MagicMock, patch
from messaging.client import RabbitMQClient, MessageHandler


@pytest.fixture
def rabbitmq_url() -> str:
    return "amqp://guest:guest@localhost:5672/"


@pytest.fixture
def client(rabbitmq_url: str) -> RabbitMQClient:
    return RabbitMQClient(rabbitmq_url, worker_count=2)


def test_client_init(client: RabbitMQClient, rabbitmq_url: str) -> None:
    assert client.rabbitmq_url == rabbitmq_url
    assert client.worker_count == 2
    assert client.connection is None
    assert client.channel is None
    assert client.exchange is None


@pytest.mark.asyncio
async def test_publish_raises_when_not_connected(client: RabbitMQClient) -> None:
    with pytest.raises(RuntimeError, match="Not connected to RabbitMQ"):
        await client.publish("test.key", {"event": "test"})


@pytest.mark.asyncio
async def test_subscribe_raises_when_not_connected(client: RabbitMQClient) -> None:
    async def handler(body: bytes) -> None:
        pass

    with pytest.raises(RuntimeError, match="Not connected to RabbitMQ"):
        await client.subscribe("test-queue", "test.key", handler)


@pytest.mark.asyncio
async def test_health_check_returns_false_when_not_connected(client: RabbitMQClient) -> None:
    result = await client.health_check()
    assert result is False


@pytest.mark.asyncio
async def test_health_check_returns_false_for_closed_connection(client: RabbitMQClient) -> None:
    mock_connection = MagicMock()
    mock_connection.is_closed = True
    client.connection = mock_connection

    result = await client.health_check()
    assert result is False


@pytest.mark.asyncio
async def test_health_check_returns_true_for_open_connection(client: RabbitMQClient) -> None:
    mock_connection = MagicMock()
    mock_connection.is_closed = False
    client.connection = mock_connection

    result = await client.health_check()
    assert result is True


@pytest.mark.asyncio
async def test_publish_sends_message(client: RabbitMQClient) -> None:
    mock_exchange = AsyncMock()
    mock_channel = AsyncMock()
    client.channel = mock_channel
    client.exchange = mock_exchange

    message = {"event": "user.created", "user_id": "123"}
    await client.publish("user.created", message)

    mock_exchange.publish.assert_called_once()
    call_kwargs = mock_exchange.publish.call_args
    assert call_kwargs.kwargs["routing_key"] == "user.created"

    import aio_pika
    published_msg = call_kwargs.args[0]
    assert isinstance(published_msg, aio_pika.Message)
    assert json.loads(published_msg.body) == message


@pytest.mark.asyncio
async def test_close_when_not_connected(client: RabbitMQClient) -> None:
    await client.close()
    assert client.connection is None


@pytest.mark.asyncio
async def test_close_closes_connection(client: RabbitMQClient) -> None:
    mock_connection = AsyncMock()
    client.connection = mock_connection

    await client.close()
    mock_connection.close.assert_called_once()


@pytest.mark.asyncio
async def test_connect_failure_raises(client: RabbitMQClient) -> None:
    with patch("aio_pika.connect_robust", new_callable=AsyncMock) as mock_connect:
        mock_connect.side_effect = ConnectionError("refused")
        with pytest.raises(ConnectionError):
            await client.connect()
