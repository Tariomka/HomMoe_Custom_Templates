# Carry-forward: Batch B verified, awaiting owner commit

Date: 2026-09-10.

**Latest owner-review follow-up:** Batch B is committed as `252a7ee` on `AD/visual_and_value_save_bugs`. The owner excluded the untracked PNG capture utility and requested its removal rather than keeping an environment-gated test that only counts files. That utility is removed; fixture enumeration is now private and the unconditional pixel regression suite is retained. Cleanup build/default tests/fresh unit coverage/layout pass, coverage 74.5% before/after, final lint zero. Only the helper, plan and handoff have unstaged tracked changes. Index was clean at follow-up start and is untouched. Earlier pending-commit/staging descriptions below are historical; backlog closure remains separate from this focused cleanup.

## 1. Session goal

Implement approved Batch B (§1.3 PNG rasterization and §1.7 Tournament Save Army), verify it and obtain owner PNG approval. All three are complete; owner commit is pending. [The plan](plans/batch-b-png-and-tournament-save-army.md) is the resumable source of truth. [The review](backlog/review-gpt-6-astra-09-07.md) still correctly records **4 fixed, 28 remaining** until the owner commits B. A/D is closed in squash `f986a91` (PR #41); do not restart it or repeat the audit.

## 2. Fixes applied

- [PNG rasterizer](../internal/services/preview_service/previewGeneratorService.go): loop count uses `int(math.Ceil(steps))`; floating-point increment and all existing stamp positions stay unchanged. Subpixel segments now paint without changing geometry, brush, dash policy or assets.
- [Game rules](../internal/services/template_generator/providers/gameRulesProvider.go): `TournamentSaveArmy` uses `tournamentRules.SaveArmy`, rather than forced true. No protected schema or writer changes.
- Checked A/D production paths match owner correction `509cc05`; the additional three closure regressions are included in squash `f986a91`. The former handoff's unstaged A/D claims are obsolete.

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
| [previewRasterFixture.go](../test/test_helpers/previewRasterFixture.go) | New fixed-layout test helper, 22 fixtures. |
| Removed review-only PNG capture utility | Was never committed; no environment-variable setup remains in test code. |
| [tournamentSaveArmySerialization_integration_test.go](../test/integration/tournamentSaveArmySerialization_integration_test.go) | New real generation/save raw-JSON matrix with structure and tournament:true positive controls. |
| [Batch B plan](plans/batch-b-png-and-tournament-save-army.md) | Owner-staged original; unstaged progress updates. |
| [Handoff](session-carry-forward.md) | Updated current state and pending owner closure. |

## 5. Tests added or updated

Before production edits, ten PNG cases failed with zero visible connector pixels; three provider false/nil cases and both serialized-false cases failed with actual true. These regressions now pass. Test files are listed in §4.

Windows Go 1.27.0 verification:

- Build, full unit suite (`-count=1`), default `go test ./test/...`, tagged integration, tagged GUI integration, vet with `integration_test`, and test-layout: PASS.
- Fresh baseline/final unit coverage: **74.5% / 74.5%**. `drawLine`, `drawConnections`, solid/dashed methods and `setTournamentRules`: 100%.
- Final golangci-lint 2.13.1 report-only run: **0 issues**, three existing unused-exclusion warnings.
- Independent Claude Opus 5 review: no blockers. Its optional short-portal gap regression was added; final full unit/coverage/lint/layout checks passed afterward. Default/tagged/GUI checks passed before that last test-only addition.
- Pixel preservation follows the additive-stamp proof plus visual spot checks; no automated comparison of all 18 image pairs is claimed.
- No protected data/schema/registry, generated Wire or GUI golden changes. No local Linux/race/benchmark/vulnerability-scan execution claimed.

## 6. Git status snapshot

Branch: `AD/visual_and_value_save_bugs`. HEAD: `f986a918e7b2ed351b2dacfddeb2ba472a8303fc`, A/D squash PR #41. Owner staged the original Batch B plan before implementation; it is `AM` now, with its index preserved and progress updates unstaged. The five tracked Go files and this handoff are modified unstaged; the three new helper/integration files are untracked. Recheck before acting. No assistant staging, unstaging, commit, push or branch switch. No release/publication.

## 7. Rejections / things the user declined

- Retain settled A/D decisions: use `PathResolutionService`, consume pending base on every Apply, no saved output path or browsing-fallback authorization.
- No protected edits, global test tags, bulk rewrite, generated Wire hand edits or fake production seams.
- Review refined the plan to preserve floating-point stamp spacing rather than resampling long chords; unit tests stay in-memory and raw JSON checks use actual nested fields instead of substring assertions.
- Tool attempts: lowercase model name failed and was corrected; the read-only Explore agent could not implement tests, so a normal implementation subagent was used. A delegated handoff update did not match its completion report; direct replacement and readback corrected it. Verify actual files rather than completion claims.

## 8. Open questions and confirmed later scope

Batch B has no open implementation or visual decision. Owner commit is pending. After it, verify contents read-only, mark §1.3/§1.7 fixed and update counts. Then Batch C (§1.4/§1.6/§1.10) requires source verification, generated/custom/imported road policy and nil EditorState questions, plus a separately approved durable plan. Do not begin C automatically.

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

1. Owner reviews/commits Batch B, including the new untracked helper and integration tests. Preserve owner staging.
2. Verify that commit read-only, finish plan Phase 4 and mark review §1.3/§1.7 fixed with updated counts.
3. Begin Batch C only after source verification, scope questions and an approved plan. Do not repeat the audit or A/D/B implementation.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md), [the Batch B plan](plans/batch-b-png-and-tournament-save-army.md) and [the review](backlog/review-gpt-6-astra-09-07.md). Branch was `AD/visual_and_value_save_bugs`, HEAD `f986a91`; recheck Git. A/D is closed. Batch B is implemented, verified and owner PNG-approved, awaiting owner commit. Windows build/unit/default/integration/GUI/vet/layout passed; coverage 74.5% before/after, lint zero. Owner-staged original plan must be preserved; implementation and progress are unstaged, three new test/helper files untracked. Confirm owner commit before marking §1.3/§1.7 fixed. Then Batch C needs scope questions and its own approved plan.
>
> Never modify protected data/schema/registry trees; proposed protected changes require owner approval and application. Preserve Windows/Linux portable paths and guarded platform code. Test nontrivial changes and measure coverage before/after. Keep multi-step work in an approved durable plan. Never stage, unstage, commit or push; preserve owner staging and do not switch branches speculatively. Never bulk-rewrite or hand-edit generated Wire output. Export only to the detected game templates directory or explicit session-only picker destination; never persist output paths or authorize browsing fallbacks. Never set global integration_test/gui/wireinject tags or add fake unit seams. Preserve later scope in §8.
