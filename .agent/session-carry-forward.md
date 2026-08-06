# Session Carry-Forward — Batch 8 (CI / security posture)

## 1. Session goal

Work through `todo/review-opus5-08-04.md` batch by batch, asking go/no-go and
clarifying questions before each batch, documenting rejections in place, and
stopping for owner review after every batch. This document covers **Batch 8 —
the CI/security-posture PR (§7.2, §7.3, §7.4, §8.3)** and the long-blocked §7.1
owner decision.

## 2. Fixes applied

- **§7.2 — workflows had no top-level `permissions:`.** Added
  `permissions:` / `contents: read` to
  [pr-validation.yml](../.github/workflows/pr-validation.yml) and
  [release.yml](../.github/workflows/release.yml), and included it in both new
  workflows. Existing per-job narrowing is untouched: `code_coverage` still adds
  `actions: read` + `pull-requests: write`, and the release publish job still
  raises to `contents: write`.
- **§7.3 — `actions/setup-go` version drift.** The composite action
  [setup-steps/action.yml](../.github/workflows/setup-steps/action.yml) moved
  from `@v6` to `@v7`, matching the direct use in `check-go-mod`. All ten jobs
  are now on `@v7`.
- **§7.4 — the `tools/` module was never verified in CI.** New path-filtered
  [tools-validation.yml](../.github/workflows/tools-validation.yml) runs
  `go mod tidy -diff` with `working-directory: tools`, and
  [tools/go.mod](../tools/go.mod) now declares `go 1.26.5` instead of `1.26.3`.
- **§8.3 — the vulnerability gate only ran on pull requests.** New
  [security-scan.yml](../.github/workflows/security-scan.yml) re-scans the
  default branch weekly (`cron: '0 6 * * 1'`).

## 3. Features added / changed

No product behaviour changed. Two new CI workflows:

**`tools-validation.yml`** — triggers on `pull_request` and `push: master`
restricted to `paths: ['tools/**', '.github/workflows/tools-validation.yml']`.
The owner asked that the `tools/` check run only when `tools/` changes, and
GitHub's `paths:` filter exists only at the workflow-trigger level, never per
job — so the step could not live inside `check-go-mod` as the review proposed.

> ⚠ **Do not add this workflow to branch protection as a required status
> check.** Path-filtered workflows never report on PRs that do not match the
> filter, which would block those PRs forever. A comment at the top of the file
> records this.

**`security-scan.yml`** — `schedule` only (Mondays 06:00 UTC),
`permissions: contents: read`, same `golang/govulncheck-action@v1` invocation as
the PR job. A dedicated workflow beat adding `schedule:` to `pr-validation.yml`
because the existing job's `if: github.event_name == 'pull_request'` would skip
it on a scheduled run. Scheduled workflows run against the default branch, which
is exactly what needs re-scanning.

> ⚠ GitHub disables scheduled workflows after 60 days of repository inactivity;
> the owner declined a `workflow_dispatch` companion trigger.

## 4. File modifications

| File | Change |
| --- | --- |
| `.github/workflows/pr-validation.yml` | M — top-level `permissions: contents: read` after the `on:` block |
| `.github/workflows/release.yml` | M — top-level `permissions: contents: read` after `concurrency:` |
| `.github/workflows/setup-steps/action.yml` | M — `actions/setup-go@v6` → `@v7` |
| `.github/workflows/tools-validation.yml` | **NEW** — path-filtered `go mod tidy -diff` for the tools module |
| `.github/workflows/security-scan.yml` | **NEW** — weekly scheduled `govulncheck` |
| `tools/go.mod` | M — `go 1.26.3` → `go 1.26.5` |
| `todo/review-opus5-08-04.md` | M — §7.1 ❌ WILL NOT FIX; §7.2/§7.3/§7.4/§8.3 ✅ FIXED; §0.3 rows for §7.1/§7.3; §12 item 8; §7.1 dropped from the open-decisions list |

## 5. Tests added or updated

**None, and none are required.** Batch 8 changed only YAML workflow files and a
`go` directive — no Go logic was added or modified, so AGENTS.md §2.3 does not
apply. Verification used `actionlint` plus the full standard suite.

| Check | Result |
| --- | --- |
| `actionlint` on the four workflow files | **exit 0, no findings** |
| `go mod tidy -diff` (root) | exit 0 |
| `go mod tidy -diff` (`tools/`) | exit 0 |
| `go build ./...` | pass |
| `go vet -tags=integration_test ./...` | pass |
| `go run ./cmd/testlayoutcheck .` | **test-layout check passed** |
| `go test -count=1 ./test/unit/...` | `unit=0` |
| `go test -tags=integration_test ./test/integration/...` | `integration=0` |
| `go test -tags='integration_test,gui' ./test/integration/gui/...` | `gui=0` |
| Unit coverage | **65.0%** — unchanged |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **42 issues** (40 `gochecknoglobals`, 2 `dupl`) — unchanged baseline |

`actionlint` is not installed locally; it was run as
`go run github.com/rhysd/actionlint/cmd/actionlint@latest -no-color <files>`.
Running it with **no** file arguments also emits five errors for
`.github/workflows/setup-steps/action.yml` — a **pre-existing false positive**:
actionlint treats every file under `.github/workflows/` as a workflow, and that
file is a composite action. Pass the four workflow paths explicitly for a clean
run.

`coverage.txt`, `coverage.html` and `lcov.info` were regenerated and are
byte-identical to the committed versions.

## 6. Git status snapshot

Branch `AD/refactoring-07-21`. **Nothing staged by the agent.**

```
 M .github/workflows/pr-validation.yml
 M .github/workflows/release.yml
 M .github/workflows/setup-steps/action.yml
 M todo/review-opus5-08-04.md
 M tools/go.mod
?? .github/workflows/security-scan.yml
?? .github/workflows/tools-validation.yml
```

Batch 7's files no longer appear, confirming the owner committed them.

## 7. Rejections / things not done

- **§7.1 ❌ WILL NOT FIX — owner ruling, recorded in the review as do-not-
  re-attempt.** The finding assumed the `push: master` trigger lets unreviewed
  code reach `master`. It does not: branch protection forbids direct pushes, so
  the trigger only fires *after* a PR merges, on code that already passed all
  ten jobs. The reduced post-merge subset (build, Windows build, unit tests,
  integration tests) exists solely to catch the rare concurrent-merge conflict;
  re-running tidy/lint/vulnerability/race/coverage there costs CI minutes for no
  signal. The proposed rename away from "PR Tests" is declined for the same
  reason.
- **§8.3's fix item 1** ("remove the PR gate so pushes are scanned") is void —
  it was a §7.1 dependency. Only the scheduled scan was implemented.
- **`workflow_dispatch` on `security-scan.yml`** — offered; the owner chose the
  schedule-only option.
- **Building/vetting the `tools/` module in CI** — offered and dropped: the
  module declares no packages of its own, only `tool` directives, so there is
  nothing to compile.
- **Putting the tools tidy step inside `check-go-mod`** (the review's literal
  fix) — replaced by the path-filtered workflow per the owner's "only trigger
  when `tools/` changes" requirement.
- **AGENTS.md §7 / §4.6.1 still do not mention `cmd/testlayoutcheck`** (added in
  Batch 7). That documentation change belongs to Batch 9 (§9.6), not here.

## 8. Open questions

None for Batch 8.

Still blocked from earlier batches: **§9.1** (public-API documentation
decision — gates Batch 9), **§2.7** (finish or remove the gladiator-arena
preview; §9.5 must agree with the outcome), **§1.8** (output-dir persistence:
`.gen.json` vs machine-local), **§2.2** (scope of extracting regeneration policy
from `app/gui/drivers/`).

## 9. Next recommended actions

Batch 9 — **Docs PR**: §9.1–§9.6, optionally §9.7, then update repository
memory. ⚠ §9.1 needs the owner's public-API decision *before* work starts, and
§9.5 must agree with whatever §2.7 decides. Fold in the AGENTS.md
`cmd/testlayoutcheck` documentation noted in §7 above.

Then batches 10–13 per review §12: duplication cleanup, coverage, product
decisions, large refactors (the last needs a plan file under `plans/` per
AGENTS.md §4.7).

## 10. Carry-forward prompt

> Read `AGENTS.md` first. Hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/`; keep everything
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chains with `;`);
> every change ships with tests and must not drop coverage; durable multi-session
> work gets a plan file under `plans/`; **never stage and never commit** — the
> owner reviews and commits.
>
> We are remediating the 46-finding review in `todo/review-opus5-08-04.md`,
> which defines 13 PR-sized batches in §12. Batches 1–8 are done: Security,
> Correctness, Durability, Input-validation, Performance, DI, Test-policy and
> CI/security-posture. Findings are marked `✅ FIXED` / `❌ WILL NOT FIX` **in
> place** in the review document (§6.1 and §7.1 are rejected with rationale) —
> do not create a separate plan file for this.
>
> Workflow for every batch, without exception: (1) ask the owner whether the
> batch should be done at all; (2) if declined, document in the review file why
> it should not be attempted in future; (3) ask all clarifying questions up
> front; (4) implement; (5) rewrite `.agent/session-carry-forward.md`; (6) stop
> and wait for owner review.
>
> Next up is Batch 9, the Docs PR: §9.1–§9.6 (+ optional §9.7). §9.1 is blocked
> on an owner decision about whether the QUICKSTART programmatic example should
> document a supported public API or be deleted — ask that first. Also fold in
> the AGENTS.md §7 / §4.6.1 update for `cmd/testlayoutcheck`, which Batch 7
> added but did not document.
>
> Baseline that must not regress: unit coverage **65.0%**,
> `golangci-lint-v2 run ./... --issues-exit-code=0` → **42 issues**,
> `go run ./cmd/testlayoutcheck .` → 0 violations, all three test suites green.
>
> Useful gotchas: `wire gen` and `golangci-lint-v2` write to stderr, so
> PowerShell shows a `NativeCommandError` even when they succeeded. Never pipe
> `go test` through `Select-Object -First N` — it kills the upstream process and
> fakes an exit code 1; redirect to a temp file and use `Select-String`. Adding
> the first test that imports a previously untested package can *lower* total
> `-coverpkg` coverage by enlarging the denominator, and CI hard-fails on any
> decrease; `cmd/` is outside `-coverpkg`. Lint workflow YAML with
> `go run github.com/rhysd/actionlint/cmd/actionlint@latest -no-color` and pass
> the workflow files explicitly — a bare run false-positives on the composite
> action `setup-steps/action.yml`.
>
> See `.agent/session-carry-forward.md` for the full handoff.
