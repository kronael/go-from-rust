# Bugs

## Open

### Four lesson pairs teach the same point twice (proposal, needs sign-off)

Audit of all 50 lessons, 2026-08-08. Each pair below is one concept split
across two files rather than two concepts:

- **09 + 10.** `09_filter.go` shows append-a-new-slice versus
  `slices.DeleteFunc` in place. `10_slice_edits.go` already covers delete,
  insert, pop, and swap-delete on the same footing; `DeleteFunc` is one more
  `slices` call in that family.
- **11 + 12.** Both teach "reslicing leaves the backing array reachable, so
  clear or clone to drop it." 11 demonstrates it with pointer elements, 12
  with a large array plus `slices.Clip`. Same claim, two props. `slices.Clone`
  is taught a third time here after 02.
- **07 + 08.** `08_pointer_elements.go` is 24 lines restating 07's
  pointer-as-optional on a slice: nil element means absent, deref panics.
  It is the repo's thinnest lesson and carries value 3.
- **15 + 16.** "Go has no tuple" split in half — 15 gives the struct and
  `[N]T` replacements, 16 gives multiple returns as a call-level feature.
  Together they are one lesson, 35 lines total.

Merging all four takes the curriculum from 50 to 46. Not done, because
numeric prefixes must stay contiguous, so every lesson after 07 renumbers:
50 README rows, every cross-reference, and the web tour's `#NN` hashes. That
is a curriculum decision, not a cleanup.

Lesser overlaps, not worth a merge: `47_struct_copies.go` opens by re-showing
01's slice aliasing before reaching its own point (trim, do not merge), and
21 and 50 both open against Rust's `Drop` while making different claims about
`defer` (ordering versus scope) — keep both.

### fakeplayground's 10s run timeout does not actually bound a run

`cmd/fakeplayground/main.go` wraps `go run` in a 10-second
`context.WithTimeout` and collects output into `bytes.Buffer`s. Killing
`go run` does not kill the compile processes it spawned, and those children
keep the stdout/stderr pipes open, so `cmd.Run()` blocks past the deadline
waiting for the copy goroutines. Observed: a cold `FAKEPLAYGROUND_GOCACHE`
made `35_http.go` compile `net/http` from scratch; the request ran 70 seconds
and then returned an empty reply instead of a timeout response after 10.

Effect: `make integration` fails the two `HTTP lesson` assertions whenever the
fakeplayground build cache is cold. Warm cache passes. Test-only tool, never
the public service, so the blast radius is a flaky local gate.

Fix direction (needs sign-off — changes process handling): put the child in
its own process group and signal the group, or use `cmd.StdoutPipe` with an
explicit read deadline rather than `bytes.Buffer` + `Run`. Pre-warming the
cache in `scripts/playtest.sh` would hide the symptom without bounding
anything.
