package dipmap

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/cheekybits/is"
)

// failEncode is a png-encode stub that always returns an error.
func failEncode(_ io.Writer, _ image.Image) error { return errors.New("encode fail") }

// stubEngineState implements EngineState for test use.
type stubEngineState struct{}

func (s stubEngineState) Dump() ([]byte, error) { return []byte(`{}`), nil }

// assertPNG checks that result begins with the PNG magic bytes (\x89PNG).
func assertPNG(t *testing.T, result []byte) {
	t.Helper()
	if len(result) < 4 {
		t.Fatalf("expected PNG bytes, got %d bytes", len(result))
	}
	is := is.New(t)
	is.Equal(result[0], byte(0x89))
	is.Equal(result[1], byte('P'))
	is.Equal(result[2], byte('N'))
	is.Equal(result[3], byte('G'))
}

// ---- SVGToPNG ---------------------------------------------------------------

func TestSVGToPNG_MagicBytes(t *testing.T) {
	is := is.New(t)
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"/>`)
	result, err := SVGToPNG(svg)
	is.NoErr(err)
	assertPNG(t, result)
}

func TestSVGToPNG_MalformedSVG_ReturnsError(t *testing.T) {
	is := is.New(t)
	_, err := SVGToPNG([]byte(`<svg><unclosed`))
	is.NotNil(err)
}

// ---- Render -----------------------------------------------------------------

func TestRender_ReturnsPNGBytes(t *testing.T) {
	is := is.New(t)
	result, err := Render(stubEngineState{})
	is.NoErr(err)
	assertPNG(t, result)
}

// TestRender_IntegrationMapSize verifies the real embedded map produces a
// substantial PNG (acceptance criterion: >50 KB).
func TestRender_IntegrationMapSize(t *testing.T) {
	is := is.New(t)
	result, err := Render(stubEngineState{})
	is.NoErr(err)
	const minBytes = 50 * 1024
	if len(result) < minBytes {
		t.Errorf("expected PNG > %d bytes, got %d", minBytes, len(result))
	}
}

// ---- LoadSVG ----------------------------------------------------------------

func TestLoadSVG_ReturnsSVGBytes(t *testing.T) {
	is := is.New(t)
	result, err := LoadSVG(stubEngineState{})
	is.NoErr(err)
	is.NotNil(result)
	if len(result) < 4 {
		t.Fatalf("expected SVG bytes, got %d bytes", len(result))
	}
}

// ---- parseSVGViewBox --------------------------------------------------------

func TestParseSVGViewBox_WithValidViewBox_ReturnsValues(t *testing.T) {
	is := is.New(t)
	svg := `<svg viewBox="10 20 300 200">`
	x, y, w, h := parseSVGViewBox(svg)
	is.Equal(x, 10.0)
	is.Equal(y, 20.0)
	is.Equal(w, 300.0)
	is.Equal(h, 200.0)
}

func TestParseSVGViewBox_NoViewBox_ReturnsZeros(t *testing.T) {
	is := is.New(t)
	x, y, w, h := parseSVGViewBox(`<svg width="100" height="100">`)
	is.Equal(x, 0.0)
	is.Equal(y, 0.0)
	is.Equal(w, 0.0)
	is.Equal(h, 0.0)
}

func TestParseSVGViewBox_MalformedViewBox_ReturnsZeros(t *testing.T) {
	is := is.New(t)
	x, y, w, h := parseSVGViewBox(`<svg viewBox="0 0 100">`)
	is.Equal(x, 0.0)
	is.Equal(y, 0.0)
	is.Equal(w, 0.0)
	is.Equal(h, 0.0)
}

// ---- parseSVGDimensions -----------------------------------------------------

func TestParseSVGDimensions_WithWidthAndHeight_ReturnsValues(t *testing.T) {
	is := is.New(t)
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="640" height="480">`)
	w, h := parseSVGDimensions(svg)
	is.Equal(w, 640)
	is.Equal(h, 480)
}

func TestParseSVGDimensions_WithPercentages_ReturnsZeros(t *testing.T) {
	is := is.New(t)
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="100%" height="100%">`)
	w, h := parseSVGDimensions(svg)
	is.Equal(w, 0)
	is.Equal(h, 0)
}

// ---- parseCoordinates -------------------------------------------------------

func TestParseCoordinates_PolygonPoints_ReturnsPairs(t *testing.T) {
	is := is.New(t)
	data := "10,20 30,40 50,60"
	xs, ys := parseCoordinates(data)
	is.Equal(len(xs), 3)
	is.Equal(xs[0], 10.0)
	is.Equal(ys[0], 20.0)
	is.Equal(xs[2], 50.0)
	is.Equal(ys[2], 60.0)
}

func TestParseCoordinates_EmptyString_ReturnsNil(t *testing.T) {
	is := is.New(t)
	xs, ys := parseCoordinates("")
	is.Equal(len(xs), 0)
	is.Equal(len(ys), 0)
}

// ---- rewriteViewBox ---------------------------------------------------------

func TestRewriteViewBox_ReplacesExistingViewBox(t *testing.T) {
	svg := `<svg viewBox="0 0 100 100">`
	newVB := `viewBox="10 20 50 50"`
	result := rewriteViewBox(svg, newVB)
	if !strings.Contains(result, newVB) {
		t.Errorf("expected new viewBox in result, got %q", result)
	}
	if strings.Contains(result, `viewBox="0 0 100 100"`) {
		t.Error("expected old viewBox to be replaced")
	}
}

func TestRewriteViewBox_InsertsWhenMissing(t *testing.T) {
	svg := `<svg xmlns="http://www.w3.org/2000/svg">`
	newVB := `viewBox="0 0 200 200"`
	result := rewriteViewBox(svg, newVB)
	if !strings.Contains(result, newVB) {
		t.Errorf("expected viewBox to be inserted, got %q", result)
	}
}

// ---- createBlankImage -------------------------------------------------------

func TestCreateBlankImage_DimensionsAndAllWhite(t *testing.T) {
	img := createBlankImage(10, 10)
	bounds := img.Bounds()
	if bounds.Dx() != 10 || bounds.Dy() != 10 {
		t.Errorf("expected 10x10, got %dx%d", bounds.Dx(), bounds.Dy())
	}
	for _, v := range img.Pix {
		if v != 0xFF {
			t.Errorf("expected all pixels 0xFF, found %02x", v)
			break
		}
	}
}

// ---- RenderZoomed -----------------------------------------------------------

// renderTestSVG is a small synthetic SVG with known province shapes for render tests.
const renderTestSVG = `<?xml version="1.0"?>
<svg xmlns="http://www.w3.org/2000/svg"
   xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape"
   viewBox="0 0 100 100" width="100" height="100">
  <g id="highlights"/>
  <polygon inkscape:label="foo" points="10,10 20,10 20,20 10,20" id="foo"/>
  <polygon inkscape:label="bar" points="50,50 80,50 80,80 50,80" id="bar"/>
  <polygon inkscape:label="baz" points="5,5 15,5 15,15 5,15" id="baz"/>
</svg>`

func TestRenderZoomed_WithMatchingProvinces_ReturnsPNG(t *testing.T) {
	is := is.New(t)
	result, err := RenderZoomed(stubEngineState{}, []byte(renderTestSVG), []string{"foo"})
	is.NoErr(err)
	assertPNG(t, result)

	// Verify output dimensions differ from full canvas (100×100).
	img, err := png.Decode(bytes.NewReader(result))
	is.NoErr(err)
	bounds := img.Bounds()
	if bounds.Dx() == 100 {
		t.Errorf("expected zoomed width != 100, got %d", bounds.Dx())
	}
}

func TestRenderZoomed_EmptyProvinceList_UsesFullCanvas(t *testing.T) {
	is := is.New(t)
	result, err := RenderZoomed(stubEngineState{}, []byte(renderTestSVG), nil)
	is.NoErr(err)
	assertPNG(t, result)
}

func TestRenderZoomed_ProvinceNotInSVG_UsesFullCanvas(t *testing.T) {
	is := is.New(t)
	result, err := RenderZoomed(stubEngineState{}, []byte(renderTestSVG), []string{"nonexistent"})
	is.NoErr(err)
	is.NotNil(result)
}

func TestRenderZoomed_MultiProvince_ExpandsBoundingBox(t *testing.T) {
	is := is.New(t)
	// "foo" spans (10,10)-(20,20); "baz" spans (5,5)-(15,15).
	// Processing "foo" then "baz" triggers x < minX and y < minY branches.
	result, err := RenderZoomed(stubEngineState{}, []byte(renderTestSVG), []string{"foo", "baz"})
	is.NoErr(err)
	is.NotNil(result)
}

func TestRenderZoomed_NoBoundsNoViewBox_UsesDefaultSize(t *testing.T) {
	is := is.New(t)
	// SVG with no viewBox and no numeric width/height; empty provinces.
	// vw and vh end up 0 → blank PNG fallback.
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`)
	result, err := RenderZoomed(stubEngineState{}, svg, nil)
	is.NoErr(err)
	assertPNG(t, result)
}

// ---- doInitFont -------------------------------------------------------------

func TestDoInitFont_MkdirTempFails_NoError(t *testing.T) {
	// If MkdirTemp fails the function returns gracefully without panicking.
	failMkdir := func(string, string) (string, error) { return "", errors.New("no dir") }
	doInitFont(failMkdir, os.WriteFile) // should not panic
}

func TestDoInitFont_WriteFileFails_NoError(t *testing.T) {
	// If WriteFile fails the function returns gracefully.
	dir, err := os.MkdirTemp("", "dipmap-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	okMkdir := func(string, string) (string, error) { return dir, nil }
	failWrite := func(string, []byte, fs.FileMode) error { return errors.New("no write") }
	doInitFont(okMkdir, failWrite) // should not panic
}

// ---- svgToPNGWith -----------------------------------------------------------

func TestSVGToPNGWith_ParseError_ReturnsError(t *testing.T) {
	is := is.New(t)
	// width="100xyz" triggers canvas.ParseSVG to return an error.
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="100xyz" height="100">`)
	_, err := svgToPNGWith(svg, png.Encode)
	is.NotNil(err)
}

func TestSVGToPNGWith_EncodeError_ReturnsError(t *testing.T) {
	is := is.New(t)
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"/>`)
	_, err := svgToPNGWith(svg, failEncode)
	is.NotNil(err)
}

// ---- renderZoomedWith -------------------------------------------------------

func TestRenderZoomedWith_NoDimensions_EncodeError_ReturnsError(t *testing.T) {
	is := is.New(t)
	// SVG with no viewBox and no width/height; empty provinces → vw=0, vh=0
	// → blank-canvas path; encode error propagates.
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`)
	_, err := renderZoomedWith(stubEngineState{}, svg, nil, failEncode)
	is.NotNil(err)
}

func TestRenderZoomedWith_ParseError_ReturnsError(t *testing.T) {
	is := is.New(t)
	// SVG has a valid province shape (for bounding-box extraction) but a
	// bad dimension unit on stroke-width causes canvas.ParseSVG to error.
	// stroke-width is not rewritten by renderZoomedWith so the error persists.
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100" height="100">` +
		`<polygon inkscape:label="foo" points="10,10 20,10 20,20 10,20" id="foo"` +
		` style="stroke-width:5badunit"/>` +
		`</svg>`)
	_, err := renderZoomedWith(stubEngineState{}, svg, []string{"foo"}, png.Encode)
	is.NotNil(err)
}

func TestRenderZoomedWith_EncodeError_ReturnsError(t *testing.T) {
	is := is.New(t)
	_, err := renderZoomedWith(stubEngineState{}, []byte(renderTestSVG), []string{"foo"}, failEncode)
	is.NotNil(err)
}
