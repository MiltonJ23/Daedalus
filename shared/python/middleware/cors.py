import logging
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

logger = logging.getLogger(__name__)


def setup_cors(app: FastAPI, allowed_origins: list[str] = None) -> None:
    if allowed_origins is None:
        allowed_origins = [
            "http://localhost:3000",
            "http://localhost:8080",
            "https://daedalus.io",
        ]

    app.add_middleware(
        CORSMiddleware,
        allow_origins=allowed_origins,
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    logger.info(f"CORS middleware configured for origins: {allowed_origins}")


def setup_logging(level: str = "INFO") -> None:
    logging.basicConfig(
        level=level,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    )
    logger.info(f"Logging configured at level {level}")
