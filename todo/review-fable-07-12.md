# Project Review — 2026-07-12 (fable)

Senior/Principal-level review of the full repository, verified against the source as of
2026-07-12 (`go1.26.3`, golangci-lint v2.12.2 run with 113 issues). This document
**supersedes and consolidates** `todo/review.md`, `todo/review-fable.md` (both 2026-06-28)
and `todo/test_observations.md`, which are deleted with this commit. Section 0 records the
disposition of every prior item so nothing is silently dropped.

Severity legend: 🔴 High (bug / correctness / user-visible), 🟠 Medium (architecture,
performance, CI gaps), 🟡 Low (readability, hygiene), ⚪ Informational / no action.

---

## 0. Disposition of prior reviews (2026-06-28)

### 0.1 Fixed since the last reviews ✅

| Prior item | Evidence |
|---|---|
| 🔴 Layering violation `internal/services` → `app/gui/constants` (fable §1.5, review §3.1) | Gone from production code; `depguard` rule in [.golangci.yml](../.golangci.yml) now denies `internal → app`. Residue: 3 **test** files still import `app/gui/constants` (see §2.7). |
| 🔴 `State.reapplyManualEdits` swallows handler error (fable §1.1) | [state.go](../app/gui/drivers/state.go#L316-L343) `handleUpdateTemplate` now checks `err`, distinguishes `ErrProvidedTemplateInvalid`, and surfaces "fix before export" status. |
| 🔴 `linq.ToMap` nil-map panic (review §1.1) | No longer present in [internal/helpers/linq/](../internal/helpers/linq/). |
| 🔴 VDF unchecked type assertions (review §1.2) | [io.go](../internal/helpers/io.go#L112-L143) `getBasePath` uses comma-ok on every assertion. |
| 🔴 Hard-coded `C:\...\USERNAME` path (review §1.3) | [io.go](../internal/helpers/io.go#L22) now uses `os.UserHomeDir()` and `%ProgramFiles(x86)%`; only a last-resort literal fallback remains (see §1.7). |
| 🟠 Duplicate preview-layout implementation (test_observations #28) | `internal/services/previewLayout.go` deleted; only [previewLayoutService.go](../internal/services/preview_service/previewLayoutService.go) remains. |
| 🟠 `oldTests/` duplicate suite (~2,500 LOC) (fable §7.1a, §8.2) | Directory deleted; unit suite reorganized per-function under [test/unit/](../test/unit/). |
| 🟠 Hand-maintained 1,365-line golden fixture (fable §7.1b) | Replaced by [defaultTemplate.json](../test/test_helpers/defaultTemplate.json) + a 38-line loader in [defaultTemplate.go](../test/test_helpers/defaultTemplate.go). |
| 🟠 CI missing lint/tidy/coverage gates (fable §7.2, review §6.1) | [pr-validation.yml](../.github/workflows/pr-validation.yml) now has go-mod-tidy check, golangci-lint v2.12.2, build, unit tests with `-coverpkg`, coverage report with decrease-fail + 60% floor, and integration+performance jobs. `.gitattributes` exists (CRLF prerequisite done). |
| 🟠 `temp_bannableItemSidValues.go` (fable §1.3) | File deleted. |
| 🟡 README stale paths (fable §8.3) | README now documents the current `app/gui` layout and data flow. |
| 🟡 Missing exit guard (fable §1.2 part) | [state.go](../app/gui/drivers/state.go#L163-L172) now has an unsaved-changes double-press guard (remaining issues in §1.2). |
| 🟡 `test_helpers` "Most likely this is random" staleness | Fixture now normalizes `GeneratorPosition` after load — intent documented. |

### 0.2 Invalidated / accepted as convention ✖

| Prior item | Reason |
|---|---|
| `this` receiver rename (fable §6.3, review §3.3) | Documented house style in [AGENTS.md](../AGENTS.md) §4.3; `recvcheck`/revive configured around it. Not re-raised. |
| `internal/registry` consolidation / code-gen (older survey idea) | Registry is a **read-only game-data directory** per AGENTS.md §2.1; restructuring is out of bounds. The 6 `gochecknoglobals` hits there should be excluded via lint config, not edits (§8.2). |
| Collapse `config`/`config_inner` (fable §2.4) | Low value; alias layer is stable, covered by tests, and touching it risks churn in read-only-adjacent schema code. Drop unless a future change forces it. |
| Most `test_observations.md` entries | They documented **defensive branches that need no action** (marshal-error guards, embedded-asset errors, derangement fallback, etc.). The testing philosophy they encode now lives in AGENTS.md §4.6. Only two actionable items are carried forward: §1.1 (road mutations) and §6.2 (io.go injectability). |

### 0.3 Carried forward ❗ — detailed in the sections below

Lost road mutations (§1.1), `os.Exit` (§1.2), gameMode load bug (§1.3), bridge-name
infinite-loop risk (§1.4), victory-condition coercion (§1.5), `UpdateTemplate` aliasing
(§1.6), `resolveGlob`/Steam fallback (§1.7), empty validators package (§2.1), god files
(§2.2–2.4), topology template-method duplication (§3.1), widget-row duplication (§3.2),
`reflect.DeepEqual` per frame (§4.1), per-frame allocations (§4.2–4.3), pixel loops (§4.4),
GUI-boundary interface (§2.5), mapper/provider coupling (§2.6), committed `output/`
artifacts (§7.6), CI gaps (§7).

---

## 1. Bugs & correctness

### 1.1 ✅ FIXED 🔴 Road mutations silently lost in `TopologyBase` (value-copy bug)

[topologyBase.go](../internal/services/template_generator/providers/topology/base/topologyBase.go#L323-L331)
(`CreateMissingPlayerConnections`) and
[topologyBase.go](../internal/services/template_generator/providers/topology/base/topologyBase.go#L415-L438)
(`CreateMissingConnections`) both do:

```go
if pz, ok := linq.FromSlice(zones).First(func(z entities.Zone) bool { ... }); ok {
    pz.Roads = append(pz.Roads, ...) // pz is a VALUE COPY — append is discarded
}
```

`linq.First` returns a copy of the `entities.Zone` struct. Appending to `pz.Roads` /
`zone.Roads` mutates the copy; the road never reaches the `zones` slice the caller keeps.
Fallback and bridge connections are therefore emitted **without their road entries**, so
generated maps can ship bridge/fallback connections with no road linking them.

**Fix** — locate by index and mutate in place:

```go
zoneIndex := slices.IndexFunc(zones, func(z entities.Zone) bool { return z.Name == zoneName })
if zoneIndex >= 0 {
    zones[zoneIndex].Roads = append(zones[zoneIndex].Roads, roadBuilder...Build())
}
```

Apply the same pattern to all three mutation sites (one in
`CreateMissingPlayerConnections`, two branches in `CreateMissingConnections`). Then add
unit tests under
`test/unit/internal/services/template_generator/providers/topology/base/topologyBase/`
asserting that after `CreateMissingConnections` returns, the affected zones in the passed
slice contain the new bridge road (this is exactly the untestable side effect
`test_observations.md` item 17 flagged — fixing the bug makes it testable).

> If investigation shows roads are *intentionally* not added (game tolerates it), delete
> the dead mutation code instead — either way the current code is wrong.

### 1.2 ✅ FIXED 🟠 `State.Exit` still calls `os.Exit(0)`; `confirmExit` never resets

[state.go](../app/gui/drivers/state.go#L162-L173):

```go
func (this *State) Exit() {
    if this.unsaved && !this.confirmExit {
        this.SetStatus("Unsaved changes exist - save first or press Exit again.", true)
        this.confirmExit = true // TODO: reset if not selected right after
        return
    }
    // this.onExit()
    os.Exit(0)
}
```

Two defects:
1. `os.Exit(0)` skips deferred cleanup and kills the Gio event loop abruptly. The
   commented-out `this.onExit()` shows the intended design was never wired.
2. `confirmExit` is never reset, so after declining once, **any** later Exit press quits
   without warning even if new unsaved changes accumulated.

**Fix**:
1. Add `onExit func()` to `State`; in [program.go](../app/gui/program.go) pass
   `func() { window.Perform(system.ActionClose) }` when constructing the state, so exit
   flows through the normal `app.DestroyEvent` path.
2. Reset the flag on any state mutation: in `UpdateState`, when
   `this.innerState.WasStateChanged()`, set `this.confirmExit = false`. Also reset in
   `Save`/`SaveAs` success paths.
3. Delete the commented line and the `os` import if it becomes unused.

### 1.3 ✅ FIXED 🟠 `generalPanel` gameMode not restored on load (author-confirmed bug)

[generalPanel.go](../app/gui/panels/generalPanel.go#L117):

```go
this.gameMode.SetSelectedIndex(0) // TODO: here is a bug where gameMode will not be loaded
```

`LoadFromState` hard-resets the game-mode dropdown to index 0 instead of selecting the
mode stored in the loaded `.gen.json`, silently discarding the user's saved choice.

**Fix** — map the state value to its dropdown index, mirroring how other dropdowns load:

```go
gameModeIndex := slices.Index(gameModeIds, state.GameMode) // build gameModeIds next to the dropdown options
if gameModeIndex < 0 {
    gameModeIndex = 0
}
this.gameMode.SetSelectedIndex(gameModeIndex)
```

Add an integration-suite check (load fixture with non-default gameMode → panel reflects
it), since panels are Gio-gated per AGENTS.md §4.6.

### 1.4 ✅ FIXED 🟠 Potential infinite loop in `CreateMissingConnections` on duplicate bridge name

[topologyBase.go](../internal/services/template_generator/providers/topology/base/topologyBase.go#L395-L399):

```go
bridgeName := fmt.Sprintf("Bridge-%s-%s", labelA, labelB)
if connNameSet[bridgeName] {
    continue // components NOT linked → same pair selected next iteration → forever
}
```

If a pre-existing connection already uses the generated bridge name (e.g. a manually
edited template reprocessed through generation), `continue` skips `adjacency.Link`, the
component structure never changes, and `GetShortestDistanceIndex` returns the same pair
every iteration — a hang, not a crash. `test_observations.md` item 3 called this out as
"review loop exit conditions"; it is still unresolved.

**Fix** — link the components before continuing, since the connection already exists:

```go
if connNameSet[bridgeName] {
    adjacency.Link(bestIndexes.X, bestIndexes.Y)
    continue
}
```

This is also semantically correct: an existing bridge *does* connect the components.

### 1.5 ✅ FIXED 🟡 Unknown victory condition silently coerced to first entry

[victoryConditions.go](../app/gui/constants/victoryConditions.go#L50):

```go
return GetVictoryConditionList()[0] // TODO: probably should return empty Victory... suck it
```

A corrupt/newer `.gen.json` with an unrecognized victory condition is silently reshaped
into the first list entry; the user never learns their setting was dropped.

**Fix** — return the comma-ok pair and let the caller surface a status message:

```go
func GetVictoryConditionByID(id string) (VictoryCondition, bool) { ... }
```

At the load call-site, on `!ok` fall back to Standard **and** call
`SetStatus("Unknown victory condition %q in file - reset to Standard.", false)`.

### 1.6 🟡 `GUIHandler.UpdateTemplate` shallow-copies the template but aliases `Variants`

[guiHandler.go](../internal/handlers/guiHandler.go#L58-L70): `newTemplate :=
*templateDto.Template` copies the struct, but `Variants` is a slice, so
`newTemplate.Variants[0].Zones = templateDto.Zones` **also mutates the caller's template**
through the shared backing array. Today the only caller immediately replaces
`lastTemplate` with the returned pointer, so nothing breaks — but the contract is a trap.

**Fix** — either document the in-place contract on the method, or make the copy real:

```go
newTemplate := *templateDto.Template
newTemplate.Variants = slices.Clone(templateDto.Template.Variants)
```

`slices.Clone` of one variant header is cheap and removes the aliasing hazard.

### 1.7 🟡 `resolveGlob` shape and Steam-path fallback

[io.go](../internal/helpers/io.go#L146-L159): `resolveGlob` returns `("", nil)` when the
glob matches nothing and `("", err)` at the bottom where `err` is whatever the *last*
`os.Stat` returned — callers cannot distinguish "not installed" from "I/O failure", and
[state.go](../app/gui/drivers/state.go#L65-L74) has to guess with a two-branch status
message. Also [io.go](../internal/helpers/io.go#L96) keeps a literal
`C:\Program Files (x86)\Steam` fallback; on non-standard installs the correct source is
the registry key `HKCU\Software\Valve\Steam\SteamPath`.

**Fix**:
1. Add `ErrTemplatesDirNotFound` to [editorErrors.go](../internal/common/editorErrors.go);
   return it from `resolveGlob` on zero usable matches; wrap real I/O errors with
   `fmt.Errorf("resolve templates glob %q: %w", pattern, err)`.
2. Optional (Windows polish): read the registry key via `golang.org/x/sys/windows/registry`
   behind `//go:build windows` before falling back to the literal path.

### 1.8 🟡 `fileExplorerDialog.handleConfirm` — `appendAssign` and complexity 43

[fileExplorerDialog.go](../app/gui/dialogs/fileExplorerDialog.go#L535):
`this.entries = append(dirs, files...)` appends onto `dirs`'s backing array (gocritic
`appendAssign`) — if `dirs` is ever retained elsewhere this aliases; and `handleConfirm`
has the highest cognitive complexity in the codebase (43, limit 30), with a nested
`overwriteActive` block of nestif complexity 15.

**Fix**: `this.entries = slices.Concat(dirs, files)`; then extract
`confirmOverwrite(gtx)` and `confirmSelection(gtx)` from `handleConfirm` so each branch is
a flat early-return function.

### 1.9 🟡 gosec G115: unchecked `int → uint8` conversions in alpha blending

[assetProvider.go](../internal/services/asset_provider/assetProvider.go#L101-L103). The math
(`(a*alpha + b*keep)/255` with `alpha+keep == 255`) cannot exceed 255, so this is a false
positive — but it should be *documented as such*, not ignored silently.

**Fix**: add `//nolint:gosec // G115: (x*alpha + y*keep)/255 ≤ 255 by construction` on the
three lines, or clamp with `min(value, 255)` if you prefer belt-and-braces.

---

## 2. Architecture

### 2.1 🟠 Validation still does not exist anywhere

Three artifacts point at the same hole:
- [generatorConfigValidator.go](../internal/validators/generatorConfigValidator.go) — the
  whole `internal/validators` package is one TODO comment.
- [editorState.go](../app/gui/models/editorState.go#L34) — `// TODO: add validator for state
  updates, e.g. to prevent invalid map sizes or player counts`.
- [state.go](../app/gui/drivers/state.go#L204) — commented-out
  `validators.ValidateEditorState(...)` call in `UpdateState`.

Out-of-range values loaded from hand-edited `.gen.json` files flow straight into
generation.

**Fix** (single PR):
1. Implement `func ValidateEditorState(state *dtos.EditorStateDto) []string` in the
   validators package: clamp/report `PlayerCount` (registry bounds), `MapSize`,
   `HeroCountMin ≤ HeroCountMax`, percent fields to 25–200, zone counts ≥ 0, known
   `GameMode`/`VictoryCondition` IDs (registry lookups).
2. Call it in `GUIHandler.LoadState` after `LoadSettingsFile` and surface joined messages
   via the existing status mechanism; clamp rather than reject so old files still load.
3. Unit tests per AGENTS.md §4.6 layout: `test/unit/internal/validators/generatorConfigValidator/validateEditorState_test.go`
   with gofakeit-fuzzed out-of-range values.
4. Delete the empty-package TODO and the commented call.

### 2.2 🟠 `drivers.State` remains a god object without a seam (408 LOC)

[state.go](../app/gui/drivers/state.go) still owns: file dialogs, save/load, generation,
debounce state machine, manual-edit reapplication, exit policy, status bar, output-path
widget. It holds `*handlers.GUIHandler` and `*mappers.GeneratorConfigMapper` concretely, so
every State-level test would drag in the full generator + preview stack (this is why
`test_observations.md` skipped unit tests for it entirely).

**Fix** (two independent steps):
1. **Interface at the boundary** — define in `app/gui/drivers/templateHandler.go`:

   ```go
   type TemplateHandler interface {
       GenerateTemplate(dtos.EditorStateDto) (dtos.TemplateLoadDto, error)
       UpdateTemplate(dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error)
       SaveTemplate(dtos.TemplateSaveDto) (string, error)
       LoadState(string) (*dtos.EditorStateDto, error)
       SaveState(dtos.EditorStateSaveDto) (string, error)
   }
   ```

   `State.handler` becomes `TemplateHandler`; `NewUIState` keeps the real impl. The
   debounce machine in `AutoRegenerate` (pure time arithmetic, already parameterized on
   `now`) then becomes unit-testable with a testify-mock handler and a fake clock.
2. **File split** by concern, same package (per AGENTS.md one-struct-per-file this is a
   method split, not a struct split): `state.go` (core + status), `stateFiles.go`
   (Load/Save/SaveAs/Exit/suggestDirectory), `stateGeneration.go`
   (Generate/AutoRegenerate/handleGenerateTemplate/applyGeneratedTemplate),
   `stateManualEdits.go` (ApplyEditedZones/reapplyManualEdits).

### 2.3 🟠 `zoneEditorDialog.go` — 1,148 LOC god file

[zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go) mixes canvas rendering,
pointer/hit-testing, geometry recomputation, snapping, and two property-panel forms.
`recomputeGeometry` alone trips funlen (67 > 60).

**Fix** — split into four files in `app/gui/dialogs/` (methods stay on the one struct):
`zoneEditorDialog.go` (struct, ctor, top-level Layout, ~300), `zoneEditorCanvas.go`
(draw* + recomputeGeometry + handlePointer, ~330), `zoneEditorConnectionProps.go` (~230),
`zoneEditorZoneProps.go` (~220). While moving `recomputeGeometry`, extract the
group-classification loop into `classifyZoneGroups()` to clear the funlen finding
honestly (no `//nolint`).

### 2.4 🟡 `layoutPanel.go` (547 LOC) and `previewPanel.getPreviewCanvasWidget` (71-line func)

[layoutPanel.go](../app/gui/panels/layoutPanel.go) declares 40+ widget fields and builds
three columns in one file; [previewPanel.go](../app/gui/panels/previewPanel.go#L148) trips
funlen. **Fix**: split layoutPanel by column (`layoutPanel.go`, `layoutPanelTopology.go`,
`layoutPanelZones.go`) and extract the canvas-image fitting math from
`getPreviewCanvasWidget` into a private `fitPreviewImage(...)` helper. Also resolve the
commented-out `MinNeutralZonesBetweenPlayers` row at
[layoutPanel.go](../app/gui/panels/layoutPanel.go#L253): the DTO field is live and mapped
into config — either re-enable the row or delete the comment and document why the option
is hidden (godox flags it).

### 2.5 🟠 `CreateTopologyVariant` — 29 lines duplicated four times (dupl ×4)

The linter confirms [crossTopology.go](../internal/services/template_generator/providers/topology/crossTopology.go#L27),
[fractalTopology.go](../internal/services/template_generator/providers/topology/fractalTopology.go#L31),
[geometricTopology.go](../internal/services/template_generator/providers/topology/geometricTopology.go#L29) and
[squareTopology.go](../internal/services/template_generator/providers/topology/squareTopology.go#L25)
share an identical 29-line body; the only variance is the `create*Layout` call.

**Fix** — template-method on `TopologyBase`:

```go
type layoutFunc func(playerLabels []string, neutralZones neutralZone.Plans) (
    allLabels []string, positions models.Positions, pairs []models.ConnectionIndexes)

func (this *TopologyBase) CreateVariantFromLayout(
    configuration config.GeneratorConfig, playerLabels []string,
    neutralZones neutralZone.Plans, tuning models.GenerationTuning,
    holdCityNeutralLabel string, buildLayout layoutFunc) entities.Variant {
    // the shared 29 lines, calling buildLayout(...)
}
```

Each service's `CreateTopologyVariant` collapses to one call passing its private layout
func. Clears all four dupl findings, and adding topology #11 no longer copies the block.

### 2.6 🟡 `GeneratorConfigMapper.FromEditorState` allocates a provider per call

[generatorConfigMapper.go](../internal/mappers/generatorConfigMapper.go#L23) constructs
`providers.NewMandatoryContentProvider()` on every mapping, and the mapper package thereby
depends on the deepest service package. `GUIHandler` already owns a
`contentProvider` instance.

**Fix**: store the provider on the mapper —

```go
type GeneratorConfigMapper struct{ contentProvider *providers.MandatoryContentProvider }
func NewConfigMapper() *GeneratorConfigMapper {
    return &GeneratorConfigMapper{contentProvider: providers.NewMandatoryContentProvider()}
}
```

(If you later want mappers fully decoupled from services, move
`CreateContentItemsFrom` into `internal/models` — but the field injection alone removes
the per-call allocation now.)

### 2.7 🟡 Unit tests for internal services import `app/gui/constants`

[getVariantsForContent_test.go](../test/unit/internal/services/content_rules/variantMappingManager/getVariantsForContent_test.go#L6),
`getVariantForContentById_test.go`, `createRuleFromSavedRule_test.go` — tests of
**internal** packages reach up into the GUI layer for SID constants. depguard covers
production code only, so the boundary quietly leaks back through tests.

**Fix**: the SIDs these tests need are game data — reference them from
`internal/registry` getters (where the values originate) or inline the literal SIDs in the
tests. Then extend the depguard rule to `test/unit/internal/**` denying `app/**`.

---

## 3. Duplicate code

### 3.1 🟠 Topology `CreateTopologyVariant` ×4 — see §2.5 (dupl-confirmed, one fix).

### 3.2 🟡 Labeled slider/checkbox row vocabulary across panels

[layoutPanel.go](../app/gui/panels/layoutPanel.go), [generalPanel.go](../app/gui/panels/generalPanel.go),
[bonusesPanel.go](../app/gui/panels/bonusesPanel.go) repeat
`widgets.NewLabeledRowWidget(theme, "Label", width, widgets.NewLabeledSliderWidget(...))`
~50 times. **Fix**: add two factories to [app/gui/widgets/](../app/gui/widgets/):
`NewSliderRow(theme, label string, width unit.Dp, slider *widget.Float, format func(float32) string)`
and `NewCheckboxRow(theme, value *widget.Bool, label, hint string)`. Panels then read like
a form schema; removes ~120 LOC. (Keep it to these two — don't build a form framework.)

### 3.3 🟡 `circlePoint` centre parameters always 0.5 (unparam ×2)

[geometryHelpers.go](../internal/services/template_generator/providers/topology/geometryHelpers.go#L12):
every caller passes `centreX=0.5, centreY=0.5`. **Fix**: drop both parameters and use a
named constant `const layoutCentre = 0.5` inside; or keep them and add one caller that
actually varies them. Prefer the former until a real need appears.

### 3.4 🟡 Duplicated registry-lookup globals

`winConditions`/`gameModes` package globals are declared in
[editorStateDto.go](../internal/dtos/editorStateDto.go#L12-L15),
[generatorConfigMapper.go](../internal/mappers/generatorConfigMapper.go#L10-L12) and
[common.go](../internal/services/template_generator/providers/common.go#L6-L22). Registry
getters return by value on every call. **Fix**: these are cheap immutable lookups — call
`registry.GetWinningConditionValues()` inline at use sites (drops the globals and clears
several `gochecknoglobals` findings), or keep one canonical alias set in
`internal/registry` accessor form. Do **not** edit the registry files themselves (§2.1
read-only rule); the globals live outside the registry and are fair game.

---

## 4. Performance

### 4.1 🟠 `reflect.DeepEqual` over the whole DTO every frame

[editorStateDto.go](../internal/dtos/editorStateDto.go#L179-L186)
`EqualsIgnoringManualEdits` copies two ~100-field structs and reflection-walks five
`[]ZoneContentRowSave` slices; [state.go](../app/gui/drivers/state.go#L233)
`AutoRegenerate` calls it (via `ResetNextStateIfStateWasNotChanged`) **every frame**, even
idle.

**Fix** (cheapest first):
1. Short-circuit: compare the scalar-only sections first (the struct already has
   `LayoutDefiningOptionsChanged` + `zoneCountOptionsChanged` precedents); only
   deep-compare the row slices when scalars match:

   ```go
   func (this *EditorStateDto) EqualsIgnoringManualEdits(other *EditorStateDto) bool {
       left, right := *this, *other
       left.ManualZones, left.ManualConnections = nil, nil
       right.ManualZones, right.ManualConnections = nil, nil
       // nil out the row slices, compare scalars with ==, then compare rows with slices.Equal
       ...
   }
   ```

2. Better: add a `revision uint64` to `models.EditorState`, incremented inside
   `UpdateCurrentState` only when the update func actually changed something (the
   equality check runs once per mutation instead of once per frame);
   `AutoRegenerate` compares `revision != lastGeneratedRevision`. This deletes the
   per-frame comparison entirely.

Add a unit test asserting the hand-rolled equality matches `reflect.DeepEqual` on
gofakeit-fuzzed pairs so the two can't drift.

### 4.2 🟡 Per-frame slice allocation for static tabs

[window.go](../app/gui/editor/window.go#L65-L75) `getTabsWidget` rebuilds
`make([]layout.FlexChild, 0)` from the same four static tabs every frame. **Fix**: build
the `[]layout.FlexChild` once in `NewWindow` (tab set never changes) or reuse a struct
field slice with `this.tabChildren = this.tabChildren[:0]`.

### 4.3 🟡 Zone editor recomputes full geometry every frame while open

[zoneEditorDialog.go](../app/gui/dialogs/zoneEditorDialog.go#L381) calls
`recomputeGeometry(side)` unconditionally per frame; that runs `BuildPreviewLayout`
(up to 500 relaxation passes for scatter layouts) plus O(edges×zones) grouping — while the
dialog just sits there. **Fix**: add `geomDirty bool` set by every mutator
(`handlePointer` press/drag/release paths, add/delete zone, property edits) and recompute
only when `geomDirty || side != this.lastSide`. The comment at L375-379 already documents
that hit-testing intentionally uses last frame's geometry, so a dirty flag is consistent
with the existing design.

### 4.4 🟡 Pixel-by-pixel brush fill in preview connector drawing

[previewGeneratorService.go](../internal/services/preview_service/previewGeneratorService.go#L94-L104)
fills the brush rect with nested `canvas.SetRGBA` loops per curve step (~15k calls per
connector). Not user-visible today (runs per generation, not per frame), but trivial to
fix: replace the inner loops with
`draw.Draw(canvas, brushRect, image.NewUniform(connectorLineColor), image.Point{}, draw.Src)`
from `image/draw`.

---

## 5. Readability & maintainability

### 5.1 🟡 Ten funlen + five gocognit functions — decompose, don't suppress

Beyond those covered above (§1.8, §2.3, §2.4), the linter flags:

| Function | Finding | Suggested split |
|---|---|---|
| [topologyBase.go](../internal/services/template_generator/providers/topology/base/topologyBase.go#L336) `CreateMissingConnections` (104 lines, cognit 39) | funlen+gocognit | Extract `buildZoneAdjacency(...)` (the block already marked `// TODO: move out to a separate function` at L363) and `appendBridgeRoads(...)` (the §1.1 fix will rewrite this anyway — do both in one PR). |
| [topologyBase.go](../internal/services/template_generator/providers/topology/base/topologyBase.go#L266) `CreateMissingPlayerConnections` (67) | funlen | Extract `spawnZoneHasConnection(zone, connNames)` and the fallback-connection builder. |
| [balancedClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/balancedClusterService.go#L112) `createSortedPairs` (51 stmts, cognit 41) | funlen+gocognit | Split scoring/sorting from pair emission. |
| [fractalTopology.go](../internal/services/template_generator/providers/topology/fractalTopology.go#L75) `createFractalLayout` (68) | funlen | Extract tier-bucketing and per-player distribution loops. |
| [geometricTopology.go](../internal/services/template_generator/providers/topology/geometricTopology.go#L75) `createGeometricLayout` (53 stmts) | funlen | Same treatment. |
| [editorStateDto.go](../internal/dtos/editorStateDto.go#L208) `DefaultPlayerZoneContentRows` (113) | funlen | Pure data — split into per-tier `defaultLowRows()`, `defaultMediumRows()`… helpers, or add a scoped `//nolint:funlen // literal table` since it has zero branching logic. |
| [position.go](../internal/models/position.go#L89) `CreateDelaunayTriangulation` (64) | funlen | Extract `circumcircleContains` / bad-triangle collection; add the missing Bowyer-Watson doc comment while there (long-standing gap). |
| [pickerDialog.go](../app/gui/dialogs/pickerDialog.go#L286) `getLeafRowWidget` (62) | funlen | Extract the trailing-control builder. |
| [buildNonAdjacentDerangement](../internal/services/template_generator/providers/topology/base/topologyBase.go#L741) (cognit 36) | gocognit | Deterministic-fallback safety net; extract the fallback into `buildDeterministicDerangement(count)` — also makes the "practically unreachable" branch (test_observations #2) directly testable. |
| [generateAllTopologies_test.go](../test/unit/internal/services/template_generator/templateGenerator/generateAllTopologies_test.go#L57) (cognit 33) | gocognit | Convert inner assertions into table-driven `t.Run` subtests per zone-field. |

### 5.2 🟡 goconst — magic rule/tier strings (7 findings)

`"Guarded"` ×14, `"Distance to town"`, `"Distance to road"`, `"Near"`, `"Next To"`,
`"Connection"`, `"Generator Default"`. These are **serialized rule identifiers** — a typo
compiles fine and silently breaks rule matching. **Fix**: define once next to their
owners: rule names in `internal/services/content_rules` (e.g.
`const RuleNameGuarded = "Guarded"`), distance names alongside
[distancePresets.go](../internal/services/content_rules/distancePresets.go), ref-type
`"Connection"` in the builders package — then use them in
[editorStateDto.go](../internal/dtos/editorStateDto.go#L223)'s default tables and
[topologyBase.go](../internal/services/template_generator/providers/topology/base/topologyBase.go#L298).

### 5.3 🟡 revive unexported-return (6 findings)

`NewItemPickerDialog`/`NewSpellPickerDialog`/`NewValueOverridePickerDialog` →
`*multiSelectPicker`; `GetVictoryConditionValues` → `victoryConditions`; `NewPairSet` →
`*pairSet`; `CreateHubZoneCandidates` → `*hubZoneCandidates`. Callers outside the package
cannot name these types (can't store them in struct fields or pass them around).
**Fix**: export the types (`MultiSelectPicker`, `PairSet`, `HubZoneCandidates`,
`VictoryConditions`) — mechanical rename, no behavior change. Also `selectedIds` →
`selectedIDs` (var-naming).

### 5.4 🟡 sloglint — default logger (3 findings)

[program.go](../app/gui/program.go#L32) / [guiHandler.go](../internal/handlers/guiHandler.go#L28)
use `slog.Error(...)` on the default logger. **Fix**: construct one
`logger := slog.New(slog.NewTextHandler(os.Stderr, nil))` in `main.go`/`StartApplication`,
pass it to `NewGuiHandler(logger)`, and call `logger.Error(...)`. This also unlocks
capturing log output in tests. If that plumbing feels heavy for three call sites, the
honest alternative is disabling sloglint's `no-global` option in `.golangci.yml` —
pick one, don't leave the warnings.

### 5.5 🟡 godox inventory — 17 TODOs; resolve or convert

Actionable ones are already covered: victory condition (§1.5), exit confirm (§1.2),
gameMode (§1.3), validators (§2.1), `move out to a separate function` (§5.1),
MinNeutralZonesBetweenPlayers row (§2.4). The remainder should be resolved cheaply or
deleted:
- [toolbar.go](../app/gui/editor/toolbar.go#L32) icon experiment comment → delete (kept "just in case" since last review).
- [zone.go](../internal/entities/template/template_variant/zone.go#L12) `make it vec2[float64]` — **read-only schema file**; the TODO itself violates the no-touch rule's spirit. Propose removal to the user; do not edit without approval.
- [previewLayout.go](../internal/models/preview/previewLayout.go#L7) Vec2 sub-pixel TODO — legitimate future work; move to an issue tracker entry and delete the comment.
- [zoneAdjacency.go](../internal/models/zoneAdjacency.go#L55), [zoneLabelProvider.go](../internal/services/zones/zoneLabelProvider.go#L41-L42), [mandatoryContentBuilder.go](../internal/services/builders/mandatory_content/mandatoryContentBuilder.go#L43), [types.go](../internal/services/builders/types.go#L3) — each is answerable in <30 min: answer it, apply or delete.

### 5.6 🟡 `getSteamPath` misc

[io.go](../internal/helpers/io.go#L40) contains a commented-out redundant condition
(`/*&& runtime.GOOS != "windows" is redundant here*/`) — fold the explanation into the
function doc comment instead of inline commented code.

---

## 6. Testing

### 6.1 🟠 `drivers.State` untested — unblock via the §2.2 interface

`test_observations.md` skipped `State` because `NewUIState` probes the host Steam install
and the concrete handler drags the whole stack. The §2.2 `TemplateHandler` interface plus
a `NewUIStateWithHandler(handler TemplateHandler, templateDir string)` constructor (a real
constructor, not a test seam — `NewUIState` calls it) makes the debounce machine, exit
guard (§1.2), and status transitions unit-testable with a fake clock. Add:
`test/unit/app/gui/drivers/state/autoRegenerate_test.go` (debounce arm/rearm/fire matrix),
`exit_test.go` (guard + reset), `updateState_test.go` (unsaved/confirmExit flags).

### 6.2 🟡 `internal/helpers/io.go` still untestable

Still zero unit tests (carried from test_observations #1). The blocker is host-state
probing inside one entry point. **Fix**: extract the pure part —
`resolveTemplatesDir(homeDir string, statFn func(string) (os.FileInfo, error), globFn func(string) ([]string, error))` — 
or simpler: accept a base-dir override parameter used by `FindOldenEraTemplatesDir`
internally, then test `getBasePath` (pure map traversal — trivially testable today, just
needs a test folder), `resolveGlob` via `t.TempDir()`, and the VDF path assembly.
`getBasePath` at minimum should be covered now; no code change needed for it beyond the
§1.7 error-shape fix.

### 6.3 ✅ FIXED 🟡 Road-mutation regression tests

After fixing §1.1, add assertions that zones passed to
`CreateMissingConnections`/`CreateMissingPlayerConnections` gain the bridge/fallback
roads. After fixing §1.4, add a test seeding a pre-existing `Bridge-A-B` connection and
asserting termination (guard with `t.Deadline`-aware context or just trust the fix — the
loop now provably progresses).

### 6.4 ⚪ Coverage posture

CI enforces a 60% floor with decrease-fail — good. The remaining structural gaps are the
ones above (State, io.go) plus Gio-gated files already routed to the integration suite.
No new blanket coverage push is warranted; target the two gaps.

---

## 7. CI/CD

### 7.1 🟠 Vulnerability scan is wired but disabled

[pr-validation.yml](../.github/workflows/pr-validation.yml#L61-L64): the
`run-vulnerability-scan` job is complete but has `if: false`. govulncheck is low-noise
(reports only reachable vulns). **Fix**: change to
`if: ${{ github.event_name == 'pull_request' }}`, and collapse the three per-tree steps
into one `go-package: ./...` step (they share the module; three runs triple the cost for
no added signal).

### 7.2 🟠 No race detector anywhere in CI

The GUI is event-driven with a generator pipeline behind it; races would ship silently.
**Fix**: add `-race` to the unit-test step (Linux runners support it out of the box):

```yaml
- name: Run unit tests
  run: go test -race -coverprofile=... -coverpkg=./internal/...,./app/... ./test/unit/...
```

If `-race` + coverage is too slow, run a separate `go test -race ./test/unit/...` job
without coverage in parallel.

### 7.3 🟡 `check-build` only builds the root package

[pr-validation.yml](../.github/workflows/pr-validation.yml#L110): `go build .` skips
`tools/` and any package not reachable from main (e.g. `*_testexports.go` files under the
`integration_test` tag). **Fix**: `go build ./...` plus
`go vet -tags=integration_test ./...` to compile-check the gated files without running
gated tests.

### 7.4 🟡 Single-OS CI for a dual-OS product

Everything runs on ubuntu-latest, yet the primary user platform (and all Steam-path logic
in [io.go](../internal/helpers/io.go)) is Windows. **Fix**: add a `windows-latest` job
running `go build ./...` and `go test ./test/unit/...` (skip lint/coverage there; guard
the Gio apt-get deps in [setup-steps/action.yml](../.github/workflows/setup-steps/action.yml)
with `if: runner.os == 'Linux'`).

### 7.5 🟡 Release workflow hardening

[release.yml](../.github/workflows/release.yml):
1. Add `-trimpath` and version injection:
   `go build -trimpath -ldflags "-s -w -X main.version=${{ github.ref_name }}"` (add a
   `var version = "dev"` in `main.go`).
2. Publish checksums: after downloading artifacts, `sha256sum dist/* > dist/checksums.txt`
   and add it to `files:`.
3. Pin third-party actions to commit SHAs (`softprops/action-gh-release@v2` at minimum;
   official `actions/*` may stay on tags per your risk appetite).
4. Add `concurrency: { group: release-${{ github.ref }}, cancel-in-progress: false }`.
5. `workflow_dispatch` builds from the branch HEAD but labels the release with the input
   tag — checkout the tag explicitly when `github.event.inputs.tag` is set:
   `with: { ref: ${{ github.event.inputs.tag || github.ref }} }`.

### 7.6 🟡 Committed generated artifacts

`git ls-files` shows [output/Colosseum/](../output/Colosseum/) (`.gen.json`, `.png`,
`.rmg.json`) is **tracked**, while `.gitignore` ignores `output/` — the files predate the
ignore rule and every generation into that folder shows up as tracked-file modifications.
**Fix**: if Colosseum is a wanted example, move it to `data/ExampleSettings/` (⚠ requires
your approval — `data/` is a protected directory) or a new `examples/` folder; otherwise
`git rm -r --cached output/` (needs your confirmation — destructive to git history view,
not to files).

### 7.7 ⚪ Dependabot

No `.github/dependabot.yml`. One 10-line file covering `gomod` + `github-actions`
(weekly) keeps Gio and action versions moving without manual sweeps.

---

## 8. Linter findings — disposition of all 113 issues

| Linter | Count | Disposition |
|---|---|---|
| `gochecknoglobals` | 50 | Mostly **intentional immutable catalogs** (registry aliases, variant mappings, distance presets, guard presets). Fix the duplicated aliases per §3.4 (~20 findings). For the rest: these are effectively constants that Go can't declare `const`; add `gochecknoglobals` exclusion paths in `.golangci.yml` for `internal/registry/`, `internal/services/content_rules/distancePresets.go`, `internal/services/connection_editor/` label tables — or convert the truly private ones to function-local lookups. Do **not** edit `internal/registry` files (read-only). |
| `godox` | 17 | Every TODO dispositioned in §5.5 (and §1.2/§1.3/§1.5/§2.1/§2.4). Target: zero. |
| `funlen` | 10 | Decomposition table in §5.1 (+§2.3, §2.4). One acceptable `//nolint` candidate: the pure data table `DefaultPlayerZoneContentRows`. |
| `goconst` | 7 | Serialized identifiers → named constants per §5.2. |
| `revive` | 7 | Export the returned types + `selectedIDs` rename per §5.3. |
| `gocognit` | 5 | Same functions as funlen; §5.1 + §1.8. |
| `dupl` | 4 | Topology template-method per §2.5. |
| `gosec` (G115) | 3 | False positive by construction; documented nolint per §1.9. |
| `sloglint` | 3 | Injected logger or explicit config opt-out per §5.4. |
| `nestif` | 3 | `handleConfirm` split (§1.8) covers 2; `topologyBase` nestif falls out of the §1.1+§5.1 rewrite. |
| `gocritic` | 2 | `appendAssign` (§1.8) and if-else→switch in topologyBase (fold into §5.1 rewrite). |
| `unparam` | 2 | `circlePoint` centre params per §3.3. |

---

## 9. Suggested execution order

1. **Bugs first, isolated PRs**: §1.1 + §1.4 + §6.3 (topology roads, one PR) → §1.3
   (gameMode) → §1.2 (exit path) → §1.5, §1.6, §1.7.
2. **Validation**: §2.1 (new package, tests included).
3. **Structure**: §2.5/§3.1 (topology template method, clears dupl) → §2.2 + §6.1 (State
   interface + tests) → §2.3, §2.4 (file splits) → §5.1 decompositions.
4. **Performance**: §4.1 (revision counter) → §4.3 (geometry dirty flag) → §4.2, §4.4.
5. **Sweeps**: §5.2 constants, §5.3 exports, §5.4 logging, §5.5 TODO purge, §3.2–3.4.
6. **CI**: §7.1, §7.2, §7.3 in one PR; §7.4, §7.5, §7.7 as follow-ups; §7.6 needs owner
   decision (protected `data/` folder or `git rm --cached`).

Every code change above must follow AGENTS.md: tests per §4.6 layout, coverage checked
before/after (§2.3), no edits under `data/`, `internal/entities/template/`,
`internal/registry/` without explicit approval, cross-platform paths only.
