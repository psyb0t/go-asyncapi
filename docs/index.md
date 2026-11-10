---
title: "Home"
description: "go-asyncapi loads an AsyncAPI 3 spec in Go and hands back servers, channels, operations and messages with every $ref followed and traits merged."
nav_order: 1
permalink: /
seo:
  type: WebSite
---

# go-asyncapi

Loads AsyncAPI 3 documents in Go. Point `asyncapi.Load` at a YAML or JSON file and you get servers, channels, operations and messages as plain structs, with every `$ref` followed, across files too, and every message and operation trait merged in. [gen-asyncapi](https://github.com/psyb0t/gen-asyncapi) is built on it.

It reads and resolves, nothing more. It does not validate messages against their schemas and it does not generate code. Payload and header schemas come back as ordered trees for you to walk.

## Install

```bash
go get github.com/psyb0t/go-asyncapi
```

The package name is `asyncapi` and the module path ends in `go-asyncapi`, so import it with an explicit name:

```go
import asyncapi "github.com/psyb0t/go-asyncapi"
```

The module declares `go 1.26.10` in its [`go.mod`](https://github.com/psyb0t/go-asyncapi/blob/main/go.mod). Its only runtime dependencies are [ctxerrors](https://github.com/psyb0t/ctxerrors) and `gopkg.in/yaml.v3`.

## Quick look

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

## Where to go next

| You want to | Read |
| --- | --- |
| Know what `Load` reads, what it skips, and how `$ref`s and traits behave | [Loading documents](loading.md) |
| Walk a payload or header schema and follow the `$ref`s inside it | [Walking schemas](schemas.md) |
| Look up a type, field or method | [API reference](api.md) |
| Match a failure with `errors.Is` | [Errors](errors.md) |
| Teach a coding agent to use the library | [Agent skill](agent-skill.md) |
| Build, lint and test the library itself | [Development](development.md) |

The Go package docs are on [pkg.go.dev](https://pkg.go.dev/github.com/psyb0t/go-asyncapi), and the source is on [GitHub](https://github.com/psyb0t/go-asyncapi).
