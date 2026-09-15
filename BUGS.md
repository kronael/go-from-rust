# Bugs

## Open

### Two lessons state a claim their output cannot show

Found by an independent factual audit, 2026-09-15. Both are correct about Go;
neither demonstrates what it asserts, which is the standard the rest of the
curriculum holds to.

- **`22_synchronization.go`.** One goroutine takes and releases every lock in
  sequence, so the output is identical with every lock deleted. `RWMutex`
  concurrency, atomic visibility, and `sync.Once` exclusion are all asserted
  and none is exercised. Showing them needs real concurrent goroutines while
  keeping output deterministic and the race detector clean.
- **`38_table_tests.go`.** The README row says `t.Errorf` records a failure and
  continues where Rust's `assert_eq!` panics. Every case passes, so the run
  prints `PASS` and the continue-on-failure behaviour never happens. Showing it
  means a deliberately failing case, which makes a teaching lesson print a test
  failure — a presentation decision, not a code fix.

Both need a design call rather than an edit, which is why neither is fixed
here.
