GOCACHE ?= /tmp/go-from-rust-cache
export GOCACHE
GOFILES := $(shell find . -name '*.go' -not -path './.git/*' -not -path './.tmp_check/*')

# Use the active toolchain's gofmt: a gofmt from an older Go on PATH cannot
# parse Go 1.27 syntax (generic methods) and reports a parse error instead.
GOFMT := $(shell go env GOROOT)/bin/gofmt

IMAGE ?= go-from-rust-web

# DOCKER may be overridden for hosts where the invoking user is in the docker
# group (then `make image DOCKER=docker`). Default is `sudo docker` so the
# target works consistently across dev hosts. sudo strips the environment, so
# DOCKER_BUILDKIT is injected via `env` after the sudo prefix.
DOCKER ?= sudo docker
DOCKER_SUDO = $(filter sudo,$(DOCKER))
DOCKER_BIN  = $(filter-out sudo,$(DOCKER))
DOCKER_BUILD = $(DOCKER_SUDO) env DOCKER_BUILDKIT=1 $(DOCKER_BIN) build

.PHONY: help build check test integration image clean playground-check

help:
	@echo "Usage:"
	@echo "    build        build lesson, web, and fakeplayground binaries"
	@echo "    check        gofmt and go vet (both build tags)"
	@echo "    test         fast unit tests"
	@echo "    integration  lesson runs, race checks, and browser playtest"
	@echo "    image        build the web Docker image ($(IMAGE))"
	@echo "    playground-check  does the public Playground speak Go 1.27 yet"
	@echo "    clean        remove build artifacts"

build:
	mkdir -p dist
	go build -o dist/go-from-rust .
	go build -tags web -o dist/go-from-rust-web .
	go build -o dist/fakeplayground ./cmd/fakeplayground

# Lessons 43 and 44 use Go 1.27 syntax the public Playground cannot compile
# yet, so their in-browser Run fails. This probes when that stops being true;
# when it reports READY, follow TODO.md to drop the caveats.
playground-check:
	@sed '1,2d' 43_generic_methods.go > $(GOCACHE)/pgcheck.go
	@if curl -fsS -X POST https://go.dev/_/compile \
	    --data-urlencode version=2 --data-urlencode withVet=false \
	    --data-urlencode body@$(GOCACHE)/pgcheck.go \
	    | grep -q 'no type parameters'; then \
	  echo "not yet: public Playground is still pre-1.27"; \
	else \
	  echo "READY: Playground compiles Go 1.27 — see TODO.md"; \
	fi

check:
	test -z "$$($(GOFMT) -l $(GOFILES))"
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
