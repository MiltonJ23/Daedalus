import logging
from typing import Optional
from fastapi import FastAPI
from opentelemetry import trace
from opentelemetry.exporter.jaeger.thrift import JaegerExporter
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.requests import RequestsInstrumentor
from opentelemetry.instrumentation.sqlalchemy import SQLAlchemyInstrumentor

logger = logging.getLogger(__name__)


def setup_tracing(service_name: str, jaeger_host: str = "jaeger:4317") -> None:
    try:
        jaeger_exporter = JaegerExporter(
            agent_host_name=jaeger_host.split(":")[0],
            agent_port=int(jaeger_host.split(":")[1]) if ":" in jaeger_host else 6831,
        )

        trace.set_tracer_provider(TracerProvider())
        trace.get_tracer_provider().add_span_processor(
            BatchSpanProcessor(jaeger_exporter)
        )

        FastAPIInstrumentor.instrument_app()
        RequestsInstrumentor().instrument()

        logger.info(f"OpenTelemetry tracing initialized for {service_name}")
    except Exception as e:
        logger.error(f"Failed to initialize tracing: {e}")
        raise
