# gofrs — Go for Rustaceans

Teaching collection of minimal, runnable Go demos for people who already know Rust.

## Project intent

Every file is a single focused demo. No abstractions, no framework, no build system.
The point is that a Rustacean can read one file and immediately understand one Go concept.

## Style (modeled on gobyexample.com)

- **One concept per file** — self-contained, under ~80 lines
- **Code-first** — minimal prose, no preamble; the code is the explanation
- **Inline annotation only** — short comments between code blocks, not at the top
- **Show output** — if it runs, the output matters; make it obvious
- **Pragmatic, dry tone** — no flair, no padding, no marketing
- **Rust→Go framing** — name the Rust thing, show the Go equivalent; don't explain both from scratch
- **No helpers, no abstractions** — if you need a 3-line helper to make the demo cleaner, the demo is wrong
- **Deliberate incompleteness is fine** — a file that won't compile (caseA.go, badidx.go) is a valid demo

## What NOT to do

- Don't add error handling that obscures the concept being demonstrated
- Don't generalize a demo into a reusable pattern
- Don't add commentary explaining what the code does (name it well instead)
- Don't make demos longer to be "more complete"
- Don't add a new file unless it demonstrates a genuinely distinct concept
