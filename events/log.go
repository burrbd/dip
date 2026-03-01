package events

import (
	"encoding/json"
	"fmt"
)

// Channel is the platform-agnostic interface for posting and reading messages
// from a chat channel. Slack and Telegram adapters implement this interface.
type Channel interface {
	// Post appends a text message to the channel identified by channelID.
	Post(channelID, text string) error
	// History returns all messages in the channel in chronological order.
	History(channelID string) ([]string, error)
}

// Write serialises payload as a JSON Envelope and posts it to channelID.
func Write(ch Channel, channelID string, eventType EventType, payload any) error {
	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("events: marshal payload: %w", err)
	}
	env := Envelope{Type: eventType, Payload: json.RawMessage(raw)}
	// Envelope contains only string and json.RawMessage fields; Marshal cannot fail.
	data, _ := json.Marshal(env)
	return ch.Post(channelID, string(data))
}

// Scan reads the channel history and returns every message that can be parsed
// as a valid Envelope, in chronological order. Messages that are not valid
// Envelopes (plain chat text, etc.) are silently skipped.
func Scan(ch Channel, channelID string) ([]Envelope, error) {
	messages, err := ch.History(channelID)
	if err != nil {
		return nil, fmt.Errorf("events: scan history: %w", err)
	}
	var envs []Envelope
	for _, msg := range messages {
		var env Envelope
		if err := json.Unmarshal([]byte(msg), &env); err != nil {
			continue
		}
		if env.Type == "" {
			continue
		}
		envs = append(envs, env)
	}
	return envs, nil
}
