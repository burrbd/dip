// Package events defines the structured event types that are serialised as JSON
// and posted to the chat channel after each game action. The channel's message
// history is the authoritative event log; state is restored by replaying events
// from the last PhaseResolved snapshot.
package events
