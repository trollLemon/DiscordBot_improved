BOT_NAME=bot
GOMANIP_NAME=gomanip
all: build


BOT_DIR=./bot/ 
GOMANIP_DIR=./gomanip/ 
CLASSIFICATION_DIR=./classification/ 


all: build

build: build-bot build-gomanip 

test: test-bot test-gomanip-docker test-classificaion 

build-bot:
	cd $(BOT_DIR) && go build -o ../../$(BOT_NAME) ./cmd/main.go

build-gomanip:
	cd $(GOMANIP_DIR) && go build -o ../../$(GOMANIP_DIR) ./...

test-bot:
	cd $(BOT_DIR) && go test ./... -cover -race -count=1 -v

test-gomanip-docker:
	cd $(GOMANIP_DIR) &&  docker build --target=tester -t gomanip-test . && docker run --rm gomanip-test   

test-classificaion:
	cd $(CLASSIFICATION_DIR) && pytest -v

compose-up:
	docker-compose -f docker/docker-compose/docker-compose.yaml up --build

.PHONY: all fmt vet build test coverage clean
