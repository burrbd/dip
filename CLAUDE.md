# CLAUDE.md

## Project Vision

This is a Diplomacy game engine intended to become a **Slack or Telegram bot**. Players submit moves via slash commands (e.g. `/move A Vie-Bud`), view the board map, and see the move history. The game is modelled as an **event log** — each turn's orders and resolutions are recorded as events, enabling replay, history browsing, and audit.

This repo is the backend engine. The bot layer (Slack/Telegram integration, event log storage, map rendering) is the next major piece of work.

---

## DATC: The Canonical Rules Reference

The **Diplomacy Adjudicator Test Cases (DATC)** by Lucas B. Kruijswijk are the canonical specification for correct Diplomacy adjudication. Every resolution rule in this engine must conform to the DATC. The full document is at: https://web.inter.nl.net/users/L.B.Kruijswijk/#5

A copy of the DATC is checked in at `DATC.txt` for offline reference.

The DATC is organised into sections:

| Section | Topic | Engine dependency |
|---------|-------|-------------------|
| 6.A | Basic checks (illegal moves, simple bounces) | Army model (partial), fleet model |
| 6.B | Coastal issues | Fleet model + coast designators |
| 6.C | Circular movement | Army model (partial), convoy |
| 6.D | Supports and dislodges | Army model (partial), fleet model, convoy, country rules |
| 6.E | Head-to-head battles and beleaguered garrison | Head-to-head algorithm, fleet model |
| 6.F | Convoys | Convoy model |
| 6.G | Convoying to adjacent provinces | Convoy model + adjacent-convoy rules |

### DATC Implementation Order

Work through the DATC in this order. Each phase has a set of integration tests in `game/resolver_integration_test.go` — uncomment the relevant block, watch it fail, implement the feature, watch it pass.

**Phase 1 — Army resolution (done/in progress)**
Tests: 6.A.11, 6.A.12, 6.C.1–6.C.3, 6.D.1–6.D.5, 6.D.9, 6.D.14, 6.D.15, 6.D.21, 6.D.25, 6.D.26, 6.D.33
These use only army moves, holds, and supports on the existing land territory map.
Uncomment the block marked `// Phase 1` in the integration test file.

**Phase 2 — Country/self-dislodgement rules**
Tests: 6.D.10, 6.D.11, 6.D.12, 6.D.13, 6.D.20
A unit may not dislodge a unit of the same power, even with foreign help. Requires tracking unit nationality and comparing it during conflict resolution.
Entry points: `game/resolver_main_phase.go`, `game/order_handler.go`.

**Phase 3 — Head-to-head battle algorithm**
Tests: 6.E.1–6.E.15
When two units move into each other's territories simultaneously, they fight a head-to-head battle. A dislodged head-to-head loser has no effect on the winner's origin territory. The current simultaneous-pass algorithm does not correctly model this; it needs a dedicated head-to-head pre-pass before general conflict resolution.
Entry points: `game/resolver_main_phase.go`.

**Phase 4 — Fleet model and coastal territories**
Tests: 6.A.1–6.A.10, 6.B.1–6.B.15, remaining 6.D (fleet-based), remaining 6.E
Fleets move along sea routes; adjacency differs from armies. Requires:
- Sea territory additions to the board (`game/order/board/territory.go`)
- `CreateFleetGraph()` analogous to `CreateArmyGraph()`
- Coast designators on coastal territories (e.g. `Spain(nc)`, `Spain(sc)`)
- `ValidateMove` switching on unit type
Entry points: `game/order/board/territory.go`, `game/order/validator.go`.

**Phase 5 — Convoy model**
Tests: 6.C.4–6.C.9, 6.D.6–6.D.8, 6.D.16, 6.D.27–6.D.32, 6.F.1–6.F.25, 6.G.1–6.G.10
An army can be convoyed across sea territories via a chain of fleets. A convoy is disrupted if any fleet in the chain is dislodged. Paradox cases (6.F.14+) require the Szykman rule (preferred by the DATC author).
Entry points: `game/resolver_main_phase.go`, `game/order_handler.go`, `game/order/order.go`.

---

## Development Philosophy

**TDD with 100% test coverage is non-negotiable.** The workflow is strictly red → green → refactor:

1. **Red** — uncomment the next DATC spec (or write a new unit test), run `go test`, confirm it fails for the right reason.
2. **Green** — write the minimal code to make it pass; don't over-engineer.
3. **Refactor** — once green, simplify and optimise. The adjudication algorithms in this domain are well-studied; feel free to read and borrow ideas from other open-source Diplomacy implementations (e.g. [jDip](https://jdip.sourceforge.net), [pydipcc](https://github.com/diplomacy/diplomacy), [godip](https://github.com/zond/godip)). Prefer clarity over cleverness, but don't leave a naïve O(n²) loop when a clean linear pass exists.

- `go test -v -cover -race ./...` must pass at all times
- Never commit with a failing test or with `focus: true` left set

The integration test file (`game/resolver_integration_test.go`) is the canonical way to specify new resolution scenarios. DATC tests are listed in DATC order with unimplemented ones commented out. To work on the next DATC case:
1. Check CLAUDE.md to see which phase is next
2. Find the corresponding commented-out spec in `game/resolver_integration_test.go`
3. Uncomment it, run tests, watch it fail
4. Implement until green
5. Refactor — simplify the algorithm, remove duplication, consult other implementations for inspiration

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

Each pass computes tentative outcomes for all conflicts simultaneously, then applies them. Passes repeat until no conflicts remain.

1. Collect all conflict groups (territorial and counter-attack)
2. Sort each group by strength descending
3. If one unit has strictly greater strength → defeats all others at that location
4. If tied → bounce all non-origin units; origin units stay
5. Apply all outcomes simultaneously, then repeat

Support strengths are calculated once before resolution begins (`OrderHandler.ApplyOrders`). Support is cut statically: any move targeting a supporter's territory cuts its support (DATC standard), unless the cutter is attacking from the territory being attacked (head-to-head — see `moveSupportCut`).

**Known gap:** Dislodgement during resolution does not retroactively cut the dislodged unit's support (DATC 6.D.17). Implementing this requires recalculating support strengths after each dislodgement. See Phase 3/4 tasks above.

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
