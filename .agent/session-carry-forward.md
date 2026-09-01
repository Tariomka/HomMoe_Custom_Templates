# Session Carry-Forward — 2026-09-01 (Batch I closed)

## 1. Session goal

Finish **Batch I** — the `EditorStateDto` rework (backlog §2.1, folding in
§1.5). Phases 6, 12 and 7 ran this session, closing the batch. **Done. Every
gate green. The plan file `.agent/plans/batch-i-editor-state-rework.md` has been
deleted**, per its own doc-lifecycle rule, now that the backlog entries are
self-contained.

> Recover it if you need the phase-by-phase record:
> `git show 938ef55:.agent/plans/batch-i-editor-state-rework.md`

## 2. Where the doctrine lives now

**Read `AGENTS.md` §4.4.1 first.** The Entity/Model/DTO layering is no longer
described in a plan file — it is in the instructions, and it is enforced by
[test/unit/architecture/dependency/layering_test.go](../test/unit/architecture/dependency/layering_test.go).

The three things most often gotten wrong:

- **The Model owns the structure.** *"Redefinition is expected in Models, but it
  should never happen in DTOs."* `EditorStateDto` is literally
  `struct { editor_state_model.EditorState }`. **A DTO embedding or carrying a
  Model is intended** — `EditorStateValidationDto.State`,
  `CastleSettingsReapplyRequestDto.Changes` and
  `ManualEditDecisionDto.ReapplyWithCastleChanges` all do it deliberately. Do not
  "fix" them.
- **`app/` may hold a Model**; only the *crossing* into `internal/` is a DTO.
  `app/` → `internal/models` is fine. `app/` → `internal/mappers`,
  `internal/services`, `internal/repositories`, `internal/validators` is not.
- **Conversion happens at exactly two seams:** `internal/handlers` (DTO ⇄ Model)
  and `internal/repositories` (Model ⇄ Entity).

The full batch history is backlog [§2.1](backlog/backlog-opus5.md); the
render-path work is [§1.5](backlog/backlog-opus5.md); the residual layering
breach is [§2.6](backlog/backlog-opus5.md).

## 3. What shipped this session

### Phase 6 — stop the per-frame whole-state clone (committed, `89c0670`)

The plan blamed five panel `Layout` paths. **Four of them were not per-frame at
all** — they are `LoadFromState`, which runs from the panel constructors and
`Window.load()`. A profile said something else entirely: **75.2 % of every
allocation** in `BenchmarkEditorWindow_TabCycling` came from
`editor_state_model.EditorState.Clone`, 97 % of that from `CloneZoneContentRows`.

The cost was the clone *mechanism*. Every row-slice clone ran
`linq.FromSlice(x).Select(f).ToSlice()`, which allocates three closures and boxes
the accumulator **before it looks at the source**, then regrows with `append`.
`Clone` runs eight such chains and six are empty on a default state — a quarter
of all allocations in the benchmark were chains projecting nothing.

- Added `linq.SelectSlice` — eager, `nil` for empty, sized once.
- Pointed the five frame-path clone helpers at it.
- `handleZoneContentDialogClicks` stopped cloning the whole state per frame;
  `openZoneContentDialog` takes a getter and reads the state on the click.

| | allocs/op | B/op | ns/op (median of 6) |
| --- | --- | --- | --- |
| Before | 12,690 | 1,045 KB | 4.05 M |
| After | **4,773** | **720 KB** | **3.57 M** |

**No per-panel view structs were built** — the profile said they would buy
nothing. **The clones that hand state *out* of the model stay** (§1.1).

### The `godox` TODO (committed, `32c92f0`)

`LoadState` returned `(*EditorStateDto, []string, error)` — the handler unpacked
`ValidateEditorState`'s result only for every caller to re-pair it. It now
returns `(*editor_state_dto.EditorStateValidationDto, error)`. Lint baseline is
**0 issues** again, not 1.

### Phase 12 — the layering gate (committed, `938ef55`)

New [layering_test.go](../test/unit/architecture/dependency/layering_test.go) in
the existing `dependency_test` package, reusing its `findImports` walker rather
than adding a third one.

| Rule | Violations | Allow-list |
| --- | --- | --- |
| Entities must not import models / dtos / services / handlers / helpers | **0** | none |
| Entities named only by repositories, models, entities, mappers, `*_helpers` | 113 files | **23 packages** |
| DTOs named only by handlers, dtos, `app/` | 6 files | **3 packages** |

Allow-lists are **per package, not per file** — 113 literals is a snapshot nobody
maintains, 23 is a list a human shrinks. The gate was **proven to trip**: dropping
`"internal/services/zones"` fails the test and prints all five offending files.
`internal/helpers/data` (`Vec2`, `Tuple`, `Adjacency`) is carved out of the
"entities must not import helpers" rule, matching the long-standing AGENTS.md
§4.4 exception; two tests pin the boundary either side of it.

### Phase 7 — docs (uncommitted)

- `README.md` — the generation flow no longer claims the GUI "collects widget
  input into `dtos.EditorStateDto`"; the entity/model/mapper layer descriptions
  and the directory tree now match the tree.
- Backlog §2.1 and §1.5 rewritten as self-contained ✅ FIXED entries; §2.6 added
  last session; header counts, batch table (I and N done, O added) and the §9
  baselines refreshed to 73.9 % and the new layering gate.
- `test_observations.md` — new entry: **the per-frame allocation budget has no
  automated guard.** The benchmark needs a GPU and never runs in CI; an
  `AllocsPerRun` assertion was considered and rejected as too flaky over a Gio
  frame.
- **408 dangling markdown links repaired.** Moving the docs to `.agent/backlog/`
  made them two levels deep while they still used single-`../` paths: 156 in
  `backlog-opus5.md`, 252 in `review-opus5-08-04.md`. Verified byte-safe —
  393 insertions, 393 deletions, no line-ending or BOM change.
- `AGENTS.md` §4.2.2's interface-placement example was factually wrong:
  `zoneLabelProviderInterface.go` moved to `zone_interfaces/` when `zones` passed
  five implementations, so it is now a rule-**2** example, not rule 1. Fixed.

**Links left dangling on purpose:** the evidence links in
`review-opus5-08-04.md` and the historical parts of `backlog-opus5.md` that name
files which genuinely no longer exist (`internal/dtos/editorStateDto.go`,
`internal/models/zoneContentRowSave.go`, `runnerHandler.go`, `backlog.md`).
Rewriting them would falsify the historical record.

## 4. File modifications this session

**Committed** (`89c0670`, `32c92f0`, `938ef55`) — Phase 6, the `LoadState`
cleanup and Phase 12. See those commits for the list.

**Uncommitted (Phase 7, docs only, no Go code):**

| File | Change |
| --- | --- |
| [README.md](../README.md) | Architecture layers 4/5/7, directory tree, generation-flow header |
| [AGENTS.md](../AGENTS.md) | §4.2.2 interface-placement example corrected |
| [.agent/backlog/backlog-opus5.md](backlog/backlog-opus5.md) | §1.5 and §2.1 rewritten as FIXED; counts, batch table, baselines; 156 links |
| [.agent/backlog/review-opus5-08-04.md](backlog/review-opus5-08-04.md) | 252 links |
| [.agent/backlog/test_observations.md](backlog/test_observations.md) | Phase 6 allocation-guard entry |
| `.agent/plans/batch-i-editor-state-rework.md` | **deleted** |
| [.agent/session-carry-forward.md](session-carry-forward.md) | this file |

## 5. Tests added or updated

- `test/unit/internal/helpers/linq/slice/selectSlice_test.go` — 5 cases
  (projection, empty → nil, nil → nil, non-aliasing, named slice type).
- `test/unit/architecture/dependency/layering_test.go` — 3 tree-scanning rules
  plus 7 predicate tests that prove the rules trip.
- `LoadState`'s signature change touched both mocks, the `guiHandler` stub, both
  `loadState` unit suites and `editorStateRoundTrip_integration_test.go`. The
  round-trip assertion got more honest — it compares `loaded.State` directly
  instead of wrapping the expectation in a throwaway DTO.

**Last full run: all green.** `go test ./test/unit/...`, `go test ./test/...`
(untagged), `go test -tags=integration_test ./test/integration/...`, GPU suite
**without `-update`** (23.7 s).

## 6. Gates

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | exit 0 |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 |
| Unit / untagged / integration / GPU (no `-update`) | pass |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **73.9 %** (floor 72.5 %) |

The two frozen fixtures under `test/test_helpers/testdata/` and the untagged
`editorStateWireFormat_integration_test.go` passed **unchanged** through the
whole batch — the `.gen.json` shape never moved.

## 7. Git status snapshot

- **Branch:** `AD/fixing_some_stuff_08-12`
- **HEAD:** `938ef55 "Batch I wip 11"` (origin is at `32c92f0`)
- **Uncommitted:** the six doc files in §4. **Nothing staged, nothing committed
  by the agent** — review, stage and commit yourself.

## 8. Rejections / things not done

- **Rejected — per-panel view structs** (the plan's own Phase 6 bullet). The
  profile showed no per-frame whole-state reads left to project. Do not
  re-propose without a measurement.
- **Rejected — collapsing the `UpdateCurrentState` → `ValidateEditorState`
  double clone.** It is the largest remaining allocation item, but it means
  letting `updateFunc` mutate through to the live `current`, undoing §1.1. With
  each clone ~9× cheaper the trade is bad.
- **Rejected — an `AllocsPerRun` regression assertion.** Too flaky over a Gio
  frame; recorded in `test_observations.md` instead.
- **Rejected (earlier, by the owner) — the nine-group DTO with struct
  conversion.** It was built in the first Phase 10, then deleted. Do not propose
  it again.
- **Not done — draining the layering allow-lists.** That is backlog §2.6 /
  batch O.
- **Reverted — the `app/` → `internal/mappers` depguard permission.** Phase 10
  added it; the rework deleted its only consumer. It is denied again.

## 9. Open questions

1. **None block the next batch.**
2. **Repo memory duplication** (`/memories/repo/conventions.md`) — flagged eight
   sessions running: ~1,300 lines, roughly four copies of the same body. Still
   needs a dedupe pass.
3. Backlog §2.6 step 2 asks a real design question: is `internal/dtos` /
   `internal/handlers` naming `entities.Zone` a breach at all, or does the
   `.rmg.json` schema vocabulary deserve a documented carve-out like
   `internal/helpers/data` has? Answer it before draining that part of the list.

## 10. Next recommended actions

1. Review and commit the Phase 7 docs.
2. **Batch O** (backlog §2.6 step 1) — give `internal/services/bonuses`,
   `pickers` and `zone_content` a model-side request/result pair with the
   handler mapping onto it, the way `internal/services/editor` got in Phase 10.
   It is the smallest step and it clears one allow-list entirely.
3. **Batch J** (backlog §2.2 branch B) — zone tier single source of truth. It
   benefits from the model layer Batch I built.

## 11. Carry-forward prompt

> Read `AGENTS.md` first — especially **§4.4.1**, which now carries the
> Entity/Model/DTO doctrine that used to live in a plan file. In one line:
> **Entity** (`internal/entities/`) is the database layer, json tags only;
> **Model** (`internal/models/`) is the service layer and **owns the structure
> and all business logic**; **DTO** (`internal/dtos/`) is the `app/` ↔
> `internal/` crossing and is thin. *"Redefinition is expected in Models, but it
> should never happen in DTOs"* — so `EditorStateDto` is literally
> `struct { editor_state_model.EditorState }`, and **a DTO embedding a Model is
> intended**. `app/` MAY hold a Model; only the crossing must be a DTO. Do not
> "fix" either. The doctrine is enforced by
> `test/unit/architecture/dependency/layering_test.go`; its two allow-lists
> **only ever shrink** — never add an entry, clean the package instead.
>
> **Batch I is closed** (backlog §2.1 and §1.5, both ✅ FIXED). Branch
> `AD/fixing_some_stuff_08-12`, HEAD `938ef55`, with the Phase 7 doc changes
> uncommitted for review. The plan file was deleted on purpose; recover it with
> `git show 938ef55:.agent/plans/batch-i-editor-state-rework.md` if you need the
> phase record. Next work is **batch O** (backlog §2.6 step 1: the three
> DTO-consuming services), then **batch J** (§2.2 branch B).
>
> The hard rules, one line each: never modify `data/`, `internal/registry/`, or
> **anything under `internal/entities/template/`**; everything must build and run
> on Windows and Linux (use `path/filepath`; chain PowerShell with `;`, never
> `&&`); every change ships with tests and unit coverage must not drop below
> 72.5 % (currently 73.9 %); the lint baseline is **0 issues**; **never stage and
> never commit** — the owner reviews, stages and commits, so **use `Move-Item`,
> never `git mv`**, and delete with `Remove-Item`, never `git rm`; never change
> where `.rmg.json` is written and never persist the output directory; never run
> a bulk in-place rewrite over the repository, and **never round-trip a `.go`
> file through `Get-Content`/`Set-Content`** — it joins every line and corrupts
> the file (this has happened once). For bulk text edits on *markdown*, use
> `[System.IO.File]::ReadAllText` / `WriteAllText` with an explicit
> `UTF8Encoding($false)` and verify insertions == deletions afterwards.
>
> Standing traps: **nil is load-bearing** on the regeneration path (nil
> `Previous` = first generation, nil `Next` = unarmed debounce) — two dereference
> bugs were fixed there, do not reintroduce them; the two frozen fixtures under
> `test/test_helpers/testdata/` plus the untagged
> `editorStateWireFormat_integration_test.go` must keep passing **unchanged**,
> comparing **parsed objects, never bytes**; `BenchmarkEditorWindow_TabCycling`
> needs a GPU and never runs in CI, so the ~4,773 allocs/op figure has no
> automated guard — re-measure by hand when touching `EditorState.Clone`,
> `linq.SelectSlice` or the clone helpers; and cap sessions at ~50 messages,
> handing off through this file.
>
> Full handoff in `./.agent/session-carry-forward.md`.
