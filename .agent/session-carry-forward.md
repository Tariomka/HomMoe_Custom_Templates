# Session carry-forward — 2026-09-04

## 1. Session goal

Take **backlog §2.6 step 2** — a *decision*, not a sweep: is `internal/dtos` /
`internal/handlers` naming `entities.Connection` a breach, or does the
`.rmg.json` vocabulary deserve a documented carve-out? The owner ruled, and the
ruling was then executed in the same session as **batch P**.

**Batch P is complete.** All five phases done, records written, every gate
green, no golden and no fixture moved.

## 2. The ruling (the durable answer to step 2)

**It is a breach. Base `internal/entities` gets NO carve-out.** Four findings
decided it — all four are recorded in backlog §2.6 step 2, and the first is a
general lesson:

1. **The backlog framing was stale.** It said "`entities.Connection` /
   `entities.RmgTemplate`". Measured, `entities.RmgTemplate` was named **nowhere**
   in `internal/dtos`, `internal/handlers` or `app/gui` — batch J had already
   removed it. The whole breach was **one type in 8 files**.
   *Re-measure a backlog item before ruling on it.*
2. **The model twin already existed** — batch J built
   `template_variant_model.Connection`, field-identical, converters exported.
3. **The seams were already half-migrated** — `connection_editor` signatures
   mixed `[]template_model.Zone` with `[]entities.Connection`, and
   `zoneEditorGeometryService` converted **mid-service**.
4. **`IsUserAdded` was editor state inside the schema mirror** as `json:"-"`,
   with `editor_state.ManualConnectionSave` already carrying a **sidecar** copy
   *because the entity could not serialize it*. The workaround was the proof.

Also recorded in `.agent/memories/settled-decisions.md` so it is not
re-litigated.

## 3. Features added / changed

No user-visible behaviour change. Structural:

- `entities.Connection` → `template_model.Connection` across `internal/dtos`,
  `internal/handlers` (+`handler_interfaces`), `internal/services/connection_editor`,
  `app/gui/{dialogs,drivers,models,panels}` and `editor_state_model`.
- `template_model/converters.go` gained singular `ToConnectionEntity` /
  `ToConnectionModel`.
- **`IsUserAdded` removed from the protected `.rmg.json` entity** (owner-approved).
  The model keeps it; `editor_state.ManualConnectionSave`'s sidecar persists it.
  `ConnectionBuilder.WithIsUserAdded()` deleted with it.
- `entityNamerAllowList` **21 → 14** entries; breach **84 files/21 packages → 64/14**.
- Four `ToConnection*` calls deleted as pure ceremony — they crossed a boundary
  that no longer exists.

## 4. File modifications

**70 files changed, +331/−605**, plus one deletion and one rename. Highlights:

| File | Change |
| --- | --- |
| [internal/entities/template/template_variant/connection.go](../internal/entities/template/template_variant/connection.go) | ⚠ **PROTECTED (§2.1), owner-approved.** `IsUserAdded` + comment removed. `git diff --stat` on that whole tree = **1 file, 5 deletions**. |
| [internal/handlers/templateHandler.go](../internal/handlers/templateHandler.go) | `ToConnectionEntities` on assign; **added the positional connection re-attach** beside the zone one. |
| [internal/models/editor_state_model/manualConnectionSave.go](../internal/models/editor_state_model/manualConnectionSave.go) | Takes/returns `[]template_model.Connection`; the save's field stays an entity. Mirrors `manualZoneSave.go`. |
| [internal/services/connection_editor/connectionEditorService.go](../internal/services/connection_editor/connectionEditorService.go) | `NewDefaultConnection` returns a model, converts the shared builder's entity, sets `IsUserAdded` by hand. |
| [test/unit/architecture/dependency/layering_test.go](../test/unit/architecture/dependency/layering_test.go) | 7 allow-list entries removed; comment states **exactly one** permanent entry (`file_service`) and names the generator explicitly as debt. |
| `test/unit/.../connectionBuilder/withIsUserAdded_test.go` | **Deleted** (`Remove-Item`). |
| `.agent/plans/batch-p-connection-follows-zone.md` | The plan, with Final Recap + Deployment Plan. **Renamed from `batch-k-`** — see §7. |
| [.agent/backlog/backlog-opus5.md](backlog/backlog-opus5.md) | §2.6 rewritten (steps 2+3 ✅, ruling, recount, retitle); header and §8 updated with batch **P**. |
| `.agent/memories/{template-model,settled-decisions}.md` | Batch P section; the no-carve-out decision. |

One test was **deleted, not migrated**:
`TestWhenSaveFlagDiffersFromEmbeddedConnectionFlag_SaveFlagWins` — the entity no
longer has a flag to disagree with, so the scenario is unconstructible. The
surviving `..._RestoresEachFlagOntoConnection` still fails if the copy is dropped.

## 5. Tests added or updated

- **Added** `TestWhenAnAppliedConnectionIsUserAdded_KeepsTheFlagThroughTheEntityRoundTrip`
  in [updateTemplate_test.go](../test/unit/internal/handlers/templateHandler/updateTemplate_test.go).
  Mutation-verified twice: once against `ToConnectionModel` (phase 1) and, more
  importantly, against the `UpdateTemplate` re-attach after the entity field was
  removed (phase 4) — proving it is **not vacuous**.
- ~35 test/mock files migrated to the model type.

**Final gate run, all green:**

| Gate | Result |
| --- | --- |
| `go build ./...` / `go vet` (both tag sets) | exit 0 |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | passed |
| `wire diff ./internal/composition/...` | exit 0 |
| `go test ./test/unit/... -count=1` + coverage | exit 0 — **74.3 %** (floor 72.5 %, unchanged) |
| `go test -tags=integration_test ./test/integration/...` | exit 0 |
| `go test -tags='integration_test,gui' ./test/integration/gui/...` | exit 0, **no `-update`** |
| `git status --short -- '*.golden'` / `testdata/*` | `goldens=0`, `fixtures=0` |
| `golangci-lint-v2 run ./...` | **0 issues** |

## 6. Git status snapshot

Branch **`AD/fixing_some_stuff_08-12`**, head `eec6904 Batch J done` — batch J
is fully committed. 70 modified files, 1 deletion, 1 untracked plan.
**Nothing was staged or unstaged by the agent (AGENTS §2.5).**

⚠ **One index oddity for the owner to resolve.** The owner had staged
`A .agent/plans/batch-k-connection-follows-zone.md`. Renaming it (see §7) left
that path as `AD` — *added in index, deleted in worktree* — while the real file
`.agent/plans/batch-p-connection-follows-zone.md` is untracked (`??`). Resolving
it means staging, which the agent must not do. `.agent/plans/batch-j-…md` also
remains staged-deleted from the previous session.

## 7. Rejections / corrections

- **The carve-out option was rejected by the owner** in favour of cleaning. The
  "dtos+handlers only" option was rejected in the framing because it would have
  left `app/` holding an entity while the DTOs held a model — the handler
  converting backwards relative to AGENTS §4.4.1 rule 2.
- ⚠ **Batch letter corrected mid-session: this is batch P, not K.** §8 of the
  backlog already used **⚠ K** for the owner-gated group (§2.2 Branch A, §2.4,
  §2.5, §6.1). Letters through O were taken. The plan file was renamed with
  `Move-Item` and all references updated. **Check §8 for a free letter before
  naming a batch.**
- A plan-authoring error was caught in phase 1: the plan sketched a nested
  `templateHandler/updateTemplate/` test folder. AGENTS §4.6 puts the folder at
  the *implementation file* name and the file at the *function* name.
- ‼ **Corrected on owner review: `template_generator` is NOT a permanent entity
  exception.** This session inherited "`file_service` + `template_generator` are
  both permanent" from the previous carry-forward and wrote it into backlog
  §2.6 step 4, `settled-decisions.md`, this file **and** the repo's
  `entityNamerAllowList` comment — without checking
  `.agent/memories/settled-decisions.md`, which already carried a ⚠ correction
  dated 2026-09-03 stating that exact claim is wrong. All five places are now
  fixed to say **exactly one permanent entry, `file_service`**, with the
  generator named explicitly as debt. **A claim inherited from a carry-forward
  is not evidence — verify it against the memories.**

## 8. Open questions

- **Backlog §2.6 step 4 is the only step left** — 64 files / 14 packages, all
  generator-side. Step 2's ruling makes it sanctioned work rather than an open
  question. It is also the batch that would let `TemplateGenerator.Generate`
  build the model directly instead of building an entity and lifting it (§2.2
  records that lift as a *migration tactic*, not a design). ⚠ **The floor is
  ONE package, `file_service`, not two** — `template_generator` and its topology
  tree are debt, per `settled-decisions.md` §"Entity confinement".
- Still unreconciled from the previous session: two `TabCycling` benchmark
  baselines disagree (~5,699 vs 6,640 allocs/op), taken on different trees.

## 9. Next recommended actions

1. **Owner reviews and commits batch P.** Follow the Deployment Plan in
   `.agent/plans/batch-p-connection-follows-zone.md` — especially step 1 (the
   protected diff is 1 file / 5 deletions and nothing else) and step 5, the
   in-app smoke test: draw a connection **by hand**, Apply, save, reload, and
   confirm it is still marked user-added.
2. Resolve the index oddity in §6 while committing.
3. Delete the transient docs once it lands: the plan and this file. **Backlog
   §2.6 is the surviving record** and is written to stand alone.
4. Take **§2.6 step 4** as its own batch (letter **Q**).

## 10. Carry-forward prompt

> Read `AGENTS.md` first. The hard rules, one line each: never modify `data/`,
> `internal/registry/` or anything under `internal/entities/template/` **without
> explicit owner approval** — `internal/entities/editor_state/` is *not*
> protected; everything must build and run on Windows and Linux
> (`path/filepath`; chain PowerShell with `;`, never `&&`); every change ships
> with tests and unit coverage must not drop below 72.5 % (currently **74.3 %**),
> lint baseline **0 issues**; **never stage and never commit** — `Move-Item` not
> `git mv`, `Remove-Item` not `git rm`; never change where `.rmg.json` is written
> and never persist the output directory; never run a bulk in-place rewrite and
> **never round-trip a `.go` file through `Get-Content`/`Set-Content`** — use
> `gofmt -r` on an explicit file list and verify insertions == deletions per file.
>
> **Batch P (Connection follows Zone, backlog §2.6 steps 2 + 3) is COMPLETE** —
> five phases, every gate green, no golden and no fixture moved. It ruled that
> base `internal/entities` gets **no vocabulary carve-out**: naming a `.rmg.json`
> schema type below the repositories is a genuine breach. `entityNamerAllowList`
> went 21 → 14 entries and the breach 84 files/21 packages → **64/14**, now
> entirely generator-side. With owner approval it removed `IsUserAdded` from the
> protected `internal/entities/template/template_variant/connection.go`. Batch P
> is **uncommitted** (70 modified files, 1 deleted, 1 untracked plan); batch J is
> committed through `eec6904` on `AD/fixing_some_stuff_08-12`.
>
> Standing traps this codebase punishes: **nil is load-bearing** three times over
> — nil `Previous` = first generation, nil `Next` = unarmed debounce, nil
> `Zone.Quality` = "infer it"; **two re-attach lines in
> `templateHandler.UpdateTemplate` are load-bearing** — `updated.Variants[0].Zones`
> carries the tier and `.Connections` carries the user-added flag across an Apply,
> both mutation-verified, both look redundant and are not; the persisted tier is
> `*int8` because `omitempty` on a plain `int8` would drop every Plastic zone
> (ordinal 0); the two frozen fixtures under `test/test_helpers/testdata/` and the
> untagged `editorStateWireFormat_integration_test.go` must keep passing unchanged
> and compare **parsed objects, never bytes**; `cmd/testlayoutcheck` matches
> test-only export names tree-wide, so grep any new accessor name first; a file
> gets `//go:build integration_test` **only** if it calls a `*_testexports.go`
> accessor; `helpers.MapSlice`/`MapPointer` preserve nil-vs-empty where
> `linq.SelectSlice` does not; and `golangci-lint --fix` wraps a long signature as
> `param,\n) Ret {` where the house style is `param) Ret {`.
>
> Three lessons from batch P worth keeping: **when migrating a type, grep its
> converters as well as its name** (three files called `ToConnectionEntities`
> without ever naming `entities.Connection` — only the compiler found them);
> **`git status` cannot verify a mutation revert** when the file is already dirty
> from the same batch, so grep for the mutation itself; and **check backlog §8 for
> a free batch letter before naming a batch** (K was already taken, this one had
> to be renamed to P mid-session).
>
> ‼ And the one that cost a review round-trip: **a claim inherited from a
> carry-forward is not evidence.** This session re-propagated "`file_service` +
> `template_generator` are both permanent entity exceptions" out of the previous
> handoff and wrote it into the backlog, the memories AND the repo's
> `entityNamerAllowList` comment — while `.agent/memories/settled-decisions.md`
> already carried a ⚠ correction, dated 2026-09-03, saying that exact claim is
> wrong. **Check `settled-decisions.md` and `architecture.md` before restating
> any "by decision" / "permanent" / "do not clean" claim.**
>
> Next up is **backlog §2.6 step 4**, the last step: 64 files in 14 packages,
> all generator-side. Step 2 already ruled it sanctioned work rather than an open
> question, and it is the batch that would let `TemplateGenerator.Generate` build
> the model directly instead of building an entity and lifting it. ⚠ **EXACTLY
> ONE of the 14 is permanent: `internal/services/file_service`** (it writes
> `.rmg.json`, and even there entities stay inside it, never in a signature).
> **The floor is one, not two.** `template_generator` and its whole topology tree
> are **DEBT** — the "the golden-template test makes the lifted entity a proof"
> argument is a *migration tactic*, explicitly rejected by the owner on
> 2026-09-03. See `.agent/memories/settled-decisions.md` §"Entity confinement".
>
> Full handoff in `./.agent/session-carry-forward.md`.
