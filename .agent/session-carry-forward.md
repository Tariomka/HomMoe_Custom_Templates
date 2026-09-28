# Carry-forward: Batch F closed, Batch G not started

Date: 2026-09-28. **Batch F (editor geometry) is closed.** The owner reviewed it, committed
it as **`b9c47a7` ("Batch F")** on `AD/editor_geometry`, and authorized the close-out with
*"Changes reviewed, you can close out this batch and update handoff for the next session"*.
The owner's review edits before committing trimmed comments only; behaviour is unchanged.

Findings §1.13, §1.14, §1.15 and §2.1 are now **marked FIXED** in the
[surviving review](backlog/review-gpt-6-astra-09-07.md). Its §9 row F is complete, and its
progress line reads **16 fixed, 16 remaining** (0 High, 10 Medium, 6 Low). The durable
record is [batch-f-editor-geometry.md](plans/batch-f-editor-geometry.md), marked CLOSED.
Consult it only if Batch F behaviour is questioned. Do not reimplement, re-verify or
re-review any of it.

**Section 8 is preserved verbatim**, contradictions and all. Its Batch C phase/engine
wording, and its "pending decisions" wording for §1.5, §1.11, §1.12 and §2.1, are
**superseded**: all were settled, implemented and marked fixed. Everything else §8 retains
as later scope remains **binding**, and governs batches from G on. Do not edit §8.

Batches A/D, B, C, E and F are closed. A/D through E are merged to `master` (Batch E in
PR #44, tagged `v0.3.9-alpha.2`). Batch F is pushed on `AD/editor_geometry` and is not yet
merged. Do not revisit or re-verify closed batches.

## 1. Session goal

Plan, implement, verify and close Batch F from review §9: §1.13, §1.14, §1.15, plus the
owner-approved optional §2.1. This closing turn marked the findings fixed and handed this
file forward to Batch G.

## 2. Fixes applied

All Batch F items are committed in `b9c47a7` and **marked FIXED**. They are closed; do not
reopen them.

- **§1.13, deterministic obstacle routing.** One shared model,
  [connectionCurveLayout.go](../internal/models/preview/connectionCurveLayout.go), computes
  every curve. Obstacle `max`/`min` is order-independent, so map iteration can no longer
  flip a curve. The review's repro gives `(350,260)` on every one of 200 builds.
- **§1.14, stale edge identity.** `ensureGeometry` rebuilds dirty or resized geometry before
  every hit test. The rebuild is wired into
  [zoneEditorCanvas.go](../app/gui/dialogs/zoneEditorCanvas.go), plus `ensureManualPositions`
  in [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go). The side is tracked in
  [zoneEditorGeometryState.go](../app/gui/dialogs/zoneEditorGeometryState.go).
- **§1.15, portal classification.** `Connection.IsEffectivePortal()`, in
  [connection.go](../internal/models/template_model/template_variant_model/connection.go),
  is the single classifier. It is used by
  [previewLayoutService.go](../internal/services/preview_service/previewLayoutService.go)
  and [connectionLineStyle.go](../app/gui/utils/connectionLineStyle.go).
- **§2.1, shared curve geometry.** `preview.ConnectionCurveLayout.Build` is the only curve
  builder for the editor, the Preview tab and the PNG export.
- **§8 dropdown note.** The type dropdown no longer overwrites unlisted types each frame
  ([zoneEditorConnectionProps.go](../app/gui/dialogs/zoneEditorConnectionProps.go)).
- **Review follow-ups** (implementation review):
  - Presses on a resize frame are resolved at the side they were aimed at.
  - The type dropdown's `WasUpdated` flag is reset whenever a connection is synced.

## 3. Features added / changed

### Delivered in Batch F, committed and settled (full table D1–D13 in the plan)

- **Exact curve agreement.** The editor, the Preview tab and the PNG export draw identical
  control points, all built by `preview.ConnectionCurveLayout.Build`.
  - Parallel spacing is **21px** everywhere; the editor's old 18px is gone.
  - Obstacle bending (clearance `radius + 8`, padding `6`, chord margin `0.08`) now applies
    in the Preview and the PNG too.
- **Bending rules.**
  - A single obstructed curve takes the shorter detour. It goes positive (up for a
    left-to-right chord, in alphabetical endpoint order) only on an exact tie.
  - N parallel obstructed curves keep slot order and split around the obstacle: even N gives
    N/2 per side, and for odd N the extra curve goes to the shorter side. Outer curves step
    21px further out.
  - The detour is a midpoint offset. It clears every obstacle on the chosen side at the
    curve's midpoint only, a documented approximation.
- **Emission order and orientation.** Curves come out grouped by pair in first-seen order,
  with `Start`/`End` in canonical endpoint order. The Preview's connection order changed
  accordingly.
- **Effective portal**, meaning a Portal type in any letter case or any placement rules:
  - drawn with the portal colour in the editor, as in the Preview;
  - drawn at `DefaultConnectionLineSmall` in the editor, while the selected curve stays
    Large;
  - drawn with the portal shape in the Preview for a lowercase `portal`.
- **Type dropdown.**
  - It shows the effective type and lists Direct, Portal, and the stored type when it is
    unlisted. An empty type reads `(none)`.
  - The type is written only on a real selection change.
  - Choosing any non-Portal type clears both placement-rule lists; Portal keeps them.
- **Hit tests** always resolve against geometry built from the live connection list.

### Settled behavior that must stay preserved

- The road tri-state is unchanged: explicit `false` is roadless, explicit `true` is roaded,
  and `nil` is roaded only for explicit Portals (`IsExplicitPortal`). Effective
  classification never feeds `HasRoad` or road policy.
- Explicit Portal flags are preserved exactly, `false` and `nil` included. With settings
  present, road policy overwrites `Road` on every non-explicit-Portal connection,
  including custom, imported, Default, empty, arena and proximity edges.
- Internal castle/object/foothold roads are independent of the between-zone road checkbox.
  Roads-off never disables them and never removes valid Portal approaches; roads-on
  restores eligible target sets rather than connector-record equality.
- A `nil` `EditorState` preserves existing roads and content: policy runs only when the
  state is present.
- Final cleanup is zone-scoped and treats nil/empty mandatory content as authoritative,
  removing only confirmed invalid `MainObject`, incident `Connection` or named
  `MandatoryContent` references.
- The arena marker is never an anchor; spawn anchoring and rebasing of shifted imported
  non-arena anchors stay intact. Sources are cloned before mutation - no shared backing
  arrays are written through.
- The editor legend stays removed by owner decision. Preview keeps its four-entry key
  (`Road`, `No road`, `Portal`, `Portal without road`).
- PNG roadless strokes apply 50% once per edge through one reusable local mask; the opaque
  raster path, dashes and clipping are unchanged.
- The output directory stays machine-detected, with an explicit session-only picker escape
  hatch. It is never persisted, and no fallback authorizes an unrelated directory.
- All Batch E behavior (effective-mode invalidation, the two-player tournament lock, guard
  propagation and Custom) is committed and settled. Do not reopen it.
- All Batch F behavior above is committed and settled. The shared curve builder stays the
  single source of curves for the editor, the Preview and the PNG; do not reintroduce
  per-renderer curve math.

## 4. File modifications

**This closing turn is documentation only, in three files:**
- [Surviving review](backlog/review-gpt-6-astra-09-07.md): §1.13, §1.14, §1.15 and §2.1
  marked `✅ FIXED` in place, each with a Progress paragraph. The §9 row F is complete, and
  the stale 2026-09-11 progress line was refreshed to 16 fixed / 16 remaining. Nothing was
  renumbered.
- [Batch F plan](plans/batch-f-editor-geometry.md): Deployment Plan marked CLOSED.
- This handoff.

**Batch F inventory, committed in `b9c47a7`:**
**Production (new):**
- `internal/models/preview/connectionCurve.go`: the curve value type.
- `internal/models/preview/connectionCurveLayout.go`: the shared builder and curve tunables.

**Production (edited):**
- `internal/models/preview/previewConnection.go`: doc comment only.
- `internal/models/template_model/template_variant_model/connection.go`: `IsEffectivePortal`.
- `internal/services/connection_editor/zoneEditorGeometryService.go`: `buildEdges` now maps
  the shared curves. The duplicate grouping, bending and tunables were removed.
- `internal/services/connection_editor/zoneEditorService.go`: `ChangeConnectionType` clears
  the rules.
- `internal/services/preview_service/previewLayoutService.go`: uses the shared curves and the
  effective classifier.
- `app/gui/utils/connectionLineStyle.go`: `DrawsAsPortal` uses the effective check.
- `app/gui/dialogs/zoneEditorCanvas.go`: portal width, `ensureGeometry`, and input handled
  before the side update.
- `app/gui/dialogs/zoneEditorGeometryState.go`: `geometrySide`.
- `app/gui/dialogs/zoneEditorDialog.go`: `ensureManualPositions` refreshes the geometry.
- `app/gui/dialogs/zoneEditorConnectionProps.go` and
  `app/gui/dialogs/zoneEditorConnectionPropertiesState.go`: dropdown options, labels and
  writeback.

**Production (deleted):** `internal/services/connection_editor/connectionPairKey.go`. The pair
key is now a `[2]string` inside the shared model.

**Test-only exports (integration_test tag):**
- `app/gui/dialogs/zoneEditorDialog_testexports.go`: `ConnectionTypeLabels` and
  `SelectedConnectionTypeLabel`.
- `app/gui/editor/window_testexports.go`: the same two methods, declared on
  `IZoneEditorDialog`.

**Test helpers:**
- `test/test_helpers/integration_common/appRunner.go`: `PressInOneFrame` and
  `ApplyManualConnectionEdit`.
- `test/test_helpers/integration_common/pointerGesture.go` (new).
- `test/test_helpers/integration_common/zoneEditorHandler.go`: batch wrappers.

**Tests:** listed in §5, plus one golden:
`test/test_helpers/integration_common/snapshot/__snapshots__/zoneEditorActions_integration_test/TestWhenADragStartsOnAZoneInAddConnectionMode_AConnectionIsCreated_5.golden`.
The Spawn-A to Spawn-B edge passes exactly over the Hub, so it now bends to the positive
side on the tie.

**Agent docs:**
- `.agent/plans/batch-f-editor-geometry.md` (new, durable).
- `.agent/backlog/test_observations.md`: portal width cannot be pixel-tested, and dropdown
  shaping is covered by the GUI suite.
- `.agent/memories/gui-and-tests.md`: batch input, seeding manual connections, the `-update`
  trap, and Dp rounding.
- This handoff.

**Not touched:** `data/`, `internal/entities/template_entity/`, `internal/registry/`, the
output path, generated Wire, topology code, dependencies, and
`.agent/backlog/owner_findings.md`.

## 5. Tests added or updated

**This closing turn added no tests and reran nothing**, because it changed documentation
only. The owner's pre-commit edits trimmed comments only. The Batch F verification below is
the evidence of record for `b9c47a7`.

**Unit tests** (each file mirrors its implementation path, one test file per public function):
- New `test/unit/internal/models/preview/connectionCurveLayout/build_test.go` (22 tests) plus
  `common_test.go`. It covers:
  - the fan and spacing;
  - reversed pairs, canonical orientation, grouped order and missing endpoints;
  - stacked endpoints, and a chord shorter than 1px;
  - single obstacles on either side, the on-chord tie, and the review repro repeated 200
    times;
  - obstacles on both sides;
  - N = 2, 3 (each side shorter), 3 (tie) and 4 obstructed;
  - the off-centre midpoint approximation, the chord margin, and the clearance boundary.
- New `test/unit/.../template_variant_model/connection/isEffectivePortal_test.go`.
- `zoneEditorService/changeConnectionType_test.go`: 5 new tests for rule clearing and
  keeping, including lowercase `portal`.
- `zoneEditorGeometryService/buildGeometry_test.go`: moved to 21px, plus a parity test.
  `common_test.go` gained `controlPoints`.
- `previewLayoutService/buildPreviewLayout_test.go`:
  - 4 new tests: bending, parity, grouped order, and canonical `Start`;
  - the lowercase-portal test was flipped to portal-shaped and renamed.

  `common_test.go` gained `obstructedManualZones`.
- `connectionLineStyle/newEditorConnectionLineStyle_test.go`: 2 new rule-only portal cases.

**Integration and GUI tests** (all tagged `integration_test && gui`; no new goldens):
- New `test/integration/gui/zoneEditorConnectionType_integration_test.go`: 10 tests, one of
  them a 4-case table.
- `zoneEditorPointer_integration_test.go`: 5 batched-press tests. All five **failed on the
  unfixed code**, and the failure modes are recorded in the plan.
- `roadStyleVisuals_integration_test.go`: 3 pixel tests for rule-only portal colours.
- `zoneEditorGeometry_integration_test.go`: 3 updated expectations.

**Baseline** (before the first edit), Windows/amd64, Go 1.27.0, empty `GOFLAGS`: every check
PASSED, coverage **74.9%**, tagged run root 3.498s / GUI 26.655s, lint 0 issues.

**Final**, after the review fixes:
- `go build ./...`: PASS.
- `-p=2` unit coverage run: PASS, **74.6%** total. By deduplicated blocks the figure is
  6911/9246, or 74.75%.
- `go test ./test/...`: PASS.
- `go test -tags='integration_test,gui' ./test/integration/...`: PASS, root 3.404s, GUI
  27.049s.
- `testlayoutcheck`: PASS. `gofmt -l`: clean.
- `golangci-lint-v2 run ./... --issues-exit-code=1`: 0 issues. One run showed 3 stale-cache
  `nolintlint` hits in untouched `editorState` tests; `golangci-lint-v2 cache clean` cleared
  them.

**Coverage decrease: accepted with the close-out.** It is above the §9 floor of 74.4%. Every
touched non-GUI function is at 100%. The drop comes from collapsing the duplicated, fully
covered curve code, and from new dialog statements that only the GUI suite exercises.
**74.6% is the new baseline for Batch G.**

**Not measured:** native Linux and Steam Deck (the WSL probe found no `go`/`pkg-config`
earlier), in-game behaviour, and any coverage-profile fingerprint.

**Reviews:**
- GPT-6 Sol plan review: REJECT on the first draft; all findings were addressed and one was
  declined with a reason.
- Claude Opus 5.5 implementation review: APPROVE WITH CHANGES. Its 1 minor and 4 nits were
  resolved, or noted where no change was made (see the plan's Phase 5).

## 6. Git status snapshot

The branch is **`AD/editor_geometry`**, HEAD **`b9c47a7` ("Batch F")**, level with
`origin/AD/editor_geometry`. The history on top of `master` (`4143bdd`) is `b9fb53e`
("Agent"), then `3e49ca2` and `4a72813` (plan snapshots), then `b9c47a7`. Batch F is not
yet merged to `master`.

The tree was clean when this closing turn began. `git status --short` now shows only this
turn's three documentation edits, all unstaged:

```text
 M .agent/backlog/review-gpt-6-astra-09-07.md
 M .agent/plans/batch-f-editor-geometry.md
 M .agent/session-carry-forward.md
```

`.agent/memories/` is gitignored and does not show.

**The assistant performed no staging, unstaging, commit, push, stash, branch switch or
worktree change in this turn.**

## 7. Rejections / things the user declined

- **Batch F, settled by the owner.**
  - Editor-only obstacle bending, and dropping obstacle bending entirely, were both
    rejected. Bending now applies everywhere.
  - The editor's 18px spacing was rejected in favour of 21px.
  - Always bending a single curve to the positive side was rejected in favour of the shorter
    detour, with positive only on a tie.
  - A display-only type change was rejected: leaving Portal clears the rules.
  - A case-sensitive preview shape check was rejected.
  - Keeping the current editor portal width was rejected: portals are drawn thinner.
  - Listing every registry type, or keeping two items with a silent Direct fallback, were
    both rejected in favour of Direct, Portal, plus the unlisted stored type.
- **Declined review items.**
  - Plan review (GPT-6 Sol): the "hidden pair members shift slots" finding was declined.
    Every member of a pair shares both endpoint names, so a pair is either fully visible or
    fully skipped.
  - Implementation review nit 3 (curves built twice per editor rebuild) was noted without
    change.
- **Do not reuse a blanket `-update` snapshot run.** It rewrote about 280 unrelated goldens,
  which were restored. Accept only the `.failure` files a plain run produces.
- **Still in force from earlier batches:**
  - Do not restore the editor legend or its deleted test.
  - Do not reopen Batch C or Batch E decisions: arena recomputation, remembered player
    counts, persisted guard preset identity, the Plastic/Bronze tables, `GeneratorConfig`
    player rejection.
  - Keep the Batch E safeguards.
  - Topology retirement (Batch K) stays out of scope.
  - No opportunistic DTO, schema, package, allocation or output-path work.
  - Never claim unobserved engine outcomes.
- **§8 verification encoding.** Use explicit UTF-8 for Git stdout and disk reads, because
  PowerShell 5.1 mis-decodes. Git stores this file LF-normalized while the worktree is
  CRLF. Before and after this edit, §8's UTF-8 SHA-256 (from `## 8.` up to `## 9.`) is
  `35D5F4F36C0492166B202DC78828A0F0594906FA3A08A899F554D8A68D118E13`.
- **Recorded, not assigned:**
  - A service comment still names the legacy rebuild entry point.
  - A zero-edge tournament count remark.
  - The defensive selected-pointer rebinding branch.
  - Double curve construction per editor rebuild.

## 8. Open questions and confirmed later scope

Batch C policy/scope approval remains settled. **Phase 3 is complete with no blockers.**
Owner Phase 3 authorization, 2026-09-10: "Changes reviewed, you can proceed".
Phase 4 remains unstarted and belongs to a new session. Do not repeat settled questions.

Non-blocking review observations, outside closeout scope:

- `ConnectionLegendRow` remains an export used only by `LegendRows` after editor legend
  removal. No behavior issue; optional simplification is not required for Phase 4.
- Pre-existing editor dropdown normalization: its non-`WasUpdated` writeback can turn
  types not offered by the Direct/Portal dropdown (for example Proximity) into Direct
  on selection. This predates Phase 3; road policy still stamps correctly. Investigate
  separately if requested, do not fix opportunistically in PNG work. See
  [writebackProps](../app/gui/dialogs/zoneEditorConnectionProps.go#L105-L119).

The owner will validate true/false/nil engine behavior after changes. Existing engine
evidence is inconclusive: schema is optional bool; shipped examples use both values
for direct and portal connections; examples do not prove defaults.

Retained owner decisions for later work:

- §2.3/O08: remove Ring, Hub, Chain, Shared Web and corresponding tournament builders; reject retired saved IDs without modifying the current document; surviving tournament fallback cases use balanced generation. Preserve Geometric Hub and shared hub-zone concepts.
- §2.4/O11: investigate live state and persistence first; compare exact-base reconstruction against regenerated untouched zones plus deltas. Cover identity, randomness, upgrades, defaults, derived fields, migration and measured size/performance. Owner approves feasibility before representation changes.
- §2.5/O12: structured entries following BonusEntry; bans carry SID, overrides SID/GuardValue, preserve semantics and Variant=-1. Legacy/UI text parsing stays at boundaries.
- §2.6/O13: section current groups except TemplateIdentity, MapSettings and SchemaOptions remain flat. Coordinate migration with §2.5 and approved §2.4 outcome; retain supported legacy loading. Section keys need planning.
- §2.7/O14: all six zone-content accessors on drivers.State, replacing panel callbacks with zone/tier selection; retain validation, dirty tracking and snapshot isolation.
- §2.8/O16: rename services/zones to zone_services and nested interfaces to zone_service_interfaces, updating imports/tests/links/generated wiring without behavior change.
- §2.9/O18: General/Layout/Bonuses named subpackages with private cohesive sections and public aliases/constructor compatibility; Preview unchanged.
- §2.10/O19: public GeneratorConfig lookup preserving Highest→Hub, unknown→nil and read-only borrowed-row semantics.
- §4.1/O20: permitted production Vec2 audit, equivalent methods, justified tested additions, numerical and hot-path behavior preserved.
- §2.2: zone-content DTO removal reopened; bonuses exception remains accepted. Replacement row/result API and full-exception versus single-seam scope still need approval.
- O09 presets remain deferred and uninvestigated.

Other pending decisions: arena/manual invalidation and effective-mode aliases (§1.5), custom guards/preset identity (§1.11), previous non-tournament player count (§1.12), GUI/PNG curve agreement (§2.1), tooling version/EOL/release tags (§6.2/§6.3/§6.5). Review §10 retains in-game bonuses/bans, hero-hire and historical preview validation. Previous tidy differences were Windows checksum EOL-only, not dependency drift.

## 9. Next recommended actions

**Batch F is closed, with no open Batch F questions and no blockers.** The next unit of
work is **Batch G, measured performance: review §3.1**, taken from the surviving review's §9
table rather than from §8. Nothing has been planned or implemented for Batch G, and no code
may be written before the gates below are cleared.

**§3.1 in brief.** The zone editor's status line re-runs graph diagnostics on every frame.
`layoutStatus` in [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go) calls
`derefConnections(this.working)`, which allocates a slice, and then
`DescribeZoneEditorGraph` ([zoneEditorHandler.go](../internal/handlers/zoneEditorHandler.go)).
Its `FindIsolatedZones` ([connectionEditorService.go](../internal/services/connection_editor/connectionEditorService.go))
is O(zones × connections). The review's proposed fix:
- cache a graph summary, invalidated by structural changes;
- keep in mind that toolbar buttons are handled before canvas input, so a cached summary
  must still be correct on the next frame;
- optionally make isolation detection O(zones + connections) with a referenced-name set;
- **measure first** with an untagged public-API benchmark at
  `test/performance/zoneEditorGraph_test.go`, and count handler calls on idle versus
  mutation frames through GUI integration.

**Owner constraint from the review:** do not reintroduce the rejected global Gio
allocation threshold, or clone-free live state pointers.

**Batch F context that affects G:** `geometryDirty` is now also cleared by `ensureGeometry`
from inside hit tests, in the middle of input handling. A graph-summary cache must **not**
piggyback on that flag. It needs its own invalidation, or it will miss changes that a hit
test has already absorbed.

Routing for the next session, in order:

1. Read [AGENTS.md](../AGENTS.md), this handoff, then review §3.1 and the §9 `G` row.
2. **Inspect the current code and tests yourself**, because the review's line numbers
   predate Batch F:
   - [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go) (`layoutStatus`,
     `derefConnections` and every `working` mutator);
   - [zoneEditorCanvas.go](../app/gui/dialogs/zoneEditorCanvas.go) (`ensureGeometry`);
   - [zoneEditorHandler.go](../internal/handlers/zoneEditorHandler.go);
   - [connectionEditorService.go](../internal/services/connection_editor/connectionEditorService.go);
   - [findIsolatedZones_test.go](../test/unit/internal/services/connection_editor/connectionEditorService/findIsolatedZones_test.go);
   - [describeZoneEditorGraph_test.go](../test/unit/internal/handlers/zoneEditorHandler/describeZoneEditorGraph_test.go);
   - the existing benchmarks under `test/performance/`.
3. **Ask the owner** the undecided questions:
   - Is measurement alone enough, or should the batch also implement the fix?
   - What measured improvement justifies the fix?
   - Is the O(zones + connections) isolation rewrite in scope?
   - Where should the summary cache live: the dialog, or behind the handler?
   - Is a GUI handler-call count acceptable as an idle-frame assertion?
4. Summarize the scope back to the owner and get approval.
5. Write a **new durable Batch G plan** under `.agent/plans/`, get an independent review
   and explicit plan approval, and capture a fresh baseline before the first edit.

**Deployment.** Nothing is deployed and nothing is authorized to be. The owner alone
stages, commits, merges and releases. Batch F still awaits the owner's merge to `master`.
No schema migration, dependency installation, Wire regeneration or output-directory
change is required or pending. Native Linux and Steam Deck execution remain unmeasured.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md). It
> is self-contained.
>
> **Batch F is CLOSED.** The owner reviewed it and committed it as `b9c47a7`. Findings
> §1.13, §1.14, §1.15 and §2.1 are marked **FIXED** in the
> [surviving review](backlog/review-gpt-6-astra-09-07.md), whose §9 table shows batch `F`
> complete, with 16 fixed and 16 remaining. Coverage of 74.6% was accepted and is the new
> baseline. Do not reimplement, re-verify or re-review Batch F, and do not reopen its
> owner decisions (D1–D13 in [the closed plan](plans/batch-f-editor-geometry.md)).
>
> **State:** branch `AD/editor_geometry`, HEAD `b9c47a7`, level with origin, not yet merged
> to `master`. The only working-tree changes are the closing turn's three documentation
> edits: the review, the plan and this handoff. The assistant performed no Git mutation;
> preserve the owner's state exactly.
>
> **Next work is Batch G, measured performance, review §3.1:** zone-editor graph
> diagnostics are rebuilt on every frame. **No Batch G plan or implementation exists yet,
> and none may be written before these gates, in order:**
> 1. Read review §3.1 and §9.
> 2. Inspect the current code yourself (list in §9 above). Batch F moved lines.
> 3. Ask the owner: measure-only or also fix; the improvement threshold; the scope of the
>    O(n) isolation rewrite; where the cache lives; and whether GUI handler-call counts are
>    acceptable.
> 4. Summarize the scope and get approval.
> 5. Write a durable Batch G plan, get an independent review and explicit plan approval.
> 6. Capture a fresh baseline before the first edit.
>
> Do not reuse `geometryDirty` for the summary cache: hit tests clear it. Do not
> reintroduce the global Gio allocation threshold, or clone-free live state pointers.
>
> **Verification on record** (Batch F final code, Windows/amd64, Go 1.27.0, empty
> `GOFLAGS`; not rerun at close-out):
> - build: PASS;
> - `-p=2` unit coverage run: PASS, 74.6%;
> - `go test ./test/...`: PASS;
> - tagged `integration_test,gui` run: PASS, root 3.404s, GUI 27.049s;
> - `testlayoutcheck`: PASS; `gofmt`: clean;
> - lint: 0 issues. If stale `nolintlint` hits appear, run `golangci-lint-v2 cache clean`
>   first.
>
> **Native Linux and Steam Deck are UNAVAILABLE.** No in-game result is claimed.
>
> **GUI snapshots:** never keep a blanket `-update`, because it rewrites about 280 goldens.
> Accept only the `.failure` files a plain run produces.
>
> **Hard rules:**
> - Never modify protected data, the template schema or the registry.
> - Keep Windows/Linux compatibility.
> - Never change or persist the machine-detected output directory.
> - Test nontrivial logic, and check coverage before and after.
> - Never stage, unstage, commit, push, stash, switch branches or manipulate worktrees;
>   preserve owner changes.
> - Never bulk-rewrite or hand-edit generated Wire.
> - Never enable global `integration_test`, `gui` or `wireinject` tags, and never
>   introduce fake unit seams.
> - Keep plans durable and resumable.
>
> **Out of scope:** Batch K topology retirement, direct `GeneratorConfig` rejection, DTO
> cleanup, schema and package work, allocation tuning beyond §3.1, and every settled
> alternative in §7. §8 stays verbatim. Its Batch C wording and its pending-decision wording
> for §1.5, §1.11, §1.12 and §2.1 are superseded; its retained later-scope decisions remain
> binding.
>
> **Preserve:**
> - explicit Portal Road `false`/`nil` and valid approaches;
> - the tri-state road display (`nil` is roaded only for explicit Portals, and effective
>   classification never feeds roads);
> - independent internal roads, nil-state content/road preservation, and source cloning;
> - the single reusable half-opacity PNG edge mask, and the Preview-only legend;
> - the shared curve builder, as the single source of curves;
> - hit tests resolving against the live connection list.
>
> This handoff contains the full continuation context.
