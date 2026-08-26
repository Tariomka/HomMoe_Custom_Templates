# Session Carry-Forward — 2026-08-26 (Batch I, Phase 10 closed)

## 1. Session goal

Finish **Phase 10** of
[.agent/plans/batch-i-editor-state-rework.md](plans/batch-i-editor-state-rework.md)
after the owner **rewrote its design**, then bring the plan document itself back
in line with the tree. **Done. All gates green; everything is committed as
`2fc6b13 "Batch I wip 8"` and the working tree is clean.**

## 2. What the owner changed, and why it matters

The agent's Phase 10 (2026-08-24) made the **Model** a thin wrapper around the
entity and the **DTO** a full structural rewrite — nine `*Dto` group structs plus
a DI'd DTO⇄Model mapper. The owner reversed it:

> *"you made the editor state model be a wrapper around the entity and the dto a
> full structural rewrite, but it should have been the other way around."*
>
> *"Redefinition is expected in Models, but it should never happen in DTOs."*

That is now **§0.7** of the plan, and it supersedes §0.6 decisions 1, 2, 3 and 7.
Decisions 4, 5 and 6 survived.

**The shape that exists today:**

- `EditorStateDto` is `struct { editor_state_model.EditorState }` — nothing else.
  **A DTO embedding a Model is intended here.** Do not "fix" it.
- Each model group wraps its entity group
  (`type CastleSettings struct { editor_state.CastleSettings }`), as do
  `BonusEntry`, `ContentRuleRow`, `ManualZoneSave`, `ManualConnectionSave`.
- **`ContentSettings` and `ManualEditSettings` cannot wrap** — their slices carry
  model element types and Go slices do not interconvert — so they are re-declared
  and carry `ToXModel` / `ToXEntity` converters.
- `ZoneContentRow` wraps the entity but **shadows `Rules`** with
  `[]ContentRuleRow`; the embedded slice is deliberately left nil.
- Converters name the layer explicitly (`ToXModel` / `ToXEntity`) because
  `ToManualZoneSaves` / `FromManualZoneSaves` already exist for a **different
  axis** — live `[]entities.Zone` ⇄ save. Confusing the two silently produces the
  wrong shape.

## 3. What shipped this session

- `ContentSettings` / `ManualEditSettings` converters, plus the element-level
  pairs they need (`ToBonusEntry*`, `ToZoneContentRow*`, `ToContentRuleRow*`,
  `ToManualZoneSave*`, `ToManualConnectionSave*`).
- `ContentRuleRow` / `ZoneContentRow` turned into real entity wrappers, with
  `Clone()` / `Normalized()` as methods delegating to the surviving
  `editor_state_helpers` entity functions.
- The entity mapper's two commented-out groups restored.
- **Cycle broken:** `defaultPlayerZoneContentRows.go` and its test folder moved
  into `editor_state_model`; `internal/common/common_zone_contents` is gone.
- `internal/models/types.go` and `internal/models/editor_state_model/types.go`
  deleted; **56 files** repointed onto `editor_state_model.ZoneContentRow` /
  `.ContentRuleRow`.
- Tests updated throughout; new per-file test folders for every new converter.

## 4. Two real bugs found in the staged code — do not reintroduce

**nil is load-bearing on the regeneration path.** `GetPreviousState()` /
`GetNextState()` return nil when absent, and the decision service reads a nil
`Previous` as "first generation" and a nil `Next` as "unarmed debounce".

- `drivers.getPreviousStateDto` / `getNextStateDto` dereferenced them
  unconditionally → panic.
- `regenerationHandler.DecideRegeneration` / `DecideManualEditReapplication` did
  the same.
- Worse than the panic: `new(this.getPreviousStateDto())` would have destroyed
  the signal anyway, because it always yields a **non-nil** pointer.

Both now go through nil-preserving helpers. **Any future refactor of these four
call sites must preserve nil.**

## 5. Plan document updates made this session

- **New §0.7** records the owner's rework and marks §0.6 decisions 1/2/3/7 void.
- **"For Future Agents"** now says to read the doctrine newest-first
  (§0.7 → §0.6 → §0.5 → §0.4) with a one-line summary of each.
- **Phase 10** — Work list and Phase Summary rewritten to describe what actually
  shipped; the benchmark table carries a warning that it measures the *deleted*
  design.
- **Phase 6** — marked next; both stale baselines flagged, with an instruction to
  re-measure from `2fc6b13`.
- **Phase 12** — refreshed: verified allow-list contents, an explicit
  "do not flag Model-in-DTO" rule, and a new decision item about the now-unused
  depguard permission.
- **Phase 7** — paths corrected, the already-done shim bullet removed, and two
  new cleanup items (AGENTS.md's dangling links, the `godox` TODO).
- **§0 / §0.5.2** — historical banners added where later rounds overtook them
  (wrapper structs, the row-type alias façade, `MapTopology`'s location).
- **69 relative links repaired.** Moving the plan to `.agent/plans/` made it two
  levels deep, so every `](../…)` became `](../../…)` and the backlog links
  became `](../backlog/…)`.

Six link targets remain dead **on purpose** — they sit in Phase 1–3 prose and name
files that genuinely no longer exist (`internal/dtos/editorStateDto.go`,
`internal/models/zoneContentRowSave.go`, the old `editorStateDto` test folder).
Rewriting them would falsify the historical record.

Also fixed outside the plan: **AGENTS.md §4.6**'s two `todo/test_observations.md`
links now point at `.agent/backlog/test_observations.md`.

## 6. Gates

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
| GPU suite, **no `-update`** | exit 0 (23.6 s) |
| `golangci-lint-v2 run ./...` | **1 issue** — see below |
| Unit coverage | **73.9 %** (floor 72.5 %) |

The single lint issue is a **`godox` TODO the owner left** in
[internal/handlers/stateHandler.go](../internal/handlers/stateHandler.go) line 29
(*"should just return EditorStateVali…"*). Ask the owner if he wants to resolve it himself
or I should do it. **Treat 1 godox issue as the current lint baseline**, not 0,
until it is cleared — Phase 7 has a bullet for it.

The frozen fixtures and the untagged `editorStateWireFormat_integration_test.go`
passed **unchanged** throughout.

## 7. Git status snapshot

- **Branch:** `AD/fixing_some_stuff_08-12`
- **HEAD:** `2fc6b13 "Batch I wip 8"` (after `7aef72f "Batch I wip 7"`)
- **Working tree clean, index empty** at the start of the doc-update work. The
  only uncommitted changes now are this file, the plan and AGENTS.md.
- Repo cleanup by the owner: plans live in `.agent/plans/`, backlog and review
  docs in `.agent/backlog/`, the review prompt in `.agent/promt_templates/`. The
  old `plans/` and `todo/` directories are gone.

## 8. Rejections / things not done

- **Rejected — the nine-group DTO with struct conversion.** It was built, then
  deleted by the owner. Do not propose it again.
- **Not done — the three remaining DTO-below-handler services** (`bonuses`,
  `pickers`, `zone_content`). Phase 10 said list, not fix; they are recorded for
  Phase 12's allow-list.
- **Not done — reverting the depguard relaxation.** `internal/mappers` is still
  permitted from `app/` but nothing uses it. Phase 12 must decide.

## 9. Open questions

1. **None block Phase 6.**
2. Phase 12 owes a decision on the unused `internal/mappers` permission in
   depguard + the architecture test.
3. **Repo memory duplication** (`/memories/repo/conventions.md`) — flagged seven
   sessions running: ~1,300 lines, roughly four copies of the same body. Still
   needs a dedupe pass.
4. **The owner's `godox` TODO.**

## 10. Next recommended actions

1. **Phase 6** — stop the per-frame whole-state clone (backlog §1.5). Re-measure
   first; every benchmark number in the plan predates the current tree.
2. Then **12 → 7**.
3. At Phase 7, remember the plan file is meant to be **deleted** once the backlog
   entries are self-contained (doc-lifecycle rule).

## 11. Carry-forward prompt

> Read `AGENTS.md` first, then `.agent/plans/batch-i-editor-state-rework.md` — and
> in that plan read the doctrine **newest-first: §0.7, then §0.6, then §0.5, then
> §0.4**, because each supersedes the ones after it. The settled shape:
> **Entity** = `internal/entities/`, database layer, json tags and
> (de)serialisation only; **Model** = `internal/models/`, service layer, **owns
> the structure and all business logic**; **DTO** = `internal/dtos/`, the
> `app/` ↔ `internal/` crossing. The owner's rule is *"redefinition is expected in
> Models, but it should never happen in DTOs"* — so `EditorStateDto` is literally
> `struct { editor_state_model.EditorState }`, and **a DTO embedding a Model is
> intended**. `app/` MAY hold a Model; only the crossing must be a DTO. Do not
> "fix" either.
>
> **Phases 1–4 and 8–10 are complete and committed** through `2fc6b13` on branch
> `AD/fixing_some_stuff_08-12`; the working tree is clean. Phases **5 and 11 are
> superseded — do not execute**. Next is **Phase 6**, then 12 → 7. No decisions
> are outstanding.
>
> The hard rules, one line each: never modify `data/`, `internal/registry/`, or
> **anything under `internal/entities/template/`**; everything must build and run
> on Windows and Linux (use `path/filepath`; chain PowerShell with `;`, never
> `&&`); every change ships with tests and unit coverage must not drop below
> 72.5 % (currently 73.9 %); **never stage and never commit** — the owner
> reviews, stages and commits, so **use `Move-Item`, never `git mv`**, and delete
> with `Remove-Item`, never `git rm`; never change where `.rmg.json` is written
> and never persist the output directory; never run a bulk in-place rewrite over
> the repository, and never round-trip a `.go` file through
> `Get-Content`/`Set-Content` (it joins every line and corrupts the file — that
> has now happened once); cap sessions at ~50 messages and hand off through this
> file.
>
> Standing traps for Phase 6: the hot path is
> `app/gui/models.EditorState.UpdateCurrentState`, which deep-`Clone()`s and then
> calls `ValidateEditorState` — that crossing was **deliberately left on the
> Model** so this phase could measure it without a conversion in the way; **every
> benchmark number in the plan is from a superseded tree, so re-measure from
> `2fc6b13`** with `BenchmarkEditorWindow_TabCycling -benchtime=50x -count=6`
> against a detached `git worktree` and read the steady-state samples; the panel
> line numbers in Phase 6's bullet list drifted in the Phase 10 sweep, so find
> those call sites by name; **nil is load-bearing** on the regeneration path
> (nil `Previous` = first generation, nil `Next` = unarmed debounce) and two
> dereference bugs there were fixed this session — do not reintroduce them; the
> lint baseline is **1 `godox` issue**, not 0, until the owner's TODO in
> `stateHandler.go` is cleared; and the two frozen fixtures under
> `test/test_helpers/testdata/` plus the untagged
> `editorStateWireFormat_integration_test.go` must keep passing **unchanged**,
> comparing **parsed objects, never bytes**. Run `golangci-lint-v2` report-only
> first and scope any `--fix` to the packages you actually want rewritten.
>
> Full handoff in `./.agent/session-carry-forward.md`.
