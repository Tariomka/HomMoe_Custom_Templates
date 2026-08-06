# Session Carry-Forward

## 1. Session goal

Finish **Batch 11** of the 46-finding review in
[todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) — findings **§6.2**
(`internal/handlers` has no mirrored unit tests) and **§6.4** (two
`app/gui/constants` catalogues at 0% coverage), after the owner-approved
expansion of first converting **every constructor-injected service under
`internal/`** to an interface so the handlers can be tested against `testify`
mocks.

**Batch 11 is COMPLETE** (owner-reviewed and committed) and **Batch 12 is
CLOSED as `❌ WILL NOT FIX`**.
[plans/batch-11-handler-coverage.md](../plans/batch-11-handler-coverage.md) is
the source of truth for Batch 11 and now carries per-phase summaries, a Final
Recap and a Deployment Plan.

### Batch 12 — §1.8, rejected (2026-08-06)

The owner rejected the finding outright; no code was written. The rationale is
recorded in place at §1.8 of
[todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) and §12 item 12 is
marked closed. Summary: the output directory is a **hard requirement of the
game**, not a user preference — Olden Era only reads templates from its own
templates directory, so a file written anywhere else is never found and the user
cannot locate it afterwards. It is also a per-*machine* value (Steam library, OS,
Proton layout), so a persisted path is invalid on any other device and goes stale
when the game moves. Per-launch auto-detection via
`helpers.FindOldenEraTemplatesDir` is self-healing and is the correct design; the
picker is a single-session escape hatch. **Do not re-propose this.**

## 2. Fixes applied

- Removed the **last** `wire.Bind` from
  [internal/composition/providerSets.go](../internal/composition/providerSets.go).
  The codebase now has **zero** `wire.Bind` calls, because every provider already
  returns an interface. `wire_gen.go` regenerated.
- [internal/services/content_rules/contentRuleService.go](../internal/services/content_rules/contentRuleService.go)
  — `NewContentRuleService` returns `IContentRuleService` instead of
  `*ContentRuleService`.
- Five `testifylint` `require-error` findings (`assert.NoError` →
  `require.NoError`) and four `modernize` findings (`interface{}` → `any`) in the
  new tests. Lint is back to **0 issues**.
- One brand-new test file picked up a duplicated `package` clause from
  `golangci-lint-v2 --fix` (`expected declaration, found 'package'`); caught by
  `go run ./cmd/testlayoutcheck .` and fixed by hand.

## 3. Features added / changed

**No production behaviour changed.** The GPU-gated GUI snapshot suite is
byte-identical, which was Phase 6's acceptance criterion. Everything added is
test infrastructure:

- **14 `testify` mocks** in `test/test_helpers/`, one per interface the handler
  tests need (the plan explicitly rules out generating one mock per interface).
- **Five mirrored handler test packages** under `test/unit/internal/handlers/`,
  taking `internal/handlers` from 0% to **97.4%**.
- **Invariant-based catalogue tests** for `app/gui/constants` — `bannableItems.go`
  and `valueOverrideSids.go` are now at **100%**, and the Phase 10 sweep also
  closed `spells.go`, `bonusOptions.go`, `gameModes.go` and
  `internal/helpers.ScaleRound`.

## 4. File modifications

### Modified

| File | Summary |
| --- | --- |
| [internal/composition/providerSets.go](../internal/composition/providerSets.go) | dropped the last `wire.Bind` |
| `internal/composition/wire_gen.go` | regenerated (3 lines) — **never hand-edit** |
| [internal/services/content_rules/contentRuleService.go](../internal/services/content_rules/contentRuleService.go) | factory returns the interface |
| [plans/batch-11-handler-coverage.md](../plans/batch-11-handler-coverage.md) | Phases 6–11 marked Complete; Final Recap + Deployment Plan written |
| [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) | §6.2 and §6.4 marked `✅ FIXED` in place; §12 item 11 marked done |
| [todo/test_observations.md](../todo/test_observations.md) | 3 new unreachable-code entries |

### Created — mocks

All in `test/test_helpers/`, `package test_helpers`, receiver `this`, comma-ok
assertions on `arguments.Get(n)`:
`connectionEditorServiceMock.go`, `contentRuleMock.go`,
`editorStateValidatorMock.go`, `fileServiceMock.go`,
`generationTuningFactoryMock.go`, `generatorConfigMapperMock.go`,
`mandatoryContentProviderMock.go`, `manualReapplyServiceMock.go`,
`previewGeneratorServiceMock.go`, `previewLayoutServiceMock.go`,
`stateHandlerMock.go`, `templateGeneratorMock.go`, `zoneClassifierMock.go`,
`zoneEditorServiceMock.go`.

### Created — tests

| Folder | Files |
| --- | --- |
| `test/unit/internal/handlers/stateHandler/` | 5 (`common_test.go` + `newStateHandler`, `loadState` ×7, `saveState` ×6, `validateEditorState` ×9) |
| `test/unit/internal/handlers/previewHandler/` | 2 (`newPreviewHandler`, `buildPreviewLayout` ×6) |
| `test/unit/internal/handlers/templateHandler/` | 6 (`common_test.go` fixture + `generateTemplate` ×6, `updateTemplate` ×10, `reapplyCastleSettings` ×2, `saveTemplate` ×7) |
| `test/unit/internal/handlers/contentRuleHandler/` | 3 (`getContentRuleEditorOptions` ×5, `describeContentRule` ×8) |
| `test/unit/internal/handlers/zoneEditorHandler/` | 14 (`common_test.go` fixture + one file per public interface method) |
| `test/unit/app/gui/constants/bannableItems/` | 5 |
| `test/unit/app/gui/constants/valueOverrideSids/` | 1 |
| `test/unit/app/gui/constants/bonusOptions/` | 3 |
| `test/unit/app/gui/constants/gameModes/` | 1 |
| `test/unit/app/gui/constants/spells/` | 4 new, alongside the pre-existing `getSpellNameAndSchool_test.go` |
| `test/unit/internal/helpers/math/` | `scaleRound_test.go` (3) |

## 5. Tests added or updated

Roughly **130 new unit tests**, all following AGENTS.md §4.6: mirrored folder
per implementation file, one `<publicFuncName>_test.go` per public function,
package `<fileName>_test`, `Test{Scenario}_{ExpectedBehavior}` names, mandatory
`// Arrange` / `// Act` / `// Assert`, `t.Parallel()` in every test and every
`t.Run`, `testify` + `gofakeit` only.

**Last full run — all green:**

| Command | Result |
| --- | --- |
| `go build ./...` | pass |
| `go vet -tags=integration_test ./...` | pass |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test -count=1 ./test/unit/...` | pass (no FAIL) |
| `go test -tags=integration_test -count=1 ./test/integration/...` | `ok` (3.5s) |
| `go test -tags='integration_test,gui' -count=1 ./test/integration/gui/...` | `ok` (1.4s), **zero snapshot diffs** |
| coverage (`-coverpkg=./internal/...,./app/...`) | **68.7%** (Phase 0 baseline 65.5%) |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **0 issues** |

`coverage.txt`, `coverage.html` and `lcov.info` were regenerated.

## 6. Git status snapshot

Branch: **`AD/refactoring-07-21`**. **Nothing staged, nothing committed** — per
AGENTS.md §2.5 the owner reviews and commits.

```
 M .agent/session-carry-forward.md
 M todo/review-opus5-08-04.md
```

All of Batch 11 (the ~70-file interface refactor plus the ~130 new tests) was
reviewed and committed by the owner, which is why it no longer appears above.
The only uncommitted work is this handoff and the Batch 12 rejection recorded in
the review document.

## 7. Rejections / things the owner declined

Carried forward from earlier phases — **do not re-propose these**:

- A separate plan file for the review itself. Findings are marked
  `✅ FIXED` / `❌ WILL NOT FIX` **in place** in
  [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md).
- `IPositionedTopologyBuilder` — dropped, not deferred: the builder is embedded
  *by value* in six topology structs, which an interface cannot express.
- One provider per topology service returning `ITopologyService` — wire keys
  providers by output type, so 12 identical bindings are a generation error.
- A `CastleFactory()` accessor on `IZoneEditorService` (option B alone) and
  duplicating the road-rebuild helper inside `ManualReapplyService` (option C
  alone) — the B+C combination was chosen instead.
- Bulk `-tags` widening (`wireinject`, blanket `integration_test`).

Decided this session:

- **§1.8 / Batch 12 rejected outright** — see §1 above and §1.8 of the review
  document. Neither persistence shape (a) nor (b) is to be built; do not
  re-propose a `preferencesRepository`, an `os.UserConfigDir` file, or an
  `OutputDirectory` field on `EditorStateDto`.
- **Hard-coding the ~31 bannable-item SIDs** into assertions was rejected; the
  catalogue tests assert structural invariants plus a few named spot checks.
- **No test-only seams were added** to production code to reach
  `validators.ValidationIssue.fix`, the Steam/registry install-discovery chain,
  or `buildShiftDerangement` — AGENTS.md §4.6 forbids it. All three are recorded
  in [todo/test_observations.md](../todo/test_observations.md).
- **No production code was deleted**, including the two dead `zoneEditorHandler`
  methods (see §8) — out of scope for a coverage batch.

## 8. Open questions

1. **Dead code.** `zoneEditorHandler.ComputeHasErrors` and
   `zoneEditorHandler.RebuildZoneConnectionRoads` are exported on the private
   struct, absent from `handler_interfaces.IZoneEditorHandler` (which is what
   `NewZoneEditorHandler` returns), and called by nobody — therefore untestable.
   Delete them, or add them to the interface?
2. **`internal/helpers/io.go`.** Should the Steam/registry install-discovery
   chain (`getVDFContent`, `getVDFFilePath`, `getSteamPath`, `getBasePath`,
   `getSteamPathFromRegistry`) get an injectable filesystem seam so it can be
   covered? That is a production API change, not a test change.

## 9. Next recommended actions

1. Answer the two open questions in §8.
2. **Item 13 — the large refactors — is all that remains in §12.** Plan first
   per AGENTS.md §4.7: §2.1 (extract filesystem policy out of
   `fileExplorerDialog.go`) → unblocks §2.5 (the same file is a god object).
   Then §2.2 (extract regeneration policy from the GUI driver), which is
   multi-session, overlaps a backlog item and **still needs an owner decision on
   refactor scope**. §2.6 (`zoneEditorDialog`'s ~58 fields) opportunistically,
   whenever the zone editor is next touched.
3. Batches 1–12 are closed. §2.2's scope is the only outstanding owner decision
   in the whole review.

## 10. Carry-forward prompt

> Read `AGENTS.md` first and follow it strictly.
>
> Hard rules, one line each: never modify `data/`, `internal/entities/template/`
> or `internal/registry/`; keep everything cross-platform (Windows + Linux,
> `path/filepath`, PowerShell chains with `;`, never `&&`); every change ships
> with tests and must not drop coverage (baseline **68.7%**, lint baseline
> **0 issues**); durable multi-session work gets a plan file under `plans/`;
> **never stage and never commit** — the owner reviews and commits.
>
> We are remediating the 46-finding review in `todo/review-opus5-08-04.md` (§12
> defines the 13 PR-sized batches). Findings are marked `✅ FIXED` /
> `❌ WILL NOT FIX` **in place** in that review document — do not create a
> separate plan file for the review itself.
>
> **Where work left off: Batches 1–12 are all closed.** Batch 11 (findings §6.2
> and §6.4) is `✅ FIXED` and committed: the whole `internal/` DI graph is
> interface-based with **zero** `wire.Bind` calls remaining, `internal/handlers`
> is at 97.4%, and total unit coverage rose 65.5% → 68.7%. Batch 12 (§1.8) is
> `❌ WILL NOT FIX` — the output directory is a per-machine requirement of the
> game, not a preference; read §1.8 of the review before anyone suggests
> persisting it again. Build, both `go vet` tag combinations, `testlayoutcheck`,
> the unit / integration / GPU-gated GUI suites and `golangci-lint-v2` are all
> green.
>
> **Next up is §12 item 13 — the large refactors**, which need a plan file per
> AGENTS.md §4.7: §2.1 (extract filesystem policy out of `fileExplorerDialog.go`,
> unblocking §2.5), then §2.2 (extract regeneration policy from the GUI driver —
> multi-session, and the owner has *not* yet decided its scope), and §2.6
> opportunistically. §2.2's scope is the only outstanding owner decision in the
> entire review.
>
> Before starting new work, ask the owner the two open questions in §8 of
> `.agent/session-carry-forward.md` (the two dead `zoneEditorHandler` methods,
> and whether `internal/helpers/io.go` should get a filesystem seam).
>
> Workflow for every batch, without exception: (1) ask the owner whether the
> batch should be done at all; (2) if declined, document in the review file why
> it should not be attempted in future; (3) ask all clarifying questions up
> front; (4) implement; (5) rewrite `.agent/session-carry-forward.md`; (6) stop
> and wait for owner review.
>
> Environment gotchas: `wire gen` writes its success banner to STDERR so
> PowerShell reports exit 1 even on success — use `wire diff` (exit 0 = current);
> `golangci-lint-v2` likewise exits 1 on stderr warnings even when it reports
> `0 issues`; `golangci-lint-v2 --fix` has duplicated the `package` clause on
> freshly created files, so always re-run `go run ./cmd/testlayoutcheck .`
> afterwards.
>
> See `.agent/session-carry-forward.md` and
> `plans/batch-11-handler-coverage.md` for the full handoff.
