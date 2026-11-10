# Go code

- **GO1** Every new failure wraps a `ctxerrors/commerr` sentinel: a dangling `$ref` or pointer is `commerr.ErrNotFound`, a malformed section is `commerr.ErrInvalidValue`, unparseable input is `commerr.ErrParseFailed`.
- **GO4** Any value that can hold a `$ref` is loaded through `Resolver.Deref` with the `Location` it was found at, and its children get locations built with `Location.Child`. Never resolve a `$ref` by hand.
- **GO5** Mappings stay `*Object` with `Keys` in document order. Never range over `Values` when order matters.
- **GO8** Doc comments and other prose describe the library's behavior, not any one consumer's: no "the generator" or "generated type".
