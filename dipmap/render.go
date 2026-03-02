// Package dipmap renders Diplomacy board images. It converts godip's SVG map
// assets to PNG and supports highlighting a neighbourhood of provinces around
// a target territory to radius n (BFS over godip's Graph.Edges()).
package dipmap

// EngineState is the minimal engine interface needed for map rendering.
// engine.Engine satisfies this interface.
type EngineState interface {
	Dump() ([]byte, error)
}

// Render converts the current board state to a PNG byte slice.
// The stub implementation returns empty bytes; Story 9 replaces this with
// real SVG-to-PNG rendering using godip's SVG assets.
func Render(_ EngineState) ([]byte, error) {
	return []byte{}, nil
}

// Highlight overlays province highlighting on an SVG byte slice.
// The stub returns the input unchanged; Story 9 replaces this with
// actual province-colouring logic.
func Highlight(svg []byte, _ []string) ([]byte, error) {
	return svg, nil
}
