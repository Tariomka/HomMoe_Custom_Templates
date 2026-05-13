# Agent Operating Instructions

These instructions apply to any AI coding agent (GitHub Copilot, Claude, etc.)
working on the **HomMoe Custom Templates** repository. Follow them strictly.

---

## 1. Project Snapshot

- **Language / Toolchain:** Go 1.25.8, single module
  `github.com/Tariomka/hommoe_custom_templates`.
- **UI:** Gio (`gioui.org v0.9.0`) — immediate-mode desktop GUI.
- **Purpose:** Generate `.rmg.json` random-map templates for *Heroes of Might
  and Magic: Olden Era* and persist editor state as `.oetgs` files.
- **Entry point:** [main.go](main.go) → [internal/gui/gui.go](internal/gui/gui.go).
- **Core generation:** [internal/services/template_generator.go](internal/services/template_generator.go).

Always read [README.md](README.md) and the relevant package before making
non-trivial changes.

---

## 2. Hard Rules — DO NOT VIOLATE

### 2.1 Read-only directories (game-data integrity)

The following folders contain **authoritative game data and the schema that
guarantees compatibility with Heroes of Might and Magic: Olden Era**. Editing,
renaming, reformatting, or "cleaning up" their contents will break the project
in production:

- [data/](data/) — including `ExampleTemplates/` and `GameData/GeneratorData/`
- [internal/models/template/](internal/models/template/) — the `.rmg.json`
  output schema

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
- Run `go test ./test/...` (and, when sensible, `go test ./...`) before
  declaring a task complete.
- Tests must also be cross-platform (no hard-coded paths, no `\` separators,
  no shell-outs that exist only on one OS).

---

## 3. Workflow Rules

### 3.1 Implementation discipline

- Make the change the user asked for — nothing more.
- No drive-by refactors, no extra docstrings/comments on untouched code,
  no speculative error handling.
- Prefer editing existing files over creating new ones.
- Do not create markdown summary files unless explicitly requested.

### 3.2 Before editing

1. Read the target file(s) and at least one caller.
2. Confirm the change does not touch the read-only directories from §2.1.
3. For Go changes, check for existing tests in [test/](test/) and plan how
   you will extend them.

### 3.3 After editing

1. Run `go build ./...` and `go test ./test/...`.
2. Report any new errors and fix them before handing back.
3. Briefly summarise: files touched, behaviour changed, tests added.

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
  - `i`, `j`, `k` for loop indices
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

| Kind of code                                              | Location                                                                                       |
| --------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| UI / rendering (Gio widgets, layouts, theming, input)     | [internal/gui/](internal/gui/)                                                                 |
| Data structs / DTOs / factory functions (no behaviour)    | [internal/models/](internal/models/) — **except** the read-only `template/` subtree (see §2.1) |
| Business logic, orchestrators, services                   | [internal/services/](internal/services/)                                                       |
| Constants, IDs, immutable lookup tables                   | [internal/constants/](internal/constants/)                                                     |
| Misc / cross-cutting utility functions                    | [internal/helpers/](internal/helpers/)                                                         |

- If a struct or function has dependencies (helper structs, private types)
  **that are not used anywhere else**, nest them in a sibling folder next to
  the dependant file rather than polluting a shared package.
- Do not introduce new top-level packages without a clear reason.

### 4.5 UI vs. business logic separation

- Code under [internal/gui/](internal/gui/) **must contain only rendering
  logic** — widget composition, layout, input handling, view state.
- All business logic (validation, generation, transformation, persistence)
  lives in [internal/services/](internal/services/) (or `models/`,
  `helpers/`, `constants/` as appropriate) and is invoked by the GUI layer.
- If you find yourself writing an `if`/`switch` in a GUI file that decides
  *what* to do (rather than *how to draw*), extract it into a service.

### 4.6 Tests

- Tests live under [test/](test/) and **mirror the structure of
  `internal/`**. A file at `internal/services/foo.go` is tested by
  `test/services/foo_test.go`; `internal/models/bar/baz.go` by
  `test/models/bar/baz_test.go`.
- Test files follow the same `camelCase` filename rule, with the `_test.go`
  suffix.
- See §2.3 for coverage requirements.

---

## 5. Session Length & Carry-Forward

To keep context windows healthy and answers high-quality, treat each session
as having a soft budget.

### 5.1 Session budget

- **Recommended length: 10–20 messages per session.**
- Around message **15**, warn the user that the session is approaching the
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
    - Summarise where work left off.
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
| Run all tests              | `go test ./test/...`                                   |
| Run with race detector     | `go test -race ./test/...`                             |
| Format Go code             | `gofmt -w .` (never run on `data/`)                    |
| Tidy modules               | `go mod tidy`                                          |

---

**TL;DR:** Don't touch [data/](data/) or
[internal/models/template/](internal/models/template/). Stay cross-platform.
Cover everything you write with tests. Cap sessions at 10–20 messages and
hand off via `./.agent/session-carry-forward.md`.
