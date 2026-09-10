# Carry-forward: Batch C Phase 2 complete, begin Phase 3

Date: 2026-09-10.

## 1. Session goal

Complete approved Batch C Phase 2 road-policy and final-content reconciliation.
Phase 1 and Phase 2 are implemented and owner-committed. Phase 3, pending editor
policy and GUI road styles, is next and has not been implemented. No further approval
question is needed.

## 2. Fixes applied

Phase 1, committed at `07ca8b6`:

- Added explicit-Portal and road classification on the connection model.
- Made public topology repair endpoint-based and road-gated while preserving graph
  repairs, guards, impossible-isolation fallback, and collision-safe edge names.

Prior settled fixes stay intact: manual compare-before-mutation/dirty tracking and
snapshot isolation; no output browsing-fallback authorization; additive PNG ceiling
sample count with existing stamp spacing; configured tournament SaveArmy (nil rules
plus tournament selector means false; normal constructor default remains true).

Phase 2, committed at `abf0d36`:

- Added `RoadPolicyService` reconciliation for authoritative non-Portal road flags,
  scoped mandatory-content validation, eligible approach restoration, and invalid
  generated-road removal.
- Kept internal castle/object/foothold roads independent of the between-zone checkbox.
  Preserved valid custom routes, opaque references, explicit Portal flags, spawn
  anchoring, and marker-only arena connector routing.
- Final generation reconciles after content and arena placement. `UpdateTemplate`
  names connections directly and runs policy only when `EditorState` is present; nil
  state preserves supplied roads, flags, and mandatory content while retaining graph
  diagnostics.
- Cloned castle-road slices before deletion/appending to prevent alias mutation and
  rebased imported non-arena anchors to actual indices. Arena markers are excluded
  from automatic anchors and never gain a new castle road.
- Regenerated Wire. The shared golden retains full equality; its only intentional
  output delta is `road: true` on `Rnd-A-B`.

## 3. Features added / changed

Implemented contract:

- With settings present, `GenerateRoads` overwrites `Road` on every non-explicit-Portal
  connection, on or off, including custom/imported, Default, empty, arena and proximity
  connections. Explicit Portal flags, including false and nil, are preserved.
- When roads are off, reconciliation removes complete zone-road segments referencing
  disabled connections. Internal roads remain; valid Portal approaches remain. Turning
  roads back on restores eligible target sets, not redundant connector record equality.
- Final reconciliation treats nil/empty mandatory content as authoritative, resolves
  names within each zone's final groups, and removes only confirmed invalid `MainObject`,
  incident `Connection`, or named `MandatoryContent` references.
- The legacy settings-free `RebuildZoneConnectionRoads` compatibility entry point
  explicitly validates with roads enabled and unknown content. It has no production
  caller; no deletion/refactor is required in this batch.

Approved full Batch C contract still pending in the visual surfaces:

- The shared display classifier is explicit false = roadless, explicit true = roaded,
  and nil = roaded only for explicit Portal. Road display classification stays separate
  from effective preview `Type` and portal placement rules.
- Preview and editor roadless colors are `#B0B0B0` for non-Portals and `#90EE90` for
  explicit Portals. Existing roaded colors stay. Selected editor edges retain the same
  road-state color and use thicker width. Legends require `Road`, `No road`, `Portal`,
  and `Portal without road`.
- Pending editor create/type changes must show checkbox policy before Apply through
  handler/service logic. Cancel must not mutate retained data. No per-edge checkbox.
- PNG roadless strokes use the existing dark color at 50% once per edge, with a reusable
  mask and no compounding among one edge's stamps. Geometry, dashes, clipping, and
  opaque markers remain unchanged. PNG work is Phase 4, separate from Phase 3 GUI work.

## 4. File modifications

Phase 2 committed files at `abf0d36`:

- Docs: [batch-c-road-policy-and-visuals.md](plans/batch-c-road-policy-and-visuals.md)
  and [session-carry-forward.md](session-carry-forward.md).
- Composition and handlers: [providerSets.go](../internal/composition/providerSets.go),
  [wire_gen.go](../internal/composition/wire_gen.go), and
  [templateHandler.go](../internal/handlers/templateHandler.go).
- Models and services: [roadReconciliationRequest.go](../internal/models/roadReconciliationRequest.go),
  [zoneEditorService.go](../internal/services/connection_editor/zoneEditorService.go),
  [zoneEditorServiceInterface.go](../internal/services/connection_editor/zoneEditorServiceInterface.go),
  [templateGenerator.go](../internal/services/template_generator/templateGenerator.go),
  [roadFactory.go](../internal/services/zones/roadFactory.go),
  [roadFactoryInterface.go](../internal/services/zones/zone_interfaces/roadFactoryInterface.go),
  [roadPolicyScope.go](../internal/services/zones/roadPolicyScope.go),
  [roadPolicyService.go](../internal/services/zones/roadPolicyService.go), and
  [roadPolicyServiceInterface.go](../internal/services/zones/zone_interfaces/roadPolicyServiceInterface.go).
- Integration and fixtures: [roadPolicyApply_integration_test.go](../test/integration/roadPolicyApply_integration_test.go),
  [defaultTemplate.json](../test/test_helpers/defaultTemplate.json),
  [roadFactoryMock.go](../test/test_helpers/roadFactoryMock.go),
  [roadPolicyServiceMock.go](../test/test_helpers/roadPolicyServiceMock.go),
  [templateGenerator.go](../test/test_helpers/templateGenerator.go), and
  [zoneEditorService.go](../test/test_helpers/zoneEditorService.go).
- Unit tests: [common_test.go](../test/unit/internal/handlers/templateHandler/common_test.go),
  [updateTemplate_test.go](../test/unit/internal/handlers/templateHandler/updateTemplate_test.go),
  [newZoneEditorService_test.go](../test/unit/internal/services/connection_editor/zoneEditorService/newZoneEditorService_test.go),
  [createSpawnZone_test.go](../test/unit/internal/services/template_generator/providers/topology/base/topologyBase/createSpawnZone_test.go),
  [generate_test.go](../test/unit/internal/services/template_generator/templateGenerator/generate_test.go),
  [createOuterZoneRoads_test.go](../test/unit/internal/services/zones/roadFactory/createOuterZoneRoads_test.go),
  [common_test.go](../test/unit/internal/services/zones/roadPolicyService/common_test.go),
  [newRoadPolicyService_test.go](../test/unit/internal/services/zones/roadPolicyService/newRoadPolicyService_test.go),
  [rebuildCastleRoads_test.go](../test/unit/internal/services/zones/roadPolicyService/rebuildCastleRoads_test.go),
  [rebuildZoneConnectionRoads_test.go](../test/unit/internal/services/zones/roadPolicyService/rebuildZoneConnectionRoads_test.go),
  and [reconcile_test.go](../test/unit/internal/services/zones/roadPolicyService/reconcile_test.go).

Current uncommitted documentation work is this rewritten handoff. No source, test,
protected-tree, generated-output, index, or branch change is part of this documentation
update.

## 5. Tests added or updated

Phase 2 added policy, factory, handler, generator, integration, arena-only,
roads-disabled, nil-state, slice-alias, malformed-anchor, custom-road, and final-output
JSON matrices. Factory mock coverage is legitimate constructor dependency coverage for
malformed anchors. Regression tests were observed failing when the implemented fixes
were removed.

Windows verification PASS:

- `go build ./...`.
- Full unit suite with `-count=1`.
- Fresh pre/post full coverage: **74.6% -> 74.9%**, followed by a final
  cache-eligible full-coverage rerun. New road-policy functions, including
  `rebaseAnchorRef`, and factory coverage are 100%; touched handler/editor entry
  points are 100%.
- Default `go test ./test/...`, tagged integration, and tagged
  `integration_test,gui` under the whole integration tree. Actual GUI execution took
  about 37 seconds.
- Test-layout and diff-whitespace checks.
- Report-only lint: 0 issues, with three existing unused-exclusion warnings.

Independent Claude Opus review approved after corrections for marker-only arena
anchoring, slice alias mutation, and a vacuous From-vs-To assertion. No Linux, race,
benchmark, or in-game validation is claimed.

## 6. Git status snapshot

Branch `AD/road_and_graph_invariants`, clean `HEAD abf0d36`. Phase 1 is the earlier
owner commit `07ca8b6`; the owner committed all Phase 2 source, tests, wiring, fixture,
plan, and prior handoff changes. This documentation rewrite intentionally leaves the
handoff modified. Do not stage, unstage, commit, push, stash, or switch branches.

## 7. Rejections / things the user declined

- Never skip cleanup because final content is nil or empty: it is authoritative during
  `Reconcile`. The settings-free legacy `RebuildZoneConnectionRoads` explicitly uses
  validation false and has no production caller.
- Never auto-anchor arena-only zones. Exclude the marker while preserving spawn
  anchoring and rebasing shifted imported non-arena anchors.
- Clone castle roads before deletion/appending; do not mutate shared backing arrays.
- Preserve explicit Portal flags even when false. Roads-off does not disable internal
  roads or valid Portal approaches.
- Do not replace restored target-set behavior with redundant connector-record equality;
  no deletion/refactor is needed for the legacy compatibility entry point.
- A/D and B are closed; no old-document lookup, validation, or closure work is
  required.
- A source-build warning from guessed nonexistent
  `internal/models/generatorConfig.go` was resolved by using the actual config path.
  One terminal command accidentally prepended `d` to `go`, failed, then succeeded on
  retry.
- No protected edits, global tags, bulk rewrites, fake seams, or hand-edited Wire.
  Wire was regenerated and its golden delta is exactly one `road: true` field.

## 8. Open questions and confirmed later scope

Batch C source verification, baseline, policy questions, scope approval and written-plan
implementation approval are done. **Phase 2 is complete; begin Phase 3 without asking approval or settled policy again.**
Owner Phase 3 authorization, 2026-09-10: "Changes reviewed, you can proceed".
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

1. Read [AGENTS.md](../AGENTS.md), this handoff, and [the active Batch C plan](plans/batch-c-road-policy-and-visuals.md).
2. Inspect Git and preserve owner changes. Phase 3 is authorized; do not ask for approval again.
3. Implement only Phase 3 from the fresh **74.9%** coverage baseline: pending create/type policy through handler/service, preview road display state, GUI colors/legends, and selected-edge styling.
4. Keep Phase 4 PNG opacity separate and leave §8 scope unchanged.
5. Run focused model/service/handler and preview tests after each first substantive edit, then the Phase 3 verification matrix. The owner, not the agent, validates engine behavior in-game.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md)
> and [the active Batch C plan](plans/batch-c-road-policy-and-visuals.md). Owner approval
> remains in force. Phase 1 (`07ca8b6`) and Phase 2 (`abf0d36`) are implemented and
> committed; begin Phase 3 without asking again. Phase 3 GUI/pending-edit work is not
> implemented. Owner Phase 3 authorization, 2026-09-10: "Changes reviewed, you can
> proceed". Begin from the fresh 74.9% coverage baseline.
>
> Never modify protected `data/`, template schema, or registry trees. Keep all paths
> cross-platform and never change or persist the game output directory. Test nontrivial
> changes and measure coverage; never stage, unstage, commit, push, stash, or switch
> branches. Never bulk-rewrite or hand-edit generated Wire. Never enable global
> `integration_test`, `gui`, or `wireinject` tags, and never add fake unit seams.
> Preserve the complete later scope in §8 verbatim. Keep business policy in the
> handler/service layer, carry road display state separately from effective preview
> type, and leave Phase 4 PNG work untouched.
