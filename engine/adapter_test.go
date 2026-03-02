package engine

import (
	"errors"
	"testing"
	"time"

	"github.com/cheekybits/is"
	"github.com/zond/godip"
	"github.com/zond/godip/orders"
)

// ---- mock game state --------------------------------------------------------

type mockAdj struct {
	phase         gamePhase
	orders        map[godip.Province]adjOrder
	units         map[godip.Province]godip.Unit
	dislodgeds    map[godip.Province]godip.Unit
	supplyCenters map[godip.Province]godip.Nation
	winner        godip.Nation
	nextAdj       gameState
	nextErr       error
	resolveErr    map[godip.Province]error
	setOrders     map[godip.Province]adjOrder
	dumpData      []byte
	dumpErr       error
}

func newMockAdj() *mockAdj {
	return &mockAdj{
		orders:    make(map[godip.Province]adjOrder),
		units:     make(map[godip.Province]godip.Unit),
		setOrders: make(map[godip.Province]adjOrder),
	}
}

func (m *mockAdj) Phase() gamePhase                            { return m.phase }
func (m *mockAdj) Orders() map[godip.Province]adjOrder         { return m.orders }
func (m *mockAdj) Units() map[godip.Province]godip.Unit        { return m.units }
func (m *mockAdj) Dislodgeds() map[godip.Province]godip.Unit   { return m.dislodgeds }
func (m *mockAdj) SoloWinner() godip.Nation                    { return m.winner }
func (m *mockAdj) Dump() ([]byte, error)                       { return m.dumpData, m.dumpErr }
func (m *mockAdj) Next() (gameState, error)                    { return m.nextAdj, m.nextErr }
func (m *mockAdj) SupplyCenters() map[godip.Province]godip.Nation {
	if m.supplyCenters == nil {
		return make(map[godip.Province]godip.Nation)
	}
	return m.supplyCenters
}

func (m *mockAdj) SetOrder(p godip.Province, o adjOrder) {
	m.setOrders[p] = o
	m.orders[p] = o
}

func (m *mockAdj) Resolve(p godip.Province) error {
	if m.resolveErr != nil {
		return m.resolveErr[p]
	}
	return nil
}

// ---- mock phase -------------------------------------------------------------

type mockPhase struct {
	typ    godip.PhaseType
	year   int
	season godip.Season
	// defaultOrderFn lets tests control what DefaultOrder returns.
	defaultOrderFn func(godip.Province) adjOrder
}

func (p *mockPhase) Type() godip.PhaseType { return p.typ }
func (p *mockPhase) Year() int             { return p.year }
func (p *mockPhase) Season() godip.Season  { return p.season }
func (p *mockPhase) DefaultOrder(prov godip.Province) adjOrder {
	if p.defaultOrderFn != nil {
		return p.defaultOrderFn(prov)
	}
	return &stubOrder{t: "Hold"}
}

// ---- mock order parser ------------------------------------------------------

type mockParser struct {
	prov  godip.Province
	order adjOrder
	err   error
}

func (mp *mockParser) Parse(_ godip.Nation, _ string) (godip.Province, adjOrder, error) {
	return mp.prov, mp.order, mp.err
}

// ---- stub order -------------------------------------------------------------

type stubOrder struct{ t string }

func (o *stubOrder) Type() godip.OrderType { return godip.OrderType(o.t) }

// ---- mockOrder: implements godip.Order for moveSucceeded edge-case tests ----

type mockOrder struct {
	typ     godip.OrderType
	targets []godip.Province
}

func (o *mockOrder) Type() godip.OrderType                                      { return o.typ }
func (o *mockOrder) DisplayType() godip.OrderType                               { return o.typ }
func (o *mockOrder) Targets() []godip.Province                                  { return o.targets }
func (o *mockOrder) Flags() map[godip.Flag]bool                                 { return nil }
func (o *mockOrder) Parse([]string) (godip.Adjudicator, error)                  { return nil, nil }
func (o *mockOrder) Options(godip.Validator, godip.Nation, godip.Province) godip.Options {
	return nil
}
func (o *mockOrder) At() time.Time                                   { return time.Time{} }
func (o *mockOrder) Validate(godip.Validator) (godip.Nation, error)  { return "", nil }
func (o *mockOrder) Corroborate(godip.Validator) []error             { return nil }
func (o *mockOrder) Execute(godip.State)                             {}

// ---- tests ------------------------------------------------------------------

func TestNew_Classical(t *testing.T) {
	is := is.New(t)
	e, err := New("classical")
	is.NoErr(err)
	is.NotNil(e)
}

func TestNew_UnknownVariant(t *testing.T) {
	is := is.New(t)
	_, err := New("unknown-variant")
	is.Err(err)
}

func TestSubmitOrder_Stages(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	parser := &mockParser{prov: "Vie", order: &stubOrder{t: "A Vie-Bud"}}
	g := &game{adj: adj, parser: parser}

	err := g.SubmitOrder("Austria", "A Vie-Bud")
	is.NoErr(err)
	is.Equal(adj.setOrders["Vie"], parser.order)
}

func TestSubmitOrder_ParseError(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	parser := &mockParser{err: errors.New("bad order")}
	g := &game{adj: adj, parser: parser}

	err := g.SubmitOrder("Austria", "bad order")
	is.Err(err)
	is.Equal(len(adj.setOrders), 0)
}

func TestResolve_ReturnsSummary(t *testing.T) {
	is := is.New(t)
	order := &stubOrder{t: "A Vie-Bud"}
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}
	adj.orders["Vie"] = order
	// Resolve() calls Next() internally; provide a valid next state.
	nextAdj := newMockAdj()
	nextAdj.phase = &mockPhase{typ: godip.Retreat, year: 1901, season: godip.Spring}
	adj.nextAdj = nextAdj

	g := &game{adj: adj, parser: &mockParser{}}
	result, err := g.Resolve()

	is.NoErr(err)
	is.Equal(result.Phase, "Movement")
	is.Equal(result.Year, 1901)
	is.Equal(len(result.Orders), 1)
	is.Equal(result.Orders[0].Province, "Vie")
	// stubOrder is not a godip.Order so moveSucceeded returns true.
	is.Equal(result.Orders[0].Success, true)
}

func TestResolve_RecordsFailure(t *testing.T) {
	is := is.New(t)
	// Use a real Move adjudicator so moveSucceeded can inspect Targets().
	moveOrder := orders.Move("vie", "bud")
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}
	adj.orders["vie"] = moveOrder
	adj.units["vie"] = godip.Unit{Type: godip.Army, Nation: "Austria"}
	// Next() returns a state where the unit is still at "vie" (move bounced).
	nextAdj := newMockAdj()
	nextAdj.phase = &mockPhase{typ: godip.Retreat, year: 1901, season: godip.Spring}
	nextAdj.units["vie"] = godip.Unit{Type: godip.Army, Nation: "Austria"}
	adj.nextAdj = nextAdj

	g := &game{adj: adj, parser: &mockParser{}}
	result, err := g.Resolve()

	is.NoErr(err)
	is.Equal(len(result.Orders), 1)
	is.Equal(result.Orders[0].Success, false)
}

func TestAdvance_CallsNext(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}
	nextAdj := newMockAdj()
	nextAdj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Fall}
	adj.nextAdj = nextAdj

	g := &game{adj: adj, parser: &mockParser{}}
	err := g.Advance()

	is.NoErr(err)
	is.Equal(g.adj, gameState(nextAdj))
}

func TestAdvance_FillsNMR(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.units["Vie"] = godip.Unit{Type: godip.Army, Nation: "Austria"}
	phase := &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}
	adj.phase = phase
	nextAdj := newMockAdj()
	nextAdj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Fall}
	adj.nextAdj = nextAdj

	g := &game{adj: adj, parser: &mockParser{}}
	err := g.Advance()

	is.NoErr(err)
	// Vienna had no order, so a default order must have been staged.
	_, filled := adj.setOrders["Vie"]
	is.Equal(filled, true)
}

func TestAdvance_SkipsEmptyRetreat(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}

	// First Next() goes to an empty retreat phase (no dislodgeds).
	retreatAdj := newMockAdj()
	retreatAdj.phase = &mockPhase{typ: godip.Retreat, year: 1901, season: godip.Spring}
	retreatAdj.dislodgeds = map[godip.Province]godip.Unit{} // empty

	// Second Next() goes to Fall movement.
	fallAdj := newMockAdj()
	fallAdj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Fall}

	adj.nextAdj = retreatAdj
	retreatAdj.nextAdj = fallAdj

	g := &game{adj: adj, parser: &mockParser{}}
	err := g.Advance()

	is.NoErr(err)
	// Should have skipped the empty retreat and landed on Fall movement.
	is.Equal(g.adj, gameState(fallAdj))
}

func TestAdvance_PropagatesNextError(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}
	adj.nextErr = errors.New("godip internal error")

	g := &game{adj: adj, parser: &mockParser{}}
	err := g.Advance()

	is.Err(err)
}

func TestSoloWinner_NoWinner(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	g := &game{adj: adj, parser: &mockParser{}}
	is.Equal(g.SoloWinner(), "")
}

func TestSoloWinner_Winner(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.winner = "France"
	g := &game{adj: adj, parser: &mockParser{}}
	is.Equal(g.SoloWinner(), "France")
}

func TestDump_ReturnsSerialisedState(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.dumpData = []byte(`{"phase":"Spring 1901 Movement"}`)
	g := &game{adj: adj, parser: &mockParser{}}

	data, err := g.Dump()
	is.NoErr(err)
	is.Equal(string(data), `{"phase":"Spring 1901 Movement"}`)
}

func TestDump_PropagatesError(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.dumpErr = errors.New("serialisation failed")
	g := &game{adj: adj, parser: &mockParser{}}

	_, err := g.Dump()
	is.Err(err)
}

func TestNewFromVariantStart_StartError(t *testing.T) {
	is := is.New(t)
	startErr := errors.New("start failed")
	_, err := newFromVariantStart(func() (gameState, error) {
		return nil, startErr
	}, "failing", &mockParser{})
	is.Err(err)
}

func TestAdvance_SkipNextError(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}

	// First Next() goes to an empty retreat phase.
	retreatAdj := newMockAdj()
	retreatAdj.phase = &mockPhase{typ: godip.Retreat, year: 1901, season: godip.Spring}
	retreatAdj.dislodgeds = map[godip.Province]godip.Unit{}
	retreatAdj.nextErr = errors.New("internal error during skip")

	adj.nextAdj = retreatAdj

	g := &game{adj: adj, parser: &mockParser{}}
	err := g.Advance()

	is.Err(err)
}

func TestAdvance_NilPhase(t *testing.T) {
	is := is.New(t)
	// adj.Phase() returns nil — fillNMR should return early without panicking.
	adj := newMockAdj()
	adj.phase = nil
	nextAdj := newMockAdj()
	nextAdj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Fall}
	adj.nextAdj = nextAdj

	g := &game{adj: adj, parser: &mockParser{}}
	err := g.Advance()

	is.NoErr(err)
	is.Equal(g.adj, gameState(nextAdj))
}

func TestAdvance_SkipsEmptyAdjustment(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Fall}

	// First Next() goes to an empty adjustment phase (no units).
	adjustAdj := newMockAdj()
	adjustAdj.phase = &mockPhase{typ: godip.Adjustment, year: 1901, season: godip.Fall}
	// units is empty by default in newMockAdj

	// Second Next() goes to Spring movement.
	springAdj := newMockAdj()
	springAdj.phase = &mockPhase{typ: godip.Movement, year: 1902, season: godip.Spring}

	adj.nextAdj = adjustAdj
	adjustAdj.nextAdj = springAdj

	g := &game{adj: adj, parser: &mockParser{}}
	err := g.Advance()

	is.NoErr(err)
	is.Equal(g.adj, gameState(springAdj))
}

func TestLoad_CreatesEngine(t *testing.T) {
	is := is.New(t)
	eng, err := Load([]byte(`{}`))
	is.NoErr(err)
	is.NotNil(eng)
}

func TestLoad_PropagatesLoaderError(t *testing.T) {
	is := is.New(t)
	failLoader := func(_ []byte) (gameState, error) {
		return nil, errors.New("deserialise failed")
	}
	_, err := loadFromSnapshot([]byte(`{}`), failLoader)
	is.Err(err)
}

func TestPhase_ReturnsFormattedString(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}
	g := &game{adj: adj, parser: &mockParser{}}
	is.Equal(g.Phase(), "Spring 1901 Movement")
}

func TestPhase_NilPhaseReturnsEmpty(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = nil
	g := &game{adj: adj, parser: &mockParser{}}
	is.Equal(g.Phase(), "")
}

func TestIsEmptyPhase_NilPhase(t *testing.T) {
	adj := newMockAdj()
	adj.phase = nil
	if isEmptyPhase(adj) {
		t.Error("isEmptyPhase with nil phase should return false")
	}
}

func TestFillNMR_NilDefaultOrder(t *testing.T) {
	// When DefaultOrder returns nil, SetOrder must NOT be called.
	adj := newMockAdj()
	adj.units["Vie"] = godip.Unit{Type: godip.Army, Nation: "Austria"}
	adj.phase = &mockPhase{
		typ: godip.Movement,
		defaultOrderFn: func(_ godip.Province) adjOrder {
			return nil // simulate no default order available
		},
	}

	fillNMR(adj)

	if len(adj.setOrders) != 0 {
		t.Errorf("expected no orders set when DefaultOrder returns nil, got %d", len(adj.setOrders))
	}
}

func TestFillNMR_FillsDislodgedUnitsInRetreat(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Retreat, year: 1901, season: godip.Spring}
	adj.dislodgeds = map[godip.Province]godip.Unit{
		"Vie": {Type: godip.Army, Nation: "Austria"},
	}

	fillNMR(adj)

	// A default order (disband) must have been staged for the dislodged unit.
	_, filled := adj.setOrders["Vie"]
	is.Equal(filled, true)
}

func TestFillNMR_SkipsDislodgedWithExistingOrder(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Retreat, year: 1901, season: godip.Spring}
	adj.dislodgeds = map[godip.Province]godip.Unit{
		"Vie": {Type: godip.Army, Nation: "Austria"},
	}
	// Already has an order.
	existingOrder := &stubOrder{t: "A Vie R Bud"}
	adj.orders["Vie"] = existingOrder

	fillNMR(adj)

	// The existing order must not be overwritten.
	is.Equal(adj.setOrders["Vie"], nil)
}

func TestFillNMR_NilDefaultOrderForDislodged(t *testing.T) {
	adj := newMockAdj()
	adj.phase = &mockPhase{
		typ: godip.Retreat,
		defaultOrderFn: func(_ godip.Province) adjOrder {
			return nil
		},
	}
	adj.dislodgeds = map[godip.Province]godip.Unit{
		"Vie": {Type: godip.Army, Nation: "Austria"},
	}

	fillNMR(adj)

	if len(adj.setOrders) != 0 {
		t.Errorf("expected no orders set when DefaultOrder returns nil, got %d", len(adj.setOrders))
	}
}

func TestFillNMR_NonRetreatPhaseIgnoresDislodgeds(t *testing.T) {
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}
	adj.dislodgeds = map[godip.Province]godip.Unit{
		"Vie": {Type: godip.Army, Nation: "Austria"},
	}

	fillNMR(adj)

	// Movement phase: dislodgeds must not get orders via fillNMR.
	if _, set := adj.setOrders["Vie"]; set {
		t.Error("fillNMR should not stage orders for dislodgeds during non-Retreat phase")
	}
}

func TestDislodgeds_ReturnsProvinceToNationMap(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.dislodgeds = map[godip.Province]godip.Unit{
		"Vie": {Type: godip.Army, Nation: "Austria"},
		"Mun": {Type: godip.Army, Nation: "Germany"},
	}
	g := &game{adj: adj, parser: &mockParser{}}

	result := g.Dislodgeds()

	is.Equal(len(result), 2)
	is.Equal(result["Vie"], "Austria")
	is.Equal(result["Mun"], "Germany")
}

func TestDislodgeds_EmptyWhenNoneDislodged(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	g := &game{adj: adj, parser: &mockParser{}}

	result := g.Dislodgeds()

	is.Equal(len(result), 0)
}

func TestSupplyCenters_ReturnsCounts(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.supplyCenters = map[godip.Province]godip.Nation{
		"Lon": "England",
		"Edi": "England",
		"Lvp": "England",
		"Par": "France",
	}
	g := &game{adj: adj, parser: &mockParser{}}

	result := g.SupplyCenters()

	is.Equal(result["England"], 3)
	is.Equal(result["France"], 1)
}

func TestSupplyCenters_EmptyWhenNoSCs(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	g := &game{adj: adj, parser: &mockParser{}}

	result := g.SupplyCenters()

	is.Equal(len(result), 0)
}

func TestUnits_ReturnsBoardUnits(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.units = map[godip.Province]godip.Unit{
		"Lon": {Type: godip.Fleet, Nation: "England"},
		"Par": {Type: godip.Army, Nation: "France"},
	}
	g := &game{adj: adj, parser: &mockParser{}}

	result := g.Units()

	is.Equal(len(result), 2)
	is.Equal(result["Lon"].Type, "Fleet")
	is.Equal(result["Lon"].Nation, "England")
	is.Equal(result["Par"].Type, "Army")
	is.Equal(result["Par"].Nation, "France")
}

func TestUnits_EmptyWhenNoUnits(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	g := &game{adj: adj, parser: &mockParser{}}

	result := g.Units()

	is.Equal(len(result), 0)
}

// ---- Resolve() error path ---------------------------------------------------

func TestResolve_PropagatesNextError(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}
	adj.nextErr = errors.New("godip internal error")

	g := &game{adj: adj, parser: &mockParser{}}
	_, err := g.Resolve()

	is.Err(err)
}

// ---- Advance() after Resolve() (advanced=true path) ------------------------

func TestAdvance_AfterResolve_ReturnsNil(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}
	nextAdj := newMockAdj()
	// Non-empty retreat: unit is dislodged, so isEmptyPhase returns false.
	nextAdj.phase = &mockPhase{typ: godip.Retreat, year: 1901, season: godip.Spring}
	nextAdj.dislodgeds = map[godip.Province]godip.Unit{
		"vie": {Type: godip.Army, Nation: "Austria"},
	}
	adj.nextAdj = nextAdj

	g := &game{adj: adj, parser: &mockParser{}}
	_, err := g.Resolve()
	is.NoErr(err)
	is.Equal(g.advanced, true)

	err = g.Advance()
	is.NoErr(err)
	is.Equal(g.advanced, false)
	is.Equal(g.adj, gameState(nextAdj))
}

func TestAdvance_AfterResolve_SkipsEmptyPhase(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}
	// Resolve() will advance to retreatAdj.
	retreatAdj := newMockAdj()
	retreatAdj.phase = &mockPhase{typ: godip.Retreat, year: 1901, season: godip.Spring}
	retreatAdj.dislodgeds = map[godip.Province]godip.Unit{} // empty → will be skipped
	// Advance() will skip retreatAdj and land on fallAdj.
	fallAdj := newMockAdj()
	fallAdj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Fall}
	adj.nextAdj = retreatAdj
	retreatAdj.nextAdj = fallAdj

	g := &game{adj: adj, parser: &mockParser{}}
	_, err := g.Resolve()
	is.NoErr(err)

	err = g.Advance()
	is.NoErr(err)
	is.Equal(g.adj, gameState(fallAdj))
}

func TestAdvance_AfterResolve_PropagatesNextError(t *testing.T) {
	is := is.New(t)
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}
	retreatAdj := newMockAdj()
	retreatAdj.phase = &mockPhase{typ: godip.Retreat, year: 1901, season: godip.Spring}
	retreatAdj.dislodgeds = map[godip.Province]godip.Unit{} // empty → skip attempted
	retreatAdj.nextErr = errors.New("skip error")
	adj.nextAdj = retreatAdj

	g := &game{adj: adj, parser: &mockParser{}}
	_, err := g.Resolve()
	is.NoErr(err)

	err = g.Advance()
	is.Err(err)
}

// ---- moveSucceeded() edge cases --------------------------------------------

func TestMoveSucceeded_MoveNoTargets_ReturnsTrue(t *testing.T) {
	// A godip.Order with Move type but zero targets hits the defensive len<2 check.
	ord := &mockOrder{typ: godip.Move, targets: nil}
	pre := map[godip.Province]godip.Unit{"vie": {Type: godip.Army, Nation: "Austria"}}
	post := map[godip.Province]godip.Unit{}
	if !moveSucceeded(ord, "vie", pre, post) {
		t.Error("expected true when Move order has no targets")
	}
}

func TestMoveSucceeded_MoveNoPreUnit_ReturnsTrue(t *testing.T) {
	// Move order staged but no unit at the source province.
	moveOrder := orders.Move("vie", "bud")
	pre := map[godip.Province]godip.Unit{} // no unit at "vie"
	post := map[godip.Province]godip.Unit{}
	if !moveSucceeded(moveOrder, "vie", pre, post) {
		t.Error("expected true when no pre-existing unit at source")
	}
}

func TestMoveSucceeded_UnitArrivedAtDestination_ReturnsTrue(t *testing.T) {
	// Unit moved from vie to bud → success.
	moveOrder := orders.Move("vie", "bud")
	pre := map[godip.Province]godip.Unit{"vie": {Type: godip.Army, Nation: "Austria"}}
	post := map[godip.Province]godip.Unit{"bud": {Type: godip.Army, Nation: "Austria"}}
	if !moveSucceeded(moveOrder, "vie", pre, post) {
		t.Error("expected true when unit arrived at destination")
	}
}
