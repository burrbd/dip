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

### Command functional tests

**Every bot command that accepts arguments requires a functional test in `bot/bot_functional_test.go`.**
Commands with no arguments should also have a functional test if they have observable side effects
(e.g. `/pause`, `/concede`). The goal is end-to-end confidence that the command is wired up,
access-controlled, and produces the expected event — without mocking the bot or engine internals.

**How they work:**

1. Build tag `//go:build functional` keeps them out of the standard `go test ./...` run.
   Run explicitly with `go test -v -tags functional ./bot/`.
2. `startedGame(t)` is the shared helper that spins up a fresh `memChannel` + `Dispatcher`,
   runs `/newgame`, `/join England` (u1), `/join France` (u2), `/start` (gm), and returns the
   dispatcher and channel. Use it whenever the test needs an in-progress game.
3. `memChannel` is the in-memory `events.Channel` implementation (msgs, dms, imgs slices).
   `events.Scan` / `events.ScanDM` work correctly against it — no mocking.
4. Call `d.Dispatch(chanCmd(...))` or `d.Dispatch(dmCmd(...))` to invoke the command and assert
   on the response string and/or events written to `memChannel`.
5. Commands that are only valid in Retreat or Adjustment phase (unreachable from a fresh Spring
   1901 game) must test the **phase-guard rejection path**: call the command in Movement phase and
   assert `err != nil`. This proves the command is wired and correctly access-controlled.

**Checklist when adding a new command:**

- [ ] Add `TestCommand_<Name>` to `bot/bot_functional_test.go`
- [ ] If the command requires arguments, verify the response references the argument (province,
      nation, duration, etc.)
- [ ] If the command writes an event, assert `hasEvent(t, ch, channelID, events.Type...)` and
      unmarshal the payload to check key fields
- [ ] If the command is GM-only or nation-only, add a second sub-test asserting that an
      unauthorized caller receives `err != nil`
- [ ] If the command is phase-restricted, use a phase-specific helper (`retreatPhaseGame`,
      `adjustmentPhaseGame`) to reach the correct phase and test the happy path. A phase-guard
      rejection test (calling the command in the wrong phase) is acceptable as *additional*
      negative-path coverage but is not a substitute for a happy-path test.

**All 24 commands in ARCHITECTURE.md have functional tests as of Story 9.** The mapping is:

| Command | Test function(s) |
|---|---|
| `/newgame` | `TestCommand_Newgame` |
| `/join` | `TestCommand_Join`, `TestCommand_Join_RejectedIfNationTaken` |
| `/start` | `TestCommand_Start`, `TestCommand_Start_RejectedByNonGM` |
| `/order` | `TestCommand_Order`, `TestCommand_Order_RejectedForWrongNation` |
| `/orders` | `TestCommand_Orders` |
| `/clear` | `TestCommand_Clear` |
| `/submit` | `TestCommand_Submit_PartialDoesNotAdvance`, `TestCommand_Submit_AllPlayersAdvancesPhase` |
| `/retreat` | `TestCommand_Retreat_RejectedOutsideRetreatPhase` (phase guard placeholder — Story 6a adds happy-path) |
| `/disband` | `TestCommand_Disband_RejectedOutsideRetreatAndAdjustmentPhase` (phase guard placeholder — Story 6a adds happy-path) |
| `/build` | `TestCommand_Build_RejectedOutsideAdjustmentPhase` (phase guard placeholder — Story 6a adds happy-path) |
| `/waive` | `TestCommand_Waive_RejectedOutsideAdjustmentPhase` (phase guard placeholder — Story 6a adds happy-path) |
| `/status` | `TestCommand_Status` |
| `/history` | `TestCommand_History_BeforeFirstResolution` |
| `/map` | `TestCommand_Map_NoArgs`, `TestCommand_Map_WithTerritoryAndRadius` |
| `/help` | `TestCommand_Help_NoArgs`, `TestCommand_Help_WithCommand` |
| `/draw` | `TestCommand_Draw_ProposesOnFirstCall`, `TestCommand_Draw_AllNationsEndGame` |
| `/concede` | `TestCommand_Concede` |
| `/pause` | `TestCommand_Pause` |
| `/resume` | `TestCommand_Resume` |
| `/extend` | `TestCommand_Extend` |
| `/force-resolve` | `TestCommand_ForceResolve` |
| `/boot` | `TestCommand_Boot` |
| `/replace` | `TestCommand_Replace` |

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
1. Add it to `go.mod` with the correct version
2. Run `go mod tidy` and `go mod vendor` (requires network; do once then commit the result)
3. Or manually create the package files under `vendor/<module-path>/` and update `vendor/modules.txt`

**godip v0.6.5 is now vendored.** `vendor/github.com/zond/godip/` contains the real library source. The previous minimal stub has been replaced.

### godip API — what changed from the stub

The original stub modelled `godip.Adjudicator` as the game-state object (with `Phase()`, `Units()`, `Next() (Adjudicator, error)`, etc.). In **real godip v0.6.5** the types are:

| Concept | Stub (old) | Real godip v0.6.5 |
|---|---|---|
| Game state | `godip.Adjudicator` | `*state.State` |
| Per-province order | `godip.Order` | `godip.Adjudicator` |
| Phase advance | `Next() (Adjudicator, error)` | `(*state.State).Next() error` (in-place) |
| Serialization | `Load([]byte) (Adjudicator, error)` | No built-in JSON; 6 raw maps via `Units()`, etc. |
| Order parser | `classical.Parser.Parse([]string)` | `classical.DATCOrder(text)` (regex, returns province + `godip.Adjudicator`) |

### engine/ internal interface shim

Because the real `godip.Adjudicator` is a per-province order (not the game state), the `engine` package defines its own internal interfaces to decouple tests from godip's concrete types:

```go
// adjOrder — minimal interface over any staged order (real or test stub)
type adjOrder interface { Type() godip.OrderType }

// gamePhase — engine view of a game phase
type gamePhase interface {
    Type() godip.PhaseType; Year() int; Season() godip.Season
    DefaultOrder(godip.Province) adjOrder
}

// gameState — engine view of the full game state
type gameState interface {
    Phase() gamePhase
    Orders() map[godip.Province]adjOrder
    Units() map[godip.Province]godip.Unit
    Dislodgeds() map[godip.Province]godip.Unit
    SupplyCenters() map[godip.Province]godip.Nation
    SetOrder(godip.Province, adjOrder)
    Resolve(godip.Province) error
    Next() (gameState, error)
    SoloWinner() godip.Nation
    Dump() ([]byte, error)
}
```

`stateWrapper` (wraps `*state.State` + `common.Variant`) and `phaseWrapper` (wraps `godip.Phase`) are the production implementations. Test files use `mockAdj`/`mockPhase` that satisfy the same interfaces without touching real godip.

### Known stubs still in engine/

- **`buildStateFromSnapshot` error paths** — `SetUnits` and `SetDislodgeds` rarely return errors in practice; these branches exist for safety but are not reachable in normal test scenarios (engine coverage sits at ~98%).

### classical.DATCOrder — text format reference

`classical.DATCOrder(text string) (godip.Province, godip.Adjudicator, error)` in
`vendor/github.com/zond/godip/variants/classical/datc.go` uses case-insensitive regex to
parse standard Diplomacy notation. Accepted formats:

| Order type | Example text |
|---|---|
| Move | `"A Vie-Bud"`, `"F Nap-Ion via convoy"` |
| Hold | `"A Vie H"`, `"A Vie Hold"` |
| Support Hold | `"A Vie S A Bud"` |
| Support Move | `"A Mar S A Par-Bur"`, `"A Mar S Par-Bur"` (unit type optional) |
| Convoy | `"F Ion C A Nap-Tun"` |
| Build | `"Build A Vie"`, `"Build F Tri"` |
| Disband | `"A Vie disband"` |
| Remove | `"remove A Vie"` |

Key properties of the returned values:
- **Province names are always lowercase** — e.g. `"A Vie-Bud"` returns source province `"vie"`, not `"Vie"`. All godip province lookups are case-sensitive and lowercase.
- **`Targets()` for a Move order returns exactly `[src, dst]`** — two elements, both lowercase provinces. This is reliable and safe to index directly.
- **`godip.Order` interface requires `At() time.Time`** — easy to miss when writing test mocks; omitting it causes a compile-time "missing method At" error.

### Resolve() + Advance() coordination

`engine.Resolve()` now calls `fillNMR()` + `Next()` internally and compares pre/post unit
positions to populate `OrderResult.Success` accurately via `moveSucceeded()`. To avoid a
double `Next()` call, the `game` struct carries an `advanced bool` flag: `Resolve()` sets
it; `Advance()` skips the main `Next()` call when `advanced=true` and only handles
empty-phase skipping.

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

### Covering error paths in functions that rarely fail

Some functions (e.g. `classical.Start`, `(*state.State).Next`) never return errors in
practice. Any wrapper that only delegates to them will have an unreachable error branch.
The established pattern is to extract a private helper that accepts the dependency as a
parameter so tests can inject a failing version:

```go
// Public entry point — delegates to the real function.
func classicalLoader(snap []byte) (gameState, error) {
    return classicalLoaderWith(snap, classical.Start)
}

// Testable helper — tests inject a failing startFn to cover the error branch.
func classicalLoaderWith(snap []byte, startFn func() (*state.State, error)) (gameState, error) { ... }
```

See `classicalLoaderWith`, `classicalStartWith`, and `nextWith` in `engine/adapter.go` for
examples of this pattern applied throughout the engine package.

---

## Legacy: Custom Adjudicator

> **Preserved but not the active development path.** Adjudication has been outsourced to godip.
> The `game/` package is kept for reference.

The `game/` package implements a partial custom Diplomacy adjudicator based on the DATC
(Diplomacy Adjudicator Test Cases). A local copy of the DATC is at `DATC.txt`.

Phase 1 (army resolution — DATC 6.A.11, 6.A.12, 6.C.1–6.C.3, 6.D.1–6.D.5 and others) is
done or in progress. Phases 2–5 (self-dislodgement, head-to-head, fleet model, convoys) are
not planned.
