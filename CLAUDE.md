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

## What's Not Yet Implemented

From the README:

1. **Hold orders with support** — `HoldSupport` structs exist but resolution logic incomplete
2. **Fleet unit type validation** — `ValidateMove` only handles army graph; fleet movement (sea routes, coasts) not modelled
3. **Invalid order handling** — malformed or illegal orders need graceful rejection surfaced to the bot layer
4. **Territory coasts** — many territories have distinct land/sea coasts (e.g. Spain NC/SC); the map doesn't model these yet
5. **Convoy orders** — `MoveConvoy` struct exists, decoder handles `C` orders, but convoy resolution is not implemented

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
