# Session Carry-Forward

## 1. Session goal

Implement backlog **batch F** — `todo/backlog-opus5.md` §5.3, *file explorer:
hidden-file toggle and pointer-driven row/scroll interactions* — under a plan
file, per AGENTS.md §2.4.

## 2. Fixes applied

- **Listing rows were unaddressable by any driver.** `getEntryRowWidget` in
  [fileExplorerDialogEntries.go](../app/gui/dialogs/fileExplorerDialogEntries.go)
  drew a clickable row with no accessibility label, so a synthetic pointer test
  had no way to find "the row called X". It now emits
  `utils.AddButtonSemantics(gtx.Ops, entry.Name, dims.Size)`, matching every
  button in `app/gui/widgets`. This is the one production line the tests
  required.
- **The new-folder row had two dismissal paths, one of them redundant.**
  [fileExplorerDialogToolbar.go](../app/gui/dialogs/fileExplorerDialogToolbar.go)
  lost its `Cancel` flex child and `"Create"` became `"Create Folder"`; the
  `cancelFolderBtn` field and its `Body` branch went with it
  ([fileExplorerDialog.go](../app/gui/dialogs/fileExplorerDialog.go)). The
  `New Folder` button is a toggle and is now the only way to put the row away —
  pinned by a new test.
- **`suggestDirectory` was a name that described nothing.** Renamed to
  `getWorkingDirectory` in
  [stateFiles.go](../app/gui/drivers/stateFiles.go) with a doc comment stating
  what it actually guarantees, and its two callers plus the `NewUIState`
  templates-dir fallback in [state.go](../app/gui/drivers/state.go) updated.

## 3. Features added / changed

All of this is test infrastructure; no user-visible behaviour changed beyond §2.

- **Semantic button lookup.** New
  [appRunnerSemantics.go](../test/test_helpers/integration_common/appRunnerSemantics.go):
  `ButtonBounds(label)` / `ClickButton(label)` resolve a widget through
  `input.Router.AppendSemantics`, and `ButtonBoundsIn(area, label)` /
  `ClickButtonIn(area, label)` scope that to a rectangle. The scoped variants
  exist because the editor behind a modal scrim is still laid out and still
  publishes semantics, so the dialog's `Save` and `Cancel` collide with the
  toolbar's — a strict "exactly one match" lookup would fail without a region.
- **`FileExplorerHandler`** — a full fluent driver
  ([fileExplorerHandler.go](../test/test_helpers/integration_common/fileExplorerHandler.go)):
  `ClickShowHidden`, `ClickRow`, `ClickBack`, `ClickNewFolder`,
  `ClickCreateFolder`, `ClickOpen`, `ClickSave`, `ClickOverwrite`,
  `ClickOverwriteCancel`, `TypeFolderName`, `Scroll`, `Dialog()`, `Close()`,
  `Editor()`. It deliberately does **not** embed `BaseHandler` — a scrim absorbs
  background clicks, so promoting the tab clicks would be a lie.
- **Fixture directories.** `BaseHandler.WithFixtureDirectory()` plus
  `WithFixtureFiles` / `WithFixtureFolders`
  ([baseHandler.go](../test/test_helpers/integration_common/baseHandler.go)).
  They seed a `t.TempDir()`, point the dialogs at it via
  `State.SetCurrentPath`, and register the toolbar `File:` mask. Files and
  folders are seeded **separately** because the earlier suffix-inference could
  not express the folder named `Custom Template.gen.json` that the refused-save
  test needs.
- **Scoped snapshot masks.** `Masker.RemoveRect`
  ([masker.go](../test/test_helpers/integration_common/snapshot/masker.go)) and
  `AppRunner.UnmaskRect` / `SnapshotsEnabled`
  ([appRunnerSnapshots.go](../test/test_helpers/integration_common/appRunnerSnapshots.go)).
  The dialog's path-bar mask sits over the settings panel, so leaving it
  registered blanked a slider row out of every later editor snapshot.
- **Observation surface.** `editor.IFileExplorerDialog` + `TopFileExplorer()` +
  `SetTemplateName()`
  ([window_testexports.go](../app/gui/editor/window_testexports.go)), and the
  lock-guarded `AppRunner` wrappers `TopFileExplorer`, `SetCurrentPath`,
  `CurrentPath`, `SetTemplateName`.

## 4. File modifications

Production:

| File | Change |
| --- | --- |
| [app/gui/dialogs/fileExplorerDialogEntries.go](../app/gui/dialogs/fileExplorerDialogEntries.go) | Rows emit `AddButtonSemantics(entry.Name, …)`. |
| [app/gui/dialogs/fileExplorerDialogToolbar.go](../app/gui/dialogs/fileExplorerDialogToolbar.go) | New-folder `Cancel` removed; `"Create"` → `"Create Folder"`. |
| [app/gui/dialogs/fileExplorerDialog.go](../app/gui/dialogs/fileExplorerDialog.go) | `cancelFolderBtn` field and its `Body` branch removed. |
| [app/gui/drivers/stateFiles.go](../app/gui/drivers/stateFiles.go) | `suggestDirectory` → `getWorkingDirectory` + doc comment. |
| [app/gui/drivers/state.go](../app/gui/drivers/state.go) | Call site updated. |

Test-only exports (all `//go:build integration_test`):

| File | Change |
| --- | --- |
| [app/gui/drivers/state_testexports.go](../app/gui/drivers/state_testexports.go) | Added `SetCurrentPath`. |
| [app/gui/editor/window_testexports.go](../app/gui/editor/window_testexports.go) | Added `IFileExplorerDialog`, `TopFileExplorer`, `SetTemplateName`. |
| [app/gui/dialogs/fileExplorerDialog_testexports.go](../app/gui/dialogs/fileExplorerDialog_testexports.go) | Added `ScrollPosition`, `NewFolderActive`; **deleted** `ClickEntry`, `ClickUp`, `ClickConfirm`, `ClickOverwriteConfirm`, `ClickOverwriteCancel`, `ClickNewFolder`, `ClickCreateFolder`, `SetNewFolderName`. `ConfirmSave` kept — `stateSaveTo_integration_test.go` still uses it. |

Test helpers:

| File | Change |
| --- | --- |
| [appRunnerSemantics.go](../test/test_helpers/integration_common/appRunnerSemantics.go) | **New.** Semantic button lookup + click, window-wide and region-scoped. |
| [fileExplorerHandler.go](../test/test_helpers/integration_common/fileExplorerHandler.go) | Rewritten from a stub into the full driver. |
| [baseHandler.go](../test/test_helpers/integration_common/baseHandler.go) | Fixture builders; `ClickLoad`/`ClickSaveTo` return `*FileExplorerHandler`. |
| [appRunner.go](../test/test_helpers/integration_common/appRunner.go) | `TopFileExplorer`, `SetCurrentPath`, `CurrentPath`, `SetTemplateName`. |
| [appRunnerSnapshots.go](../test/test_helpers/integration_common/appRunnerSnapshots.go) | `SnapshotsEnabled`, `UnmaskRect`. |
| [snapshot/masker.go](../test/test_helpers/integration_common/snapshot/masker.go) | `RemoveRect`. |
| [handlerCoordinates.go](../test/test_helpers/integration_common/handlerCoordinates.go) | Dialog panel rect, listing scroll point, `fileStatusMask`, `headerBarSlack`. |

Docs / planning:

| File | Change |
| --- | --- |
| [plans/batch-f-file-explorer.md](../plans/batch-f-file-explorer.md) | **New.** Five phases, all `Complete`, each with a Phase Summary. |
| [todo/test_observations.md](../todo/test_observations.md) | File-explorer entry rewritten; records the two remaining gaps (see §8). |
| [todo/backlog-opus5.md](../todo/backlog-opus5.md) | §5.3 marked ✅ FIXED and self-contained; §0.3 row and §8 batch **F** row marked done. |

## 5. Tests added or updated

Added — [fileExplorerDialogListing_integration_test.go](../test/integration/gui/fileExplorerDialogListing_integration_test.go)
(new, `//go:build integration_test && gui`), the five §5.3 scenarios:

- `TestWhenShowHiddenIsToggledOn_HiddenEntriesAppearInTheListing`
- `TestWhenShowHiddenIsToggledOff_HiddenEntriesDisappearAgain`
- `TestWhenARowIsClicked_ThatEntryBecomesTheSelection` — **open mode only**, per
  the batch E scope change.
- `TestWhenADirectoryRowIsClicked_TheListingDescendsIntoIt`
- `TestWhenTheListingIsScrolled_TheFirstVisibleRowAdvances` — 40 fixed-name
  fixtures so the list overflows; asserts through the new `ScrollPosition()`.

Migrated — [fileExplorerDialog_integration_test.go](../test/integration/gui/fileExplorerDialog_integration_test.go),
12 tests, none of which construct a `FileExplorerDialog` any more. All go
through the real toolbar with real pointer events and a golden per action:

- The 9 original tests, every functional assertion kept.
- `TestWhenNewFolderIsClickedTwice_TheRowIsDismissed` — **new**, pins the toggle
  now that the row's `Cancel` is gone.
- `..._TheEditorStateIsLoaded` is now a genuine round trip: save → change the
  template name → load the saved file → assert the change was undone.
- The overwrite tests no longer stub the save callback. The editor is saved, the
  game mode is changed, and the second save is confirmed onto the same target —
  so the gated write is observed as the bytes on disk staying identical.
- `frameFileExplorer` is gone; `newDialogContext` stays for the other dialog
  tests.

Goldens: 11 for the listing file + 41 for the migrated file, all generated
locally with `-update` and inspected. The suite was then re-run in validation
mode against fresh temp directories to prove the masks are complete.

**Gate results — all green:**

| Gate | Result |
| --- | --- |
| `go build ./...` | clean |
| `go vet ./...` and `go vet -tags='integration_test,gui' ./...` | clean |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test ./test/unit/... -count=1` | pass |
| `go test -tags=integration_test ./test/integration/... -count=1` | pass |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | pass |
| Unit coverage | **72.9 %** (flat; floor is 72.5 %) |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **0 issues** |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `wire diff ./internal/composition/...` | no diff |

## 6. Git status snapshot

Branch: `AD/fixing_some_stuff_08-12`. **Working tree clean.**

```
d37d1fa (HEAD -> AD/fixing_some_stuff_08-12) Batch F done
f2f7108 (origin/AD/fixing_some_stuff_08-12) Batch E done
```

The owner reviewed, staged and committed batch F as `d37d1fa` (82 files,
+1247/−256). It is **not pushed** — `origin` is still at batch E. Nothing
unstaged is inherited.

Per AGENTS.md §2.5 and the owner's standing instruction this session — *"I keep
track of changes in staged, so don't touch that"* — the agent never staged or
committed anything. Do not run any git staging command.

## 7. Rejections / things the owner declined

- Nothing was declined this session. The one design question raised — adding
  `Window.TopFileExplorer()` as a test export — was explicitly approved:
  *"the new accessor makes sense"*.

## 8. Open questions

- **10 `window_snapshot_integration_test` goldens were regenerated as a side
  effect** and are in the commit. Running `-update` over
  `./test/integration/gui/...` rewrites every golden in the package, not just
  the new ones. They differ by a few hundred bytes and passed validation both
  before and after, so nothing regressed — but next time scope `-update` with
  `-run` to the tests being added.
- Two coverage gaps are now recorded in `todo/test_observations.md` rather than
  fixed: the Windows-only `hasHiddenAttribute` branch (unit-tested only —
  setting the attribute needs a syscall no fixture builder should make), and
  `getWorkingDirectory`'s `currentPath` branch (integration-covered only, since
  nothing outside the driver writes that field).
- Nothing is blocked.

## 9. Next recommended actions

1. Push `AD/fixing_some_stuff_08-12` when the owner is ready — batches E and F
   are both on it.
2. Start **batch G** (backlog §2.3, float preview geometry). It regenerates GPU
   snapshots and needs owner review, and must land **before** §5.1.
3. Then batch **H** (§5.1, §5.2): zone-editor pointer and property-panel tests,
   against the post-§2.3 coordinates. `FileExplorerHandler` and
   `appRunnerSemantics.go` are the pattern to follow — H should extend the same
   handler framework rather than driving dialogs directly, and
   `plans/gui-handler-framework.md` governs it.

## 10. Carry-forward prompt

> Read `AGENTS.md` first, then `todo/backlog-opus5.md`.
>
> Hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/` without explicit
> approval; everything must build and run on both Windows and Linux (use
> `path/filepath`, and chain PowerShell with `;`, never `&&`); every change
> ships with tests and unit coverage must not drop below 72.5 % (currently
> 72.9 %); durable multi-session work gets a plan file under `plans/`; never
> stage and never commit — I review, stage and commit myself, so leave the
> staging area alone entirely, and delete files with `Remove-Item`, never
> `git rm`; never change where `.rmg.json` is written and never persist the
> output directory; never run a bulk in-place rewrite over the repository;
> never run CI and never generate snapshot goldens in CI.
>
> Batch **F** (backlog §5.3, file-explorer hidden-file toggle and
> pointer-driven row/scroll tests) is **complete, green and committed** as
> `d37d1fa` on `AD/fixing_some_stuff_08-12`, unpushed. The file explorer is now
> driven entirely through `integration_common.FileExplorerHandler`, which
> resolves buttons and listing rows by accessibility label and injects real
> pointer events, with a golden per action; every `Click*` test-export on the
> dialog was deleted as a result. See `plans/batch-f-file-explorer.md`.
>
> Next up is batch **G** (§2.3, float preview geometry) — it regenerates GPU
> snapshots, needs my review, and must land before §5.1. When generating
> goldens, scope `-update` with `-run` so unrelated goldens in the package are
> not rewritten. Before starting, prompt me to confirm the item and surface
> every open question first.
>
> Full handoff: `./.agent/session-carry-forward.md`.
