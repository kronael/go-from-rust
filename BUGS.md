# Bugs

## Open

### 38_table_tests.go cannot show continue-on-failure

Found by an independent factual audit, 2026-09-15. The README row says
`t.Errorf` records a failure and continues where Rust's `assert_eq!` panics.
Every case in the lesson passes, so the run prints `PASS` and the
continue-on-failure behaviour never happens.

Showing it needs a case that fails on purpose. That is blocked by the gate,
not by the lesson: `make integration` runs every lesson with `go run` and
requires exit 0, and a failing test makes `testing.Main` exit nonzero. So the
options are to teach the behaviour without running it, to exempt this one
lesson from the integration loop, or to drop the claim from the description.

That is a presentation decision, which is why it is not fixed here.
