# Carry-forward: Batch C Phase 1 complete, begin Phase 2

Date: 2026-09-10.

**Next action:** Implement Phase 2 of [the active Batch C plan](plans/batch-c-road-policy-and-visuals.md).
Owner approved implementation on 2026-09-10: “looks good, please proceed”. No code has
been committed by the agent; Phase 1 is implemented, verified and still unstaged.
Do not ask for approval again.
Owner accepted Phase 1 and explicitly requested Phase 2 on 2026-09-10:
“All looks good, you can proceed with the next Phase”. Phase 2 has not started:
the current session reached the AGENTS.md session limit. Resume directly in a
fresh session. The editor reports intervening source/test edits; re-read current
files and inspect Git before changes rather than assuming the prior verified state.
A/D and B are complete; no validation, closure or retired-document lookup is required.

## 1. Session goal

Implement approved Batch C Phase 1: shared road classifier and endpoint-based,
road-gated public graph repair. Full Batch C includes the later road reconciliation
and three visual surfaces, which remain unimplemented.

## 2. Fixes applied

Implemented Phase 1:
- [Connection model](../internal/models/template_model/template_variant_model/connection.go):
	`IsExplicitPortal()` and `HasRoad()`, explicit true/false and nil classifier.
- [Repair service](../internal/services/template_generator/providers/topology/base/topologyConnectionService.go):
	valid incident endpoints replace road references for player connectivity;
	accumulated repairs update both endpoints. Missing/self endpoints do not count.
- Repair methods now take a final `generateRoads bool`, threaded through their
	interface, TopologyBase, positioned, Chain, Ring, SharedWeb and balanced cluster
	callers. Road creation is conditional; fallback/bridge edges and guards remain.
- Repair-name collisions choose first free `-2`, `-3`, etc. suffix. Removed the
	old bridge collision behavior that merged adjacency without a real edge.

Prior settled fixes stay intact: manual compare-before-mutation/dirty tracking and
snapshot isolation; no output browsing-fallback authorization; additive PNG ceiling
sample count with existing stamp spacing; configured tournament SaveArmy (nil rules
plus tournament selector means false; normal constructor default remains true).

## 3. Features added / changed

Phase 1 changes are above. The following approved full Batch C requirements remain
the source of truth; do not mistake these later policy/visual items for implemented work:

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

- [connection.go](../internal/models/template_model/template_variant_model/connection.go): classifier.
- [topologyConnectionService.go](../internal/services/template_generator/providers/topology/base/topologyConnectionService.go): endpoint connectivity, road gate, collision-safe repair names.
- [topologyConnectionServiceInterface.go](../internal/services/template_generator/providers/topology/base/topologyConnectionServiceInterface.go): final boolean arguments.
- [topologyBase.go](../internal/services/template_generator/providers/topology/base/topologyBase.go): forward boolean arguments.
- [chainTopology.go](../internal/services/template_generator/providers/topology/chainTopology.go), [ringTopology.go](../internal/services/template_generator/providers/topology/ringTopology.go), [webTopology.go](../internal/services/template_generator/providers/topology/webTopology.go), [positionedTopologyBuilder.go](../internal/services/template_generator/providers/topology/positionedTopologyBuilder.go), [balancedClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/balancedClusterService.go): pass configuration road setting.
- Edited base test files: [createMissingConnections_test.go](../test/unit/internal/services/template_generator/providers/topology/base/topologyBase/createMissingConnections_test.go), [createMissingPlayerConnections_test.go](../test/unit/internal/services/template_generator/providers/topology/base/topologyBase/createMissingPlayerConnections_test.go).
- Edited topology regressions: [buildVariant_test.go](../test/unit/internal/services/template_generator/providers/topology/positionedTopologyBuilder/buildVariant_test.go), [createTopologyVariant_test.go](../test/unit/internal/services/template_generator/providers/topology/randomTopology/createTopologyVariant_test.go).
- New connection-model test files under [connection/](../test/unit/internal/models/template_model/template_variant_model/connection/): common, hasRoad, isExplicitPortal, clone, toConnectionModel, toConnectionEntity (all `_test.go`).
- New direct service test files under [topologyConnectionService/](../test/unit/internal/services/template_generator/providers/topology/base/topologyConnectionService/): common, createMissingConnections, createMissingPlayerConnections, createRandomPortalConnections, getBorderGuardValue (all `_test.go`).
- Updated [active plan](plans/batch-c-road-policy-and-visuals.md), [this handoff](session-carry-forward.md), and ignored [tooling memory](memories/tooling-and-shell.md).
- Refreshed ignored [coverage.txt](../coverage.txt), [coverage.html](../coverage.html), [lcov.info](../lcov.info).
- No configuration, protected trees, Wire, goldens or snapshots edited.

## 5. Tests added or updated

Added classifier matrix, adjacent cloning/mapping coverage, direct-service tests and
base/topology regressions for roadless valid graphs, repairs accumulated within/across
calls, real Random impossible isolation, invalid/self endpoints and name collisions.
Separate edge and road assertions prevent vacuous flattened expectations.

Windows verification PASS: build; affected generator/model and architecture suites;
fresh before/after full unit coverage tasks (`-count=1`), final cache-eligible full
coverage rerun after test-only additions; `go test ./test/...` including untagged
integration; test-layout checker; diff whitespace check. Coverage **74.5% → 74.6%**.
Classifier and changed/new repair functions **100%**. Report-only lint **0 issues**,
three existing unused-exclusion warnings. No gated integration/GUI, Linux, race,
benchmark or in-game validation this session. Independent Claude Opus 5 Phase 1
implementation review approved after assertion and collision fixes.

## 6. Git status snapshot

Branch `AD/road_and_graph_invariants`, HEAD `87ae8789d84178854a395e8dd34e172095a9ebc6`.
Started clean; current source edits in §4 are `M`, new test files are `??`, plan and
handoff modified. Index empty. Protected trees and generated Wire unchanged.
No agent staging, unstaging, commits, pushes, stashes or branch changes. Reinspect
current Git and preserve owner changes before continuing.

## 7. Rejections / things the user declined

- Phase 1 review rejected vacuous nil-zone endpoint tests and flattened edge/road
	expectations; corrected to independent real-graph and road assertions.
- Did not retain the old fake-adjacency collision fallback: first-free suffixes
	satisfy the approved actual-edge invariant without overwriting supplied edges.
- A guessed connection-model path did not exist; located the nested implementation.
	A promoted-field fixture lint edit produced invalid Go syntax; corrected to the
	actual Go 1.27 flat literal and updated tooling memory. One stale patch context
	was re-read before continuing. Final build/tests/lint are clean.
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
implementation approval are done. **Phase 1 is complete; begin Phase 2 without asking approval or settled policy again.**
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
3. Current Phase 1 verified coverage is 74.6%; preserve the unstaged/untracked implementation. Reconcile any owner edits since this handoff, not old completed-batch audits.
4. Begin Phase 2 road-policy service and final-content reconciliation. New classifier methods are ready for use; repair signatures already accept the road setting. Follow the plan one phase at a time, with verification and summaries.
5. Keep §8 later scope intact. In-game checks belong to the owner, not a claimed agent result.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md)
> and [active Batch C plan](plans/batch-c-road-policy-and-visuals.md). Scope AND written-plan
> implementation are approved (owner, 2026-09-10: “looks good, please proceed”). Begin
> Phase 2 without asking for approval again; Phase 1 is implemented and verified,
> but unstaged/untracked. Preserve it.
> Batch C covers §1.4/§1.6/§1.10 plus all three road-state visual surfaces. Follow the
> settled owner contract in the handoff/plan; do not reopen its policy questions.
> Phase 1 Windows verification: build/default tests/architecture/layout PASS, coverage
> 74.5% → 74.6%, lint zero. Classifier and repair methods 100%. Opus review approved.
> Inspect Git to preserve owner changes. A/D and B are closed; no old-document lookup,
> validation or closure work. Preserve all later scope in §8.
>
> Never modify protected data/schema/registry trees; proposed protected changes require owner approval and application. Preserve Windows/Linux portable paths and guarded platform code. Test nontrivial changes and measure coverage before/after. Keep multi-step work in an approved durable plan. Never stage, unstage, commit or push; preserve owner staging and do not switch branches speculatively. Never bulk-rewrite or hand-edit generated Wire output. Export only to the detected game templates directory or explicit session-only picker destination; never persist output paths or authorize browsing fallbacks. Never set global integration_test/gui/wireinject tags or add fake unit seams. Preserve later scope in §8.
