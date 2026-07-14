# gofrs

**Go for Rustaceans: twenty runnable demos and two useful compiler errors.**

The files answer Rust-shaped questions about Go. Read one, predict the output,
run it, then change one value and run it again.

This guide assumes programming and Rust fluency. It covers Go basics only where
their behavior or idiom differs from Rust.

## What you get

The demos give you a working model of:

- slices, backing arrays, aliasing, pointers, and `range`;
- the Go idioms used instead of `Option`, `match`, tuples, trait bounds, and
  Rust-style iterators;
- collection operations that need more explicit memory handling in Go;
- the compiler errors and runtime surprises most likely to catch a Rust reader.

This is lesson material, not a production collections library or a complete Go
course.

## Run it

```sh
go run .                 # arrays, slices, aliasing, and sorting
go run filter.go         # three ways to filter a slice
go run option.go         # comma-ok and nil instead of Option<T>
go run datastructures.go # Heap, Stack, Queue, and Set
```

Most demos use `//go:build ignore` so each can declare its own `main`. Run one by
passing its filename directly to `go run`.

Two demos are supposed to fail:

```sh
go run badidx.go # cannot index a pointer to a slice
go run caseA.go  # cannot dereference an int
```

`caseB.go` contrasts the second case with a slice that intentionally contains
pointers.

You need Go 1.26 or newer, matching [`go.mod`](go.mod).

## Learning path

| Demo | Rust question it answers |
|---|---|
| [`printing.go`](printing.go) | How do `println!`-style formatting and output differ? |
| [`main.go`](main.go) | How do arrays, `Vec<T>`, aliasing, and sorting translate? |
| [`demo.go`](demo.go) | What replaces `copy_from_slice` and variadic expansion? |
| [`rangekw.go`](rangekw.go) | What does `range` yield for slices, maps, strings, integers, and iterators? |
| [`ptr.go`](ptr.go) | When does Go dereference a pointer automatically? |
| [`badidx.go`](badidx.go) | Why does indexing `*[]T` fail? |
| [`filter.go`](filter.go) | What replaces `iter().filter().collect()`? |
| [`pop.go`](pop.go) | What replaces `Vec::pop`, `remove`, and `swap_remove`? |
| [`zero.go`](zero.go) | Why does deleting pointer elements need zeroing? |
| [`shrink.go`](shrink.go) | When does a short slice keep a large backing array alive? |
| [`edges.go`](edges.go) | What happens when a loop appends to the slice it visits? |
| [`option.go`](option.go) | What replaces `Option<T>`? |
| [`match.go`](match.go) | What replaces `match` and destructuring? |
| [`tuples.go`](tuples.go) | What replaces tuple values? |
| [`multireturn.go`](multireturn.go) | Are multiple return values a tuple? |
| [`ternary.go`](ternary.go) | What replaces expression-valued `if` and ternary helpers? |
| [`caseA.go`](caseA.go), [`caseB.go`](caseB.go) | When is `*v` a valid dereference? |
| [`sortstruct.go`](sortstruct.go) | What replaces `sort_by` and `sort_by_key`? |
| [`genericsort.go`](genericsort.go) | What is the Go equivalent of an `Ord`-like bound? |
| [`iters.go`](iters.go) | How do Go's push iterators differ from Rust's pull iterators? |
| [`datastructures.go`](datastructures.go) | How do generic Heap, Stack, Queue, and Set APIs look in use? |
| [`ds/ds.go`](ds/ds.go) | What do minimal teaching implementations look like? |

## Why Go

Go's omissions are part of the curriculum. Fewer constructs mean fewer ways to
express the same operation, fewer hidden control paths, and fewer decisions in
implementation and review. [Code Like Go](https://krons.fiu.wtf/lore/go) calls
this approach *remove to accelerate*.

That constraint also helps with AI-assisted code. Generated code still needs
human verification; static types, visible control flow, one formatter, and fast
compiler feedback make that review easier.

For APIs, workers, CLIs, and data plumbing, Go can be a useful middle ground:
less low-level design work than Rust, no Python or JavaScript runtime to ship,
and a smaller language surface than Java or C#.

## License

Released into the public domain under the [`UNLICENSE`](UNLICENSE).
