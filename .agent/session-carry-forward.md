# Carry-forward: begin repository-review backlog work

Date: 2026-09-08.

## 1. Session goal

Review the repository according to [the review prompt](promt_templates/review-prompt.md), then prepare a new session to implement the resulting backlog under owner approval.

**Source of truth:** [review-gpt-6-astra-09-07.md](backlog/review-gpt-6-astra-09-07.md). It contains **32 actionable items: 8 High, 17 Medium, 7 Low**, including the owner-requested work scoped on 2026-09-08, prior-item dispositions, source evidence, proposed tests, decisions, execution order (§9), and measured baselines (§11).

**Implementation has not started.** Creating this handoff is not approval to implement all findings. Start with review §1.1, then §1.2 (batch A), after confirming scope and the remaining decisions.

## 2. Fixes applied

- No application, test, configuration, dependency, or protected-data fixes were applied during the review.
- This [handoff](session-carry-forward.md) replaces obsolete operational instructions with the current starting point. It partially addresses review §7.2, but that item also covers other stale observations; do not mark the entire item fixed automatically.

## 3. Features added / changed

- Added the [review backlog](backlog/review-gpt-6-astra-09-07.md), with stable finding numbers and per-item fix/test plans.
- The owner selected **“Include verified owner items”** when asked how to treat the existing manual backlog.
- The owner explicitly selected **“Reopen the DTO removal request”** after being shown the conflict between the zone-content owner backlog and the architecture gate's accepted exception. This is review **§2.2**. The bonuses DTO exception is **not** reopened. The exact replacement API/result and whole-exception versus single-seam scope still need approval.
- No application features were changed.
- Owner-scoped additions are §2.3–§2.10 and §4.1. They cover topology retirement, compact-state investigation, typed content lists, sectioned persistence, driver accessors, package naming, panel sections, config lookup, and vector arithmetic. O09 presets remain deferred, not actionable.

## 4. File modifications

| File | State / purpose |
| --- | --- |
| [review-gpt-6-astra-09-07.md](backlog/review-gpt-6-astra-09-07.md) | Curated backlog with nine owner-scoped additions, stable IDs, acceptance criteria, and execution dependencies. |
| [owner_findings.md](backlog/owner_findings.md) | Completed O01/T01/T04/T06 source entries removed; their evidence remains in review §0. |
| [test_observations.md](backlog/test_observations.md) | Completed SaveTo observation removed as requested; its review disposition remains. |
| [session-carry-forward.md](session-carry-forward.md) | Updated to the current backlog scope and starting repository state. |

A temporary verification program under .agent/review_probe was created, executed, and deleted during the review. Its absence was checked again before this handoff. No scratch source remains. Existing root coverage reports were not overwritten; diagnostic output was directed outside the repository and is not required to resume because results are recorded in the review.

No files in the protected game-data/schema/registry trees were edited. No implementation plan exists for this backlog yet; create one only after the first batch's scope is confirmed.

## 5. Tests added or updated

**None.** The review proposes regression tests but does not implement them. Findings must be reverified on the current source before fixing; line citations describe the reviewed revision and may move.

Last measured verification, on revision `f4f4cf63f22e84040754a7231b4d1dc793070af1`:

| Check | Result |
| --- | --- |
| Go toolchain | `go1.27.0 windows/amd64` |
| Installed golangci-lint | `2.13.1`, built with Go `1.27.0`; tools module still pins `2.12.2` (§6.2) |
| `go build ./...` | PASS |
| `go vet -tags=integration_test ./...` | PASS |
| `go test ./test/... -count=1` | PASS, 183 packages, including 181 unit packages |
| `go test -tags=integration_test ./test/integration/... ./test/performance/... -count=1` | PASS; performance package reported no tests to run, not a benchmark result |
| Unit coverage with `-coverpkg=./internal/...,./app/...` | PASS, **74.4%** total statements from Go's own `cover -func` output |
| Full configured report-only lint | **0 issues**, three unused-exclusion warnings |
| `go run ./cmd/testlayoutcheck .` | PASS |
| Geometry unit, PNG unit, and untagged integration packages repeated with `-count=20` | PASS |
| Root and tools `go mod tidy -diff` | Both exit 1 solely from Windows checksum-file EOL differences; normalized content identical (§6.3) |
| Local full race / GUI / Linux / benchmarks | Not run in this review |
| Local govulncheck | Not installed/not run; PR and scheduled CI scan configuration verified, remote outcome unknown |

The review's complete coverage inventory has 282 instrumented files: 192 at 100%, 52 partial, 38 at 0%. Some GUI and other-platform files are outside that profile, not implicitly covered. Use **74.4% and zero lint issues** as the comparable Windows no-regression baseline; remeasure before implementation. Do not substitute historical 74.3%, a memory-only floor, or CI's configured 60% minimum for the current baseline.

Public-API experiments confirmed:

- §1.3: a portal with centers 100 pixels apart rendered exactly the no-edge background; Direct with the same geometry did not.
- §1.4: manual road reconstruction added roads to roadless zones.
- §1.7: `SaveArmy=false` produced `TournamentSaveArmy=true`.
- §1.9: mutating an input to `SetManualEdits` or a result of `GetManualZones` mutated stored nested position data.
- §1.13: 200 identical symmetric-obstacle geometry calls produced two different control-point Y values (300 and 400).

Other findings are source-traced, not represented as experimentally reproduced. In particular no full GUI reproduction of the revert or batched-pointer findings was run.

## 6. Git status snapshot

Checked before the 2026-09-08 documentation curation:

- Branch: **master**.
- HEAD: **6ec9e2c**, “Added New Backlog of items (#39)”, tagged **v0.3.8**.
- Local `origin/master` and `origin/HEAD` refs point at the same commit after the owner's pull. This is not a fresh remote query by the agent.
- Working tree and index were clean before these documentation edits.
- The four documents in §4 are the intended unstaged changes. No application source was changed.

A new session must inspect both index and working tree again. The agent did not stage, unstage, commit, push, or switch branches. Do not alter owner staging.

## 7. Rejections / things declined or corrected

- Review-only scope was enforced: no drive-by fixes, formatting, dependency tidy writes, or protected-data changes.
- Do not claim the portal issue fixed merely because a dashed branch exists and a test says its image differs from solid. Missing pixels satisfy that weak assertion; the public-API probe disproved the initial superficial assessment.
- Do not report accepted DTO/model layering as a fresh breach. Only the zone-content exception was explicitly reopened by the owner in this session.
- Do not propose output-path persistence, a different default export directory, live-pointer state reads, blanket test tags, or a flaky global Gio allocation threshold. Those conflict with hard rules or settled decisions.
- Proposed stale-zone-name control and generic concurrent-preview production-race findings were excluded for lack of a current demonstrated failure path. The dialog resets the zone property sync marker on addition; current preview calls are serialized.
- The first requested subagent model alias was unavailable; the display-name model worked. Use available model names, such as `Claude Opus 5 (copilot)`, for final plan/implementation review. Never use Haiku.
- A naïve custom coverage aggregation was discarded because profiles contain duplicate blocks across executables. Use Go's own total and per-file HTML percentages, not sums/averages of raw profile rows.
- Tidy failures were investigated and classified as EOL-only, not dependency drift. No files were normalized during review.

## 8. Open questions and confirmed scope

### First work: batch A, §1.1 and §1.2

1. **§1.1 no-op behavior:** should applying an identical manual layout leave a clean document clean? Recommended: mark dirty only when persisted manual state actually changes; rejected applies remain unchanged. Confirm before planning.
2. **§1.1 confirmation reset:** after any committed manual change, reset a previously armed Exit confirmation as scalar edits already do. Include clearing a manual snapshot via Revert-to-Base.
3. **§1.2 detection seam:** leaving an undetected output path empty follows the existing hard rule, but choose the narrowest real composition seam for deterministic failure tests if necessary. Do not add private test exports or persist the result. Ensure the folder picker can start browsing from an empty output path without authorizing that browsing directory as the export destination.
4. **Batch/branch scope:** confirm whether the owner wants both A items together or §1.1 first, and follow their branch preference. Do not switch branches speculatively while owner-staged changes exist.

### Later work

- §1.4/§1.10: generated versus imported/custom roads and optional nil editor-state semantics.
- §1.5: clearing incompatible manual edits versus reapplying arena policy; compare effective modes including victory-condition aliases.
- §1.7: whether omitted false means tournament army saving is disabled in the game; any required protected schema edit is owner-only.
- §1.11: custom guard numbers and ambiguous preset identity during quality reprofile.
- §1.12: whether the previous non-tournament player count should be remembered in session view state.
- §2.1: intentional editor/PNG curve differences versus required visual agreement.
- §2.2: `ContentRuleRow` versus the owner note's `ZoneContentRow`, the service result validity contract, and full DTO exception removal versus composition-only scope.
- §6.2/§6.3/§6.5: linter version, narrow EOL normalization, and release tag validation format.
- Review §10 retains in-game validation for bonuses/bans, hero-hire policy, and historical preview artifacts. No new game result was obtained.

### Confirmed owner scope for the additions

- **§2.3 / O08:** completely remove Ring, Hub, Chain, Shared Web and corresponding tournament builders. Reject saved selections of retired IDs with a clear error, without changing the current document. Route surviving tournament cases formerly using the chain fallback to the balanced builder. Preserve Geometric Hub and shared hub-zone concepts.
- **§2.4 / O11:** first investigate live state and persistence, comparing exact-base reconstruction with regenerated untouched zones plus deltas. Cover identities, randomness, upgrades, derived fields, defaults, migration, and measured size/performance. Owner approves the feasibility outcome before either representation is implemented.
- **§2.5 / O12:** follow `BonusEntry` with structured entity/model entries throughout, ban `SID` and override `SID`/`GuardValue`; preserve current semantics and output `Variant=-1`. Keep legacy/UI text parsing at boundaries.
- **§2.6 / O13:** use current groups as JSON sections except `TemplateIdentity`, `MapSettings`, and `SchemaOptions`, whose fields remain flat. One coordinated migration with §2.5 and any approved §2.4 outcome; retain supported old-state loading. Exact section keys belong in the implementation plan.
- **§2.7 / O14:** place all six zone-content getters/setters on `app/gui/drivers.State`, replacing panel callbacks with a zone/tier selection. Preserve validated dirty-tracked updates and snapshot isolation.
- **§2.8 / O16:** rename to internal/services/zone_services and its nested zone_service_interfaces package; update imports, mirrored tests, active links, and generated wiring without changing behavior.
- **§2.9 / O18:** move General/Layout/Bonuses to their named panel subpackages, with a private struct per cohesive section and public top-level aliases/constructor compatibility. Preserve GUI behavior; Preview remains unchanged.
- **§2.10 / O19:** move the lookup to a public `GeneratorConfig` method; preserve Highest→Hub, unknown→nil, and read-only borrowed-row semantics.
- **§4.1 / O20:** audit all permitted production `Vec2` use; prefer equivalent methods and add missing ones only for justified callers with tests. Preserve numerical behavior and avoid hot-path regressions.
- **O09:** not scheduled. Presets have not been investigated or planned.

## 9. Next recommended actions

1. Read [AGENTS.md](../AGENTS.md), this handoff, then review §0, §1.1–§1.2, §9, and §11. Consult [.agent/memories](memories) for settled decisions but prefer current source and the newly recorded owner decision on §2.2.
2. Inspect Git status and the staged diff without changing either. Check for edits made since the reviewed SHA. Do not rerun the whole repository review by default.
3. Re-read each target plus its callers and tests. For §1.1 begin with [stateManualEdits.go](../app/gui/drivers/stateManualEdits.go), [state.go](../app/gui/drivers/state.go), [stateFiles.go](../app/gui/drivers/stateFiles.go), and [applyEditedZones_test.go](../test/unit/app/gui/drivers/stateManualEdits/applyEditedZones_test.go).
4. Ask the concrete batch-A questions above, summarize the agreed scope, then write a durable plan under .agent/plans using the AGENTS.md template. Obtain owner approval before implementation. Use a suitable independent review of the plan.
5. Establish fresh pre-change coverage/tests. Add focused failing regressions through production APIs, then implement only approved behavior. Keep rendering in GUI and domain policy behind handlers; use existing clone/conversion helpers.
6. Verify build, unit tests, coverage, full configured lint, and test layout. Run the relevant integration/GUI suites explicitly when required. Regenerate Wire for changed constructors/provider sets; never edit generated code manually.
7. Record files, exact checks, results, and any residual limits in the plan. Present the diff to the owner. **The owner stages/commits.** Mark the review item `✅ FIXED` only after the protocol's owner-commit step is satisfied; until then record “implemented/verified, awaiting owner commit” without misrepresenting completion.
8. Continue with review §9: B (§1.3/§1.7), C roads/graph, D manual lifecycle, E modes/guards, F geometry, G performance, H (§6.2/§6.3/§6.5), I docs, J DTO work; then K topology retirement, L compact-state investigation, M coordinated persistence, N driver/panel organization, O naming/config lookup, P vector arithmetic. L may start independently but gates M's representation decisions. Coordinate overlapping files; do not renumber items or implement an entire category without plan approval.

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first, then [.agent/session-carry-forward.md](session-carry-forward.md) for the full handoff and [.agent/backlog/review-gpt-6-astra-09-07.md](backlog/review-gpt-6-astra-09-07.md) as the backlog source of truth.
>
> Begin work on the review backlog, starting with batch A: §1.1 manual edits fail to mark the document unsaved, and §1.2 failed game-directory detection falls back to an invalid export destination. First inspect current Git state and reverify source/callers/tests, ask the open scope questions, and obtain approval for a durable plan before implementing. No application fixes have been made yet. The review has 32 actionable items, with evidence/test plans and prior-item dispositions. Its nine owner additions are scoped in handoff §8; O11 is investigation-first and O09 presets are deferred. Do not repeat the full audit or mark anything fixed prematurely.
>
> Hard rules: never modify the protected data, template_entity schema, or registry trees; protected changes require explicit owner approval and owner application. Preserve Windows/Linux compatibility using portable paths and guarded platform code. Every nontrivial code change needs tests and before/after coverage; current measured Windows baseline is 74.4%, lint zero, and build/vet/default/tagged integration/layout checks passed. Multi-step work requires a durable approved plan. Never stage, unstage, commit, or push, and preserve owner-staged changes. Never bulk-rewrite the repository or hand-edit generated Wire output; format only explicit permitted files and regenerate injectors when needed. Export must use the machine-detected game templates directory or the deliberate session-only picker override, never a persisted output path or an unrelated fallback.
>
> Before the latest documentation curation, HEAD was 6ec9e2c on master, tagged v0.3.8, with a clean index and working tree. The review, two source backlogs, and handoff were then edited without staging. Recheck before acting. Temporary review verification code was deleted. No global integration_test/gui/wireinject tags, no fake unit seams for private/Gio code, and no drive-by refactors.
>
> Owner decision: §2.2 zone-content DTO removal was explicitly reopened, but its API and scope still need planning; the bonuses DTO exception remains accepted. Other in-game policy checks remain unresolved. Module tidy dry-runs failed only because Windows checksum files have CRLF, not dependency-content drift. Local GUI/race/Linux/benchmark/vulnerability outcomes were not measured in this review. Follow ask → plan → approve → implement + verify → owner commits → mark, keeping finding numbers stable.
