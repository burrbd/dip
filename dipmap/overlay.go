package dipmap

import (
	"fmt"
	"strings"
)

// Unit describes a piece on the board: its type ("Army" or "Fleet") and the
// nation that controls it (e.g. "England").
type Unit struct {
	Type   string
	Nation string
}

// nationColor maps the seven classical Diplomacy nations to their standard SVG
// fill colours.
var nationColor = map[string]string{
	"England": "#003399",
	"France":  "#3399CC",
	"Germany": "#666666",
	"Italy":   "#339933",
	"Austria": "#CC0000",
	"Russia":  "#AAAAAA",
	"Turkey":  "#CCCC00",
}

// Overlay injects army and fleet SVG glyphs at the centroid of each named
// province. Province names in units should be lowercase (godip convention).
// Provinces not found in the SVG are silently skipped. An empty units map
// returns the original SVG unchanged.
func Overlay(svg []byte, units map[string]Unit) ([]byte, error) {
	if len(units) == 0 {
		return svg, nil
	}
	s := string(svg)
	var glyphs []string
	for province, unit := range units {
		cx, cy, ok := provinceCenter(s, strings.ToLower(province))
		if !ok {
			continue
		}
		glyphs = append(glyphs, unitGlyph(cx, cy, unit))
	}
	if len(glyphs) == 0 {
		return svg, nil
	}
	layer := "<g id=\"units\">\n" + strings.Join(glyphs, "\n") + "\n</g>"
	result := strings.Replace(s, "</svg>", layer+"\n</svg>", 1)
	return []byte(result), nil
}

// provinceCenter returns the centroid (cx, cy) of the named province's polygon
// by averaging its vertex coordinates. Returns ok=false when the province is
// not found in the SVG or has no parseable coordinate data.
func provinceCenter(svg, province string) (cx, cy float64, ok bool) {
	_, data := extractProvinceShape(svg, province)
	if data == "" {
		return 0, 0, false
	}
	xs, ys := parseCoordinates(data)
	if len(xs) == 0 {
		return 0, 0, false
	}
	for _, x := range xs {
		cx += x
	}
	for _, y := range ys {
		cy += y
	}
	n := float64(len(xs))
	return cx / n, cy / n, true
}

// unitGlyph generates an SVG group containing a filled circle and a letter
// label ("A" for Army, "F" for Fleet) at the given SVG coordinates. The circle
// is filled with the nation's standard colour; unknown nations receive a dark
// grey default.
func unitGlyph(cx, cy float64, u Unit) string {
	fill := nationColor[u.Nation]
	if fill == "" {
		fill = "#333333"
	}
	label := "A"
	if strings.EqualFold(u.Type, "Fleet") {
		label = "F"
	}
	return fmt.Sprintf(
		`<g transform="translate(%.2f,%.2f)"><circle r="25" fill="%s" stroke="white" stroke-width="3"/><text text-anchor="middle" dy="0.35em" font-size="28" font-weight="bold" fill="white">%s</text></g>`,
		cx, cy, fill, label,
	)
}
