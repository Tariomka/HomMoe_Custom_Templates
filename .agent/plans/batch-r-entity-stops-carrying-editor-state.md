# Batch R — the entity stops carrying editor state, positions become `Vec2`

Backlog §2.4, re-scoped by the owner on 2026-09-05. Delete the three never-serialized
position fields from the protected `template_variant.Zone`, remove the last in-memory
model→entity→model round trip that keeps them alive, and replace `*[2]float64` with
`*data.Vec2[float64]` everywhere above the persistence seam.

The wire-format half (`.gen.json` schema v2, the versioned migration mechanism, and the
generator-stamp sidecars) is **deliberately deferred to batch S** — see the last section,
which is the complete record of every decision already made about it.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done, set its
status to `Complete` and write its **Phase Summary** (what was done, key decisions, anything
needed to continue with zero context); run the phase's **Verification Plan** and record the
result before moving on. When all phases are done, fill in **Final Recap** and **Deployment
Plan**.

### Owner decisions already taken — do not relitigate

1. **Protected edit is APPROVED**, once, for exactly this: deleting `GeneratorPosition`,
   `GeneratorRing` and `ManualPosition` from
   [internal/entities/template/template_variant/zone.go](../../internal/entities/template/template_variant/zone.go#L7-L25).
   Nothing else in `internal/entities/template/`, `internal/registry/` or `data/` may move.
2. **Spelling is `data.Vec2[float64]` everywhere**, not `models.Position`, even though the
   latter is an alias for the same type.
3. **Positions stay pointers.** `nil` means "never stamped" and is load-bearing in the
   preview layout branches. `Vec2`'s zero value `{0,0}` is a legitimate position and can
   never mean absent.
4. **The model wrapper collapses.** `editor_state_model.ManualZoneSave` goes away;
   `EditorState.ManualZones` becomes `[]template_model.Zone`. The `manualPosition` and
   `quality` sidecars survive **only** on the entity, built in the mapper, because
   `entities.Zone` cannot carry them.
5. **`data.Vec2` gets no `MarshalJSON`/`UnmarshalJSON`, ever.** It stays pure arithmetic.

### The one behavioural risk, and why it is neutral

Today a generator stamp survives a regenerate-with-manual-edits **only because
`entities.Zone` carries it** through `SetManualEdits` → `ManualZones` → `GetManualZones`
([manualZoneSave.go](../../internal/models/editor_state_model/manualZoneSave.go#L31-L47)).
Phase 1 replaces that carrier with a model zone **before** phase 2 deletes the fields, so
in-session behaviour is preserved exactly. Do not reorder those two phases. Across a
`.gen.json` save/load the stamps are already lost today (`json:"-"`), and batch S is what
fixes that.

---

## Phase 1: The model stops lowering zones to entities

Status: Complete

- [ ] `editor_state_model.EditorState.ManualZones`: `[]ManualZoneSave` → `[]template_model.Zone`
      ([editorState.go](../../internal/models/editor_state_model/editorState.go)).
- [ ] Delete the `editor_state_model.ManualZoneSave` wrapper and its file
      [manualZoneSave.go](../../internal/models/editor_state_model/manualZoneSave.go):
      `ToManualZoneSaves`, `FromManualZoneSaves`, `ToManualZoneSaveModels`,
      `ToManualZoneSaveEntities`, `Clone`, `cloneZone`.
- [ ] Move the entity conversion (zone ⇄ `editor_state.ManualZoneSave`, carrying the
      `ManualPosition` and `Quality` sidecars, including `toQualityOrdinal` /
      `fromQualityOrdinal`) into [internal/mappers/editorStateMapper.go](../../internal/mappers/editorStateMapper.go),
      which is the sanctioned Model⇄Entity seam (AGENTS.md §4.4.1 rule 4).
- [ ] `EditorState.Clone` clones manual zones with `helpers.MapSlice(zones, template_model.Zone.Clone)`.
      **`cloneZone` and its "a field added to entities.Zone must be added here" comment die
      with it** — the model zone owns its own deep copy.
- [ ] Confirm `EqualsIgnoringManualEdits` still excludes `ManualZones` and needs no change.
- [ ] Check every caller: `app/gui/models/editorState.go` (`SetManualEdits`,
      `GetManualZones`), `app/gui/drivers/stateManualEdits.go`, `internal/dtos/`,
      `internal/validators/`, and the `ManualEditSettings` mapping either side.
- [ ] Grep the rewritten files for `Unmarshal(` / `Decode(` — a type change must never land
      on a JSON decode target (the lesson from batch Q).

### Verification Plan

- `go build ./...` and `go vet -tags='integration_test,gui' ./...` exit 0.
- `go test ./test/unit/... -count=1` and `go test -tags=integration_test ./test/integration/... -count=1` pass.
- `test/integration/editorStateWireFormat_integration_test.go` passes **unchanged** — phase 1
  must not move a single byte of `.gen.json`.
- Mutation check: break the sidecar mapping in `editorStateMapper` and confirm
  `manualZoneTierPersistence_integration_test.go` fails.

### Phase Summary

**Done, but NOT verified standalone — phases 1 and 2 were verified together, deliberately.**
Phase 1 cannot be green on its own: `ToZoneEntity` copies the model's `ManualPosition`
into the entity's `json:"-"` field, so `ToManualEditSettingsEntity(ToManualEditSettingsModel(x))`
stops being identity and the all-fields wire fixture stops matching. Phase 2 removes the
cause. The *order* the plan cared about was preserved — the carrier fix landed first — only
the verification gate moved. Do not read this as licence to merge phases generally.

What changed:

- `editor_state_model.ManualZoneSave` **deleted**. `ManualEditSettings.ManualZones` and
  `EditorState.ManualZones` are now `[]template_model.Zone`.
- Four converters became two, in the same file: `ToManualZoneSaveEntities([]template_model.Zone)`
  and `ToManualZoneModels([]editor_state.ManualZoneSave)`. The old wrap/unwrap pair
  (`ToManualZoneSaveModels` / the old `ToManualZoneSaveEntities`) had no job left once the
  model stopped being a wrapper.
- **Decision changed from the plan:** the converters stayed in `editor_state_model` instead of
  moving to `internal/mappers/editorStateMapper.go`. Moving them would have split the pair
  from its identical sibling `manualConnectionSave.go` and from `ToManualEditSettingsModel`/
  `Entity`, which live in the model package and call them. The repo's established seam for
  Model⇄Entity is the model package's `ToXModel`/`ToXEntity` functions, with `internal/mappers`
  orchestrating; that pattern was followed rather than broken.
- `cloneZone`, `cloneMainObject`, `cloneRoad`, `cloneTypedRef` **deleted** — ~60 lines of
  hand-maintained entity deep-copy, along with the "a field added to entities.Zone must be
  added to cloneZone as well" comment that made it a standing maintenance hazard.
  `EditorState.Clone` now calls `template_model.Zone.Clone`, which the model owns.
- `app/gui/models/editorState.go`: `SetManualEdits` and `GetManualZones` no longer convert at
  all. They `slices.Clone`, which is the faithful equivalent of the old behaviour (fresh slice,
  fresh zone values, inner references shared) — deliberately not a deep clone, to avoid
  smuggling a behaviour change into a structural batch.

Tests: the old `manualZoneSave` folder's four files became two
(`toManualZoneSaveEntities_test.go`, `toManualZoneModels_test.go`). The 24 deep-mutation cases
from the deleted `clone_test.go` were **moved, not dropped**, to the correctly located
`test/unit/internal/models/template_model/template_variant_model/zone/clone_test.go` and
retargeted at `template_model.Zone.Clone`, plus a new `Quality` case the entity could not have.
`TestWhenSavePositionDiffersFromEmbeddedZonePosition_SavePositionWins` was **not** carried over:
it pinned the sidecar beating the embedded entity value, which phase 2 makes structurally
impossible.

## Phase 2: Delete the three fields from the protected entity

Status: Complete

- [ ] Remove `GeneratorPosition`, `GeneratorRing`, `ManualPosition` and their doc comments
      from [zone.go](../../internal/entities/template/template_variant/zone.go#L7-L25).
      **This is the one approved protected edit; verify nothing else in the file moved.**
- [ ] Drop the three copy lines from each of `ToZoneModel` and `ToZoneEntity`
      ([zone.go](../../internal/models/template_model/template_variant_model/zone.go#L86-L128)).
- [ ] [test/test_helpers/defaultTemplate.go](../../test/test_helpers/defaultTemplate.go#L31-L34)
      stamps the entity. **First determine whether `GetDefaultTemplateModel()` consumers
      depend on those stamps**; if they do, stamp on the model side in
      [defaultTemplateModel.go](../../test/test_helpers/defaultTemplateModel.go) instead of
      deleting outright.
- [ ] Rewrite the two stamp assertions in
      [generate_test.go](../../test/unit/internal/services/template_generator/templateGenerator/generate_test.go#L243-L291)
      against `Generate()`'s model output instead of the entity `generateTemplate` lowers to.
- [ ] Confirm the golden test `TestWhenDefaultConfiguration_ReturnsGoldenTemplate` still
      passes. It compares whole entity templates, so deleted fields simply stop being
      compared — **it will not fail on its own**; the assertion rewrite above is what keeps
      the coverage honest.

### Verification Plan

- `git diff --stat internal/entities/` names **exactly one file**, and `git diff data/ internal/registry/` is empty.
- `go build ./...`, full unit + integration suites pass.
- Mutation check: re-add one field to the entity and confirm no test starts depending on it.

### Phase Summary

The approved protected edit landed and nothing else in the tree moved:
`git diff --stat -- internal/entities/ data/ internal/registry/` reports **one file, 19
deletions, 0 insertions** — `internal/entities/template/template_variant/zone.go`.

- Six copy lines gone from `ToZoneModel`/`ToZoneEntity`.
- `test/test_helpers/defaultTemplate.go` no longer stamps positions onto the golden fixture;
  it did so only to match what the generator's lowered entity carried, which is now nothing.
  Checked all 24 `GetDefaultTemplateModel()` call sites first — every one treats the template
  as opaque, so no consumer needed the stamps moved to the model side.
- The two stamp assertions in `generate_test.go` now call `generator.Generate()` and assert on
  the **model** instead of on the entity `generateTemplate` lowers to. This matters: the golden
  test compares whole entity templates, so it would have kept passing while silently no longer
  comparing positions at all.

Gates after 1+2 together: `go build ./...` and `go vet` clean under both tag sets; unit,
integration and **GPU** suites all pass, the last with no `-update` and no golden moved;
`gofmt -l` empty (the three new test files were formatted from an explicit `gofmt -l` list);
`testlayoutcheck` passes; `wire diff` exit 0; `golangci-lint-v2` **0 issues**.

**Coverage 74.5 % → 74.1 %** (floor 72.5 %). This is the drop the plan predicted and it is not
a coverage hole: every function in every touched file reports 100 %. Removing ~60 lines of
fully-covered `cloneZone` machinery from a 74 %-covered codebase necessarily lowers the ratio.
No test was written to move the number back up.

## Phase 3: `*[2]float64` → `*data.Vec2[float64]` above the seam

Status: Complete

- [ ] Model: `template_variant_model.Zone.GeneratorPosition` and `.ManualPosition`
      ([zone.go](../../internal/models/template_model/template_variant_model/zone.go#L15-L16)).
- [ ] Producers: [geometricHubTopology.go](../../internal/services/template_generator/providers/topology/geometricHubTopology.go#L92-L105) (3 sites),
      [positionedTopologyBuilder.go](../../internal/services/template_generator/providers/topology/positionedTopologyBuilder.go#L48),
      [balancedClusterService.go](../../internal/services/template_generator/providers/topology/tournament_variant/balancedClusterService.go#L72-L77),
      `ZoneBuilder.WithGeneratorPosition` ([zoneBuilder.go](../../internal/services/builders/variant_content/zoneBuilder.go#L171)).
      Each currently unpacks a `models.Position` into an array literal — that noise is the
      whole point of the item and must be gone.
- [ ] Consumers: [layoutGeometry.go](../../internal/services/preview_service/layoutGeometry.go#L97-L282),
      [layoutScatter.go](../../internal/services/preview_service/layoutScatter.go#L97-L114),
      [previewLayoutService.go](../../internal/services/preview_service/previewLayoutService.go#L97-L102),
      [layoutFixed.go](../../internal/services/preview_service/layoutFixed.go).
      `p[0]`/`p[1]` become `.X`/`.Y`; **delete the three "Is this required to be copied?"
      comments** — with a value struct the answer is trivially yes and the question is dead.
- [ ] GUI: [zoneEditorCanvas.go](../../app/gui/dialogs/zoneEditorCanvas.go#L325),
      [zoneEditorDialog.go](../../app/gui/dialogs/zoneEditorDialog.go#L442-L489) including
      `manualPositions() [][2]float64`.
- [ ] The open-position API, all five declarations in step:
      `ZoneEditorService.FindOpenPosition`, `zoneEditorServiceInterface`,
      `zoneEditorHandler.FindOpenZonePosition`, `GUIHandler.FindOpenZonePosition`,
      `zoneEditorHandlerInterface`, plus `ZoneEditorServiceMock`, `TemplateHandlerMock`
      and `handlerDependenciesStub`.
- [ ] **Temporary seam:** the mapper converts `*data.Vec2[float64]` ⇄ `*[2]float64` for the
      still-array `editor_state.ManualZoneSave.ManualPosition`. This is the agreed price of a
      committable R; **batch S deletes it.** Mark it with a one-line comment saying so.
- [ ] Tests: ~20 files carry `[2]float64` literals. Do them with an explicit file list, never
      a bulk in-place rewrite, and verify insertions == deletions per file.

### Verification Plan

- `grep` for `[2]float64` returns hits **only** in `internal/entities/editor_state/`,
  `internal/mappers/` (the temporary seam) and the tests that pin the wire format.
- Full unit + integration suites pass.
- GPU suite `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`
  passes **with no `-update`**. Preview geometry is float-native since batch G, so no pixel
  may move; a moved golden means a real regression, not a rebaseline.
- `editorStateWireFormat_integration_test.go` still passes unchanged.

### Phase Summary

Type-only, ~35 files, no behaviour change and no coverage movement (74.1 % before and after).
The method was to retype the two model fields first and let the compiler enumerate every
consumer, rather than grep and guess.

The audit the plan asked for: `[2]float64` now appears in **9 places only** — the persisted
`editor_state.ManualZoneSave.ManualPosition`, the three lines of the temporary bridge, and the
five test lines that pin that seam. Everything above it is `data.Vec2[float64]`.

Worth noting:

- Several sites got shorter rather than merely different, because `Vec2` has arithmetic:
  `math.Hypot(pi[0]-pj[0], pi[1]-pj[1])` became `a.Subtract(b).Distance()`, and the
  hexagon-angle helper in `geometricHubTopology/common_test.go` lost four temporaries to
  `DotProduct`. The three "Is this required to be copied?" comments are gone.
- `balancedClusterService` stamps `new(position)` **twice** rather than sharing one pointer
  between `GeneratorPosition` and `ManualPosition` — the old array literals were two distinct
  allocations and aliasing them would have been a silent behaviour change.
- **A real trap, and the reason this phase is worth reviewing carefully:** two integration
  tests `json.Unmarshal`ed the on-disk file into `editor_state_model.EditorState`. That only
  ever worked because the model embedded the *tagged* entity groups; the moment `ManualZones`
  became `[]template_model.Zone` the decoder matched `manualZones` by name and then failed on
  `cannot unmarshal array into … data.Vec2[float64]`. Both now decode `editor_state.EditorState`.
  This is the standing "never `json.Unmarshal` into the model" rule biting from a new angle:
  it was violated in *test* code, where nothing was checking.
- Imports were fixed by hand against the compiler on the owner's instruction; an earlier
  `goimports` pass over an explicit 16-file list was the last use of that tool.

## Phase 4: Gates, coverage and records

Status: Complete

- [ ] `gofmt -l ./app ./internal ./test ./cmd` empty; `go run ./cmd/testlayoutcheck .` passes;
      `wire diff ./internal/composition/...` clean (regenerate if a constructor signature moved).
- [ ] `golangci-lint-v2 run ./...` at **0 issues**. Watch for `--fix` rewrapping signatures as
      `param,\n) Ret {` where house style is `param) Ret {`.
- [ ] Coverage ≥ **72.5 %**, and account for the delta: deleting `cloneZone` and the
      `ManualZoneSave` wrapper removes covered lines, so the figure can move in either
      direction. Baseline before starting is **74.5 %**.
- [ ] Backlog: rewrite §2.4 as the record of what was actually done, add row **R** to §8,
      update the coverage note, and fix the two drifts found while scoping —
      §9 still says 74.3 %, and §2.2's "two seams keep the entity forever" note was
      superseded by batch Q.
- [ ] `.agent/memories/`: update `architecture.md`, `template-model.md` and
      `settled-decisions.md`.

### Verification Plan

- Every gate in backlog §9 green, recorded as a table in the Final Recap.
- `git status --short` shows nothing staged and nothing committed by the agent.

### Phase Summary

All gates green (table in the Final Recap). Backlog updated: §2.4 rewritten as the record,
row **R** added to §8 with a row **S** beside it carrying the deferred scope, coverage note
and §9 baseline moved to 74.1 %, header tally to 17 done / 4 open. The two drifts found while
scoping were fixed: §9 had been left at 74.3 % since batch P, and §2.2's "two seams keep the
entity forever — do not sweep those two by reflex" now carries a superseded note, because
batch Q swept one of them and batch R deleted the fields the other paragraph described.

**A mistake worth recording:** the §8 edit used `| ✅ **Q** | §2.6 step 4 |` as its anchor, which
is a *prefix* of the batch Q row, so the replacement swallowed Q's text into the new S row.
Caught by re-listing the table rows immediately after the edit and repaired. Anchor on
something that cannot be a prefix of the line you are keeping.

## Final Recap

The entity no longer carries editor state, and positions are vectors everywhere above the
persistence seam. Backlog §2.4 asked for a type swap; the owner re-scoped it, correctly, to the
structural problem underneath — `entities.Zone` was the in-memory *carrier* for three fields it
could never serialize.

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | exit 0 / exit 0 |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 |
| `go test ./test/unit/... -count=1` | pass |
| `go test -tags=integration_test ./test/integration/... -count=1` | pass |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | pass, **no `-update`**, no golden moved |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Coverage | **74.1 %** (was 74.5, floor 72.5) |
| `git diff --stat -- internal/entities/ data/ internal/registry/` | **1 file, 19 deletions, 0 insertions** |

Net effect: one wrapper type, four converter functions and ~60 lines of hand-maintained entity
deep-copy deleted; three fields gone from the protected schema mirror; `[2]float64` reduced to
nine occurrences, all at the persistence seam. `.gen.json` is byte-identical, which is what
makes this batch safe to land ahead of the schema work.

The coverage figure went **down** 0.4 pp and that is the honest outcome: the deleted `cloneZone`
machinery was 100 % covered, and removing fully-covered code from a 74 %-covered tree
necessarily lowers the ratio. Every function in every touched file still reports 100 %. No test
was written to move the number.

## Deployment Plan

1. **Review the protected edit first.** `git diff -- internal/entities/template/template_variant/zone.go`
   must show three field blocks and their doc comments removed and **nothing else**;
   `git diff -- data/ internal/registry/` must be empty.
2. **Confirm the wire format did not move.** `test/test_helpers/testdata/editorState_v0_flat.gen.json`
   and `editorState_v1_flat.gen.json` are untouched, and
   `go test -tags=integration_test ./test/integration/ -run WireFormat -count=1` passes. Batch R
   must not change a single byte of `.gen.json`; batch S is where that changes.
3. **Smoke test in the app** — this is the part the suites cannot fully cover, because it is
   about the manual-edit snapshot:
   - generate, open the zone editor, **drag a zone**, Apply → the zone stays where you put it;
   - **add a zone by hand** → it lands on open canvas, not on top of another zone
     (`FindOpenPosition` changed shape);
   - regenerate with manual edits present → the layout is reapplied, tier colouring intact;
   - change a castle count and regenerate → counts update, manual positions survive;
   - **Save To**, then load the file in a fresh session → the manual layout comes back and the
     `.rmg.json` + `.png` land in the detected templates directory;
   - open an **existing** `.gen.json` from [output/](../../output/) → it must still load. Batch R
     does not change the format, so this must work today; if it does not, stop.
4. Commit as batch R. Then delete this plan file — **backlog §2.4 is the surviving record** and
   is written to stand alone, and the batch S scope lives in §8's S row plus the section below.

---

## Deferred to batch S — the decisions are already made

Do not re-ask these; the owner settled them on 2026-09-05.

- **Schema v2.** `editor_state.ManualZoneSave.ManualPosition` becomes
  `*data.Vec2[float64]`, so `manualZones[*].manualPosition` changes from `[0.25,0.75]` to
  `{"X":0.25,"Y":0.75}`. It is the **only** JSON array-of-floats in the 72-key schema.
  Bump `CurrentEditorStateSchemaVersion` to 2.
- **Sidecars added.** v2 also persists `generatorPosition` and `generatorRing` on
  `ManualZoneSave`, fixing the stamps that every save silently drops today. They nest under
  `manualZones`, so `persistedEditorStateFieldCount = 72` is **not** affected.
- **Migration mechanism: frozen per-version snapshots + typed migrations.** A full frozen
  copy of the v1 struct tree lives in **`internal/entities/editor_state/editor_state_v1/`**,
  never edited again, with a typed `MigrateToV2` beside it. Duplication is the feature.
- **Hook: `FileService.LoadSettingsFile`**, between decode and `ToModel`. A file whose
  `schemaVersion` is **newer than current is rejected loudly** with a dedicated error saying
  the app is too old — never a best-effort load, which would discard the newer file's data
  on the next save.
- **`data.Vec2` gets no marshaller.** The migration is what makes old files load.
- **Hard constraint discovered while scoping:** the repo decodes with `encoding/json/v2`.
  Feeding `[0.25,0.75]` to a struct field is a **hard decode error**, not a silent nil, so
  without the migration every existing user `.gen.json` — including the real files in
  [output/](../../output/) — becomes unloadable with a bare "Load failed".
- **Two fixture tests cannot survive unchanged**, and this is the one sanctioned exception
  to the frozen-fixture rule: `TestWhenTheLegacyStateFixtureIsLoaded_…` and
  `TestWhenTheCurrentStateFixtureIsLoaded_…`
  ([editorStateWireFormat_integration_test.go](../../test/integration/editorStateWireFormat_integration_test.go#L28-L84))
  `json.Unmarshal` straight into `editor_state.EditorState`. They must be rewritten to
  decode through the migration. The v0 and v1 fixture **files** stay byte-frozen as legacy
  migration inputs; add `editorState_v2_flat.gen.json` alongside them. Keep comparing
  parsed objects, never bytes.
- **Prior art: none.** The repo has never transformed a decoded old shape — the v0→v1 step
  only added a key. The dual-fixture pattern is the template to extend.
