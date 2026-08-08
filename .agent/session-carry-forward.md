# Session Carry-Forward — Batch 15, Phase 3

## 1. Session goal

Execute Batch 15 of [plans/zone-editor-state-extraction.md](../plans/zone-editor-state-extraction.md) —
extract business logic out of the zone-editor dialog family into services behind
handler seams. Phases 0-2 were already merged; this session delivered
**Phase 3 (the four extra dialogs)** and then had to recover it after a tooling
accident destroyed the working tree.

## 2. Fixes applied

- Recovered the entire Phase 3 change set from VS Code Local History (102 files)
  after a repo-wide PowerShell line-ending sweep overwrote every `.go` file.
- Re-applied three edits Local History could not supply because they postdated
  the last editor save:
  - [app/gui/dialogs/pickerDialog_testexports.go](../app/gui/dialogs/pickerDialog_testexports.go) - `entry.id` to `entry.ID`, `entry.haystack` to `entry.Haystack`.
  - [app/gui/dialogs/bonusPickerDialog_testexports.go](../app/gui/dialogs/bonusPickerDialog_testexports.go) - `this.bonusHandler` to `this.handler`.
  - [test/integration/gui/pickerDialog_integration_test.go](../test/integration/gui/pickerDialog_integration_test.go) - `composition.InitializeGuiHandler()` added at 10 picker-construction call sites.
- Regenerated [internal/composition/wire_gen.go](../internal/composition/wire_gen.go) (no history existed for it).
- Re-fixed four `revive` var-naming findings (`variantId` to `variantID`, `spellId` to `spellID`).
- Corrected `getDefaultContentRules_test.go` to assert the real behaviour (a
  Guarded default rule is always returned) rather than `assert.Empty`.

## 3. Features added / changed

Three new services, three new handler facets, all four dialogs reduced to pure
rendering:

- **`internal/services/bonuses`** (`IBonusEntryService`) - bonus validation,
  `config.BonusEntry` construction, duplicate filtering, spell-count label.
  Rationale: `bonusPickerDialog.go` had no handler at all and `app/**` may not
  import `internal/services` (depguard), so a new `IBonusHandler` was required.
- **`internal/services/zone_content`** (`IZoneContentEditorService`) - content
  rule composition/upsert/defaults, marker joining, display naming, alphabetical
  sorting, count clamping. Exposed via `IZoneContentHandler`, which **embeds**
  `IContentRuleHandler`. Rationale: widening `IContentRuleHandler` itself was
  tried and reverted - it changed `NewContentRuleHandler`'s signature and forced
  edits at 14 existing test call sites for no behavioural gain.
- **`internal/services/pickers`** (`IPickerEntryService`) - source-to-entry
  mapping, filter normalisation, grouped row model, selected-ID extraction. The
  GUI's private `pickerEntry` struct was deleted in favour of
  `dtos.PickerEntryDto`. Rationale: catalogs stay in `app/gui/constants` (which
  imports `internal/registry`) because `internal/` may not import `app/`, so the
  service accepts catalog rows as DTOs.
- `BonusPickerDialog` and `BonusesPanel` are typed on `IGuiHandler` because they
  genuinely need both the bonus and picker facets.
- `NewGuiHandler` grew to seven dependencies; `IGuiHandler` now embeds the three
  new interfaces. `providerSets.go` gained three services and three handlers.
  Still **zero `wire.Bind` calls**.

## 4. File modifications

**Modified (tracked):**

| File | Change |
| --- | --- |
| [app/gui/dialogs/bonusPickerDialog.go](../app/gui/dialogs/bonusPickerDialog.go) | logic removed; field renamed `bonusHandler` to `handler`, typed `IGuiHandler` |
| [app/gui/dialogs/bonusPickerDialog_testexports.go](../app/gui/dialogs/bonusPickerDialog_testexports.go) | follows the field rename |
| [app/gui/dialogs/pickerDialog.go](../app/gui/dialogs/pickerDialog.go) | `pickerEntry` deleted; renders `dtos.PickerRowDto`; all 4 constructors take `IPickerHandler` |
| [app/gui/dialogs/pickerDialog_testexports.go](../app/gui/dialogs/pickerDialog_testexports.go) | DTO field names |
| [app/gui/dialogs/ruleDialog.go](../app/gui/dialogs/ruleDialog.go) | rule building delegated; handler annotation widened |
| [app/gui/dialogs/zoneContent.go](../app/gui/dialogs/zoneContent.go) | display name / defaults / markers / sort / clamp delegated |
| [app/gui/dialogs/zoneContentDialog.go](../app/gui/dialogs/zoneContentDialog.go) | handler annotation widened |
| [app/gui/panels/layoutPanel.go](../app/gui/panels/layoutPanel.go) | handler annotation widened |
| [app/gui/panels/bonusesPanel.go](../app/gui/panels/bonusesPanel.go) | takes `IGuiHandler` |
| [app/gui/editor/window.go](../app/gui/editor/window.go) | `NewBonusesPanel` call updated |
| [internal/handlers/guiHandler.go](../internal/handlers/guiHandler.go) | 7 dependencies + 17 pass-throughs |
| [internal/handlers/handler_interfaces/guiHandlerInterface.go](../internal/handlers/handler_interfaces/guiHandlerInterface.go) | embeds 3 new interfaces |
| [internal/composition/providerSets.go](../internal/composition/providerSets.go) | 3 services + 3 handlers registered |
| [internal/composition/wire_gen.go](../internal/composition/wire_gen.go) | regenerated |
| [test/test_helpers/templateHandlerMock.go](../test/test_helpers/templateHandlerMock.go) | implements the full `IGuiHandler` |
| [test/integration/editorState_integration_test.go](../test/integration/editorState_integration_test.go) | `NewBonusesPanel` call updated |
| [test/integration/gui/bonusPickerDialog_integration_test.go](../test/integration/gui/bonusPickerDialog_integration_test.go) | `newBonusPicker` fixture helper, 17 call sites |
| [test/integration/gui/pickerDialog_integration_test.go](../test/integration/gui/pickerDialog_integration_test.go) | 10 call sites take a real handler |
| [test/unit/internal/handlers/guiHandler/handlerDependenciesStub_test.go](../test/unit/internal/handlers/guiHandler/handlerDependenciesStub_test.go) | stub implements the new facets |
| [test/unit/internal/handlers/guiHandler/newGuiHandler_test.go](../test/unit/internal/handlers/guiHandler/newGuiHandler_test.go) | 7-arg construction |
| [plans/zone-editor-state-extraction.md](../plans/zone-editor-state-extraction.md) | Phase 3 ticked + Phase Summary written |

**Created:** 9 DTOs in `internal/dtos/`; 3 handlers + 3 handler interfaces;
`internal/services/{bonuses,pickers,zone_content}/` (2 files each);
3 testify mocks in `test/test_helpers/`.

**Deleted:** the private `pickerEntry` struct (folded into `dtos.PickerEntryDto`).

**Deliberately untouched:** [internal/handlers/contentRuleHandler.go](../internal/handlers/contentRuleHandler.go)
and [internal/handlers/handler_interfaces/contentRuleHandlerInterface.go](../internal/handlers/handler_interfaces/contentRuleHandlerInterface.go)
- the widening attempt was reverted, they are at their original state.

## 5. Tests added or updated

New unit-test folders (AGENTS.md 4.6 layout, one file per public function):

- `test/unit/internal/services/bonuses/bonusEntryService/` (5)
- `test/unit/internal/services/zone_content/zoneContentEditorService/` (8)
- `test/unit/internal/services/pickers/pickerEntryService/` (7)
- `test/unit/internal/handlers/bonusHandler/` (6)
- `test/unit/internal/handlers/zoneContentHandler/` (9)
- `test/unit/internal/handlers/pickerHandler/` (8)
- 17 new files in `test/unit/internal/handlers/guiHandler/`

**Last verification run - all green:**

| Check | Result |
| --- | --- |
| `go build ./...` | ok |
| `go vet -tags="integration_test,gui" ./...` | ok |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test ./test/unit/... -count=1` | pass |
| `go test -tags=integration_test ./test/integration/... -count=1` | `ok 3.670s` |
| `go test -tags "integration_test,gui" ./test/integration/gui/... -count=1` | `ok 4.347s`, snapshots unchanged |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `golangci-lint-v2 run ./...` | `0 issues.` |
| Unit coverage | **72.7%** (floor 69.3%, Phase 2 left it at 70.9%) |

## 6. Git status snapshot

Branch: `AD/refactoring-07-21` (HEAD `c0499e2`, "Batch 15 Prep").
Nothing was staged or committed by the agent.

```
 M .vscode/settings.json          <- NOT touched by the agent, pre-existing
 M app/gui/dialogs/bonusPickerDialog.go
 M app/gui/dialogs/bonusPickerDialog_testexports.go
 M app/gui/dialogs/pickerDialog.go
 M app/gui/dialogs/pickerDialog_testexports.go
 M app/gui/dialogs/ruleDialog.go
 M app/gui/dialogs/zoneContent.go
 M app/gui/dialogs/zoneContentDialog.go
 M app/gui/editor/window.go
 M app/gui/panels/bonusesPanel.go
 M app/gui/panels/layoutPanel.go
 M internal/composition/providerSets.go
 M internal/composition/wire_gen.go
 M internal/handlers/guiHandler.go
 M internal/handlers/handler_interfaces/guiHandlerInterface.go
 M plans/zone-editor-state-extraction.md
 M test/integration/editorState_integration_test.go
 M test/integration/gui/bonusPickerDialog_integration_test.go
 M test/integration/gui/pickerDialog_integration_test.go
 M test/test_helpers/templateHandlerMock.go
 M test/unit/internal/handlers/guiHandler/handlerDependenciesStub_test.go
 M test/unit/internal/handlers/guiHandler/newGuiHandler_test.go
?? internal/dtos/{bonusCompositionRequest,bonusCompositionResult,contentRuleCompositionRequest,contentRuleCompositionResult,existingBonuses,pickerEntry,pickerItem,pickerRow,pickerSpell}Dto.go
?? internal/handlers/{bonusHandler,pickerHandler,zoneContentHandler}.go
?? internal/handlers/handler_interfaces/{bonus,picker,zoneContent}HandlerInterface.go
?? internal/services/{bonuses,pickers,zone_content}/
?? test/test_helpers/{bonusEntryService,pickerEntryService,zoneContentEditorService}Mock.go
?? test/unit/internal/handlers/{bonusHandler,pickerHandler,zoneContentHandler}/
?? test/unit/internal/handlers/guiHandler/  (17 new files)
?? test/unit/internal/services/{bonuses,pickers,zone_content}/
```

Phase 3 is entirely **unstaged and uncommitted** - it is the next thing to review.
The previous batch's handoff was archived to
[.agent/session-carry-forward-batch14.md](session-carry-forward-batch14.md).

## 7. Rejections / declined

- **Widening `IContentRuleHandler`** to carry the zone-content methods: written,
  then reverted. It changed `NewContentRuleHandler`'s public signature and broke
  14 existing test call sites. Replaced by a separate `IZoneContentHandler` that
  embeds it.
- **Adding `wire.Bind` calls** for the new services: not needed, providers already
  return interfaces (existing repo rule).
- **Splitting `BonusPickerDialog`'s handler into two parameters** (one bonus, one
  picker): rejected, one collaborator should be one parameter.

## 8. The corruption incident (post-mortem)

At the end of Phase 3, a repo-wide CRLF-to-LF normalisation was attempted with
`Get-ChildItem -Recurse -Filter *.go | ForEach-Object { ... [System.IO.File]::WriteAllText($p, $n, ...) }`.
The first attempt failed because `$t.Replace([char]13+[char]10, [char]10)` made
PowerShell do char arithmetic; the "corrected" second run then **overwrote roughly
1170 `.go` files with identical content**, destroying the entire working tree.

The owner reset the branch to the pre-Phase-3 commit. Phase 3 was recovered from
VS Code Local History (`%APPDATA%\Code\User\History\*\entries.json`), which held
clean revisions for 102 of the affected files. Three files had no usable
revision (they were last written by tooling, not saved through the editor) and
were re-applied by hand; `wire_gen.go` was regenerated.

**Rules adopted (also recorded in repository memory):**

- Never run a bulk in-place rewrite across the repository.
- To normalise formatting or line endings, run `gofmt -w` on an **explicit** list
  from `gofmt -l` - gofmt converts CRLF to LF and cannot mangle content.
- Local History only covers files saved through the editor, so it is never a
  complete backup.

## 9. Open questions

- None blocking. The only outstanding judgement call is whether the owner wants
  Phase 3 committed before Phase 4 starts (recommended - Phase 4 touches
  `zoneEditorDialog.go`, which Phase 3 does not).

## 10. Next recommended actions

1. **Owner reviews and commits Phase 3.** Nothing is staged.
2. **Phase 4** ([plans/zone-editor-state-extraction.md](../plans/zone-editor-state-extraction.md),
   "0.2 - make the reset button honest"): relabel "Reset to generated" to wording
   that matches what it does, delete the marker comment at
   `app/gui/dialogs/zoneEditorDialog.go` around line 213, and clear the persisted
   manual-edit snapshot on the Apply that follows a revert. Governed by owner
   decisions 6 and 7 - re-read them first. Coverage must not drop and "should
   raise it noticeably".
3. **Phase 5**: close out 2.6 in [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md).

## 11. Carry-forward prompt

> Read `AGENTS.md` first. Hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/`; keep everything
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chains with `;`
> never `&&`); every change ships with tests and must not drop coverage below
> **69.3%** (currently 72.7%); durable multi-session work gets a plan file under
> `plans/`; never stage and never commit - the owner reviews and commits; never
> change where `.rmg.json` is written and never persist the output directory
> (2.6).
>
> **NEVER run a bulk in-place rewrite over the repository** (e.g.
> `Get-ChildItem -Recurse -Filter *.go | ForEach-Object { WriteAllText }`). A
> CRLF-to-LF sweep of that shape destroyed every `.go` file last session. To fix
> formatting or line endings, run `gofmt -w` on an explicit file list produced by
> `gofmt -l`.
>
> Work continues on `plans/zone-editor-state-extraction.md`. **Phases 0-3 are
> complete**; Phase 3 is finished and fully verified (build, vet,
> testlayoutcheck, unit, integration, GPU-integration with unchanged snapshots,
> `0 issues.` lint, 72.7% coverage) but **uncommitted and unstaged** - ask the
> owner to review it before starting new work. **Start Phase 4** ("0.2 - make the
> reset button honest"); owner decisions 6 and 7 in the plan govern it and are
> settled - do not re-litigate them.
>
> Full handoff, including the file-by-file change list and the
> corruption-incident post-mortem, is in `./.agent/session-carry-forward.md`.
