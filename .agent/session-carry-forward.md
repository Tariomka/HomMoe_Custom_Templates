# Carry-forward: Phase 4 complete, Phase 5 pending

Date: 2026-09-11.

## 1. Session goal

Implement approved Phase 4 PNG opacity only. Complete: reusable per-image half-alpha
mask, one composite per roadless edge, opaque geometry/artwork preserved. Based on
owner commit `a574a9e`, with clean tree/index at start; current work is uncommitted.
Coverage **75.0% → 75.1%**, build/tests/lint/layout pass and independent review approves.
Phase 5 batch closeout and owner visual/engine review remain pending. The editor
legend stays absent; Preview legend stays unchanged.

**Reading historical sections:** the retained Phase 3 inventories/evidence below
are historical, not current pending edits. Section 8 is intentionally preserved
verbatim by owner request; its old "Phase 4 unstarted" status is superseded by this
section and the active plan. All later-scope decisions there remain in force.

Previous session goal (historical):

Complete Phase 3 closeout after the owner's review and commits through `cd2b4df`.
Done: fresh coverage **75.0%** exceeds the **74.9%** baseline; formatting, lint,
build, default/unit/integration/GUI tests and independent review pass. The owner
removed the editor dialog legend as unnecessary; preserve that decision. Phase 4
PNG opacity is next-session work, not implemented here. No Phase 3 blockers remain.

## 2. Fixes applied

Phase 4:

- [previewGeneratorService.go](../internal/services/preview_service/previewGeneratorService.go):
  lazy local `image.Alpha` reused per render. Src stamps alpha 128, Over composite
  once per roadless edge, clear only touched bounds. Opaque edges keep Src and exact
  sampling/spacing/clipping/dash geometry; marker ordering and opacity stay unchanged.
- Existing uniform brush allocation is shared across segments; no per-edge full-image
  buffers, persistent/shared scratch, new policy, output-directory or GUI changes.

Phase 3 implementation, owner-committed through `cd2b4df`:

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

Final closeout: removed the obsolete 60-pixel editor-legend exclusion from
[roadStyleVisuals_integration_test.go](../test/integration/gui/roadStyleVisuals_integration_test.go).
Canvas-wide color assertions now inspect the entire editor canvas. No production
code or goldens changed during closeout; owner edits remain intact.

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

Phase 4 now consumes projected `HasRoad`, independently of effective `Type` and
`ExplicitPortal`. Repeated stamps in one edge cannot darken it; distinct crossing
or retraced edges compound. The unchanged opaque path is pinned to pre-change
fingerprints. Player/neutral/arena artwork remains fully opaque. References below
to PNG as future Phase 4 work are historical.

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
  road-state color and use thicker width. Preview legend has `Road`, `No road`,
  `Portal`, and `Portal without road`. The editor has no legend by owner decision.
- Pending editor create/type changes must show checkbox policy before Apply through
  handler/service logic. Cancel must not mutate retained data. No per-edge checkbox.
- PNG roadless strokes use the existing dark color at 50% once per edge, with a reusable
  mask and no compounding among one edge's stamps. Geometry, dashes, clipping, and
  opaque markers remain unchanged. PNG work is Phase 4, separate from Phase 3 GUI work.

## 4. File modifications

Current Phase 4 inventory (six tracked/untracked paths):

- [previewGeneratorService.go](../internal/services/preview_service/previewGeneratorService.go): raster mask/compositing.
- [previewRasterFixture.go](../test/test_helpers/previewRasterFixture.go): explicit opaque fixtures and copied roadless variants.
- [createPreviewImage_test.go](../test/unit/internal/services/preview_service/previewGeneratorService/createPreviewImage_test.go): opaque baseline fingerprints and exact roadless/overlap/clipping/marker tests.
- [preview_raster_test.go](../test/performance/preview_raster_test.go), new: untagged allocation benchmark, 0/1/32/128 edges.
- [plan](plans/batch-c-road-policy-and-visuals.md): Phase 4 Complete with measured evidence; Phase 5 pending.
- [handoff](session-carry-forward.md): current state; complete §8 retained verbatim.

Ignored [visual comparison](tmp/phase4-visual.html) embeds four PNGs as data URIs in
HTML, produced in memory via public raster APIs; no PNG files exported. Temporary
hash/visual Go source was removed. Ignored coverage reports refreshed. No production
assets, snapshots, dependencies, protected files, GUI or generated Wire changed.

Historical Phase 3 closeout inventory (now committed by owner at `a574a9e`):

- [plan](plans/batch-c-road-policy-and-visuals.md): Phase 3 Complete, owner legend
  decision, final verification and Phase 4 boundary.
- [handoff](session-carry-forward.md): current evidence, minimal diff and next-session scope.
- [roadStyleVisuals_integration_test.go](../test/integration/gui/roadStyleVisuals_integration_test.go):
  delete obsolete editor-legend crop/helper; use the full canvas for both color censuses.

Eight reported test/helper files were normalized using an explicit `gofmt -w` list.
Besides the visual test above, these seven files remain modified in Git status but
have no normalized content diff (line endings only):

- [zoneEditorRoadPolicy_integration_test.go](../test/integration/gui/zoneEditorRoadPolicy_integration_test.go)
- [appRunnerSemantics.go](../test/test_helpers/integration_common/appRunnerSemantics.go)
- [color_test.go](../test/unit/app/gui/utils/connectionLineStyle/color_test.go)
- [newEditorConnectionLineStyle_test.go](../test/unit/app/gui/utils/connectionLineStyle/newEditorConnectionLineStyle_test.go)
- [newPreviewConnectionLineStyle_test.go](../test/unit/app/gui/utils/connectionLineStyle/newPreviewConnectionLineStyle_test.go)
- [isRoadTypeCastle_test.go](../test/unit/internal/helpers/road_helpers/roadType/isRoadTypeCastle_test.go)
- [isRoadTypeConnection_test.go](../test/unit/internal/helpers/road_helpers/roadType/isRoadTypeConnection_test.go)

No goldens or production files changed during closeout. No files created or deleted.
Index untouched. Ten paths total are shown modified, with content changes in three.

Historical implementation inventory below (already owner-committed through `cd2b4df`,
not a pending-work list). **new** means introduced during Phase 3, not untracked now.
Any reference below to the editor legend describes the earlier implementation and
is superseded by the owner's removal. Paths are relative to this handoff.

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

The preceding list is historical. Current uncommitted work is listed at the start of
this section: three content edits plus seven line-ending-only files. The owner's
implementation and snapshot changes are committed.

## 5. Tests added or updated

Phase 4 verification, 2026-09-11:

- Pre-change fingerprints: 20 connector footprints and two complete opaque canvases
  captured at `a574a9e` through the public API; all unchanged after implementation.
- Roadless tests observed red on original renderer, then green: exact half-alpha
  pixels over real background; short/long/curved/clipped direct and portal, portal
  shape without explicit Portal, gaps, overlapping stamps, crossings/retracing,
  reused-mask clearing, mixed roaded ordering, degeneracies and opaque artwork.
- Fresh pre/post full unit coverage tasks PASS; final cache-eligible rerun PASS,
  **75.0% → 75.1%**. All changed raster functions 100%; unchanged constructor 75%.
- Build, preview unit suites, default `go test ./test/...`, tagged integration,
  tagged GUI task PASS (GUI 30.052 seconds). Layout/format/diff/diagnostics PASS;
  report-only lint 0 issues, three existing unused-exclusion warnings.
- Independent Claude Opus 5 review approved without blockers. Its missing clipped/
  long-curved portal cases and misleading test comments were corrected; tests rerun.
- Public API benchmark, Windows/amd64, 60 iterations: roadless 1/32/128 direct
  edges ~10,286,600 B/op and 1,957,207 allocs/op, independent of edge count. Opaque
  ~9,795,010 B/op and 1,957,205 allocs/op. Scratch costs ~492 KB +2 allocations per
  render. Before, opaque direct 128 used 9,844,128 B/op, 1,960,275 allocs/op.
  Post-change direct roadless 1/32/128 ~33.4/37.5/48.3 ms versus opaque 128 ~35.2 ms.
  Background allocation floor intentionally untouched; full measurement in plan.
- Browser visual evidence inspected: same geometry/dashes/markers, lighter strokes,
  darker distinct-edge overlap. This is renderer validation, not engine validation.
- Windows only: no Linux, race, deployment or in-game verification claimed.

Final Phase 3 verification, after owner review, 2026-09-11:

- `go build ./...`: PASS.
- Fresh full unit coverage task (`-count=1`): PASS, **75.0%** before closeout.
- Final cache-eligible full unit coverage rerun after formatting/test adjustment:
  PASS, **75.0%**. HTML and LCOV regenerated. Baseline **74.9%**, no shortfall.
- Shared line-style constructors/Color, `StampConnectionRoad`, pending-policy service
  methods and type-change handlers: **100% statement coverage**.
- `go test ./test/...`: PASS (default suite including units and untagged integration).
- `go test -tags=integration_test ./test/integration/...`: PASS.
- `go test -tags='integration_test,gui' ./test/integration/...`: PASS, including
  actual GUI execution in **30.335 seconds**. Existing goldens passed unchanged.
- Test-layout checker, editor diagnostics, changed-file `gofmt -l`, diff whitespace:
  PASS. Report-only lint: **0 issues**, three existing unused-exclusion warnings.
- Independent Claude Opus 5 current-source review: approved, no blockers. Review
  prompted removal of the stale crop, verified by the full tagged GUI run. Optional
  cleanup/pre-existing observations are recorded in §8, not Phase 3 blockers.
- Windows only. No Linux, race, benchmark, deployment or in-game validation claimed.

Historical first-wrap-up results below are superseded by the final checks above:

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

Current: branch `AD/road_and_graph_invariants`, HEAD `a574a9e`. The owner committed
the entire Phase 3 closeout before this session. Tree/index were clean initially.
Five modified files and one untracked benchmark listed in §4; no staged changes.
No protected/GUI/Wire changes. No staging, unstaging, commit, push, stash or branch
operation by the agent. Ignored visual HTML and coverage reports are scratch only.

Historical Phase 3 snapshot:

Branch `AD/road_and_graph_invariants`, HEAD `cd2b4df`; clean tree/index at closeout
start. Owner committed implementation/review changes in `053d4fd`, `b6d4257`,
`cd2b4df`. Current unstaged files: plan, handoff, visual test and seven line-ending-only
test/helper files listed in §4 (ten paths total, three normalized content diffs).
No untracked files or staged changes. No protected-tree or PNG raster changes since
Phase 3 baseline `4923a5c`. Phase 1/2 commits remain `07ca8b6`/`abf0d36`.
No staging, unstaging, commits, pushes, stashes or branch changes by this agent.

## 7. Rejections / things the user declined

- Phase 4: an initial lowercase subagent model identifier was unavailable; corrected
  to the exact listed name and retried successfully. One golines finding corrected
  with an explicit edit; final lint zero. No unresolved tool failures.
- Final §8 comparison initially failed because PowerShell 5.1 decoded Git/file text
  differently; Python fallback was unavailable. Explicit UTF-8 for both Git stdout
  (`ProcessStartInfo.StandardOutputEncoding`) and disk confirmed exact preservation.
  No actual §8 content was changed and no tooling was installed.
- Do not turn Phase 4 into optional GUI/export cleanup or broader Phase 5 work.
  No global allocation thresholds, background optimization, output-path changes,
  per-edge image allocation or engine-default claims introduced.

- Owner on 2026-09-11: stop circling/overcomplicating; wrap up implemented Phase 3
  and record unresolved work. Do not continue unrelated coverage additions or start PNG.
- The owner subsequently reviewed/committed changes and requested bounded Phase 3
  closeout before a new Phase 4 session. That closeout is complete, not deferred.
- Owner explicitly removed the editor legend as useless. Do not restore it or the
  deleted legend test. Preview keeps its four-entry legend.
- Supplemental tests and owner snapshot changes are preserved. No unrelated coverage
  additions, optional export cleanup, new production work or Phase 4 implementation.

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

1. Read [AGENTS.md](../AGENTS.md), this handoff, and [the active Batch C plan](plans/batch-c-road-policy-and-visuals.md).
2. Inspect Git and preserve every current edit. Phase 4 implementation is complete;
  don't reimplement it, reopen approval questions or restore the editor legend.
3. Owner can review [the visual comparison](tmp/phase4-visual.html). It's ignored
  local evidence, not a release artifact or an exported game template.
4. Resume Phase 5 only: final batch-wide verification/review and record closeout,
  reusing current evidence when applicable. Baseline is now **75.1%**. Owner validates
  true/false/nil engine behavior and performs staging/commits; no agent Git mutations.
5. Preserve complete later scope in §8 verbatim. Its original Phase 4 status is
  intentionally historical and superseded above; all its later decisions remain.

Deployment: none performed. Windows source build passes; no schema migration,
dependency installation, Wire regeneration or output-directory change required.
Formal batch release awaits Phase 5 and owner engine review. Owner handles commits.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md)
> and [the active Batch C plan](plans/batch-c-road-policy-and-visuals.md). Owner approval
> remains in force. Phase 1 (`07ca8b6`), Phase 2 (`abf0d36`) and Phase 3 (through
> `a574a9e`) are owner-committed on `AD/road_and_graph_invariants`. Phase 4 is complete
> and uncommitted: PNG per-edge opacity, focused tests/benchmark, plan and handoff.
> Inspect and preserve every current edit. Build, full unit/default/integration/GUI,
> layout, formatting and lint pass; independent review approved. Coverage **75.1%**,
> all changed raster functions 100%. Mask costs ~492 KB +2 allocations per image,
> not per edge. Before/after visual evidence is in ignored `.agent/tmp/phase4-visual.html`.
> Resume Phase 5 final batch closeout only; do not reimplement Phase 4 or repeat
> approval questions. Owner engine validation remains pending. Preserve the editor
> legend's absence; Preview legend stays.
>
> Never modify protected `data/`, template schema, or registry trees. Keep all paths
> cross-platform and never change or persist the game output directory. Test nontrivial
> changes and measure coverage; never stage, unstage, commit, push, stash, or switch
> branches. Never bulk-rewrite or hand-edit generated Wire. Never enable global
> `integration_test`, `gui`, or `wireinject` tags, and never add fake unit seams.
> Keep multi-session plans durable and resumable; read the plan before changing scope.
> Preserve the complete later scope in §8 verbatim. Keep business policy in the
> handler/service layer and carry road display state separately from effective preview
> type. Phase 4 uses one reusable mask, 50% once per edge, with no compounding among
> that edge's stamps; distinct crossing edges may compound. Preserve the opaque raster
> path, sampling/spacing, clipping, curves, dashes and opaque markers. Full details and
> historical decisions are in this handoff. Section 8 is verbatim historical text;
> only its old Phase 4 status is superseded, not its later-scope decisions.
