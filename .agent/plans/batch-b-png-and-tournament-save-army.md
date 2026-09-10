# Batch B: PNG rasterization and Tournament Save Army

Fix review §1.3 and §1.7 only: make nonzero PNG connector segments visible and honor the configured tournament Save Army value. Preserve geometry, dash policy, assets, serialization schema and all export-directory rules.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../../AGENTS.md) first. [The review](../backlog/review-gpt-6-astra-09-07.md) remains the stable finding ledger. This plan governs Batch B only; do not restart A/D or repeat the audit.

## Approval and verified starting state

- Date: 2026-09-10. **Implementation approved:** owner said "looks good, you can proceed". At implementation start the owner had staged this plan; its staged contents are preserved. Progress updates below are unstaged.
- Actual branch: `AD/visual_and_value_save_bugs`. Actual HEAD: `f986a918e7b2ed351b2dacfddeb2ba472a8303fc`, `PBI Item resolution (Batches A and D) (#41)`. Index/worktree were clean before this plan. Recheck before acting; never switch branches speculatively.
- The handoff's older branch and five unstaged-file claims are superseded by the squash. Checked A/D production paths match owner correction `509cc05`; three closure regression tests and their helper cleanup are committed. The original owner commits need not be ancestors of the squash. A/D remains closed.
- Historical Windows verification: coverage 74.5%, lint zero. These are not fresh measurements for this plan; establish a new comparable baseline before implementation.
- Owner confirmed **in-game omission of `tournamentSaveArmy` means off**. Record this as owner-supplied validation, not an assistant game test. No protected schema/tag or JSON-writer modification is needed or approved.
- Follow-up owner decision: a nil `TournamentRules` block with victory-selector tournament activation uses the getter's **false**, with no new fallback. This intentionally changes the old forced-true result for that incomplete API input. Normal default configuration explicitly sets SaveArmy=true and remains true.
- Owner approved **raster sampling only**: preserve curve geometry, subdivision counts, brush/color/width, dash length/gap/phase policy, asset placement and arena behavior. No geometry consolidation, topology, road, portal-classification, panel or default-setting changes.
- Owner requires focused **before/after PNG review** before closing §1.3. Automated checks alone do not satisfy this acceptance gate. Do not update unrelated Gio goldens.
- Never edit protected data/schema/registry trees, hand-edit generated Wire output, stage/unstage/commit/push, persist output paths or authorize browsing fallbacks. No global test tags or production test-only seams. Use portable paths. No constructor/provider change is planned, so Wire regeneration is not expected.

## Phase 1: Baseline and focused reproductions

Status: Complete

- [x] Obtain explicit owner approval of this plan; record it here.
- [x] Recheck branch/HEAD/index/worktree and preserve any owner changes.
- [x] Capture fresh Windows build, unit coverage and report-only lint baselines, recording tool versions, total statement coverage and affected-function coverage. Use Go's coverage tools, not homemade sums of repeated profile blocks.
- [x] Add failing §1.3 regression cases through `CreatePreviewImage`, using a test-local implementation of the existing `IPreviewLayoutService`. Keep fixtures and helper types under tests; introduce no production seam.
- [x] Cover short/long horizontal, vertical and diagonal portal and direct connectors, including the reported endpoints `(300,350)` / `(400,350)` with midpoint control and radius 21. Use fixed layout data for edge/no-edge comparisons, so removing an edge does not relocate the nodes. Choose fitted dominant-axis chord extents on both sides of one pixel: straight paths around 96 total dominant-axis pixels for dashed subdivision and 24 for solid subdivision. Account for diagonal projection rather than using arc length alone. Initial matrix captured; additional long-axis cases are completed in Phase 2.
- [x] Verify connector-colored changed pixels outside zone artwork, not merely portal-versus-solid inequality. Cover preserved visible dash gaps, curved segments, coincident endpoint skip behavior, reverse directions and canvas clipping through the public API. Pin central fixtures to identity asset fitting and include a near-border case with nonidentity fitting. Inspect whether zero-length raster samples are reachable through valid public layouts; do not invent a seam to force a private guard. Record genuinely unreachable remainder in the existing test observations if needed.
- [x] Add failing independent configured-true/configured-false Save Army cases for checkbox and victory-selector tournament activation. Cover inactive tournament, including configured SaveArmy=true; nil-block selector activation must use false by the explicit owner decision, while the normal constructor's explicit true default stays true.
- [x] Keep unit tests pixel-only and in-memory. Capture representative pre-fix PNGs through a separate untagged integration test, skipped when `HOMMOE_BATCH_B_PNG_REVIEW_DIR` is unset (not a new test flag). Set it only for focused capture to the absolute path of `output/research/batch-b-png/<unique-run>/before`, and later the corresponding `after` directory, using portable path construction. Existing `output/` ignore rules cover these synthetic PNG diagnostics; no ignore-file edit is needed. Do not overwrite files from other work or rely on `tmp/` surviving an active owner watcher. Reuse deterministic fixture definitions for before/after capture; any shared test helper must preserve pixel-only unit behavior. These are direct image-encoding test artifacts, never `.rmg.json` files or a new application export destination; production exports remain unchanged.

### Verification Plan: Phase 1

- Before new tests or production edits: run the existing coverage task and `go tool cover '-func=coverage.txt'`; record the actual fresh baseline, with 74.5% as historical reference only. Run `go build ./...` and report-only lint.
- Run focused new regression tests against unchanged production and record the expected failures: short portal connector absent; configured false produces true. Do not weaken assertions to fit current output.
- Check every new test follows mirrored public-method layout, descriptive names, AAA, `t.Parallel()` and testify. Integration files using only production APIs remain untagged.
- Review artifacts must not touch the game installation or change app output authorization. Repository/file-service integration tests use isolated `t.TempDir()` fixtures only.

### Phase Summary: Phase 1

Fresh baseline on unchanged production: Go 1.27.0 Windows/amd64, build pass, full coverage task pass, total 74.5%. `CreatePreviewImage` 100%, `drawConnections` 92.3%, `drawLine` 90.9%, `setTournamentRules` and config mapper 100%. Report-only golangci-lint 2.13.1: zero issues, three existing unused-exclusion warnings. Before production edits, ten PNG visibility cases failed with zero pixels and three provider false/nil cases failed; both serialized false cases failed with actual true. Mapper false coverage passed. Eighteen PNGs captured at `output/research/batch-b-png/20260910-101656/before`. The zero-length raster branch is reachable when endpoint trimming collapses all points and is tested without a seam. Test-layout passed.

## Phase 2: Minimal rasterization and value fixes

Status: Complete

- [x] Change the sample loop bound in [previewGeneratorService.go](../../internal/services/preview_service/previewGeneratorService.go) to `int(math.Ceil(steps))`; preserve floating-point increment, zero-length skip, geometry, brush and dash policy. The new stamp set is a superset of the original and never overshoots, preserving existing connector pixels.
- [x] Extend [createPreviewImage_test.go](../../test/unit/internal/services/preview_service/previewGeneratorService/createPreviewImage_test.go) with the full short/long type/orientation matrix, short/long dash-gap checks, both degenerate branches, clipping/fitting and determinism. Keep existing constructor/arena tests intact.
- [x] Assign configured `tournamentRules.SaveArmy` in [gameRulesProvider.go](../../internal/services/template_generator/providers/gameRulesProvider.go), without schedule/activation/default changes.
- [x] Replace the bug-endorsing provider test with true/false checkbox/selector cases, inactive cases, nil-rule selector false and explicit constructor-default true.
- [x] Add mapper false propagation coverage in [fromEditorState_test.go](../../test/unit/internal/mappers/generatorConfigMapper/fromEditorState_test.go).
- [x] Add [tournamentSaveArmySerialization_integration_test.go](../../test/integration/tournamentSaveArmySerialization_integration_test.go): four raw-JSON cases, actual nested structure and `tournament:true` positive controls, true presence/false absence.
- [x] Exercise `composition.InitializeGuiHandler()` generation/save and the actual legacy-omission writer in isolated `t.TempDir()` directories; no GUI state/test exports/global tags.
- [x] Capture after PNGs and obtain owner approval of the short portal, longer portal gaps and short direct-edge pairs. In-app vector preview and GUI goldens unchanged.

### Verification Plan: Phase 2

- Focused PNG/provider/mapper tests and new untagged serialization integration tests pass. Confirm true, absent-false and inactive-mode contracts independently.
- New PNG checks demonstrate visible lines, genuine gaps for dashed connections, no new clipping panic and stable repeated output. Unchanged curves/assets and existing arena tests remain valid.
- Recheck diagnostics and formatting only for explicit changed files; no bulk formatting or autofix.
- Get owner acceptance of focused before/after PNGs, recording the reviewed artifact names and outcome. No unrelated golden updates.

### Phase Summary: Phase 2

Two production lines changed. Focused regressions pass. Eighteen before/after pairs are under `output/research/batch-b-png/20260910-101656`; four later fixtures are after-only and were disclosed. Owner selected **Approved** for the linked short portal, longer portal and short direct-edge PNG comparisons on 2026-09-10. No schema, default, output authorization or Wire changes. Shared fixed layouts live in a test helper; artifact capture is a separate default-skipped untagged integration test. The zero-length and failed-trimming branches both have public-API coverage.

## Phase 3: Full verification and independent review

Status: Complete

- [x] Build and fresh full unit suite (`-count=1`) pass, including the final short-dash-gap addition.
- [x] Default `go test ./test/...`, tagged integration and headless tagged GUI tasks pass. No snapshot updates or global tags.
- [x] Vet with `integration_test`, test-layout and full report-only lint pass; zero lint issues, three pre-existing unused-exclusion warnings.
- [x] Coverage before/after: 74.5% / 74.5%. Final `drawLine`, `drawConnections`, solid/dashed paths and `setTournamentRules`: 100%. New tests cover the value/path defects despite pre-existing high statement coverage.
- [x] Independent Claude Opus 5 review: no blockers. Accepted recommendation to additionally pin gaps on the short portal; final coverage/unit/lint/layout rerun passed. Pixel preservation established by additive-stamp proof and visual spot checks, not a claimed automated 18-pair pixel diff.
- [x] Diff checks clean; owner-staged original plan preserved; all implementation/test/progress updates unstaged. Branch/HEAD unchanged.
- [x] No protected data/schema/registry, Wire or GUI-golden changes.
- [x] Results are Windows-only; no local Linux, race or benchmark outcome claimed. In-game omitted-false semantics are owner-confirmed, not assistant-tested.

### Verification Plan: Phase 3

- Build/unit/default/integration/GUI/vet/layout pass, no coverage decrease against comparable baseline, zero lint issues, no unapproved image changes, and independent review has no unresolved blocker.
- Only approved source/tests/plan/progress documentation changed. No protected data/schema/registry edits, global tag changes, output-policy changes or generated-file hand edits.

### Phase Summary: Phase 3

Verification complete on Windows Go 1.27.0. Build, full units, default suite, tagged integration, tagged GUI, vet, layout and lint pass; coverage stays 74.5%. Owner PNG acceptance is recorded in Phase 2. Remaining work is owner commit and read-only closure verification, not more implementation. Final short-gap test was added after the independent review recommendation; it passed the final full unit/coverage and lint/layout checks.

## Phase 4: Owner acceptance and closure

Status: In progress

- [ ] Summarize files, behavior, test additions, measured coverage and remaining limitations for owner review.
- [ ] Wait for owner commit. The assistant must not stage, unstage, commit, push or publish.
- [ ] Verify the owner's commit read-only, then mark review §1.3 and §1.7 fixed in place and update counts. Until then, use implemented/verified pending owner commit, not fixed.
- [ ] Complete this plan and update the existing handoff only if wrapping up or requested; preserve later backlog scope and stable identifiers.

### Verification Plan: Phase 4

- Owner approved the focused PNG comparisons and committed the scoped implementation. Commit contents match the reviewed change and verified tests; Git state is reported accurately.
- Review progress reflects only committed fixes. No deployment or release is claimed without evidence.

### Phase Summary: Phase 4

Pending phase completion.

## Final Recap

Implementation, automated verification and owner PNG acceptance complete. Await owner commit; do not mark the two backlog findings fixed yet. A/D remains closed. No release/publication performed.

## Deployment Plan

Owner reviews and commits the scoped changes, then uses the existing build/release process and runs the updated application. Newly exported PNGs use corrected sampling; newly generated tournament templates honor Save Army. Existing files are not rewritten automatically. Exports remain restricted to the detected game templates directory or explicit session-only selection; no migration, remembered path or automatic publication is part of this batch. After the owner commit, verify its contents read-only and update the two stable backlog findings and counts.
