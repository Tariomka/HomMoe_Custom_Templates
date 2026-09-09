<!-- markdownlint-configure-file { "MD024": { "siblings_only": true } } -->

# Batch A: manual-state safety and export-directory detection

Fix review §1.1 and §1.2, with owner-approved inclusion of §1.8 and §1.9. Prevent silent loss of committed manual edits and exports to an unrelated fallback directory, without changing generation policy or the persistence format.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../../AGENTS.md), [the handoff](../session-carry-forward.md), and
[the review](../backlog/review-gpt-6-astra-09-07.md). The review remains the backlog
source of truth. Preserve finding numbers. Owner commits before items are marked
fixed. Never stage, unstage, commit, push, switch branches, or modify protected
data/schema/registry files. Never hand-edit generated Wire output.

## Approval and starting state

- Scope confirmed by owner on 2026-09-09. **Implementation approved 2026-09-09.**
  Owner staged the original plan for diff review; preserve the index exactly.
- Owner selected §1.1, §1.2, §1.8, and §1.9 together on the current branch.
- Branch: `AD/save_and_pathing_silent_bug`; HEAD:
  `d413c03c3424c747c78b82efc08c305c5235684c`, “Updated Backlog for future sessions (#40)”.
- Index and working tree were clean before creating this plan. The previous
  handoff's four documentation edits are now committed. No application/test
  changes between reviewed `f4f4cf63f22e84040754a7231b4d1dc793070af1` and HEAD.
- Prior measured Windows baseline: 74.4% unit statement coverage, zero lint
  issues. These are historical measurements until Phase 2 remeasures them.
- No code changes or tests have been run in this session. Read-only source,
  caller, test, composition, and Git inspection completed. One guessed source
  filename did not exist; the actual caller was found in LayoutPanel.

### Confirmed behavior

1. Dirty means the committed, persistable manual snapshot changed, not merely
   that Apply was clicked. A first snapshot is a change even if visually equal
   to the generated layout. Clearing an existing snapshot is a change.
2. Identical Apply leaves a clean document clean; it never clears existing
   dirtiness or resets an armed Exit confirmation. A changed commit sets dirty
   and resets Exit confirmation.
3. A no-template apply or `ErrProvidedTemplateInvalid` rejection changes neither
   the committed snapshot nor dirty/Exit state. Preserve the existing behavior
   of accepting an updated template accompanied by validation errors, such as
   `ErrZonesMissing`, and report those errors as before.
4. **Explicit owner decision:** every Apply consumes/discards `pendingBaseZones`,
   including rejected/no-template applies. Do not retain it for retry. Compute
   the untouched-revert predicate from the captured base and unmodified request
   before calling the handler; commit a clear only after acceptance.
5. Deep-clone zones on both `SetManualEdits` ingress and `GetManualZones` egress
   using existing `Zone.Clone`. Preserve existing nil/empty storage semantics.
6. Detection belongs to `PathResolutionService` behind
  `IPathResolutionService`; the existing filesystem handler remains a
  two-argument constructor and delegates the lookup. Keep per-launch
  detection and the existing detector algorithm, with no private test export.
7. Any failed/empty detection leaves output empty and shows an error directing
   the user to the existing folder picker. Do not open the picker automatically.
   An error accompanied by a path does not authorize export.
8. Browsing starts can fall back to an existing directory, but browsing,
   cancellation, and Reveal do not select an export destination. Only successful
   detection or explicit session-only folder confirmation does. Never persist
   output paths. Successful `.rmg.json` and preview `.png` export share the selected
   destination; export failure writes neither when no destination is selected.

### Boundaries

- Do not fix road generation (§1.4/§1.6/§1.10), modes, DTOs, compact state,
  topology retirement, or unrelated dialog ownership in this batch.
- Preserve automatic regeneration's exclusion of manual fields. Do not use
  `WasStateChanged` as manual dirty equality or change save/reset/load policy.
- Do not normalize stored collections, reorder snapshots, or change JSON tags.
  For equality only, top-level zero-length manual collections are equivalent
  where the existing serializer omits both. Nested values retain their existing
  distinctions; no generalized JSON canonicalization is proposed.
- `Connection.Length` IS persisted by `ToConnectionEntity`; include it, all
  other saved connection fields, `IsUserAdded`, zone fields, quality, generator
  stamps, and manual positions in equality. Do not compare UI revision/status.
- No new test-only production API, global test tags, global environment mutation
  in parallel tests, blanket lint fixes, or golden-image bulk updates.

## Phase 1: scope, durable plan, independent review

Status: Complete

- [x] Read repository instructions, handoff, review scope and baselines.
- [x] Inspect current Git state and compare application/test changes since review.
- [x] Reverify target sources, callers, existing tests, and composition.
- [x] Obtain owner decisions and full scope confirmation.
- [x] Write this durable plan without implementing application changes.
- [x] Independently review this plan and address actionable feedback.
- [x] Obtain explicit owner approval of the reviewed plan before Phase 2.

### Verification Plan

- Inspect Git diff: only this plan may be added by this session at this stage.
- Review must check rejection versus accepted warnings, pending-base consumption,
  compare-before-mutation, nil/empty behavior, snapshot ownership, production DI,
  deterministic failure tests, and no export from a browsing-only directory.

### Phase Summary

Owner expanded A to include both manual-lifecycle findings and selected a
composition-owned detector. Owner explicitly rejected preserving the pending
base on failed Apply. Claude Opus 5 independently reviewed the draft and found
three gaps: GUI injection, explicit comparison representation, and real-handler
no-op coverage. All three are addressed below. A second independent review of the
revised plan returned “ready for owner implementation approval” with no blockers.
Markdown diagnostics are clean. Git inspection showed only this new untracked
plan and no staged changes before owner action. Owner subsequently staged this
plan and approved implementation. No agent staging/unstaging is permitted.

## Phase 2: fresh baseline and focused regression tests

Status: Complete

- [x] Recheck index/working tree before editing; preserve any new owner changes.
- [x] Run the existing code-coverage task and record Go's `cover -func` total.
  Record toolchain/platform and per-function coverage for changed logic. If the
  fresh baseline differs from 74.4%, investigate before applying code changes.
- [x] Run build and report-only lint baseline. Do not use auto-fix lint.
- [x] Add failing regressions for §1.1/§1.8/§1.9 through existing public APIs.
- [x] Record the expected failures before fixing them. Add deterministic startup
  detection regressions in Phase 4 through the approved path-resolution contract.

### Verification Plan

- `go build ./...`
- Existing “Go: Generate code coverage report” task, then
  `go tool cover -func=coverage.txt`; comparable total at least 74.4%.
- Existing “Go: Get Linter Results” task; zero issues expected.
- Focused unit regressions for State manual Apply/Exit and EditorState Set/Get
  must fail for the expected defects, not setup/compilation errors.
- Use `-count=1` for the first measurement; avoid unnecessary uncached repeats.
  Keep before/after totals in this plan even when generated reports are replaced.

### Phase Summary

Fresh baseline on Go 1.27.0 windows/amd64: `go build ./...` passed; the coverage
task reported **74.4%**; report-only lint reported **0 issues** (three existing
unused-exclusion warnings). Focused public-API regressions failed before their
respective fixes, including nested snapshot mutation, dirty/Exit handling,
rejected Apply, real-handler revert with roads disabled, and detection failure.
The completed coverage task later measured **74.5%**.

## Phase 3: manual commit tracking, revert order, snapshot isolation

Status: Complete

- [x] In [EditorState](../../app/gui/models/editorState.go), deep-clone stored and
  returned manual zones using the existing clone primitive and preserve nil vs
  empty. Keep connection snapshot conversion behavior unchanged.
- [x] Have `SetManualEdits` and `ClearManualEdits` return whether their committed
  manual data changed, comparing the existing stored model snapshot to the
  candidate snapshot. Existing callers may ignore the result. Keep comparison
  local to manual data, rather than cloning/comparing the entire settings state.
- [x] Put persisted-field equality in the internal editor-state model as a public
  model operation with dedicated unit tests, rather than teaching GUI code the
  persistence schema. Project candidate connections through the existing
  `ToManualConnectionSaves` before comparing with stored connection saves; this
  includes `Length`, `Road`, placement rules, and `IsUserAdded`. Compare zones over
  the current model fields written by `ToManualZoneSaveEntities`, including
  generator stamps, manual positions, and quality. Do not perform entity conversion
  in GUI code or add a new conversion boundary. Pin the current field projection
  with tests so future model-only fields cannot silently alter dirty semantics.
- [x] In [manual Apply](../../app/gui/drivers/stateManualEdits.go), capture and
  consume the pending base as today, then calculate untouched-revert before
  mutation. Return early on no template.
- [x] Have `handleUpdateTemplate` report acceptance, retaining the current invalid
  template rejection and accepted-validation-warning behavior. Reapplication
  callers retain behavior and can ignore that result. The regeneration clear in
  [stateGeneration.go](../../app/gui/drivers/stateGeneration.go) intentionally
  ignores `ClearManualEdits`'s new return value; scalar `UpdateState` owns its dirty
  transition. Do not add a second dirty transition there.
- [x] Only after acceptance, clear or set the snapshot. Only a changed snapshot
  sets `unsaved=true` and `confirmExit=false`. Do not suppress existing template
  replacement/revision/status behavior for an otherwise accepted no-op.
- [x] Extend dedicated tests for the changed public methods and existing callers.
- [x] Add real-handler revert regression with separately copied top-level zone
  slices and roads-off generation so post-handler mutation cannot hide the bug.
  Do not fix or assert correctness of roads-off reconstruction in this batch.

Scope limit: connection conversion currently shares `Road *bool`; §1.9 is the
zone ingress/egress fix, not a promise of complete connection isolation. Current
GUI callers do not mutate through that pointer. Leave connection clone behavior
unchanged and record that residual boundary at handoff, rather than widening scope.

### Verification Plan

Extend these existing tests, adding one assertion/unit per named AAA test:

- [Apply](../../test/unit/app/gui/drivers/stateManualEdits/applyEditedZones_test.go):
  first snapshot dirty; zone-only and connection-only changes; identical snapshot
  clean after successful save; changed/unchanged accepted warning updates;
  untouched/edited revert; no-template/rejected applies preserve committed state;
  rejection after a pending revert consumes the base as owner requested; a mock
  mutating request roads cannot change precomputed revert identity.
- [Exit](../../test/unit/app/gui/drivers/stateFiles/exit_test.go): first Exit warns,
  second exits, changed Apply/clear re-arms warning, identical/rejected Apply does
  not reset an armed confirmation. Existing scalar behavior remains green.
- [SetManualEdits](../../test/unit/app/gui/models/editorState/setManualEdits_test.go),
  [GetManualZones](../../test/unit/app/gui/models/editorState/getManualZones_test.go),
  and [ClearManualEdits](../../test/unit/app/gui/models/editorState/clearManualEdits_test.go):
  nil/empty, equality result, input/output isolation for positions, quality/ring
  pointers, roads, content pools, main objects and nested reference slices.
  Include persisted `Length`/`IsUserAdded` connection changes. Reuse established
  clone fixtures/guards; never write fake tests for dead conversion helpers.
- [Revert integration](../../test/integration/zoneEditorRevertToBase_integration_test.go):
  real-handler untouched revert clears snapshot even when request roads mutate;
  edited revert retains snapshot. Keep production-API tests untagged.
- Real-handler no-op regression: Apply -> successful editor-state save -> Apply
  the same effective payload -> `IsUnsaved()` remains false, proving post-rebuild
  idempotence instead of relying solely on nonmutating mocks. Use the existing
  gated save integration seam in a separately tagged file if needed, not a new
  unit-only export. If road rebuilding actually changes the committed payload,
  record the failure and ask the owner before expanding into road-policy fixes;
  never hide a persisted change merely to force this assertion to pass.
- Existing save/load integration verifies dirty clears only after successful
  save and saved manual layout survives reload. Use existing gated APIs only in
  integration tests when a clean saved fixture cannot use a public unit path.
- Run focused suites; record exact results and branch coverage gaps.

### Phase Summary

Implemented and verified. Persisted-manual equality is a public model helper;
`SetManualEdits`, `ClearManualEdits`, and `ApplyEditedZones` now transition dirty
and Exit confirmation only when an accepted commit changes the snapshot.
Untouched revert identity is captured before handler mutation, while every Apply
consumes the pending base. Zone snapshots clone at ingress and egress. Focused
regressions include the accepted-warning no-op retaining an armed Exit
confirmation, real-handler save-repeat no-op, roads-off revert, and nine nested
zone mutation cases. The connection `Road` pointer/placement alias remains an
explicit out-of-scope residual boundary.

## Phase 4: composition-owned detection and explicit picker recovery

Status: Complete

- [x] Add `FindGameTemplateDirectory` to the existing `PathResolutionService` and
  `IPathResolutionService`; it delegates to `helpers.FindOldenEraTemplatesDir(false)`.
  No standalone game-template-directory service, interface, mock, or dedicated
  tests remain. No new search algorithm, caching, persistence, or
  environment-dependent UI code.
- [x] Add lookup to the existing
  [filesystem handler contract](../../internal/handlers/handler_interfaces/fileSystemHandlerInterface.go)
  and [implementation](../../internal/handlers/fileSystemHandler.go), delegating
  through the existing path-resolution collaborator. Keep the two-argument
  filesystem-handler constructor and the facade delegation-only.
- [x] Update explicit path-resolution mocks/fixtures and constructor callers;
  keep NewWindow/NewUIState signatures unchanged. Regenerate Wire, never edit
  its generated output. The regenerated output is identical to HEAD; there is
  no provider or generated-file diff.
  Known fixture fanout:
  [real filesystem fixture](../../test/test_helpers/fileSystemHandler.go),
  [handler mock](../../test/test_helpers/fileSystemHandlerMock.go), and
  [handler unit setup](../../test/unit/internal/handlers/fileSystemHandler/common_test.go).
- [x] In [NewUIState](../../app/gui/drivers/state.go), call the filesystem handler
  when lookup is enabled. Set output only for a successful nonblank result.
  On not-found, empty success, or other error, leave it empty and set an actionable
  error. Suggested copy: “Game template directory not found. Choose the game
  templates folder using the output folder picker before exporting.” Preserve
  useful underlying error details for other failures.
- [x] Do not invoke browsing-directory resolution to populate output on failure.
  Existing `ResolveStartDirectory` already supports empty starts; no picker
  production change is currently needed. Keep callback behavior; only change it if recovery tests
  expose a directly relevant defect, and document the narrow change.
- [x] Add deterministic constructor/handler regressions using the real public
  contracts and testify mocks, with no process environment changes.
- [x] Add GUI recovery tests using real input, fixture-only directories, and the
  existing explorer harness. Extend test-side harness controls only as necessary;
  do not introduce production-only setters for output or dialog confirmation.
- [x] Add test-side `NewAppRunnerWithFileSystem` in the existing
  [runner](../../test/test_helpers/integration_common/appRunner.go), retaining its
  existing test tags. Keep the default constructor delegating to the real
  composition root. Recovery tests inject an `IFileSystemHandler` decorator that
  fails lookup, uses an existing fixture directory for the empty browsing start,
  and delegates actual browsing/writing behavior to real collaborators. This is
  test infrastructure, not a production override or environment mutation. Never
  let the recovery test start with the host's real detected game destination.

### Verification Plan

- Path-resolution and handler tests exercise the real public contracts with
  deterministic mocks. Existing platform fixture tests cover helper discovery;
  do not invent a second seam or a standalone wrapper merely to test delegation.
- [Startup tests](../../test/unit/app/gui/drivers/state/newUIState_test.go): lookup
  skipped (never called), success, wrapped not-found, arbitrary error, blank result
  without error, and nonempty result accompanied by error. Assert output and status
  independently. Remove the current host-dependent unsafe fallback expectation.
  Path-plus-error is a handler contract test only: the current detector returns
  an empty path on its actual error branches; do not change it to manufacture one.
- [Export handler tests](../../test/unit/internal/handlers/templateHandler/saveTemplate_test.go):
  empty/whitespace output refuses before preview rendering or file service calls;
  success passes the chosen destination. These characterize existing correct
  behavior, not a new export-handler fix. No handler implementation change expected.
- [Explorer GUI suite](../../test/integration/gui/fileExplorerDialog_integration_test.go):
  deterministic failed lookup -> unset output -> open picker at usable fixture
  browsing start -> cancel leaves unset -> explicit folder confirmation selects
  only that session destination. Assert no JSON/PNG is written before selection;
  after selection both outputs land together in the fixture. Confirm a fresh
  session/save-load cycle does not persist the override. Never write into the
  developer's actual detected game directory from tests.
- Run “Go: Generate wire injectors”; generated graph must compile.
- Existing tests using constructor mocks must be updated explicitly, without
  silently giving every lookup a fallback destination.

### Phase Summary

Implemented and verified. `FindGameTemplateDirectory` lives on
`PathResolutionService`/`IPathResolutionService`; the filesystem handler
delegates with its unchanged two-argument constructor. Failed, blank, and
path-plus-error detection leave output unset and require explicit session-only
picker confirmation. Eight pre-fix detection regressions failed as expected;
eight GUI recovery tests use real input and fixture directories with no golden
updates. No output path is persisted, and no export occurs before selection.

## Phase 5: full verification and owner handoff

Status: In progress

- [x] Format only explicit changed permitted Go files reported by `gofmt -l`.
- [x] Run build, unit tests, coverage, lint, layout, and relevant integration/GUI
  suites. Fix batch regressions; distinguish environmental blockers from passes.
- [x] Verify comparable Windows unit coverage is at least both the fresh baseline
  and 74.4%, with new logic covered. Check Go's own coverage output, not raw row
  summation. Record any unavoidable platform/GUI gaps honestly.
- [x] Independently review implementation and tests against this approved plan.
- [x] Inspect final diff and index; confirm protected paths untouched, no global
  tags, no persisted path, no unrelated changes, and owner staging preserved.
- [x] Update this plan with exact commands/results and residual limits. In the
  review, record “implemented/verified, awaiting owner commit” only when true.
- [x] Owner committed `a9eb35c`; commit and clean starting Git state verified read-only.
- [ ] Correct and test the committed empty-snapshot revert fallthrough under §1.8,
  then verify the owner's follow-up commit before closing that finding and this plan.

### Verification Plan

- `go build ./...`
- `go test ./test/unit/... -count=1` for first full plain unit run; use cached
  repeats where applicable.
- Coverage task and `go tool cover -func=coverage.txt`.
- `go test ./test/...` (default tags; no automatic GPU tests).
- `go test -tags=integration_test ./test/integration/...`
- `go test -tags='integration_test,gui' ./test/integration/...`
- `go vet -tags=integration_test ./...`
- “Go: Get Linter Results” and “Go: Check test build-tag layout” tasks.
- `git diff --check`, diff review, status and staged-diff inspection.
- Local Linux, full race, benchmarks and vulnerability scans are not implied by
  these Windows checks. Use available Linux CI evidence after owner submission;
  do not claim Windows success proves Linux/GPU/race outcomes.

### Phase Summary

Pre-commit implementation verification completed with the results below. Owner
committed `a9eb35c` on 2026-09-09. Read-only closure inspection found the committed
Apply rewrite differs from the tested implementation: combining untouched-revert
selection with `ClearManualEdits()`'s changed result allows an empty-snapshot
revert to fall through and store the base. §1.8 closure is blocked; §1.1, §1.2,
and §1.9 are marked fixed. The working tree and index were clean before these
documentation updates. No post-commit tests or application edits were performed.

Historical pre-commit verification: Windows
Go 1.27 checks passed: `go build ./...`; `go test ./test/unit/...`; coverage task
at **74.5%**; `go test ./test/...`; tagged integration; tagged GUI integration
(GUI 27.261s, main 2.734s); `go vet -tags=integration_test ./...`; and the
test-layout check. Final lint reported 0 issues after two local funcorder fixes,
with the three baseline unused-exclusion warnings. The changed public equality/
model helpers, Set/Get/Clear manual state, Apply, commit, NewUIState/detection
helper, handler, and `FindGameTemplateDirectory` have 100% statement coverage.
Claude Opus 5's independent read-only review found no blockers. No Linux, race,
benchmark, or vulnerability-scan result is claimed.

## Final Recap

Owner commit `a9eb35c` contains this batch. §1.1, §1.2, and §1.9 are fixed.
§1.8 and final plan completion remain open because an accepted untouched revert
with no previous snapshot now stores a snapshot. Fix the fallthrough without
changing the owner-approved consume-pending-base policy, add the missing
regression, remeasure coverage and verify, then await the owner's follow-up commit.
The earlier verification ledger must not be read as verification of this later
committed rewrite. Protected data/schema/registry and output persistence remain
unchanged.

## Deployment Plan

1. Owner review and commit completed in `a9eb35c`; verified read-only.
2. Correct §1.8's empty-snapshot revert fallthrough and verify the follow-up
  before declaring this plan complete. No such correction has been applied yet.
3. Owner commits the correction; verify that commit and mark §1.8 fixed.

No data migration or output-path preference change is required. Each launch
redetects the game folder; if unavailable, users must explicitly select the
correct folder for that session before export.
