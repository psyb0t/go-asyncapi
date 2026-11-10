---
title: "Errors"
description: "The commerr sentinels go-asyncapi wraps, which failure each one means, and how to match them with errors.Is when a spec fails to load."
nav_order: 5
---

# Errors

Every error `Load` and the `Resolver` return is a [ctxerrors](https://github.com/psyb0t/ctxerrors) error. It carries the file, line and function in the library where it was created, and it wraps a sentinel from [`ctxerrors/commerr`](https://psyb0t.github.io/ctxerrors/commerr/), so `errors.Is` tells you what kind of failure it was. The message names the spec location, such as `/abs/path/asyncapi.yaml#/channels/events`.

```go
_, err := asyncapi.Load("asyncapi.yaml")

switch {
case err == nil:
case errors.Is(err, commerr.ErrNotFound):
	// a $ref or pointer that leads nowhere
case errors.Is(err, commerr.ErrInvalidValue):
	// not AsyncAPI 3, or a section with the wrong shape
case errors.Is(err, fs.ErrNotExist):
	// the spec or a file it references is missing
default:
	return err
}
```

## Sentinels

| Sentinel | When |
| --- | --- |
| `commerr.ErrNotFound` | A `$ref` or JSON pointer points at nothing: a missing key, an array index out of range, or a token inside a scalar. Also an operation or reply `channel` that does not point at an entry under `channels`. |
| `commerr.ErrInvalidValue` | The `asyncapi` field does not start with `3.`. A server, server variable, channel, parameter, message, operation, reply, reply address, trait or `correlationId` that should be a mapping is not. An operation `action` other than `send` or `receive`. A `$ref` that is not a string. A `$ref` chain longer than 64 hops, which is how a cycle shows up. |
| `commerr.ErrInvalidArgument` | A `$ref` to a remote URL, anything containing `://`. |
| `commerr.ErrRequiredFieldNotSet` | An operation without a `channel`. |
| `commerr.ErrParseFailed` | The file is not valid YAML or JSON, or its top level is not a mapping. An empty file counts as not a mapping. |
| `fs.ErrNotExist` | The spec file, or a file a `$ref` names, does not exist. This is the `os` error wrapped, not a `commerr` sentinel. |

A scalar with an explicit YAML tag that does not parse, such as `!!int nope`, fails with the `strconv` error wrapped and no sentinel.

## Wrapping in your code

The errors already carry location and context, so wrap them once more with what your code was doing and pass them up. The sentinels survive any number of `ctxerrors.Wrap` or `fmt.Errorf("...: %w", err)` layers.
