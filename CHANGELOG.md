# Changelog

## [v0.4.0] — 20260808

> Minimum Viable Go v0.4.0 — what Go does not enforce
>
> Five new lessons on guarantees Rust makes and Go does not. None of them has a fix to teach, so each shows the behaviour and stops.
>
> • The zero value cannot be forbidden — `var x T` always compiles, so no type can require initialization
> • Struct copies are silent, shallow, and unstoppable: a copied `strings.Builder` panics on the next write
> • Go generics stencil by GC shape rather than monomorphize, so specialization is a runtime type switch
> • A Go string is bytes with no UTF-8 guarantee, and decoding substitutes rather than fails
> • There are no destructors, and `defer` is function-scoped, not block-scoped
>
> Full notes: CHANGELOG.md

- Lessons 46-50 added under a new curriculum group, "what Go does not enforce."
  The five fill mappings the tour had no entry for: `46_zero_values.go`,
  `47_struct_copies.go`, `48_generic_dispatch.go`, `49_string_bytes.go`, and
  `50_no_destructors.go`.
- Each demonstrates rather than asserts, and each stays inside the gate:
  no `go vet` violation (the check vets, so copylocks could not be shown
  directly — `strings.Builder`'s address check gives 47 a real runtime failure
  instead), no finalizer timing, no benchmarks, deterministic output.
- Twelve README rows had grown from summary into a second copy of the lesson.
  The longest fell from 1028 to 592 characters, cut by whole clauses rather
  than by compressing the language.
- Dropped the Why Go claim that Go deployments skip the load-balancing
  machinery "a Rust service often accretes." It states a tradeoff as a
  universal performance claim, which this repository's own contributor rules
  forbid.
- Logged four duplicate lesson pairs in `BUGS.md` — 09+10, 11+12, 07+08, and
  15+16 each teach one concept across two files. Merging them would take the
  curriculum to 46 and renumber every lesson after 07, so it waits for a
  decision rather than shipping as a cleanup.

## [v0.3.1] — 20260808

> Minimum Viable Go v0.3.1 — the exhaustive matcher, shown
>
> Lesson 45 described a matcher that gets real compile-time exhaustiveness but never showed one. It does now.
>
> • `Match` takes one handler per variant — a new variant adds a parameter, so every call site stops compiling
> • Its cost is on screen too: neither-set and both-set need an explicit answer
> • README no longer claims Go has no exhaustiveness "for any encoding" and then describes one
>
> Full notes: CHANGELOG.md

- Lesson 45 now shows `Match[R any]` alongside the independent nil checks.
  One handler per variant is the single construction Go checks exhaustively:
  adding a variant adds a parameter, so every call site fails to compile until
  it handles the new case. The lesson prints all four states through it, which
  puts the cost on screen — presence-as-tag admits neither-set and both-set, so
  the matcher must be told what they mean, and a `return` inside a handler
  leaves only that handler.
- The lesson's README description contradicted itself: it stated Go has no
  compile-time exhaustiveness "for any encoding of this" and then described an
  encoding that has it. Rewritten around what the code now demonstrates.
- `Match` returns `(R, error)` and rejects neither-set and both-set instead of
  folding them into a catch-all handler. Exactly-one is not a compile-time
  property, so a matcher that assumes it has to say when it does not hold —
  reported, fittingly, through the same unenforced two-variant sum.
- Tightened lesson 45's inline comments and its description. Same claims, fewer
  words: the comments label the code, the description carries the depth.

## [v0.3.0] — 20260806

> Minimum Viable Go v0.3.0 — sum types without a Kind field
>
> A new lesson on modelling Rust enums in Go, plus the reason Go's own `(T, error)` is a sum type nothing enforces.
>
> • Lesson 45 — variant structs behind optional pointers; presence is the tag, so it cannot disagree with the payload
> • Independent nil checks over a switch: none-set and both-set are ordinary states, not bugs
> • Covers the alternatives — sealed interfaces, constructors, and the one encoding that gets real exhaustiveness
>
> Full notes: CHANGELOG.md

- Added lesson 45 (sum types). A Rust `enum` maps to one struct holding each
  variant behind its own optional pointer, with no separate `Kind`
  discriminant: pointer presence *is* the tag, so it can never contradict the
  payload the way a flat struct with a `Kind` field plus inlined variant fields
  can. Handling is independent nil checks rather than a type switch, because
  none-set and both-set are ordinary states whenever nothing requires exactly
  one variant.
- The lesson's description names the costs honestly — a pointer indirection and
  allocation per present variant, no compiler guarantee that only one is set —
  and the alternatives it does not show: a sealed interface (unexported marker
  method) closes the variant set at the package boundary but boxes and still
  gives no exhaustiveness; unexported fields plus constructors is what actually
  prevents both-set from outside the package; a `Match[R any]` matcher is the
  only construction that gets compile-time exhaustiveness, at the cost of
  losing `return` from the enclosing function.
- Recorded a known defect in `BUGS.md`: `cmd/fakeplayground`'s 10-second run
  timeout does not bound a run, because killing `go run` leaves its compile
  children holding the output pipes. A cold build cache makes `make
  integration` fail the two `HTTP lesson` assertions; a warm one passes.

## [v0.2.1] — 20260802

> Minimum Viable Go v0.2.1 — honest failures
>
> The two Go 1.27 lessons now explain why in-browser Run fails, and the checks that were supposed to catch mistakes actually catch them.
>
> • Lessons 43 and 44 explain the Playground is behind — but only for errors it really caused, not your typos
> • `make check` no longer passes on Go it could not even parse
> • `make playground-check` tells you when the Playground reaches 1.27, and says "cannot tell" instead of guessing when offline
>
> Full notes: CHANGELOG.md

- `make check` silently passed on unparseable Go: `gofmt` reports parse errors on
  stderr and prints nothing on stdout, so `test -z` on its output saw an empty
  string. It now fails on a nonzero `gofmt` exit as well.
- `make playground-check` reported READY whenever the request failed (network
  down, DNS, non-2xx), because nothing matched the expected error text. A failed
  request now reports "cannot tell" and exits 2. It also probes lesson 44 —
  `encoding/json/v2` availability is a separate condition from generic-method
  syntax — and writes under the project's `tmp/` rather than assuming `GOCACHE`
  exists.
- The Go 1.27 note in the tour was prepended to *every* compile error on lessons
  43 and 44, so a typo the reader introduced was reported as the Playground being
  behind. It now also requires the error to match a known 1.27 signal.
- README advertised "all 42 lessons" in two places; there are 44.

## [v0.2.0] — 20260802

> Minimum Viable Go v0.2.0 — Go 1.27, and two lessons for it
>
> The repo now builds on Go 1.27, with new lessons for generic methods and the stricter `encoding/json/v2`.
>
> • Lesson 43 — a method can declare its own type parameters; interface methods still cannot
> • Lesson 44 — `encoding/json` keeps its lenient defaults, `encoding/json/v2` rejects duplicate keys
> • Toolchain moved to Go 1.27 (rc2) — go.mod, Docker image, and the format check
> • All 42 existing lessons still run in the browser unchanged
>
> Full notes: CHANGELOG.md

- Adopted Go 1.27 (`go1.27rc2`, ahead of GA): `go.mod` declares `go 1.27` with a
  `toolchain` line, and the Docker builder moved to `golang:1.27rc2-alpine3.23`.
- Added lesson 43 (generic methods) and lesson 44 (`encoding/json/v2`), the two
  lessons that carry 1.27-only syntax. Their descriptions state plainly that the
  in-browser Run fails until the public Go Playground supports 1.27.
- Kept lessons 01–42 on 1.26-compatible syntax so every one of them still runs
  in the browser. 1.27 features were briefly added to lessons 15, 34 and 38 and
  reverted when that broke their Run button — the tour executes against the
  public Playground, which is still pre-1.27.
- `make check` and the documented validation gate now use the active toolchain's
  `gofmt`; an older `gofmt` on `PATH` cannot parse generic methods and reported a
  parse error while printing an empty list, so the check silently passed.

## [v0.1.1] — 20260726

> Minimum Viable Go v0.1.1 — a real testing lesson
>
> The table-test lesson is now a complete, runnable Go test you can lift straight into a `_test.go`.
>
> • Lesson 42 defines a real `TestParsePort` with `t.Run` subtests over a realistic port parser
> • Description explains why Go testing differs from Rust — no macros, `t.Errorf` accumulates, batteries-included `go test`
>
> Full notes: CHANGELOG.md

- Rewrote lesson 42 (table-driven tests) as a complete example: a real
  `TestParsePort(t *testing.T)` with `t.Run` subtests and `t.Errorf` over a
  realistic `parsePort` (normalize, validate, error, boundary cases); `main`
  calls `testing.Main` as a tour adapter so it runs under `go run` while the
  test body lifts verbatim into a `*_test.go`. The description now explains why
  Go testing differs from Rust: plain Go (no macros), per-case `t.Run` names,
  `t.Errorf` accumulates instead of panicking like `assert_eq!`, white-box
  same-package testing, and batteries-included `go test`.

## [v0.1.0] — 20260726

> Minimum Viable Go v0.1.0 — Go, mapped from Rust
>
> 42 small runnable Go lessons that map what you know in Rust to idiomatic Go, with a browser tour to run each one.
>
> • 42 lessons — slices, language mappings, collections, concurrency, control, services, generics
> • Interactive tour at krons.cx/go-from-rust — edit, format, and run each lesson in the browser
> • Reads on mobile — code fits with no horizontal scroll; concepts live in styled descriptions
> • Every lesson is one self-contained runnable file with a direct Rust→Go mapping
>
> Full notes: CHANGELOG.md

- 42 lessons (`01`–`42`) across seven groups: foundations and slices, language
  mappings, collections and ordering, concurrency, control and failure, services
  and data, generics and tooling.
- Interactive web tour (`-tags web`): per-lesson editor with run/format against a
  Playground stand-in; deployed at `krons.cx/go-from-rust`, discoverable from the
  krons portal, the guides list, and the go-lore essay.
- Mobile-first readability: every lesson wraps to ≤60 chars and the editor uses a
  small font, so code fits with no horizontal scroll; lesson descriptions carry
  the concept and Rust contrast, rendered with inline code chips.
- `defer` taught at lesson 21 (before its first use), `option` focused on
  representing optionality idiomatically, and struct embedding flagged as a niche
  tool rather than a default.
