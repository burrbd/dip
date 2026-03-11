# Build Plan

Ordered stories toward a working Diplomacy messenger bot. Each story is a discrete,
shippable unit of work that unblocks the next.

Completed stories have been extracted to [`PLAN_ARCHIVE.md`](PLAN_ARCHIVE.md).

---

## For agents: how to use this file

1. **Find the next story** — the first unchecked box (`- [ ]`) is the current story to work on.
2. **Start the story** — read its goal, files, and acceptance criteria before writing any code.
3. **Work TDD** — write failing tests first, then implement, then refactor.
4. **Check it off** — once all acceptance criteria are met and `go test -v -cover -race ./...`
   passes, change `- [ ]` to `- [x]` for that story, then commit.
5. **Move on** — proceed to the next unchecked story.

Never mark a story done if any test is failing or any criterion is unmet.

---

## Stories

- [x] Story 9g — Map Rendering Polish (zoom radius, labels, unit geometry, scale, z-order)
- [x] Story 9h — Pre-populate SVG with unit placeholder glyphs
- [x] Story 10a — SVG Simplification: Strip Inkscape Metadata
- [x] Story 10b — Map Rendering Fixes: Labels, Unit Colours, and Glyph Geometry
- [ ] Story 10c — Zoomed /map with Territory and Radius (deferred)
- [ ] Story 10d — Replace oksvg with tdewolff/canvas for proper text rendering
- [ ] Story 13 — Lambda / EventBridge Deployment
- [ ] Story 11 — Slack Platform Adapter
- [ ] Story 12 — WhatsApp Platform Adapter (optional)

---

### Story 9g — Map Rendering Polish

**Goal:** Fix eight rendering defects visible in the live map output: broken zoom radius,
missing province labels, missing/misplaced units, oversized units, wrong unit geometry,
and z-order issues that let units obscure labels and supply-centre markers.

**Files:** `dipmap/overlay.go`, `dipmap/overlay_test.go`, `dipmap/render.go`,
`dipmap/render_test.go`, `bot/commands.go`, `bot/commands_test.go`,
`bot/bot_functional_test.go`

---

#### SVG structure — what the godip map actually contains

Inspecting `classical.Asset("svg/map.svg")` (the vendored godip map) reveals several layers
relevant to this story. All layer names come from `inkscape:label` attributes on `<g>` elements.

**`units` layer (line 3669):** The units layer is **empty** — a self-closing `<g id="units"/>`.
godip provides the board template only; it does not embed pre-drawn unit glyphs. The
`dipmap.Overlay` injection approach (writing child elements into this group before rendering)
is therefore correct and intentional.

**`names` layer (line 2588):** Province name labels are already in the SVG as `<text>` elements
with hand-positioned `x`/`y` coordinates (some with `transform="rotate(…)"` for diagonal labels).
They fail to render because:
1. `stripStyles` removes `<style>` blocks but not inline `style=` attributes; the inline styles
   reference the font `LibreBaskerville-Bold` which oksvg cannot load, causing the text elements
   to be silently dropped.
2. Some labels use `display:inline` in their inline style; `prepareForRender` only handles
   `display:none`, so `display:inline` is left as-is and oksvg may still reject them.

The fix for Issue 2 should work with these existing `<text>` elements rather than injecting a
new layer — strip or simplify their inline styles so oksvg can render them with a fallback font.

**`province-centers` layer (line 859):** Every province has a pre-computed center marker
`<path id="<province>Center" …>` (e.g. `id="vieCenter"`, `id="mosCenter"`). These are the
concentric-circle supply-marker rings visible on the board. Their `d` attribute encodes the
exact visual centre of each province.

This is the correct data source for unit placement (Issues 3 and 5). The current
`provinceCenter` function computes centroids by averaging polygon vertices from the `provinces`
layer — an approach that fails for provinces whose polygon uses relative SVG path commands,
producing absurd coordinates (hence the phantom army near Iceland). Using the `<province>Center`
path centroid instead eliminates the polygon-averaging entirely.

**Parsing `<province>Center` paths:** Each center marker is a multi-ring concentric circle
encoded as four `m`/`c`/`z` sub-paths. The centroid is the translation point of the first
sub-path — the first pair of numbers after the leading `m` command. For example:

```
d="m 748.83,856.44 c … z m … z m … z m … z"
      ^^^^^^^ ^^^^^^^
      cx      cy  (these are the province centre coordinates)
```

Extract `cx, cy` by taking the first two numbers in the `d` string (after stripping the leading
`m`). No full path parsing is needed.

---

#### Issue 1 — `/map <territory> <n>` radius is ignored

`RenderZoomed` receives the highlighted province list but the neighbourhood BFS result is
not being passed through correctly — the bot calls `Neighborhood` but the radius `n` parsed
from the command arguments is silently lost, so every zoomed render uses the same single-province
highlight regardless of the number given.

**Fix:** Trace the `n` argument from command parsing through to `Neighborhood(graph, province, n)`
and confirm the enlarged province set is forwarded to `RenderZoomed`.

**Acceptance criteria:**
- `/map lon 1` crops to London + all immediate neighbours.
- `/map lon 3` crops to a larger region than `/map lon 1`.
- `TestCommand_Map_WithTerritoryAndRadius` is extended (or a new sub-test added) that asserts
  the zoomed image has different pixel dimensions than the single-province render.

---

#### Issue 2 — No province labels on the map

The godip SVG already has a `names` layer (line 2588) containing properly positioned `<text>`
elements for every province, including rotated labels for diagonal provinces. They do not render
because their inline `style=` attributes reference `LibreBaskerville-Bold` (unavailable to oksvg)
and oksvg silently discards text elements with unresolvable fonts.

**Fix:** In the SVG pre-processing pipeline (before passing to oksvg), rewrite the inline style
of every `<text>` element in the `names` layer to a minimal style oksvg can handle:

```
style="font-size:16px;fill:#000000"
```

Strip `font-family`, `font-variant-*`, `-inkscape-font-specification`, and similar properties.
Leave `font-size`, `fill`, `text-anchor`, and `writing-mode` intact. The `transform="rotate(…)"`
attributes must be preserved as they control label orientation.

Font size should also scale with canvas width so labels remain legible on zoomed renders:
`scaledSize = baseFontSize * (canvasWidth / naturalWidth)` where `naturalWidth` is the SVG's
viewBox width (≈1524). Apply this scaling in the pipeline step that rewrites font sizes.

**Acceptance criteria:**
- After pre-processing, `<text>` elements in the names layer have simplified `style=` attributes
  containing no `font-family` property.
- Full-board render contains visible text for at least the seven capital provinces
  (Lon, Par, Ber, Mos, Rom, Vie, Con) — verify by asserting the SVG passed to oksvg contains
  those `<text>` strings.
- A unit test covers the style-rewriting helper on a synthetic `<text>` element.

#### Combined acceptance criteria

- Issues 1 and 2 individual acceptance criteria above are met.
- Issues 3–8 are superseded by Story 9h, which replaces unit injection entirely.
- `go test -v -cover -race ./...` passes at 100% for `dipmap/` and `bot/`.
- `go test -v -tags functional ./bot/` passes.

---

### Story 9h — Pre-populate SVG with unit placeholder glyphs

**Goal:** Produce a hand-editable SVG (`dipmap/assets/map.svg`) that is our own copy of the
godip classical map, augmented with one army glyph and one fleet glyph for every province
centre. All glyphs start invisible (`display:none`). The bot's `Overlay` function is then
rewritten to activate glyphs by ID (setting `display:inline` and `fill`) rather than computing
centroids and injecting SVG at runtime.

After the agent completes this story, the owner opens `dipmap/assets/map.svg` in Inkscape and
manually adjusts glyph positions to taste. The bot code is unaffected by those manual tweaks —
it only reads IDs, not coordinates.

This story supersedes Story 9g Issues 3–8 (unit placement, unit geometry, phantom units, unit
scaling, and z-order). The generator controls layer ordering in the SVG directly, so no runtime
lifting is needed. Issues 1 and 2 from Story 9g are unaffected.

**Files:**
- `cmd/mkapsvg/main.go` — new one-off generator (keep it; it documents how the SVG was built)
- `dipmap/assets/map.svg` — new file: populated copy of the godip SVG
- `dipmap/render.go` — update asset loading to use our SVG instead of godip's
- `dipmap/overlay.go` — rewrite unit injection to use ID-based attribute setting
- `dipmap/overlay_test.go` — updated tests

---

#### Part A — Generate the populated SVG

Write `cmd/mkapsvg/main.go`. It must:

1. Load the godip classical SVG via `classical.Asset("svg/map.svg")`.

2. Find every province centre marker using the regex `id="([^"]+)Center"`. There are exactly
   81 such markers (75 base provinces + 6 coastal variants: `bul/ec`, `bul/sc`, `stp/nc`,
   `stp/sc`, `spa/nc`, `spa/sc`). Each has `d="m cx,cy …"` — the centroid is the first two
   numbers in the `d` attribute (the translation of the first `m` sub-path).

3. For each province `p` with centroid `(cx, cy)`, generate two `<g>` elements and insert them
   as children of the `<g id="units">` layer:

   ```xml
   <g id="unit-{pid}-army"  transform="translate({cx},{cy})" display="none">
     <rect x="-12" y="-12" width="24" height="24" rx="3" fill="#cccccc" stroke="#ffffff" stroke-width="2"/>
     <text x="0" y="5" text-anchor="middle" font-size="14" fill="#000000">A</text>
   </g>
   <g id="unit-{pid}-fleet" transform="translate({cx},{cy})" display="none">
     <rect x="-15" y="-9" width="30" height="18" rx="3" fill="#cccccc" stroke="#ffffff" stroke-width="2"/>
     <text x="0" y="5" text-anchor="middle" font-size="10" fill="#000000">F</text>
   </g>
   ```

   Where `{pid}` is the province name with `/` replaced by `-` (e.g. `stp/nc` → `unit-stp-nc-fleet`).
   Use a fixed placeholder fill (`#cccccc`) — the bot will set the real nation colour at render time.

4. For the three provinces that have coastal variants (`bul`, `spa`, `stp`), **also** emit army
   glyphs at the base province centre (the bot will only ever activate armies at the base, never
   at the coastal variant). For coastal variant IDs (`bul/ec`, `bul/sc`, etc.) emit **fleet only**
   (no army glyph).

5. Write the result to `dipmap/assets/map.svg` (create the directory if needed).

6. Print a summary: total centres found, total glyphs written.

**Acceptance criteria for Part A:**
- `go run ./cmd/mkapsvg/` completes without error and writes `dipmap/assets/map.svg`.
- The output SVG contains exactly 81 `unit-…-fleet` elements and 75 `unit-…-army` elements
  (coastal variant positions are fleet-only).
- Every `<g id="unit-…">` has `display="none"`.
- A unit test in `cmd/mkapsvg/main_test.go` loads the generated SVG and asserts these counts
  using a regex or XML parser. The test should run the generator against a real or synthetic
  input and check the output.

---

#### Part B — Switch asset loading and simplify Overlay

**Render.go:** Replace the call to `classical.Asset("svg/map.svg")` with a read of the embedded
`dipmap/assets/map.svg`. Use `//go:embed assets/map.svg` in `render.go` (or a new `assets.go`
file in the `dipmap` package).

**Overlay.go — unit injection rewrite:**

The current approach computes a centroid per province and injects SVG markup. Replace it with:

```go
func injectUnits(svg string, units map[string]UnitInfo) string {
    for province, u := range units {
        pid := strings.ReplaceAll(province, "/", "-")
        id := fmt.Sprintf("unit-%s-%s", pid, strings.ToLower(u.Type))
        colour := nationColour(u.Nation)
        svg = setAttr(svg, id, "display", "inline")
        svg = setAttr(svg, id, "fill", colour)   // sets fill on the <g>; children inherit
    }
    return svg
}
```

`setAttr(svg, id, attr, val string) string` — a small helper that finds
`id="<id>"` in the SVG and updates or inserts the named attribute on that element.
If the element is not found (unknown province), log to `stderr` and continue.

Remove `provinceCenter`, `extractProvinceShape`, `parseCoordinates`, `unitGlyph`, and all
centroid-computation code from `overlay.go` — it is no longer needed.

**Acceptance criteria for Part B:**
- `dipmap.Overlay` no longer contains centroid computation or glyph injection code.
- `setAttr` is tested with a synthetic SVG snippet; edge cases: attribute already present
  (must update, not duplicate), attribute absent (must add), element not found (no-op + log).
- A full-board overlay with all 22 starting units produces an SVG where 22 `unit-…` elements
  have `display="inline"` and the correct nation fill colour.
- `go test -v -cover -race ./...` passes at 100% for `dipmap/`.

---

#### Naming convention reference

| Province  | Army glyph ID          | Fleet glyph ID(s)                        |
|-----------|------------------------|------------------------------------------|
| `lon`     | `unit-lon-army`        | `unit-lon-fleet`                         |
| `mos`     | `unit-mos-army`        | `unit-mos-fleet` (never activated)       |
| `stp`     | `unit-stp-army`        | — (no base fleet; use coastal variants)  |
| `stp/nc`  | — (no army)            | `unit-stp-nc-fleet`                      |
| `stp/sc`  | — (no army)            | `unit-stp-sc-fleet`                      |
| `bul/ec`  | — (no army)            | `unit-bul-ec-fleet`                      |

Godip uses the coast-qualified name (e.g. `"stp/nc"`) as the province key when a fleet
is in a coastal position. The bot passes this key directly to `injectUnits`; the `pid`
normalisation (`/` → `-`) must happen inside `injectUnits`, not in the caller.

---

#### Combined acceptance criteria

- `go run ./cmd/mkapsvg/` produces a valid SVG with 156 unit placeholder elements.
- `dipmap.Overlay` activates exactly the right glyphs for a given game state (verified by
  counting `display="inline"` elements in the output).
- `go test -v -cover -race ./...` passes at 100% for `dipmap/` and `cmd/mkapsvg/`.
- `go test -v -tags functional ./bot/` passes.
- No centroid-computation or glyph-injection code remains in `dipmap/`.

---

### Story 10a — SVG Simplification: Strip Inkscape Metadata

**Goal:** Clean up `dipmap/assets/map.svg` to a minimal, hand-editable SVG by removing all
Inkscape/Sodipodi-specific attributes and RDF metadata while preserving everything the bot
actually uses at runtime: province shapes (`<g id="provinces">`), label positions
(`<g id="names">`), supply-centre markers (`<g id="province-centers">`), highlights layer
(`<g id="highlights"/>`), and unit placeholder glyphs (`<g id="units">`).

After this story the SVG can be opened in any text editor or Inkscape and hand-tweaked without
wading through thousands of lines of machine-generated metadata.

**Why now:** `dipmap/assets/map.svg` is ~2.2 MB because the godip source SVG was copied
verbatim. The file is full of `sodipodi:*` and `inkscape:*` attributes on every element, plus
`<metadata>` / RDF blocks and `<defs>` with Inkscape-specific markers. None of this is
consumed by the bot's render pipeline (`oksvg` ignores unknown namespaced attributes). Removing
it shrinks the file dramatically, speeds up SVG parsing, and makes hand-editing practical.

**Files:** `dipmap/assets/map.svg`, `cmd/mkapsvg/main.go`

---

#### What to strip

| Category | Example | Safe to remove? |
|---|---|---|
| RDF / Dublin Core metadata | `<metadata>…</metadata>` | Yes |
| Inkscape `<sodipodi:namedview>` | `<sodipodi:namedview …/>` | Yes |
| `inkscape:*` attributes on elements | `inkscape:connector-curvature="0"` | **Except** `inkscape:label` — see below |
| `sodipodi:*` attributes on elements | `sodipodi:nodetypes="…"` | Yes |
| `<defs>` containing only Inkscape arrowhead markers | `<defs><marker id="…">…</marker></defs>` | Yes |
| Namespace declarations for stripped prefixes | `xmlns:inkscape="…"` | Yes |
| `<style>` blocks | Already stripped at render time by `stripStyles()` | Yes — remove from source too |

**Keep `inkscape:label`:** The `extractProvinceShape` helper in `dipmap/highlight.go` locates
province polygons using `inkscape:label="<province>"` (e.g. `inkscape:label="vie"`). This
attribute is needed for the zoomed-map feature (Story 10c). Do not remove it.

#### How to implement

Two options — choose one:

**Option A — One-off script:** Write a small Go program in `cmd/stripsvg/` (or add a
`-strip` flag to `cmd/mkapsvg/`) that reads `dipmap/assets/map.svg`, applies regex / XML
cleanup, and overwrites the file. Commit the cleaned file.

**Option B — Pre-process in mkapsvg:** Add the strip pass to the existing `cmd/mkapsvg/main.go`
so that re-running the generator always produces a clean SVG. This is preferred because the tool
and the asset stay in sync.

The cleanup pass must be idempotent (running it twice produces the same output).

#### Acceptance criteria

- `dipmap/assets/map.svg` contains no `<metadata>` block, no `<sodipodi:namedview>`,
  no `sodipodi:*` or `inkscape:*` attributes (other than `inkscape:label`), no `<style>` blocks.
- File size is < 600 KB (from ~2.2 MB).
- `go test -v -cover -race ./...` still passes at 100%.
- `go test -v -tags functional ./bot/` still passes (the render pipeline is unaffected).
- A unit test in `cmd/mkapsvg/main_test.go` (or `cmd/stripsvg/main_test.go`) asserts that the
  output SVG contains no `sodipodi:` or `inkscape:` text (other than `inkscape:label=`).

---

### Story 10b — Map Rendering Fixes: Labels, Unit Colours, and Glyph Geometry

**Goal:** Fix three visible defects in the basic `/map` (no-argument) output:

1. **Province name labels are not visible** — the map renders with no text.
2. **All unit glyphs appear the same grey** — nation colours are not applied.
3. **Unit glyphs are too large; fleet glyphs are not visually distinct enough.**

**Files:** `dipmap/render.go`, `dipmap/overlay.go`, `cmd/mkapsvg/main.go`,
`dipmap/assets/map.svg`, `dipmap/render_test.go`, `dipmap/overlay_test.go`

---

#### Bug 1 — Province name labels not visible

**Root cause:** The `names` layer in `map.svg` contains `<text>` elements whose inline
`style=` attributes reference the font `LibreBaskerville-Bold`. `oksvg` cannot resolve this
font and silently drops every text element that references it.

`render.go` has a `rewriteTextStyles()` helper that is supposed to strip problematic font
properties. Verify whether it is being called in the render pipeline and whether it actually
strips `font-family`. If not, fix it. The rewritten style must reduce to a form oksvg can
handle, e.g.:

```
style="font-size:16px;fill:#000000"
```

Properties to strip: `font-family`, `font-variant-*`, `-inkscape-font-specification`,
`font-variant-ligatures`, `font-variant-caps`, `font-variant-numeric`, `font-variant-east-asian`.

Properties to keep: `font-size`, `fill`, `text-anchor`, `writing-mode`.

`transform="rotate(…)"` attributes on `<text>` elements must be **preserved** — they control
label orientation for diagonal provinces.

Font size should scale with render canvas width: `size = 16 * (canvasWidth / 1524.0)`,
clamped to a minimum of 1 px. This keeps labels legible in both the full board and any
future zoomed renders.

**Acceptance criteria for Bug 1:**
- After pre-processing, `<text>` elements contain no `font-family` in their style.
- A unit test covers the style-rewriting helper on a synthetic `<text>` element.
- The SVG passed to `oksvg` contains the province name text strings for at least the seven
  capitals: London, Paris, Berlin, Moscow, Rome, Vienna, Constantinople.

---

#### Bug 2 — Unit glyphs all appear the same grey colour

**Root cause:** `dipmap.Overlay` calls `setAttr(svg, id, "fill", colour)` which sets the
`fill` attribute on the `<g id="unit-…">` group element. However, the inner `<rect>` element
(generated by `mkapsvg`) has its own hardcoded `fill="#cccccc"` attribute. In SVG, explicit
attributes on child elements take precedence over inherited `fill` from the parent group.
The nation colour is set on the group but overridden by the rect, so all units render grey.

**Fix — two complementary changes:**

1. **In `cmd/mkapsvg/main.go`:** Remove the `fill="#cccccc"` attribute from the `<rect>`
   template so the rect inherits fill from its parent `<g>`. Set the placeholder fill on the
   `<g>` itself instead: `<g id="…" fill="#cccccc" display="none">`. Then `setAttr` on the
   group correctly overrides the placeholder.

2. **Regenerate `dipmap/assets/map.svg`** by running `go run ./cmd/mkapsvg/` and committing
   the result.

`overlay.go` itself does not need to change — the ID-based `setAttr` approach is correct.

**Verify:** After the fix, a unit test in `dipmap/overlay_test.go` should assert that when
`Overlay` is called with `{"par": {Type:"Army", Nation:"France"}}`, the returned SVG contains
`id="unit-par-army"` with both `display="inline"` and `fill="#3399CC"` (France blue).

---

#### Bug 3 — Unit glyphs too large; fleet shape not distinctive

**Root cause:** The glyph dimensions chosen in `cmd/mkapsvg/main.go` are proportional to the
full 1524-wide canvas but visually too large when rendered at typical screen sizes. Fleet
glyphs (`30×18`) look almost identical to army glyphs (`24×24`) — not narrow enough to read
as a ship.

**Fix — update glyph geometry in `cmd/mkapsvg/main.go`:**

Suggested target sizes (tune by looking at the rendered map):
- **Army:** `x="-9" y="-9" width="18" height="18" rx="2"` (reduce from 24×24 to 18×18)
- **Fleet:** `x="-14" y="-6" width="28" height="12" rx="3"` (keep wide but shorten height;
  a 2.3:1 aspect ratio reads more clearly as a hull shape)

Letter sizes: army text `font-size="11"`, fleet text `font-size="8"` (scale proportionally).

**Regenerate `dipmap/assets/map.svg`** after updating the template.

**Acceptance criteria for Bug 3:**
- A unit test in `cmd/mkapsvg/main_test.go` asserts the generated army rect width < 20 and
  fleet rect width ≥ 24 and fleet rect height < army rect height (i.e. fleets are distinctly
  wider-and-shorter than armies).

---

#### Combined acceptance criteria

- `/map` (no args) renders a JPEG that contains visible province names.
- Unit glyphs are coloured by nation (England dark blue, France light blue, Germany grey, etc.).
- Fleet glyphs are visually narrower than army glyphs.
- `go test -v -cover -race ./...` passes at 100% for `dipmap/` and `bot/`.
- `go test -v -tags functional ./bot/` passes.

---

### Story 10c — Zoomed /map with Territory and Radius (deferred)

**Goal:** Implement (or re-enable) the `/map <territory> <n>` variant that crops the board
to a province and its n-hop neighbourhood, highlighted with nation colours.

**Context:** This feature was de-scoped from the initial `/map` work because of two bugs:
1. Province highlights use a random cycling colour palette rather than nation colours.
2. The overall render had other quality issues that needed fixing in Story 10b first.

The command is currently wired but the territory+radius code path is commented out
(`bot/commands.go` always falls through to the full-board render). The functional test
`TestCommand_Map_WithTerritoryAndRadius` is skipped (`t.Skip`).

**Files:** `dipmap/highlight.go`, `dipmap/highlight_test.go`, `bot/commands.go`,
`bot/commands_test.go`, `bot/bot_functional_test.go`

---

#### What needs fixing

**Highlight colours (Bug 3 from original user report):**
The current `Highlight()` function in `dipmap/highlight.go` cycles through a fixed palette
(`#FF6B6B`, `#4ECDC4`, `#45B7D1`, …) unrelated to nations. The zoomed view should instead
highlight each province using the colour of the nation that owns a unit there, and use a
neutral highlight (e.g. `#DDDDDD` at 50% opacity) for provinces with no unit.

`Highlight` needs access to the current game state (units + supply centre ownership) to
determine per-province colours. Its signature will need to change or a new variant created.

**Re-enable the code path in `bot/commands.go`:**
- Restore argument parsing for territory and n.
- Restore the `if territory != "" && n > 0` branch calling `Neighborhood`, `Highlight`,
  and `RenderZoomed`.
- Restore `boardGraph()` and the `godipGraph` adapter.
- Restore the `graph` field on Dispatcher if needed.

**Re-enable tests:**
- Un-skip `TestCommand_Map_WithTerritoryAndRadius` in `bot/bot_functional_test.go`.
- Restore unit tests `TestDispatchMap_RejectsInvalidRadius`,
  `TestDispatchMap_RejectsHighlightError`, `TestDispatchMap_RejectsZoomError`,
  `TestDispatchMap_UsesCustomGraph` (or equivalent) in `bot/commands_test.go`.

#### Acceptance criteria

- `/map lon 1` and `/map lon 3` produce JPEG images with different dimensions (radius 3
  covers a larger bounding box).
- Province highlights are coloured by nation (where a unit is present) or neutral grey
  (empty provinces).
- All re-enabled tests pass.
- `go test -v -cover -race ./...` passes at 100% for `dipmap/` and `bot/`.
- `go test -v -tags functional ./bot/` passes.

---

### Story 10d — Replace oksvg with tdewolff/canvas for proper text rendering

**Goal:** Replace the `oksvg` + `rasterx` rendering pipeline with
`github.com/tdewolff/canvas`, which supports the full SVG text model including
font loading and `transform="rotate(…)"` on `<text>` elements. Switch map output
from JPEG to lossless PNG. Delete the workaround functions that were masking the
root cause.

**Root cause of missing labels:** `oksvg`'s `drawFuncs` map
(`vendor/github.com/srwiley/oksvg/draw.go`) has no handler for `text` or `tspan`
elements — they are silently discarded regardless of their style. The
`rewriteTextStyles` and `prepareForRender` functions in `dipmap/render.go` were
written on the incorrect assumption that simplifying `font-family` would make
`oksvg` render text. They cannot fix this; `oksvg` will never render text.
`tdewolff/canvas` handles `<text>`, `<tspan>`, and `transform="rotate(…)"` on
text elements natively, and resolves fonts from a caller-supplied font family.

**Font:** Embed `LibreBaskerville-Bold.ttf` in the binary via `//go:embed`. The
SVG's `names` layer already references `LibreBaskerville-Bold` by name; register
the embedded TTF under that exact family name so the renderer resolves it without
any style-rewriting. Libre Baskerville is published by Pablo Impallari under the
**SIL Open Font License 1.1**, which permits embedding in software and commercial
distribution. Source: https://github.com/impallari/Libre-Baskerville —
download `LibreBaskerville-Bold.ttf` and commit it to `dipmap/assets/`.

**Output format:** Switch from JPEG to PNG. `tdewolff/canvas` renders to
`image.Image`; encode with `image/png` (lossless). Update `SVGToJPEG` →
`SVGToPNG`, the `imgFn` field on `bot.Dispatcher`, and any MIME-type references
in the platform adapters.

---

#### Dependency setup (requires network — do once, then commit vendor/)

```
go get github.com/tdewolff/canvas@latest
go mod tidy
go mod vendor
```

Remove `oksvg` and `rasterx` from `go.mod` after replacing all usages. Run
`go mod tidy` again to drop them cleanly.

---

#### Files

- `dipmap/render.go` — replace rendering pipeline; delete workaround helpers
- `dipmap/render_test.go` — update tests; assert PNG magic bytes in output
- `dipmap/assets/LibreBaskerville-Bold.ttf` — new embedded font (OFL 1.1)
- `go.mod`, `go.sum`, `vendor/` — add `tdewolff/canvas`; remove `oksvg`, `rasterx`
- `bot/commands.go` — update `imgFn` from `SVGToJPEG` to `SVGToPNG`

---

#### What to delete

The following functions exist solely to work around `oksvg` limitations and must
be deleted once the renderer is replaced:

- `stripStyles` — removed `<style>` blocks because `oksvg` couldn't parse
  embedded data URIs; `tdewolff/canvas` handles CSS natively.
- `rewriteTextStyles` — attempted (ineffectively) to make `oksvg` render text;
  no longer needed.
- `prepareForRender` — converted `display:none` to `opacity:0` because `oksvg`
  does not support the CSS `display` property; `tdewolff/canvas` supports it.
- `loadClassicalSVGWith`, `renderWithLoader` — test helpers for the above; delete
  with their callers.

Keep `loadEmbeddedSVG` (or an equivalent) as the single entry point that reads the
embedded `map.svg`.

---

#### Rendering pipeline after this story

```
loadEmbeddedSVG()          read dipmap/assets/map.svg (embedded)
        │
        ▼
Overlay(svg, units)        activate unit glyphs by ID (unchanged)
        │
        ▼
SVGToPNG(svg)              tdewolff/canvas parses SVG, loads
        │                  LibreBaskerville-Bold.ttf, renders text
        │                  and paths to image.Image, encodes PNG
        ▼
[]byte PNG
```

---

#### Acceptance criteria

- `go test -v -cover -race ./...` passes at 100% for `dipmap/` and `bot/`.
- `go test -v -tags functional ./bot/` passes.
- The bytes returned by `SVGToPNG` begin with the PNG magic number
  (`\x89PNG`); a unit test asserts this.
- The SVG passed to the renderer is **not** pre-processed to strip styles or
  rewrite font references — the font is resolved via the registered font family.
- `stripStyles`, `rewriteTextStyles`, and `prepareForRender` no longer exist in
  `dipmap/render.go`; the file contains no references to `oksvg` or `rasterx`.
- `go.mod` no longer lists `github.com/srwiley/oksvg` or
  `github.com/srwiley/rasterx`.
- A render integration test loads the real `dipmap/assets/map.svg`, calls
  `SVGToPNG`, and asserts the output is a valid non-empty PNG. The test does not
  assert pixel colours (font rendering is deterministic but platform-dependent);
  asserting magic bytes and minimum file size (> 50 KB) is sufficient.

---

### Story 13 — Lambda / EventBridge Deployment

**Goal:** Refactor session scheduling to a `Scheduler` interface and wire up a Lambda entry
point so the bot runs as a stateless FaaS application with externally managed phase deadlines.

**Files:** `session/scheduler.go`, `platform/eventbridge/scheduler.go`, `cmd/lambdabot/main.go`

**Acceptance criteria:**
- `Scheduler` interface defined in `session/scheduler.go`:
  ```
  Schedule(channelID string, at time.Time) error
  Cancel(channelID string) error
  ```
- `LocalScheduler` implementation wraps `time.AfterFunc`; used in tests and server deployments
- `EventBridgeScheduler` implementation creates/deletes one-time AWS EventBridge Scheduler rules
  named by `channelID`; rule target is the Lambda function ARN (from environment variable)
- `Session.timer *time.Timer` and `Session.mu sync.Mutex` replaced by `Session.scheduler Scheduler`
- `GameStarted` and `PhaseResolved` event structs gain `DeadlineAt time.Time` (serialised as RFC3339)
- `AdvanceTurn()` gains an idempotency check: reads game channel history, no-ops if a
  `PhaseResolved` event already exists for the current phase
- `cmd/lambdabot/main.go` handles two event shapes:
  - Platform webhook payload → parse command → `bot.Dispatch`
  - `{"action": "advance_turn", "channel_id": "..."}` → `session.Load()` → `AdvanceTurn()`
- Unit tests cover `LocalScheduler` fire/cancel, idempotency guard (duplicate advance no-ops),
  and Lambda handler routing
- `go test -v -cover -race ./...` passes

---

### Story 11 — Slack Platform Adapter

**Goal:** Deploy the bot as a Slack app.

**Files:** `platform/slack/adapter.go`, `cmd/slackbot/main.go`

**Acceptance criteria:**
- Handles Slack slash command HTTP requests; parses into `bot.Command` values
- Handles Slack Events API payloads (URL verification, event dispatch)
- Posts text responses and PNG images back to Slack channels
- Implements all six `events.Channel` methods on the Slack adapter:
  - `Post` / `History` — Slack reads history via `conversations.history` API (no local store needed)
  - `SendDM` / `DMHistory` — Slack DM channel; history via `conversations.history` API
  - `PostImage` — uploads PNG to group channel via `files.upload`
  - `SendDMImage` — uploads PNG to the player's DM channel via `files.upload`
- Handles DM slash-command payloads (`channel_type = "im"`) and routes them to the order handler
- `cmd/slackbot/main.go` wires up HTTP server, Slack signing-secret verification, and `bot.Dispatch`
- Unit tests cover all Channel methods and webhook parsing
- `go test -v -cover -race ./...` passes

---

### Story 12 — WhatsApp Platform Adapter (optional)

**Goal:** Deploy the bot via the Twilio WhatsApp API or Meta Cloud API.

**Note:** WhatsApp requires a Meta Business Account (approval can take days/weeks) and has
per-conversation charges. Tackle only if Telegram/Slack do not meet deployment needs.

**Files:** `platform/whatsapp/adapter.go`, `platform/whatsapp/store.go`, `cmd/whatsappbot/main.go`

**Acceptance criteria:**
- `WhatsAppChannel` implements `events.Channel`:
  - `Post` / `History` — group messages sent via Twilio API; history backed by local JSONL file store
  - `SendDM` / `DMHistory` — 1:1 messages sent via Twilio API; history backed by local JSONL file store
  - `PostImage` — uploads PNG to Twilio Media API, posts MMS link to group
  - `SendDMImage` — uploads PNG to Twilio Media API, posts MMS link to player's 1:1 thread
- Webhook handler validates `X-Twilio-Signature` and parses `application/x-www-form-urlencoded` payloads
- `cmd/whatsappbot/main.go` reads `TWILIO_ACCOUNT_SID`, `TWILIO_AUTH_TOKEN`,
  `TWILIO_WHATSAPP_NUMBER`, `DATA_DIR`, `PORT` from env
- Unit tests cover all Channel methods using a mock Twilio API server
- `go test -v -cover -race ./...` passes
