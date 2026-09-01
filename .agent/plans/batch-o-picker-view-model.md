# Batch O — picker view-model move + layering carve-out (backlog §2.6 step 1)

Stop `internal/services` from naming DTOs. `internal/services/pickers` turns out to
be view-model logic, so it moves to `app/gui/models/` and the service, its
handler and its four DTOs are deleted. `bonuses` and `zone_content` keep their
DTOs under a written carve-out.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is
done, set its status to `Complete` and write its **Phase Summary**; run the
phase's **Verification Plan** and record the result before moving on. When all
phases are done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../../AGENTS.md) first — especially **§4.4.1** (the
Entity/Model/DTO doctrine) and **§4.5** (UI vs business logic). Hard rules that
bite here: never stage and never commit; delete with `Remove-Item`, never
`git rm`; move with `Move-Item`, never `git mv`; never round-trip a `.go` file
through `Get-Content`/`Set-Content`; unit coverage must not drop below **72.5 %**
(currently **73.9 %**); lint baseline is **0 issues**.

## 0. Decisions (settled with the owner 2026-09-01 — do not relitigate)

1. **A wrong premise was corrected before any code was written.** The first scope
   read claimed six of the eleven DTOs "never cross into `app/`". That measured
   where a type is *named*, not where it *crosses* — all six are consumed with
   `:=`, so the name never appears while the type crosses just the same
   (`ExistingBonusesDto` at `bonusPickerDialog.go:69`, `BonusCompositionResultDto`
   at `:221`, `PickerRowDto` at `pickerDialog.go:227`,
   `ContentRuleEditorOptionsDto` at `ruleDialog.go:55`,
   `ContentRuleCompositionResultDto` at `:275`, `ContentRuleDescriptionDto` at
   `:299`). **All eleven are genuine crossing DTOs.** Do not re-propose a
   "move the non-crossing ones for free" plan; there are none.
2. **`internal/services/pickers` is deleted, not modelled.** Its types are
   `Label`, `Badge`, `Trailing`, `Haystack` (lowercased search text),
   `IsGroupHeader`, `GroupMatchCount` — what the picker dialog draws, nothing
   more. The catalogues live in `app/gui`. It sat in `internal/services` only
   because §4.5 bans logic in GUI files.
3. **It lands in `app/gui/models/`**, beside the existing `EditorState` view
   state — not in `app/gui/dialogs/` as private helpers, which would push its
   tests into the GPU suite and lose unit coverage. The logic is pure, so the
   tests stay **unit** tests under `test/unit/app/gui/models/`.
4. **AGENTS.md §4.5 is left alone.** Its "extract it into a service" sentence
   reads as if it forbids this move; the owner accepts that and does not want a
   view-model carve-out written into the instructions.
5. **`bonuses` and `zone_content` keep their DTOs**, under a written
   justification in the layering allow-list. Mirroring seven types
   (`ExistingBonusesDto`, `BonusCompositionRequestDto`,
   `BonusCompositionResultDto`, `ContentRuleCompositionRequestDto`,
   `ContentRuleCompositionResultDto`, `ContentRuleEditorOptionsDto`,
   `ContentRuleDescriptionDto`, plus `ContentRuleOptionDto`,
   `ContentRuleVariantOptionDto` and `ContentRuleKey` in their wake) with a
   Model twin plus handler slice converters buys no behaviour and no clarity.
6. **The allow-list shrinks from 3 packages to 2, and stops there.** That is the
   end state for §2.6 step 1, not a way-point.

## Phase 1: Move the picker view-model into `app/gui/models/`
Status: **Complete** (2026-09-01) — uncommitted, awaiting review.

Phases 1's deletions and the `app/` rewrite are **one compile unit** — removing
`IPickerHandler` breaks `pickerDialog.go` in the same build, so they cannot be
split (the lesson from Batch I §0.6 decision 4).

- [ ] New files in [app/gui/models/](../../app/gui/models/), one struct per file
      (§4.1), package-level functions rather than an injected service (there is
      no state):
      - `pickerItem.go` — `PickerItem{Sid, Name, Category}`
      - `pickerSpell.go` — `PickerSpell{Sid, Name, School, SchoolDisplayName, Tier}`
      - `pickerEntry.go` — `PickerEntry` plus `BuildItemPickerEntries`,
        `BuildSpellPickerEntries`, `BuildValueOverridePickerEntries`,
        `NormalizePickerFilter`, `GetSelectedPickerIDs`
      - `pickerRow.go` — `PickerRow` plus `GetVisiblePickerRows` and the private
        `countGroupMatches`
- [ ] Rewrite [app/gui/dialogs/pickerDialog.go](../../app/gui/dialogs/pickerDialog.go)
      to call those functions directly and drop its `pickerHandler` field and
      constructor parameter. Update every `NewPickerDialog` call site in
      [bonusesPanel.go](../../app/gui/panels/bonusesPanel.go) and
      [bonusPickerDialog.go](../../app/gui/dialogs/bonusPickerDialog.go).
- [ ] `Remove-Item`: `internal/services/pickers/` (2 files),
      `internal/handlers/pickerHandler.go`,
      `internal/handlers/handler_interfaces/pickerHandlerInterface.go`, and the
      four DTOs `pickerItemDto.go`, `pickerSpellDto.go`, `pickerEntryDto.go`,
      `pickerRowDto.go`.
- [ ] Drop `IPickerHandler` from `IGuiHandler`
      ([guiHandlerInterface.go](../../internal/handlers/handler_interfaces/guiHandlerInterface.go)),
      the six passthroughs from [guiHandler.go](../../internal/handlers/guiHandler.go),
      and `NewPickerEntryService` / `NewPickerHandler` from
      [providerSets.go](../../internal/composition/providerSets.go). Regenerate
      with `wire gen ./internal/composition/...` — never hand-edit `wire_gen.go`.
- [ ] `Remove-Item` `test/test_helpers/pickerEntryServiceMock.go` and the picker
      methods on `TemplateHandlerMock`.
- [ ] Move the 15 test files: the 7 under
      `test/unit/internal/services/pickers/pickerEntryService/` become
      `test/unit/app/gui/models/pickerEntry/` and `.../pickerRow/` per §4.6;
      **delete** the 8 under `test/unit/internal/handlers/pickerHandler/` — they
      tested passthroughs on a handler that no longer exists.

### Verification Plan
- `go build ./...`, both `go vet` variants, `go run ./cmd/testlayoutcheck .`.
- `wire diff ./internal/composition/...` exit 0 after regeneration.
- No file under `internal/` names a picker type:
  `Select-String -Pattern 'Picker(Item|Spell|Entry|Row)' -Path internal` empty.
- Unit + untagged + integration suites pass; **GPU suite passes without
  `-update`** — the picker dialogs are snapshotted, so the rendering must be
  proven byte-identical.
- Coverage ≥ 72.5 %. The picker logic keeps unit coverage; only the deleted
  passthrough tests go.

### Phase Summary

**Landed.** `internal/services/pickers` no longer exists. The logic is four files
in [app/gui/models/](../../app/gui/models/) — `pickerItem.go`, `pickerSpell.go`,
`pickerEntry.go` (the three builders plus `NormalizePickerFilter` and
`GetSelectedPickerIDs`) and `pickerRow.go` (`GetVisiblePickerRows` +
`countGroupMatches`) — as package-level functions, since the service was
stateless and there was nothing to inject.

**Deleted:** the service (2 files), `pickerHandler.go`,
`pickerHandlerInterface.go`, the four picker DTOs, `PickerEntryServiceMock`, the
`IPickerHandler` embed in `IGuiHandler`, six passthroughs on `GUIHandler`, six on
`TemplateHandlerMock`, six on the `guiHandler` test stub (and its now-unused
`pickerCalled` field), and two wire providers. `NewGuiHandler` lost a parameter;
`wire_gen.go` was **regenerated**, not hand-edited.

**Behaviour is unchanged.** The three builders were rewritten onto
`linq.SelectSlice`, which returns `nil` rather than an empty slice for empty
input — checked against the tests first, which assert `assert.Empty` and pass
either way.

**Tests moved, not lost.** The 7 service tests became 5 under
`test/unit/app/gui/models/pickerEntry/` and 1 under `.../pickerRow/`; the 8
`pickerHandler` tests and the 6 `guiHandler` picker tests were **deleted** — they
tested passthroughs on a handler that no longer exists.
`newPickerEntryService_test.go` and `newPickerHandler_test.go` went with their
constructors, which is the whole of the −0.1 pp coverage move.

**Gates** — all green: `go build ./...`, both `go vet` variants, `gofmt -l` empty,
`go run ./cmd/testlayoutcheck .` passed, `wire diff` exit 0, unit + untagged +
integration, and the **GPU suite without `-update`** (22.9 s) — the picker
dialogs are snapshotted, so the rendering is proven byte-identical. Coverage
**73.8 %** (floor 72.5 %). `golangci-lint-v2` **0 issues**.

No file under `internal/` names a picker type any more.

**Owner review, 2026-09-01.** The reviewed tree is **staged, not committed**, on
top of `2bfc8cb`. Two changes came out of the review, both kept:

- `app/gui/dialogs/pickerDialog_testexports.go` lost 4 lines — an accessor that
  existed only to reach the picker handler. The remaining accessors
  (`EntryIDs`, `MatchingEntryIDs`, `RowCount`, `ClickEntry`, `SetSearch`,
  `ClickAdd`, `ClickCancel`) work directly off `this.entries`, so the
  `integration_test` gate on that file is still earned.
- `internal/dtos/pickerSpellDto.go` was **renamed** into
  `app/gui/models/pickerSpell.go` rather than deleted-and-recreated, so git
  records the move. Same for four of the six moved test files. Cosmetic in the
  tree, useful in the history.

All gates were re-run against the reviewed tree and are still green with the
same numbers.

## Phase 2: Write the carve-out for `bonuses` and `zone_content`
Status: In progress — the allow-list is done, the backlog is not. **Docs only,
no Go code left in this phase.**

- [x] Shrink `dtoNamerAllowList` in
      [layering_test.go](../../test/unit/architecture/dependency/layering_test.go)
      to two entries and replace its comment: these two stay by decision, not by
      debt.
- [ ] Rewrite backlog [§2.6](../backlog/backlog-opus5.md) step 1 to record that
      it is **done and closed at two entries**, not shrinking to zero. Fold in
      the §0.1 correction — all eleven DTOs cross; there were no free moves.
      Steps 2–4 of §2.6 (the 113-file entity list) are untouched by batch O and
      stay open.
- [ ] Retitle §2.6 if its heading still implies the DTO list is shrinking to
      zero, and update the §8 batch table to add row **O** as done.
- [x] Update `.agent/session-carry-forward.md`.

### Verification Plan
- `go test ./test/unit/architecture/... -count=1` passes, and fails if either
  remaining entry is removed — prove the list is still load-bearing.
- Every gate from Phase 1 still green.

### Phase Summary
_(write when phase completes)_

## Final Recap
_(write when all phases complete)_

## Deployment Plan
_(write when all phases complete)_
