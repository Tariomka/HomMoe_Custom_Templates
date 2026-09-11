# Carry-forward: Phase 3 implemented, verification closeout deferred

Date: 2026-09-11.

## 1. Session goal

Implement approved Phase 3 pending-editor policy and GUI road styles. Implementation
is done and uncommitted. The owner requested wrap-up rather than further iteration:
record unresolved items and stop. Build/unit/layout checks pass, but coverage is
74.8% against the 74.9% baseline and formatting/final review remain outstanding.
Do not restart implementation, hunt unrelated coverage, or begin Phase 4 automatically.

## 2. Fixes applied

Phase 3, current uncommitted work:

- Pending connection creation and type changes obey the current road checkbox before
  Apply through [zoneEditorHandler.go](../internal/handlers/zoneEditorHandler.go) and
  [zoneEditorService.go](../internal/services/connection_editor/zoneEditorService.go),
  reusing [roadPolicyService.go](../internal/services/zones/roadPolicyService.go).
- Type changes clone the input before mutation. GUI Apply/Cancel and retained-data
  isolation are covered in [zoneEditorRoadPolicy_integration_test.go](../test/integration/gui/zoneEditorRoadPolicy_integration_test.go).
- Preview road display state is separate from effective shape classification in
  [previewLayoutService.go](../internal/services/preview_service/previewLayoutService.go).
- Shared colors, legends and selected-edge width replace type-only road styling in
  [connectionLineStyle.go](../app/gui/utils/connectionLineStyle.go),
  [draw.go](../app/gui/utils/draw.go) and
  [zoneEditorCanvas.go](../app/gui/dialogs/zoneEditorCanvas.go).

Historical completed work follows; no reimplementation is needed.

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

Phase 3 visual contract is now implemented; PNG remains Phase 4:

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

Current Phase 3 inventory (all existing work preserved; only plan/handoff edited
during final wrap-up). Paths are relative to this handoff. New files are marked **new**.

Production and documentation:

- [plan](plans/batch-c-road-policy-and-visuals.md): implementation checklist and unresolved gates.
- [handoff](session-carry-forward.md): current status and resumable inventory.
- [.gitignore](../.gitignore): correct snapshot failure scratch directory.
- [legend.go](../app/gui/constants/legend.go): shared four-entry connection key.
- [zoneEditorCanvas.go](../app/gui/dialogs/zoneEditorCanvas.go): colors, selection width and canvas key.
- [zoneEditorConnectionProps.go](../app/gui/dialogs/zoneEditorConnectionProps.go): type-change handler call.
- [zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go): forward checkbox on creation.
- [zoneEditorDialog_testexports.go](../app/gui/dialogs/zoneEditorDialog_testexports.go): gated GUI observations.
- [window_testexports.go](../app/gui/editor/window_testexports.go): gated GUI observations.
- [previewPanel.go](../app/gui/panels/previewPanel.go): shared legend rendering.
- [colors.go](../app/gui/themes/colors.go): roadless palette.
- [connectionLineStyle.go](../app/gui/utils/connectionLineStyle.go), **new**: shared GUI color selection.
- [draw.go](../app/gui/utils/draw.go): use projected road style.
- [legendWidget.go](../app/gui/widgets/legendWidget.go), **new**: reusable legend row.
- [zoneEditorConnectionRequestDto.go](../internal/dtos/zoneEditorConnectionRequestDto.go): creation road setting.
- [zoneEditorConnectionTypeRequestDto.go](../internal/dtos/zoneEditorConnectionTypeRequestDto.go), **new**: type-change request.
- [guiHandler.go](../internal/handlers/guiHandler.go): delegate type changes.
- [zoneEditorHandlerInterface.go](../internal/handlers/handler_interfaces/zoneEditorHandlerInterface.go): new handler contract.
- [zoneEditorHandler.go](../internal/handlers/zoneEditorHandler.go): pending policy and cloned type changes.
- [previewConnection.go](../internal/models/preview/previewConnection.go): projected road/explicit-portal fields.
- [zoneEditorService.go](../internal/services/connection_editor/zoneEditorService.go): pending-policy operations.
- [zoneEditorServiceInterface.go](../internal/services/connection_editor/zoneEditorServiceInterface.go): operation contracts.
- [previewLayoutService.go](../internal/services/preview_service/previewLayoutService.go): project road status.
- [roadPolicyService.go](../internal/services/zones/roadPolicyService.go): reusable single-connection policy.
- [roadPolicyServiceInterface.go](../internal/services/zones/zone_interfaces/roadPolicyServiceInterface.go): stamping contract.

Integration tests and support:

- [roadStyleVisuals_integration_test.go](../test/integration/gui/roadStyleVisuals_integration_test.go), **new**: 14 rendering/selection/regeneration tests.
- [zoneEditorRoadPolicy_integration_test.go](../test/integration/gui/zoneEditorRoadPolicy_integration_test.go), **new**: 11 pending-policy/Apply/Cancel tests.
- [appRunnerSemantics.go](../test/test_helpers/integration_common/appRunnerSemantics.go): semantic lookup support.
- [appRunnerSnapshots.go](../test/test_helpers/integration_common/appRunnerSnapshots.go): frame capture support.
- [handlerCoordinates.go](../test/test_helpers/integration_common/handlerCoordinates.go): GUI coordinates.
- [layoutAndZonesTabHandler.go](../test/test_helpers/integration_common/layoutAndZonesTabHandler.go): road-setting interaction.
- [zoneEditorHandler.go](../test/test_helpers/integration_common/zoneEditorHandler.go): editor interaction support.
- [roadPolicyServiceMock.go](../test/test_helpers/roadPolicyServiceMock.go): stamping mock.
- [templateHandlerMock.go](../test/test_helpers/templateHandlerMock.go): preview test behavior.
- [zoneEditorServiceMock.go](../test/test_helpers/zoneEditorServiceMock.go): new service operations.
- 280 tracked goldens in the [snapshot directory](../test/test_helpers/integration_common/snapshot/__snapshots__): shared legend/edge-color changes. The exact file inventory is `git diff --name-only -- '*.golden'`; none deleted or newly untracked at wrap-up. Prior work reported inspecting all 280; not re-inspected during wrap-up.

Unit tests:

- [close_test.go](../test/unit/app/gui/drivers/dialogHost/close_test.go), **new**: dialog close behavior.
- [newTheme_test.go](../test/unit/app/gui/themes/theme/newTheme_test.go): theme coverage.
- [color_test.go](../test/unit/app/gui/utils/connectionLineStyle/color_test.go), **new**: color matrix.
- [newEditorConnectionLineStyle_test.go](../test/unit/app/gui/utils/connectionLineStyle/newEditorConnectionLineStyle_test.go), **new**: editor classification.
- [newPreviewConnectionLineStyle_test.go](../test/unit/app/gui/utils/connectionLineStyle/newPreviewConnectionLineStyle_test.go), **new**: preview classification.
- [toF32Point_test.go](../test/unit/app/gui/utils/math/toF32Point_test.go), **new**: coordinate conversion.
- [toVec2_test.go](../test/unit/app/gui/utils/math/toVec2_test.go), **new**: coordinate conversion.
- [changeZoneEditorConnectionType_test.go](../test/unit/internal/handlers/guiHandler/changeZoneEditorConnectionType_test.go), **new**: facade delegation.
- [createZoneEditorConnection_test.go](../test/unit/internal/handlers/guiHandler/createZoneEditorConnection_test.go): creation delegation.
- [handlerDependenciesStub_test.go](../test/unit/internal/handlers/guiHandler/handlerDependenciesStub_test.go): new operation stub.
- [changeZoneEditorConnectionType_test.go](../test/unit/internal/handlers/zoneEditorHandler/changeZoneEditorConnectionType_test.go), **new**: cloning/policy.
- [createZoneEditorConnection_test.go](../test/unit/internal/handlers/zoneEditorHandler/createZoneEditorConnection_test.go): checkbox policy.
- [cloneZoneContentRows_test.go](../test/unit/internal/helpers/editor_state_helpers/zoneContentRow/cloneZoneContentRows_test.go), **new**: clone coverage.
- [isRoadTypeCastle_test.go](../test/unit/internal/helpers/road_helpers/roadType/isRoadTypeCastle_test.go), **new**: reference classification.
- [isRoadTypeConnection_test.go](../test/unit/internal/helpers/road_helpers/roadType/isRoadTypeConnection_test.go), **new**: reference classification.
- [applyCastleSettingChanges_test.go](../test/unit/internal/services/connection_editor/manualReapplyService/applyCastleSettingChanges_test.go): supplemental regression coverage.
- [snapPosition_test.go](../test/unit/internal/services/connection_editor/zoneEditorGeometryService/snapPosition_test.go): supplemental geometry coverage.
- [applyConnectionRoadPolicy_test.go](../test/unit/internal/services/connection_editor/zoneEditorService/applyConnectionRoadPolicy_test.go), **new**: shared policy delegation.
- [changeConnectionType_test.go](../test/unit/internal/services/connection_editor/zoneEditorService/changeConnectionType_test.go), **new**: type/policy operation.
- [nextFreeZoneLabel_test.go](../test/unit/internal/services/connection_editor/zoneEditorService/nextFreeZoneLabel_test.go): supplemental label coverage.
- [buildPreviewLayout_test.go](../test/unit/internal/services/preview_service/previewLayoutService/buildPreviewLayout_test.go): road/type projection matrix.
- [common_test.go](../test/unit/internal/services/preview_service/previewLayoutService/common_test.go): projection fixtures.
- [stampConnectionRoad_test.go](../test/unit/internal/services/zones/roadPolicyService/stampConnectionRoad_test.go), **new**: true/false/nil and explicit Portal policy.

Ignored local coverage reports were refreshed earlier. Temporary golden inspection
programs outside the repository are disposable scratch, not source dependencies or
deployment artifacts. No source/test files were deleted; no PNG raster or Wire edits.

Historical Phase 2 inventory:

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

The preceding list is historical. Current uncommitted Phase 3 work is listed above;
wrap-up changed only the plan and this handoff, preserving source/tests and the index.

## 5. Tests added or updated

Phase 3 verification at wrap-up, 2026-09-11:

- `go build ./...`: PASS.
- `go test ./test/unit/...`: PASS (cache-eligible final run; the preceding implementation
  work reported a fresh full unit pass).
- `go run ./cmd/testlayoutcheck .`: PASS.
- `git diff --check`: PASS; one CRLF-to-LF normalization warning, not a whitespace error.
- Existing coverage report timestamp 2026-09-11 09:38:58: **74.8%**, versus **74.9%**
  baseline. Earlier intermediate result was 74.6%. This wrap-up read the report and
  did not regenerate it. `StampConnectionRoad`, `ApplyConnectionRoadPolicy`,
  `ChangeConnectionType`, and both create/type-change handler entry points show 100%.
- New GUI tests previously reported PASS, including pending changes, Cancel isolation,
  four colors and selection. No fresh full tagged GUI run during wrap-up.
- Last current-Phase-3 `go test ./test/...` result is not established in the available
  wrap-up record; the prior Phase 2 result below is historical, not a final-tree pass.
- Final report-only lint and independent review approval remain outstanding. Earlier
  review found formatting/coverage issues. No fresh Linux/race/benchmark/in-game claims.

Historical Phase 2 verification:

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

Branch `AD/road_and_graph_invariants`, HEAD `4923a5c`. Index empty at wrap-up.
Working tree contains the modified/new files above plus 280 modified tracked goldens.
New production files and tests remain untracked; tracked edits remain unstaged.
Protected data/schema/registry paths have no diff. Phase 1/2 owner commits remain
`07ca8b6`/`abf0d36`. No staging, unstaging, commits, pushes, stashes or branch changes.

## 7. Rejections / things the user declined

- Owner on 2026-09-11: stop circling/overcomplicating; wrap up implemented Phase 3
  and record unresolved work. Do not continue unrelated coverage additions or start PNG.
- Implementation is complete, but do not label the phase fully verified while the
  coverage shortfall, formatting and final review checks remain unresolved.
- All supplemental tests and snapshot changes are preserved, not reverted or cleaned
  up without owner direction. No new implementation work during this wrap-up.

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

Batch C policy/scope approval remains settled. **Phase 3 implementation is complete;
verification closeout is deferred at owner request.** No new policy questions needed.
Owner Phase 3 authorization, 2026-09-10: "Changes reviewed, you can proceed".
Remaining closeout: coverage 74.8% < 74.9%, formatting, final report-only lint/review,
and final-tree default/tagged integration checks. Phase 4 remains unstarted.

Remaining formatting files (`gofmt -l` at wrap-up; format only an explicit rechecked list):

- [appRunnerSemantics.go](../test/test_helpers/integration_common/appRunnerSemantics.go)
- [connectionLineStyle.go](../app/gui/utils/connectionLineStyle.go)
- [legendWidget.go](../app/gui/widgets/legendWidget.go)
- [roadStyleVisuals_integration_test.go](../test/integration/gui/roadStyleVisuals_integration_test.go)
- [zoneEditorRoadPolicy_integration_test.go](../test/integration/gui/zoneEditorRoadPolicy_integration_test.go)
- [close_test.go](../test/unit/app/gui/drivers/dialogHost/close_test.go)
- [color_test.go](../test/unit/app/gui/utils/connectionLineStyle/color_test.go)
- [newEditorConnectionLineStyle_test.go](../test/unit/app/gui/utils/connectionLineStyle/newEditorConnectionLineStyle_test.go)
- [newPreviewConnectionLineStyle_test.go](../test/unit/app/gui/utils/connectionLineStyle/newPreviewConnectionLineStyle_test.go)
- [toF32Point_test.go](../test/unit/app/gui/utils/math/toF32Point_test.go)
- [toVec2_test.go](../test/unit/app/gui/utils/math/toVec2_test.go)
- [cloneZoneContentRows_test.go](../test/unit/internal/helpers/editor_state_helpers/zoneContentRow/cloneZoneContentRows_test.go)
- [isRoadTypeCastle_test.go](../test/unit/internal/helpers/road_helpers/roadType/isRoadTypeCastle_test.go)
- [isRoadTypeConnection_test.go](../test/unit/internal/helpers/road_helpers/roadType/isRoadTypeConnection_test.go)

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
2. Inspect Git and preserve all changes. Review the completed Phase 3 implementation;
  do not rebuild it or repeat settled approval questions.
3. Only if asked to resume closeout: apply bounded formatting, rerun report-only lint,
  resolve or explicitly accept the 0.1 percentage-point coverage shortfall, and run
  the final default/tagged integration verification. Do not hunt unrelated coverage.
4. Complete independent review and mark the phase fully verified only with evidence.
5. Phase 4 PNG opacity is separate, unstarted work. Do not start automatically.
  Later scope in §8 remains unchanged; engine validation belongs to the owner.

Deployment: no deployment performed. For owner review, the current Windows source
build passes; no schema migration, dependency installation, Wire regeneration or
output-path change is introduced by Phase 3. Formal release awaits the recorded
closeout and remaining Batch C phases.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md)
> and [the active Batch C plan](plans/batch-c-road-policy-and-visuals.md). Owner approval
> remains in force. Phase 1 (`07ca8b6`) and Phase 2 (`abf0d36`) are implemented and
> committed. Phase 3 GUI/pending-edit implementation is done, unstaged/untracked on
> `AD/road_and_graph_invariants` at `4923a5c`. Inspect and preserve every current edit.
> The owner requested wrap-up, not more implementation. Build/unit/layout checks pass;
> coverage is 74.8% versus 74.9% baseline. Formatting, final lint/review and final-tree
> default/tagged integration verification remain recorded, not silently passed.
> Resume only requested bounded closeout; no unrelated coverage hunt or automatic Phase 4.
>
> Never modify protected `data/`, template schema, or registry trees. Keep all paths
> cross-platform and never change or persist the game output directory. Test nontrivial
> changes and measure coverage; never stage, unstage, commit, push, stash, or switch
> branches. Never bulk-rewrite or hand-edit generated Wire. Never enable global
> `integration_test`, `gui`, or `wireinject` tags, and never add fake unit seams.
> Keep multi-session plans durable and resumable; read the plan before changing scope.
> Preserve the complete later scope in §8 verbatim. Keep business policy in the
> handler/service layer, carry road display state separately from effective preview
> type, and leave Phase 4 PNG work untouched.
