# Minimum Viable Go for Rust Engineers

Go books and tutorials teach the language from zero. This repository answers a
narrower question for experienced Rust programmers: “I know how to do this in
Rust—what should I actually write in Go?”

I built it for myself. Coming from Rust the syntax was never the problem—a
week of it and you write Go that compiles. What kept catching me was the gap
underneath: reaching for `Option<T>` and finding two idioms instead of one,
expecting a clone to be deep, assuming a value that owns something cannot be
copied out from under me. Each of those is a habit Rust made correct and Go
quietly does not, and none of them appears in a tutorial that starts at hello
world.

So each numbered `.go` file is one of those collisions, small enough to read,
predict, and run inside a minute. Working through all 46 gives you practical
mappings for Go slices and allocation, language constructs, collections,
synchronization, errors, HTTP, configuration, logging, and JSON—the parts most
likely to surprise an experienced Rust user. The last five are a different
kind: guarantees Rust makes and Go does not, where there is no mapping to
learn, only a habit to unlearn.

## Run It

`go run .` runs the first lesson in `01_arrays_slices.go`. Pass a standalone
filename to run one focused lesson:

```sh
go run .
go run 05_append_capacity.go
go run 08_slice_edits.go
go run 27_errors.go
go run 29_logging.go
go run 31_http.go
```

The standalone lessons use `//go:build ignore`, so they do not collide with the
default `main`. Explicit filenames still run normally.

Run the concurrent-state lessons with the race detector as well:

```sh
go run -race 21_concurrent_maps.go
go run -race 23_channels.go
go run -race 25_barriers.go
```

The repository requires Go 1.27 (currently the `go1.27rc2` release
candidate), matching [`go.mod`](go.mod).

## Learning Path

Value scores are editorial: **5** is essential, **4** is frequent, and **3** is
narrower but still worth recognizing.

| Lesson | Value | Main point |
|---|:---:|---|
| [`01_arrays_slices.go`](01_arrays_slices.go) | 5 | Go arrays are values, so assignment copies every element; a slice is a header (pointer, length, capacity) over a backing array, so assigning it copies only that descriptor and both slices still reach the same storage. Rust arrays copy only when their elements are `Copy`. |
| [`02_copy.go`](02_copy.go) | 5 | `slices.Clone` allocates a new backing array but copies elements shallowly — unlike Rust's `to_vec` it never clones each element, so pointer-like elements stay shared; `copy` fills only min(len) elements into existing storage, returns that count, never grows the destination, and supports overlapping slices like Rust's `copy_within`. |
| [`03_range.go`](03_range.go) | 5 | `range` over a slice yields indexes and copied values (writing the copy leaves the slice unchanged — assign through the index); over a string it yields byte offsets and decoded runes, so an index can jump by more than one for multi-byte characters; over an integer `n` it counts `0..n`, like Rust's `0..n`. Rust's `for` consumes an `IntoIterator`, whereas Go `range` accepts only specific operand types. |
| [`04_iterators.go`](04_iterators.go) | 4 | A custom type exposes an `All() iter.Seq[int]` method — a push-based function that feeds values into `yield` (returning `false` when the caller breaks) — so callers can `range` over it directly; Go cannot range an arbitrary type, and `iter.Seq` is push-based where Rust's `Iterator::next` pulls. Standard helpers like `slices.Collect` consume the same `Seq`. |
| [`05_append_capacity.go`](05_append_capacity.go) | 5 | `append` reuses spare capacity when it exists, so earlier aliases keep pointing at the same storage; a full slice expression (`s[:2:2]`) caps a view so the next `append` must reallocate and detach; `slices.Grow` reserves headroom without changing length and may reallocate, so keep its return value. Rust forbids a live slice borrow across a `Vec` mutation that could move it, and Go never specifies `append`'s growth factor — never depend on it. |
| [`06_pointers.go`](06_pointers.go) | 4 | Passing a slice copies its header, not its elements, so a callee's element writes still reach the caller's backing array; a `*[]int` (roughly `&mut Vec<_>`) lets a callee replace the caller's slice value itself, though returning a new slice is usually clearer when only length or storage changes. Go also lets `p[i]` on an array pointer stand for `(*p)[i]`. |
| [`07_option.go`](07_option.go) | 5 | Go has no `Option<T>`. A function returns `(T, bool)` — the comma-ok idiom — where the bool tells a real zero from absent and forces the caller to check. `*T` represents optionality by pointer (`nil` is absent), the idiom for optional struct fields or reference semantics, but a pointer that escapes heap-allocates. Only the pointer form survives inside a container: a slice holds one value per slot, so `[]*int` is roughly `Vec<Option<&i32>>` and every element needs its own nil check, because dereferencing nil panics. `(T, bool)` is a multiple-result list, not a value, so it cannot be stored. |
| [`08_slice_edits.go`](08_slice_edits.go) | 5 | Go has no filter, pop, insert, or swap-delete built in. Append matches into a fresh slice — the analogue of Rust's `iter().filter().collect()`, and it allocates — or drop non-matches in place with `slices.DeleteFunc`, which reuses the original backing array. Pop by guarding, reading the end, `clear`-ing that slot (which writes `0`, or `nil` for pointer elements), then reslicing. `slices.Delete` and `slices.Insert` shift elements to preserve order like `Vec::remove` and `Vec::insert`; swap-delete moves the last value over the hole and zeros the old slot, which is O(1) but reorders. |
| [`09_slice_memory.go`](09_slice_memory.go) | 5 | Reslicing shortens the length and nothing else: dropped elements stay reachable through the backing array, so pointer-like values stay retained where Rust's `truncate` would drop them. `clear` writes the element type's zero value (`nil` for `*int`) before the reslice, and `slices.Delete` does that clearing for you. The same rule bites at scale — a three-element view keeps a thousand-element allocation alive, and `slices.Clip` only caps capacity without copying, so it still aliases. `slices.Clone` (like Rust's `to_vec`) copies survivors into independent storage, which is what actually lets the large array be collected. |
| [`10_range_append.go`](10_range_append.go) | 4 | `range` evaluates the operand's length once before the loop, so appends made inside the body do not extend the iteration; Rust rejects pushing to a `Vec` while a shared borrow iterates it, catching the same hazard at compile time. |
| [`11_switch.go`](11_switch.go) | 4 | Go has three switch forms: an expression switch (comma-separated cases, no implicit fall-through — use an explicit `fallthrough`), a condition-only switch that replaces an if/else-if chain, and a type switch where `v.(type)` binds the value at its concrete dynamic type; type switches are open, not exhaustive like a match on a Rust enum. |
| [`12_tuples.go`](12_tuples.go) | 5 | A stored heterogeneous Rust tuple usually becomes a named struct in Go, and a keyed literal (`Person{Name: ..., Age: ...}`) labels each field so the printed value stays meaningful; a fixed-length group of one type is an array `[N]T`. Multiple results are a different thing: they are a call-level feature, not a value. Where Rust returns one tuple, Go hands back separate results that can flow straight into another call's argument list — but there is no `(int, int)` value to store, which is why stored data needs the struct. |
| [`13_variadic.go`](13_variadic.go) | 3 | A variadic parameter arrives as a `[]T`, built from separate call arguments; the `slice...` spread passes an existing slice straight through without allocating a new one. A plain `[]T` parameter is usually closest to Rust's `&[T]`. |
| [`14_interfaces.go`](14_interfaces.go) | 5 | A non-empty interface is two words: the type word points to an `itab` (the concrete type plus its method table) and the data word points to or holds the value. Satisfaction is structural, checked at compile time on assignment: the pointer-receiver `Reset` belongs only to `*Counter`'s method set, so `&counter` satisfies `interface{ Reset() }` and a plain `Counter` does not. A call dispatches indirectly through the `itab`; devirtualization can make it direct and enable inlining, and converting an escaping value can allocate a box — so prefer concrete types on hot paths and benchmark. |
| [`15_typed_nil.go`](15_typed_nil.go) | 5 | An interface is `== nil` only when both words are nil. Assigning a nil `*problem` sets the type word to `*problem` but leaves the data word nil, so the interface is `!= nil` even though the stored pointer is nil — the classic typed-nil `error` return bug (a function returning a typed-nil pointer as `error` looks like it returned nil but does not); return `nil`, not the typed-nil pointer. Safe Rust has no analogue: it has no null references, and `Option<Box<dyn Error>>` distinguishes `None` from a present box. |
| [`16_printing.go`](16_printing.go) | 4 | `fmt.Print` writes arguments as-is, `Println` adds spaces and a trailing newline, and `Printf` uses verbs with an explicit newline; `%v`, `%+v` (field names), and `%#v` (Go syntax) do not map cleanly onto Rust's `Display` vs `Debug`. `Fprintln` and friends take an `io.Writer` to choose the destination, but stdout/stderr ordering is not guaranteed when captured separately. |
| [`17_defer.go`](17_defer.go) | 5 | `defer` schedules calls to run at function return in LIFO order — on both normal return and panic unwinding, unlike Rust's scope-bound `Drop`; a deferred call's arguments are evaluated immediately at the `defer` statement, while a deferred closure reads the variables' values at run time. |
| [`18_sorting.go`](18_sorting.go) | 4 | `slices.Sort` orders comparable elements naturally, while `slices.SortFunc` takes a comparator for custom or multi-key order — `cmp.Or` chains comparisons and returns the first nonzero, like Rust's `sort_unstable_by`. Both are unstable; reach for `SortStableFunc` when equal elements must keep their input order. |
| [`19_maps_sets.go`](19_maps_sets.go) | 5 | `map[K]V` is Go's built-in hash map and `map[T]struct{}` its zero-payload set; iteration order is deliberately unspecified, so sort `maps.Keys` (allocating and sorting a key slice) when output must be deterministic. Go has no ordered-map type — keep a `[]K` beside the map, or use a third-party collection, when insertion order is part of the data model. |
| [`20_heap.go`](20_heap.go) | 4 | `container/heap` supplies the heap algorithm and calls back into a type that provides storage plus `Len`/`Less`/`Swap`/`Push`/`Pop`; `Less` picks the direction, so this is a min-heap rather than Rust's max-by-default `BinaryHeap`. Note `heap.Pop` moves the root to the end and the adapter's `Pop` just trims that slot — it does not search. |
| [`21_concurrent_maps.go`](21_concurrent_maps.go) | 5 | Default to a typed `map` guarded by a `sync.Mutex` (like Rust's `Arc<Mutex<HashMap>>`): it keeps static types and can protect invariants spanning several operations — an unsynchronized map written while another goroutine touches it is a data race. `sync.Map` is specialized for write-once/read-many caches or disjoint-key access and stores `any`, so `Load` loses the static type. The `WaitGroup` pattern: `Add` before `go`, `defer Done`, `Wait` to join. |
| [`22_synchronization.go`](22_synchronization.go) | 5 | `RWMutex` allows many readers or one writer; typed atomics make a single value's update indivisible (sequentially consistent — no ordering argument, unlike Rust); `atomic.Pointer` publishes an immutable snapshot to readers. Prefer a `Mutex` for multi-field invariants and `RWMutex` only after measured read contention; a stored snapshot is GC-reclaimed once readers release it, so use an explicit lifetime protocol when snapshots own files or sockets. |
| [`23_channels.go`](23_channels.go) | 5 | A channel is like Rust's mpsc — send and receive copy a value — but also does synchronization: an unbuffered channel blocks each side until a peer is ready, so routing every update through one goroutine gives it sole ownership of the state. A buffered channel holds sends in FIFO order until full; `range` drains a closed channel and stops, and a receive from a closed channel returns the zero value and `false` rather than reporting disconnection. Sending copies the value (a slice copies only its header). |
| [`24_ring_buffer.go`](24_ring_buffer.go) | 4 | A generic bounded ring deque (`VecDeque`-shaped) reuses one fixed array with O(1) end operations and explicit full/empty return values; unlike a channel it is storage only — it never blocks and is not goroutine-safe. Add a `Mutex` around its methods when shared, use a channel when callers should block, and reach for a different atomic algorithm for lock-free SPSC/MPMC. |
| [`25_barriers.go`](25_barriers.go) | 4 | A barrier both counts arrivals and releases every waiter at once; Go composes it from a `WaitGroup` (the counter) plus closing a channel (the one-to-many broadcast), which keeps each role explicit. It is one-shot, unlike Rust's reusable `Barrier` — cyclic barriers need `sync.Cond`. Pick the primitive that states the rule: `Mutex` for invariants, atomics to publish a value, channels to transfer/backpressure/signal, `WaitGroup` to join, `Once` to init, `Cond` to await a repeated condition. |
| [`26_context.go`](26_context.go) | 5 | A `context.Context` threads cancellation explicitly through a call tree — conventionally the first argument, like passing a cancellation token — where Rust's std has no direct equivalent; `cancel()` closes `ctx.Done()`, the waiting goroutine wakes, and `ctx.Err()` reports why (`context.Canceled`). |
| [`27_errors.go`](27_errors.go) | 5 | Where Rust reaches for enum variants, Go creates sentinel errors with `errors.New` for stable identity, wraps them with `fmt.Errorf("...: %w", err)` to add context without losing that identity, and matches with `errors.Is`, which walks the whole `%w` chain — so it finds `strconv.ErrSyntax` even through `Atoi`'s `*NumError`. |
| [`28_panic.go`](28_panic.go) | 5 | An unrecovered panic terminates the whole process; `recover` catches it only when called directly by a deferred function in the same panicking goroutine, since goroutines do not swallow each other's panics and Go has no global handler. Panic is for violated invariants, not expected failures (use `error`); like Rust's `catch_unwind` it is an exceptional boundary — wrap each goroutine if you need recovery, and re-panic after logging when the process should still crash. Libraries such as `net/http` recover around handlers. |
| [`29_logging.go`](29_logging.go) | 5 | `log/slog` logs structured key/value fields rather than assembled strings, with a text or JSON handler configured for level, output, and attribute rewriting (the lesson drops the timestamp for deterministic output); `logger.With(...)` derives a child logger that carries shared attributes — a logger, not a tracing span. |
| [`30_json.go`](30_json.go) | 5 | `encoding/json` marshals and unmarshals exported fields by default, using struct tags to rename or omit them and skipping unexported fields, with no derive or generated code — where Rust's serde typically derives type-specific code; parsing happens at runtime (currently via cached reflection, an implementation detail), and `Marshaler`/`Unmarshaler` methods can override the default. Benchmark real payloads rather than assuming which approach is faster. |
| [`31_http.go`](31_http.go) | 5 | `http.ServeMux` routes by method and path pattern, capturing wildcard segments that `request.PathValue` reads; an `http.Client` sends requests and every response body must be closed. This lesson swaps the network for an in-memory `RoundTripper` so it runs in the Playground against the real handler — a real service wires the mux into an `http.Server`, sets timeouts, reuses clients, and calls `Shutdown` for graceful stop. |
| [`32_config_options.go`](32_config_options.go) | 4 | Functional options replace Rust's builder-plus-`#[derive(Default)]`: a plain struct whose defaults live in `New`, and `Option` values are closures that mutate a `*Config` applied in order over those defaults. Because Go's zero value is often already the default (e.g. `Verbose` false), many fields need no option at all. |
| [`33_http_middleware.go`](33_http_middleware.go) | 4 | Middleware is a plain `func(http.Handler) http.Handler` — no new type needed, the func value is the contract — that wraps a handler to add cross-cutting behavior, like Rust's `tower::Layer` wrapping one `Service` in another; chaining `count(tag(base))` runs outermost-first, and a middleware can close over shared state (guard it with a mutex or atomics under real concurrency). The lesson uses `httptest` to record the response in memory, deterministically. |
| [`34_generics.go`](34_generics.go) | 4 | Type parameters take constraints — interfaces used as bounds: `Number` is a type set (`~int` or `~float64`, where `~` admits named types like `Celsius`) that lets `+=` compile, and the built-in `comparable` bounds `==`/`!=` for lookups and map keys, matching Rust's `<T: Add + Copy>` and `<T: PartialEq>`. Type arguments are inferred from values (no turbofish), and `comparable` is wider than `Number` — it admits strings even though `Number` deliberately excludes them. |
| [`35_embedding.go`](35_embedding.go) | 4 | Embedding a struct by naming its type with no field name promotes the inner type's fields and methods to the outer struct. It is not inheritance — an outer `Start` shadows the inner one without overriding it — and it silently widens your type's API with everything the inner type exports, so prefer a plain named field and forward calls unless the promotion is deliberate: satisfying or wrapping an interface, or a real mixin like `sync.Mutex`. Rust has no inheritance either; you would hold a field and forward, or get default trait methods. |
| [`36_enums_iota.go`](36_enums_iota.go) | 5 | A typed int plus a `const` block is Go's enum: `iota` numbers the constants from 0 (adding 1 each line), a `String()` method satisfies `fmt.Stringer` (Go's `Display`, called automatically for `%v`/`%s`), and `1 << iota` gives each constant a distinct bit for flags combined with `or`/tested with `and`. Unlike Rust's real sum-type `enum`, a Go enum is just an int with no exhaustiveness checking — `Color(9)` is a valid value. |
| [`37_io.go`](37_io.go) | 5 | Implementing a single `Read` (or `Write`) method makes a type an `io.Reader`/`io.Writer` — the shape of Rust's `Read` trait — and you compose readers by wrapping, as `upperReader` wraps an inner reader to transform bytes as they stream; `io.Copy` pumps a reader into a writer in chunks until EOF, like `std::io::copy`. Note per-byte upper-casing is ASCII-only — real UTF-8 needs rune-aware decoding. |
| [`38_table_tests.go`](38_table_tests.go) | 4 | Tests are plain Go—no test or assertion macros: a named case table plus `t.Run` reports each input separately, while `t.Errorf` records failures and continues within the subtest and on to later cases instead of panicking at Rust's `assert_eq!`. A real same-package `*_test.go` can white-box unexported code, and `go test` supplies verbose output, coverage, benchmarks, and fuzzing; this runnable tour calls `testing.Main` only because every lesson must work with `go run FILE`. |
| [`39_generic_methods.go`](39_generic_methods.go) | 3 | Go 1.27 lets a method declare its own type parameters in addition to its receiver's: a `List[E]` method `Apply[F any]` is instantiated at two different `F`s off one `E` receiver, the shape of Rust's freely generic trait/impl methods. One restriction survives unchanged: an interface method still cannot declare type parameters, so a generic method can never belong to an interface contract. The in-browser Run fails until the public Playground supports Go 1.27. |
| [`40_json_v2.go`](40_json_v2.go) | 3 | Go 1.27 ships `encoding/json/v2` and `encoding/json/jsontext`, and plain `encoding/json` now runs on the v2 engine while keeping v1's lenient defaults. A duplicate object key still lets the last value win under `encoding/json`, but `encoding/json/v2` rejects it as an error (invalid UTF-8 in a string diverges the same way: v1 substitutes U+FFFD, v2 errors). Rust's serde is strict by default; Go picks its strictness per package version. The in-browser Run fails until the public Playground supports Go 1.27. |
| [`41_sum_types.go`](41_sum_types.go) | 4 | A Rust `enum` maps to one struct holding each variant behind its own optional pointer: pointer presence is the tag, so it cannot disagree with the payload the way a separate `Kind` field can. Handling is independent nil checks rather than a switch — none-set and both-set are ordinary states. `Match`'s one handler per variant is the only encoding Go checks exhaustively, since a new variant adds a parameter and every call site stops compiling; exactly-one is still not a compile-time property, so `Match` rejects neither-set and both-set at runtime through `(T, error)`. |
| [`42_zero_values.go`](42_zero_values.go) | 5 | Read the two `Ratio` cases together: `NewRatio(1, 0)` rejects a zero denominator, but `var zero Ratio` creates `Ratio{0, 0}` without calling `NewRatio`. That follows from Go's language rule that every declared value starts at its type's zero value; private fields cannot force construction through a function. The practical point is the `Counter` case: design a type's zero value to work when possible, because callers can always create it. |
| [`43_struct_copies.go`](43_struct_copies.go) | 5 | `copied := original` copies the `Doc` value: assigning a new title affects only `copied`, while changing `copied.Tags[0]` appears through both values because the copied slice still refers to the same backing array. The `strings.Builder` assignment also compiles; its later panic is a library runtime check, not a language ban on copying. Go assignment expresses neither a Rust move nor a deep `Clone`, so clone referenced data explicitly and keep values with identity behind pointers. |
| [`44_generic_dispatch.go`](44_generic_dispatch.go) | 3 | The output establishes one language fact: a type parameter converted to `any` retains its concrete dynamic type, so `%T` can report it and a type switch can choose behavior at run time. Go has no generic-specialization syntax; use constraints and methods for shared behavior, or an explicit type switch when behavior truly depends on the concrete type. GC-shape sharing, dictionaries, and devirtualization are compiler implementation details that this program cannot prove—inspect and benchmark the toolchain when they matter. |
| [`45_string_bytes.go`](45_string_bytes.go) | 5 | The same `"héllo"` value reports 6 bytes and 5 runes, and `s[1]` returns one byte rather than one character. More importantly, converting `41 ff 42` from `[]byte` to `string` succeeds even though `utf8.ValidString` says false; `range` then substitutes U+FFFD instead of returning an error. A Go string stores bytes, not guaranteed UTF-8 text, so validate text at input boundaries and choose deliberately between byte, rune, and user-perceived-character operations. |
| [`46_no_destructors.go`](46_no_destructors.go) | 4 | Watch the print order: the loop in `perFunction` finishes before either deferred `Close`, then the closes run when that function returns. Wrapping one iteration in its own function creates an earlier return point, so each close happens before the next open. Go has no Rust-style automatic `Drop`; cleanup is an explicit call, and `defer` attaches it to the surrounding function rather than the nearest block. Put the defer in the function whose lifetime matches the resource. |

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
| Drop references or retained storage | `clear`, `slices.Clip`, or `slices.Clone`; lesson 09 separates reachability from capacity. |
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

An optional `-tags web` build serves a browser tour of all 46 lessons: edit,
format, run, and navigate by hash (`#01`–`#46`). It embeds the README,
numbered lesson files, and `web/` assets into one binary; no lesson source is
duplicated.

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
  usually one obvious way to write something. Rust hands you more expressive
  power and, with it, more decisions.
- **Value types are real.** Structs and arrays are stored inline, not forced
  through boxing or indirection, so for value-heavy, allocation-light code
  throughput often lands close to Rust.
- **Concurrency just works.** Goroutines and channels reach roughly Tokio-class
  ergonomics without borrow-checking, lifetimes, or `Send + Sync` bounds to
  satisfy — you trade some compile-time guarantees for far less cognitive load.
- **The tradeoff, plainly.** Reaching for Rust often buys speed you do not need
  and compile-time memory safety the garbage collector already gives you at
  near-zero cognitive cost.

Python and TypeScript remain useful choices, but long-running systems can meet
runtime, packaging, or throughput limits. Rust, Java, and C# address broader or
lower-level needs, but can require more language, runtime, or architecture work
than a small service needs. Go is practical in the space between those tradeoffs.

This framing is based on [Code Like Go](https://krons.fiu.wtf/lore/go), which
calls the underlying discipline “remove to accelerate”: fewer language choices
leave fewer implementation and review decisions.

AI makes producing code easier without making generated code trustworthy. Go's
constrained language, visible control flow, `gofmt`, compiler, and tests can
make generated changes cheaper to inspect.

## License

Released into the public domain under the [`UNLICENSE`](UNLICENSE).
