# Minimum Viable Go for Rust Engineers

## Purpose

This repository teaches practical Go mappings to readers who already know Rust.
Every lesson is small enough for an LLM or human to read, predict, run, and
compare without framework context.

## Repository Map

- `main.go` is the default lesson for `go run .`.
- Standalone lessons use `//go:build ignore` and run by explicit filename.
- `badidx.go` intentionally fails because Go cannot index `*[]T`.
- There is no nested `ds` package, framework, or build system.

## How to Explore

1. Read `README.md` for the scope and complete learning path.
2. List the final inventory with `rg --files -g '*.go' | sort`.
3. Run the default lesson with `GOCACHE=/tmp/gofrs-go-cache go run .`.
4. Read one lesson, predict its output or compiler error, then run that file.
5. Check language claims against the [Go documentation](https://go.dev/doc/) or
   [Go specification](https://go.dev/ref/spec), not README prose alone.
6. Compare the observed output with the source and revise the prediction.
7. Inspect `badidx.go` last and confirm its compiler failure is the lesson.

## Lesson Groups

- **Slices and memory:** `main.go`, `demo.go`, `ptr.go`, `badidx.go`,
  `sliceptrs.go`, `filter.go`, `pop.go`, `zero.go`, `shrink.go`, `edges.go`.
- **Language mappings:** `printing.go`, `rangekw.go`, `option.go`, `match.go`,
  `tuples.go`, `multireturn.go`, `ternary.go`.
- **Ordering, iterators, and data structures:** `sortstruct.go`,
  `genericsort.go`, `iters.go`, `datastructures.go`.
- **Errors, logging, and panics:** `errors.go`, `logging.go`, `panic.go`.

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
export GOCACHE=/tmp/gofrs-go-cache
go test ./...
go vet ./...
go run .
for file in *.go; do
  case "$file" in main.go|badidx.go) continue ;; esac
  go run "$file"
done
test -z "$(gofmt -l *.go)"
```

Run the intentional failure separately. It must report that `*[]int` cannot be
indexed:

```sh
GOCACHE=/tmp/gofrs-go-cache go run badidx.go
```

Treat the Why Go section in `README.md` as project framing, not benchmark
evidence. Never invent universal speed ratios or turn tradeoffs into universal
performance claims.
