PYTHON ?= python3
VENV := venv
BIN := $(VENV)/bin

IMAGE := classification
COMPOSE_DEV := ../docker/docker-compose-develop/docker-compose.yaml
COMPOSE := ../docker/docker-compose/docker-compose.yaml

.DEFAULT_GOAL := help

.PHONY: help venv install test run-api run-worker \
        docker-build-api docker-build-worker docker-test \
        compose-up compose-down clean

help: ## Show this help
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

venv: ## Create the local virtualenv
	$(PYTHON) -m venv $(VENV)

install: venv ## Install API + worker deps and pytest into the venv
	$(BIN)/pip install --upgrade pip
	$(BIN)/pip install -r requirements-worker.txt
	$(BIN)/pip install pytest

test: ## Run the test suite
	$(BIN)/pytest -v

run-api: ## Run the FastAPI service locally on :8081
	$(BIN)/uvicorn app.main:app --host 0.0.0.0 --port 8081 --reload

run-worker: ## Run the Celery worker locally
	$(BIN)/celery -A app.tasks worker --loglevel=info --concurrency=1

docker-build-api: ## Build the api image
	docker build --target api -t $(IMAGE)-api .

docker-build-worker: ## Build the worker image
	docker build --target worker -t $(IMAGE)-worker .

docker-test: ## Build the tester image and run the suite in it
	docker build --target tester -t $(IMAGE)-test .
	docker run --rm $(IMAGE)-test

compose-up: ## Build and start the develop docker-compose stack
	docker compose -f $(COMPOSE_DEV) up --build

compose-down: ## Stop the develop docker-compose stack
	docker compose -f $(COMPOSE_DEV) down

clean: ## Remove the venv and Python caches
	rm -rf $(VENV) .pytest_cache
	find . -type d -name __pycache__ -prune -exec rm -rf {} +
