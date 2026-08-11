# Session carry-forward — 2026-08-11

## 1. Session goal

Compile the two still-open items of [review-opus5-08-04.md](../todo/review-opus5-08-04.md),
all of `todo/backlog.md` and the actionable gaps of
[test_observations.md](../todo/test_observations.md) into a single
`review-prompt.md`-format backlog, then fold the owner's design comment in
[runnerHandler.go](../test/test_helpers/integration_common/runnerHandler.go)
into it as well.

**This was a documentation-only session. No production or test Go code was
changed.**

## 2. Fixes applied

None — no code defects were fixed. Two documentation defects were:

- Three links in [backlog-opus5.md](../todo/backlog-opus5.md) pointed at
  `todo/backlog.md`, which the owner deleted in `cd7ad10` after accepting the
  new document. De-linked and re-worded to reference the commit instead.
- The supersession paragraph was written in the future tense ("the owner
  deletes…"); it now records what actually happened.

## 3. Features added / changed

- **[todo/backlog-opus5.md](../todo/backlog-opus5.md) — created (earlier in this
  session, already committed by the owner as `cd7ad10`).** 17 items in the
  `review-prompt.md` output format, each with evidence, failure mode, concrete
  fix plan, exact test files and owner-decision flags.
- **§5.4 and §5.5 — added (this turn, uncommitted).** The 46-line design
  comment in `runnerHandler.go` was investigated and split:
  - **§5.4 🟡** the GUI handler framework, broken into seven separable pieces
    (a) file/naming hygiene, (b) coordinate strategy, (c) mask narrowing,
    (d) per-tab/per-dialog handlers, (e) layout-shifting state, (f) the missing
    scroll seam, (g) keeping `*_testexports.go` out of the handler — with an
    explicit scope guard against building (d)–(g) speculatively.
  - **§5.5 🟡** the local-vs-CI snapshot discrepancy, which is a *different
    class of problem* (a rendering/tolerance defect, not framework design).
  - §0.4 records which sentence of the comment went where, and §8 gains
    batches **L** and **M**.

## 4. File modifications

| File | Change |
| --- | --- |
| [todo/backlog-opus5.md](../todo/backlog-opus5.md) | Added §0.4 (disposition of the `runnerHandler.go` comment), §5.4, §5.5, batches L/M in §8, a sources-table row; item count 17 → 19 (10 🟡); de-linked the deleted `backlog.md`. |
| `.agent/session-carry-forward.md` | Created (this file). |
| `/memories/repo/conventions.md` | Agent-local memory, not in the repo. Records the backlog's existence, the `GetBorderGuardValue` insight and the `runSubagent` model-name format. |

Untouched but read for evidence: `appRunner.go`, `appRunnerSnapshots.go`,
`snapshot/comparer.go`, `.github/workflows/pr-validation.yml`,
`app/gui/constants/ui.go`, `window_snapshot_integration_test.go`,
`window_tab_cycling_test.go`.

## 5. Tests added or updated

None — no Go code changed, so no suite was re-run this session. The last
verification in the terminal (`go build ./...; go vet -tags="integration_test,gui" ./...;
go run ./cmd/testlayoutcheck .; gofmt -l ./app ./internal ./test ./cmd`) exited
**0**. The baselines the backlog is written against remain: unit coverage
**72.5 %** (floor 69.3 %), lint **0 issues**, all suites green.

## 6. Git status snapshot

Branch: `AD/refactoring-07-21` (up to date with `origin/AD/refactoring-07-21`,
HEAD `594eae2`).

```
 M todo/backlog-opus5.md
```

Plus the untracked `.agent/` folder created by this document. Nothing was
staged and nothing was committed (AGENTS.md §2.5). The next session inherits one
unstaged documentation edit.

## 7. Rejections / things the user declined

Nothing was declined this turn. Standing decisions from earlier in the session
that later sessions must not re-litigate:

- The output directory is **never** persisted (AGENTS.md §2.7) — do not
  re-propose it.
- `todo/backlog.md` was deleted by the owner, deliberately; §0.2 is its record.
- §6.1 (dead `createTopologyAdjacency` branches) reverses a rollback the owner
  performed on purpose — owner-gated.
- §2.2 Branch A, §2.4 and §2.5 touch protected directories — owner-gated.

## 8. Open questions

- **§5.5 step 1 needs a CI artifact.** Whether the local/CI text difference is a
  half-rendered frame or genuine llvmpipe anti-aliasing cannot be settled from
  the workspace; it needs the `gui-snapshot-failures` artifact from a real
  failing run, or a deliberate CI run with the threshold lowered.
- **§5.4 (b)** recommends keeping coordinates literal but centralised. If the
  owner would rather add a widget-rect lookup seam to `*_testexports.go`, that
  reverses the recommendation — worth one sentence of confirmation before batch L.

## 9. Next recommended actions

1. Review §5.4/§5.5 in [backlog-opus5.md](../todo/backlog-opus5.md); commit the
   pending edit.
2. Start **batch A** (§1.4 — fatal window error to a discard handler). Two-line
   fix, no dependencies, ideal warm-up.
3. Then **batch B** (§1.2 → §1.3, hub and portal guard values). Owner already
   approved the behaviour change on 2026-08-11; expect golden-template churn.
4. Batch **L** (§5.4 a–c, §5.5) before any of §5.1–§5.3, so the new GUI tests
   are written against a settled harness and a comparison that can actually fail.

## 10. Carry-forward prompt

> Read `AGENTS.md` first — it governs everything below.
>
> Hard rules, one line each: never modify `data/`, `internal/entities/template/`
> or `internal/registry/` without explicit owner approval; keep every change
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chained with `;`
> never `&&`); every change ships with tests and must not drop unit coverage
> below 69.3 % (currently 72.5 %); durable multi-session work gets a plan file
> under `plans/`; never stage and never commit — the owner reviews and commits;
> never change where `.rmg.json` is written and never persist the output
> directory; never run a bulk in-place rewrite over the repository.
>
> Where work left off: `todo/backlog-opus5.md` is the compiled, authoritative
> backlog (19 items, 0 🔴 · 6 🟠 · 10 🟡 · 3 ⚪). It supersedes the deleted
> `todo/backlog.md` and folds in the two open items of `review-opus5-08-04.md`,
> three promoted gaps from `test_observations.md`, and the design comment in
> `test/test_helpers/integration_common/runnerHandler.go` (now §5.4 and §5.5).
> No code has been changed yet — every item is still open. §8 of that document
> gives the execution order; start with batch A (§1.4).
>
> See `./.agent/session-carry-forward.md` for the full handoff.
