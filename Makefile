GOCACHE ?= /tmp/go-from-rust-cache
export GOCACHE

.PHONY: build test check fmt vet run smoke-web clean

build:
	go build -o dist/go-from-rust .
	go build -tags web -o dist/go-from-rust-web .
	go build -o dist/fakeplayground ./cmd/fakeplayground

test:
	go test ./...

fmt:
	test -z "$$(gofmt -l *.go internal cmd)"

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
	! go run 35_bad_pointer_index.go

run:
	go run .

smoke-web: build
	scripts/smoke-web.sh

clean:
	rm -rf dist
