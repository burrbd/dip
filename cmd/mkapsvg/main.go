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

// generateSVG populates the units layer of the godip SVG with placeholder
// glyphs and returns the result along with glyph counts.
func generateSVG(raw []byte) (svg string, fleets, armies int, err error) {
	svg = string(raw)

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
func armyGlyph(pid string, cx, cy float64) string {
	return fmt.Sprintf(
		`<g id="unit-%s-army" transform="translate(%s,%s)" display="none">`+"\n"+
			`  <rect x="-12" y="-12" width="24" height="24" rx="3" fill="#cccccc" stroke="#ffffff" stroke-width="2"/>`+"\n"+
			`  <text x="0" y="5" text-anchor="middle" font-size="14" fill="#000000">A</text>`+"\n"+
			`</g>`+"\n",
		pid, fmtCoord(cx), fmtCoord(cy),
	)
}

// fleetGlyph generates the SVG markup for a hidden fleet placeholder at (cx, cy).
func fleetGlyph(pid string, cx, cy float64) string {
	return fmt.Sprintf(
		`<g id="unit-%s-fleet" transform="translate(%s,%s)" display="none">`+"\n"+
			`  <rect x="-15" y="-9" width="30" height="18" rx="3" fill="#cccccc" stroke="#ffffff" stroke-width="2"/>`+"\n"+
			`  <text x="0" y="5" text-anchor="middle" font-size="10" fill="#000000">F</text>`+"\n"+
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
