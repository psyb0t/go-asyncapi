# Testing

Refines `~/.claude/rule-details/go/writing-tests.md`.

- **TST1** Tests are black-box, in package `asyncapi_test`, and use only the exported API.
- **TST2** Table tests use `testCases` and `tc`, call `t.Parallel()` in the test and in every subtest, and assert errors with `require.ErrorIs` against the `commerr` sentinel.
- **TST3** Spec fixtures go in `testdata/`. A case that only needs a few lines of YAML writes a temp file with `writeDocument` in `malformed_test.go`, which prepends the `asyncapi`/`info` header.
- **TST4** Every new loader branch gets a test: the happy path through `testdata/kitchensink.yaml` or a small fixture, and each error branch as a row in `TestLoadRejectsMalformedSections`.
- **TST5** `make test-coverage` keeps statement coverage at 90% or above, with `-race`.
