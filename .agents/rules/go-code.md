# Go code

Refines `~/.claude/rule-details/go/wrapping-errors.md` and `~/.claude/rule-details/coding-style.md`.

- **GO1** Every new failure wraps a `ctxerrors/commerr` sentinel with `ctxerrors.Wrap` or `ctxerrors.Wrapf`, so callers can match it with `errors.Is`: a dangling `$ref` or pointer is `commerr.ErrNotFound`, a malformed section is `commerr.ErrInvalidValue`, unparseable input is `commerr.ErrParseFailed`. Errors from a call are wrapped with context on the next line. No `fmt.Errorf`, no `errors.New` outside a sentinel declaration.
- **GO2** The package never logs and never prints. It returns errors.
- **GO3** Everything exported is public API that gen-asyncapi and other importers use. Add exported names only with a doc comment that says what the value holds or what the function does, including when it returns nil or an error. Rename or remove an exported name only in a breaking release (see `releases.md`).
- **GO4** Any value that can hold a `$ref` is loaded through `Resolver.Deref` with the `Location` it was found at, and its children get locations built with `Location.Child`. Never resolve a `$ref` by hand.
- **GO5** Mappings stay `*Object` with `Keys` in document order. Never range over `Values` when order matters.
- **GO6** Finite string values get a named type with constants, as `Action` does.
- **GO7** `make lint` passes with zero issues. Lines are at most 80 columns, the `.golangci.yml` limit; split calls one argument per line rather than adding `//nolint`.
- **GO8** Doc comments and other prose describe the library's behavior, not any one consumer's: no "the generator" or "generated type".
