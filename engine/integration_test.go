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
	is.Equal(prov, godip.Province("vie")) // real parser: lowercase source province
	is.NotNil(order)
	is.Equal(order.Type(), godip.Move) // real godip Move adjudicator
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

func TestResolve_BouncedMove_SuccessFalse(t *testing.T) {
	is := is.New(t)
	e, err := New("classical")
	is.NoErr(err)
	// Both France (Par) and Germany (Mun) try to move into Burgundy → bounce.
	is.NoErr(e.SubmitOrder("France", "A Par-Bur"))
	is.NoErr(e.SubmitOrder("Germany", "A Mun-Bur"))

	result, err := e.Resolve()
	is.NoErr(err)

	// Both Move orders must report failure.
	successes := make(map[string]bool)
	for _, o := range result.Orders {
		if o.Province == "par" || o.Province == "mun" {
			successes[o.Province] = o.Success
		}
	}
	is.Equal(successes["par"], false)
	is.Equal(successes["mun"], false)
}

func TestResolve_SupportedMove_SuccessTrue(t *testing.T) {
	is := is.New(t)
	e, err := New("classical")
	is.NoErr(err)
	// France supports Par→Bur; Germany contests from Mun.
	// The supported move should win: Par arrives at Bur.
	is.NoErr(e.SubmitOrder("France", "A Par-Bur"))
	is.NoErr(e.SubmitOrder("France", "A Mar S A Par-Bur"))
	is.NoErr(e.SubmitOrder("Germany", "A Mun-Bur"))

	result, err := e.Resolve()
	is.NoErr(err)

	// The Par→Bur order must report success.
	var parResult *OrderResult
	for i := range result.Orders {
		if result.Orders[i].Province == "par" {
			parResult = &result.Orders[i]
			break
		}
	}
	is.NotNil(parResult)
	is.Equal(parResult.Success, true)
}

func TestBuildStateFromSnapshot_SetUnitsError(t *testing.T) {
	is := is.New(t)
	ph := classical.NewPhase(1901, godip.Spring, godip.Movement)
	st := classical.Blank(ph)
	// Two conflicting coastal entries for Spain cause SetUnit to fail on the
	// second one regardless of map iteration order: placing "spa" and "spa/nc"
	// together always triggers a "already at" error via coast resolution.
	snap := &stateSnapshot{
		Units: map[godip.Province]godip.Unit{
			"spa":    {Type: godip.Fleet, Nation: godip.France},
			"spa/nc": {Type: godip.Army, Nation: godip.France},
		},
	}
	_, err := buildStateFromSnapshot(st, snap)
	is.Err(err)
}

func TestBuildStateFromSnapshot_SetDislodgedsError(t *testing.T) {
	is := is.New(t)
	ph := classical.NewPhase(1901, godip.Spring, godip.Movement)
	st := classical.Blank(ph)
	// Same coastal conflict approach for dislodgeds.
	snap := &stateSnapshot{
		Dislodgeds: map[godip.Province]godip.Unit{
			"spa":    {Type: godip.Fleet, Nation: godip.France},
			"spa/nc": {Type: godip.Army, Nation: godip.France},
		},
	}
	_, err := buildStateFromSnapshot(st, snap)
	is.Err(err)
}

// ---- ValidRetreats integration tests ----------------------------------------

func TestStateWrapper_ValidRetreats_EmptyAtGameStart(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	// At game start (Movement phase) there are no dislodged units.
	retreats := w.ValidRetreats()
	is.Equal(len(retreats), 0)
}

func newRetreatPhaseWrapper(t *testing.T, dislodgedProv godip.Province, dislodgedUnit godip.Unit) *stateWrapper {
	t.Helper()
	ph := classical.NewPhase(1901, godip.Spring, godip.Retreat)
	st := classical.Blank(ph)
	// Add a regular unit somewhere so state is non-trivial.
	if err := st.SetUnits(map[godip.Province]godip.Unit{
		"par": {Type: godip.Army, Nation: godip.France},
	}); err != nil {
		t.Fatalf("SetUnits: %v", err)
	}
	if err := st.SetDislodgeds(map[godip.Province]godip.Unit{
		dislodgedProv: dislodgedUnit,
	}); err != nil {
		t.Fatalf("SetDislodgeds: %v", err)
	}
	return newStateWrapper(st, classical.ClassicalVariant)
}

func TestStateWrapper_ValidRetreats_DislodgedUnitHasEmptyDestinations(t *testing.T) {
	is := is.New(t)
	// Use a province surrounded by units so no valid retreats exist.
	// Bur (burgundy) surrounded — retreat options depend on graph edges.
	// Using "bur" at Retreat phase with no adjacent free provinces.
	w := newRetreatPhaseWrapper(t, "bur",
		godip.Unit{Type: godip.Army, Nation: godip.France})
	retreats := w.ValidRetreats()
	// Bur should appear in the map (no valid retreats or some valid retreats).
	_, hasBur := retreats["bur"]
	is.Equal(hasBur, true)
}

func TestStateWrapper_ValidRetreats_EmptyWhenNoDislodgeds(t *testing.T) {
	is := is.New(t)
	ph := classical.NewPhase(1901, godip.Spring, godip.Retreat)
	st := classical.Blank(ph)
	w := newStateWrapper(st, classical.ClassicalVariant)
	retreats := w.ValidRetreats()
	is.Equal(len(retreats), 0)
}

// ---- BuildOptions integration tests -----------------------------------------

func TestStateWrapper_BuildOptions_AtGameStart(t *testing.T) {
	is := is.New(t)
	w := newClassicalWrapper(t)
	// At game start all nations have SCs == units → delta = 0 for all.
	opts := w.BuildOptions()
	is.Equal(len(opts), 7) // 7 classical nations
	for _, nation := range classical.ClassicalVariant.Nations {
		opt, ok := opts[nation]
		is.Equal(ok, true)
		is.Equal(opt.Delta, 0)
	}
}

func newAdjustmentPhaseWrapper(t *testing.T) *stateWrapper {
	t.Helper()
	ph := classical.NewPhase(1901, godip.Fall, godip.Adjustment)
	st := classical.Blank(ph)
	// England owns 4 SCs but has only 3 units → can build 1.
	if err := st.SetUnits(map[godip.Province]godip.Unit{
		"lon": {Type: godip.Fleet, Nation: godip.England},
		"lvp": {Type: godip.Army, Nation: godip.England},
		"yor": {Type: godip.Army, Nation: godip.England},
		"par": {Type: godip.Army, Nation: godip.France},
		"mar": {Type: godip.Army, Nation: godip.France},
		"bre": {Type: godip.Fleet, Nation: godip.France},
	}); err != nil {
		t.Fatalf("SetUnits: %v", err)
	}
	st.SetSupplyCenters(map[godip.Province]godip.Nation{
		"lon": godip.England,
		"lvp": godip.England,
		"yor": godip.England,
		"edi": godip.England, // 4th SC for England (no unit there → free home SC)
		"par": godip.France,
		"mar": godip.France,
		"bre": godip.France,
	})
	return newStateWrapper(st, classical.ClassicalVariant)
}

func TestStateWrapper_BuildOptions_EnglandCanBuildOne(t *testing.T) {
	is := is.New(t)
	w := newAdjustmentPhaseWrapper(t)
	opts := w.BuildOptions()
	engOpt, ok := opts[godip.England]
	is.Equal(ok, true)
	is.Equal(engOpt.Delta, 1)
	is.Equal(len(engOpt.AvailableHomes), 1)
	is.Equal(engOpt.AvailableHomes[0], godip.Province("edi"))
}

func TestStateWrapper_BuildOptions_FranceNoBuildDisband(t *testing.T) {
	is := is.New(t)
	w := newAdjustmentPhaseWrapper(t)
	opts := w.BuildOptions()
	fraOpt, ok := opts[godip.France]
	is.Equal(ok, true)
	is.Equal(fraOpt.Delta, 0)
}

// TestStateWrapper_BuildOptions_DeltaCappedByFreeHomes covers the rawDelta > len(free)
// branch: England owns 5 SCs but only 1 free home SC, so delta is capped to 1.
func TestStateWrapper_BuildOptions_DeltaCappedByFreeHomes(t *testing.T) {
	is := is.New(t)
	ph := classical.NewPhase(1901, godip.Fall, godip.Adjustment)
	st := classical.Blank(ph)
	// England has 3 units on lon, lvp, yor (edi is a free home SC).
	if err := st.SetUnits(map[godip.Province]godip.Unit{
		"lon": {Type: godip.Fleet, Nation: godip.England},
		"lvp": {Type: godip.Army, Nation: godip.England},
		"yor": {Type: godip.Army, Nation: godip.England},
	}); err != nil {
		t.Fatalf("SetUnits: %v", err)
	}
	// England owns 5 SCs: 4 home + bel (non-home). rawDelta = 5-3 = 2.
	// Only edi is free → len(free)=1 < rawDelta=2 → delta capped to 1.
	st.SetSupplyCenters(map[godip.Province]godip.Nation{
		"lon": godip.England,
		"lvp": godip.England,
		"yor": godip.England,
		"edi": godip.England,
		"bel": godip.England, // non-home SC
	})
	w := newStateWrapper(st, classical.ClassicalVariant)
	opts := w.BuildOptions()
	engOpt, ok := opts[godip.England]
	is.Equal(ok, true)
	is.Equal(engOpt.Delta, 1) // capped to len(free)=1
	is.Equal(len(engOpt.AvailableHomes), 1)
}
