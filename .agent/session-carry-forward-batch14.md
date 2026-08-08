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
- `plans/extract-regeneration-policy.md` — the working plan. **Deleted by the
  owner in `51e5858` once the batch landed**, so this document and repository
  memory are now the only record; everything load-bearing from it is inlined in
  §7 below.

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

Branch **`AD/refactoring-07-21`**, working tree **clean** and in sync with
`origin`. Nothing was staged or committed by the agent; the owner reviewed and
landed the batch as `8bf9c95` "Batch 14 Done", then `51e5858` "Cleaned up backlog
and plans" (which deleted `plans/extract-regeneration-policy.md` and trimmed the
two `todo/` files).

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

### Proven non-defects — do NOT "fix" these

Each was investigated during Phase 3 triage and shown to be correct as written.
They look like bugs on a quick read, so they will attract a future agent:

1. **The superseded `op.InvalidateCmd` wakeup.** Gio exposes no cancellation for
   a pending invalidate, so a re-armed debounce leaves the old wakeup queued.
   The entire cost is one wasted frame.
2. **A stale `applyNextStateAt`.** Unreachable: the only branch that reads
   `DebounceDueAt` requires `Next != nil`, and the sole branch that sets `Next`
   refreshes the deadline in the same call.
3. **`NextStateLeave` vs `NextStateClear` on first generation.** Equivalent —
   every path through `handleGenerateTemplate` reaches `SnapshotCurrentState()`,
   which nils `next`.
4. **`zoneEditorHandler`'s apparently-live dead methods.** Already deleted, but
   note the trap: `ComputeHasErrors` / `RebuildZoneConnectionRoads` looked used
   under text search only because identically-named **service** methods are the
   live ones (`templateHandler` calls `IConnectionEditorService` /
   `IZoneEditorService` directly).

## 8. Open questions

None blocking. The two long-standing questions from the previous hand-off are now
**resolved**: the dead `zoneEditorHandler` methods were deleted at the owner's
instruction. The `internal/helpers/io.go` filesystem-seam question was dropped —
owner decision 9 (helpers imports from `app/` are legitimate) removed the
layering motivation, leaving only a testability argument that no longer has a
caller asking for it.

## 9. Next recommended actions

1. **Start Batch 15 in a fresh session**, driven by
   [plans/zone-editor-state-extraction.md](../plans/zone-editor-state-extraction.md).
   Scope was agreed at the end of this session and the plan is written; Phase 0
   (characterization tests) is the entry point.
2. Backlog candidates surfaced this session:
   - `handleGenerateTemplate` sits at 81.0 % coverage (pre-existing; the gap is
     the generation-failure and warning-status branches).
   - The road-distance consolidation item, now that its false "UI uses services
     directly" half has been struck.
   - The "Save As → Save To" UI honesty item.

### Batch 15 scope, as agreed

Review §2.6 — nine files, ~2 800 LOC. Owner decisions, all recorded in the plan:

- **All nine files** in scope: the five `zoneEditor*` files plus
  `bonusPickerDialog.go`, `pickerDialog.go`, `ruleDialog.go`, `zoneContent.go`.
- **Geometry to a full service** (interface + `wire` DI, as implementation #4 in
  `internal/services/connection_editor/`); **selection/drag state stays in the
  GUI**, consolidated. The review's literal wording said to move selection too;
  that half was deliberately declined — `dragPos`/`pressPos`/`addMode` are
  pointer-event view state, and a service holding them inverts the layering.
- **Extract everything identified** in the four extra dialogs (~270 LOC).
- **Characterization tests first** — two of the nine files have no test at all.
- **§0.2 is fixed by relabelling** "Reset to generated", not by retaining a
  pristine generated template (owner declined the second copy).
- **After a revert, the following Apply clears the manual-edit snapshot
  outright**, with the consequence accepted: the live template keeps showing the
  reverted-to edits until the next regeneration, then they vanish.

Two facts found while scoping that contradict existing documents: the dialog has
**67** fields, not the ~58 the review states (26 of them Gio widget handles that
cannot move), and `zoneContent.go` lives in `app/gui/dialogs/`, not
`app/gui/panels/` as both the review and the old backlog claimed.

## 10. Carry-forward prompt

> Read `AGENTS.md` first. Key hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/`; keep everything
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chains with `;`
> never `&&`); every change ships with tests and must not drop coverage below
> **69.3 %**; durable multi-session work gets a plan file under `plans/`; never
> stage and never commit — the owner reviews and commits; never change where
> `.rmg.json` is written and never persist the output directory (§2.6).
>
> Batch 14 is **complete, verified, committed and pushed** (`8bf9c95` +
> `51e5858` on `AD/refactoring-07-21`) — review §2.2 is marked `✅ FIXED` in
> place in `todo/review-opus5-08-04.md`, and the `app/`-layering backlog item is
> closed and enforced by depguard. All 46 review findings are now resolved except
> **§2.6**, which is the next batch.
>
> **Start Batch 15 from `plans/zone-editor-state-extraction.md`.** Scope and all
> seven owner decisions are already settled and recorded there — do not
> re-litigate them; begin at Phase 0 (characterization tests), because two of the
> nine files in scope have no test of any kind and the whole refactor depends on
> that safety net.
>
> Batch 14's working plan was deleted once it landed. `§7` of
> `./.agent/session-carry-forward.md` is now the record of what must **not** be
> "fixed" by a later agent; the architectural facts are in repository memory.
