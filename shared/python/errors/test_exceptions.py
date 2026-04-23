import pytest
from errors.exceptions import (
    DaedalusValidationError,
    DaedalusNotFoundError,
    DaedalusUnauthorizedError,
    DaedalusForbiddenError,
    DaedalusPaymentError,
    ErrorCode,
    is_not_found,
    is_unauthorized,
    is_forbidden,
    is_validation,
)
from http import HTTPStatus


def test_validation_error():
    error = DaedalusValidationError("Invalid input", "field: email")
    assert error.code == ErrorCode.VALIDATION_ERROR
    assert error.status == HTTPStatus.BAD_REQUEST
    assert error.to_dict()["code"] == "VALIDATION_ERROR"


def test_not_found_error():
    error = DaedalusNotFoundError("project", "123")
    assert error.code == ErrorCode.NOT_FOUND
    assert error.status == HTTPStatus.NOT_FOUND
    assert "project" in error.message


def test_unauthorized_error():
    error = DaedalusUnauthorizedError("token expired")
    assert error.code == ErrorCode.UNAUTHORIZED
    assert error.status == HTTPStatus.UNAUTHORIZED


def test_forbidden_error():
    error = DaedalusForbiddenError("insufficient permissions")
    assert error.code == ErrorCode.FORBIDDEN
    assert error.status == HTTPStatus.FORBIDDEN


def test_payment_error():
    error = DaedalusPaymentError("payment failed", "card declined")
    assert error.code == ErrorCode.PAYMENT_ERROR
    assert error.status == HTTPStatus.PAYMENT_REQUIRED


def test_error_to_dict():
    error = DaedalusNotFoundError("project", "123")
    data = error.to_dict()
    assert "code" in data
    assert "message" in data
    assert "status" in data
    assert data["code"] == "NOT_FOUND"


def test_is_not_found():
    error = DaedalusNotFoundError("project", "123")
    assert is_not_found(error)
    assert not is_not_found(DaedalusValidationError(""))


def test_is_unauthorized():
    error = DaedalusUnauthorizedError("test")
    assert is_unauthorized(error)


def test_is_forbidden():
    error = DaedalusForbiddenError("test")
    assert is_forbidden(error)


def test_is_validation():
    error = DaedalusValidationError("invalid")
    assert is_validation(error)


def test_error_inheritance():
    error = DaedalusNotFoundError("project", "123")
    assert isinstance(error, Exception)
    assert isinstance(error, DaedalusNotFoundError)
