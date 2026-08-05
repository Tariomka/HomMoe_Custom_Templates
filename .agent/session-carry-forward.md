# Session Carry-Forward — Review Remediation (`todo/review-opus5-08-04.md`)

Last updated: **2026-08-05**, after **Batch 2 of 13**.

## 1. Session goal

Work through every finding in [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md)
in the §12 execution order, one batch at a time, with owner approval before each
batch and a carry-forward update after each batch.

## 2. Fixes applied

### Batch 1 — Security PR (owner-reviewed, now committed)

- **§8.1 (🔴 reachable DoS)** — `golang.org/x/text` `v0.38.0` → `v0.39.0`.
  Closes GO-2026-5970 (infinite loop in `norm.Form.*`, reachable from
  `app/gui/program.go:21` through Gio's text handling).
- **§8.2 (🟠 present, uncalled)** — `golang.org/x/net` `v0.55.0` → `v0.56.0`.
  Closes GO-2026-5942.
- Transitive: `golang.org/x/sys` `v0.45.0` → `v0.46.0`, applied by `go mod tidy`.

### Batch 2 — Correctness PR (this batch, awaiting owner review)

- **§1.2 (🔴 data-loss)** — `SaveAs` no longer records `currentPath` when the
  write failed. `handleSaveState` now returns `bool` (mirroring
  `handleLoadState`) and the dialog callback gates the assignment on it. Without
  this, a failed *Save As* left `currentPath` pointing at a file that was never
  created, and every later `Save` silently retargeted that broken path.
  ([stateFiles.go](../app/gui/drivers/stateFiles.go))
- **§1.3 (🟠 latent nil deref)** — `WasLayoutChanged` now guards with
  `HasPreviousState()` before comparing against `previous`.
  ([editorState.go](../app/gui/models/editorState.go))
- **§1.3 follow-on (owner asked for "Both")** — `ShouldReapplyManualEdits`
  dropped its now-redundant `!this.HasPreviousState() ||` short-circuit and
  reads as `HasManualEdits() && !WasLayoutChanged()`.
- **§5.2 (🟡 dead comment)** — owner chose **"Uncomment it"** over the review's
  recommended delete. `SnapshotCurrentState` now really clears `next`, with a
  one-line comment explaining why. Traced as behaviourally a no-op today, so the
  invariant is now explicit rather than accidental.

## 3. Features added / changed

No user-visible feature changes. §1.2 changes observable behaviour only in the
failure path: after a failed *Save As* the editor now stays "unfiled" and the
next *Save* re-prompts instead of silently retrying a path that holds no file.

## 4. File modifications

| File | Change |
| --- | --- |
| [app/gui/drivers/stateFiles.go](../app/gui/drivers/stateFiles.go) | **M** — `handleSaveState` returns `bool`; `SaveAs` gates `currentPath` on it |
| [app/gui/models/editorState.go](../app/gui/models/editorState.go) | **M** — `WasLayoutChanged` nil guard; `ShouldReapplyManualEdits` simplified; `SnapshotCurrentState` clears `next` |
| [app/gui/drivers/dialogHost_testexports.go](../app/gui/drivers/dialogHost_testexports.go) | **NEW** — `//go:build integration_test`; exposes `GetTopDialog` |
| [app/gui/dialogs/fileExplorerDialog_testexports.go](../app/gui/dialogs/fileExplorerDialog_testexports.go) | **NEW** — `//go:build integration_test`; exposes `ConfirmSave` (fires `onSave` without a `layout.Context`) |
| [test/integration/stateSaveAs_integration_test.go](../test/integration/stateSaveAs_integration_test.go) | **NEW** — two §1.2 regression tests; legitimately tagged (consumes the two accessors above) |
| [test/unit/app/gui/models/editorState/wasLayoutChanged_test.go](../test/unit/app/gui/models/editorState/wasLayoutChanged_test.go) | **M** — added `TestWhenNoPreviousStateExists_ReportsLayoutNotChanged` |
| [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) | **M** — §1.2/§1.3/§5.2 marked `✅ FIXED`, §6.1 marked `❌ WILL NOT FIX`, §12 item 2 + the §6.1 summary-table row updated |
| [todo/test_observations.md](../todo/test_observations.md) | **M** — recorded why §1.2 has no unit test |
| `.agent/session-carry-forward.md` | **M** — this document |

## 5. Tests added or updated

| Test | Location |
| --- | --- |
| `TestWhenSaveAsFails_CurrentPathIsNotRecorded` | [stateSaveAs_integration_test.go](../test/integration/stateSaveAs_integration_test.go) |
| `TestWhenSaveAsSucceeds_CurrentPathIsRecorded` | [stateSaveAs_integration_test.go](../test/integration/stateSaveAs_integration_test.go) |
| `TestWhenNoPreviousStateExists_ReportsLayoutNotChanged` | [wasLayoutChanged_test.go](../test/unit/app/gui/models/editorState/wasLayoutChanged_test.go) |

Verification run after Batch 2:

| Check | Result |
| --- | --- |
| `go build ./...` | Pass |
| `go vet -tags=integration_test ./...` | Pass, no output |
| `go test -count=1 ./test/unit/...` | Pass, exit 0, no `FAIL` |
| `go test -tags=integration_test -count=1 ./test/integration/...` | `ok ... 2.431s` |
| Both new SaveAs tests, run by name with `-v` | `--- PASS` ×2 |
| Unit coverage total (AGENTS.md §2.3) | **64.7%** — unchanged from the review baseline; `editorState.go` at 100% on every changed function |
| `golangci-lint-v2 run ./app/... ./test/...` | 13 findings, **all pre-existing baseline** (2 `dupl` = §3.4, 11 `gochecknoglobals`); zero new, zero formatting findings on the new files |

No existing test needed changing to accommodate the §5.2 `next = nil`
activation — confirmation that it is the no-op it was traced to be.

Batch 1 verification, for the record: `go build ./...` pass; unit suite pass;
`govulncheck ./...` and `govulncheck -scan module` both *No vulnerabilities
found*; `go mod tidy -diff` exit 0. (`govulncheck -scan module` rejects path
patterns — run it bare, not with `./...`.)

## 6. Git status snapshot

Branch: `AD/refactoring-07-21`

```
 M app/gui/drivers/stateFiles.go
 M app/gui/models/editorState.go
 M test/unit/app/gui/models/editorState/wasLayoutChanged_test.go
 M todo/review-opus5-08-04.md
 M todo/test_observations.md
?? app/gui/dialogs/fileExplorerDialog_testexports.go
?? app/gui/drivers/dialogHost_testexports.go
?? test/integration/stateSaveAs_integration_test.go
```

Batch 1's `go.mod` / `go.sum` changes no longer appear — the owner has already
staged/committed them. Everything above is **Batch 2 only** and is unstaged;
per AGENTS.md §2.5 the agent must not stage or commit. `.agent/` is untracked
and may be added to `.gitignore` if the owner wants.

## 7. Rejections / things the owner declined

- **§6.1 — rejected outright, must not be re-attempted.** The finding asks for
  `//go:build integration_test` on
  [rmgTemplateModel_test.go](../test/integration/rmgTemplateModel_test.go),
  quoting an AGENTS.md §4.6.1 that says *"every file in these two directories
  carries the tag"*. **AGENTS.md no longer says that.** §4.6.1 was rewritten
  after the review was authored and now reads: the tag applies *if and only if*
  the file references an accessor declared in a `*_testexports.go` file, and it
  is explicitly *not* a label meaning "this is an integration test". The file
  references no such accessor (verified by grepping every accessor in
  `window_testexports.go` and `state_testexports.go` — zero hits); it only
  decodes `.rmg.json` files from `data/`. Tagging it would violate AGENTS.md.
  The review's *secondary* concern is real but unaddressed: the untagged file
  makes a plain `go test ./test/...` decode all bundled example templates, so
  the "fast" suite is slower than it looks. Fix that with `testing.Short()` or a
  sampled corpus if it ever hurts — never with a build tag.
- **Separate `plans/` file declined.** Owner chose to track progress by marking
  items `✅ FIXED` / `❌ WILL NOT FIX` in place inside
  `todo/review-opus5-08-04.md` (the convention the review's own §12 prescribes)
  rather than creating `plans/review-opus5-08-04-remediation.md`.
- **Full per-batch verification declined.** Owner chose `go build ./...` +
  `go test ./test/unit/...` as the standing per-batch gate instead of the full
  AGENTS.md §2.3 sweep. Batch 2 additionally ran vet, the integration suite,
  coverage and lint because it touched the gated tree and production code.

## 8. Open questions

Owner decisions still required before their batches can start (from review §12):

| Batch | Decision |
| --- | --- |
| 3 | §1.1 — do template JSON + preview PNG commit transactionally, or keep the documented partial-success contract? |
| 4 | §1.5 — the exact numeric ceilings for the twenty `.gen.json` int fields (and four floats). |
| 8 | §7.1 — are direct pushes to `master` intentionally permitted? |
| 9 | §9.1 — is external programmatic use a supported product surface? |
| 12 | §2.7 — finish or remove the gladiator-arena preview? §1.8 — output-directory persistence shape (a) in `.gen.json` vs (b) machine-local (review recommends (b)). |
| 13 | §2.2 — confirm scope/sequencing of the regeneration-policy refactor. |

## 9. Next recommended actions

1. Owner reviews Batch 2 — the 8 files listed in §4 above.
2. Ask the owner the **§1.1 transactionality** question; Batch 3 (Durability:
   §1.1 + §1.6 + §5.1) cannot start without it. Offer the alternative of jumping
   to an unblocked batch — **5** (performance, §4.1), **6** (DI, §2.3 + §2.4),
   **7** (test policy, §6.3 then §6.5), **10** (duplication) or **11**
   (coverage) — if the owner would rather not decide now.
3. Continue batch by batch through review §12, keeping the
   ask → implement → document → stop cadence.

## 10. Carry-forward prompt

> Read `AGENTS.md` first, then `todo/review-opus5-08-04.md`.
>
> Hard rules, one line each:
> §2.1 — never modify `data/`, `internal/entities/template/` or
> `internal/registry/`; read them freely, propose changes only.
> §2.2 — everything must build and run on both Windows and Linux; use
> `path/filepath`, no OS-specific syscalls without build tags, chain PowerShell
> commands with `;` never `&&`.
> §2.3 — every non-trivial change ships with tests in `test/`, and unit coverage
> must not drop.
> §2.4 — multi-session work gets a durable plan artifact (here: the owner chose
> to mark items `✅ FIXED` in place in the review document instead).
> §2.5 — never stage and never commit; the owner reviews and stages.
>
> Where work left off: Batches **1 (Security, §8.1 + §8.2)** and **2
> (Correctness, §1.2 + §1.3 + §5.2)** of the 13 in review §12 are done. Batch 1
> is committed; Batch 2 is unstaged on branch `AD/refactoring-07-21` and awaiting
> review — build, vet, unit suite, integration suite, coverage (64.7%, flat) and
> lint all green. **§6.1 was rejected** — read that item before touching build
> tags and do not re-open it. Next is Batch 3 (Durability, §1.1 + §1.6 + §5.1),
> **blocked** on an owner decision about template+PNG transactionality.
>
> Process the owner asked for: before starting each batch, ask whether it should
> be done at all — if not, document in the review file why it must not be
> attempted again — and ask any clarifying questions up front. After each batch,
> update `./.agent/session-carry-forward.md` and stop for owner review before
> continuing.
>
> See `./.agent/session-carry-forward.md` for the full handoff.
