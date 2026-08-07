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
Status: Complete

Owner decision: **list first, fix the clear-cut ones with tests, ask before any
user-visible change.** Follow-up decision: **implement all four (D1-D4).**

- [x] Enumerate every defect found during Phases 1–3 in the Phase Summary, each
      with: evidence, whether it is reachable today, and the proposed fix.
- [x] Fix the clear-cut correctness bugs, each with a regression test.
- [x] Ask the owner about anything that changes what the user sees.

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

Six findings. One was fixed without asking (no user-visible effect), four were
put to the owner and all four approved, and two were rejected as deliberate.

**D0 — `ListRoots` probed drive letters on every platform.** Fixed without
asking. `A:\` … `Z:\` are legal *relative file names* on Unix, so a stray file
could have been presented as a volume; the loop also cost 26 pointless `stat`
syscalls per call off Windows. `ListRoots` now returns `nil` unless
`runtime.GOOS == windowsOS`. Identical behaviour on Windows.

**D1 — an existing *folder* at the save target was offered as "overwrite".**
Reachable: `ResolveSaveTarget` appends `.gen.json`, and New Folder accepts
`foo.gen.json` as a folder name. `PathExists` answered true for it, the overwrite
prompt appeared, and confirming handed a directory path to `onSave`, which fails
at write time with a raw OS error. Fixed by adding `DirectoryExists` to
`IPathResolutionService` / `IFileSystemHandler` and checking it *before*
`PathExists`; the dialog now shows "A folder with that name already exists." and
refuses instead of prompting.

**D2 — a whitespace-only filename left *Save* enabled but did nothing.**
`confirmButtonState` tested `len(text) == 0` while `ResolveSaveTarget` trims and
reports `ok == false`, so the button looked live and clicking it was a no-op.
Fixed by deriving the disabled state from `resolveSaveTarget()` itself, which
makes the button's enablement and the save's success the *same* predicate rather
than two rules that can disagree. This also fixes D3 in the UI for free.

**D3 — Windows reserved device names were accepted.** `CON`, `PRN`, `AUX`,
`NUL`, `COM1`–`COM9`, `LPT1`–`LPT9` still resolve to the device whatever
extension follows, so "saving" a template as `NUL` reported success and left
nothing on disk. `ResolveSaveTarget` now rejects them (checking the stem before
the first dot, case-insensitively) when `runtime.GOOS == windowsOS`; the name is
not restricted elsewhere, because it is a legal file name there.

**D4 — `CreateDirectory(parent: "", …)` created the folder in the process
working directory.** Unreachable from the UI (`canModify()` requires
`currentDir != ""`), but `IDirectoryBrowserService` is a public seam now. Fixed
with a new `common_errors.ErrDirectoryParentEmpty` sentinel; the dialog's default
branch surfaces it verbatim, which is correct because it can only ever indicate a
programming error.

**Rejected — `isHidden` fails open when `entry.Info()` errors.** Deliberate: an
entry whose metadata cannot be read stays *visible* rather than silently
vanishing from the listing. Failing closed would hide files from the user with no
explanation. Do not "fix" this.

**Rejected — `loadDir` adopts an unreadable directory as `currentDir` on the very
first load.** Deliberate, so the path bar names the place the error refers to;
navigation up and the inline error both still work. It is also close to
unreachable now that `ResolveStartDirectory` guarantees an existing directory,
leaving only the permission-denied case.

**New GUI state**: `FileExplorerDialog.saveErr` holds the D1/filename rejection.
`getErrorLineWidget` prefers `listErr` and falls back to `saveErr`, so the layout
is unchanged; `loadDir` and `onEntryClicked` clear it, and `confirmSelection`
clears it on every click, so it can never go stale.

**Files touched**: `internal/common/common_errors/fileSystemErrors.go`,
`internal/services/file_system/{constants.go (new), directoryBrowserService.go,
pathResolutionService.go, pathResolutionServiceInterface.go}`,
`internal/handlers/{fileSystemHandler.go, handler_interfaces/fileSystemHandlerInterface.go}`,
`test/test_helpers/{fileSystemHandlerMock.go, pathResolutionServiceMock.go}`,
`app/gui/dialogs/fileExplorerDialog.go`, plus 5 unit-test files (2 new).

The unexported `windowsOS` constant was extracted because `goconst` flagged the
third literal `"windows"`; the reserved-name set is a `switch` rather than a map
because `gochecknoglobals` forbids the package-level variable.

**Verification**

| Check | Result |
| --- | --- |
| `go build ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test -count=1 ./test/unit/...` | 0 failures |
| `go test -tags=integration_test ./test/integration/...` | `ok ... 2.638s` |
| `go test -tags 'integration_test,gui' ./test/integration/gui/...` | `ok ... 2.375s`, zero snapshot diffs |
| `golangci-lint-v2 run ./...` | `0 issues.` |
| Total unit coverage | **69.3 %** (baseline 68.7 %) |

D1's dialog branch and D2's button state are *not* unit-covered — both need a
`layout.Context`. They are Phase 6's Save-flow and Overwrite-prompt scenarios.

---

## Phase 5: Split the dialog (§2.5)
Status: Complete

Sibling files for one struct, per AGENTS.md §4.1. Core file keeps the struct,
`newFileExplorerDialog`, `Title`, `PreferredSize`, `Body`, `getContentWidget`
and the small state helpers.

- [x] `fileExplorerDialogModes.go` — `fileDialogMode` enum + the four public
      constructors.
- [x] `fileExplorerDialogEntries.go` — `getListWidget`, `getEntryRowWidget`,
      `getErrorLineWidget`, and the now-thin `loadDir`.
- [x] `fileExplorerDialogToolbar.go` — `getHeaderWidget`, `getSaveRowWidget`,
      `getNewFolderRowWidget`, `tryCreateFolder`.
- [x] `fileExplorerDialogConfirm.go` — `getFooterWidget`, `confirmButtonState`,
      `handleConfirm`, `confirmOverwrite`, `confirmSelection`.
- [x] Pure moves only — no logic changes in this phase, so a snapshot diff means
      a mistake.

### Verification Plan
- No file in `app/gui/dialogs/` exceeds ~250 LOC.
- GPU-gated GUI suite passes with zero snapshot diffs.
- `golangci-lint-v2` issue count unchanged (watch `funcorder`, `gci`, `golines`).

### Phase Summary

A pure move. Not one statement of logic changed — every function was relocated
verbatim, which is exactly why a snapshot diff would have proved a mistake, and
there were none.

| File | LOC | Holds |
| --- | --- | --- |
| `fileExplorerDialog.go` | 221 | struct, `saveFileSuffix`, `newFileExplorerDialog`, `Title`, `PreferredSize`, `Body`, `getContentWidget`, `parentDir`, `requestNav`, `onEntryClicked`, `resolveSaveTarget`, `canModify`, `clickFor`, `resetScroll` |
| `fileExplorerDialogConfirm.go` | 156 | `getFooterWidget`, `confirmButtonState`, `handleConfirm`, `confirmOverwrite`, `confirmSelection` |
| `fileExplorerDialogEntries.go` | 130 | `loadDir`, `getListWidget`, `getEntryRowWidget`, `getErrorLineWidget` |
| `fileExplorerDialogToolbar.go` | 103 | `getHeaderWidget`, `getSaveRowWidget`, `getNewFolderRowWidget`, `tryCreateFolder` |
| `fileExplorerDialogModes.go` | 65 | `fileDialogMode` enum, `NewOpenFileDialog`, `NewSaveFileDialog`, `NewPickFolderDialog`, `NewBrowseDialog` |

750 LOC in one file became 675 across five, none over 250. The import lists tell
the same story the split does: the core file now needs only `layout`, `unit`,
`widget`, `material`, `widgets`, `handler_interfaces` and `models` — the drawing
primitives (`op`, `clip`, `paint`, `font`, `image`) and the error plumbing
(`errors`, `common_errors`) each moved to exactly one sibling.

`saveFileSuffix` stayed in the core file because both `resolveSaveTarget` (core)
and `getSaveRowWidget` (toolbar) read it.

**Verification**

| Check | Result |
| --- | --- |
| `go build ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| Largest file-explorer file | 221 LOC (target: ≤ 250) |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test -count=1 ./test/unit/...` | 0 failures |
| `go test -tags=integration_test ./test/integration/...` | `ok ... 3.607s` |
| `go test -tags 'integration_test,gui' ./test/integration/gui/...` | `ok ... 3.404s`, **zero snapshot diffs** |
| `golangci-lint-v2 run ./...` | `0 issues.` (unchanged; `funcorder` did not fire) |

---

## Phase 6: File-explorer integration scenarios
Status: Complete

`test/integration/gui/`, tags `//go:build integration_test && gui`. This closes
the "no file-explorer scenario" gap recorded in
[todo/test_observations.md](../todo/test_observations.md).

- [x] Open flow — pick a `.gen.json` through the dialog, assert the state loads.
- [x] Save flow — type a name, confirm, assert the file lands on disk.
- [x] Overwrite prompt — save over an existing file, assert the prompt gates the
      write and that cancel leaves the file untouched.
- [x] New-folder flow — create a folder, assert it appears in the listing.
- [x] Update the `fileExplorerDialog.go` entry in
      [todo/test_observations.md](../todo/test_observations.md) to reflect what
      is now covered.

Use `t.TempDir()` for every fixture; never touch the real output directory.

### Verification Plan
- `go test -tags='integration_test,gui' -count=1 ./test/integration/gui/...` passes.
- `go run ./cmd/testlayoutcheck .` passes (the `gui` tag is mandatory here).

### Phase Summary

Ten scenarios in
[test/integration/gui/fileExplorerDialog_integration_test.go](../test/integration/gui/fileExplorerDialog_integration_test.go),
all parallel, all on `t.TempDir()`.

**How they drive the dialog.** Clicking by pixel coordinate would have been
brittle and would have re-tested Gio rather than the dialog, so the tests queue
clicks with Gio's own `widget.Clickable.Click()` and then lay out a real frame.
`Clickable.update` consumes `requestClicks` before it looks at pointer input, so
the click is indistinguishable from a real one to every branch under test, while
the layout still runs for real (material.List, the editors, the footer). One
helper, `frameFileExplorer`, lays out a single frame exactly as `DialogHost`
does per vsync and returns the dialog's `done` flag.

The clickables and editors are unexported, so the accessors went into the
existing
[app/gui/dialogs/fileExplorerDialog_testexports.go](../app/gui/dialogs/fileExplorerDialog_testexports.go)
(`//go:build integration_test`) — the mechanism AGENTS.md §4.6.1 exists for.
Added: `ClickEntry`, `ClickUp`, `ClickConfirm`, `ClickOverwriteConfirm`,
`ClickOverwriteCancel`, `ClickNewFolder`, `ClickCreateFolder`, `SetFilename`,
`SetNewFolderName`, `CurrentDir`, `EntryNames`, `SelectedPath`,
`OverwriteActive`, `SaveError`, `NewFolderError`, `ConfirmDisabled`.

**The scenarios.**

| Test | Proves |
| --- | --- |
| `...OpenDialogConfirmsAFile_TheEditorStateIsLoaded` | list → row click → confirm → the pick handler installs the state (asserted on the loaded `TemplateName`) |
| `...SaveDialogIsConfirmed_TheTypedNameBecomesTheSaveTarget` | the dialog's own contract: directory + typed name + enforced suffix |
| `...SaveDialogIsConfirmedThroughTheDriver_AFileLandsInTheChosenDirectory` | the same confirm through the real `drivers.State`, handler and repository — bytes on disk |
| `...SaveTargetAlreadyExists_ConfirmDoesNotWrite` | the prompt is raised *instead of* writing, not after it |
| `...OverwriteIsCancelled_TheExistingFileIsUntouched` | cancel leaves the original bytes |
| `...OverwriteIsConfirmed_TheFileIsRewritten` | confirm actually writes |
| `...ANewFolderIsCreated_ItAppearsInTheParentListing` | toggle → name → create → Up → listed |
| `...SaveTargetIsAnExistingFolder_TheSaveIsRefused` | **D1's dialog branch** |
| `...TheFilenameIsWhitespaceOnly_TheConfirmButtonIsDisabled` | **D2's button state** |
| `...TheFilenameIsValid_TheConfirmButtonIsEnabled` | the same predicate the other way, so the test above cannot pass by always being disabled |

D1 and D2 were the two Phase 4 fixes deliberately left unit-uncovered because
both need a `layout.Context`; they are covered now.

The overwrite tests deliberately use a save callback that writes sentinel bytes
rather than the real driver. The subject is the dialog's *gating*, and a
sentinel makes "untouched" a real assertion instead of a vacuous one — a
recorder callback that never writes would pass whether or not the gate worked.

**Finding, out of scope, not fixed.** The real save path ignores the typed
filename. `FileService.SaveSettings` writes to
`filepath.Dir(pickedPath)/<TemplateName>.gen.json`, so the Save dialog's
filename field only ever chooses the *directory*; renaming the file there is
silently discarded. Normally invisible because `State.SaveAs` prefills the field
with the sanitized template name. Outside §2.1/§2.5, so it is reported to the
owner rather than changed, and the disk-level test asserts "exactly one
`*.gen.json` landed in the chosen directory" so it does not bake the quirk in.

**Verification**

| Check | Result |
| --- | --- |
| `go build ./...`, `go vet ./...`, `go vet -tags='integration_test,gui' ./...` | clean |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `go test -count=1 ./test/unit/...` | 0 failures |
| `go test -tags=integration_test ./test/integration/...` | `ok ... 2.935s` |
| `go test -tags 'integration_test,gui' ./test/integration/gui/...` | 14/14 PASS, snapshot goldens unchanged |
| `golangci-lint-v2 run ./...` | `0 issues.` |

Two notes for whoever follows. The new files were written with CRLF while the
repo is LF throughout; `gofmt` reports such a file as *entirely* unformatted, so
check line endings first when a whole-file `gofmt` diff appears. And the default
lint run carries no build tags, so gated files are never linted by it — they
were checked separately with
`golangci-lint-v2 run --build-tags 'integration_test,gui'`, which is clean apart
from `godot` objecting to the established `// X ONLY FOR INTEGRATION TEST USE`
comment convention shared with the other `*_testexports.go` files.

---

## Phase 7: Close out
Status: Complete

- [x] Full suite: build, both `go vet` tag combinations, testlayoutcheck, unit,
      integration, GPU-gated GUI, coverage, lint.
- [x] Coverage ≥ the 68.7% baseline; lint ≤ the baseline issue count.
- [x] Mark §2.1 and §2.5 `✅ FIXED` **in place** in
      [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) and update §12
      item 13.
- [x] Update repository memory.
- [x] Rewrite `.agent/session-carry-forward.md`.
- [x] Stop for owner review. Do not stage. Do not commit.

### Verification Plan
- Every command in AGENTS.md §7 Quick Reference passes.

### Phase Summary

Every §7 Quick Reference command passes. Running them surfaced one thing worth
recording — though the first diagnosis of it was **wrong**, and the correction is
the more useful lesson.

**The false alarm, and what is actually true.** `gofmt -l .` flagged six files
that `golangci-lint-v2 run ./...` had been reporting as `0 issues.` all along.
Two of them are gated `//go:build !windows`, so the Windows lint run never parses
them while CI's ubuntu run would. Reproducing with
`$env:GOOS='linux'; golangci-lint-v2 run ./...` duly reported six issues —
`gci`/`gofmt`/`golines` against [string_other.go](../internal/helpers/string_other.go)
and [hiddenAttribute_other.go](../internal/services/file_system/hiddenAttribute_other.go),
all caused by CRLF line endings in an LF repository. That looked like a live CI
breakage.

It was not. [.gitattributes](../.gitattributes) declares `*.go text eol=lf`, so
git normalises on commit: the committed blob for `string_other.go` contains
**zero** CRLF pairs (verified with `git cat-file blob` and a byte scan), CI
checks out LF, and CI lint has always been clean. The CRLF lived only in this
machine's working tree, which is why `git status` shows the files as modified
while `git diff HEAD` reports no change at all — normalisation makes the fix a
literal no-op for the repository.

The two files were still converted to LF, because a working tree that disagrees
with the repo makes `gofmt -l` and any tag-crossing lint run produce noise that
masks real findings — which is exactly what happened here. But it is **not** a bug
fix and it will not appear in the commit.

**The durable lessons**, corrected:

1. CRLF written locally never reaches the repository; `.gitattributes` handles it.
   Its only cost is local tooling noise. Do not escalate it.
2. `golangci-lint-v2 run ./...` genuinely does skip build-tag-gated files, and CI
   uses a different tag set (ubuntu, no tags) than a local Windows run. That part
   of the diagnosis stands, and `GOOS=linux` lint is a cheap way to check it — but
   run it against a *clean* tree, or line-ending noise will dominate the output.
3. `git status` reporting modified while `git diff` reports nothing is the
   signature of an attribute-normalised file. Check the blob, not the worktree.

**One real, latent finding, deliberately left alone.**
[manualCastleReapply_integration_test.go](../test/integration/manualCastleReapply_integration_test.go)
has genuinely mis-indented assertions in the committed blob (no CRLF involved).
It carries `//go:build integration_test`, so neither the local run nor CI lints
it. Real, harmless today, and a drive-by fix with no failing check behind it —
better handled by a future sweep together with the other tag-gated files.

`tmp/filecov.ps1` needed no decision: `tmp/` is gitignored and the script is gone.

A full Linux `go build ./...` cannot be run from this machine — Gio's
`gioui.org/internal/vk` needs the Linux cgo/X11 toolchain — but `./internal/...`
cross-compiles and vets cleanly, which covers both `!windows` files.

**Verification**

| Check | Result |
| --- | --- |
| `go build ./...` | clean |
| `go vet ./...` / `go vet -tags='integration_test,gui' ./...` | clean |
| `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| `wire diff ./internal/composition/...` | exit 0 (injectors current) |
| `go test -count=1 ./test/unit/...` | 0 failures |
| `go test ./test/... -count=1` (tag-free) | 0 failures |
| `go test -tags=integration_test ./test/integration/...` | `ok ... 0.627s` |
| `go test -tags 'integration_test,gui' ./test/integration/gui/...` | `ok ... 1.050s`, zero snapshot diffs |
| `golangci-lint-v2 run ./...` (Windows) | `0 issues.` |
| `golangci-lint-v2 run ./...` (`GOOS=linux`) | `0 issues.` (worktree noise cleared) |
| Total unit coverage | **69.3 %** (baseline 68.7 %) |

---

## Final Recap

Review findings **§2.1** (filesystem policy inside the GUI layer) and **§2.5**
(`fileExplorerDialog.go` is a god object) are closed, in that order, because the
review makes the split conditional on the extraction — splitting first would only
have scattered the same policy across more files.

**What moved.** `app/gui/dialogs` made sixteen direct `os`/`filepath` calls; it
now makes none, imports neither `os` nor `syscall`, and uses `filepath` only for
`Base` in display code. The policy lives in `internal/services/file_system` as
two stateless services behind `IDirectoryBrowserService` and
`IPathResolutionService`, reached from the GUI through
`handler_interfaces.IFileSystemHandler` — the flat, depguard-approved seam — and
built by its own `FileSystemSet` so the file dialogs do not drag the generator,
the repositories or the validators into their object graph. The dialog itself
went 750 → 221 LOC across five sibling files.

**The extraction paid for itself immediately.** Five real defects (D0–D4) were
sitting in code that no test could reach: drive-letter probing on Unix, an
existing *folder* offered as an overwritable file, a Save button that was enabled
for filenames that could never save, Windows reserved device names reporting
success while writing nothing, and directory creation falling back to the process
working directory. Each is now fixed with a test. Two further candidates were
examined and **rejected as deliberate** — `isHidden` failing open, and `loadDir`
adopting an unreadable directory — and are recorded so they are not "fixed" later.

**Coverage** rose 68.7 % → 69.3 % overall, with the extracted package at 92.4 %
across 19 mirrored unit-test files; the remaining 7.6 % is unreachable-by-design
OS-failure fallbacks. The dialog branches that genuinely need a `layout.Context`
(D1's refusal, D2's button state) are covered by ten new GPU-gated integration
scenarios that drive the real dialog through Gio's own `widget.Clickable.Click()`
rather than synthetic pointer coordinates.

**Three things worth carrying forward.** `Clickable.Click()` makes dialog
integration tests coordinate-free and should be the default technique for
"does this branch run" scenarios. Files behind build tags — `!windows`,
`integration_test`, `wireinject` — are skipped by the local lint run, and CI uses
a different tag set, so a `GOOS=linux` pass is a cheap sanity check at close-out.
And CRLF written into a new `.go` file is a *local* nuisance only —
`.gitattributes` normalises it away on commit — but it makes `gofmt` report the
entire file as unformatted and will happily send you chasing a CI failure that
does not exist.

**Deliberately not done, and why.** The typed filename in the Save dialog is
discarded because `FileService.SaveSettings` names the file after
`editorState.TemplateName` — confirmed by the owner as intended (review §1.1).
The defect is only that the UI still presents an editable field; that is now a
[backlog item](../todo/backlog.md) covering the read-only field, the
"Will save as:" label and the "Save As" → "Save To" rename. Mis-indented
assertions in `manualCastleReapply_integration_test.go` are real but sit behind
`integration_test`, where no check runs — left for a future sweep.

## Deployment Plan

This is a pure refactor plus five bug fixes. No schema, no output paths, no
persisted settings, no dependencies, no configuration. Nothing needs migrating
and there is nothing to roll forward or back beyond the commit itself.

1. **Review the working tree.** `git status --short` should show the files listed
   in `.agent/session-carry-forward.md` §4 and nothing else. The agent has not
   staged or committed anything (AGENTS.md §2.5).
2. **Confirm the read-only directories are untouched:**
   `git diff --stat -- data/ internal/entities/template/ internal/registry/`
   must be empty.
3. **Re-run the gate locally** (AGENTS.md §7): `go build ./...`,
   `go test -count=1 ./test/unit/...`,
   `go test -tags=integration_test ./test/integration/...`, and — because a GPU
   is present on this machine but not in CI —
   `go test -tags 'integration_test,gui' ./test/integration/gui/...`, which must
   report **zero snapshot diffs**.
4. **Optionally run the CI-equivalent lint before pushing:**
   `$env:GOOS='linux'; golangci-lint-v2 run ./... --issues-exit-code=0; Remove-Item Env:\GOOS`.
   The Windows-only run skips every `!windows` file, so this is the one gap the
   default task cannot see. Expect `0 issues.` Ignore any `gci`/`gofmt`/`golines`
   findings at `1:1` — those are local CRLF artefacts that `.gitattributes`
   normalises away on commit, not things CI will ever see.
5. **Commit and push the branch, then open the PR.** CI runs build + vet +
   testlayoutcheck on ubuntu, the unit suite on both ubuntu and windows, the race
   detector, golangci-lint, govulncheck, and a coverage-trend gate. Coverage is
   **69.3 %** against a 60.0 % floor and rising, so the trend check passes.
6. **Post-merge sanity check** (manual, one minute): launch the app, use
   *Load*, *Save As* and the output-folder picker once each. The dialogs are the
   only user-visible surface this batch touched, and the GUI snapshot suite
   already covers their rendering.

No rollback procedure is required: reverting the commit restores the previous
behaviour exactly, since nothing is written to disk in a new format or location.
