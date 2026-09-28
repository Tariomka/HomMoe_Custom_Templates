# Carry-forward: Batch G implemented, awaiting owner review

Date: 2026-09-28. **Batch G (review §3.1, measured performance) is implemented and
verified on Windows, but not committed.** The owner approved the plan and amendment G9a
with *"You can capture the baseline and proceed. Go through all of the Phases and I will
review everything at the end"*. Every change is unstaged in the working tree of
`AD/performance`. The durable record is
[batch-g-editor-graph-diagnostics.md](plans/batch-g-editor-graph-diagnostics.md): owner
decisions G1–G9a, phase summaries, before/after benchmark medians, and both reviews.

**§3.1 is NOT yet marked FIXED** in the
[surviving review](backlog/review-gpt-6-astra-09-07.md). Its fix-session protocol marks an
item only after the owner commits, so the progress line still reads 16 fixed, 16 remaining.

**Section 8 is preserved verbatim**, contradictions and all. Its Batch C phase/engine
wording, and its "pending decisions" wording for §1.5, §1.11, §1.12 and §2.1, are
**superseded**: all were settled, implemented and marked fixed. Everything else §8 retains
as later scope remains **binding**. Do not edit §8.

Batches A/D, B, C, E and F are closed and merged to `master`. Batch F went in as PR #48,
`f78caeb`. Do not revisit or re-verify closed batches.

## 1. Session goal

Plan, implement and verify Batch G from review §9: §3.1, where the zone editor's
status line rebuilt its graph diagnostics on every frame. The owner reviews everything at
the end.

## 2. Fixes applied

All are uncommitted, pending owner review.

- **§3.1, the graph summary cache.**
  [zoneEditorGraphState.go](../app/gui/dialogs/zoneEditorGraphState.go) caches the
  `ZoneEditorGraphDto` behind its own `graphDirty` flag, which is never shared with
  `geometryDirty`.
  - All seven structural mutators call `markGraphDirty`: `setEditingSet`,
    `addConnection`, `deleteConnection`, `undoSessionEdits`, `addZoneAt`, `deleteZone`
    and `applyQualityMutation`.
  - `layoutStatus` in [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go)
    rebuilds only when the flag is set, still through a fresh `derefConnections` copy.
  - Idle frames make **0** handler calls, and each structural edit makes exactly **1**.
- **The one-frame-late status line** (a pre-existing issue the owner chose to fix). The
  status line is drawn before canvas input and side-panel writeback. The key in
  [zoneEditorStatusKey.go](../app/gui/dialogs/zoneEditorStatusKey.go) captures the hint,
  both add modes, the zone and connection counts, and the dirty flag.
  `requestLateStatusRedraw` compares it at the end of `Body` and issues
  `op.InvalidateCmd` when it changed. The key is stored unconditionally after any rebuild,
  so redraws always settle, even in the error branch.

## 3. Features added / changed

### Batch G, uncommitted (decisions G1–G9a in the plan)

- **Measure first, then fix.** There is no percentage gate and no global allocation
  threshold. Acceptance is the deterministic call-count criterion (G2).
- **The O(zones + connections) `FindIsolatedZones` rewrite was measured and rejected**
  under G3. It was faster only at 24 zones (−19%) and 40 zones (−33%). It was slower at
  4 zones (+292%) and 12 zones (+67%), and added 3 allocations per call. The nested loop
  is byte-identical to HEAD.
- **New untagged benchmarks** in
  [zoneEditorGraph_test.go](../test/performance/zoneEditorGraph_test.go) run on fixed
  4/12/24/40-zone graphs (G9a, deterministic because generation is unseeded): the
  handler, the service, and a GPU-free dialog idle frame.
- **Idle-frame result:** B/op fell from 4457–18478 to 3049–3073, now flat across graph
  sizes, and allocations fell by 1–4 per frame. Full medians are in the plan.

### Delivered in Batch F, committed and settled

The shared curve builder (`preview.ConnectionCurveLayout.Build`, 21px spacing, obstacle
bending everywhere), effective-portal classification, the type dropdown and live-list hit
tests are committed, merged, and must not be reopened. The Batch F plan file is no longer
in the tree; the merged code is the reference.

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
- Batch G (once committed): the graph cache keeps its own dirty flag. Never fold it into
  `geometryDirty`, which hit tests clear mid-input. The status key must stay built from
  the full state, and never from only the fields the branch drew.

## 4. File modifications

Everything below is unstaged on `AD/performance`.

**Production (new):**
- [zoneEditorGraphState.go](../app/gui/dialogs/zoneEditorGraphState.go): the cached graph,
  `graphDirty`, `shownStatus` and `markGraphDirty`.
- [zoneEditorStatusKey.go](../app/gui/dialogs/zoneEditorStatusKey.go): the comparable
  status key.

**Production (edited):**
- [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go): embeds the state, adds
  the cached `layoutStatus`, `statusKey`, `requestLateStatusRedraw`, and 6 mutator marks.
- [zoneEditorZoneProps.go](../app/gui/dialogs/zoneEditorZoneProps.go):
  `applyQualityMutation` marks the graph (gofmt realigned one trailing comment).

**Test helpers:**
- [appRunner.go](../test/test_helpers/integration_common/appRunner.go): a private
  `newAppRunner` shared by `NewAppRunnerWithFileSystem` and the new
  `NewAppRunnerWithGuiHandler`, plus `ImmediateRedrawRequested`.
- [graphDescriptionCounter.go](../test/test_helpers/integration_common/graphDescriptionCounter.go)
  (new, `integration_test`).

**Tests:** listed in §5. No goldens changed.

**Agent docs:**
- [the Batch G plan](plans/batch-g-editor-graph-diagnostics.md): the owner committed the
  first version as `a587617`. The working copy adds the approvals, review records and phase
  summaries.
- [test_observations.md](backlog/test_observations.md): records the GUI-only statements of
  the graph cache.
- `.agent/memories/gui-and-tests.md` (gitignored): redraw assertions, handler wrappers,
  and dialog-direct limits.
- This handoff.

**Not touched:** `data/`, `internal/entities/template_entity/`, `internal/registry/`, the
output path, generated Wire, `connectionEditorService.go` (reverted byte-identical),
dependencies, the review document and `.agent/backlog/owner_findings.md`.

## 5. Tests added or updated

**Unit tests** (mirrored layout, triple-A, `t.Parallel()`):
- `connectionEditorService/findIsolatedZones_test.go`: 6 new tests. They lock the `nil`
  result, input order, `To`-only references, unknown endpoints, duplicate names and
  self-loops.
- `zoneEditorHandler/describeZoneEditorGraph_test.go`: 1 new test for a `nil` isolated
  result.

**GUI integration** (`integration_test && gui`):
[zoneEditorGraphCache_integration_test.go](../test/integration/gui/zoneEditorGraphCache_integration_test.go),
with 20 tests:
- Call counts:
  - opening the editor makes 1 call, and idle frames make 0;
  - 8 structural edits make 1 each;
  - 4 non-structural edits make 0 each.
- 2 isolation correctness tests.
- 3 window-harness redraw tests: the positive case, an empty-canvas negative control, and
  convergence.
- 1 dialog-direct error-state convergence test.

Each group was shown to fail under a temporary production mutation, reverted afterwards.

**Baseline** (unchanged code, HEAD `a587617`):
- build: PASS;
- coverage: **74.6%**;
- `go test ./test/...`: PASS;
- tagged run: PASS, root 6.453s, GUI 47.423s (machine load);
- `testlayoutcheck`: PASS; `gofmt`: clean;
- lint: 0 issues.

**Final** (Windows/amd64, Go 1.27.0, empty `GOFLAGS`):
- build: PASS;
- `-p=2` coverage: PASS, **74.5%**. Covered statements are unchanged at 6911, and the
  denominator grew 9246 → 9261 with 15 GUI-only dialog statements. **This needs owner
  acceptance.**
- `go test ./test/...`: PASS;
- tagged run: PASS, root 3.322s, GUI 32.380s;
- the new GUI tests at `-count=5`: PASS;
- `testlayoutcheck`: PASS; `gofmt -l`: clean;
- lint: **0 issues**, after fixing 4 new ones: `funcorder` ×2, `gochecknoglobals` and
  `funlen`.

**Reviews** (Claude Opus 5.5):
- Plan: APPROVE WITH CHANGES. All 12 findings were applied, and G9a was owner-approved.
- Implementation: APPROVE WITH CHANGES. The 1 minor finding and 3 of the 4 nits were
  applied. Nit 5 was declined with a reason, and one out-of-scope observation was
  recorded (see §7).

**Not measured:** native Linux, Steam Deck, the race detector (cgo is off locally), and
in-game behaviour.

## 6. Git status snapshot

The branch is **`AD/performance`**, HEAD **`a587617` ("plan")**, level with
`origin/AD/performance`, on top of `master` `f78caeb` (Batch F, #48). `git status --short`:

```text
 M .agent/backlog/test_observations.md
 M .agent/plans/batch-g-editor-graph-diagnostics.md
 M .agent/session-carry-forward.md
 M app/gui/dialogs/zoneEditorDialog.go
 M app/gui/dialogs/zoneEditorZoneProps.go
 M test/test_helpers/integration_common/appRunner.go
 M test/unit/internal/handlers/zoneEditorHandler/describeZoneEditorGraph_test.go
 M test/unit/internal/services/connection_editor/connectionEditorService/findIsolatedZones_test.go
?? app/gui/dialogs/zoneEditorGraphState.go
?? app/gui/dialogs/zoneEditorStatusKey.go
?? test/integration/gui/zoneEditorGraphCache_integration_test.go
?? test/performance/zoneEditorGraph_test.go
?? test/test_helpers/integration_common/graphDescriptionCounter.go
```

`.agent/memories/` and `tmp/` (the raw benchmark output) are gitignored.

**The assistant performed no staging, unstaging, commit, push, stash, branch switch or
worktree change.**

## 7. Rejections / things the user declined

- **Batch G.**
  - The O(z+c) isolation rewrite failed G3 and was reverted, with the numbers recorded.
  - Dialog-direct invalidate tests were planned but proved infeasible: the canvas offset
    depends on the toolbar height. The window harness was used instead, with a settle loop
    and a negative control.
  - Implementation review nit 5 (the open-editor test's Act only reads the counter) was
    declined, because opening the editor is the action.
  - Recorded, not assigned: the toolbar's "Delete selected" enabled state is still one
    frame late after a canvas selection, as at HEAD, because `hasSelection` is not in the
    status key. It is out of §3.1 scope.

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
  CRLF. Measured on 2026-09-28 before and after the Batch G rewrite, the LF-normalized
  UTF-8 SHA-256 of this file's §8 (from `## 8.` up to `## 9.`) is
  `0DDDEBF32DB51167E675EA7A147E75F186B77643E8024F26E7B262102B004AD8`. The earlier recorded
  value `35D5F4…` did not reproduce with this method for either document.
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

**Batch G is implemented and waiting on the owner.** Nothing may be reopened or reworked
without owner direction. In order:

1. **The owner reviews** the working tree (§6) and decides whether to accept coverage
   **74.5%**: covered statements are unchanged, and 15 GUI-only statements were added.
2. **Owner feedback**, if any: apply it within §3.1 scope only, rerun the affected checks
   (the plan's Phase 6 list), and update the plan's summaries.
3. **After the owner commits** (a close-out turn):
   - mark §3.1 `✅ FIXED` in place in the review with a Progress paragraph naming the
     commit;
   - update the §9 `G` row to complete;
   - set the progress line to **17 fixed, 15 remaining** (0 High, 9 Medium, 6 Low);
   - mark the plan CLOSED and refresh this handoff.

   Do not renumber, and do not edit this handoff's §8.
4. **Then the next batch from review §9**, which is the owner's choice:
   - H: CI/tooling, §6.2/§6.3/§6.5;
   - I: docs, §7.1/§7.2;
   - J: the §2.2 service boundary;
   - and so on.

   Follow the same gates: read, inspect, ask, summarize, plan, independent review, plan
   approval, baseline.

**Deployment.** Nothing is deployed and nothing is authorized to be. The owner alone
stages, commits, merges and releases. No schema migration, dependency installation, Wire
regeneration or output-directory change is required or pending. Native Linux and Steam
Deck execution remain unmeasured.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md). It
> is self-contained.
>
> **Batch G (review §3.1) is IMPLEMENTED and VERIFIED on Windows, and NOT committed.**
> Every change is unstaged on branch `AD/performance` (HEAD `a587617`, on `master`
> `f78caeb`). The durable record is
> [the Batch G plan](plans/batch-g-editor-graph-diagnostics.md), with owner decisions
> G1–G9a, phase summaries, benchmarks and reviews. §3.1 is **not** yet marked FIXED:
> that happens only after the owner commits.
>
> **What changed:**
> - The zone editor caches its status-line graph diagnostics behind its own `graphDirty`
>   flag, set by 7 structural mutators and never shared with `geometryDirty`.
> - Any status change made after the toolbar was drawn requests an immediate redraw
>   (`op.InvalidateCmd`), through a status key built from the full state.
> - The O(z+c) isolation rewrite was measured and rejected under G3.
>
> **Verification on record** (Windows/amd64, Go 1.27.0, empty `GOFLAGS`):
> - build: PASS;
> - `-p=2` coverage: 74.5% (covered statements unchanged; needs owner acceptance);
> - `go test ./test/...`: PASS;
> - tagged `integration_test,gui` run: PASS;
> - `testlayoutcheck`: PASS; `gofmt`: clean; lint: 0 issues.
>
> **Native Linux, Steam Deck and the race detector are UNAVAILABLE.** No in-game result is
> claimed.
>
> **Next:** wait for owner review. Apply feedback within §3.1 only. After the owner
> commits, close out as in §9 step 3.
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
> - hit tests resolving against the live connection list;
> - the graph cache's own dirty flag, and a status key built from the full state.
>
> This handoff contains the full continuation context.
