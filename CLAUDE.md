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

Support strengths are computed once by `OrderHandler.ApplyOrders` using order-based support-cut detection (DATC standard: any attack on a supporter's territory cuts its support, with the exception that the attacked unit cannot cut support aimed at itself). Then:

1. Collect all conflict groups simultaneously (territorial + counter-attack)
2. For each group, sort by strength descending
3. If one unit has strictly greater strength → defeats all others at that location
4. If tied → bounce all attackers; any attacker already at its origin territory is defeated (no retreat available)
5. Apply all outcomes at once, then repeat until no conflicts remain (fixed-point)

This simultaneous fixed-point approach correctly handles interdependent conflicts (circular support chains, dislodgement that cuts support, etc.).

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

All tasks follow the TDD workflow: uncomment (or add) a failing spec in
`game/resolver_integration_test.go`, then implement until green.

The DATC test cases (https://web.inter.nl.net/users/L.B.Kruijswijk/#5) are the
canonical correctness reference. All DATC tests are pre-written in the
integration test file — commented-out tests note which task unblocks them.

DATC implementation order is listed below. Start by picking the next
uncommented test group, verify it passes, then move on.

---

### Already implemented

- **Simultaneous fixed-point resolver** (`game/resolver_main_phase.go`): all
  conflict groups are resolved in each pass; passes repeat until stable.
  Covers DATC 6.A.8, 6.A.11, 6.A.12, 6.C.1–6.C.3, 6.D.1–6.D.5,
  6.D.9, 6.D.15, 6.D.21.

- **Hold-support resolution** (`game/order_handler.go` `holdStrength`):
  `HoldSupport` orders accumulate hold strength onto the defending unit,
  which the resolver uses when comparing strengths. Hold support is cut if
  the supporter is attacked (same rule as move support, without the
  "attacker from supported territory" exception — hold support doesn't have
  a directional target to exempt).

---

### 1. Country-aware self-dislodgement prohibition

**Why it matters:** DATC 6.D.10–6.D.14, 6.D.16, 6.D.20 require that a unit
cannot be dislodged by another unit of the same country, and that a power
cannot support a foreign unit to dislodge its own unit.

**DATC tests to uncomment:** 6.D.10, 6.D.11, 6.D.12, 6.D.13, 6.D.14, 6.D.16, 6.D.20

**Approach:**
- The test harness in `resolver_integration_test.go` currently uses a single
  country (`"a_country"`) for all units. Extend the harness to allow a
  per-result `country` field.
- In `OrderHandler.ApplyOrders`, skip or zero-strength any move that would
  dislodge a unit of the same country.
- In `moveStrength` / support-cut logic, void a support order if applying it
  would dislodge the supporter's own unit.

**Entry points:** `game/resolver_integration_test.go`, `game/order_handler.go`

---

### 2. Make the counter-attack conflict key separator explicit

**Why it matters:** `appendCounterAttackConflict` encodes territory pairs as
`"aaa.bbb"` and `Conflict()` detects them with `strings.Contains(key, ".")`.
This relies on the implicit invariant that no territory abbreviation contains
`"."`.

**Approach:** Define a named constant for the separator (e.g.
`const counterAttackSep = "."`), or store counter-attack conflicts in a
separate map keyed by a struct rather than a string.

**Entry point:** `game/order/board/manager.go`

---

### 3. Fleet movement + coasts (unlocks DATC 6.A.1–6.A.3, 6.A.9–6.A.10, all of 6.B, fleet variants in 6.D and 6.E)

**Why it matters:** Fleets move along sea routes and can only enter coastal
territories via the correct coast (e.g. Spain(NC) vs Spain(SC)).
`ValidateMove` currently uses only `CreateArmyGraph`; fleet orders always
pass adjacency validation.

**DATC tests to uncomment after this task:**
6.A.1, 6.A.2, 6.A.3, 6.A.4, 6.A.9, 6.A.10,
6.B.1–6.B.14,
6.D.7, 6.D.17–6.D.26, 6.D.28–6.D.31,
6.E.1–6.E.14 (head-to-head and beleaguered garrison)

**Approach:**
- Add 19 sea territories (NTH, ENG, IRI, NAO, MAO, WES, LYO, TYS, ION, AEG,
  EAS, MED, ADR, BLA, BAL, BOT, SKA, HEL, NWG, BAR) to `territory.go`.
- Extend `Territory` with a `Coast` designator (`""`, `"nc"`, `"sc"`) for
  multi-coast territories (SPA, BUL, STP).
- Build `CreateFleetGraph` with sea-route adjacency.
- Switch `ValidateMove` on `unit.Type` to use `CreateFleetGraph` for fleets.
- Fix three known territory-map bugs (see below).

**Known territory-map bugs to fix at the same time:**
- `mos.edges` has `"urk"` — should be `"ukr"` (Ukraine abbreviation).
- `bel.edges` has `"rur"` — should be `"ruh"` (Ruhr abbreviation).
- `yor.edges` has `"lvn"` (Livonia) — Yorkshire is not adjacent to Livonia;
  remove this edge.

**Entry points:** `game/order/board/territory.go`, `game/order/validator.go`

---

### 4. Convoy resolution (unlocks DATC 6.C.4–6.C.7, 6.D.6, 6.D.8, 6.D.27, all of 6.F, all of 6.G)

**Why it matters:** `MoveConvoy` is decoded but ignored during resolution.
Armies cannot cross sea territories without it.

**DATC tests to uncomment after this task:**
6.C.4, 6.C.5, 6.C.6, 6.C.7,
6.D.6, 6.D.8, 6.D.16, 6.D.27, 6.D.31,
6.F.1–6.F.24,
6.G.1–6.G.8

**Approach:**
- A convoy order `F Mid C A Lon-Bre` threads an army through a sea territory.
- During resolution, find chains of `MoveConvoy` orders that form a continuous
  route; if any fleet in the chain is dislodged the convoy fails and the army
  bounces.
- Integration specs: successful convoy; convoy disrupted by attack on fleet.

**Entry points:** `game/resolver_main_phase.go`, `game/order_handler.go`,
`game/order/order.go` (`MoveConvoy`)

---

### 5. Retreat phase (unlocks DATC 6.H.1–6.H.12)

Dislodged units must retreat to an adjacent unoccupied territory or disband.
This requires a second resolution pass after the main phase.

**Entry points:** new `game/resolver_retreat_phase.go`;
extend `board.PositionManager`

---

### 6. Build/adjustment phase (unlocks DATC 6.I.1–6.I.7)

After retreat, powers with more supply centres than units may build; powers
with fewer must disband.

**Entry points:** new `game/resolver_build_phase.go`

---

### 7. Graceful invalid-order handling

**Why it matters:** `order.Decode` and `order.Validator` return errors, but
there is no path to surface them to a caller cleanly. Invalid orders silently
become no-ops.

**Approach:**
- `game.OrderHandler` should collect validation errors per order and expose
  them alongside the resolved positions.
- Design the error type so the bot layer can report per-player which orders
  were rejected and why.

**Entry points:** `game/order_handler.go`, `game/order/validator.go`

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
