# Carry-forward: Batch F implemented, awaiting owner review and commit

Date: 2026-09-28. **Batch F (editor geometry) is implemented and verified** in the working
tree of `AD/editor_geometry`, and is **uncommitted**. The owner said they would review
everything at once after all phases were done. The durable plan is
[batch-f-editor-geometry.md](plans/batch-f-editor-geometry.md). It holds the binding owner
decisions D1–D13, the design P1–P6, per-phase summaries, every changed golden with its
reason, and the review dispositions. Read it before touching anything in this batch.

The review findings §1.13, §1.14, §1.15 and §2.1 are **not yet marked FIXED**, and the
[surviving review](backlog/review-gpt-6-astra-09-07.md)'s §9 row F is **not yet complete**.
Under the review's protocol, both happen only after the owner commits.

**Section 8 is preserved verbatim**, contradictions and all. Its Batch C phase/engine
wording and its "pending decisions" wording for §1.5/§1.11/§1.12 and §2.1 are
**superseded**: all four were settled and implemented. Everything else §8 retains as later
scope remains **binding**, and governs batches after F. Do not edit §8.

Batches A/D, B, C and E are closed and merged to `master` (Batch E in PR #44, tagged
`v0.3.9-alpha.2`). Do not revisit, re-verify or retrieve retired documents from those
batches.

## 1. Session goal

Plan and implement Batch F from review §9: §1.13 deterministic obstacle routing, §1.14
stale edge identity under batched pointer events, §1.15 portal classification, and the
owner-approved optional §2.1 shared curve geometry. The batch also absorbed the §8
dropdown-normalization note.

The owner answered every scoping question and approved the scope and the plan. The plan
passed an independent GPT-6 Sol review (first draft rejected, findings folded in). The
owner instructed the agent to run Phases 0–5 without stopping.

## 2. Fixes applied

All of these are implemented and tested, and passed an independent Claude Opus 5.5 review.
They are uncommitted and not yet marked FIXED.

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
- **§8 dropdown note.** The type dropdown no longer overwrites unlisted types each frame
  ([zoneEditorConnectionProps.go](../app/gui/dialogs/zoneEditorConnectionProps.go)).
- **Review follow-ups** (implementation review):
  - Presses on a resize frame are resolved at the side they were aimed at.
  - The type dropdown's `WasUpdated` flag is reset whenever a connection is synced.

## 3. Features added / changed

### Delivered in Batch F (all owner-decided; full table D1–D13 in the plan)

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

## 4. File modifications

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
`.agent/backlog/owner_findings.md`. The surviving review is deliberately unmarked until the
owner commits.

## 5. Tests added or updated

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

**Coverage decrease: needs owner acceptance.** It is above the §9 floor of 74.4%. Every
touched non-GUI function is at 100%. The drop comes from collapsing the duplicated, fully
covered curve code, and from new dialog statements that only the GUI suite exercises.

**Not measured:** native Linux and Steam Deck (the WSL probe found no `go`/`pkg-config`
earlier), in-game behaviour, and any coverage-profile fingerprint.

**Reviews:**
- GPT-6 Sol plan review: REJECT on the first draft; all findings were addressed and one was
  declined with a reason.
- Claude Opus 5.5 implementation review: APPROVE WITH CHANGES. Its 1 minor and 4 nits were
  resolved, or noted where no change was made (see the plan's Phase 5).

## 6. Git status snapshot

The branch is **`AD/editor_geometry`**, HEAD `4a72813` ("plan"), level with
`origin/AD/editor_geometry`. During the session the owner committed two plan snapshots,
`3e49ca2` and `4a72813`, on top of `b9fb53e` ("Agent"). Batch E already reached `master` as
`2d18f9e` (PR #44). The previous handoff's `AD/modes_and_guard_propagation` state is
historical, and its staged plan deletion is gone.

Nothing is staged. `git status --short` at hand-off shows 35 entries, all unstaged:
- **Modified:** the production, testexport, helper, test and golden files listed in §4,
  plus `.agent/backlog/test_observations.md`, this handoff, and the plan. The plan is
  tracked, and the working copy holds all phase summaries beyond the owner's `4a72813`
  snapshot.
- **Deleted:** `D internal/services/connection_editor/connectionPairKey.go`.
- **Untracked:**
  - `internal/models/preview/connectionCurve.go`
  - `internal/models/preview/connectionCurveLayout.go`
  - `test/integration/gui/zoneEditorConnectionType_integration_test.go`
  - `test/test_helpers/integration_common/pointerGesture.go`
  - `test/unit/internal/models/preview/connectionCurveLayout/`
  - `test/unit/.../connection/isEffectivePortal_test.go`

`.agent/memories/` is gitignored, so it does not show.

**The assistant performed no staging, unstaging, commit, push, stash, branch switch or
worktree change.** The one index-safe restore it ran was `git restore --worktree` on the
snapshot directory. It only reverted the goldens the agent's own `-update` run had just
rewritten, never owner content.

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
  CRLF. Before this edit, §8's UTF-8 SHA-256 (from `## 8.` up to `## 9.`) was
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

**Batch F is implemented, verified and awaiting the owner.** In order:

1. **Owner review of the whole batch.** Start from the plan's Final Recap and phase
   summaries, [batch-f-editor-geometry.md](plans/batch-f-editor-geometry.md).
2. **Owner decision on coverage.** Accept **74.6%** against the 74.9% baseline; it is above
   the §9 floor of 74.4%, and the cause is explained in §5. If it is not accepted, the only
   lever is GUI-only dialog code, which unit tests cannot reach without fake seams.
3. **Owner stages and commits.** The agent never does this.
4. **After the commit, and only then,** mark review §1.13, §1.14, §1.15 and §2.1 **FIXED**
   in place, set the §9 row `F` to complete, and refresh this handoff. Follow the review's
   own marking protocol and never renumber.
5. **Optional engine check by the owner.** Load a template whose connections carry
   placement rules, and one where the editor cleared them via a type change, and confirm
   the game still treats them as expected. Nothing about the `.rmg.json` schema changed.
6. **Then the next batch from review §9:** **G, measured performance (§3.1)**, unless the
   owner reprioritizes. §9 row P (Vec2 audit) said to coordinate with F. That is no longer
   blocking: the shared builder already uses `Vec2` operations throughout.

**Deployment.** Nothing is deployed and nothing is authorized to be. The owner alone
stages, commits and releases. No schema migration, dependency installation, Wire
regeneration or output-directory change is required or pending. Native Linux and Steam
Deck execution remain unmeasured.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first. Then read [this handoff](session-carry-forward.md)
> and the Batch F plan, [batch-f-editor-geometry.md](plans/batch-f-editor-geometry.md).
>
> **Batch F is IMPLEMENTED AND VERIFIED, NOT YET COMMITTED.** It covers review §1.13
> deterministic obstacle routing, §1.14 stale edge identity under batched pointer events,
> §1.15 portal classification, the owner-approved optional §2.1 shared curve geometry, and
> the §8 dropdown note. All phases 0–5 are complete. The owner is reviewing the whole batch
> at once. Do not reimplement, re-scope or re-review it, and do not reopen the owner
> decisions D1–D13 recorded in the plan.
>
> **State:** branch `AD/editor_geometry`, HEAD `4a72813` ("plan"), level with origin. Every
> Batch F change sits unstaged in the working tree (35 status entries; see §6). The
> assistant performed no staging, commit or push; preserve the owner's state exactly.
>
> **Your first job depends on what the owner reports:**
> - **If they committed:** mark review §1.13, §1.14, §1.15 and §2.1 **FIXED** in place, set
>   the §9 row `F` to complete, and refresh this handoff. Follow the review's own protocol
>   and never renumber.
> - **If they asked for changes:** make exactly those changes, and keep the plan's phase
>   summaries current.
> - **Either way, first confirm** whether they accept coverage at **74.6%**, against the
>   74.9% baseline and the §9 floor of 74.4%. The decrease is GUI-only, and every touched
>   non-GUI function is at 100%.
>
> After that, the next batch is **G, measured performance (§3.1)**. Follow the usual gates:
> read, inspect the code yourself, ask the owner, write a durable plan, get an independent
> review, get plan approval, capture a baseline.
>
> **Verification on record** (Windows/amd64, Go 1.27.0, empty `GOFLAGS`, final code):
> - build: PASS;
> - `-p=2` unit coverage run: PASS, 74.6%;
> - `go test ./test/...`: PASS;
> - tagged `integration_test,gui` run: PASS, root 3.404s, GUI 27.049s;
> - `testlayoutcheck`: PASS; `gofmt`: clean;
> - `golangci-lint-v2 run ./... --issues-exit-code=1`: 0 issues, after a
>   `golangci-lint-v2 cache clean` that cleared stale hits;
> - Wire: not regenerated, because no constructor changed.
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
> cleanup, schema and package work, allocation tuning, and every settled alternative in §7.
> §8 stays verbatim. Its Batch C wording and its pending-decision wording for §1.5, §1.11,
> §1.12 and §2.1 are superseded; its retained later-scope decisions remain binding.
>
> **Preserve:**
> - explicit Portal Road `false`/`nil` and valid approaches;
> - the tri-state road display (`nil` is roaded only for explicit Portals, and effective
>   classification never feeds roads);
> - independent internal roads, nil-state content/road preservation, and source cloning;
> - the single reusable half-opacity PNG edge mask, and the Preview-only legend;
> - the shared curve builder, as the single source of curves.
>
> This handoff contains the full continuation context.
