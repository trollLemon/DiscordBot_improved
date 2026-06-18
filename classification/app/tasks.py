import logging
import time
from io import BytesIO
from typing import Optional

from celery.signals import worker_process_init
from opentelemetry import trace, metrics
from opentelemetry.instrumentation.celery import CeleryInstrumentor
from opentelemetry.trace import Status, StatusCode
from PIL import Image

from app import classifier
from app.broker import celery
from app.storage import Storage
from app.telemetry import setup_otel

logger = logging.getLogger(__name__)

storage = Storage()

tracer = trace.get_tracer("classification-worker")
meter = metrics.get_meter("classification-worker")

classification_counter = meter.create_counter(
    "classification.tasks.completed",
    unit="1",
    description="Number of classification tasks completed, labelled by status",
)
prediction_duration = meter.create_histogram(
    "classification.prediction.duration",
    unit="s",
    description="Wall-clock time spent running model inference",
)

_model: Optional[classifier.Classifier] = None


def get_model() -> classifier.Classifier:
    global _model
    if _model is None:
        config = classifier.get_model_and_processor(classifier.ModelType.VITBASE)
        _model = classifier.Classifier(
            processor=config['preprocessor'],
            model=config['model'],
            logger=logger,
        )
    return _model


@worker_process_init.connect
def _warm_worker(**_):
    try:
        setup_otel("classification-worker")
        CeleryInstrumentor().instrument()
    except Exception:
        logger.exception("could not initialize opentelemetry on worker init")
    try:
        storage.ensure_bucket()
    except Exception:
        logger.exception("could not ensure storage bucket on worker init")
    get_model()


@celery.task(name="tasks.classification")
def classify(key: str) -> Optional[str]:
    with tracer.start_as_current_span("classify_image") as span:
        span.set_attribute("image.key", key)
        image = Image.open(BytesIO(storage.download(key))).convert('RGB')
        try:
            start = time.perf_counter()
            label = get_model().predict(image, lambda x: x.argmax(-1).item())
            prediction_duration.record(time.perf_counter() - start)
            span.set_attribute("classification.label", label or "")
            classification_counter.add(1, {"status": "success"})
            return label
        except Exception as e:
            span.record_exception(e)
            span.set_status(Status(StatusCode.ERROR, str(e)))
            classification_counter.add(1, {"status": "error"})
            raise RuntimeError(f"Image classification failed: {str(e)}") from e
