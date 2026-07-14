# gofrs

**Go for Rustaceans:** small, runnable Go programs for programmers who already
know Rust.

This is not a book, framework, or application. Each `.go` file isolates one
language concept and puts the Rust idea beside the Go equivalent. Read a file,
predict its output, then run it.

```sh
go run .                 # arrays, slices, aliasing, and sorting
go run filter.go         # three ways to filter a slice
go run option.go         # comma-ok and nil instead of Option<T>
go run datastructures.go # Heap, Stack, Queue, and Set
```

Most demos have `//go:build ignore` so they can each declare their own `main`
without colliding with the others. Passing a filename directly to `go run`
runs that one demo.

## Who this is for

You should already be comfortable with Rust syntax and concepts such as
`Vec<T>`, `Option<T>`, `Result<T, E>`, iterators, traits, and ownership. The
demos teach the Go differences; they do not reteach programming from zero.

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

## The idea behind the project

Go often answers a missing Rust feature with direct control flow: write the
loop, return the second value, define the small struct, check the error. That is
not merely missing syntax. It reduces the number of ways a codebase can express
the same operation.

The project follows the argument in
[Code Like Go](https://krons.fiu.wtf/lore/go): Go teaches through subtraction.
Fewer language features mean fewer hidden paths and fewer local style choices.
The result is deliberately ordinary code whose behavior is visible in the
source.

These demos make that trade concrete. They show both what Go removes and what
the programmer must write instead.

## Why Go works well with LLMs

The same constraints that help a Rust programmer learn Go also help a coding
model produce code that a human can inspect:

- **A smaller solution space.** Fewer equivalent language constructs make
  generated code more consistent.
- **One canonical format.** `gofmt` removes formatting choices from prompts,
  output, and review.
- **Explicit control flow.** Loops, error checks, and conversions remain in the
  source instead of hiding behind exceptions or implicit behavior.
- **Fast, typed feedback.** `go test`, `go vet`, and the compiler quickly reject
  invented names, wrong types, unused imports, and many incomplete edits.
- **A broad standard library.** Common HTTP, JSON, I/O, testing, and concurrency
  work needs fewer dependencies and fewer version-specific APIs to guess.
- **Readable generated code.** Humans can audit plain code more reliably than a
  dense chain of abstractions.

This is an inference from the essay's “boring on purpose” principle, not a claim
that an LLM makes Go code correct. Models still invent APIs, mishandle edge
cases, and need tests and review. Go is useful here because it shortens the
generate–compile–inspect loop and makes mistakes easier to see.

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
