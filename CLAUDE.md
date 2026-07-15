# Minimum Viable Go for Rust Engineers

## Purpose

This repository teaches practical Go mappings to readers who already know Rust.
Every lesson is small enough for an LLM or human to read, predict, run, and
compare without framework context.

## Repository Map

- `01_arrays_slices.go` is the default lesson for `go run .`.
- Standalone lessons use `//go:build ignore` and run by explicit filename.
- `35_bad_pointer_index.go` intentionally fails because Go cannot index `*[]T`.
- Numeric prefixes are the curriculum order and must stay contiguous.
- There is no nested `ds` package or framework.
- Never add repository ignore rules for generated artifacts. Keep `.ship/`,
  `dist/`, `.tmp_check/`, and similar outputs visible to Git, then delete them
  before handoff.
- `internal/webapp/` and `web.go` (`-tags web`) add an optional interactive
  browser tour; `cmd/fakeplayground/` is a Playground stand-in for tests and
  smoke, never the public service. See the README's Interactive Web Tour.

## How to Explore

1. Read `README.md` for the scope and complete learning path.
2. List the final inventory with `rg --files -g '*.go' | sort`.
3. Run the default lesson with `GOCACHE=/tmp/go-from-rust-cache go run .`.
4. Read one lesson, predict its output or compiler error, then run that file.
5. Check language claims against the [Go documentation](https://go.dev/doc/) or
   [Go specification](https://go.dev/ref/spec), not README prose alone.
6. Compare the observed output with the source and revise the prediction.
7. Inspect `35_bad_pointer_index.go` last and confirm its compiler failure.

## Lesson Groups

- **Slices and memory:** lessons 01 and 03–13.
- **Language mappings:** lessons 02 and 14–19.
- **Ordering and iterators:** lessons 20–22.
- **Collections and synchronization:** lessons 23–27.
- **Control and failure:** lessons 28–31 and 35.
- **Services and data:** lessons 32–34.

## Editing Contract

- Keep one concept per self-contained file, usually under 80 lines.
- Put code first. Never add a prose preamble at the top of a lesson.
- Use short inline comments between code blocks to map Rust concepts to Go.
- Keep output deterministic unless unspecified/nondeterministic behavior is itself the concept.
- Avoid unnecessary helpers and abstractions. If a helper hides the concept,
  keep the operation inline.
- Do not add error handling, generalization, or completeness that obscures the
  concept being taught.
- Use a deliberate compiler failure only when the failure is the useful lesson.
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
for file in *.go; do
  case "$file" in 01_arrays_slices.go|35_bad_pointer_index.go) continue ;; esac
  go run "$file"
done
go run -race 26_concurrent_maps.go
go run -race 31_channels.go
test -z "$(gofmt -l $(find . -name '*.go' -not -path './.git/*'))"
```

Run the intentional failure separately. It must report that `*[]int` cannot be
indexed:

```sh
GOCACHE=/tmp/go-from-rust-cache go run 35_bad_pointer_index.go
```

For the web tour, `make build`, `make test`, `make check`, and `make
smoke-web` cover the `-tags web` build, the `internal/webapp` unit tests, and
an `agent-browser`-driven end-to-end pass against `cmd/fakeplayground`. Never
point `PLAYGROUND_URL` at the public Playground from tests or smoke.

Treat the Why Go section in `README.md` as project framing, not benchmark
evidence. Never invent universal speed ratios or turn tradeoffs into universal
performance claims.
