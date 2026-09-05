# Session carry-forward — 2026-09-05 (batch R done, batch S next)

## 1. Session goal

Pick the next item from backlog §8. Everything left was owner-gated; the owner chose **§2.4**
and immediately re-scoped it from "swap `[2]float64` for `Vec2`" to the structural problem
underneath: `entities.Zone` was still the in-memory *carrier* for three fields it can never
serialize. Executed as **batch R**, four phases, all complete and green. The wire-format half
was deliberately split off as **batch S**, which is fully scoped but not started.

## 2. Fixes applied

No user-facing bug fixes. Two latent defects closed:

- The manual-edit snapshot round-tripped live zones through `entities.Zone`
  (`SetManualEdits` → `ManualZones` → `GetManualZones`), so a generator stamp survived a
  regenerate-with-manual-edits **only because the entity carried it**, while a `.gen.json`
  reload silently dropped it. The carrier is gone; the model keeps its own zones.
- Two integration tests were `json.Unmarshal`ing the on-disk file into
  `editor_state_model.EditorState`. That only ever worked because the model embedded the
  *tagged* entity groups. Both now decode `editor_state.EditorState`.

## 3. Features added / changed

No user-visible behaviour change. `.gen.json` and `.rmg.json` are byte-identical.

- `editor_state_model.ManualZoneSave` **deleted**; `EditorState.ManualZones` is
  `[]template_model.Zone`. Four converters became two:
  `ToManualZoneSaveEntities([]template_model.Zone)` and
  `ToManualZoneModels([]editor_state.ManualZoneSave)`.
- `cloneZone` / `cloneMainObject` / `cloneRoad` / `cloneTypedRef` **deleted** (~60 lines of
  hand-maintained entity deep-copy, plus the comment obliging future authors to extend it).
  `template_model.Zone.Clone` owns it.
- **One approved protected edit:** `GeneratorPosition`, `GeneratorRing`, `ManualPosition`
  removed from `internal/entities/template/template_variant/zone.go` — 19 deletions,
  0 insertions, nothing else in the protected trees.
- `*[2]float64` → `*data.Vec2[float64]` through the model, four topology producers,
  `ZoneBuilder`, four preview layout files, the zone-editor canvas and dialog, and the whole
  `FindOpenPosition` / `FindOpenZonePosition` chain (two interfaces, three test doubles).
  The three "Is this required to be copied?" comments are gone.

**Decision changed from the written plan:** the converters stayed in `editor_state_model`
instead of moving to `internal/mappers/editorStateMapper.go`. Moving them would have split the
pair from its identical sibling `manualConnectionSave.go` and from
`ToManualEditSettingsModel`/`Entity`, which live in the model package and call them.

## 4. File modifications

**58 modified, 3 added, 4 deleted.** Highlights:

| File | Change |
| --- | --- |
| `internal/entities/template/template_variant/zone.go` | **The protected edit.** Three fields removed. |
| `internal/models/editor_state_model/manualZoneSave.go` | Wrapper + entity clone gone; two converters plus the temporary `toPositionArray`/`fromPositionArray` bridge. |
| `internal/models/editor_state_model/manualEditSettings.go`, `editorState.go` | `ManualZones []template_model.Zone`; `Clone` uses `template_model.Zone.Clone`. |
| `internal/models/template_model/template_variant_model/zone.go` | Positions are `*data.Vec2[float64]`; six copy lines dropped from the converters. |
| `app/gui/models/editorState.go` | `SetManualEdits`/`GetManualZones` no longer convert; `slices.Clone` only. |
| `app/gui/dialogs/zoneEditorCanvas.go`, `zoneEditorDialog.go` | Vec2 drag/add/`manualPositions()`. |
| `internal/services/preview_service/*` (4 files) | Vec2 reads; `a.Subtract(b).Distance()` replaces manual `math.Hypot`. |
| topology producers (3) + `zoneBuilder.go` | Stamp the vector they already hold. |
| `internal/services/connection_editor/zoneEditorService.go` + interface, `guiHandler.go`, `zoneEditorHandler.go` + interface | `FindOpenPosition`/`FindOpenZonePosition` take and return `Vec2`. |
| `test/test_helpers/defaultTemplate.go` | Stopped stamping positions on the golden entity. |
| `.agent/backlog/backlog-opus5.md` | §2.4 rewritten as the record; §8 rows **R** and **S**; tallies, coverage note and §9 baseline refreshed; §2.2's superseded "two seams" note corrected. |
| `.agent/memories/{architecture,generator-domain,template-model}.md` | Updated to the achieved state, incl. the batch S plan. |
| `.agent/plans/batch-r-entity-stops-carrying-editor-state.md` | The plan, complete with Final Recap and Deployment Plan. |

## 5. Tests added or updated

- **Added:** `test/unit/internal/models/template_model/template_variant_model/zone/clone_test.go`
  — the 24 deep-mutation cases **moved** from the deleted wrapper test and retargeted at
  `template_model.Zone.Clone`, plus a `Quality` case the entity could not have.
  `manualZoneSave/toManualZoneSaveEntities_test.go` and `toManualZoneModels_test.go`.
- **Deleted:** `manualZoneSave/{clone,fromManualZoneSaves,toManualZoneSaves,toManualZoneSaveModels}_test.go`.
  `…_SavePositionWins` was not carried over — phase 2 makes it structurally impossible.
- **Retargeted:** the two stamp assertions in `templateGenerator/generate_test.go` now assert on
  `Generate()`'s **model** rather than the lowered entity. This matters: the golden compares
  whole entity templates, so it would have kept passing while silently no longer comparing
  positions at all.
- ~25 further test files moved onto `Vec2`.

**Final gate run, all green:**

| Gate | Result |
| --- | --- |
| `go build ./...`, `go vet` (both tag sets) | exit 0 |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | passed |
| `wire diff ./internal/composition/...` | exit 0 |
| `go test ./test/unit/... -count=1` | pass |
| `go test -tags=integration_test ./test/integration/...` | pass |
| `go test -tags='integration_test,gui' ./test/integration/gui/...` | pass, **no `-update`**, no golden moved |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Coverage | **74.1 %** (was 74.5, floor 72.5) |
| protected trees changed | 1 file / 19 deletions / 0 insertions; `data/` and `internal/registry/` untouched |

Coverage fell 0.4 pp because the deleted `cloneZone` machinery was 100 % covered; removing
fully-covered code from a 74 %-covered tree always lowers the ratio. Every function in every
touched file still reports 100 %. No test was written to move the number back.

## 6. Git status snapshot

Branch **`AD/fixing_some_stuff_08-12`**, head **`7b3ccea Batch Q done`** — the owner committed
batch Q before this session. Batch R is **uncommitted**: 58 modified, 4 deleted, 3 untracked
(`manualZoneSave/toManualZoneModels_test.go`, `…/toManualZoneSaveEntities_test.go`,
`test/unit/internal/models/template_model/template_variant_model/`).

⚠ **Two paths appear STAGED and the agent did not stage them** — no `git add` was ever run this
session. `.agent/plans/batch-r-…md` shows `AM` and `.agent/session-carry-forward.md` shows a
staged deletion `D `. Per AGENTS.md §2.5 the index was left untouched. Check before committing.

## 7. Rejections / corrections

- **The owner banned `goimports` mid-session** ("current tooling is enough to audit
  everything"). One `goimports -w` pass over an explicit 16-file list had already run; every
  import after that was fixed by hand against the compiler. Do not reach for it again.
- Self-corrections worth keeping:
  - **A backlog edit used `| ✅ **Q** | §2.6 step 4 |` as its anchor — a *prefix* of the batch Q
    row — so the replacement swallowed Q's text into the new S row.** Caught by re-listing the
    table rows immediately after the edit, and repaired. Anchor on something that cannot be a
    prefix of a line you intend to keep, and re-read the table after editing it.
  - Phase 1 could not be verified standalone (see §3 of the plan); phases 1 and 2 were gated
    together. The ordering that mattered — carrier fix before field deletion — was preserved.
  - The plan's "move the converters to `internal/mappers`" was overruled by the codebase's
    established seam; recorded in the phase summary rather than silently done differently.

## 8. Open questions

- The staged index entries in §6 — not created by the agent.
- **Manual connections still use the old wrapper shape.** `editor_state_model.ManualConnectionSave`
  wraps the entity exactly as `ManualZoneSave` used to, and `template_model.Connection` already
  carries `IsUserAdded`, so the same collapse applies. Deliberately out of scope (the owner
  scoped zones), but the asymmetry is now visible in the code. Worth a decision.
- Still unreconciled from three sessions ago: two `TabCycling` benchmark baselines disagree
  (~5,699 vs 6,640 allocs/op), taken on different trees.

## 9. Next recommended actions

1. **Review and commit batch R.** Follow the Deployment Plan in
   `.agent/plans/batch-r-entity-stops-carrying-editor-state.md` — in particular step 1 (the
   protected diff is exactly one file), step 2 (the frozen fixtures did not move), and the
   step 3 smoke test: drag a zone, add a zone by hand, regenerate with manual edits, change a
   castle count, Save To, reload a fresh session, and **open an existing `.gen.json` from
   `output/`** (it must still load — R does not change the format).
2. Delete the plan file and this file once R lands. **Backlog §2.4 is the surviving record.**
3. **Start batch S.** Full scope below.

### Batch S — editor-state schema v2 and the versioned migration

Every design decision is already made; do not re-ask. They are recorded in backlog §8's row S
and in the plan's "Deferred to batch S" section.

**Goal.** Make the persisted position a `Vec2`, persist the generator stamps, and build the
versioning machinery that lets an old `.gen.json` still load.

**Decided:**

- `editor_state.ManualZoneSave.ManualPosition` becomes `*data.Vec2[float64]`, so
  `manualZones[*].manualPosition` changes from `[0.25,0.75]` to `{"X":0.25,"Y":0.75}`. It is
  the **only** JSON array-of-floats in the 72-key schema.
- v2 **also adds `generatorPosition` and `generatorRing` sidecars** to `ManualZoneSave`, fixing
  a real latent loss: every save silently drops the generator stamps today. They nest under
  `manualZones`, so `persistedEditorStateFieldCount = 72` is **not** affected.
- `CurrentEditorStateSchemaVersion` → **2**.
- Mechanism: **frozen per-version snapshots + typed migrations.** A full frozen copy of the v1
  struct tree lives in **`internal/entities/editor_state/editor_state_v1/`**, never edited
  again, with a typed `MigrateToV2` beside it. Duplication is the feature — a later change to a
  current struct can never silently redefine what v1 meant.
- Hook: **`FileService.LoadSettingsFile`**, between decode and `ToModel`. A file whose
  `schemaVersion` is **newer than current is rejected loudly** with a dedicated error saying the
  app is too old — never a best-effort load, which would discard the newer file's data on the
  next save.
- **`data.Vec2` gets no `MarshalJSON`/`UnmarshalJSON`, ever.** The migration is what makes old
  files load. This was an explicit owner ruling; a tolerant wrapper type was offered and
  declined.
- Delete `toPositionArray`/`fromPositionArray` from `internal/models/editor_state_model/manualZoneSave.go`
  — they exist only to bridge R to S. When they are gone, `[2]float64` should vanish from the
  repository entirely.

**Constraints discovered while scoping R — read these before writing code:**

- The repo decodes with **`encoding/json/v2`**. Feeding `[0.25,0.75]` to a struct field is a
  **hard decode error**, not a silent nil. Without the migration, every existing user
  `.gen.json` — including the real files in `output/` — becomes unloadable with a bare
  "Load failed". Use those files as the acceptance test.
- `FileService.LoadSettingsFile` decodes **into a seeded default entity**
  (`editorStateMapper.NewDefaultEntity()`), so absent keys keep defaults rather than zero
  values. `test/unit/internal/repositories/editorStateRepository/load_test.go` pins this. The
  migration must not fight it.
- Writes go through `atomicFileWriter.WriteJSON` with fixed options
  (`jsontext.WithIndent("  ")`, `OmitEmptyWithLegacySemantics(true)`,
  `FormatNilSliceAsNull(true)`). Do not introduce a second write path.
- **Two fixture tests cannot survive unchanged, and this is the one sanctioned exception to the
  frozen-fixture rule:** `TestWhenTheLegacyStateFixtureIsLoaded_…` and
  `TestWhenTheCurrentStateFixtureIsLoaded_…` in
  `test/integration/editorStateWireFormat_integration_test.go` `json.Unmarshal` straight into
  `editor_state.EditorState`. They must be rewritten to decode **through the migration**. The
  fixture **files** `editorState_v0_flat.gen.json` and `editorState_v1_flat.gen.json` stay
  byte-frozen as legacy migration inputs; add `editorState_v2_flat.gen.json` beside them. Keep
  comparing parsed objects, never bytes.
- v0 carries **no** `schemaVersion` key at all and decodes at 0, so the migration gate is
  `version < 2`, not `version == 1`.
- `test/test_helpers/allFieldsEditorState.go` is the compile-time blast-radius point: the
  fixtures are the marshaled form of `NewAllFieldsEditorStateEntity()`.
- `SchemaVersion` is excluded from `EqualsIgnoringManualEdits`' change detection. Keep that
  exemption for anything version-related you add.
- **Prior art: none.** The repo has never transformed a decoded old shape — v0→v1 only added a
  key. The dual-fixture pattern is the template to extend.

## 10. Carry-forward prompt

> Read `AGENTS.md` first. The hard rules, one line each: never modify `data/`,
> `internal/registry/` or anything under `internal/entities/template/` **without explicit owner
> approval** — `internal/entities/editor_state/` is *not* protected and batch S edits it
> heavily; everything must build and run on Windows and Linux (`path/filepath`; chain
> PowerShell with `;`, never `&&`); every change ships with tests and unit coverage must not
> drop below 72.5 % (currently **74.1 %**), lint baseline **0 issues**; **never stage and never
> commit** — `Move-Item` not `git mv`, `Remove-Item` not `git rm`; never change where
> `.rmg.json` is written and never persist the output directory; never run a bulk in-place
> rewrite and **never round-trip a `.go` file through `Get-Content`/`Set-Content`**. The owner
> has **banned `goimports`** — fix imports by hand against the compiler.
>
> **Batch R is COMPLETE and uncommitted** (58 modified, 3 untracked, 4 deleted; batch Q is
> committed at `7b3ccea` on `AD/fixing_some_stuff_08-12`). ⚠ Two `.agent/` paths show as
> **staged** and no agent staged them — do not touch the index.
>
> What batch R changed that later sessions must know: `EditorState.ManualZones` is
> `[]template_model.Zone` — the `editor_state_model.ManualZoneSave` wrapper and the
> hand-written entity deep-copy are **gone**, `template_model.Zone.Clone` owns copying;
> `GeneratorPosition`, `GeneratorRing` and `ManualPosition` **no longer exist on
> `entities.Zone`** (one approved protected edit, 19 deletions); positions are
> `*data.Vec2[float64]` everywhere above the persistence seam, spelled `data.Vec2[float64]`
> and never `models.Position`. `.gen.json` is **byte-identical** — the persisted sidecar is
> still `*[2]float64` behind `toPositionArray`/`fromPositionArray`, which batch S deletes.
>
> **Next up: batch S** — editor-state schema v2 plus a versioned migration. Every decision is
> already made; see §9 of `./.agent/session-carry-forward.md` for the full brief, and backlog
> §8 row S. The short version: persisted `ManualPosition` becomes a `Vec2` (object on the
> wire), `generatorPosition`/`generatorRing` gain sidecars, `CurrentEditorStateSchemaVersion`
> → 2, migration via a **full frozen v1 snapshot** in
> `internal/entities/editor_state/editor_state_v1/` with a typed `MigrateToV2`, hooked in
> `FileService.LoadSettingsFile`, rejecting newer-than-current loudly. **`data.Vec2` never gets
> a marshaller** — owner ruling. The gate is `schemaVersion < 2` because v0 has no key at all.
>
> Standing traps: **nil is load-bearing** (nil `Previous` = first generation, nil `Next` =
> unarmed debounce, nil `Zone.Quality` = infer, nil position = never stamped); the persisted
> tier is `*int8`; **the model has no JSON tags — never `json.Unmarshal` into it, including in
> test code** (batch R found two tests doing exactly that); the repo decodes with
> `encoding/json/v2`, which **hard-fails** an array into a struct, so schema v2 without its
> migration makes every existing user `.gen.json` unloadable — test against the real files in
> `output/`; `cmd/testlayoutcheck` matches test-only export names tree-wide; a file gets
> `//go:build integration_test` **only** if it calls a `*_testexports.go` accessor;
> `helpers.MapSlice`/`MapPointer` preserve nil-vs-empty; `golangci-lint --fix` wraps as
> `param,\n) Ret {` where house style is `param) Ret {`.
>
> Lessons from R: when editing a Markdown table, **never anchor on a string that is a prefix of
> a row you intend to keep** — one such edit swallowed a neighbouring row — and re-read the
> table right after. Check `settled-decisions.md` and `architecture.md` before restating any
> "permanent" claim. The next free batch letter after S is **T**.
