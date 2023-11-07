GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
SRC_FOLDER=cmd/kbsd
BINARY_NAME=kbs
BINARY_UNIX=$(BINARY_NAME)-amd64-linux
BINARY_DARWIN=$(BINARY_NAME)-amd64-darwin
COMMIT_HASH?=$(shell git rev-parse HEAD)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION?=$(shell git describe --dirty --tags --always)
LDFLAGS?="-X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildDate=${BUILD_DATE} -s -w"

clean: 
	$(GOCLEAN)
	rm bin/$(BINARY_DARWIN)
	rm bin/$(BINARY_UNIX)

.PHONY: build
build: ## Build binary for mac
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 ${GOBUILD} -ldflags ${LDFLAGS} -o bin/${BINARY_DARWIN} ./${SRC_FOLDER}/main.go

.PHONY: build-linux
build-linux: ## Build binary for Linux
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GOBUILD} -ldflags ${LDFLAGS} -o bin/${BINARY_UNIX} ./${SRC_FOLDER}/main.go


.PHONY: run-local
run-local: ## run project local
	docker compose up --build

.PHONY: clean-local
clean-local: ## clean docker-compsoe
	docker compose down

.PHONY: test
test: ## run unit tests
	${GOCMD} test -race ./...

.PHONY: start
start: ## start kbs service + localstack
	docker compose up

.PHONY: stop
stop: ## start kbs service + localstack
	docker compose down --volumes

.PHONY: localstack/start
localstack/start: ## start localstack
	docker compose up localstack

.PHONY: localstack/startd
localstack/startd: ## start localstack
	docker compose up -d localstack

.PHONY: localstack/stop
localstack/stop:
	docker compose down --volumes

.PHONY: table/scan
table/scan:
	aws dynamodb scan --table-name kbs --endpoint-url http://localhost:4566 --region us-east-1

.PHONY: table/create
table/create:
	aws dynamodb create-table \
	--table-name kbs \
	--attribute-definitions \
		AttributeName=id,AttributeType=S \
	--key-schema \
		AttributeName=id,KeyType=HASH \
	--provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
	--endpoint-url http://localhost:4566 --region us-east-1