# Changelog

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
