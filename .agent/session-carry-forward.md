# Carry-forward: Batch G closed, next batch not started

Date: 2026-09-28. **Batch G (review §3.1, measured performance) is closed.** The owner
reviewed it, committed it as **`35e0fab` ("Batch G")** on `AD/performance`, and authorized
the close-out with *"Reviewed, you can mark the item as done (finish Phase 6) and update
the carry forward"*. The committed code matches the reviewed working tree.

**§3.1 is marked FIXED** in the
[surviving review](backlog/review-gpt-6-astra-09-07.md). Its §9 row G is complete, and
its progress line reads **17 fixed, 15 remaining** (0 High, 9 Medium, 6 Low).

The plan, [batch-g-editor-graph-diagnostics.md](plans/batch-g-editor-graph-diagnostics.md),
is marked CLOSED. **The owner will delete it once `AD/performance` reaches `master`**, so
this handoff restates everything the plan settled (§3). Do not recreate the plan, and do
not reimplement, re-verify or re-review Batch G.

**Section 8 is preserved verbatim**, contradictions and all. Its Batch C phase/engine
wording, and its "pending decisions" wording for §1.5, §1.11, §1.12 and §2.1, are
**superseded**: all were settled, implemented and marked fixed. Everything else §8 retains
as later scope remains **binding**. Do not edit §8.

Batches A/D, B, C, E, F and G are closed. A/D through F are merged to `master` (Batch F as
PR #48, `f78caeb`). Batch G is pushed on `AD/performance` and is not yet merged. Do not
revisit or re-verify closed batches.

## 1. Session goal

Plan, implement, verify and close Batch G from review §9: §3.1, where the zone editor's
status line rebuilt its graph diagnostics on every frame. This closing turn marked the
finding fixed and handed this file forward.

## 2. Fixes applied

All are committed in `35e0fab` and **marked FIXED**. They are closed; do not reopen them.

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

### Delivered in Batch G, committed and settled

The owner decisions, restated here because the plan will be deleted:

| ID | Decision |
| --- | --- |
| G1 | Measure first, then fix. |
| G2 | Acceptance is deterministic: idle frames make 0 `DescribeZoneEditorGraph` calls, and each structural edit makes exactly 1. Benchmarks are recorded with no percentage gate and no global allocation threshold. |
| G3 | The O(z+c) isolation rewrite is adopted only if it is faster at 24 and 40 zones and within ±5% at 4 and 12. |
| G4 | The cache lives in the dialog with its own dirty flag. The handler stays stateless. |
| G5 | Handler calls are counted by a test-helper wrapper embedding `IGuiHandler`, with no new test exports. |
| G6 | There is an untagged, GPU-free idle-frame benchmark. |
| G7/G8 | Late status changes (hint, both add modes, counts, diagnostics) trigger an immediate `op.InvalidateCmd`. The layout order is unchanged. |
| G9/G9a | Benchmarks run at 4/12/24/40 zones on **deterministic** fixtures, because generation uses the unseeded global `rand`. |

- **The O(zones + connections) `FindIsolatedZones` rewrite was measured and rejected**
  under G3. It was faster only at 24 zones (−19%) and 40 zones (−33%). It was slower at
  4 zones (+292%) and 12 zones (+67%), and added 3 allocations per call. The nested loop
  is unchanged. Do not re-propose it without new measurements.
- **New untagged benchmarks** in
  [zoneEditorGraph_test.go](../test/performance/zoneEditorGraph_test.go) run on fixed
  4/12/24/40-zone graphs (G9a, deterministic because generation is unseeded): the
  handler, the service, and a GPU-free dialog idle frame.
- **Idle-frame result:** B/op fell from 4457–18478 to 3049–3073, now flat across graph
  sizes, and allocations fell by 1–4 per frame. ns/op at 40 zones fell from 215030 to
  190195. The other ns/op changes were within run-to-run noise: unchanged code varied up
  to 9%.

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
- All Batch G behavior is committed and settled. The graph cache keeps its own dirty flag:
  never fold it into `geometryDirty`, which hit tests clear mid-input. Any new edit that
  replaces the zone or connection lists must call `markGraphDirty`. The status key must
  stay built from the full state, and never from only the fields the branch drew,
  otherwise the redraw never settles.

## 4. File modifications

**This closing turn is documentation only, in three files:**
- [Surviving review](backlog/review-gpt-6-astra-09-07.md): §3.1 marked `✅ FIXED` in place
  with a Progress paragraph, the §9 row G completed, and the progress line refreshed to
  17 fixed / 15 remaining. Nothing was renumbered, and the review's §8 hash is unchanged.
- [Batch G plan](plans/batch-g-editor-graph-diagnostics.md): Phase 6, Final Recap and
  Deployment Plan marked complete/CLOSED.
- This handoff.

**Batch G inventory, committed in `35e0fab`:**

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
- [the Batch G plan](plans/batch-g-editor-graph-diagnostics.md) (first version `a587617`,
  completed in `35e0fab`, closed in this turn; the owner will delete it after the merge).
- [test_observations.md](backlog/test_observations.md): records the GUI-only statements of
  the graph cache.
- `.agent/memories/gui-and-tests.md` (gitignored): redraw assertions, handler wrappers,
  and dialog-direct limits.
- This handoff.

**Not touched in Batch G:** `data/`, `internal/entities/template_entity/`,
`internal/registry/`, the output path, generated Wire, `connectionEditorService.go`
(reverted byte-identical), dependencies and `.agent/backlog/owner_findings.md`.

## 5. Tests added or updated

**This closing turn added no tests and reran nothing**, because it changed documentation
only. The Batch G verification below is the evidence of record for `35e0fab`.

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
  denominator grew 9246 → 9261 with 15 GUI-only dialog statements. **Accepted by the
  owner with the close-out. 74.5% is the new baseline for the next batch.**
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

The branch is **`AD/performance`**, HEAD **`35e0fab` ("Batch G")**, level with
`origin/AD/performance`. The history on top of `master` (`f78caeb`, Batch F #48) is
`a587617` ("plan"), then `35e0fab`. Batch G is not yet merged to `master`.

The tree was clean when this closing turn began. `git status --short` now shows only this
turn's three documentation edits, all unstaged:

```text
 M .agent/backlog/review-gpt-6-astra-09-07.md
 M .agent/plans/batch-g-editor-graph-diagnostics.md
 M .agent/session-carry-forward.md
```

`.agent/memories/` and `tmp/` (the raw benchmark output) are gitignored.

**The assistant performed no staging, unstaging, commit, push, stash, branch switch or
worktree change in this turn.**

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

**Batch G is closed, with no open Batch G questions and no blockers.** The owner merges
`AD/performance` to `master`, then deletes the Batch G plan. Neither is an agent action.

The next unit of work is the **owner's choice** among the remaining batches of the
surviving review's §9 table (15 items left):
- **H: CI/tooling hardening,** §6.2, §6.3 and §6.5. This is next in table order, and the
  items are independent configuration changes:
  - linter version alignment;
  - LF attributes for module and checksum files, with owner-approved normalization of
    those four paths only;
  - release-tag input validation.

  Each has open owner decisions recorded in its review item.
- **I: docs,** §7.1 and §7.2.
- **J: the reopened service boundary,** §2.2.
- **K–P:** topology retirement, the compact-state investigation, the persistence format,
  panel/state organization, naming and lookup, and the Vec2 audit. The binding scope for
  each is in §8's retained decisions.

Routing for the next session, in order:
1. Read [AGENTS.md](../AGENTS.md), this handoff, then the chosen items and their §9 row.
2. Inspect the current code yourself, because review line numbers predate Batches F and G.
3. Ask the owner every open decision in the items, summarize the scope, and get approval.
4. Write a new durable plan under `.agent/plans/`, get an independent review and explicit
   plan approval, and capture a fresh baseline (coverage 74.5%) before the first edit.

**Deployment.** Nothing is deployed and nothing is authorized to be. The owner alone
stages, commits, merges and releases. Batch G still awaits the owner's merge to `master`.
No schema migration, dependency installation, Wire regeneration or output-directory
change is required or pending. Native Linux and Steam Deck execution remain unmeasured.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md). It
> is self-contained.
>
> **Batch G is CLOSED.** The owner reviewed it and committed it as `35e0fab` on
> `AD/performance`, which is not yet merged to `master` (`f78caeb`). Review §3.1 is
> marked **FIXED** in the [surviving review](backlog/review-gpt-6-astra-09-07.md), whose
> §9 table shows batch `G` complete, with 17 fixed and 15 remaining. Coverage of 74.5%
> was accepted and is the new baseline. The Batch G plan is CLOSED and will be deleted by
> the owner after the merge; §3 of this handoff restates its settled decisions (G1–G9a).
> Do not reimplement, re-verify or re-review Batch G.
>
> **State:** the only working-tree changes are the closing turn's three documentation
> edits: the review, the plan and this handoff. The assistant performed no Git mutation;
> preserve the owner's state exactly.
>
> **Next work is the owner's choice of batch from review §9.** H (CI/tooling: §6.2,
> §6.3, §6.5) is next in order. No plan exists yet. Follow the gates in §9: read,
> inspect, ask, summarize, plan, independent review, plan approval, fresh baseline.
>
> **Verification on record** (Batch G final code, Windows/amd64, Go 1.27.0, empty
> `GOFLAGS`; not rerun at close-out):
> - build: PASS;
> - `-p=2` coverage: PASS, 74.5%;
> - `go test ./test/...`: PASS;
> - tagged `integration_test,gui` run: PASS, root 3.322s, GUI 32.380s;
> - `testlayoutcheck`: PASS; `gofmt`: clean; lint: 0 issues.
>
> **Native Linux, Steam Deck and the race detector are UNAVAILABLE.** No in-game result is
> claimed.
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
> **Out of scope** unless the owner selects that batch: Batch K topology retirement,
> direct `GeneratorConfig` rejection, DTO cleanup, schema and package work, allocation
> tuning, and every settled alternative in §7. §8 stays verbatim. Its Batch C wording
> and its pending-decision wording for §1.5, §1.11, §1.12 and §2.1 are superseded; its
> retained later-scope decisions remain binding.
>
> **Preserve:**
> - explicit Portal Road `false`/`nil` and valid approaches;
> - the tri-state road display (`nil` is roaded only for explicit Portals, and effective
>   classification never feeds roads);
> - independent internal roads, nil-state content/road preservation, and source cloning;
> - the single reusable half-opacity PNG edge mask, and the Preview-only legend;
> - the shared curve builder, as the single source of curves;
> - hit tests resolving against the live connection list;
> - the graph cache's own dirty flag (`markGraphDirty` on every edit that replaces the
>   zone or connection lists), and a status key built from the full state.
>
> This handoff contains the full continuation context.
