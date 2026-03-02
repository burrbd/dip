package bot

import (
	"testing"

	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/session"
	"github.com/cheekybits/is"
)

// mockEngineWithUnits is a mockEngine whose Units() returns a configured map.
type mockEngineWithUnits struct {
	mockEngine
	units map[string]engine.UnitInfo
}

func (e *mockEngineWithUnits) Units() map[string]engine.UnitInfo {
	if e.units == nil {
		return make(map[string]engine.UnitInfo)
	}
	return e.units
}

func TestAutocomplete_NilSession_ReturnsEmpty(t *testing.T) {
	result := Autocomplete(nil, "England")
	if len(result) != 0 {
		t.Errorf("expected no suggestions for nil session, got: %v", result)
	}
}

func TestAutocomplete_NoUnits_ReturnsEmpty(t *testing.T) {
	eng := &mockEngineWithUnits{
		mockEngine: mockEngine{phase: "Spring 1901 Movement", dump: []byte(`{}`)},
	}
	sess := &session.Session{
		Phase:        "Spring 1901 Movement",
		Players:      map[string]string{"u1": "England"},
		StagedOrders: make(map[string][]string),
		Submitted:    make(map[string]bool),
		Eng:          eng,
	}
	result := Autocomplete(sess, "England")
	if len(result) != 0 {
		t.Errorf("expected no suggestions when no units, got: %v", result)
	}
}

func TestAutocomplete_ReturnsHoldOrdersForNation(t *testing.T) {
	is := is.New(t)
	eng := &mockEngineWithUnits{
		mockEngine: mockEngine{phase: "Spring 1901 Movement", dump: []byte(`{}`)},
		units: map[string]engine.UnitInfo{
			"London":    {Type: "Fleet", Nation: "England"},
			"Edinburgh": {Type: "Army", Nation: "England"},
			"Paris":     {Type: "Army", Nation: "France"},
		},
	}
	sess := &session.Session{
		Phase:        "Spring 1901 Movement",
		Players:      map[string]string{"u1": "England"},
		StagedOrders: make(map[string][]string),
		Submitted:    make(map[string]bool),
		Eng:          eng,
	}
	result := Autocomplete(sess, "England")
	is.Equal(len(result), 2)
	is.Equal(result[0], "A Edinburgh H")
	is.Equal(result[1], "F London H")
}

func TestAutocomplete_OtherNationUnitsExcluded(t *testing.T) {
	eng := &mockEngineWithUnits{
		mockEngine: mockEngine{phase: "Spring 1901 Movement", dump: []byte(`{}`)},
		units: map[string]engine.UnitInfo{
			"Paris": {Type: "Army", Nation: "France"},
		},
	}
	sess := &session.Session{
		Phase:        "Spring 1901 Movement",
		Players:      map[string]string{"u1": "England"},
		StagedOrders: make(map[string][]string),
		Submitted:    make(map[string]bool),
		Eng:          eng,
	}
	result := Autocomplete(sess, "England")
	if len(result) != 0 {
		t.Errorf("expected no suggestions for other nation's units, got: %v", result)
	}
}
