# Session Carry-Forward

## Session goal

Complete the two remaining Phase 3 items in `plans/clean-architecture-refactoring.md`: injectable `GUIHandler` construction and focused internal workflow decomposition behind the unchanged GUI-facing facade.

## Fixes applied

- Converted editor-state validation to constructor-owned `validators.EditorStateValidator` and injected it into `GUIHandler`.
- Replaced content-rule package functions and mutable catalogs with `ContentRuleService`, `DistanceCatalog`, and immutable `VariantMappingCatalog`.
- Added an empty-variant guard to `NewRuleVariant` so it returns an error instead of panicking.
- Replaced connection-editor package functions with `ConnectionEditorService`, `ZoneEditorService`, and `ManualReapplyService`.
- Centralized topology-backed castle/road operations behind private `ZoneEditorService` methods; manual reapply depends on that service.
- Moved editor-row-to-mandatory-item conversion from `MandatoryContentProvider` to `mappers.MandatoryContentItemMapper`.
- Split `GUIHandler` into focused template-workflow, state-persistence, template-persistence, preview, content-rule, and zone-editor handlers while retaining every facade method.
- Added validated `GUIHandlerDependencies` injection plus production default wiring in `NewGuiHandler`.
- Made connection-editor dependency constructors replace nil mandatory collaborators with usable defaults.

## Features added / changed

- Variant catalog queries return fresh slices, including cloned nested variant tuples.
- `GUIHandler` shares one zone classifier, zone editor, and tuning factory across connection/manual services.
- Content-rule and connection-editor public APIs now expose constructors and receiver methods only.
- Per user decision, `GeneratorConfigMapper` now preserves `EditorStateDto.CityHold` exactly instead of deriving it from `VictoryCondition`; downstream `IsCityHoldMode` and game-rule generation still interpret the City Hold victory condition.
- Per user decision, preview-generator initialization remains optional: startup logs failure and preview-image saving is skipped.
- Tests can inject workflow fakes through `NewGuiHandlerWithDependencies`; production still shares mapper, validator, classifier, zone-editor, tuning, connection-editor, and file-service instances where policy must remain identical.

## File modifications

Created or renamed production owners:

- `internal/mappers/mandatoryContentItemMapper.go` — row normalization, multiplicity, SID/group/mine mapping, and content-rule application.
- `internal/services/connection_editor/connectionEditorService.go` — connection defaults and graph diagnostics.
- `internal/services/connection_editor/zoneEditorService.go` — zone mutation, castle counting, naming, positions, and road rebuilding.
- `internal/services/connection_editor/manualReapplyService.go` — castle-setting reapply policy.
- `internal/services/content_rules/contentRuleService.go`, `distanceCatalog.go`, `distanceVariation.go`, and `variantMappingCatalog.go` — content-rule ownership and catalogs.
- `internal/validators/editorStateValidator.go`, `intField.go`, and `rangedIntField.go` — validator service ownership.

Major edited/deleted files:

- `internal/handlers/guiHandler.go` — injected validator/content/connection/zone/manual services.
- `internal/handlers/guiHandlerDependencies.go` and `*Operations.go` — explicit injectable workflow contract.
- `internal/handlers/{templateWorkflow,statePersistence,templatePersistence,preview,contentRule,zoneEditor,stateValidation}Handler.go` — focused private workflow owners.
- `internal/mappers/generatorConfigMapper.go` — uses `MandatoryContentItemMapper`; copies `CityHold` directly.
- `internal/services/template_generator/providers/mandatoryContentProvider.go` — removed row mapping and content-rule dependency.
- Old `connectionEditor.go`, `zoneEditor.go`, `manualReapply.go`, content-rule manager/catalog files, and their old test ownership folders were deleted/replaced.
- `plans/clean-architecture-refactoring.md` — Phase 3 is In progress with only two unchecked handler items.

Tests were moved to implementation-matching folders under:

- `test/unit/internal/services/content_rules/`
- `test/unit/internal/services/connection_editor/{connectionEditorService,zoneEditorService,manualReapplyService}/`
- `test/unit/internal/mappers/mandatoryContentItemMapper/`

Handler and `test/integration/manualCastleReapply_integration_test.go` fixtures now call service receiver APIs.

New handler constructor/delegation tests live in `test/unit/internal/handlers/guiHandler/newGuiHandlerWithDependencies_test.go`; nil-default constructor coverage was added to the matching connection-editor service test folders.

## Tests added or updated

Latest verified state after Phase 3 completion:

- Focused content-rule, connection-editor, mapper, handler, and provider suites: passed.
- `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...`: passed at 64.6% (Phase 3 checkpoint 64.3%; Phase 2 baseline 64.2%).
- Complete unit suite: 2,104 tests passed.
- `go build ./...`: passed.
- `go test ./test/... -count=1`: passed.
- `go test -tags=integration_test ./test/integration/... -count=1`: passed.
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`: passed.
- `go test -tags=integration_test ./test/performance/... -count=1`: passed.
- Architecture dependency guard, diagnostics, handler lint, and `git diff --check`: passed.
- Broader package lint reports only eight pre-existing formatting findings in untouched tests.
- Claude Fable 5 approved final Phase 3 completion with no hard findings.

## Git status snapshot

- Branch: `AD/refactoring-07-21`, up to date with `origin/AD/refactoring-07-21` at the latest check.
- Current Phase 3 handler/dependency changes and plan/handoff updates are unstaged; new focused handler/interface and test files are untracked.
- No files were staged, unstaged, committed, or reverted by the agent.
- `coverage.txt` was regenerated by required coverage runs.
- No changes were made under `data/`, `internal/entities/template/`, or `internal/registry/`.
- Run `git status --short` for the exact inherited listing and preserve its state.

## Rejections / things the user declined

- The user chose to preserve the persisted `CityHold` flag only, declining mapper normalization from `VictoryCondition`.
- The user chose optional preview initialization, declining a startup-failing handler constructor.
- One Opus review attempt hit the prior credit limit; the required reviews were completed successfully with Claude Fable 5.
- Two patches were rejected on stale context and applied nothing; each was reapplied in smaller validated edits.
- One combined patch was rejected for duplicate file sections and reapplied in valid smaller edits; one partial-file lint invocation failed because same-package siblings were omitted, then package lint completed.

## Open questions

- Phase 3 has no blockers and is marked Complete.
- A caller can still pass a typed-nil implementation through a `GUIHandlerDependencies` interface field; only tests use this constructor today, and Fable classified reflection-based guarding as unnecessary unless external use appears.
- `ContentRuleService.GetRules` ignores the impossible built-in `NewRuleVariant` error; Fable classified this as non-blocking because the private default catalog is hardcoded non-empty.
- `HasDuplicateName` remains test-only production API and is deferred to Phase 8 cleanup.

## Next recommended actions

1. Read `AGENTS.md`, Phase 4 in the plan, and this handoff.
2. Recheck `git status --short` and `git diff --check` without altering staging.
3. Preserve the completed Phase 3 handler facade and constructor contracts while beginning Phase 4 catalog/common-type consolidation.
4. Record a fresh coverage comparison point before Phase 4 edits.
5. Follow the Phase 4 distance-context constraint: preserve content `Near = 0.1..0.25` and portal `Near = 0.075..0.35` unless the user explicitly approves a product change.

## Carry-forward prompt

Read `AGENTS.md` first, then read `plans/clean-architecture-refactoring.md` and `.agent/session-carry-forward.md`. Never modify `data/`, `internal/entities/template/`, or `internal/registry/`; preserve Windows/Linux portability; add tests for non-trivial logic and do not reduce coverage. Phase 3 is Complete and Fable-approved at 64.6% coverage: `GUIHandler` is an injectable thin facade over six focused workflow handlers, preview initialization remains optional, and mapper `CityHold` preserves only the persisted flag. Resume from the current unstaged worktree at Phase 4 only after checking status and recording its baseline. Use `.agent/session-carry-forward.md` for the full handoff.