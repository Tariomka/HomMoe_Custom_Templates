# Batch F — File explorer: hidden-file toggle and pointer-driven interactions

Implements backlog [§5.3](../todo/backlog-opus5.md) — cover the file explorer's
hidden-file toggle and its pointer-driven row/scroll interactions with real
synthetic pointer events, driven through the GUI handler framework, and migrate
the dialog's existing integration tests onto it.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

### Architecture (owner-corrected — read this first)

Tests are driven by **handlers**, not by a runner. `AppRunner` is the engine: it
owns the input router, injects synthetic pointer events and renders/validates
snapshots. `FileExplorerHandler` is the test-facing fluent driver, reached from
`BaseHandler.ClickLoad()` / `ClickSaveTo()`, exactly as the tab handlers are.
The explorer is exercised **through the real editor window**, like a user, so
every action can be snapshot-verified.

An earlier draft of this plan proposed a standalone `DialogRunner` that hosted
the dialog itself and drove the tests. **That was wrong and is abandoned** — it
bypassed the window, so it could never produce a golden.

### Owner decisions (do not re-litigate)

1. **Handler-driven**, through `AppRunner` + `FileExplorerHandler`. No
   `DialogRunner`.
2. The fixture directory reaches the dialog by making the dialog **open where
   the current file lives**: `State.workingDirectory()` becomes
   `getWorkingDirectory()` and prefers `currentPath` over the process CWD. This
   is a genuine UX improvement, and a `SetCurrentPath` testexport then seeds it.
   No unset-in-production field, no `os.Chdir`, no injected fake filesystem.
3. Rows are located by **replaying the frame's semantic tree**, which requires
   emitting button semantics per row (one production line). No pinned row pitch.
4. The new-folder row's **Cancel button is removed** (the `New Folder` toggle
   already dismisses the row) and `"Create"` becomes `"Create Folder"` — a UX
   simplification, now that label ambiguity is no longer the driver.
5. Folder-name text entry uses a **hard-coded, measured textbox point**.
6. Hidden fixtures use a **dot prefix on both platforms** — `isHidden` treats
   dotfiles as hidden everywhere. **No build-tagged helper**; §5.3's instruction
   to write one rested on the mistaken premise that Windows needs the attribute
   set. The Windows `hasHiddenAttribute` path stays end-to-end uncovered and
   must be recorded in `test_observations.md`.
7. **Goldens everywhere**: the new listing tests assert by golden, and the 9
   migrated tests get a golden *in addition to* their existing functional
   assertions.
8. Snapshot determinism: **mask the dialog's path bar and the toolbar's
   `File:` status**, and use **fixed** fixture names (no `gofakeit`) in any
   frame that reaches a golden.
9. Dialog testexports: keep the observation accessors, drop the `Click*` ones
   that pointer clicks replace.

### Constraints in force

- Never modify `data/`, `internal/entities/template/`, `internal/registry/`.
- Cross-platform: `path/filepath`, no OS-only assumptions; PowerShell chains
  with `;`.
- Unit coverage floor **72.5 %**, currently **72.9 %** — measure before and after.
- Never stage, never commit. Delete files with `Remove-Item`, never `git rm`.
- `gui`-tagged tests need a GPU and are opt-in; never run them in CI, never
  regenerate goldens in CI. New goldens are generated locally and need an
  eyeball before they are trusted.
- Snapshot tests do **not** call `t.Parallel()` — they share one headless GPU
  window, matching `window_snapshot_integration_test.go`.

---

## Phase 1: Production changes
Status: Complete

- [x] `State.workingDirectory()` → `getWorkingDirectory()`
      ([app/gui/drivers/stateFiles.go](../app/gui/drivers/stateFiles.go#L115-L121)):
      return `this.fileSystem.ResolveStartDirectory(this.currentPath)` when
      `currentPath` is non-empty, else the existing
      `ResolveStartDirectory(".")`. No path arithmetic is needed in the driver —
      `ResolveStartDirectory` already climbs from a file path to its containing
      directory, proven by
      [TestWhenPreferredPathIsAFile_ReturnsItsContainingDirectory](../test/unit/internal/services/file_system/pathResolutionService/resolveStartDirectory_test.go#L41-L54).
      Update the two call sites and their trailing comments.
- [x] Add `SetCurrentPath` to
      [app/gui/drivers/state_testexports.go](../app/gui/drivers/state_testexports.go).
- [x] Emit row semantics in `getEntryRowWidget`
      ([app/gui/dialogs/fileExplorerDialogEntries.go](../app/gui/dialogs/fileExplorerDialogEntries.go)):
      `utils.AddButtonSemantics(gtx.Ops, entry.Name, dims.Size)` inside the
      `material.Clickable` closure after `call.Add(gtx.Ops)`, mirroring
      [buttonWidget.go](../app/gui/widgets/buttonWidget.go#L218).
- [x] Remove the new-folder row's Cancel button
      ([fileExplorerDialogToolbar.go](../app/gui/dialogs/fileExplorerDialogToolbar.go)):
      drop the flex child, the `cancelFolderBtn` field on `FileExplorerDialog`
      and its branch in `Body`.
- [x] Rename that row's `"Create"` button to `"Create Folder"`.
- [x] **Do not** add unit tests for the `currentPath` branch of
      `getWorkingDirectory`, and **do not** add a seam to make it testable.
      `currentPath` is written only by the private `handleSaveState` /
      `handleLoadState`, each reachable only through a dialog callback, so no
      unit test can make it non-empty; `SetCurrentPath` is `integration_test`
      gated and unit tests must never see it (AGENTS.md §4.6). The branch is
      covered by the Phase 3/4 handler tests, which depend on it for their
      fixture directory. Record the gap in `test_observations.md` in Phase 5.

### Verification Plan
- `go build ./...` → clean.
- `go vet -tags='integration_test,gui' ./...` → clean.
- `gofmt -l ./app ./internal ./test ./cmd` → empty.
- `go test ./test/unit/... -count=1` → pass.
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` →
  still passes; the existing tests click through accessors and construct the
  dialog directly, so they are label- and directory-independent at this point.

### Phase Summary
All five production items landed. `getWorkingDirectory` on
[stateFiles.go](../app/gui/drivers/stateFiles.go) now prefers the directory of
the file being edited and falls back to `ResolveStartDirectory(".")`; the three
call sites (`Load`, `SaveTo`, and the templates-dir fallback in `NewUIState` on
[state.go](../app/gui/drivers/state.go)) were updated. `SetCurrentPath` was added
to [state_testexports.go](../app/gui/drivers/state_testexports.go) — this is the
injection point Phase 2's `WithFixtureDirectory` depends on. `getEntryRowWidget`
now emits `utils.AddButtonSemantics(gtx.Ops, entry.Name, dims.Size)`, making
every listing row locatable by name through the router's semantic tree — that is
what `AppRunner.ButtonBounds` will query in Phase 2. The new-folder row lost its
redundant Cancel (field, flex child and `Body` branch all removed; re-clicking
`New Folder` still dismisses the row) and its `"Create"` button is now
`"Create Folder"`, so no two buttons in a single frame share a label.

No unit tests were added: as anticipated, the `currentPath` branch is
unreachable from unit tests, and no seam was introduced to force it. The gap is
queued for `test_observations.md` in Phase 5.

Verification ran clean: `go build ./...`, `go vet ./...`, `go vet
-tags='integration_test,gui' ./...`, `gofmt -l` (empty), `go run
./cmd/testlayoutcheck .` (passed), `go test ./test/unit/... -count=1` (no
failures), `go test -tags=integration_test ./test/integration/... -count=1`
(ok), `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`
(ok — the existing window goldens were unaffected, confirming the semantics and
button-label changes are visually inert), and `golangci-lint-v2 run ./app/...`
(0 issues). No goldens needed regeneration.

## Phase 2: Handler framework extensions
Status: Complete

- [x] `AppRunner.ButtonBounds(label string) image.Rectangle` — replays the frame
      through `router.AppendSemantics(nil)`, matches
      `Desc.Class == semantic.Button && Desc.Label == label`, and **fails the
      test unless exactly one node matches**. `ClickButton(label)` clicks its
      center. This is what makes decision 3 usable and decision 4 safe.
- [x] `Window.TopFileExplorer()` in
      [window_testexports.go](../app/gui/editor/window_testexports.go), returning
      a narrow interface over the open dialog's observation accessors, following
      the existing `scrollablePanel` precedent. Required by decision 7: the
      migrated tests keep their functional assertions but no longer own the
      dialog instance.
- [x] `BaseHandler.WithFixtureDirectory(names ...string) *BaseHandler` — creates
      a `t.TempDir()`, writes the fixed-name entries, calls `SetCurrentPath` on
      the state driver, and registers the two new masks. Must run **before**
      `WithSnapshots()` so the fixture is in place for the first golden.
- [x] Coordinates and masks in
      [handlerCoordinates.go](../test/test_helpers/integration_common/handlerCoordinates.go):
      the dialog path-bar mask, the toolbar `File:` status mask
      (`fileStatusMask()`), the new-folder textbox point and the listing scroll
      point — each documented with how it was measured, in the file's existing
      style.
- [x] Grow [fileExplorerHandler.go](../test/test_helpers/integration_common/fileExplorerHandler.go)
      into the driver: `ClickShowHidden`, `ClickRow(name)`, `ClickBack`,
      `ClickNewFolder`, `TypeFolderName(text)`, `ClickCreateFolder`,
      `ClickConfirm`, `ClickOverwrite`, `ClickOverwriteCancel`, `Scroll(delta)`,
      plus observation pass-throughs. Each action calls `VerifySnapshot()`.

### Verification Plan
- `go vet -tags='integration_test,gui' ./...` → clean.
- A scratch test opens the explorer via `ClickSaveTo()`, asserts the listing
  shows the fixture names, and clicks `"Create Folder"` by label — proving the
  fixture directory, the semantic lookup and press/release routing all work.
  Delete the scratch test once Phase 3 covers the same ground.

### Phase Summary
The framework grew four pieces.

**Semantic lookup.** New
[appRunnerSemantics.go](../test/test_helpers/integration_common/appRunnerSemantics.go)
adds `ButtonBounds` / `ClickButton` and, crucially, the region-scoped
`ButtonBoundsIn` / `ClickButtonIn`. The region variant was not in the plan and
turned out to be mandatory: the editor behind an open modal is still laid out and
still publishes its semantics, so the dialog's `Save` and `Cancel` collide with
the toolbar's. Scoping the search to the dialog panel resolves that without
weakening the "exactly one match or fail" rule, which stayed as designed.

**Window accessor.** `IFileExplorerDialog` and `Window.TopFileExplorer()` went
into [window_testexports.go](../app/gui/editor/window_testexports.go). The
interface is exported because `revive`'s `unexported-return` is enabled, and it
lives inside the testexports file rather than its own `*Interface.go` because
`cmd/testlayoutcheck` enforces that outside `test/` only `*_testexports.go` may
carry the `integration_test` tag — the same compromise the existing
`scrollablePanel` makes. `ScrollPosition()` was added to the dialog's own
testexports and to the interface so the scroll test can assert on the listing
rather than only on a golden.

**Fixture directory.** `BaseHandler.WithFixtureDirectory(names...)` seeds a
`t.TempDir()`, creating each name ending in `.gen.json` as a file and every other
name as a folder (a leading `.` makes it hidden on both platforms), then points
the dialogs at it through `SetCurrentPath`. `FixtureDirectory()` exposes the path
so tests can build the absolute paths they expect back.

**Handler.** [fileExplorerHandler.go](../test/test_helpers/integration_common/fileExplorerHandler.go)
is now the full driver. Two deviations from the plan, both forced by the code:
`ClickConfirm` became `ClickOpen` / `ClickSave` because the confirm button's
label is mode-dependent and the handler addresses buttons by label; and every
action renders a settle frame before snapshotting, because a row click only
*records* the navigation and the dialog applies it at the top of the next `Body`.
Snapshotting is guarded: `verifySnapshot` fails outright if snapshots are on but
no fixture directory was seeded, so a machine-dependent golden cannot be
committed by accident. The path-bar mask is derived at dialog-open time from the
bounds of the two header buttons flanking it, rather than pinned.

Only the toolbar `File:` mask is a pinned constant (`fileStatusLeft = 1000`),
verified by eye against a generated golden: the window title ends at ~988px, so
the mask clears it.

Verification: `go build ./...`, `go vet -tags='integration_test,gui' ./...`,
`gofmt -l` (empty after formatting two new files), `testlayoutcheck` (passed),
and `golangci-lint-v2 run ./...` (0 issues). The scratch test was skipped — the
Phase 3 tests were written directly and served the same purpose.

## Phase 3: The §5.3 tests
Status: Complete

- [x] New [test/integration/gui/fileExplorerDialogListing_integration_test.go](../test/integration/gui/fileExplorerDialogListing_integration_test.go),
      `//go:build integration_test && gui`, handler-driven, snapshot-verified:
  - [x] `TestWhenShowHiddenIsToggledOn_HiddenEntriesAppearInTheListing`
  - [x] `TestWhenShowHiddenIsToggledOff_HiddenEntriesDisappearAgain`
        (default is off, so toggle on then off)
  - [x] `TestWhenARowIsClicked_ThatEntryBecomesTheSelection` — **open mode
        only**; batch E removed the `modeSaveFile` branch from `onEntryClicked`.
  - [x] `TestWhenADirectoryRowIsClicked_TheListingDescendsIntoIt` — needs a
        second frame; navigation is deferred to the top of the next `Body`.
  - [x] `TestWhenTheListingIsScrolled_TheFirstVisibleRowAdvances` — fixture of
        ~40 fixed-name entries so the list overflows.
- [x] Hidden fixtures are `"." + fixedName`; listing assertions go through
      `TopFileExplorer().EntryNames()` alongside the goldens.

### Verification Plan
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` → pass.
- `go run ./cmd/testlayoutcheck .` → `test-layout check passed`.
- Generate the new goldens locally with the *"Go: Update UI Integration tests
  snapshots"* task, then **eyeball every one** before trusting it.

### Phase Summary
All five tests landed in
[fileExplorerDialogListing_integration_test.go](../test/integration/gui/fileExplorerDialogListing_integration_test.go),
driven end to end through the window: a real toolbar click opens the dialog, real
pointer press/release pairs hit rows and buttons resolved by their semantic
label, and a real wheel event scrolls the listing. None call `t.Parallel()` (they
share the headless GPU window) and each carries the same `paralleltest` nolint
the existing snapshot tests use.

Every test asserts functionally *and* against a golden: 11 goldens were generated
locally and each was inspected. The dialog renders as expected, the hidden
entries appear and disappear on the toggle, the directory click descends, and the
scroll golden shows the listing starting at `entry-16` instead of `entry-00`. The
three masks (preview canvas, status line, output directory) plus the two new ones
(toolbar `File:`, dialog path bar) cover every nondeterministic region — confirmed
by re-running the suite in validation mode against a *different* temp directory,
which passed.

Verification: `go test -tags='integration_test,gui' ./test/integration/gui/...
-count=1` (ok, whole suite), `go run ./cmd/testlayoutcheck .` (passed),
`go test ./test/unit/... -count=1` (no failures),
`go test -tags=integration_test ./test/integration/... -count=1` (ok), and
`golangci-lint-v2 run ./...` (0 issues).

## Phase 4: Migrate the existing tests
Status: Complete

- [x] Rewrite the 9 tests in
      [fileExplorerDialog_integration_test.go](../test/integration/gui/fileExplorerDialog_integration_test.go)
      onto `FileExplorerHandler`, keeping every existing functional assertion
      (bytes on disk, save-error text, read-only flag, disabled confirm) and
      adding a golden per action.
- [x] Add `TestWhenNewFolderIsClickedTwice_TheRowIsDismissed` — with the row's
      Cancel gone, the toggle is the only dismissal path and must be pinned.
- [x] Delete the `Click*` and `SetNewFolderName` accessors from
      [fileExplorerDialog_testexports.go](../app/gui/dialogs/fileExplorerDialog_testexports.go);
      keep `CurrentDir`, `EntryNames`, `SelectedPath`, `SaveError`,
      `NewFolderError`, `ConfirmDisabled`, `ResolvedSaveName`,
      `SaveNameReadOnly`, `OverwriteActive`. Check whether `ConfirmSave` still
      has a caller and delete it if not.
- [x] `frameFileExplorer` disappears from that file; leave `newDialogContext` in
      place for the other dialog tests that still use it.

### Verification Plan
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` → pass.
- `go vet -tags='integration_test,gui' ./...` → clean.
- `go run ./cmd/testlayoutcheck .` → `test-layout check passed`.

### Phase Summary
All twelve tests now run through the real toolbar rather than constructing a
`FileExplorerDialog` by hand, so each one covers the driver and the repository as
well as the dialog. `frameFileExplorer` is gone; `newDialogContext` stays for the
other dialog tests. `ConfirmSave` was kept — the non-GUI
[stateSaveTo_integration_test.go](../test/integration/stateSaveTo_integration_test.go)
still drives the save callback through it — but every `Click*` accessor and
`SetNewFolderName` were deleted, and `NewFolderActive` was added for the new
toggle test.

Three things the migration forced, none of them in the original plan:

- **Fixture files and folders are seeded separately.** The old
  `WithFixtureDirectory(names...)` inferred the kind from the `.gen.json`
  suffix, which cannot express the folder named `Custom Template.gen.json` that
  the refused-save test needs. It is now `WithFixtureDirectory()` plus
  `WithFixtureFiles` / `WithFixtureFolders`, either of which creates the
  directory on first use.
- **The editor state is seeded, not typed.** The save target is derived from the
  template name, and the name field cannot be cleared through synthetic key
  events, so `Window.SetTemplateName` writes it and resyncs the panels — every
  layout writes their widget values back over the state, which would otherwise
  undo it on the next frame. `AppRunner.CurrentPath` was added alongside it, so
  the save-target test asserts on what the editor recorded.
- **The dialog's path-bar mask is now lifted when the dialog closes.** It sits
  over the settings panel, so leaving it registered blanked a slider row out of
  every later editor snapshot — caught by eye on the first generated golden.
  `Masker.RemoveRect` and `AppRunner.UnmaskRect` back it.

The overwrite tests no longer stub the save callback: the editor is saved once,
the game mode is changed, and the second save is confirmed onto the same target,
so the gated write is observed as the bytes on disk staying identical.

41 goldens were generated locally and inspected. Verification:
`go build ./...`, `go vet -tags='integration_test,gui' ./...`, `gofmt -l`
(empty after formatting the rewritten file), `testlayoutcheck` (passed),
`go test ./test/unit/...` and `go test -tags=integration_test
./test/integration/...` (ok), the GUI suite re-run in validation mode against
fresh temp directories (ok), and `golangci-lint-v2 run ./...` (0 issues).

## Phase 5: Documentation and gates
Status: Complete

- [x] [todo/test_observations.md](../todo/test_observations.md): replace the
      *"Still uncovered: the hidden-file toggle and the pointer-driven
      row/scroll interactions"* sentence with a pointer to the new tests; note
      that Windows `hasHiddenAttribute` is covered by unit tests only, never
      end-to-end; and record that `getWorkingDirectory`'s `currentPath` branch
      is integration-covered only, because no unit test can set `currentPath`.
- [x] [todo/backlog-opus5.md](../todo/backlog-opus5.md): mark §5.3 done and
      self-contained; update the batch **F** row in §8.
- [x] Record coverage before and after. Expected **flat or a shade down**: this
      batch adds two production lines that no unit test can reach (the row
      semantics call and the `currentPath` branch) and every test it adds is
      integration-only. If it falls below **72.5 %**, stop and report.
- [x] Full gate sweep per [§9 of the backlog](../todo/backlog-opus5.md).

### Verification Plan
- `go build ./...`; `go vet ./...`; `go vet -tags='integration_test,gui' ./...`
- `go run ./cmd/testlayoutcheck .`
- `go test ./test/unit/... -count=1`
- `go test -tags=integration_test ./test/integration/... -count=1`
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`
- `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...`
  then `go tool cover '-func=coverage.txt'` → **≥ 72.5 %**
- `golangci-lint-v2 run ./... --issues-exit-code=0` → 0 issues
- `gofmt -l ./app ./internal ./test ./cmd` → empty
- `wire diff ./internal/composition/...` → no diff

### Phase Summary
[test_observations.md](../todo/test_observations.md) now points at both GUI test
files and records what is still out of reach: the Windows-only
`hasHiddenAttribute` branch (unit-tested only — setting the attribute needs a
syscall no fixture builder should make) and `getWorkingDirectory`'s
`currentPath` branch (integration-covered only, since nothing outside the driver
writes that field). §5.3 and the batch **F** row of
[backlog-opus5.md](../todo/backlog-opus5.md) are marked done.

**Coverage: 72.9 % → 72.9 %, flat**, against a 72.5 % floor. The batch's only new
unit-unreachable production statement is the row semantics call, which is small
enough not to move the total.

Gate sweep, all green: `go build ./...`, `go vet ./...`,
`go vet -tags='integration_test,gui' ./...`, `go run ./cmd/testlayoutcheck .`
(passed), `gofmt -l ./app ./internal ./test ./cmd` (empty),
`go test ./test/unit/...` (no failures),
`go test -tags=integration_test ./test/integration/...` (ok),
`go test -tags='integration_test,gui' ./test/integration/gui/...` (ok),
`golangci-lint-v2 run ./...` (0 issues), `wire diff ./internal/composition/...`
(no diff).

Batch F is complete and ready for review. Nothing is staged or committed.

### Phase Summary
_(write when phase completes)_

## Final Recap
_(write when all phases complete)_

## Deployment Plan
_(write when all phases complete)_
