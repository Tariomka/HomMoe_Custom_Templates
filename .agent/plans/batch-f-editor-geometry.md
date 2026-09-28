# Batch F: editor geometry

Close review findings §1.13 (deterministic obstacle routing), §1.14 (stale edge identity
under batched pointer events) and §1.15 (portal classification), together with the
owner-approved optional §2.1 (one shared curve builder so the editor, the Preview tab and
the PNG draw identical curves) and the §8 "dropdown normalization" note. Source of the
findings: [review](../backlog/review-gpt-6-astra-09-07.md) §1.13, §1.14, §1.15, §2.1, §9 row F.

## For Future Agents
As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../../AGENTS.md) first. Never stage/commit/push/stash/switch branches.
Never touch `data/`, `internal/entities/template_entity/`, `internal/registry/`, the
output directory logic, or generated Wire by hand. No `*_testexports.go` in unit tests,
no global `integration_test`/`gui`/`wireinject` tags, no fake unit seams.

## Owner decisions (2026-09-28, binding)

| # | Question | Decision |
| --- | --- | --- |
| D1 | Curve agreement editor vs Preview/PNG | **Exact agreement required.** One shared builder produces every control point. |
| D2 | Parallel spacing | **21px** (Preview's value) everywhere; the editor's 18px goes away. |
| D3 | Obstacle bending | **On everywhere** (moves from editor-only into the shared builder; Preview tab and PNG gain it). Tunables unchanged: clearance `zoneRadius + 8px`, padding `6px`, chord margin `0.08`. |
| D4 | Single obstructed curve | **Shorter detour**; the positive-normal side wins **exact ties only**. |
| D5 | N ≥ 2 parallel curves, obstructed | Keep slot order (most negative → most positive). **Even N:** `N/2` per side. **Odd N:** the shorter-detour side (positive on a tie) gets `ceil(N/2)`, the other `floor(N/2)`; N=1 is the special case of the same rule. The innermost curve on each side gets that side's obstacle bulge (P3); each further curve steps out by 21px. Unobstructed pairs keep today's symmetric fan. |
| D6 | §2.1 | **In scope**, delivering D1–D5. |
| D7 | Goldens | Zone-editor snapshots, Preview-tab snapshots, PNG/preview-layout and editor geometry expectations **may be regenerated**; every changed golden is listed with its reason in the phase summary. |
| D8 | Portal classification | One model-level check: type equals `Portal` **case-insensitively**, **or** any `PortalPlacementRulesFrom`/`To` present. The pinned lowercase-`"portal"` preview test flips to portal-shaped. |
| D9 | Road tri-state | **Unchanged.** Explicit `false` roadless, explicit `true` roaded, `nil` roaded only for explicit Portals (`IsExplicitPortal`). Effective classification never widens the `nil` rule. |
| D10 | Editor portal width | Portals drawn with `DefaultConnectionLineSmall` like the Preview tab; the selected curve stays `DefaultConnectionLineLarge`. |
| D11 | Type change vs placement rules | A **real** dropdown change to any non-Portal type clears both placement-rule lists; choosing Portal keeps them. Loading or merely selecting never clears. Road policy restamps as today. |
| D12 | Dropdown display | Shows the **effective** type (a stored `Direct` with rules reads `Portal`). Items are `Direct`, `Portal`, plus the connection's own stored type when it is neither (and not a case variant of Portal). An empty stored type is listed as `(none)` and still saves as empty. |
| D13 | Dropdown writeback | The type is written **only** when the user changes the selection (`WasUpdated`). The silent per-frame overwrite of unlisted types (Proximity, GladiatorArena, Default, empty, lowercase portal) is removed. |

Tie direction semantics: `normal = (canonicalB − canonicalA).RotateClockwise() / distance`,
where canonical A/B are the positions of the lexicographically smaller/larger zone name.
"Positive side" = along `+normal` = left of travel from the alphabetically-first zone to
the second (up on screen for a left-to-right chord). Reversing `From`/`To` never changes it.

## Design (proposed; confirm at plan approval)

**P1 Shared curve builder location.** A model with attached logic in
`internal/models/preview/`. It is cycle-free: `preview` imports `helpers/data`, `neutral_zone`
and now `template_model`, and `template_model` does not import `preview`.
- `connectionCurve.go` — pure data `ConnectionCurve{ConnectionIndex int; Start, End, Control data.Vec2[float64]}`.
- `connectionCurveLayout.go` — `ConnectionCurveLayout{Positions map[string]data.Vec2[float64]; ZoneRadius float64}` with `Build(connections []template_model.Connection) []ConnectionCurve`; private grouping/obstacle helpers and the curve tunables live in the same file.
- Both `PreviewLayoutService.buildPreviewConnections` and `ZoneEditorGeometryService.BuildGeometry`
  call `Build` with the same positions and `layout.ZoneRadius`, so identical input yields
  identical control points by construction. No constructor changes → no Wire regeneration.

**P2 Shared emission order and orientation.** Curves are emitted in the editor's current
grouped first-seen pair order (its documented determinism contract), and `Start`/`End`
use the canonical endpoint order the Preview uses today (keeps PNG dash phase and arena
marker midpoint for reversed connections unchanged; the editor draws solid strokes so its
orientation is invisible; editor hit identity rides on `ConnectionIndex`, and the
test-exported `EdgeGeometry.From/To` still come from the connection). The Preview's
`Connections` order therefore changes from raw connection order to grouped order, which
can change which of two overlapping PNG strokes or an arena marker is painted last; this
is pinned by explicit order/orientation assertions, not only by regenerated goldens.
Visibility is uniform per pair (every member shares both endpoint names), so a pair is
either fully laid out or fully skipped and slot counts never include hidden members.

**P3 Obstacle math.** For each positioned zone centre whose chord projection ratio lies in
`(0.08, 0.92)` and whose perpendicular distance is `< clearance`, with signed offset
`s = offset · normal`: `positiveNeed = max(s + clearance + padding)`, `negativeNeed =
min(s − clearance − padding)` over all such obstacles. Shorter detour = smaller of
`positiveNeed` and `−negativeNeed`; ties → positive. A chord with squared length `< 1`
has no obstacle bulge (today's rule), and the normal's `distance < 1` clamp to 1 is kept.
The bulge is applied as the curve's **midpoint** offset (control = midpoint +
`2·bulge·normal`), exactly as the editor does today; it is a midpoint approximation, so an
obstacle projecting away from the chord centre (toward the 0.08 margin) can still be
grazed. Curve-aware clearance is **not** in scope because it would change today's
single-obstacle results, which D3 keeps. Consequences vs today, all intended by D4:
- one obstacle: identical result, except an obstacle **exactly on the chord** now bends
  positive (today it bends negative);
- multiple obstacles: the chosen side's midpoint offset clears all of them rather than
  only the deepest;
- the review's repro (endpoints `(100,350)`/`(600,350)`, obstacles `(350,340)`/`(350,360)`,
  radius 21) gives `s = ±10`, `positiveNeed = 45`, `negativeNeed = −45`, a tie → positive,
  control point **`(350, 260)`** for any map construction order.
This does **not** solve routing past zones that sit outside the chord clearance, and does
not claim to fix every overlapping-obstacle layout.

**P4 Stale edge identity (§1.14).** Hit tests resolve against geometry built from the
**current** connection list: `hitTestNode`, `hitTestEdge` and `ensureManualPositions`
first call `ensureGeometry`, which rebuilds when a mutation marked the geometry dirty
**or** the geometry was built for a different canvas side than `this.side` (tracked as the
side the geometry was built for), so `ConnectionIndex` always indexes the list it was
built from. Zone-drag snapping needs no extra call: a drag always begins with a Press whose
hit test refreshes the geometry, and afterwards only the dragged zone (which snapping
excludes) moves. Every `working` mutator already raises `geometryDirty`
(`setEditingSet`, `addConnection`, `deleteConnection`, `undoSessionEdits`, `deleteZone`,
`applyQualityMutation`); Phase 4 re-verifies that. Side-panel connection writeback edits a
record in place and never changes list identity or curve geometry. Semantic: a press that follows a
structural edit in the same frame hits the geometry that edit produced. Chosen over an
identity snapshot because `deleteZone`, Undo, Revert and quality remaps replace every
connection pointer, which a snapshot would have to drop wholesale.

**P5 Classification seam.** `template_model.Connection.IsEffectivePortal()` is the single
classifier. The preview calls it for shape; the GUI calls it through the model it already
holds for colour (`NewEditorConnectionLineStyle`) and width, exactly as it already calls
`HasRoad()`/`IsExplicitPortal()`. No new handler method.

**P6 Dropdown shaping** stays in the dialog (pure view-model shaping over model methods);
the rule-clearing business rule lives in `ZoneEditorService.ChangeConnectionType`.
`DropdownSelector.WasUpdated` is false when the current item is reselected, so reselecting
the displayed `Portal` of a stored-Direct-with-rules connection changes nothing (it does
**not** convert it to an explicit Portal). `(none)` maps to `""`. A stored lowercase
`direct` is unlisted and appears as its own item; a stored Proximity with rules reads
`Portal` and lists `Proximity` as an extra item.

## Plan review (2026-09-28)

Independent GPT-6 Sol review: **REJECT** on the first draft. All findings were folded in:
the conflict in the D5 odd-N rule was resolved, the midpoint-approximation limit is now
stated and pinned by a test, the order and orientation assertions are explicit, P4 now
also rebuilds on a canvas-side change and before `ensureManualPositions`, the batch
harness is one locked helper with a delivery proof, and the dropdown reselection semantics
are spelled out. One finding was declined with a reason: the claim that hidden pair
members shift slots cannot happen, because every member of a pair shares both endpoint
names (P2). The reviewer computed the repro control point independently as `(350, 260)`,
matching P3. Owner plan approval: **approved 2026-09-28**, with the instruction to run
Phases 0 to 5 without stopping, because the owner will review everything at once.

## Phase 0: Baseline
Status: Complete

- [x] Record `git status --short`, branch and HEAD (expected: `AD/editor_geometry`, clean, `b9fb53e` or the owner's newer commit).
- [x] `go build ./...`
- [x] `go test -p=2 -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...` then `go tool cover '-func=coverage.txt'`; record total and the per-file lines for `zoneEditorGeometryService.go`, `previewLayoutService.go`, `connection.go` (template_variant_model), `zoneEditorService.go`, `connectionLineStyle.go`.
- [x] `go test ./test/...`
- [x] `go test -tags='integration_test,gui' ./test/integration/...` (record timings)
- [x] `go run ./cmd/testlayoutcheck .`
- [x] `golangci-lint-v2 run ./... --issues-exit-code=1`

### Verification Plan
- All commands pass; baseline numbers written into the Phase Summary. Any pre-existing failure stops the batch and is reported to the owner.

### Phase Summary
Environment: Windows/amd64, Go 1.27.0, empty `GOFLAGS`. Branch `AD/editor_geometry`, HEAD
`b9fb53e`, clean tree apart from this untracked plan. Everything passed: build, the `-p=2`
unit coverage run, `go test ./test/...`, the tagged run (root integration 3.498s, GUI
26.655s), testlayoutcheck, and lint with 0 issues.
**Total coverage is 74.9%.** Per function: `obstacleBulge` 95.7%, `snapAxis` 94.7%,
`buildPreviewZones` 90.0%. Every other function in `zoneEditorGeometryService.go`,
`previewLayoutService.go`, the template `connection.go`, `connectionLineStyle.go` and
`ZoneEditorService.ChangeConnectionType` is at 100%.

## Phase 1: Shared deterministic curve builder (§1.13, §2.1)
Status: Complete

- [x] Add `internal/models/preview/connectionCurve.go` and `connectionCurveLayout.go` per P1–P3 (21px spacing, obstacle bending, D4/D5 split, grouped first-seen order, canonical orientation, skip pairs with a missing endpoint, `distance < 1` clamps to 1 as today).
- [x] Unit tests `test/unit/internal/models/preview/connectionCurveLayout/build_test.go` (`package connectionCurveLayout_test`, triple-A, `t.Parallel`, one assertion each):
  - clear chord → control on the chord midpoint;
  - two parallel unobstructed → symmetric ±21 spread; three parallel → −21/0/+21;
  - reversed pair (`B→A`) groups with `A→B` and uses the same normal;
  - grouped first-seen emission order; `ConnectionIndex` carried through;
  - missing endpoint skipped; stacked endpoints collapse onto the point;
  - one obstacle on the negative side → positive bend (exact value); on the positive side → negative bend;
  - obstacle exactly on the chord → positive (tie);
  - review repro, symmetric opposing obstacles → control point `(350, 260)`; identical result across permuted map construction and 200 repeated builds;
  - two obstacles on opposite sides with unequal depth → shorter side whose midpoint offset clears both;
  - two parallel obstructed → one curve per side at the exact bulges;
  - three parallel obstructed, positive side shorter → two positive, one negative; negative side shorter → two negative, one positive; outer curves step 21px;
  - off-centre obstacle (projection ratio ≈ 0.2) → bulge value equals the midpoint formula (pins the documented approximation);
  - obstacle inside the chord margin ignored; obstacle beyond clearance ignored; chord shorter than 1px gets no obstacle bulge.
- [x] `PreviewLayoutService.buildPreviewConnections` builds `preview.Connection` from `ConnectionCurveLayout{Positions, ZoneRadius}.Build`, keeping `Type`/`HasRoad`/`ExplicitPortal` mapping as is.
- [x] `ZoneEditorGeometryService.BuildGeometry` builds edges from the same `Build` (label point `0.25·S + 0.5·C + 0.25·E` unchanged); delete `buildEdges` curve math, `groupConnectionsByPair`, `obstacleBulge`, the moved tunables and the now-unused `connectionPairKey` type file if nothing else uses it. Hit-test and snap tunables stay.
- [x] Update `buildGeometry_test.go` expectations (21px spread, bending values) and keep tests for verbatim positions/radius, preview-layout inputs, label midpoint, missing endpoint, stacked endpoints; add one asserting editor control points equal the shared builder's for an obstructed fixture.
- [x] Update `buildPreviewLayout_test.go`: spacing expectations, a new manual-position fixture proving the Preview now bends around an obstacle, any order-dependent assertion affected by P2, plus explicit tests that Preview connections come out in grouped first-seen order and that a reversed connection keeps canonical `Start`/`End`.
- [x] Regenerate affected Preview-tab and zone-editor GUI snapshots and PNG expectations; update `zoneEditorGeometry_integration_test.go` control points (e.g. the ring spread). List every changed golden + reason.

### Verification Plan
- `go build ./...`; `go test ./test/unit/... -count=1` passes.
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` passes after snapshot regeneration; `git status --short` lists only the intended golden files.
- A dedicated test builds the review repro 200 times with freshly constructed maps and asserts a single control point.
- Coverage for `connectionCurveLayout.go` 100% of statements; editor/preview service files not below baseline.

### Phase Summary
**What changed.** A new model, `preview.ConnectionCurveLayout.Build` (plus the `preview.ConnectionCurve`
value type), is now the only source of connection control points. It implements D2–D5 and P2–P3.
Obstacle `max`/`min` does not depend on order, so a Go map's iteration order can no longer change
the result. `PreviewLayoutService.buildPreviewConnections` maps the curves onto `preview.Connection`.
`ZoneEditorGeometryService.buildEdges` maps the same curves and adds the label point.
The old editor `groupConnectionsByPair`, `obstacleBulge`, curve tunables and
`connectionPairKey.go` were deleted. The pair key is now a `[2]string`, so no private struct was needed.

**Tests.** `build_test.go` has 20 tests covering every bullet above, including the review repro.
That test builds 200 times with rotated insertion order and gets a single control point, `(350,260)`.
Editor unit expectations: the parallel spread went from 368/332 to 371/329, and the label from
y 359 to 360.5. A new parity test shows editor control points equal the builder's.
Four Preview tests were added: bending, parity, grouped order, and canonical `Start`.

**Changed goldens and GUI expectations**, all caused by D2/D4/D5:
- `zoneEditorGeometry_integration_test.go`: the ring spread moved from `272/308` to `269/311`
  (18px to 21px), and the ring label from `x 281` to `279.5` (same cause).
- In the same file, the Hub-layout `Pseudo-A-B` control moved from `x 181.03` to `199.03`.
  The hub sits exactly on the chord, so the two Spawn-A/B edges now split one per side instead
  of shifting together by ±9px.
- One snapshot:
  `zoneEditorActions_integration_test/TestWhenADragStartsOnAZoneInAddConnectionMode_AConnectionIsCreated_5.golden`.
  The new Spawn-A to Spawn-B edge passes exactly over the Hub, so the tie now bends it to the
  positive side (right on screen) instead of the negative side (left). I viewed both images to
  confirm this.
- **No** Preview-tab snapshot, PNG unit expectation or root integration expectation changed: no
  fixture has a zone within clearance of a chord in those views.

**Process note.** `-update` rewrites *every* golden, 280 files, most of them unrelated to curves.
I restored the snapshot directory from HEAD with `git restore --worktree`, which leaves the index
untouched. I then accepted only the one `.failure` a plain run produced. Future agents should do the
same: never keep a blanket `-update`.

**Verification.** `go build ./...` passes. All unit tests pass. Tagged GUI suite: PASS (23.6s).
`testlayoutcheck` passes. Coverage is measured in Phase 5.

## Phase 2: Effective portal classification (§1.15)
Status: Complete

- [x] Add `IsEffectivePortal()` to `internal/models/template_model/template_variant_model/connection.go` (case-insensitive Portal or any placement rule). `IsExplicitPortal`/`HasRoad` untouched (D9).
- [x] Tests `test/unit/.../template_variant_model/connection/isEffectivePortal_test.go`: `Portal`, `portal`, `Direct`, `Direct`+From rules, `Direct`+To rules, empty type.
- [x] `getPreviewConnectionType` uses `IsEffectivePortal()`; GladiatorArena/Proximity matching unchanged. Flip `TestWhenAPortalTypeIsLowerCased_...` to portal-shaped (rename per convention) and keep the placement-rule roadless test.
- [x] `NewEditorConnectionLineStyle`: `DrawsAsPortal = IsEffectivePortal()`, `ExplicitPortal = IsExplicitPortal()`; update `newEditorConnectionLineStyle_test.go` (rule-only portal draws as portal; rule-only portal with `nil` road stays roadless `NoRoadLine`).
- [x] `drawEdges`: portal width `DefaultConnectionLineSmall`, selected `DefaultConnectionLineLarge`, otherwise `DefaultConnectionLine`.
- [x] GUI tests (`roadStyleVisuals_integration_test.go` helpers `edgeStyleCensus`/`edgeStrokeWidth`): explicit portal, placement-rule-only roaded portal and ordinary direct edge show the same colours in editor and Preview; editor portal stroke thinner than direct; selected portal stroke larger.

### Verification Plan
- Unit + tagged GUI suites pass; `connection.go` and `connectionLineStyle.go` coverage 100% for changed functions.

### Phase Summary
**What changed.**
- `Connection.IsEffectivePortal()` returns true for a `Portal` type in any letter case, or when
  there are any placement rules. It is now the single portal check.
- `getPreviewConnectionType` uses it, so the preview shape now ignores case (D8).
- `NewEditorConnectionLineStyle.DrawsAsPortal` uses it too. `ExplicitPortal` and `HasRoad` are
  unchanged, which keeps the D9 road tri-state intact.
- In `drawEdges`, a selected edge is drawn at Large, an effective portal at Small, and
  anything else at the default width.

**Tests.**
- Unit: `isEffectivePortal_test.go` has 4 tests, one of them a 5-case table. The lowercase
  preview test now expects the portal shape and was renamed
  `..._TheProjectedEdgeIsAPortalShapedExplicitPortal`. Two new editor style cases cover a
  rule-only portal with a road (drawn as a portal) and without one (stays roadless).
- GUI, in `roadStyleVisuals_integration_test.go`: a Direct connection carrying rules is painted
  as a roaded portal in the editor, and the Preview paints it the same way. With roads off it
  keeps the roadless colour.
- The connection is seeded with the new `AppRunner.ApplyManualConnectionEdit`. It goes through
  the production `drivers.State.ApplyEditedZones`, which is also Apply's path, and simulates a
  connection loaded from a file.

**Deviation.** There is no GUI assertion for the thinner portal. `gtx.Dp` rounds both 1.5dp and
2dp to 2px at the harness's 1x scale, so width cannot tell the two apart. This is recorded in
[test_observations.md](../backlog/test_observations.md). The existing width tests
(normal 2, selected 4) still pass.

**Verification.** Build, unit and tagged suites PASS (GUI 24.5s). `testlayoutcheck` passes.

## Phase 3: Type dropdown and placement-rule clearing (§1.15, §8 note)
Status: Complete

- [x] `ZoneEditorService.ChangeConnectionType`: after setting the type, clear `PortalPlacementRulesFrom`/`To` when the new type is not Portal (case-insensitive registry compare); then road stamp as today.
- [x] Extend `test/unit/internal/services/connection_editor/zoneEditorService/changeConnectionType_test.go`: Direct clears From; Direct clears To; Proximity clears; Portal keeps both; road still stamped.
- [x] Dialog `syncPropsFromConnection`: build per-connection items (`Direct`, `Portal`, + stored unlisted type, `(none)` for empty) and select the effective type (D12); keep the item→stored-value list on dialog state.
- [x] Dialog `writebackProps`: call `ChangeZoneEditorConnectionType` only on `WasUpdated`, with the stored value of the chosen item; remove the non-`WasUpdated` assignment (D13).
- [x] GUI tests in `zoneEditorProperties_integration_test.go` (dialog built with explicit connections, real clicks): Proximity/GladiatorArena/empty/lowercase-portal survive selection + idle frames + Apply unchanged; stored-Direct-with-rules reads `Portal` and survives selection + idle frames + Apply with its rules and `Direct` type; choosing Direct on it applies a Direct connection with no rules; then choosing Portal again applies an explicit Portal with no rules; choosing Portal on a Proximity-with-rules keeps the rules; `(none)` shown for empty type.

### Verification Plan
- Unit + tagged GUI suites pass; `ChangeConnectionType` 100% covered.

### Phase Summary
**What changed.**
- `ZoneEditorService.ChangeConnectionType` now clears both placement-rule lists whenever the
  new type is not an explicit Portal.
- The dialog keeps `typeValues`, the stored value behind each dropdown row.
  - `connectionTypeOptions` lists Direct and Portal, plus the stored type when it is neither
    `Direct` nor a Portal in any letter case.
  - `connectionTypeLabel` shows an empty type as `(none)`.
  - The dropdown selects the effective type.
  - `writebackProps` writes the type only on `WasUpdated`. The per-frame fallback overwrite
    is gone.

**Tests.**
- Unit: 4 new `changeConnectionType_test.go` cases (Direct clears From, Direct clears To,
  Proximity stops being an effective portal, Portal keeps the rules). The existing road-stamp
  cases are unchanged.
- GUI: new file [zoneEditorConnectionType_integration_test.go](../../test/integration/gui/zoneEditorConnectionType_integration_test.go)
  with 9 tests, one of them a 4-case table. Every test starts from a real selection click and
  ends with Apply.
  - Unlisted Proximity, GladiatorArena, empty and lowercase portal types survive selection,
    idle frames and Apply.
  - The dropdown lists the stored Proximity, and shows `(none)` for an empty type.
  - A stored Direct with rules reads Portal and survives selection plus Apply unchanged.
  - Picking Direct clears the rules. Picking Portal afterwards stays rule-free.
  - For a Proximity with rules shown as Portal, picking Proximity clears the rules.
  - A Portal with rules keeps them.
- Two read-only testexports were added, `ConnectionTypeLabels` and
  `SelectedConnectionTypeLabel`, both declared on `editor.IZoneEditorDialog`.

**Plan deviation.** The plan's "choosing Portal on a Proximity-with-rules keeps the rules"
became "a Portal with rules keeps them on Apply". Under D12 that connection already reads
Portal, so reselecting Portal is not a change (P6). Rule retention for a real change to Portal
is unit-tested instead (`TestWhenTypeBecomesPortal_ThePlacementRulesAreKept`).

**Verification.** Build, unit and tagged suites PASS. `testlayoutcheck` passes.

## Phase 4: Batched pointer identity (§1.14)
Status: Complete

- [x] Harness (`test/test_helpers/integration_common`, not `*_testexports`): one locked `AppRunner` method that runs a leading frame, queues several press/release gestures (mouse secondary or touch primary, per gesture), then runs a single trailing frame; `ZoneEditorHandler` wrappers (right-click several edges / click several canvas points in one frame; toolbar button + canvas press in one frame).
- [x] Write the GUI tests first in `zoneEditorPointer_integration_test.go` and confirm at least the deletion-shift case **fails on current code** (record the failure). If none reproduces, stop and report to the owner before changing production code.
  - same edge right-clicked twice in one frame → exactly that one connection deleted;
  - two different edges right-clicked in one frame → exactly those two deleted;
  - right-click edge X then left-click edge Y in one frame → Y selected (by name);
  - toolbar Delete of a selected connection + queued canvas press on another edge in one frame → the pressed edge is the one selected;
  - toolbar Undo after a deletion + queued canvas press on an edge in one frame → the pressed edge (by name) is selected.
  Delivery proof: the "two different edges → exactly those two deleted" case shows both presses reached the canvas in one frame; a same-edge result is only trusted once that case passes.
- [x] Implement P4: `ensureGeometry` (dirty or side mismatch) before hit tests and in `ensureManualPositions`; update the `layoutCanvas` comment; keep the `edgeConnection` bounds guard.
- [x] Re-verify every `working` mutator raises `geometryDirty`.

### Verification Plan
- New tests fail before and pass after; full tagged GUI suite passes; `go run ./cmd/testlayoutcheck .` passes.

### Phase Summary
**Harness.**
- `integration_common.PointerGesture` (new file `pointerGesture.go`) describes one gesture:
  a mouse right click, or a primary touch tap.
- `AppRunner.PressInOneFrame` runs a leading frame, queues every press and release, and then
  runs one trailing frame, all under the lock.
- `ZoneEditorHandler` gained `PressInOneFrame`, `EdgeRightClick`, `EdgeClick`,
  `DialogButtonClick`, `DeleteSelectedButtonClick` and `UndoButtonClick`. Each gesture is
  aimed at what is on screen when it is built.
- None of this is a `*_testexports` addition.

**Red first (recorded).** On unfixed code all 5 new tests FAILED:
- Two edges in one frame: only one was deleted. The second index was out of range after the
  list shifted.
- The same edge twice: **both** portals were deleted, the second one through the stale index.
- Delete then select: nothing was selected.
- Toolbar Delete then select: `Portal-Hub-A` was reported instead of `Portal-Hub-B`.
- Undo then select: the wrong connection was selected.

The two-edge case is the delivery proof: both presses reached the canvas within one frame.

**Fix (P4).**
- `zoneEditorGeometryState` now records `geometrySide`.
- `ensureGeometry` rebuilds when the geometry is dirty or `geometrySide != side`.
  `hitTestNode`, `hitTestEdge`, `ensureManualPositions` and `layoutCanvas` all call it, which
  replaces the local `sideChanged` check.
- `recomputeGeometry` stamps `geometrySide`. The `edgeConnection` bounds guard is kept.
- The `layoutCanvas` comment now explains the new semantic.

**Mutators.** I re-verified from source that every `working` mutator raises `geometryDirty`:
`setEditingSet` (also Revert), `addConnection`, `deleteConnection`, `undoSessionEdits`,
`deleteZone` and `applyQualityMutation`.

**Verification.**
- All 5 new tests now pass.
- Tagged suites PASS (root 3.2s, GUI 24.7s). No `.failure` snapshot was produced, so the
  existing goldens are unchanged by P4.
- `go test ./test/...` passes. `testlayoutcheck` passes.

## Phase 5: Final verification and handoff
Status: Complete

- [x] `go build ./...`; unit run with coverage (`-p=2`); total ≥ Phase 0 baseline and ≥ 74.4% (review §9 floor); per-file lines for touched files.
- [x] `go test ./test/...`; `go test -tags='integration_test,gui' ./test/integration/...`; `go run ./cmd/testlayoutcheck .`; `gofmt -l` on changed files; `golangci-lint-v2 run ./... --issues-exit-code=1` zero issues.
- [x] Independent implementation review (Claude Opus 5.5 per AGENTS §3.4).
- [x] Update [session-carry-forward.md](../session-carry-forward.md). Mark review §1.13/§1.14/§1.15/§2.1 FIXED and §9 row F complete **only after the owner commits**, per the review's protocol.

### Verification Plan
- Every command above passes and is recorded with numbers.

### Phase Summary
**Final verification** (Windows/amd64, Go 1.27.0, empty `GOFLAGS`), after the review fixes:
- `go build ./...`: PASS.
- `-p=2` unit coverage run: PASS. **Total 74.6%**, against the 74.9% baseline and the §9 floor
  of 74.4%. Deduplicated blocks put the exact figure at 6911/9246, or 74.75%.
- `go test ./test/...`: PASS.
- Tagged run: root 3.404s, GUI 27.049s, PASS. No `.failure` files.
- `testlayoutcheck`: PASS. `gofmt -l`: clean.
- `golangci-lint-v2 run ./... --issues-exit-code=1`: **0 issues**. One earlier run reported 3
  `nolintlint` hits in untouched `editorState` tests. `golangci-lint-v2 cache clean` cleared
  them, so they were a stale cache.
- Wire: no constructor changed, so it was not regenerated.

**Coverage decrease: needs owner acceptance.** Every touched non-GUI function is at 100%:
`connectionCurveLayout.go` 56/56, `connection.go` 12/12, `zoneEditorService.go` 131/131,
`connectionLineStyle.go` 9/9. The drop has two sources:
- The editor's duplicated, fully covered curve code, about 60 statements, collapsed into one
  shared model, which shrank the covered numerator.
- The new statements in `zoneEditorCanvas.go`, `zoneEditorConnectionProps.go` and
  `zoneEditorDialog.go` are covered only by the GUI suite. Those files are 0% in the unit
  profile by design, and are listed in `test_observations.md`.

Batch E's accepted decrease had the same GUI-only cause.

**Implementation review (Claude Opus 5.5): APPROVE WITH CHANGES.** No blockers and no majors.
It recomputed seven builder expectations by hand and ran the layering test (PASS). Its findings
were handled as follows:
1. **[minor] Resize-frame hit testing.** Fixed: `layoutCanvas` now calls `handlePointer` before
   updating `side`, so presses resolve against the side they were aimed at.
2. **[nit] `WasUpdated` could leak to the next connection selection.** Fixed: it is reset in
   `syncPropsFromConnection`.
3. **[nit] Curves are built twice per editor rebuild.** Accepted without change. It costs the
   same as before this batch, and removing it would need a `ConnectionIndex` on
   `preview.Connection`.
4. **[nit] Test gaps.** Added four tests: N=3 tie, N=4 obstructed, `ChangeConnectionType` to
   lowercase `portal` keeping the rules, and the dropdown listing a stored `direct` as its own
   option. Several regression guards also pass on the old code, which is accepted.
5. **[nit] Stale P1 text.** Corrected: `internal/models/preview` now also imports
   `template_model`. There is no cycle and no allow-list entry.

Review §1.13/§1.14/§1.15/§2.1 must be marked FIXED, and §9 row F set to complete, **only after
the owner commits**. Neither has been done yet.

## Final Recap
Batch F delivered all three review findings plus the owner-approved optional §2.1.

**§2.1 and §1.13.** `preview.ConnectionCurveLayout.Build` is now the single source of connection
curves for the editor, the Preview tab and the PNG export.
- It spaces parallel curves 21px apart and bends them around obstacles, in every view.
- It is deterministic: a lone curve takes the shorter detour, an exact tie bends positive, and
  obstructed parallel curves split around the obstacle.

**§1.15.** `Connection.IsEffectivePortal()` is the one classifier: Portal in any letter case, or
any placement rules.
- The editor colours and widths now match the preview, and the road tri-state is untouched.
- The type dropdown shows the effective type and lists unlisted stored types.
- The type is written only when the user actually changes the selection.
- Leaving Portal clears the placement rules.
- The old overwrite that silently turned unlisted types into Direct is gone.

**§1.14.** Hit tests rebuild stale geometry first, so an edge index always matches the live
connection list.
- It is proven by 5 real-input batched-press tests, which failed before the fix.
- The fix covers deletes, toolbar Delete, Undo, and several presses queued in one frame.

**Changed golden.** Exactly one golden changed: a Spawn-A to Spawn-B edge over the Hub, which now
curves the other way.

## Deployment Plan
**CLOSED 2026-09-28.** The owner reviewed the batch, committed it as `b9c47a7` ("Batch F"),
and authorized the close-out. The owner's review edits trimmed comments only; behaviour is
unchanged. Review §1.13, §1.14, §1.15 and §2.1 are now marked FIXED, and §9 row F is
complete. Coverage of 74.6% was accepted with the close-out.

Nothing is deployed, and nothing is pending on anyone but the owner.
1. Review the working tree: `git status --short` lists the batch's files and
   `.agent/backlog/test_observations.md`.
2. Decide whether to accept the coverage figure of 74.6%, against the 74.9% baseline and the
   74.4% floor.
3. Stage and commit. Only the owner does this.
4. After the commit, the next session marks review §1.13/§1.14/§1.15/§2.1 FIXED and sets
   §9 row F complete.
5. No schema, Wire, dependency or output-path change is involved.
6. Optionally, check in-game that templates whose connections carry placement rules still load.
   Nothing written to `.rmg.json` changed, except that the rules are cleared when a user
   explicitly changes a connection's type away from Portal.
