package dipmap

import (
	"fmt"
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

// Highlight injects coloured copies of named provinces into the SVG's
// "highlights" layer (id="highlights"). Each province in the list receives a
// distinct colour from highlightPalette; colours cycle when the list is longer
// than the palette. Province names are matched case-insensitively against
// inkscape:label attributes. Provinces not found in the SVG are silently
// skipped. The original province elements are left unchanged.
func Highlight(svg []byte, provinces []string) ([]byte, error) {
	if len(provinces) == 0 {
		return svg, nil
	}
	s := string(svg)
	var elements []string
	for i, p := range provinces {
		color := highlightPalette[i%len(highlightPalette)]
		elemType, data := extractProvinceShape(s, strings.ToLower(p))
		if elemType == "" {
			continue
		}
		var elem string
		switch elemType {
		case "polygon":
			elem = fmt.Sprintf(`<polygon style="fill:%s;fill-opacity:0.7" points="%s"/>`, color, data)
		case "path":
			elem = fmt.Sprintf(`<path style="fill:%s;fill-opacity:0.7" d="%s"/>`, color, data)
		}
		elements = append(elements, elem)
	}
	if len(elements) == 0 {
		return svg, nil
	}
	// Replace the self-closing highlights group with an open group that
	// contains the coloured province copies.
	injection := strings.Join(elements, "\n")
	reHighlights := regexp.MustCompile(`(id="highlights"[^>]*)/>`)
	result := reHighlights.ReplaceAllString(s, "${1}>\n"+injection+"\n</g>")
	return []byte(result), nil
}

// extractProvinceShape finds the polygon or path element in the SVG whose
// inkscape:label matches label (already lower-cased). It returns the element
// type ("polygon" or "path") and the shape data (points or d attribute value).
func extractProvinceShape(svg, label string) (elemType, data string) {
	searchStr := `inkscape:label="` + label + `"`
	idx := strings.Index(svg, searchStr)
	if idx < 0 {
		return "", ""
	}
	// Walk backward to find the opening '<' of this element.
	start := strings.LastIndex(svg[:idx], "<")
	if start < 0 {
		return "", ""
	}
	// Walk forward to find the self-closing '/>'.
	end := strings.Index(svg[start:], "/>")
	if end < 0 {
		return "", ""
	}
	element := svg[start : start+end+2]

	rest := strings.TrimLeft(element[1:], " \t\n\r")
	switch {
	case strings.HasPrefix(rest, "polygon"):
		re := regexp.MustCompile(`points="([^"]+)"`)
		m := re.FindStringSubmatch(element)
		if len(m) > 1 {
			return "polygon", m[1]
		}
	case strings.HasPrefix(rest, "path"):
		re := regexp.MustCompile(`\bd="([^"]+)"`)
		m := re.FindStringSubmatch(element)
		if len(m) > 1 {
			return "path", m[1]
		}
	}
	return "", ""
}
