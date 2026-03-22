// Package engine wraps the godip adjudication engine, providing a simplified
// interface for the bot layer. It hides godip's internal state types and
// exposes only the operations needed to run a Diplomacy game turn.
package engine

import (
	"encoding/json"
	"fmt"

	"github.com/zond/godip"
	"github.com/zond/godip/orders"
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
	// ValidRetreats returns a map from dislodged province to valid retreat
	// destination provinces. An empty slice means the unit must disband.
	ValidRetreats() map[godip.Province][]godip.Province
	// BuildOptions returns build/disband requirements per nation.
	BuildOptions() map[godip.Nation]internalBuildOption
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

// BuildOption describes how many units a nation must build or disband and
// which home supply centres are available for new builds.
type BuildOption struct {
	Delta          int      // positive = builds available, negative = disbands required
	AvailableHomes []string // populated only when Delta > 0
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
	// ValidRetreats returns a map from dislodged province to valid retreat
	// destination provinces. An empty slice means the unit must disband.
	ValidRetreats() map[string][]string
	// BuildOptions returns build/disband requirements per nation.
	BuildOptions() map[string]BuildOption
}

// ResolutionResult summarises what happened when a phase was adjudicated.
type ResolutionResult struct {
	Phase  string // godip phase type, e.g. "Movement"
	Season string // godip season, e.g. "Spring"
	Year   int
	Orders []OrderResult
}

// OrderResult represents the outcome of a single order after adjudication.
type OrderResult struct {
	Province string
	Order    string // canonical order text (e.g. "A lon H", "F lon-nth")
	Nation   string // owning nation at the time the order was staged
	Success  bool
	IsNMR    bool   // true when the order was auto-filled (no player submission)
	Outcome  string // "success", "bounced", "dislodged", or "failed"
}

// game implements Engine around a gameState.
type game struct {
	adj      gameState
	parser   orderParser
	advanced bool // true after Resolve() has called Next(); tells Advance() to skip Next()
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
// It fills NMR orders, calls godip Next() to adjudicate, and compares unit
// positions before and after to determine per-order success.
func (g *game) Resolve() (ResolutionResult, error) {
	phase := g.adj.Phase()
	result := ResolutionResult{
		Phase:  string(phase.Type()),
		Season: string(phase.Season()),
		Year:   phase.Year(),
	}

	// Snapshot player orders, unit positions, and dislodgeds before NMR fill.
	playerOrders := g.adj.Orders()
	preUnits := g.adj.Units()
	preDislodgeds := g.adj.Dislodgeds()

	// Fill NMR so all units have orders before adjudication.
	fillNMR(g.adj)

	// Capture all orders (player + NMR) before advancing.
	allOrders := g.adj.Orders()

	// Advance the state — this is where godip adjudicates all orders.
	next, err := g.adj.Next()
	if err != nil {
		return result, err
	}
	g.adj = next
	g.advanced = true

	postUnits := g.adj.Units()
	postDislodgeds := g.adj.Dislodgeds()

	for prov, ord := range allOrders {
		_, isPlayer := playerOrders[prov]
		nation := string(preUnits[prov].Nation)
		if nation == "" {
			nation = string(preDislodgeds[prov].Nation)
		}
		succeeded := moveSucceeded(ord, prov, preUnits, postUnits)
		outcome := orderOutcome(ord, prov, succeeded, postDislodgeds)
		result.Orders = append(result.Orders, OrderResult{
			Province: string(prov),
			Order:    canonicalOrderText(prov, ord, preUnits, preDislodgeds),
			Nation:   nation,
			Success:  succeeded,
			IsNMR:    !isPlayer,
			Outcome:  outcome,
		})
	}
	return result, nil
}

// moveSucceeded reports whether an order succeeded. For Move orders it checks
// whether the unit arrived at its destination; all other order types return true.
func moveSucceeded(ord adjOrder, src godip.Province, pre, post map[godip.Province]godip.Unit) bool {
	order, ok := ord.(godip.Order)
	if !ok || order.Type() != godip.Move {
		return true
	}
	targets := order.Targets()
	if len(targets) < 2 {
		return true // defensive; real Move orders always have exactly 2 targets
	}
	preUnit, hadUnit := pre[src]
	if !hadUnit {
		return true // no pre-existing unit to track
	}
	dst := targets[1]
	// Check destination province and its super-province (for coastal variants).
	for _, p := range []godip.Province{dst, dst.Super()} {
		if u, found := post[p]; found && u.Nation == preUnit.Nation {
			return true
		}
	}
	return false
}

// canonicalOrderText formats an order as a human-readable string using the
// unit abbreviation derived from preUnits/preDislodgeds and, where possible,
// the order's Targets() for move destinations.
func canonicalOrderText(prov godip.Province, ord adjOrder, units, dislodgeds map[godip.Province]godip.Unit) string {
	unitAbbr := "A"
	u, inUnits := units[prov]
	if !inUnits {
		u = dislodgeds[prov]
	}
	if u.Type == godip.Fleet {
		unitAbbr = "F"
	}

	gOrd, ok := ord.(godip.Order)
	if !ok {
		return fmt.Sprintf("%s %s %s", unitAbbr, prov, ord.Type())
	}
	targets := gOrd.Targets()
	switch gOrd.Type() {
	case godip.Move:
		if len(targets) >= 2 {
			return fmt.Sprintf("%s %s-%s", unitAbbr, targets[0], targets[1])
		}
	case godip.Hold:
		return fmt.Sprintf("%s %s H", unitAbbr, prov)
	case godip.Support:
		if len(targets) >= 3 {
			return fmt.Sprintf("%s %s S %s-%s", unitAbbr, targets[0], targets[1], targets[2])
		}
		if len(targets) >= 2 {
			return fmt.Sprintf("%s %s S %s", unitAbbr, targets[0], targets[1])
		}
	case godip.Convoy:
		if len(targets) >= 3 {
			return fmt.Sprintf("%s %s C %s-%s", unitAbbr, targets[0], targets[1], targets[2])
		}
	case godip.Build:
		if len(targets) >= 1 {
			return fmt.Sprintf("Build %s %s", unitAbbr, targets[0])
		}
	case godip.Disband:
		if len(targets) >= 1 {
			return fmt.Sprintf("%s %s disband", unitAbbr, targets[0])
		}
	}
	return fmt.Sprintf("%s %s %s", unitAbbr, prov, ord.Type())
}

// orderOutcome returns the human-readable outcome of an order after
// adjudication: "success", "bounced", "dislodged", or "failed".
func orderOutcome(ord adjOrder, prov godip.Province, succeeded bool, postDislodgeds map[godip.Province]godip.Unit) string {
	gOrd, ok := ord.(godip.Order)
	if ok && gOrd.Type() == godip.Move {
		if succeeded {
			return "success"
		}
		return "bounced"
	}
	if _, wasDislodged := postDislodgeds[prov]; wasDislodged {
		return "dislodged"
	}
	if !succeeded {
		return "failed"
	}
	return "success"
}

// ValidRetreats returns a map from dislodged province to valid retreat destinations.
func (g *game) ValidRetreats() map[string][]string {
	result := make(map[string][]string)
	for prov, dests := range g.adj.ValidRetreats() {
		strs := make([]string, len(dests))
		for i, d := range dests {
			strs[i] = string(d)
		}
		result[string(prov)] = strs
	}
	return result
}

// BuildOptions returns build/disband requirements per nation.
func (g *game) BuildOptions() map[string]BuildOption {
	result := make(map[string]BuildOption)
	for nation, opt := range g.adj.BuildOptions() {
		homes := make([]string, len(opt.AvailableHomes))
		for i, h := range opt.AvailableHomes {
			homes[i] = string(h)
		}
		result[string(nation)] = BuildOption{
			Delta:          opt.Delta,
			AvailableHomes: homes,
		}
	}
	return result
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

// internalBuildOption is the gameState-internal equivalent of BuildOption
// using godip types.
type internalBuildOption struct {
	Delta          int
	AvailableHomes []godip.Province
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

// ValidRetreats returns a map from dislodged province to valid retreat destinations.
// An empty slice means the unit must disband (no valid retreats available).
func (w *stateWrapper) ValidRetreats() map[godip.Province][]godip.Province {
	result := make(map[godip.Province][]godip.Province)
	dislodgeds := w.st.Dislodgeds()
	if len(dislodgeds) == 0 {
		return result
	}
	// Initialise all dislodged provinces with empty slices (must disband by default).
	for prov := range dislodgeds {
		result[prov] = []godip.Province{}
	}
	// Collect the set of nations with dislodged units.
	nations := make(map[godip.Nation]bool)
	for _, unit := range dislodgeds {
		nations[unit.Nation] = true
	}
	// Use godip's Options API to find valid retreat destinations per nation.
	// state.Options returns: Province → OrderType → SrcProvince → Province(dst)
	for nation := range nations {
		opts := w.st.Options([]godip.Order{orders.MoveOrder}, nation)
		for _, orderTypeMap := range opts {
			for _, srcMap := range orderTypeMap {
				for srcVal, dstMap := range srcMap {
					srcProv := godip.Province(srcVal.(godip.SrcProvince))
					dests := make([]godip.Province, 0, len(dstMap))
					for dstVal := range dstMap {
						dests = append(dests, dstVal.(godip.Province))
					}
					result[srcProv] = dests
				}
			}
		}
	}
	return result
}

// BuildOptions returns build/disband requirements per nation for the current
// Adjustment phase.
func (w *stateWrapper) BuildOptions() map[godip.Nation]internalBuildOption {
	graph := w.st.Graph()
	currentSCs := w.st.SupplyCenters()

	scCount := make(map[godip.Nation]int)
	freeHomeSCs := make(map[godip.Nation][]godip.Province)
	for sc, nation := range currentSCs {
		scCount[nation]++
		// A home SC is one whose original graph owner matches the current owner.
		if orig := graph.SC(sc); orig != nil && *orig == nation {
			if _, _, occupied := w.st.Unit(sc); !occupied {
				freeHomeSCs[nation] = append(freeHomeSCs[nation], sc)
			}
		}
	}

	unitCount := make(map[godip.Nation]int)
	for _, unit := range w.st.Units() {
		unitCount[unit.Nation]++
	}

	result := make(map[godip.Nation]internalBuildOption)
	for _, nation := range w.variant.Nations {
		rawDelta := scCount[nation] - unitCount[nation]
		free := freeHomeSCs[nation]
		delta := rawDelta
		var available []godip.Province
		if rawDelta > 0 {
			if rawDelta > len(free) {
				delta = len(free)
			}
			available = free
		}
		result[nation] = internalBuildOption{Delta: delta, AvailableHomes: available}
	}
	return result
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
