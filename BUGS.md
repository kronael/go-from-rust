# Bugs

## Open

### README names make targets that do not exist

Confirmed at 0241549, 2026-09-26. The Interactive Web Tour section says
`make playtest` runs the browser playtest and `make full` runs the whole gate.
The Makefile has neither target: the playtest runs inside `make integration`,
and the full gate is `make check test integration`, as `CLAUDE.md` says.

### playground-check comment is stale

Confirmed at 0241549, 2026-09-26. The comment above `playground-check` says
lessons 43 and 44 use Go 1.27 syntax, but the target probes 39 and 40. It also
says to follow `TODO.md` to drop the caveats on READY, and so does the READY
message. The caveats were dropped in v0.6.0 and `TODO.md` holds no such steps.
