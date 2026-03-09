// Command slackbot is the entry point for the Slack bot. It starts an HTTP
// server that handles Slack slash commands and Events API webhooks.
//
// Environment variables:
//
//	SLACK_BOT_TOKEN      — required; Slack bot OAuth token (xoxb-...)
//	SLACK_SIGNING_SECRET — required; Slack app signing secret for request verification
//	DATA_DIR             — directory for the JSONL history store (default: ./data)
//	PORT                 — HTTP listen port (default: 8080)
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/burrbd/dip/bot"
	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/platform/slack"
	"github.com/burrbd/dip/session"
)

func main() {
	token := mustEnv("SLACK_BOT_TOKEN")
	signingSecret := mustEnv("SLACK_SIGNING_SECRET")
	dataDir := envOr("DATA_DIR", "./data")
	port := envOr("PORT", "8080")

	ch, err := slack.New(token, dataDir)
	if err != nil {
		log.Fatalf("slackbot: create channel: %v", err)
	}

	notifier := slack.NewNotifier(ch)
	d := bot.New(ch, notifier, engine.Load, engine.New)

	http.HandleFunc("/slash", makeSlashHandler(ch, d, signingSecret))
	http.HandleFunc("/events", makeEventsHandler(signingSecret))

	log.Printf("slackbot: listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("slackbot: server: %v", err)
	}
}

// makeSlashHandler returns an http.HandlerFunc that processes Slack slash
// command requests. Requests are verified using the signing secret before
// being dispatched.
func makeSlashHandler(ch *slack.Channel, d *bot.Dispatcher, signingSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("slackbot: read body: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if !slack.VerifyRequest(signingSecret, r.Header.Get("X-Slack-Request-Timestamp"), body, r.Header.Get("X-Slack-Signature")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		cmd, ok := ch.ParseSlashCommand(body)
		if !ok {
			return
		}
		resp, err := d.Dispatch(cmd)
		if err != nil {
			log.Printf("slackbot: dispatch %q: %v", cmd.Name, err)
			if postErr := ch.Post(cmd.ChannelID, "Error: "+err.Error()); postErr != nil {
				log.Printf("slackbot: post error response: %v", postErr)
			}
			return
		}
		if resp != "" {
			if postErr := ch.Post(cmd.ChannelID, resp); postErr != nil {
				log.Printf("slackbot: post response: %v", postErr)
			}
		}
	}
}

// urlVerification is the Slack Events API URL verification challenge payload.
type urlVerification struct {
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
}

// makeEventsHandler returns an http.HandlerFunc that handles Slack Events API
// payloads. Currently it handles URL verification challenges only.
func makeEventsHandler(signingSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("slackbot: events read body: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if !slack.VerifyRequest(signingSecret, r.Header.Get("X-Slack-Request-Timestamp"), body, r.Header.Get("X-Slack-Signature")) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var uv urlVerification
		if err := json.Unmarshal(body, &uv); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if uv.Type == "url_verification" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"challenge": uv.Challenge}) //nolint:errcheck
			return
		}
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("slackbot: required env var %s is not set", key)
	}
	return v
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// Notifier is satisfied by slack.Notifier.
var _ session.Notifier = (*slack.Notifier)(nil)
