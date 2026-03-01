# Bot Architecture

This document describes the design of the Diplomacy bot layer that sits on top of the
[godip](https://github.com/zond/godip) adjudication engine.

The bot is deployed as a webhook process for one or more chat platforms. Each chat channel
hosts exactly one game. The channel's message history is the audit trail — the bot posts a
structured JSON snapshot after each resolution so that state can be rebuilt on restart without
an external database.

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
  session.go         — Session struct: phase, staged orders, player map, deadline, GM user ID
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

## Event types (stored as JSON in channel)

```
GameCreated      {variant, deadline_hours, settings, gm_user_id}
PlayerJoined     {user_id, nation}
GameStarted      {initial_state: godip.Dump()}
OrderSubmitted   {user_id, nation, orders, phase}
PhaseResolved    {phase, state_snapshot: godip.Dump(), result_summary}
PhaseSkipped     {phase, reason: "no_dislodgements"|"no_sc_delta"}
NMRRecorded      {nation, phase, auto_orders}
DrawProposed     {proposer_nation}
DrawVoted        {nation, accept}
GameEnded        {result: "solo"|"draw"|"concession", winner, final_state}
```

**State restoration on bot restart:**
1. Scan channel history for last `PhaseResolved` or `GameStarted` event
2. `json.Unmarshal` snapshot → `state.Load()` — state restored
3. Scan forward for `OrderSubmitted` events after the snapshot — reload as staged orders
4. Bot is ready to accept commands or advance phase

---

## Phase flow

```
GameStarted
  └─► Spring Movement (deadline T)
        ├─ players submit orders
        ├─ T expires (or all submitted)
        ├─ NMR fill → adjudicate → PhaseResolved
        └─► Spring Retreat
              ├─ if no dislodgements → PhaseSkipped
              ├─ players submit retreats
              └─► Fall Movement (deadline T)
                    └─► Fall Retreat
                          └─► Adjustment
                                ├─ SC ownership updated
                                ├─ if any nation ≥18 SC → GameEnded (solo)
                                ├─ players submit builds/disbands
                                └─► Spring Movement (year+1)
```
