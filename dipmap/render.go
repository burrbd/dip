// Package dipmap renders Diplomacy board images. It converts godip's SVG map
// assets to PNG and supports highlighting a neighbourhood of provinces around
// a target territory to radius n (BFS over godip's Graph.Edges()).
package dipmap

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
)

//go:embed assets/map.svg
var embeddedMapSVG []byte

//go:embed assets/LibreBaskerville-Bold.ttf
var libreBaskerville []byte

var fontOnce sync.Once

// initFont extracts the embedded LibreBaskerville-Bold.ttf to a temporary
// directory and registers it with canvas's system font resolver. This runs
// exactly once per process; subsequent calls are no-ops.
func initFont() {
	fontOnce.Do(func() { doInitFont(os.MkdirTemp, os.WriteFile) })
}

// doInitFont is the testable inner body of initFont. Tests inject failing OS
// functions to cover error branches without a real timer or sync.Once.
func doInitFont(
	mkdirTemp func(string, string) (string, error),
	writeFile func(string, []byte, fs.FileMode) error,
) {
	dir, err := mkdirTemp("", "dipmap-fonts-*")
	if err != nil {
		return
	}
	path := filepath.Join(dir, "LibreBaskerville-Bold.ttf")
	if err := writeFile(path, libreBaskerville, 0644); err != nil {
		return
	}
	// CacheSystemFonts with /dev/null as cache file: stat succeeds but
	// IsRegular() returns false, so it scans the given dir. Writing to
	// /dev/null is a no-op. This populates canvas's global systemFonts.
	_ = canvas.CacheSystemFonts("/dev/null", []string{dir})
}

// EngineState is the minimal engine interface needed for map rendering.
// engine.Engine satisfies this interface.
type EngineState interface {
	Dump() ([]byte, error)
}

// Render converts the classical SVG map to a PNG byte slice using the
// embedded map.svg with pre-placed unit placeholder glyphs.
func Render(_ EngineState) ([]byte, error) {
	return svgToPNGWith(embeddedMapSVG, png.Encode)
}

// SVGToPNG converts SVG bytes to a lossless PNG using tdewolff/canvas.
// The embedded LibreBaskerville-Bold font is registered on first call so that
// the SVG names layer renders with the correct typeface.
func SVGToPNG(svg []byte) ([]byte, error) {
	return svgToPNGWith(svg, png.Encode)
}

// svgToPNGWith is the testable inner body of SVGToPNG. Tests inject a failing
// encode function to cover the error branch.
func svgToPNGWith(svg []byte, encode func(io.Writer, image.Image) error) ([]byte, error) {
	initFont()
	c, err := canvas.ParseSVG(bytes.NewReader(svg))
	if err != nil {
		return nil, fmt.Errorf("dipmap: parse SVG: %w", err)
	}
	w, h := c.Size()
	if w <= 0 || h <= 0 {
		return nil, fmt.Errorf("dipmap: SVG has zero canvas dimensions")
	}
	img := rasterizer.Draw(c, canvas.DPMM(1), canvas.DefaultColorSpace)
	var buf bytes.Buffer
	if err := encode(&buf, img); err != nil {
		return nil, fmt.Errorf("dipmap: encode PNG: %w", err)
	}
	return buf.Bytes(), nil
}

// LoadSVG returns the SVG map bytes from the embedded asset. It is used as the
// svgFn in bot.Dispatcher for the territory-zoom rendering path.
func LoadSVG(_ EngineState) ([]byte, error) {
	return embeddedMapSVG, nil
}

// RenderZoomed renders a zoomed PNG of the given provinces from the
// highlighted SVG. It computes the union bounding box of the listed province
// shapes, adds 5% diagonal padding, rewrites the SVG viewBox to crop to that
// region, and rasterises at 800 px wide preserving aspect ratio. If provinces
// is empty the full canvas viewBox is used.
func RenderZoomed(state EngineState, svg []byte, provinces []string) ([]byte, error) {
	return renderZoomedWith(state, svg, provinces, png.Encode)
}

// renderZoomedWith is the testable inner body of RenderZoomed. Tests inject a
// failing encode function to cover the error branch.
func renderZoomedWith(
	_ EngineState,
	svg []byte,
	provinces []string,
	encode func(io.Writer, image.Image) error,
) ([]byte, error) {
	initFont()
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

	var vbX, vbY, vw, vh float64
	if hasBounds {
		dx := maxX - minX
		dy := maxY - minY
		diag := math.Sqrt(dx*dx + dy*dy)
		pad := diag * 0.05
		vbX = minX - pad
		vbY = minY - pad
		vw = dx + 2*pad
		vh = dy + 2*pad
	} else {
		vbX, vbY, vw, vh = parseSVGViewBox(svgStr)
		if vw <= 0 || vh <= 0 {
			fw, fh := parseSVGDimensions(svg)
			vbX, vbY, vw, vh = 0, 0, float64(fw), float64(fh)
		}
	}

	if vw <= 0 || vh <= 0 {
		// No usable canvas dimensions — return a blank PNG.
		const defaultSize = 800
		img := createBlankImage(defaultSize, defaultSize)
		var buf bytes.Buffer
		if err := encode(&buf, img); err != nil {
			return nil, fmt.Errorf("dipmap: encode PNG: %w", err)
		}
		return buf.Bytes(), nil
	}

	// Rewrite the viewBox and width/height to crop to the computed region.
	// Setting width=vw and height=vh ensures canvas creates a vw×vh mm canvas,
	// and DPMM = 800/vw gives an 800-px-wide output.
	const outputWidth = 800.0
	newVB := fmt.Sprintf(`viewBox="%.4f %.4f %.4f %.4f"`, vbX, vbY, vw, vh)
	renderSVG := rewriteViewBox(svgStr, newVB)
	renderSVG = rewriteDimensions(renderSVG, vw, vh)

	c, err := canvas.ParseSVG(bytes.NewReader([]byte(renderSVG)))
	if err != nil {
		return nil, fmt.Errorf("dipmap: parse SVG: %w", err)
	}
	dpmm := outputWidth / vw
	img := rasterizer.Draw(c, canvas.DPMM(dpmm), canvas.DefaultColorSpace)
	var buf bytes.Buffer
	if err := encode(&buf, img); err != nil {
		return nil, fmt.Errorf("dipmap: encode PNG: %w", err)
	}
	return buf.Bytes(), nil
}

// rewriteDimensions replaces the width and height attributes in the SVG root
// element with the given values (in SVG user-unit/mm). This is used by
// RenderZoomed so the canvas dimensions match the cropped viewBox region.
func rewriteDimensions(svg string, w, h float64) string {
	wAttr := fmt.Sprintf(`width="%.4f"`, w)
	hAttr := fmt.Sprintf(`height="%.4f"`, h)
	wRe := regexp.MustCompile(`\bwidth="[^"]*"`)
	hRe := regexp.MustCompile(`\bheight="[^"]*"`)
	if wRe.MatchString(svg) {
		svg = wRe.ReplaceAllString(svg, wAttr)
	}
	if hRe.MatchString(svg) {
		svg = hRe.ReplaceAllString(svg, hAttr)
	}
	return svg
}

// rewriteViewBox replaces (or inserts) the viewBox attribute in the SVG root
// element with the new value. This is used by RenderZoomed to crop the canvas.
func rewriteViewBox(svg, newVB string) string {
	re := regexp.MustCompile(`viewBox="[^"]*"`)
	if re.MatchString(svg) {
		return re.ReplaceAllString(svg, newVB)
	}
	// No viewBox attribute — insert before the closing > of the <svg> tag.
	svgTagRe := regexp.MustCompile(`<svg\b[^>]*>`)
	return svgTagRe.ReplaceAllStringFunc(svg, func(tag string) string {
		return tag[:len(tag)-1] + " " + newVB + ">"
	})
}

// createBlankImage creates a plain white RGBA image of the given size.
func createBlankImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := range img.Pix {
		img.Pix[i] = 0xFF
	}
	return img
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
