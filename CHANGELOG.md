# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.2, 2026-10-28

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.
- Raised the minimum Go version in `go.mod` from 1.26.9 to 1.26.10, which fixes the `crypto/x509` vulnerability GO-2026-6809 found by govulncheck.
- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.1.1, 2026-10-15

- Added the `go-asyncapi` agent skill, which teaches a coding agent how to load specs and walk their schemas with this library. It installs from the psyb0t Claude Code and Codex marketplaces and from ClawHub with `openclaw skills install @psyb0t/go-asyncapi`.
- Doc comments no longer describe the library in terms of a code generator.

## v0.1.0, 2026-10-15

- Added `Load`, which reads an AsyncAPI 3 document in YAML or JSON and returns its servers, channels, operations, messages and component schemas in document order, with local and relative-file `$ref`s followed and message and operation traits merged.
- Added `Resolver` and `Location` for dereferencing `$ref`s inside schemas, and `Object`, an ordered mapping for reading schema trees.
