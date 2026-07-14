# Minimum Viable Go for Rust Engineers

Go books and tutorials teach the language from zero. This repository answers a
narrower question for experienced Rust programmers: “I know how to do this in
Rust—what should I actually write in Go?” I made it because mapping concepts I
already knew in Rust to practical Go was the hard part.

Each `.go` file contains one small standalone lesson. Twenty-three are runnable;
`badidx.go` is the single intentional compile failure. After working through all
24, you should be able to reason about Go slices and memory, translate common
Rust language patterns, use the standard ordering and collection tools, and
handle errors, structured logs, and panics without importing Rust abstractions.

## Run It

`go run .` runs the default lesson in `main.go`. Pass a standalone filename to
run one focused lesson:

```sh
go run .
go run errors.go
go run logging.go
go run panic.go
go run datastructures.go
```

The standalone lessons use `//go:build ignore`, so they do not collide with the
default `main`. Explicit filenames still run normally.

Only `badidx.go` is intended to fail:

```sh
go run badidx.go # expected: a *[]int cannot be indexed
```

The repository requires Go 1.26, matching [`go.mod`](go.mod).

## Learning Path

| Lesson | What it maps from Rust to Go |
|---|---|
| [`printing.go`](printing.go) | `print!`, `println!`, formatting verbs, and writer output |
| [`main.go`](main.go) | Array value copies versus slice header copies and shared elements |
| [`demo.go`](demo.go) | `copy_from_slice`, Go's shorter-copy rule, and variadic expansion |
| [`rangekw.go`](rangekw.go) | `range` over slices, maps, strings, and integers, plus value-copy behavior |
| [`ptr.go`](ptr.go) | Slice arguments, pointers to slice headers, and array pointer indexing |
| [`badidx.go`](badidx.go) | The compiler error from indexing `*[]T` |
| [`sliceptrs.go`](sliceptrs.go) | Nil checks and dereferencing elements of `[]*T` |
| [`filter.go`](filter.go) | Allocating, in-place, and `slices.DeleteFunc` filtering |
| [`pop.go`](pop.go) | `Vec::pop`, ordered removal, `swap_remove`, and front reslicing |
| [`zero.go`](zero.go) | Pointer retention after reslicing versus tail clearing after deletion |
| [`shrink.go`](shrink.go) | Backing-array retention, capacity, and `slices.Clone` |
| [`edges.go`](edges.go) | Appending while `range` visits its original iteration length |
| [`option.go`](option.go) | Map comma-ok results and `*T` when zero differs from absence |
| [`match.go`](match.go) | Slice splitting, type switches, and value switches |
| [`tuples.go`](tuples.go) | Named structs and fixed arrays instead of stored tuple values |
| [`multireturn.go`](multireturn.go) | Multiple result lists and why they are not tuple values |
| [`ternary.go`](ternary.go) | Branch assignment instead of expression-valued `if` |
| [`sortstruct.go`](sortstruct.go) | Struct comparators, `cmp.Or`, and unstable sorting |
| [`genericsort.go`](genericsort.go) | `cmp.Ordered`, generic sorting, and explicit struct ordering |
| [`iters.go`](iters.go) | Push iterators, early stop, and slice and map iterator adapters |
| [`datastructures.go`](datastructures.go) | Slice stacks and queues, set maps, and `container/heap` |
| [`errors.go`](errors.go) | Value-plus-error results, `%w` wrapping, and `errors.Is` |
| [`logging.go`](logging.go) | Structured `slog` fields and loggers with shared context |
| [`panic.go`](panic.go) | `panic`, `recover`, and the goroutine recovery boundary |

## Why Go

For services, workers, CLIs, and data plumbing, development and operational
simplicity usually matter more than extracting peak speed. Go offers static
types, native binaries, garbage collection, a broad standard library, and
straightforward deployment.

Python and TypeScript remain useful choices, but long-running systems can meet
runtime, packaging, or throughput limits. Rust, Java, and C# address broader or
lower-level needs, but can require more language, runtime, or architecture work
than a small service needs. Go is practical in the space between those tradeoffs.

[Code Like Go](https://krons.fiu.wtf/lore/go) calls the underlying discipline
“remove to accelerate”: fewer language choices leave fewer implementation and
review decisions.

AI makes producing code easier. It does not make generated code trustworthy.
When production rises, verification becomes the bottleneck. Go's constrained
language, visible control flow, `gofmt`, compiler, and tests can make generated
changes cheaper to inspect.

## License

Released into the public domain under the [`UNLICENSE`](UNLICENSE).
