# Batch H — Zone editor GUI tests (backlog §5.1 + §5.2)

Close backlog §5.1 (zone-editor pointer flows: drag-to-move, drag-to-connect,
snapping, right-click delete, click-to-place) and §5.2 (zone and connection
property panels: `widget.Editor` and dropdown paths) by driving the **real
window** through an extended `ZoneEditorHandler`, with a golden per action —
the same shape §5.3/batch F gave the file explorer.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is
done, set its status to `Complete` and write its **Phase Summary** (what was
done, key decisions, anything needed to continue with zero context); run the
phase's **Verification Plan** and record the result before moving on. When all
phases are done, fill in **Final Recap** and **Deployment Plan**.

### Standing constraints (AGENTS.md)

- Never modify `data/`, `internal/entities/template/` or `internal/registry/`.
- Cross-platform: `path/filepath`, PowerShell chained with `;`, never `&&`.
- Unit coverage floor **72.5 %**, currently **72.9 %**. This batch adds
  `integration_test,gui` tests, which do **not** count toward that number, so
  the expectation is **flat**, not higher. It must not drop.
- Never stage, never commit. Delete with `Remove-Item`, never `git rm`.
- Never run CI. **Goldens are generated locally on a real GPU only**, and every
  `-update` run is scoped with `-run '<TestName>'` so unrelated goldens in the
  package are not rewritten (backlog §5.4, §5.5).
- Test layout: `go run ./cmd/testlayoutcheck .` must pass. Build tags are
  per file. GPU tests carry `//go:build integration_test && gui`.

### Owner decisions (do not re-litigate)

Asked and answered before any code was written, 2026-08-20:

1. **Batch H = §5.1 + §5.2**, one plan file, six phases.
2. **Drive through `AppRunner` and the real window**, by extending
   `ZoneEditorHandler` the way `FileExplorerHandler` was extended — *not* the
   dialog-direct `newDialogContext` route that backlog §5.4(d) prescribed.
   §5.4(d)'s "reachability-only" note about this handler is hereby superseded.
3. **Switch the topology away from `Random`** before opening the zone editor,
   and **lift the preview-canvas mask** once that switch has happened. The mask
   may only be lifted after the switch.
4. **`ShufflePlayerZones` does not affect the preview canvas** — the owner's
   ruling. It *was* hard-coded `true` in `config.NewGeneratorConfig` and is
   never overridden by `GeneratorConfigMapper.FromEditorState`, so it was live
   in every GUI generation; the ruling is that it permutes labels between
   equivalent slots without changing what is drawn. **Phase 0 proves this
   empirically before any golden is committed.** If the proof fails, stop and
   ask the owner — do not silently re-mask and do not invent a fixture.
   *Superseded after Phase 0:* the owner diagnosed the flag as a bug (the
   shuffle reached the generated template but not the preview layout, so the
   two disagreed) and had it **removed outright**, so the goldens are stable
   for a second reason and the flag no longer exists anywhere in the code base.
5. **Both `Square` and `Geometric Hub`** are used, in different scenarios.
6. **Every handler action verifies a golden**, file-explorer style.
7. **The geometry tests are rewritten onto the handler too**, accepting that
   their pinned coordinates are re-derived from the real dialog geometry.
8. **The five callback-dependent tests stay dialog-direct** (see Phase 2).
9. **Accessors approved**: `Window.TopZoneEditor()` + `IZoneEditorDialog`;
   dialog `CanvasOrigin()`/`CanvasSide()` (which requires storing the canvas
   offset in production state); dialog drag/mode readbacks.

### Re-measuring coordinates

Follow backlog §5.4's procedure verbatim (temporary `*op.Ops` accessor, fresh
`input.Router`, `replay.AppendSemantics`, only `semantic.Button` bounds are
trustworthy, confirm every non-button coordinate by driving it, then delete the
probe). Constants live in
[handlerCoordinates.go](../test/test_helpers/integration_common/handlerCoordinates.go)
with a one-line note recording how each was confirmed.

---

## Phase 0: Prove determinism and measure coordinates

Status: Complete

- [x] Add a temporary `//go:build integration_test && gui` probe test that
      selects a non-`Random` topology on the Layout & Zones tab and captures
      two unmasked full-window frames of the same state; diff them. Repeat for
      `Square` and for `Geometric Hub`.
- [x] Same diff again with the zone editor dialog open, so the dialog's own
      canvas is proven stable and not just the preview panel.
- [x] If either diff shows more than anti-aliasing noise (per-pixel delta ≤ 8
      outside the status and output-directory regions), **stop and report to
      the owner** — decision 4 above is then wrong and the batch needs a new
      determinism strategy.
- [x] Measure the topology dropdown trigger and its option-row pitch on the
      Layout & Zones tab (left column, first section). Confirm by driving:
      click the row, assert `EditorStateDto.Topology`.
- [x] Measure the zone-editor dialog panel rectangle (`PreferredSize` is
      1000×720, `DialogHost` centres it) and the body inset above the canvas.
- [x] Measure the side-panel row coordinates for every editable widget:
      zone Size / Guard x / Weekly + textboxes, Quality and Castles dropdown
      triggers and option rows, connection Type / Guard zone / Guard preset /
      Weekly dropdowns, Guard value and Increment textboxes, the Advanced
      options checkbox and the three rows it reveals, and the Snap checkbox in
      the toolbar. Confirm each by driving it and asserting the state change.
- [x] Delete the probe accessor and probe test.

### Verification Plan

- The two-run diff report is recorded in the Phase Summary with the measured
  maximum per-pixel delta.
- Every measured constant has a driving confirmation named in the summary.
- `go build ./...` and `go vet -tags='integration_test,gui' ./...` clean.

### Phase Summary

Completed 2026-08-20. Two temporary probes were used and both are deleted:
`test/integration/gui/probe_integration_test.go` and
`test/test_helpers/integration_common/appRunnerProbe.go` (which carried
`DumpSemantics`, `DumpSemanticsIn` and `CaptureRaw`). `go build ./...`,
`go vet -tags='integration_test,gui' ./test/...` and
`go run ./cmd/testlayoutcheck .` are clean; `git status --short` shows only
`?? plans/`.

**Determinism — decision 4 holds.** Each pair below is two *independently
constructed* `AppRunner`s taken through the same steps, captured unmasked and
diffed outside the status-message and output-directory regions, three rounds
each:

| pair | differing px | max delta | px over delta 8 | production comparer |
| ---- | ------------ | --------- | --------------- | ------------------- |
| editor / Square, round 1 | 5797 | 154 | 44 | matches, mean 0.0016 %, changed 0.0011 % |
| editor / Square, rounds 2–3 | ~4100 | 4 | 0 | matches, changed 0.0000 % |
| editor / Geometric Hub, rounds 1–3 | ~4800 | 3 | 0 | matches, changed 0.0000 % |
| zone editor / Square, rounds 1–3 | ~3900 | 4 | 0 | matches, changed 0.0000 % |
| zone editor / Geometric Hub, rounds 1–3 | ~5300 | 4 | 0 | matches, changed 0.0000 % |

The ~4000 pixels at delta ≤ 4 are text/edge anti-aliasing spread over the whole
frame. The single outlier is round 1 of the first test executed in the process:
44 pixels at the *rounded-corner* rasterization of the Layout & Zones tab
(x 701–705 / 841–845, y 45–49) and of the topology dropdown trigger
(x 178–180 / 552–554, y 111–113 and 135–137) — a first-headless-window warm-up
artifact on corner pixels, not content. It never recurred in rounds 2 and 3 and
never appeared in any zone-editor pair.

Judged by the committed `snapshot.Comparer` (tolerance 64/255, changed-pixel
budget 0.05 % ≈ 720 px, mean budget 0.25 %) the *worst* pair scores 0.0011 %
changed pixels and 0.0016 % mean — roughly a 45× margin. **No pixel of the
preview canvas or of the zone-editor canvas ever differed by more than 4.**

That also settles `ShufflePlayerZones` empirically: the default template has two
player zones, so an active shuffle would swap `P1`/`P2` with probability ½ per
capture; twelve independent captures produced twelve identical layouts
(p ≈ 1.6 % if the shuffle were reordering what is drawn). Goldens are safe and
no `.gen.json` fixture is needed. The flag has since been removed entirely
(see owner decision 4).

**Measured coordinates (window dp, 1600 × 900, PxPerDp = 1).** Every one was
confirmed by driving it, as noted.

*Layout & Zones tab, topology dropdown*

- Trigger `(178,111)-(555,138)` → centre **(366, 124)**.
- Option rows open directly under it: first row `y 140-165`, pitch **25**,
  x `178-555`. Every row publishes `semantic.Button` with its own label, so
  rows are addressed with `ClickButtonIn(image.Rect(178,138,556,420), name)`
  rather than by coordinate.
- Confirmed by driving: selecting `"Square"` → `CurrentState().Topology ==
  "Square"`; `"Geometric Hub"` → `"GeometricHub"`. (`"Ring"` maps to
  `"Default"` — noted, not in scope.)

*Zone-editor dialog frame*

- Panel ≈ `(300,90)-(1300,810)`; title text at `(314,108)`; close `X` button
  `(1267,104)-(1286,132)`.
- Toolbar buttons, all `semantic.Button` with labels, centre y **152**:
  `Add connection` (366), `Add zone` (459), `Delete selected` (553), `Undo`
  (635), `Revert to Base` (716). Address by label.
- `Snap` checkbox `(776,136)-(843,170)` → centre **(809, 153)**; publishes a
  `semantic.CheckBox` node, no label.
- **Canvas area `(355,181)-(935,761)`** — origin **(355, 181)**, side **580**.
  Confirmed by driving: a click at `(405,232)` selected zone `Spawn-B`, and a
  click at `(640,465)` selected the `Spawn-A → Spawn-B` connection.
- Footer: `Cancel` `(1114,769)-(1172,796)`, `Apply changes`
  `(1180,769)-(1286,796)`. Address by label.

*Side panel* — labels start at x 996; text editors span x `1112-1260`
(click centre x **1186**); dropdown triggers span x `1106-1266` (centre x
**1186**); full-width buttons span x `996-1266` (centre x **1131**).

Player-spawn zone (name label 191, note row 211):

| row | bounds | click point |
| --- | ------ | ----------- |
| `Size` editor | `(1112,252)-(1260,267)` | (1186, 259) |
| `Guard x` editor | `(1112,283)-(1260,298)` | (1186, 290) |
| `Weekly +` editor | `(1112,310)-(1260,325)` | (1186, 317) |
| `Delete this zone` | `(996,341)-(1266,368)` | by label (disabled here) |

Neutral zone (no note row, so **every row sits 29 dp higher** and two dropdowns
are added — confirmed by placing `Neutral-C` with `Add zone` and selecting it):

| row | bounds | click point |
| --- | ------ | ----------- |
| `Size` editor | `(1112,223)-(1260,238)` | (1186, 230) |
| `Guard x` editor | `(1112,254)-(1260,269)` | (1186, 261) |
| `Weekly +` editor | `(1112,281)-(1260,296)` | (1186, 288) |
| `Quality` trigger | `(1106,306)-(1266,333)` | (1186, 319) |
| `Castles` trigger | `(1106,333)-(1266,360)` | (1186, 346) |
| `Delete this zone` | `(996,399)-(1266,426)` | by label |

Connection, advanced options collapsed:

| row | bounds | click point |
| --- | ------ | ----------- |
| `Type` trigger | `(1106,217)-(1266,244)` | (1186, 230) |
| `Guard zone` trigger | `(1106,244)-(1266,271)` | (1186, 257) |
| `Guard preset` trigger | `(1106,275)-(1266,302)` | (1186, 288) |
| `Guard value` editor | `(1112,308)-(1260,323)` | (1186, 315) |
| `Weekly +` trigger | `(1106,333)-(1266,360)` | (1186, 346) |
| `Increment` editor | `(1112,366)-(1260,381)` | (1186, 373) |
| `Advanced options` checkbox | `(996,396)-(1266,430)` | (1131, 413) |
| `Delete this connection` | `(996,443)-(1266,470)` | by label |

Connection with `Advanced options` checked (confirmed by clicking (1131,413)
and seeing the three rows appear):

| row | bounds | click point |
| --- | ------ | ----------- |
| `Match group` editor | `(1112,439)-(1260,454)` | (1186, 446) |
| `Guard escape` checkbox | `(996,463)-(1266,497)` | (1131, 480) |
| `Sim turn squad` checkbox | `(996,503)-(1266,537)` | (1131, 520) |
| `Delete this connection` | `(996,550)-(1266,577)` | by label |

**Two layout rules the handler must respect (both observed, not assumed):**

1. An **open dropdown pushes every row below it down** — opening `Guard preset`
   moved `Guard value` from y 308 to y 460 and `Delete this connection` from
   y 443 to y 702. Option rows are `semantic.Button`s with their labels
   (`Default (30000)`, `Weakest (10000)`, `Low (22000)`, `Medium (34000)`,
   `High (46000)`, `Very High (58000)`, pitch 25 starting at y 304), so a
   dropdown is always driven as *click trigger → `ClickButtonIn` the option by
   label*, and no other row coordinate may be used while a dropdown is open.
2. The zone panel has **two coordinate variants**, player-spawn and neutral,
   differing by the 29 dp note row. `handlerCoordinates.go` therefore carries
   both sets and the handler picks by whether the selected zone is a player
   spawn.


## Phase 1: Accessors and handler build-out

Status: Not started

- [ ] Production: store the canvas offset alongside the existing `side` in
      `zoneEditorGeometryState`, assigned in
      [layoutCanvas](../app/gui/dialogs/zoneEditorCanvas.go#L32). One
      assignment, no behaviour change.
- [ ] `zoneEditorDialog_testexports.go`: `CanvasOrigin()`, `CanvasSide()`, and
      the drag/mode readbacks a golden cannot express (dragged zone, pending
      connection source, snap-guide state).
- [ ] `window_testexports.go`: `IZoneEditorDialog` + `TopZoneEditor()`,
      mirroring `IFileExplorerDialog` / `TopFileExplorer` exactly, including the
      comment explaining why the contract is declared there.
- [ ] `LayoutAndZonesTabHandler.SelectTopology(name)` — opens the dropdown,
      clicks the row, commits, snapshots, clears `isRandomTopology` and lifts
      the preview-canvas mask.
- [ ] `ZoneEditorHandler`: canvas actions (`ClickCanvasAt`, `DragZone`,
      `DragFromZoneTo`, `RightClickEdge`, `ClickEmptyCanvas`), toolbar actions
      by label, side-panel actions (`TypeZoneSize`, `SelectZoneQuality`, …),
      a `Dialog()` observation surface, and a golden per action.
- [ ] Canvas-to-window coordinate mapping in one place, backed by a permanent
      calibration test that presses a zone's mapped point and asserts that zone
      became the selection.

### Verification Plan

- `go run ./cmd/testlayoutcheck .` → `test-layout check passed`.
- `go build ./...`, `go vet ./...`, `go vet -tags='integration_test,gui' ./...`
  all clean; `gofmt -l` empty on touched files.
- The calibration test passes, proving the mapping rather than asserting it.

### Phase Summary

_(write when phase completes)_

## Phase 2: Rewrite the existing zone-editor tests onto the handler

Status: Not started

- [ ] Move the behaviour tests in
      [zoneEditorDialog_integration_test.go](../test/integration/gui/zoneEditorDialog_integration_test.go)
      onto the handler: button labels, Undo, Revert to Base, Apply, delete
      selected, mode toggles, selection rendering — asserting through
      `Window.CurrentState()` / the state driver and a golden per action.
- [ ] Re-derive the pins in
      [zoneEditorGeometry_integration_test.go](../test/integration/gui/zoneEditorGeometry_integration_test.go)
      from the real dialog geometry and assert them through
      `TopZoneEditor()`. Keep them **exact float** assertions: their job is to
      guard batch G's single-rounding property, so a pin that rounds is a
      regression in the test, not a tidier number.
- [ ] Keep dialog-direct, unchanged, the five tests the window cannot reach:
      the two `RevertToBaseCannotRegenerate` cases, the two Apply/RevertToBase
      flag cases, `TestWhenRevertToBaseFailed_TheApplyReportsNoRevert`, and
      `TestWhenZoneEditorDialogRenders_UsesHandlerProvidedOptions`. Add a
      one-line comment on the file's fixture explaining why they stay.
- [ ] Generate the goldens locally, `-run`-scoped per test.

### Verification Plan

- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`
  passes with no `-update`.
- A second run of the same command passes, proving the goldens are stable
  rather than freshly written.

### Phase Summary

_(write when phase completes)_

## Phase 3: §5.1 pointer flows

Status: Not started

- [ ] `zoneEditorPointer_integration_test.go`
      (`//go:build integration_test && gui`):
      - `TestWhenAZoneIsDraggedToANewPosition_TheAppliedLayoutRecordsIt`
      - `TestWhenAZoneIsDraggedNearAGuide_ItSnapsToTheGuide`
      - `TestWhenADragStartsOnAZoneInAddConnectionMode_AConnectionIsCreated`
      - `TestWhenADragEndsOnEmptyCanvas_NoConnectionIsCreated`
      - `TestWhenACurveIsRightClicked_ThatConnectionIsDeleted`
      - `TestWhenAddZoneModeIsArmedAndEmptyCanvasIsClicked_AZoneIsPlaced`
- [ ] Drags must clear the 6 px dead zone
      (`zoneDragDeadZonePx`) or they are only selections.
- [ ] Snap assertions use **fractional** post-batch-G coordinates: the grid
      correction is no longer rounded away.
- [ ] Goldens generated locally, `-run`-scoped.

### Verification Plan

- The gui suite passes twice in a row with no `-update`.
- Each test asserts a state change, not only a golden: a golden alone cannot
  distinguish "the drag did the right thing" from "nothing happened".

### Phase Summary

_(write when phase completes)_

## Phase 4: §5.2 property panels

Status: Not started

- [ ] `zoneEditorProperties_integration_test.go`
      (`//go:build integration_test && gui`), covering all five approved
      groups:
      - zone textboxes: Size, Guard x, Weekly +
      - zone dropdowns: Quality, Castles (the `ApplyZoneEditorQuality` reprofile
        path)
      - connection guard value typed, and a non-numeric value rejected
      - connection dropdowns: Type, Guard zone, Guard preset, Weekly
      - the Advanced options checkbox and the Match group / Guard escape /
        Sim turn squad rows it reveals
- [ ] Note in the file that **the zone name is a read-only label**, so §5.2's
      `TestWhenAZoneNameIsTyped_…` was not written and why.
- [ ] Focus each field with a click before typing, per `InputText`'s contract.
- [ ] Goldens generated locally, `-run`-scoped.

### Verification Plan

- The gui suite passes twice in a row with no `-update`.
- `Size` clamping (0.1–2.0, two decimals) and the non-numeric rejection are
  asserted on committed state, not on widget text.

### Phase Summary

_(write when phase completes)_

## Phase 5: Documentation and gates

Status: Not started

- [ ] [todo/test_observations.md](../todo/test_observations.md): delete the
      "Still uncovered: the property panels' `widget.Editor`/dropdown paths and
      the pointer flows … still future work" sentence and replace it with a
      pointer to the new tests; record anything that genuinely stayed
      uncovered, with the reason.
- [ ] [todo/backlog-opus5.md](../todo/backlog-opus5.md): mark §5.1 and §5.2
      ✅ DONE and self-contained, update the header item counts and the §8
      batch table with row **H**.
- [ ] Correct backlog §5.4(d)'s claim that the zone-editor handler is
      reachability-only, since batch H made it a driving handler.
- [ ] Full gate run (see below) and a carry-forward document.

### Verification Plan

Every gate green, matching the batch-G baseline:

- `go build ./...`; `go vet ./...`; `go vet -tags='integration_test,gui' ./...`
- `go run ./cmd/testlayoutcheck .`
- `go test ./test/unit/... -count=1`
- `go test -tags=integration_test ./test/integration/... -count=1`
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`
- Unit coverage ≥ 72.5 % (expected flat at 72.9 %)
- `golangci-lint-v2 run ./... --issues-exit-code=0` → 0 issues
- `gofmt -l ./app ./internal ./test ./cmd` empty

### Phase Summary

_(write when phase completes)_

## Final Recap

_(write when all phases complete)_

## Deployment Plan

_(write when all phases complete)_
