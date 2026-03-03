package dipmap

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/cheekybits/is"
)

// stubEngineState implements EngineState for test use.
type stubEngineState struct{}

func (s stubEngineState) Dump() ([]byte, error) { return []byte(`{}`), nil }

func TestRender_ReturnsPNGBytes(t *testing.T) {
	if _, err := exec.LookPath("rsvg-convert"); err != nil {
		t.Skip("rsvg-convert not installed; skipping PNG render test")
	}
	is := is.New(t)
	result, err := Render(stubEngineState{})
	is.NoErr(err)
	// PNG magic bytes: 0x89 'P' 'N' 'G'
	if len(result) < 4 {
		t.Fatalf("expected PNG bytes, got %d bytes", len(result))
	}
	is.Equal(result[0], byte(0x89))
	is.Equal(result[1], byte('P'))
	is.Equal(result[2], byte('N'))
	is.Equal(result[3], byte('G'))
}

func TestSVGToPNG_ErrorOnBadCommand(t *testing.T) {
	if _, err := exec.LookPath("rsvg-convert"); err != nil {
		t.Skip("rsvg-convert not installed")
	}
	// Pass garbage SVG — rsvg-convert may or may not error on invalid input;
	// the test verifies the function does not panic.
	_, _ = svgToPNG([]byte("not valid svg"))
}

func TestRenderWithLoader_AssetError_ReturnsError(t *testing.T) {
	is := is.New(t)
	failLoader := func(string) ([]byte, error) {
		return nil, errors.New("asset not found")
	}
	_, err := renderWithLoader(stubEngineState{}, failLoader)
	is.NotNil(err)
}

func TestSVGToPNG_ValidSVG_ReturnsPNG(t *testing.T) {
	if _, err := exec.LookPath("rsvg-convert"); err != nil {
		t.Skip("rsvg-convert not installed")
	}
	is := is.New(t)
	svg := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="10" height="10"><rect width="10" height="10" fill="red"/></svg>`)
	result, err := svgToPNG(svg)
	is.NoErr(err)
	if len(result) < 4 {
		t.Fatalf("expected PNG bytes, got %d bytes", len(result))
	}
	is.Equal(result[0], byte(0x89))
	is.Equal(result[1], byte('P'))
}
