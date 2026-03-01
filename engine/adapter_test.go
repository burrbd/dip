package engine

import (
	"errors"
	"testing"

	"github.com/cheekybits/is"
	"github.com/zond/godip"
	"github.com/zond/godip/variants"
)

// ---- mock adjudicator -------------------------------------------------------

type mockAdj struct {
	phase      godip.Phase
	orders     map[godip.Province]godip.Order
	units      map[godip.Province]godip.Unit
	dislodgeds map[godip.Province]godip.Unit
	winner     godip.Nation
	nextAdj    godip.Adjudicator
	nextErr    error
	resolveErr map[godip.Province]error
	setOrders  map[godip.Province]godip.Order
	dumpData   []byte
	dumpErr    error
}

func newMockAdj() *mockAdj {
	return &mockAdj{
		orders:    make(map[godip.Province]godip.Order),
		units:     make(map[godip.Province]godip.Unit),
		setOrders: make(map[godip.Province]godip.Order),
	}
}

func (m *mockAdj) Phase() godip.Phase                        { return m.phase }
func (m *mockAdj) Orders() map[godip.Province]godip.Order    { return m.orders }
func (m *mockAdj) Units() map[godip.Province]godip.Unit      { return m.units }
func (m *mockAdj) Dislodgeds() map[godip.Province]godip.Unit { return m.dislodgeds }
func (m *mockAdj) SoloWinner() godip.Nation                  { return m.winner }
func (m *mockAdj) Dump() ([]byte, error)                     { return m.dumpData, m.dumpErr }
func (m *mockAdj) Next() (godip.Adjudicator, error)          { return m.nextAdj, m.nextErr }

func (m *mockAdj) SetOrder(p godip.Province, o godip.Order) {
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
	defaultOrderFn func(godip.Province) godip.Order
}

func (p *mockPhase) Type() godip.PhaseType { return p.typ }
func (p *mockPhase) Year() int             { return p.year }
func (p *mockPhase) Season() godip.Season  { return p.season }
func (p *mockPhase) DefaultOrder(prov godip.Province) godip.Order {
	if p.defaultOrderFn != nil {
		return p.defaultOrderFn(prov)
	}
	return &stubOrder{t: "Hold"}
}

// ---- mock order parser ------------------------------------------------------

type mockParser struct {
	prov  godip.Province
	order godip.Order
	err   error
}

func (mp *mockParser) Parse(_ godip.Nation, _ string) (godip.Province, godip.Order, error) {
	return mp.prov, mp.order, mp.err
}

// ---- stub order -------------------------------------------------------------

type stubOrder struct{ t string }

func (o *stubOrder) Type() godip.OrderType      { return godip.OrderType(o.t) }
func (o *stubOrder) Flags() map[godip.Flag]bool { return nil }

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

	g := &game{adj: adj, parser: &mockParser{}}
	result, err := g.Resolve()

	is.NoErr(err)
	is.Equal(result.Phase, "Movement")
	is.Equal(result.Year, 1901)
	is.Equal(len(result.Orders), 1)
	is.Equal(result.Orders[0].Province, "Vie")
	is.Equal(result.Orders[0].Success, true)
}

func TestResolve_RecordsFailure(t *testing.T) {
	is := is.New(t)
	order := &stubOrder{t: "A Vie-Bud"}
	adj := newMockAdj()
	adj.phase = &mockPhase{typ: godip.Movement, year: 1901, season: godip.Spring}
	adj.orders["Vie"] = order
	adj.resolveErr = map[godip.Province]error{"Vie": errors.New("bounce")}

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
	is.Equal(g.adj, godip.Adjudicator(nextAdj))
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
	is.Equal(g.adj, godip.Adjudicator(fallAdj))
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

func TestNewFromVariant_StartError(t *testing.T) {
	is := is.New(t)
	v := variants.Variant{
		Name: "failing",
		Start: func() (godip.Adjudicator, error) {
			return nil, errors.New("start failed")
		},
	}
	_, err := newFromVariant(v, "failing", &mockParser{})
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
	is.Equal(g.adj, godip.Adjudicator(nextAdj))
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
	is.Equal(g.adj, godip.Adjudicator(springAdj))
}

func TestLoad_CreatesEngine(t *testing.T) {
	is := is.New(t)
	eng, err := Load([]byte(`{}`))
	is.NoErr(err)
	is.NotNil(eng)
}

func TestLoad_PropagatesLoaderError(t *testing.T) {
	is := is.New(t)
	failLoader := func(_ []byte) (godip.Adjudicator, error) {
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
		defaultOrderFn: func(_ godip.Province) godip.Order {
			return nil // simulate no default order available
		},
	}

	fillNMR(adj)

	if len(adj.setOrders) != 0 {
		t.Errorf("expected no orders set when DefaultOrder returns nil, got %d", len(adj.setOrders))
	}
}
