# Extract filesystem policy out of the GUI (review §2.1 + §2.5)

Move every `os`/`filepath` decision out of `app/gui/dialogs/fileExplorerDialog.go`
into two services under `internal/services/file_system/`, reach them from the
dialog through a new standalone `IFileSystemHandler`, then split the remaining
dialog into sibling files by responsibility. Closes review findings **§2.1** and
**§2.5** (§12 item 13, first half).

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is
done, set its status to `Complete` and write its **Phase Summary** (what was
done, key decisions, anything needed to continue with zero context); run the
phase's **Verification Plan** and record the result before moving on. When all
phases are done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../AGENTS.md) first. Never stage, never commit — the owner
reviews and commits.

## Owner decisions (already made — do not re-litigate)

| Question | Decision |
| --- | --- |
| Batch scope | **§2.1 + §2.5 together**, one batch. |
| Seam | **Standalone `IFileSystemHandler`**, *not* embedded into `IGuiHandler`. |
| Wiring | **Add a 3rd parameter to `NewUIState`** (~25 call sites, nearly all tests). |
| Service granularity | **Two services** in `internal/services/file_system/`. |
| What else moves | `hasHiddenAttr` (build-tagged pair), `listWindowsDrives`, `parentDir`, `matchesFilter` — all four move. |
| Defects | **List first, fix clear-cut, ask before any user-visible change.** |
| Integration tests | Open flow, Save flow, Overwrite prompt, New-folder flow (New-folder included; hidden-file toggle excluded). |
| §2.5 split | **Four new sibling files** + the core: `...Modes.go`, `...Entries.go`, `...Toolbar.go`, `...Confirm.go`. |

## Starting facts

- `app/gui/dialogs/fileExplorerDialog.go`: **750 LOC, 33 funcs, 29 fields,
  20.4% covered**. Eight functions touch `os`/`filepath`: `confirmSelection`,
  `loadDir`, `isHidden`, `parentDir`, `resolveSaveTarget`, `tryCreateFolder`,
  `listWindowsDrives`, `resolveInitialDir`. All but `loadDir` (65.8%),
  `isHidden` (66.7%) and `resolveInitialDir` (47.4%) are at **0%**.
- §2.1's prerequisite (§1.1) is satisfied: `FolderPermission = 0o755` /
  `FilePermission = 0o644` already live in
  [internal/common/constants/permissions.go](../internal/common/constants/permissions.go)
  and the dialog already uses `FolderPermission`.
- `drivers.State` holds one handler field, `handler handler_interfaces.IGuiHandler`,
  set by `NewUIState(handler, findTemplateDir)`.
- depguard rule `no-services-from-app` forbids `app/**` importing
  `internal/services`, so the handler seam is mandatory, not stylistic.
- `internal/services/file_system/` does not exist yet.

---

## Phase 0: Baseline
Status: Complete

- [x] Record `go test -count=1 ./test/unit/...` result.
- [x] Record total coverage and the per-file numbers for
      `app/gui/dialogs/fileExplorerDialog.go` and `app/gui/drivers/stateFiles.go`.
- [x] Record `golangci-lint-v2 run ./... --issues-exit-code=0` issue count.
- [x] Run the GPU-gated GUI suite once so later snapshot diffs are attributable.

### Verification Plan
- All four numbers written into the Phase Summary below.

### Phase Summary

Baseline recorded 2026-08-06 on branch `AD/refactoring-07-21`, working tree
containing only the uncommitted doc edits from the previous batch.

| Metric | Baseline |
| --- | --- |
| `go build ./...` | clean |
| `go test -count=1 ./test/unit/...` | **0 failures** (no `FAIL`/`--- FAIL` lines) |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| Total unit coverage (`-coverpkg=./internal/...,./app/...`) | **68.7 %** |
| `app/gui/dialogs/fileExplorerDialog.go` | **59 / 289 statements = 20.4 %** |
| `app/gui/drivers/stateFiles.go` | **17 / 57 statements = 29.8 %** |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **0 issues.** |
| `go test -tags 'integration_test,gui' -count=1 ./test/integration/gui/...` | `ok ... 2.169s` |

Notes for later phases:

- The merged `coverage.txt` contains **duplicate blocks** (one set per test
  package under `-coverpkg`), so naive summing over-counts. Per-file numbers
  above were computed with [tmp/filecov.ps1](../tmp/filecov.ps1), which
  de-duplicates by block key and takes the max hit count. Re-run it as
  `./tmp/filecov.ps1 -Files 'dialogs/fileExplorerDialog.go','drivers/stateFiles.go'`
  after regenerating `coverage.txt`. `tmp/` is scratch — delete the script in
  Phase 7 if the owner does not want it kept.
- `golangci-lint-v2` still exits **1** while printing `0 issues.` (stderr
  warnings about exclusion rules). Judge by the printed count, not the exit code.
- The GUI snapshot suite is green *before* any change, so any snapshot diff
  appearing in Phase 3 or Phase 5 is caused by this batch and must be treated as
  a regression, not a stale baseline.

---

## Phase 1: The `internal/services/file_system/` package
Status: Complete

New package with **two** services (AGENTS.md §4.2.2: fewer than 5
implementations → interfaces live in the same package, `*Interface.go`).

- [x] `internal/models/directoryEntry.go` — `DirectoryEntry{Name, Path string; IsDir bool}`.
- [x] `internal/services/file_system/directoryBrowserService.go` +
      `directoryBrowserServiceInterface.go` — `IDirectoryBrowserService`:
      - `ListEntries(directory string, filterSuffixes []string, showHidden bool) ([]models.DirectoryEntry, error)`
        (folds in the old `isHidden` + `matchesFilter`, and the dirs-then-files
        case-insensitive sort)
      - `ListRoots() []models.DirectoryEntry` (was `listWindowsDrives`)
      - `CreateDirectory(parent, name string) (string, error)`
- [x] `internal/services/file_system/pathResolutionService.go` +
      `pathResolutionServiceInterface.go` — `IPathResolutionService`:
      - `ResolveStartDirectory(preferred string) string` (was `resolveInitialDir`)
      - `ParentDirectory(current string) string` (was `parentDir`)
      - `ResolveSaveTarget(directory, name, requiredSuffix string) (string, bool)`
      - `PathExists(path string) bool` / overwrite detection
- [x] `hiddenAttribute_windows.go` / `hiddenAttribute_other.go` — the moved
      build-tagged `hasHiddenAttr` pair, unchanged in behaviour.
- [x] Factories return the **interface** (AGENTS.md §4.2.2).
- [x] Full mirrored unit tests under
      `test/unit/internal/services/file_system/{directoryBrowserService,pathResolutionService}/`,
      one file per public method, all using `t.TempDir()` so they stay
      cross-platform.

Constraints: no Gio import anywhere in this package; `path/filepath` only, never
literal separators.

### Verification Plan
- `go build ./...` passes.
- `go test -count=1 ./test/unit/internal/services/file_system/...` passes.
- `go run ./cmd/testlayoutcheck .` prints `test-layout check passed`.
- Coverage of the two new files is ≥ 90%.

### Phase Summary

All four verification checks pass: `go build ./...` clean,
`go test -count=1 ./test/unit/internal/services/file_system/...` **ok** for both
packages, `test-layout check passed`, and package coverage **92.4 %** (well over
the 90 % bar).

**Files created**

| File | Note |
| --- | --- |
| [internal/models/directoryEntry.go](../internal/models/directoryEntry.go) | Pure data struct, no methods → no unit tests (AGENTS.md §4.6 scope rules). |
| [internal/common/common_errors/fileSystemErrors.go](../internal/common/common_errors/fileSystemErrors.go) | `ErrDirectoryNameEmpty`, `ErrDirectoryNameInvalid`. |
| [internal/services/file_system/directoryBrowserServiceInterface.go](../internal/services/file_system/directoryBrowserServiceInterface.go) | `IDirectoryBrowserService`. |
| [internal/services/file_system/directoryBrowserService.go](../internal/services/file_system/directoryBrowserService.go) | Also holds the now-package-private `isHidden`, `matchesFilter`, `compareByLowercaseName`. |
| [internal/services/file_system/pathResolutionServiceInterface.go](../internal/services/file_system/pathResolutionServiceInterface.go) | `IPathResolutionService`. |
| [internal/services/file_system/pathResolutionService.go](../internal/services/file_system/pathResolutionService.go) | |
| [internal/services/file_system/hiddenAttribute_windows.go](../internal/services/file_system/hiddenAttribute_windows.go) | `//go:build windows`, renamed `hasHiddenAttr` → `hasHiddenAttribute`. |
| [internal/services/file_system/hiddenAttribute_other.go](../internal/services/file_system/hiddenAttribute_other.go) | `//go:build !windows`. |
| `test/unit/internal/services/file_system/directoryBrowserService/` | `common_test.go` + 4 method files. |
| `test/unit/internal/services/file_system/pathResolutionService/` | 5 method files. |

**Key decisions taken during implementation**

1. **Error wording stays in the GUI.** `tryCreateFolder` used to set the literal
   strings `"Enter a folder name."` / `"Invalid folder name."`. Those are UI copy,
   and staticcheck ST1005 forbids capitalised, punctuated error strings, so the
   service returns the two sentinels from `common_errors` and Phase 3 maps them
   back to the same user-visible text. `common_errors` is chosen because
   `app/gui/drivers` already imports it, whereas depguard's `no-services-from-app`
   would block importing sentinels declared in `file_system`.
2. **`CreateDirectory` owns the trim and the name validation**, not just `Mkdir`.
   That keeps the whole "can this name escape the parent" question on one side of
   the seam.
3. **`ResolveSaveTarget` takes `requiredSuffix` as a parameter** rather than
   baking in `.gen.json`; the constant stays a GUI concern. An empty suffix means
   "use the name verbatim", which the tests pin.
4. **The suffix filter no longer needs a mode check.** The old `loadDir` applied
   `matchesFilter` only in `modeOpenFile`; since `filterSuffixes` is non-empty
   *only* in open mode, passing it unconditionally is behaviour-identical and one
   fewer piece of dialog state crossing the seam. Subdirectories are still never
   filtered.
5. **Both services are stateless structs** (`struct{}`), so a single instance can
   serve every dialog and wire can provide them as plain singletons.

**Behaviour deliberately preserved verbatim**, including the parts Phase 4 will
re-examine as defects: `isHidden` still fails open when `entry.Info()` errors,
and `ListEntries` still returns dot-prefixed entries when `showHidden` is set.

**Uncovered remainder (7.6 %)** is only unreachable-by-design fallbacks:
`os.Getwd()` failing, `entry.Info()` failing, and `info.Sys()` not being a
`*syscall.Win32FileAttributeData`. A test that sets `FILE_ATTRIBUTE_HIDDEN` would
need a `_windows_test.go` file, whose name would violate the §4.6
`<functionName>_test.go` rule, so it was not added.

**Lint gotcha hit again:** every new file tripped `gci`/`gofmt`/`golines` at
`1:1`. `golangci-lint-v2 run ./... --issues-exit-code=0 --fix` cleared them and
this time did **not** duplicate any `package` clause — re-verified with
`testlayoutcheck` plus a rebuild. Two real findings were fixed by hand first:
`unparam` on the two test helpers' unused `string` returns, and `godoclint`
wanting `[os.DirEntry]` in the `DirectoryEntry` doc comment.

Nothing outside the new files was touched, so `go test ./test/unit/...` as a
whole and the GUI snapshot suite are unaffected at this point.

---

## Phase 2: The `IFileSystemHandler` seam
Status: Complete

- [x] `internal/handlers/handler_interfaces/fileSystemHandlerInterface.go` —
      `IFileSystemHandler`, the union of the two service interfaces as consumed
      by the GUI. **Not** embedded into `IGuiHandler` (owner decision).
- [x] `internal/handlers/fileSystemHandler.go` — thin delegation, receiver
      `this`, constructor returns the interface.
- [x] Register both services and the handler in
      [internal/composition/providerSets.go](../internal/composition/providerSets.go)
      and add an `InitializeFileSystemHandler` injector to `wire.go`.
- [x] `wire gen ./internal/composition/...`, then `wire diff` to confirm.
      **Never** pass `-tags=wireinject` to build/test.
- [x] Check the `concrete-handlers-only-at-gui-composition-roots` depguard rule
      still holds for the new composition root.
- [x] `test/test_helpers/fileSystemHandlerMock.go`.
- [x] Mirrored unit tests under `test/unit/internal/handlers/fileSystemHandler/`.

### Verification Plan
- `wire diff ./internal/composition/...` exits 0.
- `go build ./...`, `go vet -tags=integration_test ./...` pass.
- `go test -count=1 ./test/unit/internal/handlers/...` passes.

### Phase Summary

All verification steps pass: `wire diff` exits **0**, `go build ./...` and
`go vet -tags=integration_test ./...` are clean, `go test ./test/unit/internal/handlers/...`
reports **0 failures**, `testlayoutcheck` passes and lint is back to **0 issues.**

**Files created**

| File | Note |
| --- | --- |
| [internal/handlers/handler_interfaces/fileSystemHandlerInterface.go](../internal/handlers/handler_interfaces/fileSystemHandlerInterface.go) | 7 methods; the flattened union of both service interfaces. |
| [internal/handlers/fileSystemHandler.go](../internal/handlers/fileSystemHandler.go) | Private struct `fileSystemHandler`, pure delegation, zero logic. |
| [test/test_helpers/fileSystemHandlerMock.go](../test/test_helpers/fileSystemHandlerMock.go) | For Phase 3, when `drivers.State` gains the dependency. |
| [test/test_helpers/directoryBrowserServiceMock.go](../test/test_helpers/directoryBrowserServiceMock.go) | |
| [test/test_helpers/pathResolutionServiceMock.go](../test/test_helpers/pathResolutionServiceMock.go) | |
| `test/unit/internal/handlers/fileSystemHandler/` | `common_test.go` + 7 method files. |

**Files modified**

- [internal/composition/providerSets.go](../internal/composition/providerSets.go) —
  added the `file_system` import and a new **`FileSystemSet`**.
- [internal/composition/wire.go](../internal/composition/wire.go) —
  added `InitializeFileSystemHandler()`.
- [internal/composition/wire_gen.go](../internal/composition/wire_gen.go) —
  regenerated; the new injector is a clean three-liner with no shared
  collaborators.

**Key decisions**

1. **A separate `FileSystemSet`, not an addition to `GuiHandlerSet`.** The two
   graphs share nothing, so keeping them disjoint means `InitializeFileSystemHandler`
   does not drag in the generator, the repositories or the validators. This is
   the structural payoff of the owner's "standalone seam" decision.
2. **The interface is flat, not embedded.** `IFileSystemHandler` restates the
   seven methods rather than embedding `file_system.IDirectoryBrowserService`
   and `IPathResolutionService`, because embedding them would force `app/gui` to
   import `internal/services` and trip depguard's `no-services-from-app`. The
   handler package is precisely the boundary that translation happens at.
3. **`fileSystemHandler` is unexported** and the constructor returns
   `handler_interfaces.IFileSystemHandler`, matching `previewHandler` and the
   rest of the package (AGENTS.md §4.2.2 factory rule).
4. **Two extra service mocks were needed.** The handler is tested against mocked
   services rather than the real filesystem, so `DirectoryBrowserServiceMock` and
   `PathResolutionServiceMock` joined the 14 mocks Batch 11 introduced (now 17).

**depguard check:** the new composition root is reached the same way as the
existing one — `app/gui` will call `composition.InitializeFileSystemHandler()` in
Phase 3 and never name `internal/handlers` directly, so
`concrete-handlers-only-at-gui-composition-roots` continues to hold. Confirmed by
lint still reporting `0 issues.`

**Reminder that held again:** `wire gen` prints its success banner to STDERR and
PowerShell surfaces that as a `NativeCommandError`; the run succeeded. Judge by
`wire diff` (exit 0), never by `wire gen`'s exit code.

Nothing in `app/` was touched yet, so the GUI still constructs dialogs exactly as
before and the snapshot suite is still untouched at this point.

---

## Phase 3: Rewire the dialog and the driver
Status: Complete

- [x] `NewUIState(guiHandler, fileSystemHandler, findTemplateDir)` — new field on
      `drivers.State`.
- [x] `editor.NewWindow` takes both handlers; `app/gui/program.go` supplies
      `composition.InitializeGuiHandler()` and
      `composition.InitializeFileSystemHandler()`.
- [x] The four dialog constructors take the handler as their first parameter.
- [x] Delete `os` and `path/filepath` from
      `app/gui/dialogs/fileExplorerDialog.go` entirely; delete
      `fileExplorerHidden_windows.go` / `_other.go`.
- [x] Replace the private `fileEntry` struct with `models.DirectoryEntry`
      (`app` → `internal/models` is an allowed direction).
- [x] Route `stateFiles.go`'s direct `os.Getwd()` calls through
      `ResolveStartDirectory` so the driver stops doing filesystem policy too.
- [x] Update all ~25 `NewUIState` call sites using the new mock.

### Verification Plan
- `grep` finds **zero** `"os"` / `"path/filepath"` imports under
  `app/gui/dialogs/` and `app/gui/drivers/`.
- Full unit + integration + GPU-gated GUI suites pass with **zero snapshot diffs**
  (this phase is strictly behaviour-preserving; defects are Phase 4).

### Phase Summary

The GUI no longer touches the filesystem. `app/gui/dialogs/` and
`app/gui/drivers/` import neither `os` nor `path/filepath` (the only surviving
`"os"` under `app/gui/` is `program.go`'s process exit, which is out of scope).

**Production changes**

- `app/gui/dialogs/fileExplorerDialog.go` — holds
  `fileSystem handler_interfaces.IFileSystemHandler` as its first field; the
  private `fileEntry` struct is gone in favour of `models.DirectoryEntry`; the
  four public constructors and `newFileExplorerDialog` take the handler first;
  `loadDir`, `parentDir`, `onEntryClicked`, `resolveSaveTarget`,
  `tryCreateFolder`, `confirmSelection` and `getEntryRowWidget` all delegate.
  `isHidden`, `matchesFilter`, `listWindowsDrives` and `resolveInitialDir` were
  deleted (their logic now lives in `internal/services/file_system/`).
  `tryCreateFolder` maps `common_errors.ErrDirectoryNameEmpty` /
  `ErrDirectoryNameInvalid` back to the existing UI copy via `errors.Is`, so the
  user-visible messages are byte-identical to before.
- Deleted `app/gui/dialogs/fileExplorerHidden_windows.go` and
  `fileExplorerHidden_other.go`.
- `app/gui/drivers/state.go` — new `fileSystem` field;
  `NewUIState(handler, fileSystem, findTemplateDir)`.
- `app/gui/drivers/stateFiles.go` — `Load`, `SaveAs`, `PickOutputDir` and
  `RevealOutputDir` pass the handler to the dialogs and start from
  `this.workingDirectory()`.
- `app/gui/editor/window.go` / `app/gui/program.go` — the second handler is
  threaded from `composition.InitializeFileSystemHandler()`.

**Behaviour changes to disclose** (both user-invisible in any reachable
scenario, but recorded because they are not literal no-ops):

1. `suggestDirectory()` was **deleted**. It was reachable only after `os.Getwd()`
   already failed, at which point its own final `os.Getwd()` failed too and it
   returned `""`. It is replaced by `workingDirectory()` =
   `ResolveStartDirectory(".")`, which always yields an existing directory.
2. `state.go`'s `os.Getwd()` fallback likewise became `state.workingDirectory()`.

**Test changes** — 24 `NewUIState` call sites and 1 `editor.NewWindow` call site
updated. Rather than force every driver test to script a mock, a new
`test/test_helpers/fileSystemHandler.go` exposes `NewFileSystemHandler()`, which
builds the handler over the **real** services so the driver behaves exactly as in
production; `FileSystemHandlerMock` stays reserved for tests that assert on the
interaction itself. Integration tests use
`composition.InitializeFileSystemHandler()`.

**Verification**

| Check | Result |
| --- | --- |
| `go build ./...` | clean |
| `go vet -tags='integration_test,gui' ./...` | clean |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test -count=1 ./test/unit/...` | 0 failures |
| `go test -tags=integration_test ./test/integration/...` | `ok ... 2.523s` |
| `go test -tags 'integration_test,gui' ./test/integration/gui/...` | `ok ... 2.435s`, **zero snapshot diffs** |
| `golangci-lint-v2 run ./...` | `0 issues.` |
| Total unit coverage | **69.3 %** (baseline 68.7 %) |

One lint fix was needed: godoclint required `[os.Getwd]` link syntax in the new
`workingDirectory` doc comment.

Per-file coverage moved as expected — `fileExplorerDialog.go` fell from
59/289 = 20.4 % to 29/217 = 13.4 % because the *covered* filesystem logic left
the file for the 92.4 %-covered `file_system` package, leaving behind almost
nothing but Gio widget code that only the integration suite can reach. Total
coverage rose regardless. Phase 6's GUI scenarios are what will lift this file.
`stateFiles.go` went 17/57 = 29.8 % → 16/48 = 33.3 %.

---

## Phase 4: Defects
Status: Not started

Owner decision: **list first, fix the clear-cut ones with tests, ask before any
user-visible change.**

- [ ] Enumerate every defect found during Phases 1–3 in the Phase Summary, each
      with: evidence, whether it is reachable today, and the proposed fix.
- [ ] Fix the clear-cut correctness bugs, each with a regression test.
- [ ] Ask the owner about anything that changes what the user sees.

Known candidates to assess (not yet confirmed as defects):
1. `confirmSelection` uses `os.Stat` and treats an existing **directory** as an
   overwritable file — saving would then target a directory path.
2. `resolveSaveTarget` rejects `.`/`..`/separators but not Windows reserved
   device names (`CON`, `NUL`, `COM1`…) or trailing dots/spaces.
3. `tryCreateFolder` with `currentDir == ""` would `filepath.Join("", name)` and
   create the folder in the process CWD — currently unreachable via `canModify()`.
4. `isHidden` swallows the `entry.Info()` error and fails open.
5. `loadDir` adopts an unreadable directory as `currentDir` on first load.

### Verification Plan
- Every accepted fix has a failing-before/passing-after unit test.
- Anything rejected is recorded here with the reason.

### Phase Summary
_(write when phase completes)_

---

## Phase 5: Split the dialog (§2.5)
Status: Not started

Sibling files for one struct, per AGENTS.md §4.1. Core file keeps the struct,
`newFileExplorerDialog`, `Title`, `PreferredSize`, `Body`, `getContentWidget`
and the small state helpers.

- [ ] `fileExplorerDialogModes.go` — `fileDialogMode` enum + the four public
      constructors.
- [ ] `fileExplorerDialogEntries.go` — `getListWidget`, `getEntryRowWidget`,
      `getErrorLineWidget`, and the now-thin `loadDir`.
- [ ] `fileExplorerDialogToolbar.go` — `getHeaderWidget`, `getSaveRowWidget`,
      `getNewFolderRowWidget`, `tryCreateFolder`.
- [ ] `fileExplorerDialogConfirm.go` — `getFooterWidget`, `confirmButtonState`,
      `handleConfirm`, `confirmOverwrite`, `confirmSelection`.
- [ ] Pure moves only — no logic changes in this phase, so a snapshot diff means
      a mistake.

### Verification Plan
- No file in `app/gui/dialogs/` exceeds ~250 LOC.
- GPU-gated GUI suite passes with zero snapshot diffs.
- `golangci-lint-v2` issue count unchanged (watch `funcorder`, `gci`, `golines`).

### Phase Summary
_(write when phase completes)_

---

## Phase 6: File-explorer integration scenarios
Status: Not started

`test/integration/gui/`, tags `//go:build integration_test && gui`. This closes
the "no file-explorer scenario" gap recorded in
[todo/test_observations.md](../todo/test_observations.md).

- [ ] Open flow — pick a `.gen.json` through the dialog, assert the state loads.
- [ ] Save flow — type a name, confirm, assert the file lands on disk.
- [ ] Overwrite prompt — save over an existing file, assert the prompt gates the
      write and that cancel leaves the file untouched.
- [ ] New-folder flow — create a folder, assert it appears in the listing.
- [ ] Update the `fileExplorerDialog.go` entry in
      [todo/test_observations.md](../todo/test_observations.md) to reflect what
      is now covered.

Use `t.TempDir()` for every fixture; never touch the real output directory.

### Verification Plan
- `go test -tags='integration_test,gui' -count=1 ./test/integration/gui/...` passes.
- `go run ./cmd/testlayoutcheck .` passes (the `gui` tag is mandatory here).

### Phase Summary
_(write when phase completes)_

---

## Phase 7: Close out
Status: Not started

- [ ] Full suite: build, both `go vet` tag combinations, testlayoutcheck, unit,
      integration, GPU-gated GUI, coverage, lint.
- [ ] Coverage ≥ the 68.7% baseline; lint ≤ the baseline issue count.
- [ ] Mark §2.1 and §2.5 `✅ FIXED` **in place** in
      [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) and update §12
      item 13.
- [ ] Update repository memory.
- [ ] Rewrite `.agent/session-carry-forward.md`.
- [ ] Stop for owner review. Do not stage. Do not commit.

### Verification Plan
- Every command in AGENTS.md §7 Quick Reference passes.

### Phase Summary
_(write when phase completes)_

---

## Final Recap
_(write when all phases complete)_

## Deployment Plan
_(write when all phases complete)_
