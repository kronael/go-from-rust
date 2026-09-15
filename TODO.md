# TODO

## SIMD lesson — needs a decision

`GOEXPERIMENT=simd` is recognized by Go 1.27 but the public Playground cannot
set it, so a real `simd/archsimd` lesson could never Run in the tour. Two
options:

1. A concept lesson whose executed code is plain Go. It runs everywhere and
   teaches the mapping, but the SIMD itself stays on the page.
2. A self-hosted runner behind `PLAYGROUND_URL` that sets the experiment.
   `cmd/fakeplayground` already executes submissions and inherits
   `os.Environ()`, but it binds loopback only and has no sandbox — exposing it
   as-is is remote code execution, so this option means writing a sandbox
   first.

Nothing else is queued.
