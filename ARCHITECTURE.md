# Bot Architecture

This document describes the design of the Diplomacy bot layer that sits on top of the
[godip](https://github.com/zond/godip) adjudication engine.

The bot is deployed as a **stateless function** (AWS Lambda recommended; a long-running webhook
server is also supported). Each invocation rebuilds all necessary state from the event log — no
warm in-process state is required between invocations. Each chat channel hosts exactly one game.
The channel's message history is the audit trail — the bot posts a structured JSON snapshot after
each resolution so that state can be rebuilt on restart without an external database.

---

## Component breakdown

```
cmd/
  slackbot/          — Slack entry point (slash commands + Events API webhook)
  telegrambot/       — Telegram entry point (Bot API webhook)
  whatsappbot/       — WhatsApp entry point (Business API / Twilio)

bot/
  commands.go        — platform-agnostic command router + access control
  autocomplete.go    — generate valid orders / province choices for current state
  formatter.go       — format resolution results, board state, history as text

engine/
  adapter.go         — thin wrapper around godip state.State
  phases.go          — phase advance, NMR DefaultOrder() fill, phase-skip logic
  parser.go          — text order → godip Adjudicator via classical.Parser
  winner.go          — solo win / draw detection (polls SoloWinner after Fall Adjustment)

session/
  session.go         — Session struct: phase, staged orders, player map, scheduler, GM user ID
  scheduler.go       — Scheduler interface; LocalScheduler (time.AfterFunc) and EventBridgeScheduler
  store.go           — serialize/deserialize Session to/from snapshot JSON (godip Dump/Load)
  lifecycle.go       — turn advance: collect → NMR fill → adjudicate → snapshot → notify

events/
  types.go           — event type constants + structs
  log.go             — write structured JSON event to channel; scan channel history for events
  replay.go          — rebuild game state: find last PhaseResolved snapshot, apply pending orders

dipmap/
  render.go          — SVG → PNG conversion using godip SVG assets
  highlight.go       — highlight a province set
  neighborhood.go    — BFS expansion: given territory + radius n, return all provinces within n hops
                       n=0 → just the territory; n=1 → territory + adjacent; n=2 → +adjacent-of-adjacent

platform/
  slack/adapter.go
  telegram/adapter.go
  whatsapp/adapter.go
```

Note: the package is named `session/` (not `game/`) to avoid collision with the existing
engine package at `game/`.

---

## Channel interface

Defined in `events/log.go`. Platform adapters (Slack, Telegram, WhatsApp) must implement
all four methods.

```
Post(channelID, text string) error          — post message to group channel
History(channelID string) ([]string, error) — read group channel message history
SendDM(userID, text string) error           — send private message to a player
DMHistory(userID string) ([]string, error)  — read a player's DM thread history
```

`Post` / `History` operate on the shared game channel. `SendDM` / `DMHistory` operate on
private per-player threads. Order submission uses `SendDM`/`DMHistory` so that staged orders
are never visible to other players in the game channel.

---

## Scheduler interface

Defined in `session/scheduler.go`. Decouples deadline management from the execution model.

```
Schedule(channelID string, at time.Time) error — set (or overwrite) a one-time deadline
Cancel(channelID string) error                 — cancel the pending deadline
```

**Implementations:**

| Type | Backed by | Use case |
|---|---|---|
| `LocalScheduler` | `time.AfterFunc` | Long-running server / unit tests |
| `EventBridgeScheduler` | AWS EventBridge Scheduler | Lambda deployment |

`EventBridgeScheduler` creates a one-time schedule rule named by `channelID`. The rule
target is the Lambda function ARN (from environment). On `Cancel` it deletes the rule.

`Session.scheduler Scheduler` replaces the former `Session.timer *time.Timer` +
`Session.mu sync.Mutex` fields.

---

## What godip provides vs. what the bot layer owns

| Concern | godip | Bot layer |
|---|---|---|
| Movement adjudication | ✓ | |
| Retreat resolution | ✓ | |
| Build/disband processing | ✓ | |
| NMR `DefaultOrder()` | ✓ | |
| Phase sequencing (`Next()`) | ✓ | |
| Valid retreat destinations | ✓ | |
| DATC correctness | ✓ | |
| Win / draw detection | | ✓ |
| Deadlines / timers | | ✓ |
| Player notifications | | ✓ |
| Access control | | ✓ |
| Command parsing | | ✓ |
| Map rendering | SVG assets | ✓ render + post |
| Autocomplete | | ✓ |
| Event log / persistence | | ✓ |
| GM tools | | ✓ |

---

## Design decisions

| Decision | Choice | Rationale |
|---|---|---|
| **Press** | Out of scope — chat IS the press | Players negotiate in the channel; bot handles orders and adjudication only |
| **Event log** | Snapshot per turn | Bot posts `godip.Dump()` JSON after each resolution; bot restart restores from last snapshot |
| **GM role** | Game creator is GM | `/newgame` issuer gets admin privileges for that game |
| **Multi-game** | One game per channel | No game-ID needed in commands; simpler UX |
| **Order secrecy** | DM-based submission | `/order` commands go to the bot via DM; the player's DM thread is the per-player order log; the game channel only receives resolved state so opponents cannot read pending orders |
| **Scheduler** | `Scheduler` interface | Decouples deadline logic from execution model; `LocalScheduler` wraps `time.AfterFunc` for server deployments; `EventBridgeScheduler` creates one-time EventBridge Scheduler rules for Lambda |

---

## Full command set

| Category | Command | Phase | Who |
|---|---|---|---|
| Setup | `/newgame [settings]` | — | Anyone |
| Setup | `/join [country]` | — | Anyone |
| Setup | `/start` | — | GM |
| Movement | `/order <order-text>` | Movement | Own nation |
| Movement | `/orders` | Movement | Own nation |
| Movement | `/clear [order]` | Movement | Own nation |
| Movement | `/submit` | Movement | Own nation |
| Retreat | `/retreat <unit> <province>` | Retreat | Own nation |
| Retreat | `/disband <unit>` | Retreat | Own nation |
| Adjustment | `/build <unit-type> <province>` | Adjustment | Own nation |
| Adjustment | `/disband <unit>` | Adjustment | Own nation |
| Adjustment | `/waive` | Adjustment | Own nation |
| Info | `/map [territory [n]]` | Any | Anyone |
| Info | `/status` | Any | Anyone |
| Info | `/history <turn>` | Any | Anyone |
| Info | `/help [command]` | Any | Anyone |
| Draw | `/draw` | Any | Own nation |
| Draw | `/concede` | Any | Own nation |
| GM | `/pause` | Any | GM |
| GM | `/resume` | Any | GM |
| GM | `/extend <duration>` | Any | GM |
| GM | `/force-resolve` | Any | GM |
| GM | `/boot <nation>` | Any | GM |
| GM | `/replace <nation> <user>` | Any | GM |

`/map Vienna 1` shows Vienna and all adjacent territories; `/map Vienna 2` extends one hop
further. Implemented via BFS over godip's `Graph.Edges()` to radius `n`.

---

## Phase management

godip has 3 distinct phases. Retreat and Adjustment may be skipped entirely.

| Phase | Trigger | Commands |
|---|---|---|
| Movement | Start of Spring / Fall | `/order`, `/submit` |
| Retreat | After any dislodgement | `/retreat`, `/disband` |
| Adjustment | After Fall SC count | `/build`, `/disband`, `/waive` |

`Phase.DefaultOrder()` fills holds for NMR in Movement; unordered retreat units are
auto-disbanded by godip's `PostProcess`.

---

## Event types (stored as JSON)

Events are split between the shared game channel (visible to all players) and private
player DM threads (visible only to that player).

**Game channel events:**
```
GameCreated     {variant, deadline_hours, settings, gm_user_id}
PlayerJoined    {user_id, nation}
GameStarted     {initial_state: godip.Dump(), deadline_at: RFC3339}
PhaseResolved   {phase, state_snapshot: godip.Dump(), result_summary, deadline_at: RFC3339}
PhaseSkipped    {phase, reason: "no_dislodgements"|"no_sc_delta"}
NMRRecorded     {nation, phase, auto_orders}
DrawProposed    {proposer_nation}
DrawVoted       {nation, accept}
GameEnded       {result: "solo"|"draw"|"concession", winner, final_state}
```

**Player DM events** (private, one thread per player):
```
OrderSubmitted  {user_id, nation, orders, phase}
```

`deadline_at` is the absolute UTC time (RFC3339) at which the current phase resolves.
Any Lambda invocation can re-derive the deadline from the most recent `GameStarted` or
`PhaseResolved` event without carrying in-process timer state.

**State restoration on bot restart:**
1. Scan game channel history for last `PhaseResolved` or `GameStarted` event
2. `json.Unmarshal` snapshot → `state.Load()` — state restored
3. Scan forward for any `PhaseSkipped` / `NMRRecorded` events after the snapshot
4. For each player nation, read that player's DM thread for `OrderSubmitted` events for
   the current phase — reload as staged orders
5. Bot is ready to accept commands or advance phase

---

## Phase flow

```
GameStarted
  └─► Spring Movement (deadline T, deadline_at stored in event)
        ├─ players submit orders via DM
        ├─ T expires (or all nations submitted)
        ├─ NMR fill → adjudicate → PhaseResolved (deadline_at for next phase)
        └─► Spring Retreat
              ├─ if no dislodgements → PhaseSkipped
              ├─ players submit retreats via DM
              └─► Fall Movement (deadline T)
                    └─► Fall Retreat
                          └─► Adjustment
                                ├─ SC ownership updated
                                ├─ if any nation ≥18 SC → GameEnded (solo)
                                ├─ players submit builds/disbands via DM
                                └─► Spring Movement (year+1)
```

**Phase resolution triggers** — whichever fires first:

1. **All nations submit** — each DM invocation checks all player DM threads after staging an
   order; if all nations have submitted, calls `AdvanceTurn()` inline and calls
   `scheduler.Cancel(channelID)`.
2. **Deadline fires** — EventBridge (or `time.AfterFunc` on a server) triggers the handler
   with `{action: "advance_turn", channelID: "..."}`, which calls `session.Load()` then
   `AdvanceTurn()`.

`AdvanceTurn()` is **idempotent**: it checks for an existing `PhaseResolved` event for the
current phase before resolving, and no-ops if one is found. This guards against duplicate
invocations from concurrent Lambda instances or a race between trigger 1 and trigger 2.
