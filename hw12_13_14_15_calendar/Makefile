BIN := "./bin"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
cdir = $(shell pwd)

.PHONY: build
build:
	mkdir -p bin
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./...

.PHONY: test
test:
	go test -race ./...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

.PHONY: lint
lint: install-lint-deps
	golangci-lint run ./...

launch:
	./bin/calendar -config=./configs/calendar.yml & \
	./bin/calendar_scheduler -config=./configs/scheduler.yml & \
	./bin/calendar_sender -config=./configs/sender.yml

.PHONY: run
run:
	make build && make launch

.PHONY: version
version: build
	$(BIN) version

.PHONY: migrate
migrate:
	./scripts/migrations.sh

up:
	docker-compose -f deployments/docker-compose.yml up -d
	
down:
	docker-compose -f deployments/docker-compose.yml down

.PHONY: generate
generate:
	protoc \
	-I /usr/local/include \
	-I ${cdir}/service/schema \
	--go_out=plugins=grpc:${cdir}/service/server \
	--grpc-gateway_out=logtostderr=true:${cdir}/service/server \
	${cdir}/service/schema/calendar.proto


# protoc -I /usr/local/include -I ${cdir}/service/grpc/schema --go_out=plugins=grpc:${cdir}/service/grpc/server ${cdir}/service/grpc/schema/calendar.proto