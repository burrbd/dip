# PLAN_QA.md — Local QA REPL Implementation Plan

## Overview

Add a standalone `go run ./cmd/qabot` entry point that provides a readline-style REPL for
exercising the complete bot layer without any external platform. State is held entirely in
memory; players are switched with a `/as <Nation|gm>` meta-command; phases are advanced
manually with `/force-resolve`. No new external dependencies are required.

---

## Codebase Observations Relevant to Implementation

**`events.Channel` interface (5 methods in current codebase):**
```go
Post(channelID, text string) error
History(channelID string) ([]string, error)
SendDM(userID, text string) error
DMHistory(userID string) ([]string, error)
PostImage(channelID string, data []byte) error
```

**`bot.Dispatcher` constructor:**
```go
bot.New(ch events.Channel, notifier session.Notifier, loader session.EngineLoader, newEng bot.EngineFactory) *bot.Dispatcher
```

**`session.Notifier` interface:**
```go
Notify(channelID, message string) error
```

**`session.EngineLoader`** is `func(snapshot []byte) (engine.Engine, error)` — use `engine.Load`.

**`bot.EngineFactory`** is `func(variant string) (engine.Engine, error)` — use `engine.New`.

**No `Scheduler` interface exists yet** (planned for Story 13). The session's deadline timer is
controlled by `session.New(..., deadlineHours int, ...)`. When `deadlineHours <= 0`,
`startDeadline()` returns immediately. However, `handleNewGame` hardcodes `DeadlineHours: 24`
in the `GameCreated` event. The resulting 24-hour timer fires into the in-memory channel and is
harmless in practice (see Edge Cases).

**`memChannel` in `bot/bot_functional_test.go`** is behind `//go:build functional`. The local
platform package must implement its own equivalent — it cannot import from `bot_test`.

**DM command routing** (established pattern from functional tests):
- DM commands: `IsDM=true`, `ChannelID="dm_"+userID`, `GameChannelID=gameChannelID`
- Channel commands: `IsDM=false`, `ChannelID=gameChannelID`

**DM commands** (must be dispatched with `IsDM=true`):
`order`, `orders`, `clear`, `submit`, `retreat`, `disband`, `build`, `waive`

**Channel commands** (dispatched with `IsDM=false`): everything else.

---

## Files to Create

### 1. `platform/local/channel.go`

Implements `events.Channel` for in-memory use. Thread-safe. Tracks message and image cursors
so the REPL can print only new output after each command dispatch.

```go
package local

import "sync"

type Channel struct {
    mu   sync.Mutex
    msgs map[string][]string // channelID → messages
    dms  map[string][]string // userID → DM messages
    imgs map[string][][]byte // channelID → image byte slices
}

func NewChannel() *Channel

// events.Channel interface:
func (c *Channel) Post(channelID, text string) error
func (c *Channel) History(channelID string) ([]string, error)
func (c *Channel) SendDM(userID, text string) error
func (c *Channel) DMHistory(userID string) ([]string, error)
func (c *Channel) PostImage(channelID string, data []byte) error

// REPL cursor helpers (not part of the interface):
func (c *Channel) MessagesSince(channelID string, cursor int) []string
func (c *Channel) MessageCount(channelID string) int
func (c *Channel) DMsSince(userID string, cursor int) []string
func (c *Channel) DMCount(userID string) int
func (c *Channel) ImagesSince(channelID string, cursor int) [][]byte
func (c *Channel) ImageCount(channelID string) int
```

All fields protected by `c.mu`. Cursor methods return `slice[cursor:]` copies.

### 2. `platform/local/channel_test.go`

Unit tests (no build tag — included in `go test ./...`):
- `Post` / `History` round-trip
- `SendDM` / `DMHistory` round-trip
- `PostImage` stores bytes retrievable via `ImagesSince`
- `MessagesSince` / `DMsSince` / `ImagesSince` return only entries after the cursor
- Concurrent `Post` calls do not race (run with `-race`)

### 3. `platform/local/notifier.go`

A thin `session.Notifier` that posts notifications to the in-memory channel.

```go
package local

import "github.com/burrbd/dip/events"

type Notifier struct{ ch *Channel }

func NewNotifier(ch *Channel) *Notifier

// Notify posts "[notify] <message>" to channelID via ch.Post.
func (n *Notifier) Notify(channelID, message string) error
```

Add a small `notifier_test.go` verifying that `Notify` posts to the channel.

### 4. `cmd/qabot/main.go`

The REPL entry point. ~150 lines.

**Key constants:**
```go
const gameChannelID = "local-game"

var dmCommands = map[string]bool{
    "order": true, "orders": true, "clear": true, "submit": true,
    "retreat": true, "disband": true, "build": true, "waive": true,
}
```

**`main()` flow:**
1. Construct `ch := local.NewChannel()`
2. Construct `notifier := local.NewNotifier(ch)`
3. Construct `d := bot.New(ch, notifier, engine.Load, engine.New)`
4. Set `activeUser := "gm"`
5. Enter `bufio.Scanner` readline loop

**Loop body:**
```
line = strings.TrimSpace(scanner.Text())
if line == "" { continue }
if !strings.HasPrefix(line, "/") { print usage hint; continue }

tokens = strings.Fields(line)
cmdName = strings.TrimPrefix(tokens[0], "/")
args = tokens[1:]

if cmdName == "as" { handle /as locally; continue }

msgCursor = ch.MessageCount(gameChannelID)
dmCursor  = ch.DMCount(activeUser)
imgCursor = ch.ImageCount(gameChannelID)

cmd = buildCommand(cmdName, args, activeUser, gameChannelID)
resp, err = d.Dispatch(cmd)

if err != nil  { fmt.Printf("Error: %v\n", err) }
else if resp != "" { fmt.Println(resp) }

for _, msg := range ch.MessagesSince(gameChannelID, msgCursor) {
    fmt.Printf("[channel] %s\n", msg)
}
for _, dm := range ch.DMsSince(activeUser, dmCursor) {
    fmt.Printf("[dm:%s] %s\n", activeUser, dm)
}
for _, imgBytes := range ch.ImagesSince(gameChannelID, imgCursor) {
    f, _ := os.CreateTemp("", "dip-map-*.png")
    f.Write(imgBytes); f.Close()
    fmt.Printf("Map saved to %s\n", f.Name())
}
```

**`/as` handler:**
```
if len(args) == 0 { print current player; continue }
activeUser = strings.ToLower(args[0])
fmt.Printf("Now acting as: %s\n", activeUser)
```

**`buildCommand` helper:**
```go
func buildCommand(name string, args []string, activeUser, gameChannelID string) bot.Command {
    if dmCommands[name] {
        return bot.Command{
            Name: name, Args: args, UserID: activeUser,
            ChannelID: "dm_" + activeUser, IsDM: true,
            GameChannelID: gameChannelID,
        }
    }
    return bot.Command{
        Name: name, Args: args, UserID: activeUser,
        ChannelID: gameChannelID,
    }
}
```

**Prompt:** Print `[<displayName(activeUser)>] > ` before each `Scan()` call.
`displayName` capitalises the first letter (`"england"` → `"England"`, `"gm"` → `"gm"`).

---

## Player / UserID Convention

| `/as` argument | `activeUser` sent to bot |
|---|---|
| `gm` | `"gm"` |
| `England` | `"england"` |
| `France` | `"france"` |
| `Austria` | `"austria"` |
| *(any nation)* | lowercase of the name |

`/join England` dispatched with `UserID="england"` → bot stores `PlayerJoined{UserID:"england", Nation:"England"}`. Subsequent `/order` DMs from `UserID="england"` resolve to nation England. The GM's `UserID="gm"` matches the `GMUserID` in the `GameCreated` event.

---

## Example Session

```
[gm] > /newgame
Game created. You are the GM. ...
[channel] {"type":"GameCreated",...}

[gm] > /as England
Now acting as: England

[England] > /join England
Joined as England.
[channel] {"type":"PlayerJoined",...}

[England] > /as France
Now acting as: France

[France] > /join France
...

[France] > /as gm
Now acting as: gm

[gm] > /start
Game started! Spring 1901 Movement.
[channel] {"type":"GameStarted",...}

[gm] > /as England
Now acting as: England

[England] > /order A Lon H
Order staged: A Lon H

[England] > /as gm
Now acting as: gm

[gm] > /force-resolve
Phase resolved.
[channel] {"type":"PhaseResolved",...}

[gm] > /map
Map posted.
Map saved to /tmp/dip-map-1234567890.png

[gm] > /status
Phase: Fall 1901 Movement ...
```

---

## Acceptance Criteria

1. `go build ./cmd/qabot` succeeds with no errors.
2. `go run ./cmd/qabot` starts with prompt `[gm] > `.
3. Full game flow works end-to-end: `/newgame` → `/join` (×2–7 nations) → `/start` →
   `/order` → `/force-resolve` → `/map`.
4. `/as` switches produce the correct authorization context:
   - `/order` DM from `"england"` is accepted for nation England.
   - `/order` DM from `"gm"` returns a "not a player" error.
   - `/force-resolve` from `"england"` returns a "GM only" error.
5. `/map` writes a valid PNG (`\x89PNG` header) to a temp file and prints its path.
6. `go test -v -cover -race ./...` still passes — new `platform/local` tests pass;
   no existing tests broken.
7. No modifications to any existing package (`engine/`, `events/`, `session/`, `bot/`,
   `dipmap/`, `platform/slack/`, `platform/telegram/`).
8. No new external dependencies — only stdlib + existing project packages.

---

## Edge Cases and Notes

**24-hour deadline timer:** `handleNewGame` hardcodes `DeadlineHours: 24`. A 24-hour
`time.AfterFunc` goroutine is started on `/start`. In a typical QA session this never fires.
If it does, it calls `AdvanceTurn()` on the in-memory channel — the channel never blocks and
the only visible effect is extra messages printed on the next REPL command. Documented in a
code comment in `main.go`. When Story 13 introduces the `Scheduler` interface, a
`NoOpScheduler` can be added to `platform/local/` to eliminate this entirely.

**Ctrl+D / EOF:** `scanner.Scan()` returns `false`; loop exits cleanly. No special signal
handling needed for a QA tool.

**Concurrent access:** The 24-hour timer callback and the REPL main goroutine run
concurrently. `Channel.mu` protects all slice mutations, so the race detector is satisfied.

**`/history` before first resolution:** Returns a bot error ("no history found"); REPL prints
the error and continues — correct behaviour.

**Image cursor:** Each `/map` call appends one PNG to `ch.imgs[gameChannelID]`. The cursor
ensures only the image from the current dispatch is saved to a new temp file.

**Story 13 forward-compatibility:** The `NoOpScheduler` is not needed yet. When Story 13
lands, add `platform/local/scheduler.go` implementing the new `Scheduler` interface with
empty `Schedule`/`Cancel` methods and wire it into `main.go`.

---

## Implementation Sequence

1. `platform/local/channel.go` + `platform/local/channel_test.go`
2. `platform/local/notifier.go` + `platform/local/notifier_test.go`
3. `cmd/qabot/main.go`
4. Run `go test -v -cover -race ./...` — must pass
5. Smoke-test with `go run ./cmd/qabot` interactively
