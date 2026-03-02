package bot

import (
	"fmt"
	"sort"

	"github.com/burrbd/dip/session"
)

// Autocomplete returns suggested order strings for the given nation's units
// in the current phase. Each suggestion is a hold order for one of the
// nation's units, formatted as "<unit-type-prefix> <province> H".
// Returns nil when sess is nil or the nation has no units on the board.
func Autocomplete(sess *session.Session, nation string) []string {
	if sess == nil || sess.Eng == nil {
		return nil
	}
	units := sess.Eng.Units()
	var suggestions []string
	for prov, info := range units {
		if info.Nation != nation {
			continue
		}
		prefix := "A"
		if info.Type == "Fleet" {
			prefix = "F"
		}
		suggestions = append(suggestions, fmt.Sprintf("%s %s H", prefix, prov))
	}
	sort.Strings(suggestions)
	return suggestions
}
