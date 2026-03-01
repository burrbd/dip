package engine

import "github.com/zond/godip"

// orderParser parses a text order string for a given nation into a source
// province and a godip Order ready to be staged on the adjudicator.
type orderParser interface {
	Parse(nation godip.Nation, orderText string) (godip.Province, godip.Order, error)
}
