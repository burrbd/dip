# CLAUDE.md

## Project Vision

This repo is a **Diplomacy messenger bot** — players submit moves via slash commands
(e.g. `/order A Vie-Bud`), view the board map, and see move history. Each channel hosts
exactly one game. The game state is modelled as an **event log**: structured JSON snapshots
posted to the channel after each phase resolution serve as the audit trail and persistence
layer (no external database required).

**Adjudication is handled by [godip](https://github.com/zond/godip)**, which provides complete
DATC compliance across movement, retreat, and adjustment phases. This repo's job is the bot layer
on top of godip: command parsing, session management, event logging, map rendering, and platform
integration (Slack, Telegram).

---

## Build Plan

See [`PLAN.md`](PLAN.md) for the ordered list of stories and current progress.

---

## Architecture

See [`ARCHITECTURE.md`](ARCHITECTURE.md) for the full component breakdown, data flow, command
set, and design decisions.

Package layout at a glance:

```
engine/     — godip wrapper (adapter, parser, phase advance, win detection)
events/     — structured JSON event log (write, scan, replay)
session/    — game lifecycle (Session struct, serialization, turn advance, deadlines)
bot/        — platform-agnostic command router, formatter, autocomplete
dipmap/     — SVG → PNG rendering, province highlighting, BFS neighborhood
platform/   — Slack and Telegram adapters
cmd/        — entry points (slackbot, telegrambot)

game/       — [LEGACY] partial custom adjudicator (preserved, not active development)
```

---

## Development Philosophy

**TDD with 100% test coverage is non-negotiable.** The workflow is red → green → refactor:

1. **Red** — write a failing test for the story's acceptance criteria first
2. **Green** — write the minimal code to make it pass
3. **Refactor** — simplify and remove duplication

- `go test -v -cover -race ./...` must pass at all times
- Never commit with a failing test

---

## Running Tests

```bash
go test -v -cover -race ./...
```

CI runs this same command via CircleCI (`.circleci/config.yml`).

---

## Testing Conventions

- Use `github.com/cheekybits/is` for assertions: `is.NoErr(err)`, `is.Equal(a, b)`, `is.NotNil(v)`
- Unit tests live alongside the code they test (`_test.go` files in same package or `_test` package)
- Mock interfaces inline in test files
- Table-driven tests for parsers and validators

---

## Package Conventions

- **Interfaces over concrete types** — keep packages decoupled and testable
- **Value semantics for orders** — order structs passed by value, never mutated
- **No external database** — channel message history is the persistence layer; state rebuilt via event replay

---

## Module

```
module github.com/burrbd/dip
go 1.12
```

Key dependencies:
- `github.com/zond/godip` — Diplomacy adjudication engine (DATC-compliant, all phases)
- `gonum.org/v1/gonum` — graph library (used in legacy `game/` package)
- `github.com/cheekybits/is` — test assertions

---

## Legacy: Custom Adjudicator

> **Preserved but not the active development path.** Adjudication has been outsourced to godip.
> The `game/` package is kept for reference.

The `game/` package implements a partial custom Diplomacy adjudicator based on the DATC
(Diplomacy Adjudicator Test Cases). A local copy of the DATC is at `DATC.txt`.

Phase 1 (army resolution — DATC 6.A.11, 6.A.12, 6.C.1–6.C.3, 6.D.1–6.D.5 and others) is
done or in progress. Phases 2–5 (self-dislodgement, head-to-head, fleet model, convoys) are
not planned.
