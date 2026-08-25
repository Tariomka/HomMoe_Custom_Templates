# Session Carry-Forward — 2026-08-24 (Batch I, Phase 10)

## 1. Session goal

Execute **Phase 10** of
[plans/batch-i-editor-state-rework.md](../plans/batch-i-editor-state-rework.md)
— rebuild `EditorStateDto` as a logic-free consumer contract and make the
editor-state handler crossings speak it. **Done, all gates green, unstaged and
uncommitted.** Phase 11 was folded in and is now marked superseded.

## 2. Read this before touching anything

Plan **§0.5 overrides §0.4**, and the new **§0.6** (written this session)
settles Phase 10. The doctrine is unchanged:

- **Entity** = `internal/entities/`, DB layer, no logic beyond (de)serialisation.
- **Model** = `internal/models/`, service layer, owns the business logic.
- **DTO** = `internal/dtos/`, the `app/` ↔ `internal/` crossing, no logic.
- Direction **DTO → Model → Entity**, never the reverse. Two conversion seams:
  the handler (DTO ⇄ Model) and the repository/FileService (Model ⇄ Entity).
- **`app/` MAY hold a Model.** Only the *crossing* must be a DTO.

**Phase order is now 8 → 9 → 10 → 6 → 12 → 7.** Phases 5 and 11 are superseded.
**Next is Phase 6** (stop the per-frame whole-state clone, backlog §1.5).

## 3. Owner decisions taken this session (§0.6 of the plan)

1. **DTO shape — nine group DTOs, converted as structs.** Not 72 flat fields.
2. **Leaf types reused through the §0.5.2 alias façade**, never copied.
3. **`app/` may import `internal/mappers`** — depguard relaxed, AGENTS.md §4.4
   amended.
4. **Phase 11 folded into Phase 10.**
5. **Regeneration envelopes split in two** — DTO crosses, model type goes to the
   service.
6. **Every crossing swaps except `ValidateEditorState`** (per-frame path).
7. **Full DI for the new mapper**, threaded from the composition root.

## 4. What shipped

### The trick worth remembering

`EditorStateDto` embeds nine mirror structs. **Go ignores struct tags when
deciding conversion identity**, so the mapper is
`editor_state_dto.MapSettingsDto(state.MapSettings)` — nine conversions per
direction — and **a field added to an entity group but not mirrored in its DTO
does not compile.** That is a compile-time drift guard; it is why the DTO is
grouped rather than flat, and it is why no reflection drift test was needed.
Field **order** matters as much as names and types.

### Created

- `internal/dtos/editor_state_dto/` — 10 new files: the nine group DTOs plus
  `castleSettingChangesDto.go`.
- `internal/mappers/editorStateMapper.go` + `…MapperInterface.go` —
  `ToDto`/`ToModel`, nil-safe `ToDtoPointer`/`ToModelPointer`, and the
  `CastleSettingChanges` pair.
- `internal/models/regeneration/` — `DecisionRequest`, `Decision`,
  `ManualEditDecision`, `NextStateAction`.
- `test/unit/internal/mappers/editorStateMapper/` — new coverage.
- `common_test.go` in the four `test/unit/app/gui/drivers/*` folders (a shared
  `newEditorStateMapper()` helper).

### Deleted

- `test/unit/internal/dtos/editor_state_dto/` — `NewDefaultEditorStateDto` and
  the four delegation shims are gone, so the tests had no subject. A
  behaviour-free struct needs none per §4.6.

### Changed shape

- Eight crossings now take/return `EditorStateDto`: `LoadState`, `SaveState`,
  `GenerateTemplate`, `UpdateTemplate`, `ReapplyCastleSettings`,
  `DecideRegeneration`, `DecideManualEditReapplication`, `GetZoneEditorOptions`.
- `NewStateHandler`, `NewTemplateHandler` (3rd arg), `NewZoneEditorHandler`
  (2nd arg), `NewRegenerationHandler` all gained the mapper.
- `drivers.NewUIState(handler, fileSystem, regeneration, mapper, findTemplateDir)`
  and `editor.NewWindow(handler, fileSystem, regeneration, mapper)`.
- `drivers.State` gained `GetStateDto()`.
- `internal/services/editor` **no longer imports `internal/dtos` at all**.

## 5. Three things that will bite the next agent

1. **`internal/models` cannot import `editor_state_model`.** There is a real
   cycle: `internal/models` ← `common_zone_contents` ← `editor_state_model`. The
   regeneration types were planned for `internal/models` and had to move to
   `internal/models/regeneration/`. `dtos.NextStateAction` is now an **alias** of
   `regeneration.NextStateAction`, which is why every `dtos.NextStateClear` in
   `app/` still compiles.
2. **`ValidateEditorState` is deliberately still on the Model.** It is called
   from `app/gui/models.EditorState.UpdateCurrentState` on every panel write.
   `EditorStateValidationDto` therefore still carries a Model. Do not "fix" it
   without a measurement.
3. **I corrupted a file with PowerShell once this session** —
   `Get-Content x | … | Set-Content -NoNewline` joins every line into one
   (`package editorimport …`). Recovered with `git restore` on that single file.
   The repo memory already warns about this; it is real. Use the edit tools.

## 6. Tests

Suites all green. The **frozen fixtures and the untagged
`editorStateWireFormat_integration_test.go` passed unchanged**, as expected —
the DTO is off the persistence path entirely.

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | exit 0 |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 |
| `go test ./test/unit/...` | exit 0 |
| `go test ./test/...` (untagged) | exit 0 |
| `go test -tags=integration_test ./test/integration/...` | exit 0 |
| GPU suite, **no `-update`** | exit 0 (40.3 s) |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **73.7 %** — unchanged, floor 72.5 % |

**Benchmark** (`BenchmarkEditorWindow_TabCycling`, `-benchtime=50x -count=6`,
baseline from a detached worktree at `2106c70`): **+24 allocs/op (9,520 →
9,544, +0.25 %)**, +2.1 % B/op, ns/op inside the noise. The crossing is free
because `ValidateEditorState` was excluded from it.

## 7. Git status snapshot

- **Branch:** `AD/fixing_some_stuff_08-12`
- **HEAD:** `2106c70 "PR review"` — Phase 9 is already committed, and the owner
  reworked it after the previous handoff was written.
- **Index empty. Nothing staged, nothing committed.** ~100 paths: 82 modified,
  1 deleted, 17 new. No `git mv`, no `git rm`.

## 8. Rejections / things not done

- **Rejected — 72 flat DTO fields.** It is the "can silently drop one" shape the
  owner rejected for the entity mapper in Phase 9; the grouped form gets a
  compile-time guard instead.
- **Rejected — embedding the entity groups in the DTO.** One-line mapper, but
  `internal/dtos` would name an Entity, which §0.5.4's namer list and Phase 12's
  checker forbid.
- **Not done — the remaining DTO-below-handler breaches.** Phase 10 said list,
  not fix. Three services remain (below).
- **Not done — deleting `EditorStateValidationDto`'s Model.** Deliberate.

## 9. Open questions

1. **None block Phase 6.**
2. §0.5.5 still owes an answer on the Phase 6 frame path; the benchmark above is
   a useful new datum (the crossing costs 0.25 % allocs).
3. **Repo memory duplication** (`/memories/repo/conventions.md`) — now flagged
   six sessions running: ~1,280 lines, roughly four copies of the same body.
   Still needs a real dedupe pass.

## 10. Next recommended actions

1. **Review and commit Phase 10.**
2. **Phase 6** — stop the per-frame whole-state clone (backlog §1.5).
3. Then **12 → 7**. Phase 12's checker must encode the §0.5.4 Entity-namer list
   (`repositories`, `models`, `entities`, `helpers/*_helpers`, `mappers`) and the
   DTO rule, seeded with an allow-list holding the three services below.
4. **Residual audit for Phase 12's allow-list** — DTOs still used below the
   handler layer:

   | Package | DTOs it names |
   | --- | --- |
   | `internal/services/bonuses` | `BonusCompositionRequestDto`, `BonusCompositionResultDto`, `ExistingBonusesDto` |
   | `internal/services/pickers` | `PickerItemDto`, `PickerSpellDto`, `PickerEntryDto`, `PickerRowDto` |
   | `internal/services/zone_content` | `ContentRuleCompositionRequestDto`, `ContentRuleCompositionResultDto` |

   Each wants the treatment the regeneration service just got: a model-side
   request/result pair, mapped at the handler.

## 11. Carry-forward prompt

> Read `AGENTS.md` first, then `plans/batch-i-editor-state-rework.md` — and in
> that plan read **§0.6, then §0.5, then §0.4**, because each overrides the one
> after it. The layering doctrine, as settled: **Entity** = database layer,
> `internal/entities/`, no logic beyond json tags and (de)serialisation;
> **Model** = service layer, `internal/models/`, owns all business logic;
> **DTO** = consumer layer, `internal/dtos/`, no logic, the `app/` ↔ `internal/`
> contract. Dependency direction is DTO → Model → Entity, never the reverse.
> Conversion happens at exactly two seams: the handler (DTO ⇄ Model) and the
> repository/FileService (Model ⇄ Entity). **`app/` MAY hold a Model** — AGENTS.md
> §4.4 says so explicitly, and as of this session `app/` may also import
> `internal/mappers`, because the consumer has to build the crossing DTO itself.
> Do not "fix" either.
>
> **Phases 1–4 and 8–10 are complete.** Phases 1–9 are committed through
> `2106c70` on branch `AD/fixing_some_stuff_08-12`; **Phase 10 is finished and
> green but sits unstaged in the working tree** — review and commit it before
> starting new work. Phases **5 and 11 are superseded — do not execute**. Next is
> **Phase 6**, then 12 → 7. No decisions are outstanding.
>
> The hard rules, one line each: never modify `data/`, `internal/registry/`, or
> **anything under `internal/entities/template/`**; everything must build and run
> on Windows and Linux (use `path/filepath`; chain PowerShell with `;`, never
> `&&`); every change ships with tests and unit coverage must not drop below
> 72.5 % (currently 73.7 %); **never stage and never commit** — the owner
> reviews, stages and commits, so **use `Move-Item`, never `git mv`**, and delete
> with `Remove-Item`, never `git rm`; never change where `.rmg.json` is written
> and never persist the output directory; never run a bulk in-place rewrite over
> the repository, and never round-trip a `.go` file through
> `Get-Content`/`Set-Content` (it joins every line and corrupts the file — it
> happened this session); cap sessions at ~50 messages and hand off through this
> file.
>
> Standing traps for Phase 6: the per-frame path is
> `app/gui/models.EditorState.UpdateCurrentState`, which deep-`Clone()`s and then
> calls `ValidateEditorState` — that crossing was **deliberately left on the
> Model** in Phase 10 precisely so Phase 6 could measure it without a DTO
> conversion in the way; `internal/models` **cannot** import
> `editor_state_model` (real cycle via `common_zone_contents`), which is why the
> regeneration types live in `internal/models/regeneration/`; measure with
> `BenchmarkEditorWindow_TabCycling -benchtime=50x -count=6` against a detached
> `git worktree` baseline and read the steady-state samples, and the current
> figure to beat is **9,544 allocs/op, 927 KB/op**; the two frozen fixtures under
> `test/test_helpers/testdata/` and the untagged
> `editorStateWireFormat_integration_test.go` must keep passing **unchanged**,
> and they compare **parsed objects, never bytes**. Run `golangci-lint-v2`
> report-only first and scope any `--fix` to the packages you actually want
> rewritten.
>
> Full handoff in `./.agent/session-carry-forward.md`.
