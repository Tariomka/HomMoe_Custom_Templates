# Test Observations - untestable / integration-only code registry

Referenced by AGENTS.md §4.6. Records code that cannot be exercised through
public APIs in unit tests, so per-file coverage gaps here are intentional.

## Gio-UI-heavy code (integration-suite territory, no unit tests)

- app/gui/widgets/buttonWidget.go - all button constructors (and the private
  `addButtonSemantics` helper added 2026-07-12 for button-position debug
  logging) need a `layout.Context` + `material.Theme` text shaper to lay out;
  covered indirectly by the integration/performance suites that render the
  full editor window. The semantic-op replay path itself IS unit-tested via
  test/unit/app/gui/utils/buttonPositionLogger/ (headless ops that mirror
  `addButtonSemantics`), and was verified end-to-end against a real
  `NewButtonWidget` layout during development.

- app/gui/widgets/sliderRowWidget.go - `NewSliderRowWidget` (added 2026-07-13,
  review item §3.2) is a thin composition of `NewLabeledRowWidget` +
  `NewLabeledSliderWidget` whose returned closure needs a `layout.Context` +
  text shaper; covered indirectly by the integration/performance suites. The
  formatter funcs it receives ARE unit-tested (test/unit/app/gui/utils/string/
  *Formatter_test.go).

- app/gui/dialogs/fileExplorerDialog.go - `handleConfirm` / `confirmOverwrite` /
  `confirmSelection` need `layout.Context` + `widget.Clickable` click routing;
  the integration suite currently has NO file-explorer scenario (open/save/
  overwrite flows are unexercised). Noted 2026-07-12 during review item §1.8
  (behavior-preserving split of `handleConfirm`); synthetic-click coverage via
  the test/performance AppRunner pattern is possible future work.

- app/gui/dialogs/zoneEditorDialog.go + zoneEditorCanvas.go + zoneEditorSnap.go +
  zoneEditorConnectionProps.go + zoneEditorZoneProps.go - the Manual Zone Editor
  (one struct, method-split across five files in review item §2.3, 2026-07-12).
  Everything runs off `layout.Context` frames: canvas drawing/hit-testing uses
  the previous frame's geometry, property panels are `widget.Editor`/dropdown
  driven, and pointer flows (drag-to-connect, zone drag + snapping) need
  synthetic pointer events. The extracted `groupConnectionsByPair` and the snap
  helpers are private and only reachable through `Body`. Zone/connection
  *business* logic is unit-tested in internal/services/connection_editor;
  dialog interaction coverage would need the test/performance AppRunner
  synthetic-click pattern (future work — no integration scenario exists yet).

- app/gui/panels/layoutPanel.go + layoutPanelTopology.go + layoutPanelZones.go
  and previewPanel.go - Layout/Preview panels (layoutPanel method-split by
  column in review item §2.4, 2026-07-12; previewPanel's canvas closures
  promoted to private funcs/methods the same day). Pure Gio rendering:
  section/widget builders need `layout.Context` + text shaper, click handlers
  need `widget.Clickable` routing, and `LoadFromState`/`SaveToState` round-trips
  are exercised end-to-end by the integration suite (window save/load
  scenarios drive the tabs' SaveToState/LoadFromState closures). The state
  values they marshal are validated by the unit-tested
  internal/validators/editorStateValidator.

## app/gui/drivers.State (partially unit-tested since review item §2.2)

Unit tests use `NewUIStateWithHandler` + `test_helpers.TemplateHandlerMock`.
Still unit-untestable (dialog-callback or Gio territory):

- state.go - `NewUIState` (probes the disk for the game templates dir and
  builds the real `GUIHandler`/preview stack) and `GetOutputPathWidget`
  (returns a Gio widget). Covered by the integration suite.
- stateFiles.go - `handleSaveState` / `handleLoadState` success paths and
  `suggestDirectory` are only reachable through file-dialog callbacks
  (`Load`/`SaveAs` pick handlers); unit tests assert the dialogs open, the
  integration suite exercises the load/save flows via the
  `integration_test`-gated `SaveStateToFile`/`LoadStateFromFile` exports.
- stateFiles.go - `PickOutputDir` / `RevealOutputDir` only open dialogs whose
  behavior lives in the dialog implementations.
- stateGeneration.go - `reapplyManualEdits` castle-change branch requires a
  generation-then-castle-option-change sequence entangled with the real
  mapper; exercised by the integration suite's manual-edit scenarios.

