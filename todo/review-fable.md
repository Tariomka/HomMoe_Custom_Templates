# Codebase Review — HomMoe Custom Templates

Reviewer role: Senior/Principal engineer, full-codebase review.
Scope: `main.go`, `app/gui/**`, `internal/**`, `test/**`, CI/tooling, repo hygiene.
Goal: after applying everything below the app behaves the same or better, is easier to
navigate for a new developer, and has far less duplication and accidental coupling.

Legend: 🔴 must fix (bug/violation) · 🟠 should fix (design/robustness) · 🟡 improvement (maintainability/perf) · ⚪ nice-to-have.

Every finding was verified against the actual source; line numbers are current as of this review.

---

## 0. Executive summary

The project is in decent shape overall: layering (gui → handlers → services → entities) exists and is
mostly respected, the read-only wire schema is cleanly isolated in `internal/entities/template`, the
registry/constants split (pure SIDs vs display names) is a good idea, and the test suite around the
generator is substantial. The Gio usage is competent, including the non-trivial debounced
auto-regeneration and the in-app file explorer.

The main problems, in order of importance:

1. **One real layering violation** (`internal/services/content_rules` → `app/gui/constants`) that
   inverts the dependency direction of the whole architecture.
2. **`drivers.State` is a god object** (file I/O, generation orchestration, debouncing, manual-edit
   lifecycle, dialog opening, status bar — 455 lines, ~25 responsibilities) and it *bypasses the
   handler layer* by owning its own `mappers.GeneratorConfigMapper`.
3. **Swallowed errors and dead safety features**: `reapplyManualEdits` discards the handler error;
   the `unsaved` flag is never set to `true` so the documented "unsaved `*` marker" feature and the
   commented-out exit guard are both dead; `os.Exit` is called from UI code.
4. **A latent Windows path bug** in `internal/helpers/io.go` (`filepath.Join("C:", ...)` produces a
   *drive-relative* path on Go ≥ 1.20) plus an unchecked type assertion that panics on a corrupt
   Steam VDF.
5. **Large god files** (`zoneEditorDialog.go` 1078, `previewLayout.go` 1024, `topologyBase.go` 710+,
   `layoutPanel.go` 451, `state.go` 455) and **systematic duplication** across the ten topology
   services and across panels/dialogs.
6. **Per-frame waste** in the UI loop (full state snapshot + `reflect.DeepEqual` every frame,
   widget-slice rebuilds, zone-editor geometry recomputation without a dirty flag).
7. **CI is build+test only** — no vet/lint/race/coverage/vulncheck, and `gofmt` gating is blocked by
   CRLF (no `.gitattributes`).
8. **Repo hygiene**: a committed `hommoe_custom_templates.exe`, generated artifacts at the root and
   in `output/`, a stale `test/unit/services/oldTests` tree, an empty `internal/validators` package,
   a `temp_`-prefixed registry file that is actually live production code.

None of these require a rewrite. The refactoring plan below is incremental and behavior-preserving.

---

## 1. Verified bugs and correctness issues

### 1.1 🔴 `State.reapplyManualEdits` swallows the handler error
[app/gui/drivers/state.go](../app/gui/drivers/state.go#L436-L451)

```go
dto, _ := this.handler.UpdateTemplate(dtos.TemplateUpdateDto{ ... })
this.lastTemplate = dto.Template
```

`GUIHandler.UpdateTemplate` returns `common.ErrZonesMissing` when a manual connection references a
missing zone. That error is silently discarded, so the user gets no status message and a broken
template can be exported. (Note: this is *not* a nil-pointer crash — `TemplateLoadDto` is a value and
`dto.Template` aliases the input template — but the error loss is real.)

**Fix (exact):**

```go
dto, err := this.handler.UpdateTemplate(dtos.TemplateUpdateDto{
    Template:    this.lastTemplate,
    Zones:       append([]entities.Zone(nil), this.manualZones...),
    Connections: append([]entities.Connection(nil), this.manualConnections...),
    Config:      this.GetGeneratorConfig(),
})
if err != nil {
    this.SetStatus(fmt.Sprintf("Manual edits could not be fully reapplied: %v.", err), true)
}
this.lastTemplate = dto.Template
this.connectionsModified = true
```

The same pattern must be audited everywhere `_` receives an error: `grep -n ", _ :=" app/gui internal`
and justify or fix each hit.

### 1.3 🔴 `State.Exit` calls `os.Exit(0)` with the safety check commented out
[app/gui/drivers/state.go](../app/gui/drivers/state.go#L223-L229)

```go
func (this *State) Exit() {
    // if this.unsaved { ... }
    os.Exit(0)
}
```

`os.Exit` from a UI driver skips deferred cleanup and (per your own perf-test notes) has already
killed test processes when the Exit button was accidentally clicked during calibration.

**Fix:** after 1.2 makes `unsaved` meaningful, restore the guard and route shutdown through Gio
instead of `os.Exit`:

1. Give `State` an `onExit func()` callback (injected from `program.go`).
2. In `program.go`, pass `func() { window.Perform(system.ActionClose) }` so the normal
   `app.DestroyEvent` path runs.
3. `Exit()` becomes:

```go
func (this *State) Exit() {
    if this.unsaved {
        this.SetStatus("Unsaved changes exist - save first or press Exit again.", true)
        this.confirmExit = true // second press actually exits
        return
    }
    this.onExit()
}
```

(A tiny two-press confirm avoids building a confirm dialog; if you prefer a dialog, open one via
`this.dialogs.Open(...)` — the infrastructure already exists.)

### 1.5 🟠 Panics in library code
- [internal/services/previewAssets.go](../internal/services/previewAssets.go#L44) (twice): a failed
  embedded-PNG decode panics inside `sync.Once` during preview rendering. Embedded assets failing is
  a build-time bug, so a panic is *defensible*, but it will take the whole app down mid-frame with
  no message. Wrap: make `loadPreviewAssets() (*previewAssets, error)`, propagate through
  `RenderPreviewImage`/`WritePreviewPNG` (both already return `error`), and keep a `sync.OnceValues`
  memoization.
- [internal/models/neutralZoneProfile.go](../internal/models/neutralZoneProfile.go#L33):
  `panic("invalid quality")` on an enum the UI can never produce today — but the manual zone editor
  writes qualities too. Return `(NeutralZoneProfile, error)` or, minimally, fall back to the low
  profile and log.

### 1.6 🟠 `GUIHandler.UpdateTemplate` conflates "warning" with "error" and mutates its input
[internal/handlers/guiHandler.go](../internal/handlers/guiHandler.go#L45-L73)

- It writes into `templateDto.Template.Variants[0]` in place and *also* returns the same pointer in
  the DTO. Callers cannot tell whether to trust the returned value or the argument; combined with
  1.1, an "error" result still contains a fully updated template. Pick one contract:
  **pure-function style** — deep-copy the template, mutate the copy, return it; or **command style**
  — return only `error` and document in-place mutation. Recommended: keep in-place (cheap) but
  change the signature to `UpdateTemplate(dto dtos.TemplateUpdateDto) error` and delete the
  misleading `TemplateLoadDto` return.
- The `ComputeHasErrors → ErrZonesMissing` result is a *validation warning* about user data, not a
  failure of the operation. Model it as a typed result:

```go
type UpdateResult struct{ DanglingConnections []string }
func (h *GUIHandler) UpdateTemplate(dto dtos.TemplateUpdateDto) (UpdateResult, error)
```

- Line 52 TODO ("might not be needed") on `RebuildZoneConnectionRoads`: it *is* needed — it
  self-heals stale roads from already-saved `.gen.json` files. Replace the TODO with that sentence.

### 1.7 🟡 `go.mod` Go version vs CI mismatch
[go.mod](../go.mod#L3) says `go 1.26.3`; [.github/workflows/pr-validation.yml](../.github/workflows/pr-validation.yml)
pins `go-version: '1.25.8'`. `go build` with a toolchain older than the `go` directive fails (or
auto-downloads a toolchain, depending on `GOTOOLCHAIN`). Align them: either bump the workflow to
`1.26.x` or lower the directive to the minimum you actually need (`go 1.25`) — the directive should
be a *minimum*, not your current toolchain patch version.

### 1.8 🟡 `internal/registry/temp_bannableItemSidValues.go` is live code wearing a "delete me" sign
The file starts with `// remove/update`, is named `temp_…`, yet defines
`GetBannableItemSidValues()` which **is the production registry** used by
`app/gui/constants/bannableItems.go`. Rename the file to `bannableItemSidValues.go` (there is
currently no other file by that name — the memory of one is stale), delete the `// remove/update`
comment, and add the standard doc comment the other registry files have. Zero behavior change.

### 1.9 🟡 `internal/validators` is an empty package
[internal/validators/generatorConfigValidator.go](../internal/validators/generatorConfigValidator.go)
contains only `package validators` and a TODO. Empty packages confuse newcomers and tooling.
Either delete the directory (recommended until you actually build it) or implement the minimal
validator that `State.UpdateState`'s TODO (state.go line ~254) is asking for:

```go
package validators

func ValidateEditorState(s *dtos.EditorStateDto) []string {
    var problems []string
    if s.PlayerCount < 2 || s.PlayerCount > 8 {
        problems = append(problems, "player count must be 2-8")
    }
    if s.HeroCountMin > s.HeroCountMax {
        problems = append(problems, "hero count min exceeds max")
    }
    // clamp/flag zone counts, percents 25-200, etc.
    return problems
}
```

and call it from `UpdateState`, surfacing problems via `SetStatus`.

### 1.10 ⚪ Rejected findings (checked, **not** bugs — do not "fix")
For the record, these look suspicious but are correct:

- `zoneEditorDialog.deleteZone` (line ~999) `kept := connections[i]; append(..., &kept)` — each
  iteration creates a fresh variable; the pointers are valid. (Could still be simplified to
  `&connections[i]`, but it is not a bug.)
- `linq.QueryMap.ToMap()` ([internal/helpers/linq/map.go](../internal/helpers/linq/map.go#L74)) does
  `make(...)` correctly. It has **zero callers**, though — see §4.6 dead code.
- `buildNonAdjacentDerangement` (topologyBase.go tail) has a deterministic shift fallback after 100
  attempts — it never returns a partial result.
- `relaxPasses` (previewLayout.go ~686) already early-exits on `!moved`.
- `GUIHandler.UpdateTemplate` *does* check `len(Variants) == 0` before indexing.

---

## 2. Architecture (clean-architecture evaluation)

### 2.1 Current vs target layering

Current dependency edges (arrows = imports):

```
main → app/gui(program) → editor → drivers ─┬→ internal/handlers → internal/services → internal/{models,entities,registry,helpers}
                                            ├→ internal/mappers            (⚠ bypasses handlers)
                                            ├→ internal/models/config      (⚠ domain type in UI)
app/gui/dialogs → internal/models/config/config_inner                      (⚠ inner domain enum in UI)
internal/services/content_rules → app/gui/constants                        (🔴 INVERTED EDGE)
internal/mappers → internal/services/template_generator/providers          (⚠ mapper depends on provider)
```

Target:

```
app/gui (panels/dialogs/drivers/editor)  — knows only dtos + handler interface
        ↓
internal/handlers                        — the ONLY door into internal/*
        ↓
internal/services (+providers/builders)  — domain logic
        ↓
internal/{models, entities, registry, constants, helpers}  — leaf layers, no upward imports
```

### 2.2 🔴 Kill the inverted edge: `content_rules` → `app/gui/constants`
[internal/services/content_rules/variantMappingManager.go](../internal/services/content_rules/variantMappingManager.go#L6)
imports `app/gui/constants` for `constants.ContentIds.{DragonUtopia,PandoraBox,MontyHall}` — the
import even carries a `// TODO: This should not exist`.

**Fix (exact, mechanical):**
1. `app/gui/constants/contentIds.go` currently holds SID+display-name catalog entries. Move the
   *SID-bearing* struct down: create `internal/constants/contentIds.go` (package `constants` already
   exists at `internal/constants` with `topologies.go`) containing the `ContentIds` value verbatim.
2. Change `app/gui/constants/contentIds.go` to alias/re-export it
   (`var ContentIds = service_constants.ContentIds`) so no GUI call-site changes.
3. Point `variantMappingManager.go` at `internal/constants` and delete the TODO.
4. Enforce it forever with a `depguard` rule in `.golangci.yml` (see §7.2):
   deny `app/**` from any `internal/**` package.

### 2.3 🟠 `drivers.State` is a god object — split it and stop bypassing the handler
`State` currently owns: default-state construction, template-dir discovery, output path editor
widget, status bar, dialog host, save/load of `.gen.json`, save of `.rmg.json`+PNG, generation,
debounced auto-regeneration, manual-edit memory/persistence/reapply, and app exit. It also holds its
own `mappers.GeneratorConfigMapper` and `GetGeneratorConfig()` maps DTO→config **in the GUI layer**,
duplicating what `GUIHandler.GenerateTemplate` does internally.

**Step 1 — remove the mapper from the GUI.** `State.GetGeneratorConfig()` exists only to feed
`TemplateUpdateDto.Config`. Change `TemplateUpdateDto` to *not* carry a `*config.GeneratorConfig`;
instead carry the `dtos.EditorStateDto` (which the GUI legitimately owns) and let the **handler**
map it:

```go
// internal/dtos/templateUpdateDto.go
type TemplateUpdateDto struct {
    Template    *entities.RmgTemplate
    Zones       []entities.Zone
    Connections []entities.Connection
    State       *EditorStateDto // was: Config *config.GeneratorConfig
}

// guiHandler.go
if templateDto.State != nil {
    cfg := this.mapper.FromEditorState(*templateDto.State)
    templateDto.Template.MandatoryContent = this.contentProvider.CreateContentsForZones(cfg, ...)
}
```

Then delete `State.mapper`, `State.GetGeneratorConfig`, and the
`internal/mappers` + `internal/models/config` imports from
[app/gui/drivers/state.go](../app/gui/drivers/state.go#L15-L20). The GUI now touches only
`dtos` + `entities` (entities are needed for the zone editor) + the handler.

**Step 2 — split the file along its existing seams** (no public API change; methods keep the same
receiver so this is pure file reorganization):

| New file in `app/gui/drivers/` | Moves from state.go |
|---|---|
| `state.go` (~120 lines) | struct, `NewUIState`, accessors (`Status`, `IsUnsaved`, `OutputPath`, `Dialogs`, `LastTemplate`), `SetStatus`, `UpdateState`, `Exit` |
| `stateFiles.go` | `Load`, `SaveAs`, `Save`, `LoadStateFromFile`, `SaveStateToFile`, `handleLoadState`, `handleSaveState`, `handleSaveTemplate`, `SaveTemplate`, `PickOutputDir`, `RevealOutputDir` |
| `stateGeneration.go` | `Generate`, `AutoRegenerate`, `performAutoRegen`, `applyGeneratedTemplate`, `snapshotGeneratedState`, `clearGeneratedState`, `lastTemplateZoneAndConnectionCount`, `autoRegenDebounce` const |
| `stateManualEdits.go` | `ApplyEditedZones`, `rememberManualEdits`, `discardManualEdits`, `syncManualEditsToDto`, `restoreManualEdits`, `reapplyManualEdits`, `shouldReapplyManualEdits` |

**Step 3 (optional, later) —** extract the debounce into a reusable value type so it can be unit
tested without a `State`:

```go
// app/gui/drivers/debounce.go
type debouncer struct {
    pending  *dtos.EditorStateDto
    deadline time.Time
}
func (d *debouncer) Observe(now time.Time, current *dtos.EditorStateDto,
    equal func(a, b *dtos.EditorStateDto) bool) (fire bool, wakeAt time.Time, schedule bool)
```

`AutoRegenerate` shrinks to dispatch logic; the state machine gets table-driven tests.

### 2.4 🟠 Introduce a handler interface at the GUI boundary
`State` holds `*handlers.GUIHandler` concretely, which makes every `State` test drag the full
generator stack in. Define the port where it is consumed (Go convention):

```go
// app/gui/drivers/handler.go
type TemplateHandler interface {
    GenerateTemplate(dtos.EditorStateDto) (dtos.TemplateLoadDto, error)
    UpdateTemplate(dtos.TemplateUpdateDto) error            // after §1.6
    SaveTemplate(dtos.TemplateSaveDto) (string, error)
    LoadState(path string) (*dtos.EditorStateDto, error)
    SaveState(path string, s dtos.EditorStateDto) error
}
```

`NewUIState` keeps constructing the real one (`handlers.NewGuiHandler()` satisfies it implicitly),
but tests can inject a fake and the compile-time surface between layers becomes explicit. This is
the single highest-leverage change for making `drivers` unit-testable.

### 2.5 🟠 GUI dialogs importing `config_inner`
[app/gui/dialogs/bonusPickerDialog.go](../app/gui/dialogs/bonusPickerDialog.go) imports
`internal/models/config/config_inner` for `BonusEntry`/`BonusPresetType`. Two acceptable fixes:

- **Cheap (recommended):** import the alias package `internal/models/config` instead —
  `config.BonusEntry`, `config.BonusTownPortalFree`, … already exist in
  [internal/models/config/types.go](../internal/models/config/types.go). Rule: *nothing outside
  `internal/models/config` may import `config_inner`*. Add a depguard rule.
- Thorough: move bonus serialization types into `internal/dtos` since `EditorStateDto.BonusesJSON`
  already exposes them at the DTO layer. Bigger churn, same effect.

### 2.6 🟡 Collapse the `config` / `config_inner` split
`internal/models/config/types.go` is 38 lines of pure aliases; `config_inner` exists (per git
history) to break an import cycle that no longer exists. Two-package indirection for one logical
model confuses newcomers ("which one do I import?").

**Fix:** move every file from `config_inner/` up into `config/`, delete `types.go`, run
`go build ./...` and fix the ~10 import paths (`config_inner.` → `config.`). Pure rename, zero
behavior change. Do this *after* 2.5 so only one layer needs touching.

### 2.7 🟡 `internal/mappers` depends on a provider
[internal/mappers/generatorConfigMapper.go](../internal/mappers/generatorConfigMapper.go#L24) calls
`providers.NewMandatoryContentProvider().CreateContentItemsFrom(rows)` — a *mapper* invoking a
*generation provider*, and constructing it per call. `CreateContentItemsFrom` is a pure
rows→`[]MandatoryContentItem` conversion. Move that function (and its private helpers) into the
mappers package or into `internal/models` next to `ZoneContentRowSave` (it is the natural home for
"how does a saved row become a content item"). Then `mappers` depends only on
`dtos`/`models`/`registry` and the dependency arrow points strictly downward.

### 2.8 🟡 Formalize the topology plug-in seam
All ten topology services share the same shape but there is no interface; `topologyProvider.go`
switches on the enum and calls concrete types. Add:

```go
// internal/services/template_generator/providers/topology/topologyService.go
type Service interface {
    CreateTopologyVariant(cfg config.GeneratorConfig, playerLabels []string,
        neutralZones models.NeutralZonePlans, tuning models.GenerationTuning,
        holdCityLabel string) entities.Variant
}
```

and register implementations in a map keyed by `config.MapTopology`:

```go
var services = map[config.MapTopology]Service{
    config.TopologyDefault:  NewRingTopologyService(),
    config.TopologyChain:    NewChainTopologyService(),
    // ...
}
```

`topologyProvider.go` becomes a lookup + tournament short-circuit. Your own memory notes list six
files to touch when adding a topology; this removes one of them and gives a compile error instead of
a silent fall-through when a new enum value is missed (add an explicit
`if svc, ok := services[t]; !ok { return defaultService }`). Same treatment for the
`tournament_variant` cluster services (interface already exists in
`tournament_variant/interface.go` — use it in the switch).

### 2.9 ⚪ Things that are *right* — keep them
- `internal/entities/types.go` aliasing the deep template packages into one import — good; keep.
- Registry pattern (private struct + getter): verbose but safe and greppable. Do **not** replace it
  with a stringly-typed map or a codegen framework; the only fixes needed are the `temp_` rename
  (§1.8) and doc comments. Consistency > cleverness here.
- The `//go:build integration_test` gating of test exports is a clean solution; keep the scope rule.
- `dtos.ManualZoneSave`/`ManualConnectionSave` wrappers: slightly verbose but the correct way to
  persist `json:"-"` runtime fields without polluting the wire schema. Keep.

---

## 3. God files — concrete decomposition plans

These are pure file moves within the same package (methods keep receivers); each is a
zero-behavior-change commit that makes review and navigation dramatically easier.

### 3.1 🟠 `app/gui/dialogs/zoneEditorDialog.go` (1078 lines → 4 files)

| File | Contents |
|---|---|
| `zoneEditorDialog.go` (~300) | struct, constructor, `Title`, `PreferredSize`, `Body`, toolbar/footer/status layout, apply/cancel/reset plumbing |
| `zoneEditorCanvas.go` (~330) | `layoutCanvas`, `handlePointer`, press/drag/release, `addZoneAt`, `moveDraggedZone`, `hitTestNode`, `hitTestEdge`, `recomputeGeometry`, `obstacleBulge`, drawing helpers |
| `zoneEditorConnectionProps.go` (~230) | side-panel connection editing: `propertyRows`, `syncPropsFromConnection`, `writebackProps`, `addConnection`, `deleteConnection` |
| `zoneEditorZoneProps.go` (~220) | `zonePropertyRows`, `syncZoneProps`, `writebackZoneProps`, `selectZone`, `selectedZoneRef`, `zoneByName`, `deleteZone`, `resetToOriginal` |

While splitting, extract the magic numbers into named constants at the top of `zoneEditorCanvas.go`:
`const (edgeHitSlopPx = 9.0; nodeHitSlopPx = 4.0; nodeHitSlopSelectedPx = 8.0; bulgeGapPx = 18.0)`.

### 3.2 🟠 `internal/services/previewLayout.go` (1024 lines → 6 files, same package)

| File | Contents |
|---|---|
| `previewLayout.go` (~120) | `PreviewLayout` structs, constants block, `BuildPreviewLayout` dispatcher, `buildPreviewConnections` |
| `previewLayoutManual.go` | `layoutManualPositions` + parallel-gap logic |
| `previewLayoutRings.go` | `layoutBalancedRings`, `allHaveRing` |
| `previewLayoutFixed.go` | `layoutFixedPositions`, `isFixedGeometryTopology`, `allHavePosition` |
| `previewLayoutScatter.go` | `layoutScatter`, `relaxPasses`, `isScatterTopology` |
| `previewLayoutHub.go` | `layoutRingOrHub`, `layoutMultiHub`, `connectedComponents`, angle/centroid helpers |

The dispatch-order comment (manual > rings > fixed > scatter > ring/hub) should live as a doc
comment on `BuildPreviewLayout` — it is the single most important invariant in the file and today it
only exists in your notes.

### 3.3 🟠 `topology/base/topologyBase.go` (710+ lines → 4 files, same package)

| File | Contents |
|---|---|
| `topologyBase.go` (~200) | struct, embeds, `CreateVariant`, label/name plumbing |
| `zoneFactory.go` | `CreateSpawnZone`, `CreateNeutralZone`, `CreateHubZone`, `createPlayerSpawnCastle`, `createPlayerOwnedCastles`, `createPlayerUnclaimedCastles`, `CreateNeutralZoneCastles`, `createAbandonedOutposts`, hold-city helpers |
| `connectionFactory.go` | `CreateRandomPortalConnections`, `CreateMissingPlayerConnections`, `CreateMissingConnections`, `GetBorderGuardValue`, `buildNonAdjacentDerangement` |
| `roadFactory.go` | `createOuterZoneRoads`, `CreateConnectorZoneRoads`, `buildSideContentLimits` |

Also fix here: 🟡 `CreateNeutralZoneCastles` mutates its `castleCount` parameter for the hold-city
case (line ~571) — use a local `effectiveCount` for readability; and 🟡 the module-level
`var resourceContentPool = registry.GetResourcesContentPoolValues()` globals are fine functionally
(read-only) but move them into the struct during the split so the base has no package-level state.

### 3.4 🟡 `app/gui/panels/layoutPanel.go` (451) and `app/gui/drivers/state.go` (455)
`state.go` split is specified in §2.3. For `layoutPanel.go`, split by UI column:
`layoutPanel.go` (struct + Load/SaveToState + `GetPanelWidget`), `layoutPanelLeft.go`
(topology/connectivity/sizes/difficulty section builders), `layoutPanelZones.go`
(zones + advanced tier sections + hub sliders). Delete the commented-out
`sldMinNeutralBetween` row (line ~259–262) — the *slider field* is live (it maps to
`MinNeutralZonesBetweenPlayers`), only the dead commented row should go, along with the
"Investigate this" TODO once you confirm the wired slider works.

---

## 4. Duplication and simplification

### 4.1 🟠 Topology services: extract the shared variant-assembly skeleton
Every service (`ring`, `chain`, `hub`, `random`, `circles`, `square`, `geometric`, `cross`,
`fractal`, plus 3 tournament cluster services) repeats:

1. `CreateOrderedZoneLabels(...)`
2. a `createZones` loop that is 90% identical: *"if label is a player label → CreateSpawnZone else →
   CreateNeutralZone(plan, …, label == holdCityLabel)"*
3. a `createConnections` loop building `NewConnectionBuilder().WithName(prefix-A-B).WithFrom…`
   with only the name prefix and guard-match-group prefix differing
4. optional `CreateRandomPortalConnections` / `CreateMissingPlayerConnections` appends
5. `CreateVariant(...)`

**Fix — add two helpers to `topology/base` and rewrite services to use them:**

```go
// base/zoneAssembly.go
type ZoneSpec struct {
    Config        config.GeneratorConfig
    Tuning        models.GenerationTuning
    PlayerLabels  []string
    NeutralPlans  models.NeutralZonePlans
    HoldCityLabel string
    ConnNamesFor  func(label string) []string
}

// CreateZonesForLabels builds one zone per ordered label, choosing spawn vs
// neutral exactly like every topology service does today.
func (b *TopologyBase) CreateZonesForLabels(spec ZoneSpec, orderedLabels []string) []entities.Zone

// base/connectionAssembly.go
// CreateDirectConnection builds the standard direct connection both endpoints
// share in every topology; prefix distinguishes e.g. "Ring"/"Rnd"/"Chain".
func (b *TopologyBase) CreateDirectConnection(prefix, from, to string,
    guardValue int) entities.Connection
```

Measured against the current files this removes roughly 40–60 lines from each of the ten services
(≈450 lines total) and — more importantly — means a change to the standard connection shape
(e.g. a new guard field) happens in one place. The geometric family already proves the pattern
works: `circles/square/geometric/cross` reuse `RandomTopologyService.createZones/createConnections`
via `(allLabels, positions, pairs)`; this generalizes that reuse to the structured topologies too.

### 4.2 🟠 Panels: extract the row-building vocabulary
`generalPanel.go`, `layoutPanel.go`, `bonusesPanel.go` repeat the same three-line motifs dozens of
times: *labeled slider row with a `RoundedRangeString` value*, *labeled checkbox row*, *section with
rows*. `widgets.NewLabeledRowWidget` exists but each call site still assembles
label-width/format/slider by hand. Add two thin factories to `app/gui/widgets`:

```go
// widgets/sliderRow.go
func NewSliderRow(theme *material.Theme, label string, labelWidth unit.Dp,
    f *widget.Float, format func(v float32) string) layout.Widget

// widgets/checkboxRow.go  (variant of existing labeledCheckboxRow with hint text)
func NewCheckboxRow(theme *material.Theme, cb *widget.Bool, label, hint string) layout.Widget
```

Then panel section functions become declarative lists. Expected reduction: ~120 lines across the
three panels, and a new control means one line, not five.

### 4.3 🟠 The bonuses panel "list section" pattern
[app/gui/panels/bonusesPanel.go](../app/gui/panels/bonusesPanel.go) implements
add-button + item rows + per-row remove ✕ **four times** (bonuses, banned items, banned magics,
value overrides). Extract once:

```go
// app/gui/widgets/listSection.go
type ListSection struct {
    Title   string
    AddBtn  widget.Clickable
    Removes []widget.Clickable // resized to item count each frame
}
// Process returns the index clicked for removal (or -1) and whether Add was clicked.
func (s *ListSection) Process(gtx layout.Context) (removeIdx int, addClicked bool)
func (s *ListSection) Layout(theme *material.Theme, itemCount int,
    row func(i int) layout.Widget, emptyHint string) layout.Widget
```

This also centralizes the "poll Clickables **before** layout" rule (your `processClicks` gotcha)
in one audited place instead of four hand-rolled copies.

### 4.4 🟡 Picker row rendering
`pickerDialog.go` (leaf rows) and `bonusPickerDialog.go` (spell rows) both render
"colored dot + name + secondary text + trailing control". Extract
`widgets.NewEntryRowWidget(theme, dot color.NRGBA, name, secondary string, trailing layout.Widget) layout.Widget`
and use it in both; the spell-school color and bonus-dot color already come from
`themes/accents.go`, so no new coupling.

### 4.5 🟡 Zone editor property rows
`propertyRows` (connection props) and `zonePropertyRows` duplicate the
"label + dropdown/editor flex row" assembly (~75 lines). After the §3.1 split, extract within the
dialogs package: `func propertyRow(theme *material.Theme, label string, control layout.Widget) layout.Widget`.

### 4.6 🟡 Dead code inventory (delete all of it)
- `linq.QueryMap.ToMap`, and any other `linq` methods without callers —
  `internal/helpers/linq` is used in ~4 production sites; run
  `staticcheck -unused` and delete unreferenced methods. Alternatively (recommended): replace the
  four production uses (`zoneLabelProvider.go`, `neutralZonePlan.go`, `topologyBase.go`) with plain
  loops/`slices` calls and demote `linq` to a test-only helper or remove it. A custom LINQ layer is a
  heavy dependency for four call sites, and it produced the `First`-returns-pointer aliasing subtlety
  in `topologyBase.go` (~line 322) that only works by luck of Go slice semantics.
- Commented-out icon experiment in [app/gui/editor/toolbar.go](../app/gui/editor/toolbar.go#L33-L51) (19 lines).
- Commented-out registry vars atop `variantMappingManager.go` (lines 10–13).
- Commented-out exit guard (resolved properly by §1.2/§1.3).
- The commented-out "Min neutrals between players" row in `layoutPanel.go` (§3.4).
- `app/gui/drivers/state_testexports.go` is a placeholder with only a TODO — delete until needed.
- TODO in `app/gui/constants/victoryConditions.go:50` ("probably should return empty Victory...
  suck it") — resolve it: return the zero `VictoryCondition` and log; and clean up the comment.

### 4.7 ⚪ Builder boilerplate
The 8 builders in `internal/services/builders/variant_content` are repetitive but *explicit and
correct* — the fluent API reads well at call sites and mirrors the C# reference. Do **not**
genericize them (a `Builder[T]` erases field-level discoverability). Only fix: defensive-copy the
two slice-accepting setters if you keep sharing config slices (`WithGuardedContentPool`,
`WithMainObjects`) — one `slices.Clone` each — *or* document "caller must not retain the slice".
`cloneContentItems` in the mandatory-content provider already exists because this bit once; make
the builders safe at the source instead.

---

## 5. Performance (UI hot path)

Ordered by measured impact on a 60 fps frame budget; none of these change behavior.

### 5.1 🟠 Stop snapshotting + deep-comparing the whole DTO every frame
`Window.Layout` → `save()` (every tab's `SaveToState`, every frame) → `State.AutoRegenerate` →
`EqualsIgnoringManualEdits` → **two struct copies + `reflect.DeepEqual` over a 100+-field DTO with
five `[]ZoneContentRowSave` slices, every frame**, even when idle.

**Fix — a change counter, not comparisons:**

1. Give `State` a `revision uint64`. `UpdateState` (the single mutation funnel — after §1.2 all
   mutations already flow through it or through panel `SaveToState`) increments it *only when it
   changed something* (it already computes that for `unsaved`; same comparison, so compute once and
   use it for both).
2. Panels: make `SaveToState` report `changed bool` (each panel already knows which widget changed
   this frame via `.Update(gtx)`; where it currently blindly copies, compare the single field it
   writes). `Window.save()` ORs the results and bumps the revision once.
3. `AutoRegenerate` compares `revision != lastGeneratedRevision` — an integer compare — and only
   falls back to `EqualsIgnoringManualEdits` when deciding *debounce vs immediate* (i.e. at most
   once per actual change, not per frame).

Interim cheap win if step 2 is too invasive at first: replace `reflect.DeepEqual` with a generated
or hand-written `Equals` method — the DTO has only slices of comparable structs, so
`slices.Equal`+field compares are ~20 lines and an order of magnitude faster. But do the revision
counter; it also gives you `unsaved` (§1.2) for free.

### 5.2 🟠 Zone editor: dirty-flag the geometry
`ZoneEditorDialog.recomputeGeometry` runs `BuildPreviewLayout` (which can include 500 relax passes
for scatter topologies) plus O(edges×zones) `obstacleBulge` **every frame while the dialog is open**.
Add `geomDirty bool`, set it in every mutator (`moveDraggedZone`, `addZoneAt`, `deleteZone`,
`addConnection`, `deleteConnection`, `writeback*Props`, `resetToOriginal`, and on canvas-size
change), and guard:

```go
if this.geomDirty || side != this.lastSide {
    this.recomputeGeometry(side)
    this.geomDirty = false
    this.lastSide = side
}
```

Keep the existing event-ordering invariant (pointer handling *before* recompute) — the dirty flag
does not disturb it because pointer handlers are exactly the places that set the flag. Also hoist
the per-frame `groups` map in `recomputeGeometry` to a struct field cleared with `clear(m)`.

### 5.3 🟡 Cache per-frame string formatting and widget slices
- `zoneEditorDialog.layoutStatus`: build the status string only when
  zones/connections/mode/hint change (they are all already tracked); cache in a field.
- `bonusesPanel.buildWidgets` rebuilds 4 slices of closures per frame. `material.List` is already
  index-based — replace the pre-built `[]layout.Widget` with a `func(index int)` dispatch over the
  section boundaries (`if index < len(header)+len(bonuses) …`), or at minimum reuse a persistent
  slice with `s = s[:0]`.
- `Window.getTabsWidget` allocates the `[]layout.FlexChild` every frame for a tab strip that
  changes never — build it once in `NewWindow` and rebuild only on tab-set change.
- `ZoneContentDialog.persist()` fires `onApply` (→ full rows clone → state write) every frame while
  the dialog is open. Track `changed` from the section widgets' click processing and only persist
  on change frames.

### 5.4 ⚪ Non-issues (measured/structural — skip)
Generation itself (`TemplateGenerator.Generate`) runs at human interaction frequency (debounced
300 ms) on ≤ ~40 zones; its O(n²) label loops are irrelevant. Don't optimize the generator until a
profile says otherwise.

---

## 6. Error handling, robustness, and API polish

### 6.1 🟠 Adopt one error-wrapping convention
`internal/common/editorErrors.go` has good sentinels; services mostly return raw `err`. Convention
to apply mechanically (low risk, high debuggability):

- I/O boundaries wrap with context: `fmt.Errorf("load settings %q: %w", path, err)` in
  `settingsFileLoader.go`, `templateWriter.go`, `helpers/io.go`, `previewRenderer.go` (file writes).
- Domain failures use the existing sentinels (`common.Err…`), matched with `errors.Is`.
- No panics outside `main`/`init`-time invariants (§1.5).

### 6.2 🟠 Guard the frame loop
[app/gui/program.go](../app/gui/program.go): a panic in any widget kills the app with work lost
(and given §1.5's panics, this is reachable). Wrap the frame:

```go
case app.FrameEvent:
    gtx := app.NewContext(&ops, event)
    func() {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("panic in frame: %v\n%s", r, debug.Stack())
                windowLayout.State().SetStatus("Internal error - see log.", true)
            }
        }()
        windowLayout.Layout(gtx, theme)
    }()
    event.Frame(gtx.Ops)
```

One recover, at the outermost sensible point, logging the stack — not scattered recovers.

### 6.3 🟡 Receiver naming: retire `this`
Every method in the repo uses `(this *T)`. It compiles, but it fights Go tooling/convention
(golangci-lint's default `revive`/`stylecheck` flag it; every new Go contributor will trip). Do a
mechanical rename per type (`gopls rename` or `gofmt -r` per file): `State → s`,
`ZoneEditorDialog → d`, `TopologyBase → b`, `GUIHandler → h`, `Query → q`, etc. Purely cosmetic —
schedule it as one atomic PR with no other changes so it doesn't pollute diffs forever. If you
consciously prefer `this`, then instead disable the corresponding linter rules explicitly in
`.golangci.yml` and document the choice in `.github/AGENTS.md`; the worst option is the current
implicit inconsistency with the ecosystem.

### 6.4 🟡 Godoc the exported surface
`internal/services`, `internal/models`, and `internal/registry` export dozens of undocumented
symbols. Priorities: `BuildPreviewLayout` (dispatch invariants, §3.2), the Bowyer-Watson code in
[internal/models/position.go](../internal/models/position.go) (150 lines of geometry with no
explanation — document algorithm, inputs, complexity), `GeneratorConfig` field semantics, every
registry getter (one line each: what game concept the SIDs identify). Enforce with `revive`'s
`exported` rule scoped to `internal/services/...` once the backlog is cleared.

### 6.5 🟡 Long parameter lists
`CreateSpawnZone`/`CreateNeutralZone`/`createZones` variants take 6–9 positional parameters of
which 4 are always `(config, tuning, plans, holdCityLabel)`. The `ZoneSpec` struct from §4.1 fixes
the worst offenders as a side effect; apply the same idea to `CreateHubZone` if it grows again.

---

## 7. Testing and CI

### 7.1 🟠 Test-suite consolidation
- **`test/unit/services/oldTests/` (~2 500 lines)**: ~70 % duplicates the newer suites. But five
  tests are *unique and valuable*: the all-topologies matrix stress test, city-hold flag
  propagation, abandoned outposts, player-owned castles, and the tournament cluster-isolation
  floodfill. **Plan:** move those five into
  `test/unit/services/template_generator/templateGenerator/` (rewriting assertions to testify +
  `t.Run` subtests), then delete the whole `oldTests` tree in the same PR. Do not keep a
  "deprecated" directory — it will never be cleaned later.
- **`test/helpers/defaultTemplate.go` (1365 lines)**: replace the hand-maintained struct literal
  with a golden JSON file. Add `test/testdata/golden/default.rmg.json` (generate it once from the
  current fixture to guarantee equivalence), a loader in `test/helpers`, and a
  `-update` flag pattern:

```go
var update = flag.Bool("update", false, "rewrite golden files")

func AssertGoldenTemplate(t *testing.T, got *entities.RmgTemplate, name string) {
    t.Helper()
    path := filepath.Join("testdata", "golden", name+".rmg.json")
    if *update { writeGolden(t, path, got) }
    want := readGolden(t, path)
    require.Equal(t, want, got) // testify prints a usable diff
}
```

  Benefits: schema changes become a reviewed one-command regen instead of a 1365-line hand edit,
  and failures print field-level diffs.
- **Coverage gaps (add, in this order):**
  1. `internal/mappers` — a `FromEditorState` field-mapping test (fill every DTO field with
     `gofakeit`, assert the config; this is the class of bug — a forgotten mapping line — most
     likely to actually happen here).
  2. `internal/handlers` — table tests over the 5 methods with a stubbed filesystem path
     (`t.TempDir()`); they are thin, so this is an afternoon.
  3. `dtos` round-trip: `EditorStateDto` marshal→unmarshal→`EqualsIgnoringManualEdits` equality,
     incl. manual-edit saves (this locks the `.gen.json` format).
  4. `drivers.State` debounce state machine — possible today with a fake clock since
     `AutoRegenerate(now)` already takes time as a parameter; becomes trivial after §2.4's
     interface.
  5. `helpers/io.go` — after §1.4's refactor to lazy functions, test `getBasePath` with fixture
     VDF maps (no Steam needed) and add corrupt-VDF cases.

### 7.2 🟠 CI pipeline
[.github/workflows/pr-validation.yml](../.github/workflows/pr-validation.yml) currently: deps →
build → 3 test runs. Missing, in priority order:

1. **`.gitattributes`** (new file, root — prerequisite for any format gate):

```gitattributes
* text=auto
*.go text eol=lf
*.md text eol=lf
*.json text eol=lf
*.yml text eol=lf
*.png binary
*.exe binary
```

   Then one-time renormalize: `git add --renormalize . && git commit -m "normalize line endings"`.
   (Working-tree files stay CRLF on Windows checkouts; gofmt checks run on LF in CI.)
2. **Race detector** on unit tests: `go test -race -count=1 ./...` — the perf suite already proved
   the editor window is shared across goroutines; `-race` on Linux CI needs no cgo setup.
3. **golangci-lint** with a minimal, high-signal config (new `.golangci.yml`):

```yaml
run:
  timeout: 5m
  build-tags: [integration_test]
linters:
  enable: [govet, staticcheck, errcheck, ineffassign, unused, gofmt, depguard]
linters-settings:
  errcheck: { check-type-assertions: true }
  depguard:
    rules:
      no-ui-from-internal:
        files: ["**/internal/**"]
        deny:
          - pkg: github.com/Tariomka/hommoe_custom_templates/app
            desc: internal/* must not import app/*
      config-inner-private:
        files: ["!**/internal/models/config/**"]
        deny:
          - pkg: github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner
            desc: import internal/models/config instead
```

   The two depguard rules make §2.2 and §2.5 permanent.
4. **`go mod tidy` gate**: `go mod tidy && git diff --exit-code go.mod go.sum`.
5. **govulncheck**: `golang/govulncheck-action@v1`.
6. Coverage: collect (`-coverprofile`) and *report* first; only add a threshold gate after §7.1
   closes the gaps, or the gate blocks unrelated PRs.
7. Fix the Go-version mismatch (§1.7) in the same PR.

### 7.3 🟡 Test style
Standardize on testify + table-driven `t.Run` subtests (the newer suites already do). When
migrating the five `oldTests` (§7.1), convert their manual `if err != nil { t.Fatalf }` chains.

---

## 8. Repo hygiene

| Item | Action |
|---|---|
| 🔴 `hommoe_custom_templates.exe` committed at root | `git rm --cached hommoe_custom_templates.exe`; ensure `.gitignore` has `*.exe` (it does — the file predates it) |
| 🟠 `Custom Template.gen.json`, `Initial.gen.json` at root | Move to `data/ExampleSettings/` (they are user-facing examples) or `test/testdata/fixtures/` if only tests use them; update any hardcoded paths |
| 🟠 `output/Colosseum/*` committed | Generated artifacts: `git rm -r --cached output/` and add `output/` to `.gitignore` |
| 🟡 `README.md` stale paths | Still references the old `internal/gui/` layout in the structure section — regenerate the tree from the current repo |
| 🟡 No `ARCHITECTURE.md` | Add one page: the layer diagram from §2.1, the data flow `EditorStateDto → GeneratorConfig → RmgTemplate → .rmg.json + preview PNG`, the topology plug-in checklist (currently only in your private notes), and the `integration_test` tag rules. This is the single best onboarding artifact you can add |
| ⚪ `app/tui/readme.md`, `app/web/readme.md` | Placeholders are fine; add one line each stating they are intentional future work so nobody deletes them |
| ⚪ `todo/` | Keep review docs here (this file), but move actionable items to issues once triaged |
| ⚪ `go.mod` test-only deps | `testify`/`gofakeit` in `require` is normal Go — **no action** (ignore any advice to create a `tools.go` for them; that pattern is for CLI tools, not test libraries) |

---

## 9. Prioritized execution plan

Each phase is independently shippable; the app behaves identically (or strictly better) after each.

**Phase 1 — correctness & safety (small diffs, do first)**
1. §1.1 propagate `UpdateTemplate` error in `reapplyManualEdits`.
2. §1.4 io.go: drive-relative paths, unchecked assertion, lazy env reads (+ tests §7.1.5).
3. §1.2 + §1.3 `unsaved` flag + exit guard via `onExit` callback.
4. §6.2 frame-loop recover; §1.5 de-panic previewAssets/neutralZoneProfile.
5. §1.7 Go version alignment; §8 remove committed exe + output artifacts.

**Phase 2 — architecture (medium diffs, mechanical)**
6. §2.2 move ContentIds down; add depguard.
7. §2.3 step 1: remove mapper/config from `drivers.State`; slim `TemplateUpdateDto`.
8. §2.5 + §2.6 config_inner consolidation.
9. §2.7 move `CreateContentItemsFrom` out of providers.
10. §2.4 `TemplateHandler` interface; §1.6 `UpdateTemplate` contract cleanup.
11. §1.8 rename temp_ registry file; §1.9 delete or implement validators; §4.6 dead-code sweep.

**Phase 3 — structure & duplication (large but zero-risk file moves)**
12. §3.1–§3.4 god-file splits (one PR per file).
13. §4.1 topology zone/connection assembly helpers; §2.8 topology Service interface + registry map.
14. §4.2–§4.5 widget/row extraction in panels and dialogs.

**Phase 4 — performance**
15. §5.1 revision counter replacing per-frame DeepEqual.
16. §5.2 zone-editor geometry dirty flag.
17. §5.3 cached strings/slices; per-change `persist()`.

**Phase 5 — tests & CI**
18. §7.2 .gitattributes → renormalize → golangci-lint + race + tidy + govulncheck in CI.
19. §7.1 oldTests migration + deletion; golden-file fixture; mappers/handlers/dtos/debounce tests.

**Phase 6 — polish**
20. §6.3 receiver rename (single atomic PR).
21. §6.4 godoc backlog; §8 ARCHITECTURE.md + README refresh.

---

## Appendix A — findings investigated and rejected

Documented so future reviews don't re-litigate them:

| Claim | Verdict |
|---|---|
| `zoneEditorDialog.deleteZone` stores pointers to a loop variable → dangling | **False.** Per-iteration variable; pointers are valid Go. |
| `linq.ToMap()` writes to a nil map | **False.** Uses `make`; the old review note is stale. It *is* dead code though (§4.6). |
| `buildNonAdjacentDerangement` can return a partial slice after 100 failed tries | **False.** Deterministic shift fallback exists. |
| `relaxPasses` lacks convergence early-exit | **False.** `if !moved { break }` present. |
| `GUIHandler.UpdateTemplate` indexes `Variants[0]` unchecked | **False.** Length-checked on entry. |
| `reapplyManualEdits` nil-derefs `dto.Template` on error | **False as stated** (value DTO), but the swallowed error is real (§1.1). |
| Move testify/gofakeit out of go.mod `require` | **Rejected.** Test-only deps in `require` is standard Go module behavior. |
| Replace registry pattern with dynamic/codegen registry | **Rejected.** Current pattern is verbose but safe and greppable; only naming/doc fixes needed. |
| Genericize the fluent builders | **Rejected.** Explicit builders are the point; only slice-copy safety needed (§4.7). |
