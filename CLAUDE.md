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
go 1.21
```

Key dependencies:
- `github.com/zond/godip` — Diplomacy adjudication engine (DATC-compliant, all phases)
- `gonum.org/v1/gonum` — graph library (used in legacy `game/` package)
- `github.com/cheekybits/is` — test assertions

---

## Environment Constraints

**No network access.** `go get` and the module proxy do not work. All dependencies must be present in `vendor/`. If a new dependency is needed:
1. Add it to `go.mod` with any pseudo-version (e.g. `v0.0.0-20200101000000-000000000000`)
2. Add it to `vendor/modules.txt` with `## explicit` annotation
3. Create the package files under `vendor/<module-path>/`

**godip is currently a stub.** The real godip library is not yet downloaded. `vendor/github.com/zond/godip/` contains a minimal stub that provides the correct interface signatures but zero adjudication logic. It is sufficient to compile and run all engine tests. When real godip becomes available:
- Replace `vendor/github.com/zond/godip/` with the actual source
- Update `vendor/modules.txt` with the correct version
- Update `go.mod` and `go.sum` with the real checksum

**godip.Adjudicator interface** (defined in the stub at `vendor/github.com/zond/godip/godip.go`) is the key seam. Any package interacting with game state should depend on this interface, not on `*classicalState` or other concrete types.

---

## Testing Conventions (additions)

- Tests that need access to unexported helpers (e.g. `fillNMR`, `isEmptyPhase`, `newFromVariant`) use `package engine` (not `package engine_test`) so they share the package namespace. This is the standard Go pattern for white-box testing.
- 100% coverage means **all branches**, not just all lines — use `go tool cover -func` to check per-function coverage when the summary isn't 100%.

### Timer tests and the race detector

`time.AfterFunc` callbacks run in a new goroutine. The race detector will flag any
unsynchronized access to shared state between the test goroutine and the callback
goroutine. The established pattern:

- Add a `sync.Mutex` (`mu`) to any struct that owns a `*time.Timer`.
- Lock `mu` whenever reading or writing the timer field — including in tests that
  set `s.timer` directly.
- Use `sync.WaitGroup` (or a channel) to wait for the callback goroutine to finish
  completely before making post-condition assertions.

### Timer callback coverage

Anonymous function literals passed to `time.AfterFunc` are counted as separate
statements by the coverage tool. If the timer never fires during a test, the closure
body is marked uncovered even though the outer function is 100% covered. The fix:
extract the callback to a named method (e.g. `onDeadline()`) and call that from
`time.AfterFunc`. The named method can then be called directly in a unit test
without needing a real timer.

### Recovering engine.Engine from events.Rebuild

`events.Rebuild` returns `events.EngineState` (which only has `SubmitOrder`). If
the calling code needs the full `engine.Engine` after rebuilding, capture the engine
reference in the loader closure before passing the closure to `Rebuild`:

```go
var eng engine.Engine
events.Rebuild(ch, channelID, func(snap []byte) (events.EngineState, error) {
    e, err := loader(snap)
    eng = e   // capture the concrete engine.Engine
    return e, err
})
// eng is now the restored engine.Engine
```

### Covering error paths in stub-backed functions

Stub functions (e.g. `classical.Load`) never return errors by design. Any wrapper
that only delegates to the stub will have an unreachable error branch. The established
pattern is to extract a private helper that accepts the dependency as a parameter:

```go
// Public entry point — always succeeds with the real stub.
func Load(snap []byte) (Engine, error) { return loadFromSnapshot(snap, classical.Load) }

// Testable helper — tests inject a failing loader to cover the error branch.
func loadFromSnapshot(snap []byte, loader func([]byte) (godip.Adjudicator, error)) (Engine, error) { ... }
```

See `loadFromSnapshot` and `newFromVariant` in `engine/adapter.go` for examples.

---

## Legacy: Custom Adjudicator

> **Preserved but not the active development path.** Adjudication has been outsourced to godip.
> The `game/` package is kept for reference.

The `game/` package implements a partial custom Diplomacy adjudicator based on the DATC
(Diplomacy Adjudicator Test Cases). A local copy of the DATC is at `DATC.txt`.

Phase 1 (army resolution — DATC 6.A.11, 6.A.12, 6.C.1–6.C.3, 6.D.1–6.D.5 and others) is
done or in progress. Phases 2–5 (self-dislodgement, head-to-head, fleet model, convoys) are
not planned.
