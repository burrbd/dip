// Package engine wraps the godip adjudication engine, providing a simplified
// interface for the bot layer. It hides godip's internal state types and
// exposes only the operations needed to run a Diplomacy game turn.
package engine

import (
	"encoding/json"
	"fmt"

	"github.com/zond/godip"
	"github.com/zond/godip/state"
	"github.com/zond/godip/variants/classical"
	"github.com/zond/godip/variants/common"
)

// adjOrder is the minimal interface our engine needs from an order.
// Satisfied by real godip.Adjudicator and by stub orders used in tests.
type adjOrder interface {
	Type() godip.OrderType
}

// gamePhase is the engine-internal phase interface.
type gamePhase interface {
	Type() godip.PhaseType
	Year() int
	Season() godip.Season
	DefaultOrder(godip.Province) adjOrder
}

// gameState is the engine-internal game state interface.
// Satisfied by stateWrapper (real *state.State) and by test mocks.
type gameState interface {
	Phase() gamePhase
	Orders() map[godip.Province]adjOrder
	Units() map[godip.Province]godip.Unit
	Dislodgeds() map[godip.Province]godip.Unit
	SupplyCenters() map[godip.Province]godip.Nation
	SetOrder(godip.Province, adjOrder)
	Resolve(godip.Province) error
	Next() (gameState, error)
	SoloWinner() godip.Nation
	Dump() ([]byte, error)
}

// Load restores an Engine from a JSON snapshot produced by Dump.
func Load(snapshot []byte) (Engine, error) {
	return loadFromSnapshot(snapshot, classicalLoader)
}

// loadFromSnapshot restores an Engine using loader to deserialise the snapshot.
// Separated from Load so tests can inject a failing loader.
func loadFromSnapshot(snapshot []byte, loader func([]byte) (gameState, error)) (Engine, error) {
	gs, err := loader(snapshot)
	if err != nil {
		return nil, fmt.Errorf("engine: load snapshot: %w", err)
	}
	return &game{adj: gs, parser: &classicalOrderParser{}}, nil
}

// classicalLoader deserialises a JSON snapshot into a classical game state.
func classicalLoader(snapshot []byte) (gameState, error) {
	return classicalLoaderWith(snapshot, classical.Start)
}

// classicalLoaderWith is the testable core of classicalLoader.
// startFn is injected so tests can simulate errors from classical.Start.
func classicalLoaderWith(snapshot []byte, startFn func() (*state.State, error)) (gameState, error) {
	var snap stateSnapshot
	if err := json.Unmarshal(snapshot, &snap); err != nil || snap.Year == 0 {
		// Empty or unparseable snapshot: start a fresh classical game.
		st, err := startFn()
		if err != nil {
			return nil, err
		}
		return newStateWrapper(st, classical.ClassicalVariant), nil
	}
	ph := classical.NewPhase(snap.Year, snap.Season, snap.PhaseType)
	return buildStateFromSnapshot(classical.Blank(ph), &snap)
}

// buildStateFromSnapshot loads units, supply centres, and dislodgeds from snap
// into st. Extracted so tests can trigger SetUnits/SetDislodgeds error paths.
func buildStateFromSnapshot(st *state.State, snap *stateSnapshot) (gameState, error) {
	if len(snap.Units) > 0 {
		if err := st.SetUnits(snap.Units); err != nil {
			return nil, fmt.Errorf("engine: set units: %w", err)
		}
	}
	if len(snap.SupplyCenters) > 0 {
		st.SetSupplyCenters(snap.SupplyCenters)
	}
	if len(snap.Dislodgeds) > 0 {
		if err := st.SetDislodgeds(snap.Dislodgeds); err != nil {
			return nil, fmt.Errorf("engine: set dislodgeds: %w", err)
		}
	}
	return newStateWrapper(st, classical.ClassicalVariant), nil
}

// UnitInfo holds the type and owning nation of a unit on the board.
type UnitInfo struct {
	Type   string
	Nation string
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
	// Dislodgeds returns a map from province name to nation name for all
	// dislodged units. Used by the bot to validate retreat and disband orders
	// during the Retreat phase.
	Dislodgeds() map[string]string
	// SupplyCenters returns the number of supply centers owned by each nation.
	SupplyCenters() map[string]int
	// Units returns all units on the board keyed by province name.
	Units() map[string]UnitInfo
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

// game implements Engine around a gameState.
type game struct {
	adj    gameState
	parser orderParser
}

// New creates an Engine for the named Diplomacy variant (currently "classical").
func New(variant string) (Engine, error) {
	start, err := lookupVariantStart(variant)
	if err != nil {
		return nil, err
	}
	return newFromVariantStart(start, variant, &classicalOrderParser{})
}

// newFromVariantStart starts a game from the given start function and parser.
// Separated from New so that tests can inject a failing start function.
func newFromVariantStart(start func() (gameState, error), name string, p orderParser) (Engine, error) {
	gs, err := start()
	if err != nil {
		return nil, fmt.Errorf("engine: start %s: %w", name, err)
	}
	return &game{adj: gs, parser: p}, nil
}

// lookupVariantStart returns the start function for the given variant name.
func lookupVariantStart(name string) (func() (gameState, error), error) {
	switch name {
	case "classical":
		return func() (gameState, error) {
			return classicalStartWith(classical.Start)
		}, nil
	default:
		return nil, fmt.Errorf("engine: unknown variant %q", name)
	}
}

// classicalStartWith is the testable core of the classical start function.
// startFn is injected so tests can simulate errors from classical.Start.
func classicalStartWith(startFn func() (*state.State, error)) (gameState, error) {
	st, err := startFn()
	if err != nil {
		return nil, err
	}
	return newStateWrapper(st, classical.ClassicalVariant), nil
}

// SubmitOrder parses the order text and stages it on the game state.
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

// Dislodgeds returns a map from province name to nation name for all dislodged
// units, so the bot can validate retreat and disband orders.
func (g *game) Dislodgeds() map[string]string {
	result := make(map[string]string)
	for prov, unit := range g.adj.Dislodgeds() {
		result[string(prov)] = string(unit.Nation)
	}
	return result
}

// SupplyCenters returns the number of supply centers owned by each nation.
func (g *game) SupplyCenters() map[string]int {
	result := make(map[string]int)
	for _, nation := range g.adj.SupplyCenters() {
		result[string(nation)]++
	}
	return result
}

// Units returns all units on the board keyed by province name.
func (g *game) Units() map[string]UnitInfo {
	result := make(map[string]UnitInfo)
	for prov, unit := range g.adj.Units() {
		result[string(prov)] = UnitInfo{
			Type:   string(unit.Type),
			Nation: string(unit.Nation),
		}
	}
	return result
}

// ---- stateWrapper: adapts *state.State to gameState -------------------------

// stateSnapshot is the JSON format used to persist and restore game state.
type stateSnapshot struct {
	Year          int                             `json:"year"`
	Season        godip.Season                    `json:"season"`
	PhaseType     godip.PhaseType                 `json:"phase_type"`
	Units         map[godip.Province]godip.Unit   `json:"units"`
	SupplyCenters map[godip.Province]godip.Nation `json:"supply_centers"`
	Dislodgeds    map[godip.Province]godip.Unit   `json:"dislodgeds"`
}

// stateWrapper wraps *state.State to implement gameState.
type stateWrapper struct {
	st      *state.State
	variant common.Variant
}

func newStateWrapper(st *state.State, v common.Variant) *stateWrapper {
	return &stateWrapper{st: st, variant: v}
}

func (w *stateWrapper) Phase() gamePhase {
	p := w.st.Phase()
	if p == nil {
		return nil
	}
	return phaseWrapper{p}
}

func (w *stateWrapper) Orders() map[godip.Province]adjOrder {
	result := make(map[godip.Province]adjOrder)
	for p, o := range w.st.Orders() {
		result[p] = o // godip.Adjudicator satisfies adjOrder (has Type())
	}
	return result
}

func (w *stateWrapper) Units() map[godip.Province]godip.Unit {
	return w.st.Units()
}

func (w *stateWrapper) Dislodgeds() map[godip.Province]godip.Unit {
	return w.st.Dislodgeds()
}

func (w *stateWrapper) SupplyCenters() map[godip.Province]godip.Nation {
	return w.st.SupplyCenters()
}

func (w *stateWrapper) SetOrder(p godip.Province, o adjOrder) {
	// Only real godip.Adjudicator values can be staged for adjudication.
	if adj, ok := o.(godip.Adjudicator); ok {
		_ = w.st.SetOrder(p, adj)
	}
	// Stub orders from tests (not godip.Adjudicator) are silently ignored.
}

func (w *stateWrapper) Resolve(_ godip.Province) error {
	// Real adjudication is performed inside Next(). This is a no-op stub
	// so that engine.Resolve() can gather order summaries before advancing.
	return nil
}

func (w *stateWrapper) Next() (gameState, error) {
	return w.nextWith(w.st.Next)
}

// nextWith is the testable core of Next, accepting an advance function so tests
// can inject failures (since state.State.Next never errors in practice).
func (w *stateWrapper) nextWith(advance func() error) (gameState, error) {
	if err := advance(); err != nil {
		return nil, err
	}
	return w, nil // *state.State is mutated in-place by Next()
}

func (w *stateWrapper) SoloWinner() godip.Nation {
	if w.variant.SoloWinner != nil {
		return w.variant.SoloWinner(w.st)
	}
	return ""
}

func (w *stateWrapper) Dump() ([]byte, error) {
	ph := w.st.Phase()
	snap := stateSnapshot{
		Year:          ph.Year(),
		Season:        ph.Season(),
		PhaseType:     ph.Type(),
		Units:         w.st.Units(),
		SupplyCenters: w.st.SupplyCenters(),
		Dislodgeds:    w.st.Dislodgeds(),
	}
	return json.Marshal(snap)
}

// ---- phaseWrapper: adapts godip.Phase to gamePhase -------------------------

// phaseWrapper wraps godip.Phase to implement gamePhase.
type phaseWrapper struct{ godip.Phase }

func (pw phaseWrapper) DefaultOrder(p godip.Province) adjOrder {
	o := pw.Phase.DefaultOrder(p)
	if o == nil {
		return nil
	}
	return o // godip.Adjudicator satisfies adjOrder
}
