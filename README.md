# gofrs

> **In the age of AI, code is cheap. Code you can trust is not.**

**Go for Rustaceans, one runnable file at a time.**

`gofrs` shows what Go does where a Rustacean expects `Option`, `match`, tuples,
iterators, or `Vec` methods. Each file asks one question, prints the answer, and
stays small enough to read in one pass.

It is the executable companion to
[Code Like Go](https://krons.fiu.wtf/lore/go): boring on purpose.

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

## Why Go now

Across three randomized field experiments involving 4,867 developers, the
combined estimate found **26.08%** more completed tasks with an AI assistant,
although individual results were noisy. More code is coming. But when
researchers expanded a popular code benchmark with stronger tests, reported LLM
pass rates fell by up to **28.9%**. Generation got cheaper; verification did
not. ([productivity study](https://doi.org/10.1287/mnsc.2025.00535),
[correctness study](https://arxiv.org/abs/2305.01210))

Go was designed for people who must “write—and read and debug and maintain”
large systems. Its original targets included slow builds, uncontrolled
dependencies, inconsistent language subsets, and poor program understanding.
Those are now AI coding problems too.
([Go at Google](https://go.dev/talks/2012/splash.article))

- **One canonical format.** [`gofmt`](https://go.dev/blog/gofmt) makes generated
  diffs about logic, not style.
- **Fast, typed feedback.** The compiler, `go test`, and `go vet` reject many bad
  guesses before review.
- **Explicit control flow.** Errors, loops, and conversions stay visible in the
  source.
- **Fewer moving parts.** A broad standard library and the
  [Go 1 compatibility promise](https://go.dev/doc/go1compat) reduce dependency
  and version churn.
- **Deliberately ordinary code.** Fewer ways to be clever make unfamiliar code
  cheaper to inspect and own.

The claim is not that every LLM writes Go better than every other language. The
claim is narrower and more useful:

> **Go's advantage in the AI age is not that models can write it. It is that
> humans can still read it.**

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
