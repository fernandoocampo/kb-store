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
# LDFLAGS?="-X main.Version=${VERSION} -X main.CommitHash=${COMMIT_HASH} -X main.BuildDate=${BUILD_DATE} -s -w"
LDFLAGS?="-X github.com/fernandoocampo/kb-store/apps/kbs/internal/setups.Version=${VERSION} -X github.com/fernandoocampo/kb-store/apps/kbs/internal/setups.CommitHash=${COMMIT_HASH} -X github.com/fernandoocampo/kb-store/apps/kbs/internal/setups.BuildDate=${BUILD_DATE} -s -w"
DOCKER_REPO?=fdocampo
DOCKER_CONTAINER=kbs

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

.PHONY: build-image
build-image: ## Build container image
	DOCKER_BUILDKIT=0 podman build --build-arg appVersion=${VERSION} --build-arg commitHash=${COMMIT_HASH} --no-cache -t $(DOCKER_REPO)/$(DOCKER_CONTAINER) .

.PHONY: run-container-local
run-container-local: ## Run new container local
	@podman run --rm -it -p 8080:8080 $(DOCKER_REPO)/$(DOCKER_CONTAINER)

.PHONY: podman-login
podman-login: ## login into docker registry
	@podman login --username $(DOCKER_USERNAME) --verbose docker.io

.PHONY: push-tag
push-tag: ## Push a new docker image tag for kbs app.
	@podman push $(DOCKER_REPO)/$(DOCKER_CONTAINER)

.PHONY: goreleaser-snapshot
goreleaser-snapshot: ## run a "local-only" release to see if it works using the release command
	@goreleaser release --snapshot --clean

.PHONY: goreleaser-check
goreleaser-check: ## verify your .goreleaser.yaml is valid
	@goreleaser check

.PHONY: goreleaser-darwin-amd64
goreleaser-darwin-amd64: # build the binary only for a given Darwin amd64,
	GOOS=darwin GOARCH=amd64 goreleaser build --snapshot --clean --single-target

.PHONY: goreleaser-linux-amd64
goreleaser-linux-amd64: # build the binary only for a given linux amd64,
	GOOS=linux GOARCH=amd64 goreleaser build --snapshot --clean --single-target

.PHONY: goreleaser-release
goreleaser-release: # release project
	@goreleaser release

.PHONY: run-local
run-local: ## run project local
	@podman-compose --file docker-compose.yml up --build

.PHONY: clean-local
clean-local: ## clean kbs containers
	@podman-compose --file docker-compose.yml down --volumes

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