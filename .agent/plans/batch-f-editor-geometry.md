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
`internal/models/preview/` (no cycle: `preview` imports only `helpers/data` and
`neutral_zone`; `template_model` does not import `preview`):
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
Status: Not started

- [ ] Record `git status --short`, branch and HEAD (expected: `AD/editor_geometry`, clean, `b9fb53e` or the owner's newer commit).
- [ ] `go build ./...`
- [ ] `go test -p=2 -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...` then `go tool cover '-func=coverage.txt'`; record total and the per-file lines for `zoneEditorGeometryService.go`, `previewLayoutService.go`, `connection.go` (template_variant_model), `zoneEditorService.go`, `connectionLineStyle.go`.
- [ ] `go test ./test/...`
- [ ] `go test -tags='integration_test,gui' ./test/integration/...` (record timings)
- [ ] `go run ./cmd/testlayoutcheck .`
- [ ] `golangci-lint-v2 run ./... --issues-exit-code=1`

### Verification Plan
- All commands pass; baseline numbers written into the Phase Summary. Any pre-existing failure stops the batch and is reported to the owner.

### Phase Summary
_(write when phase completes)_

## Phase 1: Shared deterministic curve builder (§1.13, §2.1)
Status: Not started

- [ ] Add `internal/models/preview/connectionCurve.go` and `connectionCurveLayout.go` per P1–P3 (21px spacing, obstacle bending, D4/D5 split, grouped first-seen order, canonical orientation, skip pairs with a missing endpoint, `distance < 1` clamps to 1 as today).
- [ ] Unit tests `test/unit/internal/models/preview/connectionCurveLayout/build_test.go` (`package connectionCurveLayout_test`, triple-A, `t.Parallel`, one assertion each):
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
- [ ] `PreviewLayoutService.buildPreviewConnections` builds `preview.Connection` from `ConnectionCurveLayout{Positions, ZoneRadius}.Build`, keeping `Type`/`HasRoad`/`ExplicitPortal` mapping as is.
- [ ] `ZoneEditorGeometryService.BuildGeometry` builds edges from the same `Build` (label point `0.25·S + 0.5·C + 0.25·E` unchanged); delete `buildEdges` curve math, `groupConnectionsByPair`, `obstacleBulge`, the moved tunables and the now-unused `connectionPairKey` type file if nothing else uses it. Hit-test and snap tunables stay.
- [ ] Update `buildGeometry_test.go` expectations (21px spread, bending values) and keep tests for verbatim positions/radius, preview-layout inputs, label midpoint, missing endpoint, stacked endpoints; add one asserting editor control points equal the shared builder's for an obstructed fixture.
- [ ] Update `buildPreviewLayout_test.go`: spacing expectations, a new manual-position fixture proving the Preview now bends around an obstacle, any order-dependent assertion affected by P2, plus explicit tests that Preview connections come out in grouped first-seen order and that a reversed connection keeps canonical `Start`/`End`.
- [ ] Regenerate affected Preview-tab and zone-editor GUI snapshots and PNG expectations; update `zoneEditorGeometry_integration_test.go` control points (e.g. the ring spread). List every changed golden + reason.

### Verification Plan
- `go build ./...`; `go test ./test/unit/... -count=1` passes.
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` passes after snapshot regeneration; `git status --short` lists only the intended golden files.
- A dedicated test builds the review repro 200 times with freshly constructed maps and asserts a single control point.
- Coverage for `connectionCurveLayout.go` 100% of statements; editor/preview service files not below baseline.

### Phase Summary
_(write when phase completes)_

## Phase 2: Effective portal classification (§1.15)
Status: Not started

- [ ] Add `IsEffectivePortal()` to `internal/models/template_model/template_variant_model/connection.go` (case-insensitive Portal or any placement rule). `IsExplicitPortal`/`HasRoad` untouched (D9).
- [ ] Tests `test/unit/.../template_variant_model/connection/isEffectivePortal_test.go`: `Portal`, `portal`, `Direct`, `Direct`+From rules, `Direct`+To rules, empty type.
- [ ] `getPreviewConnectionType` uses `IsEffectivePortal()`; GladiatorArena/Proximity matching unchanged. Flip `TestWhenAPortalTypeIsLowerCased_...` to portal-shaped (rename per convention) and keep the placement-rule roadless test.
- [ ] `NewEditorConnectionLineStyle`: `DrawsAsPortal = IsEffectivePortal()`, `ExplicitPortal = IsExplicitPortal()`; update `newEditorConnectionLineStyle_test.go` (rule-only portal draws as portal; rule-only portal with `nil` road stays roadless `NoRoadLine`).
- [ ] `drawEdges`: portal width `DefaultConnectionLineSmall`, selected `DefaultConnectionLineLarge`, otherwise `DefaultConnectionLine`.
- [ ] GUI tests (`roadStyleVisuals_integration_test.go` helpers `edgeStyleCensus`/`edgeStrokeWidth`): explicit portal, placement-rule-only roaded portal and ordinary direct edge show the same colours in editor and Preview; editor portal stroke thinner than direct; selected portal stroke larger.

### Verification Plan
- Unit + tagged GUI suites pass; `connection.go` and `connectionLineStyle.go` coverage 100% for changed functions.

### Phase Summary
_(write when phase completes)_

## Phase 3: Type dropdown and placement-rule clearing (§1.15, §8 note)
Status: Not started

- [ ] `ZoneEditorService.ChangeConnectionType`: after setting the type, clear `PortalPlacementRulesFrom`/`To` when the new type is not Portal (case-insensitive registry compare); then road stamp as today.
- [ ] Extend `test/unit/internal/services/connection_editor/zoneEditorService/changeConnectionType_test.go`: Direct clears From; Direct clears To; Proximity clears; Portal keeps both; road still stamped.
- [ ] Dialog `syncPropsFromConnection`: build per-connection items (`Direct`, `Portal`, + stored unlisted type, `(none)` for empty) and select the effective type (D12); keep the item→stored-value list on dialog state.
- [ ] Dialog `writebackProps`: call `ChangeZoneEditorConnectionType` only on `WasUpdated`, with the stored value of the chosen item; remove the non-`WasUpdated` assignment (D13).
- [ ] GUI tests in `zoneEditorProperties_integration_test.go` (dialog built with explicit connections, real clicks): Proximity/GladiatorArena/empty/lowercase-portal survive selection + idle frames + Apply unchanged; stored-Direct-with-rules reads `Portal` and survives selection + idle frames + Apply with its rules and `Direct` type; choosing Direct on it applies a Direct connection with no rules; then choosing Portal again applies an explicit Portal with no rules; choosing Portal on a Proximity-with-rules keeps the rules; `(none)` shown for empty type.

### Verification Plan
- Unit + tagged GUI suites pass; `ChangeConnectionType` 100% covered.

### Phase Summary
_(write when phase completes)_

## Phase 4: Batched pointer identity (§1.14)
Status: Not started

- [ ] Harness (`test/test_helpers/integration_common`, not `*_testexports`): one locked `AppRunner` method that runs a leading frame, queues several press/release gestures (mouse secondary or touch primary, per gesture), then runs a single trailing frame; `ZoneEditorHandler` wrappers (right-click several edges / click several canvas points in one frame; toolbar button + canvas press in one frame).
- [ ] Write the GUI tests first in `zoneEditorPointer_integration_test.go` and confirm at least the deletion-shift case **fails on current code** (record the failure). If none reproduces, stop and report to the owner before changing production code.
  - same edge right-clicked twice in one frame → exactly that one connection deleted;
  - two different edges right-clicked in one frame → exactly those two deleted;
  - right-click edge X then left-click edge Y in one frame → Y selected (by name);
  - toolbar Delete of a selected connection + queued canvas press on another edge in one frame → the pressed edge is the one selected;
  - toolbar Undo after a deletion + queued canvas press on an edge in one frame → the pressed edge (by name) is selected.
  Delivery proof: the "two different edges → exactly those two deleted" case shows both presses reached the canvas in one frame; a same-edge result is only trusted once that case passes.
- [ ] Implement P4: `ensureGeometry` (dirty or side mismatch) before hit tests and in `ensureManualPositions`; update the `layoutCanvas` comment; keep the `edgeConnection` bounds guard.
- [ ] Re-verify every `working` mutator raises `geometryDirty`.

### Verification Plan
- New tests fail before and pass after; full tagged GUI suite passes; `go run ./cmd/testlayoutcheck .` passes.

### Phase Summary
_(write when phase completes)_

## Phase 5: Final verification and handoff
Status: Not started

- [ ] `go build ./...`; unit run with coverage (`-p=2`); total ≥ Phase 0 baseline and ≥ 74.4% (review §9 floor); per-file lines for touched files.
- [ ] `go test ./test/...`; `go test -tags='integration_test,gui' ./test/integration/...`; `go run ./cmd/testlayoutcheck .`; `gofmt -l` on changed files; `golangci-lint-v2 run ./... --issues-exit-code=1` zero issues.
- [ ] Independent implementation review (Claude Opus 5.5 per AGENTS §3.4).
- [ ] Update [session-carry-forward.md](../session-carry-forward.md). Mark review §1.13/§1.14/§1.15/§2.1 FIXED and §9 row F complete **only after the owner commits**, per the review's protocol.

### Verification Plan
- Every command above passes and is recorded with numbers.

### Phase Summary
_(write when phase completes)_

## Final Recap
_(write when all phases complete)_

## Deployment Plan
_(write when all phases complete)_
