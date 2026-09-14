# Minimum Viable Go for Rust Engineers

## Purpose

This repository teaches practical Go mappings to readers who already know Rust.
Every lesson is small enough for an LLM or human to read, predict, run, and
compare without framework context.

## Repository Map

- `01_arrays_slices.go` is the default lesson for `go run .`.
- Standalone lessons use `//go:build ignore` and run by explicit filename.
- Numeric prefixes are the curriculum order and must stay contiguous.
- There is no nested `ds` package or framework.
- Generated artifacts (`dist/`, `.tmp_check/`, `.ship/`, `tmp/`) and local
  per-machine settings are covered by `.gitignore`. Keep it lean, and still
  delete stray outputs before handoff.
- `internal/webapp/` and `web.go` (`-tags web`) add an optional interactive
  browser tour; `cmd/fakeplayground/` is a Playground stand-in for tests and
  smoke, never the public service. See the README's Interactive Web Tour.

## How to Explore

1. Read `README.md` for the scope and complete learning path.
2. List the final inventory with `rg --files -g '*.go' | sort`.
3. Run the default lesson with `GOCACHE=/tmp/go-from-rust-cache go run .`.
4. Read one lesson, predict its output, then run that file.
5. Check language claims against the [Go documentation](https://go.dev/doc/) or
   [Go specification](https://go.dev/ref/spec), not README prose alone.
6. Compare the observed output with the source and revise the prediction.
7. Inspect `35_http.go` last and trace the handler through the client call.

## Lesson Groups

- **Foundations and slices:** lessons 01–13.
- **Language mappings:** lessons 14–21.
- **Collections and ordering:** lessons 22–24.
- **Concurrency:** lessons 25–30.
- **Control and failure:** lessons 31–32.
- **Services and data:** lessons 33–37.
- **Generics and tooling:** lessons 38–44.
- **Type modeling:** lesson 45.
- **What Go does not enforce:** lessons 46–50. Each shows a guarantee Rust
  makes and Go does not, then the Go idiom that lives with it where one
  exists.

## Editing Contract

- Keep one concept per self-contained file, usually under 80 lines.
- Put code first. Never add a prose preamble at the top of a lesson.
- Use short inline comments between code blocks to map Rust concepts to Go.
- Keep output deterministic unless unspecified/nondeterministic behavior is itself the concept.
- Avoid unnecessary helpers and abstractions. If a helper hides the concept,
  keep the operation inline.
- Do not add error handling, generalization, or completeness that obscures the
  concept being taught.
- Add no new file unless it demonstrates a genuinely distinct concept.
- Keep the tone pragmatic, dry, and free of marketing language.

## Validation

Use a writable cache:

```sh
export GOCACHE=/tmp/go-from-rust-cache
go test ./...
go test -tags web .
go vet ./...
go vet -tags web .
go run .
for file in [0-9][0-9]_*.go; do
  case "$file" in 01_arrays_slices.go) continue ;; esac
  go run "$file"
done
go run -race 25_concurrent_maps.go
go run -race 27_channels.go
go run -race 29_barriers.go
test -z "$($(go env GOROOT)/bin/gofmt -l $(find . -name '*.go' -not -path './.git/*'))"
```

`make integration` runs the lesson runs, race checks, and the browser playtest
against `cmd/fakeplayground`; `make check test integration` is the full local
gate. Never point `PLAYGROUND_URL` at the public Playground from tests.

Treat the Why Go section in `README.md` as project framing, not benchmark
evidence. Never invent universal speed ratios or turn tradeoffs into universal
performance claims.
