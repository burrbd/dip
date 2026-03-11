package dipmap

import (
	"fmt"
	"os"
	"regexp"
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

// nationColour returns the SVG fill colour for nation, or a dark-grey default
// for unrecognised nations.
func nationColour(nation string) string {
	if c := nationColor[nation]; c != "" {
		return c
	}
	return "#333333"
}

// Overlay activates the pre-placed unit placeholder glyphs in the SVG for
// each province in units. It sets fill to the nation colour and stroke to
// "#ffffff" on the matching <g id="unit-<province>-<type>"> element. Both
// attributes cascade to the child <rect> via SVG inheritance, making the
// coloured unit box visible. Hidden glyphs carry fill="none" stroke="none"
// so they are transparent; tdewolff/canvas ignores display:none.
//
// Province names in units should use godip's lowercase convention (e.g. "vie",
// "stp/nc"). The "/" separator in coastal names is normalised to "-" to form
// the element ID (e.g. "unit-stp-nc-fleet").
//
// Provinces not found in the SVG are logged to stderr and skipped silently.
// An empty units map returns the original SVG unchanged.
func Overlay(svg []byte, units map[string]Unit) ([]byte, error) {
	if len(units) == 0 {
		return svg, nil
	}
	s := injectUnits(string(svg), units)
	return []byte(s), nil
}

// injectUnits activates pre-placed unit glyphs by setting fill to the nation
// colour and stroke to "#ffffff" on each matching element. Both cascade into
// the child <rect> via SVG inheritance (the rect carries no explicit fill or
// stroke, so it inherits from the <g> parent).
func injectUnits(svg string, units map[string]Unit) string {
	for province, u := range units {
		pid := strings.ReplaceAll(province, "/", "-")
		id := fmt.Sprintf("unit-%s-%s", pid, strings.ToLower(u.Type))
		colour := nationColour(u.Nation)
		svg = setAttr(svg, id, "fill", colour)
		svg = setAttr(svg, id, "stroke", "#ffffff")
	}
	return svg
}

// setAttr finds the element with the given id in svg and updates or inserts
// the named attribute with val. If the element is not found, a warning is
// written to stderr and svg is returned unchanged.
func setAttr(svg, id, attr, val string) string {
	marker := `id="` + id + `"`
	idx := strings.Index(svg, marker)
	if idx < 0 {
		fmt.Fprintf(os.Stderr, "dipmap: setAttr: element id=%q not found\n", id)
		return svg
	}

	// Find the start of the opening tag (last '<' before the id attribute).
	tagStart := strings.LastIndex(svg[:idx], "<")
	if tagStart < 0 {
		return svg
	}

	// Find the end of the opening tag (next '>' after tagStart).
	tagEnd := strings.Index(svg[tagStart:], ">")
	if tagEnd < 0 {
		return svg
	}
	tagEnd = tagStart + tagEnd + 1

	tag := svg[tagStart:tagEnd]

	// If attr already present in the tag, replace its value.
	attrPat := regexp.MustCompile(`\b` + regexp.QuoteMeta(attr) + `="[^"]*"`)
	var newTag string
	if attrPat.MatchString(tag) {
		newTag = attrPat.ReplaceAllString(tag, attr+`="`+val+`"`)
	} else {
		// Insert attr before the closing '>' or '/>'.
		if strings.HasSuffix(tag, "/>") {
			newTag = tag[:len(tag)-2] + ` ` + attr + `="` + val + `"/>`
		} else {
			newTag = tag[:len(tag)-1] + ` ` + attr + `="` + val + `">`
		}
	}

	return svg[:tagStart] + newTag + svg[tagEnd:]
}
