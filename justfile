version := `git describe --tags --always --dirty 2>/dev/null || echo dev`
ldflags := "-s -w -X main.version=" + version

# build binary to dist/dotld
build:
    go build -ldflags "{{ldflags}}" -o dist/dotld .

# build and install to ~/.local/bin
dev: build
    cp dist/dotld ~/.local/bin/dotld

# run tests
test:
    go run gotest.tools/gotestsum@latest --format testdox ./...

# vet source code
lint:
    go vet ./...

# format source code
format:
    gofmt -w .

# tag and publish a release (patch, minor, major)
release bump="patch":
    ./scripts/release.sh --bump {{bump}}

# dry-run a release
release-dry:
    ./scripts/release.sh --dry-run

# build and run e2e tests in Docker
e2e:
    docker build --no-cache -f Dockerfile.test -t dotld-e2e .
    docker run --rm --env-file .env dotld-e2e

# shell into the e2e container for debugging
e2e-shell:
    docker build --no-cache -f Dockerfile.test -t dotld-e2e .
    docker run --rm -it --env-file .env dotld-e2e /bin/bash

# generate .env from 1Password
env:
    ./scripts/env.sh

# remove build artifacts
clean:
    rm -rf dist

# remove dotld binary and config from local machine
uninstall:
    rm -f ~/.local/bin/dotld
    rm -rf ~/.config/dotld
