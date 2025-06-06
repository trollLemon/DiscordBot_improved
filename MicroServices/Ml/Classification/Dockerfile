FROM python:3.13.3-slim AS tester

WORKDIR /app

COPY . .
RUN pip install --no-cache-dir -r requirements.txt
RUN pip install --no-cache-dir pytest
CMD ["pytest", "-v"]


FROM python:3.13.3-slim

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt
RUN addgroup --system --gid 101 celery && \
    adduser --system --uid 100 --ingroup celery celery

# Create a directory for Hugging Face cache and give ownership to celery
RUN mkdir -p /app/.cache/huggingface && chown -R celery:celery /app/.cache

# Set environment variables for cache
ENV HF_HOME=/app/.cache/huggingface
COPY . .

USER celery
EXPOSE 8081
