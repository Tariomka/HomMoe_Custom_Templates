# Session Carry-Forward

## Session goal

Complete Phase 6 of `plans/clean-architecture-refactoring.md`: standardize invariant-rich entity construction, extract positioned topology composition, and centralize topology capabilities without changing generated behavior or reducing coverage.

## Fixes applied

- Added builder options for runtime `Connection.IsUserAdded` and mandatory-content include lists, then migrated direct production construction in editor, repair, tournament, and mapper paths.
- Added an AST architecture guard for explicit and elided builder-owned entity literals, including slice, array, map, pointer, and nested forms.
- Extracted positioned topology assembly and removed Circles' duplicate pipeline through a typed zone decorator.
- Unified unknown topology descriptor, capability, preview, and provider fallback on Ring.
- Restored explicit provider service dispatch so future hub-capable descriptors cannot silently select the wrong generator.
- Moved zone-editor preview aggregate construction behind the handler boundary to preserve `app/gui` dependency rules.

## Features added / changed

- Added passive `TopologyCapabilities` and `TopologyLayoutKind` models and registered immutable capabilities with topology descriptors.
- Added `PositionedTopologyBuilder`, `PositionedTopologyLayoutBuilder`, and `PositionedTopologyZoneDecorator`.
- Random, Circles, Square, Geometric, Cross, and Fractal use the shared positioned pipeline for zones, connections, portals, isolation repair, connectivity repair, and final variants.
- Circles decorates player rings as `0` and neutral rings with the domain's one-based plan tier.
- Geometric Hub intentionally remains on its custom `TopologyBase` graph because central-hub portal topology does not fit the generic pair pipeline.
- Preview layout classification now consumes descriptor `LayoutKind`/ring metadata; unknown values preserve historical Ring behavior.
- `GeneratorConfig.IsHubCityToHold` remains an explicit Hub/GeometricHub check because importing the common descriptor catalog into `config` would create a package cycle.

## File modifications

Edited production files:

- `app/gui/dialogs/zoneEditorCanvas.go` - sends editable zones/connections through the preview handler instead of constructing a variant in the GUI.
- `app/gui/panels/layoutPanelTopology.go` - consumes shared topology capabilities.
- `app/gui/panels/layoutPanelZones.go` - consumes shared topology capabilities.
- `internal/common/common_topologies/topologies.go` - registers capabilities and Ring fallback.
- `internal/dtos/previewLayoutRequestDto.go` - supports editor zone/connection preview input with documented precedence.
- `internal/handlers/previewHandler.go` - assembles editor preview variants through `VariantBuilder`.
- `internal/mappers/mandatoryContentItemMapper.go` - uses the mandatory-content builder for SID/include-list construction.
- `internal/models/config/config_inner/mapTopology.go` - removes the old topology helper.
- `internal/models/topologyDescriptor.go` - carries topology capabilities.
- `internal/services/builders/mandatory_content/mandatoryContentItemBuilder.go` - adds include-list construction.
- `internal/services/builders/variant_content/connectionBuilder.go` - adds runtime user-added state construction.
- `internal/services/connection_editor/connectionEditorService.go` - uses `ConnectionBuilder`.
- `internal/services/preview_service/layoutGeometry.go` - removes duplicate topology classification helpers.
- `internal/services/preview_service/previewLayoutService.go` - dispatches through shared capabilities.
- `internal/services/template_generator/providers/mandatoryContentProvider.go` - consumes shared hub capability metadata.
- `internal/services/template_generator/providers/topology/base/topologyBase.go` - builds repair connections through `ConnectionBuilder`.
- `internal/services/template_generator/providers/topology/circlesTopology.go` - uses shared positioned assembly and ring decoration.
- `internal/services/template_generator/providers/topology/crossTopology.go` - uses shared positioned assembly.
- `internal/services/template_generator/providers/topology/fractalTopology.go` - uses shared positioned assembly.
- `internal/services/template_generator/providers/topology/geometricTopology.go` - uses shared positioned assembly.
- `internal/services/template_generator/providers/topology/randomTopology.go` - delegates positioned assembly.
- `internal/services/template_generator/providers/topology/squareTopology.go` - uses shared positioned assembly.
- `internal/services/template_generator/providers/topology/tournament_variant/balancedClusterService.go` - uses `ConnectionBuilder`.
- `internal/services/template_generator/providers/topologyProvider.go` - keeps explicit algorithm dispatch and Ring fallback.
- `internal/services/zones/zoneLabelProvider.go` - corrects the hub-capable hold-city comment.

Created production files:

- `internal/models/topologyCapabilities.go`
- `internal/models/topologyLayoutKind.go`
- `internal/services/template_generator/providers/topology/positionedTopologyBuilder.go`
- `internal/services/template_generator/providers/topology/positionedTopologyLayoutBuilder.go`
- `internal/services/template_generator/providers/topology/positionedTopologyZoneDecorator.go`

Deleted production file:

- `internal/services/template_generator/providers/topology/layoutFunc.go` - replaced by the named positioned layout callback.

Edited or created tests:

- `test/unit/architecture/construction/construction_test.go`
- `test/unit/internal/common/common_topologies/topologies/getTopologyCapabilities_test.go`
- `test/unit/internal/common/common_topologies/topologies/getTopologyDescriptorFromType_test.go`
- `test/unit/internal/handlers/guiHandler/buildPreviewLayout_test.go`
- `test/unit/internal/services/builders/mandatory_content/mandatoryContentItemBuilder/withIncludeList_test.go`
- `test/unit/internal/services/builders/variant_content/connectionBuilder/withIsUserAdded_test.go`
- `test/unit/internal/services/preview_service/previewLayoutService/buildPreviewLayout_test.go`
- `test/unit/internal/services/template_generator/providers/topology/circlesTopology/createTopologyVariant_test.go`
- `test/unit/internal/services/template_generator/providers/topology/positionedTopologyBuilder/buildVariant_test.go`
- `test/unit/internal/services/template_generator/providers/topology/positionedTopologyBuilder/newPositionedTopologyBuilder_test.go`
- `test/unit/internal/services/template_generator/providers/topology/positionedTopologyBuilder/newPositionedTopologyBuilderWithCreationServices_test.go`
- `test/unit/internal/services/template_generator/providers/topologyProvider/createTopologyVariant_test.go`

Plan and handoff:

- `plans/clean-architecture-refactoring.md` - Phase 6 marked Complete with decisions and verification.
- `.agent/session-carry-forward.md` - replaced with this Phase 6 handoff.

## Tests added or updated

Latest verified state after final Opus review repairs:

- `go build ./...`: passed.
- `go test -count=1 ./test/unit/...`: passed.
- `go test -count=1 ./test/...`: passed without integration tags.
- Coverage command: passed at 64.9%, equal to the 64.9% Phase 6 baseline.
- Topology/provider unit packages with `-count=20`: passed.
- `go test -tags=integration_test -count=1 ./test/integration/...`: passed.
- `go test -tags='integration_test,gui' -count=1 ./test/integration/gui/...`: passed.
- Headed `BenchmarkEditorWindow_TabCycling`, 20 iterations: passed at about 5.55 ms/op, 2.92 MB/op, 8,964 allocs/op.
- Focused topology-provider test lint: zero issues.
- Diagnostics on changed code, protected-path check, and `git diff --check`: passed.
- Broad lint still reports pre-existing debt in unrelated GUI/tests and the existing 61-line `CreateMissingConnections`; no introduced lint issue remains.
- Claude Opus 5 final review reported no blocking findings. Its medium provider-dispatch, AST-guard bypass, and Circles parity findings were fixed and focused/full tests rerun.

## Git status snapshot

- Branch: `AD/refactoring-07-21`, synchronized with its origin before Phase 6 began.
- All Phase 6 changes are unstaged; new model/builder/test files are untracked.
- No files were staged, committed, reverted, or stashed by the agent.
- `coverage.txt` was regenerated by the required coverage checks.
- No change exists under `data/`, `internal/entities/template/`, or `internal/registry/`.
- The exact final `git status --short` listing is the production/test file list above plus the modified plan; run it again before Phase 7 to detect external changes.

## Rejections / things the user declined

- No user proposal was declined in Phase 6.
- Capability-gated provider service dispatch was rejected after review because `UsesHub` describes generated content, not a unique generator algorithm; explicit service cases are safer.
- Reusing `PositionedTopologyBuilder` for Geometric Hub was rejected because its central hub, hub portals, and custom graph require a distinct assembly pipeline.
- Moving `IsHubCityToHold` onto the shared catalog was rejected because `config` importing `common_topologies` creates a package cycle.
- Removing `UsesGeneratorPosition` metadata was deferred; it remains passive descriptor information even though preview dispatch now trusts `LayoutKind` plus runtime position presence.
- Existing unrelated repository-wide lint debt was not changed.

## Open questions

- Phase 6 has no blocker and is marked Complete.
- `GetTopologyDescriptorFromIndex` intentionally retains its historical out-of-range Random fallback while type/capability/provider fallbacks use Ring; revisit only with a user-visible selection behavior task.
- The construction guard scans production `app` and `internal` imports of the root `internal/entities` package. There are currently no direct schema-subpackage constructions; extending it to arbitrary schema subpackage aliases is a future hardening option.
- Phase 7 should reassess `GUIHandler` and `TopologyBase`; do not pull Phase 8 compatibility deletion forward without checking callers.

## Next recommended actions

1. Read `AGENTS.md`, Phase 7 in the plan, and this handoff.
2. Run `git status --short` and preserve or commit the complete unstaged Phase 6 worktree as directed by the user.
3. Record a fresh Phase 7 coverage baseline before editing.
4. Start Phase 7 with the `TopologyBase` and `GUIHandler` coordinator reassessment; keep GUI domain decisions behind handlers.
5. Preserve explicit topology algorithm dispatch and the protected-path rules.

## Carry-forward prompt

Read `AGENTS.md` first, then read `plans/clean-architecture-refactoring.md` and `.agent/session-carry-forward.md`. Never modify `data/`, `internal/entities/template/`, or `internal/registry/`; preserve Windows/Linux portability; add tests for non-trivial logic and do not reduce coverage. Phase 6 is Complete and Opus-reviewed at 64.9% coverage: builder-owned entity construction is enforced, positioned topology assembly is shared across compatible topologies, topology capabilities drive preview/GUI/content behavior, unknown topologies fall back to Ring, and explicit provider cases select generator algorithms. The full Phase 6 worktree remains unstaged. Resume at Phase 7 only after checking status and recording its baseline; use `.agent/session-carry-forward.md` for the full handoff.
