// Package engine wraps the godip adjudication engine, providing a simplified
// interface for the bot layer. It hides godip's internal state types and
// exposes only the operations needed to run a Diplomacy game turn.
package engine

import (
	"fmt"

	"github.com/zond/godip"
	"github.com/zond/godip/variants"
	"github.com/zond/godip/variants/classical"
)

// Load restores an Engine from a JSON snapshot produced by Dump.
func Load(snapshot []byte) (Engine, error) {
	return loadFromSnapshot(snapshot, classical.Load)
}

// loadFromSnapshot restores an Engine using loader to deserialise the snapshot.
// Separated from Load so tests can inject a failing loader.
func loadFromSnapshot(snapshot []byte, loader func([]byte) (godip.Adjudicator, error)) (Engine, error) {
	adj, err := loader(snapshot)
	if err != nil {
		return nil, fmt.Errorf("engine: load snapshot: %w", err)
	}
	return &game{adj: adj, parser: classical.Parser}, nil
}

// Engine is the public interface for interacting with a running Diplomacy game.
// All methods operate on the current game phase.
type Engine interface {
	// SubmitOrder parses orderText and stages it for nation's unit.
	SubmitOrder(nation, orderText string) error
	// Resolve adjudicates all staged orders and returns a summary of outcomes.
	Resolve() (ResolutionResult, error)
	// Advance fills any missing (NMR) orders, calls godip Next(), and skips
	// empty retreat or adjustment phases.
	Advance() error
	// SoloWinner returns the nation that has achieved a solo victory, or ""
	// if no solo winner exists yet.
	SoloWinner() string
	// Dump serialises the current game state to JSON for event-log storage.
	Dump() ([]byte, error)
	// Phase returns the current game phase as a human-readable string
	// (e.g. "Spring 1901 Movement").
	Phase() string
}

// ResolutionResult summarises what happened when a phase was adjudicated.
type ResolutionResult struct {
	Phase  string
	Year   int
	Orders []OrderResult
}

// OrderResult represents the outcome of a single order after adjudication.
type OrderResult struct {
	Province string
	Order    string
	Success  bool
}

// game implements Engine around a godip Adjudicator.
type game struct {
	adj    godip.Adjudicator
	parser orderParser
}

// New creates an Engine for the named Diplomacy variant (currently "classical").
func New(variant string) (Engine, error) {
	v, err := lookupVariant(variant)
	if err != nil {
		return nil, err
	}
	return newFromVariant(v, variant, classical.Parser)
}

// newFromVariant starts a game from the given Variant and parser. Separated
// from New so that tests can inject a variant with a failing Start function.
func newFromVariant(v variants.Variant, name string, p orderParser) (Engine, error) {
	adj, err := v.Start()
	if err != nil {
		return nil, fmt.Errorf("engine: start %s: %w", name, err)
	}
	return &game{adj: adj, parser: p}, nil
}

// lookupVariant returns the Variant for the given name.
func lookupVariant(name string) (variants.Variant, error) {
	switch name {
	case "classical":
		return classical.ClassicalVariant, nil
	default:
		return variants.Variant{}, fmt.Errorf("engine: unknown variant %q", name)
	}
}

// SubmitOrder parses the order text and stages it on the adjudicator.
func (g *game) SubmitOrder(nation, orderText string) error {
	prov, order, err := g.parser.Parse(godip.Nation(nation), orderText)
	if err != nil {
		return fmt.Errorf("engine: parse order: %w", err)
	}
	g.adj.SetOrder(prov, order)
	return nil
}

// Resolve adjudicates all staged orders and returns a summary of outcomes.
func (g *game) Resolve() (ResolutionResult, error) {
	phase := g.adj.Phase()
	result := ResolutionResult{
		Phase: string(phase.Type()),
		Year:  phase.Year(),
	}

	for prov, ord := range g.adj.Orders() {
		err := g.adj.Resolve(prov)
		result.Orders = append(result.Orders, OrderResult{
			Province: string(prov),
			Order:    string(ord.Type()),
			Success:  err == nil,
		})
	}
	return result, nil
}

// Dump serialises the current game state to JSON.
func (g *game) Dump() ([]byte, error) {
	return g.adj.Dump()
}

// Phase returns the current game phase as a human-readable string,
// e.g. "Spring 1901 Movement". Returns "" if the phase is nil.
func (g *game) Phase() string {
	phase := g.adj.Phase()
	if phase == nil {
		return ""
	}
	return fmt.Sprintf("%s %d %s", phase.Season(), phase.Year(), phase.Type())
}
