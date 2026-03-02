package engine

import (
	"fmt"

	"github.com/zond/godip"
	"github.com/zond/godip/variants/classical"
)

// orderParser parses a text order string for a given nation into a source
// province and an adjOrder ready to be staged on the game state.
type orderParser interface {
	Parse(nation godip.Nation, orderText string) (godip.Province, adjOrder, error)
}

// classicalOrderParser is the production parser for the Classical variant.
// It uses classical.DATCOrder to convert player order text (e.g. "A Vie-Bud")
// into a real godip.Adjudicator that can be staged for adjudication.
type classicalOrderParser struct{}

func (classicalOrderParser) Parse(_ godip.Nation, orderText string) (godip.Province, adjOrder, error) {
	prov, order, err := classical.DATCOrder(orderText)
	if err != nil {
		return "", nil, fmt.Errorf("invalid order %q: %w", orderText, err)
	}
	return prov, order, nil
}
