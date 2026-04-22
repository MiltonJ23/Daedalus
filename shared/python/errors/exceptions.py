from enum import Enum
from typing import Optional
from http import HTTPStatus


class ErrorCode(str, Enum):
    VALIDATION_ERROR = "VALIDATION_ERROR"
    NOT_FOUND = "NOT_FOUND"
    UNAUTHORIZED = "UNAUTHORIZED"
    FORBIDDEN = "FORBIDDEN"
    CONFLICT = "CONFLICT"
    INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
    PAYMENT_ERROR = "PAYMENT_ERROR"
    SERVICE_UNAVAILABLE = "SERVICE_UNAVAILABLE"
    BAD_GATEWAY = "BAD_GATEWAY"


class DaedalusError(Exception):
    def __init__(
        self,
        code: ErrorCode,
        message: str,
        details: Optional[str] = None,
        status: int = HTTPStatus.INTERNAL_SERVER_ERROR,
    ):
        self.code = code
        self.message = message
        self.details = details
        self.status = status
        super().__init__(f"[{code}] {message}")

    def to_dict(self):
        return {
            "code": self.code.value,
            "message": self.message,
            "details": self.details,
            "status": self.status,
        }


class DaedalusValidationError(DaedalusError):
    def __init__(self, message: str, details: Optional[str] = None):
        super().__init__(
            code=ErrorCode.VALIDATION_ERROR,
            message=message,
            details=details,
            status=HTTPStatus.BAD_REQUEST,
        )


class DaedalusNotFoundError(DaedalusError):
    def __init__(self, resource: str, resource_id: str):
        super().__init__(
            code=ErrorCode.NOT_FOUND,
            message=f"{resource} not found",
            details=f"ID: {resource_id}",
            status=HTTPStatus.NOT_FOUND,
        )


class DaedalusUnauthorizedError(DaedalusError):
    def __init__(self, reason: str):
        super().__init__(
            code=ErrorCode.UNAUTHORIZED,
            message="Unauthorized access",
            details=reason,
            status=HTTPStatus.UNAUTHORIZED,
        )


class DaedalusForbiddenError(DaedalusError):
    def __init__(self, reason: str):
        super().__init__(
            code=ErrorCode.FORBIDDEN,
            message="Access forbidden",
            details=reason,
            status=HTTPStatus.FORBIDDEN,
        )


class DaedalusConflictError(DaedalusError):
    def __init__(self, message: str):
        super().__init__(
            code=ErrorCode.CONFLICT,
            message=message,
            status=HTTPStatus.CONFLICT,
        )


class DaedalusPaymentError(DaedalusError):
    def __init__(self, message: str, details: Optional[str] = None):
        super().__init__(
            code=ErrorCode.PAYMENT_ERROR,
            message=message,
            details=details,
            status=HTTPStatus.PAYMENT_REQUIRED,
        )


class DaedalusServiceUnavailableError(DaedalusError):
    def __init__(self, service: str):
        super().__init__(
            code=ErrorCode.SERVICE_UNAVAILABLE,
            message=f"{service} service temporarily unavailable",
            status=HTTPStatus.SERVICE_UNAVAILABLE,
        )


class DaedalusBadGatewayError(DaedalusError):
    def __init__(self, details: str):
        super().__init__(
            code=ErrorCode.BAD_GATEWAY,
            message="Bad gateway",
            details=details,
            status=HTTPStatus.BAD_GATEWAY,
        )


def is_not_found(error: Exception) -> bool:
    return isinstance(error, DaedalusNotFoundError)


def is_unauthorized(error: Exception) -> bool:
    return isinstance(error, DaedalusUnauthorizedError)


def is_forbidden(error: Exception) -> bool:
    return isinstance(error, DaedalusForbiddenError)


def is_validation(error: Exception) -> bool:
    return isinstance(error, DaedalusValidationError)
