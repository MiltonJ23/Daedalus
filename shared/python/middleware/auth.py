import logging
import os
import jwt
from typing import Optional
from datetime import datetime, timedelta
from fastapi import HTTPException, Depends
from fastapi.security import HTTPBearer, HTTPAuthCredentials
from ..errors.exceptions import DaedalusUnauthorizedError

logger = logging.getLogger(__name__)
security = HTTPBearer()


class TokenManager:
    def __init__(self, secret: str, algorithm: str = "HS256"):
        self.secret = secret
        self.algorithm = algorithm

    def create_access_token(
        self, user_id: str, role: str, plan: str, expires_delta: Optional[timedelta] = None
    ) -> str:
        if expires_delta is None:
            expires_delta = timedelta(minutes=30)

        expire = datetime.utcnow() + expires_delta
        payload = {
            "user_id": user_id,
            "role": role,
            "plan": plan,
            "exp": expire,
            "iat": datetime.utcnow(),
        }

        encoded_jwt = jwt.encode(payload, self.secret, algorithm=self.algorithm)
        return encoded_jwt

    def verify_token(self, token: str) -> dict:
        try:
            payload = jwt.decode(token, self.secret, algorithms=[self.algorithm])
            return payload
        except jwt.ExpiredSignatureError:
            raise DaedalusUnauthorizedError("token expired")
        except jwt.InvalidTokenError:
            raise DaedalusUnauthorizedError("invalid token")


async def get_current_user(
    credentials: HTTPAuthCredentials = Depends(security),
    jwt_secret: str = Depends(lambda: _get_jwt_secret()),
) -> dict:
    token = credentials.credentials
    if not token:
        raise HTTPException(status_code=401, detail="Missing token")

    try:
        tm = TokenManager(jwt_secret)
        return tm.verify_token(token)
    except Exception as e:
        raise HTTPException(status_code=401, detail=str(e))


def _get_jwt_secret() -> str:
    secret = os.environ.get("JWT_SECRET")
    if not secret:
        raise HTTPException(status_code=500, detail="JWT_SECRET not configured")
    return secret
