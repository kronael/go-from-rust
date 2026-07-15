---
status: shipped
---

# Interactive Web Tour

## Goal

Ship a production-ready interactive version of the 35-lesson curriculum without
forking the official Go Tour, introducing a JavaScript build system, duplicating
lesson source, or executing untrusted code on this server.

## Deliverables

### 1. Embedded lesson catalog

- **Files**: `01_arrays_slices.go`, `web.go`, `internal/webapp/lessons.go`,
  `internal/webapp/lessons_test.go`
- **Accept**: `go run .` still runs lesson 01; `go run -tags web .` starts the
  web server; the web build embeds all numbered `.go` files and derives their
  value and question from the README learning table; missing, duplicate, or
  non-contiguous lessons fail startup.
- **Notes**: the root numbered files remain the only source of lesson code.
  Enumerate embedded files matching `[0-9][0-9]_*.go`. Match each file to
  exactly one README row of the literal form
  ``| [`NN_name.go`](NN_name.go) | V | Q |`` where link text and target equal
  the filename, `NN` is the two-digit lesson ID, `V` parses as an integer, and
  `Q` is retained as the question string. Sort by `NN`, require IDs 1 through
  the row count without gaps, and reject unmatched files, unmatched rows,
  duplicate IDs, duplicate filenames, or unequal link text/targets. The loaded
  catalog is immutable and safe for concurrent reads.
  `V` is the README's editorial importance score (3–5); it never controls
  ordering and is displayed as `Value V/5` in the lesson header and navigation.

### 2. Guarded web API

- **Files**: `internal/webapp/server.go`, `internal/webapp/server_test.go`
- **Accept**: `/api/lessons` returns the ordered embedded catalog;
  `/api/format` runs `go/format`; `/api/run` proxies the current Playground
  protocol without local execution; `/health` and `/ready` return useful
  status; unsupported methods and malformed or oversized input return stable
  JSON errors. `/api/format` accepts `{"body":"<source>"}` and returns
  `{"body":"<gofmt source>"}` with HTTP 200.
- **Notes**: `/api/run` accepts JSON `{"body":"<source>"}` up to 128 KiB,
  then POSTs `application/x-www-form-urlencoded` fields `version=2`, `body`, and
  `withVet=true` to `PLAYGROUND_URL` (default
  `https://go.dev/_/compile`). A successful upstream response must be HTTP 200
  JSON and is relayed unchanged up to 1 MiB. The upstream client timeout is 12
  seconds. Permit at most 8 concurrent cache misses and 120 upstream requests
  per rolling minute. Cache at most 256 successful results for 10 minutes by
  source hash; cache hits consume neither concurrency nor run budget. Attach
  `PLAYGROUND_USER_AGENT`, defaulting to
  `go-from-rust/1 (+https://github.com/kronael/go-from-rust)`. Log every request
  with a request ID. `/health` reports only process liveness; `/ready` reports
  whether the embedded catalog and handlers initialized and never probes the
  upstream service. Every API failure uses
  `{"error":{"code":"<machine_code>","message":"<human message>"}}`.
  Unsupported methods return 405 `method_not_allowed`; malformed JSON returns
  400 `invalid_request`; source above 128 KiB returns 413 `source_too_large`;
  unformattable source returns 422 `format_failed`; an occupied 8-run semaphore
  returns 503 `runner_busy` with `Retry-After: 1`; an exhausted run budget
  returns 429 `run_budget_exhausted` with `Retry-After: 60`; upstream timeout
  returns 504 `upstream_timeout`; and upstream status, size, or JSON failures
  return 502 `upstream_failed`. Cache eviction is deterministic LRU after
  expired entries are removed. The semaphore, rolling timestamp budget, and
  LRU cache are process-global fields on one server instance; mutexes protect
  budget and cache state under concurrent `net/http` handlers.
  `web.go` embeds `README.md`, `[0-9][0-9]_*.go`, and `web/*` in the same
  binary and passes that filesystem into `internal/webapp`. GET `/` serves the
  embedded `web/index.html`; GET `/static/app.js` and `/static/styles.css`
  serve the corresponding embedded assets with correct content types and
  cache headers. Unknown paths return 404; hash routing requires no fallback.

### 3. Dependency-free browser experience

- **Files**: `web/index.html`, `web/app.js`, `web/styles.css`
- **Accept**: a reader can navigate all lessons, edit code, run, format, reset,
  move previous/next, deep-link by lesson number, and retain edits in local
  storage; compiler errors and the intentional failure are clearly rendered;
  keyboard operation, narrow screens, light mode, and dark mode remain usable.
- **Notes**: route entirely in the browser with hash fragments `#01` through
  `#35`; the server needs no per-lesson routes. Use semantic HTML and vanilla
  JavaScript; no npm, CDN, framework, analytics, cookies, or remote font
  dependency.

### 4. Build and deployment

- **Files**: `Makefile`, `Dockerfile`, `.dockerignore`
- **Accept**: `make build`, `make test`, and `make check` are deterministic;
  the multi-stage container runs as non-root, includes CA certificates, exposes
  port 8080, and checks `/health`.
- **Notes**: pin the build stage to `golang:1.26.0-alpine3.23` and the runtime
  stage to `alpine:3.23.3`; keep the runtime image minimal. The server listens
  on `ADDR`, default `127.0.0.1:3999`; the container sets `ADDR=:8080`. The
  runtime image copies only the compiled binary because lessons and web assets
  are embedded.

### 5. Documentation and browser proof

- **Files**: `README.md`, `CLAUDE.md`, `specs/1/web-tour.md`,
  `cmd/fakeplayground/main.go`, `scripts/smoke-web.sh`
- **Accept**: local use, environment variables, Playground policy, security
  boundary, container deployment, and web validation are documented.
  `make smoke-web` starts the in-repository fake Playground and web server on
  dynamic loopback ports, then uses the `agent-browser` CLI with a fresh session
  to assert navigation, edited-code persistence after reload, formatting,
  reset, successful execution, lesson 35 expected compiler failure, and the
  mobile layout at 390x844. The script exits nonzero on any failed assertion
  and never contacts the public Playground.

## Constraints

- Keep the title **Minimum Viable Go for Rust Engineers**.
- Keep every numbered lesson independently runnable by filename.
- Do not fork or copy the official Tour frontend.
- Do not execute submitted programs locally in production or development.
- Do not add third-party Go or browser dependencies.
- The official Playground upstream must be configurable and replaceable in tests.
- Do not send test traffic to the public Playground.
- HTTP server limits are 5 seconds to read headers, 10 seconds to read a
  request, 20 seconds to write a response, and 60 seconds idle.

## Verification

- [x] `make build`
- [x] `make test`
- [x] `make check`
- [x] `go run .`
- [x] all 34 runnable lessons execute and lesson 35 fails as intended
- [x] web API tests use an in-process fake Playground
- [x] browser smoke passes against the built server and fake Playground
- [x] container image builds and its health check succeeds
- [x] no tracked generated lesson copy exists
