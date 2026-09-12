# Batch E: Effective modes, guard propagation, and tournament players

<!-- markdownlint-configure-file {"MD024": {"siblings_only": true}} -->

Fix review findings §1.5, §1.11 and §1.12 only: invalidate incompatible manual snapshots on effective-mode transitions, propagate neutral-quality edits to incident preset guards, and enforce two-player tournament settings.

## Authority and approval gate

- Date: 2026-09-12. Owner confirmed the product decisions and full scope below through three question rounds.
- **Written plan APPROVED for a subsequent implementation session.** No production code, tests, generated files, configuration or existing reports have changed during discovery/planning. No tests were run. This session must not begin implementation.
- Read [AGENTS.md](../../AGENTS.md), the [current handoff](../session-carry-forward.md) and the [surviving review](../backlog/review-gpt-6-astra-09-07.md). This plan becomes the working record for Batch E after approval.
- Batches A/D, B and C are closed. Accept the owner's commit/engine record; do not retrieve retired documents, re-review those batches, or schedule a separate closed-batch verification. Normal regression suites for actual Batch E changes are required.
- The handoff's §8 is verbatim historical text. Do not edit it. Its later-scope decisions remain binding, but its Batch C phase/engine wording is superseded.
- Approval record, 2026-09-12: owner selected **"Approve the plan for subsequent implementation"** after reading the written-plan approval question and independent-review result. This is explicit approval of this plan, beyond prior scope confirmation. Next session starts with Phase 1's remaining read-only worktree inspection and fresh coverage baseline, then Phase 2. No code is authorized in the current discovery/planning session.
- Independent plan review: **APPROVED by Claude Opus 5, 2026-09-12**, after one revision. All five initial blockers resolved: remove unapproved direct-generator rejection; separate count validation from visible discard outcomes; invalidate before generation failure can consume transitions; rebind working/selected pointers; specify ordered Default matching. No implementation review has occurred.

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
| Live validation | [GUI state model](../../app/gui/models/editorState.go#L47-L52) validates a clone but discards warning metadata. Keep the validator responsible for count correction only. A domain-model transition operation compares old stored modes against validated modes and detects invalid tournament count in the raw requested candidate; it clears manual snapshots on the validated clone and reports actual discard/correction. Invoke it before assigning current state, independently of generation snapshots. Return the transient outcome to [driver UpdateState](../../app/gui/drivers/state.go#L135-L140), which explicitly calls its dirty/Exit setter for an actual discard even if `WasStateChanged` is false, and queues the warning. No business policy is placed in the GUI. |
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

Status: In progress

- [x] Read required instructions, handoff and surviving review; inspect Batch E callers and existing tests without modifying code.
- [x] Resolve product choices and receive full scope confirmation.
- [x] Draft this durable plan, with no implementation authorization implied.
- [x] Obtain independent Claude Opus 5 plan review; record verdict and address blockers here. Re-review approved all five revisions on 2026-09-12.
- [x] Obtain explicit owner approval of the written plan; record authorization above. Approved for a subsequent session on 2026-09-12.
- [ ] After approval, inspect worktree/index read-only and record inherited changes without disturbing them. Do not look up closed-batch documents or commits.
- [ ] Run fresh baseline unit coverage and record total/per-target function coverage before the first Go edit.

### Verification Plan

- Discovery verification is source inspection only; no runtime claim. Check plan links and scope against confirmed answers.
- After implementation authorization, run the existing Go coverage task or equivalent: `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...`, then `go tool cover '-func=coverage.txt'`. Record exact output summary here before changing Go code.
- Historical comparison context only: the handoff reports Windows coverage 75.1% on 2026-09-11. Do not claim that as a fresh measurement or use the original review's older 74.4% to excuse a decrease.

### Phase Summary

Read-only discovery and scope questions completed on 2026-09-12. Owner chose full snapshot invalidation instead of arena reconciliation, numeric preset identity with explicit Custom, no previous-count restoration, and automatic invalid tournament-load correction with warning/dirty state. Independent revised-plan review and explicit owner plan approval are complete. No code, tests, generated files or coverage artifacts changed. Worktree inspection and baseline remain for the subsequent implementation session.

## Phase 2: Effective modes and two-player state lifecycle

Status: Not started

- [ ] Add effective editor-state predicates and include boolean transitions in layout comparison. Keep config behavior equivalent, aliases intact and unrelated rules non-layout-defining.
- [ ] Implement tournament-count validation/correction only, including fix ordering and report-only behavior. Implement snapshot clearing and its structured outcome in the domain transition/load operation, not the validator fix.
- [ ] Apply live invalidation during state update, before any generation attempt: compare stored versus validated effective modes, and raw requested versus validated tournament count. Clear on the validated clone, then assign state and deliver actual discard/correction to the driver. Thus a failed generation cannot consume the invalidation; no change to generation's snapshot-on-failure policy is needed. Cover before-first-generation and failure/retry explicitly.
- [ ] Explicitly mark an actual manual discard dirty/re-arm Exit using its outcome, independently of `WasStateChanged` and its intentional manual-field exclusion.
- [ ] Update load to inspect original structured state after the single disk read: the existing non-fixing load validation is followed by one fixing validation pass, whose warnings are used. Set dirty/re-arm Exit only for tournament count corrections and retain atomic failed-load behavior.
- [ ] Carry a small session-only notification outcome through automatic regeneration so new discard/correction warnings remain visible. Clear pending notices on document replacement/reset; do not redesign global status history or suppress generation errors. Test failed generation/retry and subsequent-document isolation.
- [ ] Fix/disable the player slider for both effective tournament aliases; leaving mode retains two. Preserve selector reset behavior and non-tournament range behavior. Test several idle save/render frames: no stale 3-8 value is repeatedly submitted and no repeated correction warning is queued.
- [ ] Verify the application generation invariant through the existing validate-before-map handler. Leave direct GeneratorConfig generation, provider fallbacks and topology implementations unchanged.
- [ ] Add focused unit and integration tests in lockstep with changes.

### Verification Plan

- Extend [layout comparison tests](../../test/unit/internal/models/editor_state_model/editorState/layoutDefiningOptionsChanged_test.go), [regeneration decisions](../../test/unit/internal/services/editor/regenerationDecisionService/decideRegeneration_test.go), [manual reapplication decisions](../../test/unit/internal/services/editor/regenerationDecisionService/decideManualEditReapplication_test.go), [GUI-model updates](../../test/unit/app/gui/models/editorState/updateCurrentState_test.go) and [driver updates](../../test/unit/app/gui/drivers/state/updateState_test.go).
- Test each mode both directions with zone-only, connection-only, full and absent snapshots; two-player tournament transitions; both aliases true and removal of just one alias; alias swaps with unchanged effective booleans; simultaneous tournament/arena changes; unrelated rule changes; first generation and failure/retry. Test input clone isolation, dirty and Exit behavior separately.
- Extend [validator tests](../../test/unit/internal/validators/editorStateValidator/validateEditorState_test.go), [handler validation](../../test/unit/internal/handlers/stateHandler/validateEditorState_test.go), [handler load](../../test/unit/internal/handlers/stateHandler/loadState_test.go) and affected facade tests. Enumerate 3-8 and below/above-range counts for both tournament aliases; two is a no-op. Non-tournament behavior stays unchanged; no-fix validation and failed load preserve input/current document.
- Extend [generation handler tests](../../test/unit/internal/handlers/templateHandler/generateTemplate_test.go) and real application-handler integration for both aliases, graph/player/rules consistency, warnings and caller ownership after correction. Keep direct provider compatibility assertions and direct generator 2-8 player tests unchanged; direct GeneratorConfig validation is not this batch's contract.
- Extend [state integration](../../test/integration/editorState_integration_test.go) and appropriate existing Exit/load suites. Load malformed-count tournament settings with a manual snapshot; verify correction, discard, dirty/Exit, warning after actual next-frame regeneration, and unchanged file bytes until explicit save. Save/reload the corrected file; it must load cleanly at two. Valid two-player loads preserve snapshots.
- Exercise General UI via real frame/input integration: checkbox-origin loaded state, victory-selector alias, fixed/disabled player control, uncheck while a selector still activates tournament, leaving tournament then changing count, and next-frame writeback. Use existing gated GUI harness, no global tags or new fake unit exports.
- Run focused unit/integration/GUI suites and `go build ./...`; record exact results. Verify required state warnings are visible, not merely produced by an unused return value.

### Phase Summary

Written plan approved; implementation and fresh baseline deferred to the next session.

## Phase 3: Quality propagation and explicit Custom display

Status: Not started

- [ ] Extend the existing quality-edit handler contract with current zone/connection context; reuse the mutation response. Update interfaces, facade, mocks and direct callers together, with no unrelated DTO changes.
- [ ] Add model-based zone-plus-incident-connection editing in the existing service package. Capture old effective quality/preset before reprofile, resolve new endpoint quality afterward, and change only matching incident guard values on an actual quality change.
- [ ] Clone zones/connections before mutation using existing deep-clone facilities; preserve order, names, endpoint identity and nil/empty semantics. Keep the current zone-only reprofile API and road finalization behavior for other callers.
- [ ] Inject the tier collaborator where needed; regenerate Wire through the existing task. Update explicit test constructors/helpers, never generated source by hand.
- [ ] Consume the whole mutation using the index-preserving pointer replacement described above. Rebind selected, clear `syncedFor`, finish old zone-pointer use before replacing zones, and refresh geometry/property state. Test that editing a connection after a quality mutation affects the retained working collection, not a discarded allocation.
- [ ] Add explicit Custom display and exact-match synchronization. Use a nonnumeric Custom item with safe index mapping; choosing it must not change the guard number.
- [ ] Add unit, handler, GUI and save/reload regressions as changes land.

### Verification Plan

- New service public method gets its own lower-camel test file under the existing mirrored service folder. Extend [zone reprofile tests](../../test/unit/internal/services/connection_editor/zoneEditorService/applyNeutralZoneQuality_test.go), [quality handler tests](../../test/unit/internal/handlers/zoneEditorHandler/applyZoneEditorQuality_test.go), [facade tests](../../test/unit/internal/handlers/guiHandler/applyZoneEditorQuality_test.go) and constructor tests when affected.
- Matrix: all six named presets, with a dedicated Bronze Default 15,000 to Gold Default 25,000 test and ordered first-match assertion; upgrades/downgrades; stronger unchanged opposite endpoint; equal endpoints; reversed endpoint order; hub/player/unknown/nil-quality inference and existing missing-endpoint fallback; Plastic/Bronze unchanged mapping; multiple parallel incident edges; self-edge if supported by the existing API; unrelated edges; empty/nil collections.
- Preserve unmatched zero, negative and arbitrary custom numbers; exact typed numeric matches are presets. Explicitly test the generated Plastic 10,000 exception. Test same-quality/castle-only requests do not recalculate. Whole-connection comparisons must pin type, Road true/false/nil, placement data, GuardZone, weekly increments and all other fields.
- Test source/result mutation isolation for nested zones and connections. No protected type changes, new Clone ownership semantics or production test-only seams.
- Extend [zone property GUI integration](../../test/integration/gui/zoneEditorProperties_integration_test.go): named preset -> quality change -> incident new value; custom value preserved/displayed; exact typed match becomes preset; select Custom without numeric change; reselect/reopen; no stale widgets on subsequent frames; Apply vs Cancel. Use real inputs/visible output with existing harness rather than new broad test exports.
- Add persistence round trip through existing public handlers: edited quality and numeric guard survive save/reload; reopened GUI infers the same label/Custom; no new fields in either persisted format. Use test-owned temporary paths only.
- Run focused unit/integration/GUI suites, build and generated-wiring checks. Check actual GUI labels/value results rather than only nonzero pixels or changed hashes.

### Phase Summary

Pending.

## Phase 4: Whole Batch E verification and owner handback

Status: Not started

- [ ] Verify the combined real flow: manual edits -> arena/tournament change -> new generated graph; invalid tournament load -> correction -> frame/save/reload; pending quality edits -> Apply/Cancel -> regeneration/persistence.
- [ ] Run final build, fresh unit tests, comparable coverage, default/integration/GUI suites, test-layout checker and report-only lint. No blanket snapshot updates or auto-fix across the repository.
- [ ] Compare before/after total coverage and changed-function coverage. Cover every reachable changed branch/logical unit; record genuine GUI-only/unreachable limitations instead of fake unit seams.
- [ ] Obtain independent Claude Opus 5 implementation review limited to Batch E and its regression interactions. Address blockers and rerun affected checks.
- [ ] Inspect diff scope and owner's index read-only; protected trees/output path/opaque raster/retired topology work remain untouched except generated Wire produced by its generator as required.
- [ ] Present owner-facing behavior and verification results. Owner performs all staging/commits and acceptance. Do not mark review findings fixed before the owner's commit protocol is satisfied.
- [ ] Complete Final Recap and Deployment Plan, and update this plan with remaining blockers if work stops partway.

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

Pending.

## Verification ledger

| Check | Baseline | Final |
| --- | --- | --- |
| Build / unit / coverage | Not run; next-session baseline required before edits | Pending |
| Focused domain / handler / GUI / persistence | Source inspected only | Pending |
| Default / tagged integration / GUI | Not run | Pending |
| Layout / report-only lint / formatting / Wire | Not run | Pending |
| Native Linux execution | Not run | Pending or explicitly unavailable |
| Independent plan review | Claude Opus 5 APPROVED revised plan, 2026-09-12 | N/A |
| Independent implementation review | N/A | Pending |

Discovery tooling notes: initial subagent requests used unsupported lowercase model IDs; retrying with exact display name `GPT-5.6 Terra (copilot)` succeeded. Two guessed source paths did not exist; confirmed source/search results, not those guesses, inform this plan. No repository changes resulted from those failed reads/calls. Existing tooling memory already documents exact model-name requirements.

## Final Recap

Pending implementation and verification. Discovery resolved all requested product choices; only this plan was created. The revised plan passed independent review and the owner approved subsequent-session implementation. No claim of fixed findings, successful tests or runtime behavior is made yet.

## Deployment Plan

No deployment is authorized or performed in this planning session. After implementation approval and successful verification:

1. Owner reviews the scoped changes, warnings and UI behavior; performs staging/commits and any release steps.
2. Build/package through the existing platform workflow. No dependency installation, saved-schema migration or output-directory change is planned; Wire regeneration is build-time only if required.
3. On application launch, retain per-machine game-directory detection. Existing invalid tournament settings are corrected in memory with warning/dirty state; the file changes only on explicit Save.
4. Communicate that mode transitions discard manual edits and exiting tournament leaves two players. Guard Custom values remain numeric and are not given hidden persistent preset identities.
5. Record platform-specific verification and unresolved acceptance items honestly. Do not infer new engine results from Batch C's completed owner report.
