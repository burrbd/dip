package dipmap

import (
	"regexp"
	"strings"
)

// highlightPalette is the ordered list of colours assigned to highlighted
// provinces. Colours cycle if more provinces are requested than palette entries.
var highlightPalette = []string{
	"#FF6B6B",
	"#4ECDC4",
	"#45B7D1",
	"#96CEB4",
	"#FFEAA7",
	"#DDA0DD",
	"#98D8C8",
}

// Highlight modifies the fill colour of named provinces in an Inkscape SVG.
// Each province in the list receives a distinct colour from highlightPalette;
// colours cycle if the list is longer than the palette. Province names are
// matched case-insensitively against the inkscape:label attribute. SVG bytes
// that do not contain a matching element are returned unchanged.
func Highlight(svg []byte, provinces []string) ([]byte, error) {
	if len(provinces) == 0 {
		return svg, nil
	}
	result := append([]byte(nil), svg...)
	for i, p := range provinces {
		color := highlightPalette[i%len(highlightPalette)]
		label := regexp.QuoteMeta(strings.ToLower(p))
		// Match inkscape:label="province" followed by whitespace and style="fill:#RRGGBB".
		re := regexp.MustCompile(`(inkscape:label="` + label + `"\s+style="fill:)#[0-9a-fA-F]{3,8}`)
		result = re.ReplaceAll(result, []byte("${1}"+color))
	}
	return result, nil
}
