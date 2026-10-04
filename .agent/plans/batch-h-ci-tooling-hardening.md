# Batch H: CI/tooling hardening (review §6.2, §6.3, §6.5)

Close review items §6.2 (one linter version), §6.3 (LF policy for module/checksum files) and
§6.5 (release tag interpolated into shell source) of
[review-gpt-6-astra-09-07.md](../backlog/review-gpt-6-astra-09-07.md). Configuration-only:
no Go production code or tests change.

## For Future Agents
As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

**Starting state** (2026-10-04): branch `AD/pbi_resolution`, HEAD `9d302a5` ("init",
AGENTS.md only) on top of `master` `ec42217` (Batches G and J, PR #49). Clean tree.

**Hard rules in force:** never stage, unstage, commit, push, stash, switch branches or touch
worktrees (the owner does all Git mutation; the only exception is H4's working-copy refresh of
the four verified-unmodified module/checksum files); never touch `data/`,
`internal/entities/template_entity/`, `internal/registry/`, generated Wire or the output
directory; no bulk rewrites; never dispatch or publish a release to test anything.

### Owner decisions (2026-10-04)

| ID | Decision |
| --- | --- |
| H1 | Single linter version **v2.13.1** (matches CI and the local binary). |
| H2 | `tools/go.mod` is the single source; the CI lint job reads the version from it. |
| H3 | Bump via `go get …@v2.13.1` + `go mod tidy` in `tools/` only; accept whatever indirect bumps MVS forces; no deliberate upgrade of `wire`, `gcov2lcov` or anything else. |
| H4 | Add `eol=lf` for `go.mod`/`go.sum`; the agent refreshes only the four working copies (no staging). |
| H5 | Accepted tag format: `^v[0-9]+\.[0-9]+(\.[0-9]+)?(-[0-9A-Za-z.]+)?$` (every existing tag matches, including `v0.1` and `v0.3.9-alpha.1`). |
| H6 | Validation is a bash step in a new `validate` job of `release.yml`; the tag reaches shell only through `env`. |
| H7 | A tag with a `-` suffix is published as a GitHub prerelease; otherwise `prerelease: false` as today. |
| H8 | Checkout `ref` unchanged; no explicit tag-existence check; concurrency group unchanged. |
| H9 | Local tag-case check is a one-off, throwaway run through Git for Windows bash (`C:\Program Files\Git\bin\bash.exe`). Nothing is committed for it. No WSL use. |
| H10 | No baseline/rerun of build, unit tests, coverage, testlayoutcheck: no Go code changes (owner override of AGENTS.md §2.3/§3.3 for this batch). |
| H11 | Windows-only local verification; Linux tidy and workflow execution are verified by the owner's PR CI. |
| H12 | Independent plan and implementation review by gpt-6.1-sol. |

**Plan review** (GPT-6.1 Sol): APPROVE WITH CHANGES; all 4 findings applied (fail-fast
non-empty lint version lookup, H4 exception stated, step id/output mappings, actionlint static
check). **Owner plan approval:** pending.

**Out of scope:** `wire`/`gcov2lcov` versions, other workflow hardening (static
`matrix`/`env` interpolations in `run:`), README/docs (Batch I), `.golangci.yml`, any Go code.

## Phase 1: LF policy for module/checksum files (§6.3)
Status: Not started

- [ ] In [.gitattributes](../../.gitattributes) add `go.mod text eol=lf` and
      `go.sum text eol=lf` (slash-less patterns match at any depth, so `tools/` is covered).
      Reword the comment so it no longer claims `*.go` covers modules.
- [ ] Refresh exactly `go.mod`, `go.sum`, `tools/go.mod`, `tools/go.sum`: confirm each is
      unmodified (`git status --short -- <path>` empty) and `i/lf` in `git ls-files --eol`,
      then `Remove-Item <path>` and `git checkout -- <path>`. Index is already LF, so no
      content diff results; the explicit attribute overrides `core.autocrlf`, and no index
      normalization or staging is needed.

### Verification Plan
- `git ls-files --eol go.mod go.sum tools/go.mod tools/go.sum` → all `i/lf w/lf attr/text eol=lf`.
- `git status --short` → only `.gitattributes` (plus this plan) modified.
- `go mod tidy -diff` in root and in `tools/` → both exit 0, no output.

### Phase Summary
_(write when phase completes)_

## Phase 2: Single linter version (§6.2)
Status: Not started

- [ ] In `tools/`: `go get github.com/golangci/golangci-lint/v2@v2.13.1`, then `go mod tidy`.
      Record which other requirements moved (expected: indirect only).
- [ ] In [pr-validation.yml](../../.github/workflows/pr-validation.yml) `run-gci-lint` job,
      add a step (id `lint-version`, `working-directory: tools`, `shell: bash`) before the
      action. The lookup is a standalone assignment so a failure fails the step, and an empty
      value is rejected (an empty `version` would make the action pick the latest release):
      ```bash
      version=$(go list -m -f '{{.Version}}' github.com/golangci/golangci-lint/v2)
      if [[ -z "$version" ]]; then
        echo "::error::golangci-lint version not found in tools/go.mod"
        exit 1
      fi
      echo "version=$version" >> "$GITHUB_OUTPUT"
      ```
      Set the action's `version: ${{ steps.lint-version.outputs.version }}`. Keep
      `install-mode: binary` and the explicit `--disable=godox,dupl,unparam,gochecknoglobals`.
      One-line comment noting the version lives in `tools/go.mod`.

### Verification Plan
- `git diff tools/go.mod` → `golangci-lint/v2 v2.13.1`; `wire` and `gcov2lcov` direct
  requirements unchanged unless MVS forced them (record if so).
- `go mod tidy -diff` in `tools/` → exit 0.
- In `tools/`: `go list -m -f '{{.Version}}' github.com/golangci/golangci-lint/v2` → `v2.13.1`.
- In `tools/`: `go build -o <tmp>/golangci-lint.exe github.com/golangci/golangci-lint/v2/cmd/golangci-lint`
  then `go version -m <tmp>/golangci-lint.exe` → module version `v2.13.1`; delete the binary.
- Root `go mod tidy -diff` → exit 0 (root module untouched).
- `golangci-lint-v2 run ./... --issues-exit-code=0` → 0 issues.
- actionlint (see Phase 3) also covers `pr-validation.yml`.

### Phase Summary
_(write when phase completes)_

## Phase 3: Release tag validation (§6.5)
Status: Not started

- [ ] In [release.yml](../../.github/workflows/release.yml) add a first job `validate`
      (`runs-on: ubuntu-latest`, inherits `contents: read`, no checkout) with job outputs
      `tag: ${{ steps.release-tag.outputs.tag }}` and
      `prerelease: ${{ steps.release-tag.outputs.prerelease }}`. Its single step
      (`id: release-tag`, `shell: bash`) sets
      `env: TAG: ${{ github.event.inputs.tag || github.ref_name }}` and runs:
      ```bash
      if [[ ! "$TAG" =~ ^v[0-9]+\.[0-9]+(\.[0-9]+)?(-[0-9A-Za-z.]+)?$ ]]; then
        echo "::error::Release tag must match vMAJOR.MINOR[.PATCH][-PRERELEASE]"
        exit 1
      fi
      prerelease=false
      if [[ "$TAG" == *-* ]]; then prerelease=true; fi
      echo "tag=$TAG" >> "$GITHUB_OUTPUT"
      echo "prerelease=$prerelease" >> "$GITHUB_OUTPUT"
      ```
      The rejected value is not echoed (avoids workflow-command injection in logs).
- [ ] `build`: add `needs: validate`; in "Build binary" add `VERSION: ${{ needs.validate.outputs.tag }}`
      to `env` and use `-X main.version=${VERSION}` inside the existing double-quoted
      `-ldflags`. Checkout `ref` unchanged (H8).
- [ ] `release`: `needs: [validate, build]`; `tag_name` and `name` from
      `${{ needs.validate.outputs.tag }}`; `prerelease: ${{ needs.validate.outputs.prerelease == 'true' }}`.
- [ ] No other changes: triggers, dispatch input, concurrency, permissions, artifact and
      checksum steps stay as they are.

### Verification Plan
- Throwaway (not committed, deleted afterwards): run the exact step body through
  `C:\Program Files\Git\bin\bash.exe` with `GITHUB_OUTPUT` pointing at a temp file, per case:
  - accept with `prerelease=false`: `v0.1`, `v0.3.8`, `v1.0.0`, `v10.20.30`;
  - accept with `prerelease=true`: `v0.3.4-beta`, `v0.3.6-beta.2`, `v0.3.9-alpha.1`;
  - reject (exit 1, temp output empty): empty, `v1`, `1.0.0`, `v1.0.0.0`, `v1.0.0-`,
    ` v1.0.0`, `v1.0.0 `, `v1.0.0"`, `v1.0.0'`, `v1.0.0;id`, `v1.0.0$(id)`, `` v1.0.0`id` ``,
    `v1.0.0\nv1.0.1` (embedded newline), `v1.0.0-beta_1`, `refs/tags/v1.0.0`.
  - Every existing `git tag --list` entry is accepted.
- `git diff .github/workflows/release.yml` reviewed: no `${{ … }}` of the tag remains
  inside any `run:` block (`Select-String` for `inputs.tag|ref_name` only hits `env:`,
  checkout `ref`, and nothing in `run:`).
- Static validation: `go run github.com/rhysd/actionlint/cmd/actionlint@latest .github/workflows/release.yml .github/workflows/pr-validation.yml`
  (temporary download into the module cache; no repo file changes) → no findings in the
  edited jobs. Pre-existing findings in untouched lines are recorded, not fixed.
- Manually inspect the output plumbing: step `id` ↔ job `outputs` ↔ `needs.validate.outputs.*`
  in both `build` and `release`, and `needs` lists.
- The release workflow is not executed by PR CI. Runtime confirmation is the owner's next
  legitimate release; never dispatch or publish a test release.

### Phase Summary
_(write when phase completes)_

## Phase 4: Review and close-out
Status: Not started

- [ ] Independent implementation review by gpt-6.1-sol; apply or record findings.
- [ ] After the owner reviews/commits: mark §6.2, §6.3, §6.5 `✅ FIXED` in place in the
      review, complete §9 row H, refresh the progress line (expected 21 fixed / 11 remaining);
      never renumber; §8 untouched.
- [ ] Update `.agent/session-carry-forward.md`.

### Verification Plan
- Review verdict recorded; `git status --short` lists only the intended files.

### Phase Summary
_(write when phase completes)_

## Final Recap
_(write when all phases complete: summary of the entire piece of work)_

## Deployment Plan
_(write when all phases complete: step-by-step deployment instructions)_
