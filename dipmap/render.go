// Package dipmap renders Diplomacy board images. It converts godip's SVG map
// assets to JPEG and supports highlighting a neighbourhood of provinces around
// a target territory to radius n (BFS over godip's Graph.Edges()).
package dipmap

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"github.com/zond/godip/variants/classical"
)

// EngineState is the minimal engine interface needed for map rendering.
// engine.Engine satisfies this interface.
type EngineState interface {
	Dump() ([]byte, error)
}

// jpegEncode wraps jpeg.Encode at quality 85 to match the
// func(io.Writer, image.Image) error signature used for encoder injection.
func jpegEncode(w io.Writer, img image.Image) error {
	return jpeg.Encode(w, img, &jpeg.Options{Quality: 85})
}

// classicalSVGOnce guards one-time loading of the style-stripped godip SVG.
var (
	classicalSVGOnce  sync.Once
	classicalSVGBytes []byte
	classicalSVGErr   error
)

// loadClassicalSVGWith loads the SVG from assetFn, strips <style> blocks, and
// returns the resulting bytes. It is the testable core used by classicalSVG.
func loadClassicalSVGWith(assetFn func(string) ([]byte, error)) ([]byte, error) {
	raw, err := assetFn("svg/map.svg")
	if err != nil {
		return nil, fmt.Errorf("dipmap: load SVG asset: %w", err)
	}
	return stripStyles(raw), nil
}

// classicalSVG returns the style-stripped godip classical SVG, loading and
// decompressing it at most once per process lifetime. Thread-safe.
func classicalSVG() ([]byte, error) {
	classicalSVGOnce.Do(func() {
		classicalSVGBytes, classicalSVGErr = loadClassicalSVGWith(classical.Asset)
	})
	return classicalSVGBytes, classicalSVGErr
}

// stripStyles removes all <style>…</style> blocks from svg. oksvg cannot parse
// embedded data URIs (e.g. @font-face with base64 src) inside style blocks.
func stripStyles(svg []byte) []byte {
	return regexp.MustCompile(`(?s)<style[^>]*>.*?</style\s*>`).ReplaceAll(svg, nil)
}

// Render converts godip's classical SVG map to a JPEG byte slice.
// It uses the cached style-stripped SVG so classical.Asset decompression
// only runs once per process.
func Render(state EngineState) ([]byte, error) {
	return renderWith(state, classicalSVG)
}

// renderWith is the testable core of Render; svgLoader is injected in tests.
func renderWith(_ EngineState, svgLoader func() ([]byte, error)) ([]byte, error) {
	svg, err := svgLoader()
	if err != nil {
		return nil, err
	}
	return svgToJPEG(svg)
}

// renderWithLoader is a test helper that bypasses the process-level cache and
// calls assetFn directly. It delegates to renderWith.
func renderWithLoader(state EngineState, assetFn func(string) ([]byte, error)) ([]byte, error) {
	return renderWith(state, func() ([]byte, error) { return loadClassicalSVGWith(assetFn) })
}

// LoadSVG returns the cached classical SVG map bytes. It is used as the svgFn
// in bot.Dispatcher for the territory-zoom rendering path.
func LoadSVG(_ EngineState) ([]byte, error) {
	return classicalSVG()
}

// loadSVGWith is the testable core of LoadSVG with an injectable asset loader.
func loadSVGWith(_ EngineState, assetFn func(string) ([]byte, error)) ([]byte, error) {
	return loadClassicalSVGWith(assetFn)
}

// svgToJPEG converts SVG bytes to JPEG using a pure-Go renderer. It parses the
// SVG viewBox (falling back to explicit width/height) to determine output
// dimensions and rasterises entirely in-process without external binaries.
func svgToJPEG(svg []byte) ([]byte, error) {
	return svgToJPEGWith(svg, jpegEncode)
}

// SVGToJPEG is the exported form of svgToJPEG. It is used as the default imgFn
// for bot.Dispatcher so callers can convert a pre-processed SVG (with overlay
// and highlight layers already applied) to a JPEG byte slice.
func SVGToJPEG(svg []byte) ([]byte, error) { return svgToJPEG(svg) }

// svgToJPEGWith is the testable core of svgToJPEG with an injectable encoder.
func svgToJPEGWith(svg []byte, encoderFn func(io.Writer, image.Image) error) ([]byte, error) {
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
	icon, err := oksvg.ReadIconStream(bytes.NewReader(svg))
	if err != nil {
		return nil, fmt.Errorf("dipmap: parse SVG: %w", err)
	}
	icon.SetTarget(0, 0, float64(w), float64(h))
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// Fill white before drawing so that JPEG encoding (no alpha channel)
	// does not produce black where the SVG has a transparent background.
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	scanner := rasterx.NewScannerGV(w, h, img, img.Bounds())
	raster := rasterx.NewDasher(w, h, scanner)
	icon.Draw(raster, 1.0)
	var buf bytes.Buffer
	if err := encoderFn(&buf, img); err != nil {
		return nil, fmt.Errorf("dipmap: encode image: %w", err)
	}
	return buf.Bytes(), nil
}

// RenderZoomed renders a zoomed JPEG of the given provinces from the
// highlighted SVG. It computes the union bounding box of the listed province
// shapes, adds 5% diagonal padding, and rasterises at 800 px wide preserving
// aspect ratio. If provinces is empty the full canvas viewBox is used.
func RenderZoomed(state EngineState, svg []byte, provinces []string) ([]byte, error) {
	return renderZoomedWith(state, svg, provinces, jpegEncode)
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
		return nil, fmt.Errorf("dipmap: encode image: %w", err)
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
