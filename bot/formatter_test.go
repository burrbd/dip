package bot

import (
	"strings"
	"testing"

	"github.com/burrbd/dip/engine"
	"github.com/cheekybits/is"
)

func TestFormatResult_IncludesPhase(t *testing.T) {
	r := engine.ResolutionResult{Phase: "Spring 1901 Movement"}
	out := FormatResult(r)
	if !strings.Contains(out, "Spring 1901 Movement") {
		t.Errorf("expected phase in output, got: %q", out)
	}
}

func TestFormatResult_IncludesOrderCount(t *testing.T) {
	r := engine.ResolutionResult{
		Phase: "Spring 1901 Movement",
		Orders: []engine.OrderResult{
			{Province: "lon", Order: "Move", Success: true},
			{Province: "par", Order: "Hold", Success: false},
		},
	}
	out := FormatResult(r)
	if !strings.Contains(out, "2") {
		t.Errorf("expected order count 2 in output, got: %q", out)
	}
}

func TestFormatResult_MarksSuccessAndFailure(t *testing.T) {
	r := engine.ResolutionResult{
		Phase: "Spring 1901 Movement",
		Orders: []engine.OrderResult{
			{Province: "lon", Order: "Move", Success: true},
			{Province: "par", Order: "Hold", Success: false},
		},
	}
	out := FormatResult(r)
	if !strings.Contains(out, "succeeded") {
		t.Errorf("expected 'succeeded' in output, got: %q", out)
	}
	if !strings.Contains(out, "failed") {
		t.Errorf("expected 'failed' in output, got: %q", out)
	}
}

func TestFormatResult_EmptyOrders(t *testing.T) {
	is := is.New(t)
	r := engine.ResolutionResult{Phase: "Spring 1901 Movement"}
	out := FormatResult(r)
	is.NotNil(out) // must not panic; empty orders is valid
}

func TestFormatStatus_IncludesPhase(t *testing.T) {
	out := FormatStatus("Spring 1901 Movement", map[string]string{"u1": "England"}, map[string]bool{})
	if !strings.Contains(out, "Spring 1901 Movement") {
		t.Errorf("expected phase in output, got: %q", out)
	}
}

func TestFormatStatus_IncludesNationAndSubmitted(t *testing.T) {
	players := map[string]string{"u1": "England", "u2": "France"}
	submitted := map[string]bool{"England": true}
	out := FormatStatus("Spring 1901 Movement", players, submitted)
	if !strings.Contains(out, "England") {
		t.Errorf("expected England in output, got: %q", out)
	}
	if !strings.Contains(out, "submitted") {
		t.Errorf("expected 'submitted' in output, got: %q", out)
	}
	if !strings.Contains(out, "France") {
		t.Errorf("expected France in output, got: %q", out)
	}
	if !strings.Contains(out, "pending") {
		t.Errorf("expected 'pending' in output, got: %q", out)
	}
}

func TestFormatStatus_EmptyPlayers(t *testing.T) {
	is := is.New(t)
	out := FormatStatus("Spring 1901 Movement", map[string]string{}, map[string]bool{})
	is.NotNil(out)
}
