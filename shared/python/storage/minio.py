import logging
from typing import Optional
from aioboto3 import Session

logger = logging.getLogger(__name__)


class MinIOClient:
    def __init__(
        self,
        endpoint: str,
        access_key: str,
        secret_key: str,
        use_ssl: bool = False,
        bucket_3d: str = "daedalus-3d-assets",
        bucket_exports: str = "daedalus-exports",
    ):
        self.endpoint = endpoint
        self.access_key = access_key
        self.secret_key = secret_key
        self.use_ssl = use_ssl
        self.bucket_3d = bucket_3d
        self.bucket_exports = bucket_exports
        self.session = Session()

    async def upload_file(
        self, bucket: str, object_name: str, file_path: str
    ) -> None:
        try:
            async with self.session.client(
                "s3",
                endpoint_url=self._endpoint_url(),
                aws_access_key_id=self.access_key,
                aws_secret_access_key=self.secret_key,
            ) as client:
                with open(file_path, "rb") as f:
                    await client.put_object(Bucket=bucket, Key=object_name, Body=f)
            logger.debug(f"Uploaded {object_name} to {bucket}")
        except Exception as e:
            logger.error(f"Failed to upload file: {e}")
            raise

    async def download_file(
        self, bucket: str, object_name: str, file_path: str
    ) -> None:
        try:
            async with self.session.client(
                "s3",
                endpoint_url=self._endpoint_url(),
                aws_access_key_id=self.access_key,
                aws_secret_access_key=self.secret_key,
            ) as client:
                obj = await client.get_object(Bucket=bucket, Key=object_name)
                with open(file_path, "wb") as f:
                    f.write(await obj["Body"].read())
            logger.debug(f"Downloaded {object_name} from {bucket}")
        except Exception as e:
            logger.error(f"Failed to download file: {e}")
            raise

    async def get_presigned_url(
        self, bucket: str, object_name: str, expiration: int = 3600
    ) -> str:
        try:
            async with self.session.client(
                "s3",
                endpoint_url=self._endpoint_url(),
                aws_access_key_id=self.access_key,
                aws_secret_access_key=self.secret_key,
            ) as client:
                url = await client.generate_presigned_url(
                    "get_object",
                    Params={"Bucket": bucket, "Key": object_name},
                    ExpiresIn=expiration,
                )
            return url
        except Exception as e:
            logger.error(f"Failed to get presigned URL: {e}")
            raise

    async def delete_object(self, bucket: str, object_name: str) -> None:
        try:
            async with self.session.client(
                "s3",
                endpoint_url=self._endpoint_url(),
                aws_access_key_id=self.access_key,
                aws_secret_access_key=self.secret_key,
            ) as client:
                await client.delete_object(Bucket=bucket, Key=object_name)
            logger.debug(f"Deleted {object_name} from {bucket}")
        except Exception as e:
            logger.error(f"Failed to delete object: {e}")
            raise

    async def object_exists(self, bucket: str, object_name: str) -> bool:
        try:
            async with self.session.client(
                "s3",
                endpoint_url=self._endpoint_url(),
                aws_access_key_id=self.access_key,
                aws_secret_access_key=self.secret_key,
            ) as client:
                await client.head_object(Bucket=bucket, Key=object_name)
            return True
        except Exception as e:
            if "NoSuchKey" in str(e) or "404" in str(e):
                return False
            logger.error(f"Failed to check object existence: {e}")
            raise

    async def list_objects(self, bucket: str, prefix: str = "") -> list[str]:
        objects = []
        try:
            async with self.session.client(
                "s3",
                endpoint_url=self._endpoint_url(),
                aws_access_key_id=self.access_key,
                aws_secret_access_key=self.secret_key,
            ) as client:
                paginator = client.get_paginator("list_objects_v2")
                async for page in paginator.paginate(Bucket=bucket, Prefix=prefix):
                    if "Contents" in page:
                        objects.extend([obj["Key"] for obj in page["Contents"]])
            return objects
        except Exception as e:
            logger.error(f"Failed to list objects: {e}")
            raise

    async def health_check(self) -> bool:
        try:
            async with self.session.client(
                "s3",
                endpoint_url=self._endpoint_url(),
                aws_access_key_id=self.access_key,
                aws_secret_access_key=self.secret_key,
            ) as client:
                await client.head_bucket(Bucket=self.bucket_3d)
            return True
        except Exception as e:
            logger.error(f"Health check failed: {e}")
            return False

    def _endpoint_url(self) -> str:
        protocol = "https" if self.use_ssl else "http"
        return f"{protocol}://{self.endpoint}"
