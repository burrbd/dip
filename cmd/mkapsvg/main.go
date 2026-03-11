// Command mkapsvg generates dipmap/assets/map.svg — a copy of the godip
// classical SVG map pre-populated with army and fleet placeholder glyphs for
// every province centre. All glyphs start with display="none" so the board
// appears clean until the bot activates them by ID.
//
// Usage:
//
//	go run ./cmd/mkapsvg/
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/zond/godip/variants/classical"
)

// coastalVariants lists the coast-qualified province names that have a fleet
// position but no army position.
var coastalVariants = map[string]bool{
	"bul/ec": true,
	"bul/sc": true,
	"stp/nc": true,
	"stp/sc": true,
	"spa/nc": true,
	"spa/sc": true,
}

// baseProvinces are the three provinces whose coastal variants are separate
// entries but that also have a base army position.
var baseProvinces = map[string]bool{
	"bul": true,
	"spa": true,
	"stp": true,
}

// exitFn is the function called on fatal error. It is a variable so that tests
// can replace it with a non-exiting stub.
var exitFn = os.Exit

func main() {
	outPath := filepath.Join("dipmap", "assets", "map.svg")
	if err := runWith(classical.Asset, outPath); err != nil {
		fmt.Fprintf(os.Stderr, "mkapsvg: %v\n", err)
		exitFn(1)
	}
}

// runWith is the testable core of main. It loads the godip SVG via assetFn,
// generates the populated map, and writes it to outPath.
func runWith(assetFn func(string) ([]byte, error), outPath string) error {
	raw, err := assetFn("svg/map.svg")
	if err != nil {
		return fmt.Errorf("load godip SVG: %w", err)
	}
	out, fleets, armies, err := generateSVG(raw)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := os.WriteFile(outPath, []byte(out), 0644); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}

	total := fleets + armies
	fmt.Printf("mkapsvg: wrote %s\n", outPath)
	fmt.Printf("  province centres found: %d fleet glyphs, %d army glyphs, %d total placeholders\n",
		fleets, armies, total)
	return nil
}

// stripInkscape removes Inkscape/Sodipodi-specific metadata and attributes
// from an SVG string. The result retains everything the bot render pipeline
// uses (province shapes, label positions, supply-centre markers, highlights
// layer, unit placeholder glyphs) but strips the machine-generated cruft that
// makes the file hard to read or hand-edit.
//
// Removed:
//   - <metadata>…</metadata> blocks (RDF / Dublin Core)
//   - <sodipodi:namedview> elements
//   - <defs>…</defs> blocks (Inkscape arrowhead markers and embedded images)
//   - <style>…</style> blocks (already stripped at render time by stripStyles)
//   - sodipodi:* attributes on any element
//   - inkscape:* attributes on any element except inkscape:label
//   - Namespace declarations for sodipodi, dc, cc, rdf, svg prefixes
//
// Preserved:
//   - xmlns:inkscape (needed for inkscape:label= to be valid XML)
//   - inkscape:label="…" attributes (used by extractProvinceShape in highlight.go)
//
// The function is idempotent: running it twice produces the same output.
func stripInkscape(svg string) string {
	// Remove <metadata>…</metadata> blocks.
	svg = regexp.MustCompile(`(?s)<metadata[^>]*>.*?</metadata\s*>`).ReplaceAllString(svg, "")

	// Remove <sodipodi:namedview> (self-closing or with body).
	svg = regexp.MustCompile(`(?s)<sodipodi:namedview[^>]*/\s*>`).ReplaceAllString(svg, "")
	svg = regexp.MustCompile(`(?s)<sodipodi:namedview[^>]*>.*?</sodipodi:namedview\s*>`).ReplaceAllString(svg, "")

	// Remove <defs>…</defs> blocks (contain only Inkscape markers and embedded
	// images that are not required by the bot's render pipeline).
	svg = regexp.MustCompile(`(?s)<defs[^>]*>.*?</defs\s*>`).ReplaceAllString(svg, "")

	// Remove <style>…</style> blocks.
	svg = regexp.MustCompile(`(?s)<style[^>]*>.*?</style\s*>`).ReplaceAllString(svg, "")

	// Remove sodipodi:* attributes on elements.
	svg = regexp.MustCompile(`\s+sodipodi:[a-zA-Z_-]+="[^"]*"`).ReplaceAllString(svg, "")

	// Remove inkscape:* attributes except inkscape:label. Go's regexp does not
	// support lookaheads, so we use a two-step approach: protect inkscape:label
	// by replacing it with a placeholder, strip all remaining inkscape:* attrs,
	// then restore the placeholder.
	const labelPlaceholder = "INKSCAPE_LABEL_PLACEHOLDER="
	svg = strings.ReplaceAll(svg, "inkscape:label=", labelPlaceholder)
	svg = regexp.MustCompile(`\s+inkscape:[a-zA-Z_-]+="[^"]*"`).ReplaceAllString(svg, "")
	svg = strings.ReplaceAll(svg, labelPlaceholder, "inkscape:label=")

	// Remove namespace declarations for prefixes that are no longer used.
	// xmlns:inkscape is intentionally kept so that inkscape:label= remains
	// valid XML; its presence does not trigger the forbidden-attribute test
	// because "xmlns:inkscape" does not contain the substring "inkscape:".
	svg = regexp.MustCompile(`\s+xmlns:(?:sodipodi|dc|cc|rdf|svg)="[^"]*"`).ReplaceAllString(svg, "")

	// Collapse runs of three or more newlines (blank lines left by removed
	// blocks) to a single blank line.
	svg = regexp.MustCompile(`\n{3,}`).ReplaceAllString(svg, "\n\n")

	return svg
}

// stripProvincesGroup removes the <g id="provinces"> layer from the SVG.
// That layer fills every province shape with solid black (#000000), obscuring
// all other layers when rendered by tdewolff/canvas. Province borders are
// already rendered by the foreground group; the provinces layer is not needed
// by the bot's render pipeline.
func stripProvincesGroup(svg string) string {
	re := regexp.MustCompile(`(?s)<g\b[^>]*\bid="provinces"[^>]*>`)
	loc := re.FindStringIndex(svg)
	if loc == nil {
		return svg
	}
	start := loc[0]
	pos := loc[1]
	depth := 1
	openRe := regexp.MustCompile(`<g\b`)
	closeRe := regexp.MustCompile(`</g>`)
	for depth > 0 && pos < len(svg) {
		openIdx := openRe.FindStringIndex(svg[pos:])
		closeIdx := closeRe.FindStringIndex(svg[pos:])
		switch {
		case openIdx != nil && (closeIdx == nil || openIdx[0] < closeIdx[0]):
			depth++
			pos += openIdx[1]
		case closeIdx != nil:
			depth--
			pos += closeIdx[1]
		default:
			return svg
		}
	}
	return svg[:start] + svg[pos:]
}

// fixSVGDimensions replaces percentage width/height on the SVG root element
// with the numeric values from its viewBox. tdewolff/canvas resolves a
// percentage width/height against a 1 mm parent, producing a 1×1 px output.
// Using explicit pixel values gives a full-resolution raster image.
func fixSVGDimensions(svg string) string {
	vbRe := regexp.MustCompile(`viewBox="0 0 ([0-9.]+) ([0-9.]+)"`)
	m := vbRe.FindStringSubmatch(svg)
	if m == nil {
		return svg
	}
	w, h := m[1], m[2]
	svg = regexp.MustCompile(`\bwidth="[^"]*%[^"]*"`).ReplaceAllString(svg, `width="`+w+`"`)
	svg = regexp.MustCompile(`\bheight="[^"]*%[^"]*"`).ReplaceAllString(svg, `height="`+h+`"`)
	return svg
}

// generateSVG populates the units layer of the godip SVG with placeholder
// glyphs and returns the result along with glyph counts.
func generateSVG(raw []byte) (svg string, fleets, armies int, err error) {
	svg = stripInkscape(string(raw))
	svg = stripProvincesGroup(svg)
	svg = fixSVGDimensions(svg)

	// Find all province centre markers: <path id="<name>Center" d="m cx,cy …"/>
	centers, err := parseProvinceCenters(svg)
	if err != nil {
		return "", 0, 0, err
	}

	// Build glyph markup for each province.
	var glyphs strings.Builder
	for _, c := range centers {
		pid := strings.ReplaceAll(c.province, "/", "-")
		isCoastal := coastalVariants[c.province]

		if !isCoastal {
			// Army glyph for every non-coastal province.
			glyphs.WriteString(armyGlyph(pid, c.cx, c.cy))
			armies++
		}
		// Fleet glyph for every province (including coastal variants).
		glyphs.WriteString(fleetGlyph(pid, c.cx, c.cy))
		fleets++
	}

	// Inject glyphs into the units layer (self-closing <g id="units"/>).
	reUnits := regexp.MustCompile(`<g\s+id="units"\s*/>`)
	replacement := `<g id="units">` + "\n" + glyphs.String() + `</g>`
	if !reUnits.MatchString(svg) {
		// Also try with attributes in different order.
		reUnits = regexp.MustCompile(`<g[^>]*\bid="units"[^>]*/\s*>`)
	}
	svg = reUnits.ReplaceAllString(svg, replacement)

	return svg, fleets, armies, nil
}

// center holds the parsed coordinates for one province centre marker.
type center struct {
	province string
	cx, cy   float64
}

// centerRe matches id="<name>Center" elements and captures the province name
// and the d attribute value.
var centerRe = regexp.MustCompile(`id="([^"]+)Center"[^>]*d="([^"]+)"`)

// centerReAlt matches in the opposite attribute order (d before id).
var centerReAlt = regexp.MustCompile(`d="([^"]+)"[^>]*id="([^"]+)Center"`)

// firstNumberPair extracts the first two numbers from a string (the cx, cy
// translation of the first m sub-path in a province centre path).
var firstNumberPair = regexp.MustCompile(`[-+]?[0-9]*\.?[0-9]+`)

// parseProvinceCenters finds all <path id="<name>Center"> elements in the SVG
// and returns their province names and centre coordinates.
func parseProvinceCenters(svg string) ([]center, error) {
	seen := map[string]bool{}
	var centers []center

	// Try both attribute orderings: id first, then d; and d first, then id.
	addMatches := func(re *regexp.Regexp, nameGroup, dGroup int) {
		for _, m := range re.FindAllStringSubmatch(svg, -1) {
			prov := strings.ToLower(m[nameGroup])
			if seen[prov] {
				continue
			}
			dAttr := m[dGroup]
			nums := firstNumberPair.FindAllString(dAttr, 2)
			if len(nums) < 2 {
				continue
			}
			// firstNumberPair only matches valid float strings; ParseFloat
			// cannot fail here, but we assign to _ to satisfy the compiler.
			cx, _ := strconv.ParseFloat(nums[0], 64)
			cy, _ := strconv.ParseFloat(nums[1], 64)
			seen[prov] = true
			centers = append(centers, center{province: prov, cx: cx, cy: cy})
		}
	}

	addMatches(centerRe, 1, 2)
	addMatches(centerReAlt, 2, 1)

	if len(centers) == 0 {
		return nil, fmt.Errorf("no province centre markers found in SVG")
	}
	return centers, nil
}

// armyGlyph generates the SVG markup for a hidden army placeholder at (cx, cy).
// The fill is placed on the <g> element (not the <rect>) so that Overlay's
// setAttr call on the group propagates the nation colour to the rect via SVG
// inheritance.  The text is pinned to black (#000000) with an explicit fill so
// it does not change colour when the nation fill is applied to the group.
func armyGlyph(pid string, cx, cy float64) string {
	return fmt.Sprintf(
		`<g id="unit-%s-army" transform="translate(%s,%s)" display="none" fill="#cccccc">`+"\n"+
			`  <rect x="-9" y="-9" width="18" height="18" rx="2" stroke="#ffffff" stroke-width="2"/>`+"\n"+
			`  <text x="0" y="5" text-anchor="middle" font-size="11" fill="#000000">A</text>`+"\n"+
			`</g>`+"\n",
		pid, fmtCoord(cx), fmtCoord(cy),
	)
}

// fleetGlyph generates the SVG markup for a hidden fleet placeholder at (cx, cy).
// Fleet glyphs use a wider, shorter rect (2.3:1 aspect ratio) to read clearly
// as a ship hull shape and distinguish them from army squares.
func fleetGlyph(pid string, cx, cy float64) string {
	return fmt.Sprintf(
		`<g id="unit-%s-fleet" transform="translate(%s,%s)" display="none" fill="#cccccc">`+"\n"+
			`  <rect x="-14" y="-6" width="28" height="12" rx="3" stroke="#ffffff" stroke-width="2"/>`+"\n"+
			`  <text x="0" y="5" text-anchor="middle" font-size="8" fill="#000000">F</text>`+"\n"+
			`</g>`+"\n",
		pid, fmtCoord(cx), fmtCoord(cy),
	)
}

// fmtCoord formats a coordinate value, omitting the decimal part if it is zero.
func fmtCoord(v float64) string {
	if v == float64(int(v)) {
		return strconv.Itoa(int(v))
	}
	return strconv.FormatFloat(v, 'f', 2, 64)
}
