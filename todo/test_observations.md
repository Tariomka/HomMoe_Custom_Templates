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
