---
title: "Agent skill"
description: "Install the go-asyncapi agent skill in Claude Code, Codex or OpenClaw so a coding agent knows how to load AsyncAPI 3 specs with the library."
nav_order: 6
---

# Agent skill

The repository ships an agent skill in [`.agents/skills/go-asyncapi`](https://github.com/psyb0t/go-asyncapi/tree/main/.agents/skills/go-asyncapi). It teaches a coding agent how to use go-asyncapi when it writes Go: when to reach for the library and when not to, how to load a spec, how to walk payload schemas with the resolver, the type reference, and the error sentinels. It works in any agent that reads `.agents/skills/`, and it installs natively in the clients below.

The skill only guides edits to your own Go code. The library it describes reads the spec file you pass to `Load` and the files that spec references, writes nothing, and makes no network calls.

## Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-asyncapi@psyb0t
```

The plugin manifest is [`.agents/.claude-plugin/plugin.json`](https://github.com/psyb0t/go-asyncapi/blob/main/.agents/.claude-plugin/plugin.json).

## Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-asyncapi@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repository. The plugin manifest is [`.agents/.codex-plugin/plugin.json`](https://github.com/psyb0t/go-asyncapi/blob/main/.agents/.codex-plugin/plugin.json).

## OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-asyncapi
```

It declares the `go` binary as a requirement.
