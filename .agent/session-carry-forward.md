# Session Carry-Forward — 2026-09-01 (Batch O, phase 1 done)

## 1. Session goal

Close **Batch I** (phase 7 docs), then start **Batch O** — backlog §2.6 step 1,
stopping `internal/services` from naming DTOs. Batch I is committed. Batch O
**phase 1 is complete and staged**; **phase 2 is docs-only and unfinished**.

Plan: [.agent/plans/batch-o-picker-view-model.md](plans/batch-o-picker-view-model.md).

## 2. Read this before touching anything

**`AGENTS.md` §4.4.1** carries the Entity/Model/DTO doctrine (it moved out of the
Batch I plan file, which was deleted). The three things most often gotten wrong:

- **The Model owns the structure.** *"Redefinition is expected in Models, but it
  should never happen in DTOs."* `EditorStateDto` is literally
  `struct { editor_state_model.EditorState }`. **A DTO carrying a Model is
  intended** — do not "fix" `EditorStateValidationDto.State`,
  `CastleSettingsReapplyRequestDto.Changes` or
  `ManualEditDecisionDto.ReapplyWithCastleChanges`.
- **`app/` may hold a Model**; only the *crossing* into `internal/` is a DTO.
  `app/` → `internal/models` is fine; `app/` → `internal/mappers`,
  `internal/services`, `internal/repositories`, `internal/validators` is not.
- The doctrine is enforced by
  [test/unit/architecture/dependency/layering_test.go](../test/unit/architecture/dependency/layering_test.go).
  **Its allow-lists only ever shrink — never add an entry, clean the package.**

## 3. The mistake I made, so you don't repeat it

Scoping batch O, I reported that six of the eleven DTOs "never cross into
`app/`" and could therefore be moved to `internal/models/` for free. **That was
wrong.** I had grepped for where each type is *named* in `app/`, not where it
*crosses*. All six are consumed with `:=`, so the identifier never appears while
the type crosses just the same:

| Type | Crossing call site |
| --- | --- |
| `ExistingBonusesDto` | `bonusPickerDialog.go:69` — `summary := handler.DescribeExistingBonuses(...)` |
| `BonusCompositionResultDto` | `bonusPickerDialog.go:221` — `result := this.handler.BuildBonusEntries(...)` |
| `PickerRowDto` | `pickerDialog.go:227` — `model := this.pickerHandler.GetVisiblePickerRows(...)` |
| `ContentRuleEditorOptionsDto` | `ruleDialog.go:55` — `options := contentRuleHandler.GetContentRuleEditorOptions(...)` |
| `ContentRuleCompositionResultDto` | `ruleDialog.go:275` — `result := ...ComposeContentRule(...)` |
| `ContentRuleDescriptionDto` | `ruleDialog.go:299` — `...DescribeContentRule(...).DisplayText` |

**All eleven are genuine crossing DTOs.** When measuring whether a type crosses a
boundary, follow the *method signatures the other side calls*, never a grep for
the type name. This is recorded as §0.1 of the batch-O plan; do not re-propose a
"move the non-crossing ones for free" plan, because there are none.

## 4. Owner decisions for batch O (settled — do not relitigate)

1. **`internal/services/pickers` is deleted, not modelled.** Its types were
   `Label`, `Badge`, `Trailing`, `Haystack`, `IsGroupHeader`, `GroupMatchCount` —
   what the picker dialog draws. It sat in `internal/services` only because
   AGENTS §4.5 bans logic in GUI files.
2. **It lands in `app/gui/models/`**, not as private helpers in
   `app/gui/dialogs/` — that would push its tests into the GPU suite and lose
   unit coverage.
3. **AGENTS.md §4.5 is left alone.** Its "extract it into a service" sentence
   reads as if it forbids this move; the owner accepts that and does **not** want
   a view-model carve-out written into the instructions. Do not add one.
4. **`bonuses` and `zone_content` keep their DTOs**, under a written
   justification in the allow-list. Mirroring seven types with Model twins plus
   handler slice converters buys no behaviour and no clarity.
5. **The DTO allow-list ends at two entries.** That is the end state for §2.6
   step 1, not a way-point.

## 5. What shipped this session

### Batch I phase 7 — docs (committed, `1aff58f` + `0b3d1dd`)

- `README.md` — generation flow no longer claims the GUI "collects widget input
  into `dtos.EditorStateDto`"; layer descriptions and directory tree corrected.
- Backlog §2.1 and §1.5 rewritten as self-contained ✅ FIXED entries; counts,
  batch table and §9 baselines refreshed.
- `test_observations.md` — new entry: the per-frame allocation budget has **no
  automated guard** (the benchmark needs a GPU and never runs in CI; an
  `AllocsPerRun` assertion was rejected as too flaky over a Gio frame).
- **408 dangling markdown links repaired** (156 backlog + 252 review doc), verified
  byte-safe: 393 insertions, 393 deletions, no line-ending or BOM drift.
- `AGENTS.md` §4.2.2's interface-placement example was factually wrong and is fixed.
- **`.agent/plans/batch-i-editor-state-rework.md` was deleted** per its own
  doc-lifecycle rule. Recover with
  `git show 938ef55:.agent/plans/batch-i-editor-state-rework.md`.

### Batch O phase 1 — the picker move (staged, not committed)

`internal/services/pickers` no longer exists. The logic is four files in
[app/gui/models/](../app/gui/models/) as **package-level functions** (the service
was stateless, so there was nothing to inject):

| File | Holds |
| --- | --- |
| `pickerItem.go` | `PickerItem` |
| `pickerSpell.go` | `PickerSpell` |
| `pickerEntry.go` | `PickerEntry`, the three builders, `NormalizePickerFilter`, `GetSelectedPickerIDs` |
| `pickerRow.go` | `PickerRow`, `GetVisiblePickerRows`, `countGroupMatches` |

**Deleted:** the service (2 files), `pickerHandler.go`,
`pickerHandlerInterface.go`, the four picker DTOs, `PickerEntryServiceMock`, the
`IPickerHandler` embed in `IGuiHandler`, six passthroughs on `GUIHandler`, six on
`TemplateHandlerMock`, six on the `guiHandler` test stub (plus its now-unused
`pickerCalled` field), and two wire providers. `NewGuiHandler` lost a parameter;
`wire_gen.go` was **regenerated**, never hand-edited.

**One behavioural note.** The three builders now use `linq.SelectSlice`, which
returns `nil` rather than an empty slice for empty input. Checked against the
tests *first* — they assert `assert.Empty`, which passes either way.

**Tests moved, not lost.** The 7 service tests became 5 under
`test/unit/app/gui/models/pickerEntry/` and 1 under `.../pickerRow/`. The 8
`pickerHandler` and 6 `guiHandler` picker tests were **deleted** — they tested
passthroughs on a handler that no longer exists.
`newPickerEntryService_test.go` and `newPickerHandler_test.go` went with their
constructors, which is the entire −0.1 pp coverage move.

**Owner review changes, both kept:** `pickerDialog_testexports.go` lost 4 lines
(an accessor that only existed to reach the picker handler), and
`internal/dtos/pickerSpellDto.go` was **renamed** into
`app/gui/models/pickerSpell.go` rather than deleted-and-recreated so git records
the move — same for four of the six moved test files.

## 6. Gates (re-run against the reviewed, staged tree)

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | exit 0 |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 |
| Unit / untagged / integration | pass |
| GPU suite, **no `-update`** | pass (22.9 s) |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **73.8 %** (floor 72.5 %) |

The GPU suite passing without `-update` is the load-bearing one here: the picker
dialogs are snapshotted, so the rendering is proven byte-identical.

## 7. Git status snapshot

- **Branch:** `AD/fixing_some_stuff_08-12`
- **HEAD:** `2bfc8cb "Quick update"`
- **Everything for batch O phase 1 is STAGED, not committed** — the owner staged
  it during review. **Do not unstage it** (AGENTS §2.5): the owner stages and
  commits, the agent never does.
- The only unstaged files are this one and
  `.agent/plans/batch-o-picker-view-model.md`.

## 8. Rejections / things not done

- **Rejected — the "move six DTOs for free" plan.** It was based on a bad
  measurement (§3). There are no free moves.
- **Rejected — Model twins for `bonuses` and `zone_content`.** Seven types plus
  handler slice converters for no behaviour change.
- **Rejected — a view-model carve-out sentence in AGENTS §4.5.**
- **Rejected — putting the picker logic in `app/gui/dialogs/` as private
  helpers.** It would have moved its tests into the GPU suite.
- **Not done — §2.6 steps 2–4** (the 113-file entity allow-list). Untouched by
  batch O and still open.

## 9. Open questions

1. **None block phase 2.** It is pure documentation.
2. **Repo memory duplication** (`/memories/repo/conventions.md`) — flagged nine
   sessions running: ~1,300 lines, roughly four copies of the same body. Still
   needs a dedupe pass.
3. §2.6 step 2 still asks a real design question before that list can be drained:
   is `internal/dtos` / `internal/handlers` naming `entities.Zone` a breach at
   all, or does the `.rmg.json` schema vocabulary deserve a documented carve-out
   like `internal/helpers/data` already has?

## 10. Next recommended actions

1. **Finish batch O phase 2** — three doc edits, no Go code:
   - Rewrite backlog §2.6 step 1 as **done and closed at two entries** (not
     shrinking to zero), folding in the §3 correction.
   - Retitle §2.6 if its heading implies the DTO list drains to zero.
   - Add row **O** to the §8 batch table as done.
2. Commit, then pick up **batch J** (backlog §2.2 branch B — zone tier single
   source of truth), which benefits from the model layer Batch I built.

## 11. Carry-forward prompt

> Read `AGENTS.md` first — especially **§4.4.1**, which carries the
> Entity/Model/DTO doctrine. In one line: **Entity** (`internal/entities/`) is
> the database layer, json tags only; **Model** (`internal/models/`) is the
> service layer and **owns the structure and all business logic**; **DTO**
> (`internal/dtos/`) is the `app/` ↔ `internal/` crossing and is thin.
> *"Redefinition is expected in Models, but it should never happen in DTOs"* — so
> `EditorStateDto` is literally `struct { editor_state_model.EditorState }`, and
> **a DTO carrying a Model is intended**. `app/` MAY hold a Model; only the
> crossing must be a DTO. Do not "fix" either. The doctrine is enforced by
> `test/unit/architecture/dependency/layering_test.go`; its allow-lists **only
> ever shrink** — never add an entry, clean the package instead.
>
> Then read `.agent/plans/batch-o-picker-view-model.md`. **Phase 1 is complete
> and STAGED (not committed) on top of `2bfc8cb`** — do not unstage it. **Phase 2
> is documentation only**: rewrite backlog §2.6 step 1 as *done and closed at two
> allow-list entries* rather than shrinking to zero, retitle §2.6 if its heading
> implies otherwise, and add row **O** to the §8 batch table. §2.6 steps 2–4 (the
> 113-file entity list) are out of scope and stay open. After that, batch J
> (§2.2 branch B).
>
> **A measurement mistake from last session, worth internalising:** to decide
> whether a type crosses a layer boundary, follow the *method signatures the
> other side calls* — never grep for the type name. Six DTOs consumed with `:=`
> looked like they never reached `app/` and did.
>
> The hard rules, one line each: never modify `data/`, `internal/registry/`, or
> **anything under `internal/entities/template/`**; everything must build and run
> on Windows and Linux (use `path/filepath`; chain PowerShell with `;`, never
> `&&`); every change ships with tests and unit coverage must not drop below
> 72.5 % (currently 73.8 %); the lint baseline is **0 issues**; **never stage and
> never commit** — the owner reviews, stages and commits, so **use `Move-Item`,
> never `git mv`**, and delete with `Remove-Item`, never `git rm`; never change
> where `.rmg.json` is written and never persist the output directory; never run
> a bulk in-place rewrite over the repository, and **never round-trip a `.go`
> file through `Get-Content`/`Set-Content`** — it joins every line and corrupts
> the file. For bulk text edits on *markdown*, use
> `[System.IO.File]::ReadAllText`/`WriteAllText` with an explicit
> `UTF8Encoding($false)`, then verify insertions == deletions.
>
> Standing traps: **nil is load-bearing** on the regeneration path (nil
> `Previous` = first generation, nil `Next` = unarmed debounce); the two frozen
> fixtures under `test/test_helpers/testdata/` plus the untagged
> `editorStateWireFormat_integration_test.go` must keep passing **unchanged**,
> comparing **parsed objects, never bytes**; the picker dialogs are snapshotted,
> so the GPU suite must pass **without `-update`**; and
> `BenchmarkEditorWindow_TabCycling` needs a GPU and never runs in CI, so the
> ~4,773 allocs/op figure has no automated guard — re-measure by hand when
> touching `EditorState.Clone`, `linq.SelectSlice` or the clone helpers. Cap
> sessions at ~50 messages and hand off through this file.
>
> Full handoff in `./.agent/session-carry-forward.md`.
