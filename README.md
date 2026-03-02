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

## Running tests

```bash
go test -v -cover -race ./...
```
