# Bugs

## Open

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
