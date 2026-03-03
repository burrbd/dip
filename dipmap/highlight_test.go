package dipmap

import (
	"strings"
	"testing"

	"github.com/cheekybits/is"
)

// minimalSVG is a small SVG fragment with two province elements for testing.
const minimalSVG = `<svg xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape">
<polygon
   inkscape:label="vie"
   style="fill:#000000;fill-opacity:1;stroke:#000000" points="0,0"/>
<polygon
   inkscape:label="lon"
   style="fill:#000000;fill-opacity:1;stroke:#000000" points="1,1"/>
<polygon
   inkscape:label="par"
   style="fill:#000000;fill-opacity:1;stroke:#000000" points="2,2"/>
</svg>`

func TestHighlight_ChangesMatchingProvinceFill(t *testing.T) {
	is := is.New(t)
	result, err := Highlight([]byte(minimalSVG), []string{"vie"})
	is.NoErr(err)
	// The highlighted province should no longer have the original black fill.
	s := string(result)
	is.Equal(strings.Contains(s, `inkscape:label="vie"`), true)
	is.Equal(strings.Contains(s, `inkscape:label="vie"`+"\n   style=\"fill:#000000"), false)
}

func TestHighlight_DoesNotChangeNonMatchingProvince(t *testing.T) {
	is := is.New(t)
	result, err := Highlight([]byte(minimalSVG), []string{"vie"})
	is.NoErr(err)
	// lon is not in the highlight list — its fill must stay #000000.
	s := string(result)
	is.Equal(strings.Contains(s, `inkscape:label="lon"`+"\n   style=\"fill:#000000"), true)
}

func TestHighlight_DifferentProvincesGetDifferentColors(t *testing.T) {
	is := is.New(t)
	result, err := Highlight([]byte(minimalSVG), []string{"vie", "lon"})
	is.NoErr(err)
	s := string(result)
	// Both provinces should have been changed.
	is.Equal(strings.Contains(s, `inkscape:label="vie"`+"\n   style=\"fill:#000000"), false)
	is.Equal(strings.Contains(s, `inkscape:label="lon"`+"\n   style=\"fill:#000000"), false)
	// And they should be different colours (palette has distinct entries).
	idx1 := strings.Index(s, `inkscape:label="vie"`)
	idx2 := strings.Index(s, `inkscape:label="lon"`)
	seg1 := s[idx1 : idx1+60]
	seg2 := s[idx2 : idx2+60]
	is.Equal(seg1 == seg2, false)
}

func TestHighlight_CaseInsensitiveProvinceMatch(t *testing.T) {
	is := is.New(t)
	// Pass "VIE" (uppercase) — should still match inkscape:label="vie".
	result, err := Highlight([]byte(minimalSVG), []string{"VIE"})
	is.NoErr(err)
	s := string(result)
	is.Equal(strings.Contains(s, `inkscape:label="vie"`+"\n   style=\"fill:#000000"), false)
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
	is.Equal(string(result), minimalSVG)
}

func TestHighlight_PaletteWrapsAround(t *testing.T) {
	is := is.New(t)
	// Request more provinces than palette entries — should not panic.
	provinces := []string{"vie", "lon", "par", "vie", "lon", "par", "vie", "lon", "par"}
	_, err := Highlight([]byte(minimalSVG), provinces)
	is.NoErr(err)
}
