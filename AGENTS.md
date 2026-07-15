# Agent Operating Instructions

These instructions apply to any AI coding agent (GitHub Copilot, Claude, etc.)
working on the **HomMoe Custom Templates** repository. Follow them strictly.

---

## 1. Project Snapshot

- **Language / Toolchain:** Go 1.26.3, single module `github.com/Tariomka/hommoe_custom_templates`.
- **UI:** Gio (`gioui.org v0.9.0`) — immediate-mode desktop GUI.
- **Purpose:** Generate `.rmg.json` random-map templates for *Heroes of Might
  and Magic: Olden Era* and persist editor state as `.gen.json` files.
- **Entry point:** [main.go](main.go) → [app/gui/program.go](app/gui/program.go) (`StartApplication`).
- **Core generation:** [internal/services/template_generator/templateGenerator.go](internal/services/template_generator/templateGenerator.go).

Always read [README.md](README.md) and the relevant package before making
non-trivial changes.

---

## 2. Hard Rules — DO NOT VIOLATE

### 2.1 Read-only directories (game-data integrity)

The following folders contain **authoritative game data and the schema that
guarantees compatibility with Heroes of Might and Magic: Olden Era**. Editing,
renaming, reformatting, or "cleaning up" their contents will break the project
in production:

- [data/](data/) - including `ExampleTemplates/` and `GameData/` and `Images/`
- [internal/entities/template/](internal/entities/template/) - the `.rmg.json`
  output schema
- [internal/registry/](internal/registry) - game map generation template values/constants

You **MUST NOT**:

- Modify, delete, rename, reorder fields, change JSON tags, or "tidy"
  formatting in any file under those paths.
- Add new fields, remove unused ones, or change types of existing ones.
- Auto-format JSON files in `data/`.

You **MAY**:

- Read those files freely to understand the schema and game expectations.
- Reference their types from other packages.
- Propose changes in chat — but only the user may approve and apply them, and
  only after explicit confirmation.

### 2.2 Cross-platform compatibility (Windows + Linux)

The project must build and run on **both Windows and Linux**. Therefore:

- Use `path/filepath` (never hard-code `/` or `\`).
- Use `filepath.Join`, `filepath.Separator`, `os.UserConfigDir`, etc.
- Avoid OS-specific syscalls without a build-tag-guarded fallback.
- No CRLF-only assumptions; read files as bytes/strings and let Go normalize.
- Avoid shell-specific commands in code; if a tool/script is needed, provide
  both `.ps1` (Windows) and `.sh` (Linux) variants or use a Go program.
- Do not introduce dependencies that are Windows- or Linux-only without
  build tags (`//go:build windows` / `//go:build linux`).
- When suggesting terminal commands to the user, remember the workspace's
  default shell is **PowerShell on Windows** — chain with `;`, never `&&`.

### 2.3 Test coverage

- **Every new piece of non-trivial logic must ship with tests.** If you add
  or modify a function, add or update a test in [test/](test/) to cover it.
- If you discover existing untested code adjacent to your change, add tests
  for it as part of the change.
- **Every code change (new feature, fix, refactor) must check unit test code
  coverage.** Run the *"Go: Generate code coverage report"* task (or the command
  below) before and after the change and verify that every branch and logical
  unit of the code in question is covered, and that total coverage did not drop:

  ```powershell
  go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...
  go tool cover '-func=coverage.txt'
  ```

- Run `go test ./test/... -count=1` before declaring a task complete.
- The integration and performance suites are gated behind the `integration_test`
  build tag and are skipped by a plain `go test ./...`; run them explicitly
  with `go test -tags=integration_test ./test/integration/... ./test/performance/...`
  (see §4.6.1). Never make `integration_test` a global/default test tag.
- Tests must also be cross-platform (no hard-coded paths, no `\` separators,
  no shell-outs that exist only on one OS).

### 2.4 Real Work Planning

Turn planning into a durable, resumable artifact. The plan file - not the conversation -
is the source of truth: it records what to do, what's done, how it was verified,
and how to deploy. Any future agent can resume from it with zero prior context.

Use when planning multi-step / multi-session work that may outlive the
current session. Skip for trivial single-session tasks.

---

## 3. Workflow Rules

### 3.1 Implementation discipline

- Do not write the plan until scope is fully understood.
  Relentlessly ask the user questions until you both share a complete
  understanding with no gaps — treat an unasked question as a future bug.
- Don't stop at the first round; keep going until no ambiguity, assumption,
  or open decision remains. Probe edges: scope boundaries (in/out),
  dependencies, constraints, success criteria, data, environments,
  deployment, failure cases.
- Surface every assumption for the user to confirm. If an answer opens a
  new unknown, ask the follow-up — drill down recursively.
- Use `AskUserQuestion` for concrete choices. When done, summarize the full
  scope back and only proceed once the user confirms nothing is missing.

- Make the change the user asked for — nothing more.
- No drive-by refactors, no extra docstrings/comments on untouched code,
  no speculative error handling.
- Prefer editing existing files over creating new ones.
- Do not create markdown summary files unless explicitly requested.

### 3.2 Before editing

1. Read the target file(s) and at least one caller.
2. Confirm the change does not touch the read-only directories from §2.1.
3. Create a plan, see §2.4 for plan requirements and §4.7 for details.
4. For Go changes, check for existing tests in [test/](test/) and plan how
   you will extend them.

### 3.3 After editing

1. Run `go build ./...` and `go test ./test/...`.
2. If you touched editor internals or the gated suites, also run
   `go test -tags=integration_test ./test/integration/... ./test/performance/...`
   (see §4.6.1).
3. Report any new errors and fix them before handing back.
4. Briefly summarize: files touched, behaviour changed, tests added.

### 3.4 Picking the right models for workflows and subagents

When orchestrating subagents, pick the model per task using these ratings
(0–10, higher is better):

- **Cost** — relative cost to the user (higher = cheaper to run).
- **Intelligence** — how hard a problem you can hand the model unsupervised.
- **Taste** — code quality, API design, UI/UX and other subjective decisions.

| model           | cost | intelligence | taste |
|-----------------|------|--------------|-------|
| claude-fable-5  | 3    | 9            | 9     |
| gpt-5.6-sol     | 7    | 7            | 6     |
| gpt-5.6-terra   | 6    | 7            | 5     |
| gpt-5.5         | 5    | 6            | 5     |
| claude-opus-4.8 | 4    | 7            | 8     |
| sonnet-5        | 5    | 5            | 7     |

Application directives:

- These are defaults, not limits: if a cheaper model's output doesn't meet
  standards, rerun or redo the work with a smarter model without asking.
  Judge the output, not the price tag.
- Don't let cost prevent you from using the right model for the job.
  Instead, take advantage of cheaper options to gather information and try
  things before moving the work to a more expensive option.
- Anything user-facing (UI, API design, copy) or project-maintainability
  related requires taste > 7.
- Review of plans/implementations must be done by fable-5 or opus-4.8;
  optionally add gpt-5.5 as an extra independent perspective.
- **Never use Haiku models.**
- Match model to task shape: use cheap, high-cost-rating models (gpt-5.5,
  sonnet-5) for read-only exploration, searching, summarizing, and
  mechanical/repetitive edits; reserve fable-5/opus-4.8 for design
  decisions, tricky debugging, and final review.
- Parallelize independent exploration across cheap subagents rather than
  serializing everything through one expensive model.
- Give each subagent a self-contained brief (goal, constraints, expected
  output format) — subagents are stateless, and a weaker model with a
  precise brief beats a stronger model with a vague one.
- Escalate at most once per task; if two model tiers fail the same task,
  the brief is the problem — rewrite it instead of burning more runs.

---

## 4. Code Style

These rules are mandatory for any Go file the agent creates or substantially
edits. Do **not** retroactively rewrite untouched files just to comply, but
any file you *do* touch must leave the repo in conformance.

### 4.1 File & struct layout

- **One struct per file.** A file defines exactly one primary struct (its
  methods, constructors, and tightly-bound private helpers may live with it).
- **File name == struct name in `camelCase`** (lower-camel), with `.go`
  extension. Example: a struct `ZoneContentManager` lives in
  `zoneContentManager.go`. Test files mirror this: `zoneContentManager_test.go`.
- Interfaces, enums, and small value types that are not the primary struct
  belong in their own appropriately-named file (or in `internal/models/` /
  `internal/constants/` per §4.4).

### 4.2 Naming

- **No single-letter variables** and **no cryptic abbreviations**. Use
  descriptive names (`zoneIndex`, `playerCount`, `templatePath`).
- **Allowed exceptions** — only the well-established Go idioms:
  - `i`, `j` for loop indices only, but if a single for loop better use `index`
  - `err` for errors
  - `ok` for the comma-ok idiom
  - `ctx` for `context.Context`
  - Standard short receiver names are **not** allowed — see §4.3.

### 4.3 Method receivers

- The receiver for any method attached to a struct **must be named `this`**.

  ```go
  func (this *ZoneContentManager) Load(path string) error { ... }
  ```

- Apply this consistently across every method of the struct (Go vet warns on
  inconsistent receiver names, so do not mix `this` with anything else).

### 4.4 Package layout

Place new code in the package whose responsibility matches its role:

| Kind of code                                              | Location                                                                                                                                                                  |
| --------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| UI / rendering (Gio widgets, layouts, theming, input)     | [app/gui/](app/gui/)                                                                                                                                                      |
| Data structs / DTOs / factory functions (no behaviour)    | [internal/models/](internal/models/) + [internal/dtos/](internal/dtos/); read-only `.rmg.json` schema in [internal/entities/](internal/entities/) (`template/`, see §2.1) |
| Business logic, orchestrators, services                   | [internal/services/](internal/services/) + [internal/handlers/](internal/handlers/)                                                                                       |
| Constants, IDs, immutable lookup tables                   | [internal/constants/](internal/constants/) + [internal/registry/](internal/registry/)                                                                                     |
| Misc / cross-cutting utility functions                    | [internal/helpers/](internal/helpers/)                                                                                                                                    |

- If a struct or function has dependencies (helper structs, private types)
  **that are not used anywhere else**, nest them in a sibling folder next to
  the dependant file rather than polluting a shared package.
- Do not introduce new top-level packages without a clear reason.

### 4.5 UI vs. business logic separation

- Code under [app/gui/](app/gui/) **must contain only rendering
  logic** — widget composition, layout, input handling, view state.
- All business logic (validation, generation, transformation, persistence)
  lives in [internal/services/](internal/services/) (or `models/`,
  `helpers/`, `constants/` as appropriate) and is invoked by the GUI layer.
- If you find yourself writing an `if`/`switch` in a GUI file that decides
  *what* to do (rather than *how to draw*), extract it into a service.

### 4.6 Tests

**Unit test layout** — unit tests live under [test/unit/](test/unit/) and
mirror the full repository path of the implementation file:

- Each implementation `.go` file gets its **own folder** named after the file
  (without extension), located at `test/unit/<full/impl/path>/<fileName>/`.
- Each **public function or method** of that file gets its **own test file**
  named `<functionName>_test.go` (lower-camel).
- Example — `internal/services/settingsFileLoader.go` with functions
  `LoadSettingsFile` and `SaveSettingsFile`:

  ```
  test/unit/internal/services/settingsFileLoader/loadSettingsFile_test.go
  test/unit/internal/services/settingsFileLoader/saveSettingsFile_test.go
  ```

- The Go package for such a folder is `<fileName>_test`
  (e.g. `package settingsFileLoader_test`).

**Test naming** — `Test{Scenario}_{ExpectedBehavior}`, e.g.
`TestWhenFileIsMissing_ReturnsError`, `TestWhenTopologyIsHubAndSpoke_CreatesHubZone`.
For nested conditions X→Y→Z, name after the deepest relevant condition (Z);
if ambiguous, add the next-higher condition. Do not encode requirement IDs or
method names in test names — the folder/file already identify the method.

**Test body — triple-A**: every test contains `// Arrange`, `// Act`,
`// Assert` sections in that order. Pre-execution guard assertions
(e.g. `require.NoError` on setup) may live in the Arrange section.

**One unit per test**: each test verifies a single unit — a return value, or a
single expected interaction (mock called with expected values). Generally one
assertion per test; a single `assert.Equal` on a whole struct/slice counts as
one assertion. If a scenario allows several independent assertions, write a
separate test per assertion. Table-driven tests are allowed when each case
runs in a named `t.Run` subtest whose name follows the same
`{Scenario}_{ExpectedBehavior}` convention.

**Libraries**: only `testify` (`assert`, `require`, `mock`) for assertions and
mocking; use `gofakeit` for fuzzed input data wherever possible.

**Scope rules**:

- Pure data structs (no methods/logic) need no tests.
- [internal/registry/](internal/registry/) needs no tests (game-data constants).
- Private code is tested indirectly through public entry points — never add
  helpers/seams to implementation code just to make it testable. If code is
  unreachable through public APIs, record it in
  [todo/test_observations.md](todo/test_observations.md) instead.
- Code that is exercised indirectly by other tests still requires its **own**
  test folder with dedicated tests, so coverage can be assessed per file.
- `*_testexports.go` files (`//go:build integration_test`) must never be used
  by or tested in unit tests (see §4.6.1).
- Gio-UI-heavy code (widgets, dialogs, panels, window/event-loop code) that
  requires a `layout.Context`/window is covered by the integration suite, not
  unit tests; list such files in [todo/test_observations.md](todo/test_observations.md).

See §2.3 for coverage requirements.

### 4.6.1 The `integration_test` build tag (integration & performance only)

Some out-of-package tests need access to `editor.Window` internals (tab count,
selected tab, dialog state, programmatic load/save). Because those tests live in
a **different directory** than the `editor` package, the standard
`export_test.go` mechanism cannot reach them, and exposing the accessors as
normal methods would leak them into the production API.

The accessors therefore live in `*_testexports.go` files (e.g.
[app/gui/editor/window_testexports.go](app/gui/editor/window_testexports.go)),
guarded by `//go:build integration_test`. They compile **only** when the
`integration_test` tag is passed, so production builds (`go build ./...`) never
include them.

**Scope — this tag is for integration and performance tests ONLY:**

- The tag is used exclusively by the suites under [test/integration/](test/integration/)
  and [test/performance/](test/performance/). Every test file in those two
  directories carries `//go:build integration_test` at the top.
- **Do NOT** run the whole suite with the tag, and **do NOT** set it as a global
  `go.testTags`/`go.buildTags`. A normal `go test ./...` must stay tag-free; the
  two gated directories then compile to "[no test files]" and are skipped rather
  than failing.
- Only files under those two directories may reference the `integration_test`
  accessors. If another package needs an internal, do not widen this tag — add a
  test beside the code (`package X_test` in the same directory) instead.

**Running them:**

```powershell
# Default run — everything EXCEPT the gated dirs (no tag):
go test ./test/... -count=1

# Integration + performance only (tag scoped to these two dirs):
go test -tags=integration_test ./test/integration/... ./test/performance/... -count=1
```

In VS Code use the tasks in [.vscode/tasks.json](.vscode/tasks.json): *"go: test
(default, no integration_test)"* and *"go: test integration+performance
(integration_test)"*. gopls is configured with `-tags=integration_test` for
**analysis only** so the gated files still get IntelliSense — that does not cause
them to run.

### 4.7 Writing the Plan

Save to `plans/<descriptive-name>.md` in the repository root (create `plans/` if needed).
Use this self-documenting template:

```markdown
# <Work Title>

<1-2 sentence goal and scope.>

## For Future Agents
As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

## Phase 1: <Title>
Status: Not started   <!-- Not started | In progress | Complete -->

- [ ] <concrete, actionable item>
- [ ] <concrete, actionable item>

### Verification Plan
- <command/check the agent can run autonomously, with expected result>

### Phase Summary
_(write when phase completes)_

## Phase 2: <Title>
Status: Not started
- [ ] <actionable item>
### Verification Plan
- <autonomous check>
### Phase Summary
_(write when phase completes)_

## Final Recap
_(write when all phases complete: summary of the entire piece of work)_

## Deployment Plan
_(write when all phases complete: step-by-step deployment instructions)_
```

---

## 5. Session Length & Carry-Forward

To keep context windows healthy and answers high-quality, treat each session
as having a soft budget.

### 5.1 Session budget

- **Recommended length: <20 messages per session.**
- Around message **18**, warn the user that the session is approaching the
  recommended limit.
- At message **20** (or sooner if context feels saturated, tools start
  failing, or summaries become lossy), **stop taking new work** and produce a
  carry-forward document instead.

### 5.2 Carry-forward document

When the limit is reached — or when the user asks to wrap up — write a
**carry-forward prompt** that the next session can be started with. Save it
to `./.agent/session-carry-forward.md`

(Create the `./.agent/` folder if missing. This folder is for agent scratch
notes and is safe to add to `.gitignore` if the user wants.)

The document **must** include the following sections, in this order:

1. **Session goal** — one-line description of what the user asked for.
2. **Fixes applied** — bullet list, each linking to the changed file(s).
3. **Features added / changed** — bullet list with brief rationale.
4. **File modifications** — full list of files created/edited/deleted with
   one-line summaries.
5. **Tests added or updated** — list, plus pass/fail status of the last
   `go test ./test/...` run.
6. **Git status snapshot** — output (or summary) of `git status --short` and
   the current branch name. Note any unstaged or untracked work the next
   session will inherit.
7. **Rejections / things the user declined** — what was proposed but not
   merged, and why.
8. **Open questions** — anything blocked on user input or unresolved.
9. **Next recommended actions** — ordered to-do list for the next session.
10. **Carry-forward prompt** — a ready-to-paste prompt for the next agent
    session that re-establishes context. It must:
    - Reference this `AGENTS.md` file ("Read `AGENTS.md` first").
    - Re-state the §2 hard rules in one sentence each.
    - Summarize where work left off.
    - Point to `./.agent/session-carry-forward.md` for the full handoff.

The carry-forward document should be **self-contained** — a fresh agent with
no prior memory must be able to resume work from it alone.

### 5.3 During the session

- Keep a running mental (or `manage_todo_list`) checklist of items to fold
  into the carry-forward so the final write-up is fast and lossless.
- If a tool call rejects/errors or the user declines a suggestion, note it
  immediately so it lands in the *Rejections* section.

---

## 6. Communication Style

- Be brief. 1–3 sentences for simple answers; expand only for genuine
  complexity.
- Use Markdown. Wrap symbols in backticks; link files using
  `[path](path)` or `[path](path#L10-L20)` (workspace-relative, never
  inside backticks).
- No emojis unless the user asks.
- Never name internal tools to the user ("I'll run the tests", not "I'll
  use `runTests`").

---

## 7. Quick Reference

| Task                       | Command (Windows PowerShell / Linux bash)              |
| -------------------------- | ------------------------------------------------------ |
| Build                      | `go build ./...`                                       |
| Run GUI                    | `go run .`                                             |
| Run all tests              | `go test ./test/... -count=1`                          |
| Run integration/perf tests | `go test -tags=integration_test ./test/integration/... ./test/performance/... -count=1` |
| Run with race detector     | `go test -race ./test/...`                             |
| Unit test coverage report  | `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...` then `go tool cover '-func=coverage.txt'` (see §2.3; VS Code task *"Go: Generate code coverage report"*) |
| Lint (report only)         | `golangci-lint-v2 run ./... --issues-exit-code=0` (VS Code task *"Go: Get Linter Results"*) |
| Lint (auto-fix)            | `golangci-lint-v2 run ./... --issues-exit-code=0 --fix` (VS Code task *"Go: Run Linter"*; clears gci/gofmt/golines formatting findings — re-run to verify) |
| Format Go code             | `gofmt -w .` (never run on `data/`)                    |
| Tidy modules               | `go mod tidy`                                          |

---

**TL;DR:** Don't touch [data/](data/) or
[internal/entities/template/](internal/entities/template/). Stay cross-platform.
Cover everything you write with tests. Cap sessions at 17–20 messages and
hand off via `./.agent/session-carry-forward.md`.
