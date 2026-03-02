package bot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/burrbd/dip/engine"
)

// FormatResult renders a ResolutionResult as a human-readable plain-text string.
func FormatResult(r engine.ResolutionResult) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Phase %s resolved. %d orders processed.\n", r.Phase, len(r.Orders))
	for _, o := range r.Orders {
		status := "succeeded"
		if !o.Success {
			status = "failed"
		}
		fmt.Fprintf(&sb, "  %s %s: %s\n", o.Province, o.Order, status)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// FormatStatus renders the current game status as plain text, listing each
// nation and whether their orders have been submitted for the current phase.
func FormatStatus(phase string, players map[string]string, submitted map[string]bool) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Phase: %s\n", phase)

	// Collect unique nation names and sort for deterministic output.
	seen := make(map[string]bool)
	for _, nation := range players {
		seen[nation] = true
	}
	nations := make([]string, 0, len(seen))
	for n := range seen {
		nations = append(nations, n)
	}
	sort.Strings(nations)

	for _, n := range nations {
		sub := "pending"
		if submitted[n] {
			sub = "submitted"
		}
		fmt.Fprintf(&sb, "  %s: %s\n", n, sub)
	}
	return strings.TrimRight(sb.String(), "\n")
}
