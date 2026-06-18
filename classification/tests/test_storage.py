from app.storage import Storage


class FakeResponse:
    def __init__(self, data):
        self._data = data
        self.closed = False
        self.released = False

    def read(self):
        return self._data

    def close(self):
        self.closed = True

    def release_conn(self):
        self.released = True


class FakeMinioClient:
    def __init__(self, existing_buckets=()):
        self.buckets = set(existing_buckets)
        self.objects = {}
        self.removed = []
        self.last_response = None

    def bucket_exists(self, bucket):
        return bucket in self.buckets

    def make_bucket(self, bucket):
        self.buckets.add(bucket)

    def put_object(self, bucket, key, data, length, content_type="application/octet-stream"):
        self.objects[(bucket, key)] = (data.read(), length, content_type)

    def get_object(self, bucket, key):
        self.last_response = FakeResponse(self.objects[(bucket, key)][0])
        return self.last_response

    def remove_object(self, bucket, key):
        self.removed.append((bucket, key))
        self.objects.pop((bucket, key), None)


def make_storage(**kwargs):
    client = FakeMinioClient(**kwargs)
    return Storage(client=client, bucket="images"), client


def test_ensure_bucket_creates_when_missing():
    storage, client = make_storage()
    storage.ensure_bucket()
    assert "images" in client.buckets


def test_ensure_bucket_noop_when_exists():
    storage, client = make_storage(existing_buckets=["images"])
    storage.ensure_bucket()
    assert client.buckets == {"images"}


def test_upload_records_length_and_content_type():
    storage, client = make_storage(existing_buckets=["images"])
    storage.upload("k.png", b"hello", "image/png")
    data, length, content_type = client.objects[("images", "k.png")]
    assert data == b"hello"
    assert length == 5
    assert content_type == "image/png"


def test_upload_download_roundtrip():
    storage, _ = make_storage(existing_buckets=["images"])
    storage.upload("k.png", b"payload")
    assert storage.download("k.png") == b"payload"


def test_download_releases_connection():
    storage, client = make_storage(existing_buckets=["images"])
    storage.upload("k.png", b"payload")
    storage.download("k.png")
    assert client.last_response.closed
    assert client.last_response.released


def test_delete_removes_object():
    storage, client = make_storage(existing_buckets=["images"])
    storage.upload("k.png", b"payload")
    storage.delete("k.png")
    assert ("images", "k.png") in client.removed
