# Batch 15 — Zone-editor and dialog state extraction (review §2.6)

Break up the undifferentiated state in the zone-editor dialog and the four other
oversized dialogs: move pure geometry and non-GUI logic into `internal/services`,
consolidate the interaction state that legitimately stays in the GUI, and make
the "Reset to generated" button honest. Closes the last open finding in
[todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md).

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../AGENTS.md) first. This batch touches only `app/gui/**` and
`internal/services/**` — none of the protected directories (`data/`,
`internal/entities/template/`, `internal/registry/`). Never stage, never commit.

**This batch is large (~2 800 LOC across 9 files) and will span sessions.** Phases
1–4 are independently shippable; do not start a phase you cannot finish.

---

## Owner decisions (already made — do not re-litigate)

1. **Scope is all nine files** — the five `zoneEditor*` files *plus*
   `bonusPickerDialog.go`, `pickerDialog.go`, `ruleDialog.go`, `zoneContent.go`.
2. **Pure geometry moves to a service; selection/drag state stays in the GUI**
   but is consolidated into a named struct. The review's literal wording ("extract
   the pure-geometry *and selection* state into `internal/services/`") was
   **deliberately not followed** for selection: those 15 fields are pointer-event
   view state (`dragPos`, `pressPos`, `zoneDragName`, `addMode`, `pendingFrom`,
   `hint`…), and a service holding `image.Point` drag offsets inverts the layering
   instead of fixing it.
3. **The geometry code is a full service** — interface + constructor + `wire` DI
   registration, following every existing service convention.
4. **Extract everything identified** in the four extra dialogs (~270 LOC of
   candidates), not just the high-branching parts.
5. **Characterization tests come first** (Phase 0). Two of the nine files have no
   test of any kind; the refactor must be verifiable before it starts.
6. **§0.2 "Reset to generated" is fixed by RELABELLING**, not by retaining a
   pristine generated template. The owner explicitly rejected adding a second
   retained template copy.
7. **After the revert button is used, the following Apply clears the persisted
   manual-edit snapshot outright.** The owner accepted the stated consequence:
   the live template keeps showing the reverted-to edits until the next
   regeneration, at which point they disappear.

---

## Starting facts (verified 2026-08-07 at `51e5858`)

### Line counts

| File | LOC |
| --- | ---: |
| `app/gui/dialogs/zoneEditorDialog.go` | 507 |
| `app/gui/dialogs/zoneEditorCanvas.go` | 479 |
| `app/gui/dialogs/bonusPickerDialog.go` | 434 |
| `app/gui/dialogs/pickerDialog.go` | 371 |
| `app/gui/dialogs/ruleDialog.go` | 314 |
| `app/gui/dialogs/zoneContent.go` | 299 |
| `app/gui/dialogs/zoneEditorConnectionProps.go` | 164 |
| `app/gui/dialogs/zoneEditorSnap.go` | 156 |
| `app/gui/dialogs/zoneEditorZoneProps.go` | 95 |

`zoneContent.go` is in **`app/gui/dialogs/`**, not `app/gui/panels/` — the review
and the old backlog both said `panels/`. They were wrong.

### `ZoneEditorDialog` field census — 67 fields, not the ~58 the review claimed

| Bucket | Count | Notes |
| --- | ---: | --- |
| (a) Gio widget handles | 26 | Must stay in the GUI. Nothing to do here. |
| (b) Pure geometry / layout | 8 | `positions`, `previewZones`, `radius`, `side`, `geometrySide`, `edges`, `snapGuideX`, `snapGuideY` |
| (c) Selection / drag / interaction | 15 | Stays in GUI, gets consolidated (Phase 2) |
| (d) Domain data | 8 | `zones`, `originalZones`, `playerZones`, `topology`, `tuning`, `generateRoads`, `working`, `original` |
| (e) Lifecycle flags | 1 | `geometryDirty` |
| (f) Other | 9 | 5 embedded state structs, `onApply`, 2 handlers, `guardPresetValues` |

The five embedded state structs declare **zero methods** — all 42 methods hang off
`ZoneEditorDialog` itself, which is the actual reason the state reads as one blob.

### The extractable geometry — ~160 LOC, all currently at 0 % coverage

| Method | LOC | File |
| --- | ---: | --- |
| `recomputeGeometry` | 65 | canvas |
| `obstacleBulge` | 37 | canvas |
| `hitTestEdge` | 25 | canvas |
| `groupConnectionsByPair` | 19 | canvas |
| `hitTestNode` | 13 | canvas |
| `otherZoneGuides` | 11 | snap |
| `gridStep` | 3 | snap |

All are **PURE** (no `layout.Context`, no drawing, no widget access). Note
`recomputeGeometry` calls `previewHandler` for the layout — that dependency comes
with it.

### Non-GUI candidates in the four extra dialogs (~270 LOC)

| File | Candidate | LOC |
| --- | --- | ---: |
| `bonusPickerDialog.go` | `buildEntries` validation + `BonusEntry` mapping | ~53 |
| | duplicate/spell-ID extraction in constructor | ~12 |
| | `spellCountLabel`, `isNumeric` | ~18 |
| `pickerDialog.go` | item/spell/value-override source→entry mapping | ~47 |
| | filtered/grouped row-model build inside `getRowWidgets` | ~28 |
| | `selectedIDs` ordered extraction | ~8 |
| `ruleDialog.go` | `upsertFromEditor` + `buildRuleFromEditor` | ~43 |
| `zoneContent.go` | `rowDisplayName`, `defaultContentRules`, `ruleMarkers` | ~42 |
| | alphabetical mapping sort in constructor | ~17 |
| | count clamping / row creation in `Add` | ~13 |

### Existing services to extend, not duplicate

`internal/services/connection_editor/` already holds `ConnectionEditorService`
(97), `ManualReapplyService` (200) and `ZoneEditorService` (284) with their
interfaces. The new geometry service is the **fourth** implementation in that
package, so under AGENTS.md §4.2.2 its interface stays **in-package** (the
threshold for a `_interfaces` subpackage is 5). Register the provider in
`EditorSet` in [providerSets.go](../internal/composition/providerSets.go#L53-L60).

### Reaching the service from `app/` — use the handler the dialog already holds

`app/` must not import `internal/services` (depguard `no-services-from-app`,
widened in Batch 14). `NewZoneEditorDialog` **already receives**
`zoneHandler handler_interfaces.IZoneEditorHandler`. Add the geometry methods to
that existing interface and delegate from `zoneEditorHandler`, rather than
threading a new dependency through
[layoutPanelZones.go](../app/gui/panels/layoutPanelZones.go#L101-L123). This is
the same flat-handler pattern Batch 14 used for `IRegenerationHandler`.

### Test coverage today — the main risk

- **No unit tests** reference any of the nine files.
- Only three GUI snapshot tests touch them, all initial-render only:
  - `TestWhenZoneEditorDialogRenders_UsesHandlerProvidedOptions`
    ([zoneEditorDialog_integration_test.go](../test/integration/gui/zoneEditorDialog_integration_test.go))
  - `TestWhenManageRulesDialogHasVariantRule_RendersContent` and
    `TestWhenZoneContentDialogRenders_PreservesSavedRules`
    ([contentRuleDialogs_integration_test.go](../test/integration/gui/contentRuleDialogs_integration_test.go))
- **`bonusPickerDialog.go` and `pickerDialog.go` have no test of any kind.**

### §0.2 — why a true "reset to generated" was rejected

Traced this session. `handleConnectionEditorClick` passes
`State.lastTemplate.Variants[0].Zones/.Connections` into the dialog;
`handleUpdateTemplate` then **overwrites** `lastTemplate` with the edited version
([stateManualEdits.go](../app/gui/drivers/stateManualEdits.go#L26-L55)).
`EditorState.previous` holds only the input `EditorStateDto`, never generated
output. **The pristine generated layout is destroyed on first apply and cannot be
reconstructed.** Making the label true would require a new retained template copy
on `State`; the owner declined. Hence the relabel.

The in-code marker to remove is at
[zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go#L213):
`// Reset only resets current edits, not all manual edits, need to fix eventually. This is a todo, Just don't want to trigger the linter`
— it is the only TODO/FIXME/HACK-class comment in the non-test Go tree, so
deleting it should keep that property.

---

## Phase 0: Characterization safety net
Status: Complete

Pin current behaviour *before* changing anything. Every test here must pass
against unmodified code; if one does not, you have found a real bug — stop and
report it rather than encoding it.

- [x] Record baselines: unit coverage %, `golangci-lint-v2 run ./...` issue count,
      and the full green run of every suite in AGENTS.md §7.
- [x] Add GUI snapshot tests under `test/integration/gui/` for
      `bonusPickerDialog` and `pickerDialog` (currently untested): initial render
      of each, plus one render with a non-empty selection.
- [x] Extend the zone-editor snapshot coverage beyond initial render: one
      snapshot with a connection selected, one with a zone selected, one with
      snap enabled and a drag guide active.
- [x] Add unit tests for the geometry methods **in place**, before they move, by
      exercising them through the dialog's public surface where reachable. Where
      they are not reachable, note it and rely on the snapshots — do **not** add
      test-only seams to production code (AGENTS.md §4.6).
- [x] Capture a manual-behaviour note in this file for the four interactions that
      snapshots cannot assert: add-connection drag, add-zone placement, delete
      selected, revert.

### Deviation: numeric geometry pins instead of pixel goldens

The golden-image machinery in
[test/test_helpers/integration_common/appRunnerSnapshots.go](test/test_helpers/integration_common/appRunnerSnapshots.go)
renders through `AppRunner.captureScreenshot`, which lays out **the whole editor
window** (`this.App.Layout`). There is no per-dialog capture path, and reaching
the zone editor through a real window means going through
`GenerateTemplate`, whose zone list is not deterministic — the resulting goldens
would be flaky for reasons that have nothing to do with geometry.

The safety net is therefore built on **exact numeric geometry assertions** over a
deterministic fixture instead. Every zone in the fixture carries a
`ManualPosition`, which makes `dispatchClusterLayout` take the
`layoutManualPositions` branch: canvas position is `round(p × side)` verbatim.
At `side = 700` the preview metrics scale is exactly `1.0`, so the zone radius
settles at the unscaled maximum of `38` and every expected coordinate is an
integer that can be written down by hand.

This is *stronger* than a pixel diff for the Phase 1 criterion "must not move a
single pixel": the control points, label midpoints, hit-test results and snap
results are asserted to the unit, and a mismatch names the exact quantity that
changed instead of a percentage of differing pixels.

**Phase 1 must therefore read its verification bullet as: the Phase 0 geometry
assertions in
[test/integration/gui/zoneEditorGeometry_integration_test.go](test/integration/gui/zoneEditorGeometry_integration_test.go)
still pass unchanged.** Do not weaken or re-baseline them; if one fails, the
extraction changed behaviour.

### Manual-behaviour note: interactions the harness cannot drive

`layoutCanvas` registers its pointer area *inside* an `op.Offset` push, and every
mutation is driven by `pointer.Press` / `Drag` / `Release` events routed to
`&this.canvasTag`. The dialog tests lay out a bare `layout.Context` and cannot
inject coordinate-addressed pointer events into that transform, so two of the
four interactions stay manual. The other two turned out to be reachable through
the toolbar `Clickable`s and are now automated:

| Interaction | Status |
| --- | --- |
| Delete selected (connection **and** zone) | **Automated** — `TestWhenTheSelectedConnectionIsDeleted_ItLeavesTheWorkingSet`, `TestWhenTheSelectedZoneIsDeleted_ItsConnectionsGoWithIt` |
| Revert ("Reset to generated") | **Automated** — `TestWhenTheEditorIsResetToGenerated_TheDeletedConnectionComesBack`, `..._TheSelectionIsCleared` |
| Add-connection drag | **Manual** — see checklist below |
| Add-zone placement | **Manual** — see checklist below |

Run this by hand (`go run .` → Zones tab → *Edit zones manually*) after any change
to `zoneEditorCanvas.go`, `zoneEditorSnap.go` or `zoneEditorDialog.go`:

1. **Add-connection drag.** Click *Add connection*; the button highlights and the
   status line reads "Add mode: press a zone and drag to another to connect."
   Press a zone and drag: a straight rubber band follows the cursor and the
   source zone gets a highlight ring. Release on a different zone: a new curve
   appears with a dot at its midpoint (the user-added marker), the side panel
   switches to that connection's properties, and add mode **stays on** so the
   next drag works without re-clicking. Release on empty canvas or on the source
   zone: nothing is created. Click empty canvas: add mode turns off.
2. **Add-zone placement.** Click *Add zone*; the status line reads "Add zone
   mode: click an empty spot to place a zone." Click empty canvas: a new neutral
   zone appears exactly where clicked (clamped to 4 %–96 % of the canvas), it
   becomes the selection, the status line reads "Added *N* — connect it with
   'Add connection'.", and every other zone stays put (`ensureManualPositions`
   froze them). Repeat clicks keep placing zones. Click an existing zone: add
   zone mode turns off. Exhaust the label pool: the status line reads "Zone label
   pool exhausted - cannot add more zones." and nothing is placed.
3. **Drag to move + snap.** Drag a zone a few pixels: nothing happens until the
   pointer passes the 6 px dead zone, then the zone follows. With *Snap* on, a
   faint dot grid appears behind the graph, the zone's edges/centre stick to
   nearby grid intersections and to other zones' edge/centre extensions, and a
   thin green guide line is drawn across the canvas for each held axis. With
   *Snap* off, no grid, no guides, free movement.
4. **Right-click delete.** Right-click directly on a connection curve: it is
   removed immediately, without selecting it first.

### Verification Plan
- `go test -tags 'integration_test,gui' ./test/integration/gui/... -count=1` passes
  with the new snapshots committed as the baseline.
- `go run ./cmd/testlayoutcheck .` prints `test-layout check passed`.
- Coverage recorded in the Phase Summary as the number Phase 4 must not drop below.

### Phase Summary

**Baselines recorded at HEAD `c0499e2`, all re-confirmed after the new tests:**

| Check | Result |
| --- | --- |
| `go build ./...` | clean |
| `go test ./test/unit/... -count=1` | pass |
| **Unit coverage total** | **69.3 %** ← the number Phase 4 must not drop below |
| Per-tree statement coverage | `internal` 96.1 %, `app` 19.5 %, `app/gui/dialogs` **2.0 %** |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test -tags=integration_test ./test/integration/...` | pass |
| `go test -tags 'integration_test,gui' ./test/integration/gui/...` | pass |
| `go test '-bench=.' -run=xxx ./test/performance/...` | pass |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **`0 issues.`** |
| GUI integration test count | **19 → 61** |

**What was added:**

- [app/gui/dialogs/bonusPickerDialog_testexports.go](app/gui/dialogs/bonusPickerDialog_testexports.go)
  and [app/gui/dialogs/pickerDialog_testexports.go](app/gui/dialogs/pickerDialog_testexports.go)
  — accessors for the two dialogs that had no test of any kind.
- [test/integration/gui/bonusPickerDialog_integration_test.go](test/integration/gui/bonusPickerDialog_integration_test.go)
  (17 tests) and
  [test/integration/gui/pickerDialog_integration_test.go](test/integration/gui/pickerDialog_integration_test.go)
  (10 tests) — every bonus kind, every validation message, duplicate detection,
  spell sub-picker round-trip, search filtering, flat vs grouped row counts, and
  the `sid=…` value-override formatting.
- [app/gui/dialogs/zoneEditorDialog_testexports.go](app/gui/dialogs/zoneEditorDialog_testexports.go)
  — geometry, hit-test, snap and selection accessors plus toolbar click helpers.
- [test/integration/gui/zoneEditorGeometry_integration_test.go](test/integration/gui/zoneEditorGeometry_integration_test.go)
  (13 tests) — the Phase 1 safety net, described under *Deviation* above.
- Six tests appended to
  [test/integration/gui/zoneEditorDialog_integration_test.go](test/integration/gui/zoneEditorDialog_integration_test.go)
  — connection-selected render, zone-selected render, snap-guide-overlay render,
  delete-connection, delete-zone, and the two revert cases.

**Exact values now pinned** (fixture: zones at `(140,350)`, `(560,350)`,
`(350,140)`; `side = 700`; radius `38`):

- Manual positions land verbatim on the canvas.
- Parallel edges spread symmetrically: control points `(350,368)` and `(350,332)`
  (± the 18 px `bulgeGap`, halved around the pair's midline).
- Grouping is first-seen order: `ab`, `ba`, `ac` — the reversed `B→A` connection
  is canonicalised into the `A–B` bucket.
- Label midpoint of the first edge: `(350,359)`.
- A zone 14 px off the chord pushes the curve to control point `(350,274)`.
- `hitTestNode` hits at the centre, misses at `radius + 1`.
- `hitTestEdge` hits on the curve, misses 38 px off it.
- `gridStep` = `2 × radius / 7`.
- Snap off: `(200,355)` unchanged. Snap on: `(200,355) → (201,350)` — the Y axis
  holds a zone guide, the X axis falls through to the grid.

**Naming trap discovered:** [cmd/testlayoutcheck](cmd/testlayoutcheck) matches
test-only exports **by identifier name across the whole tree**, not by receiver
type. Declaring `func (this *ZoneEditorDialog) Zones()`, `Connections()` or
`ZoneRadius()` made the checker flag 46 unrelated unit tests that merely
reference `variant.Zones` or `layout.ZoneRadius`. The accessors were renamed to
`EditedZones`, `EditedConnectionNames` and `CanvasZoneRadius`. **Any new
`*_testexports.go` accessor must use a name that does not collide with a common
field name.**


---

## Phase 1: The geometry service
Status: Complete

- [x] Create `internal/services/connection_editor/zoneEditorGeometryService.go`
      plus `zoneEditorGeometryServiceInterface.go` (interface in-package — this is
      implementation #4, under the §4.2.2 threshold of 5).
- [x] Move the seven pure methods listed above. Keep them pure: inputs are zones,
      connections, positions and sizes; outputs are geometry values. No
      `layout.Context`, no Gio types. `image.Point` is stdlib and fine.
- [x] `recomputeGeometry` depends on the preview layout — inject that dependency
      through the constructor rather than reaching for a handler.
- [x] Register the provider in `EditorSet`; run the *"Go: Generate wire injectors"*
      task and verify with `wire diff` (exit 0).
- [x] Extend `IZoneEditorHandler` + `zoneEditorHandler` with delegating methods so
      the dialog reaches the service without importing `internal/services`.
- [x] Update the canvas and snap files to call through the handler.
- [x] Full unit tests for every moved method — this is the coverage win that
      justifies the batch (~160 LOC moving from 0 % to covered).

### Verification Plan
- `go build ./...`; `go vet ./...` and `go vet -tags='integration_test,gui' ./...`.
- `wire diff ./internal/composition/...` exits 0.
- `go test ./test/unit/... -count=1` passes; the new service's coverage is ≥ 90 %.
- **The Phase 0 geometry assertions still pass unchanged** (see the Phase 0
  *Deviation* section — numeric pins replace pixel goldens). Geometry extraction
  must not move a single pixel. If an assertion changes, the extraction changed
  behaviour.

### Phase Summary

The canvas geometry now lives in `internal/services/connection_editor` behind
`IZoneEditorGeometryService` (`BuildGeometry`, `HitTestNode`, `HitTestEdge`,
`GridStep`, `SnapPosition`), constructed with `IPreviewLayoutService` injected.
Supporting value types went to `internal/models` (`ZoneEditorEdge`,
`ZoneEditorGeometry`, `ZoneEditorSnapResult`) and `internal/dtos`
(`ZoneEditorGeometryRequestDto`, `ZoneEditorHitTestRequestDto`,
`ZoneEditorSnapRequestDto`). `app/gui/dialogs/connectionEdgeGeometry.go` and
`connectionPairKey.go` were deleted; `zoneEditorCanvas.go` and
`zoneEditorSnap.go` are now thin call-throughs on `IZoneEditorHandler`, which
gained five delegating methods (mirrored on `GUIHandler`). The provider is in
`EditorSet` and `wire diff` exits 0.

Key decisions worth carrying forward:

- **Pointer identity is preserved by index, not by pointer.** The service works
  on connection *values*, so `ZoneEditorEdge` carries a `ConnectionIndex` into
  the slice handed to `BuildGeometry`; the dialog resolves it back to
  `this.working[i]` through the nil-guarded `edgeConnection` helper. Selection
  comparison and `deleteConnection` therefore behave exactly as before.
- **`float32` stayed at the GUI boundary.** The service speaks
  `data.Vec2[float64]`; `zoneEditorDialog_testexports.go` converts with
  `toCanvasPoint`, so the Phase 0 `EdgeGeometry` shape and every numeric pin
  were untouched — which is what makes them a genuine regression gate.
- **The dead `err != nil` branch in `recomputeGeometry` was dropped** (the
  preview handler never returned a non-nil error).
- **`ZoneEditorDialog` lost its `previewHandler` field and constructor
  parameter**; `layoutPanelZones.go` and two integration tests were updated.
- **`snapBool` stayed in the GUI.** `SnapPosition` only guards on
  `zoneRadius <= 0`; the "is snapping switched on" check is view state.
- **Three test-only exports were renamed** to dodge the `testlayoutcheck`
  name-collision trap: `HitTestNode`/`HitTestEdge`/`GridStep` on
  `ZoneEditorDialog` became `HitTestCanvasNode`/`HitTestCanvasEdge`/
  `CanvasGridStep`, because the checker matches identifiers by name across the
  whole tree and the new service methods share those names.

Verification: `go build ./...`, `go vet ./...` and
`go vet -tags='integration_test,gui' ./...` clean; `wire diff` exit 0;
`go run ./cmd/testlayoutcheck .` → `test-layout check passed`;
`go test ./test/unit/... -count=1` all `ok`;
`go test -tags=integration_test ./test/integration/... -count=1` `ok`;
`go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` `ok`
— **the Phase 0 geometry pins pass unchanged, so the extraction moved no
pixels**; `gofmt -l` clean apart from the two known permanent entries;
`golangci-lint-v2 run ./...` → `0 issues.`

Coverage: every function of the new service is ≥ 92.9 % (most 100 %), the five
new handler methods are 100 %, and total unit coverage rose from **69.3 % to
71.0 %**.

---

## Phase 2: Consolidate the interaction state
Status: Complete

- [x] Group the 15 selection/drag/interaction fields into one named struct that
      stays in `app/gui/dialogs/` — it is view state, by owner decision 2.
- [x] Give that struct the small number of methods that operate on it alone
      (select/clear/begin-drag/end-drag), so the reset path has one obvious place
      to call instead of assigning eleven fields inline as
      [resetToOriginal](../app/gui/dialogs/zoneEditorDialog.go#L365-L383) does now.
- [x] Do **not** move it to `internal/services`.
- [x] Re-home the geometry fields (bucket b) that Phase 1 made redundant; delete
      any that the service now owns.

### Verification Plan
- Snapshots unchanged; unit and integration suites green.
- `ZoneEditorDialog`'s direct field count is materially lower than 67 — record the
  new number in the Phase Summary.

### Phase Summary

`app/gui/dialogs/zoneEditorCanvasState.go` — the struct that mixed geometry with
interaction — is **deleted** and replaced by two focused ones:

- [zoneEditorInteractionState.go](../app/gui/dialogs/zoneEditorInteractionState.go)
  holds the 12 selection/drag/mode fields plus the 13 methods that operate on
  them alone: `selectConnection`, `selectZoneNamed`, `clearSelection`,
  `hasSelection`, `toggleAddConnectionMode`, `toggleAddZoneMode`,
  `exitAddModes`, `beginConnectionDrag`, `finishConnectionDrag`,
  `beginZoneDrag`, `endZoneDrag`, `zoneDragLeftDeadZone` and `reset`. The magic
  `6` px drag threshold is now the named `zoneDragDeadZonePx`. It stays in
  `app/gui/dialogs/` — nothing moved to `internal/services` (owner decision 2).
- [zoneEditorGeometryState.go](../app/gui/dialogs/zoneEditorGeometryState.go)
  collapses `positions`, `previewZones`, `radius` and `edges` into the single
  `models.ZoneEditorGeometry` the Phase 1 service already returns, leaving
  `geometry`, `side` and `geometryDirty`. `recomputeGeometry` is now one
  assignment instead of four.

`geometrySide` is **gone**: `layoutCanvas` computes `sideChanged := this.side != side`
before overwriting `this.side`, which is the same test with one field fewer.

The eleven-assignment tail of `resetToOriginal` is now
`this.reset()` + the two property-panel sync caches, and the twelve other inline
mutation sites (both toolbar toggles, `onPress`, `onRelease`, `addConnection`,
`deleteConnection`, the toolbar's `hasSelection`) call the named methods.

**Field count: 67 → 63.** The reduction is entirely bucket (b): six geometry
fields became two. Bucket (c) was *consolidated*, not removed — that is what
owner decision 2 asked for; those 15 fields now live behind an intention-revealing
API instead of being assigned ad hoc from nine different call sites. 26 of the
remaining 63 are Gio widget handles that cannot move.

Behaviour was preserved deliberately rather than assumed: `exitAddModes` clears
`pendingFrom`/`dragging` at the add-zone call site where they are provably
already zero, and leaves `hint` alone (only the two toolbar toggles clear it, as
before); `deleteZone` still nils `selected` without touching `selectedZone`
unless the names match.

Verification: `go build ./...`, `go vet ./...` and `go vet -tags='integration_test,gui' ./...`
clean; `go run ./cmd/testlayoutcheck .` → `test-layout check passed`; unit suite
green; `go test -tags=integration_test ./test/integration/...` ok 2.169s;
`go test -tags 'integration_test,gui' ./test/integration/gui/...` ok 2.952s with
**every Phase 0 geometry pin and every snapshot unchanged**;
`gofmt -l` clean apart from the two permanently-CRLF files;
`golangci-lint-v2 run ./...` → `0 issues.`

Five new GUI integration tests cover the add-mode toggles, which had no test at
all: entering add-connection mode, leaving it on a second click, add-zone
cancelling add-connection (both directions asserted separately), and reset
clearing the armed mode. They use two new test-only accessors,
`AddConnectionModeActive` / `AddZoneModeActive` (name-collision-checked against
the whole tree first).

Unit coverage reads **70.9 %** against Phase 1's 71.0 %. Nothing became less
tested: the 0.1 pp is denominator dilution from adding ~30 statements to
`app/gui/dialogs`, a package the unit suite deliberately does not cover (Gio view
code is covered by the GUI integration suite). It remains well above the 69.3 %
Phase 0 floor.

---

## Phase 3: The four extra dialogs
Status: Complete

Extract **everything identified** in the candidates table (owner decision 4). Work
one file at a time and keep each file's snapshot green before moving on.

- [x] `bonusPickerDialog.go` — validation, `BonusEntry` mapping, duplicate/spell-ID
      extraction, `spellCountLabel`, `isNumeric`.
- [x] `pickerDialog.go` — source→entry mapping, filtering/grouping row model,
      `selectedIDs`.
- [x] `ruleDialog.go` — `upsertFromEditor`, `buildRuleFromEditor`.
- [x] `zoneContent.go` — `rowDisplayName`, `defaultContentRules`, `ruleMarkers`,
      the constructor's alphabetical sort, `Add`'s clamping.
- [x] Decide placement per AGENTS.md §4.4 as you go: content-rule logic belongs
      near `internal/services/content_rules/`; picker/bonus mapping may warrant its
      own package. Note the choice and its reason in the Phase Summary.
- [x] Unit tests for every extracted unit.

### Verification Plan
- Snapshots for all four dialogs unchanged.
- Each extracted unit has its own `test/unit/` folder per §4.6.
- `golangci-lint-v2 run ./...` still reports `0 issues.`

### Phase Summary

All four dialogs are now pure rendering code; every decision they used to make
lives in a service behind a handler seam.

**Placement decisions and their reasons**

- **`bonusPickerDialog.go` → `internal/services/bonuses` (`IBonusEntryService`),
  exposed through a new `IBonusHandler` folded into `IGuiHandler`.** The dialog
  had no handler at all and `app/**` may not import `internal/services`
  (depguard), so a handler facet had to be introduced rather than reused. The
  service owns validation, `config.BonusEntry` construction, duplicate
  filtering and the spell-count label; the three request/result/existing shapes
  travel as DTOs (`BonusCompositionRequestDto`, `BonusCompositionResultDto`,
  `ExistingBonusesDto`).
- **`ruleDialog.go` + `zoneContent.go` → `internal/services/zone_content`
  (`IZoneContentEditorService`), exposed through a new `IZoneContentHandler`
  that *embeds* `IContentRuleHandler`.** Widening `IContentRuleHandler` itself
  was tried and reverted: it forced edits at 14 existing test call sites and
  changed `NewContentRuleHandler`'s signature for no behavioural gain. Embedding
  keeps that constructor and its tests untouched, and because both GUI views
  already sat on a single content-rule seam only the field's *annotation* had to
  widen (`ruleDialog.go`, `zoneContent.go`, `zoneContentDialog.go`,
  `layoutPanel.go`). The service owns rule composition, upsert, defaults, marker
  joining, display naming, alphabetical sorting and count clamping.
- **`pickerDialog.go` → `internal/services/pickers` (`IPickerEntryService`)
  behind `IPickerHandler`.** The GUI's `pickerEntry` struct was deleted in favour
  of `dtos.PickerEntryDto`, and the visible-row model moved wholesale into
  `dtos.PickerRowDto`. The catalogs stay in `app/gui/constants` (which imports
  `internal/registry`) because `internal/` must not import `app/`; the service
  therefore accepts catalog rows as DTOs (`PickerItemDto`, `PickerSpellDto`,
  plain SID strings) instead of reaching for them itself.
- **`BonusPickerDialog` and `BonusesPanel` are typed on `IGuiHandler`, not on a
  single facet**, because they legitimately need both the bonus and the picker
  facets; splitting them would have meant two constructor parameters for one
  collaborator.

Wiring: three services added to `EditorSet` and three handlers to `HandlerSet`
in `providerSets.go`, with `NewGuiHandler` growing to seven dependencies. Still
zero `wire.Bind` calls.

**Verification** — `go build ./...`; `go vet -tags="integration_test,gui" ./...`;
`go run ./cmd/testlayoutcheck .` → `test-layout check passed`;
`go test ./test/unit/... -count=1` clean; `go test -tags=integration_test
./test/integration/...` ok; `go test -tags "integration_test,gui"
./test/integration/gui/...` ok with **all snapshots unchanged**;
`golangci-lint-v2 run ./...` → `0 issues.`; coverage **72.7 %**
(Phase 0 floor 69.3 %, Phase 2 left it at 70.9 %).

**Incident note** — a repo-wide PowerShell CRLF→LF rewrite run at the end of this
phase corrupted every `.go` file in the working tree. The branch was reset by the
owner and the phase was recovered from VS Code Local History (102 files), then
re-verified end to end. Do not run bulk in-place rewrites across the repository;
use `gofmt -w` on an explicit file list instead.

---

## Phase 4: §0.2 — make the reset button honest
Status: Complete

Owner decisions 6 and 7 govern this phase. Re-read them before starting.

- [x] Relabel `"Reset to generated"` to wording that matches what it does —
      restore the state the dialog was opened with. `"Revert changes"` or
      `"Discard edits"`.
- [x] Delete the in-code marker comment at
      [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go#L213). Verify
      afterwards that the non-test Go tree still contains zero TODO/FIXME/HACK
      comments.
- [x] Track that the revert button was used during this dialog session.
- [x] On Apply **after a revert**, clear the persisted manual-edit snapshot
      (`ManualZones` / `ManualConnections`) outright. This needs a signal from the
      dialog to `ApplyEditedZones`, whose callback is currently
      `func([]entities.Zone, []entities.Connection)` — extend it, or add a
      sibling method on the driver. Prefer whichever keeps the GUI free of policy.
- [x] Put the *decision* ("was this apply a post-revert apply, and should the
      snapshot be cleared?") in a service or the driver, not in the dialog —
      consistent with Batch 14's split.
- [x] Test the accepted consequence explicitly: after revert → Apply, the live
      template still shows the reverted-to edits, and after the next regeneration
      they are gone. That is intended behaviour, and the test should say so.
- [x] Update `todo/review-opus5-08-04.md` §0.2's disposition row: it currently
      says "Owner-deferred… Agents must not action it."

### Verification Plan
- New unit tests cover both the revert-then-apply and apply-without-revert paths.
- A GUI snapshot confirms the new button label.
- Full suite green.

### Phase Summary

**What was done.** The button now reads **"Revert changes"**; the marker comment
is gone (a repo-wide grep confirms the non-test Go tree still has zero
TODO/FIXME/HACK-class comments); the dialog tracks `revertUsed`, set in
`resetToOriginal`; and a post-revert Apply clears the persisted manual snapshot.

**The seam.** `onApply` was widened from
`func([]entities.Zone, []entities.Connection)` to
`func(dtos.ZoneEditorApplyDto)` — a new DTO in
[zoneEditorApplyDto.go](../internal/dtos/zoneEditorApplyDto.go) carrying
`Zones`, `Connections` and `RevertUsed`. The alternative the plan offered (a
sibling driver method) was rejected: it would force the dialog to *choose* which
method to call, which is exactly the policy decision Batch 14 pushed out of the
GUI. With a DTO the dialog only reports a fact.

**Where the decision lives.** In the driver —
[`State.ApplyEditedZones`](../app/gui/drivers/stateManualEdits.go) — not in a new
service. The policy is one branch (`if request.RevertUsed { ClearManualEdits() }`),
and a service for one branch is over-engineering; the plan explicitly allowed
either. Owner decision 7 was implemented literally: once revert is pressed in a
session, the following Apply clears the snapshot outright, even if the user then
edited further before applying.

**Safety check.** Clearing the snapshot cannot make a later regeneration wipe the
template: `regenerationDecisionService.DecideManualEditReapplication` early-returns
an empty decision when `!current.HasManualEdits()`, so the `ClearManualEdits`
branch is taken and `reapplyManualEdits(nil, nil)` is never called.

**Deviation from the Verification Plan.** "A GUI snapshot confirms the new button
label" was not achievable as written — per the Phase 0 Deviation there is no
per-dialog golden-image path. The label is instead asserted through Gio
**semantics**: one laid-out frame, then `input.Router.AppendSemantics(nil)`
filtered to `semantic.Button` nodes. Two tests: one asserts `"Revert changes"` is
present, one asserts `"Reset to generated"` is absent. The helper guards with
`require.NotEmpty` on the label list so neither can pass vacuously.

**Tests.** 13 `ApplyEditedZones` call sites migrated to the DTO. Added: 2 unit
tests in
[applyEditedZones_test.go](../test/unit/app/gui/drivers/stateManualEdits/applyEditedZones_test.go)
(snapshot dropped / live template keeps the reverted-to zones); 3 untagged
integration tests in
[zoneEditorRevert_integration_test.go](../test/integration/zoneEditorRevert_integration_test.go)
covering the accepted consequence end-to-end, including a comment stating the
disappearing act is intended; 4 GUI tests (2 label, 2 `RevertUsed` round-trip).
Three existing tests were renamed `…_TheEditorIsResetToGenerated_…` →
`…_TheEditorIsReverted_…`.

**Trap worth remembering.** The `ClickReset()` testexport was deliberately *not*
renamed to `ClickRevert()`: `cmd/testlayoutcheck` matches test-only export
accessors **by identifier name across the whole tree**, in both directions, so
colliding names silently mis-gate files.

**Verification results.** `go build ./...`, both `go vet` tag combinations,
`go run ./cmd/testlayoutcheck .` (`test-layout check passed`), unit, integration
(`ok 2.292s`) and GPU-gated GUI (`ok 2.849s`) suites all green.
`gofmt -l ./app ./internal ./test ./cmd` is empty (both new files needed an
explicit `gofmt -w` — they were written with CRLF).
`golangci-lint-v2 run ./...` reports `0 issues.` Coverage held at **72.7%**
(floor 69.3%); the phase's new code is GUI-layer and driver-layer, so it moves
the needle very little either way.

**Superseded — see Phase 4b.** The owner tested this build and the behaviour was
wrong in practice. Everything above about `RevertUsed` and the "Revert changes"
button is history, not current state.

---

## Phase 4b: Split the revert into Undo and Revert to Base
Status: Complete

Phase 4 shipped, the owner ran it, and found it half-working: pressing "Revert
changes" cleared the manual snapshot but nothing on screen changed, because no
regeneration was triggered — the edits only disappeared after a *separate*
Generate click. The owner asked for two distinct actions instead of one.

- [x] Rename "Revert changes" to **"Undo"** — one-shot restore of the current
  editing session, purely in-session.
- [x] Add **"Revert to Base"** — drop every stored manual edit, regenerate
  immediately, and re-seed the open editor with the fresh layout.
- [x] Remove `RevertUsed` and `ZoneEditorApplyDto` entirely.
- [x] Migrate every test and call site.

### Owner decisions (settled — do not re-litigate)
1. **Revert to Base rerolls.** It clears manual edits and regenerates. It does
   **not** restore the pristine layout the edits were originally made on — that
   layout is retained nowhere, and no new state is introduced to retain it.
2. **Undo is one-shot**, not a step-by-step undo stack. It resets the whole
   editing session at once (the Phase 4 behaviour, relabelled).
3. **Revert to Base keeps the dialog open**, showing the newly generated base.
4. **Undo is purely in-session.** Apply after an Undo re-commits the
   previously-applied edits unchanged; only Revert to Base clears the stored
   snapshot.

**On the record:** with reroll + stay-open, Revert to Base fires immediately and
**Cancel cannot take it back**. That is inherent to decision 3, and accepted.

### Verification Plan
- Unit tests prove `RevertZonesToBase` drops the snapshot, regenerates, returns
  the new zones, and reports failure without leaking stale zones.
- An untagged integration test reproduces the reported defect end-to-end: after
  `RevertZonesToBase` the live template must already carry no `ManualPosition`,
  with **no** separate `Generate()` call.
- GUI tests assert both button labels and both behaviours.
- Full suite, lint, format and coverage green.

### Phase Summary

**The seam.** `dtos.ZoneEditorApplyDto` was deleted and replaced by
[zoneEditorZonesDto.go](../internal/dtos/zoneEditorZonesDto.go), a plain
`{Zones, Connections}` pair used **bidirectionally** — out of the editor on
Apply, back into the open editor on Revert to Base. With `RevertUsed` gone the
two payloads are structurally identical, so two near-duplicate DTOs would have
been worse. The dialog gained a second callback,
`onRevertToBase func() (dtos.ZoneEditorZonesDto, bool)`, making
`NewZoneEditorDialog` an 8-parameter constructor.

**Driver side.** `State.handleGenerateTemplate` now returns `bool`. Without it a
failed reroll would leave `lastTemplate` holding the *edited* template with
`hasTemplateVariants()` still true, and `RevertZonesToBase` would have handed the
edited zones back labelled as "base". `RevertZonesToBase` calls
`ClearManualEdits()` **before** generating — the other order lets
`DecideManualEditReapplication` stamp the edits straight back onto the fresh
template.

**Dialog side.** `resetToOriginal` split into `undoSessionEdits()` (the old body
minus the `revertUsed` flag) and `revertToBase()`. A new `setEditingSet` helper
installs a zone/connection list as *both* the working copy and the Undo
baseline, so a revert to base also moves the point Undo returns to; it is used by
the constructor and by `revertToBase`. `funcorder` required it to sit after the
exported methods.

**Tests.** `applyEditedZones_test.go` rewritten (8 tests, the two `RevertUsed`
ones deleted); new
[revertZonesToBase_test.go](../test/unit/app/gui/drivers/stateManualEdits/revertZonesToBase_test.go)
(7 tests);
[zoneEditorRevert_integration_test.go](../test/integration/zoneEditorRevert_integration_test.go)
replaced by
[zoneEditorRevertToBase_integration_test.go](../test/integration/zoneEditorRevertToBase_integration_test.go)
(4 untagged tests, one of which is the reported defect); GUI file reworked —
`ClickReset` → `ClickUndo` at 4 sites, new `ClickRevertToBase`, label tests now
assert `"Undo"` and `"Revert to Base"`, plus 4 new behaviour tests driven by a
stub `onRevertToBase`.

**Name-collision check.** `ClickUndo`, `ClickRevertToBase`, `RevertZonesToBase`
and `ZoneEditorZonesDto` were all grepped tree-wide before use — see the Phase 4
trap note.

**Verification results.** `go build ./...`, both `go vet` tag combinations,
`go run ./cmd/testlayoutcheck .` (`test-layout check passed`), unit, integration
(`ok 3.727s`) and GPU-gated GUI (`ok 4.109s`) suites all green.
`gofmt -l ./app ./internal ./test ./cmd` empty (the two new files needed an
explicit `gofmt -w` — CRLF again). `golangci-lint-v2 run ./...` reports
`0 issues.` Coverage **72.5%** (floor 69.3%).

**Superseded by Phase 4c.** The owner tested this and rejected the immediacy:
the reroll must not touch the live template until Apply.


---

## Phase 4c: Defer the revert until Apply

Status: Complete

Owner report after testing Phase 4b, verbatim: *"when Revert to Base is clicked,
the template is regenerated even before clicking Apply, so not just the preview
shows the base generated template, but also when cancel is clicked, base is
still shown in preview panel"*. This overturns the consequence the owner had
accepted in Phase 4b decision 3 — keeping the dialog open was right, committing
on the spot was not.

- [x] `State.RevertZonesToBase` replaced by `State.PreviewBaseZones`, which
  generates through `handler.GenerateTemplate` directly and **commits nothing**
- [x] `dtos.ZoneEditorZonesDto` gained `RevertToBase bool` so the apply
  direction can report that the session reverted
- [x] `State.ApplyEditedZones` decides: an untouched base clears the manual
  snapshot, a base the user edited afterwards is stored as a normal snapshot
- [x] `handleGenerateTemplate`'s Phase 4b `bool` return reverted to void (dead
  code once the reroll stopped going through it)
- [x] Dialog tracks `revertedToBase` and reports it on Apply; hint text no
  longer claims immediacy
- [x] Unit, integration and GUI tests reworked for the deferred semantics

### Verification Plan

- `go build ./...` and `go vet -tags="integration_test,gui" ./...` — silent
- `go run ./cmd/testlayoutcheck .` — `test-layout check passed`
- `go test ./test/unit/... -count=1`,
  `go test -tags=integration_test ./test/integration/... -count=1`,
  `go test -tags "integration_test,gui" ./test/integration/gui/... -count=1` — all ok
- `golangci-lint-v2 run ./... --issues-exit-code=0` — `0 issues.`
- Coverage ≥ 69.3%

### Phase Summary

**Why the DTO carries a flag.** The driver has to tell two applies apart that
look identical from outside: applying an untouched fresh base (→
`ClearManualEdits`, so later regenerations stay clean) and applying edits made
*on top of* that base (→ `SetManualEdits`, so they persist and are saved).
Always storing a snapshot would pin the base layout forever, which is the very
§0.2 bug the phase set out to fix; never storing one would silently discard
post-revert edits. The dialog reports only the bare fact that a revert happened;
the *policy* — comparing the applied zones against `State.pendingBaseZones` with
`reflect.DeepEqual` — lives in the driver, keeping the GUI free of business
logic (the constraint that made Batch 14 reject a sibling-method design).

**Why `PreviewBaseZones` bypasses `handleGenerateTemplate`.** The latter commits
(`applyGeneratedTemplate` → `setLastTemplate` + `SnapshotCurrentState`) and
reapplies manual edits. Generation itself never consults `ManualZones` /
`ManualConnections` — reapplication happens afterwards — so calling the handler
directly yields a genuine base without touching anything the user can see.

**Why no template is stashed.** Applying the base runs the normal
`handleUpdateTemplate(baseZones, baseConnections)` path, which rebuilds the
layout from `lastTemplate`; zones carry their own content and the non-variant
fields are state-derived and identical, so a second stored template would be
redundant.

**Why a stale `pendingBaseZones` is harmless.** It is only read when
`request.RevertToBase` is true, which only a dialog that actually reverted can
set; it is cleared at the top of every `ApplyEditedZones` and overwritten by
every `PreviewBaseZones`. A failed reroll leaves `revertedToBase` false, so a
failed revert cannot make the next Apply drop the user's edits — covered by
`TestWhenRevertToBaseFailed_TheApplyReportsNoRevert`.

**Tests.** `revertZonesToBase_test.go` → `previewBaseZones_test.go` (7 tests,
semantics inverted: the live template and the stored snapshot must be
*untouched*); `applyEditedZones_test.go` gained 3 tests for the `RevertToBase`
branch; `zoneEditorRevertToBase_integration_test.go` rewritten (6 untagged
tests, two of which reproduce the owner's exact complaint — nothing changes on
preview, everything changes on apply, with no separate `Generate()`); the GUI
file gained 3 tests asserting the flag round-trip.

**Verification results.** `go build ./...`, both `go vet` tag combinations and
`go run ./cmd/testlayoutcheck .` clean; unit suite ok, integration `ok 3.593s`,
GPU-gated GUI `ok 4.161s`; `gofmt -l ./app ./internal ./test ./cmd` empty (the
new files needed an explicit `gofmt -w` — CRLF again); `golangci-lint-v2` reports
`0 issues.`; coverage **72.5%** (floor 69.3%).


---

## Phase 5: Close out
Status: Complete

- [x] Full suite: build, both `go vet` tag combinations, testlayoutcheck,
      `wire diff`, unit, integration, GPU-gated GUI, performance, coverage, lint.
      Also `gofmt -l .` and a `GOOS=linux` lint pass.
- [x] Coverage ≥ the Phase 0 baseline. This batch should *raise* it noticeably —
      ~430 LOC of currently-uncovered logic becomes unit-testable.
- [x] Mark §2.6 `✅ FIXED` **in place** in
      [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) and update §12.
      With §2.6 closed, **all 46 review findings are resolved** — say so.
- [x] Update `todo/test_observations.md`: it currently lists `zoneEditorDialog.go`
      and its canvas/snap/property files as untestable Gio territory. Much of that
      is no longer true.
- [x] Update repository memory.
- [x] Rewrite `.agent/session-carry-forward.md`.
- [x] Stop for owner review. Do not stage. Do not commit.

### Verification Plan
- Every command in AGENTS.md §7 Quick Reference passes.
- `golangci-lint-v2 run ./...` reports `0 issues.`

### Phase Summary

| Check | Result |
| --- | --- |
| `go build ./...` | clean |
| `go vet ./...` / `go vet -tags="integration_test,gui" ./...` | clean |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 — generated code current |
| `go test ./test/unit/... -count=1` | ok |
| `go test -tags=integration_test ./test/integration/... -count=1` | `ok 3.593s` |
| `go test -tags "integration_test,gui" ./test/integration/gui/... -count=1` | `ok 4.161s` |
| `go test "-bench=." -run=xxx ./test/performance/... -benchtime=20x` | PASS, 10 benchmarks |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `golangci-lint-v2 run ./...` | `0 issues.` |
| `GOOS=linux CGO_ENABLED=0` lint of `./internal/... ./cmd/...` | `0 issues.` |
| Coverage | **72.5%** (Phase 0 floor 69.3%) |

**Benchmarks** (Intel Core Ultra 7 165H, 20 iterations): preview layout
3.5 µs (default) to 4.36 ms (random large); generation 20.6 µs (tournament) to
83.8 µs (hub-and-spoke). No regression signal — the Phase 1 extraction moved
geometry behind an interface without a measurable cost.

**On the `GOOS=linux` pass.** A whole-tree `GOOS=linux` lint cannot typecheck
`app/gui` from Windows: Gio's X11/EGL backend needs cgo, so
`gioui.org/internal/gl` fails to import. That is a cross-compilation limitation,
not a finding. The portable half of the tree (`./internal/...`, `./cmd/...`) was
linted under `GOOS=linux CGO_ENABLED=0` and reports `0 issues.`; CI lints the
whole tree on Ubuntu.

**Documentation closed out.** [review-opus5-08-04.md](../todo/review-opus5-08-04.md)
§2.6 marked `✅ FIXED` in place with the evidence retained beneath it, §0.2
rewritten for the final Undo / Revert-to-Base design, and §12 item 13 records the
batch and states that **all 46 review findings are now resolved**.
[test_observations.md](../todo/test_observations.md) no longer claims the zone
editor is untestable Gio territory — the geometry is unit-tested in
`connection_editor`, the revert semantics in `drivers`, and the dialog itself is
driven by the GPU-gated GUI suite; only the property-panel editors and pointer
flows remain uncovered.

## Final Recap

Seven phases (0, 1, 2, 3, 4, 4b, 4c, 5) closed review §2.6 and §0.2 together.

- **Phase 0** pinned the zone editor's geometry numerically and by snapshot, so
  every later phase could prove it moved no pixels. Those pins passed unchanged
  through all of it.
- **Phase 1** extracted the geometry into
  `internal/services/connection_editor` behind the `IZoneEditorHandler` the
  dialog already held — ~160 LOC that was at 0 % coverage is now ≥ 92.9 %.
  Pointer identity survives by carrying a `ConnectionIndex` rather than a
  pointer; `float32` stayed at the GUI boundary.
- **Phase 2** replaced the `zoneEditorCanvasState` blob with a geometry struct
  and an interaction struct, the latter exposing 13 intention-revealing methods
  over its 12 fields. 67 → 63 fields.
- **Phase 3** emptied the four sibling dialogs into `internal/services/bonuses`,
  `internal/services/pickers` and `internal/services/zone_content`, adding three
  handler facets. Coverage reached 72.7 %.
- **Phases 4 → 4b → 4c** fixed §0.2 across three rounds of owner testing: a
  relabelled button, then a real reroll, then — because the reroll committed
  before Apply and survived Cancel — a preview that only commits on Apply.
  Final shape: **Undo** (session-scoped, one-shot) and **Revert to Base**
  (`State.PreviewBaseZones` shows a fresh manual-edit-free layout, commits
  nothing; `State.ApplyEditedZones` clears the manual snapshot only if the
  applied zones still match the previewed base). All policy in the driver, none
  in the dialog.
- **Phase 5** verified everything and closed the review.

Final field count: **60, of which 26 are Gio widget handles that cannot move.**
Coverage **69.3 % → 72.5 %**. `golangci-lint-v2`: `0 issues.`

The two things a future agent should not undo: `dtos.ZoneEditorApplyDto` /
`RevertUsed` / `State.RevertZonesToBase` are deleted deliberately, and the
`handleGenerateTemplate` `bool` return added in Phase 4b was removed again in
Phase 4c because the preview path no longer goes through it.

## Deployment Plan

Nothing to deploy — this is a refactor with no schema, no output-path and no
persisted-state change. `.gen.json` and `.rmg.json` formats are untouched, and
the game's templates directory is still resolved per launch by
`FindOldenEraTemplatesDir`.

1. Owner reviews the working tree and runs `go run .` to confirm the zone
   editor's Undo / Revert-to-Base behaviour by hand.
2. Owner stages and commits on `AD/refactoring-07-21` (the agent staged and
   committed nothing, per AGENTS.md §2.5).
3. Push and let CI run the Ubuntu lint and the untagged suites; the GPU-gated
   `gui` tests do not run there by design.
4. No migration, no user action, no rollback step beyond reverting the commit.

---

## Known traps (carried from Batches 13–14 — read before starting)

- **Snapshot tests are the safety net for this entire batch.** A changed snapshot
  during Phases 1–3 means the refactor changed rendering, which it must not. Do
  not run the *"Go: Update UI Integration tests snapshots"* task to make a failure
  go away — investigate it.
- **GPU-gated tests need `-tags 'integration_test,gui'`** and are excluded from
  catch-all runs by design. Never make `gui` a default tag.
- **PowerShell mangles a bare `-bench=.`** — `go test -bench=. ./test/performance`
  prints `[no test files]` and runs nothing. Quote it: `go test "-bench=."`.
- **`wire gen` writes its success banner to STDERR**, which PowerShell surfaces as
  an error. Judge by `wire diff` (exit 0).
- **`golangci-lint-v2 run ./...` skips build-tag-gated files** on Windows. CI lints
  on ubuntu with no tags. `gofmt -l .` ignores tags entirely and is the widest
  tripwire. Two files (`dialogHost_testexports.go`, `wire.go`) are permanent
  CRLF-only `gofmt -l` noise — do not "fix" them.
- **After `--fix` over brand-new files**, re-run `testlayoutcheck` and check line 2
  for a duplicated `package` clause.
- **Gio dialog tests**: use the public `widget.Clickable.Click()` plus one laid-out
  frame. `Clickable.update` drains `requestClicks` before consulting pointer input,
  so no coordinates are needed. Reserve `AppRunner.ClickAt` for genuinely geometric
  behaviour — which, in this batch, the canvas drag tests genuinely are.
- **`app/` may not import `internal/services`** — depguard `no-services-from-app`
  also denies repositories, mappers and validators. Reach services through a
  handler interface.

## Final Recap
_(write when all phases complete)_

## Deployment Plan
_(write when all phases complete)_
