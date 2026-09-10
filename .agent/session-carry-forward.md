# Carry-forward: Start Batch C

Date: 2026-09-10.

**Next session:** Start Batch C (§1.4/§1.6/§1.10), covering road policy and graph connectivity. A/D and B are complete; they require no further validation or closure work. Use this handoff and the current review for context, not documents from completed batches.

## 1. Session goal

Prepare the next session to address Batch C without reopening completed work. [The review](backlog/review-gpt-6-astra-09-07.md) records **6 fixed, 26 remaining** (3 High, 16 Medium, 7 Low). Batch C needs source verification specific to its findings, owner decisions on road handling and approval of its implementation approach.

## 2. Fixes applied

- [PNG rasterizer](../internal/services/preview_service/previewGeneratorService.go): loop count uses `int(math.Ceil(steps))`; floating-point increment and all existing stamp positions stay unchanged. Subpixel segments now paint without changing geometry, brush, dash policy or assets.
- [Game rules](../internal/services/template_generator/providers/gameRulesProvider.go): `TournamentSaveArmy` uses `tournamentRules.SaveArmy`, rather than forced true. No protected schema or writer changes.

## 3. Features added / changed

- Owner confirmed omitted `tournamentSaveArmy` means off in-game. This is owner-supplied validation, not an assistant game test.
- Owner selected nil `TournamentRules` plus tournament victory selector => false, no new fallback. Normal constructor default explicitly sets true and is unchanged.
- Owner approved focused before/after PNG pairs: short portal, longer portal and short direct edge. Eighteen paired images and four after-only cases are under `output/research/batch-b-png/20260910-101656/{before,after}`. These are ignored synthetic image diagnostics, not game exports.
- Review-only artifact capture was removed at owner request. Pixel assertions run unconditionally and stay in-memory; future GUI flow and saved-preview golden coverage is separate work.
- In-app vector rendering, GUI goldens and all output-directory authorization/persistence behavior remain unchanged.

## 4. File modifications

| File | Change |
| --- | --- |
| [previewGeneratorService.go](../internal/services/preview_service/previewGeneratorService.go) | One-line additive sample-count fix. |
| [gameRulesProvider.go](../internal/services/template_generator/providers/gameRulesProvider.go) | One-line configured Save Army assignment. |
| [createPreviewImage_test.go](../test/unit/internal/services/preview_service/previewGeneratorService/createPreviewImage_test.go) | Raster matrix, visibility, gaps, fitting/clipping, degeneracy and determinism. |
| [createGameRules_test.go](../test/unit/internal/services/template_generator/providers/gameRulesProvider/createGameRules_test.go) | True/false, activation aliases, inactive, nil and default cases. |
| [fromEditorState_test.go](../test/unit/internal/mappers/generatorConfigMapper/fromEditorState_test.go) | False propagation. |
| [previewRasterFixture.go](../test/test_helpers/previewRasterFixture.go) | Fixed-layout test helper, 22 fixtures; enumeration privatized in owner cleanup commit. |
| Removed review-only PNG capture utility | Was never committed; no environment-variable setup remains in test code. |
| [tournamentSaveArmySerialization_integration_test.go](../test/integration/tournamentSaveArmySerialization_integration_test.go) | New real generation/save raw-JSON matrix with structure and tournament:true positive controls. |
| [Review](backlog/review-gpt-6-astra-09-07.md) | §1.3/§1.7 fixed and counts updated; closure update unstaged. |
| [Handoff](session-carry-forward.md) | Current committed state and next-session scope; closure update unstaged. |

## 5. Tests added or updated

Before production edits, ten PNG cases failed with zero visible connector pixels; three provider false/nil cases and both serialized-false cases failed with actual true. These regressions now pass. Test files are listed in §4.

Windows Go 1.27.0 verification:

- Build, full unit suite (`-count=1`), default `go test ./test/...`, tagged integration, tagged GUI integration, vet with `integration_test`, and test-layout: PASS.
- Fresh baseline/final unit coverage: **74.5% / 74.5%**. `drawLine`, `drawConnections`, solid/dashed methods and `setTournamentRules`: 100%.
- Final golangci-lint 2.13.1 report-only run: **0 issues**, three existing unused-exclusion warnings.
- Independent Claude Opus 5 review: no blockers. Its optional short-portal gap regression was added; final full unit/coverage/lint/layout checks passed afterward. Default/tagged/GUI checks passed before that last test-only addition.
- Pixel preservation follows the additive-stamp proof plus visual spot checks; no automated comparison of all 18 image pairs is claimed.
- No protected data/schema/registry, generated Wire or GUI golden changes. No local Linux/race/benchmark/vulnerability-scan execution claimed.
- After capture removal/helper cleanup: build, default `go test ./test/...`, fresh unit coverage task, layout and lint passed; 74.5% before/after. Final focused pixel tests after helper function ordering: 36 passed. GUI was not rerun for removal of a non-GUI capture utility. Owner commits were checked read-only; this final wrap-up changes documentation only.

## 6. Git status snapshot

Last observed branch: `AD/visual_and_value_save_bugs`. HEAD: `a4c5fa3a05184b43f2ef1e3c0a7d4a55b74c55d2`. Implementation is committed; documentation-only edits were pending at wrap-up. Inspect current Git state before Batch C solely to preserve owner changes and staging, not to revalidate completed batches. Never stage, unstage, commit, push or switch branches speculatively.

## 7. Rejections / things the user declined

- Retain settled A/D decisions: use `PathResolutionService`, consume pending base on every Apply, no saved output path or browsing-fallback authorization.
- No protected edits, global test tags, bulk rewrite, generated Wire hand edits or fake production seams.
- Owner rejected retaining the environment-gated PNG capture test: it only checked errors/file count, not image correctness. It was removed entirely, along with capture-only setup; fixture enumeration is private. Retain real unconditional pixel assertions. Future GUI flow and saved-preview goldens are separate work, not implemented here.
- Retain floating-point stamp spacing rather than resampling long chords; unit tests stay in-memory and raw JSON checks use actual nested fields instead of substring assertions.
- Carry-forward documents must be self-contained and forward-looking. Do not reference completed-batch planning documents or assign previous-batch checking, validation or closure to the next session. Completed-batch documents are disposable; unresolved corrections normally stay in the same session.

## 8. Open questions and confirmed later scope

Batch C (§1.4/§1.6/§1.10) is not started: verify its current sources and tests, ask about generated/custom/imported road handling and optional nil EditorState behavior, then agree the scope and obtain approval before implementation. Preserve the accepted impossible-isolation connectivity fallback while keeping roads independent from graph connectivity. No previous-batch follow-up is required.

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

1. Read [AGENTS.md](../AGENTS.md), this handoff and review §1.4/§1.6/§1.10. Inspect current Git state to preserve owner changes.
2. Trace Batch C's road rebuild, connectivity repair and foothold-content paths, their callers and existing tests. Do not repeat the full audit.
3. Resolve generated versus custom/imported road handling when roads are disabled, stale foothold-target handling and behavior when `EditorState` is nil. Preserve required graph repairs independently of road creation.
4. Confirm Batch C scope with the owner, prepare its new durable implementation plan and obtain approval. Then capture a fresh coverage baseline and implement with dedicated regression coverage.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md) and [the review](backlog/review-gpt-6-astra-09-07.md). Start Batch C (§1.4/§1.6/§1.10): verify its current source paths and tests, ask about generated/custom/imported road policy and nil EditorState behavior, confirm scope and obtain approval before implementation. Preserve graph connectivity repair independently of road creation. Last observed Windows coverage was 74.5%, lint zero; measure a fresh baseline for C. Inspect current Git state to preserve owner changes. A/D and B are complete and require no validation or closure work. This handoff contains the required carry-forward decisions; do not seek documents from completed batches.
>
> Never modify protected data/schema/registry trees; proposed protected changes require owner approval and application. Preserve Windows/Linux portable paths and guarded platform code. Test nontrivial changes and measure coverage before/after. Keep multi-step work in an approved durable plan. Never stage, unstage, commit or push; preserve owner staging and do not switch branches speculatively. Never bulk-rewrite or hand-edit generated Wire output. Export only to the detected game templates directory or explicit session-only picker destination; never persist output paths or authorize browsing fallbacks. Never set global integration_test/gui/wireinject tags or add fake unit seams. Preserve later scope in §8.
