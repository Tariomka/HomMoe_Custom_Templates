# Session Carry-Forward

## 1. Session goal

Implement backlog **batch G** — `todo/backlog-opus5.md` §2.3, *preview geometry
is integer-only although a `Vec2` already exists* — under a plan file, per
AGENTS.md §2.4.

## 2. Fixes applied

- **Two contradictory rounding rules in the same package.** `commitPositions`
  rounded (`image.Pt(int(math.Round(...)))`) while `canvasMetrics.center()`
  truncated, both in
  [layoutGeometry.go](../internal/services/preview_service/layoutGeometry.go).
  Both are gone: positions are committed as raw `float64` and `center()` returns
  `data.NewVec2(this.cx, this.cy)`. The same `image.Pt(int(math.Round(...)))`
  pattern was removed from
  [layoutRingHub.go](../internal/services/preview_service/layoutRingHub.go) and
  [layoutBalancedRings.go](../internal/services/preview_service/layoutBalancedRings.go).
- **Pointer input was truncated before hit-testing.**
  [zoneEditorCanvas.go](../app/gui/dialogs/zoneEditorCanvas.go) did
  `image.Pt(int(e.Position.X), int(e.Position.Y))`, so hit-testing and snapping
  saw the pixel the pointer landed in rather than where it actually was. The
  `f32` position now flows through as `data.Vec2[float64]`.
- **`helpers.GetPointOnQuadraticBezierCurve` was dead.** Its three call sites all
  migrated to the pre-existing float variant `GetVectorOnQuadraticBezierCurve`,
  so it was deleted from [math.go](../internal/helpers/math.go).
- **Duplicated selection-ring drawing.** `drawNodes` had two identical blocks
  that would each have needed the new rounding; extracted as
  `drawSelectionRing`.

## 3. Features added / changed

No user-visible behaviour changed, but two outputs shifted by design (see §8).

- **The layout pipeline is `float64` end to end.** `preview.Layout.Positions` /
  `ZoneRadius`, `preview.Zone.Center`, `preview.Connection.Start/Ctrl/End`,
  `ZoneEditorGeometry.Positions` / `ZoneRadius`, `ZoneEditorEdge.MidPoint` and
  `ZoneEditorSnapResult.Position` are all float. Leaving any of them integer
  would have re-quantised one step downstream and defeated the item.
- **Rounding happens exactly once per output.**
  [app/gui/utils/draw.go](../app/gui/utils/draw.go) is the single rounding site
  for the Gio canvases; the pixel loop in
  [assetProvider.go](../internal/services/asset_provider/assetProvider.go) is the
  one for the generated PNG. `drawAsset` deliberately keeps a **fractional**
  centre — the bilinear sample resolves it, so sub-pixel placement survives into
  the image.
- **`utils.ToF32Point`** — new [app/gui/utils/math.go](../app/gui/utils/math.go),
  the `Vec2[float64]` → `f32.Point` bridge. It lives in `app/gui/utils` and not
  beside `Vec2` so that no package under `internal/` imports Gio, which is still
  true.
- **`zoneEditorGeometryService` lost conversions rather than gaining them.**
  `HitTestNode`, `HitTestEdge`, `SnapPosition`, `GridStep`, `buildEdges`,
  `obstacleBulge` and `otherZoneGuides` already computed in `float64` via
  `math.Hypot`; the change mostly deletes the round-trip through `image.Point`.
  `GridStep` no longer passes through an int, and `SnapPosition` no longer
  rounds its result.

**Constraint discovered:** `internal/models` imports `internal/models/preview`,
so `preview` **cannot** import `models`. The `preview` types therefore spell the
type `data.Vec2[float64]` while everything outside spells it `models.Position`.
`models.Position = data.Vec2[float64]` is a type alias, so they are identical —
this is cosmetic only, but it explains the mixed spelling.

## 4. File modifications

Models and DTOs:

| File | Change |
| --- | --- |
| [internal/models/preview/previewLayout.go](../internal/models/preview/previewLayout.go) | `Positions map[string]data.Vec2[float64]`, `ZoneRadius float64`. |
| [internal/models/preview/previewZone.go](../internal/models/preview/previewZone.go) | `Center` → `data.Vec2[float64]`. |
| [internal/models/preview/previewConnection.go](../internal/models/preview/previewConnection.go) | `Start`, `Ctrl`, `End` → `data.Vec2[float64]`. |
| [internal/models/zoneEditorGeometry.go](../internal/models/zoneEditorGeometry.go) | `Positions map[string]Position`, `ZoneRadius float64`. |
| [internal/models/zoneEditorEdge.go](../internal/models/zoneEditorEdge.go) | `MidPoint` → `data.Vec2[float64]`. |
| [internal/models/zoneEditorSnapResult.go](../internal/models/zoneEditorSnapResult.go) | `Position` → `models.Position`. |
| [internal/dtos/zoneEditorHitTestRequestDto.go](../internal/dtos/zoneEditorHitTestRequestDto.go), [zoneEditorSnapRequestDto.go](../internal/dtos/zoneEditorSnapRequestDto.go) | Float positions and radius. |

Services and handlers:

| File | Change |
| --- | --- |
| [layoutGeometry.go](../internal/services/preview_service/layoutGeometry.go) | Both rounding rules removed; `generatorCoords` still converts the protected `*[2]float64` at the read site. |
| [layoutRingHub.go](../internal/services/preview_service/layoutRingHub.go), [layoutBalancedRings.go](../internal/services/preview_service/layoutBalancedRings.go) | Float placement, float `ZoneRadius`. |
| [previewLayoutService.go](../internal/services/preview_service/previewLayoutService.go) | `buildPreviewConnections` takes float positions; ctrl point computed with `Subtract`/`Add`/`MultiplyScalar`. |
| [previewGeneratorService.go](../internal/services/preview_service/previewGeneratorService.go), [assetFitter.go](../internal/services/preview_service/assetFitter.go) | Float draw calls; rounding moved to the brush rect. |
| [assetProvider.go](../internal/services/asset_provider/assetProvider.go) | `DrawPlayerZone` / `DrawNeutralZone` / `DrawArenaMarker` / `drawAsset` take float centres. |
| [zoneEditorGeometryService.go](../internal/services/connection_editor/zoneEditorGeometryService.go) + its interface | Fully float; `image` import dropped. |
| [zoneEditorHandler.go](../internal/handlers/zoneEditorHandler.go), [guiHandler.go](../internal/handlers/guiHandler.go) + [zoneEditorHandlerInterface.go](../internal/handlers/handler_interfaces/zoneEditorHandlerInterface.go) | Forwarding layers follow the new signatures. |
| [internal/helpers/math.go](../internal/helpers/math.go) | `CalculatePointTowards` takes/returns `Vec2[float64]`; `GetPointOnQuadraticBezierCurve` **deleted**. |

GUI:

| File | Change |
| --- | --- |
| [app/gui/utils/math.go](../app/gui/utils/math.go) | `ToF32Point`. |
| [app/gui/utils/draw.go](../app/gui/utils/draw.go) | The single rounding site; `DrawConnection` / `DrawPreviewZone` / `drawCurve` / `drawOffsetText` take floats. |
| [app/gui/dialogs/zoneEditorCanvas.go](../app/gui/dialogs/zoneEditorCanvas.go) | Pointer truncation dropped; `drawSelectionRing` extracted. |
| [app/gui/dialogs/zoneEditorSnap.go](../app/gui/dialogs/zoneEditorSnap.go), [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go), [zoneEditorInteractionState.go](../app/gui/dialogs/zoneEditorInteractionState.go) | Float drag/press positions and normalised manual positions. |
| [app/gui/dialogs/zoneEditorDialog_testexports.go](../app/gui/dialogs/zoneEditorDialog_testexports.go) | `ZonePositions`, `CanvasZoneRadius`, `HitTest*`, `SnapDraggedPosition` follow the new types. |
| [app/gui/panels/previewPanel.go](../app/gui/panels/previewPanel.go) | **No edit needed** — it already forwarded `ZoneRadius`, which changed type underneath it. |

Docs / planning:

| File | Change |
| --- | --- |
| [plans/batch-g-float-preview-geometry.md](../plans/batch-g-float-preview-geometry.md) | **New.** Five phases, all `Complete`, each with a Phase Summary, plus Final Recap and Deployment Plan. |
| [todo/backlog-opus5.md](../todo/backlog-opus5.md) | §2.3 marked ✅ and self-contained; §8 batch **G** row and the header counts updated. |

## 5. Tests added or updated

Added:

- [test/unit/internal/helpers/math/calculatePointTowards_test.go](../test/unit/internal/helpers/math/calculatePointTowards_test.go)
  and [getVectorOnQuadraticBezierCurve_test.go](../test/unit/internal/helpers/math/getVectorOnQuadraticBezierCurve_test.go) —
  **neither function had any test at all**, and `CalculatePointTowards` changed
  signature, so AGENTS.md §2.3 required them.
- `TestWhenTwoZonesAreLessThanAPixelApart_TheirCentresDiffer` in
  [buildPreviewLayout_test.go](../test/unit/internal/services/preview_service/previewLayoutService/buildPreviewLayout_test.go) —
  the test named by the backlog item. Two manual zones 0.3 px apart, a case the
  integer layout collapsed onto one pixel.
- `TestWhenOnlyTheGridIsInReach_TheDraggedPositionKeepsTheFractionalCorrection`
  in [snapPosition_test.go](../test/unit/internal/services/connection_editor/zoneEditorGeometryService/snapPosition_test.go) —
  pins that the grid correction is no longer rounded away.

Updated: the `previewLayoutService`, `zoneEditorGeometryService`, `guiHandler`,
`zoneEditorHandler`, `previewHandler`, `assetProvider` and `previewLayoutCache`
unit folders; [zoneEditorGeometry_integration_test.go](../test/integration/gui/zoneEditorGeometry_integration_test.go)
and [zoneEditorDialog_integration_test.go](../test/integration/gui/zoneEditorDialog_integration_test.go);
and the [zoneEditorGeometryServiceMock](../test/test_helpers/zoneEditorGeometryServiceMock.go) /
[templateHandlerMock](../test/test_helpers/templateHandlerMock.go).

Assertion policy applied as decided: **exact literals in integration tests**,
**`InDelta` at 1e-9 in unit tests**.

**No goldens were regenerated** — see §8.

**Gate results — all green:**

| Gate | Result |
| --- | --- |
| `go build ./...` | clean |
| `go vet ./...` and `go vet -tags='integration_test,gui' ./...` | clean |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test ./test/unit/... -count=1` | pass |
| `go test -tags=integration_test ./test/integration/... -count=1` | pass |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | pass |
| Unit coverage | **72.9 %** (flat; floor is 72.5 %) |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **0 issues** |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `wire diff ./internal/composition/...` | no diff |

## 6. Git status snapshot

Branch: `AD/fixing_some_stuff_08-12`. **Working tree clean.**

```
c09dc8b (HEAD -> AD/fixing_some_stuff_08-12) Batch G done
3b0e689 (origin/AD/fixing_some_stuff_08-12) Docs
d37d1fa Batch F done
```

The owner reviewed, staged and committed batch G as `c09dc8b`, making his own
edits to roughly eighteen of the touched files before committing. It is **not
pushed** — `origin` is at `3b0e689`, one commit behind. Nothing unstaged is
inherited.

Per AGENTS.md §2.5 and the owner's standing instruction — *"leave the staging
area alone entirely"* — the agent never staged or committed anything, and used
`Remove-Item` rather than `git rm`. Do not run any git staging command.

## 7. Rejections / things the owner declined

- Nothing was declined this session. Seven scope questions were asked up front
  and all seven answered before any code was written; they are recorded verbatim
  in the plan file's *Owner decisions* section and must not be re-litigated.
- The owner approved a **wider** scope than §2.3 literally specifies: all four
  downstream fields (`preview.Zone.Center`, `Connection.Start/Ctrl/End`,
  `ZoneEditorEdge.MidPoint`, `ZoneEditorSnapResult.Position`) were floated too.

## 8. Open questions

- **The backlog's snapshot prediction was wrong, in our favour.** §2.3 warned
  that the GPU suite would need `-update` and owner eyeballing. It did not: the
  suite passed unchanged on the first run, because the preview canvas is masked
  by the harness and the zone-editor handler takes no snapshots. **No `-update`
  was ever run and none of the 62 goldens moved.** The §5.3 lesson (scope
  `-update` with `-run`) was therefore never exercised and still stands untested.
- **Two values changed, both intended and both verified by hand.** A snap that
  returned `x = 201` now returns `x = 200 + 6/7` — the grid step is `2·38/7` and
  the leading edge at `x = 162` puts the 15th grid line 6/7 px right of 200. The
  ring zero-angle slot yields `47.99999999999997` instead of `48`, so that one
  assertion became `InDeltaSlice`.
- **`.gen.json` gains sub-pixel precision.** `zoneEditorDialog.go` normalises
  zone centres into the persisted `manualPosition` field, which now carries
  fractional values. The field was already `float64`, so old and new files stay
  mutually readable — but the owner accepted this side effect knowingly, and it
  is worth remembering if a diff of saved files ever looks noisy.
- `Zone.GeneratorPosition *[2]float64` was **not** touched; converting it is
  backlog §2.4, owner-gated, and now unblocked by this batch.
- Nothing is blocked.

## 9. Next recommended actions

1. Push `AD/fixing_some_stuff_08-12` when the owner is ready — batches E, F and
   G are all on it and `origin` is one commit behind.
2. Start **batch H** (backlog §5.1 + §5.2, zone-editor pointer and
   property-panel tests). It was gated on §2.3 and is now unblocked; write it
   against the post-batch-G **fractional** coordinates, not the old integers.
3. Batch **I** (§2.1, `EditorStateDto` rework) is the next large one and needs
   its own `plans/` file — multi-phase, twelve packages.

## 10. Carry-forward prompt

> Read `AGENTS.md` first, then `todo/backlog-opus5.md`.
>
> Hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/` without explicit
> approval; everything must build and run on both Windows and Linux (use
> `path/filepath`, and chain PowerShell with `;`, never `&&`); every change
> ships with tests and unit coverage must not drop below 72.5 % (currently
> 72.9 %); durable multi-session work gets a plan file under `plans/`; never
> stage and never commit — I review, stage and commit myself, so leave the
> staging area alone entirely, and delete files with `Remove-Item`, never
> `git rm`; never change where `.rmg.json` is written and never persist the
> output directory; never run a bulk in-place rewrite over the repository;
> never run CI and never generate snapshot goldens in CI.
>
> Batch **G** (backlog §2.3, float preview geometry) is **complete, green and
> committed** as `c09dc8b` on `AD/fixing_some_stuff_08-12`, unpushed. The
> preview and zone-editor pipeline is now `float64` end to end and rounds
> exactly once — at `app/gui/utils/draw.go` for the Gio canvases and at the
> pixel loop in `assetProvider` for the PNG. No goldens moved.
>
> Next up is batch **H** (backlog §5.1 and §5.2 — zone-editor pointer tests and
> property-panel tests). It was gated on §2.3 and is now unblocked. Write it
> against the post-batch-G **fractional** coordinates: pointer positions are no
> longer truncated, snapped positions keep their fractional grid correction, and
> `SnapPosition` returns e.g. `200 + 6/7` where it used to return `201`. When
> generating goldens, scope `-update` with `-run` so unrelated goldens in the
> package are not rewritten. Before starting, prompt me to confirm the item and
> surface every open question first.
>
> Full handoff: `./.agent/session-carry-forward.md`.
