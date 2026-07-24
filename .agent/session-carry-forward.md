# Session Carry-Forward

## Session goal

Complete Phase 2 of `plans/clean-architecture-refactoring.md`: restore model, generic-algorithm, classification, catalog, and tuning ownership while preserving preview, content, manual-edit, and zone-editor behavior.

## Fixes applied

- Moved generation-specific radial plan placement from `internal/models/position.go` to `internal/services/template_generator/providers/topology/position_layout/positionLayoutService.go`.
- Moved Delaunay triangulation, bounds, and closest-cross-component calculations to `internal/helpers/geometry/` as pure functions.
- Removed `internal/models/zoneAdjacency.go`; generic undirected links, BFS distances, and connected components now live in `internal/helpers/graph/adjacency.go`.
- Moved entity-derived zone quality and connection-guard classification to `internal/services/zones/zoneClassifier.go`.
- Reduced `internal/helpers/zone_helpers/zoneNameType.go` to lexical name/type parsing and routed GUI classification through focused handler port methods.
- Moved registry-backed neutral profile assembly from the model to `internal/common/common_zones/neutralZoneProfile.go`, returning fresh slices.
- Moved generator-config-to-tuning construction from the model to the cycle-free `internal/services/template_generator/generation_tuning/generationTuningFactory.go` package.

## Features added / changed

- Delaunay edge endpoints are normalized and edge output is lexicographically sorted. This is the only intentional behavior change and removes map-order nondeterminism.
- `models.Position`, `ConnectionIndexes`, `Positions`, `neutral_zone.Profile`, `neutral_zone.Quality`, and `models.GenerationTuning` now retain passive data or small value semantics only.
- `GUIHandler.GetZoneQuality` and `GetZoneConnectionGuardQuality` preserve the Phase 1 GUI boundary; `app/gui` still imports no internal services, mappers, or validators.
- `common_zones.GetNeutralZoneProfile` clones or concatenates all pool slices so a caller cannot mutate future catalog results.
- Topology test fixtures use `test_helpers.NewGenerationTuning`, which delegates to the production factory.
- Phase 2 is complete in `plans/clean-architecture-refactoring.md`.

## File modifications

Production files created:

- `internal/common/common_zones/neutralZoneProfile.go` — neutral profile catalog factory.
- `internal/helpers/geometry/closestAcrossComponents.go` — closest cross-component position pair.
- `internal/helpers/geometry/delaunayTriangulation.go` — deterministic Bowyer-Watson triangulation.
- `internal/helpers/geometry/positionBounds.go` — component-wise bounds.
- `internal/helpers/graph/adjacency.go` — generic adjacency, link, BFS, and components.
- `internal/services/template_generator/generation_tuning/generationTuningFactory.go` — tuning factory.
- `internal/services/template_generator/providers/topology/position_layout/positionLayoutService.go` — tier placement service.
- `internal/services/zones/zoneClassifier.go` — entity and guard classification.

Production files edited or deleted:

- `app/gui/dialogs/zoneEditorConnectionProps.go`, `zoneEditorZoneProps.go`, and `app/gui/interfaces/backendInterface.go` — route classification through the zone-editor port.
- `internal/handlers/guiHandler.go` — owns classifier/tuning factory and exposes focused quality methods.
- `internal/helpers/zone_helpers/zoneNameType.go` — lexical helpers only.
- `internal/models/generationTuning.go`, `neutral_zone/neutralZoneProfile.go`, `neutral_zone/neutralZoneQuality.go`, and `position.go` — remove misplaced construction/classification/algorithm logic.
- `internal/models/zoneAdjacency.go` — deleted after graph extraction.
- `internal/services/connection_editor/connectionEditor.go`, `manualReapply.go`, and `zoneEditor.go` — use classifier, common profile factory, and tuning factory.
- `internal/services/preview_service/previewLayoutService.go` and `template_generator/providers/mandatoryContentProvider.go` — constructor-owned classifier.
- `internal/services/template_generator/templateGenerator.go` — constructor-owned tuning factory.
- `internal/services/template_generator/providers/topology/base/topologyBase.go`, `circlesTopology.go`, `randomTopology.go`, and `tournament_variant/balancedClusterService.go` — use graph, geometry, layout, profile, and tuning owners.
- `internal/services/zones/zoneLabelProvider.go` — generic graph operations.
- `plans/clean-architecture-refactoring.md` — Phase 2 marked complete with verification summary.

Test support and tests created or relocated:

- `test/test_helpers/generationTuning.go` and `templateHandlerMock.go` — production-factory tuning helper and new GUI port methods.
- `test/unit/internal/common/common_zones/neutralZoneProfile/` — profile catalog and fresh-slice tests.
- `test/unit/internal/helpers/geometry/` — Delaunay, bounds, and closest-component tests.
- `test/unit/internal/helpers/graph/adjacency/` — adjacency operation tests.
- `test/unit/internal/services/template_generator/generation_tuning/generationTuningFactory/` — constructor and creation tests.
- `test/unit/internal/services/template_generator/providers/topology/position_layout/positionLayoutService/` — service constructor and coordinate tests.
- `test/unit/internal/services/zones/zoneClassifier/` — quality and guard classification matrices.
- `test/unit/internal/handlers/guiHandler/getZoneQuality_test.go` and `getZoneConnectionGuardQuality_test.go` — GUI boundary delegation.
- Existing model/helper ownership tests were deleted after relocation.
- Existing handler, connection-editor, topology-base, and integration assertions were updated to use the new owners.
- All topology and tournament `createTopologyVariant_test.go` / `createClusterVariant_test.go` fixtures plus `geometricHubTopology/common_test.go` now use `test_helpers.NewGenerationTuning`.

## Tests added or updated

- Geometry covers 0/1/2 points, triangle winding, normalized deterministic edges, square determinism, collinear positions, duplicate positions, bounds, and disconnected components.
- Graph covers node allocation, symmetric linking, BFS shortest distances, disconnected graphs, and connected components.
- Zone classification covers player, neutral tiers, hub/highest, missing/unknown layouts, empty pools, unknown data, and connection guard policy.
- Profile tests cover every quality catalog plus fresh returned slices.
- Tuning factory tests cover percentage mapping and disabled optional spawns.
- `go build ./...`: passed.
- Required coverage command: passed at 64.2%, above the Phase 1 checkpoint of 64.1%.
- `go test ./test/unit/... -count=1`: passed.
- `go test ./test/... -count=1`: passed.
- `go test -tags=integration_test ./test/integration/... -count=1`: passed.
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`: passed headlessly.
- `go test -tags=integration_test ./test/performance/... -count=1`: passed; no tests to run.
- Architecture guard, workspace diagnostics, and `git diff --check`: passed.
- Independent Claude Fable 5 review: approved with no hard bugs.

## Git status snapshot

- Branch: `AD/refactoring-07-21`, up to date with `origin/AD/refactoring-07-21` before these Phase 2 edits.
- Phase 2 changes are unstaged; new owner packages and their tests are untracked.
- No files were staged or unstaged by the agent.
- `coverage.txt` was regenerated by the required coverage command.
- No changes were made under `data/`, `internal/entities/template/`, or `internal/registry/`.
- Use `git status --short` for the exact inherited listing; preserve the current staged/unstaged state unless the user explicitly requests staging.

## Rejections / things the user declined

- None.
- The first Fable review invocation used the repository shorthand model identifier, which was unavailable; it was retried with the exact registered `Claude Fable 5 (copilot)` name and completed.
- The first common-profile combined patch missed an import anchor and applied nothing; the move was split into smaller validated patches.
- The first tuning-factory location introduced a package cycle through the parent `template_generator` package. The focused check caught it immediately; the factory moved to the dependency-only `generation_tuning` subpackage.
- `goimports` was unavailable. Compiler-guided import cleanup and `gofmt` were used.
- A PowerShell import insertion initially assumed CRLF while the test files used LF, leaving missing helper imports; the same focused suites caught it, and a newline-agnostic repair passed.

## Open questions

- `geometry.GetPositionBounds` intentionally preserves the old non-empty-input precondition. It is now public; consider documenting or guarding empty input only if Phase 3 or a new caller needs that contract.
- `geometry.FindClosestAcrossComponents` faithfully compares the first component against all others because the topology loop repeatedly links that result. Do not reinterpret it as an all-pairs global search without behavior tests.
- Stateless `connection_editor` functions construct classifiers/factories per operation. Phase 3 service conversion is the appropriate place to make those constructor-owned dependencies.
- Full-variant zone count remains the manual-editor tuning convention; generation excludes an extra hub. This pre-existing semantic question is still deferred to Phase 3.
- The portal `Near` rationale remains undocumented; preserve `0.075..0.35` unless explicitly changed.

## Next recommended actions

1. Read `AGENTS.md`, this handoff, and Phase 3 in `plans/clean-architecture-refactoring.md`.
2. Ask the user to review/push the completed Phase 2 checkpoint before beginning Phase 3 if they have not already done so.
3. Start Phase 3 with the smallest cohesive-service conversion, preserving the Phase 1 GUI boundary and Phase 2 ownership rules.
4. Keep protected directories read-only and rerun build, coverage, default, tagged, and review gates.

## Carry-forward prompt

Read `AGENTS.md` first, then read `plans/clean-architecture-refactoring.md` and `.agent/session-carry-forward.md`. Phase 2 is complete and independently reviewed: position layout belongs to a topology service; geometry and graph algorithms belong to helpers; entity quality belongs to `zones.ZoneClassifier`; neutral profiles belong to `common_zones`; tuning construction belongs to `template_generator/generation_tuning`; coverage is 64.2% and all required suites pass. Never modify `data/`, `internal/entities/template/`, or `internal/registry/`; preserve Windows/Linux portability; add tests for non-trivial logic and do not reduce coverage. Preserve the zero-forbidden-import GUI architecture guard and current unstaged worktree state. Resume with Phase 3 only after the user confirms the Phase 2 checkpoint is accepted; use `.agent/session-carry-forward.md` for the full handoff.
