---
title: "Loading documents"
description: "What asyncapi.Load reads from an AsyncAPI 3 spec, what it leaves out, and how it follows $refs across files, merges traits and orders results."
nav_order: 2
---

# Loading documents

`asyncapi.Load(path)` reads one AsyncAPI 3 document and every local file it references, and returns a [`*asyncapi.Document`](api.md#document). This page covers what ends up in that document and the rules `Load` follows on the way. The code is in [`document.go`](https://github.com/psyb0t/go-asyncapi/blob/main/document.go).

```go
document, err := asyncapi.Load("asyncapi.yaml")
if err != nil {
	return err
}

fmt.Println(document.Title, document.Version, len(document.Channels))
```

## Input

- YAML or JSON. JSON is valid YAML, so both go through the same `yaml.v3` decoder. YAML anchors and aliases are expanded.
- The top level must be a mapping, or `Load` fails with `commerr.ErrParseFailed`. An empty file fails the same way.
- The `asyncapi` field must start with `3.`. A 2.x document, or one with no `asyncapi` field, fails with `commerr.ErrInvalidValue`.
- A mapping that repeats a key keeps the last value at the position of the first one.

Scalars decode to `string`, `int64`, `float64`, `bool` or `nil`. A scalar of any other YAML type, such as a timestamp, stays a `string`.

## What it reads

| Section | Ends up in |
| --- | --- |
| `info.title`, `info.version`, `info.description` | `Document.Title`, `Document.Version`, `Document.Description` |
| `defaultContentType` | `Document.DefaultContentType`, and every message that sets no `contentType` |
| `servers` | `Document.Servers`: name, host, pathname, protocol, protocol version, description and variables with their enum and default. `URL()` fills the variables with their defaults; `URLTemplate()` leaves the `{placeholders}` in. |
| `channels` | `Document.Channels`: name, address, description, parameters (enum, default, description) and messages. A null address, which AsyncAPI uses for a dynamic channel, stays an empty string. |
| `operations` | `Document.Operations`: name, action, summary, description, channel, messages and reply. |
| operation `reply` | `Operation.Reply`: channel, messages and the runtime expression of a dynamic reply address, such as `$message.header#/replyTo`. |
| messages | `Document.Messages`: key, name, title, summary, description, content type, the deprecated flag, payload and headers schemas, the `correlationId` location, and bindings. |
| `components.messages` | `Document.Messages`, including messages no channel or operation uses. |
| `components.schemas` | `Document.ComponentSchemas`, in document order. |
| other `components` entries | Only through the `$ref`s that point into them. |

The full field list is on the [API reference](api.md).

## What it does not read

`Load` models the parts of a spec you need to send and receive messages. It does not copy these into the structs:

- `info` beyond title, version and description: contact, license, terms of service, tags, external docs;
- the document `id`;
- security schemes and `security` lists on servers and operations;
- `tags` and `externalDocs` on any object;
- server, channel and operation bindings (message bindings are read);
- server `title` and `summary`, channel `title`, `summary` and `servers`;
- message `examples`, `correlationId.description`, and operation `title` and `bindings`.

Everything is still in the parsed tree. `document.Resolver.Lookup(document.Root)` returns the root mapping as an [`*asyncapi.Object`](api.md#object), and [Walking schemas](schemas.md) shows how to move through it.

## Order

Every list comes back in the order the document writes it. `Document.Messages` holds each distinct message once, in the order `Load` first meets it: channel messages first, then operation messages that no channel listed, then component messages nobody referenced.

## One value, one pointer

A message or channel is identified by its [`Location`](api.md#location), the file and JSON pointer where its definition lives. Two `$ref`s that land on the same definition give the same `*Message`, so you can compare pointers. An operation's `channel` must point at one of the entries under `channels`, either directly or through the same `$ref` that entry uses. A channel defined only under `components.channels` and never listed in `channels` is not found.

## Message keys

`Message.Key` is the last token of the message's location. For a message under `components.messages` that is its key there. A message written inline inside a channel gets the channel name as a prefix, `<channel>_<key>`, so two channels that both define an inline `ping` give `alpha_ping` and `beta_ping`.

## Operation messages and replies

- An operation with a `messages` list gets exactly those messages. An operation without one gets every message of its channel.
- A reply with a `channel` uses that channel. A reply without one replies on the operation's channel.
- A reply's `messages` follow the same rule as the operation's: the listed messages, or every message of the reply channel.
- The reply itself can be a `$ref`, for example into `components.replies`.

## `$ref` resolution

- Local refs (`#/components/messages/ping`) and refs to other files (`common.yaml#/components/schemas/id`) both work. A relative file path resolves against the directory of the file that holds the `$ref`, not the working directory. Absolute paths work too.
- Each file is read once and cached by its absolute path.
- Chains are followed: a `$ref` to a `$ref` keeps going. After 64 hops `Load` gives up with `commerr.ErrInvalidValue`, which is how a cycle fails.
- A `$ref` whose target is missing fails with `commerr.ErrNotFound`. A `$ref` that is not a string fails with `commerr.ErrInvalidValue`.
- A `$ref` to a URL, anything with `://` in it, fails with `commerr.ErrInvalidArgument`. The library never opens a network connection. Download the document and reference it by path.

`Load` follows the refs that make up the document's structure: servers, server variables, channels, parameters, messages, operations, replies, reply addresses, traits, `correlationId`, and the top of each payload and headers schema. `$ref`s deeper inside a schema stay as they are, and you follow them with the resolver, as [Walking schemas](schemas.md) shows.

## Traits

Message and operation `traits` are merged in list order, following the AsyncAPI traits merge mechanism: JSON Merge Patch where a property the message or operation already sets is never overridden. Nested mappings merge key by key with the same rule, so a trait can add a binding without replacing the ones the message has.

A top-level property that comes from a trait keeps the location of the trait, so a `$ref` in a trait's payload resolves against the trait's file. `Message.PayloadAt` and `Message.HeadersAt` point there.

## Schema formats

A payload can be a plain schema or a Multi Format Schema Object with `schemaFormat` and `schema`. Without a `schemaFormat` key the payload itself is the schema. When `schemaFormat` is empty or names an AsyncAPI schema (`application/vnd.aai.asyncapi...`), JSON Schema (`application/schema+json...` or `application/schema+yaml...`) or OpenAPI (`application/vnd.oai.openapi...`), the inner `schema` becomes `Message.Payload`. The check ignores case.

Any other format, such as Avro or Protobuf, leaves `Payload` nil and puts the format string in `Message.UnsupportedFormat`. Headers follow the same rule, but an unsupported headers format just leaves `Headers` nil.

## Errors

Every failure wraps a `commerr` sentinel you can match with `errors.Is`. The full list is on [Errors](errors.md).
