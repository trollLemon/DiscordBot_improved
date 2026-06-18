from io import BytesIO

import pytest
from PIL import Image

from app import tasks


class FakeClassifier:
    def __init__(self, label=None, exc=None):
        self.label = label
        self.exc = exc
        self.seen_image = None

    def predict(self, image, decider):
        self.seen_image = image
        if self.exc is not None:
            raise self.exc
        return self.label


def png_bytes():
    buf = BytesIO()
    Image.new('RGB', (8, 8), 'blue').save(buf, format='PNG')
    return buf.getvalue()


def test_classify_returns_label(monkeypatch):
    fake = FakeClassifier(label='cat')
    monkeypatch.setattr(tasks, 'get_model', lambda: fake)
    monkeypatch.setattr(tasks.storage, 'download', lambda key: png_bytes())

    result = tasks.classify('k.png')

    assert result == 'cat'
    assert fake.seen_image.mode == 'RGB'


def test_classify_downloads_by_key(monkeypatch):
    fake = FakeClassifier(label='dog')
    seen = []
    monkeypatch.setattr(tasks, 'get_model', lambda: fake)

    def fake_download(key):
        seen.append(key)
        return png_bytes()

    monkeypatch.setattr(tasks.storage, 'download', fake_download)

    tasks.classify('abc.png')

    assert seen == ['abc.png']


def test_classify_wraps_model_error(monkeypatch):
    fake = FakeClassifier(exc=ValueError('boom'))
    monkeypatch.setattr(tasks, 'get_model', lambda: fake)
    monkeypatch.setattr(tasks.storage, 'download', lambda key: png_bytes())

    with pytest.raises(RuntimeError) as excinfo:
        tasks.classify('k.png')

    assert 'boom' in str(excinfo.value)


def test_get_model_is_cached(monkeypatch):
    calls = []

    class Sentinel:
        pass

    def fake_loader(model_type):
        calls.append(model_type)
        return {'model': Sentinel(), 'preprocessor': Sentinel()}

    monkeypatch.setattr(tasks, '_model', None)
    monkeypatch.setattr(tasks.classifier, 'get_model_and_processor', fake_loader)

    first = tasks.get_model()
    second = tasks.get_model()

    assert first is second
    assert len(calls) == 1
