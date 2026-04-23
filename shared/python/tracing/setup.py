import logging
from fastapi import FastAPI
from opentelemetry import trace
from opentelemetry.exporter.jaeger.thrift import JaegerExporter
from opentelemetry.sdk.resources import Resource, SERVICE_NAME
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.requests import RequestsInstrumentor
from opentelemetry.instrumentation.sqlalchemy import SQLAlchemyInstrumentor

logger = logging.getLogger(__name__)


def setup_tracing(service_name: str, app: FastAPI, jaeger_host: str = "jaeger:6831") -> None:
    try:
        host_parts = jaeger_host.split(":")
        agent_host = host_parts[0]
        agent_port = int(host_parts[1]) if len(host_parts) > 1 else 6831

        jaeger_exporter = JaegerExporter(
            agent_host_name=agent_host,
            agent_port=agent_port,
        )

        provider = TracerProvider(
            resource=Resource.create({SERVICE_NAME: service_name})
        )
        provider.add_span_processor(BatchSpanProcessor(jaeger_exporter))
        trace.set_tracer_provider(provider)

        FastAPIInstrumentor.instrument_app(app)
        RequestsInstrumentor().instrument()

        logger.info(f"OpenTelemetry tracing initialized for {service_name}")
    except Exception as e:
        logger.error(f"Failed to initialize tracing: {e}")
        raise
