# Prompt: In-Depth Repository Review

Paste everything below the line into a fresh agent session to produce a full,
principal-level review of this repository. Use a high-taste, high-intelligence
model (per AGENTS.md §3.4: fable-class or opus-class; never Haiku).

---

Read `AGENTS.md` first and follow it strictly for the entire session. You are
acting as a **senior/principal engineer performing an exhaustive review of the
whole repository**. Your ONLY deliverable is a review document — do **not** fix
anything, do not refactor, do not "improve while you're in there". Every
finding must be actionable later by a separate agent under the per-item
protocol (ask → plan → approve → implement + verify → owner commits → mark).

## Ground rules

1. **Read-only review.** The only file you create or modify is the review
   document: `todo/review-<model>-<MM-DD>.md`.
   Scratch verification programs are allowed only if deleted before you finish.
2. **Respect AGENTS.md hard rules** when judging code: `data/`,
   `internal/entities/template/`, and `internal/registry/` are read-only game
   data — findings inside them may only propose lint-config exclusions or
   owner-approved changes, never direct edits. `this` receivers, file-per-struct,
   camelCase file names, and the `test/unit` mirror layout are house style, not
   findings. The `integration_test` build tag is deliberately scoped to
   `test/integration/` + `test/performance/` — do not flag it.
3. **Verify every claim against the source.** No finding may be based on memory,
   prior reviews, or assumption. Quote the offending code, link the exact file
   and line range (`[file.go](../path/file.go#L10-L20)` — paths relative to
   `todo/`), and re-read the file immediately before writing the finding.
   If you cannot reproduce or point at concrete evidence, the finding does not
   go in the document.
4. **Audit prior reviews first.** If earlier review documents or observation
   files exist (e.g. `todo/review*.md`, `todo/test_observations.md`,
   `todo/backlog.md`), give a disposition for **every** prior item: Fixed (with
   evidence), Invalidated/accepted-as-convention (with reason), or Carried
   forward (re-verified, with a pointer to the new section). Nothing may be
   silently dropped. State explicitly which documents yours supersedes.
5. **Check the session-history notes** (`/memories/repo/`, `todo/*.md`, `.agent/*.md`)
   for known-stale findings, deliberate owner decisions (e.g. dead code kept on
   purpose), and items already marked "owner's responsibility" — do not
   re-report those as new findings; list them in the disposition section.

## Evidence gathering (run all of these, record versions and numbers)

Run and capture — these numbers anchor the review and become the baselines the
fix sessions will be held against:

- `go version`; toolchain versions of linters used.
- `go build ./...` and `go vet -tags=integration_test ./...`.
- `go test ./test/... -count=1` (default suite) and
  `go test -tags=integration_test ./test/integration/... ./test/performance/... -count=1`.
- Coverage: `go test -count=1 '-coverpkg=./internal/...,./app/...'
  '-coverprofile=coverage.txt' ./test/unit/...` then
  `go tool cover '-func=coverage.txt'` — record the total AND per-file gaps.
- Lint: `golangci-lint-v2 run ./... --issues-exit-code=0` — record the total
  issue count and per-linter breakdown (this becomes the lint baseline).
- `govulncheck ./...` (or confirm the CI job covers it).
- `git ls-files` cross-checked against `.gitignore` (committed artifacts,
  binaries, generated output); `git log` for suspicious recent churn.
- `go mod tidy` dry-run / `go.mod` vs CI Go-version consistency; note every
  module in the repo (`tools/` is separate).

## Review dimensions (cover ALL of these; add more if the code warrants)

1. **Bugs & correctness** — value-copy mutations lost, aliasing of slices/maps
   across API boundaries, swallowed errors, infinite-loop risks, unchecked type
   assertions/conversions, off-by-one, nil handling, silent coercion of invalid
   input, state machines that never reset, `os.Exit` in library/GUI paths,
   goroutine/data-race hazards, load/save round-trip fidelity (every UI field
   restored on load?), cross-platform path handling.
2. **Architecture** — layering violations (imports flowing the wrong way:
   `internal` → `app`, UI logic in services or business logic in GUI files),
   god objects/files (LOC counts), missing seams for testability (interfaces at
   handler/driver boundaries), packages that exist but are empty or dead,
   duplicated responsibilities between packages, per-call allocations of
   stateless collaborators.
3. **Duplicate code** — confirmed by `dupl`/manual diff; near-duplicates across
   topology/panel/builder families; repeated widget-row or provider vocabulary
   that wants a factory; duplicated lookup globals.
4. **Performance** — per-frame work in the GUI loop (allocations, reflection,
   geometry recomputation, pixel loops), missing dirty flags/caches, O(n²)
   where n can grow. Distinguish measured vs. reasoned findings.
5. **Readability & maintainability** — oversized functions (funlen/gocognit
   lists with a decomposition table: function → suggested split), magic
   strings/numbers wanting named constants, TODO/godox inventory with a
   disposition for each, misleading names, stale comments.
6. **Testing** — untested packages with coverage numbers, untestable code and
   WHY (missing seam vs. genuinely UI-bound), fixture quality, flakiness risks
   (run suspicious suites with `-count=20`), test-layout violations of AGENTS.md
   §4.6, test code importing layers it shouldn't.
7. **CI/CD** — missing gates (lint, race, coverage trend/floor, vuln scan,
   multi-OS, tag-gated suites compile check), release workflow hardening
   (pinned actions by SHA, checksums, `-trimpath`, version injection,
   concurrency), Dependabot, committed artifacts, `.gitattributes`/EOL policy,
   Go-version drift between `go.mod` and workflows.
8. **Linter disposition** — a table covering **every** current linter finding:
   linter, count, disposition (fix via §X / config exclusion / accepted). The
   fix sessions re-baseline against this table.
9. **Security & dependencies** — govulncheck results, risky file-system or
   registry access, injection surfaces in templates/JSON, unpinned or
   unmaintained dependencies.
10. **Docs & DX** — README/QUICKSTART accuracy vs. actual layout, stale paths,
    onboarding gaps, missing platform notes.

## Output format (this structure is the contract)

- **Header**: date, model, toolchain versions, lint issue count, which prior
  documents this supersedes.
- **Severity legend**: 🔴 High (bug/correctness/user-visible) · 🟠 Medium
  (architecture, performance, CI gaps) · 🟡 Low (readability, hygiene) ·
  ⚪ Informational.
- **§0 Disposition of prior reviews**: three subsections — Fixed ✅ (item +
  evidence link), Invalidated/accepted ✖ (item + reason), Carried forward ❗
  (item + new section pointer).
- **Numbered topical sections** (§1 Bugs … §7 CI/CD, extend as needed). Every
  finding gets its own `### N.M <severity> <title>` subsection containing, in
  order:
  1. Evidence — file/line links + quoted code snippet.
  2. Why it's wrong — the failure mode a user or maintainer hits.
  3. **Fix** — concrete, implementable instructions (code sketch where
     non-obvious), including the exact test files to add per AGENTS.md §4.6
     and any "if investigation shows X, do Y instead" branches.
  4. Any owner-decision flag (⚠ protected dir, destructive action, style call)
     so the fix session knows to ask first.
- **Linter-disposition table** (§8-style) covering the full current count.
- **Suggested execution order** (§9-style): group findings into safe PR-sized
  batches, bugs first, note which items block others and which need owner
  decisions.
- Keep item numbering stable once written — fix sessions reference items as
  `§N.M` and mark them `✅ FIXED` in place.

## Quality bar

- Be thorough (for example: ~50 findings across 9 dimensions with per-item fix
  plans). Depth beats breadth: one verified 🔴 with a correct fix plan
  outweighs ten vague 🟡s.
- Actively hunt for the classes of bug that survived previous reviews: value
  semantics (struct copies mutated and discarded), aliasing (shallow copies
  sharing slices), UI state not surviving save/load round-trips, error paths
  that silently continue, and concurrency around lazily-initialized globals.
- For every 🔴/🟠, ask "what test would have caught this?" and put that test in
  the fix plan.
- Mark anything you checked and found to be FINE that a reader might expect to
  be flagged (false-positive protection for future reviews) in a short
  "verified non-issues" list.
- End with the measured baselines (coverage %, lint count and composition,
  test pass state) so fix sessions can enforce no-regression.

When finished: report the document path, the finding counts per severity, and
the top three highest-impact items. Do not start fixing anything.
