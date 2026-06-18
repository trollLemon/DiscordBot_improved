import logging
import uuid
from contextlib import asynccontextmanager
from io import BytesIO

from celery.result import AsyncResult
from fastapi import FastAPI, File, HTTPException, UploadFile
from fastapi.concurrency import run_in_threadpool
from fastapi.responses import JSONResponse
from PIL import Image, UnidentifiedImageError
from opentelemetry import trace, metrics
from opentelemetry.instrumentation.celery import CeleryInstrumentor
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor

from app.broker import celery
from app.storage import Storage
from app.telemetry import setup_otel

logger = logging.getLogger(__name__)

CONTENT_TYPE_EXTENSIONS = {"image/jpeg": "jpeg", "image/png": "png"}

storage = Storage()

def _get_trace_id() -> str:
    span = trace.get_current_span()
    ctx = span.get_span_context()
    if ctx and ctx.trace_id:
        return format(ctx.trace_id, "032x")
    return ""


tracer = trace.get_tracer("classification-service")
meter = metrics.get_meter("classification-service")

enqueued_counter = meter.create_counter(
    "classification.images.enqueued",
    unit="1",
    description="Number of images accepted and enqueued for classification",
)
rejected_counter = meter.create_counter(
    "classification.images.rejected",
    unit="1",
    description="Number of images rejected before enqueueing",
)


@asynccontextmanager
async def lifespan(app: FastAPI):
    try:
        setup_otel("classification-service")
        CeleryInstrumentor().instrument()
        storage.ensure_bucket()
    except Exception:
        logger.exception("could not initialize required resources")
    yield


app = FastAPI(lifespan=lifespan)
FastAPIInstrumentor.instrument_app(app)

@app.get("/readyz")
async def readyz():
    if storage is not None:
        return {"status": "ready"}
    else:
        raise HTTPException(status_code=503, detail="not ready")


@app.get("/healthz")
def health():
    return JSONResponse(status_code=200, content={"status": "ok"})


@app.post("/api/v1/images")
async def enqueue_classification(file: UploadFile = File(...)):
    with tracer.start_as_current_span("enqueue_image") as span:
        span.set_attribute("image.content_type", file.content_type or "")

        extension = CONTENT_TYPE_EXTENSIONS.get(file.content_type)
        if extension is None:
            span.set_attribute("image.rejected_reason", "unsupported_type")
            rejected_counter.add(1, {"reason": "unsupported_type"})
            return JSONResponse(status_code=400, content={"detail": f"Unsupported file type. {file.content_type} is not supported", "traceId": _get_trace_id()})

        content = await file.read()

        try:
            Image.open(BytesIO(content)).verify()
        except (UnidentifiedImageError, OSError, SyntaxError) as e:
            span.set_attribute("image.rejected_reason", "corrupt")
            span.record_exception(e)
            rejected_counter.add(1, {"reason": "corrupt"})
            logger.error(f"failed to open image: {e}")
            return JSONResponse(status_code=400, content={"detail": "Could not read image data", "traceId": _get_trace_id()})

        key = f"{uuid.uuid4().hex}.{extension}"
        span.set_attribute("image.key", key)
        span.set_attribute("image.size_bytes", len(content))
        await run_in_threadpool(storage.upload, key, content, file.content_type)

        result = celery.send_task("tasks.classification", args=[key])
        span.set_attribute("classification.job_id", result.task_id)
        enqueued_counter.add(1, {"content_type": file.content_type})

        return JSONResponse(status_code=201, content={"jobId": result.task_id})


@app.get("/api/v1/images/classifications/{task_id}")
def get_result(task_id: str):
    result = AsyncResult(task_id, app=celery)

    if result.state == "SUCCESS":
        return JSONResponse(status_code=200, content={"Class": result.result})
    elif result.state == "FAILURE":
        return JSONResponse(status_code=400, content={"detail": f"Task {task_id} failed: {str(result.result)}"})
    elif result.state == "PENDING":
        return JSONResponse(status_code=202, content={"detail": f"Task {task_id} pending"})
    elif result.state == "REVOKED":
        return JSONResponse(status_code=400, content={"detail": f"Task {task_id} revoked"})
    elif result.state == "RETRY":
        return JSONResponse(status_code=400, content={"detail": f"Task {task_id} will be retried"})
    else:
        return JSONResponse(status_code=404, content={"detail": f"Task {task_id} not found"})
