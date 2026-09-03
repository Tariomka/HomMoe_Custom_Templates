# Session carry-forward — Batch J phase 3

## 1. Session goal

Execute **phase 3** of `.agent/plans/batch-j-zone-tier-source-of-truth.md`:
persist a zone's tier in `.gen.json` so it survives a save/load round trip
instead of being reverse-engineered from the zone's content pools.

## 2. Scope decision taken up front (owner-approved)

Phase 3 could not be done as written. The tier's **write** path —
`ApplyNeutralZoneQuality` and `NewDefaultNeutralZone` — is reachable only through
the zone editor, whose DTOs and dialog still spoke `entities.Zone`, which has
nowhere to put a tier. A manual re-tier therefore could not reach the file at
all, and `carryZoneTiersByName` would have re-attached the **stale planned** tier
by name.

As phase 3's own last checklist item anticipated, the **zone-editor half of phase
4 came first**. Offered A (persist only, ship a known-wrong case), B (pull the
zone-editor chain forward, ~30 prod files) and C (full phase-4 sweep first);
owner chose **B**, plus moving the dialog's Quality dropdown onto `ResolveQuality`
now rather than in phase 4.

## 3. Features added / changed

- **`editor_state.ManualZoneSave.Quality *int8 json:"quality,omitempty"`** — the
  persisted tier. Raw ordinal because an entity may not import `internal/models`
  (§4.4.1 rule 3); **pointer** because the enum is `iota - 1`, so `omitempty` on a
  value field would drop every Plastic zone (ordinal 0) and reload it as "never
  recorded". `editor_state_model.toQualityOrdinal` / `fromQualityOrdinal` convert,
  both preserving nil.
- **The zone-editor chain speaks `template_model.Zone`**: nine DTOs
  (`zoneEditorZonesDto`, `…GeometryRequestDto`, `…ConnectionRequestDto`,
  `…QualityRequestDto`, `…RemoveRequestDto`, `…MutationDto`, `templateUpdateDto`,
  `castleSettingsReapplyRequestDto`, and `ZoneEditorNeutralZoneRequestDto`'s
  return), `IZoneEditorHandler`, `ITemplateHandler.ReapplyCastleSettings`, all
  four `connection_editor` services and their interfaces, `ZoneEditorDialog`,
  `drivers.State.handleUpdateTemplate`, and
  `models.EditorState.SetManualEdits` / `GetManualZones`.
  Connections stayed `entities.Connection` — they carry no tier.
- **`templateHandler.carryZoneTiersByName` is DELETED.** `UpdateTemplate` now
  rebuilds roads on the model zones, maps the *template* through the entity for a
  clean copy, then **re-attaches the applied model zone slice positionally**.
  Exact, where the name lookup was neither.
- **Four consumers stopped inferring** and now call `ResolveQuality`:
  `ManualReapplyService.SetNeutralZoneCastleCount` + `neutralCastleTarget`
  (**this is the §1 `Unknown → Plastic` fix** — `GetNeutralZoneProfile(Unknown)`
  returns the Lowest profile, so an unclassifiable zone whose castles were rebuilt
  silently got Plastic city stats), `ZoneTierService.GetGuardQuality` /
  `GetConnectionGuardQuality` (both now take `[]template_model.Zone`),
  `MandatoryContentProvider.CreateContentsForZones` (whose `Unknown` branch
  attached **no** rows), and `zoneEditorHandler.GetZoneQuality` (the dialog's
  Quality dropdown). **`GetQuality` still takes the entity** — a raw `.rmg.json`
  has no recorded tier and never will.
- `road_helpers.IsRoadTypeConnection` / `IsRoadTypeCastle` now take
  `template_model.Road`. `ZoneEditorService` was their only caller, so an
  entity-typed twin would have been dead weight.

## 4. File modifications

101 modified files + 1 new; full detail in the plan's **Phase 3 → Phase Summary**.

- **New**: `test/integration/manualZoneTierPersistence_integration_test.go`.
- **Edited (production)**: `internal/entities/editor_state/manualZoneSave.go`;
  `internal/models/editor_state_model/manualZoneSave.go`;
  `internal/models/template_model/converters.go` (re-exports `ToZoneModel`,
  `ToMainObjectModels`, `ToRoadModel(s)`); `internal/helpers/road_helpers/`;
  `internal/services/zones/zoneTierService.go` + interface; all four
  `connection_editor` services + interfaces; `mandatoryContentProvider` +
  interface; `templateHandler`, `zoneEditorHandler`, `guiHandler` + two handler
  interfaces; eight `internal/dtos/*.go`; `app/gui/dialogs/zoneEditorDialog.go`,
  `zoneEditorZoneProps.go`, both testexports; `app/gui/drivers/stateManualEdits.go`;
  `app/gui/models/editorState.go`; `app/gui/panels/layoutPanelZones.go`;
  `app/gui/editor/window_testexports.go`.
- **Edited (tests)**: seven mocks in `test/test_helpers/`, five files under
  `test/integration/`, ~50 under `test/unit/`.
- **Deleted**: nothing. **Renamed**: nothing.

## 5. Tests added or updated

- **`manualZoneTierPersistence_integration_test.go`** drives the real save/load
  seam: **Plastic** survives the round trip, **Gold** does too, an unrecorded tier
  writes **no** `quality` key and loads back as nil. Carries
  `//go:build integration_test` because it uses the `SaveStateToFile` /
  `LoadStateFromFile` test-only exports — the one and only reason §4.6.1 allows
  the tag.
  ⚠ **Mutation-verified**: making `toQualityOrdinal` drop the zero ordinal — the
  exact shape of the `omitempty` bug — fails the Plastic test and **only** that
  test.
- Unit round-trip tests on `ToManualZoneSaves` / `FromManualZoneSaves` for both
  the recorded-Plastic and the absent-quality cases.
- `ApplyNeutralZoneQuality` and `NewDefaultNeutralZone` each gained a test
  asserting the tier they **record**, alongside the pre-existing ones asserting
  the profile they **stamp** — those now infer explicitly via `ToZoneEntity`, so
  they still test what they always tested.
- `TestWhenZoneQualityIsRequested_ReturnsTheClassifiersQuality` →
  `..._ReturnsTheResolvedQuality` (its subject changed).
- The gladiator-arena fixture now clears `Quality` explicitly, because
  `NewDefaultNeutralZone` records it and those tests exercise the inference
  fallback on purpose.

## 6. Git status snapshot

Branch `AD/fixing_some_stuff_08-12`, working tree **clean**. The owner reviewed
and committed the batch as **`b210640 Batch J wip 4`**. Nothing was staged by
this session. The review edits were **style-only** — house brace style
(`Quality: x}` rather than a trailing comma and a lone `}`) and flattened
promoted-field composite literals, which is the form `embedlit` wants. No
behaviour was changed in review.

## 7. Verification

| Gate | Result |
| --- | --- |
| `go build ./...`, `go vet ./...`, `go vet -tags='integration_test,gui' ./...` | clean |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 — signatures moved, the graph did not |
| unit / untagged / integration | pass |
| **GPU suite, no `-update`** | **pass** — no pixel moved, no golden modified |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **74.3 %** (unchanged from phase 2.5, floor 72.5 %) |
| Frozen fixtures | unmodified — `omitempty` + absent is exactly the legacy shape |

**Why no pixel moved:** phase 2's 792-configuration sweep proved inference and
the plan agree on every zone the generator emits, and `ApplyNeutralZoneQuality`
stamps the pools its recorded tier implies. So `ResolveQuality` and `GetQuality`
return the same answer for every zone reachable today — the correction is latent
by design, and fires only for a zone inference cannot classify.

## 8. Rejections / deliberate deviations

- **Option A was rejected** (persist only, leave the dialog on entities): it would
  have shipped a zone whose hand-picked tier is silently replaced by the stale
  planned one.
- **`previewLayoutRequestDto` was NOT re-typed**, and `PreviewLayoutService` still
  calls `GetQuality`. Deliberate: it is the one seam whose change moves pixels,
  and it belongs with phase 4's golden regeneration.
- **The generator and topology tree were left on entities.** `Generate` builds the
  entity and lifts it (phase 2.5), which is what makes the golden-template test
  meaningful.
- **`BenchmarkEditorWindow_TabCycling` was not re-measured.** Phase 3 changes no
  per-frame path; `GetCurrentState` and the panel reads are untouched.

## 9. Open questions

- Should the preview colour a zone by its **recorded** tier (phase 4), or is
  inference the intended behaviour there? Changing it will move goldens, so it
  wants an explicit decision rather than a sweep.
- Phase 4 should state outright whether the **topology tree stays on entities**.
  The plan's original "~100 production files move onto `template_model`" reads
  like a reflex; the golden test argues the generator should keep building the
  entity.

## 10. Next recommended actions

1. **Phase 4**, in this order: `preview_service` first (the only visually
   meaningful gap — `PreviewLayoutService.GetQuality` → `ResolveQuality`, plus
   `previewLayoutRequestDto`), then decide per seam for `file_service`, the
   generator and the topology tree.
   ⚠ This is the first change in the batch that **will** move goldens. Regenerate
   with `-update`, then check `git status` for ` M *.golden` and restore any that
   were not yours.
2. **Phase 5** — shrink `entityNamerAllowList` in
   `test/unit/architecture/dependency/layering_test.go`. Phase 3 removed
   `internal/dtos`' zone-editor entity references, so re-measure the whole list
   rather than the three packages the plan names. Only ever remove entries.
3. Backlog §2.2 → ✅ DONE record; §8 gets row **J**; refresh the coverage figure
   in the three places it is quoted.

## 11. Carry-forward prompt

> Read `AGENTS.md` first — especially **§4.4.1**, the Entity/Model/DTO doctrine:
> **Entity** (`internal/entities/`) is the wire format, json tags only; **Model**
> (`internal/models/`) owns the structure and all business logic; **DTO**
> (`internal/dtos/`) is the thin `app/` ↔ `internal/` crossing. A DTO carrying a
> Model is intended and `app/` may hold a Model — do not "fix" either.
>
> Then read `.agent/plans/batch-j-zone-tier-source-of-truth.md` and do **phase
> 4**. Read **§0 and §0b — where they disagree, §0b wins.** Phases 1, 2, 2.5 and
> 3 are complete and **committed** (`b210640`). Read phase 3's Phase Summary
> before touching anything: it records that the **zone-editor half of phase 4 was
> already done there**, so phase 4's checklist is mostly ticked and what remains
> is `preview_service`, `file_service`, the generator and the topology tree.
>
> The template layer lives in `internal/models/template_model/`, a full mirror of
> `internal/entities/template/`. The zone tier is
> `template_model.Zone.Quality *neutral_zone.Quality`, nil meaning "not recorded,
> infer it", and it now persists to `.gen.json` as
> `editor_state.ManualZoneSave.Quality *int8`. `IZoneTierService.ResolveQuality`
> takes the **model** zone and prefers the recorded tier; `GetQuality` takes the
> **entity** and always infers — keep both, a raw `.rmg.json` has no recorded
> tier.
>
> Hard rules, one line each: never modify `data/`, `internal/registry/`, or
> anything under `internal/entities/template/` — `internal/entities/editor_state/`
> is *not* protected; everything must build and run on Windows and Linux
> (`path/filepath`; chain PowerShell with `;`, never `&&`); every change ships
> with tests and unit coverage must not drop below 72.5 % (currently **74.3 %**);
> the lint baseline is **0 issues**; **never stage and never commit** — use
> `Move-Item` never `git mv`, `Remove-Item` never `git rm`; never change where
> `.rmg.json` is written and never persist the output directory; never run a bulk
> in-place rewrite, and **never round-trip a `.go` file through
> `Get-Content`/`Set-Content`** — use `gofmt -r` on an explicit file list and
> verify insertions == deletions per file.
>
> Standing traps: **nil is load-bearing** three times over — nil `Previous` =
> first generation, nil `Next` = unarmed debounce, nil `Quality` = "infer it";
> the persisted tier is `*int8` because `omitempty` on a plain `int8` would
> silently drop every Plastic zone (ordinal 0), and there is a mutation-verified
> guard for exactly that in
> `test/integration/manualZoneTierPersistence_integration_test.go`; the two frozen
> fixtures under `test/test_helpers/testdata/` plus the untagged
> `editorStateWireFormat_integration_test.go` must keep passing **unchanged**,
> comparing **parsed objects, never bytes**; `cmd/testlayoutcheck` matches
> test-only export names tree-wide, so grep any new accessor name before adding
> it; a file gets `//go:build integration_test` **only** if it calls a
> `*_testexports.go` accessor; and `helpers.MapSlice` / `helpers.MapPointer`
> preserve nil-vs-empty where `linq.SelectSlice` does not — never swap them
> inside a converter.
>
> ⚠ **Unlike every phase so far, phase 4 WILL move pixels** — the preview colours
> a zone by inference today. Regenerate goldens with `-update`, then check
> `git status` for ` M *.golden` and `git restore` any that were not yours.
>
> Cap sessions at ~50 messages and hand off through this file. Full handoff in
> `./.agent/session-carry-forward.md`.
