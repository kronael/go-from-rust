# Minimum Viable Go for Rust Engineers

Go books and tutorials teach the language from zero. This repository answers a
narrower question for experienced Rust programmers: “I know how to do this in
Rust—what should I actually write in Go?” I made it because mapping concepts I
already knew in Rust to practical Go was the hard part.

Each numbered `.go` file contains one small standalone lesson. Thirty are
runnable; `31_bad_pointer_index.go` is the single intentional compile failure.
After working through all 31, you should be able to reason about Go slices and
allocation, translate common Rust language patterns, and use the concurrency,
error, logging, JSON, and collection idioms expected in practical Go.

## Run It

`go run .` runs the first lesson in `01_arrays_slices.go`. Pass a standalone
filename to run one focused lesson:

```sh
go run .
go run 06_append_capacity.go
go run 10_slice_edits.go
go run 25_errors.go
go run 29_logging.go
```

The standalone lessons use `//go:build ignore`, so they do not collide with the
default `main`. Explicit filenames still run normally.

Only `31_bad_pointer_index.go` is intended to fail:

```sh
go run 31_bad_pointer_index.go # expected: a *[]int cannot be indexed
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
| [`23_data_structures.go`](23_data_structures.go) | 4 | Which built-ins replace stacks, queues, sets, and heaps? |
| [`24_defer.go`](24_defer.go) | 5 | How do deferred calls differ from lexical `Drop`? |
| [`25_errors.go`](25_errors.go) | 5 | How are sentinel errors created, wrapped, and inspected? |
| [`26_panic.go`](26_panic.go) | 5 | When does a panic kill the process, and where can it recover? |
| [`27_channels.go`](27_channels.go) | 5 | What does closing and receiving from a channel mean? |
| [`28_context.go`](28_context.go) | 5 | How is cancellation passed through service code? |
| [`29_logging.go`](29_logging.go) | 5 | How does standard structured logging carry fields and context? |
| [`30_json.go`](30_json.go) | 5 | How do exported fields, tags, and omission replace serde derives? |
| [`31_bad_pointer_index.go`](31_bad_pointer_index.go) | 3 | Why can Go not index through `*[]T`? |

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
