GOMANIP_NAME=gomanip
DOCKER_IMAGE=gomanip:latest
DOCKER_CONTAINER=gomanip-server
PORT=8080

build:
	go build -o$(GOMANIP_NAME) ./...

test:
	docker build --target=tester -t gomanip-test . && \
	docker run --rm -v "$(shell pwd)":/out gomanip-test \
		sh -c 'go test -v -coverprofile=/out/coverage.out -race -count=1 ./...' && \
	go tool cover -func=coverage.out

docker-build:
	docker build --target application -t $(DOCKER_IMAGE) .

docker-run: docker-build
	docker run --rm --name $(DOCKER_CONTAINER) \
		-p $(PORT):8080 \
		-e REDIS_ADDR=localhost:6379 \
		$(DOCKER_IMAGE) --address 0.0.0.0 --port 8080

compose-up:
	docker compose up --build

compose-down:
	docker compose down

.PHONY: all fmt vet build test coverage clean docker-build docker-run compose-up compose-down
