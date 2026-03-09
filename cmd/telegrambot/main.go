// Command telegrambot is the entry point for the Telegram bot. It starts an
// HTTP server that receives Telegram Bot API webhook updates, parses them into
// bot commands, and posts responses back to the originating chat.
//
// Environment variables:
//
//	TELEGRAM_BOT_TOKEN  — required; Telegram Bot API token
//	DATA_DIR            — directory for the JSONL history store (default: ./data)
//	PORT                — HTTP listen port (default: 8080)
package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/burrbd/dip/bot"
	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/platform/telegram"
	"github.com/burrbd/dip/session"
)

func main() {
	token := mustEnv("TELEGRAM_BOT_TOKEN")
	dataDir := envOr("DATA_DIR", "./data")
	port := envOr("PORT", "8080")

	store, err := telegram.NewStore(dataDir)
	if err != nil {
		log.Fatalf("telegrambot: create store: %v", err)
	}

	ch := telegram.New(token, store)
	notifier := telegram.NewNotifier(ch)
	d := bot.New(ch, notifier, engine.Load, engine.New)

	http.HandleFunc("/webhook", makeWebhookHandler(ch, d))

	log.Printf("telegrambot: listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("telegrambot: server: %v", err)
	}
}

// makeWebhookHandler returns an http.HandlerFunc that processes Telegram
// webhook updates. It is extracted for testability.
func makeWebhookHandler(ch *telegram.Channel, d *bot.Dispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("telegrambot: read body: %v", err)
			return
		}
		cmd, ok := ch.ParseUpdate(body)
		if !ok {
			return
		}
		resp, err := d.Dispatch(cmd)
		if err != nil {
			log.Printf("telegrambot: dispatch %q: %v", cmd.Name, err)
			if postErr := ch.Post(cmd.ChannelID, "Error: "+err.Error()); postErr != nil {
				log.Printf("telegrambot: post error response: %v", postErr)
			}
			return
		}
		if resp != "" {
			if postErr := ch.Post(cmd.ChannelID, resp); postErr != nil {
				log.Printf("telegrambot: post response: %v", postErr)
			}
		}
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("telegrambot: required env var %s is not set", key)
	}
	return v
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Notifier is satisfied by telegram.Notifier.
var _ session.Notifier = (*telegram.Notifier)(nil)
