# Session Carry-Forward — Batch 7 (Test-policy PR)

## 1. Session goal

Work through `todo/review-opus5-08-04.md` batch by batch, asking go/no-go and
clarifying questions before each batch, documenting rejections in place, and
stopping for owner review after every batch. This document covers **Batch 7 —
the Test-policy PR (§6.3 + §6.5)**.

## 2. Fixes applied

- **§6.3** — the three `internal/services/content_rules` unit tests no longer
  import `app/gui/constants`. They build `models.SidMapping` values from
  package-local SID string consts (`dragon_utopia`, `pandora_box`,
  `monty_hall`, `watchtower`). The owner chose literals over importing
  `internal/registry`, so the tests are decoupled from both layers.
- **§6.5a** — new depguard scope `test-unit-internal-no-gui` in
  [.golangci.yml](../.golangci.yml) denies `app/*` from
  `**/test/unit/internal/**`. Proven effective by temporarily reintroducing a
  GUI import (the linter flagged it) and then deleting the probe file.
- **§6.5b** — new [cmd/testlayoutcheck](../cmd/testlayoutcheck/main.go)
  enforces AGENTS.md §4.6.1 / §4.6.2 build-tag placement, wired into the
  `check-build` job of [pr-validation.yml](../.github/workflows/pr-validation.yml)
  as `go run ./cmd/testlayoutcheck .`. Current tree: **zero violations**.

## 3. Features added / changed

**The tag check is inverted relative to the review text — deliberately.**
§6.5 asked CI to fail unless the first non-blank line of every
`test/integration/**` and `test/performance/*_test.go` file is a build
constraint. That rests on the pre-2026-08-05 wording of AGENTS.md §4.6.1, the
same wording §6.1 was rejected over. Implemented literally it would fail
instantly on `test/integration/rmgTemplateModel_test.go` and
`test/performance/template_generation_test.go`, both correctly untagged. The
checker therefore enforces the rule **as it reads today**, failing when:

| Rule constant | Fails when |
| --- | --- |
| `RuleMissingIntegrationTag` | a `_test.go` file calls a `*_testexports.go` accessor without requiring `integration_test` |
| `RuleTaggedUnitTest` | any file under `test/unit/**` requires `integration_test` |
| `RuleMissingGuiTag` | any file under `test/integration/gui/**` does not require `gui` |
| `RuleTaggedProductionFile` | any file outside `test/` other than `*_testexports.go` requires `integration_test` |

**Implementation notes.** `TestLayoutChecker.Check(root)` walks the tree once,
parses every `.go` file with `go/parser`, extracts the `//go:build` line with
`go/build/constraint`, and resolves accessor names from the AST of every
`*_testexports.go` file — no grepping, no hard-coded accessor list. Skipped
directories: anything starting with `.`, plus `data`, `output`, `tmp`.

**Two false positives found while building it, both now regression-tested:**

- A naive `Eval(tag => tag != "integration_test")` reports `//go:build
  !wireinject` and `//go:build !windows` as *requiring* `integration_test`
  (`wire_gen.go`, `io_other.go`, `string_other.go`). Fixed by first checking
  that the tag is actually mentioned in the expression (`mentionsTag`).
- `test/test_helpers/integration_common/*.go` are tagged helper files, not
  production code, so the production rule is scoped to files outside `test/`.

## 4. File modifications

**New (6 files, 2 new top-level folders):**

| File | Purpose |
| --- | --- |
| `cmd/testlayoutcheck/main.go` | CLI entry point; exit 1 on violations, 2 on error |
| `cmd/testlayoutcheck/checker/testLayoutChecker.go` | `TestLayoutChecker` + the four rules |
| `cmd/testlayoutcheck/checker/goFile.go` | `goFile` — parsed file + constraint helpers |
| `cmd/testlayoutcheck/checker/violation.go` | `Violation` (pure data struct) |
| `test/unit/cmd/testlayoutcheck/checker/testLayoutChecker/check_test.go` | 13 tests |
| `test/unit/cmd/testlayoutcheck/checker/testLayoutChecker/newTestLayoutChecker_test.go` | constructor test |

**Modified (6):** `.golangci.yml` (depguard scope),
`.github/workflows/pr-validation.yml` (CI step), the three `content_rules`
unit tests, and `todo/review-opus5-08-04.md` (§0.3 rows, §6.3, §6.5, §12).

## 5. Tests added or updated

**New — 14 tests** in `test/unit/cmd/testlayoutcheck/checker/testLayoutChecker/`:
one happy path, one per rule, both false-positive regressions
(`!wireinject` / `!windows`, and the tagged `test_helpers` package), the
`!integration_test` inversion, skipped directories, violation path formatting,
unparseable source, and a missing root. Fixtures are written into `t.TempDir()`,
so the suite is cross-platform. Rule identities are asserted through the
exported `checker.Rule*` constants.

**Updated — 3 tests**: the §6.3 files. Assertions are unchanged; only the SID
source moved.

**Verification — all green:**

| Check | Result |
| --- | --- |
| `go build ./...` | pass |
| `go vet -tags=integration_test ./...` | pass |
| `go run ./cmd/testlayoutcheck .` | **test-layout check passed** |
| `go test -count=1 ./test/unit/...` | `unit=0` |
| `go test -tags=integration_test ./test/integration/...` | `integration=0` |
| `go test -tags='integration_test,gui' ./test/integration/gui/...` | `gui=0` |
| Unit coverage | **65.0%** — identical to the Batch 6 baseline |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **42 issues** (40 `gochecknoglobals`, 2 `dupl`) — unchanged baseline |

Coverage is unaffected because `-coverpkg=./internal/...,./app/...` does not
include `./cmd/...`, and `app/gui/constants` stays in the denominator via the
`mapSizes` / `victoryConditions` / `lookupSid` unit tests.

`coverage.txt`, `coverage.html` and `lcov.info` were regenerated (gitignored).

## 6. Git status snapshot

Branch `AD/refactoring-07-21`. **Nothing staged by the agent.** 6 modified
files plus the untracked `cmd/` and `test/unit/cmd/` folders — see §4. Batch
6's files no longer appear in `git status`, confirming the owner committed them.

## 7. Rejections / things not done

- **Widening the depguard scope to `test/integration`** — the owner initially
  chose this, then withdrew it when shown that four of the seven files there
  (`editorState`, `manualCastleReapply`, `stateExit`, `stateSaveAs`)
  legitimately import `app/gui/drivers` and friends. Depguard matches by path,
  not by intent, so a wider rule produces only false positives. Recorded in the
  review under §6.5 as **do not re-attempt**.
- **The literal §6.5 tag check** — not implemented; see §3 above and the §6.5
  deviation notes in the review.
- **A bash CI step** — the owner overrode the review's "keep it in CI, not in
  Go" so the check is cross-platform and locally runnable per AGENTS.md §2.2.
- **A VS Code task for the checker** — not added, to keep the diff inside the
  requested scope. Worth considering next to the existing lint/test tasks; the
  command is `go run ./cmd/testlayoutcheck .`.
- **§6.4** (the two 0%-coverage catalogues) — belongs to Batch 11, not touched.

## 8. Open questions

None for Batch 7.

Still blocked from earlier batches: **§7.1** (are direct pushes to master
intentional?), **§9.1** (public-API documentation decision), **§2.7**
(finish or remove the gladiator-arena preview), **§1.8** (output-dir
persistence: `.gen.json` vs machine-local), **§2.2** (scope of extracting
regeneration policy from `app/gui/drivers/`).

New, non-blocking: AGENTS.md §7 (Quick Reference) and §4.6.1 do not yet mention
`cmd/testlayoutcheck`. That is a documentation change and belongs to Batch 9
(§9.6), not here.

## 9. Next recommended actions

Batch 8 — **CI/security-posture PR**: §7.2 (top-level `permissions:`),
§7.3 (`actions/setup-go` version drift), §7.4 (the `tools/` module is never
built, tested, linted or tidy-checked in CI), §8.3 (scheduled vulnerability
scan). ⚠ §7.1 (direct pushes to `master`) needs the owner's policy decision
first.

Then batches 9–13 per review §12: docs, duplication cleanup, coverage, product
decisions, large refactors.

## 10. Carry-forward prompt

> Read `AGENTS.md` first. Hard rules, one line each: never modify `data/`,
> `internal/entities/template/` or `internal/registry/`; keep everything
> cross-platform (Windows + Linux, `path/filepath`, PowerShell chains with `;`);
> every change ships with tests and must not drop coverage; durable multi-session
> work gets a plan file under `plans/`; **never stage and never commit** — the
> owner reviews and commits.
>
> We are remediating the 46-finding review in `todo/review-opus5-08-04.md`,
> which defines 13 PR-sized batches in §12. Batches 1–7 are done: Security,
> Correctness, Durability, Input-validation, Performance, DI and Test-policy.
> Findings are marked `✅ FIXED` / `❌ WILL NOT FIX` **in place** in the review
> document — do not create a separate plan file for this.
>
> Workflow for every batch, without exception: (1) ask the owner whether the
> batch should be done at all; (2) if declined, document in the review file why
> it should not be attempted in future; (3) ask all clarifying questions up
> front; (4) implement; (5) rewrite `.agent/session-carry-forward.md`; (6) stop
> and wait for owner review.
>
> Next up is Batch 8, the CI/security-posture PR: §7.2, §7.3, §7.4, §8.3.
> §7.1 is blocked on an owner decision about direct pushes to `master`.
>
> Useful gotchas: `wire gen` and `golangci-lint-v2` write to stderr, so
> PowerShell shows a `NativeCommandError` even when they succeeded. Never pipe
> `go test` through `Select-Object -First N` — it kills the upstream process and
> fakes an exit code 1; redirect to a temp file and use `Select-String`. Adding
> the first test that imports a previously untested package can *lower* total
> `-coverpkg` coverage by enlarging the denominator, and CI hard-fails on any
> decrease; `cmd/` is outside `-coverpkg`, so tooling added there does not move
> the number. Run `go run ./cmd/testlayoutcheck .` before handing back — CI now
> fails on build-tag misplacement.
>
> See `.agent/session-carry-forward.md` for the full handoff.
