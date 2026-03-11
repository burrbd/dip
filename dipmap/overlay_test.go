package dipmap

import (
	"strings"
	"testing"

	"github.com/cheekybits/is"
)

// ---- test SVG fixtures ------------------------------------------------------

// unitSVG is a minimal SVG that contains pre-placed unit glyph elements
// matching the format produced by cmd/mkapsvg after the Story 10b Bug 2 fix:
// fill is on the <g> element (not the <rect>) so that setAttr on the group
// propagates the nation colour via SVG inheritance.
const unitSVG = `<?xml version="1.0"?>
<svg xmlns="http://www.w3.org/2000/svg"
   xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
   viewBox="0 0 100 100">
  <g id="units">
    <g id="unit-foo-army" transform="translate(50,50)" display="none" fill="#cccccc">
      <rect x="-9" y="-9" width="18" height="18" rx="2" stroke="#ffffff" stroke-width="2"/>
      <text x="0" y="5" text-anchor="middle" font-size="11" fill="#000000">A</text>
    </g>
    <g id="unit-foo-fleet" transform="translate(50,50)" display="none" fill="#cccccc">
      <rect x="-14" y="-6" width="28" height="12" rx="3" stroke="#ffffff" stroke-width="2"/>
      <text x="0" y="5" text-anchor="middle" font-size="8" fill="#000000">F</text>
    </g>
    <g id="unit-bar-army" transform="translate(20,20)" display="none" fill="#cccccc">
      <rect x="-9" y="-9" width="18" height="18" rx="2" stroke="#ffffff" stroke-width="2"/>
      <text x="0" y="5" text-anchor="middle" font-size="11" fill="#000000">A</text>
    </g>
  </g>
</svg>`

// ---- Overlay ----------------------------------------------------------------

func TestOverlay_EmptyUnits_ReturnsSVGUnchanged(t *testing.T) {
	is := is.New(t)
	result, err := Overlay([]byte(unitSVG), nil)
	is.NoErr(err)
	is.Equal(string(result), unitSVG)
}

func TestOverlay_ActivatesArmyGlyph(t *testing.T) {
	is := is.New(t)
	units := map[string]Unit{
		"foo": {Type: "Army", Nation: "Austria"},
	}
	result, err := Overlay([]byte(unitSVG), units)
	is.NoErr(err)
	s := string(result)

	// The army glyph for "foo" must now have display="inline".
	if !strings.Contains(s, `id="unit-foo-army"`) {
		t.Fatal("expected unit-foo-army glyph in SVG")
	}
	if !strings.Contains(s, `display="inline"`) {
		t.Error("expected display=inline on activated glyph")
	}
	// Nation colour must be set.
	if !strings.Contains(s, `fill="#CC0000"`) {
		t.Error("expected Austria's colour #CC0000")
	}
}

func TestOverlay_ActivatesFleetGlyph(t *testing.T) {
	is := is.New(t)
	units := map[string]Unit{
		"foo": {Type: "Fleet", Nation: "England"},
	}
	result, err := Overlay([]byte(unitSVG), units)
	is.NoErr(err)
	s := string(result)
	if !strings.Contains(s, `id="unit-foo-fleet"`) {
		t.Fatal("expected unit-foo-fleet glyph in SVG")
	}
	if !strings.Contains(s, `display="inline"`) {
		t.Error("expected display=inline on activated fleet glyph")
	}
	if !strings.Contains(s, `fill="#003399"`) {
		t.Error("expected England's colour #003399")
	}
}

func TestOverlay_UnknownProvince_LogsAndContinues(t *testing.T) {
	is := is.New(t)
	// Province "nonexistent" has no glyph in the SVG. The function must not
	// return an error; it logs to stderr and skips the province.
	units := map[string]Unit{
		"nonexistent": {Type: "Army", Nation: "France"},
	}
	result, err := Overlay([]byte(unitSVG), units)
	is.NoErr(err)
	// SVG must be otherwise unchanged — no display=inline added.
	if strings.Contains(string(result), `display="inline"`) {
		t.Error("expected no display=inline when province not found")
	}
}

func TestOverlay_CoastalVariant_NormalisesSlash(t *testing.T) {
	// Province "stp/nc" should map to id "unit-stp-nc-fleet".
	svg := `<svg><g id="unit-stp-nc-fleet" display="none"/></svg>`
	units := map[string]Unit{
		"stp/nc": {Type: "Fleet", Nation: "Russia"},
	}
	result, err := Overlay([]byte(svg), units)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(result), `display="inline"`) {
		t.Error("expected display=inline after activating stp/nc fleet")
	}
}

func TestOverlay_AllStartingUnits_ActivatesCorrectCount(t *testing.T) {
	is := is.New(t)
	// Build a synthetic SVG with 22 unit glyph placeholders.
	var glyphs strings.Builder
	startingUnits := map[string]Unit{
		"lon": {Type: "Fleet", Nation: "England"},
		"edi": {Type: "Fleet", Nation: "England"},
		"lvp": {Type: "Army", Nation: "England"},
		"bre": {Type: "Fleet", Nation: "France"},
		"par": {Type: "Army", Nation: "France"},
		"mar": {Type: "Army", Nation: "France"},
		"kie": {Type: "Fleet", Nation: "Germany"},
		"ber": {Type: "Army", Nation: "Germany"},
		"mun": {Type: "Army", Nation: "Germany"},
		"nap": {Type: "Fleet", Nation: "Italy"},
		"rom": {Type: "Army", Nation: "Italy"},
		"ven": {Type: "Army", Nation: "Italy"},
		"tri": {Type: "Fleet", Nation: "Austria"},
		"vie": {Type: "Army", Nation: "Austria"},
		"bud": {Type: "Army", Nation: "Austria"},
		"ank": {Type: "Fleet", Nation: "Turkey"},
		"con": {Type: "Army", Nation: "Turkey"},
		"smy": {Type: "Army", Nation: "Turkey"},
		"sev": {Type: "Fleet", Nation: "Russia"},
		"mos": {Type: "Army", Nation: "Russia"},
		"war": {Type: "Army", Nation: "Russia"},
		"stp/sc": {Type: "Fleet", Nation: "Russia"},
	}
	// Generate placeholder glyphs for all units.
	for prov, u := range startingUnits {
		pid := strings.ReplaceAll(prov, "/", "-")
		t2 := strings.ToLower(u.Type)
		glyphs.WriteString(`<g id="unit-` + pid + `-` + t2 + `" display="none"/>` + "\n")
	}
	svg := `<svg xmlns="http://www.w3.org/2000/svg">` + glyphs.String() + `</svg>`

	result, err := Overlay([]byte(svg), startingUnits)
	is.NoErr(err)

	s := string(result)
	count := strings.Count(s, `display="inline"`)
	if count != 22 {
		t.Errorf("expected 22 display=inline elements, got %d", count)
	}
}

// TestOverlay_NationColour_FranceArmy verifies that Overlay sets the correct
// France blue fill (#3399CC) on the <g> element for a French army (Bug 2
// acceptance criterion: nation colour must appear on unit-par-army).
func TestOverlay_NationColour_FranceArmy(t *testing.T) {
	svg := `<svg><g id="unit-par-army" fill="#cccccc" display="none"><rect stroke="#ffffff"/></g></svg>`
	units := map[string]Unit{
		"par": {Type: "Army", Nation: "France"},
	}
	result, err := Overlay([]byte(svg), units)
	if err != nil {
		t.Fatal(err)
	}
	s := string(result)
	if !strings.Contains(s, `id="unit-par-army"`) {
		t.Fatal("expected unit-par-army in result")
	}
	if !strings.Contains(s, `display="inline"`) {
		t.Error("expected display=inline on activated glyph")
	}
	if !strings.Contains(s, `fill="#3399CC"`) {
		t.Errorf("expected France blue fill=#3399CC, got: %s", s)
	}
}

// ---- setAttr ----------------------------------------------------------------

func TestSetAttr_UpdatesExistingAttribute(t *testing.T) {
	is := is.New(t)
	svg := `<g id="unit-vie-army" display="none" fill="#cccccc"/>`
	result := setAttr(svg, "unit-vie-army", "display", "inline")
	is.Equal(result, `<g id="unit-vie-army" display="inline" fill="#cccccc"/>`)
}

func TestSetAttr_InsertsNewAttribute(t *testing.T) {
	svg := `<g id="unit-vie-army" display="none"/>`
	result := setAttr(svg, "unit-vie-army", "fill", "#CC0000")
	if !strings.Contains(result, `fill="#CC0000"`) {
		t.Errorf("expected fill attribute in result: %q", result)
	}
}

func TestSetAttr_ElementNotFound_ReturnsUnchanged(t *testing.T) {
	is := is.New(t)
	svg := `<g id="unit-vie-army" display="none"/>`
	result := setAttr(svg, "unit-nonexistent-army", "display", "inline")
	is.Equal(result, svg)
}

func TestSetAttr_SelfClosingTag_InsertsAttribute(t *testing.T) {
	// Element with no existing fill; self-closing form.
	svg := `<g id="unit-lon-fleet"/>`
	result := setAttr(svg, "unit-lon-fleet", "fill", "#003399")
	if !strings.Contains(result, `fill="#003399"`) {
		t.Errorf("expected fill attribute in self-closing tag: %q", result)
	}
}

func TestSetAttr_OpenTag_InsertsAttribute(t *testing.T) {
	// Element in open-tag form (not self-closing).
	svg := `<g id="unit-lon-fleet" display="none">` + "\n  <rect/>\n</g>"
	result := setAttr(svg, "unit-lon-fleet", "fill", "#003399")
	if !strings.Contains(result, `fill="#003399"`) {
		t.Errorf("expected fill attribute in open tag: %q", result)
	}
}

func TestSetAttr_NoTagStartBeforeID_ReturnsUnchanged(t *testing.T) {
	// id= appears in the string but there is no '<' before it.
	svg := `id="unit-vie-army" display="none"/>`
	result := setAttr(svg, "unit-vie-army", "display", "inline")
	if result != svg {
		t.Errorf("expected svg unchanged when no tag start found, got %q", result)
	}
}

func TestSetAttr_NoTagEnd_ReturnsUnchanged(t *testing.T) {
	// id= appears inside a tag that has no closing '>'.
	svg := `<g id="unit-vie-army" display="none"`
	result := setAttr(svg, "unit-vie-army", "display", "inline")
	if result != svg {
		t.Errorf("expected svg unchanged when no tag end found, got %q", result)
	}
}

// ---- nationColour -----------------------------------------------------------

func TestNationColour_KnownNation_ReturnsColour(t *testing.T) {
	is := is.New(t)
	is.Equal(nationColour("England"), "#003399")
	is.Equal(nationColour("France"), "#3399CC")
	is.Equal(nationColour("Austria"), "#CC0000")
}

func TestNationColour_UnknownNation_ReturnsDefault(t *testing.T) {
	is := is.New(t)
	is.Equal(nationColour("Narnia"), "#333333")
	is.Equal(nationColour(""), "#333333")
}
