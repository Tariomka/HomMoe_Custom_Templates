# Session carry-forward — Batch J phase 4

## 1. Session goal

Execute **phase 4** of `.agent/plans/batch-j-zone-tier-source-of-truth.md`: move
the last consumers off the `TemplateMapper.ToEntity` bridges so the preview
colours a zone by the tier that was **recorded** for it rather than by the tier
its content pools imply, and decide explicitly which seams keep the entity.

Phases 1, 2, 2.5 and 3 were already complete and committed (`b210640`). Phase 3
had pulled the zone-editor half of phase 4 forward, so what remained was
`preview_service`, `file_service`, the generator and the topology tree.

## 2. Fixes applied

- **The preview stopped inferring.**
  [previewLayoutService.go](internal/services/preview_service/previewLayoutService.go)'s
  `buildPreviewZones` now calls `ResolveQuality(zone)` instead of
  `GetQuality(zone)`. This is the last consumer in the batch to move.
- **`ZoneEditorGeometryService` stopped erasing the tier.**
  [zoneEditorGeometryService.go](internal/services/connection_editor/zoneEditorGeometryService.go)
  used to round-trip its model zones through `ToZoneEntities` to synthesise a
  template for the layout service, which silently dropped `Quality`. It now
  builds a model variant.
- **`templateHandler.SaveTemplate` renders before it flattens.**
  [templateHandler.go](internal/handlers/templateHandler.go) hands the preview
  generator the **model** and only then maps to the entity for the file service.
  Flatten-first would have thrown the tier away before the PNG could read it.

## 3. Features added / changed

- **`preview_service` names no entity type at all.** `BuildPreviewLayout` and
  `CreatePreviewImage` (real generator *and*
  [nullPreviewGeneratorService.go](internal/services/preview_service/nullPreviewGeneratorService.go))
  take `*template_model.Template`; all five layout strategies — ring/hub,
  scatter, fixed geometry, balanced rings, manual positions — plus
  [layoutGeometry.go](internal/services/preview_service/layoutGeometry.go)'s
  shared predicates take `[]template_model.Zone` / `[]template_model.Connection`.
- **`PreviewLayoutRequestDto` moved onto the model and then shed two fields.**
  `Zones` / `.Connections` — the "editor-only preview when `Template` is nil"
  branch — had **no production caller**: `previewPanel` always passes a template
  and the zone editor goes through `ZoneEditorGeometryService`. The only reader
  was the synthesis branch inside `previewHandler` itself. Owner confirmed the
  removal, so both fields, the branch and the three tests that existed solely to
  cover it are gone. The DTO is now `Template`, `Topology`, `CanvasSide`, and
  `BuildPreviewLayout` is a single forwarding line.
- **`NewPreviewHandler` lost its `mappers.ITemplateMapper` argument.** The
  handler no longer converts anything; the zones-only branch builds a
  one-variant `template_model.Template` literal rather than running the entity
  `VariantBuilder`. `wire_gen.go` was regenerated with `wire gen`.
- **Two seams were decided to stay on the entity, deliberately:**
  - `file_service.SaveTemplateWithPreview` writes `.rmg.json`, one of the two
    places §0b.16 keeps the wire format forever.
  - The generator and the **whole topology tree** keep building entities.
    `Generate` assembles the entity from the providers and lifts it with
    `TemplateMapper.ToModel` before stamping tiers, and that ordering is what
    makes `TestWhenDefaultConfiguration_ReturnsGoldenTemplate` a proof rather
    than an argument. No topology has a tier to carry — `stampPlannedZoneTiers`
    derives all of them afterwards. This answers the previous carry-forward's
    open question; do not revisit it by reflex.

- **`BuildPreviewLayout` no longer returns an error.** With the synthesis branch
  gone it could only ever return `nil`, and the layout service has no failure
  mode, so the whole chain returns a bare `dtos.PreviewLayoutDto`:
  `IPreviewHandler`, `previewHandler`, `GUIHandler`, `TemplateHandlerMock` and
  the guiHandler stub. That emptied
  [previewLayoutCache.go](app/gui/models/previewLayoutCache.go) too — its `build`
  callback and its "a failed build is not cached" retry existed only to carry
  that error, so `Get` now takes `func() preview.Layout` and returns
  `preview.Layout`. `previewPanel` lost the unreachable error branch it fed.

## 4. File modifications

31 modified files, nothing created, nothing deleted, nothing renamed.

- **Production (18)**: `internal/composition/wire_gen.go` (regenerated);
  `internal/dtos/previewLayoutRequestDto.go`;
  `internal/handlers/previewHandler.go`, `templateHandler.go`, `guiHandler.go`,
  `handler_interfaces/previewHandlerInterface.go`;
  `internal/services/connection_editor/zoneEditorGeometryService.go`;
  all ten files of `internal/services/preview_service/`;
  `app/gui/models/previewLayoutCache.go`, `app/gui/panels/previewPanel.go`.
- **Test helpers (3)**: `previewLayoutServiceMock.go`,
  `previewGeneratorServiceMock.go`, `templateHandlerMock.go`.
- **Tests (10)**: `test/performance/preview_layout_test.go`; the `guiHandler`
  (×3), `previewHandler` (×2), `templateHandler`, `zoneEditorGeometryService` and
  `previewLayoutCache` unit tests; all four preview_service unit-test files.

## 5. Tests added or updated

- **New**, in
  `test/unit/internal/services/preview_service/previewLayoutService/buildPreviewLayout_test.go`:
  - `TestWhenZoneCarriesARecordedTier_ColoursItWithThatTierInsteadOfInferring` —
    a zone with `Sides` layout and a `_t2_` guarded pool infers as `QualityLow`;
    recording `QualityHigh` on it makes the preview report High.
    ⚠ **Mutation-verified**: reverting `buildPreviewZones` to
    `GetQuality(ToZoneEntity(zone))` fails this test and **only** this test.
  - `TestWhenZoneCarriesNoRecordedTier_ColoursItWithTheInferredTier` — pins the
    nil-Quality fallback so the inference branch cannot rot.
  - A `lowTierZone` helper in that package's `common_test.go`.
- Updated expectations, all of them because the seam's *type* changed rather
  than its behaviour: `previewHandler`'s mock assertions and constructor call,
  `templateHandler.SaveTemplate`'s preview-generator expectation (now the model),
  `zoneEditorGeometryService`'s laid-out-template assertion, and the whole
  preview_service unit suite (swept with `gofmt -r`).
- `test/performance/preview_layout_test.go` dropped a `ToEntity` of the whole
  generated template from its per-case setup — a strict reduction.

Suite status at handoff: unit / untagged / integration / **GPU** all exit 0.

## 6. Git status snapshot

Branch `AD/fixing_some_stuff_08-12`. 31 modified files, **zero staged**, **zero
`*.golden` modified**, no untracked files added. Nothing was committed.

## 7. Verification

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 (regenerated, never hand-edited) |
| unit / untagged / integration | pass |
| **GPU suite, no `-update`** | **pass** — no pixel moved, no golden modified |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **74.3 %** (unchanged, floor 72.5 %) |

**⚠ The plan predicted this phase would move pixels. It did not.** Two facts
combine: phase 2's 792-configuration sweep proved inference and the plan agree on
every zone the *generator* emits; and every zone the *editor* re-tiers goes
through `ApplyNeutralZoneQuality`, which stamps the content pools the new tier
implies, so inference lands on the same answer. The correction stays latent — it
fires only for a zone whose tier is recorded but whose pools do not imply it,
which no path reachable today produces.

**The pass is not vacuous.** `BaseHandler` masks the preview canvas interior by
default because the shipped default topology is Random, but
`LayoutAndZonesTabHandler.SelectTopology` **lifts that mask** for any
deterministic topology. Every snapshot after that point compares the canvas pixel
for pixel, including the zone-editor suites that re-tier a zone through the
dropdown.

## 8. Rejections / deliberate deviations

- **Goldens were not regenerated.** The suite is green without `-update` and no
  golden is stale; regenerating would rewrite ten files with anti-aliasing noise
  and no signal.
- **`file_service` was not moved.** It is the `.rmg.json` writer.
- **The generator and topology tree were not moved.** See §3.
- **`PreviewLayoutRequestDto.Zones` / `.Connections` were deleted** after the
  owner confirmed the sweep had proved them unreachable. Three tests went with
  them: `TestWhenOnlyZonesAreProvided_LaysOutASynthesizedTemplate`,
  `TestWhenOnlyConnectionsAreProvided_LaysOutASynthesizedTemplate` and
  `TestWhenRequestContainsZones_BuildsPreviewTemplate`. Nothing was recorded in
  `test_observations.md` — the code is gone, not merely untested.
- **The preview layout error return was deleted** on the same reasoning, at the
  owner's request. `TestWhenBuildFails_ReturnsTheError` and
  `TestWhenBuildFails_IsNotCached` were removed with the failure mode they
  described — the cache can no longer be handed a build that fails.
- **`golangci-lint --fix` restyled two signatures** as `param,\n) Ret {`; both
  were hand-restyled to the repo's `param) Ret {` house form.

## 9. Open questions

- Phase 5 will find `internal/dtos` still on `entityNamerAllowList`: five files
  (`templateUpdateDto`, `zoneEditorGeometryRequestDto`, `zoneEditorMutationDto`,
  `zoneEditorRemoveRequestDto`, `zoneEditorZonesDto`) name `entities.Connection`,
  which phase 3 deliberately left alone because a connection carries no tier.
  Moving connections purely to shorten a list needs an explicit decision.

## 10. Next recommended actions

1. **Phase 5** — re-measure `entityNamerAllowList` in
   `test/unit/architecture/dependency/layering_test.go` across the whole list,
   not just the three packages the plan names. **Only ever remove entries**, and
   prove the shrink is real by re-adding a removed entry and watching the test
   fail. `internal/dtos` will not come off yet (see §9).
2. Backlog `.agent/backlog/backlog-opus5.md`: §2.2 becomes a ✅ DONE record with
   the behaviour deltas; §8 gets row **J**; refresh the coverage figure in the
   three places it is quoted.
3. Fill in the plan's **Final Recap** and **Deployment Plan**.

## 11. Carry-forward prompt

> Read `AGENTS.md` first — especially **§4.4.1**, the Entity/Model/DTO doctrine:
> **Entity** (`internal/entities/`) is the wire format, json tags only; **Model**
> (`internal/models/`) owns the structure and all business logic; **DTO**
> (`internal/dtos/`) is the thin `app/` ↔ `internal/` crossing. A DTO carrying a
> Model is intended and `app/` may hold a Model — do not "fix" either.
>
> Then read `.agent/plans/batch-j-zone-tier-source-of-truth.md` and do **phase
> 5**, the last one. Read **§0 and §0b — where they disagree, §0b wins.** Phases
> 1, 2, 2.5 and 3 are committed (`b210640`); **phase 4 is complete but
> uncommitted** — read its Phase Summary before touching anything.
>
> Phase 5 is a measurement, not a sweep: shrink `entityNamerAllowList` in
> `test/unit/architecture/dependency/layering_test.go` by re-measuring the whole
> list, **only ever removing entries**, and prove the shrink is real by re-adding
> one and watching the test fail. `internal/dtos` will **not** come off — five
> DTOs still name `entities.Connection`, which phase 3 deliberately left alone
> because a connection carries no tier. Then update
> `.agent/backlog/backlog-opus5.md` (§2.2 → ✅ DONE, §8 row **J**, coverage figure
> in three places) and write the plan's Final Recap and Deployment Plan.
>
> The template layer lives in `internal/models/template_model/`, a full mirror of
> `internal/entities/template/`. The zone tier is
> `template_model.Zone.Quality *neutral_zone.Quality`, nil meaning "not recorded,
> infer it", persisted as `editor_state.ManualZoneSave.Quality *int8`.
> `IZoneTierService.ResolveQuality` takes the **model** zone and prefers the
> recorded tier; `GetQuality` takes the **entity** and always infers — keep both,
> a raw `.rmg.json` has no recorded tier. Two seams keep the entity **on
> purpose**: `file_service` (it writes `.rmg.json`) and the generator plus the
> whole topology tree (`Generate` lifts the assembled entity, which is what makes
> the golden-template test a proof). Do not sweep them.
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
> inside a converter. New in phase 4: `golangci-lint --fix` wraps a long
> signature as `param,\n) Ret {`, but the house style is `param) Ret {` —
> restyle by hand after a `--fix`.
>
> Cap sessions at ~50 messages and hand off through this file. Full handoff in
> `./.agent/session-carry-forward.md`.
