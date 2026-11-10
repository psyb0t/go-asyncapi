---
title: "API reference"
description: "Every exported type, field, function and method in go-asyncapi: Load, Document, Server, Channel, Operation, Message, Object, Location and Resolver."
nav_order: 4
---

# API reference

Everything `github.com/psyb0t/go-asyncapi` exports, by type. The same reference, generated from the doc comments, is on [pkg.go.dev](https://pkg.go.dev/github.com/psyb0t/go-asyncapi). What each field is filled from is on [Loading documents](loading.md).

## Load

```go
func Load(path string) (*Document, error)
```

Reads the AsyncAPI 3 document at `path` and every local file it references, and returns the loaded document. On any failure it returns nil and an error that wraps a [`commerr` sentinel](errors.md).

## Document

A loaded AsyncAPI 3 document.

| Field | Type | Holds |
| --- | --- | --- |
| `AsyncAPIVersion` | `string` | The `asyncapi` field, such as `3.0.0`. |
| `Title` | `string` | `info.title`. |
| `Version` | `string` | `info.version`. |
| `Description` | `string` | `info.description`. |
| `DefaultContentType` | `string` | `defaultContentType`. Messages without a `contentType` already carry it. |
| `Servers` | `[]*Server` | The `servers` entries, in document order. |
| `Channels` | `[]*Channel` | The `channels` entries, in document order. |
| `Operations` | `[]*Operation` | The `operations` entries, in document order. |
| `Messages` | `[]*Message` | Every distinct message, in first-seen order, component messages included. |
| `ComponentSchemas` | `[]SchemaComponent` | The `components.schemas` entries, in document order. |
| `Resolver` | `*Resolver` | Holds every loaded file and dereferences `$ref`s inside schemas. |
| `Root` | `Location` | The location of the document's top level: the absolute file path and an empty pointer. |

## Server

One entry of the `servers` object.

| Field or method | Holds |
| --- | --- |
| `Name` | The key under `servers`. |
| `Host` | `host`, with `{variable}` placeholders as written. |
| `Pathname` | `pathname`, with placeholders as written. |
| `Protocol` | `protocol`, such as `wss` or `kafka`. |
| `ProtocolVersion` | `protocolVersion`. |
| `Description` | `description`. |
| `Variables` | `[]*ServerVariable`, in document order. |
| `URLTemplate() string` | `protocol://host` followed by the pathname, with the placeholders left in. |
| `URL() string` | The same URL with every variable replaced by its default. |

### ServerVariable

| Field | Holds |
| --- | --- |
| `Name` | The variable name, the text inside `{}`. |
| `Default` | `default`. |
| `Enum` | `enum` as `[]string`. Non-string items are dropped. |
| `Description` | `description`. |

## Channel

One entry of the `channels` object.

| Field | Holds |
| --- | --- |
| `Name` | The key under `channels`. |
| `Address` | `address`, or `""` for a null address. |
| `Description` | `description`. |
| `Parameters` | `[]*Parameter`, in document order. |
| `Messages` | `[]*Message`, in document order. |
| `Location` | Where the channel definition lives, after `$ref`s. |

### Parameter

| Field | Holds |
| --- | --- |
| `Name` | The key under `parameters`. |
| `Description` | `description`. |
| `Enum` | `enum` as `[]string`. Non-string items are dropped. |
| `Default` | `default`. |

## Operation

One entry of the `operations` object.

| Field | Holds |
| --- | --- |
| `Name` | The key under `operations`. |
| `Action` | `ActionSend` or `ActionReceive`. |
| `Summary` | `summary`, after traits. |
| `Description` | `description`, after traits. |
| `Channel` | The `*Channel` the operation's `channel` points at. |
| `Messages` | The listed messages, or every message of the channel when the operation lists none. |
| `Reply` | `*Reply`, or nil without a `reply`. |

### Action

```go
type Action string

const (
	ActionSend    Action = "send"
	ActionReceive Action = "receive"
)
```

Any other action fails `Load` with `commerr.ErrInvalidValue`.

### Reply

| Field | Holds |
| --- | --- |
| `Channel` | The reply's channel, or the operation's channel when the reply names none. |
| `Messages` | The listed messages, or every message of `Channel`. |
| `AddressLocation` | The runtime expression of a dynamic reply address, such as `$message.header#/replyTo`, or `""`. |

## Message

A message definition. A message referenced from several places is one `*Message`.

| Field | Type | Holds |
| --- | --- | --- |
| `Key` | `string` | The `components.messages` key, or `<channel>_<key>` for a message defined inline in a channel. Unique within a document. |
| `Name` | `string` | `name`. |
| `Title` | `string` | `title`. |
| `Summary` | `string` | `summary`. |
| `Description` | `string` | `description`. |
| `ContentType` | `string` | `contentType`, or the document's `defaultContentType`. |
| `Deprecated` | `bool` | `deprecated`. |
| `Payload` | `any` | The payload schema, usually an `*Object`. Nil without a payload or when the schema format is not supported. |
| `PayloadAt` | `Location` | Where `Payload` lives, for resolving `$ref`s inside it. |
| `UnsupportedFormat` | `string` | The payload's `schemaFormat` when it is not AsyncAPI, JSON Schema or OpenAPI, or `""`. |
| `Headers` | `any` | The headers schema, or nil. |
| `HeadersAt` | `Location` | Where `Headers` lives. |
| `CorrelationIDLocation` | `string` | The runtime expression of the `correlationId`, such as `$message.payload#/id`, or `""`. |
| `Bindings` | `*Object` | The message bindings after traits, or nil. |
| `Location` | `Location` | Where the message definition lives, after `$ref`s. |

Every string field and `Bindings` include what traits merged in.

## SchemaComponent

One `components.schemas` entry.

| Field | Type | Holds |
| --- | --- | --- |
| `Name` | `string` | The key under `components.schemas`. |
| `Value` | `any` | The schema as written, possibly a `$ref`. |
| `At` | `Location` | Where the entry lives. |

## Object

```go
type Object struct {
	Keys   []string
	Values map[string]any
}
```

A JSON or YAML mapping that keeps its keys in document order. Values are `string`, `int64`, `float64`, `bool`, `nil`, `[]any` or `*Object`. Range over `Keys` when order matters.

| Method | Returns |
| --- | --- |
| `Get(key string) (any, bool)` | The value under `key` and whether it exists. |
| `String(key string) string` | The string under `key`, or `""` when it is missing or not a string. |
| `Object(key string) *Object` | The mapping under `key`, or nil. |

All three accept a nil receiver.

## Location

```go
type Location struct {
	File    string
	Pointer string
}
```

A value inside a loaded document: an absolute file path and a JSON pointer such as `/components/schemas/user`. Two `$ref`s that resolve to the same location point at the same value, so it works as a map key.

| Method | Returns |
| --- | --- |
| `Child(key string) Location` | The location of `key` inside this one, with `~` and `/` escaped. |
| `LastToken() string` | The last pointer token, unescaped, or `""` at the root. |
| `String() string` | `file#/pointer`. |

## Resolver

Loads documents and follows `$ref`s between and inside them. `Load` creates one and keeps it in `Document.Resolver`.

| Function or method | Does |
| --- | --- |
| `NewResolver() *Resolver` | Returns a resolver with no documents loaded. |
| `Document(path string) (any, error)` | Returns the root value of the YAML or JSON file at `path`, reading and caching it by absolute path on first use. |
| `Lookup(location Location) (any, error)` | Returns the raw value at `location` without following `$ref`s. A pointer that leads nowhere is `commerr.ErrNotFound`. |
| `Deref(value any, location Location) (any, Location, error)` | Follows `$ref` chains starting at `value`, which lives at `location`, and returns the first value that is not a `$ref` and where it lives. Gives up after 64 hops. |

## Located

```go
type Located struct {
	Value any
	At    Location
}
```

A value merged in from a trait together with the location of the trait. `Load` uses it while merging and unwraps it before filling the structs, so it does not appear in a loaded `Document`.
