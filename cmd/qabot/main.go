// Command qabot is a standalone QA REPL for exercising the Diplomacy bot layer
// without any external platform. State is held entirely in memory.
//
// Usage:
//
//	go run ./cmd/qabot
//
// Meta-commands:
//
//	/as <Nation|gm>   — switch the active player
//
// All other commands are dispatched to the bot dispatcher as normal bot.Commands.
// DM commands (order, orders, clear, submit, retreat, disband, build, waive) are
// automatically routed with IsDM=true.
//
// Note on the 24-hour deadline timer: when /start is called, a 24-hour time.AfterFunc
// is started. In a typical QA session this never fires, but if it does it simply calls
// AdvanceTurn() on the in-memory channel — harmless extra messages will appear on the
// next REPL command. When Story 13 introduces the Scheduler interface, a NoOpScheduler
// can be added to platform/local/ to eliminate this entirely.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/burrbd/dip/bot"
	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/platform/local"
)

const gameChannelID = "local-game"

// dmCommands lists commands that must be dispatched as DMs (IsDM=true).
var dmCommands = map[string]bool{
	"order":   true,
	"orders":  true,
	"clear":   true,
	"submit":  true,
	"retreat": true,
	"disband": true,
	"build":   true,
	"waive":   true,
}

func main() {
	ch := local.NewChannel()
	notifier := local.NewNotifier(ch)
	d := bot.New(ch, notifier, engine.Load, engine.New)

	activeUser := "gm"
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("[%s] > ", displayName(activeUser))
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "/") {
			fmt.Println("Commands must start with /. Try /help or /as <Nation|gm>.")
			continue
		}

		tokens := strings.Fields(line)
		cmdName := strings.TrimPrefix(tokens[0], "/")
		args := tokens[1:]

		if cmdName == "as" {
			if len(args) == 0 {
				fmt.Printf("Current player: %s\n", displayName(activeUser))
			} else {
				activeUser = strings.ToLower(args[0])
				fmt.Printf("Now acting as: %s\n", displayName(activeUser))
			}
			continue
		}

		msgCursor := ch.MessageCount(gameChannelID)
		dmCursor := ch.DMCount(activeUser)
		imgCursor := ch.ImageCount(gameChannelID)

		cmd := buildCommand(cmdName, args, activeUser, gameChannelID)
		resp, err := d.Dispatch(cmd)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else if resp != "" {
			fmt.Println(resp)
		}

		for _, msg := range ch.MessagesSince(gameChannelID, msgCursor) {
			fmt.Printf("[channel] %s\n", msg)
		}
		for _, dm := range ch.DMsSince(activeUser, dmCursor) {
			fmt.Printf("[dm:%s] %s\n", activeUser, dm)
		}
		for _, imgBytes := range ch.ImagesSince(gameChannelID, imgCursor) {
			f, err := os.CreateTemp("", "dip-map-*.png")
			if err != nil {
				fmt.Printf("Error saving map: %v\n", err)
				continue
			}
			f.Write(imgBytes)
			f.Close()
			fmt.Printf("Map saved to %s\n", f.Name())
		}
	}
}

// buildCommand constructs a bot.Command from the parsed input. DM commands are
// routed with IsDM=true and ChannelID set to the player's DM channel.
func buildCommand(name string, args []string, activeUser, gameChannelID string) bot.Command {
	if dmCommands[name] {
		return bot.Command{
			Name:          name,
			Args:          args,
			UserID:        activeUser,
			ChannelID:     "dm_" + activeUser,
			IsDM:          true,
			GameChannelID: gameChannelID,
		}
	}
	return bot.Command{
		Name:      name,
		Args:      args,
		UserID:    activeUser,
		ChannelID: gameChannelID,
	}
}

// displayName capitalises the first letter of the user identifier for display.
// "england" → "England", "gm" → "gm".
func displayName(user string) string {
	if len(user) == 0 {
		return user
	}
	return strings.ToUpper(user[:1]) + user[1:]
}
