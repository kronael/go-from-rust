GOCACHE ?= /tmp/go-from-rust-cache
export GOCACHE
GOFILES := $(shell find . -name '*.go' -not -path './.git/*' -not -path './.tmp_check/*')

.PHONY: build test check playtest full clean

build:
	mkdir -p dist
	go build -o dist/go-from-rust .
	go build -tags web -o dist/go-from-rust-web .
	go build -o dist/fakeplayground ./cmd/fakeplayground

test:
	go test ./...
	go test -tags web .
	go run .
	for file in [0-9][0-9]_*.go; do \
		case "$$file" in \
			01_arrays_slices.go) continue ;; \
		esac; \
		go run "$$file"; \
	done
	CGO_ENABLED=1 CC="$(CC)" go run -race 24_concurrent_maps.go
	CGO_ENABLED=1 CC="$(CC)" go run -race 26_channels.go
	CGO_ENABLED=1 CC="$(CC)" go run -race 28_barriers.go

check:
	test -z "$$(gofmt -l $(GOFILES))"
	go vet ./...
	go vet -tags web .

playtest: build
	scripts/playtest.sh

full: check test playtest

clean:
	rm -rf dist .tmp_check
