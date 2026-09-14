#!/bin/sh
# Playtests the built web tour end to end against an in-repo fake
# Playground, never the public one. Exits nonzero on any failed assertion.
set -eu

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
WORK_DIR="$REPO_ROOT/.tmp_check/playtest"
FAKE_LOG="$WORK_DIR/fakeplayground.log"
WEB_LOG="$WORK_DIR/web.log"
SESSION="playtest-$$"
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
  echo "playtest: timed out waiting for $log_file to report a listen address" >&2
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

not_contains() {
	! contains "$1" "$2"
}

# fakeplayground bounds each run at 10 seconds. A cold build cache makes
# the HTTP lesson exceed that on its first compile, so warm the cache the
# server will use before any assertion depends on a run finishing.
FAKE_GOCACHE="${FAKEPLAYGROUND_GOCACHE:-${TMPDIR:-/tmp}/fakeplayground-gocache}"
GOCACHE="$FAKE_GOCACHE" go build -o /dev/null "$REPO_ROOT/31_http.go"

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

wait_value_contains() {
	selector="$1"
	needle="$2"
	tries=200
	while [ "$tries" -gt 0 ]; do
		contains "$(ab get value "$selector")" "$needle" && return 0
		tries=$((tries - 1))
		sleep 0.1
	done
	return 1
}

is_visible() {
	[ "$(ab is visible "$1")" = "true" ]
}

has_tokens() {
	selector="$1"
	count="$(ab eval "document.querySelectorAll('$selector').length")"
	[ "$count" -gt 0 ]
}

wide_shell_is_capped() {
	width="$(ab eval 'Math.round(document.querySelector("main").getBoundingClientRect().width)')"
	[ "$width" -le 1440 ]
}

highlight_is_escaped() {
	[ "$(ab eval 'document.querySelector("#editor-highlight-code b") === null')" = "true" ]
}

highlight_scroll_is_synced() {
	script='(() => {
  const editor = document.querySelector("#editor");
  const highlight = document.querySelector("#editor-highlight");
  editor.scrollLeft = 80;
  editor.dispatchEvent(new Event("scroll"));
  return editor.scrollLeft === highlight.scrollLeft;
})()'
	synced="$(ab eval "$script")"
	[ "$synced" = "true" ]
}

editor_tab_leaves() {
	ab eval 'document.querySelector("#editor").focus()' >/dev/null
	ab press Tab >/dev/null
	[ "$(ab eval 'document.activeElement !== document.querySelector("#editor")')" = "true" ]
}

mobile_panels_fit_viewport() {
	editor_min_height="$(ab eval 'parseFloat(getComputedStyle(document.querySelector(".editor-section")).minHeight)')"
	output_min_height="$(ab eval 'parseFloat(getComputedStyle(document.querySelector(".output-section")).minHeight)')"
	[ "$(awk "BEGIN { print ($editor_min_height < 340 && $output_min_height < 160) }")" -eq 1 ]
}

ab open "$BASE_URL/#01" >/dev/null
assert "loads lesson 01 title" contains "$(ab get text '#lesson-title')" "01_arrays_slices.go"
assert "shows lesson pane" is_visible ".lesson-pane"
assert "shows editor" is_visible "#editor"
assert "shows output" is_visible "#output"
assert "highlights Go keywords" has_tokens ".tok-keyword"
assert "highlights Go types" has_tokens ".tok-type"
assert "highlights Go strings" has_tokens ".tok-string"
assert "highlights Go comments" has_tokens ".tok-comment"
assert "hides imports by default" not_contains "$(ab get value '#editor')" 'import "fmt"'
assert "shows matching line numbers" contains "$(ab get text '#line-numbers')" "10"
assert "Tab leaves the editor" editor_tab_leaves

ab click "#next-btn" >/dev/null
assert "next button navigates to lesson 02" contains "$(ab get url)" "#02"
assert "lesson 02 finishes rendering" wait_contains "#lesson-title" "02_copy.go"
assert "grouped imports start hidden" not_contains "$(ab get value '#editor')" '"slices"'
ab click "#imports-btn" >/dev/null
assert "imports toggle reveals imports" contains "$(ab get value '#editor')" '"slices"'
ab click "#imports-btn" >/dev/null
assert "imports toggle hides imports again" not_contains "$(ab get value '#editor')" '"slices"'
ab click "#run-btn" >/dev/null
assert "hidden imports are restored for Run" wait_contains "#status-msg" "Run complete"
assert "lesson 02 output is captured" contains "$(ab get text '#output-code')" "Clone creates"

ab open "$BASE_URL/#16" >/dev/null
assert "loads printing lesson" wait_contains "#lesson-title" "16_printing.go"

ab fill "#editor" 'package main

// import "os"
var example = `import "strings"`

func main() {}
' >/dev/null
assert "import text in comments stays visible" contains "$(ab get value '#editor')" 'import "os"'
assert "import text in raw strings stays visible" contains "$(ab get value '#editor')" 'import "strings"'

ab fill "#editor" 'package main

import "os"

func main() { _ = os.Stdout }
' >/dev/null
assert "new top-level import is hidden" not_contains "$(ab get value '#editor')" 'import "os"'
ab click "#imports-btn" >/dev/null
assert "new top-level import replaces hidden imports" contains "$(ab get value '#editor')" 'import "os"'
assert "replaced hidden imports discard the old import" not_contains "$(ab get value '#editor')" 'import "fmt"'
ab click "#imports-btn" >/dev/null

ab eval 'document.querySelector("#editor").focus()' >/dev/null
ab press Control+Enter >/dev/null
assert "Ctrl+Enter still runs the lesson" wait_contains "#status-msg" "Run complete"

ab fill "#editor" 'package main

func main() {}
' >/dev/null
ab reload >/dev/null
assert "edited source persists across reload" wait_value_contains "#editor" "func main() {}"

ab click "#format-btn" >/dev/null
assert "format button reports success" wait_contains "#status-msg" "Formatted"

ab click "#reset-btn" >/dev/null
assert "reset restores original source" contains "$(ab get value '#editor')" "type point struct"

ab click "#run-btn" >/dev/null
assert "run completes" wait_contains "#status-msg" "Run complete"
assert "run captures stdout" contains "$(ab get text '#output-code')" "Printf:"
assert "run captures stderr" contains "$(ab get text '#output-code')" "Fprintln: stderr"

ab fill "#editor" 'package main

import "fmt"

func main() {
	// <b>Verify the Run button uses edited source.</b>
	fmt.Println("clicked run")
}
' >/dev/null
assert "highlight follows edits" contains "$(ab get text '#editor-highlight-code')" "clicked run"
assert "edited comments stay highlighted" has_tokens ".tok-comment"
assert "edited strings stay highlighted" has_tokens ".tok-string"
assert "highlighted source stays escaped" highlight_is_escaped
ab click "#run-btn" >/dev/null
assert "Run click executes edited source" wait_contains "#status-msg" "Run complete"
assert "Run click shows edited output" contains "$(ab get text '#output-code')" "clicked run"

ab fill "#editor" 'package main

func main( {
' >/dev/null
ab click "#run-btn" >/dev/null
assert "Run click reports compile failure" wait_contains "#status-msg" "Compile failed"
assert "compile failure reaches output" contains "$(ab get text '#output-code')" "syntax error"

ab open "$BASE_URL/#31" >/dev/null
ab click "#run-btn" >/dev/null
assert "HTTP lesson completes" wait_contains "#status-msg" "Run complete"
assert "HTTP lesson serves and calls its handler" \
	contains "$(ab get text '#output-code')" "response: 200 hello Ana"

ab set viewport 390 844 >/dev/null
ab open "$BASE_URL/#01" >/dev/null
assert "editor stays visible at mobile width" is_visible "#editor"
assert "run button stays visible at mobile width" is_visible "#run-btn"
assert "highlight scroll follows editor" highlight_scroll_is_synced

ab set viewport 740 390 >/dev/null
assert "mobile panels shrink on short landscape screens" mobile_panels_fit_viewport

ab set viewport 2048 1048 >/dev/null
assert "tour shell stays capped at wide width" wide_shell_is_capped

if [ "$FAIL" -ne 0 ]; then
  echo "playtest: one or more assertions failed" >&2
  exit 1
fi
echo "playtest: all assertions passed"
