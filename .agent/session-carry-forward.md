# Session carry-forward — Batch J phase 2.5

## 1. Session goal

Execute **phase 2.5** of `.agent/plans/batch-j-zone-tier-source-of-truth.md`:
replace `models.GeneratedTemplate` plus its `map[zoneName]Quality` side-car with
`internal/models/template_model/`, a full one-for-one mirror of
`internal/entities/template/`, so the zone tier becomes a field on the zone.

## 2. Fixes applied

- `BuildPreviewLayout` no longer hands the layout service the caller's template
  pointer — [previewHandler.go](../internal/handlers/previewHandler.go).
- `UpdateTemplate` no longer aliases the caller's template through a shallow
  variant clone; it maps to a fresh entity —
  [templateHandler.go](../internal/handlers/templateHandler.go).
- `helpers.MapSlice` preserves nil-vs-empty, which `linq.SelectSlice` does not —
  [slice.go](../internal/helpers/slice.go). Without it `ContentPools: []` would
  have become `null` on disk.

## 3. Features added / changed

- **`internal/models/template_model/`** — 30 types across six `_model`
  subpackages, `types.go` aliasing all of them, `converters.go` re-exporting the
  five converters outside seams need. No JSON tags anywhere in the tree.
- **`template_model.Zone.Quality *neutral_zone.Quality`** is now the tier. nil =
  "not recorded, infer it"; the pointer is load-bearing because the enum is
  `iota - 1`.
- **`mappers.TemplateMapper` / `ITemplateMapper`**; `EditorStateEntityMapper`
  renamed `EditorStateMapper`.
- **depguard `template-model-inner-private`** in
  [.golangci.yml](../.golangci.yml), proven to fire.
- `Generate()` returns `*template_model.Template`; `ResolveQuality` and
  `PlaceArena` lost their tier-index parameters; `models.GeneratedTemplate`,
  `TemplateLoadDto.ZoneTiers`, `State.lastZoneTiers` and `GetLastZoneTiers()`
  are deleted.

## 4. File modifications

Full detail is in the plan's **Phase 2.5 → Phase Summary**; the shape of it:

- **New**: `internal/models/template_model/**` (36 files),
  `internal/mappers/templateMapper.go` + its interface,
  `test/test_helpers/allFieldsTemplate.go`, `defaultTemplateModel.go`,
  `test/unit/internal/mappers/templateMapper/**` (3 files).
- **Renamed with `Move-Item`**: `internal/mappers/editorStateEntityMapper*.go` →
  `editorStateMapper*.go`; `test/unit/internal/mappers/editorStateEntityMapper/`
  → `editorStateMapper/`.
- **Deleted with `Remove-Item`**: `internal/models/generatedTemplate.go`,
  `test/unit/app/gui/drivers/state/getLastZoneTiers_test.go`.
- **Edited**: the generator, the arena provider, the tier service (+ both
  interfaces), four DTOs, `templateHandler`, `previewHandler`, `drivers.State`
  and its two driver files, `layoutPanelZones.go`, `providerSets.go`,
  regenerated `wire_gen.go`, plus roughly thirty test files.

## 5. Tests added or updated

- **Round-trip drift guard**: `ToEntity(ToModel(x)) == x` over a gofakeit-fuzzed
  all-fields template. Verified by mutation (dropping one converter field fails
  it).
- `resolveQuality_test.go` rewritten for the pointer contract: recorded wins /
  the lowest tier is still "recorded" / nil infers / unclassifiable ⇒ Unknown.
- `generateZoneTiers_test.go` rewritten to assert on
  `generated.Variants[0].Zones[n].Quality`.
- Two tests renamed because their subject changed (see §7).
- One test deleted: the driver no longer preserves a tier index, so
  `TestWhenEditsAreApplied_TheGeneratedZoneTiersAreKept` had no subject; the
  behaviour now lives in `templateHandler.carryZoneTiersByName`.

Last runs: unit / untagged / integration / **GPU (no `-update`)** all exit 0.
Lint 0 issues. Coverage **74.3 %** (was 73.9 %, floor 72.5 %).

## 6. Git status snapshot

Branch `AD/fixing_some_stuff_08-12`. **Nothing was staged by this session.** The
already-staged entries are the owner's phase-2 staging; the `AD` / ` D` markers
are phase-2-staged files this phase deleted or moved
(`internal/models/generatedTemplate.go`, `getLastZoneTiers_test.go`, the three
`editorStateEntityMapper` files). Untracked: `internal/models/template_model/`,
the four new `internal/mappers/*.go`, the two new test helpers and the new
`templateMapper` test folder.

## 7. Rejections / deliberate deviations

- **No per-subpackage test folders** under
  `test/unit/internal/models/template_model/**`. Thirty folders of
  near-identical converter tests would restate the round-trip guard with weaker
  assertions; the converters are covered transitively and §4.6 exempts pure data
  structs. Say the word if you want them anyway — they are mechanical.
- **The mapper was not threaded into `app/`.** `templateHandler.SaveTemplate`
  and `previewHandler.BuildPreviewLayout` flatten internally instead, so
  `NewUIState` / `NewWindow` keep their argument lists. Batch I already tried the
  other way and it was reverted.
- `TestWhenUpdateSucceeds_ReturnedTemplateIsProvidedTemplateInstance` was
  asserting pointer aliasing that no longer exists; renamed
  `..._ReturnedTemplateCarriesTheAppliedZones`.
  `TestWhenTemplateIsProvided_LaysOutThatTemplate` compares the name, not the
  identity.

## 8. Open questions

- Do you want the per-subpackage `template_model` test folders (§7)?
- `BenchmarkEditorWindow_TabCycling` was **not** re-measured. Phase 2.5 touches
  the generation and save paths, not the per-frame path, so a regression is
  unlikely — but the two figures in the older docs (4,773 and 6,640) disagree
  with each other and phase 2 measured **5,699 allocs/op**. Worth one manual run
  before phase 3 if you want a clean baseline.

## 9. Next recommended actions

1. Review phase 2.5, stage and commit it.
2. **Phase 3** — persist the tier in `.gen.json` (`*int8` on the entity,
   `*neutral_zone.Quality` on the model). It is unblocked. This is also where the
   `Unknown → Plastic` silent down-tier described in the plan's §1 actually gets
   fixed.
3. Phase 4 — sweep the production and test files still naming `entities.*`, and
   delete `templateHandler.carryZoneTiersByName`.
4. Phase 5 — shrink the layering allow-lists and record the outcome.

## 10. Carry-forward prompt

> Read `AGENTS.md` first — especially **§4.4.1**, the Entity/Model/DTO doctrine:
> **Entity** (`internal/entities/`) is the wire format, json tags only; **Model**
> (`internal/models/`) owns the structure and all business logic; **DTO**
> (`internal/dtos/`) is the thin `app/` ↔ `internal/` crossing. A DTO carrying a
> Model is intended and `app/` may hold a Model — do not "fix" either.
>
> Then read `.agent/plans/batch-j-zone-tier-source-of-truth.md` and do **phase
> 3**. Read **§0 and §0b — where they disagree, §0b wins.** Phases 1, 2 and 2.5
> are complete and uncommitted. The template layer now lives in
> `internal/models/template_model/`, a full mirror of
> `internal/entities/template/`; the zone tier is
> `template_model.Zone.Quality *neutral_zone.Quality`, nil meaning "not recorded,
> infer it". Read phase 2.5's Phase Summary before touching anything — it records
> the per-type embed/redefine rule, the `converters.go` function-re-export trick
> (type aliases cross a package boundary, functions do not), and the two
> behavioural changes that renamed tests.
>
> Hard rules, one line each: never modify `data/`, `internal/registry/`, or
> anything under `internal/entities/template/` — `internal/entities/editor_state/`
> is *not* protected and phase 3 changes it; everything must build and run on
> Windows and Linux (`path/filepath`; chain PowerShell with `;`, never `&&`);
> every change ships with tests and unit coverage must not drop below 72.5 %
> (currently **74.3 %**); the lint baseline is **0 issues**; **never stage and
> never commit** — use `Move-Item` never `git mv`, `Remove-Item` never `git rm`;
> never change where `.rmg.json` is written and never persist the output
> directory; never run a bulk in-place rewrite, and **never round-trip a `.go`
> file through `Get-Content`/`Set-Content`** — use `gofmt -r` on an explicit file
> list and verify insertions == deletions.
>
> Standing traps: **nil is load-bearing** three times over — nil `Previous` =
> first generation, nil `Next` = unarmed debounce, nil `Quality` = "infer it";
> phase 3's persisted tier must be `*int8` because `omitempty` on a plain `int8`
> would silently drop every Plastic zone; the two frozen fixtures under
> `test/test_helpers/testdata/` plus the untagged
> `editorStateWireFormat_integration_test.go` must keep passing **unchanged**,
> comparing **parsed objects, never bytes**; the picker and zone-editor dialogs
> are snapshotted, so the GPU suite must pass **without `-update`**;
> `cmd/testlayoutcheck` matches test-only export names tree-wide, so grep any new
> accessor name before adding it; and `helpers.MapSlice` / `helpers.MapPointer`
> preserve nil-vs-empty where `linq.SelectSlice` does not — never swap them
> inside a converter. Cap sessions at ~50 messages and hand off through this file.
>
> Full handoff in `./.agent/session-carry-forward.md`.
