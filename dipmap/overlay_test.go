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
