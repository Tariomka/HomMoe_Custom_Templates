# Session Carry-Forward

## 1. Session goal

Remediate **Batch 11** of the 46-finding review in
[todo/review-opus5-08-04.md](todo/review-opus5-08-04.md) — review findings
**§6.2** (`internal/handlers` has no mirrored unit tests) and **§6.4** (two
`app/gui/constants` catalogues at 0% coverage). Because the five handlers held
**concrete** struct dependencies that cannot be mocked, the owner expanded the
scope: first convert **every constructor-injected service under `internal/`** to
an interface (full closure), then write the handler tests against `testify`
mocks.

Work is tracked in [plans/batch-11-handler-coverage.md](plans/batch-11-handler-coverage.md),
which is the source of truth. **Phases 0–5 are complete; Phases 6–11 remain.**

## 2. Fixes applied

- Two `wastedassign` lint regressions inherited from Batch 10 fixed in
  [app/gui/panels/bonusesPanel.go](app/gui/panels/bonusesPanel.go)
  (`label := ""` → `var label string` in the `BonusSpell` and
  `BonusStartingItem` branches of `bonusDisplayName`). Lint is back to **0 issues**.
- Regenerated the stale `coverage.txt` (baseline total coverage: **65.5%**).
- **Typed-nil hazard eliminated:** `providePreviewGenerator` used to return a nil
  `*PreviewGeneratorService`, which would have become a non-nil interface holding
  a nil pointer once the field was interface-typed. It now returns a null object,
  and the `if this.previewGenerator != nil` guard is gone from
  [internal/handlers/templateHandler.go](internal/handlers/templateHandler.go).
- **Private-access leak removed:** `ManualReapplyService` no longer reaches into
  `ZoneEditorService.castleFactory` / `.rebuildCastleRoads`.

## 3. Features added / changed

All changes are DI refactors that unblock mock-based handler testing.

- **AGENTS.md §4.2.2 "Interface placement"** (authored by the owner, added
  between §4.2.1 and §4.3): <5 implementation files needing interfaces → same
  package; ≥5 → `{singular package name}_interfaces` subpackage; spanning
  packages / cycle-breaking → `internal/interfaces/`. Factories return the
  **interface**, and the **broadest** one when an implementation satisfies several.
- **Phase 1 — 6 leaf interfaces:** `IFileService`, `IEditorStateValidator`,
  `IGeneratorConfigMapper`, `IMandatoryContentItemMapper`,
  `IGenerationTuningFactory`, `ITopologyConnectionService`.
- **Phase 2 — `internal/services/zones/zone_interfaces/`** (5 interfaces):
  `ICastleFactory`, `IRoadFactory`, `IZoneFactory`, `IZoneClassifier`,
  `IZoneLabelProvider`. The old `zones/zoneLabelProviderInterface.go` was deleted
  and its `wire.Bind` dropped. Owner decisions: `createPlayerSpawnCastle` /
  `createAbandonedOutposts` promoted to public; `ZoneEditorService.CastleFactory`
  renamed to unexported `castleFactory`.
- **Phase 3 — `topology/topology_interfaces/ITopologyService`** plus
  `topology/topologyServiceAssertions.go` with 12 compile-time assertions. The 12
  concrete constructors are **unchanged** because wire keys providers by output
  type and 12 providers returning one interface is a multiple-bindings error.
  `IPositionedTopologyBuilder` was **dropped** — the builder is embedded *by
  value*, which an interface cannot express.
- **Phase 4 — `providers/provider_interfaces/`** (7 interfaces) +
  `template_generator.ITemplateGenerator`. The `TopologyVariantCreator` func type
  was relocated from `providers` into `provider_interfaces` to break an import
  cycle, and the original file deleted.
- **Phase 5 — `connection_editor` + `preview_service` interfaces and the null
  object.** Owner decision (B+C combination) for the `ManualReapplyService`
  blocker: it now takes its own `ICastleFactory`, and
  `ZoneEditorService.rebuildCastleRoads` was promoted to public
  `RebuildCastleRoads` and put on `IZoneEditorService` — so no `CastleFactory()`
  accessor and **no `wire.Bind`** were needed. `NewManualReapplyService`'s
  parameters are now `(zoneEditor, castleFactory, zoneClassifier, tuningFactory)`.
  All five handler structs are now interface-typed.

Net effect: exactly **one** `wire.Bind` remains in the codebase
(`IContentRuleService` ← `*ContentRuleService`).

## 4. File modifications

### Created (untracked)

| File | Summary |
| --- | --- |
| [plans/batch-11-handler-coverage.md](plans/batch-11-handler-coverage.md) | The 12-phase plan; **read this first**. |
| [internal/services/file_service/fileServiceInterface.go](internal/services/file_service/fileServiceInterface.go) | `IFileService` (3 methods). |
| [internal/validators/editorStateValidatorInterface.go](internal/validators/editorStateValidatorInterface.go) | `IEditorStateValidator`. |
| [internal/mappers/generatorConfigMapperInterface.go](internal/mappers/generatorConfigMapperInterface.go) | `IGeneratorConfigMapper`. |
| [internal/mappers/mandatoryContentItemMapperInterface.go](internal/mappers/mandatoryContentItemMapperInterface.go) | `IMandatoryContentItemMapper`. |
| [internal/services/template_generator/generation_tuning/generationTuningFactoryInterface.go](internal/services/template_generator/generation_tuning/generationTuningFactoryInterface.go) | `IGenerationTuningFactory`. |
| [internal/services/template_generator/providers/topology/base/topologyConnectionServiceInterface.go](internal/services/template_generator/providers/topology/base/topologyConnectionServiceInterface.go) | `ITopologyConnectionService` (4 methods). |
| `internal/services/zones/zone_interfaces/` (5 files) | `ICastleFactory`, `IRoadFactory`, `IZoneFactory`, `IZoneClassifier`, `IZoneLabelProvider`. |
| `internal/services/template_generator/providers/topology/topology_interfaces/topologyServiceInterface.go` | `ITopologyService`. |
| [internal/services/template_generator/providers/topology/topologyServiceAssertions.go](internal/services/template_generator/providers/topology/topologyServiceAssertions.go) | 12 compile-time `ITopologyService` assertions. |
| `internal/services/template_generator/providers/provider_interfaces/` (8 files) | 7 provider interfaces + the relocated `TopologyVariantCreator`. |
| [internal/services/template_generator/templateGeneratorInterface.go](internal/services/template_generator/templateGeneratorInterface.go) | `ITemplateGenerator`. |
| [internal/services/connection_editor/connectionEditorServiceInterface.go](internal/services/connection_editor/connectionEditorServiceInterface.go) | `IConnectionEditorService` (4 methods). |
| [internal/services/connection_editor/zoneEditorServiceInterface.go](internal/services/connection_editor/zoneEditorServiceInterface.go) | `IZoneEditorService` (10 methods). |
| [internal/services/connection_editor/manualReapplyServiceInterface.go](internal/services/connection_editor/manualReapplyServiceInterface.go) | `IManualReapplyService` (2 methods). |
| [internal/services/preview_service/previewLayoutServiceInterface.go](internal/services/preview_service/previewLayoutServiceInterface.go) | `IPreviewLayoutService`. |
| [internal/services/preview_service/previewGeneratorServiceInterface.go](internal/services/preview_service/previewGeneratorServiceInterface.go) | `IPreviewGeneratorService`. |
| [internal/services/preview_service/nullPreviewGeneratorService.go](internal/services/preview_service/nullPreviewGeneratorService.go) | Null-object generator; `CreatePreviewImage` returns `nil`. |
| [test/unit/internal/services/connection_editor/zoneEditorService/rebuildCastleRoads_test.go](test/unit/internal/services/connection_editor/zoneEditorService/rebuildCastleRoads_test.go) | 3 tests for the newly public method. |
| `test/unit/internal/services/preview_service/nullPreviewGeneratorService/` (2 files) | 3 tests for the null object. |

### Deleted

- `internal/services/zones/zoneLabelProviderInterface.go` — relocated into
  `zone_interfaces/`.
- `internal/services/template_generator/providers/topologyVariantCreator.go` —
  relocated into `provider_interfaces/`.

### Edited (highlights; see the `git status` snapshot in §6 for the full list)

- [AGENTS.md](AGENTS.md) — new §4.2.2.
- [app/gui/panels/bonusesPanel.go](app/gui/panels/bonusesPanel.go) — lint fix.
- [internal/composition/previewGeneratorProvider.go](internal/composition/previewGeneratorProvider.go) — null object on asset failure.
- [internal/composition/topologyServiceProvider.go](internal/composition/topologyServiceProvider.go) — returns `ITopologyServiceLookup`.
- [internal/composition/providerSets.go](internal/composition/providerSets.go) — `wire.Bind` for `IZoneLabelProvider` removed.
- `internal/composition/wire_gen.go` — regenerated (never hand-edit).
- All five handlers: [templateHandler.go](internal/handlers/templateHandler.go),
  [zoneEditorHandler.go](internal/handlers/zoneEditorHandler.go),
  [previewHandler.go](internal/handlers/previewHandler.go),
  [stateHandler.go](internal/handlers/stateHandler.go), plus `guiHandler`'s graph.
- ~45 further service/provider/topology files switched to interface parameters
  and interface-returning factories.
- 8 `test/test_helpers/*.go` builders and 7 `common_test.go` helpers updated for
  the new return types and the new `NewManualReapplyService` signature.

## 5. Tests added or updated

**Added (6 tests in 3 files):**
- `rebuildCastleRoads_test.go` — rebuilds castle roads for the current main
  objects; drops stale castle roads; preserves connection roads.
- `nullPreviewGeneratorService/createPreviewImage_test.go` — returns no image for
  a real template and for a nil template.
- `nullPreviewGeneratorService/newNullPreviewGenerator_test.go` — constructs.

**Updated:** the `common_test.go` helpers in `guiHandler`,
`manualReapplyService`, `fileService`, `gladiatorArenaProvider`,
`mandatoryContentProvider`, `zoneFactory`, and
`previewGeneratorService/createPreviewImage_test.go` — all mechanical
return-type / signature changes, no behavioural edits.

**Last run status — all green:**

| Command | Result |
| --- | --- |
| `go build ./...` | pass |
| `go vet -tags=integration_test ./...` | pass |
| `go vet -tags='integration_test,gui' ./test/...` | pass |
| `go test -count=1 ./test/unit/...` | pass (no FAIL) |
| `go test -tags=integration_test -count=1 ./test/integration/...` | `ok` (3.1s) |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `golangci-lint-v2 run ./...` | **0 issues** |

Not re-run this session: the GPU-gated UI suite
(`go test -tags='integration_test,gui' ./test/integration/gui/...`) — snapshot
verification is Phase 6's job. Coverage has not been re-measured since the
**65.5%** baseline; no production behaviour changed, so it is expected to hold.

## 6. Git status snapshot

Branch: **`AD/refactoring-07-21`**. **Nothing is staged** — the working tree
holds 68 modified, 2 deleted and ~20 untracked paths, all of them this session's
work. Per AGENTS.md §2.5 nothing was staged or committed; the owner reviews and
commits.

Untracked directories the next session inherits: `plans/`,
`internal/services/zones/zone_interfaces/`,
`internal/services/template_generator/providers/provider_interfaces/`,
`internal/services/template_generator/providers/topology/topology_interfaces/`,
`test/unit/internal/services/preview_service/nullPreviewGeneratorService/`.

Re-run `git status --short` at the start of the next session for the exact list.

## 7. Rejections / things the owner declined

- **A separate plan file for the review itself** — declined at the outset;
  findings are marked `✅ FIXED` / `❌ WILL NOT FIX` **in place** in
  [todo/review-opus5-08-04.md](todo/review-opus5-08-04.md).
- **`IPositionedTopologyBuilder`** — dropped, not deferred:
  `PositionedTopologyBuilder` is embedded *by value* in six topology structs, so
  an interface cannot replace it.
- **One provider per topology service returning `ITopologyService`** — rejected;
  wire keys providers by output type, so 12 identical bindings are a generation
  error. The owner chose the interface + concrete factories + assertions file.
- **Adding a `CastleFactory()` accessor to `IZoneEditorService`** (option B
  alone) and **duplicating the road-rebuild helper inside
  `ManualReapplyService`** (option C alone) — both rejected in favour of the B+C
  combination described in §3.
- **Bulk `-tags` widening** (`wireinject`, blanket `integration_test`) — never
  attempted; forbidden by AGENTS.md §4.6.1 / §4.6.3.

## 8. Open questions

None blocking. Two decisions the next session will need to make (or confirm with
the owner) inside Phase 7:

1. **Mock breadth.** The plan says to generate **only** the mocks the handler
   tests actually need, not one per interface. Confirm before writing 20+ mocks.
2. **Phase 9 spot checks.** The plan explicitly forbids hard-coding the ~31 SIDs;
   invariants plus 2–3 named spot checks only.

## 9. Next recommended actions

1. **Phase 6 — wire.** Review [internal/composition/providerSets.go](internal/composition/providerSets.go)
   against the new interface-returning factories, run
   `wire gen ./internal/composition/...`, confirm a `wire diff` is clean, then run
   the GPU-gated UI snapshot suite and confirm **zero** snapshot diffs.
2. **Phase 7 — mocks.** Add `test/test_helpers/<interfaceName>Mock.go` files,
   copying the style of
   [test/test_helpers/contentRuleServiceMock.go](test/test_helpers/contentRuleServiceMock.go)
   (`type XMock struct { mock.Mock }`, receiver `this`,
   `arguments := this.Called(...)`, `value, _ := arguments.Get(0).(T)`).
3. **Phase 8 — §6.2.** Mirrored test folders for all five handlers. Cover
   explicitly: empty/whitespace path → `ErrNoOutputPath`; nil state →
   `ErrNothingToSave`; `fixIssues=false` returns warnings but an unmodified
   state; `AdvancedMode` true/false each zeroing the correct field set;
   `BuildPreviewLayout` with nil template + nil zones + nil connections.
4. **Phase 9 — §6.4.** `test/unit/app/gui/constants/bannableItems/` (5 files:
   `getBannableItemsWithExclusions`, `findBannableItem`, `getBannedItemLabel`,
   `sidToDisplayName`, `compareBannableItems`) and
   `test/unit/app/gui/constants/valueOverrideSids/getValueOverrideSidsWithExclusions_test.go`.
5. **Phase 10.** Regenerate the coverage profile, find the remaining 0% non-Gio
   files, test the easy/pure ones, record the untestable ones in
   [todo/test_observations.md](todo/test_observations.md).
6. **Phase 11.** Full suite; mark §6.2 and §6.4 `✅ FIXED` **in place** in
   [todo/review-opus5-08-04.md](todo/review-opus5-08-04.md) and tick §12 item 11;
   update repository memory; rewrite this file; stop for owner review.

## 10. Carry-forward prompt

> Read `AGENTS.md` first. Hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/`; keep everything
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chains with `;`,
> never `&&`); every change ships with tests and must not drop coverage;
> durable multi-session work gets a plan file under `plans/`; **never stage and
> never commit** — the owner reviews and commits.
>
> We are remediating **Batch 11** of the 46-finding review in
> `todo/review-opus5-08-04.md` (§12 defines the 13 PR-sized batches). Findings
> are marked `✅ FIXED` / `❌ WILL NOT FIX` **in place** in that review document —
> do not create a separate plan file for the review itself.
>
> Batch 11 closes review findings §6.2 (no mirrored unit tests for
> `internal/handlers`) and §6.4 (two `app/gui/constants` catalogues at 0%). The
> owner expanded it: every constructor-injected service under `internal/` is
> being converted to an interface first, so the handlers can be tested against
> `testify` mocks.
>
> **Where work left off:** Phases 0–5 of `plans/batch-11-handler-coverage.md`
> are **Complete** with full Phase Summaries written. The whole `internal/` DI
> graph is now interface-based (only one `wire.Bind` remains,
> `IContentRuleService`). `go build`, both `go vet` tag combinations, the unit
> suite, the integration suite, `go run ./cmd/testlayoutcheck .` and
> `golangci-lint-v2` are all green. **Next up is Phase 6 (wire regeneration +
> GUI snapshot verification), then Phases 7–11.**
>
> Read `plans/batch-11-handler-coverage.md` — it is the source of truth — and
> `.agent/session-carry-forward.md` for the full handoff, including the owner
> decisions, the rejected options (do not re-propose them) and the environment
> gotchas (`wire gen` writes its success line to STDERR; freshly created files
> have once picked up a stray duplicate `package` clause).
>
> Workflow for every batch, without exception: (1) ask the owner whether the
> batch should be done at all; (2) if declined, document in the review file why
> it should not be attempted in future; (3) ask all clarifying questions up
> front; (4) implement; (5) rewrite `.agent/session-carry-forward.md`; (6) stop
> and wait for owner review.
