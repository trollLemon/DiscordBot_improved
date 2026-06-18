# syntax=docker/dockerfile:1

FROM python:3.13.3-slim AS api

WORKDIR /app
COPY requirements-api.txt .
RUN pip install --no-cache-dir -r requirements-api.txt

COPY app/ /app/app/

RUN addgroup --system --gid 101 app && \
    adduser --system --uid 100 --ingroup app app
USER app
EXPOSE 8081

ENTRYPOINT ["uvicorn", "app.main:app"]
CMD ["--host", "0.0.0.0", "--port", "8081", "--workers", "4"]


FROM python:3.13.3-slim AS worker

WORKDIR /app
COPY requirements-api.txt requirements-worker.txt ./
RUN pip install --no-cache-dir -r requirements-worker.txt

COPY app/ /app/app/

RUN addgroup --system --gid 101 celery && \
    adduser --system --uid 100 --ingroup celery celery

RUN mkdir -p /app/.cache/huggingface && chown -R celery:celery /app/.cache
ENV HF_HOME=/app/.cache/huggingface

USER celery

ENTRYPOINT ["celery", "-A", "app.tasks", "worker"]
CMD ["--loglevel=info", "--concurrency=1"]


FROM worker AS tester

USER root
COPY pytest.ini ./
COPY tests/ /app/tests/
RUN pip install --no-cache-dir pytest
CMD ["pytest", "-v"]
