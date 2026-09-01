# Session Carry-Forward — 2026-09-01 (Batch O closed, Batch J planned)

## 1. Session goal

Finish **Batch O phase 2** (documentation only), then scope and plan **Batch J**
(backlog §2.2 branch B — zone tier single source of truth).

Both are done. **No Go code was written this session** — batch O phase 2 was
docs, and batch J is planned but not started.

Plans: [.agent/plans/batch-o-picker-view-model.md](plans/batch-o-picker-view-model.md)
(complete), [.agent/plans/batch-j-zone-tier-source-of-truth.md](plans/batch-j-zone-tier-source-of-truth.md)
(written, phase 1 not started).

## 2. Read this before touching anything

**`AGENTS.md` §4.4.1** carries the Entity/Model/DTO doctrine. The three things
most often gotten wrong:

- **The Model owns the structure.** *"Redefinition is expected in Models, but it
  should never happen in DTOs."* `EditorStateDto` is literally
  `struct { editor_state_model.EditorState }`. **A DTO carrying a Model is
  intended** — do not "fix" it.
- **`app/` may hold a Model**; only the *crossing* into `internal/` is a DTO.
  `app/` → `internal/models` is fine; `app/` → `internal/mappers`,
  `internal/services`, `internal/repositories`, `internal/validators` is not.
- Enforced by
  [test/unit/architecture/dependency/layering_test.go](../test/unit/architecture/dependency/layering_test.go).
  **Its allow-lists only ever shrink — never add an entry, clean the package.**

**The measurement lesson, still worth internalising.** To decide whether a type
crosses a layer boundary, follow the *method signatures the other side calls* —
never grep for the type name. Six DTOs consumed with `:=` looked like they never
reached `app/` and did. This is now recorded permanently in backlog §2.6 step 1.

## 3. What shipped this session

### Batch O phase 2 — docs (unstaged, needs a commit)

Phase 1 was **committed by the owner** as `85a7d76 "Batch O wip"` between
sessions. Phase 2 made the backlog agree with the code:

- **§2.6 step 1 rewritten** as ✅ DONE, closed at **two allow-list entries, not
  zero**: `pickers` was view-model logic and was deleted into `app/gui/models/`;
  `bonuses` and `zone_content` keep their DTOs by owner decision. The §0.1
  measurement correction is folded in.
- **§2.6 heading kept** — it names the *113 entity* files and never implied the
  DTO list drains to zero; a one-line note now says which half is closed.
- The "6 files in 3 packages" evidence line records the close at **4 files in 2
  packages**. Entity counts (113 / 85 / 11 / 11 / 6) were **re-measured and are
  unchanged** — the deleted picker DTOs named no entity, and the new
  `app/gui/models` files name none either.
- **§8 row O marked ✅ done**, scoped to "§2.6 step 1", stating that steps 2–4
  stay open.
- Coverage refreshed **73.9 % → 73.8 %** in all three places that quote it.
- Plan file closed out: phase 1 and 2 checkboxes ticked, Phase Summaries, Final
  Recap and Deployment Plan written.

**Verified the allow-list is still load-bearing**: dropping
`internal/services/bonuses` makes `TestWhenDtoConsumersAreScanned_…` fail naming
`bonusEntryService.go` and `bonusEntryServiceInterface.go`; restored immediately.

### Batch J — planned only, no code

[.agent/plans/batch-j-zone-tier-source-of-truth.md](plans/batch-j-zone-tier-source-of-truth.md),
five phases, ~31 production files. Two exploration passes mapped the whole
`entities.Zone` surface before any decision was taken.

## 4. Owner decisions for batch J (settled — do not relitigate)

1. **Branch B, not A.** No `Quality` field on `entities.Zone`. Branch A is not a
   shortcut to fall back on when phase 4 gets tedious.
2. **`models.QualifiedZone`, and it embeds `entities.Zone`.** Field promotion
   keeps every `zone.Name` / `zone.Layout` compiling, turning most of phase 4
   into a type swap. Name chosen over `PlannedZone` (misleading — the type also
   carries hand-created and `.rmg.json`-loaded zones) and `TieredZone`.
3. **The store behind the wrapper is the generator.** `Variant.Zones` is
   protected and is rewritten on *both* the generate and apply-back paths, so a
   wrapper with no store is the classifier with extra steps. `Generate` returns
   the tiers it planned; `drivers.State` carries the index.
4. **`Generate()` changes signature** to `(*models.GeneratedTemplate, []string)`;
   the ~130 test call sites take the mechanical `actual` → `actual.Template`
   edit. A `GenerateWithTiers()` second entry point was **rejected**.
5. **`IZoneTierService` replaces `IZoneClassifier.GetQuality`** at all 8
   consumers and absorbs `GetGuardQuality` / `GetConnectionGuardQuality`.
   `ZoneClassifier` becomes the private fallback and **can never be deleted** —
   a raw `.rmg.json` load has no recorded tier.
6. **Persisted tier is nullable.** `omitempty` on a plain `int8` would silently
   drop every Plastic zone (`QualityLowest` is 0). Entity stores `*int8` (it may
   not import `internal/models/neutral_zone`); model exposes
   `*neutral_zone.Quality`.
7. **The 9 zone-editor DTOs carry `[]models.QualifiedZone`.**
8. **Output changes are approved.** Zones the classifier calls `Unknown` will
   start getting planned mandatory-content rows, arena eligibility, castle
   city-guard values and connection guard defaults. **Every delta must be
   enumerated in the phase summary that causes it.**
9. **One batch, phased, reviewed per phase.**

## 5. Two traps found while scoping batch J

- **`PreviewLayoutService` bypasses DI** — `NewPreviewLayoutService()` hard-builds
  its own `NewZoneClassifier()`, so it is not the wire singleton and would
  silently keep inferring. Phase 1 fixes it (~60 test call sites follow).
- **`Unknown` → Plastic is a silent down-tier** —
  `GetNeutralZoneProfile(QualityUnknown)` returns the Lowest profile, so an
  unclassifiable zone whose castles are rebuilt gets Plastic city stats today.

## 6. Gates

Nothing executable changed this session. Re-verified after the doc edits:

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go test ./test/unit/... -count=1` | pass (no FAIL) |
| `go test ./test/unit/architecture/... -count=1` | pass, and fails correctly when an allow-list entry is dropped |

Standing baselines from batch O: coverage **73.8 %** (floor 72.5 %), lint **0
issues**, GPU suite passes **without `-update`**.

## 7. Git status snapshot

- **Branch:** `AD/fixing_some_stuff_08-12`, ahead of origin by 5.
- **HEAD:** `85a7d76 "Batch O wip"` (batch O phase 1, committed by the owner).
- **Unstaged, needs review and a commit:**
  - `.agent/backlog/backlog-opus5.md` — §2.6 step 1, §8 row O, coverage figures
  - `.agent/plans/batch-o-picker-view-model.md` — phase 2 closed out
  - `.agent/session-carry-forward.md` — this file
  - `.agent/plans/batch-j-zone-tier-source-of-truth.md` — new, untracked
- The agent has not touched the index (AGENTS §2.5).

## 8. Rejections / things not done

- **Rejected — `GenerateWithTiers()` as a second entry point.** It would exist
  only to dodge test churn; §3.1 calls that a speculative abstraction.
- **Rejected — `PlannedZone` / `TieredZone` as the wrapper name.**
- **Rejected — branch A** (a `json:"-"` `Quality` on the protected schema).
- **Rejected — persistence-only tier store**, and **rejected — a side index
  without a wrapper**. The owner chose the wrapper plus a generator-sourced
  index.
- **Not done — §2.6 steps 2–4** (the 113-file entity list). Still open. Batch J
  phase 5 may shrink it for free, but that is not its purpose.
- **Not done — any batch J code.** Deliberate: batch O's docs are uncommitted,
  and mixing two batches in one working tree makes review harder.

## 9. Open questions

1. **None block batch J phase 1.** It is fully specified.
2. **Repo memory duplication** (`/memories/repo/conventions.md`) — flagged ten
   sessions running: ~1,300 lines, roughly four copies of the same body. Still
   needs a dedupe pass.
3. §2.6 step 2 still asks a real design question before that list can drain: is
   `internal/dtos` / `internal/handlers` naming `entities.Zone` a breach at all,
   or does the `.rmg.json` schema vocabulary deserve a documented carve-out like
   `internal/helpers/data` already has? **Batch J phase 4 partly answers it in
   practice** — the 9 DTOs stop naming `entities.Zone` because they carry the
   wrapper instead.

## 10. Next recommended actions

1. **Review and commit the batch O phase 2 docs** (four files in §7).
2. **Start batch J phase 1** — `IZoneTierService` over the existing classifier.
   Pure indirection: generated output must stay byte-identical and the GPU suite
   must pass **without `-update`**. It is the safety net that proves the seam
   before any behaviour moves in phase 2.
3. Phases 2–5 in order, each reviewed and committed on its own.

## 11. Carry-forward prompt

> Read `AGENTS.md` first — especially **§4.4.1**, the Entity/Model/DTO doctrine.
> In one line: **Entity** (`internal/entities/`) is the database layer, json tags
> only; **Model** (`internal/models/`) is the service layer and **owns the
> structure and all business logic**; **DTO** (`internal/dtos/`) is the `app/` ↔
> `internal/` crossing and is thin. *"Redefinition is expected in Models, but it
> should never happen in DTOs"* — **a DTO carrying a Model is intended**, and
> `app/` MAY hold a Model. Do not "fix" either. Enforced by
> `test/unit/architecture/dependency/layering_test.go`; its allow-lists **only
> ever shrink** — never add an entry, clean the package instead.
>
> **Batch O is finished.** Phase 1 is committed (`85a7d76`); phase 2 was docs and
> is unstaged awaiting your commit. Do not reopen it: the DTO allow-list ends at
> **two entries** (`bonuses`, `zone_content`) by decision, not debt.
>
> **Your work is `.agent/plans/batch-j-zone-tier-source-of-truth.md` — read it
> and start phase 1.** Batch J gives the zone tier a single source of truth
> (backlog §2.2 branch B). Phase 1 is pure indirection: put an
> `IZoneTierService` over the existing `ZoneClassifier`, swap all 8 consumers
> onto it, and fix `PreviewLayoutService`, which hard-constructs its own
> classifier instead of taking the wire singleton (~60 test call sites follow).
> **Generated output must be byte-identical and the GPU suite must pass without
> `-update`** — phase 1 is the safety net that proves the seam before behaviour
> moves in phase 2. §0 of that plan holds nine settled owner decisions; do not
> relitigate them.
>
> **A measurement lesson worth internalising:** to decide whether a type crosses
> a layer boundary, follow the *method signatures the other side calls* — never
> grep for the type name. Six DTOs consumed with `:=` looked like they never
> reached `app/` and did.
>
> The hard rules, one line each: never modify `data/`, `internal/registry/`, or
> **anything under `internal/entities/template/`** — batch J is branch B
> precisely so it never has to (`internal/entities/editor_state/` is *not*
> protected); everything must build and run on Windows and Linux (use
> `path/filepath`; chain PowerShell with `;`, never `&&`); every change ships
> with tests and unit coverage must not drop below 72.5 % (currently 73.8 %);
> the lint baseline is **0 issues**; **never stage and never commit** — the owner
> reviews, stages and commits, so **use `Move-Item`, never `git mv`**, and delete
> with `Remove-Item`, never `git rm`; never change where `.rmg.json` is written
> and never persist the output directory; never run a bulk in-place rewrite over
> the repository, and **never round-trip a `.go` file through
> `Get-Content`/`Set-Content`** — it joins every line and corrupts the file. For
> bulk text edits on *markdown*, use
> `[System.IO.File]::ReadAllText`/`WriteAllText` with an explicit
> `UTF8Encoding($false)`, then verify insertions == deletions.
>
> Standing traps: **nil is load-bearing** on the regeneration path (nil
> `Previous` = first generation, nil `Next` = unarmed debounce) and it becomes
> load-bearing again in batch J phase 3 (nil `Quality` = "not recorded", because
> `omitempty` on a plain int8 would silently erase every Plastic zone); the two
> frozen fixtures under `test/test_helpers/testdata/` plus the untagged
> `editorStateWireFormat_integration_test.go` must keep passing **unchanged**,
> comparing **parsed objects, never bytes**; the picker and zone-editor dialogs
> are snapshotted, so the GPU suite must pass **without `-update`**; the test
> layout checker matches test-only export names **tree-wide**, so grep any new
> accessor name before adding it; and `BenchmarkEditorWindow_TabCycling` needs a
> GPU and never runs in CI, so the ~4,773 allocs/op figure has no automated
> guard — re-measure by hand when touching `EditorState.Clone`,
> `linq.SelectSlice` or the clone helpers. Cap sessions at ~50 messages and hand
> off through this file.
>
> Full handoff in `./.agent/session-carry-forward.md`.
