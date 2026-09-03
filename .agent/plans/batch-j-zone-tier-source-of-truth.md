# Batch J — zone tier single source of truth (backlog §2.2 branch B)

A neutral zone's tier is decided once at plan time, flattened into layout/pools/
castles, and then reverse-engineered by `ZoneClassifier` everywhere it is needed.
This batch records the tier instead: the generator hands back what it planned on
the zone itself, a `template_model` layer carries it through the editor,
`.gen.json` persists it, and a single `IZoneTierService` answers the question —
falling back to inference only for a template loaded from a raw `.rmg.json`.

> **The wrapper design changed at the phase 2 review.** The original plan used a
> standalone `models.QualifiedZone` plus a name-keyed tier index. That was
> rejected; see **§0b** for what replaced it and which §0 decisions it
> supersedes. Where the two disagree, **§0b wins**.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is
done, set its status to `Complete` and write its **Phase Summary** (what was
done, key decisions, anything needed to continue with zero context); run the
phase's **Verification Plan** and record the result before moving on. When all
phases are done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../../AGENTS.md) first — especially **§4.4.1** (Entity/Model/DTO
doctrine) and **§2.1** (protected directories). Then read **§0 and §0b**: where
they disagree, **§0b wins**, because it records the phase 2 review. Hard rules
that bite here: **`internal/entities/template/**` is protected and this batch
must not touch a single byte of it**; never stage and never commit; move with
`Move-Item`, never `git mv`; delete with `Remove-Item`, never `git rm`; never
round-trip a `.go` file through `Get-Content`/`Set-Content`; unit coverage must
not drop below **72.5 %** (currently **74.3 %**); lint baseline is **0 issues**.

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

## 0b. Amendment (settled with the owner 2026-09-02, after the phase 2 review)

The phase 2 review rejected `models.GeneratedTemplate`. The verdict, in the
owner's words: it *"shifts the goal post"* rather than moving away from using the
template entity structs directly, unlike what was already done for
`EditorState`. A tier index keyed by zone name is a side-car; the tier is a
property **of the zone** and belongs on a zone model. Phase 2.5 corrects this
before phase 3, because phases 3-5 and several backlog items build on top of it.

10. **`models.GeneratedTemplate` is deleted.** The tier rides on the zone, not in
    a parallel map. This **supersedes §0.4's return type**: `Generate()` returns
    `(*template_model.Template, []string)`. §0.4's *other* half stands — the test
    call sites take the churn, and no second production entry point exists.
11. **`internal/models/template_model/` mirrors `internal/entities/template/`
    one for one** — the whole schema, not just the three types the tier needs.
    Same package structure, every subpackage suffixed `_model`:

    | Entity package | Model package |
    | --- | --- |
    | `entities/template` | `models/template_model` |
    | `entities/template/template_common` | `models/template_model/template_common_model` |
    | `entities/template/template_content` | `models/template_model/template_content_model` |
    | `entities/template/template_layout` | `models/template_model/template_layout_model` |
    | `entities/template/template_override` | `models/template_model/template_override_model` |
    | `entities/template/template_rule` | `models/template_model/template_rule_model` |
    | `entities/template/template_variant` | `models/template_model/template_variant_model` |

    `template_model/types.go` re-exports every subpackage type exactly as
    [entities/template/types.go](../../internal/entities/template/types.go) does,
    so callers write `template_model.Zone` and never name an inner package. A
    depguard rule **`template-model-inner-private`** enforces it, mirroring the
    existing `template-inner-private`. This **supersedes §0.2** —
    `QualifiedZone` is never built — and **supersedes §0.3**: the model is its
    own store, so `drivers.State` needs no side index.
12. **The model mirrors the entity, it does not wrap it.** `editor_state_model`
    already shows three shapes and `template_model` reuses all three, choosing
    per type by one question: *does any field need re-typing?*
    - **Re-types children → no embedding.** The model declares its own fields
      and the converters map each one. Precedents:
      [editorState.go](../../internal/models/editor_state_model/editorState.go),
      [contentSettings.go](../../internal/models/editor_state_model/contentSettings.go),
      [manualEditSettings.go](../../internal/models/editor_state_model/manualEditSettings.go).
      **`Template`, `Variant`, `Zone`, `Connection` and every other composite
      type take this shape** — `RmgTemplate` is explicitly *not* embedded in
      `template_model.Template`.
    - **All-scalar → embed the entity.** Precedent:
      [templateIdentity.go](../../internal/models/editor_state_model/templateIdentity.go).
      Free promotion, trivial converter, nothing to keep in sync.
    - **`ZoneContentRow` shadow** — embed and shadow one collection, nil'ing the
      embedded copy. Reserved for a mostly-scalar struct that still needs the
      entity value itself for a helper call. Use it only if such a case turns up.

    `Zone` additionally gains the tier. **Model types carry no JSON tags and no
    (un)marshalling**: serialization is the entity's job (§4.4.1). The three
    entity types with decode behaviour — `StringList`, `BonusList` and
    `GameRules.UnmarshalJSONFrom` — keep it, and their models must not duplicate
    it.
13. **The tier field is `Quality *neutral_zone.Quality`**, nil meaning "not
    recorded, infer it". This is not stylistic: `Quality` is `iota - 1`, so
    `QualityUnknown` is −1 and **the zero value is `QualityLowest`**. A plain
    value field would make every zero-valued zone literal silently claim Plastic
    — the precise bug this batch exists to kill. It also matches §0.6 and the
    existing `ManualPosition *[2]float64`.
14. **The protected schema stays untouched in 2.5.** Moving `GeneratorPosition`,
    `GeneratorRing`, `ManualPosition` and `IsUserAdded` out of
    `internal/entities/template/**` and onto `template_model` is approved **in
    principle** and is the right end state — they are runtime concerns squatting
    in a serialization schema — but it happens in **its own reviewed phase**,
    after every reader goes through `template_model`. Until then AGENTS §2.1
    applies in full: not one byte of `internal/entities/template/**` changes.
15. **Mappers.** `EditorStateEntityMapper` is renamed `EditorStateMapper` (and
    `IEditorStateEntityMapper` → `IEditorStateMapper`). A new `TemplateMapper`
    owns `RmgTemplate(Entity) ⇄ Template(Model)`. `Template(Model) →
    EditorState(Model)` lives in **`EditorStateMapper`**, not in `TemplateMapper`:
    `RmgTemplate` is the source of truth, so each consumer state pulls *from* it.
    A future TUI or Web state adds its own mapper rather than another branch
    inside `TemplateMapper`. The full chain is
    `RmgTemplate(Entity) ⇄ Template(Model) ⇄ EditorState(Model) ⇄ EditorState(Entity)`.
16. **Phase 2.5 lands the whole model package, but sweeps no consumer.** All 30
    schema types get their model, the mappers get written, and the
    generator / handler / `drivers.State` seam moves onto `template_model`.
    The 100 production and 266 test files that name `entities.*` are **not**
    rewritten here — that stays phase 4, which `template_model.Zone` shrinks
    since it replaces the `QualifiedZone` swap that phase was written around.

    ⚠ **Without embedding there is no free promotion**, so a seam that is not
    migrated cannot silently keep compiling. That is a feature: the compiler
    lists the work. Bridge any seam phase 2.5 does not intend to move with
    `TemplateMapper.ToEntity` and note it, rather than widening the phase. The
    entity is genuinely required at exactly two places — reading and writing
    `.rmg.json` — and those two keep it forever.

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

## Phase 2: the generator records what it planned
Status: **Complete** (2026-09-02) — reviewed, staged, **and partly superseded by
phase 2.5.** The `planZoneTiers` derivation and the `ResolveQuality` precedence
rule survive; `models.GeneratedTemplate`, `TemplateLoadDto.ZoneTiers` and
`State.lastZoneTiers` do not. Read this phase as history, not as a spec.

Behaviour deltas start here, but only **one**: the gladiator arena. Everything
else is plumbing.

### How the tier gets out of the generator (settled — do not re-derive)

The scoping pass found a much cheaper route than threading a collector through
every topology service. **Every** neutral zone name in the repository is built by
`constants.GetNeutralZoneNameFor(plan.Label)` — including the tournament cluster
services — and labels are unique because zone names must be. So the index can be
derived in `TemplateGenerator.Generate` from the `neutral_zone.Plans` it already
holds, with **zero changes to `ZoneFactory`, `TopologyBase` or any topology**:

```go
// planZoneTiers records the tier the generator chose per zone name. Hub zones
// are always built from the Highest profile; spawn zones have no tier.
func planZoneTiers(
    neutralZones neutral_zone.Plans,
    variant entities.Variant) map[string]neutral_zone.Quality
```

Build a `label-name → quality` map from the plans, then walk `variant.Zones` and
record only names that actually exist: hubs (`zone_helpers.IsZoneNameHub`) as
`QualityHighest` — matching `CreateHubZone`, which always uses the Highest
profile — and anything matching a plan name at its planned quality. Spawn zones
are deliberately absent. Compute it **before** `PlaceArena` (which mutates zones
but adds none).

⚠ **Comma-ok is mandatory on every lookup.** A missing key yields `Quality(0)`,
which is `QualityLowest` — a silent down-tier, the same class of bug as the
`omitempty` trap in phase 3. Always `quality, ok := tiers[name]`, and fall back
to inference when `!ok`. This also makes a `nil` map safe, which matters because
test call sites will pass one.

### The precedence rule lives in exactly one place

Add to `IZoneTierService`:

```go
ResolveQuality(zone entities.Zone, zoneTiers map[string]neutral_zone.Quality) neutral_zone.Quality
```

Recorded tier wins; inference is the fallback. Phases 3 and 4 reuse it rather
than re-implementing the precedence at each consumer.

### `QualifiedZone` is deferred to phase 4

The wrapper has **no consumer** until the editor chain moves, and phase 2's one
consumer (the arena provider) wants the map, not a wrapper — it mutates
`variant.Zones` in place by index, which a wrapper only gets in the way of.
Building it now would be an abstraction ahead of its caller (§3.1). The design
is unchanged; only the file's arrival moves. Its checklist item lives in phase 4.

### Checklist

- [x] `internal/models/generatedTemplate.go` — `GeneratedTemplate{Template
      *entities.RmgTemplate; ZoneTiers map[string]neutral_zone.Quality}`.
- [x] `TemplateGenerator.Generate()` returns
      `(*models.GeneratedTemplate, []string)`, building the index via
      `planZoneTiers`. Update `templateGeneratorInterface.go` and
      `test_helpers/templateGeneratorMock.go`.
- [x] `IZoneTierService.ResolveQuality` + its unit tests (recorded wins,
      unrecorded infers, nil map infers).
- [x] `PlaceArena(configuration, variant, zoneTiers)` — the provider resolves
      through `ResolveQuality` in `findRichestNeutralZoneIndex` and
      `mapNeutralZoneQualities`. Update the interface and the 12 call sites in
      `placeArena_test.go`. **This is the behaviour delta**: a zone that
      inference called `Unknown` scored −1 and could never win the arena; with
      its planned tier it can.
- [x] Carry the index to the GUI: `dtos.TemplateLoadDto` gains `ZoneTiers`,
      `templateHandler.GenerateTemplate` fills it, and `drivers.State` stores it
      beside `lastTemplate` (`setLastTemplate` is the only writer — keep it that
      way so `templateRevision` stays correct).
- [x] Tests: the index contains a planned neutral at its planned tier, a hub at
      Highest and **no** spawn entry; plus an arena test proving a previously
      `Unknown` zone now wins.

### The ~73 test call sites

`Generate()` has **one** production caller and ~73 test ones, nearly all
`actual, _ := generator.Generate()` followed by `actual.Variants[...]`. They
cannot be rewritten with `gofmt -r` alone, because the fix needs a second
statement, and **a PowerShell text sweep over `.go` files is forbidden**.

The safe route: all eight files in
`test/unit/internal/services/template_generator/templateGenerator/` are one
package, so add a **test-local** helper there —

```go
func generateTemplate(generator ...) (*entities.RmgTemplate, []string)
```

— and rewrite with one AST-safe pass per file:
`gofmt -r 'generator.Generate() -> generateTemplate(generator)'`. This is *not*
the rejected `GenerateWithTiers()`: that was a **production** entry point added
to dodge churn. A test-package unpack helper is ordinary test hygiene and the
repo already uses `common_test.go` helpers. Tests that assert on tiers call
`Generate()` directly. The two `test/performance` sites are edited by hand.

**Expect generated output to move here** — but only for the arena, and only for
zones inference called `Unknown`. Enumerate every delta: which topology, which
zone, and why the recorded tier is the correct answer. A delta you cannot
explain is a bug.

### Verification Plan
- Full gate set from phase 1.
- Diff a generated `.rmg.json` for each topology before/after; every difference
  must map to a zone inference called `Unknown`.
- `BenchmarkEditorWindow_TabCycling` re-measured **by hand** (needs a GPU, never
  runs in CI) if the clone path or `State` shape changed — baseline ~4,773
  allocs/op.

### Phase Summary

**Landed as designed.** `Generate()` now returns
`*models.GeneratedTemplate` — the template plus
`ZoneTiers map[string]neutral_zone.Quality` — and the index rides
`dtos.TemplateLoadDto` into `drivers.State`, which keeps it beside
`lastTemplate`. `IZoneTierService.ResolveQuality(zone, zoneTiers)` is the single
place the precedence rule lives: recorded tier wins, inference is the fallback,
and the lookup is comma-ok so a missing key never reads back as `QualityLowest`.
`PlaceArena` is its only consumer so far.

**The index needed no topology changes**, exactly as scoped. `planZoneTiers`
lives in `templateGenerator.go` beside `Generate`, builds a
`label-name → quality` map from the `neutral_zone.Plans` the generator already
holds, then walks `variant.Zones` recording hubs at `QualityHighest` and any zone
matching a plan name at its planned quality. `ZoneFactory`, `TopologyBase` and
every topology service were untouched.

**⚠ Correction to decision §0.8: phase 2 moves no generated output at all.** The
delta was expected to be "the arena, for zones inference called `Unknown`". A
throwaway sweep over **792 configurations** — 11 topologies × {0,2,4,6} neutral
zones × {2,4,8} players × simple/advanced tiers × tournament on/off ×
{0,1,2} castles — compared the arena placement computed from the recorded index
against the placement computed from inference alone, and additionally compared
`GetQuality(zone)` against the recorded tier for **every** zone of every
generated variant. Both counts were **zero**: no arena moved, and inference and
the plan agree on every zone the generator produces. The generator never emits a
zone its own content pools cannot classify, so the `Unknown` case the arena delta
was predicted for does not exist on the generation path. The scratch test was
deleted after the measurement.

That is a stronger result than the plan expected, not a weaker one: it means the
recorded tier is a *faithful* replacement for inference on generated templates,
and the correction it buys is latent — it fires only for zones inference cannot
classify, which is what a manually re-tiered zone (phase 3) or a raw `.rmg.json`
(phase 4) produces. Two new unit tests pin that latent behaviour down directly
against the provider:
`TestWhenTheRichestZoneCannotBeInferred_TheRecordedTierWinsTheArena` and
`TestWhenConnectionEndpointsCannotBeInferred_TheRecordedTiersPickTheConnection`.
**Consequence for phase 3:** it, not phase 2, is where the `Unknown → Plastic`
silent down-tier described in §1 actually gets fixed.

**The ~73 test call sites went the planned route.** A test-local
`generateTemplate(generator)` unpack helper in
`test/unit/.../templateGenerator/common_test.go`, then one
`gofmt -r 'a.Generate() -> generateTemplate(a)'` pass per file over an explicit
list of the eight files — an AST rewrite, never a text sweep. The diff came back
**71 insertions / 71 deletions**, i.e. line-for-line, which is the check that the
rewrite touched nothing else. `placeArena_test.go`'s 12 sites took the same
treatment with `a.PlaceArena(b, c) -> a.PlaceArena(b, c, nil)` (12/12), and the
two `test/performance` sites plus the four `templateHandler` mock returns were
edited by hand.

**`QualifiedZone` was not built**, per the design note — phase 2 has no consumer
for it and the arena provider wants the map, since it mutates `variant.Zones` in
place by index. Its checklist item stays in phase 4.

**One decision the plan left open.** `State.handleUpdateTemplate` (the manual-edit
apply path) calls `setLastTemplate(dto.Template, this.lastZoneTiers)` — it
*preserves* the index rather than clearing it, because a manual edit reshapes the
template that is already loaded and the tiers planned for it still describe it.
Nothing reads `State.lastZoneTiers` yet, so this cannot misbehave in phase 2;
phase 3 is where a manually re-tiered zone must **overwrite** its entry, and the
preserve-by-default behaviour is what makes that a one-line change instead of a
restoration.

### Verification results (2026-09-02)

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 (no DI change — `Generate` is called through an interface the graph already provides) |
| Unit / untagged / integration | pass |
| **GPU suite, no `-update`** | **pass (24.7 s)** — no pixel moved |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **73.9 %** (was 73.8 %, floor 72.5 %) |
| Per-function coverage of new code | `planZoneTiers`, `ResolveQuality`, `PlaceArena`, `findArenaConnectionIndex`, `findRichestNeutralZoneIndex`, `mapNeutralZoneQualities`, `setLastTemplate`, `applyGeneratedTemplate`, `GetLastZoneTiers` — **100 % each** |
| Output diff, 792 configurations | **0 arena moves, 0 inference/plan disagreements** (see summary) |
| `BenchmarkEditorWindow_TabCycling` | **5,699 allocs/op** (3 runs: 5,702 / 5,698 / 5,699; 2.31–2.88 M ns/op) |

⚠ **The benchmark baseline in this plan was already stale.** It cites ~4,773
allocs/op, while backlog §1.4 records **6,640** after the clone batch that
followed that measurement. The number now measures 5,699 and is stable to ±4
across runs. Phase 2 adds two small maps per *generation* (the benchmark cycles
tabs and generates roughly once), so it cannot account for a ~900 allocs/op
shift in either direction; the two documented baselines simply disagree with each
other. Someone should reconcile them — do not read this row as a phase 2
regression.

---

## Phase 2.5: `template_model` — mirror the schema, put the tier on the zone
Status: **Complete** (2026-09-02) — uncommitted, awaiting review.

Phase 2 solved the right problem the wrong way. `models.GeneratedTemplate` keeps
the application talking in `entities.RmgTemplate` and bolts a
`map[zoneName]Quality` alongside it, so three disjoint fields have to be dragged
around together (`lastTemplate`, `lastZoneTiers`, `templateRevision`) and every
tier read is a name lookup instead of a field access.

This phase does for the template what was already done for the editor state:
**`internal/models/template_model/` mirrors `internal/entities/template/` in
full**, subpackage for subpackage, so the service layer owns the structure and
the entity goes back to being nothing but the `.rmg.json` wire format. Once the
zone is a model, the tier is simply a field on it, and the map has nowhere left
to live.

The payoff is not only the tier. `Template` becomes extendable without touching
the protected schema, which is what §0b.14's field migration and several backlog
items depend on.

Read **§0b** first — it holds the seven decisions this phase implements. Nothing
here is open for redesign.

### The shape, stated once

`Template` **does not embed** `RmgTemplate`; it mirrors it, exactly as
`editor_state_model.EditorState` mirrors `editor_state.EditorState`:

```go
// internal/models/template_model/template.go
type Template struct {
    Name                string
    GameMode            string
    Description         string
    DisplayWinCondition string
    SizeX               int
    SizeZ               int

    ValueOverrides []ValueOverride
    Orientation    *Orientation
    Border         *Border
    GameRules      GameRules
    GlobalBans     *GlobalBans
    Variants       []Variant

    ZoneLayouts        []ZoneLayoutDef
    MandatoryContent   []MandatoryContent
    ContentCountLimits []ContentCountLimit
    ContentPools       []ContentPool
    ContentLists       []ContentList
}

// internal/models/template_model/template_variant_model/zone.go
type Zone struct {
    // ... every field of template_variant.Zone, composites re-typed ...
    Quality *neutral_zone.Quality   // nil = not recorded, infer it
}
```

All-scalar types (`Noise`, `TypedRef`, `PlacementRule`, `ElevationMode`, …) embed
their entity instead, per §0b.12. Which shape a given type takes is decided by
one question, not by taste: **does any of its fields need re-typing?**

### The inventory

30 types, mirrored one for one. Nothing is skipped — a partial mirror would leave
callers switching between `template_model.X` and `entities.Y` mid-expression,
which is the confusion this phase exists to end.

| Model package | Types |
| --- | --- |
| `template_model` | `Template` |
| `template_common_model` | `PlacementRule` |
| `template_content_model` | `ContentCountLimit`, `ContentLimit`, `ContentList`, `ContentPool`, `MandatoryContent`, `MandatoryContentItem`, `WeightedContent` |
| `template_layout_model` | `AmbientPickupDistribution`, `ElevationMode`, `GuardedEncounterResourceFractions`, `ZoneLayoutDef` |
| `template_override_model` | `ValueOverride` |
| `template_rule_model` | `Bonus`, `BonusList`, `GameRules`, `GlobalBans`, `WinConditions` |
| `template_variant_model` | `Border`, `Connection`, `EncounterHolesSettings`, `MainObject`, `Noise`, `Orientation`, `Road`, `StringList`, `TypedRef`, `Variant`, `Zone` |

One struct per file, file named after the struct in camelCase (§4.1), converters
beside the struct they convert — the `editor_state_model` layout verbatim.

### Why this deletes the map rather than hiding it

The name→quality map does not vanish from the source; it stops being an **API**.
`planZoneTiers` currently returns it to `Generate`, which returns it to the
handler, the DTO and `drivers.State`. After this phase the same derivation is a
**local variable inside one function** — the converter that lifts the generated
`entities.Variant` into a `template_model.Variant` and stamps `Quality` on each
zone as it goes. It crosses no boundary and no type carries it.

### Checklist

Order that keeps the tree compiling at every step.

- [x] `internal/models/template_model/` and its six `_model` subpackages — all
      30 types from the inventory above, one struct per file, each with its
      `To*Model` / `To*Entity` converters. Shape per type by §0b.12. **No JSON
      tags on any model type.**
- [x] `internal/models/template_model/types.go` re-exporting every subpackage
      type, mirroring
      [entities/template/types.go](../../internal/entities/template/types.go).
- [x] `.golangci.yml` gains `template-model-inner-private` beside
      `template-inner-private`:

      ```yaml
      template-model-inner-private:
        files:
          - "!**/internal/models/template_model/**"
          - "!$test"
        deny:
          - pkg: github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/
            desc: import internal/models/template_model instead
      ```

      Note the exclusion is the **model tree only**, tighter than the entity
      rule's `!**/internal/entities/**` — the subpackages import each other, but
      no other model package has any business reaching inside.
- [x] `internal/mappers/templateMapper.go` + `ITemplateMapper`:
      `ToModel(entity) Template`, `ToEntity(model) RmgTemplate`. Register it in
      `providerSets.go` and regenerate — never hand-edit `wire_gen.go`.
- [x] Rename `EditorStateEntityMapper` → `EditorStateMapper` and
      `IEditorStateEntityMapper` → `IEditorStateMapper`: both files, the struct,
      the constructor, `test_helpers` mock, and the unit-test folder
      `test/unit/internal/mappers/editorStateEntityMapper/` →
      `.../editorStateMapper/`. **`Move-Item`, never `git mv`.** Regenerate wire.
- [x] `TemplateGenerator.Generate()` returns `(*template_model.Template, []string)`.
      `planZoneTiers` stops being a returned index and becomes the local
      derivation inside the entity→model conversion (see above). `PlaceArena`
      runs **after** the conversion, on the model variant.
- [x] `IZoneTierService.ResolveQuality(zone template_model.Zone) neutral_zone.Quality`
      — the map parameter is gone; nil `Quality` falls back to inference.
      `GetQuality` keeps taking the **entity** zone: a template loaded from a raw
      `.rmg.json` has no recorded tier and never will.
- [x] `IGladiatorArenaProvider.PlaceArena(configuration, variant *template_model.Variant)`
      — no tier parameter. It still mutates `variant.Zones` in place by index.
- [x] `dtos.TemplateLoadDto.Template` becomes `*template_model.Template` and
      **`ZoneTiers` is deleted**. `templateHandler.GenerateTemplate` fills it.
- [x] `templateHandler.UpdateTemplate` needs a **temporary wrap-by-name seam**:
      its `TemplateUpdateDto.Zones` is still `[]entities.Zone` until phase 4
      re-types the 9 DTOs, so it must carry the previous model zones' `Quality`
      across by name. Confine it to that one function, comment it as phase-4
      scaffolding, and do not let it spread. This is deliberately the *only*
      place a name lookup survives.
- [x] `drivers.State`: `lastTemplate *template_model.Template`; **delete**
      `lastZoneTiers`, `GetLastZoneTiers()` and `setLastTemplate`'s second
      parameter. **`templateRevision` stays** — it is the preview cache key
      ([previewLayoutCache.go](../../app/gui/models/previewLayoutCache.go)),
      answering "was the pointer replaced", which is a different question from
      "what tier is this zone". Do not remove it as part of the tidy-up.
- [x] Bridge, do not widen: every consumer the compiler now breaks either moves
      to `template_model` **because it is on the seam** (the handler, the DTO,
      `drivers.State`) or gets a `TemplateMapper.ToEntity` call and a note that
      phase 4 removes it. `preview_service`, `file_service` and the GUI panels
      are expected to take the bridge, not the migration.
- [x] **Delete `internal/models/generatedTemplate.go`.**
- [x] Tests: new folders under `test/unit/internal/models/template_model/**` and
      `test/unit/internal/mappers/templateMapper/` per §4.6;
      `resolveQuality_test.go` reshaped to the pointer contract (recorded wins /
      nil infers); `getLastZoneTiers_test.go` deleted; the generator tier tests
      assert on `generated.Variants[0].Zones[n].Quality` instead of an index.
- [x] The `generateTemplate(generator)` unpack helper in the templateGenerator
      test package becomes a **`TemplateMapper.ToEntity`** call. That keeps all
      eight rewritten test files — including the golden-template assertion —
      compiling and asserting against `entities.RmgTemplate` unchanged, which is
      also the strongest possible proof that the round trip is lossless.

### Verification Plan

- Full gate set from phase 1.
- **Generated output must be byte-identical.** This is a pure restructuring:
  `TemplateMapper.ToEntity(Generate())` must equal what `Generate()` produced
  before, for every topology. `TestWhenDefaultConfiguration_ReturnsGoldenTemplate`
  is the primary guard; re-run the 792-configuration sweep from phase 2 if any
  doubt remains.
- **Round-trip test**: `ToModel` → `ToEntity` on a fully populated template is
  the identity. Build the fixture with `gofakeit` so a field added to the schema
  and forgotten in the mapper fails the test instead of passing silently — a
  30-type hand-written mapper is exactly where a dropped field hides.
- `golangci-lint-v2 run ./...` proves `template-model-inner-private` fires:
  temporarily import a `_model` subpackage from outside the tree, confirm the
  issue, then revert. A rule that was never seen to fail is not a rule.
- `wire diff ./internal/composition/...` exit 0 after the mapper rename and the
  new provider.
- **GPU suite without `-update`** — no pixel may move.
- Coverage ≥ 72.5 %. ⚠ 30 mostly-mechanical converter pairs will **drag the
  percentage down** unless they are tested; budget for that rather than
  discovering it at the end.
- `go run ./cmd/testlayoutcheck .` — and grep any new test-only accessor name
  tree-wide before adding it.

### Phase Summary

**Landed as designed.** `internal/models/template_model/` now mirrors
`internal/entities/template/` one for one — all 30 types across the six
`_model` subpackages, `types.go` re-exporting every one of them, and no JSON tag
anywhere in the tree. `Template` does not embed `RmgTemplate`; the per-type rule
from §0b.12 decided each shape mechanically, and it split the **13 embedded**
all-scalar types (`PlacementRule`, `WeightedContent`,
`AmbientPickupDistribution`, `ElevationMode`,
`GuardedEncounterResourceFractions`, `ValueOverride`, `Bonus`, `GlobalBans`,
`WinConditions`, `Noise`, `TypedRef`, `Orientation`, `EncounterHolesSettings`)
from the composites that re-type a child and are therefore declared in full.
`BonusList`, `StringList`, `ContentList` and `ContentPool` are named slice/map
types, which cannot embed at all, so they are re-declared with the same shape
and **without** the tolerant decode behaviour — that stays on the entity, where
the wire format is.

**The map is gone, not hidden.** `planZoneTiers` became
`stampPlannedZoneTiers(neutralZones, zones []template_model.Zone)`, a package
function in `templateGenerator.go` that writes `Quality` onto each zone and
returns nothing. `Generate()` builds the entity exactly as before, hands it to
`TemplateMapper.ToModel`, stamps the tiers, then runs `PlaceArena` on the model
variant. Building the entity first and lifting it — rather than assembling the
model field by field from the providers, which all still return entities — is
what makes the output guaranteed-identical rather than argued-identical.
`models.GeneratedTemplate`, `TemplateLoadDto.ZoneTiers`, `State.lastZoneTiers`
and `GetLastZoneTiers()` are all deleted.

**`ResolveQuality` lost its parameter** and reads `zone.Quality` directly, nil
meaning "infer it". `GetQuality` still takes the entity, so the fallback path
converts the one zone it needs — a generation-time cost, not a per-frame one.

**Four DTOs moved, three seams bridged.** `TemplateLoadDto`, `TemplateSaveDto`,
`TemplateUpdateDto` and `PreviewLayoutRequestDto` now carry
`*template_model.Template`. The two places that genuinely need the wire format
convert **inside the handler**, not in `app/`: `templateHandler.SaveTemplate`
flattens before the preview generator and the file service, and
`previewHandler.BuildPreviewLayout` flattens before the layout service. That
kept `app/` free of the mapper entirely and left `NewUIState`'s argument list
alone — the alternative was threading a mapper through `program.go` →
`editor.NewWindow` → `drivers.NewUIState`, which batch I already tried and the
owner reverted. The one place `app/` still converts is
`layoutPanelZones.handleConnectionEditorClick`, which flattens the variant for
the zone editor dialog with `template_model.ToZoneEntities`; that is a model
package function, which §4.4.1 explicitly allows `app/` to call.

**Functions do not cross a package boundary the way type aliases do.** That is
the one thing `types.go` cannot solve: `template_model.Zone` is an alias and
works everywhere, but `ToZoneEntities` lives in `template_variant_model`, which
`template-model-inner-private` forbids naming. So `converters.go` re-exports the
five converters the seams outside the tree actually need — zones, connections
and `ToMainObjectModel` for the arena provider. Everything else stays internal.

**The wrap-by-name seam is one function, `carryZoneTiersByName`**, at the bottom
of `templateHandler.go` and commented as phase-4 scaffolding. It exists only
because `TemplateUpdateDto.Zones` is still `[]entities.Zone`; when the zone
editor moves onto `template_model` it is deleted, not generalised. Note that the
driver no longer preserves anything itself — `handleUpdateTemplate` just calls
`setLastTemplate(dto.Template)` — so the unit test that used to assert the
driver kept the index was **deleted**, its subject having moved into the
handler.

**Two behavioural facts fell out of the restructure, both improvements:**

- `UpdateTemplate` no longer aliases the caller's template. It used to copy the
  struct and `slices.Clone` the variants, which still shared zone slices with
  the source; now it maps to a fresh entity. The unit test
  `TestWhenUpdateSucceeds_ReturnedTemplateIsProvidedTemplateInstance` was
  asserting that aliasing, so it was renamed
  `..._ReturnedTemplateCarriesTheAppliedZones` and pointed at the applied slice.
  `templateHandler`'s existing `..._LeavesTheSourceTemplateUntouched` covers the
  new guarantee.
- `BuildPreviewLayout` no longer hands the layout service the caller's pointer,
  so `TestWhenTemplateIsProvided_LaysOutThatTemplate` compares the name instead
  of the identity.

**Nil versus empty is preserved deliberately.** Every slice converter goes
through the new `helpers.MapSlice`, which returns nil for nil and an empty slice
for empty — `linq.SelectSlice` collapses both to nil, which would have turned
`ContentPools: []` into `null` on disk. `helpers.MapPointer` does the same for
the optional pointer fields. There is a dedicated round-trip test for it.

**The round-trip guard is the real test of the 30 converter pairs.**
`test_helpers.NewAllFieldsTemplate()` fuzzes a whole `RmgTemplate` with a pinned
`gofakeit` seed and fills the four loosely-typed corners (`ContentPool`,
`ContentList`, and both `PlacementRule.Args` sites) by hand, because gofakeit
cannot invent a value for an `any`. `ToEntity(ToModel(x)) == x` then covers every
converter at once. **Verified by mutation**: dropping
`Connection.GuardMatchGroup` from `ToConnectionEntity` fails it, and it passes
again once restored. This is also why coverage went *up* rather than down — the
mechanical converters are all exercised.

**Test-layout deviation, flagged not buried.** Per-subpackage test folders under
`test/unit/internal/models/template_model/**` were **not** created: thirty
folders of near-identical `To*Model`/`To*Entity` tests would restate the
round-trip guard thirty times with weaker assertions. The converters are covered
transitively through `test/unit/internal/mappers/templateMapper/`, and the
package holds pure data structs plus their conversions, which §4.6 already
exempts from per-file tests. If the owner wants the folders anyway, they are
mechanical to add.

### Verification results (2026-09-02)

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 (regenerated, never hand-edited) |
| Unit / untagged / integration | pass |
| **GPU suite, no `-update`** | **pass (29.9 s)** — no pixel moved |
| `golangci-lint-v2 run ./...` | **0 issues** |
| **`template-model-inner-private` proven to fire** | a temporary `_model` import in `internal/mappers` produced the depguard issue; reverted |
| Generated output | `TestWhenDefaultConfiguration_ReturnsGoldenTemplate` green *through* `TemplateMapper.ToEntity`, so the round trip reproduces the previous bytes |
| Round-trip guard | verified by mutation (dropped field ⇒ failure) |
| Unit coverage | **74.3 %** (was 73.9 %, floor 72.5 %) |

⚠ **Not re-measured:** `BenchmarkEditorWindow_TabCycling`. Phase 2.5 adds a
model⇄entity conversion on the *generation* and *save* paths, not on the
per-frame path, and `GetLastTemplate` is still a bare pointer read. The figure
to compare against is phase 2's **5,699 allocs/op** — not the 4,773 or 6,640 the
older docs cite, which disagree with each other.

---

## Phase 3: persist the tier in `.gen.json`
Status: **Complete** (2026-09-03) — uncommitted, awaiting review.

⚠ **Scope decision taken at the start of this phase, with the owner.** The
checklist below could not be executed as written. Persisting the tier is indeed
one conversion at one seam, but the *write* path — `ApplyNeutralZoneQuality` and
`NewDefaultNeutralZone` — is reached only through the zone editor, whose DTOs
and dialog still spoke `entities.Zone`, which has nowhere to put a tier. A
manual re-tier therefore could not reach the file at all, and
`carryZoneTiersByName` would have re-attached the *stale planned* tier by name.
As this phase's own last checklist item anticipated, **the zone-editor half of
phase 4's DTO re-typing came first**. `preview_service`, `file_service`, the
generator and the whole topology tree keep their `TemplateMapper.ToEntity`
bridges — those are still phase 4.

- [x] `internal/entities/editor_state/manualZoneSave.go` gains a nullable tier.
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
- [x] Model + mapper + GUI round trip: `ToManualZoneSaves` / `FromManualZoneSaves`
      carry the tier; `EditorState.SetManualEdits` / `GetManualZones` speak
      `[]template_model.Zone`. Per §0b.15 the `Template(Model) →
      EditorState(Model)` direction lives in **`EditorStateMapper`**, so
      `template_model` never imports `editor_state_model`.
- [x] The write path must actually record: `ApplyNeutralZoneQuality` and
      `NewDefaultNeutralZone` currently flatten Quality into the profile and
      forget it. **This is where the `Unknown → Plastic` down-tier from §1 is
      actually fixed** — phase 2 proved the generator never emits an
      unclassifiable zone, so a manually re-tiered zone is the first one that
      needs the recorded value.
- [x] `templateHandler.UpdateTemplate`'s temporary wrap-by-name seam from phase
      2.5 must **not** grow to cover the manual re-tier. It needed to, so the
      DTO re-typing came first and `carryZoneTiersByName` is **deleted**.
- [x] Backward compatibility: a `.gen.json` with no `quality` loads as `nil` and
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

**The tier now has a home on disk and a single path to it.**
`editor_state.ManualZoneSave` gained `Quality *int8 json:"quality,omitempty"`,
and `editor_state_model` converts it to and from `*neutral_zone.Quality` through
`toQualityOrdinal` / `fromQualityOrdinal`. Both directions preserve nil, which is
the whole point: the enum is `iota - 1`, so a value field with `omitempty` would
drop every Plastic zone (ordinal 0) and it would load back as "never recorded".
The entity keeps the raw ordinal rather than the enum because
`internal/entities` may not import `internal/models` — §4.4.1 rule 3, enforced by
the layering gate.

**The zone-editor chain moved onto `template_model.Zone`** — the scope change
described above. Nine DTOs (`zoneEditorZonesDto`, `…GeometryRequestDto`,
`…ConnectionRequestDto`, `…QualityRequestDto`, `…RemoveRequestDto`,
`…MutationDto`, `templateUpdateDto`, `castleSettingsReapplyRequestDto`, plus
`ZoneEditorNeutralZoneRequestDto`'s return), `IZoneEditorHandler`,
`ITemplateHandler.ReapplyCastleSettings`, all four `connection_editor` services
and their interfaces, `ZoneEditorDialog`, `drivers.State.handleUpdateTemplate`
and `models.EditorState.SetManualEdits`/`GetManualZones`. Connections stayed
`entities.Connection` throughout — they carry no tier, and moving them would
have been churn for nothing.

**`carryZoneTiersByName` is deleted, not generalised.** `UpdateTemplate` now
rebuilds roads on the model zones, maps the *template* through the entity for a
clean copy, and then **re-attaches the applied model zone slice** to the result.
That is positional and exact, where the name lookup was neither — and it was
about to become actively wrong, since it would have restored a re-tiered zone's
*previous* planned tier.

**Three consumers stopped inferring** and now ask `ResolveQuality`, which reads
the recorded tier and only falls back to inference when there is none:

- `ManualReapplyService.SetNeutralZoneCastleCount` and `neutralCastleTarget` —
  **this is the `Unknown → Plastic` fix promised in §1.**
  `GetNeutralZoneProfile(QualityUnknown)` returns the Lowest profile, so an
  unclassifiable zone whose castles were rebuilt used to silently get Plastic
  city stats. A recorded tier is never `Unknown`.
- `ZoneTierService.GetGuardQuality` / `GetConnectionGuardQuality`, so a
  hand-picked tier decides the guard values of the connections touching it.
- `MandatoryContentProvider.CreateContentsForZones`, whose `Unknown` branch used
  to attach **no** mandatory-content rows at all.
- `zoneEditorHandler.GetZoneQuality`, which is what the dialog's Quality
  dropdown displays (owner's call, pulled forward from phase 4).

`GetQuality` still takes the **entity** and still exists: a template loaded from
a raw `.rmg.json` has no recorded tier and never will.

**No output moved.** Phase 2's 792-configuration sweep proved inference and the
plan agree on every zone the generator emits, and `ApplyNeutralZoneQuality`
stamps the pools its recorded tier implies, so `ResolveQuality` and `GetQuality`
return the same answer for every zone reachable today. The correction is latent
by design — it fires for a zone inference cannot classify, which is exactly what
this batch set out to make impossible to get wrong. The GPU suite passing
**without `-update`** is the evidence.

**`road_helpers` changed type rather than growing a twin.** `IsRoadTypeConnection`
and `IsRoadTypeCastle` now take `template_model.Road`; `ZoneEditorService` was
their only caller, so a second entity-typed pair would have been dead weight.

**Guards, and one of them was mutation-verified.** The new
`test/integration/manualZoneTierPersistence_integration_test.go` drives the real
save/load seam: a **Plastic** zone survives the round trip, a **Gold** one does
too, an unrecorded tier writes **no** `quality` key, and it loads back as nil.
Making `toQualityOrdinal` drop the zero ordinal — the exact shape of the
`omitempty` bug — fails the Plastic test and **only** that test. The file carries
`//go:build integration_test` because it uses the `SaveStateToFile` /
`LoadStateFromFile` test-only exports, which is the one and only reason §4.6.1
allows the tag. Unit-level round-trip tests cover the converter pair, and
`ApplyNeutralZoneQuality` / `NewDefaultNeutralZone` each gained a test asserting
the tier they record, alongside the pre-existing ones that assert the profile
they stamp (those now infer explicitly via `ToZoneEntity`, so they still test
what they always tested).

**The frozen fixtures were not touched**, and did not need to be: the new field
is `omitempty` and absent from both, which is exactly the legacy shape.
`git status` shows them unmodified, and no golden moved.

**Left for phase 4**, deliberately: `preview_service`, `file_service`, the
generator and the topology tree still speak entities behind
`TemplateMapper.ToEntity`; `previewLayoutRequestDto` still carries
`[]entities.Zone`; `PreviewLayoutService` still calls `GetQuality`, so the
preview colours a zone by inference rather than by its recorded tier. That is a
visual-only gap and the reason no pixel moved.

### Verification results (2026-09-03)

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 (no DI change — signatures moved, the graph did not) |
| Unit / untagged / integration | pass |
| **GPU suite, no `-update`** | **pass** — no pixel moved, no golden modified |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **74.3 %** (unchanged from phase 2.5, floor 72.5 %) |
| Frozen fixtures | unmodified in `git status` |
| Plastic round-trip guard | **verified by mutation** (dropping the zero ordinal fails it, and only it) |

---

## Phase 4: the rest of the sweep
Status: **Complete** (2026-09-03) — uncommitted, awaiting review.

⚠ **Phase 3 already took the zone-editor half of this phase** — read its summary
first. The 9 DTOs, `IZoneEditorHandler`, all four `connection_editor` services,
`ZoneEditorDialog` and the `drivers`/`app/gui/models` seam are done, and
`carryZoneTiersByName` is gone. What remains is everything phase 3 deliberately
left on a `TemplateMapper.ToEntity` bridge.

- [x] ~~`template_model.Zone` already carries the tier as of phase 2.5, so phase 4
      has no wrapper to build.~~ (This item replaced the former
      `internal/models/qualifiedZone.go` task — see §0b.11.)
- [x] The 9 DTOs carry `[]template_model.Zone` — done in phase 3, except
      `previewLayoutRequestDto`, which is still `[]entities.Zone`.
- [x] `ZoneEditorDialog.zones` / `originalZones`, `selectedZoneRef` /
      `zoneByName`, `zonePropertyRows` / `syncZoneProps` / `writebackZoneProps`
      and the Quality dropdown — done in phase 3.
- [x] Handlers and the `connection_editor` services take the model zone, and the
      temporary wrap-by-name seam is deleted — done in phase 3.
- [x] `*_testexports.go` accessors follow — `EditedZones()` and the
      `IZoneEditorDialog` contract moved in phase 3.
- [x] **`preview_service` is the one that still matters visually.**
      `PreviewLayoutService` calls `GetQuality` on the entity zone, so a zone the
      user re-tiered is coloured by inference rather than by its recorded tier.
      Move `BuildPreviewLayout` and `previewLayoutRequestDto` onto
      `template_model` and switch that call to `ResolveQuality`. ⚠ This is the
      first change in the batch that **can move pixels** — expect to regenerate
      goldens, and enumerate every zone whose colour changes.
- [x] `file_service`, the generator and the topology tree: decide per seam
      whether they move or keep their bridge. The generator builds the entity
      and lifts it (phase 2.5), which is what makes the golden test meaningful,
      so the topology tree arguably should **stay** on entities — say so
      explicitly rather than sweeping it by reflex.
- [x] ⚠ **Name-collision trap** for any new testexport: the layout checker
      matches test-only exports by identifier name tree-wide, so grep first —
      the testexports side always yields. (No new testexport was needed.)

### Verification Plan
- Full gate set, including `go vet -tags='integration_test,gui' ./...`.
- **The GPU suite will need `-update` for the preview change** — unlike every
  earlier phase. Check `git status` for ` M *.golden` afterwards and restore any
  that were not yours.
- Coverage ≥ 72.5 %.

### Phase Summary

**The preview now colours a zone by the tier that was recorded for it.**
`PreviewLayoutService.buildPreviewZones` calls `ResolveQuality(zone)` instead of
`GetQuality(zone)`, which is the last consumer in the batch to stop inferring.
That single line is the whole behavioural point of the phase; everything else in
the diff is the type change that made it reachable.

**`preview_service` speaks `template_model` end to end.** `BuildPreviewLayout`
and `CreatePreviewImage` (both the real generator and the null one) take
`*template_model.Template`, and all five layout strategies — ring/hub, scatter,
fixed geometry, balanced rings, manual positions — plus `layoutGeometry.go`'s
shared predicates now take `[]template_model.Zone` and
`[]template_model.Connection`. The package names no entity type at all any more.
The sweep was done with `gofmt -r` on an explicit file list and verified with
`git diff --numstat`: insertions equalled deletions on every file, so the AST
rewrite touched nothing but the type names. Imports were fixed separately.

**Three bridges came out, and one deliberately stayed:**

- `previewHandler` no longer holds a `TemplateMapper` at all. It used to flatten
  the request template with `ToEntity`; now it forwards the model, and the
  zones-only branch builds a one-variant `template_model.Template` literal
  instead of running the entity `VariantBuilder`. `NewPreviewHandler` lost its
  second argument and `wire_gen.go` was regenerated.
- `templateHandler.SaveTemplate` renders the preview **from the model** and only
  then flattens for the file service. Order matters for the reason this phase
  exists: flattening first would have thrown the tier away before the preview
  could read it.
- `ZoneEditorGeometryService.BuildGeometry` no longer round-trips its model zones
  through `ToZoneEntities` to synthesise a template for the layout service. It
  builds a model variant, lifting only the connections (which are still
  entities). That round trip was silently erasing `Quality`, so the zone editor's
  own canvas would have kept inferring even after this change.
- **`file_service` stays on the entity, permanently.** `SaveTemplateWithPreview`
  hands `*entities.RmgTemplate` to the template repository, which serializes it
  to `.rmg.json`. That is one of the two seams §0b.16 says keeps the wire format
  forever; converting there would be conversion for its own sake.

**The generator and the topology tree stay on entities — stated outright, as the
checklist asked.** `Generate` assembles the entity from the providers and lifts
it with `TemplateMapper.ToModel` before stamping tiers, and that ordering is
exactly what makes `TestWhenDefaultConfiguration_ReturnsGoldenTemplate` a proof
rather than an argument: the golden compares the bytes the providers produced,
not the bytes a model would round-trip to. Moving ~40 topology and builder files
onto the model would buy nothing — no topology has a tier to carry, since
`stampPlannedZoneTiers` derives every one of them from the plans afterwards — and
would cost the golden its meaning. This is the answer to the carry-forward's open
question, and it should not be revisited by reflex in a later batch.

**⚠ The plan predicted this phase would move pixels. It did not, and the reason
is worth recording rather than celebrating.** The GPU suite passes **without
`-update`**, and no `*.golden` is modified in `git status`. Two facts combine:

1. Every zone the *generator* emits has a recorded tier that inference agrees
   with — phase 2's 792-configuration sweep measured exactly that, zero
   disagreements.
2. Every zone the *editor* re-tiers goes through `ApplyNeutralZoneQuality`, which
   stamps the content pools its new tier implies. Inference reads those pools, so
   it lands on the same answer.

So the correction stays latent, as it has since phase 2: it fires only for a zone
whose tier was recorded but whose pools do not imply it — a state no path
reachable today produces. That is the batch working as intended, not the change
being a no-op.

**The goldens that cover the preview were genuinely exercised**, so this is not a
vacuous pass. `BaseHandler` masks the preview canvas interior by default because
the shipped default topology is Random, but `LayoutAndZonesTabHandler.SelectTopology`
**lifts that mask** as soon as a deterministic topology is picked. Every snapshot
taken after that compares the canvas pixel for pixel — including the zone-editor
suites, which re-tier a zone through the dropdown. Only the Random-topology
snapshots keep the canvas masked, and those could never have shown a tier change
either way.

**The real guard is therefore a unit test, and it is mutation-verified.**
`TestWhenZoneCarriesARecordedTier_ColoursItWithThatTierInsteadOfInferring` gives a
zone `Sides` layout and a `_t2_` guarded pool — which infers as `QualityLow` —
then records `QualityHigh` on it and asserts the preview zone comes back High.
Reverting `buildPreviewZones` to `GetQuality(ToZoneEntity(zone))` fails that test
and **only** that test; restoring it passes. Its sibling,
`TestWhenZoneCarriesNoRecordedTier_ColoursItWithTheInferredTier`, pins the
fallback so the nil branch cannot rot.

**Input for phase 5, measured rather than assumed:** `internal/dtos` still names
`internal/entities` in five files — `templateUpdateDto`,
`zoneEditorGeometryRequestDto`, `zoneEditorMutationDto`,
`zoneEditorRemoveRequestDto` and `zoneEditorZonesDto` — and in every case it is
`entities.Connection`, which phase 3 deliberately left alone because a connection
carries no tier. So `internal/dtos` cannot come off `entityNamerAllowList` yet,
and phase 5 should say that rather than move connections just to shorten a list.

**One dead path found, and removed.** `PreviewLayoutRequestDto.Zones` /
`.Connections` — the "editor-only preview when Template is nil" branch — had no
production caller: `previewPanel` always passes a template, and the zone editor
goes through `ZoneEditorGeometryService`. The only reader was the synthesis
branch inside `previewHandler` itself, and the only writers were unit tests. Both
fields, the branch and the three tests that existed solely to cover it are gone,
so `BuildPreviewLayout` is now a single forwarding line. The DTO is three fields:
`Template`, `Topology`, `CanvasSide`.

**And with the branch gone, so is the error return.** `BuildPreviewLayout` could
only ever return `nil` — the layout service has no failure mode — so the whole
chain now returns a bare `dtos.PreviewLayoutDto`: `IPreviewHandler`,
`previewHandler`, `GUIHandler`, the `TemplateHandlerMock` and the guiHandler
stub. That in turn emptied `models.PreviewLayoutCache.Get`, whose `build` callback
and "a failed build is not cached" retry existed only to carry that error; it now
takes `func() preview.Layout` and returns `preview.Layout`. `previewPanel` lost
its unreachable error branch, and the cache's two failure tests went with the
failure mode they described.

### Verification results (2026-09-03)

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 (regenerated after `NewPreviewHandler` lost its mapper; never hand-edited) |
| Unit / untagged / integration | pass (exit 0) |
| **GPU suite, no `-update`** | **pass** — no pixel moved, no golden modified (see summary for why) |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **74.3 %** (unchanged, floor 72.5 %) |
| Recorded-tier preview guard | **verified by mutation** — reverting to `GetQuality` fails it, and only it |
| `git status` | 27 modified files, **zero** `*.golden`, nothing staged |

⚠ **Not re-measured:** `BenchmarkEditorWindow_TabCycling`. The per-frame path is
untouched; the preview layout is rebuilt only when
`(templateRevision, topology, canvasSide)` changes. `BenchmarkPreviewLayoutService_BuildPreviewLayout`
did lose one whole `ToEntity` of the generated template per case in its setup,
which is a strict reduction, not a regression.

---

## Phase 5: shrink the allow-list, record the outcome
Status: **Complete** (2026-09-03)

- [x] Remove `internal/dtos` from `entityNamerAllowList` in
      [layering_test.go](../../test/unit/architecture/dependency/layering_test.go)
      **if** the 9 DTOs no longer name `entities.Zone`. Same check for
      `app/gui/dialogs` and `internal/handlers`. **Only ever remove entries** —
      if one will not come off, leave it and say why here.
      ⚠ `template_model` (§0b.11) makes far more of this list removable than the
      original plan assumed: any package that moved onto the model in phase 4
      stops naming an entity outright. Re-measure the whole list here rather
      than checking only the three named above.
- [x] Backlog: §2.2 becomes a ✅ DONE record with the behaviour deltas; §8 gets
      row **J**; refresh the coverage figure everywhere it is quoted (three
      places); update §2.6's file counts if the allow-list moved.
- [x] Update `.agent/session-carry-forward.md`.

### Verification Plan
- `go test ./test/unit/architecture/... -count=1` passes, and **fails if a
  removed entry is re-added as a violation** — prove the shrink is real.
- Every gate from every prior phase still green.

### Phase Summary

**Two entries came off: `app/gui/editor` and `internal/services/preview_service`.**
The whole list was re-measured rather than spot-checked, as the checklist
demanded — every `.go` file under `app/`, `internal/` and `cmd/` was scanned for
an `internal/entities` import and grouped by directory, then filtered against
the permitted-namer prefixes. That yields **84 files in 21 packages**, down from
the 113 files in 23 packages the gate was seeded with in batch I. The 21 groups
matched the 21 surviving allow-list entries exactly, which is the cross-check
that nothing was missed in either direction.

`preview_service` is the interesting one and it is phase 4's doing: the package
now names no entity type at all, because `BuildPreviewLayout`, all five layout
strategies and `layoutGeometry.go`'s shared predicates take `template_model`.
`app/gui/editor` is the cheap one — `window.go` had already stopped naming an
entity in an earlier phase and nobody had re-measured.

**The shrink was proven by mutation, twice.** Adding
`"…/internal/entities"` plus a `var _ entities.Zone` to
`previewLayoutService.go` made
`TestWhenEntityConsumersAreScanned_OnlyPermittedPackagesNameAnEntity` fail with
exactly `internal/services/preview_service/previewLayoutService.go`; the same
edit in `app/gui/editor/window.go` produced exactly that file. Both were reverted
and `git status` confirms neither production file is modified. So the two
removals are real constraints now, not bookkeeping.

**⚠ Nothing else came off, and three refusals are worth naming rather than
retrying:**

- **`internal/dtos` stays.** All five files name `entities.Connection` —
  `templateUpdateDto`, `zoneEditorGeometryRequestDto`, `zoneEditorMutationDto`,
  `zoneEditorRemoveRequestDto`, `zoneEditorZonesDto`. Phase 3 left connections
  on the entity deliberately: a connection carries no tier, so moving it buys
  nothing but a shorter list. Whether the `.rmg.json` vocabulary deserves a
  documented carve-out the way `internal/helpers/data` has one is backlog §2.6
  step 2, a decision — not a sweep.
- **`app/gui/dialogs`, `app/gui/drivers`, `app/gui/models` stay**, for the same
  reason: their zones are `template_model.Zone` now, their connections and
  templates are not.
- **`internal/services/file_service` and the whole `template_generator` tree
  stay permanently**, per §0b.16 and phase 4's summary. Those two seams own the
  wire format on purpose. An agent reading only the shrinking list would
  eventually try to close them; the backlog §2.6 step 4 text now says outright
  that they are exempt by decision, not by debt.

**The backlog is the surviving record.** §2.2 is now a ✅ DONE entry carrying the
design (pointer-typed tier and why, the two-query tier service and why both
survive, the persisted `*int8`), the enumerated behaviour deltas (there are none,
and the 792-configuration measurement that establishes it), the two permanent
entity seams, and the dead code phase 4 removed. §8 gained row **J**. The
coverage figure was refreshed in all three places it is quoted (the header
baseline, the §8 coverage note, the §9 gate table). §2.6's title, area table and
steps 3–4 were rewritten to the new counts.

### Verification results (2026-09-03)

| Gate | Result |
| --- | --- |
| `go build ./...` | exit 0 |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 |
| Unit / untagged / integration | pass (exit 0) |
| **GPU suite, no `-update`** | **pass** — no golden modified |
| `golangci-lint-v2 run ./...` | **0 issues** |
| Unit coverage | **74.3 %** (unchanged, floor 72.5 %) |
| Allow-list shrink | **verified by mutation** — re-importing an entity in either removed package fails the rule with exactly that file |
| `git status` | one modified file (`layering_test.go`), zero `*.golden`, nothing staged |

---

## Final Recap

**A neutral zone's tier is now recorded where the zone is, and asked for in one
place.** Before this batch the tier was decided once on `neutral_zone.Plan`,
flattened into layout/pools/castles, thrown away, and then reverse-engineered by
`ZoneClassifier` at eight consumers — so every feature that edited a zone's
content could silently flip its tier, and the inference rules had to stay in
lockstep with the profile catalogue with nothing enforcing it.

What landed, in five phases:

1. **`IZoneTierService` absorbed the classifier and `ZoneClassifier` was
   deleted** (not wrapped — a wrapper would have parked the real logic in
   another type forever). Six injection sites moved onto it, including
   `PreviewLayoutService`, which had been bypassing DI by hard-building its own
   classifier and would otherwise have kept inferring for the rest of the batch.
2. **The generator records what it planned.** `planZoneTiers` builds the
   label→quality map from the `neutral_zone.Plans` the generator already holds;
   no topology, factory or builder changed.
3. **2.5 — `internal/models/template_model/` mirrors the whole `.rmg.json`
   schema**, and the tier rides *on the zone* as
   `Quality *neutral_zone.Quality` rather than in a side-car map. This replaced
   both the original `models.QualifiedZone` design and the phase 2 tier index;
   `models.GeneratedTemplate` was deleted with them.
4. **`.gen.json` persists the tier** as
   `ManualZoneSave.Quality *int8 json:"quality,omitempty"`, nil-preserving in
   both directions, with a mutation-verified guard on the Plastic (ordinal 0)
   case that `omitempty` on a value field would have silently dropped.
5. **The sweep** moved the zone editor, the handlers, the `connection_editor`
   services and `preview_service` onto the model; `file_service` and the
   generator/topology tree keep the entity permanently and on the record.
6. **The allow-list shrank by two packages, measured and mutation-proven.**

**The behavioural result is deliberately latent, and that was measured, not
assumed.** A 792-configuration sweep found zero arena moves and zero
inference-versus-plan disagreements; every gate, including the GPU suite, passed
without `-update` at every phase. The generator never emits a zone its own pools
cannot classify, and `ApplyNeutralZoneQuality` stamps the pools a re-tier
implies, so recorded and inferred answers coincide everywhere reachable today.
The correction fires for the state this batch set out to make unrepresentable —
a zone whose tier is known but whose content does not imply it — and the
`Unknown → Plastic` silent down-tier through `GetNeutralZoneProfile` is gone
because a recorded tier is never `Unknown`.

**Nothing under `internal/entities/template/**` changed.** Branch A of backlog
§2.2 was never needed and stays closed.

Coverage 72.9 % → **74.3 %**; lint held at **0** throughout; no golden moved.

## Deployment Plan

This batch ships as source only — no data file, no schema file, no build flag and
no configuration changes with it.

1. **Review the working tree.** `git status --short` should show the phase 5
   change plus whatever of phase 4 is still uncommitted; **nothing staged**, and
   **zero `*.golden`**. Per AGENTS §2.5 the author stages and commits, not the
   agent.
2. **Run the gate set from a clean tree**, in this order:

   ```powershell
   go build ./...
   go vet ./...
   go vet -tags='integration_test,gui' ./...
   gofmt -l ./app ./internal ./test ./cmd
   go run ./cmd/testlayoutcheck .
   wire diff ./internal/composition/...
   go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...
   go tool cover '-func=coverage.txt'
   go test ./test/... -count=1
   go test -tags=integration_test ./test/integration/... -count=1
   go test -tags='integration_test,gui' ./test/integration/gui/... -count=1
   golangci-lint-v2 run ./... --issues-exit-code=0
   ```

   Expected: everything clean, coverage **≥ 74.3 %**, lint **0 issues**, and the
   GPU suite green **without** `-update`. A golden that moves here is a real
   rendering change and must be explained, not regenerated.
3. **Commit and push.** CI re-runs the same gates on Ubuntu; the GUI job uses
   Xvfb + Mesa llvmpipe against goldens generated locally on a real GPU, and the
   two-gate comparer already tolerates that difference.
4. **Migration: none, in either direction.** `quality` is a new `omitempty` key
   in `.gen.json`. An old file simply has no key, which loads as `nil` = "infer
   it" — exactly the pre-batch behaviour. A new file read by an old build hits an
   unknown field, which the loader already ignores. The two frozen fixtures under
   `test/test_helpers/testdata/` are unmodified for that reason.
5. **`.rmg.json` is byte-identical.** Not one byte of
   `internal/entities/template/**` changed, and `TestWhenDefaultConfiguration_ReturnsGoldenTemplate`
   still compares the bytes the providers produced. The game reads the same file
   it always did, from the same auto-detected templates directory.
6. **Smoke test in the app**, since the visible surfaces did move types even
   though they did not move pixels: generate a template, open the zone editor,
   re-tier a neutral zone through the Quality dropdown, Apply, save the
   `.gen.json`, reload it, and confirm the dropdown still shows the chosen tier.
   That is the one path the whole batch exists to make reliable.
7. **Rollback** is a plain revert of the batch's commits. No data written by the
   new build is unreadable by the old one.

## After this batch

- Backlog §2.6 step 2 is the next decision: whether `internal/dtos` /
  `internal/handlers` naming `entities.Connection` and `entities.RmgTemplate` is
  a breach at all, or whether the `.rmg.json` vocabulary earns a documented
  carve-out. Do not start it as a sweep.
- Two benchmark baselines disagree in the record (phase 2's ~5,699 allocs/op
  versus backlog §1.4's 6,640 after the clone batch). Neither is a batch J
  regression; someone should reconcile them.
