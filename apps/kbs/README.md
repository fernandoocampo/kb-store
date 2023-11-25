[![Go Report Card](https://goreportcard.com/badge/github.com/fernandoocampo/kb-store/apps/kbs)](https://goreportcard.com/report/github.com/fernandoocampo/kb-store/apps/kbs) ![Quality](https://github.com/fernandoocampo/kb-store/actions/workflows/kbs-app-quality.yml/badge.svg?branch=main)

# kbs service

This is a service that provides an API for kbs administration.

## How to build?

### raw build to use locally
from project folder run below commands, it will output binaries in `./bin/` folder

* just build with current operating system
```sh
make build
```

* build for a linux distro operating system
```sh
make build-linux
```

### goreleaser build to use locally
from project folder run below commands, it will output binaries in `./dist/` folder

* verify your `.goreleaser.yml` file.

```sh
make goreleaser-check
```

* run goreleaser `"local only"` to test the release
```sh
make goreleaser-snapshot
```

* release only for darwin amd64 (snapshot).
```sh
make goreleaser-darwin-amd64
```

* release only for linux amd64 (snapshot).
```sh
make goreleaser-linux-amd64
```

* release with goreleaser.
```sh
make goreleaser-release
```

### to use locally as a container

* build image
```sh
make build-image
```

* run container based on in the image created above.
```sh
make run-container-local
```

## How to run a test environment quickly?

1. make sure you have docker-compose installed.
2. run the docker compose.
```sh
docker-compose up --build
```

or run this shortcut

```sh
make run-local
```

3. once you are done using the environment follow these steps.

    * ctrl + c
    * make clean-local

## How to test?

from project folder run the following command

```sh
make test
```
