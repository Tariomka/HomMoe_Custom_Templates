# Batch G â€” Float Preview Geometry

Move the preview/zone-editor layout pipeline from integer pixels to `float64`
throughout, so that positions, radii and Bezier control points are rounded
exactly once â€” at the final draw call â€” instead of being quantised by every
layout strategy before anything consumes them.

Backlog item: `todo/backlog-opus5.md` Â§2.3. Execution slot: batch **G** in Â§8,
which must land **before** Â§5.1 (batch H, zone-editor pointer tests) because
those tests will be written against the post-Â§2.3 coordinates.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

### Standing constraints (AGENTS.md)

- `data/`, `internal/entities/template/` and `internal/registry/` are **read-only**.
  In particular `Zone.GeneratorPosition *[2]float64` stays exactly as it is â€”
  converting it is backlog Â§2.4 and is owner-gated. Convert at the read site
  (`generatorCoords`).
- Windows + Linux both must build; `path/filepath` only; chain PowerShell with `;`.
- Never stage, never commit, never `git rm` (use `Remove-Item`).
- Unit coverage floor **72.5 %**; the figure entering this batch is **72.9 %**.
- Never run a bulk in-place rewrite over the repository.
- Scope any golden `-update` with `-run` so unrelated goldens are not rewritten.

### Owner decisions taken before work started (2026-08-19)

All seven were confirmed explicitly; do not re-litigate them.

1. Batch G proceeds, with this plan file.
2. **All four downstream integer fields become float**: `preview.Zone.Center`,
   `preview.Connection.Start/Ctrl/End`, `ZoneEditorEdge.MidPoint`,
   `ZoneEditorSnapResult.Position`. Leaving any of them integer would re-quantise
   one step after `Layout.Positions` and defeat the item.
3. `app/gui/utils/draw.go` becomes the **single rounding site** for the Gio
   preview canvas: `DrawConnection` / `DrawPreviewZone` take floats and round or
   convert to `f32.Point` internally.
4. `Vec2[float64]` â†’ `f32.Point` conversion lives in **`app/gui/utils`**, not in
   `internal/helpers/data`. No package under `internal/` imports Gio today and
   that stays true.
5. Pointer input is passed through as float: drop the
   `image.Pt(int(e.Position.X), int(e.Position.Y))` truncation in
   `zoneEditorCanvas.go` so hit-testing and snapping see the real position.
6. Assertion policy: **exact literals in integration tests**, **tight
   `InDelta` (1e-9) in unit tests**.
7. Goldens: run the GPU suite and regenerate **only the goldens that actually
   fail**, scoped with `-run`. The preview canvas is masked and the zone-editor
   handler takes no snapshots, so few or none are expected to move.

### Known side effect, accepted

`zoneEditorDialog.go` normalises zone centres into the persisted
`manualPosition` field of `.gen.json`. Those values gain sub-pixel precision.
The field is already `float64`, so old and new files stay mutually readable.

---

## Phase 1: Float the preview and editor model types
Status: Complete

- [x] `internal/models/preview/previewLayout.go`: `Positions map[string]models.Position`, `ZoneRadius float64`.
- [x] `internal/models/preview/previewZone.go`: `Center models.Position`.
- [x] `internal/models/preview/previewConnection.go`: `Start`, `Ctrl`, `End` â†’ `models.Position`.
- [x] `internal/models/zoneEditorGeometry.go`: `Positions map[string]models.Position`, `ZoneRadius float64`.
- [x] `internal/models/zoneEditorEdge.go`: `MidPoint` â†’ `data.Vec2[float64]`.
- [x] `internal/models/zoneEditorSnapResult.go`: `Position` â†’ `models.Position`.
- [x] Add whatever `Vec2` arithmetic the conversion needs (e.g. `Length`) to `internal/helpers/data/vec2.go`; do **not** add a Gio-typed method.
- [x] `internal/services/preview_service/`: stop rounding in `commitPositions`, `canvasMetrics.center()`, `layoutRingHub.go`, `layoutBalancedRings.go`. `generatorCoords` keeps reading the protected `*[2]float64` and converting at the read site.
- [x] Confirm the two conflicting rules (`math.Round` in `commitPositions` vs. truncation in `center()`) are both gone rather than merely unified.

### Verification Plan
- `go build ./...` clean.
- `go vet ./...` clean.
- Compile errors outside `internal/models` and `internal/services/preview_service` are expected at this point and are Phase 2/3 work.

### Phase Summary

All six model types are float. **Key constraint discovered:** `internal/models`
imports `internal/models/preview`, so `preview` cannot import `models` — the
`preview` types therefore use `data.Vec2[float64]` directly while everything
outside `preview` uses the `models.Position` alias. They are the same type
(`models.Position = data.Vec2[float64]` is a type alias), so this is purely
cosmetic at the call sites.

No new `Vec2` arithmetic was needed: `Subtract`/`Add`/`MultiplyScalar`/
`DivideScalar`/`SquaredLength`/`DotProduct` already covered every case, with
`math.Hypot` used where a true length was required.

Both conflicting rounding rules are gone: `commitPositions` stores raw floats
and `canvasMetrics.center()` returns `data.NewVec2(this.cx, this.cy)` with no
truncation. `layoutRingHub.go` and `layoutBalancedRings.go` lost every
`image.Pt(int(math.Round(...)))`. `generatorCoords` still reads the protected
`*[2]float64` and converts at the read site, untouched.

## Phase 2: Float through the zone-editor geometry service and handlers
Status: Complete

- [x] `internal/services/connection_editor/zoneEditorGeometryService.go`: `HitTestNode`, `HitTestEdge`, `SnapPosition`, `GridStep`, `buildEdges`, `obstacleBulge`, `otherZoneGuides` take/return float geometry. Most of this deletes conversions â€” the bodies already compute in `float64` via `math.Hypot`.
- [x] Its interface file and the DTOs that carry positions/radius.
- [x] `internal/handlers/guiHandler.go` and `internal/handlers/zoneEditorHandler.go` forwarding layers.
- [x] `test/test_helpers/*Mock.go` for any interface whose signature moved.

### Verification Plan
- `go build ./...` clean.
- `go vet -tags='integration_test,gui' ./...` clean.

### Phase Summary

`zoneEditorGeometryService.go` is fully float and, as predicted, the change was
mostly *deletion* of conversions — the bodies already computed in `float64`.
`GridStep(zoneRadius float64) float64` now returns `zoneRadius * 2.0 / 7`
without an intermediate int. `SnapPosition` no longer rounds its result, which
is the behavioural change the batch exists for.

Also converted: the service interface, `ZoneEditorHitTestRequestDto`,
`ZoneEditorSnapRequestDto`, `zoneEditorHandler.go`, `guiHandler.go`, the handler
interface, and both `test/test_helpers/templateHandlerMock.go` and
`test/test_helpers/zoneEditorGeometryServiceMock.go`. Six files lost a
now-unused `image` import.

**Dead code found and removed:** `helpers.GetPointOnQuadraticBezierCurve` had
three call sites, all of which migrated to the pre-existing float variant
`GetVectorOnQuadraticBezierCurve`, so it was deleted.

## Phase 3: Round only at the draw boundaries
Status: Complete

- [x] `app/gui/utils/draw.go`: `DrawConnection`/`DrawPreviewZone` accept float radius and float centres; build rectangles and `f32.Point`s internally. Add the local `Vec2[float64]` â†’ `f32.Point` helper here.
- [x] `app/gui/panels/previewPanel.go`: pass floats straight through.
- [x] `app/gui/dialogs/zoneEditorCanvas.go`: drop the pointer truncation; float node centres, selection outlines (`ZoneRadius + 4`) and label offsets; round only where Gio demands `image.Rectangle`.
- [x] `app/gui/dialogs/zoneEditorSnap.go` and `zoneEditorDialog.go` (normalised manual positions).
- [x] `internal/services/preview_service/previewGeneratorService.go` + `internal/services/preview_service/assetFitter.go` + `internal/services/asset_provider/assetProvider.go`: round immediately before `image.Point` / `image.Rectangle` / pixel loops.
- [x] `app/gui/dialogs/zoneEditorDialog_testexports.go`: `EditedZones`/`CanvasZoneRadius` follow the new types.

### Verification Plan
- `go build ./...` and `go vet -tags='integration_test,gui' ./...` clean.
- `go run ./cmd/testlayoutcheck .` prints `test-layout check passed`.
- Launch is not required; the GUI suite in Phase 5 exercises the canvases.

### Phase Summary

`app/gui/utils/draw.go` is the single rounding site for the Gio preview canvas,
with a comment saying so. The `Vec2[float64]` → `f32.Point` helper landed as
`utils.ToF32Point` in the new `app/gui/utils/math.go`; `internal/` still imports
no Gio package.

`previewPanel.go` needed no edit — it already forwarded `layout.ZoneRadius`,
which simply changed type underneath it.

In `zoneEditorCanvas.go` the pointer truncation is gone: the raw `f32` pointer
position becomes a `data.Vec2[float64]` so hit-testing and snapping see the
exact position rather than the pixel it lands in. Two duplicated selection-ring
blocks would each have needed the new rounding, so they were extracted into
`drawSelectionRing`.

`previewGeneratorService.go`, `assetFitter.go` and `assetProvider.go` round only
at the pixel loop; `assetProvider.drawAsset` deliberately keeps a fractional
centre because the bilinear sample resolves it, preserving sub-pixel placement.

## Phase 4: Update the geometry tests
Status: Complete

- [x] Unit: `test/unit/internal/services/preview_service/previewLayoutService/buildPreviewLayout_test.go` (~40 assertions) â†’ tight `InDelta` (1e-9).
- [x] Unit: the four `test/unit/internal/services/connection_editor/zoneEditorGeometryService/*` folders.
- [x] Unit: the `guiHandler` and `zoneEditorHandler` folders, plus `test/unit/app/gui/models/previewLayoutCache/get_test.go`.
- [x] New test named by Â§2.3: `TestWhenTwoZonesAreLessThanAPixelApart_TheirCentresDiffer`, in the `buildPreviewLayout` folder.
- [x] Integration: `test/integration/gui/zoneEditorGeometry_integration_test.go` â€” update the ~12 pins (radius `38`; `A:(140,350) B:(560,350) C:(350,140)`; control points `(350,368)`/`(350,332)`; label midpoint `(350,359)`; obstacle control `(350,274)`; snap outputs `(200,355)`/`(201,350)`) to **exact** new literals. Do not relax to `InDelta`.
- [x] Integration: `test/integration/gui/zoneEditorDialog_integration_test.go` snapping setup.

### Verification Plan
- `go test ./test/unit/... -count=1` passes.
- `go test -tags=integration_test ./test/integration/... -count=1` passes.

### Phase Summary

All listed suites converted and green. Two expectations changed *value*, both
because rounding was removed and both verified by hand:

- Snap: dragging to `(200, 355)` against a guide zone now returns
  `x = 200 + 6/7` instead of `201`. The grid step is `2*38/7`; the leading edge
  at `x = 162` puts the 15th grid line 6/7 px right of 200. Both the unit test
  and `zoneEditorGeometry_integration_test.go` were repinned.
- Ring layout: `TestWhenZeroAngleZoneIsSet_RotatesThatZoneToFirstRingSlot`
  produced `y = 47.99999999999997` where the integer pipeline gave `48`. Per
  decision 6 this became an `InDeltaSlice` at 1e-9.

Tests added beyond the plan: `CalculatePointTowards` and
`GetVectorOnQuadraticBezierCurve` had **no** tests at all and
`CalculatePointTowards` changed signature, so
`test/unit/internal/helpers/math/calculatePointTowards_test.go` and
`getVectorOnQuadraticBezierCurve_test.go` were written (AGENTS §2.3).

The named sub-pixel test `TestWhenTwoZonesAreLessThanAPixelApart_TheirCentresDiffer`
places two manual zones 0.3 px apart — a case the integer layout collapsed onto
one pixel — and asserts their centres differ.

## Phase 5: Gates, coverage and scoped goldens
Status: Complete

- [x] `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`; regenerate only failing goldens, each scoped with `-run`.
- [x] Coverage task; record the figure and compare against the 72.9 % baseline (floor 72.5 %).
- [x] `golangci-lint-v2 run ./... --issues-exit-code=0` â†’ 0 issues.
- [x] `gofmt -l ./app ./internal ./test ./cmd` â†’ empty (fix via `gofmt -w` on an **explicit** file list only).
- [x] `wire diff ./internal/composition/...` â†’ no diff.
- [x] Update `todo/backlog-opus5.md`: mark Â§2.3 done and self-contained, and the Â§8 batch **G** row.

### Verification Plan
- Every row of the backlog Â§9 baseline table reproduced green.

### Phase Summary

Every gate green:

| Gate | Result |
| --- | --- |
| `go build ./...` | clean |
| `go vet ./...` | clean |
| `go vet -tags='integration_test,gui' ./...` | clean |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test ./test/unit/... -count=1` | pass |
| `go test -tags=integration_test ./test/integration/... -count=1` | pass |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | pass |
| Unit coverage | **72.9 %** — unchanged from baseline, floor 72.5 % |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | 0 issues |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `wire diff ./internal/composition/...` | no diff |

**No goldens were regenerated.** The GPU suite passed on the first run, which
confirms the prediction in owner decision 7: the preview canvas is masked by the
harness and the zone-editor handler takes no snapshots, so the sub-pixel shift
is invisible to the 62 stored goldens. No `-update` was ever run.

One linter finding surfaced and was fixed: `testifylint` flagged
`assert.Equal(t, fixtureZoneRadius, geometry.ZoneRadius)` in
`buildGeometry_test.go` once the radius became a float; it is now `InDelta`.

## Final Recap

The preview and zone-editor layout pipeline is `float64` end to end. Positions,
zone radii, Bezier control points, edge midpoints, snap results and pointer
input all stay fractional from the layout strategies through the services,
DTOs and handlers, and are rounded exactly once at the draw boundary —
`app/gui/utils/draw.go` for the Gio canvas and the pixel loop in
`assetProvider` for the generated PNG.

The two contradictory rounding rules that motivated the item
(`math.Round` in `commitPositions` versus truncation in `canvasMetrics.center()`)
no longer exist. `helpers.GetPointOnQuadraticBezierCurve` became dead and was
removed, and `helpers.CalculatePointTowards` / `GetVectorOnQuadraticBezierCurve`
gained their first unit tests.

Gio stays confined to `app/`: the `Vec2[float64]` → `f32.Point` bridge is
`app/gui/utils.ToF32Point`, not a method on `Vec2`.

Behavioural deltas, both intended: snapped positions keep their fractional grid
correction, and normalised `manualPosition` values written to `.gen.json` gain
sub-pixel precision. The field is already `float64`, so old and new files remain
mutually readable.

Not touched, as required: `data/`, `internal/entities/template/`,
`internal/registry/`, and `Zone.GeneratorPosition` (backlog §2.4, owner-gated).
Nothing was staged or committed.

## Deployment Plan

This is a pure library/UI refactor with no schema, output-path or configuration
change, so deployment is an ordinary build.

1. Review the working tree (`git status --short`) and stage/commit as you see
   fit — the agent left the staging area untouched.
2. Re-run the gate table above on the review machine, in particular
   `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`,
   which needs a GPU and must not be run in CI.
3. `go build ./...` on both Windows and Linux; the change introduces no
   OS-specific code and no new dependency.
4. Ship the binary as usual. No migration is required: `.rmg.json` output is
   byte-compatible and existing `.gen.json` files load unchanged.
5. Batch **H** (backlog §5.1, zone-editor pointer tests) is now unblocked and
   should be written against the post-§2.3 fractional coordinates.
