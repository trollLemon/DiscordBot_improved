import os
from io import BytesIO

from minio import Minio


class Storage:
    """S3-compatible object store for image blobs."""

    def __init__(self, client: Minio = None, bucket: str = None):
        self.bucket = bucket or os.getenv("S3_BUCKET", "images")
        self._client = client or Minio(
            os.getenv("S3_ENDPOINT", "seaweedfs:8333"),
            access_key=os.getenv("S3_ACCESS_KEY", "classification"),
            secret_key=os.getenv("S3_SECRET_KEY", "classification-secret"),
            secure=os.getenv("S3_SECURE", "false").lower() == "true",
        )

    def ensure_bucket(self) -> None:
        if not self._client.bucket_exists(self.bucket):
            self._client.make_bucket(self.bucket)

    def upload(self, key: str, data: bytes, content_type: str = "application/octet-stream") -> None:
        self._client.put_object(self.bucket, key, BytesIO(data), length=len(data), content_type=content_type)

    def download(self, key: str) -> bytes:
        response = self._client.get_object(self.bucket, key)
        try:
            return response.read()
        finally:
            response.close()
            response.release_conn()

    def delete(self, key: str) -> None:
        self._client.remove_object(self.bucket, key)
