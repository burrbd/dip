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

- [x] Story 1 — godip Engine Adapter
- [x] Story 1a — Real Order Parsing & Resolve Accuracy
- [x] Story 2 — Event Log
- [x] Story 3 — Session Management
- [x] Story 4 — Bot Command Router + Game Setup
- [x] Story 5 — Movement Phase Commands
- [x] Story 6 — Retreat & Adjustment Commands
- [ ] Story 6a — Happy-path Functional Tests for Retreat & Adjustment Commands
- [x] Story 7 — Info Commands
- [x] Story 8 — Draw & GM Commands
- [x] Story 9 — Map Rendering
- [ ] Story 10 — Telegram Platform Adapter
- [ ] Story 11 — Slack Platform Adapter
- [ ] Story 12 — WhatsApp Platform Adapter (optional)
- [ ] Story 13 — Lambda / EventBridge Deployment

---

### Story 1 — godip Engine Adapter

**Goal:** Thin wrapper around godip that hides its internal state types from the rest of the bot.

**Files:** `engine/adapter.go`, `engine/phases.go`, `engine/parser.go`, `engine/winner.go`

**Acceptance criteria:**
- `engine.New(variant)` creates a godip `state.State` and returns an opaque `Engine`
- `Engine.SubmitOrder(nation, orderText)` parses order text and stages the order
- `Engine.Resolve()` returns a pre-advance summary of staged orders (province, order text, success)
- `Engine.Advance()` fills NMR orders via `DefaultOrder()`, calls godip `Next()`, skips empty retreat/adjustment phases
- `Engine.SoloWinner()` returns the winning nation or empty string
- `Engine.Dump()` / `engine.Load()` serialise and restore game state via a JSON snapshot (`stateSnapshot`)
- All public functions have unit tests

> **Known gaps (see Story 1a):**
> - `classicalOrderParser` is a text tokenizer only — it does not produce real `godip.Adjudicator`
>   orders. Staged orders from `SubmitOrder` are silently dropped by `stateWrapper.SetOrder` and
>   have no effect on real adjudication. NMR `DefaultOrder()` fills (which produce real
>   `godip.Adjudicator` holds) are unaffected.
> - `OrderResult.Success` in `Resolve()` is always `true` because godip adjudicates inside
>   `Next()` with no per-province API beforehand.

---

### Story 1a — Real Order Parsing & Resolve Accuracy

**Goal:** Wire up real godip order parsing so that player-submitted orders are actually
adjudicated, and make `Resolve()` report accurate success/failure after advancing.

**Files:** `engine/parser.go`, `engine/adapter.go`

**Acceptance criteria:**
- `classicalOrderParser.Parse` converts order text into a real `godip.Adjudicator` using the
  `github.com/zond/godip/orders` package (e.g. `orders.Move`, `orders.Hold`, `orders.Support`)
- `stateWrapper.SetOrder` stages the real adjudicator; the type-assertion guard can be removed
  once the parser produces real adjudicators
- `Engine.Resolve()` compares unit positions before and after `Next()` (or uses godip's
  resolution result map) to populate `OrderResult.Success` accurately
- All existing tests continue to pass; new tests cover a move that bounces (Success=false)
  and a supported move that succeeds (Success=true)

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
  calls `engine.Load()`, then replays subsequent `OrderSubmitted` events as staged orders
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

> **Note (refactor pending — Story 13):** The current implementation uses `timer *time.Timer` +
> `sync.Mutex`. Before Story 10, this must be refactored to a `Scheduler` interface
> (`LocalScheduler` / `EventBridgeScheduler`) as described in ARCHITECTURE.md. `AdvanceTurn()`
> must also gain an idempotency check (no-op if `PhaseResolved` already exists for the current
> phase). These changes are tracked in Story 13.

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
- `/start` validates 2–7 players joined, posts `GameStarted` with `engine.Dump()` snapshot,
  starts Spring Movement deadline
- `bot.Format*` helpers render results and board state as plain text
- Unit tests use a mock channel/session

---

### Story 5 — Movement Phase Commands

**Goal:** Players submit and manage orders during the Movement phase.

**Files:** `bot/commands.go` (extended)

**Commands:** `/order <order-text>`, `/orders`, `/clear [order]`, `/submit`

**Acceptance criteria:**
- `/order` is accepted only via DM to the bot (not in the game channel); parses order text via
  `engine.Parser`, validates it belongs to the caller's nation, stages it
- `/orders` (via DM) lists the caller's staged orders for the current phase
- `/clear` (via DM) removes one or all of the caller's staged orders
- `/submit` (via DM) marks the caller's orders as final; after each `/submit`, the handler reads
  all player DM threads to check whether every nation has submitted — if so, `AdvanceTurn()`
  fires immediately and the scheduler deadline is cancelled
- Posts `OrderSubmitted` event to the player's DM thread (not the game channel)
- Game channel receives no `OrderSubmitted` events — only the `PhaseResolved` result
- Unit tests cover invalid orders, wrong phase, DM-only enforcement, and early resolution trigger

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

> **Known gap (see Story 6a):** The functional tests for `/retreat`, `/disband`, `/build`, and
> `/waive` only exercise the phase-guard rejection path (calling the command in the wrong phase).
> They do not test the actual command behaviour in the correct phase. Story 6a adds proper
> happy-path functional tests using phase-specific game helpers.

---

### Story 6a — Happy-path Functional Tests for Retreat & Adjustment Commands

**Goal:** Replace phase-guard-only functional tests with end-to-end tests that exercise
`/retreat`, `/disband`, `/build`, and `/waive` in the correct game phase.

**Files:** `bot/bot_functional_test.go`

**New helpers:**

`retreatPhaseGame(t)` — spins up a full 7-nation game, submits Spring 1901 orders that cause a
dislodgement (Austria `A Vie-Tri` + `A Bud S A Vie-Tri`, all others hold), and force-resolves.
Returns a dispatcher in Spring 1901 Retreat phase with Italy's `F Tri` dislodged.

`adjustmentPhaseGame(t)` — spins up a full 7-nation game, submits `F Lon-NTH` for England in
Spring 1901 (all others hold), force-resolves through Spring, submits `F NTH-NOR` in Fall 1901
(all others hold), and force-resolves through Fall. Returns a dispatcher in Winter 1901
Adjustment phase with England owning Norway (+1 SC, needs to build).

**Acceptance criteria:**
- `retreatPhaseGame(t)` and `adjustmentPhaseGame(t)` helpers exist and produce a dispatcher in
  the correct phase (verified by reading the current phase from the rebuilt session)
- `TestCommand_Retreat` — happy path: Italy submits `/retreat F Tri Adr` (or another valid
  destination); no error returned; an `OrderSubmitted` event is recorded for Italy
- `TestCommand_Disband_InRetreatPhase` — happy path: Italy submits `/disband F Tri`; no error;
  `OrderSubmitted` event recorded
- `TestCommand_Build` — happy path: England submits `/build F Lon`; no error; `OrderSubmitted`
  event recorded
- `TestCommand_Waive` — happy path: France submits `/waive` (no build needed, so just confirm
  no error and the order is staged); or pick a nation that needs to waive
- Existing phase-guard tests (`TestCommand_Retreat_RejectedOutsideRetreatPhase`, etc.) are
  kept alongside the new happy-path tests — they remain valid as negative-path coverage
- `go test -v -tags functional ./bot/` passes

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

### Story 10 — Telegram Platform Adapter

**Goal:** Deploy the bot as a Telegram bot.

**Files:** `platform/telegram/adapter.go`, `platform/telegram/store.go`, `cmd/telegrambot/main.go`

**Acceptance criteria:**
- Handles Telegram Bot API webhook updates; parses `/command` messages into `bot.Command` values
- Posts text responses and PNG images back to Telegram chats via Bot API
- Implements all five `events.Channel` methods on the Telegram adapter:
  - `Post` / `History` — group chat messages; history backed by local JSONL file store
    (Telegram Bot API does not expose historical messages)
  - `SendDM` / `DMHistory` — private chat messages; history backed by local JSONL file store
  - `PostImage` — sends PNG to group chat via `sendPhoto`
- Handles private chat (`chat.type = "private"`) update payloads and routes them to the order handler
- `cmd/telegrambot/main.go` reads `TELEGRAM_BOT_TOKEN`, `DATA_DIR`, `PORT` from env; wires up
  HTTP server, webhook registration, and `bot.Dispatch`
- Unit tests cover all Channel methods using a mock Telegram API server
- `go test -v -cover -race ./...` passes

---

### Story 11 — Slack Platform Adapter

**Goal:** Deploy the bot as a Slack app.

**Files:** `platform/slack/adapter.go`, `cmd/slackbot/main.go`

**Acceptance criteria:**
- Handles Slack slash command HTTP requests; parses into `bot.Command` values
- Handles Slack Events API payloads (URL verification, event dispatch)
- Posts text responses and PNG images back to Slack channels
- Implements all five `events.Channel` methods on the Slack adapter:
  - `Post` / `History` — Slack reads history via `conversations.history` API (no local store needed)
  - `SendDM` / `DMHistory` — Slack DM channel; history via `conversations.history` API
  - `PostImage` — uploads PNG via `files.upload`
- Handles DM slash-command payloads (`channel_type = "im"`) and routes them to the order handler
- `cmd/slackbot/main.go` wires up HTTP server, Slack signing-secret verification, and `bot.Dispatch`
- Unit tests cover all Channel methods and webhook parsing
- `go test -v -cover -race ./...` passes

---

### Story 12 — WhatsApp Platform Adapter (optional)

**Goal:** Deploy the bot via the Twilio WhatsApp API or Meta Cloud API.

**Note:** WhatsApp requires a Meta Business Account (approval can take days/weeks) and has
per-conversation charges. Tackle only if Telegram/Slack do not meet deployment needs.

**Files:** `platform/whatsapp/adapter.go`, `platform/whatsapp/store.go`, `cmd/whatsappbot/main.go`

**Acceptance criteria:**
- `WhatsAppChannel` implements `events.Channel`:
  - `Post` / `History` — group messages sent via Twilio API; history backed by local JSONL file store
  - `SendDM` / `DMHistory` — 1:1 messages sent via Twilio API; history backed by local JSONL file store
  - `PostImage` — uploads PNG to Twilio Media API, posts MMS link to group
- Webhook handler validates `X-Twilio-Signature` and parses `application/x-www-form-urlencoded` payloads
- `cmd/whatsappbot/main.go` reads `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN`,
  `TWILIO_WHATSAPP_NUMBER`, `DATA_DIR`, `PORT` from env
- Unit tests cover all Channel methods using a mock Twilio API server
- `go test -v -cover -race ./...` passes

---

### Story 13 — Lambda / EventBridge Deployment

**Goal:** Refactor session scheduling to a `Scheduler` interface and wire up a Lambda entry
point so the bot runs as a stateless FaaS application with externally managed phase deadlines.

**Files:** `session/scheduler.go`, `platform/eventbridge/scheduler.go`, `cmd/lambdabot/main.go`

**Acceptance criteria:**
- `Scheduler` interface defined in `session/scheduler.go`:
  ```
  Schedule(channelID string, at time.Time) error
  Cancel(channelID string) error
  ```
- `LocalScheduler` implementation wraps `time.AfterFunc`; used in tests and server deployments
- `EventBridgeScheduler` implementation creates/deletes one-time AWS EventBridge Scheduler rules
  named by `channelID`; rule target is the Lambda function ARN (from environment variable)
- `Session.timer *time.Timer` and `Session.mu sync.Mutex` replaced by `Session.scheduler Scheduler`
- `GameStarted` and `PhaseResolved` event structs gain `DeadlineAt time.Time` (serialised as RFC3339)
- `AdvanceTurn()` gains an idempotency check: reads game channel history, no-ops if a
  `PhaseResolved` event already exists for the current phase
- `cmd/lambdabot/main.go` handles two event shapes:
  - Platform webhook payload → parse command → `bot.Dispatch`
  - `{"action": "advance_turn", "channel_id": "..."}` → `session.Load()` → `AdvanceTurn()`
- Unit tests cover `LocalScheduler` fire/cancel, idempotency guard (duplicate advance no-ops),
  and Lambda handler routing
- `go test -v -cover -race ./...` passes
