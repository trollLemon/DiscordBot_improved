from io import BytesIO

import pytest
from fastapi.testclient import TestClient
from PIL import Image

from app import main


class FakeAsyncResult:
    def __init__(self, task_id="generated-task-id"):
        self.task_id = task_id


class FakeResult:
    def __init__(self, state, result=None):
        self.state = state
        self.result = result


class FakeStorage:
    def __init__(self):
        self.uploaded = []

    def upload(self, key, data, content_type="application/octet-stream"):
        self.uploaded.append((key, data, content_type))

    def ensure_bucket(self):
        pass


def png_bytes():
    buf = BytesIO()
    Image.new("RGB", (10, 10), "red").save(buf, format="PNG")
    return buf.getvalue()


@pytest.fixture
def client():
    return TestClient(main.app)


@pytest.fixture
def storage(monkeypatch):
    fake = FakeStorage()
    monkeypatch.setattr(main, "storage", fake)
    return fake


@pytest.fixture
def sent(monkeypatch):
    calls = []

    def fake_send_task(name, args):
        calls.append({"name": name, "args": args})
        return FakeAsyncResult("job-123")

    monkeypatch.setattr(main.celery, "send_task", fake_send_task)
    return calls


def test_enqueue_uploads_and_returns_job_id(client, storage, sent):
    resp = client.post("/api/v1/images", files={"file": ("x.png", png_bytes(), "image/png")})

    assert resp.status_code == 201
    assert resp.json() == {"jobId": "job-123"}

    assert len(sent) == 1
    assert sent[0]["name"] == "tasks.classification"

    key = sent[0]["args"][0]
    assert key.endswith(".png")
    assert storage.uploaded[0][0] == key
    assert storage.uploaded[0][2] == "image/png"


def test_enqueue_rejects_unsupported_type(client, storage, sent):
    resp = client.post("/api/v1/images", files={"file": ("x.txt", b"hello", "text/plain")})

    assert resp.status_code == 400
    assert "text/plain" in resp.json()["detail"]
    assert storage.uploaded == []
    assert sent == []


def test_enqueue_rejects_corrupt_image(client, storage, sent):
    resp = client.post("/api/v1/images", files={"file": ("x.png", b"not a real png", "image/png")})

    assert resp.status_code == 400
    assert resp.json()["detail"] == "Could not read image data"
    assert storage.uploaded == []
    assert sent == []


@pytest.mark.parametrize(
    "state,expected_status",
    [
        ("PENDING", 202),
        ("REVOKED", 400),
        ("RETRY", 400),
        ("STARTED", 404),
    ],
)
def test_get_result_states(client, monkeypatch, state, expected_status):
    monkeypatch.setattr(main, "AsyncResult", lambda task_id, app: FakeResult(state))

    resp = client.get("/api/v1/images/classifications/abc")

    assert resp.status_code == expected_status


def test_get_result_success(client, monkeypatch):
    monkeypatch.setattr(main, "AsyncResult", lambda task_id, app: FakeResult("SUCCESS", "tabby cat"))

    resp = client.get("/api/v1/images/classifications/abc")

    assert resp.status_code == 200
    assert resp.json() == {"Class": "tabby cat"}


def test_get_result_failure(client, monkeypatch):
    monkeypatch.setattr(main, "AsyncResult", lambda task_id, app: FakeResult("FAILURE", "kaboom"))

    resp = client.get("/api/v1/images/classifications/abc")

    assert resp.status_code == 400
    assert "kaboom" in resp.json()["detail"]


def test_health(client):
    resp = client.get("/healthz")

    assert resp.status_code == 200
    assert resp.json() == {"status": "ok"}
