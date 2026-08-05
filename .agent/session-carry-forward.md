# Session Carry-Forward — Review Remediation, Batch 3 (Durability)

## 1. Session goal

Work through `todo/review-opus5-08-04.md` batch by batch, asking go/no-go and
clarifying questions before each batch, and stopping for owner review after each.
This document covers **Batch 3 — Durability PR (§1.1 + §1.6 + §5.1)**.

## 2. Fixes applied

- **§1.1 🔴 non-atomic persistent writes.** Every persisted user file is now
  written to `{directory}/TEMP-{name}{ext}`, `Sync`ed, closed, and only then
  renamed onto the destination. A failed encode leaves the previous file
  byte-identical and removes the temporary file.
- **§1.6 🟠 PNG truncation + swallowed close error.** The preview goes through
  the same writer, and the writer's `defer` uses a named return so a `Close`
  error wins when the encode succeeded.
- **§5.1 🟡 disagreeing directory permissions.** One `FolderPermission = 0o755`
  constant, consumed by both the repositories and the in-app file explorer
  (which previously used `0o750`).

## 3. Features added / changed

- **New repository layer.** The pre-existing but unimplemented
  [fileRepositoryInterface.go](../internal/repositories/fileRepositoryInterface.go)
  now has three real implementations. Each owns its extension and encoding;
  a shared private `atomicFileWriter` owns the temp-file/rename mechanics.
- **`FileService` is now a controller.** It only decides the directory and the
  requested file name and delegates; the writer sanitizes the name and falls
  back to `Generated_Template`. `SaveTemplate` and `SavePreviewImage` were
  replaced by a single `SaveTemplateWithPreview`, which skips the preview when
  the image is `nil` and preserves the partial-success contract (a preview
  failure still returns the template path).
- **New `internal/constants` package** — the home AGENTS.md §4.4 prescribes for
  constants, which did not exist yet.
- **⚠ Behaviour change 1 (owner-approved):** editor-state saves now create
  missing directories instead of failing.
- **⚠ Behaviour change 2 (owner-approved):** the Save As dialog's file name is
  discarded; the state is written as `{TemplateName}.gen.json`. `SaveState`
  returns the path actually written and `drivers.State` records **that**.

## 4. File modifications

**Created**

- `internal/constants/filePermissions.go` — `FolderPermission`, `FileReadWritePermission`.
- `internal/repositories/atomicFileWriter.go` — private writer: `MkdirAll` →
  TEMP file → `Sync` → `Close` → bounded `os.Rename` retry (5 × 20 ms) → discard
  temp on any failure.
- `test/unit/internal/repositories/{editorStateRepository,templateRepository,previewRepository}/`
  — `new*_test.go`, `load_test.go`, `save_test.go` in each.
- `test/unit/internal/services/file_service/fileService/common_test.go` —
  generic `mockFileRepository[T]` + `newServiceWithMocks()`.
- `test/unit/internal/services/file_service/fileService/saveTemplateWithPreview_test.go`.

**Edited**

- `internal/repositories/editorStateRepository.go` / `templateRepository.go` /
  `previewRepository.go` — implemented (`.gen.json` / `.rmg.json` / `.png`).
  Only the editor state implements `Load`; the other two return
  `common_errors.ErrNotImplemented`.
- `internal/services/file_service/fileService.go` — rewritten as a controller;
  `SaveSettings` now returns `(string, error)`.
- `internal/handlers/stateHandler.go` — returns the path `SaveSettings` reports.
- `internal/handlers/templateHandler.go` — builds the preview image (when a
  generator exists) and makes one `SaveTemplateWithPreview` call.
- `app/gui/drivers/stateFiles.go` — `handleSaveState` records the returned path;
  it no longer returns a `bool` (the gating moved inside it).
- `app/gui/dialogs/fileExplorerDialog.go` — uses `FolderPermission`.
- `internal/composition/providerSets.go` + `wire_gen.go` — three repository
  providers registered; injectors regenerated with `wire gen`.
- `test/unit/architecture/dependency/dependency_test.go` — `internal/constants`
  added to the `app/*` import allow-list.
- `test/unit/internal/handlers/guiHandler/{loadState,saveState}_test.go` — use
  the returned path; the "missing directory" test now asserts creation.
- `test/integration/{editorState,window_render,stateSaveAs}_integration_test.go`
  — read back the derived `{TemplateName}.gen.json` path.
- `todo/review-opus5-08-04.md`, `todo/test_observations.md` — documentation.

**Deleted**

- `test/unit/internal/services/file_service/fileService/saveTemplate_test.go`
- `test/unit/internal/services/file_service/fileService/savePreviewImage_test.go`

  Both tested methods that no longer exist. Their filesystem-behaviour cases
  moved to the repository test folders; their naming/routing cases moved to
  `saveTemplateWithPreview_test.go`. Both files are recoverable from git.

## 5. Tests added or updated

- Nine new repository test files, including the §1.1 regression cases:
  pre-existing destination byte-identical after a failed save, and no `TEMP-*`
  residue. Fault injection uses `math.NaN()` (JSON) and a zero-sized
  `image.RGBA` (PNG) — no test-only seam was needed.
- `FileService` tests rewritten against mocked `IFileRepository[T]`.

**Last run — all green:**

- `go build ./...` ✓ · `go vet -tags=integration_test ./...` ✓
- `go test -count=1 ./test/unit/...` → exit 0
- `go test -tags=integration_test -count=1 ./test/integration/...` → exit 0
- Coverage **64.8%** (was 64.7%). All new production files at 100% except
  `atomicFileWriter.encodeToTemporaryFile` at 77.8% (Close/Sync error branches).
- `golangci-lint-v2 run ./...` → **42 issues, all pre-existing categories**
  (`dupl` 2 = §3.4, `gochecknoglobals` 40).

## 6. Git status snapshot

Branch **`AD/refactoring-07-21`**. Nothing staged. 23 modified, 2 deleted,
4 untracked paths — see §4. (An accidental `git rm --cached` during this session
was immediately reverted with `git restore --staged`; the index is clean.)

## 7. Rejections / things the owner declined

- Batch 2's **§6.1** (build tag on `rmgTemplateModel_test.go`) — ❌ WILL NOT FIX,
  documented in the review; it contradicts the current AGENTS.md §4.6.1.
- The review's own §1.1 design (an `AtomicFileWriter` inside `file_service`) was
  **replaced** by the owner with the repository layer described above.
- Transactional template+PNG commit — **declined**; partial success is kept.
- A separate `plans/` file for this remediation — declined; findings are marked
  in place inside the review document.

## 8. Open questions

None blocking Batch 3. Blocking later batches:

- **Batch 4 / §1.5** — the exact numeric ceilings for the 20 unbounded ints and
  4 floats are a product decision.
- **Batch 8 / §7.1** — are direct pushes to `master` intentional?
- **Batch 9 / §9.1** — the public-API decision behind the QUICKSTART example.
- **Batch 12** — §2.7 (finish or remove the gladiator-arena preview) and §1.8
  (where the output directory is persisted).
- **Batch 13 / §2.2** — scope of extracting regeneration policy from `drivers`.

## 9. Next recommended actions

1. Owner reviews Batch 3 and confirms the two behaviour changes in §3.
2. **Batch 4 — Input-validation PR** (§1.5 + §1.7). Ask for the ceilings first.
3. Batch 5 — Performance (§4.1, preview layout cache).
4. Batch 6 — DI (§2.3, §2.4).
5. Batches 7–13 in the order given by review §12.

## 10. Carry-forward prompt

> Read `AGENTS.md` first. Its hard rules: never modify `data/`,
> `internal/entities/template/` or `internal/registry/`; keep everything
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chains with `;`);
> cover every change with tests and check coverage; never stage or commit
> anything; cap the session around 50 messages and hand off in
> `./.agent/session-carry-forward.md`.
>
> We are working through `todo/review-opus5-08-04.md` in the batches listed in
> its §12. Batches 1 (security), 2 (correctness) and 3 (durability) are done —
> Batch 3 introduced `internal/repositories` with an atomic file writer and
> turned `FileService` into a controller. Findings are marked `✅ FIXED` /
> `❌ WILL NOT FIX` in place inside the review document.
>
> Next is **Batch 4 — Input-validation PR (§1.5 + §1.7)**, which is blocked on
> the owner choosing the numeric ceilings. Before starting any batch: ask
> whether it should be done at all, document the rationale in the review file if
> it is declined, ask every clarifying question up front, implement, rewrite this
> carry-forward, then stop for review.
>
> Full handoff: `./.agent/session-carry-forward.md`.
