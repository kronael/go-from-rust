# Changelog

## [v0.1.0] — 20260726

> Minimum Viable Go v0.1.0 — Go, mapped from Rust
>
> 42 small runnable Go lessons that map what you know in Rust to idiomatic Go, with a browser tour to run each one.
>
> • 42 lessons — slices, language mappings, collections, concurrency, control, services, generics
> • Interactive tour at krons.cx/go-from-rust — edit, format, and run each lesson in the browser
> • Reads on mobile — code fits with no horizontal scroll; concepts live in styled descriptions
> • Every lesson is one self-contained runnable file with a direct Rust→Go mapping
>
> Full notes: CHANGELOG.md

- 42 lessons (`01`–`42`) across seven groups: foundations and slices, language
  mappings, collections and ordering, concurrency, control and failure, services
  and data, generics and tooling.
- Interactive web tour (`-tags web`): per-lesson editor with run/format against a
  Playground stand-in; deployed at `krons.cx/go-from-rust`, discoverable from the
  krons portal, the guides list, and the go-lore essay.
- Mobile-first readability: every lesson wraps to ≤60 chars and the editor uses a
  small font, so code fits with no horizontal scroll; lesson descriptions carry
  the concept and Rust contrast, rendered with inline code chips.
- `defer` taught at lesson 21 (before its first use), `option` focused on
  representing optionality idiomatically, and struct embedding flagged as a niche
  tool rather than a default.
