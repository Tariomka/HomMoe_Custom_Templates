# Session Carry-Forward — 2026-08-23 (Batch I, Phase 8 + doctrine settlement)

## 1. Session goal

Resume Batch I after the layering pivot: settle the four blocking owner
decisions, then execute **Phase 8** of
[plans/batch-i-editor-state-rework.md](../plans/batch-i-editor-state-rework.md)
— restore the dependency direction in the entity layer. **Both done. Phase 8 is
complete, green, and committed by the owner as `54191b7 "Batch I wip 5"`.**

## 2. The doctrine is now settled — read this before anything else

The previous session recorded plan **§0.4** ("a Model never leaves the backend;
`app/` must not see one"). That is **superseded**. The owner's committed
AGENTS.md §4.4 amendment (`a264d0b`) says the opposite and wins:

> `app/` **may hold a Model** as its working state. Only the **crossing** into
> `internal/` must be a DTO. The GUI maps *stored Model → request DTO* on the
> way in and *response DTO → Model and stores it* on the way back.

All four blocking decisions are answered and written up as **§0.5 of the plan**
(§0.5.1–§0.5.5), which **overrides §0.4 wherever they conflict**:

| # | Decision | Effect |
| --- | --- | --- |
| §0.5.1 | `app/` keeps the Model; DTO is the crossing | **Phase 4 is correct as committed.** Phase 11 rewritten from "move `app/` off the Model" to "pass DTOs at the handler call sites". Phase 12 **drops** its `app/ → models` prohibition. |
| §0.5.2 | Entity leaf types + helper **functions**, no wrapper structs | Wrappers rejected: entity groups hold *slices*, and Go slices don't interconvert, so wrappers would force wrap/unwrap at ~85 sites. |
| §0.5.3 | Template façade in **zero-churn** form | `internal/entities/template/types.go` is the canonical alias home; `internal/entities/types.go` re-aliases it. All **426** importers untouched. |
| §0.5.4 | Entity ⇄ Model mapping goes in `internal/mappers/` | Needs a §0.4 carve-out: permitted Entity namers are `repositories`, `models`, `entities`, `helpers/*_helpers`, `mappers`. Phase 12's checker must encode **that** list. |
| §0.5.5 | Phase 6 / frame-path cost | **Deferred — benchmark first.** Largely moot now: `app/` keeps the Model, so no per-frame DTO⇄Model conversion is introduced. §0.4.3's update flowchart overstates the cost. |

**Phase order is unchanged: 8 → 9 → 10 → 11 → 6 → 12 → 7.** Phase 5 stays
superseded by Phase 9.

## 3. Phase 8 — what shipped

Entities no longer depend on models. Verified mechanically, not by eye:
`go list -f '{{.ImportPath}} => {{.Imports}}' ./internal/entities/...` over all
**9** entity packages reports **zero** imports of `internal/models`,
`internal/dtos`, `internal/services`, `internal/handlers` or `internal/helpers`.

- **Five behaviour-free leaf types** added to `internal/entities/editor_state/`:
  `MapTopology`, `BonusPresetType`, `BonusEntry`, `ContentRuleRow`,
  `ZoneContentRow`. json tags copied **verbatim** — a renamed type must not move
  a key.
- **Behaviour became functions**, per the `road_helpers`/`zone_helpers`
  convention: `internal/helpers/config_helpers` (`GetHash`, `GetString`,
  `IsResource`) and `internal/helpers/editor_state_helpers`
  (`CloneZoneContentRow(s)`, `NormalizeZoneContentRow`, `CloneContentRuleRow(s)`).
- **Alias façades absorbed the churn.** `internal/models/config/types.go`
  re-points `MapTopology`/`BonusEntry`/`BonusPresetType` and all 22 constants at
  the entity types — **none of the 182 files importing `config` changed**.
  `internal/models/{zoneContentRow,contentRuleRow}.go` re-export the two rows, so
  the 77 referencing files only lost the `Save` suffix.
- **The rename was semantic, never a text sweep** — `vscode_renameSymbol`,
  290 edits across ~70 files (§2.6 forbids bulk in-place rewrites).
- **Six production call sites** changed behaviourally, all method → function.

## 4. File modifications

**Created**

- `internal/entities/editor_state/{mapTopology,bonusPresetType,bonusEntry,contentRuleRow,zoneContentRow}.go`
- `internal/helpers/config_helpers/{bonusEntry,bonusPresetType}.go`
- `internal/helpers/editor_state_helpers/{zoneContentRow,contentRuleRow}.go`
- `internal/models/{zoneContentRow,contentRuleRow}.go` — alias façades
- `internal/entities/template/types.go` — canonical alias home (§0.5.3)

**Deleted**

- `internal/models/config/config_inner/{mapTopology,bonusEntry,bonusPresetType}.go`
- `internal/models/{zoneContentRowSave,contentRuleRowSave}.go`
- four emptied `test/unit/internal/models/...` folders

**Edited (agent)** — `internal/entities/editor_state/{contentSettings,generationSettings}.go`
(now import **nothing**), `internal/entities/types.go`,
`internal/models/config/{types,generatorConfig}.go`,
`internal/models/editor_state_model/editorStateModel.go`,
`internal/mappers/mandatoryContentItemMapper.go`,
`internal/services/bonuses/bonusEntryService.go`,
`app/gui/dialogs/{bonusPickerDialog,zoneContentDialog}.go`,
`app/gui/panels/bonusesPanel.go`, plus the ~70 files touched by the semantic
rename, plus the plan.

**Edited (owner, in `54191b7`)** — `internal/entities/template/rmgTemplate.go`
(simplified to the new unqualified names), renamed
`common_topologies/topologies.go` → `topologyDescriptors.go`, **deleted
`app/gui/utils/models.go`** (`CloneRuleRows`, now redundant against
`editor_state_helpers`) and its test folder, and applied lint/whitespace fixes
across `content_rules`, `zone_content` and the dialogs.

> **`internal/entities/template/` is the owner's.** The agent added `types.go`
> there with explicit approval; `rmgTemplate.go` and every `template`
> subpackage remain off-limits (§2.1).

## 5. Tests

**Relocated to mirror the new impl paths (§4.6):**

```text
test/unit/internal/helpers/config_helpers/bonusEntry/getHash_test.go
test/unit/internal/helpers/config_helpers/bonusPresetType/{getString,isResource}_test.go
test/unit/internal/helpers/editor_state_helpers/zoneContentRow/{cloneZoneContentRow,normalizeZoneContentRow}_test.go
test/unit/internal/helpers/editor_state_helpers/contentRuleRow/cloneContentRuleRow_test.go
test/unit/internal/entities/editor_state/{zoneContentRow,contentRuleRow}/*Json_test.go
```

The two json round-trip files were kept and moved to the **entity** path: they
characterise the persisted format, which outweighs §4.6's "pure data structs
need no tests".

**Last full gate run — all green:**

| Gate | Result |
| --- | --- |
| `go build ./...`, `go vet -tags='integration_test,gui' ./...` | exit 0 |
| `gofmt -l .` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test ./test/unit/...` / `./test/...` / `-tags=integration_test ./test/integration/...` | pass |
| GPU suite `-tags='integration_test,gui'`, **no `-update`** | pass |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **73.7 %** (floor 72.5 %) — unchanged |

Re-verified after the owner's `54191b7`: build exit 0, unit suite pass,
coverage 73.7 %.

## 6. Git status snapshot

- **Branch:** `AD/fixing_some_stuff_08-12`
- **HEAD:** `54191b7 "Batch I wip 5"` — contains all of Phase 8 plus the owner's
  own edits. Previous commit `a264d0b "Batch I wip 4"` carried Phase 4 and the
  AGENTS.md §4.4 amendment.
- **Working tree: clean. Index: empty.** Nothing is inherited.

> **Process note — an agent mistake worth not repeating.** The agent used
> `git mv` to relocate the test files, which **stages**, violating §2.5. The
> owner absorbed it while committing. **Use `Move-Item`, never `git mv`.** The
> root cause was an opening `git status --short | Group-Object …` whose output
> was swallowed, so a dirty index was read as a clean tree — verify index state
> with `git diff --cached --name-status`, not a grouped summary.

## 7. Rejections / things declined

- **Rejected — plan §0.4's "`app/` must not see a Model".** Superseded by
  AGENTS.md §4.4 and §0.5.1.
- **Rejected — reversing Phase 4.** It was committed instead; its direction is
  now the correct one.
- **Rejected — wrapper structs** embedding the entity rows (§0.5.2), and a
  wrapper for `BonusEntry` specifically: helper functions won on the slice
  argument.
- **Rejected — the 426-file `entities.X` → `template.X` rename** inside Batch I
  (§0.5.3); the zero-churn façade was taken instead.
- **Rejected — moving `config.MapTopology`/`BonusEntry` to `internal/common/`**;
  they went to `internal/entities/editor_state/` with `config` re-aliasing.

## 8. Open questions

1. **None block Phase 9.** Its one former decision (mapping location) is
   answered: `internal/mappers/` (§0.5.4).
2. **Phase 6's benchmark still owes an answer** for §0.5.5, but only when
   Phase 6 runs.
3. **Repository-memory duplication** (`/memories/repo/conventions.md`) — flagged
   four sessions running, ~1234 lines of roughly four copies. Its Batch I
   content is now **wrong** in places (it predates §0.4 *and* §0.5) and should be
   corrected or dropped.

## 9. Next recommended actions

1. **Phase 9** — make the persisted shape an Entity. Add
   `internal/entities/editor_state/editorStateEntity.go` with `SchemaVersion` +
   the 9 groups; retype `EditorStateRepository` to read/write the **Entity** and
   return/accept a **Model**; put the mapping in
   `internal/mappers/editorStateEntityMapper.go`; delete `NewEditorStateDto` and
   `(*EditorStateDto).Model()`; make `file_service` speak Models only.
2. Then **10 → 11 → 6 → 12 → 7**.
3. Consider pulling **Phase 12's checker forward** and running it in allow-list
   mode from Phase 9 onward. Phase 8 was verified with a one-off `go list`
   pipeline; that pipeline is a ready-made basis for the real gate, and §0.4 was
   violated silently for four phases precisely because nothing checked it.

## 10. Carry-forward prompt

> Read `AGENTS.md` first, then `plans/batch-i-editor-state-rework.md` — and in
> that plan read **§0.5 before §0.4**, because §0.5 overrides §0.4 wherever they
> conflict. The layering doctrine, as finally settled: **Entity** =
> database layer, `internal/entities/`, no logic beyond json tags, may be
> embedded in a Model; **Model** = service layer, `internal/models/`, owns all
> business logic; **DTO** = consumer layer, `internal/dtos/`, no logic, the
> `app/` ↔ `internal/` contract. Dependency direction is DTO → Model → Entity,
> never the reverse. Conversion happens at exactly two seams: the handler
> (DTO ⇄ Model) and the repository (Model ⇄ Entity). **Crucially, `app/` MAY
> hold a Model** — AGENTS.md §4.4 says so explicitly; only the *crossing* into
> `internal/` must be a DTO. Do not "fix" that.
>
> **Phases 1–8 are complete and committed** through `54191b7 "Batch I wip 5"` on
> branch `AD/fixing_some_stuff_08-12`; the working tree and index are clean.
> Phase 5 is **superseded — do not execute**. Next is **Phase 9**, then
> 10 → 11 → 6 → 12 → 7. No decisions are outstanding.
>
> The hard rules, one line each: never modify `data/`,
> `internal/registry/`, or **anything under `internal/entities/template/`** —
> `rmgTemplate.go` and the `template_*` subpackages are the owner's alone;
> everything must build and run on Windows and Linux (use `path/filepath`; chain
> PowerShell with `;`, never `&&`); every change ships with tests and unit
> coverage must not drop below 72.5 % (currently 73.7 %); **never stage and never
> commit** — the owner reviews, stages and commits, so **use `Move-Item`, never
> `git mv`**, and delete with `Remove-Item`, never `git rm`; never change where
> `.rmg.json` is written and never persist the output directory; never run a bulk
> in-place rewrite over the repository — use `vscode_renameSymbol` for renames;
> cap sessions at ~50 messages and hand off through this file.
>
> Standing traps for Phase 9: `MarshalJSON`'s alias type must be **locally
> declared and must not embed the type whose method it is**, or it recurses
> forever; `UnmarshalJSON` must **merge into the existing receiver**, because the
> repository unmarshals *over* a defaults-seeded value and that is how absent
> keys keep their defaults; `omitempty` already conflates nil and empty on disk,
> so that distinction is in-memory only and must be characterised, not "fixed";
> every round-trip gate compares **parsed objects, never bytes**, because key
> order moves freely. Go 1.27 allows **promoted fields as composite-literal
> keys** and the linter's `embedlit` rule enforces that flat form, but eliding
> only the type (`Group: {…}`) does **not** compile. Run `golangci-lint-v2`
> report-only first and scope any `--fix` to the packages you actually want
> rewritten.
>
> Full handoff in `./.agent/session-carry-forward.md`.
