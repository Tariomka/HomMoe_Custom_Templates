# Session Carry-Forward — 2026-08-22 (Batch I, Phase 4 + layering pivot)

## 1. Session goal

Execute **Phase 4** of [plans/batch-i-editor-state-rework.md](../plans/batch-i-editor-state-rework.md)
— swap the runtime editor-state type from `EditorStateDto` to `EditorState`.
Done and green. **Then the owner's review found the whole batch was built on a
wrong premise**, and the rest of the session went into documenting the pivot.
No code was changed after the correction.

## 2. The correction — read this before anything else

**`EditorStateDto` is not, and was never meant to be, the persisted `.gen.json`
shape. That is an Entity's job.** The owner's ruling, now recorded as
**[§0.4 of the plan](../plans/batch-i-editor-state-rework.md)**:

| Layer | Type | Declared in | May be used by | Logic |
| --- | --- | --- | --- | --- |
| **Database** | Entity | `internal/entities/` | `internal/repositories/` only | none beyond (de)serialisation |
| **Service** | Model | `internal/models/` | services, validators, mappers, repositories | **all business logic** |
| **Consumer** | DTO | `internal/dtos/` | `internal/handlers/`, and in `app/` only at handler call sites | none |

- An Entity **may be embedded in or held as a field of a Model** — then it is the
  *Model* being used, which is allowed. An Entity is never passed around alone
  above the repository.
- A Model **never leaves the backend**. `app/` must not see one.
- A DTO is the frontend↔backend contract and crosses only the handler boundary.
- **Dependency direction is DTO → Model → Entity, never the reverse.**

**Request/response flow.** The UI creates a request DTO; conversion happens at
exactly two seams — the handler (DTO ⇄ Model) and the repository (Model ⇄ Entity):

```
save:  app/gui --DTO--> handler --Model--> service --Model--> repository --Entity--> disk
load:  disk --Entity--> repository --Model--> service --Model--> handler --DTO--> app/gui
```

Standard C#/Java enterprise layering. References the owner gave:
[Baeldung — Entities vs DTOs](https://www.baeldung.com/java-entity-vs-dto),
[DevsDaily — Entity, Model, ViewModel and DTO in C#](https://devsdaily.com/understanding-entity-model-viewmodel-and-dto-in-csharp/).

### What is wrong in the tree right now

1. **The persisted shape is a DTO.** `EditorStateRepository` is typed
   `IFileRepository[editor_state_dto.EditorStateDto]`, so the storage boundary
   speaks the consumer type.
2. **Phase 4 pushed `app/` onto the Model — backwards.** The GUI must hold the
   DTO. Phase 4's grep gate encodes the inverted rule.
3. **The dependency arrow is reversed in the entity layer.** Verified this
   session: [contentSettings.go](../internal/entities/editor_state/contentSettings.go)
   imports `internal/models` **and** `internal/models/config`;
   [generationSettings.go](../internal/entities/editor_state/generationSettings.go)
   imports `internal/models/config`. `models.ZoneContentRowSave` /
   `ContentRuleRowSave` have `Clone()` and `Normalized()`, so they are Models
   being used as Entity fields.
4. **AGENTS.md §4.4 contradicts the doctrine** — it permits `app/` to use
   `internal/models/` and `internal/entities/` for "data typing".
5. **The breach is repository-wide.** **36 files under `app/`** import
   `internal/models` or `internal/entities`. Phase 4 added 8; **28 predate this
   batch** (`config`, `neutral_zone`, `preview`, `entities.Zone`/`Connection`).
   That residue is explicitly out of scope for Batch I — Phase 12 records it as
   its own backlog item.

## 3. Plan changes made this session (documentation only)

[plans/batch-i-editor-state-rework.md](../plans/batch-i-editor-state-rework.md),
**added, not rewritten**, as the owner asked:

- **§0.4 Layering doctrine** + **§0.4.1 request/response flow** (with a mermaid
  sequence diagram) + **§0.4.2 what this batch got wrong**.
- A **⚠ banner in "For Future Agents"** so nobody reads Phases 1–5 without the
  correction.
- **Phase 5 marked `Superseded — do not execute`**, with a blockquote explaining
  that every item is still wanted but belongs on the Entity. Its body was left
  intact so the reasoning is not lost.
- **Phase 7** marked as running last.
- A new **"Correction phases"** part: a *Decision required* block plus
  **Phases 8–12**.
- **§0.4.3 Editor-state flows in the target state** — three mermaid flowcharts
  (save, load, update) showing the post-Phase-12 shape with layer subgraphs, the
  type on every arrow, and the request/response paths distinguished. It doubles
  as the review checklist: exactly two conversion seams (handlers DTO ⇄ Model,
  repositories Model ⇄ Entity); `app/` never names a Model or Entity;
  `repositories` never names a DTO; `services`/`validators` never name a DTO;
  the Entity never leaves the repository.
- **Revised phase order: 8 → 9 → 10 → 11 → 6 → 12 → 7.** Phase 6 must follow
  Phase 11 or its view structs are invalidated.

The five new phases in one line each:

| Phase | Purpose |
| --- | --- |
| **8** | Restore the dependency direction in the entity layer — entities import no model package. |
| **9** | Make the persisted shape an **Entity** (`EditorStateEntity` + `SchemaVersion` + migration); repository maps Entity ⇄ Model. Carries superseded Phase 5 forward. |
| **10** | Rebuild `EditorStateDto` as a flat, logic-free consumer contract; handlers map DTO ⇄ Model via `internal/mappers/`. |
| **11** | Move `app/` off the Model — the GUI holds DTOs; state logic goes behind handler calls. Reverses Phase 4's direction. |
| **12** | Enforce the layering mechanically in `cmd/testlayoutcheck`, amend AGENTS.md §4.4, log the 28-file residue. |

## 4. Decisions the owner still owes

These block the phases they sit in. **Do not start those phases without answers.**

1. **The unstaged Phase 4 diff — keep or revert?** (Plan: *Decision required
   before Phase 8*.) Option A keeps it and lets Phases 10–11 rewrite ~80 of the
   same files twice. Option B reverts to `0ca3b6d` and starts at Phase 8.
   **Recommendation: B** — at `0ca3b6d` the GUI and handlers already name
   `EditorStateDto` at exactly the call sites §0.4 wants; only its *shape* is
   wrong there, and Phase 10 fixes shape in one package. If reverting, **save
   [editorStateRoundTrip_integration_test.go](../test/integration/editorStateRoundTrip_integration_test.go)
   and the two new DTO unit tests first** — they are still wanted.
   **Only the owner may revert; the agent must not `git restore` this tree.**
2. **Phase 8** — `config.MapTopology` / `config.BonusEntry` live in
   `internal/models/config`, which the generator uses heavily. Move the enums to
   `internal/common/` (clean, touches the generator) or declare entity-local
   counterparts and convert at the repository seam?
3. **Phase 9** — where does Entity ⇄ Model mapping live? AGENTS.md §4.4 says
   `internal/mappers/`; §0.4 says entities are used only in repositories.
   Recommendation: keep it private to `internal/repositories/`.
4. **Phase 11** — `Clone` and `EqualsIgnoringManualEdits` are on the frame path.
   Route them through handlers literally (a DTO⇄Model conversion per frame,
   fighting Phase 6), or let `app/` keep a local DTO copy helper as view-layer
   plumbing and route only the four comparison methods? The update flowchart in
   **§0.4.3** is drawn against the second option and marks what changes under the
   first; `IStateComparisonHandler` in that diagram is an **invented name** —
   no such interface exists yet. Decide this **with Phase 6's benchmark in view**,
   not separately.

## 5. Phase 4 — what actually shipped (all green, all unstaged)

101 files, **+663/−665**; net-negative because the four Phase 3 shims and their
five delegation tests went away. Its *direction* is superseded by Phase 11, but
it is complete and passing, so it is a coherent thing to review or revert as a
unit.

- `NewEditorStateDto(model)` and `(*EditorStateDto).Model()` added; `file_service`
  became the model↔DTO seam. **Phase 9 deletes both** — they exist only to make
  the DTO a persistence shell.
- `app/gui`, handlers + `handler_interfaces`, validators, mappers,
  `services/editor` and the carrier DTOs all moved onto the model.
- The owner renamed `EditorStateModel` → **`EditorState`** mid-review; the
  factory is still `NewDefaultEditorStateModel`.
- **New tests:** [editorStateRoundTrip_integration_test.go](../test/integration/editorStateRoundTrip_integration_test.go)
  (untagged — production APIs only, no GPU) proving non-aliasing on load, on
  save, and a full round-trip; plus `newEditorStateDto_test.go` and
  `model_test.go`.
- **Deleted:** the four shim delegation tests.
- **Naming:** production parameters named `stateDto` that now hold a model were
  renamed; ~55 test-local `dto`/`stateDto` identifiers were deliberately left
  (renaming them would shadow the driver `state` in most closures) — recorded in
  [todo/test_observations.md](../todo/test_observations.md).

| Gate | Result |
| --- | --- |
| `go build ./...`, `go vet ./...`, `go vet -tags='integration_test,gui' ./...` | exit 0 |
| `gofmt -l .` | empty |
| `go run ./cmd/testlayoutcheck .` | passed |
| unit / untagged / tagged-integration suites | pass |
| GPU suite (`-tags='integration_test,gui'`) | pass, **no `-update`** |
| Unit coverage | **73.7 %** (floor 72.5 %) |
| `golangci-lint-v2 run ./...` | **0 issues** |

## 6. Git status snapshot

- **Branch:** `AD/fixing_some_stuff_08-12`
- **HEAD:** `0ca3b6d "Batch I wip 3"` — unchanged. Phases 1–3 are committed;
  Phase 4 and this session's plan edits are **entirely unstaged**.
- **`git status --short`:** ~105 entries — 101 modified, 4 deleted (shim tests),
  3 untracked (new test files), plus `plans/…` and `todo/test_observations.md`.
- The one **staged** entry, `.agent/session-carry-forward.md`, was already staged
  when the session began and **was not touched** — only its working-tree content
  was rewritten.
- `coverage.txt` / `coverage.html` / `lcov.info` regenerate identically.

**Nothing was staged or committed by the agent.**

## 7. Rejections / things the user declined

- **Rejected — the entire premise of Phases 1–5:** that `EditorStateDto` is the
  file format. Superseded by §0.4.
- **Rejected — Phase 4's direction:** `app/` must hold DTOs, not Models.
- **Declined earlier in the session:** moving `EditorStateSaveDto` /
  `EditorStateValidationDto` out of `internal/dtos/editor_state_dto/`. They are
  DTOs for `EditorState` *operations*. Still stands, and is now reinforced by
  §0.4 — they are consumer-layer types and belong in `internal/dtos/`.
- **No code changes were made after the correction**, at the owner's instruction
  ("I'm still reviewing the damage").

## 8. Open questions

1. **The four decisions in §4 above.** All block their phases.
2. **Repository-memory duplication** (`/memories/repo/conventions.md`) — flagged
   three sessions running, ~1234 lines of roughly four copies of one body. Its
   Batch I content is now partly **wrong** and should be corrected or dropped
   during the dedupe.

## 9. Next recommended actions

1. Answer decision §4.1 (keep or revert Phase 4). Nothing else can start cleanly.
2. Answer §4.2, then run **Phase 8** — entities must import no model package.
   It is the bottom of the stack; nothing above can be laid out correctly first.
3. **Phase 9**, **10**, **11** in order, then **6**, **12**, **7**.
4. Consider pulling **Phase 12's checker forward** and running it in
   allow-list mode from Phase 8 onward, so each phase's progress is measured
   rather than asserted. §0.4 was violated silently for four phases precisely
   because nothing checked it.

## 10. Carry-forward prompt

> Read `AGENTS.md` first, then `plans/batch-i-editor-state-rework.md` — and in
> that plan read **§0.4 before anything else**. Batch I was built on a wrong
> premise: `EditorStateDto` was treated as the persisted `.gen.json` shape, but
> that is an **Entity's** job. The owner's layering doctrine, now §0.4: **Entity
> = database layer**, used only in `internal/repositories/`, no logic beyond
> JSON (de)serialisation, and may be embedded in a Model; **Model = service
> layer**, owns all business logic and **never leaves the backend**; **DTO =
> consumer layer**, the `app/gui/` ↔ `internal/` contract, no logic, usable only
> in `internal/handlers/` and at handler call sites in `app/`. Dependency
> direction is DTO → Model → Entity, never the reverse. Conversion happens at
> exactly two seams: the handler (DTO ⇄ Model) and the repository (Model ⇄
> Entity).
>
> Phases 1–3 are committed as `0ca3b6d "Batch I wip 3"` on branch
> `AD/fixing_some_stuff_08-12`. **Phase 4 is complete and green but unstaged, and
> its direction is backwards** — it moved `app/` onto the Model. Phase 5 is
> **superseded, do not execute**. The corrective work is **Phases 8–12**, and the
> revised order is **8 → 9 → 10 → 11 → 6 → 12 → 7**.
>
> **Four owner decisions are outstanding and block their phases** — see §4 of
> `./.agent/session-carry-forward.md` and the *Decision required before Phase 8*
> block in the plan. The first is whether to keep or revert the unstaged Phase 4
> diff; the recommendation is revert, but **only the owner may do that** — never
> `git restore` a tree you did not create.
>
> The hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/` without explicit
> approval; everything must build and run on Windows and Linux (use
> `path/filepath`; chain PowerShell with `;`, never `&&`); every change ships
> with tests and unit coverage must not drop below 72.5 % (currently 73.7 %);
> durable multi-session work gets a plan file under `plans/`; **never stage and
> never commit** — the owner reviews, stages and commits, so leave the staging
> area alone entirely, and delete with `Remove-Item`, never `git rm`; never
> change where `.rmg.json` is written and never persist the output directory;
> never run a bulk in-place rewrite over the repository; cap sessions at ~50
> messages and hand off through this file.
>
> Standing traps that survive the pivot: `MarshalJSON`'s alias type must be
> **locally declared and must not embed the type whose method it is**, or it
> recurses forever; `UnmarshalJSON` must **merge into the existing receiver**,
> because the repository unmarshals *over* a defaults-seeded value and that is
> how absent keys keep their defaults; `omitempty` already conflates nil and
> empty on disk, so that distinction is in-memory only and must be
> characterised, not "fixed"; every round-trip gate compares **parsed objects,
> never bytes**, because key order moves freely. Go 1.27 allows **promoted
> fields as composite-literal keys** and the linter's `embedlit` rule enforces
> that flat form, but eliding only the type (`Group: {…}`) does **not** compile.
> Run `golangci-lint-v2` report-only first and scope any `--fix` to the packages
> you actually want rewritten.
>
> Full handoff in `./.agent/session-carry-forward.md`.
