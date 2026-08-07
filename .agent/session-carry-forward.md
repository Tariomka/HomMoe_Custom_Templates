# Session carry-forward — Batch 14 (review §2.2, extract regeneration policy)

## 1. Session goal

Remediate **§2.2 "Generation and manual-edit policy lives in the GUI driver"** from
[todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) — the last open item in
that review requiring an owner scope decision — and close the overlapping
`app/`-layering backlog item.

## 2. Fixes applied

- Regeneration + manual-edit-reapplication policy moved out of the GUI into a
  pure service: [internal/services/editor/regenerationDecisionService.go](../internal/services/editor/regenerationDecisionService.go).
- [app/gui/drivers/stateGeneration.go](../app/gui/drivers/stateGeneration.go) —
  `AutoRegenerate` reduced to a dispatcher.
- [app/gui/models/editorState.go](../app/gui/models/editorState.go) — nine
  decision methods removed; it is a snapshot store again.
- Three genuinely dead exported methods pruned (see §3).
- Pre-existing committed over-indentation in
  [test/integration/manualCastleReapply_integration_test.go](../test/integration/manualCastleReapply_integration_test.go)
  fixed with `gofmt -w`.
- Two inaccurate [todo/backlog.md](../todo/backlog.md) items corrected/closed.

## 3. Features added / changed

- **`IRegenerationHandler` seam** — new flat two-method handler interface plus
  `RegenerationSet` / `InitializeRegenerationHandler`. Flat (not embedding the
  service interface) so `app/gui` never imports `internal/services`, which
  depguard forbids.
- **Deterministic debounce** — the 300 ms window moved into the service, with
  `now` as a parameter. No clock seam needed; timing tests are exact.
- **`ManualEditDecisionDto` collapsed to one nullable pointer**
  (`ReapplyWithCastleChanges *CastleSettingChanges`, nil = drop the edits).
  Owner's suggestion. The previous `bool` + value pair made two illegal states
  representable and relied on a doc comment to prevent them.
- **Anti-aliasing accessors** — `EditorState.GetPreviousState()` /
  `GetNextState()` return pointers to *copies*, so callers cannot mutate the
  snapshots that change detection compares against (review §1.3 bug class).
- **depguard widened** — `no-services-from-app` now also denies
  `internal/repositories`, `internal/mappers`, `internal/validators` from
  `app/**`. Pure regression-proofing; produced zero new findings.
- **Dead code pruned**: `EditorState.ResetPreviousState`,
  `zoneEditorHandler.ComputeHasErrors`, `zoneEditorHandler.RebuildZoneConnectionRoads`.

## 4. File modifications

**Created (13)**
- `internal/dtos/nextStateAction.go` — `NextStateAction` enum.
- `internal/dtos/regenerationDecisionRequestDto.go`, `regenerationDecisionDto.go`,
  `manualEditDecisionDto.go`.
- `internal/services/editor/regenerationDecisionServiceInterface.go` +
  `regenerationDecisionService.go`.
- `internal/handlers/handler_interfaces/regenerationHandlerInterface.go` +
  `internal/handlers/regenerationHandler.go`.
- `test/test_helpers/regenerationHandler.go` (real handler, not a mock).
- `test/unit/internal/services/editor/regenerationDecisionService/{decideRegeneration,decideManualEditReapplication}_test.go`.
- `test/unit/app/gui/models/editorState/{getPreviousState,getNextState}_test.go`.
- `test/unit/internal/composition/wire/initializeRegenerationHandler_test.go`.
- `plans/extract-regeneration-policy.md` — the authoritative plan, with full
  per-phase summaries.

**Modified (key)**
- `app/gui/drivers/state.go` (new `regeneration` field, `autoRegenDebounce` const
  removed), `stateGeneration.go`, `app/gui/models/editorState.go`,
  `app/gui/editor/window.go`, `app/gui/program.go`.
- `internal/composition/{providerSets.go,wire.go,wire_gen.go}` (regenerated).
- `internal/handlers/zoneEditorHandler.go` (two dead methods removed).
- `.golangci.yml` (depguard widening).
- ~31 `drivers.NewUIState` call sites across unit + integration tests, plus
  `test/test_helpers/integration_common/appRunner.go`.
- `todo/review-opus5-08-04.md`, `todo/backlog.md`, `todo/test_observations.md`.

**Deleted (10)** — 9 obsolete `test/unit/app/gui/models/editorState/` tests whose
subjects moved to the service, plus `resetPreviousState_test.go`.

## 5. Tests added or updated

16 new service tests (9 decision + 7 manual-edit), 9 accessor tests, 1 wire
smoke test. All suites green at hand-off:

| Suite | Result |
| --- | --- |
| `go test ./test/unit/...` | no FAIL |
| `go test -tags=integration_test ./test/integration/...` | `ok … 2.782s` |
| `go test -tags 'integration_test,gui' ./test/integration/gui/...` | `ok … 3.388s`, zero snapshot diffs |
| `go test "-bench=." ./test/performance` | `ok … 0.199s` |
| Total unit coverage | **69.3 %** — equal to the Phase 0 baseline |
| `golangci-lint-v2 run ./...` (Windows **and** `GOOS=linux`) | `0 issues.` |
| `go run ./cmd/testlayoutcheck .` | passed |
| `wire diff ./internal/composition/...` | exit 0 |

## 6. Git status snapshot

Branch **`AD/refactoring-07-21`**. Nothing was staged or committed by the agent.
The owner had already staged most of the batch (`M ` entries); the agent's Phase 4
edits are unstaged (` M .golangci.yml`, ` M todo/backlog.md`,
` M test/integration/manualCastleReapply_integration_test.go`, plus `MM` on
`todo/review-opus5-08-04.md` and `test/integration/stateSaveAs_integration_test.go`).

Permanent `gofmt -l .` noise remains on `app/gui/drivers/dialogHost_testexports.go`
and `internal/composition/wire.go` — both **CRLF-only**, invisible in the committed
blob. Do not "fix" them; it only creates empty `git diff HEAD` churn.

## 7. Rejections / things declined

- **A defensive clear of `applyNextStateAt`** — written, then reverted. The stale
  state it guarded is unreachable, and it could not be tested without a
  test-only seam. **Do not re-add it.**
- **Making the reapply ordering structural** by threading the decision through
  `applyGeneratedTemplate`'s signature — declined as signature churn; the DTO
  and an explicit comment already carry it.
- **Shrinking `autoRegenerate_test.go`** (the plan called for it) — declined
  after reading: those 8 tests already assert *driver* outcomes, not decision
  logic, so they were already the dispatch tests the plan wanted.
- **Moving `WasStateChanged()`** off `EditorState` — reverted. It drives the
  unsaved marker, not generation.

## 8. Open questions

None blocking. The two long-standing questions from the previous hand-off are now
**resolved**: the dead `zoneEditorHandler` methods were deleted at the owner's
instruction. The `internal/helpers/io.go` filesystem-seam question was dropped —
owner decision 9 (helpers imports from `app/` are legitimate) removed the
layering motivation, leaving only a testability argument that no longer has a
caller asking for it.

## 9. Next recommended actions

1. Review and commit this batch.
2. **§2.6** is the only review item left: `zoneEditorDialog`'s ~58 fields and the
   oversized files (`zoneEditorDialog.go` 507 LOC, `zoneEditorCanvas.go` 479,
   `bonusPickerDialog.go` 434, `pickerDialog.go` 371, `ruleDialog.go` 314,
   `zoneContent.go` 299). Explicitly low priority — do it opportunistically the
   next time the zone editor is touched.
3. Backlog candidates surfaced this session:
   - `handleGenerateTemplate` sits at 81.0 % coverage (pre-existing; the gap is
     the generation-failure and warning-status branches).
   - The road-distance consolidation item, now that its false "UI uses services
     directly" half has been struck.
   - The "Save As → Save To" UI honesty item.

## 10. Carry-forward prompt

> Read `AGENTS.md` first. Key hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/`; keep everything
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chains with `;`
> never `&&`); every change ships with tests and must not drop coverage below
> **69.3 %**; durable multi-session work gets a plan file under `plans/`; never
> stage and never commit — the owner reviews and commits; never change where
> `.rmg.json` is written and never persist the output directory (§2.6).
>
> Batch 14 is **complete and verified** — review §2.2 is marked `✅ FIXED` in
> place in `todo/review-opus5-08-04.md`, and the `app/`-layering backlog item is
> closed and enforced by depguard. All 46 review findings are now resolved except
> §2.6 (zone-editor size/field count), which is deliberately opportunistic.
>
> Full detail of what was done and why — including four proven non-defects that
> must NOT be "fixed" by a later agent — is in
> `plans/extract-regeneration-policy.md`. The complete hand-off is
> `./.agent/session-carry-forward.md`.
