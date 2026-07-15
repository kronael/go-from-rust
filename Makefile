GOCACHE ?= /tmp/go-from-rust-cache
export GOCACHE
GOFILES := $(shell find . -name '*.go' -not -path './.git/*' -not -path './.tmp_check/*')

.PHONY: build test check fmt vet run smoke-web clean

build:
	mkdir -p dist
	go build -o dist/go-from-rust .
	go build -tags web -o dist/go-from-rust-web .
	go build -o dist/fakeplayground ./cmd/fakeplayground

test:
	go test ./...
	go test -tags web ./...

fmt:
	test -z "$$(gofmt -l $(GOFILES))"

vet:
	go vet ./...
	go vet -tags web ./...

check: fmt vet test
	go run .
	for file in [0-9][0-9]_*.go; do \
		case "$$file" in \
			01_arrays_slices.go|35_bad_pointer_index.go) continue ;; \
		esac; \
		go run "$$file"; \
	done
	go run -race 26_concurrent_maps.go
	go run -race 31_channels.go
	@output="$$(go run 35_bad_pointer_index.go 2>&1)"; status=$$?; \
		test "$$status" -ne 0; \
		printf '%s\n' "$$output" | grep -q 'cannot index'

run:
	go run .

smoke-web: build
	scripts/smoke-web.sh

clean:
	rm -rf dist
