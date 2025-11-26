import os
import tempfile

from celery.result import AsyncResult

from Broker.broker import celery
from PIL import Image
from fastapi import FastAPI, UploadFile, File
from fastapi.responses import JSONResponse
from io import BytesIO

app = FastAPI()


@app.post("/api/v1/images")
async def enqueue_classification(file: UploadFile = File(...)):
    if file.content_type not in ["image/jpeg", "image/png"]:
        return JSONResponse(status_code=400, content={"detail": f"Unsupported file type. {file.content_type} is not supported"})

    content = await file.read()

    image = Image.open(BytesIO(content))

    temp_dir = tempfile.mkdtemp(dir='/app/shared')
    image_path = os.path.join(temp_dir, file.filename)
    image.save(image_path)  # store image later for processing

    result = celery.send_task("tasks.classification", args=[image_path])#classify.delay(image_path)

    return JSONResponse(status_code=202, content={"jobId": result.task_id})


@app.get("/api/v1/images/classifications/{task_id}")
def get_result(task_id: str):
    result = AsyncResult(task_id, app=celery)

    if result.state == "SUCCESS":
        value = result.result
        return JSONResponse(status_code=200, content={"Class": value})
    elif result.state == "FAILURE":
        exc = result.result
        exc_str = str(exc)
        return JSONResponse(
            status_code=400,
            content={
                "detail": f"Task {task_id} failed: {exc_str}",
            }
        )
    elif result.state == "PENDING":
        return JSONResponse(status_code=202, content={"detail": f"Task {task_id} pending"})
    elif result.state == "REVOKED":
        return JSONResponse(status_code=400, content={"detail": f"Task {task_id} revoked"})
    elif result.state == "RETRY":
        return JSONResponse(status_code=400, content={"detail": f"Task {task_id} will be retried"})
    else:
        return JSONResponse(status_code=404, content={"detail": f"Task {task_id} not found"})
