# Test Observations - untestable / integration-only code registry

Referenced by AGENTS.md §4.6. Records code that cannot be exercised through
public APIs in unit tests, so per-file coverage gaps here are intentional.

## Gio-UI-heavy code (integration-suite territory, no unit tests)

- app/gui/dialogs/fileExplorerDialog.go - `handleConfirm` / `confirmOverwrite` /
  `confirmSelection` need `layout.Context` + `widget.Clickable` click routing;
  the integration suite currently has NO file-explorer scenario (open/save/
  overwrite flows are unexercised). Noted 2026-07-12 during review item §1.8
  (behavior-preserving split of `handleConfirm`); synthetic-click coverage via
  the test/performance AppRunner pattern is possible future work.

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

