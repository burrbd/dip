# CLAUDE.md

## Project Vision

This repo is a **Diplomacy messenger bot** — players submit moves via slash commands (e.g. `/order A Vie-Bud`), view the board map, and see move history. Each channel hosts exactly one game. The game state is modelled as an **event log**: structured JSON snapshots posted to the channel after each phase resolution serve as the audit trail and persistence layer (no external database required).

**Adjudication is handled by [godip](https://github.com/zond/godip)**, which provides complete DATC compliance across movement, retreat, and adjustment phases. This repo's job is the bot layer on top of godip: command parsing, session management, event logging, map rendering, and platform integration (Slack, Telegram).

The full bot architecture is documented in [`ARCHITECTURE.md`](ARCHITECTURE.md).

---

## Build Plan

Ordered stories toward a working messenger bot. Work through them in sequence — each story is a discrete, shippable unit that unblocks the next.

### Story 1 — godip Engine Adapter
**Goal:** Thin wrapper around godip that hides its internal state types from the rest of the bot.

**Files:** `engine/adapter.go`, `engine/phases.go`, `engine/parser.go`, `engine/winner.go`

**Acceptance criteria:**
- `engine.NewGame(variant)` creates a godip `state.State` and returns an opaque `Engine`
- `Engine.SubmitOrder(nation, orderText)` parses via `classical.Parser` and stages the order
- `Engine.Resolve()` calls godip's adjudicator and returns a result summary
- `Engine.Advance()` calls `Next()`, fills NMR orders via `DefaultOrder()`, skips empty retreat/adjustment phases
- `Engine.SoloWinner()` returns the winning nation or empty string
- All public functions have unit tests

---

### Story 2 — Event Log
**Goal:** Write and read structured JSON events from the chat channel history.

**Files:** `events/types.go`, `events/log.go`, `events/replay.go`

**Event types:** `GameCreated`, `PlayerJoined`, `GameStarted`, `OrderSubmitted`, `PhaseResolved`, `PhaseSkipped`, `NMRRecorded`, `DrawProposed`, `DrawVoted`, `GameEnded`

**Acceptance criteria:**
- `events.Write(channelID, event)` posts a JSON-encoded event message to the channel
- `events.Scan(channelID)` reads channel history and returns typed events in order
- `events.Rebuild(channelID)` finds the last `PhaseResolved` or `GameStarted` snapshot, calls `state.Load()`, then replays subsequent `OrderSubmitted` events as staged orders
- Unit tests cover serialization round-trips and replay logic

---

### Story 3 — Session Management
**Goal:** Own the lifecycle of a single game within a channel.

**Files:** `session/session.go`, `session/store.go`, `session/lifecycle.go`

**Acceptance criteria:**
- `Session` struct holds: current phase, staged orders map, player→nation map, deadline timer, GM user ID
- `session.Load(channelID)` rebuilds session from event log (via `events.Rebuild`)
- `session.AdvanceTurn()` runs: collect staged orders → NMR fill → `Engine.Resolve()` → post `PhaseResolved` event → notify players → `Engine.Advance()` → start deadline timer for next phase
- Deadline timer fires `AdvanceTurn()` automatically when it expires
- Unit tests cover timer cancellation and state transitions

---

### Story 4 — Bot Command Router + Game Setup
**Goal:** Platform-agnostic command dispatch with access control; implement game setup commands.

**Files:** `bot/commands.go`, `bot/formatter.go`

**Commands:** `/newgame [settings]`, `/join [country]`, `/start`

**Acceptance criteria:**
- `bot.Dispatch(cmd)` routes to the correct handler based on command name and current phase
- Access control: `/start` and GM commands require the caller to be the GM
- `/newgame` posts `GameCreated` event, sets GM to caller
- `/join` posts `PlayerJoined` event; rejects if game already started or nation taken
- `/start` validates 2–7 players joined, posts `GameStarted` with `godip.Dump()` snapshot, starts Spring Movement deadline
- `bot.Format*` helpers render results and board state as plain text
- Unit tests use a mock channel/session

---

### Story 5 — Movement Phase Commands
**Goal:** Players submit and manage orders during the Movement phase.

**Files:** `bot/commands.go` (extended)

**Commands:** `/order <order-text>`, `/orders`, `/clear [order]`, `/submit`

**Acceptance criteria:**
- `/order` parses order text via `engine.Parser`, validates it belongs to the caller's nation, stages it
- `/orders` lists the caller's staged orders for the current phase
- `/clear` removes one or all of the caller's staged orders
- `/submit` marks the caller's orders as final; once all nations have submitted, `AdvanceTurn()` fires immediately
- Posts `OrderSubmitted` event after each `/order`
- Unit tests cover invalid orders, wrong phase, and early resolution trigger

---

### Story 6 — Retreat & Adjustment Commands
**Goal:** Handle the Retreat and Adjustment phases.

**Files:** `bot/commands.go` (extended)

**Commands:** `/retreat <unit> <province>`, `/disband <unit>`, `/build <unit-type> <province>`, `/waive`

**Acceptance criteria:**
- Retreat commands only accepted during Retreat phase; adjustment commands only during Adjustment phase
- `/retreat` and `/disband` validated against godip's valid retreat destinations
- `/build` validates supply-centre ownership and build slot availability
- `/waive` stages a waive order for one available build
- Auto-disband via godip `PostProcess` for unordered retreat units
- Unit tests cover phase-guard rejections and NMR auto-fill

---

### Story 7 — Info Commands
**Goal:** Players and observers can inspect game state at any time.

**Files:** `bot/commands.go` (extended), `bot/autocomplete.go`

**Commands:** `/map [territory [n]]`, `/status`, `/history <turn>`, `/help [command]`

**Acceptance criteria:**
- `/status` shows current phase, year, SC counts, and order submission status per nation
- `/history <turn>` fetches the `PhaseResolved` event for that turn from the event log and formats results
- `/map` with no args posts the current board PNG (via `dipmap.Render`)
- `/map Vienna 1` highlights Vienna and all adjacent provinces (BFS radius 1 via `dipmap.Neighborhood`)
- `/help` lists all commands; `/help <command>` shows detailed usage
- `bot.Autocomplete(session, nation)` returns valid order strings for the nation's units

---

### Story 8 — Draw & GM Commands
**Goal:** End-game draw mechanics and game-master admin tools.

**Files:** `bot/commands.go` (extended)

**Commands:** `/draw`, `/concede`, `/pause`, `/resume`, `/extend <duration>`, `/force-resolve`, `/boot <nation>`, `/replace <nation> <user>`

**Acceptance criteria:**
- `/draw` proposes a draw; posts `DrawProposed`; resolves to `GameEnded` when all remaining nations vote yes via `/draw`
- `/concede` ends the game immediately with that nation conceding; posts `GameEnded`
- `/pause` cancels the deadline timer; `/resume` restarts it
- `/extend <duration>` adds time to the current deadline
- `/force-resolve` triggers `AdvanceTurn()` immediately (GM only)
- `/boot <nation>` removes a player; their orders are NMR'd each turn going forward
- `/replace <nation> <user>` transfers a nation to a new player

---

### Story 9 — Map Rendering
**Goal:** Convert godip's SVG map assets to PNG and post to channel.

**Files:** `dipmap/render.go`, `dipmap/highlight.go`, `dipmap/neighborhood.go`

**Acceptance criteria:**
- `dipmap.Render(state)` converts godip's SVG for the current board state to a PNG byte slice
- `dipmap.Highlight(svg, provinces)` colours a set of provinces distinctly
- `dipmap.Neighborhood(graph, territory, n)` returns all provinces within `n` hops via BFS over `graph.Edges()`; `n=0` returns only the territory itself
- Unit tests cover BFS boundary cases (n=0, n=1, disconnected graph)

---

### Story 10 — Slack Platform Adapter
**Goal:** Deploy the bot as a Slack app.

**Files:** `platform/slack/adapter.go`, `cmd/slackbot/main.go`

**Acceptance criteria:**
- Handles Slack slash command HTTP requests; parses into `bot.Command` values
- Handles Slack Events API payloads (URL verification, event dispatch)
- Posts text responses and PNG images back to Slack channels
- `cmd/slackbot/main.go` wires up HTTP server, Slack signing-secret verification, and `bot.Dispatch`

---

### Story 11 — Telegram Platform Adapter
**Goal:** Deploy the bot as a Telegram bot.

**Files:** `platform/telegram/adapter.go`, `cmd/telegrambot/main.go`

**Acceptance criteria:**
- Handles Telegram Bot API webhook updates; parses `/command` messages into `bot.Command` values
- Posts text responses and PNG images back to Telegram chats via Bot API
- `cmd/telegrambot/main.go` wires up HTTP server and `bot.Dispatch`

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

## Architecture

```
cmd/
  slackbot/          — Slack entry point (Story 10)
  telegrambot/       — Telegram entry point (Story 11)

bot/
  commands.go        — platform-agnostic command router + access control (Stories 4–8)
  autocomplete.go    — generate valid orders for current state (Story 7)
  formatter.go       — format results and board state as text (Story 4)

engine/
  adapter.go         — thin wrapper around godip state.State (Story 1)
  phases.go          — phase advance, NMR fill, phase-skip logic (Story 1)
  parser.go          — text order → godip via classical.Parser (Story 1)
  winner.go          — solo win / draw detection (Story 1)

session/
  session.go         — Session struct: phase, staged orders, player map, deadline, GM (Story 3)
  store.go           — serialize/deserialize via godip Dump/Load (Story 3)
  lifecycle.go       — turn advance orchestration (Story 3)

events/
  types.go           — event type constants + structs (Story 2)
  log.go             — write/read JSON events from channel (Story 2)
  replay.go          — rebuild state from snapshot + pending orders (Story 2)

dipmap/
  render.go          — SVG → PNG via godip SVG assets (Story 9)
  highlight.go       — colour a set of provinces (Story 9)
  neighborhood.go    — BFS expansion to radius n (Story 9)

platform/
  slack/adapter.go   — Slack slash commands + Events API (Story 10)
  telegram/adapter.go — Telegram Bot API (Story 11)

game/               — [LEGACY] Partial custom adjudicator (preserved, not active development)
  order/board/      — Territory map, unit types, position state machine
  order/            — Order structs, decoder, validator
```

### Data flow for a turn (bot layer)

```
Slash command string (e.g. "/order A Vie-Bud")
  → platform adapter     — parse into bot.Command
  → bot.Dispatch()       — access control + route to handler
  → session.Session      — stage order / advance turn
  → engine.Engine        — submit to godip, resolve, advance phase
  → events.Write()       — post PhaseResolved snapshot to channel
  → dipmap.Render()      — post updated board PNG to channel
```

---

## Package Conventions

- **Interfaces over concrete types** — keep packages decoupled and testable
- **Value semantics for orders** — order structs passed by value, never mutated
- **No external database** — channel message history is the persistence layer; state rebuilt via event replay

---

## Testing Conventions

- Use `github.com/cheekybits/is` for assertions: `is.NoErr(err)`, `is.Equal(a, b)`, `is.NotNil(v)`
- Unit tests live alongside the code they test (`_test.go` files in same package or `_test` package)
- Mock interfaces inline in test files
- Table-driven tests for parsers and validators

---

## Legacy: Custom Adjudicator

> **This work is preserved but is no longer the active development path.** Adjudication has been outsourced to godip. The `game/` package code below is kept for reference and may be useful for understanding Diplomacy adjudication concepts.

The `game/` package implements a partial custom Diplomacy adjudicator based on the DATC (Diplomacy Adjudicator Test Cases) by Lucas B. Kruijswijk. The full DATC document is at https://web.inter.nl.net/users/L.B.Kruijswijk/#5 and a local copy is at `DATC.txt`.

### Completed work

- **Phase 1 — Army resolution (done/in progress)**
  Tests: 6.A.11, 6.A.12, 6.C.1–6.C.3, 6.D.1–6.D.5, 6.D.9, 6.D.14, 6.D.15, 6.D.21, 6.D.25, 6.D.26, 6.D.33

### Remaining phases (not planned)

- Phase 2 — Country/self-dislodgement rules (6.D.10–6.D.13, 6.D.20)
- Phase 3 — Head-to-head battle algorithm (6.E.1–6.E.15)
- Phase 4 — Fleet model and coastal territories (6.A.1–6.A.10, 6.B.1–6.B.15)
- Phase 5 — Convoy model (6.C.4–6.C.9, 6.F.1–6.F.25, 6.G.1–6.G.10)

### Key types (legacy)

| Type | Package | Purpose |
|---|---|---|
| `Territory` | `board` | 56-territory map with adjacency edges |
| `Unit` | `board` | Army or Fleet, owned by a country |
| `PositionManager` | `board` | Tracks unit positions through a turn |
| `Move`, `Hold`, `MoveSupport`, `HoldSupport`, `MoveConvoy` | `order` | Typed order structs |
| `OrderHandler` | `game` | Applies orders + calculates support strengths |
| `ResolveOrders` | `game` | Iterative conflict resolution algorithm |

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
