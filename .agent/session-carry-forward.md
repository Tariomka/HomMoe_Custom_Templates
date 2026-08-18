# Session Carry-Forward

## 1. Session goal

Implement backlog **batch E** — `todo/backlog-opus5.md` §4.1, *"Save As" is really
"Save To": the UI offers a filename it then discards*.

## 2. Fixes applied

- **The dialog no longer offers a filename it cannot honour.** The save row in
  [fileExplorerDialogToolbar.go](../app/gui/dialogs/fileExplorerDialogToolbar.go)
  is now a **read-only** textbox labelled `"Will save as:"`, populated with the
  resolved name instead of a `filename.gen.json` hint. The typed value used to be
  silently dropped by `FileService.SaveSettings`, which always writes
  `{TemplateName}.gen.json`.
- **A blank template name no longer produces a second lie.**
  [stateFiles.go](../app/gui/drivers/stateFiles.go) previously appended
  `.gen.json` unconditionally, so an unnamed template prefilled `".gen.json"`
  while `atomicFileWriter.resolveFileName` would actually write
  `Generated_Template.gen.json`. `SaveTo` now appends the suffix **only** when
  the sanitized, trimmed name is non-empty, so an unnamed template resolves to
  no filename at all.
- **Clicking a file row no longer silently retargets the save.**
  `onEntryClicked`'s `modeSaveFile` branch in
  [fileExplorerDialog.go](../app/gui/dialogs/fileExplorerDialog.go) copied the
  clicked entry into the filename field. With a read-only preview that would
  change the target with no way to undo it, so the branch was removed. Row
  clicks now mean something in **open mode only**.

## 3. Features added / changed

- **`hasResolvedSaveName()` guard** on `FileExplorerDialog`. The owner
  explicitly rejected reusing the existing `resolveSaveTarget` `!ok` predicate
  for this — the "nothing to save under" case deserves its own, readable
  predicate. It drives two things:
  - `confirmButtonState` disables the confirm button, so the user cannot press
    Save and get a file named after nothing.
  - `getSaveRowWidget` renders `missingSaveNameMessage`
    (`"Template name is required."`) inline under the row in
    `themes.ColorsBase.Error`, mirroring the new-folder row's error pattern.
    The owner chose an inline message over a hint or a bare empty field, so the
    disabled button has a stated reason.
- **Vocabulary made honest**: toolbar button `"Save As"` → `"Save To"`
  ([toolbar.go](../app/gui/editor/toolbar.go), field `buttonSaveAs` →
  `buttonSaveTo`), dialog title `"Save File"` → `"Save To"`, driver method
  `State.SaveAs` → `State.SaveTo`. **`NewSaveFileDialog`, `modeSaveFile` and
  `onSave` deliberately keep their names** (owner decision): they name the
  *explorer mode and callback*, not the toolbar action.

## 4. File modifications

Production:

| File | Change |
| --- | --- |
| [app/gui/dialogs/fileExplorerDialog.go](../app/gui/dialogs/fileExplorerDialog.go) | Added `missingSaveNameMessage` const and `hasResolvedSaveName()`; removed the `modeSaveFile` branch from `onEntryClicked`. |
| [app/gui/dialogs/fileExplorerDialogToolbar.go](../app/gui/dialogs/fileExplorerDialogToolbar.go) | `getSaveRowWidget` is a vertical Flex: read-only textbox + conditional inline error. Dropped the `fmt` import and the hint. |
| [app/gui/dialogs/fileExplorerDialogConfirm.go](../app/gui/dialogs/fileExplorerDialogConfirm.go) | `confirmButtonState` disables confirm when no name resolved. |
| [app/gui/dialogs/fileExplorerDialogModes.go](../app/gui/dialogs/fileExplorerDialogModes.go) | `defaultName` → `resolvedName`, title `"Save To"`, doc comment rewritten. |
| [app/gui/drivers/stateFiles.go](../app/gui/drivers/stateFiles.go) | `SaveAs` → `SaveTo` with the conditional-suffix fix; `Save`'s fallback updated. |
| [app/gui/editor/toolbar.go](../app/gui/editor/toolbar.go) | `buttonSaveAs` → `buttonSaveTo`, label `"Save To"`, calls `SaveTo`. |
| [app/gui/dialogs/fileExplorerDialog_testexports.go](../app/gui/dialogs/fileExplorerDialog_testexports.go) | `SetFilename` **removed**; added `ResolvedSaveName()` and `SaveNameReadOnly()`. |

Test helpers:

| File | Change |
| --- | --- |
| [test/test_helpers/integration_common/baseHandler.go](../test/test_helpers/integration_common/baseHandler.go) | `ClickSaveAs()` → `ClickSaveTo()`. |
| [test/test_helpers/integration_common/handlerCoordinates.go](../test/test_helpers/integration_common/handlerCoordinates.go) | `saveAsButtonX` → `saveToButtonX` (value 200 unchanged — x=200 is inside the button under either label). |
| 10 `.golden` files under `snapshot/__snapshots__/window_snapshot_integration_test/` | Regenerated locally for the new button label. |

Docs:

| File | Change |
| --- | --- |
| [QUICKSTART.md](../QUICKSTART.md) | Toolbar list now says `Save To`; the save/load section explains that the filename is derived and that an unnamed template cannot be saved. |
| [README.md](../README.md) | Step 4 now says **Save To** picks the folder only. |
| [todo/test_observations.md](../todo/test_observations.md) | `stateFiles.go` entries renamed to `SaveTo` / `stateSaveTo_integration_test.go`; recorded the now-unreachable whitespace branch. |
| [todo/backlog-opus5.md](../todo/backlog-opus5.md) | §4.1 marked ✅ DONE and rewritten to be **self-contained** (original finding folded into a `<details>` block); §8 batch table row E marked done, row F annotated with the row-click change. |

## 5. Tests added or updated

Renamed (via `Move-Item`, never `git mv` — that stages):

- `test/integration/stateSaveAs_integration_test.go` → [stateSaveTo_integration_test.go](../test/integration/stateSaveTo_integration_test.go);
  `newSaveAsProbe` → `newSaveToProbe`, both test names updated.
- `test/unit/app/gui/drivers/stateFiles/saveAs_test.go` → [saveTo_test.go](../test/unit/app/gui/drivers/stateFiles/saveTo_test.go).

Added:

- `TestWhenSaveToOpens_TheDialogPreviewsTheResolvedFilename` — drives the whole
  chain (template name → `SanitizeFilename` → suffix → dialog field) and asserts
  `"  Jebus: Cross  "` previews as `"Jebus_ Cross.gen.json"`.
- `TestWhenTheTemplateIsUnnamed_SaveToPreviewsNoFilename` — a whitespace-only
  template name previews nothing.
- `TestWhenTheSaveDialogIsLaidOut_TheNameFieldIsReadOnly`.
- `TestWhenNoNameWasResolved_TheConfirmButtonIsDisabled`.
- `TestWhenAFileRowIsClickedInSaveMode_TheResolvedNameIsUnchanged` — pins the
  removed `onEntryClicked` branch.

Rewritten / renamed:

- `..._TheTypedNameBecomesTheSaveTarget` → `..._TheResolvedNameBecomesTheSaveTarget`,
  and the driver-level save test — both now pass the name through the
  constructor instead of `SetFilename`.
- `TestWhenTheFilenameIsValid_TheConfirmButtonIsEnabled` →
  `TestWhenANameWasResolved_TheConfirmButtonIsEnabled`.
- `TestWhenNoCurrentPathExists_SaveOpensSaveAsDialog` → `...SaveOpensSaveToDialog`.
- `TestHandlerDialogs_SaveAsOpensFileExplorer` → `..._SaveToOpensFileExplorer`.

Deleted:

- `TestWhenTheFilenameIsWhitespaceOnly_TheConfirmButtonIsDisabled` — that state
  is no longer reachable through any production path (only by calling
  `NewSaveFileDialog` directly), per the backlog's explicit conditional
  instruction. Recorded in `todo/test_observations.md`.

**Gate results — all green:**

| Gate | Result |
| --- | --- |
| `go build ./...` | clean |
| `go vet -tags='integration_test,gui' ./...` | clean |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test ./test/unit/... -count=1` | pass |
| `go test -tags=integration_test ./test/integration/... -count=1` | pass |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | pass |
| Unit coverage | **72.9 %** (unchanged; floor is 72.5 %) |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **0 issues** |
| `gofmt -l ./app ./internal ./test ./cmd` | empty |
| `wire diff ./internal/composition/...` | no diff |

## 6. Git status snapshot

Branch: `AD/fixing_some_stuff_08-12`. **Nothing staged, nothing committed** —
per AGENTS.md §2.5 the owner reviews and commits.

```
 M QUICKSTART.md
 M README.md
 M app/gui/dialogs/fileExplorerDialog.go
 M app/gui/dialogs/fileExplorerDialogConfirm.go
 M app/gui/dialogs/fileExplorerDialogModes.go
 M app/gui/dialogs/fileExplorerDialogToolbar.go
 M app/gui/dialogs/fileExplorerDialog_testexports.go
 M app/gui/drivers/stateFiles.go
 M app/gui/editor/toolbar.go
 M test/integration/gui/fileExplorerDialog_integration_test.go
 M test/integration/handlerDialogReachability_integration_test.go
 D test/integration/stateSaveAs_integration_test.go
 M test/test_helpers/integration_common/baseHandler.go
 M test/test_helpers/integration_common/handlerCoordinates.go
 M ...10 window_snapshot_integration_test .golden files
 D test/unit/app/gui/drivers/stateFiles/saveAs_test.go
 M test/unit/app/gui/drivers/stateFiles/save_test.go
 M todo/backlog-opus5.md
 M todo/test_observations.md
?? test/integration/stateSaveTo_integration_test.go
?? test/unit/app/gui/drivers/stateFiles/saveTo_test.go
```

The two `D` + two `??` pairs are the renames. All of it is inherited unstaged
work for the next session if the owner has not committed yet.

## 7. Rejections / things the owner declined

- **Reusing the `resolveSaveTarget` `!ok` predicate** to disable the confirm
  button. Recommended, but declined in favour of an explicit
  `hasResolvedSaveName()` guard in the dialog.
- **A hint, or a silently empty field**, for the unnamed-template case. The
  owner asked for an inline message under the row instead.
- **Renaming `NewSaveFileDialog` / `modeSaveFile` / `onSave`.** Kept — they name
  the explorer mode and callback, not the toolbar action.
- **A plan file under `plans/`** was waived for this batch.

## 8. Open questions

- **Eyeball the 10 regenerated goldens before committing.** They passed both
  before and after regeneration — the two-glyph label change fell under the
  changed-pixel gate — so they were regenerated to keep the reference truthful,
  not because a test failed. A visual check confirms nothing else moved.
- Nothing else is blocked.

## 9. Next recommended actions

1. Owner reviews and commits batch E.
2. Start **batch F** (backlog §5.3): file-explorer hidden-file toggle and
   pointer-driven row/scroll tests —
   `TestWhenShowHiddenIsToggledOn_HiddenEntriesAppearInTheListing`,
   `TestWhenShowHiddenIsToggledOff_HiddenEntriesDisappearAgain`,
   `TestWhenARowIsClicked_ThatEntryBecomesTheSelection`,
   `TestWhenADirectoryRowIsClicked_TheListingDescendsIntoIt`, plus a
   build-tagged cross-platform hidden-file fixture helper in
   `test/test_helpers/`.
   ⚠ **Batch E changed F's scope**: the save-mode row-click behaviour is gone,
   so `TestWhenARowIsClicked_ThatEntryBecomesTheSelection` applies to **open
   mode only**. Confirm with the owner before starting.
3. Then batch **G** (§2.3, float preview geometry) — it regenerates GPU
   snapshots and needs owner review, and must land **before** §5.1.

## 10. Carry-forward prompt

> Read `AGENTS.md` first, then `todo/backlog-opus5.md`.
>
> Hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/` without explicit
> approval; everything must build and run on both Windows and Linux (use
> `path/filepath`, and chain PowerShell with `;`, never `&&`); every change
> ships with tests and unit coverage must not drop below 72.5 % (currently
> 72.9 %); durable multi-session work gets a plan file under `plans/`; never
> stage and never commit — the owner reviews and commits, and files are deleted
> with `Remove-Item`, never `git rm`; never change where `.rmg.json` is written
> and never persist the output directory; never run a bulk in-place rewrite over
> the repository; never run CI and never generate snapshot goldens in CI.
>
> Batch **E** (backlog §4.1, "Save As" → "Save To") is **complete and green**,
> unstaged on `AD/fixing_some_stuff_08-12`. The save dialog now shows a
> read-only, derived filename, refuses to save an unnamed template with an
> inline reason, and no longer lets a row click retarget the save. All gates
> pass; the 10 window goldens were regenerated locally for the new button label
> and want a visual check.
>
> Next up is batch **F** (§5.3, file-explorer hidden-file toggle and
> pointer-driven row/scroll tests). Note that batch E removed the save-mode
> row-click behaviour, so F's row-click selection test now applies to open mode
> only. Before starting, prompt me to confirm the item and surface every open
> question first.
>
> Full handoff: `./.agent/session-carry-forward.md`.
