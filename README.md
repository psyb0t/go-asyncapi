# go-asyncapi

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-asyncapi.svg)](https://pkg.go.dev/github.com/psyb0t/go-asyncapi)
[![CI](https://github.com/psyb0t/go-asyncapi/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-asyncapi/actions/workflows/pipeline.yml)
[![coverage](https://raw.githubusercontent.com/psyb0t/go-asyncapi/badges/coverage.svg)](https://github.com/psyb0t/go-asyncapi/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-asyncapi/badges/version.svg)](https://github.com/psyb0t/go-asyncapi/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-asyncapi/badges/license.svg)](https://github.com/psyb0t/go-asyncapi/blob/main/LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-asyncapi/badges/importers.svg)](https://github.com/psyb0t/go-asyncapi/blob/badges/importers.md)

Loads AsyncAPI 3 documents in Go. Point it at a YAML or JSON file and you get servers, channels, operations and messages as plain structs, with every `$ref` followed, across files too, and every trait merged in. [gen-asyncapi](https://github.com/psyb0t/gen-asyncapi) is built on it.

## Install

```sh
go get github.com/psyb0t/go-asyncapi
```

## Usage

```go
document, err := asyncapi.Load("asyncapi.yaml")
if err != nil {
	return err
}

for _, server := range document.Servers {
	fmt.Println(server.Name, server.URL()) // production wss://eu.example.com:443/v2/ws
}

for _, operation := range document.Operations {
	for _, message := range operation.Messages {
		fmt.Println(operation.Action, operation.Channel.Address, message.Name)
	}
}
```

Everything comes back in the order the document lists it. A message used by several channels is one `*Message`, so you can compare pointers. Traits are already merged: a trait fills what the message or operation leaves out and never overrides what it sets.

Payload and header schemas stay as they are in the document, as `*asyncapi.Object` trees that keep key order, and `document.Resolver.Deref` follows the `$ref`s inside them. Errors wrap [`ctxerrors/commerr`](https://github.com/psyb0t/ctxerrors) sentinels, so `errors.Is(err, commerr.ErrNotFound)` catches a dangling `$ref`. A `$ref` to a remote URL is an error. The library never touches the network.

## Documentation

The full docs live at **[psyb0t.github.io/go-asyncapi](https://psyb0t.github.io/go-asyncapi/)**.

| Page | What's in it |
| --- | --- |
| [Loading documents](https://psyb0t.github.io/go-asyncapi/loading/) | What `Load` reads and skips, `$ref` rules, traits, schema formats, ordering |
| [Walking schemas](https://psyb0t.github.io/go-asyncapi/schemas/) | `Object` trees, `Resolver.Deref`, `Location.Child`, visiting each schema once |
| [API reference](https://psyb0t.github.io/go-asyncapi/api/) | Every exported type, field and method |
| [Errors](https://psyb0t.github.io/go-asyncapi/errors/) | Each `commerr` sentinel and the failure it stands for |
| [Agent skill](https://psyb0t.github.io/go-asyncapi/agent-skill/) | Installing the skill in Claude Code, Codex and OpenClaw |
| [Development](https://psyb0t.github.io/go-asyncapi/development/) | Make targets, the dev image, tests, CI and releases |

## Agent integrations

The [skill](https://github.com/psyb0t/go-asyncapi/tree/main/.agents/skills/go-asyncapi) teaches a coding agent how to use go-asyncapi when it writes Go. In Claude Code:

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-asyncapi@psyb0t
```

Codex and OpenClaw install it too. See [Agent skill](https://psyb0t.github.io/go-asyncapi/agent-skill/).

## Development

Everything runs in the `Dockerfile.dev` image. You need Docker and Make.

```sh
make help           # list targets
make lint           # shellcheck, shfmt, go fix, golangci-lint
make test-coverage  # go test -race with a 90% coverage gate
make sec            # gitleaks, govulncheck, semgrep
make docs           # build the docs site locally
```

More on [Development](https://psyb0t.github.io/go-asyncapi/development/).

## Sponsor

<a href="https://www.buymeacoffee.com/psyb0t"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-blue.png" alt="Buy Me a Coffee" height="60" width="217"></a>

## License

[MIT](https://github.com/psyb0t/go-asyncapi/blob/main/LICENSE).
