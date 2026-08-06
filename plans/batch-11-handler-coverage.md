# Batch 11 — Service interfaces + handler/catalogue coverage

Close review findings §6.2 (`internal/handlers` has no mirrored unit tests) and
§6.4 (two `app/gui/constants` catalogues at 0%). Because the five handlers hold
**concrete** struct dependencies they cannot be mocked, the owner chose to first
convert every constructor-injected service under `internal/` to an interface
(full closure), then write the tests against `testify` mocks.

Source of truth for the findings themselves stays
[todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) — mark §6.2 / §6.4
`✅ FIXED` in place there when this plan completes.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Hard rules that apply throughout (AGENTS.md): never modify `data/`,
`internal/entities/template/` or `internal/registry/`; stay cross-platform;
**never stage and never commit** — the owner reviews and commits.

## Owner decisions already taken (do not re-litigate)

| Question | Decision |
| --- | --- |
| Do §6.2 + §6.4 at all | Yes, both |
| Handler scope | All five (`state`, `preview`, `template`, `contentRule`, `zoneEditor`) |
| Mock vs real dependencies | Refactor to interfaces first, then mock |
| Refactor boundary | **Full closure** of constructor-injected services under `internal/` |
| Interface surface | Full public surface of each struct |
| Constructor return type | Return the interface, not `*Struct` |
| `normalizeInactiveNeutralCounts` | Test current behaviour only; no behaviour change |
| Preview-generator nil guard | Null-object `NullPreviewGeneratorService` in `preview_service` |
| Catalogue assertions | Invariants + a few named spot checks (no ~31 hard-coded SIDs) |
| Other 0% code | Scan, and add tests for anything pure/easy |
| Interface-placement rules | New AGENTS.md §4.2.2 |

## Interface placement (derived from the new AGENTS.md §4.2.2 rules)

Rule recap: **<5** implementation files in a package needing interfaces → same
package; **≥5** → `{singular package}_interfaces` subpackage; cross-package
implementations or cycle-breaking → `internal/interfaces/<sub>`.

| Package | Impl files needing an interface | Placement | Interfaces |
| --- | ---: | --- | --- |
| `internal/services/file_service` | 1 | same package | `IFileService` |
| `internal/validators` | 1 | same package | `IEditorStateValidator` |
| `internal/mappers` | 2 | same package | `IGeneratorConfigMapper`, `IMandatoryContentItemMapper` |
| `internal/services/preview_service` | 2 | same package | `IPreviewLayoutService`, `IPreviewGeneratorService` |
| `internal/services/connection_editor` | 3 | same package | `IConnectionEditorService`, `IZoneEditorService`, `IManualReapplyService` |
| `internal/services/template_generator` | 1 | same package | `ITemplateGenerator` |
| `.../template_generator/generation_tuning` | 1 | same package | `IGenerationTuningFactory` |
| `.../providers/topology/base` | 1 | same package | `ITopologyConnectionService` |
| `.../providers/topology/tournament_variant` | 4 | same package | `IClusterService` — **already exists, no change** |
| `internal/services/zones` | **5** | `zones/zone_interfaces` | `ICastleFactory`, `IRoadFactory`, `IZoneFactory`, `IZoneClassifier`, `IZoneLabelProvider` (**relocated**) |
| `.../providers` | **7** | `providers/provider_interfaces` | `IContentLimitProvider`, `IGameRulesProvider`, `IGladiatorArenaProvider`, `IMandatoryContentProvider`, `ITopologyProvider`, `ITopologyServiceLookup`, `IZoneLayoutProvider` |
| `.../providers/topology` | **13** | `topology/topology_interfaces` | `ITopologyService` (**one** interface for all 12 services), `IPositionedTopologyBuilder` |

**21 new interface files, 1 relocation.** Not 34 types-worth: all twelve topology
services expose the single method
`CreateTopologyVariant(config.GeneratorConfig, []string, neutral_zone.Plans, models.GenerationTuning, string) entities.Variant`,
so they share one interface exactly like the four cluster services share
`IClusterService`.

### Out of scope (not constructor-injected)

`asset_provider.AssetProvider` (built inside `NewPreviewGenerator`),
`content_rules.DistanceCatalog` / `VariantMappingCatalog` (built inside
`NewContentRuleService`), `position_layout.PositionLayoutService`, all
`builders/*` fluent builders, and every pure data/DTO constructor.

### Known hazards

- `TopologyProvider.ShufflePlayerZones(bool) *TopologyProvider` is fluent — its
  return type becomes `ITopologyProvider`.
- Relocating `IZoneLabelProvider` out of `zones` touches `providerSets.go`
  (`wire.Bind`), `base.NewTopologyConnectionService`, `NewTemplateGenerator` and
  all 17 topology/cluster constructors.
- `wire.Bind` for `IZoneLabelProvider` and `IContentRuleService` disappears once
  those constructors return interfaces.
- Never pass `-tags=wireinject` to build/test (AGENTS.md §4.6.3).
- `git mv` stages; follow with `git restore --staged <path>`.

---

## Phase 0: Baseline and rules
Status: Complete

- [x] Add AGENTS.md §4.2.2 "Interface placement" with the three rules, the
      owner's examples, and the "factory returns the broadest interface" rule.
- [x] Regenerate the coverage baseline (`coverage.txt` at repo root is stale —
      it predates `app/gui/constants/spellLabel.go`).
- [x] Record baseline total coverage and `golangci-lint-v2` issue count here.

### Verification Plan
- `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...` then `go tool cover '-func=coverage.txt'` — succeeds, no missing-file error.
- `golangci-lint-v2 run ./... --issues-exit-code=0` — record the count.

### Phase Summary
AGENTS.md §4.2.2 added between §4.2.1 and §4.3 (three placement rules, three
worked examples, factory-returns-broadest-interface rule).

Baseline captured after regenerating the stale profile: **total coverage 65.5%**,
`go build ./...` green, unit suite green. Lint had regressed from Batch 10's
reported 0 to **2 `wastedassign` findings** in
[app/gui/panels/bonusesPanel.go](../app/gui/panels/bonusesPanel.go) (`label := ""`
immediately overwritten in both `BonusSpell` and `BonusStartingItem` branches —
fallout from Batch 10's §3.3 spell-label extraction). Both changed to
`var label string`; lint back to 0.

## Phase 1: Leaf-package interfaces
Status: Complete

Packages with <5 implementation files — interface lives beside the
implementation, constructor returns the interface.

- [x] `file_service.IFileService` (3 methods)
- [x] `validators.IEditorStateValidator` (1)
- [x] `mappers.IGeneratorConfigMapper` (1), `mappers.IMandatoryContentItemMapper` (1)
- [x] `generation_tuning.IGenerationTuningFactory` (1)
- [x] `base.ITopologyConnectionService` (4)
- [x] Update every consumer field/parameter to the interface type.

### Verification Plan
- `go build ./...` and `go vet -tags=integration_test ./...` pass.
- `go test -count=1 ./test/unit/...` passes.

### Phase Summary
Six `*Interface.go` files created beside their implementations; the six
constructors now return the interface instead of `*Struct`.

Consumers switched to interface fields/params: `stateHandler` (fileService,
editorValidator), `templateHandler` (mapper, fileService), `zoneEditorHandler`
(mapper, tuningFactory), `connection_editor.ManualReapplyService`
(tuningFactory), `template_generator.TemplateGenerator` (tuningFactory),
`mappers.NewConfigMapper` (contentItemMapper), `base.TopologyBase` (field +
`NewTopologyBase` param), `composition.provideTopologyServices`, all 12 topology
service constructors, all 4 cluster service constructors,
`NewPositionedTopologyBuilder`, and four `test/test_helpers` builders plus
`test/unit/.../fileService/common_test.go`.

Key detail: `ITopologyConnectionService` renames the `zones []entities.Zone`
parameter to `allZones` — the `base` package imports a package named `zones`, so
the original name shadows it in an interface declaration.

`wire gen ./internal/composition/...` regenerated `wire_gen.go` cleanly; no
`wire.Bind` changes were needed because the constructors now yield the interface
type directly.

Verification: `go build ./...` clean, `go vet -tags=integration_test ./...`
silent, `go test -count=1 ./test/unit/...` exit 0 (architecture construction and
dependency suites included).

### Phase Summary
_(write when phase completes)_

## Phase 2: `zones/zone_interfaces`
Status: Complete

- [x] Create `internal/services/zones/zone_interfaces/` with `ICastleFactory` (6),
      `IRoadFactory` (2), `IZoneFactory` (3), `IZoneClassifier` (3).
- [x] Relocate `IZoneLabelProvider` into it and delete
      `internal/services/zones/zoneLabelProviderInterface.go`.
- [x] All five `zones` constructors return their interface; drop the
      `IZoneLabelProvider` `wire.Bind`.
- [x] Update all consumers (`connection_editor`, `providers`, `topology`,
      `tournament_variant`, `preview_service`, `template_generator`).

### Verification Plan
- `go build ./...` passes; `grep` finds no remaining `*zone_services.ZoneFactory` etc. outside `zones`.
- `go test -count=1 ./test/unit/...` passes.

### Phase Summary
`internal/services/zones/zone_interfaces/` created with five interface files;
`internal/services/zones/zoneLabelProviderInterface.go` deleted (relocated
verbatim). All five constructors now return their interface, and the
`IZoneLabelProvider` `wire.Bind` was removed from `ZoneSet`.

**Owner decision taken mid-phase (option A):** `ZoneFactory` called two
*private* `CastleFactory` methods across files, which no interface can express.
`createPlayerSpawnCastle` and `createAbandonedOutposts` were promoted to
`CreatePlayerSpawnCastle` / `CreateAbandonedOutposts` and added to
`ICastleFactory`, so `ICastleFactory` has **six** methods, not four. Behaviour is
unchanged. A sweep of `internal/` confirmed this was the only such case.

**Second owner decision:** `ZoneEditorService.CastleFactory` was exported only so
the same-package `ManualReapplyService` could reach it. Renamed to
`castleFactory` (same package, so no access is lost); four call sites in
[manualReapplyService.go](../internal/services/connection_editor/manualReapplyService.go)
and one struct literal updated.

Consumers migrated: 13 topology services, 4 cluster services,
`base.TopologyBase`, `base.TopologyConnectionService`, `TemplateGenerator`,
`composition.provideTopologyServices`, `zoneEditorHandler`,
`ConnectionEditorService`, `ManualReapplyService`, `ZoneEditorService`,
`PreviewLayoutService`, `GladiatorArenaProvider`, `MandatoryContentProvider`,
three `test/test_helpers` builders and one unit-test helper.

`PreviewLayoutService` keeps an import of the concrete `zones` package: its
`zoneClassifier` is **self-constructed** inside `NewPreviewLayoutService()`, not
constructor-injected, so it is outside this refactor's boundary. Its field type
is now `zone_interfaces.IZoneClassifier`, so it can still be swapped later.

Verification: `go build ./...` clean; `go vet -tags=integration_test ./...` and
`go vet -tags='integration_test,gui' ./test/...` silent;
`go test -count=1 ./test/unit/...` exit 0; `go run ./cmd/testlayoutcheck .`
prints `test-layout check passed`; `golangci-lint-v2 run ./... --fix` reports
**0 issues**. (A first unit-suite run reported a spurious `[build failed]` for
`test/unit/app/gui/drivers/state`; it passed in isolation and on re-run — stale
build cache, not a code problem.)

## Phase 3: `topology/topology_interfaces`
Status: Complete

- [x] `ITopologyService` — one interface, `CreateTopologyVariant`; assert all 12
      services satisfy it.
- [x] ~~`IPositionedTopologyBuilder`~~ — **dropped, see summary.**
- [x] ~~All 12 topology constructors return the interface~~ — **owner chose
      option B: constructors stay concrete, see summary.**

### Verification Plan
- `go build ./...` passes.
- `go test -count=1 ./test/unit/...` and `go test -tags=integration_test -count=1 ./test/integration/...` pass.

### Phase Summary
**Blocker found and resolved by owner decision (option B).** The plan called for
one shared `ITopologyService` returned by all 12 constructors. That is
impossible: `NewTopologyServiceLookup` takes 12 distinct concrete parameters and
**wire keys providers by output type**, so 12 providers returning the same
interface would collide with a multiple-bindings error.

The owner chose to **declare the interface but keep the factories concrete**:

- [topology_interfaces/topologyServiceInterface.go](../internal/services/template_generator/providers/topology/topology_interfaces/topologyServiceInterface.go)
  declares `ITopologyService` and documents *why* the implementations stay
  concrete.
- [topologyServiceAssertions.go](../internal/services/template_generator/providers/topology/topologyServiceAssertions.go)
  holds all 12 `var _ ITopologyService = (*XTopologyService)(nil)` assertions in
  one place, so the shared contract is enforced at compile time without
  touching 12 files. All 12 compile clean.
- `NewTopologyServiceLookup` and the 12 constructors are **unchanged**.

**`IPositionedTopologyBuilder` was dropped.** `PositionedTopologyBuilder` is not
constructor-injected — it is **embedded by value** (`PositionedTopologyBuilder`,
initialised with `*NewPositionedTopologyBuilder(...)`) in `circlesTopology`,
`crossTopology`, `fractalTopology`, `geometricTopology`, `randomTopology` and
`squareTopology`. A value-embedded struct cannot be replaced by an interface
without restructuring all six services, and it is outside this batch's
"constructor-injected services" boundary.

Net effect on §6.2: none. Handlers never touch topology services directly —
they reach them only through `ITemplateGenerator` (Phase 4) and
`ITopologyServiceLookup` (Phase 4), both of which remain mockable.

Verification: `go build ./...` clean.

## Phase 4: `providers/provider_interfaces`
Status: Complete

- [x] Seven interfaces: `IContentLimitProvider`, `IGameRulesProvider`,
      `IGladiatorArenaProvider`, `IMandatoryContentProvider`, `ITopologyProvider`,
      `ITopologyServiceLookup`, `IZoneLayoutProvider`.
- [x] `ShufflePlayerZones` returns `ITopologyProvider`.
- [x] `template_generator.ITemplateGenerator` + `NewTemplateGenerator` takes the
      provider interfaces.

### Verification Plan
- `go build ./...` passes; `go test -count=1 ./test/unit/...` passes.

### Phase Summary
Created `internal/services/template_generator/providers/provider_interfaces/`
(package `provider_interfaces`) with the seven interfaces above plus
`topologyVariantCreator.go`. The `TopologyVariantCreator` func type was
**relocated** from `providers` into `provider_interfaces` and the original file
deleted, because `ITopologyServiceLookup` returns that type and keeping it in
`providers` would create the cycle `provider_interfaces → providers →
provider_interfaces`. `providers.TopologyServiceLookup` now uses
`provider_interfaces.TopologyVariantCreator` for both fields, the `byTopology`
map and both method return types.

`internal/services/template_generator/templateGeneratorInterface.go` declares
`ITemplateGenerator` (`SetConfiguration`, `Generate`). `TemplateGenerator`'s six
provider fields and the matching `NewTemplateGenerator` parameters are now the
`provider_interfaces.I*` types, and the constructor returns `ITemplateGenerator`
— which removed the `providers` import from `templateGenerator.go` entirely.

All seven provider constructors now return their interface. Consumers switched:
`internal/composition/topologyServiceProvider.go` (`provideTopologyServices` →
`ITopologyServiceLookup`), `internal/handlers/templateHandler.go`
(`templateGenerator` → `ITemplateGenerator`, `contentProvider` →
`IMandatoryContentProvider`; the `providers` import became
`provider_interfaces`), the three `test/test_helpers` builders
(`templateGenerator.go`, `topologyProvider.go`, `topologyServiceLookup.go`) and
three unit-test `common_test.go` helpers (`gladiatorArenaProvider`,
`mandatoryContentProvider`, `guiHandler`).

Wire regenerated with `wire gen ./internal/composition/...` (writes its success
line to STDERR, so PowerShell reports a non-zero exit — that is expected).
Verified green: `go build ./...`, `go vet -tags=integration_test ./...`,
`go vet -tags='integration_test,gui' ./test/...`,
`go test -count=1 ./test/unit/...`, `go run ./cmd/testlayoutcheck .`
(`test-layout check passed`) and `golangci-lint-v2 run ./...` (0 issues after
one `--fix` pass for gci/gofmt/golines on the new files).

Gotcha for future agents: a stray duplicate `package providerinterfaces` line
appeared above the real package clause in the newly created
`provider_interfaces/topologyVariantCreator.go`; if `go build` reports "found
packages X and Y in <dir>", check the top of freshly created files first.

## Phase 5: Handlers + null-object preview generator
Status: Complete

- [x] `preview_service.IPreviewLayoutService`, `IPreviewGeneratorService`.
- [x] New `internal/services/preview_service/nullPreviewGeneratorService.go` —
      `CreatePreviewImage` returns `nil` (`FileService.SaveTemplateWithPreview`
      already early-returns on a nil image).
- [x] `composition.providePreviewGenerator` returns the null object on
      asset-load failure instead of a nil pointer.
- [x] Delete the `if this.previewGenerator != nil` guard in
      `templateHandler.SaveTemplate`.
- [x] `connection_editor` interfaces; all five handler structs + constructors
      switch to interface-typed fields.

### Verification Plan
- `go build ./...` passes.
- A unit test proves `SaveTemplate` with the null generator still writes the template and returns its path.

### Phase Summary
**Owner decision (blocker resolved).** `ManualReapplyService` reached into
`ZoneEditorService` internals — `this.zoneEditor.castleFactory` (4 sites) and
`this.zoneEditor.rebuildCastleRoads` (3 sites) — which made an
`IZoneEditorService` field impossible. The owner chose a **B+C combination**:
`ManualReapplyService` now takes its own `zone_interfaces.ICastleFactory`
(no duplicated logic, no private-field access), and
`ZoneEditorService.rebuildCastleRoads` was promoted to the public
`RebuildCastleRoads` and added to `IZoneEditorService` — so no `CastleFactory()`
accessor was needed and `NewZoneEditorService` returns the interface without any
`wire.Bind`. `NewManualReapplyService`'s parameter list is now
`(zoneEditor, castleFactory, zoneClassifier, tuningFactory)`.

New interfaces, all declared in-package (3 and 2 implementation files
respectively, so AGENTS.md §4.2.2 rule 1 applies):
- `internal/services/connection_editor/connectionEditorServiceInterface.go` →
  `IConnectionEditorService` (4 methods).
- `internal/services/connection_editor/zoneEditorServiceInterface.go` →
  `IZoneEditorService` (10 methods, including the promoted `RebuildCastleRoads`).
- `internal/services/connection_editor/manualReapplyServiceInterface.go` →
  `IManualReapplyService` (2 methods).
- `internal/services/preview_service/previewLayoutServiceInterface.go` →
  `IPreviewLayoutService` (`BuildPreviewLayout`).
- `internal/services/preview_service/previewGeneratorServiceInterface.go` →
  `IPreviewGeneratorService` (`CreatePreviewImage`).

New implementation: `internal/services/preview_service/nullPreviewGeneratorService.go`
(`NullPreviewGeneratorService` + `NewNullPreviewGenerator() IPreviewGeneratorService`).
`providePreviewGenerator` now takes `IPreviewLayoutService`, returns
`IPreviewGeneratorService`, and degrades to the null object on asset-load
failure — which removed the typed-nil hazard entirely and let the
`previewGenerator != nil` guard (and the `image` import) disappear from
`templateHandler.SaveTemplate`.

All five constructors now return interfaces:
`NewConnectionEditorService`, `NewZoneEditorService`, `NewManualReapplyService`,
`NewPreviewLayoutService`, `NewPreviewGenerator` (still `(IPreviewGeneratorService, error)`).

Consumers switched: `templateHandler` (`connectionEditor`, `zoneEditor`,
`manualReapply`, `previewGenerator`), `zoneEditorHandler` (`connectionEditor`,
`zoneEditor`), `previewHandler` (`previewLayout`),
`providers.MandatoryContentProvider.zoneEditor`,
`test/test_helpers/zoneEditorService.go`, and the `common_test.go` helpers in
`manualReapplyService`, `guiHandler` and `previewGeneratorService`.

Tests added (§2.3): `rebuildCastleRoads_test.go` (3 tests — rebuilds castle
roads for the current main objects, drops stale castle roads, preserves
connection roads) under the existing `zoneEditorService` folder, and a new
`test/unit/internal/services/preview_service/nullPreviewGeneratorService/`
folder with `createPreviewImage_test.go` (2 tests) and
`newNullPreviewGenerator_test.go` (1 test). The "SaveTemplate with the null
generator" assertion from the verification plan is deferred to the handler
tests in Phase 8, where a mocked `IPreviewGeneratorService` returning `nil`
covers the same path directly.

Verified green: `go build ./...`, `go vet -tags=integration_test ./...`,
`go vet -tags='integration_test,gui' ./test/...`,
`go test -count=1 ./test/unit/...`,
`go test -tags=integration_test -count=1 ./test/integration/...`,
`go run ./cmd/testlayoutcheck .` and `golangci-lint-v2 run ./...` (0 issues).
`wire gen ./internal/composition/...` reports `unchanged` — the injector graph
already matched after the earlier regeneration in this phase.

## Phase 6: Wire regeneration
Status: Not started

- [ ] Update `providerSets.go`: drop the two now-redundant `wire.Bind` calls.
- [ ] `wire gen ./internal/composition/...`, then `wire diff ./internal/composition/...` (exit 0).
- [ ] Confirm `wire_gen.go` is committed and not hand-edited.

### Verification Plan
- `wire diff ./internal/composition/...` exits 0.
- `go run .` starts (or at minimum `go build ./...` + the GUI integration suite passes).
- `go test -tags='integration_test,gui' -count=1 ./test/integration/gui/...` — zero snapshot diff.

### Phase Summary
_(write when phase completes)_

## Phase 7: Mocks
Status: Not started

- [ ] One `testify/mock` per interface in `test/test_helpers/`, file named
      `<interfaceName>Mock.go` in lower-camel, style copied from
      `contentRuleServiceMock.go` (receiver `this`, `arguments.Get(0).(T)` with
      the comma-ok form).
- [ ] Only the mocks the handler tests actually need — do not pre-build all 21.

### Verification Plan
- `go vet ./test/...` passes; `go run ./cmd/testlayoutcheck .` prints `test-layout check passed`.

### Phase Summary
_(write when phase completes)_

## Phase 8: §6.2 handler unit tests
Status: Not started

One folder per implementation file, one `<publicMethod>_test.go` per method,
package `<fileName>_test`, `t.Parallel()` everywhere, AAA sections, one logical
assertion per test, `testify` + `gofakeit` only.

- [ ] `test/unit/internal/handlers/stateHandler/` — `loadState_test.go`,
      `saveState_test.go`, `validateEditorState_test.go`. Cover: empty/whitespace
      path → `ErrNoOutputPath`; loader error propagated; nil state →
      `ErrNothingToSave`; `fixIssues=false` returns warnings with an unmodified
      state; `AdvancedMode` true/false each zeroing the correct field set.
- [ ] `test/unit/internal/handlers/previewHandler/buildPreviewLayout_test.go` —
      nil template + nil zones + nil connections; synthesised one-variant template.
- [ ] `test/unit/internal/handlers/templateHandler/` — `generateTemplate_test.go`
      (`ErrNoTemplateName`, `ErrGeneratedTemplateInvalid`, warning concatenation),
      `updateTemplate_test.go` (`ErrProvidedTemplateInvalid`, `ErrZonesMissing`,
      content rebuild only when `EditorState != nil`),
      `reapplyCastleSettings_test.go`, `saveTemplate_test.go`.
- [ ] `test/unit/internal/handlers/contentRuleHandler/` —
      `getContentRuleEditorOptions_test.go` (Variant option appended only when
      variants exist), `describeContentRule_test.go` (nil rule → invalid baseline).
- [ ] `test/unit/internal/handlers/zoneEditorHandler/` — one file per public
      method (14).

### Verification Plan
- `go test -count=1 ./test/unit/internal/handlers/...` passes.
- `go run ./cmd/testlayoutcheck .` passes.
- Coverage for `internal/handlers/*.go` rises; total does not drop.

### Phase Summary
_(write when phase completes)_

## Phase 9: §6.4 catalogue unit tests
Status: Not started

- [ ] `test/unit/app/gui/constants/bannableItems/` — one file per public function:
      `getBannableItemsWithExclusions_test.go`, `findBannableItem_test.go`,
      `getBannedItemLabel_test.go`, `sidToDisplayName_test.go`,
      `compareBannableItems_test.go`.
      Invariants: no empty `Sid`/`Name`/`Category`; SIDs unique; sorted by
      category then name; excluded SIDs absent; caller's slice not mutated;
      category set is exactly the six known values. Plus 2–3 named-SID spot checks.
- [ ] `test/unit/app/gui/constants/valueOverrideSids/getValueOverrideSidsWithExclusions_test.go`
      — sorted, no duplicates, no empty SIDs, exclusions removed, input not mutated.

### Verification Plan
- `go test -count=1 ./test/unit/app/gui/constants/...` passes.
- Both files move off 0.0% in `go tool cover -func`.

### Phase Summary
_(write when phase completes)_

## Phase 10: Remaining 0% scan
Status: Not started

- [ ] From the fresh profile, list every non-Gio production file still at 0% or
      far below 80%.
- [ ] Add tests for the pure/easy ones; record genuinely untestable or
      GUI-bound ones in [todo/test_observations.md](../todo/test_observations.md).

### Verification Plan
- `go test -count=1 ./test/unit/...` passes; coverage strictly above the Phase 0 baseline.

### Phase Summary
_(write when phase completes)_

## Phase 11: Close out
Status: Not started

- [ ] Full suite: build, vet, testlayoutcheck, unit, integration, GUI integration,
      coverage, lint.
- [ ] Mark §6.2 and §6.4 `✅ FIXED` in place in `todo/review-opus5-08-04.md`, and
      §12 item 11 done.
- [ ] Update repository memory with the new interface conventions.
- [ ] Rewrite `.agent/session-carry-forward.md`.

### Verification Plan
- Every command in AGENTS.md §7 Quick Reference passes; coverage ≥ Phase 0 baseline; lint issue count ≤ baseline.

### Phase Summary
_(write when phase completes)_

## Final Recap
_(write when all phases complete)_

## Deployment Plan
_(write when all phases complete)_
