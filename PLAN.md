# Build Plan

Ordered stories toward a working Diplomacy messenger bot. Each story is a discrete,
shippable unit of work that unblocks the next.

---

## For agents: how to use this file

1. **Find the next story** — the first unchecked box (`- [ ]`) is the current story to work on.
2. **Start the story** — read its goal, files, and acceptance criteria before writing any code.
3. **Work TDD** — write failing tests first, then implement, then refactor.
4. **Check it off** — once all acceptance criteria are met and `go test -v -cover -race ./...`
   passes, change `- [ ]` to `- [x]` for that story, then commit.
5. **Move on** — proceed to the next unchecked story.

Never mark a story done if any test is failing or any criterion is unmet.

---

## Stories

- [ ] Story 1 — godip Engine Adapter
- [ ] Story 2 — Event Log
- [ ] Story 3 — Session Management
- [ ] Story 4 — Bot Command Router + Game Setup
- [ ] Story 5 — Movement Phase Commands
- [ ] Story 6 — Retreat & Adjustment Commands
- [ ] Story 7 — Info Commands
- [ ] Story 8 — Draw & GM Commands
- [ ] Story 9 — Map Rendering
- [ ] Story 10 — Slack Platform Adapter
- [ ] Story 11 — Telegram Platform Adapter

---

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

**Event types:** `GameCreated`, `PlayerJoined`, `GameStarted`, `OrderSubmitted`, `PhaseResolved`,
`PhaseSkipped`, `NMRRecorded`, `DrawProposed`, `DrawVoted`, `GameEnded`

**Acceptance criteria:**
- `events.Write(channelID, event)` posts a JSON-encoded event message to the channel
- `events.Scan(channelID)` reads channel history and returns typed events in order
- `events.Rebuild(channelID)` finds the last `PhaseResolved` or `GameStarted` snapshot,
  calls `state.Load()`, then replays subsequent `OrderSubmitted` events as staged orders
- Unit tests cover serialization round-trips and replay logic

---

### Story 3 — Session Management

**Goal:** Own the lifecycle of a single game within a channel.

**Files:** `session/session.go`, `session/store.go`, `session/lifecycle.go`

**Acceptance criteria:**
- `Session` struct holds: current phase, staged orders map, player→nation map, deadline timer, GM user ID
- `session.Load(channelID)` rebuilds session from event log (via `events.Rebuild`)
- `session.AdvanceTurn()` runs: collect staged orders → NMR fill → `Engine.Resolve()` →
  post `PhaseResolved` event → notify players → `Engine.Advance()` → start deadline timer
- Deadline timer fires `AdvanceTurn()` automatically when it expires
- Unit tests cover timer cancellation and state transitions

---

### Story 4 — Bot Command Router + Game Setup

**Goal:** Platform-agnostic command dispatch with access control; game setup commands.

**Files:** `bot/commands.go`, `bot/formatter.go`

**Commands:** `/newgame [settings]`, `/join [country]`, `/start`

**Acceptance criteria:**
- `bot.Dispatch(cmd)` routes to the correct handler based on command name and current phase
- Access control: `/start` and GM commands require the caller to be the GM
- `/newgame` posts `GameCreated` event, sets GM to caller
- `/join` posts `PlayerJoined` event; rejects if game already started or nation taken
- `/start` validates 2–7 players joined, posts `GameStarted` with `godip.Dump()` snapshot,
  starts Spring Movement deadline
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
- `/submit` marks the caller's orders as final; once all nations have submitted,
  `AdvanceTurn()` fires immediately
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
- `/history <turn>` fetches the `PhaseResolved` event for that turn and formats results
- `/map` with no args posts the current board PNG (via `dipmap.Render`)
- `/map Vienna 1` highlights Vienna and all adjacent provinces (BFS radius 1 via `dipmap.Neighborhood`)
- `/help` lists all commands; `/help <command>` shows detailed usage
- `bot.Autocomplete(session, nation)` returns valid order strings for the nation's units

---

### Story 8 — Draw & GM Commands

**Goal:** End-game draw mechanics and game-master admin tools.

**Files:** `bot/commands.go` (extended)

**Commands:** `/draw`, `/concede`, `/pause`, `/resume`, `/extend <duration>`,
`/force-resolve`, `/boot <nation>`, `/replace <nation> <user>`

**Acceptance criteria:**
- `/draw` proposes a draw; posts `DrawProposed`; resolves to `GameEnded` when all remaining
  nations vote yes via `/draw`
- `/concede` ends the game immediately; posts `GameEnded`
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
- `dipmap.Neighborhood(graph, territory, n)` returns all provinces within `n` hops via BFS
  over `graph.Edges()`; `n=0` returns only the territory itself
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
