package engine

// Integration tests that exercise the real *state.State code paths via
// stateWrapper and phaseWrapper. These tests call New("classical") or
// classicalLoader directly to reach the godip-backed implementations.

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/cheekybits/is"
	"github.com/zond/godip"
	"github.com/zond/godip/orders"
	"github.com/zond/godip/state"
	"github.com/zond/godip/variants/classical"
	"github.com/zond/godip/variants/common"
)

// newClassicalWrapper creates a fresh classical stateWrapper for tests.
func newClassicalWrapper(t *testing.T) *stateWrapper {
	t.Helper()
	st, err := classical.Start()
	if err != nil {
		t.Fatalf("classical.Start: %v", err)
	}
	return newStateWrapper(st, classical.ClassicalVariant)
}

func TestStateWrapper_Phase(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	ph := w.Phase()
	is.NotNil(ph)
	is.Equal(ph.Type(), godip.Movement)
	is.Equal(ph.Year(), 1901)
	is.Equal(ph.Season(), godip.Spring)
}

func TestStateWrapper_Units(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	units := w.Units()
	// Classical 1901 starts with 22 units (3 per power × 7 powers + Russia has 4).
	is.Equal(len(units), 22)
}

func TestStateWrapper_SupplyCenters(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	scs := w.SupplyCenters()
	// Real godip tracks only home SCs at game start (22 owned, neutrals untracked).
	is.Equal(len(scs), 22)
}

func TestStateWrapper_Dislodgeds(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	// No dislodged units at game start.
	is.Equal(len(w.Dislodgeds()), 0)
}

func TestStateWrapper_Orders_Empty(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	is.Equal(len(w.Orders()), 0)
}

func TestStateWrapper_SetOrder_RealAdjudicator(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	// Hold is a real godip.Adjudicator; SetOrder should accept it.
	w.SetOrder("vie", orders.Hold("vie"))
	staged := w.Orders()
	_, ok := staged["vie"]
	is.Equal(ok, true)
}

func TestStateWrapper_SetOrder_StubOrderIsNoOp(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	// A stubOrder is not a godip.Adjudicator; SetOrder should be a no-op.
	w.SetOrder("vie", &stubOrder{t: "Hold"})
	is.Equal(len(w.Orders()), 0)
}

func TestStateWrapper_Resolve(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	// Resolve is a no-op stub; it must return nil.
	err := w.Resolve("vie")
	is.NoErr(err)
}

func TestStateWrapper_Next(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	// Advance without any orders — godip applies NMR holds for all units.
	next, err := w.Next()
	is.NoErr(err)
	is.NotNil(next)
	// State is mutated in-place; Next returns the same wrapper.
	is.Equal(next, gameState(w))
}

func TestStateWrapper_SoloWinner_NoWinner(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	is.Equal(w.SoloWinner(), godip.Nation(""))
}

func TestStateWrapper_Dump_RoundTrip(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)

	data, err := w.Dump()
	is.NoErr(err)
	is.NotNil(data)

	var snap stateSnapshot
	err = json.Unmarshal(data, &snap)
	is.NoErr(err)
	is.Equal(snap.Year, 1901)
	is.Equal(snap.Season, godip.Spring)
	is.Equal(snap.PhaseType, godip.Movement)
	is.Equal(len(snap.Units), 22)
}

func TestPhaseWrapper_DefaultOrder(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	ph := w.Phase()
	is.NotNil(ph)
	// Vienna has an Austrian army; DefaultOrder should return a Hold adjudicator.
	def := ph.DefaultOrder("vie")
	is.NotNil(def)
	is.Equal(def.Type(), godip.Hold)
}

func TestPhaseWrapper_DefaultOrder_SeaProvince(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	ph := w.Phase()
	// Real godip returns a Hold for any province regardless of occupation.
	def := ph.DefaultOrder("nth") // North Sea
	is.NotNil(def)
	is.Equal(def.Type(), godip.Hold)
}

func TestClassicalOrderParser_Parse(t *testing.T) {
	is := is.New(t)
	p := classicalOrderParser{}
	prov, order, err := p.Parse("Austria", "A Vie-Bud")
	is.NoErr(err)
	is.Equal(prov, godip.Province("Vie-Bud")) // stub parser: second token is province
	is.NotNil(order)
}

func TestClassicalOrderParser_TooFewTokens(t *testing.T) {
	is := is.New(t)
	p := classicalOrderParser{}
	_, _, err := p.Parse("Austria", "Vie")
	is.Err(err)
}

func TestClassicalLoader_EmptySnapshot(t *testing.T) {
	is := is.New(t)
	gs, err := classicalLoader([]byte(`{}`))
	is.NoErr(err)
	is.NotNil(gs)
	// Should default to Spring 1901 Movement.
	is.Equal(gs.Phase().Year(), 1901)
}

func TestClassicalLoader_FullSnapshot(t *testing.T) {
	is := is.New(t)
	// Dump a fresh state, then reload it.
	w := newClassicalWrapper(t)
	data, err := w.Dump()
	is.NoErr(err)

	gs, err := classicalLoader(data)
	is.NoErr(err)
	ph := gs.Phase()
	is.Equal(ph.Year(), 1901)
	is.Equal(ph.Season(), godip.Spring)
	is.Equal(ph.Type(), godip.Movement)
	is.Equal(len(gs.Units()), 22)
}

func TestClassicalLoader_BadJSON(t *testing.T) {
	is := is.New(t)
	// Invalid JSON should fall back to a fresh start (no error).
	gs, err := classicalLoader([]byte(`not json`))
	is.NoErr(err)
	is.NotNil(gs)
}

func TestLookupVariantStart_Classical(t *testing.T) {
	is := is.New(t)
	start, err := lookupVariantStart("classical")
	is.NoErr(err)
	gs, err := start()
	is.NoErr(err)
	is.NotNil(gs)
}

func TestStateWrapper_Phase_NilPhase(t *testing.T) {
	// classical.Blank(nil) creates a state with nil phase.
	st := classical.Blank(nil)
	w := newStateWrapper(st, classical.ClassicalVariant)
	ph := w.Phase()
	if ph != nil {
		t.Errorf("expected nil gamePhase for nil-phase state, got %v", ph)
	}
}

func TestStateWrapper_NextWith_Error(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	advanceErr := errors.New("next failed")
	_, err := w.nextWith(func() error { return advanceErr })
	is.Err(err)
}

func TestStateWrapper_SoloWinner_NilVariantFn(t *testing.T) {
	is := is.New(t)
	st, _ := classical.Start()
	// A variant with no SoloWinner function should return the empty nation.
	w := newStateWrapper(st, common.Variant{})
	is.Equal(w.SoloWinner(), godip.Nation(""))
}

func TestPhaseWrapper_DefaultOrder_NonMovementPhase(t *testing.T) {
	// During Retreat/Adjustment phases, godip's DefaultOrder returns nil.
	ph := classical.NewPhase(1901, godip.Spring, godip.Retreat)
	pw := phaseWrapper{ph}
	def := pw.DefaultOrder("vie")
	if def != nil {
		t.Errorf("expected nil adjOrder for non-Movement phase, got %v", def)
	}
}

func TestClassicalLoaderWith_StartError(t *testing.T) {
	is := is.New(t)
	startErr := errors.New("start failed")
	failStart := func() (*state.State, error) { return nil, startErr }
	_, err := classicalLoaderWith([]byte(`{}`), failStart)
	is.Err(err)
}

func TestClassicalStartWith_Error(t *testing.T) {
	is := is.New(t)
	startErr := errors.New("start failed")
	_, err := classicalStartWith(func() (*state.State, error) { return nil, startErr })
	is.Err(err)
}


func TestParsedOrder_Type(t *testing.T) {
	is := is.New(t)
	o := &parsedOrder{orderText: "A Vie-Bud"}
	is.Equal(o.Type(), godip.OrderType("A Vie-Bud"))
}
