package dipmap

import (
	"testing"

	"github.com/cheekybits/is"
)

// stubEngineState implements EngineState for test use.
type stubEngineState struct{}

func (s stubEngineState) Dump() ([]byte, error) { return []byte(`{}`), nil }

func TestRender_ReturnsByteSlice(t *testing.T) {
	is := is.New(t)
	result, err := Render(stubEngineState{})
	is.NoErr(err)
	is.NotNil(result)
}

func TestHighlight_ReturnsSVGUnchanged(t *testing.T) {
	is := is.New(t)
	svg := []byte("<svg/>")
	result, err := Highlight(svg, []string{"Vienna", "Budapest"})
	is.NoErr(err)
	is.Equal(string(result), "<svg/>")
}

func TestHighlight_EmptyProvinces(t *testing.T) {
	is := is.New(t)
	svg := []byte("<svg/>")
	result, err := Highlight(svg, nil)
	is.NoErr(err)
	is.Equal(string(result), "<svg/>")
}
