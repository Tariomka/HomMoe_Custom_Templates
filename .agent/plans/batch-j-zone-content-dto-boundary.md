# Batch J: zone-content service DTO boundary (review §2.2)

Remove the whole `internal/services/zone_content` DTO exception. The service speaks
models only, the handlers own every DTO ⇄ Model conversion, and the
`dtoNamerAllowList` entry is deleted. The `internal/services/bonuses` exception stays.
There is no behaviour change.

## For Future Agents
As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../../AGENTS.md), [the handoff](../session-carry-forward.md) and review
[§2.2](../backlog/review-gpt-6-astra-09-07.md) first. Branch at planning time:
`AD/performance`, HEAD `33bff1f`, clean tree. Batch G (`35e0fab`) is not yet merged to
`master` (`f78caeb`). The owner stages and commits everything; never stage, commit or
push.

## 0. Settled owner decisions (2026-09-28)

| ID | Decision |
| --- | --- |
| J1 | Scope is the **whole** zone-content exception: all four DTO-bearing service methods move to models, and `internal/services/zone_content` is removed from `dtoNamerAllowList`. The bonuses entry is untouched. |
| J2 | The composition unit is `editor_state_model.ContentRuleRow` (not `ZoneContentRow`). The service returns `(ContentRuleRow, bool)` comma-ok. |
| J3 | `ContentRuleKey` (5 consts) and `ContentRuleEditorKind` (3 consts) move out of `internal/dtos` into the new model package, because §4.4 lets models import only entities, helpers and registry. Their `dtos` files are deleted with no alias left behind. |
| J4 | New package `internal/models/content_rule_model/`, one type per file (§4.1). |
| J5 | `ContentRuleComposition` is **flat**: `Key, Name, DistanceNames, DistanceIndex, IsGuarded, IsSoloEncounter, VariantIDs, VariantIndex`. `ContentRuleCompositionRequestDto` becomes `struct{ content_rule_model.ContentRuleComposition }`. |
| J6 | New `ContentRuleDescription{Key, DisplayText, Marker, VariantLabel, Valid, SavedRule}`. `ContentRuleDescriptionDto` becomes `struct{ content_rule_model.ContentRuleDescription }`. |
| J7 | New `ContentRuleOption{Key, Name}`. `ContentRuleOptionDto` embeds it and keeps `Description, Marker, EditorKind, EditorLabel`. `GetDefaultContentRules` takes `[]ContentRuleOption`, so today's "only when Guarded is offered" check is kept. |
| J8 | `ContentRuleCompositionResultDto{Rule, Valid}` stays at the handler boundary. The handler interfaces (`IZoneContentHandler`, `IContentRuleHandler`, `IGuiHandler`) keep their signatures. `ContentRuleEditorOptionsDto` and `ContentRuleVariantOptionDto` are unchanged. |
| J9 | The option catalogue and describe logic in `contentRuleHandler` stay where they are. That this is business logic in a handler is recorded as an observation only. |

**Plan review (Claude Opus 5.5, 2026-09-28):** APPROVE WITH CHANGES; all 10 findings
applied (missed `TemplateHandlerMock.ComposeContentRule`, keyed embedded literals, `ok`
assertions, handler test names, reproducible mutation check, tagged `go vet`, git-status
and grep scopes, §8 hash baseline). No cycles, selector ambiguity, Wire or
testlayoutcheck issues found. **Owner plan approval: 2026-09-28**, "Approved, proceed, finish all of the phases, then I
will review". Phase 4 stops at the owner review; close-out waits for the owner's commit.

**Out of scope:** behaviour changes, GUI goldens, Wire (no constructor signature
changes), `data/`, `internal/entities/template_entity/`, `internal/registry/`, the
output path, the bonuses service, dependencies, and `ContentRuleVariantOptionDto`.

## 1. Inventory at planning time

Service methods that name DTOs ([zoneContentEditorService.go](../../internal/services/zone_content/zoneContentEditorService.go),
[its interface](../../internal/services/zone_content/zoneContentEditorServiceInterface.go)):

| Method today | After |
| --- | --- |
| `ComposeContentRule(dtos.ContentRuleCompositionRequestDto) dtos.ContentRuleCompositionResultDto` | `ComposeContentRule(content_rule_model.ContentRuleComposition) (editor_state_model.ContentRuleRow, bool)` |
| `GetDefaultContentRules(dtos.ContentRuleEditorOptionsDto) []ContentRuleRow` | `GetDefaultContentRules([]content_rule_model.ContentRuleOption) []ContentRuleRow` |
| `GetContentRuleMarkers([]dtos.ContentRuleDescriptionDto) string` | `GetContentRuleMarkers([]content_rule_model.ContentRuleDescription) string` |
| `GetContentRowDisplayName(string, []dtos.ContentRuleDescriptionDto) string` | `GetContentRowDisplayName(string, []content_rule_model.ContentRuleDescription) string` |

Other production references to the moved/reshaped types: `internal/handlers/contentRuleHandler.go`,
`zoneContentHandler.go`, `guiHandler.go`, `handler_interfaces/contentRuleHandlerInterface.go`,
`handler_interfaces/zoneContentHandlerInterface.go`, `app/gui/dialogs/ruleDialog.go`.

Test references: `test/test_helpers/templateHandlerMock.go`,
`test/test_helpers/zoneContentEditorServiceMock.go`, and the unit folders
`handlers/contentRuleHandler`, `handlers/guiHandler`, `handlers/zoneContentHandler`,
`services/zone_content/zoneContentEditorService`. Only GUI flow:
`test/integration/gui/contentRuleDialogs_integration_test.go`. Re-grep
`ContentRuleKey|ContentRuleEditorKind|ContentRuleOptionDto|ContentRuleComposition|ContentRuleDescription`
before editing; line numbers drift.

## Phase 0: Baseline
Status: Complete

- [x] Confirm branch/HEAD with `git status --short`; the only expected entry is `?? .agent/plans/` (this plan).
- [x] Record the review's LF-normalized UTF-8 SHA-256 of §8 (from `## §8` up to `## §9`) for the Phase 4 check.
- [x] `go build ./...`: PASS.
- [x] Coverage: `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...` then `go tool cover '-func=coverage.txt'`. Record the total (expected 74.5%) and the per-function lines for `zoneContentEditorService.go`, `zoneContentHandler.go`, `contentRuleHandler.go`.
- [x] `go test ./test/... -count=1`: PASS.
- [x] `go test -tags='integration_test,gui' ./test/integration/... -count=1`: PASS.
- [x] `go run ./cmd/testlayoutcheck .`: PASS; `gofmt -l` clean on `app internal test cmd`.
- [x] `golangci-lint-v2 run ./... --issues-exit-code=0`: 0 issues.

### Verification Plan
- All commands above succeed; numbers recorded in the Phase Summary.

### Phase Summary
Windows/amd64, Go 1.27.0, empty `GOFLAGS`, HEAD `33bff1f`. `git status --short` showed
` M .agent/plans/batch-j-zone-content-dto-boundary.md` (the owner had staged the plan;
left as is). Review §8 LF hash: `3849AE8FF1B3A826744EF2C5AF6DACA97FD0028CD8A229AB827AA9937DC88425`.
Build PASS. Coverage **74.5%**, covered 6909 / 9261 statements (unique blocks, max count).
Every function in `zoneContentEditorService.go` and `zoneContentHandler.go` is 100%;
`contentRuleHandler.go` is 100% except `contentRuleKeyFromName` at 71.4%.
`./test/...` PASS; tagged run PASS (root 3.429s, GUI 26.150s); testlayoutcheck PASS; gofmt
clean; lint 0 issues.

## Phase 1: Model package and DTO reshape (no service signature change yet)
Status: Complete

- [x] Create `internal/models/content_rule_model/`:
  - `contentRuleKey.go`: `type ContentRuleKey string` + the 5 constants, names and values identical to today's `dtos` ones.
  - `contentRuleEditorKind.go`: `type ContentRuleEditorKind string` + the 3 constants, identical names and values.
  - `contentRuleOption.go`: `ContentRuleOption{Key ContentRuleKey; Name string}`.
  - `contentRuleComposition.go`: flat `ContentRuleComposition` (J5).
  - `contentRuleDescription.go`: `ContentRuleDescription` (J6), `SavedRule editor_state_model.ContentRuleRow`.
- [x] Delete `internal/dtos/contentRuleKey.go` and `internal/dtos/contentRuleEditorKind.go`.
- [x] Reshape `ContentRuleOptionDto` (embed + 4 fields), `ContentRuleCompositionRequestDto` (embed only), `ContentRuleDescriptionDto` (embed only). Leave the result, editor-options and variant-option DTOs as they are.
- [x] Update every production reference to compile: `contentRuleHandler.go` (catalogue literals, `contentRuleKeyFromName`, describe), `ruleDialog.go` (the `switch` on `EditorKind`, the Guarded key check, the request literal), and the service, which switches every `dtos.ContentRuleKey*` to `content_rule_model.*` and every `request.Option.*` to `request.*` (in this phase it imports both `dtos` and `content_rule_model`).
- [x] ~~Composite literals cannot set promoted fields~~ **Superseded:** Go 1.27 accepts promoted fields in composite literals, and lint's `modernize/embedlit` rejects the keyed embedded form. Final literals are flat (`dtos.ContentRuleOptionDto{Key: ..., Name: ..., Description: ...}`); only literals that wrap an existing model value keep the embedded key (`dtos.ContentRuleCompositionRequestDto{ContentRuleComposition: composition}`). The embedded-field-first rule in declarations stands. Original text: every literal of the three reshaped DTOs uses the keyed embedded field, e.g. `dtos.ContentRuleOptionDto{ContentRuleOption: content_rule_model.ContentRuleOption{Key: ..., Name: ...}, Description: ...}` and `dtos.ContentRuleCompositionRequestDto{ContentRuleComposition: content_rule_model.ContentRuleComposition{Key: option.Key, Name: option.Name, ...}}`. In the DTO declarations the embedded field comes first, followed by a blank line (`embeddedstructfieldcheck`). Expect `golines` reflows.
- [x] Update test literals (about 40, including `guiHandler/handlerDependenciesStub_test.go` and `contentRuleHandler/getContentRuleEditorOptions_test.go`) and the two mocks so everything compiles. No assertion semantics change in this phase.

### Verification Plan
- `go build ./...` PASS; `go vet ./...`, `go vet -tags=integration_test ./...` and `go vet -tags='integration_test,gui' ./...` PASS.
- `go test ./test/unit/... -count=1` PASS.
- Grep over `app internal test cmd`: no `dtos.ContentRuleKey` / `dtos.ContentRuleEditorKind`.

### Phase Summary
Models, key/kind move and DTO reshape done as listed; `internal/dtos/contentRuleKey.go` and
`contentRuleEditorKind.go` deleted. To avoid rewriting the same files twice, the service and
`zoneContentHandler` test migrations were folded into Phase 2, so Phase 1's checks ran together
with Phase 2's: build, `go vet` under no tag / `integration_test` / `integration_test,gui` all
PASS, unit PASS, grep clean. Adjacent gap closed: `contentRuleKeyFromName` (71.4% at baseline)
now has a table test over all five rule names in
`contentRuleHandler/describeContentRule_test.go`. `contentRuleHandler.go` keeps its original
function shapes (an interim helper extraction for `funlen` was reverted once literals went flat).
`handler_interfaces/`, `guiHandler.go` and `guiHandler/handlerDependenciesStub_test.go` (zero-value
DTO literals only) needed no change.

## Phase 2: Service speaks models; handlers own conversion
Status: Complete

- [x] Change the four service methods and the interface to the "After" signatures in §1. Delete `validRule`. Remove the `dtos` import from both service files.
- [x] `zoneContentHandler`:
  - `ComposeContentRule` passes `request.ContentRuleComposition` and builds `dtos.ContentRuleCompositionResultDto{Rule, Valid}` from the comma-ok pair.
  - `GetDefaultContentRules` maps `GetContentRuleEditorOptions(content).Rules` to `[]ContentRuleOption` (`option.ContentRuleOption`).
  - `describeContentRules` returns `[]content_rule_model.ContentRuleDescription` (`DescribeContentRule(...).ContentRuleDescription`).
- [x] `ZoneContentEditorServiceMock`: new signatures; `ComposeContentRule` returns `(row, arguments.Bool(1))`.
- [x] `TemplateHandlerMock.ComposeContentRule` (it delegates to the real service): `rule, valid := zone_content.NewZoneContentEditorService().ComposeContentRule(request.ContentRuleComposition)`, returning `dtos.ContentRuleCompositionResultDto{Rule: rule, Valid: valid}`.
- [x] Service unit tests (`test/unit/internal/services/zone_content/zoneContentEditorService/`), one assertion per test, AAA, `t.Parallel()`, gofakeit:
  - `composeContentRule_test.go`: migrate every existing case to the model and split rule/ok assertions where a test now checks both. Cases: unknown key → not ok and zero row; distance-to-road and distance-to-town valid → row with name and distance; distance index `-1` and `len` → not ok; guarded `true` and `false` → non-nil pointer with that value; solo encounter `true` and `false` likewise; variant valid → `VariantID` pointer to the selected ID; variant index `-1` and `len` → not ok; guarded → a single `assert.Equal` on the whole `ContentRuleRow{Name, IsGuarded: &value}` (other fields empty/nil); **new:** each valid kind (distance-to-road, distance-to-town, guarded, solo, variant) → `ok` is true, as named `t.Run` subtests or one test per kind.
  - `getDefaultContentRules_test.go`, `getContentRuleMarkers_test.go`, `getContentRowDisplayName_test.go`: migrate to model inputs, same scenarios.
- [x] Handler unit tests (`test/unit/internal/handlers/zoneContentHandler/`):
  - `composeContentRule_test.go` replaces `TestWhenAContentRuleIsComposed_ReturnsTheServiceResult` with `TestWhenAContentRuleIsComposed_TheEmbeddedCompositionReachesTheService`, `TestWhenTheServiceAcceptsTheComposition_ReturnsAValidResult` (DTO `{rule, true}`) and `TestWhenTheServiceRejectsTheComposition_ReturnsAnInvalidResult` (DTO `{rule, false}`: the handler passes the service's rule through unchanged; it does not blank it).
  - `getDefaultContentRules_test.go` replaces `TestWhenDefaultRulesAreRequested_TheEditorOptionsAreHandedToTheService` with `TestWhenDefaultRulesAreRequested_OnlyKeyAndNameReachTheService` (option DTOs carry gofakeit `Description`, `Marker`, `EditorKind`, `EditorLabel`) and keeps a test that the service result is returned. No empty-input mock case unless its expected slice matches the handler's nil-vs-empty output exactly.
  - `getContentRuleMarkers_test.go`, `getContentRowDisplayName_test.go`: the descriptions reach the service unwrapped, in order.
- [x] `contentRuleHandler` and `guiHandler` unit tests: literal updates only, unless coverage shows a new branch.

### Verification Plan
- `go build ./...` PASS; `go vet -tags=integration_test ./...` and `go vet -tags='integration_test,gui' ./...` PASS; `go test ./test/unit/... -count=1` PASS.
- `grep internal/dtos internal/services/zone_content` returns nothing.

### Phase Summary
Service, interface, `zoneContentHandler` and both mocks changed as listed; `validRule` deleted.
`GetDefaultContentRules` and `describeContentRules` convert with `linq` in the handler. Service
compose tests now: unknown key rejected / empty row; distance index `-1` and past-the-end
rejected for both distance kinds; exact-row checks for road, town, guarded and solo (true and
false, other checkbox set to the opposite value), variant; variant `-1` / past-the-end rejected;
all five valid kinds accepted. Handler tests use `AssertCalled` for the unwrapped composition and
Key/Name-only options, and exact DTO equality for accept/reject. Build, tagged vet, unit PASS;
service files import no `dtos`.

## Phase 3: Architecture gate and full verification
Status: Complete

- [x] Remove `"internal/services/zone_content"` from `dtoNamerAllowList` in [layering_test.go](../../test/unit/architecture/dependency/layering_test.go) and update its comment to describe the single remaining (bonuses) service. Nothing else in that file changes.
- [x] Grep repo docs (`AGENTS.md`, `README.md`, `QUICKSTART.md`, `.agent/backlog/`) for claims that zone_content is an accepted DTO exception; list them for close-out rather than editing the review early.
- [x] Run the full gate set from Phase 0 and compare with the baseline. Coverage must be ≥ 74.5%; if it moves, explain it through the per-file lines.
- [x] `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` PASS with **no** golden changes (no `.failure` files; never run `-update`).
- [x] Mutation check: record `Get-FileHash internal/services/zone_content/zoneContentEditorService.go`; add `_ "github.com/Tariomka/hommoe_custom_templates/internal/dtos"` to its imports; `go test -count=1 -run TestWhenDtoConsumersAreScanned ./test/unit/architecture/dependency/` must FAIL naming that file; remove the import, confirm the hash matches, rerun to PASS.

### Verification Plan
- build, unit, `./test/...`, tagged integration + GUI, testlayoutcheck, gofmt, lint 0: all PASS.
- `git status --short` lists changes only under `app/gui/dialogs/ruleDialog.go`, `internal/{dtos,handlers,handlers/handler_interfaces,models/content_rule_model,services/zone_content}/`, `test/test_helpers/`, `test/unit/{architecture/dependency,internal/handlers,internal/services/zone_content}/` and `.agent/plans/`.

### Phase Summary
Allow-list entry removed; comment now names bonuses as the single remaining service and records
the zone_content removal. Mutation check: hash `DBC75AEF…B162` before and after; with the blank
`dtos` import the DTO layering test FAILED naming `zoneContentEditorService.go`, then PASSED.
Final gates (Windows/amd64, Go 1.27.0): build PASS; coverage **74.5%**, covered 6915 / 9263
(baseline 6909 / 9261; +2 statements from the handler conversions, +6 covered, every function in
the three touched files 100% including `contentRuleKeyFromName`); `./test/...` PASS; tagged run
PASS (root 4.103s, GUI 33.104s), no `.failure` files; testlayoutcheck PASS; gofmt clean (gofmt -w
ran on the explicit list of the 5 new model files); lint 0 after `--fix` (34 findings, all in this
batch's files: gofmt/gci on the new files, golines, and 18 `modernize/embedlit`). `git status`
matches the allowed set; `handler_interfaces/` untouched.

Docs to update at close-out (not edited now): review §2.2, the §2 architecture-inventory
bullet ("The owner has reopened the zone-content DTO exception…"), §9 row J and the progress
line; `owner_findings.md` O15 entry (lines ~90-92). The review's §0 "Memory/history
invalidations" sentence is historical and stays.

## Phase 4: Review and close-out
Status: In progress (awaiting owner review and commit)

- [x] Independent implementation review (Claude Opus 5.5). Apply or explicitly decline each finding.
  APPROVE, 0 Blocker/Major/Minor, 3 Nits, all applied: empty-row assertions for out-of-range
  compositions (`TestWhenAnIndexIsOutOfRange_ReturnsAnEmptyRule`), a doc line on
  `fakeComposition`, and the `handlerDependenciesStub_test.go` note in Phase 1. Affected unit
  packages, testlayoutcheck, gofmt and lint (0) re-run PASS afterwards.
- [ ] Hand to the owner for review and commit. **Stop here until the owner commits.**
- [ ] After the owner commits: mark §2.2 `✅ FIXED` in place in the review with a Progress paragraph, complete its §9 row J, refresh the progress line (18 fixed / 14 remaining), update the architecture-inventory sentence about the reopened exception, and record J9 as "recorded, not assigned". Do not renumber anything.
- [ ] Rewrite the handoff per AGENTS.md §5.2, preserving its §8 verbatim.

### Verification Plan
- The review's §8 LF-normalized hash equals the Phase 0 value; the handoff's §8 hash equals `0DDDEBF32DB51167E675EA7A147E75F186B77643E8024F26E7B262102B004AD8`.

### Phase Summary
_(write when phase completes)_

## Final Recap
_(write when all phases complete)_

## Deployment Plan
_(write when all phases complete)_
