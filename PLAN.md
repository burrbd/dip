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
- [x] Story 9a — Mobile Map: Viewport Zoom and Lambda-Safe SVG→PNG
- [x] Story 9b — Unit Overlay: Draw Armies and Fleets on the Map
- [ ] Story 9c — Real SVG Rasterisation (oksvg + rasterx)
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
dislodgement (Italy `A Ven-Tri` + `A Rom S A Ven-Tri`, all others hold), and force-resolves.
Returns a dispatcher in Spring 1901 Retreat phase with Austria's `F Tri` dislodged.

Classical starting positions relevant to this scenario:
- Austria: A Vie, A Bud, **F Tri** (the unit that gets dislodged)
- Italy: **A Ven**, **A Rom**, F Nap (A Ven attacks Tri, A Rom supports)
- Tri's neighbours: Ven, Tyr, Vie, Adr, Alb — F Tri may retreat to any unoccupied one
  (Ven is vacated when A Ven moves out; Vie has Austria's A Vie so is blocked)
- Safe retreat choice for tests: `Ven` (vacant after Italy's move) or `Adr`

`adjustmentPhaseGame(t)` — spins up a full 7-nation game, submits `F Lon-NTH` for England in
Spring 1901 (all others hold), force-resolves through Spring (no dislodgements → Spring Retreat
skipped), submits `F NTH-NOR` in Fall 1901 (all others hold), and force-resolves through Fall.
Returns a dispatcher in Winter 1901 Adjustment phase with England owning Norway (+1 SC → 4 SCs,
3 units → 1 build slot available).

**Acceptance criteria:**
- `retreatPhaseGame(t)` and `adjustmentPhaseGame(t)` helpers exist and produce a dispatcher in
  the correct phase (verified by reading the current phase from the rebuilt session)
- `TestCommand_Retreat` — happy path: Austria submits `/retreat F Tri Ven`; no error returned;
  an `OrderSubmitted` event is recorded for Austria
- `TestCommand_Disband_InRetreatPhase` — happy path: Austria submits `/disband F Tri`; no error;
  `OrderSubmitted` event recorded
- `TestCommand_Build` — happy path: England submits `/build F Lon`; no error; `OrderSubmitted`
  event recorded
- `TestCommand_Waive` — happy path: England submits `/waive` (has 1 available build slot from
  Norway, chooses to waive it); no error; order staged
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

### Story 9a — Mobile Map: Viewport Zoom and Lambda-Safe SVG→PNG

**Goal:** Make `/map territory n` useful on small smartphone screens by cropping the
rendered image to the bounding box of the highlighted neighbourhood, and replace the
`rsvg-convert` shell-out with a pure-Go SVG→PNG renderer that works inside AWS Lambda.

**Background:**

`/map Vienna 2` currently highlights the neighbourhood provinces on the full Europa
board. On a smartphone the highlighted region is tiny and unreadable. The fix is to
compute the bounding box of the highlighted province shapes, expand it by a small
margin, set a `viewBox` attribute on the SVG `<svg>` element to that rectangle, and
then rasterise — giving a zoomed-in PNG that fills the image with only the relevant
region.

The current rasterisation step calls the system binary `rsvg-convert` via `os/exec`.
This works on a developer machine but is unavailable in the AWS Lambda execution
environment (the sandbox contains only a minimal set of binaries). It also makes the
binary a required system dependency, complicating Docker image builds and CI.

**Files:** `dipmap/render.go`, `dipmap/render_test.go`

**New dependency — pure-Go SVG rasteriser:**

Replace `rsvg-convert` with a vendored pure-Go SVG renderer.
Recommended library: `github.com/srwiley/oksvg` (SVG parsing) +
`github.com/srwiley/rasterx` (anti-aliased rasterisation) — both are MIT-licensed
and have no C dependencies.

Because the environment has no network access, the packages must be vendored manually:
1. Copy `github.com/srwiley/oksvg` and `github.com/srwiley/rasterx` source trees into
   `vendor/github.com/srwiley/oksvg/` and `vendor/github.com/srwiley/rasterx/`.
2. Add both to `go.mod` with the correct version tag.
3. Add entries to `vendor/modules.txt` following the existing pattern.

Alternative if oksvg proves inadequate for godip's SVG dialect: vendor
`github.com/tdewolff/canvas` (more complete SVG support, pure Go, MIT-licensed) using
the same manual-vendoring process.

**Acceptance criteria:**

- `svgToPNG` no longer shells out to `rsvg-convert`; it uses the vendored pure-Go
  renderer to convert SVG bytes → `image.RGBA` → PNG bytes entirely in-process.
- `go test -v -cover -race ./dipmap/` passes without `rsvg-convert` installed.
- `Render(state)` (no neighbourhood args) continues to render the full board at its
  natural dimensions.
- New function `RenderZoomed(state EngineState, svg []byte, provinces []string) ([]byte, error)`:
  - Receives the highlighted SVG produced by `Highlight`.
  - Computes the union bounding box of all `points` / `d` data for the listed provinces
    by parsing the coordinate sequences already extracted by `extractProvinceShape`.
  - Adds a configurable padding (e.g. 5% of the diagonal) around the bounding box.
  - Rewrites the `<svg … viewBox="…">` attribute to that padded bounding box.
  - Rasterises at a fixed output width (e.g. 800 px) preserving aspect ratio.
  - Returns a PNG byte slice.
- `dipmap.Neighborhood` result is passed through `Highlight` then `RenderZoomed` when
  `/map territory n` is called with `n > 0`; the full-board render is used when `n == 0`
  or no territory is given.
- Unit tests cover:
  - `RenderZoomed` with a minimal synthetic SVG — result is a valid PNG whose pixel
    dimensions reflect the zoomed bounding box, not the original canvas size.
  - `RenderZoomed` with an empty province list falls back to the full canvas viewBox.
  - `svgToPNG` round-trip produces non-empty PNG bytes (integration test, skipped if
    the renderer returns an unsupported-element error on the real godip SVG — document
    the skip reason).
- `go test -v -cover -race ./...` passes.

**`os/exec` audit — no other callers found:**

A codebase-wide search confirms that `os/exec` is used **only** in `dipmap/render.go`
(and its test). No other packages shell out to external binaries. This story eliminates
the last remaining `os/exec` dependency.

---

### Story 10 — Telegram Platform Adapter

**Goal:** Deploy the bot as a Telegram bot.

**Files:** `platform/telegram/adapter.go`, `platform/telegram/store.go`, `cmd/telegrambot/main.go`

**Acceptance criteria:**
- Handles Telegram Bot API webhook updates; parses `/command` messages into `bot.Command` values
- Posts text responses and PNG images back to Telegram chats via Bot API
- Implements all six `events.Channel` methods on the Telegram adapter:
  - `Post` / `History` — group chat messages; history backed by local JSONL file store
    (Telegram Bot API does not expose historical messages)
  - `SendDM` / `DMHistory` — private chat messages; history backed by local JSONL file store
  - `PostImage` — sends PNG to group chat via `sendPhoto`
  - `SendDMImage` — sends PNG to a player's private chat via `sendPhoto`
- Handles private chat (`chat.type = "private"`) update payloads and routes them to the order handler
- `cmd/telegrambot/main.go` reads `TELEGRAM_BOT_TOKEN`, `DATA_DIR`, `PORT` from env; wires up
  HTTP server, webhook registration, and `bot.Dispatch`
- Unit tests cover all Channel methods using a mock Telegram API server
- `go test -v -cover -race ./...` passes

---

### Story 9b — Unit Overlay: Draw Armies and Fleets on the Map

**Goal:** Display army and fleet positions on the rendered map. Province
centroids are computed from each province's polygon bounding box; a labelled
circle glyph is injected into the SVG before rasterisation.

**Files:** `dipmap/overlay.go`, `dipmap/overlay_test.go`, `dipmap/render.go`,
`bot/commands.go`, `bot/commands_test.go`

**Acceptance criteria:**

- New `dipmap.Unit` struct: `{ Type string; Nation string }` where `Type` is
  `"Army"` or `"Fleet"` and `Nation` is e.g. `"England"`.
- New function `dipmap.Overlay(svg []byte, units map[string]Unit) ([]byte, error)`:
  - Computes the centroid of each province's polygon from its points/d data.
  - Injects an SVG `<circle>` + `<text>` glyph for each unit, coloured by
    nation (standard Diplomacy palette). Unknown nations default to `#333333`.
  - Army label `"A"`, Fleet label `"F"`.
  - Provinces not found in the SVG are silently skipped.
  - Empty unit map returns the original SVG unchanged.
- `dipmap.SVGToPNG(svg []byte) ([]byte, error)` — exported wrapper over the
  internal `svgToPNG`; default `pngFn` for `bot.Dispatcher`.
- `bot.Dispatcher.handleMap` pipeline (both full-board and zoomed):
  1. `svgFn` — load raw SVG
  2. `overlayFn` — inject unit glyphs (converts `engine.UnitInfo` → `dipmap.Unit`)
  3a. Full-board (`n == 0` or no territory): `pngFn(svg)` → PNG
  3b. Zoomed (`n > 0`): `highlightFn` → `renderZoomedFn` → PNG
- `renderFn` removed from `Dispatcher`; replaced by `pngFn func([]byte)([]byte, error)`.
- New injectable `overlayFn` field on `Dispatcher`.
- Unit tests cover:
  - `Overlay` with known province → SVG contains `<g id="units">` and glyph
  - `Overlay` with unknown province → SVG unchanged
  - `Overlay` with empty map → SVG unchanged
  - `provinceCenter` valid / unknown / no-numeric-coords cases
  - `unitGlyph` for Army+known nation and Fleet+unknown nation
  - Bot: `TestDispatchMap_RejectsSVGLoadError`, `TestDispatchMap_RejectsOverlayError`,
    `TestDispatchMap_RejectsPNGError`, `TestDispatchMap_OverlaysUnitsOnMap`
- `go test -v -cover -race ./...` passes with 100% dipmap and bot coverage.

---

### Story 9c — Real SVG Rasterisation (oksvg + rasterx)

**Goal:** Replace the blank-white stub renderer in `svgToPNGWith` with a real
pure-Go SVG rasteriser so that the map image posted to the channel shows actual
board geography (land, sea, province borders) and unit glyphs.

**Pre-requisite:** vendor `oksvg` + `rasterx` + `golang.org/x/image` (which
`oksvg` depends on transitively). Run once with network access:

```bash
go get github.com/srwiley/rasterx@latest
go get github.com/srwiley/oksvg@latest
go mod tidy
go mod vendor
```

**Files:** `dipmap/render.go`, `dipmap/render_test.go`

**Acceptance criteria:**

- `svgToPNGWith` replaces `image.NewRGBA` + blank fill with an `oksvg` render
  pass: parse the SVG with `oksvg.ReadIconStream`, create an `rasterx` scanner
  backed by an `image.RGBA`, call `icon.Draw(scanner, 1.0)`, then encode via
  `encoderFn`.
- `Render(state)` posts a visually correct full-board PNG (province polygons
  filled, borders visible).
- `RenderZoomed` posts a zoomed crop with highlighted provinces and unit glyphs
  visible.
- If `oksvg` returns an unsupported-element error for the real godip SVG, the
  integration test is skipped with a documented reason (the existing
  `TestRender_ReturnsPNGBytes` skip pattern).
- A new `TestSVGToPNGWith_RendersContent` test decodes the returned PNG and
  asserts that at least one non-white pixel exists (proves the rasteriser fired,
  not the stub). This test is skipped if oksvg errors on the test SVG.
- All existing tests continue to pass; `go test -v -cover -race ./...` at 100%
  for `dipmap`.

---

### Story 11 — Slack Platform Adapter

**Goal:** Deploy the bot as a Slack app.

**Files:** `platform/slack/adapter.go`, `cmd/slackbot/main.go`

**Acceptance criteria:**
- Handles Slack slash command HTTP requests; parses into `bot.Command` values
- Handles Slack Events API payloads (URL verification, event dispatch)
- Posts text responses and PNG images back to Slack channels
- Implements all six `events.Channel` methods on the Slack adapter:
  - `Post` / `History` — Slack reads history via `conversations.history` API (no local store needed)
  - `SendDM` / `DMHistory` — Slack DM channel; history via `conversations.history` API
  - `PostImage` — uploads PNG to group channel via `files.upload`
  - `SendDMImage` — uploads PNG to the player's DM channel via `files.upload`
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
  - `SendDMImage` — uploads PNG to Twilio Media API, posts MMS link to player's 1:1 thread
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
