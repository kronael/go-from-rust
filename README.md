# Minimum Viable Go for Rust Engineers

Go books and tutorials teach the language from zero. This repository answers a
narrower question for experienced Rust programmers: “I know how to do this in
Rust—what should I actually write in Go?” I made it because mapping concepts I
already knew in Rust to practical Go was the hard part.

Each numbered `.go` file contains one small standalone lesson. Thirty-four are
runnable; `35_bad_pointer_index.go` is the single intentional compile failure.
Working through all 35 gives you practical mappings for Go slices and
allocation, language constructs, collections, synchronization, errors,
logging, and JSON—the parts most likely to surprise an experienced Rust user.

## Run It

`go run .` runs the first lesson in `01_arrays_slices.go`. Pass a standalone
filename to run one focused lesson:

```sh
go run .
go run 06_append_capacity.go
go run 10_slice_edits.go
go run 29_errors.go
go run 33_logging.go
```

The standalone lessons use `//go:build ignore`, so they do not collide with the
default `main`. Explicit filenames still run normally.

Run the concurrent-state lessons with the race detector as well:

```sh
go run -race 26_concurrent_maps.go
go run -race 31_channels.go
```

Only `35_bad_pointer_index.go` is intended to fail:

```sh
go run 35_bad_pointer_index.go # expected: a *[]int cannot be indexed
```

The repository requires Go 1.26, matching [`go.mod`](go.mod).

## Learning Path

Value scores are editorial: **5** is essential, **4** is frequent, and **3** is
narrower but still worth recognizing. The removed ternary helper scored 2.

| Lesson | Value | Rust question it answers |
|---|:---:|---|
| [`01_arrays_slices.go`](01_arrays_slices.go) | 5 | When do assignments copy values or alias elements? |
| [`02_printing.go`](02_printing.go) | 4 | What replaces `print!`, `println!`, Display, and Debug? |
| [`03_copy.go`](03_copy.go) | 5 | How does `copy` differ from `copy_from_slice` and `copy_within`? |
| [`04_variadic.go`](04_variadic.go) | 4 | How are variadic calls declared and expanded? |
| [`05_range.go`](05_range.go) | 5 | What does `range` yield for each built-in type? |
| [`06_append_capacity.go`](06_append_capacity.go) | 5 | When does `append` alias or allocate? |
| [`07_pointers.go`](07_pointers.go) | 4 | When are slice and array pointers useful? |
| [`08_slice_pointers.go`](08_slice_pointers.go) | 3 | How do nil checks and dereferencing work in `[]*T`? |
| [`09_filter.go`](09_filter.go) | 4 | What replaces `filter().collect()` and in-place retain? |
| [`10_slice_edits.go`](10_slice_edits.go) | 5 | What are the pop, delete, insert, and swap-delete idioms? |
| [`11_delete_tail.go`](11_delete_tail.go) | 5 | How do reslicing and deletion affect retained references? |
| [`12_slice_memory.go`](12_slice_memory.go) | 5 | Do `Clip`, `Clone`, and subslices release allocations? |
| [`13_range_append.go`](13_range_append.go) | 4 | What happens when a range loop appends to its slice? |
| [`14_option.go`](14_option.go) | 5 | What replaces `Option<T>` in common APIs? |
| [`15_match.go`](15_match.go) | 4 | What replaces matching, destructuring, and enum dispatch? |
| [`16_tuples.go`](16_tuples.go) | 4 | What stores heterogeneous or fixed-size grouped values? |
| [`17_multiple_returns.go`](17_multiple_returns.go) | 5 | Why are multiple results syntax rather than a tuple type? |
| [`18_interfaces.go`](18_interfaces.go) | 5 | How do implicit interfaces and method sets replace traits? |
| [`19_typed_nil.go`](19_typed_nil.go) | 5 | Why can an interface containing a nil pointer be non-nil? |
| [`20_generic_sort.go`](20_generic_sort.go) | 3 | What does `cmp.Ordered` actually constrain? |
| [`21_struct_sort.go`](21_struct_sort.go) | 4 | How are custom and multi-key orderings expressed? |
| [`22_iterators.go`](22_iterators.go) | 4 | How do push iterators and iterator adapters work? |
| [`23_maps_sets.go`](23_maps_sets.go) | 5 | How do maps, sets, and deterministic key order work? |
| [`24_deque.go`](24_deque.go) | 4 | What replaces `VecDeque`, and how does a bounded ring work? |
| [`25_heap.go`](25_heap.go) | 4 | How does `container/heap` replace `BinaryHeap`? |
| [`26_concurrent_maps.go`](26_concurrent_maps.go) | 5 | When should a map use `Mutex` or `sync.Map`? |
| [`27_synchronization.go`](27_synchronization.go) | 5 | When should code use `Mutex`, `RWMutex`, atomics, or no spinlock? |
| [`28_defer.go`](28_defer.go) | 5 | How do deferred calls differ from lexical `Drop`? |
| [`29_errors.go`](29_errors.go) | 5 | How are sentinel errors created, wrapped, and inspected? |
| [`30_panic.go`](30_panic.go) | 5 | When does a panic kill the process, and where can it recover? |
| [`31_channels.go`](31_channels.go) | 5 | When should channels transfer ownership, and what does closing mean? |
| [`32_context.go`](32_context.go) | 5 | How is cancellation passed through service code? |
| [`33_logging.go`](33_logging.go) | 5 | How does standard structured logging carry fields and context? |
| [`34_json.go`](34_json.go) | 5 | Does JSON use generated code, cached reflection, or repeated parsing? |
| [`35_bad_pointer_index.go`](35_bad_pointer_index.go) | 3 | Why can Go not index through `*[]T`? |

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
| `PLAYGROUND_USER_AGENT` | `go-from-rust/1 (+https://github.com/kronael/go-from-rust)` | User-Agent sent to the compile service |

**Playground policy and security boundary.** The server never executes
submitted code itself. `/api/run` forwards the submitted source to the
Playground compile protocol at `PLAYGROUND_URL`, bounded by a 128 KiB request
cap, a 12-second upstream timeout, 8 concurrent in-flight compiles, a 120
requests/minute budget, and a 256-entry 10-minute result cache. Point
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

`make smoke-web` builds the binaries, starts `cmd/fakeplayground` and the web
server on loopback, and drives the tour end to end with the `agent-browser`
CLI (navigation, edit persistence, format, reset, run, lesson 35's expected
failure, and mobile layout).

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
