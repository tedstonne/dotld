VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build dev test lint format release release-patch release-minor release-major release-dry e2e clean env

build:
	go build -ldflags "$(LDFLAGS)" -o dist/dotld .

dev: build
	cp dist/dotld ~/.local/bin/dotld

test:
	go run gotest.tools/gotestsum@latest --format testdox ./...

lint:
	go vet ./...

format:
	gofmt -w .

release: release-patch

release-patch:
	./scripts/release.sh --bump patch

release-minor:
	./scripts/release.sh --bump minor

release-major:
	./scripts/release.sh --bump major

release-dry:
	./scripts/release.sh --dry-run

e2e: env
	docker build -f Dockerfile.test -t dotld-e2e .
	docker run --rm --env-file .env dotld-e2e

env:
	./scripts/env.sh

clean:
	rm -rf dist
