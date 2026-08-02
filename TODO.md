# TODO

## Switch when the public Playground supports Go 1.27

The tour runs code on the public Go Playground
(`defaultPlaygroundURL` in `internal/webapp/server.go`), which is still
pre-1.27. Lessons 43 and 44 use 1.27-only syntax, so their in-browser Run
fails; everything else runs. Check with:

```sh
make playground-check
```

When it reports READY:

1. Drop the "Needs Go 1.27 — this lesson's in-browser Run fails …" sentence
   from the lesson 43 and 44 rows in `README.md`.
2. Delete `go127Note` / `needsGo127` and their use in `extractOutput`
   (`web/app.js`) — the fallback message is then dead code.
3. Optionally fold the 1.27 features back into the lessons they belong to,
   now that Run works: the struct-literal promoted-field key into
   `15_tuples.go` (it was reverted for exactly this reason), and the
   `encoding/json` v1-versus-v2 contrast into `34_json.go` if lesson 44
   then reads as redundant.
4. Re-run `make check test integration` and redeploy.

Not blocked on the Playground: pointing `PLAYGROUND_URL` at any 1.27-capable
backend also makes Run work — no code change needed.

## Move off the release candidate

`go.mod` pins `toolchain go1.27rc2` and `Dockerfile` builds on
`golang:1.27rc2-alpine3.23`. When Go 1.27 goes GA, bump both to the final
release and drop the "currently the `go1.27rc2` release candidate" wording
from `README.md`.

## SIMD lesson — needs a decision

`GOEXPERIMENT=simd` is recognized by 1.27 but the public Playground cannot
set it, so a real `simd/archsimd` lesson could never Run there. Options: a
concept lesson whose executed code is plain Go, or a self-hosted sandboxed
runner (`PLAYGROUND_URL`) that sets the experiment. `cmd/fakeplayground`
executes submissions and inherits `os.Environ()`, but binds loopback only and
has no sandbox — exposing it as-is would be remote code execution.
