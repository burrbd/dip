[![CircleCI](https://circleci.com/gh/burrbd/dip/tree/master.svg?style=svg)](https://circleci.com/gh/burrbd/dip/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/burrbd/dip)](https://goreportcard.com/report/github.com/burrbd/dip)

# dip

A Diplomacy messenger bot for Slack and Telegram. Players submit moves via slash commands (e.g. `/order A Vie-Bud`), view the board map, and see move history. Each channel hosts one game.

Adjudication is handled by [godip](https://github.com/zond/godip) (DATC-compliant). Game state is persisted as a JSON event log in the channel — no external database required.

## Packages

| Package | Role |
|---|---|
| `engine/` | godip wrapper — order parsing, phase advance, win detection |
| `events/` | structured JSON event log — write, scan, replay |
| `session/` | game lifecycle — turns, deadlines, serialization |
| `bot/` | platform-agnostic command router and formatter |
| `dipmap/` | SVG → PNG map rendering with province highlighting |
| `platform/` | Slack and Telegram adapters |
| `cmd/` | entry points (`slackbot`, `telegrambot`) |

## Running the QA bot

The QA bot is a local REPL for exercising the full bot layer in your terminal without needing
Slack, Telegram, or any external service. It uses the same event-log architecture as the real
bot: every command posts structured JSON events to an in-memory channel, and state is rebuilt
from that event log on each dispatch — exactly as it would be against Slack or Telegram. The
only difference from production is that the channel history is not persisted to disk, so
restarting the REPL starts a fresh game.

```bash
go run ./cmd/qabot
```

This starts a prompt: `[gm] > `

### Switching players

Use the `/as` meta-command to switch the active player. All subsequent commands are sent as
that player.

```
/as gm          — act as the Game Master (default)
/as England     — act as the England player
/as France      — act as France
```

### Command routing

Commands that involve secret order submission are automatically sent as DMs (so opponents
cannot see them): `order`, `orders`, `clear`, `submit`, `retreat`, `disband`, `build`, `waive`.

All other commands (`newgame`, `join`, `start`, `status`, `map`, `force-resolve`, etc.) are
sent to the shared game channel.

### Example session

```
[gm] > /newgame
[gm] > /as England
[England] > /join England
[England] > /as France
[France] > /join France
[France] > /as gm
[gm] > /start
[gm] > /as England
[England] > /order A Lon H
[England] > /submit
[England] > /as gm
[gm] > /force-resolve
[gm] > /map
Map saved to /tmp/dip-map-3456789012.jpg
[gm] > /status
```

Map images are written to a temp file and the path is printed. Exit with Ctrl-D.

## Running tests

```bash
go test -v -cover -race ./...
```
