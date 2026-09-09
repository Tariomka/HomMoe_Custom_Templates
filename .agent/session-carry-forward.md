# Carry-forward: Batch A and D complete

Date: 2026-09-09.

## 1. Session goal

Complete review §1.1, §1.2, §1.8 and §1.9, verify the owner's committed correction, and prepare the next session. [The review](backlog/review-gpt-6-astra-09-07.md) remains the backlog source of truth: **4 fixed, 28 remaining** (5 High, 16 Medium, 7 Low), from 32 original findings.

## 2. Fixes applied

- [Manual Apply](../app/gui/drivers/stateManualEdits.go): changed accepted snapshots mark dirty and reset Exit confirmation; rejected/no-template/identical applies preserve the appropriate flags. Revert identity is computed before handler mutation.
- [GUI editor state](../app/gui/models/editorState.go) and [manual-state equality](../internal/models/editor_state_model/manualEditSettings.go): deep-clone manual zones on ingress/egress and compare persistable snapshot data.
- [Startup state](../app/gui/drivers/state.go), [filesystem handler](../internal/handlers/fileSystemHandler.go), and [path resolution](../internal/services/file_system/pathResolutionService.go): failed lookup leaves export unset; only detection or explicit session folder selection authorizes a destination.
- Owner committed the batch in `a9eb35c`, then corrected untouched-revert fallthrough in `509cc05`. An accepted untouched revert now returns even if there was no manual snapshot to clear. Both commits were checked read-only.

## 3. Features added / changed

- No-op Apply does not dirty a clean document or reset an armed Exit confirmation. Creating the first manual snapshot counts as a persisted change even if visually identical to generation.
- Every Apply consumes the pending base, including rejected/no-template applies, by explicit owner decision.
- `FindGameTemplateDirectory` belongs to existing `PathResolutionService` and `IPathResolutionService`; no separate detector service remains. No path persistence, detector-algorithm changes, or browsing fallback exports.
- Closure added three regression tests and documentation only; no production code was changed after `509cc05` by the assistant.

## 4. File modifications

Committed production and test work is in `a9eb35c` and `509cc05`. Five files remain unstaged from the closure:

| File | Change |
| --- | --- |
| [applyEditedZones_test.go](../test/unit/app/gui/drivers/stateManualEdits/applyEditedZones_test.go) | Assistant-added empty-snapshot revert tests; remove always-true helper parameter reported by lint. |
| [exit_test.go](../test/unit/app/gui/drivers/stateFiles/exit_test.go) | Assistant-added empty-snapshot revert confirmation test. |
| [review-gpt-6-astra-09-07.md](backlog/review-gpt-6-astra-09-07.md) | Mark §1.8 fixed in `509cc05`; current progress count. |
| [batch-a-manual-state-and-export-directory.md](plans/batch-a-manual-state-and-export-directory.md) | Completed plan with correction verification and uncommitted test disclosure. |
| [session-carry-forward.md](session-carry-forward.md) | Current handoff. |

Protected data/schema/registry trees are unchanged. Wire was regenerated during implementation and the final graph is identical to its original form; no separate provider remains.

## 5. Tests added or updated

Three new tests prove that an untouched revert with no previous snapshot stores no snapshot, leaves a clean document saved, and preserves an armed Exit confirmation. These tests are assistant-added and still await owner review/commit. The specific empty-snapshot case is unit-tested; existing real-handler integration tests cover revert with a previous snapshot.

Verification on Go 1.27.0 Windows/amd64 after the owner's correction:

- Build, full unit suite (`-count=1`), default `go test ./test/...`, tagged integration, tagged GUI integration, vet with `integration_test`, and test-layout check: PASS.
- Focused test run: 29 passed; final focused packages also passed after test-helper cleanup.
- Coverage task before/after closure tests: **74.5% / 74.5%**, `ApplyEditedZones` 100%. Original pre-batch baseline was 74.4%.
- Final report-only lint: **0 issues**, three existing unused-exclusion warnings. An always-true test-helper parameter finding was corrected without production edits.
- No GUI goldens changed. Recovery tests export only into isolated fixture directories.
- No local Linux, full race, benchmark, or vulnerability-scan result is claimed.

Independent Claude Opus 5 review accepted the corrected branch and regression tests with no application blocker. Existing connection `Road` pointer/placement aliases remain outside the explicitly zone-only isolation scope.

## 6. Git status snapshot

Branch: `AD/save_and_pathing_silent_bug`. HEAD: `509cc05` (`Quick fix`), following `a9eb35c` (`Batch A and D done`). Index was clean and working tree clean before closure tests; the five files in §4 are now unstaged. Recheck before acting. No agent staging, unstaging, commits, pushes, or branch switches occurred. Release/publication has not been verified.

## 7. Rejections / things the user declined

- Separate detector service rejected: use `PathResolutionService`.
- Retaining a pending base after rejected Apply rejected: consume on every Apply.
- No persisted output path, unrelated fallback export, protected edits, global test tags, bulk rewrite, or generated Wire hand edits.
- The owner briefly combined revert selection with the clear-change result; closure found the fallthrough, and owner corrected it in `509cc05`.
- Do not repeat the full audit or restart completed A/D work. No Batch B implementation has started.

## 8. Open questions and confirmed later scope

Next is Batch B: §1.3 PNG rasterization and §1.7 Tournament SaveArmy. Reverify sources/tests, ask remaining scope questions, create a durable plan, and obtain approval. Confirm whether omitted false disables army saving in-game; any protected schema change is owner-approved and owner-applied only. PNG expectation changes need review; do not bulk-update unrelated GUI goldens.

Retained owner decisions for later work:

- §2.3/O08: remove Ring, Hub, Chain, Shared Web and corresponding tournament builders; reject retired saved IDs without modifying the current document; surviving tournament fallback cases use balanced generation. Preserve Geometric Hub and shared hub-zone concepts.
- §2.4/O11: investigate live state and persistence before implementation, comparing exact-base reconstruction against regenerated untouched zones plus deltas. Cover identity, randomness, upgrades, defaults, derived fields, migration and measured size/performance; owner approves feasibility first.
- §2.5/O12: structured entries following BonusEntry; bans carry SID, overrides SID/GuardValue, preserve semantics and Variant=-1. Legacy/UI text parsing stays at boundaries.
- §2.6/O13: section current groups except TemplateIdentity, MapSettings and SchemaOptions remain flat. Coordinate migration with §2.5 and any approved §2.4 outcome; retain supported legacy loading. Section keys still require planning.
- §2.7/O14: all six zone-content accessors on drivers.State, replacing panel callbacks with zone/tier selection; retain validation, dirty tracking and snapshot isolation.
- §2.8/O16: rename services/zones to zone_services and nested interfaces to zone_service_interfaces, updating imports/tests/links/generated wiring without behavior change.
- §2.9/O18: General/Layout/Bonuses named subpackages with private cohesive sections and public aliases/constructor compatibility; Preview unchanged.
- §2.10/O19: public GeneratorConfig lookup preserving Highest→Hub, unknown→nil and read-only borrowed-row semantics.
- §4.1/O20: permitted production Vec2 audit, equivalent methods, justified tested additions, numerical and hot-path behavior preserved.
- §2.2: zone-content DTO removal reopened; bonuses exception remains accepted. Replacement row/result API and full-exception versus single-seam scope still need approval.
- O09 presets remain deferred and uninvestigated.

Other pending decisions: generated/imported road policy and nil EditorState (§1.4/§1.10); arena/manual invalidation and effective mode aliases (§1.5); custom guards/preset identity (§1.11); remembering non-tournament player count (§1.12); GUI/PNG curve agreement (§2.1); tooling version/EOL/release-tag rules (§6.2/§6.3/§6.5). Review §10 retains in-game bonuses/bans, hero-hire and historical preview validation. Tidy dry-run differences previously measured were Windows checksum EOL-only, not dependency drift.

## 9. Next recommended actions

1. Owner reviews/commits the two additional regression-test edits and three documentation updates.
2. Recheck Git and read the backlog, then ask/plan/approve Batch B (§1.3/§1.7); establish a fresh baseline before changes.
3. Follow backlog §9 afterward, skipping completed A/D. Keep stable finding numbers and distinguish verified code from in-game outcomes.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [this handoff](session-carry-forward.md) and [the review](backlog/review-gpt-6-astra-09-07.md). Batch A/D (§1.1/§1.2/§1.8/§1.9) is fixed in owner commits `a9eb35c` and `509cc05`. Current branch is `AD/save_and_pathing_silent_bug`; recheck HEAD/index/worktree. Two assistant-added regression-test files and three docs are unstaged pending owner review/commit. Tests/build/GUI/vet/layout passed; Windows coverage is 74.5%, lint zero. Begin Batch B only after source verification, scope questions and an approved durable plan; §1.7 needs the in-game omitted-false decision. Do not repeat the audit or restart A/D.
>
> Never modify protected data/schema/registry trees; proposed protected changes require owner approval and owner application. Preserve Windows/Linux portable paths and guarded platform code. Test nontrivial changes and measure coverage before/after. Keep multi-step work in an approved durable plan. Never stage, unstage, commit or push; preserve owner staging and do not switch branches speculatively. Never bulk-rewrite or hand-edit generated Wire output. Export only to the detected game templates directory or explicit session-only picker destination; never persist output paths or authorize browsing fallbacks. Never set global integration_test/gui/wireinject tags or add fake unit seams. Preserve later scope in §8, including reopened zone-content DTO removal, accepted bonuses exception, investigation-first O11, and deferred O09.
