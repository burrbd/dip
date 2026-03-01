// Package session manages the lifecycle of a single Diplomacy game session
// within a chat channel. It owns the deadline timer, staged orders, player→nation
// mapping, GM identity, and phase transitions.
package session
