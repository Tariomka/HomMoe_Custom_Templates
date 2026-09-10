# Carry-forward: Batch C approved, begin Phase 1

Date: 2026-09-10.

**Next action:** Implement Phase 1 of [the active Batch C plan](plans/batch-c-road-policy-and-visuals.md).
Owner approved implementation on 2026-09-10: “looks good, please proceed”. No code has
changed; the planning session reached its length limit. Do not ask for approval again.
A/D and B are complete; no validation, closure or retired-document lookup is required.

## 1. Session goal

Investigate review §1.4/§1.6/§1.10, resolve owner road policy, measure a fresh baseline,
and prepare an approved-scope plan. Owner expanded scope to distinguish roaded and
roadless connections in the Preview panel, manual editor and exported PNG.

## 2. Fixes applied

None this session. Source verification confirms all three Batch C findings remain.
Prior settled fixes stay intact: manual compare-before-mutation/dirty tracking and
snapshot isolation; no output browsing-fallback authorization; additive PNG ceiling
sample count with existing stamp spacing; configured tournament SaveArmy (nil rules
plus tournament selector means false; normal constructor default remains true).

## 3. Features added / changed

No implemented changes. Owner approved these Batch C requirements:

- Checkbox means roads **between zones**. With EditorState present, set Road true/on or false/off on every non-explicit-Portal connection, including custom/imported, Default, empty, arena and proximity. On overwrites imported false; off overwrites true. Off removes whole zone-road segments referencing these connections, not edges.
- Internal castle/object/foothold roads stay generated even off. Preserve explicit Portal flags and valid approaches, including Road=false portals; generated portals keep true. Portal placement rules alone do NOT grant road-policy exemption.
- Reconcile foothold additions/removals/count changes with final mandatory content and actual main objects. Remove only confirmed invalid MainObject, incident Connection and named MandatoryContent references; preserve valid custom internal roads and opaque types. No prefix-only deletion, no new provenance metadata.
- Nil EditorState preserves supplied roads/flags and existing mandatory content, with direct EnsureConnectionNames and normal graph checks but NO road cleanup.
- Repair graph connectivity from edges, including accumulated repairs. Preserve impossible-isolation fallback independently of road materialization.
- All three visual surfaces: explicit Road false is roadless, true roaded, nil roaded ONLY for explicit Portal. Owner will test this interpretation in-game.
- GUI roadless non-portals #B0B0B0; roadless explicit portals #90EE90. Keep roaded colors. Selected editor edges retain road color (even roaded) and use thicker width. Legends: Road / No road / Portal / Portal without road.
- PNG roadless strokes: existing color at 50% once per edge; stamps within one edge do not accumulate, intersections of separate edges may. Preserve geometry/dashes, clipping and opaque markers. Reusable mask, no per-edge full-image allocations.
- Pending editor create/type changes show the checkbox policy before Apply via handler/service logic. Cancel stays non-mutating. No per-edge checkbox requested.
- Road display classification is separate from existing effective preview Type. Keep placement-rule-derived portal shapes and existing type comparison behavior; road policy uses explicit Portal only. §1.15 unification remains deferred.

## 4. File modifications

- Created [active Batch C plan](plans/batch-c-road-policy-and-visuals.md): approved requirements, concrete boundaries, phased tests and independent review corrections.
- Updated [this handoff](session-carry-forward.md) to record implementation approval and the Phase 1 starting point.
- Updated [settled decisions](memories/settled-decisions.md) with durable owner road/visual policy.
- Refreshed ignored [coverage.txt](../coverage.txt), [coverage.html](../coverage.html) and [lcov.info](../lcov.info) with the fresh baseline task.
- No production, tests, configuration, protected trees, Wire or goldens edited.

## 5. Tests added or updated

None added. Fresh Windows Go 1.27.0 full unit coverage run (`-count=1`) PASS;
Go `cover -func` reports **74.5%**. Report-only lint: **0 issues**, three known
unused-exclusion warnings. GOFLAGS empty. Existing tests lack the newly agreed
road-state contracts despite high statement coverage.

No build/default/tagged/GUI/Linux/race/benchmark run this session. Last recorded
`go test ./test/...` outcome was PASS in the completed previous session, not a fresh
Batch C run. No implementation result or in-game correctness is claimed.

Claude Opus 5 reviewed design and written plan. Written-plan verdict: approve with
corrections. Incorporated both blockers (whole-template golden flag updates without
weakening equality; direct naming on nil-state path and revised old mock expectation)
and clarified new type-change handler operation and unchanged shape classification.

## 6. Git status snapshot

Baseline was measured on `master` at `5b1feb0eb315a4df54ca1ea9e1defd070f6061d8`.
Latest branch `AD/road_and_graph_invariants`, HEAD `5240389aa0a8eb6a0a04d0552369bc2437540e60`.
Owner changed branch and committed plan/handoff; the agent did neither. Read-only
diff from the baseline confirms only those two documentation files changed. Working
tree and staging were clean before recording approval. Final `git status --short`
reports only this handoff modified; staged diff is empty. Approval was saved in the
plan, handoff and memory. Reinspect current Git rather than assuming their status.
Inspect again before implementation; preserve all owner staging/changes. No agent
staging, unstaging, commits, pushes, stashes or speculative branch changes.

## 7. Rejections / things the user declined

- Original zero-roads-anywhere interpretation rejected: internal roads and portal approaches stay; only non-Portal border roads follow the checkbox.
- Do not infer display state from zone-road references or default every nil to roaded.
- Scope was initially withheld pending visual expansion, then approved in full.
- Reject review suggestions contrary to owner: nil-state cleanup, orange selected roaded edges, deleting portal approaches for false flags, adding arena-marker castle roads, and relying solely on final generator cleanup for public repair correctness.
- No protected edits, global tags, bulk rewrites, fake seams or Wire hand edits.
- Keep PathResolutionService, pending-base consumption on every Apply and unpersisted session-only export selection. Browsing fallbacks never authorize export.
- Keep unconditional in-memory PNG assertions; no environment-gated capture-only utility. Keep original floating-point stamp spacing, no resampling disguise.
- Two exploration calls rejected lowercase model IDs; retried successfully using exact `GPT-5.6 Terra (copilot)` display name (already documented in tooling memory).
- Completed plans are owner-deleted. Future handoffs must not link them or delegate their validation/closure; surviving decisions must be self-contained.

## 8. Open questions and confirmed later scope

Batch C source verification, baseline, policy questions, scope approval and written-plan
implementation approval are done. **Begin Phase 1; do not ask approval or settled policy again.**
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

1. Read AGENTS, this handoff and [active plan](plans/batch-c-road-policy-and-visuals.md).
2. Inspect current Git and preserve owner changes. Approval is already granted.
3. If implementation sources changed since the recorded baseline, measure a new C baseline and reconcile current sources only. The latest owner commit was documentation-only. Do not repeat completed-batch audits or settled questions.
4. Begin Phase 1 classifier/graph regressions, then implement/verify the plan one phase at a time; update its checkboxes and summaries.
5. Keep §8 later scope intact. In-game checks belong to the owner, not a claimed agent result.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md)
> and [active Batch C plan](plans/batch-c-road-policy-and-visuals.md). Scope AND written-plan
> implementation are approved (owner, 2026-09-10: “looks good, please proceed”). Begin
> Phase 1 without asking for approval again; no implementation has started yet.
> Batch C covers §1.4/§1.6/§1.10 plus all three road-state visual surfaces. Follow the
> settled owner contract in the handoff/plan; do not reopen its policy questions.
> Fresh Windows baseline at master 5b1feb0: full unit coverage PASS, 74.5%, lint zero.
> Inspect Git to preserve owner changes. A/D and B are closed; no old-document lookup,
> validation or closure work. Preserve all later scope in §8.
>
> Never modify protected data/schema/registry trees; proposed protected changes require owner approval and application. Preserve Windows/Linux portable paths and guarded platform code. Test nontrivial changes and measure coverage before/after. Keep multi-step work in an approved durable plan. Never stage, unstage, commit or push; preserve owner staging and do not switch branches speculatively. Never bulk-rewrite or hand-edit generated Wire output. Export only to the detected game templates directory or explicit session-only picker destination; never persist output paths or authorize browsing fallbacks. Never set global integration_test/gui/wireinject tags or add fake unit seams. Preserve later scope in §8.
