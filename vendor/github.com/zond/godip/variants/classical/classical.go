// Package classical provides the standard Classical Diplomacy variant and its
// order parser.
package classical

import (
	"fmt"
	"strings"

	"github.com/zond/godip"
	"github.com/zond/godip/variants"
)

// ClassicalVariant is the standard seven-player Diplomacy variant.
var ClassicalVariant = variants.Variant{
	Name:  "Classical",
	Start: startClassical,
}

// Parser is the order parser for the Classical variant.
var Parser OrderParser = classicalParser{}

// OrderParser parses a text order string for a given nation into a province
// and godip Order.
type OrderParser interface {
	Parse(nation godip.Nation, orderText string) (godip.Province, godip.Order, error)
}

// classicalParser is the concrete parser for Classical Diplomacy orders.
type classicalParser struct{}

// Parse parses a Classical Diplomacy order string such as "A Vie-Bud" or
// "F Nth C A Yor-Bre". It returns the source province and the corresponding
// godip Order.
func (classicalParser) Parse(nation godip.Nation, orderText string) (godip.Province, godip.Order, error) {
	parts := strings.Fields(orderText)
	if len(parts) < 2 {
		return "", nil, fmt.Errorf("invalid order %q: too few tokens", orderText)
	}
	// First token is unit type (A/F), second is province.
	src := godip.Province(parts[1])
	order := &parsedOrder{orderText: orderText}
	return src, order, nil
}

// parsedOrder is a minimal Order implementation returned by the stub parser.
type parsedOrder struct {
	orderText string
}

func (o *parsedOrder) Type() godip.OrderType      { return godip.OrderType(o.orderText) }
func (o *parsedOrder) Flags() map[godip.Flag]bool { return nil }

// startClassical returns a new Classical game state. This is a stub
// implementation; the real godip library provides a fully-adjudicating state.
func startClassical() (godip.Adjudicator, error) {
	return &classicalState{
		phase:  &classicalPhase{typ: godip.Movement, year: 1901, season: godip.Spring},
		orders: make(map[godip.Province]godip.Order),
		units:  make(map[godip.Province]godip.Unit),
	}, nil
}

// classicalState is a stub Adjudicator for the Classical variant.
type classicalState struct {
	phase      godip.Phase
	orders     map[godip.Province]godip.Order
	units      map[godip.Province]godip.Unit
	dislodgeds map[godip.Province]godip.Unit
	winner     godip.Nation
}

func (s *classicalState) Phase() godip.Phase                        { return s.phase }
func (s *classicalState) Orders() map[godip.Province]godip.Order    { return s.orders }
func (s *classicalState) Units() map[godip.Province]godip.Unit      { return s.units }
func (s *classicalState) Dislodgeds() map[godip.Province]godip.Unit { return s.dislodgeds }
func (s *classicalState) SetOrder(p godip.Province, o godip.Order)  { s.orders[p] = o }
func (s *classicalState) Resolve(_ godip.Province) error            { return nil }
func (s *classicalState) SoloWinner() godip.Nation                  { return s.winner }

func (s *classicalState) Next() (godip.Adjudicator, error) {
	next := &classicalState{
		phase:  s.phase,
		orders: make(map[godip.Province]godip.Order),
		units:  make(map[godip.Province]godip.Unit),
		winner: s.winner,
	}
	return next, nil
}

func (s *classicalState) Dump() ([]byte, error) {
	return []byte(`{}`), nil
}

// classicalPhase is a stub Phase implementation.
type classicalPhase struct {
	typ    godip.PhaseType
	year   int
	season godip.Season
}

func (p *classicalPhase) Type() godip.PhaseType           { return p.typ }
func (p *classicalPhase) Year() int                       { return p.year }
func (p *classicalPhase) Season() godip.Season            { return p.season }
func (p *classicalPhase) DefaultOrder(_ godip.Province) godip.Order {
	return &parsedOrder{orderText: "Hold"}
}
