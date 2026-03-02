package engine

import (
	"fmt"
	"strings"

	"github.com/zond/godip"
)

// orderParser parses a text order string for a given nation into a source
// province and an adjOrder ready to be staged on the game state.
type orderParser interface {
	Parse(nation godip.Nation, orderText string) (godip.Province, adjOrder, error)
}

// classicalOrderParser is the production parser for the Classical variant.
// It uses a simple tokeniser that extracts the source province from the order
// text (e.g. "A Vie-Bud" → province "Vie").  Full order semantics are
// delegated to the real godip adjudicator via the stateWrapper.
type classicalOrderParser struct{}

func (classicalOrderParser) Parse(_ godip.Nation, orderText string) (godip.Province, adjOrder, error) {
	parts := strings.Fields(orderText)
	if len(parts) < 2 {
		return "", nil, fmt.Errorf("invalid order %q: too few tokens", orderText)
	}
	// First token is unit type (A/F), second is province.
	src := godip.Province(parts[1])
	return src, &parsedOrder{orderText: orderText}, nil
}

// parsedOrder is a minimal adjOrder returned by classicalOrderParser.
type parsedOrder struct{ orderText string }

func (o *parsedOrder) Type() godip.OrderType { return godip.OrderType(o.orderText) }
