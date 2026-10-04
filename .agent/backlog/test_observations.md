# Test Observations - untestable / integration-only code registry

Referenced by AGENTS.md §4.6. Records code that cannot be exercised through
public APIs in unit tests, so per-file coverage gaps here are intentional.

## Gio-UI-heavy code (integration-suite territory, no unit tests)

- app/gui/dialogs/zoneEditorCanvas.go `drawEdges` - Batch F draws effective portals at
  `DefaultConnectionLineSmall` (1.5dp), matching the Preview tab. At the headless
  harness's 1x scale `gtx.Dp` rounds both 1.5dp and 2dp to 2px, so a pixel test
  cannot tell a portal from a direct edge by width. The colour half of the same
  change is pixel-tested in roadStyleVisuals_integration_test.go.
  The type dropdown shaping (`connectionTypeOptions`, `connectionTypeLabel`) is
  covered by test/integration/gui/zoneEditorConnectionType_integration_test.go.

- app/gui/dialogs/zoneEditorZoneProps.go and zoneEditorConnectionProps.go - Batch E
  quality-mutation installation and Custom preset synchronization are covered by
  test/integration/gui/zoneEditorProperties_integration_test.go. They remain 0% in
  the unit profile; owner accepted total coverage 75.1% -> 74.9% on 2026-09-12
  (historical figures).
  The non-nil selected-index rebinding branch is defensive: selecting a zone clears
  the connection selection, so current real inputs cannot enter that branch.

- app/gui/dialogs/zoneEditorDialog.go, zoneEditorGraphState.go and
  zoneEditorStatusKey.go - Batch G's status-line graph cache (`layoutStatus`,
  `statusKey`, `markGraphDirty` in the seven structural mutators,
  `requestLateStatusRedraw` at the end of `Body`). They are covered by
  test/integration/gui/zoneEditorGraphCache_integration_test.go, which counts
  handler calls per frame and reads the aux router's wakeup. In the unit profile
  they added 15 uncovered statements when Batch G landed (2026-09-28, historical:
  covered 6911 unchanged, denominator 9246 -> 9261, reported total 74.6% -> 74.5%).

- app/gui/program.go - `StartApplication`, `eventLoop` and
  `getAndConfigureWindow` are the Gio bootstrap: they create a real
  `app.Window`, block in `app.Main`, read `os.Args` process-wide and call
  `os.Exit` from the event loop, so none of them can be driven from a unit test.
  The `app.DestroyEvent` error branch (which since backlog item §1.4 also writes
  the failure to stderr, because the default `slog` handler discards unless
  `-with-logging` is passed) is reachable only when Gio itself fails to create
  or maintain the window - no GPU, or missing X/Wayland libraries - which cannot
  be provoked in-process. Making it testable would mean returning the error to
  main.go instead of exiting; the owner deferred that (see §1.4).

- app/gui/widgets/buttonWidget.go - all button constructors need a
  `layout.Context` + `material.Theme` text shaper to lay out; covered indirectly
  by the integration/performance suites that render the full editor window.
  The button-position debug semantics they emit come from the exported
  `utils.AddButtonSemantics` (app/gui/utils/buttonPositionLogger.go), which needs
  only ops, a label and dimensions. The semantic-op replay path IS unit-tested via
  test/unit/app/gui/utils/buttonPositionLogger/ (headless ops equivalent to what
  the buttons emit), and was verified end-to-end against a real
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
  test/integration/gui/fileExplorerDialog_integration_test.go and, since batch F
  (§5.3), test/integration/gui/fileExplorerDialogListing_integration_test.go.
  Both drive the real toolbar and inject real pointer events through
  `integration_common.FileExplorerHandler`, which resolves every button and
  listing row by its accessibility label, and both compare a golden image per
  action: open-and-load, save target resolution, save through the real state
  driver, the overwrite prompt (gated write, cancel, confirm), new-folder
  creation and dismissal, the existing-folder refusal, the disabled-confirm
  predicate, the hidden-file toggle in both directions, row selection, directory
  descent and wheel scrolling. All filesystem *policy* (listing, filtering,
  hidden entries, roots, path resolution, reserved names) moved to
  internal/services/file_system and is unit-tested there; what is left in the
  dialog is rendering and click wiring only.
  Still uncovered end-to-end: the Windows-only `hasHiddenAttribute` branch of
  the hidden-entry filter. The toggle tests exercise the dot-prefix rule, which
  is what both platforms share; the file-attribute rule is unit-tested only,
  because setting it needs a syscall no fixture builder should be making.
  Since batch E ("Save As" -> "Save To") the save-name field is read-only and
  fed only by `State.SaveTo`, which sanitizes and trims before appending the
  suffix. `resolveSaveTarget`'s whitespace-only rejection is therefore no longer
  reachable through any production path - only by calling
  `dialogs.NewSaveFileDialog` directly - so its test was dropped in favour of
  `TestWhenNoNameWasResolved_TheConfirmButtonIsDisabled`, which covers the state
  the UI can actually reach (an unnamed template).

- app/gui/dialogs/zoneEditorDialog.go with its canvas, snap, property-panel,
  geometry, and transient-state sibling files - the Manual Zone Editor
  (one primary dialog struct with rendering methods and UI state split by
  responsibility).
  Much less of this is untestable than it was. Since 2026-08-08 (review item
  §2.6, Batch 15) ALL the geometry moved out of the dialog -
  `BuildGeometry`, the obstacle bulge, edge/node hit-testing, the other-zone
  guides and the grid step are unit-tested in internal/services/connection_editor
  (>=92.9%, most 100%, as measured 2026-08-08), and `groupConnectionsByPair` now
  lives with the shared curve builder in internal/models/preview
  (connectionCurveLayout.go). The dialog's
  canvas/snap files are thin call-throughs on `IZoneEditorHandler`. The revert
  semantics live in `drivers.State` (`PreviewBaseZones`, `ApplyEditedZones`) and
  are unit-tested under test/unit/app/gui/drivers/stateManualEdits.
  What is left in the dialog really is rendering and click wiring, and the
  `integration_test,gui` suite drives it through real frames with
  `widget.Clickable.Click` and the `zoneEditorDialog_testexports.go` accessors:
  Undo, Revert to Base (success, failure, flag round-trip on Apply), Apply,
  delete-selected, the add-connection/add-zone mode toggles, button labels read
  back off the semantics tree, plus the Phase 0 numeric geometry pins that guard
  the extraction.
  Since 2026-08-11 (batch H) the synthetic-pointer and keyboard paths are driven
  too, through the `integration_common` handlers: the pointer flows in
  test/integration/gui/zoneEditorPointer_integration_test.go (zone drag +
  Apply round-trip, snapping, drag-to-connect and the drag that ends on empty
  canvas, right-click delete of a curve, and placing a zone from Add zone mode)
  and the property panels in
  test/integration/gui/zoneEditorProperties_integration_test.go (the zone Size /
  Guard x / Weekly + editors including `Size` clamping and rounding, the neutral
  Quality and Castles dropdowns with their `ApplyZoneEditorQuality` reprofile,
  the connection guard value typed and rejected, the Type / Guard zone / Guard
  preset / Weekly dropdowns, and the Advanced options checkbox with the Match
  group, Guard escape and Sim turn squad rows it reveals).
  Still uncovered, with reasons:
  - The zone **name** row is a read-only `material.Body1` label and the dialog
    offers no rename, so there is no typing path to drive.
  - `integration_common`'s `zoneEditorZone*Y` row coordinates were measured on a
    zone whose note wraps to one line, which a player spawn and a neutral zone
    do but the shared `Hub` does not, so the Hub's property rows cannot be
    clicked through the handler. `zoneRowY` only compensates for neutral zones.
    The rows themselves are the same code for every zone, so the gap costs no
    coverage - it is a harness limitation to fix if a Hub-specific behaviour
    ever needs driving.

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

Unit tests use `NewUIState(handler, fileSystem, regeneration, findTemplateDir)`
with the `test_helpers` handler mocks. Game templates folder detection is
injected through `IFileSystemHandler.FindGameTemplateDirectory`, so every
detection branch of `NewUIState` is unit-tested
(test/unit/app/gui/drivers/state/newUIState_test.go).
Still unit-untestable (dialog-callback or Gio territory):

- state.go - `GetOutputPathWidget` (returns a Gio widget). Covered by the
  integration suite.
- stateFiles.go - `handleLoadState` and the first `handleSaveState` of a
  document are only reachable through file-dialog callbacks (`Load`/`SaveTo`
  pick handlers), which are what establish `currentPath`; afterwards the public
  `Save` calls `handleSaveState` directly, but a unit test cannot set
  `currentPath` to get there. Unit tests assert the dialogs open; the
  integration suite exercises the load/save flows via the
  `integration_test`-gated `SaveStateToFile`/`LoadStateFromFile` exports.
- stateFiles.go - `PickOutputDir` / `RevealOutputDir` only open dialogs whose
  behavior lives in the dialog implementations.
- stateFiles.go - `getWorkingDirectory`'s `currentPath != ""` branch: nothing
  outside the driver can set `currentPath`, which is written only by the private
  `handleSaveState` / `handleLoadState`, so a unit test always measures the
  fallback. Batch F covered the branch end-to-end instead, through the
  `integration_test`-gated `State.SetCurrentPath` that the GUI suite's fixture
  directories are seeded with.
- stateManualEdits.go - `reapplyManualEdits`' castle-change branch is NOT a
  testability limitation any more: it is reached from the public `Generate` when
  the mocked `IRegenerationHandler.DecideManualEditReapplication` returns a
  non-nil `ReapplyWithCastleChanges` with at least one change flag set and the
  template has variants, and it then calls the
  mocked `IGuiHandler.ReapplyCastleSettings`. No unit test covers it yet; the
  integration suite's manual-edit scenarios do. Recorded, not assigned.
- test/test_helpers/integration_common - the `integration_test`-tagged helpers
  (the app runner, snapshots, run mode, tab and dialog handlers, coordinates)
  need `editor.Window` + a headless GPU context, so §4.6 forbids unit tests;
  they are exercised by the gated integration/performance suites (snapshot
  capture/validation via `window_snapshot_integration_test.go`). The untagged
  `snapshot/` subpackage (comparer, difference, masker, store) has dedicated
  unit tests under `test/unit/test/test_helpers/integration_common/snapshot/`.

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

- internal/helpers/io.go - the game templates folder discovery chain. Fixture
  tests in test/unit/internal/helpers/io/ cover `FindOldenEraTemplatesDir` per
  platform: on Windows the user-profile `my_map_templates` glob (temporary
  `USERPROFILE`), on other platforms Steam's `libraryfolders.vdf` and the Proton
  prefix (temporary `HOME`), success and not-found. Each test skips on the other
  platform, so a single run only measures one side. Still host-dependent:
  io_windows.go `getSteamPathFromRegistry` and the Windows Steam-path fallbacks in
  `getSteamPath` (`ProgramFiles(x86)`, the hard-coded default), which read the
  real registry/environment; io_other.go is the non-Windows no-op registry
  lookup. Covering the registry needs an injectable seam that does not exist
  today.

- internal/services/template_generator/providers/topology/base/topologyConnectionService.go -
  `buildShiftDerangement`: reached only after `buildNonAdjacentDerangement`
  fails 100 consecutive randomized attempts, which cannot be forced without
  seeding control over `math/rand` inside production code. Deterministic
  fallback, purely defensive; do not add a seam to reach it.

- Earlier backlog's Batch I Phase 4 (2026-08-22; not the review's Batch I) - roughly 55 test-local identifiers named `dto`
  or `stateDto` now hold an `EditorStateModel` rather than a DTO, mostly as the
  closure parameter of `UpdateState` / `UpdateCurrentState` (e.g.
  test/unit/app/gui/drivers/state/, test/unit/app/gui/models/editorState/,
  test/unit/internal/handlers/guiHandler/). They were deliberately left as-is:
  in most of those closures the enclosing scope already binds `state` to the
  driver `State`, so a blind rename to `state` would shadow it, and the gain is
  purely cosmetic. Production-side names were fixed in the same phase. Rename
  them opportunistically when a file is edited for another reason.

- Earlier backlog's Batch I Phase 6 (2026-08-31; not the review's Batch I) - **the
  per-frame allocation budget has no automated guard.** The phase cut
  `BenchmarkEditorWindow_TabCycling` from ~12,690 to ~4,773 allocs/op (historical
  figures from that date), but nothing fails if it climbs back: the benchmark
  needs a GPU and carries both the `integration_test` and `gui` tags, while the
  CI performance job runs with `integration_test` only, so this benchmark never
  runs in CI. (GUI integration tests themselves do run in CI: the PR workflow's
  Mesa/Xvfb job runs `./test/integration/gui/...` on every pull request.) An
  `allocs/op` assertion was considered and
  rejected - `testing.AllocsPerRun` over a Gio frame is dominated by rendering
  and font shaping, so a threshold tight enough to catch a regression in
  `EditorState.Clone` would be far too flaky to keep. The recorded figures in
  backlog §1.5 are the reference; re-measure by hand when touching
  `EditorState.Clone`, `linq.SelectSlice` or the clone helpers in
  internal/models/editor_state_model/ and internal/helpers/editor_state_helpers/.

- Converter move into the template mapper (2026-09-07) -
  internal/models/template_model/template_variant_model/zone.go `ToZoneModels`
  and `ToZoneEntities` have **no callers anywhere**, production or test, since
  the mapper grew its own `GetZoneModelList` / `GetZoneEntityList`. They are
  left at 0 % deliberately: writing a test for a function nothing calls would
  buy coverage and hide the fact that it is dead. Delete them, or give them a
  caller - do not paper over them with a test. Their singular forms
  (`ToZoneModel` / `ToZoneEntity`) are alive: `internal/models/editor_state_model`
  reaches them through template_model/converters.go, and cannot use the mapper
  instead because mappers already import models.

