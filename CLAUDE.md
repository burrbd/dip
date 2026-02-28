# CLAUDE.md

## Project Vision

This is a Diplomacy game engine intended to become a **Slack or Telegram bot**. Players submit moves via slash commands (e.g. `/move A Vie-Bud`), view the board map, and see the move history. The game is modelled as an **event log** — each turn's orders and resolutions are recorded as events, enabling replay, history browsing, and audit.

This repo is the backend engine. The bot layer (Slack/Telegram integration, event log storage, map rendering) is the next major piece of work.

---

## Development Philosophy

**TDD with 100% test coverage is non-negotiable.** Every change starts with a failing test. No production code is written without a corresponding test that drove it.

- Write the test first — watch it fail
- Write the minimal code to make it pass
- Refactor if needed, keeping tests green
- `go test -v -cover -race ./...` must pass at all times

The integration test file (`game/resolver_integration_test.go`) is the canonical way to specify new resolution scenarios. Add a `spec` entry with a human-readable `description`, the orders, and expected positions — then implement until it passes.

The `focus` field on `spec` can be set to `true` to run only that scenario during development — remember to unset it before committing.

---

## Running Tests

```bash
go test -v -cover -race ./...
```

CI runs this same command via CircleCI (`.circleci/config.yml`).

---

## Architecture

```
game/order/board/   — Territory map, unit types, position state machine
game/order/         — Order data structures, string decoder, validator
game/               — Order application (OrderHandler) and conflict resolution (ResolveOrders)
```

### Data flow for a turn

```
Raw order string (e.g. "A Vie-Bud")
  → order.Decode()          — parse into typed order struct
  → order.Validator         — validate adjacency and ownership
  → game.OrderHandler       — apply orders to PositionManager (move/hold + support counts)
  → game.ResolveOrders()    — resolve conflicts iteratively until stable
  → board.PositionManager   — query final positions and defeated units
```

### Key types

| Type | Package | Purpose |
|---|---|---|
| `Territory` | `board` | 56-territory map with adjacency edges |
| `Unit` | `board` | Army or Fleet, owned by a country |
| `PositionManager` | `board` | Tracks unit positions through a turn; records Move/Hold/Bounce/Defeated events |
| `Move`, `Hold`, `MoveSupport`, `HoldSupport`, `MoveConvoy` | `order` | Typed order structs |
| `Set` | `order` | Collection of all orders for one turn |
| `Validator` | `order` | Graph-based adjacency validation |
| `OrderHandler` | `game` | Applies orders + calculates support strengths |
| `ResolveOrders` | `game` | Iterative conflict resolution algorithm |

### PositionEvent state machine

```
UnitPlaced → Moved | Held
Moved      → Bounced (on conflict loss)
Held       → Defeated (on conflict loss)
```

### Conflict resolution algorithm (`ResolveOrders`)

1. Find first conflict (multiple units in same territory, or counter-attack)
2. Sort by strength descending
3. If one unit has strictly greater strength → defeats all others at that location
4. If tied → bounce all attackers; defeat any attacker already at that territory (nowhere to retreat)
5. Repeat until no conflicts remain

**Known limitation:** This iterative one-conflict-at-a-time model is fundamentally incomplete. Many Diplomacy scenarios involve *interdependent* conflicts: whether A wins at X depends on whether B's support holds, which depends on whether B gets cut, which depends on A. Sequential resolution produces wrong outcomes for these cases regardless of ordering. The correct approach is a **simultaneous fixed-point algorithm**: in each pass, recompute all support strengths and tentative outcomes across the whole board, then repeat until the state stabilises. See the DATC (Diplomacy Adjudicator Test Cases) for the canonical test suite and algorithm specification.

---

## Package Conventions

- **Interfaces over concrete types** — `Manager`, `validator`, `simpleGraph` keep packages decoupled and testable
- **Value semantics for orders** — order structs are passed by value, never mutated
- **History-based position tracking** — `PositionManager` keeps a history slice per unit enabling bounce-back
- **Graph adjacency for validation** — gonum undirected graph (`CreateArmyGraph`) represents territory connectivity

---

## Testing Conventions

- Use `github.com/cheekybits/is` for assertions: `is.NoErr(err)`, `is.Equal(a, b)`, `is.NotNil(v)`
- Unit tests live alongside the code they test (`_test.go` files in same package or `_test` package)
- Integration tests live in `game/resolver_integration_test.go` using the `spec` table pattern
- Mock interfaces inline in test files (see `validator_test.go` for `mockGraph`)
- Table-driven tests for decoders and validators

---

## Open Tasks

Tasks are listed in priority order. Each one follows the TDD workflow: write a failing integration spec in `game/resolver_integration_test.go` (or a unit test alongside the relevant file), then implement until green.

---

### 1. Replace iterative resolver with a simultaneous fixed-point algorithm

**Why it matters:** The current `ResolveOrders` loop finds one conflict, resolves it, and repeats. Conflicts with circular support dependencies (e.g. A supports B into X, B's support is cut only if A fails) produce wrong outcomes because the algorithm resolves them in sequence rather than simultaneously.

**Approach:**
- Each resolution pass should compute tentative outcomes for *all* conflicts at once, treating unresolved conflicts as "pending"
- Support strengths are only counted from units not yet cut/disrupted
- Repeat passes until no outcomes change (fixed point)
- The DATC (https://web.inter.nl.net/users/L.B.Kruijswijk/#5) defines the canonical algorithm and 60+ test cases; use those as specs

**Entry points:** `game/resolver.go` (`ResolveOrders`), `game/handler.go` (`OrderHandler`)

---

### 2. Implement hold-support resolution

**Why it matters:** `HoldSupport` order structs are decoded and validated but `ResolveOrders` never applies them — a held unit always defends with strength 1 regardless of support.

**Approach:**
- In `OrderHandler.ApplyOrders`, accumulate hold-support strength onto the holding unit (same pattern as move-support)
- Add integration specs covering: supported hold repels unsupported attack; supported hold loses to stronger attack

**Entry points:** `game/handler.go`, `game/order/order.go` (`HoldSupport`)

---

### 3. Make the counter-attack conflict key separator explicit

**Why it matters:** `appendCounterAttackConflict` encodes the pair of territories as `"aaa.bbb"` and `Conflict()` detects it with `strings.Contains(key, ".")`. This works because no territory abbreviation contains `"."`, but the coupling is implicit.

**Approach:** Either (a) define a named constant for the separator and document the invariant, or (b) keep counter-attack and territorial conflict groups in separate maps so the type is structural rather than encoded in the key string.

**Entry point:** `game/order/board/manager.go`

---

### 4. Model territory coasts and implement fleet movement

**Why it matters:** Fleets move along sea routes and can only enter coastal territories via the correct coast (e.g. Spain(NC) vs Spain(SC)). `ValidateMove` currently only calls `CreateArmyGraph`; fleet orders always pass adjacency validation.

**Approach:**
- Extend `Territory` to carry a coast designator (e.g. `Coast string`, values `""`, `"nc"`, `"sc"`)
- Build a fleet/sea adjacency graph (`CreateFleetGraph`) analogous to `CreateArmyGraph`
- Switch `ValidateMove` on unit type to use the correct graph
- Add DATC coast-movement specs as integration tests

**Entry points:** `game/order/board/territory.go`, `game/order/validator.go`

---

### 5. Implement convoy resolution

**Why it matters:** `MoveConvoy` is decoded but ignored during resolution. Armies can't cross sea territories without it.

**Approach:**
- A convoy order `F Mid C A Lon-Bre` means the fleet threads the army through its sea territory
- During resolution, a chain of convoy orders forms a route; if any fleet in the chain is dislodged the convoy fails and the army bounces
- Add integration specs for: successful convoy; convoy disrupted by attack on fleet

**Entry points:** `game/resolver.go`, `game/handler.go`, `game/order/order.go` (`MoveConvoy`)

---

### 6. Graceful invalid-order handling

**Why it matters:** `order.Decode` and `order.Validator` return errors, but there is no path to surface these to a caller (bot layer) cleanly. Invalid orders currently cause silent no-ops or panics.

**Approach:**
- `game.OrderHandler` should collect validation errors per order and expose them alongside the resolved positions
- Design the error type so the bot layer can report per-player which orders were rejected and why

**Entry points:** `game/handler.go`, `game/order/validator.go`

---

## Future: Bot Integration

The planned interface layer will sit on top of this engine:

- **Slack**: slash commands (`/dip move A Vie-Bud`, `/dip map`, `/dip history`) via Slack's Events API or slash command webhooks
- **Telegram**: inline commands via the Bot API
- **Event log**: each turn is a sequence of events (OrderSubmitted, TurnResolved, UnitDefeated, etc.) stored in a database and replayable
- **Map rendering**: generate an image of the board state for posting to the channel after each resolution
- **Multi-player**: each country is controlled by a different Slack/Telegram user; the engine already models `Country` on units and orders
- **Turn management**: deadline-based turn advancement, order collection phase, resolution phase

New code for the bot layer should live in a new top-level package (e.g. `cmd/slackbot/` or `cmd/telegrambot/`) and import the engine packages — don't add bot concerns to the engine packages.

---

## Module

```
module github.com/burrbd/dip
go 1.12
```

Key dependencies:
- `gonum.org/v1/gonum` — graph library for territory adjacency
- `github.com/cheekybits/is` — test assertions
