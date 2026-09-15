# TODO

## Fold the 1.27 features back into their home lessons

Now that the public Playground runs Go 1.27, lessons 39 and 40 exist only
because their syntax used to fail in the browser. Consider folding the
struct-literal promoted-field key into `12_tuples.go` (it was reverted for
exactly that reason) and the `encoding/json` v1-versus-v2 contrast into
`30_json.go`, if lesson 40 then reads as redundant.

## SIMD lesson — needs a decision

`GOEXPERIMENT=simd` is recognized by 1.27 but the public Playground cannot
set it, so a real `simd/archsimd` lesson could never Run there. Options: a
concept lesson whose executed code is plain Go, or a self-hosted sandboxed
runner (`PLAYGROUND_URL`) that sets the experiment. `cmd/fakeplayground`
executes submissions and inherits `os.Environ()`, but binds loopback only and
has no sandbox — exposing it as-is would be remote code execution.
