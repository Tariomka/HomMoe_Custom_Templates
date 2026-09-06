# Batch S — editor-state schema v2 and the versioned migration

Move the persisted editor state to schema **v2**: `manualZones[*].manualPosition` becomes a
`Vec2` object, `generatorPosition` / `generatorRing` gain sidecars that today are silently
dropped on every save, and a frozen-snapshot migration keeps every existing `.gen.json`
loadable. `.rmg.json` is not touched.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

## Decisions already settled — do not re-ask

Taken by the owner on 2026-09-05 (batch R hand-off) and 2026-09-05 (this session):

- `editor_state.ManualZoneSave.ManualPosition` becomes `*data.Vec2[float64]`. On the wire
  `[0.25,0.75]` becomes `{"X":0.25,"Y":0.75}`. It is the only JSON array-of-floats in the schema.
- v2 also persists `generatorPosition` and `generatorRing` on `ManualZoneSave`. They nest under
  `manualZones`, so `persistedEditorStateFieldCount = 72` is **unchanged**.
- `CurrentEditorStateSchemaVersion` → **2**.
- Mechanism: **frozen per-version snapshots + typed migrations.** Duplication is the feature.
- **The freeze stops at the editor-state boundary** (owner, this session). `editor_state_v1`
  copies the 17 `internal/entities/editor_state` structs and no more; its `ManualZoneSave.Zone`
  and `ManualConnectionSave.Connection` keep naming the live `entities.Zone` / `entities.Connection`.
  Those are the game's `.rmg.json` vocabulary — they version with the game, not with our editor
  state, and copying protected types into an unprotected package cuts against AGENTS.md §2.1.
- **The version sniff-and-dispatch lives in a separate migrator** under the package that owns
  loading: `internal/services/file_service/editor_state_migrator/`. `FileService.LoadSettingsFile`
  seeds the default entity and hands it to `migrator.Load(filePath, &entity)`. The migrator probes
  the file's `schemaVersion` with `os.Open` and then delegates to **one of two repositories**:
  `IFileRepository[editor_state.EditorState]` for v2, or
  `IFileRepository[editor_state_v1.EditorState]` plus `MigrateToV2` for anything older.
  `EditorStateRepository` is unchanged from before batch S. Rationale: json/v2 hard-fails a v1
  payload decoded into the v2 entity, so the version has to be read before any decode - but that
  is a *dispatch* decision, not a reason to change what a repository is.
- A file whose `schemaVersion` is **newer than current is rejected loudly** with a dedicated
  error. Never a best-effort load, which would silently discard the newer file's data on the
  next save.
- **`data.Vec2` gets no `MarshalJSON`/`UnmarshalJSON`, ever.** The migration is what makes old
  files load.
- The gate is `schemaVersion < 2`, not `== 1`: v0 carries no `schemaVersion` key and decodes at 0,
  and its `manualPosition` is the same array shape as v1's.
- `MigrateToV2` stamps `SchemaVersion = CurrentEditorStateSchemaVersion` on its result.
- The frozen structs are **data only** — `MigrateToV2` lives in the service, not beside them, so
  `internal/entities/` keeps carrying no logic (AGENTS.md §4.4.1).

## Constraints

- **The seeded default is load-bearing.** `LoadSettingsFile` seeds
  `editorStateMapper.NewDefaultEntity()` so an absent key keeps its default instead of collapsing
  to a zero value; `NewDefaultEditorStateModel` seeds `PlayerZoneContentRows` among ~40 scalars.
  The legacy path needs a **v1-shaped seed**, built by down-converting the default entity.
  `test/unit/internal/repositories/editorStateRepository/load_test.go` pins the repository half of
  this behaviour.
- Writes go through `atomicFileWriter.WriteJSON` with fixed options (`jsontext.WithIndent("  ")`,
  `OmitEmptyWithLegacySemantics(true)`, `FormatNilSliceAsNull(true)`). Do not add a second write
  path. Under legacy semantics a non-nil pointer is never "empty", so `&Vec2{0,0}` and `new(0)`
  are written — which is what keeps "nil = never stamped" distinguishable from "stamped at origin".
- **Two fixture tests cannot survive unchanged**, the one sanctioned exception to the
  frozen-fixture rule: `TestWhenTheLegacyStateFixtureIsLoaded_…` and
  `TestWhenTheCurrentStateFixtureIsLoaded_…` in
  [test/integration/editorStateWireFormat_integration_test.go](../../test/integration/editorStateWireFormat_integration_test.go).
  They `json.Unmarshal` straight into `editor_state.EditorState`; they must decode through the
  migration. `editorState_v0_flat.gen.json` and `editorState_v1_flat.gen.json` stay **byte-frozen**
  as migration inputs; add `editorState_v2_flat.gen.json` beside them. Compare parsed objects,
  never bytes.
- `test/test_helpers/allFieldsEditorState.go` is the blast radius: the fixtures are the marshaled
  form of `NewAllFieldsEditorStateEntity()`.
- `SchemaVersion` is excluded from `EqualsIgnoringManualEdits`. Keep that exemption for anything
  version-related.
- **Prior art: none.** v0→v1 only added a key; the repo has never transformed a decoded old shape.

## Phase 1: Freeze the v1 snapshot

Status: Complete

Done *before* anything changes, so the snapshot is a copy of the shipped v1 shape rather than a
reconstruction of it.

- [x] Create `internal/entities/editor_state/editor_state_v1/` (package `editor_state_v1`) as a
      verbatim copy of the 17 files in `internal/entities/editor_state/`: `bonusEntry`,
      `bonusPresetType`, `castleSettings`, `contentRuleRow`, `contentSettings`, `editorState`,
      `gameRuleSettings`, `generationSettings`, `manualConnectionSave`, `manualEditSettings`,
      `manualZoneSave`, `mapSettings`, `neutralZoneSettings`, `playerSettings`, `schemaOptions`,
      `templateIdentity`, `zoneContentRow`.
- [x] `ManualZoneSave.Zone` and `ManualConnectionSave.Connection` keep naming the live
      `entities.Zone` / `entities.Connection`; nothing else is shared.
- [x] Replace each file's doc comment with a one-line note that the package is a frozen snapshot
      and must never be edited; add a package doc on `editorState.go` saying the same and naming
      the migration.
- [x] Drop `CurrentEditorStateSchemaVersion` from the copy — the constant only makes sense for the
      current version. Add `SchemaVersion = 1` in its place as the version this snapshot describes.

### Verification Plan

- `go build ./...` exits 0; `gofmt -l ./internal` is empty.
- `go test ./test/unit/... -count=1` still passes untouched (nothing references the new package yet).

### Phase Summary

17 files created under `internal/entities/editor_state/editor_state_v1/`, package
`editor_state_v1`. Data only — no logic, so `internal/entities/` keeps its no-logic rule.
`topology.MapTopology` is referenced live rather than copied, consistent with the
editor-state-boundary ruling and required for the struct conversions phase 3 relies on.
`BonusPresetType`'s ordinals are re-declared rather than aliased, because those ordinals are the
wire values a v1 file carries.

The snapshot's `SchemaVersion` constant is `1`; the current package keeps
`CurrentEditorStateSchemaVersion`.

**Trap for the next agent:** newly created `.go` files land with CRLF on Windows and `gofmt -l`
flags all of them. Fix with `gofmt -w (gofmt -l ./internal)` — an explicit list, per AGENTS.md
§2.6, never a bulk rewrite.

Gates: `go build ./...` 0, `go vet` on the new package 0, `gofmt -l ./internal` empty,
`go test ./test/unit/... -count=1` pass.

## Phase 2: Move the live schema to v2

Status: Complete

- [x] `internal/entities/editor_state/manualZoneSave.go`: `ManualPosition *data.Vec2[float64]`,
      new `GeneratorPosition *data.Vec2[float64] json:"generatorPosition,omitempty"` and
      `GeneratorRing *int json:"generatorRing,omitempty"`. Refresh the doc comment — it still
      claims `entities.Zone` omits the position with `json:"-"`, which batch R made false.
- [x] `CurrentEditorStateSchemaVersion` → 2.
- [x] `internal/models/editor_state_model/manualZoneSave.go`: carry the two stamps in both
      directions; **delete `toPositionArray` / `fromPositionArray`**. Confirm `[2]float64` no
      longer appears anywhere in the repository.
- [x] `test/test_helpers/allFieldsEditorState.go`: stamp `GeneratorPosition` and `GeneratorRing`
      on the fixture zone so the new keys are actually covered.

### Verification Plan

- `go build ./...` exits 0.
- `grep` for `[2]float64` across `app/`, `internal/`, `test/` returns nothing.
- The wire-format fixture tests are expected to FAIL here (v1 fixtures against a v2 decoder);
  phase 4 is what fixes them. Record which ones fail so phase 4 can confirm it fixed exactly those.

### Phase Summary

The live `ManualZoneSave` now reads `Zone` / `GeneratorPosition` / `GeneratorRing` /
`ManualPosition` / `Quality`, all pointers but the zone. The converters use
`helpers.ClonePointer` in both directions, so the save is a real snapshot rather than an alias of
the live zone's pointer — there is a test for that.

`[2]float64` now survives in exactly one place, `editor_state_v1/manualZoneSave.go`, which is the
frozen shape by design. `toPositionArray` / `fromPositionArray` are gone.

Tests: three converter tests retargeted at `Vec2`; nine added covering the stamps in both
directions, including ring 0 (a pointer to the innermost ring is not "unstamped") and the
snapshot-not-alias property. `test/integration/editorState_integration_test.go` moved off array
indexing onto `.X` / `.Y`.

**Expected failures at the end of this phase — exactly four**, all in
`test/integration/editorStateWireFormat_integration_test.go` plus its neighbour:
`TestWhenTheCurrentStateFixtureIsLoaded_…`, `TestWhenTheLegacyStateFixtureIsLoaded_…`,
`TestWhenTheLegacyStateFixtureIsMappedToAState_…` and
`TestWhenTheAllFieldsStateIsWritten_ItMatchesTheCurrentFixture`. Phase 4 must fix exactly these
and no others. The third one is a useful confirmation that json/v2 really does hard-fail an
array into a struct — it is a bare `json.Unmarshal` of the v0 fixture.

Gates: `go build ./...` 0, `go vet ./...` 0, `gofmt -l` empty, `go test ./test/unit/...` pass.

## Phase 3: The migrator and its wiring

Status: Complete

- [x] `internal/repositories/legacyEditorStateRepository.go`: `LegacyEditorStateRepository`
      implementing `IFileRepository[editor_state_v1.EditorState]`. `Load` decodes the frozen v1
      shape; `Save` refuses, because writing a v1 file would produce something this build could
      not round-trip. **`EditorStateRepository` is untouched.**
- [x] `internal/services/file_service/editor_state_migrator/editorStateMigrator.go`: probes the
      version with `os.Open` + `json.UnmarshalRead` off a private minimal struct (never the current
      `SchemaOptions`, which is free to move), rejects newer-than-current, then delegates to the
      current or the legacy repository.
- [x] `…/editorStateMigratorInterface.go`: `IEditorStateMigrator`.
- [x] `…/unsupportedSchemaVersionError.go`: typed error naming both versions and telling the user
      to update the app.
- [x] `…/v1ToV2.go`: `MigrateToV2` plus the seed down-conversion. Flat groups cross via Go struct
      conversion, which is the compile-time tripwire the freeze exists for; `ContentSettings` and
      `ManualEditSettings` need explicit mapping because their slice element types are distinct
      named types.
- [x] `FileService`: keeps `IFileRepository[editor_state.EditorState]` for `SaveSettings` and gains
      `IEditorStateMigrator`; `LoadSettingsFile` is seed → `migrator.Load` → map.
- [x] `internal/composition/providerSets.go` + `wire gen ./internal/composition/...`.
- [x] Test doubles for the migrator; the repository doubles are unchanged.

### Verification Plan

- `go build ./...` and `go vet ./...` exit 0 under both tag sets.
- `wire diff ./internal/composition/...` exits 0.
- `go test ./test/unit/... -count=1` passes.

### Phase Summary

`FileService.LoadSettingsFile` is three lines: seed the default entity, `migrator.Load`, map to
model. The migrator owns both repositories and picks one after probing the version. A load opens
the file twice, once for the probe and once for the decode — that is the deliberate cost of
leaving the repositories as plain typed decoders.

**Reworked on owner review.** The first attempt turned `EditorStateRepository` into a byte reader
behind a bespoke `IEditorStateRepository` and put the migration in `internal/services/`. The owner
rejected that shape: a repository stays a typed decoder, so the legacy shape gets its **own**
repository under the same generic `IFileRepository[T]`, and the migrator belongs **inside**
`file_service`, which is the package that owns loading.

**Owner decision also taken this phase.** The enforced layering test refused the new package:
`entityNamerPrefixes` allowed only `internal/repositories/`, `internal/models/`,
`internal/entities/` and `internal/mappers/`, and the allow-list beneath it is annotated "never
add an entry". The owner chose a **permitted prefix** rather than an exception — reading a
`.gen.json` whose shape depends on its own `schemaVersion` is entity work by definition. The
prefix is `internal/services/file_service/editor_state_migrator/`; the allow-list is untouched and
still holds exactly `internal/services/file_service`.

Notes worth carrying:

- The version probe is a private struct in the migrator, not `editor_state.SchemaOptions`. Sharing
  the current struct would let a rename of the version key break every old file at once.
- `newV1Seed` down-converts the seeded default so the legacy path keeps a default for every
  omitted key. It deliberately drops manual edits: defaults never hold any.
- `toManualConnectionSave` and `toContentRuleRow` are plain struct conversions; `BonusEntry` and
  `ZoneContentRow` cannot be, because their nested types are distinct named types. That asymmetry
  is the tripwire working as intended.
- `LegacyEditorStateRepository.Save` returns an error rather than writing. `IFileRepository[T]`
  demands the method; emitting a v1 file would create something this build could not read back.
- **Trap:** a bare `json.Marshal` in a test writes nil slices as `[]`, not `null` — json/v2's
  default. The real writer passes `FormatNilSliceAsNull(true)`. A round-trip test built on bare
  `json.Marshal` therefore fails on every nil slice in `entities.Zone`. Assert on hand-written
  payloads, or go through the repository.

Gates: build 0, `go vet` 0 under no tag / `integration_test` / `integration_test,gui`,
`gofmt -l` empty, `wire diff` 0, `go run ./cmd/testlayoutcheck .` passed,
`go test ./test/unit/... -count=1` pass.

## Phase 4: Fixtures, tests and the gates

Status: Complete

- [x] Add `test/test_helpers/testdata/editorState_v2_flat.gen.json`, written from the updated
      all-fields entity. v0 and v1 stay byte-identical — verify with `git diff --stat`.
- [x] Rewrite the two sanctioned fixture tests to decode through the migration service. Expected
      value is the all-fields entity with the generator stamps cleared, because no v1 file can
      carry them.
- [x] Point `TestWhenTheAllFieldsStateIsWritten_ItMatchesTheCurrentFixture` at the v2 fixture.
- [x] New integration coverage: a v1 file's `manualPosition` array survives as a `Vec2`; a v2 file
      round-trips the generator stamps; a file claiming version 3 is rejected with the typed error;
      **every real `.gen.json` in [output/](../../output/) still loads**.
- [x] Unit tests per AGENTS.md §4.6 layout for every new public function: the service's `Decode`,
      `MigrateToV2`, the error's `Error`, the repository's `Load`, and `FileService.LoadSettingsFile`
      against the two new mocks.
- [x] Full gate run: build, vet, `gofmt -l`, `go run ./cmd/testlayoutcheck .`, `wire diff`, unit,
      integration, GUI integration (**no `-update`**), `golangci-lint-v2 run ./...` at 0 issues,
      coverage ≥ 72.5 %.
- [x] Update `.agent/backlog/backlog-opus5.md` §8 row S and the `.agent/memories/` notes that
      describe the persisted shape (`template-model.md` lines about the `[2]float64` bridge,
      `architecture.md` on the load path and the schema version).

### Verification Plan

- Every gate above green, with coverage recorded before and after.
- `git status --short` shows nothing under `data/` or `internal/registry/`, and no change to
  `internal/entities/template/`.

### Phase Summary

The v2 fixture was generated by a one-shot `main.go` under `tmp/_genfixture/` — a leading
underscore makes the Go tool ignore the directory, so it never joined the build — run through the
**real** `EditorStateRepository.Save`, then deleted. `git status` confirms v0 and v1 did not move.

The four failures phase 2 predicted are the four that phase 4 fixed, and nothing else broke.
`TestWhenTheLegacyStateFixtureIsMappedToAState_…` was replaced rather than repaired: with the
migration in place the state lands at v2 during the decode, so asserting it after `ToModel` no
longer tests anything. `TestWhenTheCurrentStateFixtureIsLoaded_…` now means v2 and the old v1
assertion lives on as `TestWhenTheV1StateFixtureIsLoaded_…`.

Two coverage holes found and closed: `toV1BonusEntry` was unreachable because the default seed
carries no bonuses (now driven by a seeded target), and the legacy-decode error branch needed a
file that passes the version probe but fails the v1 decode.

After the phase 3 rework the migrator's tests write real temp files instead of passing byte
slices, because `Load` takes a path now. `test/unit/internal/repositories/editorStateRepository/`
is back to its pre-batch-S content, and `…/legacyEditorStateRepository/` mirrors it.

Lint notes for whoever follows:

- Every doc comment in the frozen package had to start with its type name (`revive` + `godoclint`),
  so the "frozen, never edit" line reads `// TypeName is the frozen v1 snapshot - …`.
- `musttag` fires on `json.Marshal` of anything containing a `data.Vec2`, because `Vec2` has no
  json tags. That is deliberate — `X`/`Y` **is** the wire form — so the call carries a nolint
  with that reason.
- Go 1.27 accepts promoted fields as composite-literal keys, and `modernize`'s `embedlit`
  **requires** the flat form. Writing `editor_state.EditorState{ContentSettings: …{…}}` is a lint
  error; write `editor_state.EditorState{Bonuses: …}`.

Final gate run:

| Gate | Result |
| --- | --- |
| `go build ./...` | 0 |
| `go vet` — no tag, `integration_test`, `integration_test,gui` | 0 |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | passed |
| `wire diff ./internal/composition/...` | 0 |
| `go test ./test/unit/... -count=1` | pass |
| `go test -tags=integration_test ./test/integration/...` | pass |
| `go test -tags='integration_test,gui' ./test/integration/gui/...` | pass, **no `-update`**, no golden moved |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Coverage | **74.3 %** (was 74.1, floor 72.5); every function in the new packages 100 % |
| Protected trees | `data/`, `internal/registry/`, `internal/entities/template/` all untouched |

## Final Recap

Batch S gave the editor state a real schema version and the machinery to move between versions,
and used it to fix a silent data loss.

**What changed on disk.** `manualZones[*].manualPosition` went from `[0.25,0.75]` to
`{"X":0.25,"Y":0.75}`; `generatorPosition` and `generatorRing` are written for the first time;
`schemaVersion` is `2`. Nothing else in the 72-key schema moved, and `.rmg.json` was not touched.

**What changed in the code.**

- `internal/entities/editor_state/editor_state_v1/` — 17 frozen files, data only, never to be
  edited. It stops at the editor-state boundary: `ManualZoneSave.Zone` still names the live
  `entities.Zone`, because that type versions with the game rather than with this file format.
- `internal/services/file_service/editor_state_migrator/` — the version probe, the dispatch, the
  typed `UnsupportedSchemaVersionError` and `MigrateToV2`. Settings groups cross by plain Go struct
  conversion, so a future field added to a current group breaks `v1ToV2.go` at compile time
  instead of quietly changing what a v1 file means.
- `LegacyEditorStateRepository` is a second `IFileRepository[T]`, for `editor_state_v1.EditorState`.
  `EditorStateRepository` is exactly what it was before batch S. `FileService.LoadSettingsFile` is
  seed → `migrator.Load` → map, and the migrator picks the repository.
- The array⇄Vec2 bridge is gone. `[2]float64` survives in exactly one place, the frozen snapshot.

**The three owner decisions taken during the work**, all recorded above: the freeze boundary; the
migrator's shape and home (a second repository under the same generic interface, and the migrator
inside `file_service` rather than a standalone service that turns the repository into a byte
reader); and adding the migrator to `entityNamerPrefixes` rather than to the allow-list — it is an
entity-speaking layer by nature, not an exception. The allow-list is still exactly
`internal/services/file_service` and still shrinks only.

**Behaviour the user will notice.** Reloading a saved manual layout now restores the generator
stamps, so the preview draws the real layout instead of falling back to a computed one. Opening a
file written by a newer build fails with a message telling them to update, rather than
half-loading and destroying the unknown keys on the next save.

## Deployment Plan

1. **Confirm the protected trees are untouched.** `git diff -- data/ internal/registry/
   internal/entities/template/` must be empty. Batch S needed no protected edit at all.
2. **Confirm the frozen fixtures did not move.** `git status --short -- test/test_helpers/testdata/`
   must show only the new `editorState_v2_flat.gen.json` as untracked. If `_v0_` or `_v1_` appears,
   stop — they are migration inputs and must stay byte-identical.
3. **Re-run the gates** in the table above, including the GUI suite **without `-update`**.
4. **Smoke test in the app** — this is the half the suites cannot cover, because it is about real
   files on the user's disk:
   - Open an **existing pre-v2 `.gen.json`** from [output/](../../output/). It must load. Check a
     handful of settings against what the file says.
   - Open the zone editor on it, drag a zone, Apply, **Save To**, then reload in a fresh session:
     the manual layout comes back **and the preview is the real layout, not a scatter** — that is
     the generator stamps surviving, which is new.
   - Look at the saved file: `manualPosition` is an object, `generatorPosition` and
     `generatorRing` are present, `schemaVersion` is 2.
   - Hand-edit a copy to `"schemaVersion": 3` and open it: the load must fail with the
     "saved by a newer version of the editor" message, and the session must be left alone.
   - Generate, change a castle count, regenerate: counts update, manual positions survive.
5. **One-way door to call out in the commit message.** A file saved by this build cannot be opened
   by any earlier build — the older code has no migration and will hard-fail on the object-shaped
   `manualPosition`. That is the intended cost of versioning; it is why the refusal is loud.
6. Commit as batch S. Then delete this plan file — **backlog §2.4 and §8 row S are the surviving
   record**, and both are written to stand alone.
