# Session Carry-Forward

## 1. Session goal

Finish **Batch 13** of the 46-finding review in
[todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) — findings **§2.1**
(the file-explorer dialog implements filesystem policy inside the GUI layer) and
**§2.5** (`fileExplorerDialog.go` is a god object), which the review makes
strictly sequential: splitting before extracting would only scatter the same
policy across more files.

**Batch 13 is COMPLETE.** All seven phases of
[plans/extract-filesystem-policy.md](../plans/extract-filesystem-policy.md) are
`Status: Complete` and that plan — not this document — is the source of truth
for the work; it carries per-phase summaries, a Final Recap and a Deployment
Plan. Phases 1–4 were reviewed and committed by the owner in earlier sessions
(`02bbb67`, `c0e35b1`); Phase 5 is staged awaiting review; Phases 6 and 7 are
unstaged.

## 2. Fixes applied

Five latent defects, all found *because* the extraction made previously
unreachable logic testable. Each landed with a regression test. Full evidence in
the plan's Phase 4 summary.

- **D0** — `ListRoots` probed `A:\`…`Z:\` on every platform. `A:\` is a legal
  relative filename on Unix, so a stray file could be shown as a volume; it also
  cost 26 pointless `stat` calls per listing.
- **D1** — an existing **folder** at the save target was offered as an
  overwritable file, then handed a directory path to `onSave`. Fixed by adding
  `DirectoryExists` and checking it *before* `PathExists`.
- **D2** — a whitespace-only filename left *Save* enabled but doing nothing.
  Fixed by deriving the button state from `resolveSaveTarget()` itself, so
  enablement and success are one predicate rather than two that can disagree.
- **D3** — Windows reserved device names (`CON`, `NUL`, `COM1`…) reported a
  successful save and wrote nothing.
- **D4** — `CreateDirectory` with an empty parent created the folder in the
  process working directory.

Two further candidates were examined and **rejected as deliberate** — do not
"fix" them later: `isHidden` fails **open** so an unreadable entry stays visible
rather than silently vanishing, and `loadDir` adopts an unreadable directory as
`currentDir` so the path bar names the place the error refers to.

## 3. Features added / changed

- **`internal/services/file_system`** — two stateless services behind
  `IDirectoryBrowserService` (`ListEntries`, `ListRoots`, `CreateDirectory`) and
  `IPathResolutionService` (`ResolveStartDirectory`, `ParentDirectory`,
  `ResolveSaveTarget`, `PathExists`, `DirectoryExists`). All the policy the GUI
  used to own.
- **`handler_interfaces.IFileSystemHandler`** — the GUI-facing seam, a **flat**
  eight-method union. Flat on purpose: embedding the service interfaces would
  force `app/gui` to import `internal/services` and trip depguard's
  `no-services-from-app`. Built by its own `FileSystemSet` /
  `InitializeFileSystemHandler`, disjoint from `GuiHandlerSet`, so the file
  dialogs do not drag in the generator, repositories or validators.
- **`app/gui/dialogs` makes zero `os` calls** — it imports neither `os` nor
  `syscall`, and uses `filepath` only for `Base` in display code.
- **The dialog split 750 → 221 LOC** across five files:
  [fileExplorerDialog.go](../app/gui/dialogs/fileExplorerDialog.go) (struct,
  `IDialog`, navigation),
  [fileExplorerDialogConfirm.go](../app/gui/dialogs/fileExplorerDialogConfirm.go)
  (footer, button state, overwrite),
  [fileExplorerDialogEntries.go](../app/gui/dialogs/fileExplorerDialogEntries.go)
  (load + list + rows),
  [fileExplorerDialogToolbar.go](../app/gui/dialogs/fileExplorerDialogToolbar.go)
  (header, save row, new-folder row) and
  [fileExplorerDialogModes.go](../app/gui/dialogs/fileExplorerDialogModes.go)
  (mode enum + the four constructors). Every move was verbatim; the GUI snapshot
  suite produced zero diffs, which is the evidence that it was a pure split.
- **Ten GUI integration scenarios** drive the real dialog end to end, including
  D1's and D2's branches, which no unit test can reach.

## 4. File modifications

Committed in earlier sessions (Phases 1–4), listed for completeness:
`internal/models/directoryEntry.go`,
`internal/common/common_errors/fileSystemErrors.go`,
`internal/services/file_system/*` (7 files incl.
`hiddenAttribute_{windows,other}.go`),
`internal/handlers/fileSystemHandler.go`,
`internal/handlers/handler_interfaces/fileSystemHandlerInterface.go`,
`internal/composition/{providerSets.go,wire.go,wire_gen.go}`,
`app/gui/program.go`, `app/gui/drivers/*`, three mocks in `test/test_helpers/`,
and 19 unit-test files.

Uncommitted now (`git status --short` against `c0e35b1`):

| File | State | Change |
| --- | --- | --- |
| [app/gui/dialogs/fileExplorerDialog.go](../app/gui/dialogs/fileExplorerDialog.go) | staged `M` | reduced to the struct, constructor, `IDialog` methods and navigation |
| [app/gui/dialogs/fileExplorerDialogConfirm.go](../app/gui/dialogs/fileExplorerDialogConfirm.go) | staged `A` | 156 LOC, moved verbatim |
| [app/gui/dialogs/fileExplorerDialogEntries.go](../app/gui/dialogs/fileExplorerDialogEntries.go) | staged `A` | 130 LOC, moved verbatim |
| [app/gui/dialogs/fileExplorerDialogToolbar.go](../app/gui/dialogs/fileExplorerDialogToolbar.go) | staged `A` | 103 LOC, moved verbatim |
| [app/gui/dialogs/fileExplorerDialogModes.go](../app/gui/dialogs/fileExplorerDialogModes.go) | staged `A` | 65 LOC, mode enum + four constructors |
| [app/gui/dialogs/fileExplorerDialog_testexports.go](../app/gui/dialogs/fileExplorerDialog_testexports.go) | unstaged `M` | +16 integration accessors (`//go:build integration_test`) |
| [test/integration/gui/fileExplorerDialog_integration_test.go](../test/integration/gui/fileExplorerDialog_integration_test.go) | unstaged `A` | new, 10 scenarios (`integration_test && gui`) |
| [plans/extract-filesystem-policy.md](../plans/extract-filesystem-policy.md) | both | Phases 5–7 summaries, Final Recap, Deployment Plan |
| [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) | unstaged `M` | §2.1 and §2.5 marked `✅ FIXED` in place; §12 item 13 updated |
| [todo/backlog.md](../todo/backlog.md) | staged `M` | new "Save As is really Save To" item (see §7) |
| [todo/test_observations.md](../todo/test_observations.md) | staged `M` | `fileExplorerDialog.go` entry rewritten to record what is now covered |

Two further files — `internal/helpers/string_other.go` and
`internal/services/file_system/hiddenAttribute_other.go` — show as modified in
`git status` but produce **no diff**: only their working-tree line endings were
normalised to LF. See §7.

## 5. Tests added or updated

- 19 unit-test files under `test/unit/internal/services/file_system/` and
  `test/unit/internal/handlers/fileSystemHandler/` (committed). Extracted
  package coverage **92.4 %**; the 7.6 % remainder is unreachable-by-design
  OS-failure fallbacks.
- 10 new scenarios in
  [fileExplorerDialog_integration_test.go](../test/integration/gui/fileExplorerDialog_integration_test.go):
  open→confirm installs state; save-target resolution; a save through the real
  `drivers.State` that lands bytes on disk; the three overwrite paths
  (refuse-to-write, cancel, confirm); folder creation; D1's folder refusal; and
  D2's disabled/enabled confirm pair. The overwrite tests write sentinel bytes
  (`"original"` → `"rewritten"`) so "untouched" is a real assertion rather than a
  vacuous one.

**Technique worth reusing:** these queue clicks through Gio's own public
`widget.Clickable.Click()`. `Clickable.update` drains `requestClicks` *before* it
consults pointer input, so the click is indistinguishable from a real one to the
code under test while the layout still runs for real — no coordinates, no
calibration, no prior-frame requirement. Far more robust than the pixel-based
`AppRunner.ClickAt` used by the benchmarks; reserve that one for genuinely
geometric behaviour (drag, scroll, hit-testing).

Last full run — all green:

| Suite | Result |
| --- | --- |
| `go test -count=1 ./test/unit/...` | 0 failures |
| `go test ./test/... -count=1` (tag-free) | 0 failures |
| `go test -tags=integration_test ./test/integration/...` | `ok … 0.627s` |
| `go test -tags 'integration_test,gui' ./test/integration/gui/...` | `ok … 1.050s`, 14/14, zero snapshot diffs |
| `go build` · both `go vet` tag sets · `testlayoutcheck` · `wire diff` | clean / passed / exit 0 |
| `golangci-lint-v2 run ./...` | `0 issues.` |
| Total unit coverage | **69.3 %** (baseline 68.7 %) |

## 6. Git status snapshot

Branch **`AD/refactoring-07-21`**, HEAD `c0e35b1 "Batch 13-Part 2 Done"`
(= `origin/AD/refactoring-07-21`, so nothing is unpushed). Nothing was staged or
committed by the agent; the staged entries in §4 were staged by the owner as part
of reviewing Phase 5.

The next session inherits: Phase 5 staged and under review, Phases 6–7 unstaged,
and the two line-ending-only working-tree changes from §7.

## 7. Rejections / things the user declined

- **The `{TemplateName}.gen.json` save path is intended behaviour, not a bug.**
  The agent reported that `FileService.SaveSettings` writes
  `<dir>/{TemplateName}.gen.json` and therefore discards the filename typed into
  the Save dialog. The owner confirmed this is deliberate (already recorded at
  review §1.1) — **do not "fix" the service.** The real defect is that the UI
  still offers an editable filename field. At the owner's instruction that became
  a [backlog item](../todo/backlog.md): make the field read-only, relabel it
  `"Will save as:"`, and rename the action `"Save As"` → `"Save To"` in both UI
  text and code. Two open decisions are recorded inside that item — whether
  `NewSaveFileDialog`/`modeSaveFile`/`onSave` should follow the rename (they name
  the explorer *mode*, not the toolbar action), and what should trigger the
  disabled-confirm test once a user can no longer type a whitespace filename.
- **A CRLF "CI failure" that turned out not to exist.** `gofmt -l .` flagged six
  files the Windows lint run reported clean, and a `GOOS=linux` lint reproduced
  six `gci`/`gofmt`/`golines` issues, which looked like a live CI break. It was
  not: [.gitattributes](../.gitattributes) declares `*.go text eol=lf`, so git
  normalises on commit — the committed blobs contain **zero** CRLF (verified by
  byte-scanning `git cat-file blob`) and CI has always been clean. Two files were
  still converted to LF to stop the local noise masking real findings, which is
  why they show as modified with an empty diff. **Lesson: `git status` modified +
  `git diff` empty is the signature of an attribute-normalised file; check the
  blob, not the worktree.**
- **Four tag-gated files left formatting-dirty on purpose.** `wire.go`,
  `dialogHost_testexports.go`, `stateSaveAs_integration_test.go` (all CRLF-only,
  i.e. non-issues) and `manualCastleReapply_integration_test.go` (genuinely
  over-indented assertions in the committed blob). All sit behind
  `integration_test`/`wireinject`, so neither the local run nor CI lints them.
  Real but latent; a drive-by fix with no failing check behind it.

## 8. Open questions

Both predate this batch and are still unanswered:

1. Delete `zoneEditorHandler.ComputeHasErrors` / `RebuildZoneConnectionRoads`, or
   add them to `handler_interfaces.IZoneEditorHandler`? They are exported on the
   private struct, absent from the interface, and called by nobody.
2. Should `internal/helpers/io.go`'s Steam/registry install-discovery chain get an
   injectable filesystem seam so it can be covered? The new
   `IPathResolutionService` is now a natural home for it.

And one that blocks the next batch:

3. **§2.2 scope** — extracting regeneration policy out of the GUI driver is
   multi-session, overlaps a backlog item, and is now the **only** owner decision
   left open in the entire review.

## 9. Next recommended actions

1. Review Phases 6–7 (the testexports accessors, the integration test, and the
   documentation updates), then stage and commit Batch 13. The agent never
   stages or commits.
2. Decide §2.2's scope, then write `plans/extract-regeneration-policy.md` before
   touching any code (AGENTS.md §4.7).
3. Optionally schedule the "Save As → Save To" backlog item; it is small,
   self-contained and touches only the GUI.
4. Optionally sweep the four tag-gated formatting-dirty files from §7.

Remaining in the review after Batch 13: **§2.2** (regeneration policy) and
**§2.6** (`zoneEditorDialog`'s ~58 fields, plus the still-oversized
`zoneEditorDialog.go` 507 / `zoneEditorCanvas.go` 479 / `bonusPickerDialog.go`
434 / `pickerDialog.go` 371 / `ruleDialog.go` 314 / `zoneContent.go` 299 — low
priority, opportunistic).

## 10. Carry-forward prompt

> Read [AGENTS.md](../AGENTS.md) first and follow it strictly. The hard rules, one
> line each: never modify `data/`, `internal/entities/template/` or
> `internal/registry/`; keep everything cross-platform (Windows + Linux,
> `path/filepath`, PowerShell chains with `;` and never `&&`); every change ships
> with tests and must not drop coverage (currently 69.3 %, baseline 68.7 %);
> durable multi-session work gets a plan file under `plans/` before any code;
> **never stage and never commit** — the owner reviews and commits; and never
> change where `.rmg.json` is written nor persist the output directory.
>
> We are remediating the 46-finding review in
> [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md); §12 defines the
> PR-sized batches and findings are marked `✅ FIXED` / `❌ WILL NOT FIX`
> **in place** in that document. For each batch: ask whether it should be done at
> all; if declined, record in the review file why it should not be attempted
> again; ask every clarifying question up front; implement; rewrite
> `.agent/session-carry-forward.md`; then stop and wait for review.
>
> **Batch 13 (§2.1 + §2.5) is complete** — filesystem policy now lives in
> `internal/services/file_system` behind `handler_interfaces.IFileSystemHandler`,
> `fileExplorerDialog.go` is split into five files, five latent defects are fixed
> and ten GUI integration scenarios cover the dialog. Phase 5 is staged under
> owner review; Phases 6–7 are unstaged. See
> [plans/extract-filesystem-policy.md](../plans/extract-filesystem-policy.md) for
> the full record and
> [.agent/session-carry-forward.md](session-carry-forward.md) for the handoff,
> including three open questions and the deliberate rejections you must not
> re-litigate (notably: writing `{TemplateName}.gen.json` and ignoring the typed
> filename is **intended**).
>
> Next up is **§2.2** — extract regeneration policy out of the GUI driver. Its
> scope is the last owner decision outstanding in the review, so ask before
> planning anything.
