package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/cheekybits/is"
	"github.com/zond/godip/variants/classical"
)

// ---- parseProvinceCenters ---------------------------------------------------

func TestParseProvinceCenters_ReturnsExpectedCount(t *testing.T) {
	is := is.New(t)
	raw, err := classical.Asset("svg/map.svg")
	is.NoErr(err)

	centers, err := parseProvinceCenters(string(raw))
	is.NoErr(err)
	// The godip classical map has exactly 81 province centre markers
	// (75 base provinces + 6 coastal variants).
	is.Equal(len(centers), 81)
}

func TestParseProvinceCenters_CoordinatesNonZero(t *testing.T) {
	is := is.New(t)
	raw, err := classical.Asset("svg/map.svg")
	is.NoErr(err)

	centers, err := parseProvinceCenters(string(raw))
	is.NoErr(err)
	for _, c := range centers {
		if c.cx == 0 && c.cy == 0 {
			t.Errorf("province %q has zero coordinates", c.province)
		}
	}
}

func TestParseProvinceCenters_NoCentersFound_ReturnsError(t *testing.T) {
	_, err := parseProvinceCenters(`<svg><rect/></svg>`)
	if err == nil {
		t.Error("expected error when no province centres are found")
	}
}

// ---- generateSVG -----------------------------------------------------------

func TestGenerateSVG_GlyphCounts(t *testing.T) {
	is := is.New(t)
	raw, err := classical.Asset("svg/map.svg")
	is.NoErr(err)

	_, fleets, armies, err := generateSVG(raw)
	is.NoErr(err)

	// 81 province centres → 81 fleet glyphs.
	is.Equal(fleets, 81)
	// 75 base provinces (81 − 6 coastal variants) → 75 army glyphs.
	is.Equal(armies, 75)
}

func TestGenerateSVG_AllGlyphsHideByDefault(t *testing.T) {
	is := is.New(t)
	raw, err := classical.Asset("svg/map.svg")
	is.NoErr(err)

	svg, _, _, err := generateSVG(raw)
	is.NoErr(err)

	// All unit glyphs must use fill="none" stroke="none" to hide them.
	// tdewolff/canvas ignores display:none, so we use transparent fill/stroke.
	if !strings.Contains(svg, `fill="none" stroke="none"`) {
		t.Error("expected unit glyphs to be hidden via fill=none stroke=none")
	}
	// No glyph should have display="none" (canvas ignores it).
	reDisplayNone := regexp.MustCompile(`id="unit-[^"]*"[^>]*display="none"`)
	if reDisplayNone.MatchString(svg) {
		t.Error("unexpected display=none on unit glyph; use fill=none stroke=none instead")
	}
}

func TestGenerateSVG_CoastalVariantsFleetOnly(t *testing.T) {
	is := is.New(t)
	raw, err := classical.Asset("svg/map.svg")
	is.NoErr(err)

	svg, _, _, err := generateSVG(raw)
	is.NoErr(err)

	for cv := range coastalVariants {
		pid := strings.ReplaceAll(cv, "/", "-")
		armyID := `id="unit-` + pid + `-army"`
		fleetID := `id="unit-` + pid + `-fleet"`

		if strings.Contains(svg, armyID) {
			t.Errorf("coastal variant %q should not have an army glyph", cv)
		}
		if !strings.Contains(svg, fleetID) {
			t.Errorf("coastal variant %q missing fleet glyph", cv)
		}
		is.NotNil(svg)
	}
}

func TestGenerateSVG_BaseProvincesHaveBothGlyphs(t *testing.T) {
	is := is.New(t)
	raw, err := classical.Asset("svg/map.svg")
	is.NoErr(err)

	svg, _, _, err := generateSVG(raw)
	is.NoErr(err)

	for base := range baseProvinces {
		armyID := `id="unit-` + base + `-army"`
		fleetID := `id="unit-` + base + `-fleet"`

		if !strings.Contains(svg, armyID) {
			t.Errorf("base province %q missing army glyph", base)
		}
		if !strings.Contains(svg, fleetID) {
			t.Errorf("base province %q missing fleet glyph", base)
		}
		is.NotNil(svg)
	}
}

// ---- main -------------------------------------------------------------------

func TestMain_WritesToDefaultPath(t *testing.T) {
	// main() writes to a relative path "dipmap/assets/map.svg". Change to a
	// temp directory so the file is written there and not into the source tree.
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })

	main()

	outPath := filepath.Join(dir, "dipmap", "assets", "map.svg")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("expected main() to write %s: %v", outPath, err)
	}
	if !strings.Contains(string(data), `id="units"`) {
		t.Error("expected units layer in generated SVG")
	}
}

func TestMain_ErrorCallsExit(t *testing.T) {
	// Replace exitFn so os.Exit(1) doesn't kill the test process.
	var gotCode int
	orig := exitFn
	exitFn = func(code int) { gotCode = code }
	defer func() { exitFn = orig }()

	// Change to a temp dir where writing will fail (read-only root).
	// Instead, inject a failing assetFn by temporarily swapping classic.Asset
	// is not possible, so we change directory to a path where mkdir will fail.
	origWD, _ := os.Getwd()
	dir := t.TempDir()
	// Create a file where the output dir should be.
	blocker := filepath.Join(dir, "dipmap")
	if err := os.WriteFile(blocker, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWD)

	main()

	if gotCode != 1 {
		t.Errorf("expected exit code 1, got %d", gotCode)
	}
}

// ---- runWith ----------------------------------------------------------------

func TestRunWith_WritesFileToOutputPath(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	outPath := filepath.Join(dir, "sub", "map.svg")
	err := runWith(classical.Asset, outPath)
	is.NoErr(err)
	data, err := os.ReadFile(outPath)
	is.NoErr(err)
	if !strings.Contains(string(data), `id="units"`) {
		t.Error("expected units layer in generated SVG")
	}
}

func TestRunWith_AssetError_ReturnsError(t *testing.T) {
	failAsset := func(string) ([]byte, error) { return nil, os.ErrNotExist }
	err := runWith(failAsset, "/tmp/unused.svg")
	if err == nil {
		t.Error("expected error when asset loading fails")
	}
}

func TestRunWith_GenerateSVGError_ReturnsError(t *testing.T) {
	// Return an SVG with no province centres so generateSVG returns an error.
	emptyAsset := func(string) ([]byte, error) { return []byte(`<svg></svg>`), nil }
	err := runWith(emptyAsset, "/tmp/unused.svg")
	if err == nil {
		t.Error("expected error when SVG has no province centres")
	}
}

func TestRunWith_MkdirError_ReturnsError(t *testing.T) {
	// Create a file where the output directory should be so MkdirAll fails.
	dir := t.TempDir()
	blocker := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	err := runWith(classical.Asset, filepath.Join(blocker, "map.svg"))
	if err == nil {
		t.Error("expected error when output directory cannot be created")
	}
}

func TestRunWith_WriteFileError_ReturnsError(t *testing.T) {
	// Create a directory where the output file should be so WriteFile fails.
	dir := t.TempDir()
	outDir := filepath.Join(dir, "out")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatal(err)
	}
	// Create a directory named "map.svg" so os.WriteFile fails (it's a dir, not a file).
	mapDir := filepath.Join(outDir, "map.svg")
	if err := os.MkdirAll(mapDir, 0755); err != nil {
		t.Fatal(err)
	}
	err := runWith(classical.Asset, mapDir)
	if err == nil {
		t.Error("expected error when output path is a directory")
	}
}

// ---- generateSVG (alt regex branch) ----------------------------------------

func TestGenerateSVG_AltRegex_UnitsLayerHasExtraAttributes(t *testing.T) {
	is := is.New(t)
	// SVG where the units layer has attributes before id= so the primary regex
	// (<g\s+id="units"\s*/>) does not match; the fallback must be used.
	svg := `<svg><path id="vieCenter" d="m 100.00,200.00 c 1,2 z"/><g style="fill:none" id="units"/></svg>`
	result, fleets, armies, err := generateSVG([]byte(svg))
	is.NoErr(err)
	is.Equal(fleets, 1)
	is.Equal(armies, 1)
	if !strings.Contains(result, `id="unit-vie-army"`) {
		t.Error("expected unit-vie-army glyph in alt-regex output")
	}
}

// ---- parseProvinceCenters (edge cases) -------------------------------------

func TestParseProvinceCenters_DuplicateProvince_CountedOnce(t *testing.T) {
	is := is.New(t)
	// Same province appears twice (once in each attribute order) — seen check
	// must prevent duplicates.
	svg := `<svg>
  <path id="vieCenter" d="m 100,200 c 1,2 z"/>
  <path d="m 100,200 c 1,2 z" id="vieCenter"/>
</svg>`
	centers, err := parseProvinceCenters(svg)
	is.NoErr(err)
	is.Equal(len(centers), 1)
}

func TestParseProvinceCenters_InsufficientNumbers_Skipped(t *testing.T) {
	is := is.New(t)
	// A center whose d attribute contains fewer than two numbers is skipped.
	svg := `<svg>
  <path id="vieCenter" d="m z"/>
  <path id="lonCenter" d="m 50,60 c 1,2 z"/>
</svg>`
	centers, err := parseProvinceCenters(svg)
	is.NoErr(err)
	// Only "lon" should be present (vie was skipped).
	is.Equal(len(centers), 1)
	is.Equal(centers[0].province, "lon")
}

// ---- fmtCoord ---------------------------------------------------------------

func TestFmtCoord_Integer_NoDecimal(t *testing.T) {
	is := is.New(t)
	is.Equal(fmtCoord(100.0), "100")
	is.Equal(fmtCoord(0.0), "0")
}

func TestFmtCoord_Fractional_TwoDecimalPlaces(t *testing.T) {
	is := is.New(t)
	is.Equal(fmtCoord(748.83), "748.83")
	is.Equal(fmtCoord(12.50), "12.50")
}

// ---- stripInkscape ----------------------------------------------------------

func TestStripInkscape_RemovesMetadata(t *testing.T) {
	svg := `<svg><metadata id="m"><rdf:RDF/></metadata><rect id="r"/></svg>`
	result := stripInkscape(svg)
	if strings.Contains(result, "<metadata") {
		t.Error("expected <metadata> block to be removed")
	}
	if !strings.Contains(result, `id="r"`) {
		t.Error("expected non-metadata element to be preserved")
	}
}

func TestStripInkscape_RemovesSodipodiNamedview(t *testing.T) {
	svg := `<svg><sodipodi:namedview id="nv" units="px"/><rect id="r"/></svg>`
	result := stripInkscape(svg)
	if strings.Contains(result, "sodipodi:namedview") {
		t.Error("expected sodipodi:namedview to be removed")
	}
	if !strings.Contains(result, `id="r"`) {
		t.Error("expected other elements to be preserved")
	}
}

func TestStripInkscape_RemovesDefs(t *testing.T) {
	svg := `<svg><defs><marker id="m"/></defs><rect id="r"/></svg>`
	result := stripInkscape(svg)
	if strings.Contains(result, "<defs") {
		t.Error("expected <defs> block to be removed")
	}
	if !strings.Contains(result, `id="r"`) {
		t.Error("expected non-defs element to be preserved")
	}
}

func TestStripInkscape_RemovesStyleBlocks(t *testing.T) {
	svg := `<svg><style>body{font-family:X}</style><rect id="r"/></svg>`
	result := stripInkscape(svg)
	if strings.Contains(result, "<style") {
		t.Error("expected <style> block to be removed")
	}
}

func TestStripInkscape_RemovesSodipodiAttributes(t *testing.T) {
	svg := `<svg><path d="M 0,0" sodipodi:nodetypes="cc" id="p"/></svg>`
	result := stripInkscape(svg)
	if strings.Contains(result, "sodipodi:") {
		t.Errorf("expected sodipodi: attributes to be removed, got: %q", result)
	}
	if !strings.Contains(result, `id="p"`) {
		t.Error("expected path to be preserved")
	}
}

func TestStripInkscape_RemovesInkscapeAttributesExceptLabel(t *testing.T) {
	svg := `<svg><g inkscape:label="foo" inkscape:groupmode="layer" id="g1"/></svg>`
	result := stripInkscape(svg)
	if !strings.Contains(result, `inkscape:label="foo"`) {
		t.Error("expected inkscape:label to be preserved")
	}
	if strings.Contains(result, `inkscape:groupmode`) {
		t.Error("expected inkscape:groupmode to be removed")
	}
}

func TestStripInkscape_RemovesNamespaceDeclarations(t *testing.T) {
	svg := `<svg xmlns:sodipodi="http://s.f.net" xmlns:dc="http://dc" xmlns:inkscape="http://i.o">` +
		`<rect id="r"/></svg>`
	result := stripInkscape(svg)
	if strings.Contains(result, "xmlns:sodipodi") {
		t.Error("expected xmlns:sodipodi to be removed")
	}
	if strings.Contains(result, "xmlns:dc") {
		t.Error("expected xmlns:dc to be removed")
	}
	// xmlns:inkscape must be kept (needed for inkscape:label= XML validity).
	if !strings.Contains(result, "xmlns:inkscape") {
		t.Error("expected xmlns:inkscape to be preserved")
	}
}

func TestStripInkscape_Idempotent(t *testing.T) {
	is := is.New(t)
	raw, err := classical.Asset("svg/map.svg")
	is.NoErr(err)
	once := stripInkscape(string(raw))
	twice := stripInkscape(once)
	is.Equal(once, twice)
}

// TestStripInkscape_GeneratedSVGHasNoForbiddenAttributes asserts that the SVG
// produced by generateSVG contains no sodipodi: or inkscape: text other than
// inkscape:label= (Story 10a acceptance criterion).
func TestStripInkscape_GeneratedSVGHasNoForbiddenAttributes(t *testing.T) {
	is := is.New(t)
	raw, err := classical.Asset("svg/map.svg")
	is.NoErr(err)

	svg, _, _, err := generateSVG(raw)
	is.NoErr(err)

	// No sodipodi: attributes must remain.
	if strings.Contains(svg, "sodipodi:") {
		t.Error("expected no sodipodi: in generated SVG")
	}

	// No inkscape: attributes except inkscape:label= must remain.
	// Replace all inkscape:label= occurrences and check nothing else remains.
	svgNoLabel := strings.ReplaceAll(svg, "inkscape:label=", "PLACEHOLDER")
	if strings.Contains(svgNoLabel, "inkscape:") {
		t.Error("expected no inkscape: attributes other than inkscape:label= in generated SVG")
	}
}

// ---- glyph geometry ---------------------------------------------------------

// TestGlyphGeometry_Dimensions asserts the army and fleet rect sizes required
// by Story 10b Bug 3: army width < 20, fleet width >= 24, fleet height < army height.
func TestGlyphGeometry_Dimensions(t *testing.T) {
	armySVG := armyGlyph("test", 0, 0)
	fleetSVG := fleetGlyph("test", 0, 0)

	widthRe := regexp.MustCompile(`width="([0-9]+)"`)
	heightRe := regexp.MustCompile(`height="([0-9]+)"`)

	armyWm := widthRe.FindStringSubmatch(armySVG)
	fleetWm := widthRe.FindStringSubmatch(fleetSVG)
	armyHm := heightRe.FindStringSubmatch(armySVG)
	fleetHm := heightRe.FindStringSubmatch(fleetSVG)

	if armyWm == nil || fleetWm == nil || armyHm == nil || fleetHm == nil {
		t.Fatal("could not extract rect dimensions from glyph SVG")
	}

	armyW, _ := strconv.Atoi(armyWm[1])
	fleetW, _ := strconv.Atoi(fleetWm[1])
	armyH, _ := strconv.Atoi(armyHm[1])
	fleetH, _ := strconv.Atoi(fleetHm[1])

	if armyW >= 20 {
		t.Errorf("army rect width %d must be < 20", armyW)
	}
	if fleetW < 24 {
		t.Errorf("fleet rect width %d must be >= 24", fleetW)
	}
	if fleetH >= armyH {
		t.Errorf("fleet rect height %d must be < army rect height %d", fleetH, armyH)
	}
}

// TestGlyphGeometry_FillOnGroup asserts that the <g> element carries
// fill="none" stroke="none" (the hidden state) and that the inner <rect>
// carries no explicit fill or stroke (so it inherits from the group). When
// Overlay activates a glyph it sets fill=nationColor stroke=#ffffff on the
// group, both of which cascade into the rect via SVG inheritance.
func TestGlyphGeometry_FillOnGroup(t *testing.T) {
	for _, tc := range []struct {
		name string
		svg  string
	}{
		{"army", armyGlyph("test", 0, 0)},
		{"fleet", fleetGlyph("test", 0, 0)},
	} {
		// The <g> opening tag must carry fill="none" stroke="none".
		gRe := regexp.MustCompile(`<g\b[^>]*>`)
		gTag := gRe.FindString(tc.svg)
		if !strings.Contains(gTag, `fill="none"`) {
			t.Errorf("%s: expected fill=none on <g> tag, got: %q", tc.name, gTag)
		}
		if !strings.Contains(gTag, `stroke="none"`) {
			t.Errorf("%s: expected stroke=none on <g> tag, got: %q", tc.name, gTag)
		}
		// The <rect> must NOT carry its own fill or stroke attributes (so it inherits).
		rectRe := regexp.MustCompile(`<rect\b[^>]*/?>`)
		rectTag := rectRe.FindString(tc.svg)
		if strings.Contains(rectTag, "fill=") {
			t.Errorf("%s: expected no fill on <rect>, got: %q", tc.name, rectTag)
		}
		if strings.Contains(rectTag, `stroke="`) {
			t.Errorf("%s: expected no stroke color on <rect>, got: %q", tc.name, rectTag)
		}
	}
}

// ---- stripGroupByID ---------------------------------------------------------

func TestStripGroupByID_RemovesTargetGroup(t *testing.T) {
	svg := `<svg><g id="provinces"><path id="p"/></g><rect id="r"/></svg>`
	result := stripGroupByID(svg, "provinces")
	if strings.Contains(result, `id="provinces"`) {
		t.Error("expected provinces group to be removed")
	}
	if !strings.Contains(result, `id="r"`) {
		t.Error("expected sibling element to be preserved")
	}
}

func TestStripGroupByID_GroupNotFound_ReturnsUnchanged(t *testing.T) {
	svg := `<svg><g id="other"><path/></g></svg>`
	result := stripGroupByID(svg, "provinces")
	if result != svg {
		t.Error("expected SVG unchanged when group not found")
	}
}

func TestStripGroupByID_NestedGroups_RemovesCorrectly(t *testing.T) {
	svg := `<svg><g id="provinces"><g id="inner"><path/></g></g><g id="keep"/></svg>`
	result := stripGroupByID(svg, "provinces")
	if strings.Contains(result, `id="provinces"`) || strings.Contains(result, `id="inner"`) {
		t.Error("expected provinces group and all nested content to be removed")
	}
	if !strings.Contains(result, `id="keep"`) {
		t.Error("expected sibling group to be preserved")
	}
}

func TestStripGroupByID_UnterminatedGroup_ReturnsUnchanged(t *testing.T) {
	// Group is opened but never closed (no </g>). The fallback default branch
	// in the depth-tracking loop must return the original svg.
	svg := `<svg><g id="provinces"><path id="p"/>`
	result := stripGroupByID(svg, "provinces")
	if result != svg {
		t.Error("expected SVG unchanged when group has no closing tag")
	}
}

func TestGenerateSVG_StripsHiddenGroups(t *testing.T) {
	is := is.New(t)
	raw, err := classical.Asset("svg/map.svg")
	is.NoErr(err)

	svg, _, _, err := generateSVG(raw)
	is.NoErr(err)

	for _, id := range []string{"provinces", "supply-centers", "province-centers"} {
		if strings.Contains(svg, `id="`+id+`"`) {
			t.Errorf("expected group id=%q to be stripped from generated SVG", id)
		}
	}
}

// ---- fixImpassableFill ------------------------------------------------------

func TestFixImpassableFill_ReplacesImpassableStripes(t *testing.T) {
	svg := `<path style="fill:url(#impassableStripes);stroke:#000"/>`
	result := fixImpassableFill(svg)
	if strings.Contains(result, "impassableStripes") {
		t.Error("expected impassableStripes reference to be replaced")
	}
	if !strings.Contains(result, "fill:#d4d0ad") {
		t.Error("expected fill to be replaced with background beige #d4d0ad")
	}
}

func TestFixImpassableFill_ReplacesPattern1827(t *testing.T) {
	svg := `<rect style="fill:url(#pattern1827);fill-opacity:0.05"/>`
	result := fixImpassableFill(svg)
	if strings.Contains(result, "pattern1827") {
		t.Error("expected pattern1827 reference to be replaced")
	}
	if !strings.Contains(result, "fill:none") {
		t.Error("expected fill:url(#pattern1827) to become fill:none")
	}
}

func TestFixImpassableFill_NoMatch_ReturnsUnchanged(t *testing.T) {
	svg := `<rect style="fill:#d4d0ad"/>`
	result := fixImpassableFill(svg)
	if result != svg {
		t.Error("expected SVG unchanged when no url() fill references present")
	}
}

func TestGenerateSVG_NoImpassableStripes(t *testing.T) {
	is := is.New(t)
	raw, err := classical.Asset("svg/map.svg")
	is.NoErr(err)

	svg, _, _, err := generateSVG(raw)
	is.NoErr(err)

	if strings.Contains(svg, "impassableStripes") {
		t.Error("expected impassableStripes url() reference to be replaced in generated SVG")
	}
}

// ---- flattenTextTspan -------------------------------------------------------

func TestFlattenTextTspan_FlattensSimpleCase(t *testing.T) {
	svg := `<text transform="rotate(-8)" id="Paris" x="-10" y="-20">` +
		`<tspan x="400" y="500">Paris</tspan></text>`
	result := flattenTextTspan(svg)
	if strings.Contains(result, "<tspan") {
		t.Error("expected <tspan> to be removed")
	}
	if !strings.Contains(result, `x="400"`) || !strings.Contains(result, `y="500"`) {
		t.Errorf("expected tspan coordinates in result, got: %s", result)
	}
	if !strings.Contains(result, `transform="rotate(-8)"`) {
		t.Error("expected transform from <text> to be preserved")
	}
	if !strings.Contains(result, ">Paris<") {
		t.Error("expected text content to be preserved")
	}
}

func TestFlattenTextTspan_MultilineAttributes(t *testing.T) {
	// tspan opening tag spans multiple lines (as Inkscape emits it).
	svg := "<text id=\"foo\" x=\"-89\" y=\"79\"><tspan\n   id=\"ts1\"\n   y=\"1059\"\n   x=\"287\">Berlin</tspan></text>"
	result := flattenTextTspan(svg)
	if strings.Contains(result, "<tspan") {
		t.Error("expected <tspan> to be removed")
	}
	if !strings.Contains(result, `x="287"`) || !strings.Contains(result, `y="1059"`) {
		t.Errorf("expected tspan coordinates used, got: %s", result)
	}
	if !strings.Contains(result, "Berlin") {
		t.Error("expected text content preserved")
	}
}

func TestFlattenTextTspan_NoTspan_ReturnsUnchanged(t *testing.T) {
	// Text element with direct content (no tspan) — must be left alone.
	svg := `<text x="10" y="20">London</text>`
	result := flattenTextTspan(svg)
	if result != svg {
		t.Errorf("expected unchanged, got: %s", result)
	}
}

func TestFlattenTextTspan_MultipleTspans_JoinsContent(t *testing.T) {
	// Two-tspan label like "ENGLISH CHANNEL" — text joined, first tspan coords used.
	svg := `<text transform="rotate(-3)" id="EC" x="10" y="38">` +
		`<tspan x="222" y="781">ENGLISH</tspan>` +
		`<tspan x="232" y="801">CHANNEL</tspan></text>`
	result := flattenTextTspan(svg)
	if strings.Contains(result, "<tspan") {
		t.Error("expected all tspan elements removed")
	}
	if !strings.Contains(result, "ENGLISH CHANNEL") {
		t.Errorf("expected joined text, got: %s", result)
	}
	if !strings.Contains(result, `x="222"`) || !strings.Contains(result, `y="781"`) {
		t.Errorf("expected first tspan coords, got: %s", result)
	}
}

func TestFlattenTextTspan_TspanMissingCoords_ReturnsUnchanged(t *testing.T) {
	// tspan has no x/y — cannot determine text position, leave as is.
	svg := `<text x="10" y="20"><tspan id="ts">Moscow</tspan></text>`
	result := flattenTextTspan(svg)
	if result != svg {
		t.Errorf("expected unchanged when tspan has no x/y, got: %s", result)
	}
}

func TestFlattenTextTspan_TextWithoutXY_InsertsFromTspan(t *testing.T) {
	// <text> has no x/y — the else branches insert them from the first tspan.
	svg := `<text transform="rotate(-5)" id="vie"><tspan x="600" y="700">Vienna</tspan></text>`
	result := flattenTextTspan(svg)
	if strings.Contains(result, "<tspan") {
		t.Error("expected <tspan> to be removed")
	}
	if !strings.Contains(result, `x="600"`) || !strings.Contains(result, `y="700"`) {
		t.Errorf("expected tspan coordinates inserted, got: %s", result)
	}
	if !strings.Contains(result, "Vienna") {
		t.Error("expected text content preserved")
	}
}

func TestFlattenTextTspan_GeneratedSVGHasNoTspan(t *testing.T) {
	is := is.New(t)
	raw, err := classical.Asset("svg/map.svg")
	is.NoErr(err)

	svg, _, _, err := generateSVG(raw)
	is.NoErr(err)

	if strings.Contains(svg, "<tspan") {
		t.Error("expected no <tspan> elements in generated SVG after flattening")
	}
}
