# Backlog

Small future-work items moved out of code comments (godox purge, review §5.5).

- **Preview sub-pixel precision**: `internal/models/preview/previewLayout.go` —
  `Layout.Positions` uses `map[string]image.Point`; switching to a `Vec2`
  (float64) type would give sub-pixel precision in preview rendering and zone
  editor geometry. Related: `Zone.GeneratorPosition *[2]float64` in the
  read-only schema would ideally share that Vec2 type.

- remove the [2]float from the template entities for example entities/template/template_variant/zone.go. It's easier to use Vec2 instead

- **`createTopologyAdjacency` dead Chain/Ring branches**:
  `internal/services/zones/zoneLabelProvider.go` — the `case TopologyChain` and
  `case TopologyRing, TopologyCircles` branches (plus the `isIsolated` guard
  they use) are unreachable: the only caller, `GetHoldCityLabel`, gates on
  `IsHubCityToHold()` (= `Topology == HubAndSpoke && IsCityHoldMode()`), so the
  switch always takes `default`. Verified 2026-07 (review §5.5): single private
  call site, single production caller (`templateGenerator.Generate`), branches
  never reachable in any commit (all three symbols born together in `bb50aab`),
  0% test coverage on them. A removal was implemented and verified
  (build/tests/lint green, coverage 64.1→64.2) but ROLLED BACK by the owner to
  keep the topology-aware adjacency as a starting point. Decide eventually:
  - either delete the branches (pure ratio win, zero behavior change), or
  - start using them: extend hold-city (or other adjacency-based features) to
    Chain/Ring/Circles topologies, which would also fix the `default` branch
    modelling Hub & Spoke as a sequential ring instead of its real star graph.

- need to add untracked zone tier property to Entities and/or method. This will be the source of truth for all tier related operations.
  As currently there is no reading of rmg.json files, deserialization doesn't matter so the entity can assume the tier from the content currently inside the zone.
  Or maybe Zone doesn't actually need a property, neutralZone.Quality can be used to track and infer this information as well as neutralZone.Profile can have this property.

- need to use common (either from commons or models) values for template generation

- either random portals or connections in general to hub have incorrect guard values(hub values are not being applied, might be because they are never calculated. hub(s) probably need to be inside the neutral zone list? maybe not, because neutral zone generation will not work as expected. need to think whats the cause and how to fix this issue)

- rework EditorStateDto, the dto content should be in entities, it should be wrapped in a model and model should be embedded to the dto.

- types inside entities should be moved to template, template package should be renamed to template_entity

- **"Save As" is really "Save To": make the UI say so.** Writing editor state as
  `{TemplateName}.gen.json` is *intended* behaviour (review §1.1, owner-approved:
  `FileService.SaveSettings` passes `filepath.Dir(path)` and
  `editorState.TemplateName` to the repository). The defect is only that the UI
  still offers a file **name** field the user can type into and whose value is
  then silently discarded — the dialog picks a *directory*, nothing more.
  Found while writing the Batch 13 GUI integration tests (2026-08). To fix:
  - [fileExplorerDialogToolbar.go](../app/gui/dialogs/fileExplorerDialogToolbar.go#L36-L51) —
    `getSaveRowWidget`: pass `readonly = true` to
    [`NewTextboxWidget`](../app/gui/widgets/textboxWidget.go#L18) and relabel
    `"Save as:"` → `"Will save as:"`, so the row reads as a preview of the
    resolved name instead of an input.
  - [toolbar.go](../app/gui/editor/toolbar.go#L71) — button `"Save As"` →
    `"Save To"`, field `buttonSaveAs` → `buttonSaveTo`.
  - [stateFiles.go](../app/gui/drivers/stateFiles.go#L35) — `State.SaveAs` →
    `State.SaveTo` (plus its `Save` fallback call site).
  - [fileExplorerDialogModes.go](../app/gui/dialogs/fileExplorerDialogModes.go#L29-L42) —
    dialog title `"Save File"` → `"Save To"`. Decide whether
    `NewSaveFileDialog` / `modeSaveFile` / `onSave` follow: they name the
    *explorer mode*, not the toolbar action, so they may legitimately stay.
  - Docs: [QUICKSTART.md](../QUICKSTART.md#L48) (also §L109) and
    [README.md](../README.md#L154) both advertise `Save As`.
  - Tests: rename `test/integration/stateSaveAs_integration_test.go` and
    `test/unit/app/gui/drivers/stateFiles/saveAs_test.go` with their subjects.
    In [fileExplorerDialog_integration_test.go](../test/integration/gui/fileExplorerDialog_integration_test.go)
    the scenarios that type a name (via the `SetFilename` testexport) must
    become "the field shows the resolved name and cannot be edited"; the
    disabled/enabled-confirm pair keyed on a whitespace filename needs a new
    trigger, since the user can no longer produce that state.
