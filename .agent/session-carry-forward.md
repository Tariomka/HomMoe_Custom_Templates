# Session Carry-Forward — Review Remediation (`todo/review-opus5-08-04.md`)

Last updated: **2026-08-05**, after **Batch 1 of 13**.

## 1. Session goal

Work through every finding in [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md)
in the §12 execution order, one batch at a time, with owner approval before each
batch and a carry-forward update after each batch.

## 2. Fixes applied

- **§8.1 (🔴 reachable DoS)** — `golang.org/x/text` `v0.38.0` → `v0.39.0`
  ([go.mod](../go.mod)). Closes GO-2026-5970 (infinite loop in `norm.Form.*`,
  reachable from `app/gui/program.go:21` through Gio's text handling).
- **§8.2 (🟠 present, uncalled)** — `golang.org/x/net` `v0.55.0` → `v0.56.0`
  ([go.mod](../go.mod)). Closes GO-2026-5942.
- Transitive: `golang.org/x/sys` `v0.45.0` → `v0.46.0`, applied by `go mod tidy`.

## 3. Features added / changed

None. Batch 1 is a dependency-only change; no production Go source was touched.

## 4. File modifications

| File | Change |
| --- | --- |
| [go.mod](../go.mod) | x/text v0.39.0, x/net v0.56.0, x/sys v0.46.0 |
| [go.sum](../go.sum) | Checksums for the three upgrades |
| [todo/review-opus5-08-04.md](../todo/review-opus5-08-04.md) | §8.1, §8.2 and §12 item 1 marked `✅ FIXED` in place with evidence |
| `.agent/session-carry-forward.md` | This document (new) |

## 5. Tests added or updated

None — the review explicitly states "No Go test applies" for §8.1.

Verification run after the change:

| Check | Result |
| --- | --- |
| `go build ./...` | Pass |
| `go test -count=1 ./test/unit/...` | Pass, **0** `FAIL` packages |
| `govulncheck ./...` (symbol scan) | **No vulnerabilities found** (was 1 affected) |
| `govulncheck -scan module` | **No vulnerabilities found** (was 2) |
| `go mod tidy -diff` (root) | Exit 0 |

Not run this batch (owner chose "build + unit tests only" as the per-batch gate):
coverage report, lint, integration/GUI suites.

## 6. Git status snapshot

Branch: `AD/refactoring-07-21`. Tree was clean before the change.

Unstaged after Batch 1:

```
 M go.mod
 M go.sum
 M todo/review-opus5-08-04.md
?? .agent/
```

Nothing is staged — per AGENTS.md §2.5 the agent must not stage or commit.
`.agent/` is untracked and may be added to `.gitignore` if the owner wants.

## 7. Rejections / things the owner declined

- **Separate `plans/` file declined.** Owner chose to track progress by marking
  items `✅ FIXED` in place inside `todo/review-opus5-08-04.md` (the convention
  the review's own §12 prescribes) rather than creating
  `plans/review-opus5-08-04-remediation.md`.
- **Full per-batch verification declined.** Owner chose `go build ./...` +
  `go test ./test/unit/...` as the standing per-batch gate instead of the full
  AGENTS.md §2.3 sweep. Coverage must still be checked on batches that change
  Go code (§2.3) — raise this with the owner when Batch 2 starts.

## 8. Open questions

Owner decisions still required before their batches can start (from review §12):

| Batch | Decision |
| --- | --- |
| 3 | §1.1 — do template JSON + preview PNG commit transactionally, or keep the documented partial-success contract? |
| 4 | §1.5 — the exact numeric ceilings for the twenty `.gen.json` int fields. |
| 8 | §7.1 — are direct pushes to `master` intentionally permitted? |
| 9 | §9.1 — is external programmatic use a supported product surface? |
| 12 | §2.7 — finish or remove the gladiator-arena preview? §1.8 — output-directory persistence shape (a) in `.gen.json` vs (b) machine-local. |
| 13 | §2.2 — confirm scope/sequencing of the regeneration-policy refactor. |

## 9. Next recommended actions

Batch order is review §12. Batch 1 is done; next is:

1. **Batch 2 — Correctness PR** (cheap, isolated, no owner decision needed):
   §1.2 `SaveAs` records the path on failure, §1.3 `WasLayoutChanged` nil guard,
   §6.1 missing `integration_test` build tag on
   `test/integration/rmgTemplateModel_test.go`, §5.2 dead commented-out line in
   `SnapshotCurrentState`. Plus the two named regression tests
   (`saveAs_test.go`, `wasLayoutChanged_test.go`).
2. Batch 3 — Durability (§1.1 + §1.6 + §5.1) — **blocked on the §1.1 owner decision.**
3. Batch 4 — Input validation (§1.5 + §1.7) — **blocked on the §1.5 ceilings decision.**
4. Batches 5–13 as listed in review §12.

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
> Where work left off: Batch 1 of 13 (Security, §8.1 + §8.2) is complete and
> verified — `x/text` v0.39.0, `x/net` v0.56.0, `x/sys` v0.46.0, both
> `govulncheck` scans clean, build and unit suite green, changes unstaged on
> branch `AD/refactoring-07-21`. Next up is Batch 2 (Correctness: §1.2, §1.3,
> §6.1, §5.2).
>
> Process the owner asked for: before starting each batch, ask whether it should
> be done at all — if not, document in the review file why it must not be
> attempted again — and ask any clarifying questions up front. After each batch,
> update `./.agent/session-carry-forward.md` and stop for owner review before
> continuing.
>
> See `./.agent/session-carry-forward.md` for the full handoff.
