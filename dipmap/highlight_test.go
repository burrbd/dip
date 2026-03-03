package dipmap

import (
	"strings"
	"testing"

	"github.com/cheekybits/is"
)

// minimalSVG is a small SVG that mirrors the structure of godip's classical map:
// a hidden "provinces" layer with province shapes, and a visible empty
// "highlights" layer into which Highlight injects coloured copies.
const minimalSVG = `<svg xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape">
<g style="display:none"
   inkscape:label="provinces"
   id="provinces"
   inkscape:groupmode="layer">
<polygon
   inkscape:label="vie"
   style="fill:#000000;fill-opacity:1;stroke:#000000" points="0,0 10,0 10,10"/>
<polygon
   inkscape:label="lon"
   style="fill:#000000;fill-opacity:1;stroke:#000000" points="20,20 30,20 30,30"/>
<path
   inkscape:label="par"
   style="fill:#000000;fill-opacity:1;stroke:#000000"
   d="M 40,40 L 50,40 L 50,50 Z"/>
</g>
<g style="display:inline"
   inkscape:label="highlights"
   id="highlights"
   inkscape:groupmode="layer" />
</svg>`

func TestHighlight_InjectsPolygonIntoHighlightsLayer(t *testing.T) {
	is := is.New(t)
	result, err := Highlight([]byte(minimalSVG), []string{"vie"})
	is.NoErr(err)
	s := string(result)
	// A new polygon with the highlight colour must appear in the output.
	is.Equal(strings.Contains(s, `fill:#FF6B6B`), true)
	// The injected element must carry the province's original points data.
	is.Equal(strings.Contains(s, `points="0,0 10,0 10,10"`), true)
}

func TestHighlight_InjectsPathIntoHighlightsLayer(t *testing.T) {
	is := is.New(t)
	// "par" is a <path> element in minimalSVG.
	result, err := Highlight([]byte(minimalSVG), []string{"par"})
	is.NoErr(err)
	s := string(result)
	is.Equal(strings.Contains(s, `fill:#FF6B6B`), true)
	is.Equal(strings.Contains(s, `d="M 40,40 L 50,40 L 50,50 Z"`), true)
}

func TestHighlight_OriginalProvinceElementUnchanged(t *testing.T) {
	is := is.New(t)
	result, err := Highlight([]byte(minimalSVG), []string{"vie"})
	is.NoErr(err)
	s := string(result)
	// The original element in the provinces layer must still have fill:#000000.
	is.Equal(strings.Contains(s, `inkscape:label="vie"`+"\n   style=\"fill:#000000"), true)
}

func TestHighlight_DoesNotInjectNonMatchingProvince(t *testing.T) {
	is := is.New(t)
	result, err := Highlight([]byte(minimalSVG), []string{"vie"})
	is.NoErr(err)
	s := string(result)
	// lon was not requested — its colour must not appear in the highlights.
	is.Equal(strings.Contains(s, `points="20,20 30,20 30,30"`), true) // original unchanged
	// Only one coloured element injected (for vie, palette[0]=#FF6B6B).
	is.Equal(strings.Count(s, `fill:#FF6B6B`), 1)
}

func TestHighlight_DifferentProvincesGetDifferentColors(t *testing.T) {
	is := is.New(t)
	result, err := Highlight([]byte(minimalSVG), []string{"vie", "lon"})
	is.NoErr(err)
	s := string(result)
	// vie → palette[0], lon → palette[1]
	is.Equal(strings.Contains(s, `fill:#FF6B6B`), true)
	is.Equal(strings.Contains(s, `fill:#4ECDC4`), true)
}

func TestHighlight_CaseInsensitiveProvinceMatch(t *testing.T) {
	is := is.New(t)
	result, err := Highlight([]byte(minimalSVG), []string{"VIE"})
	is.NoErr(err)
	s := string(result)
	is.Equal(strings.Contains(s, `fill:#FF6B6B`), true)
}

func TestHighlight_EmptyProvinceList_ReturnsOriginal(t *testing.T) {
	is := is.New(t)
	result, err := Highlight([]byte(minimalSVG), nil)
	is.NoErr(err)
	is.Equal(string(result), minimalSVG)
}

func TestHighlight_ProvinceNotInSVG_NoChange(t *testing.T) {
	is := is.New(t)
	result, err := Highlight([]byte(minimalSVG), []string{"xyz"})
	is.NoErr(err)
	// No provinces found → unchanged.
	is.Equal(string(result), minimalSVG)
}

func TestHighlight_PaletteWrapsAround(t *testing.T) {
	is := is.New(t)
	// Request more provinces than palette entries — should not panic.
	provinces := []string{"vie", "lon", "par", "vie", "lon", "par", "vie", "lon", "par"}
	_, err := Highlight([]byte(minimalSVG), provinces)
	is.NoErr(err)
}

func TestExtractProvinceShape_UnknownLabel_ReturnsEmpty(t *testing.T) {
	is := is.New(t)
	elemType, data := extractProvinceShape(minimalSVG, "xyz")
	is.Equal(elemType, "")
	is.Equal(data, "")
}

func TestExtractProvinceShape_MissingPoints_ReturnsEmpty(t *testing.T) {
	is := is.New(t)
	// An element with the label but no points/d attribute.
	svg := `<svg><polygon inkscape:label="bad" style="fill:#000"/></svg>`
	elemType, data := extractProvinceShape(svg, "bad")
	is.Equal(elemType, "")
	is.Equal(data, "")
}

func TestExtractProvinceShape_NoPrecedingAngle_ReturnsEmpty(t *testing.T) {
	is := is.New(t)
	// Label appears at the very start — no '<' before it.
	svg := `inkscape:label="orphan" points="0,0"/>`
	elemType, data := extractProvinceShape(svg, "orphan")
	is.Equal(elemType, "")
	is.Equal(data, "")
}

func TestExtractProvinceShape_NoSelfClose_ReturnsEmpty(t *testing.T) {
	is := is.New(t)
	// Element never self-closes (uses open + close tag form).
	svg := `<polygon inkscape:label="open" points="0,0"></polygon>`
	elemType, data := extractProvinceShape(svg, "open")
	is.Equal(elemType, "")
	is.Equal(data, "")
}
