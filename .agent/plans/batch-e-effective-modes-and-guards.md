# Batch E: Effective modes, guard propagation, and tournament players

<!-- markdownlint-configure-file {"MD024": {"siblings_only": true}} -->

Fix review findings §1.5, §1.11 and §1.12 only: invalidate incompatible manual snapshots on effective-mode transitions, propagate neutral-quality edits to incident preset guards, and enforce two-player tournament settings.

## Authority and approval gate

- **Current status, 2026-09-12: ALL FOUR PHASES COMPLETE.** The combined cross-phase flow is covered by six new real-input GUI tests, every Windows gate passes, and independent Claude Opus 5 reviews approved the plan, Phase 2, the Phase 3 implementation and the Phase 4 tests. Native Linux execution is recorded **UNAVAILABLE**, which Phase 4's verification plan explicitly permits. The owner accepted the implementation: *"Changes reviewed, all seems fine, you can finish up with Phase 4 and update Session Carry Forward"*, and separately accepted unit coverage moving from 75.1% to 74.9% for integration-covered GUI helpers. Phase 3 and partial Phase 4 are committed at `1b4658c` on `AD/modes_and_guard_propagation`; the only code change left in the worktree is one unstaged GUI test file. **No review finding is marked fixed yet** - that happens only after the owner commits this final verification and documentation work, per the surviving review's protocol. Every dated phase-status statement below is historical evidence, not current status.
- Date: 2026-09-12. Owner confirmed the product decisions and full scope below through three question rounds.
- *Historical, Phase 2 era:* **Written plan APPROVED; Phases 1 and 2 COMPLETE.** The owner committed the plan in `b9dd616` on `AD/modes_and_guard_propagation` and asked for Phase 1 completion, then for Phase 2 implementation in the following session. Phase 2 was implemented, fully verified and independently approved, and the owner has since committed it. Phases 3 and 4 were unstarted at that point; both are complete now.
- Read [AGENTS.md](../../AGENTS.md), the [current handoff](../session-carry-forward.md) and the [surviving review](../backlog/review-gpt-6-astra-09-07.md). This plan becomes the working record for Batch E after approval.
- Batches A/D, B and C are closed. Accept the owner's commit/engine record; do not retrieve retired documents, re-review those batches, or schedule a separate closed-batch verification. Normal regression suites for actual Batch E changes are required.
- The handoff's §8 is verbatim historical text. Do not edit it. Its later-scope decisions remain binding, but its Batch C phase/engine wording is superseded.
- Approval record, 2026-09-12: owner selected **"Approve the plan for subsequent implementation"** after reading the written-plan approval question and independent-review result. This is explicit approval of this plan, beyond prior scope confirmation. Subsequent instructions: **"the plan is commited, proceed with finishing Phase 1 and update carry forward"**, then **resume and finish Phase 2**. All four phases are now complete; do not repeat settled product questions, the approved plan review or any closed phase's work.
- Independent plan review: **APPROVED by Claude Opus 5, 2026-09-12**, after one revision. All five initial blockers resolved: remove unapproved direct-generator rejection; separate count validation from visible discard outcomes; invalidate before generation failure can consume transitions; rebind working/selected pointers; specify ordered Default matching.
- Independent Phase 2 implementation review: **APPROVED by Claude Opus 5, 2026-09-12**, scoped to Phase 2 and its regression interactions, after one blocker was resolved (see the Phase 2 summary). The Phase 3 implementation review and the Phase 4 review of the combined-flow tests were also **APPROVED by Claude Opus 5** on 2026-09-12.

## For Future Agents

As work proceeds, mark checkboxes when items complete. For every phase, record its status, changed files, decisions, exact verification results and blockers before moving on. Record baseline and final coverage using Go's own reporting, not averages or profile-block sums. Keep this plan self-contained and resumable. Fill Final Recap and Deployment Plan at completion. Stop and ask if a change requires expanding the approved scope.

## Confirmed product contract

### Effective-mode invalidation

- Effective tournament is `Tournament || VictoryCondition == Tournament`.
- Effective arena is `GladiatorArena || VictoryCondition == FinalBattle`; the UI calls the latter victory choice Guardian Arena.
- A false-to-true or true-to-false change in **either** effective boolean clears the **entire** manual zone/connection snapshot and regenerates. This includes tournament transitions when player count already equals two.
- Warn when actual manual edits are discarded. No new confirmation dialog. A discard is a persisted change and must set dirty/re-arm Exit.
- Changing checkbox/selector representation without changing either effective boolean preserves manual edits. Unrelated victory changes, arena timing, tournament timing/points and Save Army changes also preserve edits when both booleans remain unchanged.
- Preserve the current selector's checkbox-reset behavior. Judge the resulting effective modes, not selector labels alone. Do not normalize additional persisted checkbox fields merely for equivalence.
- Existing other layout-defining settings retain their behavior. No manual-delta migration, suspended snapshot or restoration on toggling back.
- Effective-mode transitions regenerate immediately through the existing layout-defining predicate, rather than using the 300 ms non-layout debounce.
- **No arena reconciliation:** no new cleanup/restoration/relocation pass after ordinary Apply or quality edits. Valid loaded states retain manual snapshots; no retrospective repair of stale arena markers or historical graphs is attempted.

### Incident guard recalculation

- Trigger only when a neutral zone's **resolved quality actually changes**. Same-quality requests and castle-only changes do not recalculate connection guards.
- Before modifying the zone, determine every incident connection's old effective quality using the existing stronger-endpoint rule. Player/hub, nil quality inference, unknown and missing-endpoint fallbacks remain those of the existing tier service.
- Infer identity by exact numeric match in that old quality's existing preset table. Carry the matching named tier (Default, Weakest, Low, Medium, High, Very High) to the new stronger-endpoint table after reprofile.
- Match the ordered preset list and take its first exact match. Default is a named tier, not a fixed number: Bronze Default 15,000 becomes Gold Default 25,000. Test this separately from the other named tiers.
- `GuardZone` controls placement only, not which table wins. If the opposite endpoint remains stronger, values remain unchanged. Example: Bronze/Silver Medium 24,000 becomes Gold Medium 48,000 when Bronze rises to Gold; lowering a weaker endpoint next to Gold does not change Gold's Medium value.
- All incident connection types participate, including Direct, Portal and arena. Parallel connections remain independent. Unrelated connections and every field other than the affected `GuardValue` remain unchanged.
- Unmatched numeric values remain unchanged and are **Custom**. Typing a value exactly equal to a table entry counts as that preset, including after reopening/reloading. No remembered selection or persisted preset metadata.
- Keep current tables exactly: Plastic and Bronze share Bronze's preset table. In particular, a generated Plastic guard of 10,000 does not match that table and stays Custom. Do not fix this table discrepancy in Batch E.
- Explicit Custom display must refresh after typing, selecting/reselecting a connection and reopening/loading. Selecting Custom preserves the number; choosing a named preset writes that preset's number. Invalid numeric text retains the last valid model value; do not add unrelated numeric validation.
- Quality propagation is pending dialog work until Apply. Cancel must leave the live template and retained manual snapshot unchanged. Save/reload retains only the numeric result; preset display is re-inferred.

### Tournament player count and loading

- Under either effective tournament alias, force player count to two and disable the existing player control. Keep its displayed and stored values aligned across frames.
- Leaving effective tournament mode unlocks the control but leaves two players. **No remembered previous player count**, session or persistent.
- Load with effective tournament and **any count other than two**: correct to two, clear all manual zones/connections, warn about correction and any actual discard, mark unsaved and re-arm Exit. No automatic file rewrite; saving remains explicit.
- Valid two-player tournament loads preserve manual edits. Other loads retain current correction/dirty behavior. Failed loads leave the current document unchanged.
- Validation without fixes reports the problem without mutating its input. Application generation with editor-state input uses the same corrected state so rules, labels, graph and exported player count cannot disagree.
- Warnings must survive the first automatic regeneration rather than being immediately replaced by a generic success message.

## Non-negotiable boundaries

- No writes to protected data, template schema or registry trees. No persisted format changes, DTO cleanup/removal, package moves, topology retirement/redesign or GUI/PNG geometry consolidation.
- A narrowly necessary extension of the existing quality-edit request/response contract is functional plumbing for §1.11, not permission for a DTO refactor. Reuse the existing mutation response where possible; keep services model/scalar-based and conversions at the handler boundary.
- Preserve explicit Portal road false/nil and valid approaches; internal roads independent of the between-zone checkbox; nil-state roads/content preservation; source cloning; unchanged opaque raster path with one reusable local half-opacity mask; Preview-only legend; machine-detected, never-persisted output directory.
- Never stage, unstage, commit, push, stash, switch branches or manipulate worktrees. Leave owner-staged changes alone. No bulk rewrites, global test tags, fake unit seams or hand-edited generated Wire.
- Windows and Linux compatibility remains required. Do not claim Linux/Steam Deck execution from Windows or from historical owner engine reports.

## Current source observations and implementation map

| Area | Current behavior and target seam |
| --- | --- |
| Effective aliases | [GeneratorConfig](../../internal/models/config/generatorConfig.go#L97-L105) has the two predicates; [EditorState](../../internal/models/editor_state_model/editorState.go#L98-L106) omits them from layout comparison. Add model-level effective predicates and compare booleans, with equivalence tests against config semantics. Do not introduce reverse model dependencies. |
| Regeneration | [decision service](../../internal/services/editor/regenerationDecisionService.go#L72-L91) preserves loaded snapshots when no previous generation exists. [driver generation](../../app/gui/drivers/stateGeneration.go#L82-L127) reapplies snapshots after generation and overwrites status. Retain the existing generation/reapply structure while enforcing mode invalidation and preserving new warnings. |
| Live validation | [GUI state model](../../app/gui/models/editorState.go#L47-L52) validates a clone but discards warning metadata. Keep the validator responsible for count correction only. A domain-model transition operation compares old stored modes against validated modes and detects invalid tournament count in the raw requested candidate; it clears manual snapshots on the validated clone and reports actual discard/correction. Invoke it before assigning current state, independently of generation snapshots. Return the transient outcome to [driver UpdateState](../../app/gui/drivers/state.go#L118-L123), which explicitly calls its dirty/Exit setter for an actual discard even if `WasStateChanged` is false, and queues the warning. No business policy is placed in the GUI. |
| Load | [state handler](../../internal/handlers/stateHandler.go#L28-L40) supports `fixIssues=false`; [validation](../../internal/handlers/stateHandler.go#L56-L75) clones before fixes. [driver load](../../app/gui/drivers/stateFiles.go#L95-L120) currently always marks loaded state clean. Call existing `LoadState(path, false)` once then existing `ValidateEditorState(raw, true)`, using only the latter's warnings. Use a model-level operation on raw and corrected state to identify tournament count correction and clear its snapshot, without comparing modes against the document being replaced. The driver consumes that structured outcome to set dirty/Exit and notify. No new load API, warning parsing, duplicate warnings, extra disk read or persisted metadata. |
| Player UI | [General panel](../../app/gui/panels/generalPanel.go#L104-L200) loads, rounds and saves a 2-8 slider without a tournament gate. Ensure the slider value itself becomes two, not just the outgoing DTO, and disable input under either alias. |
| Validator | [EditorStateValidator](../../internal/validators/editorStateValidator.go) currently has independent range checks. Add a tournament-count issue/fix that changes the count only; warning-producing domain transition/load policy owns manual-snapshot clearing. Ensure general range fixes cannot overwrite the two-player result; test fix ordering and extreme invalid values. |
| Application generator | [GenerateTemplate](../../internal/handlers/templateHandler.go#L60-L77) validates before mapping. Prove corrected editor-state input produces two-player labels, graph and rules. Do not change the separate direct GeneratorConfig contract in [TemplateGenerator](../../internal/services/template_generator/templateGenerator.go#L64-L103); defensive rejection there was not owner-approved. |
| Topology coordination | [TopologyProvider](../../internal/services/template_generator/providers/topologyProvider.go) dispatches tournament only for two labels. Keep topology implementations/fallback algorithms and direct generator configuration behavior unchanged. Test that the application handler supplies a valid two-player configuration. Direct provider/generator fallback behavior is outside this batch, not redesigned here; Batch K owns topology retirement. |
| Arena | [arena provider](../../internal/services/template_generator/providers/gladiatorArenaProvider.go#L44-L64) stamps generated variants; [UpdateTemplate](../../internal/handlers/templateHandler.go#L79-L108) replaces them with edited collections. Use invalidation, not a new arena reconciliation call. |
| Quality contract | [quality request](../../internal/dtos/zoneEditorQualityRequestDto.go) contains one zone; [handler](../../internal/handlers/zoneEditorHandler.go#L76-L84) returns one zone. Extend the request narrowly with zone/connection context and player names; return [existing mutation DTO](../../internal/dtos/zoneEditorMutationDto.go). Update facade/interface/mocks together. |
| Quality service | [ZoneEditorService](../../internal/services/connection_editor/zoneEditorService.go#L123-L159) reprofiles content only. Add a cohesive model-based operation for zone-plus-incident-guard editing, preserving the current zone-only method for existing callers. Inject the existing tier service and regenerate Wire if the constructor changes. |
| Guard rules | [tier service](../../internal/services/zones/zoneTierService.go#L53-L94) provides old/new effective quality; [preset tables](../../internal/common/common_connections/guardStrength.go) provide numeric-to-named-tier matching. Keep matching/recalculation business logic outside Gio. Add a focused public helper only if it has a real production caller and dedicated tests. |
| Pending GUI edits | [zone properties](../../app/gui/dialogs/zoneEditorZoneProps.go#L72-L101) must pass current working collections and consume the atomic result. This operation preserves connection order and count: capture the selected working index, rebuild working pointers from returned values and re-point selected by that index (or nil if none); set `syncedFor = nil` and `geometryDirty = true`. Reassign `this.zones` only after the last use of the old `zone` pointer, then reset zone-property synchronization. No canvas restructuring. |
| Preset display | [connection properties](../../app/gui/dialogs/zoneEditorConnectionProps.go#L67-L161) currently ignores failed preset matching. Append Custom last, outside the numeric values slice, preserving existing preset indices and bounds checks. Recompute display after numeric writeback without per-frame rewriting of the value on an unchanged selection. |

Implementation details may be adjusted within these seams, but a changed product rule, broader API redesign or persistence requirement needs owner approval. New methods get their own mirrored unit-test file; a pure request/result type does not need artificial tests.

## Phase 1: Approval and comparable baseline

Status: Complete

- [x] Read required instructions, handoff and surviving review; inspect Batch E callers and existing tests without modifying code.
- [x] Resolve product choices and receive full scope confirmation.
- [x] Draft this durable plan, with no implementation authorization implied.
- [x] Obtain independent Claude Opus 5 plan review; record verdict and address blockers here. Re-review approved all five revisions on 2026-09-12.
- [x] Obtain explicit owner approval of the written plan; record authorization above. Approved for a subsequent session on 2026-09-12.
- [x] After approval, inspect worktree/index read-only and record inherited changes without disturbing them. Branch `AD/modes_and_guard_propagation`, HEAD and plan commit `b9dd616` (`Init`); clean worktree, no staged paths, no inherited changes. No closed-batch documents or commits inspected.
- [x] Run fresh baseline unit coverage and record total/per-target function coverage before the first Go edit. PASS: 186 unit packages, 75.1% total statements, Windows/amd64 Go 1.27.0; target detail below.

### Verification Plan

- Discovery verification is source inspection only; no runtime claim. Check plan links and scope against confirmed answers.
- After implementation authorization, run the existing Go coverage task or equivalent: `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...`, then `go tool cover '-func=coverage.txt'`. Record exact output summary here before changing Go code.
- Historical comparison context only: the handoff reports Windows coverage 75.1% on 2026-09-11. Do not claim that as a fresh measurement or use the original review's older 74.4% to excuse a decrease.

### Phase Summary

Completed on 2026-09-12. Owner chose full snapshot invalidation instead of arena reconciliation, numeric preset identity with explicit Custom, no previous-count restoration, and automatic invalid tournament-load correction with warning/dirty state. Independent revised-plan review and explicit owner plan approval are complete; plan commit `b9dd616` was verified read-only on `AD/modes_and_guard_propagation`.

Fresh baseline: the existing coverage task ran `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...`, generated HTML and LCOV, and reported all **186 packages passing**. `go tool cover '-func=coverage.txt'` succeeded and reported **75.1%** total statements. `GOFLAGS` was empty. No application build, default/full suite, tagged integration/GUI suite, lint, Wire generation or Linux execution was run for this Phase 1 closeout. No production/test changes were made; only plan/handoff edits and ignored coverage reports result. Phase 2 remains unstarted.

### Baseline target function coverage

Measured from Go's function report, not per-test-package percentages. Existing 100% statement coverage does not establish the new mode/guard contracts; add the specified value/branch regressions during implementation.

| Current target | Function(s) | Baseline |
| --- | --- | --- |
| [Domain editor state](../../internal/models/editor_state_model/editorState.go) | `LayoutDefiningOptionsChanged`, `EqualsIgnoringManualEdits`, `HasManualEdits` | 100.0% each |
| [Config mode predicates](../../internal/models/config/generatorConfig.go) | `IsTournamentMode`, `IsGladiatorArenaMode` | 100.0% each |
| [Regeneration decisions](../../internal/services/editor/regenerationDecisionService.go) | `DecideRegeneration`, `DecideManualEditReapplication` | 100.0% each |
| [Validator](../../internal/validators/editorStateValidator.go) | `Validate` and every currently reported private validation/descriptor function | 100.0% each |
| [State handler](../../internal/handlers/stateHandler.go) | `LoadState`, `ValidateEditorState` | 100.0% each |
| [Generation handler](../../internal/handlers/templateHandler.go) | `GenerateTemplate` | 100.0% |
| [Quality handler](../../internal/handlers/zoneEditorHandler.go) | `ApplyZoneEditorQuality` | 100.0% |
| [Facade](../../internal/handlers/guiHandler.go) | `LoadState`, `ValidateEditorState`, `GenerateTemplate`, `ApplyZoneEditorQuality` | 100.0% each |
| [GUI state model](../../app/gui/models/editorState.go) | `UpdateCurrentState`, `HasManualEdits` | 100.0% each |
| [State driver](../../app/gui/drivers/state.go) | `UpdateState`, `flagAsUnsaved` | 100.0% each |
| [Load driver](../../app/gui/drivers/stateFiles.go) | `handleLoadState` | 0.0%; integration-only callback path |
| [Generation driver](../../app/gui/drivers/stateGeneration.go) | `handleGenerateTemplate`; `applyGeneratedTemplate`, `clearGeneratedState` | 81.0%; 100.0% each respectively |
| [Quality service](../../internal/services/connection_editor/zoneEditorService.go) | `NewZoneEditorService`, `ApplyNeutralZoneQuality` | 100.0% each |
| [Tier service](../../internal/services/zones/zoneTierService.go) | `ResolveQuality`, `GetGuardQuality`, `GetConnectionGuardQuality` | 100.0% each |
| [Preset lookup](../../internal/common/common_connections/guardStrength.go) | `GetGuardStrengthListForQuality`, `GetGuardStrengthForQuality` | 100.0% each |
| [Connection properties](../../app/gui/dialogs/zoneEditorConnectionProps.go) | `syncPropsFromConnection`, `writebackProps`, `guardPresetItems`, `matchGuardLabel` | 0.0% each; GUI integration territory |
| [Zone properties](../../app/gui/dialogs/zoneEditorZoneProps.go) | `syncZoneProps`, `writebackZoneProps` | 0.0% each; GUI integration territory |
| [General panel](../../app/gui/panels/generalPanel.go) | No function entries in this unit profile | Not instrumented, **not** a measured 0% or 100%; use GUI integration |

The coverage task regenerated ignored [coverage.txt](../../coverage.txt), [coverage.html](../../coverage.html) and [lcov.info](../../lcov.info). Baseline profile SHA-256: `44FE40F418607CC691CAE32DFD0D48439A9C620B5333B5A2ADF559CBF25EEBDF`. These reports are local artifacts, not tracked baseline files. The recorded figures survive report replacement. Before Phase 2 edits, check source/toolchain/coverage scope still matches this baseline; refresh it only if inputs changed or evidence is unavailable. This is Batch E's pre-change baseline, not a reopening of earlier batches.

Closeout checks passed: full index unchanged; handoff §8 exact UTF-8 text including line endings preserved; only this plan and the handoff modified in tracked files; local document links valid; whitespace check and editor diagnostics clean. The owner receives two unstaged documentation changes, not a code implementation.

## Phase 2: Effective modes and two-player state lifecycle

Status: Complete

Resume check, 2026-09-12: clean worktree/index on `AD/modes_and_guard_propagation`,
HEAD `8f852be0218180d8f5629e2580f61669af3bdb2a`. No Go/module/VS Code input
differences from `b9dd616`; Windows/amd64 Go 1.27.0, empty GOFLAGS and the recorded
coverage SHA-256 still match. Reuse Phase 1's 75.1% baseline without rerunning it.
No inherited source changes. Implementation and verification now underway.

- [x] Add effective editor-state predicates and include boolean transitions in layout comparison. Keep config behavior equivalent, aliases intact and unrelated rules non-layout-defining.
- [x] Implement tournament-count validation/correction only, including fix ordering and report-only behavior. Implement snapshot clearing and its structured outcome in the domain transition/load operation, not the validator fix.
- [x] Apply live invalidation during state update, before any generation attempt: compare stored versus validated effective modes, and raw requested versus validated tournament count. Clear on the validated clone, then assign state and deliver actual discard/correction to the driver. Thus a failed generation cannot consume the invalidation; no change to generation's snapshot-on-failure policy is needed. Cover before-first-generation and failure/retry explicitly.
- [x] Explicitly mark an actual manual discard dirty/re-arm Exit using its outcome, independently of `WasStateChanged` and its intentional manual-field exclusion.
- [x] Update load to inspect original structured state after the single disk read: the existing non-fixing load validation is followed by one fixing validation pass, whose warnings are used. Set dirty/re-arm Exit only for tournament count corrections and retain atomic failed-load behavior.
- [x] Carry a small session-only notification outcome through automatic regeneration so new discard/correction warnings remain visible. Clear pending notices on document replacement/reset; do not redesign global status history or suppress generation errors. Test failed generation/retry and subsequent-document isolation.
- [x] Fix/disable the player slider for both effective tournament aliases; leaving mode retains two. Preserve selector reset behavior and non-tournament range behavior. Test several idle save/render frames: no stale 3-8 value is repeatedly submitted and no repeated correction warning is queued.
- [x] Verify the application generation invariant through the existing validate-before-map handler. Leave direct GeneratorConfig generation, provider fallbacks and topology implementations unchanged.
- [x] Add focused unit and integration tests in lockstep with changes.

### Verification Plan

- Extend [layout comparison tests](../../test/unit/internal/models/editor_state_model/editorState/layoutDefiningOptionsChanged_test.go), [regeneration decisions](../../test/unit/internal/services/editor/regenerationDecisionService/decideRegeneration_test.go), [manual reapplication decisions](../../test/unit/internal/services/editor/regenerationDecisionService/decideManualEditReapplication_test.go), [GUI-model updates](../../test/unit/app/gui/models/editorState/updateCurrentState_test.go) and [driver updates](../../test/unit/app/gui/drivers/state/updateState_test.go).
- Test each mode both directions with zone-only, connection-only, full and absent snapshots; two-player tournament transitions; both aliases true and removal of just one alias; alias swaps with unchanged effective booleans; simultaneous tournament/arena changes; unrelated rule changes; first generation and failure/retry. Test input clone isolation, dirty and Exit behavior separately.
- Extend [validator tests](../../test/unit/internal/validators/editorStateValidator/validateEditorState_test.go), [handler validation](../../test/unit/internal/handlers/stateHandler/validateEditorState_test.go), [handler load](../../test/unit/internal/handlers/stateHandler/loadState_test.go) and affected facade tests. Enumerate 3-8 and below/above-range counts for both tournament aliases; two is a no-op. Non-tournament behavior stays unchanged; no-fix validation and failed load preserve input/current document.
- Extend [generation handler tests](../../test/unit/internal/handlers/templateHandler/generateTemplate_test.go) and real application-handler integration for both aliases, graph/player/rules consistency, warnings and caller ownership after correction. Keep direct provider compatibility assertions and direct generator 2-8 player tests unchanged; direct GeneratorConfig validation is not this batch's contract.
- Extend [state integration](../../test/integration/editorState_integration_test.go) and appropriate existing Exit/load suites. Load malformed-count tournament settings with a manual snapshot; verify correction, discard, dirty/Exit, warning after actual next-frame regeneration, and unchanged file bytes until explicit save. Save/reload the corrected file; it must load cleanly at two. Valid two-player loads preserve snapshots.
- Exercise General UI via real frame/input integration: checkbox-origin loaded state, victory-selector alias, fixed/disabled player control, uncheck while a selector still activates tournament, leaving tournament then changing count, and next-frame writeback. Use existing gated GUI harness, no global tags or new fake unit exports.
- Run focused unit/integration/GUI suites and `go build ./...`; record exact results. Verify required state warnings are visible, not merely produced by an unused return value.

### Phase 2 verification results

All commands below are the final post-fix run on **Windows/amd64, Go 1.27.0, empty `GOFLAGS`**, 2026-09-12:

- `go build ./...` — **PASS**.
- `go test ./test/unit/... -count=1` — **PASS, 187 packages**. The count rose from the 186-package baseline because `ModeTransitionOutcome` brought its own mirrored test folder.
- `go test ./test/...` — **PASS** (untagged default run).
- `go test -tags='integration_test,gui' ./test/integration/...` — **PASS**, covering both the root integration suite and the GPU-backed GUI suite.
- `go run ./cmd/testlayoutcheck .` — **PASS**.
- Coverage equivalent of the existing task: `go test '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...` then `go tool cover '-func=coverage.txt'` — **PASS at 75.1% total statements, unchanged from the Phase 1 baseline** despite the new code. HTML and LCOV reports were refreshed. New profile SHA-256: `34DD607508B5BE7F21C8BC466580C9722AAD061C6E273E0D338013B52C7C7187`.
- Report-only lint, `golangci-lint-v2 run ./... --issues-exit-code=1` — **PASS with zero issues**, after reordering the two new unexported `EditorState` methods below the exported ones (`funcorder`) and formatting an explicit `gofmt -l` file list. No auto-fix sibling task and no bulk rewrite were used.
- **Native Linux execution: UNRUN.** A WSL Ubuntu environment exists, but `go` and `pkg-config` are unavailable inside it and installing toolchains is out of scope. Cross-platform correctness rests on source review only; Phase 4 still owes a real Linux run.

Changed-function coverage against the Phase 1 table: `IsEffectiveTournament`, `IsEffectiveGladiatorArena`, `ApplyModeTransition`, `ModeTransitionOutcome.MergedWith`, `validateTournamentPlayerCount`, `UpdateCurrentState`, `OverrideStateFromLoad`, `UpdateState` and the three notice functions in [stateNotices.go](../../app/gui/drivers/stateNotices.go) all report **100.0%**. [handleLoadState](../../app/gui/drivers/stateFiles.go) and its new `describeLoadedState` helper remain **0.0% in the unit profile** — they are driver callbacks reached only through the integration harness, where they are now exercised directly. The General panel still contributes no entries to this profile; its new locking behavior is proven by integration frames only.

### Phase Summary

Completed on 2026-09-12; findings §1.5 and §1.12 are implemented, §1.11 remains Phase 3 work.

**Domain.** [EditorState](../../internal/models/editor_state_model/editorState.go) gained `IsEffectiveTournament`, `IsEffectiveGladiatorArena`, the `TournamentPlayerCount` constant and `ApplyModeTransition`, with private `effectiveModesChanged`/`tournamentCountCorrected` helpers; `LayoutDefiningOptionsChanged` now compares the two resolved booleans rather than the checkbox/selector fields, so alias-only changes keep the layout. The transient result type lives in the new [modeTransitionOutcome.go](../../internal/models/editor_state_model/modeTransitionOutcome.go) with `ManualEditsDiscarded`, `TournamentCountCorrected` and `MergedWith`. `ApplyModeTransition` takes the previous document (nil on load) and the raw requested state, clears the snapshot on an actual transition or correction, and reports only what the user can see.

**Validation.** [EditorStateValidator](../../internal/validators/editorStateValidator.go) appends `validateTournamentPlayerCount` **last** in `Validate`, so the general range fix cannot clamp back over the two players tournament mode requires. Its fix changes the count and nothing else; snapshot clearing stays with the domain transition.

**State lifecycle.** [UpdateCurrentState](../../app/gui/models/editorState.go) now validates the edited clone, applies the transition to it and returns the outcome before storing it, so no generation attempt can ever observe an incompatible snapshot. The new `OverrideStateFromLoad` compares the raw file state against its corrected result, which is the only place a corrected count is still visible. [UpdateState](../../app/gui/drivers/state.go) flags the document unsaved when `WasStateChanged` is false but a discard or correction happened, and holds the outcome in a new `pendingOutcome` field that `Reset` clears.

**Notices.** The new [stateNotices.go](../../app/gui/drivers/stateNotices.go) owns `stateTransitionNotice`, `noteStateTransition` and `takePendingNotice`. A notice is shown immediately, kept pending across the automatic regeneration that follows, and consumed exactly once by a **successful** generation in [handleGenerateTemplate](../../app/gui/drivers/stateGeneration.go); a failed generation leaves it for the retry.

**Load.** [handleLoadState](../../app/gui/drivers/stateFiles.go) reads once through the existing `LoadState(path, false)`, then runs `ValidateEditorState(raw, true)` and uses only its warnings. It marks the document unsaved solely for a tournament count correction, leaves the file on disk untouched and composes its status through the new `describeLoadedState`. Failed loads still leave the current document alone.

**Player control.** [GeneralPanel](../../app/gui/panels/generalPanel.go) gained `getAllowedPlayerCount`, `isPlayerCountLocked` and `getPlayerCountRowWidget`, which renders the existing slider row through `gtx.Disabled()` under either alias. `SaveToState` moves the slider itself back to two rather than only the outgoing value, so idle frames cannot resubmit a stale 3-8 count or re-announce the same correction.

**Review blocker and fix.** Claude Opus 5's implementation review rejected the first cut: a no-op `UpdateState` merged an empty outcome and rewrote the status, wiping a generation error or a just-saved message on the very next idle frame. `noteStateTransition` now returns before merging when the incoming outcome is empty. The regression is pinned end to end - load, failed generation, idle frames, save, idle frames, retry - including the error flag itself. The review approved after that fix.

**Tests.** New unit folders cover `IsEffectiveTournament`, `IsEffectiveGladiatorArena`, `ApplyModeTransition`, `MergedWith` and `OverrideStateFromLoad`; existing suites for layout comparison, both regeneration decisions, `UpdateCurrentState`, `UpdateState`, the validator and both state-handler entry points were extended. Two new tagged integration files carry the behavior the unit profile cannot reach: [effectiveModes_integration_test.go](../../test/integration/effectiveModes_integration_test.go) (invalid and valid tournament loads under both aliases, discard, dirty/Exit, unchanged file bytes, save/reload, notice survival across regeneration, failure/save/retry ordering, document isolation, and two-player generation through the application handler) and [generalPanelTournament_integration_test.go](../../test/integration/generalPanelTournament_integration_test.go) (real slider drags, both aliases, disabled control, unchecking the rule under a tournament victory, leaving tournament, and idle-frame stability). Shared helpers gained `AppRunner.SetStatus` plus `DragPlayerCountToMaximum`, `SelectVictoryCondition` and `ToggleConditionRule` on the General tab handler. No new build tags, global tags or unit-test seams were introduced.

**Untouched, as required.** Protected data/schema/registry trees, the output path, the opaque raster path, generated Wire, topology implementations and the direct `GeneratorConfig` contract are unchanged. No dependencies, snapshots or goldens changed. All Phase 2 work is unstaged/untracked and awaits the owner.

## Phase 3: Quality propagation and explicit Custom display

Status: Complete

Resume, 2026-09-12: owner requested Phase 3. Worktree and index are clean on
`AD/modes_and_guard_propagation`; Phase 2 has been committed since the handoff.
Approved product decisions remain closed. Implementation starts with a cloned
zone/connection mutation and ordered preset remapping, followed by dialog pointer
rebinding and explicit Custom synchronization. Protected paths remain untouched.

- [x] Extend the existing quality-edit handler contract with current zone/connection context; reuse the mutation response. Update interfaces, facade, mocks and direct callers together, with no unrelated DTO changes.
- [x] Add model-based zone-plus-incident-connection editing in the existing service package. Capture old effective quality/preset before reprofile, resolve new endpoint quality afterward, and change only matching incident guard values on an actual quality change.
- [x] Clone zones/connections before mutation using existing deep-clone facilities; preserve order, names, endpoint identity and nil/empty semantics. Keep the current zone-only reprofile API and road finalization behavior for other callers.
- [x] Inject the tier collaborator where needed; regenerate Wire through the existing task. Update explicit test constructors/helpers, never generated source by hand.
- [x] Consume the whole mutation using the index-preserving pointer replacement described above. Rebind selected, clear `syncedFor`, finish old zone-pointer use before replacing zones, and refresh geometry/property state. Test that editing a connection after a quality mutation affects the retained working collection, not a discarded allocation.
- [x] Add explicit Custom display and exact-match synchronization. Use a nonnumeric Custom item with safe index mapping; choosing it must not change the guard number.
- [x] Add unit, handler, GUI and save/reload regressions as changes land.

### Verification Plan

- New service public method gets its own lower-camel test file under the existing mirrored service folder. Extend [zone reprofile tests](../../test/unit/internal/services/connection_editor/zoneEditorService/applyNeutralZoneQuality_test.go), [quality handler tests](../../test/unit/internal/handlers/zoneEditorHandler/applyZoneEditorQuality_test.go), [facade tests](../../test/unit/internal/handlers/guiHandler/applyZoneEditorQuality_test.go) and constructor tests when affected.
- Matrix: all six named presets, with a dedicated Bronze Default 15,000 to Gold Default 25,000 test and ordered first-match assertion; upgrades/downgrades; stronger unchanged opposite endpoint; equal endpoints; reversed endpoint order; hub/player/unknown/nil-quality inference and existing missing-endpoint fallback; Plastic/Bronze unchanged mapping; multiple parallel incident edges; self-edge if supported by the existing API; unrelated edges; empty/nil collections.
- Preserve unmatched zero, negative and arbitrary custom numbers; exact typed numeric matches are presets. Explicitly test the generated Plastic 10,000 exception. Test same-quality/castle-only requests do not recalculate. Whole-connection comparisons must pin type, Road true/false/nil, placement data, GuardZone, weekly increments and all other fields.
- Test source/result mutation isolation for nested zones and connections. No protected type changes, new Clone ownership semantics or production test-only seams.
- Extend [zone property GUI integration](../../test/integration/gui/zoneEditorProperties_integration_test.go): named preset -> quality change -> incident new value; custom value preserved/displayed; exact typed match becomes preset; select Custom without numeric change; reselect/reopen; no stale widgets on subsequent frames; Apply vs Cancel. Use real inputs/visible output with existing harness rather than new broad test exports.
- Add persistence round trip through existing public handlers: edited quality and numeric guard survive save/reload; reopened GUI infers the same label/Custom; no new fields in either persisted format. Use test-owned temporary paths only.
- Run focused unit/integration/GUI suites, build and generated-wiring checks. Check actual GUI labels/value results rather than only nonzero pixels or changed hashes.

### Phase Summary

Complete; implementation and verification details follow.

### Phase 3 completed work and evidence

Completed 2026-09-12. `ApplyNeutralZoneQualityEdit` takes the new model request,
deep-clones both collections, captures ordered preset indices before reprofile and
resolves the new stronger-endpoint table afterward. Same-quality/castle-only edits
never capture presets; custom numbers are untouched. Missing target returns clones
unchanged. The existing zone-only method and road policy remain intact. The tier
service is injected; Wire generation passed and a second generation was unchanged.
DTO/interface/facade/mocks and constructors were updated together.

The dialog installs the full mutation, rebinds selected by working index and clears
property/geometry synchronization. Custom is appended outside the numeric values
slice and the displayed selection is derived after numeric writeback, never persisted.
Choosing Custom on an exact preset retains the number and displays that preset again.

Dedicated service tests: 46 passing cases, including all named tiers, dedicated
Bronze Default -> Gold Default, stronger endpoints, inference/fallbacks, all edge
types, parallel/self edges, road true/false/nil, custom/Plastic 10,000, cloning and
no-op requests. Current tables have no duplicate values: ordered lookup is explicit
in code; precedence cannot be experimentally distinguished with current data.
Handler/facade tests pass. GUI adds 21 real-input regressions for remapping, Custom,
exact/invalid typing, idle frames, collection replacement, Apply/Cancel including
existing manual state, reopen and save/load. Two handler persistence regressions pass.
Selection is nil during reachable zone editing; the non-nil rebinding branch is
defensive and source-reviewed, not claimed as directly exercised by GUI inputs.

## Phase 4: Whole Batch E verification and owner handback

Status: Complete

- [x] Verify the combined real flow: manual edits -> arena/tournament change -> new generated graph; invalid tournament load -> correction -> frame/save/reload; pending quality edits -> Apply/Cancel -> regeneration/persistence.
- [x] Run final build, fresh unit tests, comparable coverage, default/integration/GUI suites, test-layout checker and report-only lint. No blanket snapshot updates or auto-fix across the repository.
- [x] Compare before/after total coverage and changed-function coverage. Cover every reachable changed branch/logical unit; record genuine GUI-only/unreachable limitations instead of fake unit seams.
- [x] Obtain independent Claude Opus 5 implementation review limited to Batch E and its regression interactions. Address blockers and rerun affected checks.
- [x] Inspect diff scope and owner's index read-only; protected trees/output path/opaque raster/retired topology work remain untouched except generated Wire produced by its generator as required.
- [x] Present owner-facing behavior and verification results. Owner performs all staging/commits and acceptance. Do not mark review findings fixed before the owner's commit protocol is satisfied.
- [x] Complete Final Recap and Deployment Plan, and update this plan with remaining blockers if work stops partway.

### Verification Plan

- `go build ./...`
- `go test ./test/unit/... -count=1` once freshly for final validation; subsequent repeats without code changes may use cache.
- Coverage task/equivalent with `-coverpkg=./internal/...,./app/...` and `go tool cover -func`; use Go's authoritative total. Final coverage must not drop below the freshly measured comparable baseline; explain denominator changes and seek approval, never hide them.
- `go test ./test/...`
- `go test -tags=integration_test ./test/integration/...`
- `go test -tags='integration_test,gui' ./test/integration/...` with actual GPU execution when available. Record unavailable environments as unrun, not passing. Tags are per invocation only.
- `go run ./cmd/testlayoutcheck .`
- Existing report-only lint task; preserve zero reported issues, distinguish existing unused-exclusion warnings.
- Format only an explicit changed-file list identified by `gofmt -l`, excluding protected/generated files. Do not use whole-tree `gofmt -w` or report-only lint's auto-fix sibling task.
- Wire generation task after constructor/provider changes; inspect generated output, build it and record generation failures honestly.
- Cross-platform source review plus Linux build/tests in an available supported Linux environment/CI. Do not install toolchains, mutate remote state or claim a native build solely from source inspection.
- Repository diff/whitespace and plan-link checks are read-only. Never alter the index to produce a clean status.

### Phase Summary

Completed 2026-09-12. Phase 4 added the combined cross-phase coverage the earlier phases
deliberately left open, reran every Windows gate over the whole batch, and closed the
owner handback.

**Combined flow.** Six real-input tests were added to the single changed code file,
[zoneEditorProperties_integration_test.go](../../test/integration/gui/zoneEditorProperties_integration_test.go)
(+171 lines). They Apply a remapped Silver Medium to Gold Medium incident edge and then
switch into Guardian Arena or Tournament through the General victory selector, across
idle frames and a dialog reopen. They assert positive outcomes, not merely absence: the
arena main object exists on the regenerated Hub graph, the old Neutral-C zone and its
remapped edge are gone, the discard notice survives the generation status that follows,
and the tournament transition replaces both the Hub layout and the manual zone. The
existing load-correction, frame, save/reload and Apply/Cancel persistence suites already
cover their own halves and were not duplicated. The dirty-flag side is unit-covered;
these GUI tests check the visible warning.

**Gates.** Build, fresh unit run with coverage, the untagged default suite, the tagged
integration plus GPU GUI suite, the test-layout checker, `gofmt -l` on the changed file
and report-only lint at `--issues-exit-code=1` all pass. Exact commands and results are
in the verification ledger below. No constructor or provider changed, so the committed
generated Wire is already current and was deliberately not regenerated.

**Coverage.** 74.9% total statements against the 75.1% baseline. This is the owner's
explicitly accepted exception for integration-covered GUI helpers, not a restored
baseline. All three new service functions and both quality handler/facade methods are
100.0%.

**Linux.** Recorded **UNAVAILABLE**, exactly as this phase's verification plan allows.
The WSL Ubuntu probe run this turn produced empty output and exit code 1 for both `go`
and `pkg-config`. Nothing was installed, and no native Linux, Steam Deck or in-game
result is claimed.

**Review and handback.** Claude Opus 5 approved the six new tests independently, after
earlier approvals of the plan, Phase 2 and the Phase 3 implementation. The owner reviewed
and accepted the implementation. No finding is marked fixed yet: the surviving review's
protocol requires the owner's commit of this final work first, and this session performed
no Git mutation of any kind.

## Verification ledger

### Batch-final results, 2026-09-12

Windows/amd64, Go 1.27.0, empty `GOFLAGS`. Clean start at owner commit `1b4658c`
("Phase 3 and partial 4") on `AD/modes_and_guard_propagation`, branch in sync with
origin. The only code change in the worktree is the one unstaged GUI test file.

- `go build ./...` — PASS.
- `go test -p=2 -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...` — PASS. Bounded compile parallelism stays because unbounded Windows linkers were interrupted (`0xc000013a`) previously; that incomplete profile was discarded then and is not reused. This is the baseline over the reviewed production code, which has not changed since.
- `go tool cover '-func=coverage.txt'` — **74.9%** total statements against the recorded 75.1% baseline. The owner explicitly accepted this decrease; it is an accepted exception, not a restored baseline. All three new service functions and both quality handler/facade methods report 100.0%. The three new GUI helpers report 0.0% here and are integration-covered, with the defensive pointer branch noted in Phase 3.
- `go test ./test/...` — PASS, untagged default run.
- `go test -tags='integration_test,gui' ./test/integration/...` — PASS; root integration 3.297s, GUI integration 28.558s. The single invocation covers both the root suite and the GPU-backed GUI suite.
- `go run ./cmd/testlayoutcheck .` — PASS. `gofmt -l` over the changed file — clean.
- `golangci-lint-v2 run ./... --issues-exit-code=1` — PASS, zero issues.
- Wire: no constructor or provider changed, so the committed generated output is already current and was deliberately not regenerated.
- Independent Claude Opus 5 review of the six new combined tests: **APPROVED**; earlier plan, Phase 2 and Phase 3 implementation reviews were also approved.
- **Native Linux / Steam Deck: UNAVAILABLE.** The probe `wsl.exe -d Ubuntu -- sh -lc 'command -v go; command -v pkg-config'` returned empty output and exit code 1 this turn. Nothing was installed, and no native Linux, Steam Deck or engine result is claimed. Cross-platform correctness rests on source review, which this phase's verification plan explicitly permits recording as an unavailable environment.
- **No coverage profile fingerprint is claimed for this run.** The reports were regenerated but their SHA-256 was not measured. The Phase 3 value `23B78F2CF15B73902398D8A21A755FD684D3668486E56D8CFA770BE2BE7F9F31` is historical and no longer describes the current file.
- No protected/output/topology/schema writes, no dependency, snapshot or golden changes, no owner-findings edits, and no Git mutations of any kind.

Historical: Phase 3's closeout recorded the same 74.9% at owner commit `f80f7ca`, with
lint clean after wrapping two test literals and the same Linux outcome. The transient
75.2% reading had a different delegated coverage scope and was never a comparable
full baseline.

The Phase 2 column is scoped to Phase 2's changes; the batch-final column is the Phase 4
rerun over the whole batch, Phase 3 included.

| Check | Baseline | Phase 2 result | Batch final (Phase 4) |
| --- | --- | --- | --- |
| Build | Not run in Phase 1 | PASS, `go build ./...` | PASS |
| Fresh unit / coverage | PASS, 186 packages, 75.1%; Windows/amd64 Go 1.27.0 at `b9dd616`, 2026-09-12 | PASS, 187 packages, 75.1% unchanged; profile SHA-256 `34DD6075...C7C7187` | PASS at **74.9%**, owner-accepted exception; no fingerprint measured this run |
| Focused domain / handler / GUI / persistence | Source inspected only | PASS; new/extended unit folders plus two tagged integration files | PASS; six added combined real-input GUI regressions |
| Default / tagged integration / GUI | Not run | PASS, `go test ./test/...` and `-tags='integration_test,gui' ./test/integration/...` | PASS; root 3.297s, GUI 28.558s |
| Layout / report-only lint / formatting / Wire | Not run | PASS layout checker; lint zero issues at `--issues-exit-code=1`; explicit `gofmt -l` list only; no Wire change required | PASS layout checker; lint zero issues; `gofmt -l` clean; no Wire change required |
| Native Linux execution | Not run | UNRUN; WSL Ubuntu lacks `go`/`pkg-config`, nothing installed | **UNAVAILABLE**; probe empty, exit 1; nothing installed, nothing claimed |
| Independent plan review | Claude Opus 5 APPROVED revised plan, 2026-09-12 | N/A | N/A |
| Independent implementation review | N/A | Claude Opus 5 APPROVED Phase 2 after the notice idle-frame fix, 2026-09-12 | APPROVED, batch-wide, including the six new tests |

Discovery tooling notes: initial subagent requests used unsupported lowercase model IDs; retrying with exact display name `GPT-5.6 Terra (copilot)` succeeded. Two guessed source paths did not exist; confirmed source/search results, not those guesses, inform this plan. No repository changes resulted from those failed reads/calls. Existing tooling memory already documents exact model-name requirements.

## Final Recap

Batch E is complete. Findings §1.5, §1.11 and §1.12 all have implemented, tested and
independently reviewed behavior, and the owner accepted the implementation on 2026-09-12.

- **§1.5** — effective tournament and effective arena booleans drive layout comparison. A
  transition in either clears the entire manual snapshot before any generation can observe
  it, warns on an actual discard, marks the document unsaved and re-arms Exit. Alias-only
  representation changes preserve edits. No arena reconciliation pass was added.
- **§1.12** — effective tournament forces and locks two players, and leaving it keeps two
  with no remembered count. Invalid tournament loads correct to two, clear incompatible
  snapshots, warn, become unsaved and re-arm Exit without rewriting the file, and the
  warning survives the automatic regeneration that follows.
- **§1.11** — an actual neutral-quality change remaps exact preset guard values on every
  incident edge through the old and new stronger-endpoint tables, carrying the named tier
  including Default. Unmatched numbers stay put and display an explicit Custom entry
  derived from the number, never persisted.

Phase 4 supplied the combined cross-phase GUI regressions, reran every Windows gate green
and closed the handback. Coverage sits at the owner-accepted 74.9%. Native Linux is
unavailable and unrun. The three findings remain **unmarked** in the surviving review
until the owner commits this final verification and documentation work; only then are
§1.5, §1.11 and §1.12 marked per that review's protocol.

No in-game outcome is claimed. Do not reopen accepted product decisions, reimplement any
phase, or extend into Batch K, topology retirement, direct GeneratorConfig rejection,
geometry consolidation, DTO cleanup, schema or package work.

## Deployment Plan

Nothing is deployed and the assistant made no Git change. Phase 3 and partial Phase 4 are
committed at `1b4658c`; what remains is one unstaged GUI test file plus this plan and the
handoff.

1. Owner reviews the remaining unstaged changes, then stages and commits them. Only after
   that commit are review findings §1.5, §1.11 and §1.12 marked fixed, following the
   surviving review's protocol.
2. Owner builds and packages through the existing platform workflow. No dependency
   installation, saved-schema migration, output-directory change or Wire regeneration is
   required; the committed generated wiring is current.
3. On launch, per-machine game-directory detection is unchanged. Existing invalid
   tournament settings are corrected in memory with a warning and unsaved state; the file
   changes only on explicit Save.
4. Communicate that mode transitions discard manual edits, that leaving tournament leaves
   two players, and that unmatched guard numbers display as Custom without acquiring any
   hidden persistent preset identity.
5. Record platform verification honestly: Windows verified, native Linux and Steam Deck
   unavailable and unrun. Do not infer engine results from Batch C's owner report.
