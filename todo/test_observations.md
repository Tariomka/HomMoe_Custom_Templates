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

- app/gui/dialogs/fileExplorerDialog.go and its `fileExplorerDialog*.go`
  siblings - `handleConfirm` / `confirmOverwrite` / `confirmSelection`,
  `confirmButtonState` and `tryCreateFolder` need `layout.Context` +
  `widget.Clickable` click routing, so they have no unit tests. Since 2026-08-07
  (review item §2.1/§2.5) they ARE covered end-to-end by
  test/integration/gui/fileExplorerDialog_integration_test.go, which drives real
  frames and queues clicks with `widget.Clickable.Click` through the
  `integration_test`-gated accessors in `fileExplorerDialog_testexports.go`:
  open-and-load, save target resolution, save through the real state driver,
  the overwrite prompt (gated write, cancel, confirm), new-folder creation, the
  existing-folder refusal and the disabled-confirm predicate. All filesystem
  *policy* (listing, filtering, hidden entries, roots, path resolution, reserved
  names) moved to internal/services/file_system and is unit-tested there; what
  is left in the dialog is rendering and click wiring only.
  Still uncovered: the hidden-file toggle and the pointer-driven row/scroll
  interactions (owner decision - excluded from the scenario set).

- app/gui/dialogs/zoneEditorDialog.go with its canvas, snap, property-panel,
  geometry, and transient-state sibling files - the Manual Zone Editor
  (one primary dialog struct with rendering methods and UI state split by
  responsibility).
  Everything runs off `layout.Context` frames: canvas drawing/hit-testing uses
  the previous frame's geometry, property panels are `widget.Editor`/dropdown
  driven, and pointer flows (drag-to-connect, zone drag + snapping) need
  synthetic pointer events. The extracted `groupConnectionsByPair` and the snap
  helpers are private and only reachable through `Body`. Zone/connection
  *business* logic is unit-tested in internal/services/connection_editor. The
  integration suite renders the real dialog and verifies handler-provided
  options; pointer interactions still need the test/performance AppRunner
  synthetic-click pattern as future work.

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

Unit tests use `NewUIState(handler, false)` + `test_helpers.TemplateHandlerMock`.
Still unit-untestable (dialog-callback or Gio territory):

- state.go - `GetOutputPathWidget`
- state.go - the `templateDir == ""` fallback inside `NewUIState`: whether
  `FindOldenEraTemplatesDir` succeeds depends on whether the game is installed
  on the machine running the tests, so only the branch that is true locally is
  ever measured.
  (returns a Gio widget). Covered by the integration suite.
- stateFiles.go - `handleSaveState` / `handleLoadState` success paths and
  `suggestDirectory` are only reachable through file-dialog callbacks
  (`Load`/`SaveAs` pick handlers); unit tests assert the dialogs open, the
  integration suite exercises the load/save flows via the
  `integration_test`-gated `SaveStateToFile`/`LoadStateFromFile` exports.
- stateFiles.go - `SaveAs`'s "only record `currentPath` when the write
  succeeded" rule (review item §1.2) cannot be unit tested: the decision lives
  in the closure handed to `dialogs.NewSaveFileDialog`, which is stored in the
  unexported `onSave` field and normally fires only from `confirmSelection` /
  `confirmOverwrite`, both of which need a `layout.Context`. The regression is
  covered by `test/integration/stateSaveAs_integration_test.go` through the
  `integration_test`-gated `DialogHost.GetTopDialog` and
  `FileExplorerDialog.ConfirmSave` exports.
- stateFiles.go - `PickOutputDir` / `RevealOutputDir` only open dialogs whose
  behavior lives in the dialog implementations.
- stateGeneration.go - `reapplyManualEdits` castle-change branch requires a
  generation-then-castle-option-change sequence entangled with the real
  mapper; exercised by the integration suite's manual-edit scenarios.
- test/test_helpers/integration_common - the `integration_test`-tagged files
  (`appRunner.go`, `appRunnerSnapshots.go`, `runMode.go`, `tabCalibration.go`)
  need `editor.Window` + a headless GPU context, so §4.6 forbids unit tests;
  they are exercised by the gated integration/performance suites (snapshot
  capture/validation via `window_snapshot_integration_test.go`). The untagged
  helpers (`snapshotComparer.go`, `snapshotMasker.go`, `snapshotStore.go`) have
  dedicated unit tests under `test/unit/test/test_helpers/integration_common/`.

- internal/services/template_generator/providers/topology/base/topologyConnectionService.go -
  private connection, portal, repair, guard, and road policy is reachable through
  the public `TopologyBase` methods and covered by that file's mirrored unit-test
  folder. Do not add test-only exports or a duplicate public service solely for
  per-file coverage attribution.

## Unreachable defensive branches (unit-test coverage gaps by design)

- internal/services/template_generator/providers/topology/geometricHubLayout.go -
  `connectInteriorStables` early-return for `len(stables) == 0`: the growth
  ladder in `distributeGeometricHubSlots` only assigns interiors after every
  gap holds 2 stables, so a hexagon with interiors always has both flanking
  stables. The guard is purely defensive; do not add seams to reach it.

- internal/composition/previewGeneratorProvider.go - the `err != nil` branch of
  `providePreviewGenerator`: `preview_service.NewPreviewGenerator` only fails
  when the `go:embed`-ed preview assets cannot be decoded, which cannot happen
  in a build that compiled. The branch exists to keep the injector error-free
  (a broken asset set degrades to "no preview images" instead of failing
  construction); reaching it would require an injectable asset provider seam.

- internal/repositories/atomicFileWriter.go - the struct is private to the
  package and has no test folder of its own. It is exercised through the three
  repositories, which each have a mirror folder under
  test/unit/internal/repositories/. Its remaining uncovered lines are the
  `Close` and `Sync` error branches of `encodeToTemporaryFile`: making either
  fail needs a genuinely full filesystem or a test-only seam in production
  code, both of which AGENTS.md 4.6 rules out. Review item 1.6's requested
  "close failure is propagated" test is therefore not written; the truncation
  half of that item is covered by
  `TestWhenEncodingFailsOverAnExistingPreview_LeavesTheDestinationUntouched`.

- internal/helpers/io.go - `getVDFContent`, `getVDFFilePath`, `getSteamPath`,
  `getBasePath`, and internal/helpers/io_windows.go -
  `getSteamPathFromRegistry`: this is the Steam/Olden-Era install discovery
  chain. It reads the Windows registry and the real Steam `libraryfolders.vdf`
  from the host filesystem, so its result depends entirely on whether the
  machine running the tests has Steam and the game installed. Covering it needs
  an injectable filesystem/registry seam that does not exist today; the public
  entry points that use it are covered through their error paths instead.

- internal/services/template_generator/providers/topology/base/topologyConnectionService.go -
  `buildShiftDerangement`: reached only after `buildNonAdjacentDerangement`
  fails 100 consecutive randomized attempts, which cannot be forced without
  seeding control over `math/rand` inside production code. Deterministic
  fallback, purely defensive; do not add a seam to reach it.
