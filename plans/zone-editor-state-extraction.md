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
Status: Not started

- [ ] Group the 15 selection/drag/interaction fields into one named struct that
      stays in `app/gui/dialogs/` — it is view state, by owner decision 2.
- [ ] Give that struct the small number of methods that operate on it alone
      (select/clear/begin-drag/end-drag), so the reset path has one obvious place
      to call instead of assigning eleven fields inline as
      [resetToOriginal](../app/gui/dialogs/zoneEditorDialog.go#L365-L383) does now.
- [ ] Do **not** move it to `internal/services`.
- [ ] Re-home the geometry fields (bucket b) that Phase 1 made redundant; delete
      any that the service now owns.

### Verification Plan
- Snapshots unchanged; unit and integration suites green.
- `ZoneEditorDialog`'s direct field count is materially lower than 67 — record the
  new number in the Phase Summary.

### Phase Summary
_(write when phase completes)_

---

## Phase 3: The four extra dialogs
Status: Not started

Extract **everything identified** in the candidates table (owner decision 4). Work
one file at a time and keep each file's snapshot green before moving on.

- [ ] `bonusPickerDialog.go` — validation, `BonusEntry` mapping, duplicate/spell-ID
      extraction, `spellCountLabel`, `isNumeric`.
- [ ] `pickerDialog.go` — source→entry mapping, filtering/grouping row model,
      `selectedIDs`.
- [ ] `ruleDialog.go` — `upsertFromEditor`, `buildRuleFromEditor`.
- [ ] `zoneContent.go` — `rowDisplayName`, `defaultContentRules`, `ruleMarkers`,
      the constructor's alphabetical sort, `Add`'s clamping.
- [ ] Decide placement per AGENTS.md §4.4 as you go: content-rule logic belongs
      near `internal/services/content_rules/`; picker/bonus mapping may warrant its
      own package. Note the choice and its reason in the Phase Summary.
- [ ] Unit tests for every extracted unit.

### Verification Plan
- Snapshots for all four dialogs unchanged.
- Each extracted unit has its own `test/unit/` folder per §4.6.
- `golangci-lint-v2 run ./...` still reports `0 issues.`

### Phase Summary
_(write when phase completes)_

---

## Phase 4: §0.2 — make the reset button honest
Status: Not started

Owner decisions 6 and 7 govern this phase. Re-read them before starting.

- [ ] Relabel `"Reset to generated"` to wording that matches what it does —
      restore the state the dialog was opened with. `"Revert changes"` or
      `"Discard edits"`.
- [ ] Delete the in-code marker comment at
      [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go#L213). Verify
      afterwards that the non-test Go tree still contains zero TODO/FIXME/HACK
      comments.
- [ ] Track that the revert button was used during this dialog session.
- [ ] On Apply **after a revert**, clear the persisted manual-edit snapshot
      (`ManualZones` / `ManualConnections`) outright. This needs a signal from the
      dialog to `ApplyEditedZones`, whose callback is currently
      `func([]entities.Zone, []entities.Connection)` — extend it, or add a
      sibling method on the driver. Prefer whichever keeps the GUI free of policy.
- [ ] Put the *decision* ("was this apply a post-revert apply, and should the
      snapshot be cleared?") in a service or the driver, not in the dialog —
      consistent with Batch 14's split.
- [ ] Test the accepted consequence explicitly: after revert → Apply, the live
      template still shows the reverted-to edits, and after the next regeneration
      they are gone. That is intended behaviour, and the test should say so.
- [ ] Update `todo/review-opus5-08-04.md` §0.2's disposition row: it currently
      says "Owner-deferred… Agents must not action it."

### Verification Plan
- New unit tests cover both the revert-then-apply and apply-without-revert paths.
- A GUI snapshot confirms the new button label.
- Full suite green.

### Phase Summary
_(write when phase completes)_

---

## Phase 5: Close out
Status: Not started

- [ ] Full suite: build, both `go vet` tag combinations, testlayoutcheck,
      `wire diff`, unit, integration, GPU-gated GUI, performance, coverage, lint.
      Also `gofmt -l .` and a `GOOS=linux` lint pass.
- [ ] Coverage ≥ the Phase 0 baseline. This batch should *raise* it noticeably —
      ~430 LOC of currently-uncovered logic becomes unit-testable.
- [ ] Mark §2.6 `✅ FIXED` **in place** in
      [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) and update §12.
      With §2.6 closed, **all 46 review findings are resolved** — say so.
- [ ] Update `todo/test_observations.md`: it currently lists `zoneEditorDialog.go`
      and its canvas/snap/property files as untestable Gio territory. Much of that
      is no longer true.
- [ ] Update repository memory.
- [ ] Rewrite `.agent/session-carry-forward.md`.
- [ ] Stop for owner review. Do not stage. Do not commit.

### Verification Plan
- Every command in AGENTS.md §7 Quick Reference passes.
- `golangci-lint-v2 run ./...` reports `0 issues.`

### Phase Summary
_(write when phase completes)_

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
