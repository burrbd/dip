// Package bot contains the platform-agnostic command router. It receives
// parsed commands from platform adapters, enforces access control (only the
// assigned player may submit orders for their nation), and delegates to the
// session and engine packages.
package bot
