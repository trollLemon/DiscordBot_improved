APP_NAME=bot
all: build

fmt:
	go fmt ./...

vet:
	go vet ./...

build: fmt vet
	go build -o $(APP_NAME) ./cmd/*


compose-up:
	@echo "Starting Docker Compose services"
	@echo "If you are not in the docker group you will need to run this with sudo"
	docker-compose -f docker-compose/docker-compose.yaml up --build

test-coverage:
	go test -v ./... -cover -race -count=1

clean:
	go clean
	rm -f $(APP_NAME)

.PHONY: all fmt vet build test coverage clean
