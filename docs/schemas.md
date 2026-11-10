---
title: "Walking schemas"
description: "Read go-asyncapi payload and header schemas as ordered Object trees and follow the $refs inside them with Resolver.Deref and Location.Child."
nav_order: 3
---

# Walking schemas

`Load` resolves the structure of a document, but it leaves payload and header schemas as they are in the file. You get them as [`*asyncapi.Object`](api.md#object) trees that keep key order, and you follow the `$ref`s inside them yourself with the document's [`Resolver`](api.md#resolver). This is what a code generator, a doc tool or a validator built on go-asyncapi does.

## Where schemas live

Each schema comes with the [`Location`](api.md#location) it was found at. Keep the two together.

| Schema | Value | Location |
| --- | --- | --- |
| Message payload | `message.Payload` | `message.PayloadAt` |
| Message headers | `message.Headers` | `message.HeadersAt` |
| A `components.schemas` entry | `component.Value` | `component.At` |

`message.Payload` has already been dereferenced once by `Load`, so it is the schema itself, not a `$ref`. A `components.schemas` value is raw and may still be a `$ref`.

## Values in the tree

A schema is one of these Go values:

| Go type | JSON or YAML |
| --- | --- |
| `*asyncapi.Object` | a mapping, with `Keys` in document order and `Values` by key |
| `[]any` | a sequence |
| `string`, `int64`, `float64`, `bool`, `nil` | a scalar |

A JSON Schema can be a bare `true` or `false`, so a payload can be a `bool` instead of an `*Object`.

`Object.Get(key)` returns a value and whether the key exists. `Object.String(key)` returns the string under a key or `""`. `Object.Object(key)` returns the mapping under a key or nil. All three are safe to call on a nil `*Object`, so chains like `object.Object("properties").Object("id")` do not panic on a missing key.

## Following `$ref`s

`document.Resolver.Deref(value, location)` follows a `$ref` chain and returns the target value and the target's own location. A value that is not a `$ref` comes back unchanged with the location you passed in.

```go
message := document.Messages[0]

schema, at, err := document.Resolver.Deref(message.Payload, message.PayloadAt)
if err != nil {
	return err
}

object, ok := schema.(*asyncapi.Object)
if !ok {
	return nil // a boolean schema such as `payload: true`
}

properties := object.Object("properties")
if properties == nil {
	return nil // not an object schema: oneOf, allOf, a primitive and so on
}

for _, name := range properties.Keys {
	property, propertyAt, err := document.Resolver.Deref(
		properties.Values[name],
		at.Child("properties").Child(name),
	)
	if err != nil {
		return err
	}

	fmt.Println(name, propertyAt)
	_ = property // *asyncapi.Object, or a bool schema
}
```

Build the location of a nested value from the location `Deref` returned, with `Child`, one key or array index at a time. A `$ref` resolves against the file its value came from, and the location is how the resolver knows which file that is. If you pass the parent's location from before the `Deref`, a relative `$ref` in another file resolves against the wrong directory.

`Child` escapes `~` and `/` in a key as JSON Pointer requires, so a property called `a/b` is safe to pass as is.

## Visiting each schema once

Two `$ref`s that land on the same definition return the same `Location`. `Location` is a comparable struct of two strings, so it works as a map key:

```go
seen := map[asyncapi.Location]bool{}

visit := func(value any, at asyncapi.Location) error {
	_, target, err := document.Resolver.Deref(value, at)
	if err != nil {
		return err
	}

	if seen[target] {
		return nil
	}

	seen[target] = true
	// handle the schema here

	return nil
}
```

`Location.String()` renders the location as `file#/pointer`, which also makes a usable key or a log field. `Location.LastToken()` returns the last pointer token unescaped, which is the schema name for a `components.schemas` entry.

## Reading without following

`document.Resolver.Lookup(location)` returns the raw value at a location without following any `$ref`. `document.Resolver.Lookup(document.Root)` gives you the whole root mapping, which is how you reach anything `Load` does not copy into the structs, such as `info.contact` or a channel's bindings.

```go
root, err := document.Resolver.Lookup(document.Root)
if err != nil {
	return err
}

contact := root.(*asyncapi.Object).Object("info").Object("contact")
fmt.Println(contact.String("email"))
```

## Using the resolver alone

`asyncapi.NewResolver()` returns an empty resolver, and `Resolver.Document(path)` loads and caches a file and returns its root value. You can use the two to read any YAML or JSON file as an ordered tree with `$ref` support, without `Load` and without any AsyncAPI checks.
