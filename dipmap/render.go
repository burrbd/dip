// Package dipmap renders Diplomacy board images. It converts godip's SVG map
// assets to PNG and supports highlighting a neighbourhood of provinces around
// a target territory to radius n (BFS over godip's Graph.Edges()).
package dipmap

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/zond/godip/variants/classical"
)

// EngineState is the minimal engine interface needed for map rendering.
// engine.Engine satisfies this interface.
type EngineState interface {
	Dump() ([]byte, error)
}

// Render converts godip's classical SVG map to a PNG byte slice using a
// pure-Go renderer. The state parameter is available for future unit-overlay
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

// LoadSVG returns the raw classical SVG map bytes. It is used as the svgFn
// in bot.Dispatcher for the territory-zoom rendering path.
func LoadSVG(state EngineState) ([]byte, error) {
	return loadSVGWith(state, classical.Asset)
}

// loadSVGWith is the testable core of LoadSVG with an injectable asset loader.
func loadSVGWith(_ EngineState, assetFn func(string) ([]byte, error)) ([]byte, error) {
	svg, err := assetFn("svg/map.svg")
	if err != nil {
		return nil, fmt.Errorf("dipmap: load SVG asset: %w", err)
	}
	return svg, nil
}

// svgToPNG converts SVG bytes to PNG using a pure-Go renderer. It parses the
// SVG viewBox (falling back to explicit width/height) to determine output
// dimensions and rasterises entirely in-process without external binaries.
func svgToPNG(svg []byte) ([]byte, error) {
	return svgToPNGWith(svg, png.Encode)
}

// SVGToPNG is the exported form of svgToPNG. It is used as the default pngFn
// for bot.Dispatcher so callers can convert a pre-processed SVG (with overlay
// and highlight layers already applied) to a PNG byte slice.
func SVGToPNG(svg []byte) ([]byte, error) { return svgToPNG(svg) }

// svgToPNGWith is the testable core of svgToPNG with an injectable PNG encoder.
func svgToPNGWith(svg []byte, encoderFn func(io.Writer, image.Image) error) ([]byte, error) {
	svgStr := string(svg)
	// Prefer viewBox for natural dimensions (godip SVGs use width="100%").
	_, _, vw, vh := parseSVGViewBox(svgStr)
	w, h := int(vw), int(vh)
	if w <= 0 || h <= 0 {
		w, h = parseSVGDimensions(svg)
	}
	if w <= 0 || h <= 0 {
		w, h = 600, 400
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	var buf bytes.Buffer
	if err := encoderFn(&buf, img); err != nil {
		return nil, fmt.Errorf("dipmap: encode PNG: %w", err)
	}
	return buf.Bytes(), nil
}

// RenderZoomed renders a zoomed PNG of the given provinces from the
// highlighted SVG. It computes the union bounding box of the listed province
// shapes, adds 5% diagonal padding, and rasterises at 800 px wide preserving
// aspect ratio. If provinces is empty the full canvas viewBox is used.
func RenderZoomed(state EngineState, svg []byte, provinces []string) ([]byte, error) {
	return renderZoomedWith(state, svg, provinces, png.Encode)
}

// renderZoomedWith is the testable core of RenderZoomed with an injectable encoder.
func renderZoomedWith(_ EngineState, svg []byte, provinces []string, encoderFn func(io.Writer, image.Image) error) ([]byte, error) {
	svgStr := string(svg)

	var minX, minY, maxX, maxY float64
	hasBounds := false

	for _, p := range provinces {
		_, data := extractProvinceShape(svgStr, strings.ToLower(p))
		if data == "" {
			continue
		}
		xs, ys := parseCoordinates(data)
		for i := range xs {
			x, y := xs[i], ys[i]
			if !hasBounds {
				minX, maxX, minY, maxY = x, x, y, y
				hasBounds = true
			} else {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	var vw, vh float64
	if hasBounds {
		dx := maxX - minX
		dy := maxY - minY
		diag := math.Sqrt(dx*dx + dy*dy)
		pad := diag * 0.05
		vw = dx + 2*pad
		vh = dy + 2*pad
	} else {
		_, _, vw, vh = parseSVGViewBox(svgStr)
		if vw <= 0 || vh <= 0 {
			fw, fh := parseSVGDimensions([]byte(svgStr))
			vw, vh = float64(fw), float64(fh)
		}
	}

	const outputWidth = 800
	outH := outputWidth
	if vw > 0 && vh > 0 {
		outH = int(float64(outputWidth) * vh / vw)
	}

	img := image.NewRGBA(image.Rect(0, 0, outputWidth, outH))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	var buf bytes.Buffer
	if err := encoderFn(&buf, img); err != nil {
		return nil, fmt.Errorf("dipmap: encode PNG: %w", err)
	}
	return buf.Bytes(), nil
}

// parseSVGViewBox extracts the viewBox x, y, w, h values from an SVG string.
// Returns zero values if no valid viewBox is found.
func parseSVGViewBox(svg string) (x, y, w, h float64) {
	re := regexp.MustCompile(`viewBox="([^"]+)"`)
	m := re.FindStringSubmatch(svg)
	if len(m) < 2 {
		return
	}
	parts := strings.Fields(m[1])
	if len(parts) != 4 {
		return
	}
	x, _ = strconv.ParseFloat(parts[0], 64)
	y, _ = strconv.ParseFloat(parts[1], 64)
	w, _ = strconv.ParseFloat(parts[2], 64)
	h, _ = strconv.ParseFloat(parts[3], 64)
	return
}

// parseSVGDimensions extracts pixel width and height from numeric SVG
// width/height attributes. Returns zero values if not found or if the
// attributes use percentage values.
func parseSVGDimensions(svg []byte) (w, h int) {
	s := string(svg)
	wRe := regexp.MustCompile(`\bwidth="([0-9]+(?:\.[0-9]+)?)"`)
	hRe := regexp.MustCompile(`\bheight="([0-9]+(?:\.[0-9]+)?)"`)
	if m := wRe.FindStringSubmatch(s); len(m) > 1 {
		f, _ := strconv.ParseFloat(m[1], 64)
		w = int(f)
	}
	if m := hRe.FindStringSubmatch(s); len(m) > 1 {
		f, _ := strconv.ParseFloat(m[1], 64)
		h = int(f)
	}
	return
}

// parseCoordinates extracts all (x, y) numeric pairs from polygon points or
// SVG path d attribute data. Numbers are extracted in order and paired
// sequentially, so both "x,y x,y" and "M x y L x y" formats work correctly.
func parseCoordinates(data string) (xs, ys []float64) {
	re := regexp.MustCompile(`[-+]?[0-9]*\.?[0-9]+`)
	matches := re.FindAllString(data, -1)
	for i := 0; i+1 < len(matches); i += 2 {
		x, err1 := strconv.ParseFloat(matches[i], 64)
		y, err2 := strconv.ParseFloat(matches[i+1], 64)
		if err1 == nil && err2 == nil {
			xs = append(xs, x)
			ys = append(ys, y)
		}
	}
	return
}
