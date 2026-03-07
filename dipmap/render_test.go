package dipmap

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io"
	"testing"

	"github.com/cheekybits/is"
)

// stubEngineState implements EngineState for test use.
type stubEngineState struct{}

func (s stubEngineState) Dump() ([]byte, error) { return []byte(`{}`), nil }

// failEncoder is a PNG encoder that always returns an error.
func failEncoder(_ io.Writer, _ image.Image) error { return errors.New("encode fail") }

// ---- Render -----------------------------------------------------------------

func TestRender_ReturnsPNGBytes(t *testing.T) {
	is := is.New(t)
	result, err := Render(stubEngineState{})
	is.NoErr(err)
	if len(result) < 4 {
		t.Fatalf("expected PNG bytes, got %d bytes", len(result))
	}
	// PNG magic bytes: 0x89 'P' 'N' 'G'
	is.Equal(result[0], byte(0x89))
	is.Equal(result[1], byte('P'))
	is.Equal(result[2], byte('N'))
	is.Equal(result[3], byte('G'))
}

func TestRenderWithLoader_AssetError_ReturnsError(t *testing.T) {
	is := is.New(t)
	failLoader := func(string) ([]byte, error) {
		return nil, errors.New("asset not found")
	}
	_, err := renderWithLoader(stubEngineState{}, failLoader)
	is.NotNil(err)
}

// ---- LoadSVG ----------------------------------------------------------------

func TestLoadSVG_ReturnsSVGBytes(t *testing.T) {
	is := is.New(t)
	result, err := LoadSVG(stubEngineState{})
	is.NoErr(err)
	is.NotNil(result)
	// SVG must start with an XML header or <svg element.
	if len(result) < 4 {
		t.Fatalf("expected SVG bytes, got %d bytes", len(result))
	}
}

func TestLoadSVGWith_AssetError_ReturnsError(t *testing.T) {
	is := is.New(t)
	failLoader := func(string) ([]byte, error) {
		return nil, errors.New("asset not found")
	}
	_, err := loadSVGWith(stubEngineState{}, failLoader)
	is.NotNil(err)
}

// ---- svgToPNG ---------------------------------------------------------------

func TestSVGToPNG_ValidSVGWithViewBox_ReturnsPNG(t *testing.T) {
	is := is.New(t)
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 80"><rect width="100" height="80" fill="white"/></svg>`)
	result, err := svgToPNG(svg)
	is.NoErr(err)
	if len(result) < 4 {
		t.Fatalf("expected PNG bytes, got %d bytes", len(result))
	}
	is.Equal(result[0], byte(0x89))
	is.Equal(result[1], byte('P'))
}

func TestSVGToPNG_ValidSVGWithWidthHeight_ReturnsPNG(t *testing.T) {
	is := is.New(t)
	// No viewBox — falls back to width/height.
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10"><rect width="10" height="10" fill="red"/></svg>`)
	result, err := svgToPNG(svg)
	is.NoErr(err)
	if len(result) < 4 {
		t.Fatalf("expected PNG bytes, got %d bytes", len(result))
	}
	is.Equal(result[0], byte(0x89))
}

func TestSVGToPNG_NoDimensions_ReturnsDefaultSizePNG(t *testing.T) {
	is := is.New(t)
	// No viewBox, no width/height → falls back to 600×400.
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg"><rect/></svg>`)
	result, err := svgToPNG(svg)
	is.NoErr(err)
	is.NotNil(result)
}

func TestSVGToPNG_ExportedWrapper_ReturnsPNG(t *testing.T) {
	is := is.New(t)
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"/>`)
	result, err := SVGToPNG(svg)
	is.NoErr(err)
	is.NotNil(result)
}

func TestSVGToPNGWith_EncoderError_ReturnsError(t *testing.T) {
	is := is.New(t)
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 10 10"/>`)
	_, err := svgToPNGWith(svg, failEncoder)
	is.NotNil(err)
}

func TestSVGToPNGWith_MalformedXML_ReturnsError(t *testing.T) {
	is := is.New(t)
	// Malformed XML causes oksvg.ReadIconStream to return an XML syntax error.
	_, err := svgToPNGWith([]byte(`<svg><unclosed`), png.Encode)
	is.NotNil(err)
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
	// viewBox present but with wrong number of parts.
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
	// Percentage values are rejected by the regex.
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
	if len(result) < 4 {
		t.Fatalf("expected PNG bytes, got %d bytes", len(result))
	}
	is.Equal(result[0], byte(0x89))
	is.Equal(result[1], byte('P'))

	// Verify the output dimensions differ from the original canvas (100×100).
	img, err := png.Decode(bytes.NewReader(result))
	is.NoErr(err)
	bounds := img.Bounds()
	// The zoomed view of "foo" (10×10 box + padding) at 800px wide
	// should not produce a 100px wide image.
	if bounds.Dx() == 100 {
		t.Errorf("expected zoomed width != 100, got %d", bounds.Dx())
	}
}

func TestRenderZoomed_EmptyProvinceList_UsesFullCanvas(t *testing.T) {
	is := is.New(t)
	// Empty provinces → falls back to full canvas viewBox → valid PNG.
	result, err := RenderZoomed(stubEngineState{}, []byte(renderTestSVG), nil)
	is.NoErr(err)
	if len(result) < 4 {
		t.Fatalf("expected PNG bytes, got %d bytes", len(result))
	}
	is.Equal(result[0], byte(0x89))
}

func TestRenderZoomed_ProvinceNotInSVG_UsesFullCanvas(t *testing.T) {
	is := is.New(t)
	// Provinces listed but not found in SVG → no bounds → full canvas fallback.
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
	// SVG with no viewBox and no numeric width/height; provinces empty.
	// vw and vh end up 0 → outH = outputWidth (800).
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg"></svg>`)
	result, err := renderZoomedWith(stubEngineState{}, svg, nil, png.Encode)
	is.NoErr(err)
	is.NotNil(result)
}

func TestRenderZoomedWith_EncoderError_ReturnsError(t *testing.T) {
	is := is.New(t)
	_, err := renderZoomedWith(stubEngineState{}, []byte(renderTestSVG), []string{"foo"}, failEncoder)
	is.NotNil(err)
}

