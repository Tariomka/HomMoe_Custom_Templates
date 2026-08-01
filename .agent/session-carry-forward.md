# Session Carry-Forward

## Session goal

Complete Phase 4 of `plans/clean-architecture-refactoring.md`: consolidate common passive types and distance catalogs without changing generated output, persisted state, or coverage.

## Fixes applied

- Moved passive catalog types from `internal/common` to [internal/models](../internal/models).
- Privatized map-size slices behind copy-returning getters in [internal/common/mapSizes.go](../internal/common/mapSizes.go).
- Added one context-aware distance owner in [internal/common/common_distances/distancePresets.go](../internal/common/common_distances/distancePresets.go).
- Removed duplicate placement/content distance types and the unused `app/gui/constants/roadDistances.go`.
- Preserved content `Near = 0.1..0.25` and portal `Near = 0.075..0.35`.

## Features added / changed

- Added passive `DistancePreset`, `MapSize`, `GuardStrength`, `GuardWeeklyIncrement`, `TopologyDescriptor`, and `TopologyDescriptors` models.
- `ContentRuleService` consumes content distance presets; `PlacementRuleBuilder` consumes portal placement presets.
- Added `PlacementRuleBuilder.BuildNearCrossroadsRule` for topology portal rules.
- Map-size and distance catalog getters return fresh slices or construct fresh values.
- Production usage confirmed `contentIds.go`, `includeListIds.go`, and `contentItemGroups.go` remain GUI-only; themes were unchanged.

## File modifications

Created production files:

- `internal/common/common_distances/distancePresets.go`
- `internal/models/distancePreset.go`
- `internal/models/guardStrength.go`
- `internal/models/guardWeeklyIncrement.go`
- `internal/models/mapSize.go`
- `internal/models/topologyDescriptor.go`
- `internal/models/topologyDescriptors.go`

Edited production files:

- `app/gui/constants/mapSizes.go`
- `app/gui/panels/generalPanel.go`
- `app/gui/panels/layoutPanel.go`
- `internal/common/common_connections/guardStrength.go`
- `internal/common/common_connections/guardWeeklyIncrement.go`
- `internal/common/common_topologies/topologies.go`
- `internal/common/mapSizes.go`
- `internal/services/builders/placement_rule/placementRuleBuilder.go`
- `internal/services/content_rules/distanceCatalog.go`
- `internal/services/content_rules/ruleDistanceToRoad.go`
- `internal/services/content_rules/ruleDistanceToTown.go`
- `internal/services/template_generator/providers/mandatoryContentProvider.go`
- `internal/services/template_generator/providers/topology/base/topologyBase.go`
- `internal/services/template_generator/providers/topology/geometricHubTopology.go`

Deleted production files:

- `app/gui/constants/roadDistances.go`
- `internal/services/builders/placement_rule/distance.go`
- `internal/services/content_rules/distanceVariation.go`

Plan and handoff:

- `plans/clean-architecture-refactoring.md` - Phase 4 marked Complete with verification summary.
- `.agent/session-carry-forward.md` - replaced with this Phase 4 handoff.

Test changes:

- Created `test/unit/internal/common/common_distances/distancePresets/` with four public-method test files.
- Created `test/unit/internal/common/mapSizes/` with three implementation-matching test files.
- Created `test/unit/internal/services/builders/placement_rule/placementRuleBuilder/buildNearCrossroadsRule_test.go`.
- Updated all tests under `test/unit/internal/common/common_connections/guardStrength/` and `guardWeeklyIncrement/getGuardWeeklyIncrements_test.go` for model-owned types.
- Updated `test/unit/app/gui/constants/mapSizes/getMapSize_test.go` and `getMapSizes_test.go` for getter-only catalogs.
- Updated placement builder tests under `test/unit/internal/services/builders/placement_rule/placementRuleBuilder/` for `models.DistancePreset`.
- Updated content-rule tests under `test/unit/internal/services/content_rules/{contentRuleService,distanceCatalog,ruleDistanceToRoad,ruleDistanceToTown}/` for `models.DistancePreset`.
- Deleted old `test/unit/internal/constants/mapSizes/` tests after moving them to the implementation-matching common path.
- Deleted `test/unit/internal/services/builders/placement_rule/distance/tryGetDistanceFrom_test.go` with the duplicate lookup owner.

## Tests added or updated

Latest verified state after the final no-global distance refactor:

- Unit coverage command: passed at 64.7%, equal to the Phase 4 baseline.
- Complete unit suite: passed.
- `go build ./...`: passed.
- `go test ./test/... -count=1`: passed.
- Tagged integration: passed.
- Tagged headless GUI integration: passed.
- Tagged performance packages: passed.
- Focused new distance package lint: zero issues.
- Diagnostics, architecture guard, `git diff --check`, and protected-path status check: passed.
- Broader touched-package lint reported only pre-existing global/style findings outside the new distance package.
- Claude Fable 5 approved Phase 4 with no hard findings.

## Git status snapshot

- Branch: `AD/refactoring-07-21`, up to date with `origin/AD/refactoring-07-21` at the final check.
- All Phase 4 changes are unstaged; new model/catalog/test files are untracked.
- No files were staged, committed, reverted, or stashed by the agent.
- `coverage.txt` was regenerated by the required coverage checks.
- No change appears under `data/`, `internal/entities/template/`, or `internal/registry/`.
- Use `git status --short` for the exact inherited listing before Phase 5.

## Rejections / things the user declined

- No new user rejection occurred in Phase 4.
- The existing decision to preserve both different `Near` bounds was honored; no product behavior change was proposed or applied.
- One exploration subagent call used an unavailable model spelling and made no changes; it was retried successfully.
- One combined test patch failed on stale context and applied nothing; it was reapplied with exact local context.
- One PowerShell `gofmt` command passed wildcard literals and failed without changing files; expanded paths then formatted successfully.

## Open questions

- Phase 4 has no blocker and is marked Complete.
- `BuildNearCastleRule` remains production-unused, as before Phase 4; Phase 8 cleanup owns removal of test-only production API.
- Existing internal tests import GUI `ContentIDs`; production ownership remains correct, and any test-only decoupling is outside Phase 4.
- Existing mutable GUI `GameModes` and pre-existing common private globals were not introduced by Phase 4 and remain outside its touched catalog contract.

## Next recommended actions

1. Read `AGENTS.md`, Phase 5 in the plan, and this handoff.
2. Recheck `git status --short` and preserve the unstaged Phase 4 worktree.
3. Record a fresh Phase 5 coverage baseline before editing.
4. Add the Phase 5 neutral/hub characterization required immediately before that phase.
5. Consolidate neutral-like zone creation without changing hub castles, mandatory content, roads, naming, or manual reapply behavior.

## Carry-forward prompt

Read `AGENTS.md` first, then read `plans/clean-architecture-refactoring.md` and `.agent/session-carry-forward.md`. Never modify `data/`, `internal/entities/template/`, or `internal/registry/`; preserve Windows/Linux portability; add tests for non-trivial logic and do not reduce coverage. Phase 4 is Complete and Fable-approved at 64.7% coverage: common packages declare no public structs, map-size catalogs return copies, one context-aware distance owner preserves content `Near = 0.1..0.25` and portal `Near = 0.075..0.35`, and GUI content/include/group catalogs remain presentation-only. Resume from the current unstaged worktree at Phase 5 only after checking status and recording its baseline. Use `.agent/session-carry-forward.md` for the full handoff.