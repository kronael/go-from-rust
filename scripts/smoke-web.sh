#!/bin/sh
# Smoke-tests the built web tour end to end against an in-repo fake
# Playground, never the public one. Exits nonzero on any failed assertion.
set -eu

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
WORK_DIR="$REPO_ROOT/.tmp_check/smoke-web"
FAKE_LOG="$WORK_DIR/fakeplayground.log"
WEB_LOG="$WORK_DIR/web.log"
SESSION="smoke-web-$$"
FAIL=0

mkdir -p "$WORK_DIR"
: > "$FAKE_LOG"
: > "$WEB_LOG"

cleanup() {
  agent-browser --session "$SESSION" close >/dev/null 2>&1 || true
  [ -n "${WEB_PID:-}" ] && kill "$WEB_PID" >/dev/null 2>&1 || true
  [ -n "${FAKE_PID:-}" ] && kill "$FAKE_PID" >/dev/null 2>&1 || true
}
trap cleanup EXIT

wait_for_addr() {
  log_file="$1"
  tries=50
  while [ "$tries" -gt 0 ]; do
    addr="$(grep -o 'listening on [^ ]*' "$log_file" 2>/dev/null | tail -1 | cut -d' ' -f3)"
    if [ -n "$addr" ]; then
      echo "$addr"
      return 0
    fi
    tries=$((tries - 1))
    sleep 0.2
  done
  echo "smoke-web: timed out waiting for $log_file to report a listen address" >&2
  return 1
}

assert() {
	desc="$1"
	shift
	if "$@" >/dev/null; then
		echo "ok - $desc"
  else
    echo "FAIL - $desc"
    FAIL=1
	fi
}

contains() {
	case "$1" in
		*"$2"*) return 0 ;;
		*) return 1 ;;
	esac
}

"$REPO_ROOT/dist/fakeplayground" >"$FAKE_LOG" 2>&1 &
FAKE_PID=$!
FAKE_ADDR="$(wait_for_addr "$FAKE_LOG")"

ADDR=127.0.0.1:0 PLAYGROUND_URL="http://$FAKE_ADDR/compile" \
  "$REPO_ROOT/dist/go-from-rust-web" >"$WEB_LOG" 2>&1 &
WEB_PID=$!
WEB_ADDR="$(wait_for_addr "$WEB_LOG")"
BASE_URL="http://$WEB_ADDR"

curl -sf "$BASE_URL/health" >/dev/null

ab() { agent-browser --session "$SESSION" "$@"; }

wait_contains() {
	selector="$1"
	needle="$2"
	tries=200
	while [ "$tries" -gt 0 ]; do
		contains "$(ab get text "$selector")" "$needle" && return 0
		tries=$((tries - 1))
		sleep 0.1
	done
	return 1
}

is_visible() {
	[ "$(ab is visible "$1")" = "true" ]
}

ab open "$BASE_URL/#01" >/dev/null
assert "loads lesson 01 title" contains "$(ab get text '#lesson-title')" "01_arrays_slices.go"

ab click "#next-btn" >/dev/null
assert "next button navigates to lesson 02" contains "$(ab get url)" "#02"

ab fill "#editor" 'package main

func main() {}
' >/dev/null
ab reload >/dev/null
assert "edited source persists across reload" contains "$(ab get value '#editor')" "func main() {}"

ab click "#format-btn" >/dev/null
assert "format button reports success" wait_contains "#status-msg" "Formatted"

ab click "#reset-btn" >/dev/null
assert "reset restores original source" contains "$(ab get value '#editor')" "type point struct"

ab click "#run-btn" >/dev/null
assert "run completes" wait_contains "#status-msg" "Run complete"
assert "run captures stdout" contains "$(ab get text '#output-code')" "Printf:"
assert "run captures stderr" contains "$(ab get text '#output-code')" "Fprintln: stderr"

ab open "$BASE_URL/#35" >/dev/null
ab click "#run-btn" >/dev/null
assert "lesson 35 completes" wait_contains "#status-msg" "Expected compiler error"
assert "lesson 35 run surfaces the intentional compiler failure" \
	contains "$(ab get text '#output-code')" "cannot index"
assert "lesson 35 labels the compiler failure as expected" \
	contains "$(ab get text '#status-msg')" "Expected compiler error"

ab set viewport 390 844 >/dev/null
ab open "$BASE_URL/#01" >/dev/null
assert "editor stays visible at mobile width" is_visible "#editor"
assert "run button stays visible at mobile width" is_visible "#run-btn"

if [ "$FAIL" -ne 0 ]; then
  echo "smoke-web: one or more assertions failed" >&2
  exit 1
fi
echo "smoke-web: all assertions passed"
