# Batch I — `EditorStateDto` rework (backlog §2.1, folding in §1.5)

Split the 72-field god-DTO into behaviour-free entity groups under a model that
owns its behaviour, make `EditorStateDto` a thin versioned persistence shell,
switch the runtime type from the DTO to the model, and stop the panels
deep-cloning the whole state every frame.

Backlog items closed by this work: [§2.1](../todo/backlog-opus5.md) and
[§1.5](../todo/backlog-opus5.md).

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../AGENTS.md) first. Hard rules that bite in this batch:
never touch `data/`, `internal/entities/template/` or `internal/registry/`;
never stage and never commit; delete with `Remove-Item`, never `git rm`;
never run a bulk in-place rewrite over the repository; chain PowerShell with `;`,
never `&&`; unit coverage must not drop below **72.5 %** (currently 72.8 %).

---

## 0. Design decisions (settled with the owner 2026-08-21 — do not relitigate)

| Decision | Choice |
| --- | --- |
| How groups attach | **Anonymous embedding** of the 9 entity groups into `EditorStateModel`. Field promotion means `state.MapSize` keeps compiling at ~1,000 call sites, and `encoding/json` flattens anonymous embedded structs for free. |
| json tags | Live on the **entity leaf fields**, same as `internal/entities/template/` already does for `.rmg.json`. |
| Entity behaviour | Entities are **strictly behaviour-free value groups**. All logic lives on the model. Per AGENTS.md §4.6 pure data structs need no unit tests. |
| Runtime type | The **model**, everywhere. `EditorStateDto` appears only at the load/save boundary (`internal/repositories`, `internal/services/file_service`). |
| DTO shape | `EditorStateDto { SchemaVersion int; EditorState editor_state_model.EditorStateModel }` — a **named** field, not embedded. |
| Wire shape | Stays **flat**. `MarshalJSON`/`UnmarshalJSON` on the DTO hoist the model's fields using an internal anonymous-embedding alias, so `schemaVersion` sits as a sibling of the existing 72 keys. |
| `schemaVersion` | Always written as `1`. On load, `0` is normalised to `1` through an explicit migration hook. |
| Entity package | `internal/entities/editor_state/` (package `editor_state`). **Note:** AGENTS.md §4.4's table points non-`.rmg.json` data structs at `internal/models/`; the owner decided on 2026-08-11 and re-confirmed on 2026-08-21 to put these under `internal/entities/` instead. That departure is deliberate — do not "fix" it. |
| Model package | `internal/models/editor_state_model/` (package `editor_state_model`) — renamed from `editor_state` to avoid a package-name collision with the entities. |
| DTO package | The three `editorState*Dto.go` files move into the existing `internal/dtos/editor_state_dto/`, alongside its current internals. |
| Content defaults | `DefaultPlayerZoneContentRows` and its helpers move to a new `internal/common/common_zone_contents/`, with `Get`/`get` prefixes, mirroring [neutralZoneProfile.go](../internal/common/common_zones/neutralZoneProfile.go). |
| §1.5 | **In scope**, as its own phase, using option 3 (per-panel view structs) in `app/gui/models/`. Acceptance: any measurable allocation improvement, recorded here. |
| Round-trip guard | A committed golden `.gen.json` with every field populated, under `test/test_helpers/testdata/`. |

### 0.1 Field → group map (72 fields, 9 groups)

Verify the count is exactly 72 against
[editorStateDto.go](../internal/dtos/editorStateDto.go#L15-L94) before writing
any struct; a prior automated census reported 76 and was wrong.

| Group (file) | Fields |
| --- | --- |
| `templateIdentity.go` (2) | `TemplateName`, `GameMode` |
| `mapSettings.go` (2) | `MapSize`, `ExperimentalMapSizes` |
| `playerSettings.go` (4) | `PlayerCount`, `HeroCountMin`, `HeroCountMax`, `HeroCountIncrement` |
| `neutralZoneSettings.go` (11) | `NeutralZoneCount`, `SpawnAbandonedOutposts`, `AbandonedOutpostCount`, `NeutralLowestNoCastleCount`, `NeutralLowestCastleCount`, `NeutralLowNoCastleCount`, `NeutralLowCastleCount`, `NeutralMediumNoCastleCount`, `NeutralMediumCastleCount`, `NeutralHighNoCastleCount`, `NeutralHighCastleCount` |
| `castleSettings.go` (10) | `AdvancedMode`, `PlayerOwnedCastles`, `PlayerZoneCastles`, `NeutralZoneCastles`, `HubZoneCastles`, `NeutralLowestCastlesPerZone`, `NeutralLowCastlesPerZone`, `NeutralMediumCastlesPerZone`, `NeutralHighCastlesPerZone`, `MatchPlayerCastleFactions` |
| `generationSettings.go` (15) | `PlayerZoneSize`, `NeutralZoneSize`, `HubZoneSize`, `GuardRandomization`, `Topology`, `RandomPortals`, `MaxPortalConnections`, `SpawnRemoteFootholds`, `RemoteFootholdCount`, `GenerateRoads`, `NoDirectPlayerConn`, `ResourceDensityPercent`, `StructureDensityPercent`, `NeutralStackStrengthPercent`, `BorderGuardStrengthPercent` |
| `gameRuleSettings.go` (16) | `VictoryCondition`, `FactionLawXpPercent`, `AstrologyXpPercent`, `LostStartCity`, `LostStartCityDay`, `LostStartHero`, `CityHold`, `CityHoldDays`, `GladiatorArena`, `GladiatorArenaDaysDelayStart`, `GladiatorArenaCountDay`, `Tournament`, `TournamentFirstTournamentDay`, `TournamentInterval`, `TournamentPointsToWin`, `TournamentSaveArmy` |
| `contentSettings.go` (10) | `BannedItems`, `BannedMagics`, `ValueOverridesText`, `Bonuses`, `PlayerZoneContentRows`, `LowestNeutralContentRows`, `LowNeutralContentRows`, `MediumNeutralContentRows`, `HighNeutralContentRows`, `HubZoneContentRows` |
| `manualEditSettings.go` (2) | `ManualZones`, `ManualConnections` |

Owner rulings folded in: `AdvancedMode` sits with `castleSettings` because it
gates castle diffing; all three zone sizes sit in `generationSettings`; there is
**no** `hubZoneSettings` group — `HubZoneSize` → `generationSettings`,
`HubZoneCastles` → `castleSettings`.

### 0.2 Phase ordering rationale

The build stays green at every phase boundary. Phase 3 has the DTO embed the
model **anonymously and temporarily**, so field promotion survives and no
consumer's *field accesses* change; its promoted **methods** do change shape,
which is why Phase 3 also keeps DTO-signature shims (hazard 1). Only after
Phase 4 moves every consumer onto the model does Phase 5 give the DTO its named
field, at which point nothing reads DTO fields any more.

### 0.3 Known hazards

Items 1–6 were found by an independent review of this plan's first draft and
had each falsified something it originally claimed. Treat them as load-bearing.

1. **Promoted methods change signatures, so Phase 3 needs shims.** With the model
   embedded, `dto.Clone()` resolves to the model's method and returns an
   `EditorStateModel`, and `EqualsIgnoringManualEdits`/`LayoutDefiningOptionsChanged`/
   `DiffCastleSettings` start demanding `*EditorStateModel`. Today's callers assign
   the clone to a DTO and pass DTO pointers —
   [editorState.go](../app/gui/models/editorState.go#L26) and
   [regenerationDecisionService.go](../internal/services/editor/regenerationDecisionService.go#L29).
   Phase 3 must therefore keep **explicit DTO-signature shim methods** that shadow
   the promoted ones and delegate inward; Phase 4 deletes them.
2. **`omitempty` already conflates nil and empty — do not "restore" a
   distinction that was never persisted.** All six content-row slices and both
   manual-edit slices carry `omitempty`
   ([editorStateDto.go](../internal/dtos/editorStateDto.go#L85-L94)), so nil and
   `[]` both serialise to an absent key, and an absent key then keeps whatever
   `NewDefaultEditorStateDto` seeded. `EqualsIgnoringManualEdits` distinguishes
   nil from empty **in memory only**. Batch I must preserve this behaviour
   exactly and characterise it in a test — it must not try to make the
   distinction survive the file.
3. **JSON key order will change, so byte-diffing is not a valid gate.**
   `encoding/json` emits fields in declaration/embedding order, and regrouping
   reorders them. The format stays flat with the same key set; only the textual
   order moves. Every round-trip assertion must compare **parsed objects**, never
   bytes. `.gen.json` is only ever read by `json.Unmarshal`, so the reordering is
   cosmetic — but the owner's existing files will be rewritten in the new order
   on their next save.
4. **Phase 1 is not a pure package-clause move.**
   [editorStateDto.go](../internal/dtos/editorStateDto.go#L4) imports
   `editor_state_dto` and qualifies `editor_state_dto.ManualZoneSave` /
   `CastleSettingChanges`; once it lives *in* that package those qualifiers must
   be stripped. Four sibling DTOs in the parent package also embed or carry
   `EditorStateDto` and need the new import —
   [regenerationDecisionRequestDto.go](../internal/dtos/regenerationDecisionRequestDto.go#L13),
   [templateUpdateDto.go](../internal/dtos/templateUpdateDto.go#L13),
   [castleSettingsReapplyRequestDto.go](../internal/dtos/castleSettingsReapplyRequestDto.go#L11)
   and `manualEditDecisionDto.go`.
5. **Carrier DTOs are part of the runtime swap.** Those same four, plus
   `EditorStateSaveDto` and `EditorStateValidationDto`, hold editor state inside
   another DTO. Phase 4 changes what they carry; it is not only an interface
   signature change. `stateHandler` forwards the very value that `file_service`
   persists ([stateHandler.go](../internal/handlers/stateHandler.go#L27)), so the
   model→DTO wrap on save and DTO→model unwrap on load need a named owner.
6. **The two reflection guards need *different* treatment.**
   [clone_test.go](../test/unit/internal/dtos/editorStateDto/clone_test.go#L120-L132)
   already walks recursively through `assertNoSharedStorage` and keeps working
   unchanged. [equalsIgnoringManualEdits_test.go](../test/unit/internal/dtos/editorStateDto/equalsIgnoringManualEdits_test.go#L211-L230)
   enumerates **top-level fields only** and mutates them by index — after
   grouping it will fail loudly on a struct kind rather than fail silently. It
   must be rewritten to enumerate recursive leaf paths, ignore only the two
   manual-edit leaves, exclude `SchemaVersion`, and mutate the model rather than
   the DTO shell.
7. **Named literal.** `NewDefaultEditorStateDto` is the only keyed literal with
   ~32 named fields; it becomes a nested literal and is the one place the
   compiler will not help after promotion.
8. **Load seeds defaults.** [editorStateRepository.go](../internal/repositories/editorStateRepository.go#L26-L29)
   unmarshals *over* `NewDefaultEditorStateDto()`, which is how absent keys keep
   their defaults. `UnmarshalJSON` in Phase 5 must preserve that merge semantic —
   it must not zero the receiver first.
9. **`MarshalJSON` recursion.** The alias type inside `MarshalJSON` must be a
   locally declared type that embeds the model — not the DTO — or the marshaller
   calls itself forever.
10. **Import cycle — the shared types cannot stay in `editor_state_dto`.**
    Found while starting Phase 2; the plan's first draft missed it entirely.
    `ManualZoneSave`, `ManualConnectionSave` and `CastleSettingChanges` live in
    `internal/dtos/editor_state_dto/`, but Phase 2's `manualEditSettings.go`
    must hold the first two and Phase 3's model must return the third. Combined
    with the settled DTO shape (`editor_state_dto` → `editor_state_model`) and
    the model → entities edge, that closes a cycle:

    ```
    editor_state_dto ──embeds──▶ editor_state_model ──embeds──▶ entities/editor_state
           ▲                                                             │
           └──────────── ManualZoneSave / CastleSettingChanges ◀─────────┘
    ```

    Only the third edge is removable, so those three types **must** move out.
    **Owner ruling (2026-08-21):** split them by behaviour.
    - The two pure structs (fields + json tags only, no methods) go to
      `internal/entities/editor_state/`, keeping entities behaviour-free per §0.
    - `CastleSettingChanges` (it has `Any()`) and two wrapper models embedding
      the pure structs — `ManualZoneSaveModel`, `ManualConnectionSaveModel`,
      carrying `Clone()` — go to `internal/models/editor_state_model/`.
    - The four `To…`/`From…` converters keep returning the **entity** slice, so
      no assignment site needs an unwrap loop. The six private deep-clone
      helpers for `entities.Zone`/`entities.Connection` stay private in
      `editor_state_model`.

    This runs as **Phase 2 step 0**, before the group extraction, and partly
    undoes Phase 1's consolidation — that is expected, not a regression.
    `internal/entities` must remain the base layer: it imports nothing from
    `models`, `dtos` or `helpers`, and `entities/editor_state` must not either.

---

## Phase 1: Baseline guard + DTO package move
Status: Complete   <!-- Not started | In progress | Complete -->

Freeze today's on-disk format as a committed artefact **before any code changes**,
then perform the pure package move on its own so the noisy import diff never
mixes with a behaviour change.

- [x] Confirm the field count is exactly 72 and record any discrepancy here.
      **72 confirmed, no discrepancy** — counted by hand on
      [editorStateDto.go](../internal/dtos/editorStateDto.go) (24 + 36 + 4 + 6 + 2)
      and re-confirmed mechanically by the fixture's top-level key count assertion.
- [x] Record the **before** coverage number (AGENTS.md §2.3) so every later phase
      has a baseline to compare against. **Baseline = 73.6 %** on a clean tree.
      This supersedes the 72.8 % quoted in the carry-forward (the Go 1.27 batch
      raised it); the 72.5 % floor still governs.
- [x] Write a throwaway program (or a temporary test) that builds an
      `EditorStateDto` with **every** field set to a distinctive non-zero value,
      including all six content-row slices, bonuses, manual zones and manual
      connections, and marshal it with the **current** code.
      Builder kept permanently as
      [allFieldsEditorState.go](../test/test_helpers/allFieldsEditorState.go)
      (`NewAllFieldsEditorStateDto`); the throwaway `cmd/tmpfixturegen` generator
      wrote the fixture through the **real** `EditorStateRepository.Save` and was
      then deleted.
- [x] Add the result as `test/test_helpers/testdata/editorState_v0_flat.gen.json`
      — a new **unstaged** file; do not stage it and do not commit it. Name it
      `_v0_` because it has no `schemaVersion`: it is the legacy-shape fixture
      that Phase 5's migration hook must keep loading.
- [x] Add a golden round-trip test that loads the fixture and asserts every one of
      the 72 fields survives with its distinctive value. Compare **parsed values,
      never bytes** (hazard 3). This test must stay green, unchanged in intent,
      through every later phase.
      Added as
      [editorStateWireFormat_integration_test.go](../test/integration/editorStateWireFormat_integration_test.go)
      — **untagged** (production APIs only, no `*_testexports.go`, no GPU), so it
      runs in a plain `go test ./test/...`. Three gates: fixture unmarshals into a
      **zero** DTO and equals the builder; the fixture has exactly 72 top-level
      keys; marshalling the builder equals the fixture **as parsed maps**.
- [x] Move [editorStateDto.go](../internal/dtos/editorStateDto.go),
      [editorStateSaveDto.go](../internal/dtos/editorStateSaveDto.go) and
      [editorStateValidationDto.go](../internal/dtos/editorStateValidationDto.go)
      into `internal/dtos/editor_state_dto/`.
- [x] Strip the now-invalid self-qualifiers inside the moved file
      (`editor_state_dto.ManualZoneSave` → `ManualZoneSave`, likewise
      `CastleSettingChanges`) and drop the self-import (hazard 4).
      Four sites, exactly as surveyed.
- [x] Add the new import to the **three** sibling DTOs left in `internal/dtos/`
      that carry editor state: `regenerationDecisionRequestDto.go`,
      `templateUpdateDto.go`, `castleSettingsReapplyRequestDto.go`.
      **Plan correction:** `manualEditDecisionDto.go` needs **no** change — it
      only references `*editor_state_dto.CastleSettingChanges` and already
      imports that package. §0.3 hazard 4 lists it in error. The compiler
      confirmed this: after the move it reported exactly these three files.
- [x] Update imports and qualifiers (`dtos.EditorStateDto` →
      `editor_state_dto.EditorStateDto`) across the 12 production packages and the
      test tree. Do this file by file — **no bulk in-place rewrite** (AGENTS.md §2.6).
      **104 files, 509 sites.** Owner approved the mechanism: substitution of the
      five fixed, unambiguous tokens over an **explicit enumerated file list**
      (§2.6's own "explicit list from `gofmt -l`" pattern), never a repo-wide sweep.
- [x] Move the mirrored unit-test folders to
      `test/unit/internal/dtos/editor_state_dto/…` per AGENTS.md §4.6.
      `test/unit/internal/dtos/editorStateDto/` →
      `test/unit/internal/dtos/editor_state_dto/editorStateDto/` (7 files).

### Verification Plan
- `go build ./...` clean; `go vet ./...` and
  `go vet -tags='integration_test,gui' ./...` clean.
- `go run ./cmd/testlayoutcheck .` prints `test-layout check passed`.
- `go test ./test/unit/... -count=1` and
  `go test -tags=integration_test ./test/integration/... -count=1` pass.
- The new golden test passes.
- Coverage report re-run; total ≥ the number recorded above and ≥ 72.5 %.
- Diff review: apart from the fixture and its test, every change must be an
  import, a package clause, a qualifier, or the self-qualifier strip — **no
  behaviour change**.

### Phase Summary
**Complete.** All verification gates green:

| Gate | Result |
| --- | --- |
| `go build ./...` | clean |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test ./test/... -count=1` (untagged, incl. the new golden test) | exit 0 |
| `go test -tags=integration_test ./test/integration/... -count=1` | exit 0 |
| `go test ./test/unit/... -count=1` | exit 0, 164 packages |
| Coverage | **73.7 %** (baseline 73.6 %, floor 72.5 %) |
| `gofmt -l .` | empty |
| `golangci-lint-v2 run ./...` | 50 issues, **all pre-existing `funcorder`**; down from 71 |

What exists now:

- `internal/dtos/editor_state_dto/` owns `editorStateDto.go`, `editorStateSaveDto.go`
  and `editorStateValidationDto.go` alongside the three types that were already
  there. `internal/dtos/` keeps only the carrier DTOs, which now qualify.
- The wire format is frozen by
  [editorState_v0_flat.gen.json](../test/test_helpers/testdata/editorState_v0_flat.gen.json),
  written through the **real** `EditorStateRepository.Save`, and guarded by three
  untagged tests in
  [editorStateWireFormat_integration_test.go](../test/integration/editorStateWireFormat_integration_test.go).
  The builder [allFieldsEditorState.go](../test/test_helpers/allFieldsEditorState.go)
  is permanent — Phase 5 regenerates the `_v1_` fixture from the same builder.

Decisions and gotchas for whoever picks this up:

- **The golden test is untagged on purpose.** It touches production APIs only, so
  it runs in a plain `go test ./test/...`. Do not add `integration_test` to it.
- **Zero-value unmarshal is deliberate.** The fixture is loaded into a *zero*
  `EditorStateDto`, not over `NewDefaultEditorStateDto()`. That is what proves the
  file alone carries all 72 values; loading over defaults would let a dropped
  field pass by inheriting its default. Do not "fix" this to match
  `editorStateRepository.Load`.
- **Every bool in the fixture is `true`** for the same reason — a `false` would
  pass vacuously against a zero struct.
- **Comparisons are parsed objects, never bytes** (hazard 3), so Phase 2–3
  regrouping may reorder keys freely without touching this test.
- The 72-key count is now asserted mechanically, so Phase 2's field→group split
  cannot silently lose or add a key.
- Three files — [state.go](../app/gui/drivers/state.go),
  [editorState.go](../app/gui/models/editorState.go),
  [layoutPanelZones.go](../app/gui/panels/layoutPanelZones.go) — carry extra
  reflow in the diff. They were **already** non-`gofmt`-clean at HEAD (single-line
  function bodies), and the longer `editor_state_dto.` qualifier pushed those
  lines past the `golines` limit, so a reflow was unavoidable. It was applied by
  the project's own *"Go: Run Linter"* `--fix` task, not by an ad-hoc `gofmt -w`.
- `coverage.txt`, `coverage.html` and `lcov.info` are regenerated tracked
  artefacts, expected in the diff.
- Nothing is staged. Nothing is committed.

---

## Phase 2: Extract the entity groups
Status: Complete

- [x] **Step 0 — break the import cycle (hazard 10).** Move `ManualZoneSave` and
      `ManualConnectionSave` (pure structs) to `internal/entities/editor_state/`;
      move `CastleSettingChanges`, the two wrapper models and the four converters
      to `internal/models/editor_state_model/`. Update consumers and move the
      three mirrored unit-test folders. Must land before the groups are written.
- [x] Create `internal/entities/editor_state/` with the 9 files from §0.1, one
      struct per file, camelCase file names (AGENTS.md §4.1).
- [x] Move each field **with its existing json tag verbatim**. Do not rename, do
      not reorder within a group, do not change a type.
- [x] Create `internal/common/common_zone_contents/` and move
      `DefaultPlayerZoneContentRows` → `GetDefaultPlayerZoneContentRows`, with the
      four private helpers renamed to a `get` prefix.
- [x] Move its unit test to the mirrored
      `test/unit/internal/common/common_zone_contents/…` folder.
- [x] Leave `EditorStateDto` untouched this phase — the entities exist but nothing
      uses them yet, so the phase is purely additive.

### Verification Plan
- `go build ./...` clean, full unit + integration suites pass unchanged.
- `go run ./cmd/testlayoutcheck .` passes.
- Manual check: the 9 structs together declare exactly 72 tagged fields, and the
  tag set is byte-identical to the DTO's.
- Coverage ≥ 72.5 % and not below the Phase 1 baseline.
- **Owner gate:** the §0.1 field→group table is signed off before Phase 3 starts.

### Phase Summary

**Step 0 — the import cycle (hazard 10).** Writing the groups would have closed
`editor_state_dto → editor_state_model → entities/editor_state → editor_state_dto`.
The first edge is a settled §0 decision and the second is the point of the batch,
so only the third was removable. Owner ruling: behaviour-free structs
(`ManualZoneSave`, `ManualConnectionSave`) go to `internal/entities/editor_state/`;
`CastleSettingChanges`, the two wrapper models (`ManualZoneSaveModel`,
`ManualConnectionSaveModel`, anonymous embeds with `Clone()`) and the four
converters go to `internal/models/editor_state_model/`. Converters return the
**entity** slice; the six deep-clone helpers stay **private** in
`editor_state_model`. This partly undoes Phase 1's consolidation, and it makes the
layering rule load-bearing: **`internal/entities` must import nothing from
`models`/`dtos`/`helpers`.**

**Groups.** `internal/entities/editor_state/` now holds the 9 structs from §0.1
(`templateIdentity` 2, `mapSettings` 2, `playerSettings` 4, `neutralZoneSettings`
11, `castleSettings` 10, `generationSettings` 15, `gameRuleSettings` 16,
`contentSettings` 10, `manualEditSettings` 2) — **72 tagged fields**, every tag
verbatim.

**Tag verification.** The plan called for comparing the tag set against the DTO at
HEAD, but Phase 1 is now committed (`586fc18 "Batch I wip"`) and that path no
longer exists. Compared instead against the frozen wire-format fixture
`test/test_helpers/testdata/editorState_v0_flat.gen.json` — a stronger check:
entity tag count 72, fixture key count 72, sets equal, 0 duplicates.

**Content rows.** `GetDefaultPlayerZoneContentRows` (+ four `get`-prefixed private
helpers and `const guardedRuleName`) moved to
`internal/common/common_zone_contents/`; both callers (`editorStateDto.go`,
`app/gui/dialogs/zoneContentDialog.go`) updated; the unit test moved to the
mirrored folder. `editorStateDto.go` net `14 +/112 -`.

**Gates.** `go build` clean; `go vet -tags='integration_test,gui'` clean;
`gofmt -l .` empty; `testlayoutcheck` passed; unit, default (frozen golden test
included) and tagged-integration suites all exit 0 — re-run after the lint pass.
Coverage **73.6 %** (floor 72.5 %). Lint `--fix` took 85 issues → **21, all
pre-existing `funcorder`**; `git status --short` counts were identical before and
after, confirming no drive-by files. Overall diff: 32 files, +85/−1007.

**Gotchas worth remembering.** (a) `embedlit`/modernize collapses
`Model{Embedded: X{...}}` to `Model{{...}}`, which orphaned the
`entities/editor_state` import in both `clone_test.go` files. (b) PowerShell
`-replace '…' , '$1' + $type` parses `$1ManualZoneSave` as a *group name*, so the
replacement silently no-ops — use a separate blanket pass. (c) The plan's hazard 4
is also inaccurate: it lists `manualEditDecisionDto.go` as a Phase 1 change that
was never needed.

**Next:** Phase 3 is blocked on the owner gate — the §0.1 field→group table must be
signed off first.

---

## Phase 3: Introduce the model, DTO delegates
Status: Not started

- [ ] Create `internal/models/editor_state_model/editorStateModel.go` —
      `EditorStateModel` anonymously embedding the 9 entity structs.
- [ ] Move the behaviour onto it: `Clone`, `LayoutDefiningOptionsChanged`,
      `DiffCastleSettings`, `EqualsIgnoringManualEdits`, `HasManualEdits`, the
      private comparison helpers, and a `NewDefaultEditorStateModel` factory
      carrying the current defaults.
- [ ] Reduce `EditorStateDto` to `struct { editor_state_model.EditorStateModel }`
      — **anonymous, temporarily** — so every `state.MapSize` selector and the flat
      json output are unchanged and the whole repository still compiles untouched.
- [ ] Keep **DTO-signature shim methods** on the DTO that shadow the promoted
      ones (hazard 1): `Clone() EditorStateDto`,
      `EqualsIgnoringManualEdits(*EditorStateDto) bool`,
      `LayoutDefiningOptionsChanged(*EditorStateDto) bool`,
      `DiffCastleSettings(*EditorStateDto) CastleSettingChanges`. Each delegates
      to the embedded model. Without these, Phase 3 does **not** build.
- [ ] `NewDefaultEditorStateDto` delegates to `NewDefaultEditorStateModel`.
- [ ] Move the DTO behaviour tests to
      `test/unit/internal/models/editor_state_model/editorStateModel/<method>_test.go`,
      one file per public method (AGENTS.md §4.6).
- [ ] Leave the recursive clone guard alone; rewrite the equality guard to walk
      recursive **leaf paths** instead of top-level fields (hazard 6), and add a
      case proving it still trips when a field is added to a nested group but
      omitted from the comparison.

### Verification Plan
- `go build ./...` clean; full unit + integration suites pass **with no test
  changes other than the moves and the equality guard rewrite**.
- Phase 1's golden test passes untouched — this is the proof the key set is
  intact. Expect the key **order** to change; that is allowed (hazard 3).
- Coverage ≥ 72.5 % and not below the Phase 1 baseline.

### Phase Summary
_(write when phase completes)_

---

## Phase 4: Swap the runtime type to the model
Status: Not started

Pure type-name substitution. Because fields stay promoted, **no field-access site
changes** — only signatures and variable types.

- [ ] `app/gui/models.EditorState`: `current`/`previous`/`next` become
      `*editor_state_model.EditorStateModel`; update `GetCurrentState`,
      `OverrideState`, `SetNextState`, `UpdateCurrentState`, `GetPreviousState`,
      `GetNextState`.
- [ ] `internal/handlers/handler_interfaces`: `IStateValidationHandler` and the
      other DTO-carrying signatures take/return the model.
- [ ] `internal/validators`: `ValidationIssue.Fix(*EditorStateModel)`; the 45
      validated fields and 46 issue categories are unchanged.
- [ ] `internal/mappers/generatorConfigMapper.go`: `FromEditorState` takes the
      model. All 69 field reads stay as written.
- [ ] `app/gui/drivers`, `app/gui/editor`, `app/gui/panels`,
      `internal/services/editor`: update the types they pass through.
- [ ] `internal/repositories` and `internal/services/file_service` keep the
      **DTO** — they are the load/save boundary. Name the owner of the conversion
      explicitly: `file_service` wraps model→DTO on save and unwraps DTO→model on
      load, and nothing above it sees a DTO (hazard 5).
- [ ] Convert the **carrier DTOs** that hold editor state inside another DTO to
      carry the model: `EditorStateSaveDto`, `EditorStateValidationDto`,
      `RegenerationDecisionRequestDto`, `TemplateUpdateDto`,
      `CastleSettingsReapplyRequestDto`, `ManualEditDecisionDto`. This is the item
      most likely to be under-estimated — it is not just interface signatures.
- [ ] Delete the Phase 3 shim methods once every caller is on the model.
- [ ] Update every mock under [test/test_helpers](../test/test_helpers) and the
      affected unit/integration tests.

### Verification Plan
- `go build ./...`, both `go vet` variants, `testlayoutcheck` all clean.
- Full unit + integration suites pass.
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`
  passes **without `-update`** — the GPU suite proves the GUI still behaves.
- Phase 1's golden test still passes.
- Round-trip non-aliasing test: load a file, mutate the returned model, reload,
  assert the second read is unaffected.
- `grep` check: `EditorStateDto` appears only in `internal/dtos/editor_state_dto`,
  `internal/repositories` and `internal/services/file_service`.
- Coverage ≥ 72.5 % and not below the Phase 1 baseline.

### Phase Summary
_(write when phase completes)_

---

## Phase 5: Versioned persistence shell
Status: Not started

- [ ] Reshape the DTO to
      `struct { SchemaVersion int; EditorState editor_state_model.EditorStateModel }`
      — named field, no embedding.
- [ ] Add `MarshalJSON` using a **locally declared** alias type that anonymously
      embeds the model (never the DTO — hazard 9) so the output stays flat with
      `schemaVersion` as a sibling key. Always write `1`.
- [ ] Add `UnmarshalJSON` that **merges into the existing receiver** (preserving
      the seed-defaults-then-overlay semantic of
      [editorStateRepository.go](../internal/repositories/editorStateRepository.go#L26-L29)),
      then calls an explicit migration hook normalising `SchemaVersion == 0` → `1`.
- [ ] Point `internal/repositories` and `internal/services/file_service` at the
      new field.
- [ ] Generate the current-writer golden as `editorState_v1_flat.gen.json` and
      keep **both** unstaged fixtures: `_v0_` proves legacy files still load,
      `_v1_` proves the current writer's output round-trips.
- [ ] Add tests: legacy fixture loads and lands at version 1; the two fixtures
      have the same key set apart from `schemaVersion` (**compare parsed objects,
      not bytes**); an absent content-row key still yields the seeded default.
- [ ] Add a characterisation test pinning today's `omitempty` behaviour — nil and
      empty both serialise to an absent key and both reload as the seeded default
      (hazard 2). This documents the limitation rather than pretending to fix it.

### Verification Plan
- All gates from Phase 4, plus: marshal the default state and compare the
  **parsed object** against the pre-batch output — the only permitted delta is
  the added `schemaVersion` key. Key order may differ.
- Load each real `.gen.json` in [output/](../output/) and assert no error and a
  non-default template name. These are the owner's own legacy files and are the
  realest migration evidence available.
- Coverage ≥ 72.5 % and not below the Phase 1 baseline.

### Phase Summary
_(write when phase completes)_

---

## Phase 6: Stop the per-frame whole-state clone (backlog §1.5)
Status: Not started

Baseline to beat, from batch D on `BenchmarkEditorWindow_TabCycling`:
**3.01 M ns/op · 1,435 KB/op · 6,640 allocs/op** (pre-batch-D was 4,676 allocs/op).
Acceptance: any measurable improvement, recorded below.

- [ ] Re-measure the baseline on this machine first — the batch D numbers are from
      a different tree and must not be compared across hardware.
- [ ] Add read-only per-panel view structs in `app/gui/models/`, each carrying
      only what its panel renders.
- [ ] Convert the five cloning `Layout` paths:
      [bonusesPanel.go](../app/gui/panels/bonusesPanel.go#L63),
      [generalPanel.go](../app/gui/panels/generalPanel.go#L105),
      [layoutPanel.go](../app/gui/panels/layoutPanel.go#L131) and
      [layoutPanelZones.go](../app/gui/panels/layoutPanelZones.go#L113) and
      [#L129](../app/gui/panels/layoutPanelZones.go#L129).
- [ ] Keep the clone on any path that hands state **out** of the model — the
      §1.1 aliasing fix must not be undone.
- [ ] Unit tests for each new view struct's projection.

### Verification Plan
- `go test -v -tags='integration_test,gui' -bench=BenchmarkEditorWindow_TabCycling -run=xxx ./test/performance/... -benchmem -benchtime=20x -timeout=120s`
  — record before/after ns/op, B/op, allocs/op in the Phase Summary. Run locally
  on the real GPU; never in CI.
- GPU integration suite passes without `-update`. If a golden genuinely must
  change, `-run`-scope it tightly enough not to match neighbouring tests.
- Full unit + integration suites pass; coverage ≥ 72.5 %.

### Phase Summary
_(write when phase completes)_

---

## Phase 7: Cleanup, docs and gates
Status: Not started

- [ ] Remove any delegation shim left from Phases 3–5.
- [ ] Update the architecture description in [README.md](../README.md#L233),
      which still says the GUI collects input into `dtos.EditorStateDto`.
- [ ] Rewrite backlog [§2.1](../todo/backlog-opus5.md) and
      [§1.5](../todo/backlog-opus5.md) as self-contained ✅ FIXED entries; update
      the header counts and the §8 batch table row **I**.
- [ ] Update [test_observations.md](../todo/test_observations.md) if any gap
      opened or closed.
- [ ] Delete this plan file once the backlog entries are self-contained
      (doc-lifecycle rule), using `Remove-Item`.
- [ ] Write `./.agent/session-carry-forward.md` per AGENTS.md §5.2.

### Verification Plan
- Every gate green in one run: `go build ./...`, both `go vet` variants,
  `go run ./cmd/testlayoutcheck .`, unit, integration, GPU integration,
  `golangci-lint-v2 run ./... --issues-exit-code=0` at **0 issues**,
  `gofmt -l ./app ./internal ./test ./cmd` empty, coverage ≥ 72.5 %.
- `git status --short` reviewed and reported to the owner. **Nothing staged,
  nothing committed.**

### Phase Summary
_(write when phase completes)_

---

## Final Recap
_(write when all phases complete)_

## Deployment Plan
_(write when all phases complete)_
