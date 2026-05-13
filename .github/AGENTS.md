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

## 4. Session Length & Carry-Forward

To keep context windows healthy and answers high-quality, treat each session
as having a soft budget.

### 4.1 Session budget

- **Recommended length: 10–20 messages per session.**
- Around message **15**, warn the user that the session is approaching the
  recommended limit.
- At message **20** (or sooner if context feels saturated, tools start
  failing, or summaries become lossy), **stop taking new work** and produce a
  carry-forward document instead.

### 4.2 Carry-forward document

When the limit is reached — or when the user asks to wrap up — write a
**carry-forward prompt** that the next session can be started with. Save it
to:

```
./.agent/session-carry-forward.md
```

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

### 4.3 During the session

- Keep a running mental (or `manage_todo_list`) checklist of items to fold
  into the carry-forward so the final write-up is fast and lossless.
- If a tool call rejects/errors or the user declines a suggestion, note it
  immediately so it lands in the *Rejections* section.

---

## 5. Communication Style

- Be brief. 1–3 sentences for simple answers; expand only for genuine
  complexity.
- Use Markdown. Wrap symbols in backticks; link files using
  `[path](path)` or `[path](path#L10-L20)` (workspace-relative, never
  inside backticks).
- No emojis unless the user asks.
- Never name internal tools to the user ("I'll run the tests", not "I'll
  use `runTests`").

---

## 6. Quick Reference

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
