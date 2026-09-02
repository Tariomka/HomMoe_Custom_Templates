# Session Carry-Forward — 2026-09-02 (Batch J phase 1 done, phase 2 designed)

## 1. Session goal

Close **Batch O phase 2** (docs), then plan and start **Batch J** (backlog §2.2
branch B — zone tier single source of truth).

Batch O is finished and committed. Batch J's plan is written, **phase 1 is
complete and reviewed**, and **phase 2 is fully designed but not started** — it
was deliberately not begun because the session hit its ~50-message budget and
starting a ~73-call-site refactor there would hand over a half-migrated tree.

Plan: [.agent/plans/batch-j-zone-tier-source-of-truth.md](plans/batch-j-zone-tier-source-of-truth.md).

## 2. Read this before touching anything

**`AGENTS.md` §4.4.1** carries the Entity/Model/DTO doctrine:

- **The Model owns the structure.** *"Redefinition is expected in Models, but it
  should never happen in DTOs."* `EditorStateDto` is literally
  `struct { editor_state_model.EditorState }`. **A DTO carrying a Model is
  intended** — do not "fix" it.
- **`app/` may hold a Model**; only the *crossing* into `internal/` is a DTO.
- Enforced by
  [test/unit/architecture/dependency/layering_test.go](../test/unit/architecture/dependency/layering_test.go).
  **Its allow-lists only ever shrink — never add an entry, clean the package.**

**Measuring a boundary crossing:** follow the *method signatures the other side
calls*, never a grep for the type name. Six DTOs consumed with `:=` once looked
like they never reached `app/` and did.

## 3. What shipped this session

### Batch O phase 2 — docs (committed)

Backlog §2.6 step 1 rewritten as ✅ DONE, **closed at two allow-list entries, not
zero**; §8 row O marked done and scoped to "step 1"; coverage figures refreshed.
Entity counts re-measured and unchanged. §2.6 steps 2–4 stay open.

### Batch J phase 1 — the tier service (uncommitted at time of writing)

`IZoneTierService` is the only way the application asks for a zone's tier. Six
injection sites moved onto it (mandatory content, gladiator arena, connection
editor, manual reapply, zone editor handler, preview), field named `tierService`.

**`PreviewLayoutService`'s DI bypass is fixed** — it hard-built its own
classifier and was never the wire singleton. It now takes the service, and
`wire_gen.go` proves it.

**`ZoneClassifier` is deleted, by owner decision at review.** The first cut had
the tier service *wrapping* it; the owner asked for the service to own the logic
natively so phase 2 adds the recorded-tier path to the real type rather than
bolting it onto a forwarder. `zoneClassifier.go` and `zoneClassifierInterface.go`
are gone; their bodies (`GetQuality`, `GetGuardQuality`,
`GetConnectionGuardQuality` + `getCenterQuality` / `getTreasureQuality` /
`getSidesQuality`) now live on `ZoneTierService`, which is an empty struct with a
no-arg constructor. The whole classifier unit suite **moved intact** to
`test/unit/internal/services/zones/zoneTierService/`; the mock became
`ZoneTierServiceMock`.

An `iface` lint hit (two identical interfaces) is what first exposed the
redundancy — worth remembering that the linter caught a real design smell.

## 4. Owner decisions for batch J (settled — do not relitigate)

1. **Branch B**, not A. No `Quality` field on `entities.Zone`, ever.
2. **`models.QualifiedZone` embeds `entities.Zone`** and adds `Quality`. Field
   promotion keeps every `zone.Name` compiling, which is what makes phase 4 a
   type swap instead of a rewrite. Name settled over `PlannedZone`/`TieredZone`.
3. **The store behind the wrapper is the generator.** `Variant.Zones` is
   protected and rewritten on both the generate and apply-back paths, so a
   wrapper with no store is just inference with extra steps.
4. **`Generate()` changes signature**; ~73 test call sites take the churn. A
   production `GenerateWithTiers()` was **rejected**.
5. **`IZoneTierService` owns tier questions**; the classifier is deleted (see §3).
6. **Persisted tier is nullable.** `omitempty` on a plain `int8` would silently
   erase every Plastic zone (`QualityLowest` is 0). Entity stores `*int8`; model
   exposes `*neutral_zone.Quality`.
7. **The 9 zone-editor DTOs carry `[]models.QualifiedZone`.**
8. **Output changes are approved**, and every delta must be enumerated in the
   phase summary that causes it.
9. **One batch, phased, reviewed per phase.**

## 5. Phase 2 is designed — read the plan, do not re-derive it

The plan's phase 2 section now carries the full implementation design found by
reading the generation path. The three things that took the analysis:

1. **The index needs no topology changes.** Every neutral zone name in the repo
   comes from `constants.GetNeutralZoneNameFor(plan.Label)` — tournament cluster
   services included — so `Generate` can derive `map[zoneName]Quality` from the
   `neutral_zone.Plans` it already holds, plus hubs at `QualityHighest`. No
   `ZoneFactory` / `TopologyBase` / topology signature changes at all.
2. **Comma-ok is mandatory on every tier lookup.** A missing key yields
   `Quality(0)` = `QualityLowest`, a silent down-tier — the same bug class as the
   `omitempty` trap. It also makes a `nil` map safe, which matters because tests
   will pass one.
3. **The ~73 `Generate()` call sites cannot be done by `gofmt -r` alone** (the
   fix needs a second statement) and **a PowerShell text sweep over `.go` files
   is forbidden**. The route is a *test-local* `generateTemplate(generator)`
   helper in the one test package that holds eight of those files, then one
   `gofmt -r` pass per file. That is not the rejected production entry point.

`QualifiedZone` was **moved to phase 4**, where its consumers are. Phase 2's only
consumer, the arena provider, wants the map — it mutates `variant.Zones` in place
by index. Building the wrapper in phase 2 would be an abstraction ahead of its
caller (§3.1). The design is unchanged; only the file's arrival moved.

**Phase 2 has exactly one behaviour delta:** the gladiator arena. A zone that
inference called `Unknown` scored −1 and could never win the arena; with its
planned tier it can.

## 6. Gates (batch J phase 1, re-run after the classifier absorption)

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 (regenerated, never hand-edited) |
| Unit / untagged / integration | pass |
| **GPU suite, no `-update`** | **pass (23.9 s)** — no pixel moved |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **73.8 %** (floor 72.5 %) |

## 7. Git status snapshot

- **Branch:** `AD/fixing_some_stuff_08-12`.
- Batch O phase 1 is committed as `85a7d76 "Batch O wip"`; batch O phase 2 docs
  were committed after that.
- **Batch J phase 1 is in the working tree, partly staged by the owner during
  review.** The agent has not touched the index (AGENTS §2.5) — do not unstage.
- `.agent/plans/batch-j-zone-tier-source-of-truth.md` and this file are the only
  doc changes outstanding.

## 8. Rejections / things not done

- **Rejected — a production `GenerateWithTiers()`.** Dodging test churn is not a
  reason for a second entry point.
- **Rejected — `PlannedZone` / `TieredZone`** as the wrapper name.
- **Rejected — branch A** (a `json:"-"` `Quality` on the protected schema).
- **Rejected — persistence-only tier store**, and **a side index with no wrapper**.
- **Rejected — keeping `ZoneClassifier`** as a private collaborator behind the
  service (owner asked for outright absorption).
- **Not done — batch J phases 2–5.** Phase 2 is designed, not written.
- **Not done — §2.6 steps 2–4** (the 113-file entity list). Still open.

## 9. Open questions

1. **None block phase 2.** It is specified down to function signatures.
2. **Repo memory duplication** (`/memories/repo/conventions.md`) — flagged eleven
   sessions running: ~1,300 lines, roughly four copies of the same body. Still
   needs a dedupe pass.
3. §2.6 step 2's design question stays open, though batch J phase 4 answers part
   of it in practice: the 9 DTOs stop naming `entities.Zone` because they carry
   the wrapper instead.

## 10. Next recommended actions

1. **Review and commit batch J phase 1** if not already done.
2. **Do batch J phase 2** exactly as the plan's phase 2 section specifies.
   Order that keeps the tree compiling: `GeneratedTemplate` → `Generate` +
   `planZoneTiers` → mock and test call sites → `ResolveQuality` → `PlaceArena`
   → `TemplateLoadDto` / `drivers.State` → new tests.
3. Phases 3–5 in order, each reviewed and committed on its own.

## 11. Carry-forward prompt

> Read `AGENTS.md` first — especially **§4.4.1**, the Entity/Model/DTO doctrine.
> **Entity** (`internal/entities/`) is the database layer, json tags only;
> **Model** (`internal/models/`) is the service layer and **owns the structure
> and all business logic**; **DTO** (`internal/dtos/`) is the `app/` ↔
> `internal/` crossing and is thin. *"Redefinition is expected in Models, but it
> should never happen in DTOs"* — **a DTO carrying a Model is intended**, and
> `app/` MAY hold a Model. Do not "fix" either. Enforced by
> `test/unit/architecture/dependency/layering_test.go`; its allow-lists **only
> ever shrink**.
>
> **Then read `.agent/plans/batch-j-zone-tier-source-of-truth.md` and do phase
> 2.** Batch J gives the zone tier a single source of truth (backlog §2.2 branch
> B). Phase 1 is **complete**: `IZoneTierService` now owns the inference outright
> and `ZoneClassifier` is deleted. **Phase 2 is fully designed in the plan — read
> its "How the tier gets out of the generator" and "The ~73 test call sites"
> sections and follow them; do not re-derive the design.** §0 holds nine settled
> owner decisions; do not relitigate them.
>
> Three things phase 2 will punish you for forgetting: **comma-ok on every tier
> lookup** (a missing key yields `Quality(0)` = `QualityLowest`, a silent
> down-tier); **never sweep `.go` files with PowerShell text replacement** — use
> `gofmt -r` on an explicit file list, since an AST rewrite cannot mangle a file;
> and **`QualifiedZone` belongs to phase 4**, not phase 2, because phase 2 has no
> consumer for it.
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
> `Previous` = first generation, nil `Next` = unarmed debounce) and again in
> phase 3 (nil `Quality` = "not recorded"); the two frozen fixtures under
> `test/test_helpers/testdata/` plus the untagged
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
