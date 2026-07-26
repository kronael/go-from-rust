GOCACHE ?= /tmp/go-from-rust-cache
export GOCACHE
GOFILES := $(shell find . -name '*.go' -not -path './.git/*' -not -path './.tmp_check/*')

IMAGE ?= go-from-rust-web

# DOCKER may be overridden for hosts where the invoking user is in the docker
# group (then `make image DOCKER=docker`). Default is `sudo docker` so the
# target works consistently across dev hosts. sudo strips the environment, so
# DOCKER_BUILDKIT is injected via `env` after the sudo prefix.
DOCKER ?= sudo docker
DOCKER_SUDO = $(filter sudo,$(DOCKER))
DOCKER_BIN  = $(filter-out sudo,$(DOCKER))
DOCKER_BUILD = $(DOCKER_SUDO) env DOCKER_BUILDKIT=1 $(DOCKER_BIN) build

.PHONY: help build check test integration image clean

help:
	@echo "Usage:"
	@echo "    build        build lesson, web, and fakeplayground binaries"
	@echo "    check        gofmt and go vet (both build tags)"
	@echo "    test         fast unit tests"
	@echo "    integration  lesson runs, race checks, and browser playtest"
	@echo "    image        build the web Docker image ($(IMAGE))"
	@echo "    clean        remove build artifacts"

build:
	mkdir -p dist
	go build -o dist/go-from-rust .
	go build -tags web -o dist/go-from-rust-web .
	go build -o dist/fakeplayground ./cmd/fakeplayground

check:
	test -z "$$(gofmt -l $(GOFILES))"
	go vet ./...
	go vet -tags web .

test:
	go test ./...
	go test -tags web .

integration: build
	go run .
	for file in [0-9][0-9]_*.go; do \
		case "$$file" in \
			01_arrays_slices.go) continue ;; \
		esac; \
		go run "$$file"; \
	done
	CGO_ENABLED=1 CC="$(CC)" go run -race 25_concurrent_maps.go
	CGO_ENABLED=1 CC="$(CC)" go run -race 27_channels.go
	CGO_ENABLED=1 CC="$(CC)" go run -race 29_barriers.go
	scripts/playtest.sh

image:
	$(DOCKER_BUILD) -t $(IMAGE) .

clean:
	rm -rf dist .tmp_check
