# Minimum Viable Go for Rust Engineers

Go books and tutorials teach the language from zero. This repository answers a
narrower question for experienced Rust programmers: “I know how to do this in
Rust—what should I actually write in Go?” I created it for myself as I was
struggling to map concepts I already knew in Rust to practical Go.

Each numbered `.go` file contains one small runnable lesson. Working through all
37 gives you practical mappings for Go slices and allocation, language
constructs, collections, synchronization, errors, HTTP, configuration, logging,
and JSON—the parts most likely to surprise an experienced Rust user.

## Run It

`go run .` runs the first lesson in `01_arrays_slices.go`. Pass a standalone
filename to run one focused lesson:

```sh
go run .
go run 05_append_capacity.go
go run 10_slice_edits.go
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

| Lesson | Value | Main point |
|---|:---:|---|
| [`01_arrays_slices.go`](01_arrays_slices.go) | 5 | Arrays copy all elements; assigned slices still describe shared backing storage. |
| [`02_copy.go`](02_copy.go) | 5 | `Clone` creates a new backing array; `copy` fills existing storage without growing it. |
| [`03_range.go`](03_range.go) | 5 | `range` yields indexes and copied values; strings yield byte indexes and decoded runes. |
| [`04_iterators.go`](04_iterators.go) | 4 | A custom type exposes `All() iter.Seq` so callers can range over generated values. |
| [`05_append_capacity.go`](05_append_capacity.go) | 5 | `append` reuses spare capacity when possible, so aliases may remain connected. |
| [`06_pointers.go`](06_pointers.go) | 4 | Passing a slice copies its header; element writes still reach the same backing array. |
| [`07_option.go`](07_option.go) | 5 | Comma-ok returns `(T, bool)`; `*T` is nullable but adds indirection and may escape. |
| [`08_pointer_elements.go`](08_pointer_elements.go) | 3 | A nil pointer element must be checked before it is dereferenced. |
| [`09_filter.go`](09_filter.go) | 4 | Build a filtered copy with `append`, or remove matches in place with `DeleteFunc`. |
| [`10_slice_edits.go`](10_slice_edits.go) | 5 | Slice stacks and vectors use reslicing, `Delete`, `Insert`, or manual swap-delete. |
| [`11_delete_tail.go`](11_delete_tail.go) | 5 | Reslicing hides elements but does not clear references in the backing array. |
| [`12_slice_memory.go`](12_slice_memory.go) | 5 | Subslices and `Clip` retain their backing array; `Clone` creates independent storage. |
| [`13_range_append.go`](13_range_append.go) | 4 | `range` fixes its iteration count before the loop, even if the slice grows. |
| [`14_switch.go`](14_switch.go) | 4 | Go has expression, condition-only, and dynamic-type switches. |
| [`15_tuples.go`](15_tuples.go) | 4 | Stored mixed values use structs; fixed homogeneous groups use arrays. |
| [`16_multiple_returns.go`](16_multiple_returns.go) | 5 | Go returns multiple results directly; they are not a tuple value. |
| [`17_variadic.go`](17_variadic.go) | 3 | A variadic parameter receives a slice; `...` expands a slice into arguments. |
| [`18_interfaces.go`](18_interfaces.go) | 5 | Interfaces are satisfied implicitly; value and pointer method sets differ. |
| [`19_typed_nil.go`](19_typed_nil.go) | 5 | An interface is nil only when both its dynamic type and value are absent. |
| [`20_printing.go`](20_printing.go) | 4 | `fmt` prints values with verbs; `io.Writer` variants choose the destination. |
| [`21_sorting.go`](21_sorting.go) | 4 | `Sort` handles natural order; `SortFunc` supplies custom and multi-key order. |
| [`22_maps_sets.go`](22_maps_sets.go) | 5 | Maps are built in; sets use `map[T]struct{}`, and deterministic output sorts keys. |
| [`23_heap.go`](23_heap.go) | 4 | `container/heap` adds heap operations to a type that defines storage and ordering. |
| [`24_concurrent_maps.go`](24_concurrent_maps.go) | 5 | Use a typed map with `Mutex` by default; `sync.Map` serves specific access patterns. |
| [`25_synchronization.go`](25_synchronization.go) | 5 | Mutexes protect mutable invariants; atomics update one value or publish an immutable snapshot. |
| [`26_channels.go`](26_channels.go) | 5 | Channels combine value transfer, blocking, synchronization, and close notification. |
| [`27_ring_buffer.go`](27_ring_buffer.go) | 4 | A bounded `VecDeque` mapping reuses fixed storage; blocking and concurrent access are separate choices. |
| [`28_barriers.go`](28_barriers.go) | 4 | A barrier counts arrivals, then broadcasts release; Go composes those roles from simpler primitives. |
| [`29_context.go`](29_context.go) | 5 | `context` carries cancellation through calls and reports why work stopped. |
| [`30_defer.go`](30_defer.go) | 5 | `defer` runs at function return in reverse order; arguments are captured immediately. |
| [`31_errors.go`](31_errors.go) | 5 | Create stable errors, wrap them with `%w`, and inspect the chain with `errors.Is`. |
| [`32_panic.go`](32_panic.go) | 5 | An unrecovered panic ends the process; deferred recovery is local to the panicking goroutine. |
| [`33_logging.go`](33_logging.go) | 5 | `slog` writes key-value fields and derives loggers with shared fields. |
| [`34_json.go`](34_json.go) | 5 | JSON uses exported fields and tags by default; each input is parsed at runtime. |
| [`35_http.go`](35_http.go) | 5 | `ServeMux` routes requests; clients send them and response bodies must be closed. |
| [`36_config_options.go`](36_config_options.go) | 4 | Functional options apply closures over a defaults struct, replacing the builder pattern. |
| [`37_http_middleware.go`](37_http_middleware.go) | 4 | Middleware is `func(http.Handler) http.Handler`; wrap handlers to add cross-cutting behavior. |

### Slice Operations

The [Go Wiki's SliceTricks](https://go.dev/wiki/SliceTricks) shows the underlying
append, copy, and reslice patterns. This guide groups them by the behavior that
matters and uses the current standard library where it is clearer:

| Need | Use |
|---|---|
| Copy or concatenate | `copy`, `slices.Clone`, or `slices.Concat`; lesson 02 distinguishes supplied from new storage. |
| Reserve or constrain capacity | `slices.Grow` or a full slice expression; lesson 05 shows both. |
| Filter or deduplicate | A loop, `slices.DeleteFunc`, or `slices.Compact` for adjacent duplicates; lesson 09 compares allocation. |
| Insert, delete, pop, or swap-delete | `slices.Insert`, `slices.Delete`, or reslicing; lesson 10 shows order preservation and tail clearing. |
| Drop references or retained storage | `clear`, `slices.Clip`, or `slices.Clone`; lessons 11–12 separate reachability from capacity. |
| Reverse, shuffle, batch, or slide a window | `slices.Reverse`, `rand.Shuffle`, `slices.Chunk`, or `s[i:i+n]`; use lesson 27's ring for a bounded queue. |

Start with a loop or the standard `slices` and `maps` packages. [samber/lo](https://lo.samber.dev/)
is useful when operations absent from the standard library—such as `Map` or
`GroupBy`—make application code clearer. Generics do not inherently add a
runtime penalty. A particular helper can still add callback, allocation, or
missed-inlining costs, so benchmark hot paths instead of assuming either loops
or generic helpers win. For runtime internals, [this slice internals article](https://themsaid.com/slice-internals-in-go)
is a useful implementation guide; growth rules are not language guarantees and
can change between Go releases.

## Interactive Web Tour

An optional `-tags web` build serves a browser tour of all 37 lessons: edit,
format, run, and navigate by hash (`#01`–`#37`). It embeds the README,
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
test helper refuses any listener that resolves to a non-loopback address and
executes code in a temporary directory; it is not a sandbox and must never be
exposed as a public service.

**Container deployment.** The multi-stage `Dockerfile` builds a static binary
and ships only that binary in an `alpine` runtime image, running as a non-root
user with `ADDR=:8080` and a `/health` healthcheck:

```sh
docker build -t go-from-rust-web .
docker run -p 8080:8080 go-from-rust-web
```

`make playtest` builds the binaries, starts `cmd/fakeplayground` and the web
server on loopback, and drives the tour end to end with `agent-browser`. It
checks imports, editing and running lessons, failure handling, and responsive
layouts.

`make full` runs formatting, vet, every Go lesson, race checks, unit tests, and
the browser playtest.

## Why Go

For services, workers, CLIs, and data plumbing, development and operational
simplicity usually matter more than extracting peak speed. Go offers static
types, native binaries, garbage collection, a broad standard library, and
straightforward deployment.

Coming from Rust specifically:

- **Minimal by design.** Go has few features and little ceremony, so there is
  usually one obvious way to write something and the language stays out of the
  way. Rust hands you more expressive power and, with it, more decisions.
- **Operational simplicity.** A Go service ships as one fast static binary, and
  because a single instance does a lot with goroutines and a broad stdlib, most
  deployments skip the load-balancing and packaging machinery a Rust service
  often accretes.
- **Value types are real.** Structs and arrays are stored inline, not forced
  through boxing or indirection, so for value-heavy, allocation-light code
  throughput often lands close to Rust.
- **Concurrency just works.** Goroutines and channels reach roughly Tokio-class
  ergonomics without borrow-checking, lifetimes, or `Send + Sync` bounds to
  satisfy, and without manual memory management — you trade some compile-time
  guarantees for far less cognitive load.
- **The tradeoff, plainly.** Reaching for Rust often buys speed you do not need
  and compile-time memory safety the garbage collector already gives you at
  near-zero cognitive cost. For most applications, Go is the cheaper path to the
  same outcome.

Python and TypeScript remain useful choices, but long-running systems can meet
runtime, packaging, or throughput limits. Rust, Java, and C# address broader or
lower-level needs, but can require more language, runtime, or architecture work
than a small service needs. Go is practical in the space between those tradeoffs.

This framing is based on [Code Like Go](https://krons.fiu.wtf/lore/go), which
calls the underlying discipline “remove to accelerate”: fewer language choices
leave fewer implementation and review decisions.

AI makes producing code easier without making generated code trustworthy. As
production rises, verification becomes the bottleneck. Go's constrained
language, visible control flow, `gofmt`, compiler, and tests can make generated
changes cheaper to inspect.

## License

Released into the public domain under the [`UNLICENSE`](UNLICENSE).
