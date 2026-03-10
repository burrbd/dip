package dipmap

import (
	"strings"
	"testing"

	"github.com/cheekybits/is"
)

// noCoordsSVG is a synthetic SVG whose province element has no numeric
// coordinate data, exercising the len(xs)==0 branch in provinceCenter.
const noCoordsSVG = `<?xml version="1.0"?>
<svg xmlns="http://www.w3.org/2000/svg"
   xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape">
  <polygon inkscape:label="weird" points="abc def" id="weird"/>
</svg>`

// scForegroundSVG is a synthetic SVG that includes a supply-centers foreground
// copy group before the province shapes, so tests can verify that Overlay
// re-injects it after the units layer.
const scForegroundSVG = `<?xml version="1.0"?>
<svg xmlns="http://www.w3.org/2000/svg"
   xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
   viewBox="0 0 100 100">
  <g inkscape:label="supply-centers foreground copy"><circle r="5" cx="15" cy="15"/></g>
  <polygon inkscape:label="foo" points="10,10 20,10 20,20 10,20" id="foo"/>
</svg>`

// ---- Overlay ----------------------------------------------------------------

func TestOverlay_EmptyUnits_ReturnsSVGUnchanged(t *testing.T) {
	is := is.New(t)
	result, err := Overlay([]byte(renderTestSVG), nil)
	is.NoErr(err)
	is.Equal(string(result), renderTestSVG)
}

func TestOverlay_KnownProvince_InjectsGlyph(t *testing.T) {
	is := is.New(t)
	units := map[string]Unit{
		"foo": {Type: "Army", Nation: "Austria"},
	}
	result, err := Overlay([]byte(renderTestSVG), units)
	is.NoErr(err)
	s := string(result)
	if !strings.Contains(s, `id="units"`) {
		t.Error(`expected <g id="units"> layer in SVG`)
	}
	if !strings.Contains(s, `fill="#CC0000"`) {
		t.Error("expected Austria's colour in unit glyph")
	}
	if !strings.Contains(s, ">A<") {
		t.Error("expected Army label 'A' in unit glyph")
	}
}

func TestOverlay_FleetUnit_InjectsFleetLabel(t *testing.T) {
	is := is.New(t)
	units := map[string]Unit{
		"foo": {Type: "Fleet", Nation: "England"},
	}
	result, err := Overlay([]byte(renderTestSVG), units)
	is.NoErr(err)
	if !strings.Contains(string(result), ">F<") {
		t.Error("expected Fleet label 'F' in unit glyph")
	}
}

func TestOverlay_UnknownProvince_ReturnsSVGUnchanged(t *testing.T) {
	is := is.New(t)
	// Province name is not in the SVG → no glyph generated → no units layer.
	units := map[string]Unit{
		"nonexistent": {Type: "Army", Nation: "France"},
	}
	result, err := Overlay([]byte(renderTestSVG), units)
	is.NoErr(err)
	if strings.Contains(string(result), `id="units"`) {
		t.Error(`unexpected <g id="units"> when no province is found`)
	}
}

// ---- provinceCenter ---------------------------------------------------------

func TestProvinceCenter_ValidProvince_ReturnsCentroid(t *testing.T) {
	is := is.New(t)
	// "foo" has points "10,10 20,10 20,20 10,20" → centroid = (15, 15).
	cx, cy, ok := provinceCenter(renderTestSVG, "foo")
	is.Equal(ok, true)
	is.Equal(cx, 15.0)
	is.Equal(cy, 15.0)
}

func TestProvinceCenter_UnknownProvince_ReturnsFalse(t *testing.T) {
	is := is.New(t)
	_, _, ok := provinceCenter(renderTestSVG, "nonexistent")
	is.Equal(ok, false)
}

func TestProvinceCenter_NoNumericCoordinates_ReturnsFalse(t *testing.T) {
	is := is.New(t)
	// Province is found but its points attribute contains no parseable numbers.
	_, _, ok := provinceCenter(noCoordsSVG, "weird")
	is.Equal(ok, false)
}

// ---- liftSupplyCenterForeground ---------------------------------------------

func TestLiftSupplyCenterForeground_Found_RemovesGroupAndReturnsIt(t *testing.T) {
	is := is.New(t)
	result, lifted := liftSupplyCenterForeground(scForegroundSVG)
	if strings.Contains(result, `inkscape:label="supply-centers foreground copy"`) {
		t.Error("expected foreground copy group removed from result")
	}
	if !strings.Contains(lifted, `inkscape:label="supply-centers foreground copy"`) {
		t.Error("expected foreground copy group in lifted string")
	}
	if !strings.Contains(lifted, `<circle`) {
		t.Error("expected circle element inside lifted group")
	}
	// Province shape must still be in result so Overlay can compute centroids.
	if !strings.Contains(result, `inkscape:label="foo"`) {
		t.Error("expected province shape preserved in result")
	}
	is.NotNil(result)
}

func TestLiftSupplyCenterForeground_LabelNotFound_ReturnsSVGUnchanged(t *testing.T) {
	result, lifted := liftSupplyCenterForeground(renderTestSVG)
	if result != renderTestSVG {
		t.Error("expected svg unchanged when label not found")
	}
	if lifted != "" {
		t.Errorf("expected empty lifted, got %q", lifted)
	}
}

func TestLiftSupplyCenterForeground_NoGBeforeLabel_ReturnsSVGUnchanged(t *testing.T) {
	// Label present but no <g before it — should return unchanged.
	svg := `<svg>` + `inkscape:label="supply-centers foreground copy"` + `</svg>`
	result, lifted := liftSupplyCenterForeground(svg)
	if result != svg {
		t.Error("expected svg unchanged when no <g before label")
	}
	if lifted != "" {
		t.Errorf("expected empty lifted, got %q", lifted)
	}
}

func TestLiftSupplyCenterForeground_UnclosedGroup_ReturnsSVGUnchanged(t *testing.T) {
	// Group is opened but never closed — findGroupEnd returns -1.
	svg := `<svg><g inkscape:label="supply-centers foreground copy"><unclosed</svg>`
	result, lifted := liftSupplyCenterForeground(svg)
	if result != svg {
		t.Error("expected svg unchanged when group end not found")
	}
	if lifted != "" {
		t.Errorf("expected empty lifted, got %q", lifted)
	}
}

// ---- findGroupEnd -----------------------------------------------------------

func TestFindGroupEnd_SimpleGroup_ReturnsCorrectOffset(t *testing.T) {
	s := `<g id="x"><rect/></g>extra`
	end := findGroupEnd(s)
	if end != len(`<g id="x"><rect/></g>`) {
		t.Errorf("expected %d, got %d", len(`<g id="x"><rect/></g>`), end)
	}
}

func TestFindGroupEnd_NestedGroups_ReturnsOutermostClose(t *testing.T) {
	s := `<g><g><rect/></g></g>after`
	end := findGroupEnd(s)
	if end != len(`<g><g><rect/></g></g>`) {
		t.Errorf("expected %d, got %d", len(`<g><g><rect/></g></g>`), end)
	}
}

func TestFindGroupEnd_SelfClosingInnerGroup_NotCounted(t *testing.T) {
	s := `<g id="outer"><g id="sc"/></g>tail`
	end := findGroupEnd(s)
	if end != len(`<g id="outer"><g id="sc"/></g>`) {
		t.Errorf("expected %d, got %d", len(`<g id="outer"><g id="sc"/></g>`), end)
	}
}

func TestFindGroupEnd_NonGroupTag_IgnoredInDepthCount(t *testing.T) {
	// <google> starts with <g but is not a group tag, so depth is unaffected.
	s := `<g id="outer"><google>text</google></g>after`
	end := findGroupEnd(s)
	if end != len(`<g id="outer"><google>text</google></g>`) {
		t.Errorf("expected %d, got %d", len(`<g id="outer"><google>text</google></g>`), end)
	}
}

func TestFindGroupEnd_UnclosedGroup_ReturnsMinusOne(t *testing.T) {
	s := `<g id="x"><rect/>`
	end := findGroupEnd(s)
	if end != -1 {
		t.Errorf("expected -1, got %d", end)
	}
}

// ---- Overlay (with supply-centre foreground lifting) ------------------------

func TestOverlay_LiftsForegroundAboveUnits(t *testing.T) {
	units := map[string]Unit{
		"foo": {Type: "Army", Nation: "England"},
	}
	result, err := Overlay([]byte(scForegroundSVG), units)
	if err != nil {
		t.Fatal(err)
	}
	s := string(result)
	unitsIdx := strings.Index(s, `id="units"`)
	foreIdx := strings.Index(s, `supply-centers foreground copy`)
	if unitsIdx < 0 {
		t.Fatal("expected units layer in result")
	}
	if foreIdx < 0 {
		t.Fatal("expected supply-centers foreground copy in result")
	}
	if foreIdx < unitsIdx {
		t.Error("expected supply-centers foreground copy to appear AFTER units layer")
	}
}

// ---- unitGlyph --------------------------------------------------------------

func TestUnitGlyph_Army_KnownNation_UsesNationColour(t *testing.T) {
	glyph := unitGlyph(100, 200, Unit{Type: "Army", Nation: "France"})
	if !strings.Contains(glyph, ">A<") {
		t.Error("expected Army label 'A'")
	}
	if !strings.Contains(glyph, "#3399CC") {
		t.Error("expected France's colour #3399CC")
	}
}

func TestUnitGlyph_Fleet_UnknownNation_UsesDefaultColour(t *testing.T) {
	glyph := unitGlyph(100, 200, Unit{Type: "Fleet", Nation: "Narnia"})
	if !strings.Contains(glyph, ">F<") {
		t.Error("expected Fleet label 'F'")
	}
	if !strings.Contains(glyph, "#333333") {
		t.Error("expected default colour #333333 for unknown nation")
	}
}
