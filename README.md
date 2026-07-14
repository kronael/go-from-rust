# gofrs

**Go for Rustaceans: a minimal, runnable guide for senior Rust programmers.**

Each root `.go` file teaches one Go concept by comparing it with the Rust concept
you already know. Read one file, run it, and inspect the output.

There are no programming basics, framework conventions, or course machinery.
The goal is to make a senior Rust developer productive in Go without teaching
programming again.

## What it contains

- Standalone demos for slices, pointers, ranges, sorting, generics, iterators,
  multiple returns, and common collection operations.
- Direct translations for Rust concepts such as `Vec`, `Option`, `match`,
  tuples, trait bounds, and iterators.
- Printed output for every runnable lesson.
- Intentional compiler errors where the error is the lesson.
- Minimal teaching implementations of Heap, Stack, Queue, and Set.

This is not a production collections library. The code is deliberately small,
direct, and incomplete where completeness would hide the concept.

## What you get from it

After working through the demos, you can:

- open unfamiliar Go code and follow its data flow;
- write Go without translating Rust syntax literally;
- predict slice aliasing, backing-array retention, and `range` behavior;
- use Go's explicit replacements for Rust conveniences;
- review small Go changes and recognize the sharp edges;
- decide when Go is the better default than Rust or Python.

## Why Go

Most application code is plumbing: HTTP, databases, queues, files, schemas,
serialization, retries, and deployment. It needs to be quick to build, easy to
operate, and fast enough that infrastructure does not become the product.

Rust gives you control over allocation, layout, lifetimes, and every performance
edge. Most services do not need that control. Python removes initial friction,
but interpreter overhead can turn throughput limits and CPU use into deployment
problems.

Use Rust when control is the product. Use Go when the product is APIs, workers,
CLIs, and data plumbing.

Go takes the middle path:

- a small language and fast build loop;
- garbage collection instead of lifetime management;
- goroutines and a production-ready standard library;
- native code in one deployable binary;
- enough performance headroom to keep a simple deployment for longer.

There is no universal “Go is 100× faster” ratio. In one production example,
Stream reported that Go was typically **40× faster than Python** for its
workload. Its Python ranking code took three days to build and another two weeks
to optimize; the Go version took four days and needed no further performance
work. The useful result is not a benchmark trophy. It is less time tuning code
and fewer machines doing the same job.
([case study](https://go.dev/solutions/stream))

The design attitude follows
[Code Like Go](https://krons.fiu.wtf/lore/go): boring, explicit, and small.

## Run it

```sh
go run .                 # arrays, slices, aliasing, and sorting
go run filter.go         # three ways to filter a slice
go run option.go         # comma-ok and nil instead of Option<T>
go run datastructures.go # Heap, Stack, Queue, and Set
```

Most demos have `//go:build ignore` so they can each declare their own `main`
without colliding with the others. Passing a filename directly to `go run`
runs that one demo.

You need Go 1.26 or newer, matching [`go.mod`](go.mod). Check your installation:

```sh
go version
```

## How to learn from it

1. Pick one file from the learning path below.
2. Read the Rust comparison in its comments.
3. Predict the output or compiler error.
4. Run only that file with `go run <file>.go`.
5. Change one value and run it again.

The compiler-error demos are intentional:

```sh
go run badidx.go # cannot index a pointer to a slice
go run caseA.go  # cannot dereference an int
go run caseB.go  # the corrected pointer-element case
```

`badidx.go` and `caseA.go` are successful lessons when they fail to compile.

## Learning path

### 1. Slices and values

| Demo | Rust question it answers |
|---|---|
| [`main.go`](main.go) | How do arrays, `Vec<T>`, aliasing, and sorting translate? |
| [`demo.go`](demo.go) | What replaces `copy_from_slice` and variadic expansion? |
| [`rangekw.go`](rangekw.go) | What does `range` yield for slices, maps, strings, integers, and iterators? |
| [`ptr.go`](ptr.go) | When does Go dereference a pointer automatically? |
| [`badidx.go`](badidx.go) | Why does indexing `*[]T` fail? |

### 2. Slice operations and memory

| Demo | Rust question it answers |
|---|---|
| [`filter.go`](filter.go) | What replaces `iter().filter().collect()`? |
| [`pop.go`](pop.go) | What replaces `Vec::pop`, `remove`, and `swap_remove`? |
| [`zero.go`](zero.go) | Why does deleting pointer elements need zeroing? |
| [`shrink.go`](shrink.go) | When does a short slice keep a large backing array alive? |
| [`edges.go`](edges.go) | What happens when a loop appends to the slice it visits? |

### 3. Rust conveniences Go omits

| Demo | Rust question it answers |
|---|---|
| [`option.go`](option.go) | What replaces `Option<T>`? |
| [`match.go`](match.go) | What replaces `match` and destructuring? |
| [`tuples.go`](tuples.go) | What replaces tuple values? |
| [`multireturn.go`](multireturn.go) | Are multiple return values a tuple? |
| [`ternary.go`](ternary.go) | What replaces expression-valued `if` and ternary helpers? |
| [`caseA.go`](caseA.go), [`caseB.go`](caseB.go) | When is `*v` a valid dereference? |

### 4. Ordering, generics, and iteration

| Demo | Rust question it answers |
|---|---|
| [`sortstruct.go`](sortstruct.go) | What replaces `sort_by` and `sort_by_key`? |
| [`genericsort.go`](genericsort.go) | What is the Go equivalent of an `Ord`-like bound? |
| [`iters.go`](iters.go) | How do Go's push iterators differ from Rust's pull iterators? |

### 5. Small data structures

| Demo | Rust question it answers |
|---|---|
| [`datastructures.go`](datastructures.go) | How do generic Heap, Stack, Queue, and Set APIs look in use? |
| [`ds/ds.go`](ds/ds.go) | What do minimal teaching implementations look like? |

`ds` is lesson material, not a production collections library.

## Rust to Go map

| Rust | Go | Main difference |
|---|---|---|
| `[T; N]` | `[N]T` | Arrays are values; their length is part of the type. |
| `Vec<T>` | `[]T` | A slice is a small header over a backing array. |
| `Option<T>` | `(T, bool)` or `*T` | Go uses comma-ok or `nil`, depending on the boundary. |
| `Result<T, E>` | `(T, error)` | Errors are ordinary return values. |
| `?` | `if err != nil { ... }` | Error propagation stays explicit. |
| `match` | `switch` or a type switch | There is no destructuring or exhaustiveness check. |
| tuple | struct, array, or multiple return values | Multiple returns are not storable tuple values. |
| `Iterator` | `iter.Seq[T]` | Go iterators call `yield`; consumers do not call `next`. |
| trait bound | interface or type set | Interfaces are satisfied implicitly. |
| `BinaryHeap<T>` | `container/heap` or a local type | The standard library has no generic heap container. |
| `HashSet<T>` | `map[T]struct{}` | Sets are conventionally represented with a map. |
| `VecDeque<T>` | slice, ring, or local queue | The standard library has no generic deque. |
| operator overloading | none | Operators keep their built-in meanings. |

## Repository rules

- One concept per demo.
- Keep demos self-contained and roughly under 80 lines.
- Put the Rust comparison next to the relevant Go code.
- Print observable output when the demo runs.
- Prefer explicit code over helper abstractions.
- Keep deliberate compiler failures when the error is the lesson.

There is no framework or build system to learn. The Go toolchain is the runner.

## License

Released into the public domain under the [`UNLICENSE`](UNLICENSE).
