# Minimum Viable Go for Rust Engineers

Go books and tutorials teach the language from zero. This repository answers a
narrower question for experienced Rust programmers: “I know how to do this in
Rust—what should I actually write in Go?” I made it because mapping concepts I
already knew in Rust to practical Go was the hard part.

Each numbered `.go` file contains one small runnable lesson. Working through all
35 gives you practical mappings for Go slices and allocation, language
constructs, collections, synchronization, errors, HTTP, logging, and JSON—the
parts most likely to surprise an experienced Rust user.

## Run It

`go run .` runs the first lesson in `01_arrays_slices.go`. Pass a standalone
filename to run one focused lesson:

```sh
go run .
go run 04_append_capacity.go
go run 08_slice_edits.go
go run 31_errors.go
go run 33_logging.go
go run 35_http.go
```

The standalone lessons use `//go:build ignore`, so they do not collide with the
default `main`. Explicit filenames still run normally.

Run the concurrent-state lessons with the race detector as well:

```sh
go run -race 24_concurrent_maps.go
go run -race 26_channels.go
go run -race 28_barriers.go
```

The repository requires Go 1.26, matching [`go.mod`](go.mod).

## Learning Path

Value scores are editorial: **5** is essential, **4** is frequent, and **3** is
narrower but still worth recognizing.

| Lesson | Value | Rust question it answers |
|---|:---:|---|
| [`01_arrays_slices.go`](01_arrays_slices.go) | 5 | What is a slice header, and why does assignment alias elements? |
| [`02_copy.go`](02_copy.go) | 5 | When should I use `slices.Clone` instead of `copy`? |
| [`03_range.go`](03_range.go) | 5 | What does `range` yield for each built-in type? |
| [`04_append_capacity.go`](04_append_capacity.go) | 5 | When does `append` alias or allocate? |
| [`05_pointers.go`](05_pointers.go) | 4 | When are slice and array pointers useful? |
| [`06_pointer_elements.go`](06_pointer_elements.go) | 3 | How do nil checks and dereferencing work in `[]*T`? |
| [`07_filter.go`](07_filter.go) | 4 | What replaces `filter().collect()` and in-place retain? |
| [`08_slice_edits.go`](08_slice_edits.go) | 5 | What are the pop, delete, insert, and swap-delete idioms? |
| [`09_delete_tail.go`](09_delete_tail.go) | 5 | How do reslicing and deletion affect retained references? |
| [`10_slice_memory.go`](10_slice_memory.go) | 5 | Do `Clip`, `Clone`, and subslices release allocations? |
| [`11_range_append.go`](11_range_append.go) | 4 | What happens when a range loop appends to its slice? |
| [`12_option.go`](12_option.go) | 5 | What replaces `Option<T>` in common APIs? |
| [`13_switch.go`](13_switch.go) | 4 | Which three `switch` forms replace `match` and if chains? |
| [`14_tuples.go`](14_tuples.go) | 4 | What stores heterogeneous or fixed-size grouped values? |
| [`15_multiple_returns.go`](15_multiple_returns.go) | 5 | Why are multiple results syntax rather than a tuple type? |
| [`16_variadic.go`](16_variadic.go) | 3 | How are variadic calls declared and expanded? |
| [`17_interfaces.go`](17_interfaces.go) | 5 | How do implicit interfaces and method sets replace traits? |
| [`18_typed_nil.go`](18_typed_nil.go) | 5 | Why can an interface containing a nil pointer be non-nil? |
| [`19_printing.go`](19_printing.go) | 4 | What replaces `print!`, `println!`, Display, and Debug? |
| [`20_sorting.go`](20_sorting.go) | 4 | How do natural, custom, and multi-key orderings work? |
| [`21_iterators.go`](21_iterators.go) | 4 | How do push iterators and iterator adapters work? |
| [`22_maps_sets.go`](22_maps_sets.go) | 5 | How do maps, sets, and deterministic key order work? |
| [`23_heap.go`](23_heap.go) | 4 | How does `container/heap` replace `BinaryHeap`? |
| [`24_concurrent_maps.go`](24_concurrent_maps.go) | 5 | When should a map use `Mutex` or `sync.Map`? |
| [`25_synchronization.go`](25_synchronization.go) | 5 | When should code use `Mutex`, `RWMutex`, atomics, or no spinlock? |
| [`26_channels.go`](26_channels.go) | 5 | When should channels own state, and what does closing mean? |
| [`27_ring_buffer.go`](27_ring_buffer.go) | 4 | How does a bounded ring differ from a channel? |
| [`28_barriers.go`](28_barriers.go) | 4 | How do goroutines wait until every worker is ready? |
| [`29_context.go`](29_context.go) | 5 | How is cancellation passed through service code? |
| [`30_defer.go`](30_defer.go) | 5 | How do deferred calls differ from lexical `Drop`? |
| [`31_errors.go`](31_errors.go) | 5 | How are sentinel errors created, wrapped, and inspected? |
| [`32_panic.go`](32_panic.go) | 5 | When does a panic kill the process, and where can it recover? |
| [`33_logging.go`](33_logging.go) | 5 | How does standard structured logging carry fields and context? |
| [`34_json.go`](34_json.go) | 5 | Does JSON use generated code, cached reflection, or repeated parsing? |
| [`35_http.go`](35_http.go) | 5 | How do I serve and call HTTP without a framework? |

## Interactive Web Tour

An optional `-tags web` build serves a browser tour of all 35 lessons: edit,
format, run, and navigate by hash (`#01`–`#35`). It embeds the README,
numbered lesson files, and `web/` assets into one binary; no lesson source is
duplicated. The production web server never executes submitted code locally.

```sh
go build -tags web -o go-from-rust-web .
./go-from-rust-web            # listens on 127.0.0.1:3999
```

Environment variables:

| Variable | Default | Purpose |
|---|---|---|
| `ADDR` | `127.0.0.1:3999` | Listen address (container sets `:8080`) |
| `PLAYGROUND_URL` | `https://go.dev/_/compile` | Compile-service endpoint `/api/run` proxies to |

**Playground policy and security boundary.** The server never executes
submitted code itself. `/api/run` forwards the submitted source to the
Playground compile protocol at `PLAYGROUND_URL`, bounded by a 128 KiB request
cap, a 12-second upstream timeout, 8 concurrent compiles, a 120-request
one-minute budget, and a 256-entry 10-minute result cache. Point
`PLAYGROUND_URL` at `cmd/fakeplayground` for local development and tests. That
loopback-only test helper executes code in a temporary directory; it is not a
sandbox and must never be exposed as a public service.

**Container deployment.** The multi-stage `Dockerfile` builds a static binary
and ships only that binary in an `alpine` runtime image, running as a non-root
user with `ADDR=:8080` and a `/health` healthcheck:

```sh
docker build -t go-from-rust-web .
docker run -p 8080:8080 go-from-rust-web
```

`make playtest` builds the binaries, starts `cmd/fakeplayground` and the web
server on loopback, and drives the tour end to end with `agent-browser`. It
checks syntax highlighting, navigation, edit persistence, format, reset, Run
success and failure, the HTTP lesson, mobile layout, and the wide-screen cap.

`make full` runs formatting, vet, every Go lesson, race checks, unit tests, and
the browser playtest.

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
