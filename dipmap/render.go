// Package dipmap renders Diplomacy board images. It converts godip's SVG map
// assets to PNG and supports highlighting a neighbourhood of provinces around
// a target territory to radius n (BFS over godip's Graph.Edges()).
package dipmap

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/zond/godip/variants/classical"
)

// EngineState is the minimal engine interface needed for map rendering.
// engine.Engine satisfies this interface.
type EngineState interface {
	Dump() ([]byte, error)
}

// Render converts godip's classical SVG map to a PNG byte slice using
// rsvg-convert. The state parameter is available for future unit-overlay
// rendering; the current implementation renders the base map only.
func Render(state EngineState) ([]byte, error) {
	return renderWithLoader(state, classical.Asset)
}

// renderWithLoader is the testable core of Render with an injectable asset loader.
func renderWithLoader(_ EngineState, assetFn func(string) ([]byte, error)) ([]byte, error) {
	svg, err := assetFn("svg/map.svg")
	if err != nil {
		return nil, fmt.Errorf("dipmap: load SVG asset: %w", err)
	}
	return svgToPNG(svg)
}

// svgToPNG converts SVG bytes to PNG by piping through rsvg-convert.
func svgToPNG(svg []byte) ([]byte, error) {
	cmd := exec.Command("rsvg-convert", "--format=png")
	cmd.Stdin = bytes.NewReader(svg)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("dipmap: rsvg-convert: %w", err)
	}
	return out, nil
}
