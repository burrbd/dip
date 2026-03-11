package main

import (
	"os"
	"path/filepath"
	"regexp"
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

	// All unit placeholders must start with display="none".
	reVisible := regexp.MustCompile(`id="unit-[^"]*"[^>]*display="inline"`)
	if reVisible.MatchString(svg) {
		t.Error("expected all unit glyphs to have display=none, but found display=inline")
	}
	// At least one glyph with display="none" must be present.
	if !strings.Contains(svg, `display="none"`) {
		t.Error("expected at least one display=none in generated SVG")
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
