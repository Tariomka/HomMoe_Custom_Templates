# Batch J — zone tier single source of truth (backlog §2.2 branch B)

A neutral zone's tier is decided once at plan time, flattened into layout/pools/
castles, and then reverse-engineered by `ZoneClassifier` everywhere it is needed.
This batch records the tier instead: the generator hands back what it planned, a
`models.QualifiedZone` wrapper carries it through the editor, `.gen.json`
persists it, and a single `IZoneTierService` answers the question — falling back
to inference only for a template loaded from a raw `.rmg.json`.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is
done, set its status to `Complete` and write its **Phase Summary** (what was
done, key decisions, anything needed to continue with zero context); run the
phase's **Verification Plan** and record the result before moving on. When all
phases are done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../../AGENTS.md) first — especially **§4.4.1** (Entity/Model/DTO
doctrine) and **§2.1** (protected directories). Hard rules that bite here:
**`internal/entities/template/**` is protected and this batch must not touch a
single byte of it**; never stage and never commit; move with `Move-Item`, never
`git mv`; delete with `Remove-Item`, never `git rm`; never round-trip a `.go`
file through `Get-Content`/`Set-Content`; unit coverage must not drop below
**72.5 %** (currently **73.8 %**); lint baseline is **0 issues**.

## 0. Decisions (settled with the owner 2026-09-01 — do not relitigate)

1. **Branch B, not A.** `entities.Zone` gets no `Quality` field. Backlog §2.2
   branch A (a `json:"-"` field on the protected schema) is **not** in scope and
   is not to be re-proposed as a shortcut when phase 4 gets tedious.
2. **The wrapper is `models.QualifiedZone`, and it embeds `entities.Zone`.**
   Embedding is deliberate: field promotion means every existing `zone.Name`,
   `zone.Layout`, `zone.GuardedContentPool` keeps compiling, which turns most of
   phase 4 from a rewrite into a type swap. §4.4.1 blesses this — "a Model group
   embeds its entity group and adds the behaviour". The name is `QualifiedZone`
   because the enum is `Quality`; `PlannedZone` was rejected as misleading (the
   type also carries hand-created and `.rmg.json`-loaded zones that never had a
   `neutral_zone.Plan`).
3. **The wrapper needs a store behind it, and the store is the generator.**
   `Variant.Zones []Zone` is protected and is written on **both** the generation
   path (`WithZones`) and the apply-back path
   (`templateHandler.UpdateTemplate` → `newTemplate.Variants[0].Zones = ...`),
   and the editor is opened from `lastTemplate.Variants[0].Zones`. So the
   working slice is rebuilt from protected zones repeatedly and the tier must be
   re-attached each time. A wrapper with no store is just the classifier with
   extra steps. `Generate` therefore returns the tiers it planned and
   `drivers.State` carries the index beside `lastTemplate`.
4. **`Generate()` changes signature** to `(*models.GeneratedTemplate, []string)`,
   and the ~130 test call sites get the mechanical `actual` → `actual.Template`
   edit. A second `GenerateWithTiers()` entry point was **rejected** — it would
   exist only to avoid test churn, which §3.1 calls a speculative abstraction.
5. **`IZoneTierService` replaces `IZoneClassifier.GetQuality` at all 8
   consumers**, and takes `GetGuardQuality` / `GetConnectionGuardQuality` with
   it — they are tier queries over the same data. **`ZoneClassifier` is deleted
   outright** (owner call at phase 1 review): the tier service owns the
   inference natively rather than delegating to a wrapped classifier, so that it
   gains the recorded-tier path *alongside* the old behaviour instead of on top
   of another type. Inference itself can never go away — a template loaded from
   a raw `.rmg.json` has no recorded tier — but it is now the service's own
   fallback branch, not a separate collaborator.
6. **The persisted tier is nullable, and `nil` means "not recorded".** A plain
   `int8` with `omitempty` would **silently drop every Plastic zone**
   (`QualityLowest` is 0) back to "absent" — the exact bug this item is about.
   This mirrors `ManualPosition *[2]float64`, which already exists in the same
   struct for the same reason. The **entity** stores `*int8` (it may not import
   `internal/models/neutral_zone`; see phase 3) and the **model** exposes
   `*neutral_zone.Quality`. `internal/entities/editor_state/` is **not**
   protected; only `internal/entities/template/` is.
7. **The 9 zone-editor DTOs carry `[]models.QualifiedZone`.** Doctrine-consistent
   (§4.4.1: a DTO carrying a Model is intended), the tier crosses as one value
   rather than a parallel field, and it drops `internal/dtos` out of the entity
   allow-list for free.
8. **Output changes are approved.** Zones the classifier currently calls
   `QualityUnknown` will start getting their planned mandatory-content rows,
   arena eligibility, castle city-guard values and connection guard defaults.
   That correction is the point of the item. **Every delta must be enumerated in
   the phase summary that causes it** — an unexplained golden move is a bug
   until proven otherwise.
9. **One batch, phased, reviewed per phase.** ~31 production files plus tests.

## 1. The eight consumers (the thing being fixed)

| Consumer | What the tier decides today |
| --- | --- |
| [previewLayoutService.go](../../internal/services/preview_service/previewLayoutService.go) | `preview.Zone.Quality` → canvas fill colour and the PNG asset name (`neutral_high`, …) |
| [zoneEditorZoneProps.go](../../app/gui/dialogs/zoneEditorZoneProps.go) | which entry the Quality dropdown shows |
| [manualReapplyService.go](../../internal/services/connection_editor/manualReapplyService.go) | `SetNeutralZoneCastleCount`: which profile's city-guard values and building SIDs rebuild the castles; and `neutralCastleTarget`: which advanced castle-count option applies |
| [mandatoryContentProvider.go](../../internal/services/template_generator/providers/mandatoryContentProvider.go) | `CreateContentsForZones`: which mandatory-content row set is attached after a manual edit (`Unknown` → **no rows**) |
| [gladiatorArenaProvider.go](../../internal/services/template_generator/providers/gladiatorArenaProvider.go) | which zone wins the arena, and which N–N connection becomes a `GladiatorArena` (`Unknown` scores −1, so it never wins) |
| [zoneEditorHandler.go](../../internal/handlers/zoneEditorHandler.go) | pass-through for the dialog |
| [connectionEditorService.go](../../internal/services/connection_editor/connectionEditorService.go) | `NewDefaultConnection`'s `GuardValue` |
| [zoneEditorService.go](../../internal/services/connection_editor/zoneEditorService.go) | **writer, not reader** — it is handed a Quality and stamps the profile. It must start *recording* the tier it stamps. |

Two traps found while scoping, both load-bearing:

- **`PreviewLayoutService` bypasses DI.** `NewPreviewLayoutService()` hard-builds
  its own `NewZoneClassifier()`, so it is not the wire singleton and would
  silently keep using raw inference. Fixed in phase 1.
- **`Unknown` → Plastic is a silent down-tier.** `GetNeutralZoneProfile(QualityUnknown)`
  returns the Lowest profile, so an unclassifiable zone whose castles are rebuilt
  gets Plastic city stats today. Recording the tier removes that path; say so in
  the phase 2 summary.

---

## Phase 1: `IZoneTierService` absorbs the classifier
Status: **Complete** (2026-09-01) — uncommitted, awaiting review.

Pure indirection. The tier is still inferred, so **generated output must be
byte-identical** — this phase is the safety net that proves the seam is correct
before any behaviour moves.

- [x] New `internal/services/zones/zoneTierService.go` + its interface per
      §4.2.2 (the package already has `zone_interfaces/`, so it goes there):
      `GetQuality(zone entities.Zone) neutral_zone.Quality`,
      `GetGuardQuality(...)`, `GetConnectionGuardQuality(...)`. It **owns** the
      inference; `zoneClassifier.go` and `IZoneClassifier` are deleted.
- [x] Swap the 8 consumers onto `IZoneTierService`. Nothing outside the zones
      package names a classifier any more.
- [x] `NewPreviewLayoutService()` takes the tier service as a parameter instead of
      constructing a classifier. Update `providerSets.go`, regenerate with
      `wire gen ./internal/composition/...` — never hand-edit `wire_gen.go`.
      Update the ~60 `NewPreviewLayoutService()` call sites in
      `test/unit/internal/services/preview_service/previewLayoutService/`.
- [x] Unit tests for the new service per §4.6 (one folder per file, one file per
      public method).

### Verification Plan
- `go build ./...`, both `go vet` variants, `go run ./cmd/testlayoutcheck .`.
- `wire diff ./internal/composition/...` exit 0 after regeneration.
- Unit + untagged + integration suites pass; **GPU suite passes without
  `-update`** — no pixel may move in a pure-indirection phase.
- Coverage ≥ 72.5 %.

### Phase Summary

**Landed.** `IZoneTierService` is now the only way the application asks for a
zone's tier. Six injection sites moved onto it — `MandatoryContentProvider`,
`GladiatorArenaProvider`, `ConnectionEditorService`, `ManualReapplyService`,
`zoneEditorHandler` and `PreviewLayoutService` — and the field is named
`tierService` everywhere, so no consumer still says "classifier".

**`PreviewLayoutService`'s DI bypass is fixed.** It took no constructor argument
and hard-built its own `NewZoneClassifier()`, so it was never the wire
singleton. It now takes `IZoneTierService`, and `wire_gen.go` shows
`preview_service.NewPreviewLayoutService(iZoneTierService)`. This was the trap
that would have made phase 2 silently keep inferring on the preview path.

**One deviation from the plan, settled at review.** The plan had the tier service
*wrapping* `ZoneClassifier`. Two things pushed against that. First, with the
service delegating one-for-one the two interfaces had **identical method sets**,
which the `iface` linter correctly flagged as redundant against a 0-issue
baseline. Second — the owner's call — a service that only forwards to a
classifier means phase 2 would bolt the recorded-tier path onto a *wrapper*,
leaving the real logic parked in another type forever.

So **`ZoneClassifier` is gone**. `zoneClassifier.go` and
`zoneClassifierInterface.go` were deleted and their bodies now live on
`ZoneTierService` verbatim — `GetQuality`, `GetGuardQuality`,
`GetConnectionGuardQuality` plus the three private branches `getCenterQuality`,
`getTreasureQuality`, `getSidesQuality`. `NewZoneTierService()` takes no
arguments and the struct is empty, exactly as `ZoneClassifier` was. Phase 2 adds
the recorded-tier lookup as a branch **in front of** that inference, in the type
that owns it. Consequences:

- `test_helpers/zoneClassifierMock.go` was **renamed** (`Move-Item`) to
  `zoneTierServiceMock.go`, type `ZoneTierServiceMock`. It mocks the one
  surviving interface, and `zoneEditorHandler`'s fixture — the one place that
  genuinely needs to fake a tier — keeps using it (`fixture.zoneClassifier` →
  `fixture.tierService`).
- The **whole classifier unit suite moved** to
  `test/unit/internal/services/zones/zoneTierService/` per §4.6 (the folder is
  named after the implementation file). Nothing was thinned: the table-driven
  `getQuality` suite, the seven `getGuardQuality` cases and the three
  `getConnectionGuardQuality` cases all came across intact, against the service.
- A short-lived `test_helpers.NewZoneTierService()` shim was deleted again once
  the constructor lost its parameter — a wrapper around a no-arg constructor is
  noise, and going direct returns most of those test files to nearly their
  original text.

**Mechanical edits were done with `gofmt -r`**, not text substitution — an AST
rewrite cannot mangle a `.go` file the way a `Get-Content`/`Set-Content` round
trip can. That covered the ~60 `NewPreviewLayoutService()` call sites in one file
plus the constructor swaps across roughly twenty test files.

**No behaviour changed.** The inference code moved verbatim and every gate
agrees.

### Verification results (2026-09-01)

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

---

## Phase 2: `QualifiedZone` + the generator records what it planned
Status: Not started

Behaviour deltas start here.

- [ ] `internal/models/qualifiedZone.go` — `QualifiedZone` embedding
      `entities.Zone` plus `Quality neutral_zone.Quality`, with the helpers the
      later phases need (wrap/unwrap a slice, look up by name). Keep it small;
      add helpers when a caller needs one, not before.
- [ ] `internal/models/generatedTemplate.go` — `GeneratedTemplate{Template
      *entities.RmgTemplate; ZoneTiers map[string]neutral_zone.Quality}`.
- [ ] Stamp the tier where the factory currently *consumes* it:
      `ZoneFactory.CreateNeutralZone` takes `Quality` and converts it to a
      Profile at [zoneFactory.go](../../internal/services/zones/zoneFactory.go);
      `CreateHubZone` always uses `QualityHighest`; spawn zones are
      `QualityUnknown`. Carry those out to the topology layer and into
      `GeneratedTemplate.ZoneTiers`.
- [ ] `TemplateGenerator.Generate()` returns
      `(*models.GeneratedTemplate, []string)`. Update `templateHandler.go:65`
      and the ~130 test call sites (`actual, _ :=` → `actual.Template`).
- [ ] `IZoneTierService` gains a recorded-tier lookup that is preferred over
      inference; `drivers.State` carries the index beside `lastTemplate`.
- [ ] Tests: a test per consumer proving the tier it reads is the tier the
      generator planned (backlog §2.2 asks for exactly this).

**Expect goldens and generated output to move here.** Enumerate every delta —
which topology, which zone, which of the five decisions from §1 changed, and why
the recorded tier is the correct answer. A delta you cannot explain is a bug.

### Verification Plan
- Full gate set from phase 1.
- Diff a generated `.rmg.json` for each topology before/after; every difference
  must map to a zone inference called `Unknown`.
- `BenchmarkEditorWindow_TabCycling` re-measured **by hand** (needs a GPU, never
  runs in CI) if the clone path or `State` shape changed — baseline ~4,773
  allocs/op.

### Phase Summary
_(write when phase completes)_

---

## Phase 3: persist the tier in `.gen.json`
Status: Not started

- [ ] `internal/entities/editor_state/manualZoneSave.go` gains a nullable tier.
      **It cannot be typed `neutral_zone.Quality`**: that enum lives in
      `internal/models/neutral_zone`, and an entity importing `internal/models`
      is a defect under §4.4.1 rule 3 (the layering gate would fail on it). So
      the entity stores the raw ordinal —

      ```go
      Quality *int8 `json:"quality,omitempty"`
      ```

      — and `internal/models/editor_state_model` converts to and from
      `neutral_zone.Quality` at the mapper seam, which is where conversion
      belongs anyway. `internal/entities/editor_state/` is **not** protected;
      only `internal/entities/template/` is.
- [ ] Model + mapper + GUI round trip: `ToManualZoneSaves` / `FromManualZoneSaves`
      carry the tier; `EditorState.SetManualEdits` / `GetManualZones` speak
      `[]models.QualifiedZone`.
- [ ] The write path must actually record: `ApplyNeutralZoneQuality` and
      `NewDefaultNeutralZone` currently flatten Quality into the profile and
      forget it.
- [ ] Backward compatibility: a `.gen.json` with no `quality` loads as `nil` and
      falls back to inference. Add a test that proves it, and one that proves a
      **Plastic** zone survives a save/load round trip (the `omitempty` trap).

**The two frozen fixtures**
([editorState_v0_flat.gen.json](../../test/test_helpers/testdata/editorState_v0_flat.gen.json),
[editorState_v1_flat.gen.json](../../test/test_helpers/testdata/editorState_v1_flat.gen.json))
and the untagged `editorStateWireFormat_integration_test.go` must keep passing
**unchanged**, comparing **parsed objects, never bytes**. Do not regenerate them.

### Verification Plan
- Full gate set.
- The two frozen fixtures pass untouched; `git status` shows them unmodified.
- Round-trip test: Plastic in → Plastic out, no inference involved.

### Phase Summary
_(write when phase completes)_

---

## Phase 4: the wrapper through the editor chain
Status: Not started

The big mechanical phase (~24 files). Embedding is what makes it survivable.

- [ ] The 9 DTOs carry `[]models.QualifiedZone`:
      `zoneEditorZonesDto`, `zoneEditorGeometryRequestDto`,
      `zoneEditorConnectionRequestDto`, `zoneEditorQualityRequestDto`,
      `zoneEditorRemoveRequestDto`, `zoneEditorMutationDto`,
      `templateUpdateDto`, `castleSettingsReapplyRequestDto`,
      `previewLayoutRequestDto`.
- [ ] `ZoneEditorDialog.zones` / `originalZones` become `[]models.QualifiedZone`;
      `selectedZoneRef` / `zoneByName` return `*models.QualifiedZone`;
      `zonePropertyRows` / `syncZoneProps` / `writebackZoneProps` follow. The
      Quality dropdown reads the carried tier instead of classifying.
- [ ] Handlers and the `connection_editor` services take the wrapper.
      Wrap/unwrap **only** at the protected `Variant.Zones` field, in exactly two
      places: `WithZones(...)` on the way in and
      `newTemplate.Variants[0].Zones = ...` on the way out.
- [ ] `*_testexports.go` accessors follow. ⚠ **Name-collision trap**: the layout
      checker matches test-only exports by identifier name tree-wide, so grep any
      new accessor name across the repo first — the testexports side always
      yields.

### Verification Plan
- Full gate set, including `go vet -tags='integration_test,gui' ./...`.
- **GPU suite without `-update`** — phase 4 is a type swap, so no pixel may move.
  If one does, phase 2 leaked into phase 4.
- Coverage ≥ 72.5 %.

### Phase Summary
_(write when phase completes)_

---

## Phase 5: shrink the allow-list, record the outcome
Status: Not started

- [ ] Remove `internal/dtos` from `entityNamerAllowList` in
      [layering_test.go](../../test/unit/architecture/dependency/layering_test.go)
      **if** the 9 DTOs no longer name `entities.Zone`. Same check for
      `app/gui/dialogs` and `internal/handlers`. **Only ever remove entries** —
      if one will not come off, leave it and say why here.
- [ ] Backlog: §2.2 becomes a ✅ DONE record with the behaviour deltas; §8 gets
      row **J**; refresh the coverage figure everywhere it is quoted (three
      places); update §2.6's file counts if the allow-list moved.
- [ ] Update `.agent/session-carry-forward.md`.

### Verification Plan
- `go test ./test/unit/architecture/... -count=1` passes, and **fails if a
  removed entry is re-added as a violation** — prove the shrink is real.
- Every gate from every prior phase still green.

### Phase Summary
_(write when phase completes)_

---

## Final Recap
_(write when all phases complete)_

## Deployment Plan
_(write when all phases complete)_
