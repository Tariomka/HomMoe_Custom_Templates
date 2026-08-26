# Heroes of Might and Magic: Olden Era - Custom Templates

A desktop GUI for designing and generating `.rmg.json` random map templates for
**Heroes of Might and Magic: Olden Era**.

---

## 1. Overview

This is a highly performant, fully offline, multi-platform, native Golang application,
rendered using immediate-mode rendering library GioUI.
Most of the projects don't really accomplish all of those points so I strive to keep the project
simple, well structured and most importantly performant, so you need to do it as well. Clean and
concise architecture is also an aspiration so that the codebase would be understandable and extendable.

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
- Avoid shell-specific commands in code; if a tool/script is needed, use a Go program.
- Do not introduce dependencies that are Windows- or Linux-only without build tags (`//go:build windows` / `//go:build linux`).

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

- Run `go test ./test/unit/... -count=1` before declaring a task complete.
- The integration and ui test suites are gated behind the `integration_test`
  and `integration_test,gui` build tags respectfully and are skipped by a
  plain `go test ./...`; run them explicitly with
  `go test -tags=integration_test ./test/integration/...` and
  `go test -tags='integration_test,gui' ./test/integration/gui/...` respectfully.
  (see §4.6.1 and §4.6.2). Never make `integration_test` or `gui` a
  global/default test tag.
- **Build tags are applied per file, never per directory.** Tag only the files
  that genuinely need the tag; a test that compiles and passes without a tag
  must not carry one just because a sibling file in the same directory does.
- Tests must also be cross-platform (no hard-coded paths, no `\` separators,
  no shell-outs that exist only on one OS).

### 2.4 Real Work Planning

Turn planning into a durable, resumable artifact. The plan file - not the conversation -
is the source of truth: it records what to do, what's done, how it was verified,
and how to deploy. Any future agent can resume from it with zero prior context.

Use when planning multi-step / multi-session work that may outlive the
current session. Skip for trivial single-session tasks.

### 2.5 Staging and Committing

You **MUST NOT** stage any changes you do and/or commit them to origin or any other brach.
If you notice staged changes, **NEVER** unstage them - it is done by the author
after he reviews and ensures the changes are correct.

### 2.6 Bulk rewrites

- Never run a bulk in-place rewrite across the repository.
- To normalize formatting or line endings, run `gofmt -w` on an **explicit** list
  from `gofmt -l` - gofmt converts CRLF to LF and cannot mangle content.

### 2.7 Output path is a hard requirement of the game

Heroes of Might and Magic: Olden Era only reads random-map templates from **its
own templates directory**. A `.rmg.json` written anywhere else is not merely in
an unusual place — **the game will never find it**, and the user is left with a
file they cannot locate or place correctly after the fact.

The directory is therefore **not** a user preference. It is a property of the
*machine*: it differs per device with the Steam library location, the OS, and
the Proton/native layout, and it changes whenever the game is moved or
reinstalled. Per-launch auto-detection via
[internal/helpers/io.go](internal/helpers/io.go) (`FindOldenEraTemplatesDir`) is
the correct design precisely because it is self-healing.

You **MUST NOT**:

- Change where `.rmg.json` (or its preview `.png`) is written, or add any
  behaviour that lets a template land outside the detected templates directory
  by default.
- **Persist the output directory** in any form — not in `EditorStateDto` /
  `.gen.json`, not in an `os.UserConfigDir()` preferences file, not anywhere
  else. A stored path is invalid on any other machine and goes stale on the
  current one. Do not re-propose it.
- Rework `outputPath` on `drivers.State` into a "remembered setting". The
  in-app folder picker is a deliberate **single-session escape hatch** for
  layouts the detector does not recognise yet — nothing more.

If detection turns out to be inadequate on some platform, the fix is to improve
`FindOldenEraTemplatesDir`, **never** to store a path.

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

1. Run `go build ./...` and `go test ./test/unit/...`.
2. If you touched editor internals or the gated suites, also run
   `go test -tags='integration_test,gui' ./test/integration/...` (see §4.6.1).
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
| claude-opus-5   | 5    | 8            | 9     |
| claude-fable-5  | 2    | 9            | 10    |
| gpt-5.6-sol     | 6    | 7            | 6     |
| kimi-k3         | 7    | 7            | 7     |
| gpt-5.6-terra   | 7    | 7            | 6     |
| grok-4.6        | 8    | 6            | 5     |
| claude-opus-4.8 | 4    | 7            | 7     |
| gpt-5.5         | 4    | 6            | 4     |
| sonnet-5        | 4    | 3            | 5     |

Application directives:

- These are defaults, not limits: if a cheaper model's output doesn't meet
  standards, rerun or redo the work with a smarter model without asking.
  Judge the output, not the price tag.
- Don't let cost prevent you from using the right model for the job.
  Instead, take advantage of cheaper options to gather information and try
  things before moving the work to a more expensive option.
- Anything user-facing (UI, API design, copy) or project-maintainability
  related requires taste > 7.
- Review of plans/implementations must be done by opus-5 preferably
  (use fable-5 sparingly as it is much more costly);
  optionally add gpt-5.6-terra as an extra independent perspective.
- **Never use Haiku models.**
- Match model to task shape: use cheap, high-cost-rating models (gpt-5.6-terra,
  kimi-k3, grok-4.6, gpt-5.5) for read-only exploration, searching, summarizing, and
  mechanical/repetitive edits; reserve opus-5/fable-5 for design
  decisions, tricky debugging, and final review.
- Parallelize independent exploration and/or action execution
  (like running tests) across cheap subagents rather than serializing
  everything through one expensive model.
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
  - `i`, `j` for loop indices only, but if a single for loop, better use `index`
  - `x` in simple `linq.Select` queries (single returns with no closures, example `return x.Label`)
  - `err` for errors
  - `ok` for the comma-ok idiom
  - `ctx` for `context.Context`
  - Standard short receiver names are **not** allowed — see §4.3.

### 4.2.1 Interfaces

- Interface types **must use `I` prefix** (`IDialog`, `IPanel`, `IBackend`).

  ```go
  type IBackend interface {
    ITemplateWorkflowHandler
    IStatePersistenceHandler
    IStateValidationHandler
    IPreviewHandler
    IContentRuleHandler
    IZoneEditorHandler
  }
  ```

- Interface file names **must use `Interface` suffix** (`dialogInterface.go`,
  `panelInterface.go`, `backendInterface.go`).

- Interfaces must be separated from the implementation files.

- Apply this consistently across every interface file.

### 4.2.2 Interface placement

Where an interface file lives is decided by counting the **concrete
implementation files in that package that require an interface** — not the
number of interfaces:

1. **Fewer than 5** → declare the interface in the **same package**, in its own
   `*Interface.go` file.
2. **5 or more** → create a `{singular package name}_interfaces` subpackage and
   put the interface files there.
3. **Spanning packages** → if one interface is implemented by concrete types in
   **more than one package**, or the interface exists to break a circular
   dependency, put it in a subpackage under `internal/interfaces/`.

Examples of each:

1. [.../tournament_variant/clusterServiceInterface.go](internal/services/template_generator/providers/topology/tournament_variant/clusterServiceInterface.go)
   — one interface for four implementations, all in that package;
   [internal/services/zones/zoneLabelProviderInterface.go](internal/services/zones/zoneLabelProviderInterface.go)
   — one interface for one implementation. Likewise
   [internal/services/connection_editor/](internal/services/connection_editor/)
   holds three interfaces for its three implementations.
2. [internal/handlers/](internal/handlers/) has six implementation files, so its
   contracts live in
   [internal/handlers/handler_interfaces/](internal/handlers/handler_interfaces/)
   (six files, eight interfaces).
3. [app/gui/interfaces/](app/gui/interfaces/) holds the shared `IDialog` /
   `IPanel` contracts: `drivers` and `panels` both implement them, and declaring
   them in either package would create a `drivers`↔`panels` and
   `drivers`↔`dialogs` cycle.

**Factory return type.** When an interface exists for an implementation, that
implementation's factory function returns the **interface**, not a pointer to the
struct — unless doing so breaks existing functionality. When one implementation
satisfies several interfaces, return the **broadest** one (e.g.
`handlers.NewGuiHandler(...)` returns `handler_interfaces.IGuiHandler`, which
embeds the other handler interfaces).

### 4.3 Method receivers

- The receiver for any method attached to a struct **must be named `this`**.

  ```go
  func (this *ZoneContentManager) Load(path string) error { ... }
  ```

- Apply this consistently across every method of the struct (Go vet warns on
  inconsistent receiver names, so do not mix `this` with anything else).

### 4.4 Package layout

Place new code in the package whose responsibility matches its role:

| Kind of code                                                            | Location                                                                                                                                       |
| ----------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| UI / rendering (Gio widgets, layouts, theming, input)                   | [app/gui/](app/gui/)                                                                                                                           |
| Serializable objects (entities)                                         | [internal/entities/](internal/entities/); read-only `.rmg.json` schema in [internal/entities/template/](internal/entities/template) (see §2.1) |
| Minimal data transfer objects                                           | [internal/dtos/](internal/dtos/)                                                                                                               |
| Data structs with attached logic (might have factory functions as well) | [internal/models/](internal/models/)                                                                                                           |
| Data mappers and converters                                             | [internal/mappers/](internal/mappers/)                                                                                                         |
| Orchestrators / entry points                                            | [internal/handlers/](internal/handlers/)                                                                                                       |
| Business logic, services                                                | [internal/services/](internal/services/)                                                                                                       |
| Game data values (exclusively)                                          | [internal/registry/](internal/registry/) (read-only, see §2.1)                                                                                 |
| Constants, IDs, immutable lookup tables (non game values)               | [internal/common/](internal/common/)                                                                                                           |
| Misc / cross-cutting utility functions                                  | [internal/helpers/](internal/helpers/)                                                                                                         |
| Validators                                                              | [internal/validators/](internal/validators/)                                                                                                   |

- If a struct or function has dependencies (helper structs, private types)
  **that are not used anywhere else**, nest them in a sibling folder next to
  the dependant file rather than polluting a shared package.
- Do not introduce new top-level packages without a clear reason.
- `internal/entities/` is the base data layer and as such should never import other
  packages (except `internal/helpers/data/`); `internal/models/` is the second layer -
  it can import `internal/entities/`, `internal/helpers/`, `internal/registry/`;
  `internal/dtos/` is the highest layer - it is used for transferring data from and
  to `app/.../` - it can import `internal/entities/`, `internal/models/` and any data
  related additional packages (`internal/helpers/`, `internal/registry/`, `internal/common/`, etc.).
  - In theory, if `app/` would be detached from the project, it would have a separate model
    implementation, but here `app/.../` can use the `internal/models/` instead of duplicating code.
  - Object traversal flow for request(from app)-response(from internal handlers):

    ```mermaid
    sequenceDiagram
        participant GUI as app/{user interface implementation}
        participant H as internal/handlers
        participant S as internal/services
        participant R as internal/repositories
        participant Disk as .gen.json / .rmg.json

        Note over GUI: maps stored Model → DTO
        GUI->>H: request DTO
        Note over H: maps DTO → Model
        H->>S: Model
        Note over S: save: maps Model → Entity
        S->>R: Entity
        R->>Disk: Entity
        Disk-->>R: Entity
        R-->>S: Entity
        Note over S: load: maps Entity → Model
        S-->>H: Model
        Note over H: maps Model → DTO
        H-->>GUI: response DTO
        Note over GUI: maps DTO → Model and stores it
    ```

- Packages outside `internal/` must enter internal functionality
  (services, validators, etc.) through `internal/handlers/`; data accessing
  (`internal/registry/`, `internal/common/`), data typing (`internal/models/`, `internal/dtos/`)
  and usage of helpers (`internal/helpers`, `internal/mappers`) is permitted.

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

**Parallel test**: each unit test must contain `t.Parallel()` at the start of
the test and inside every `t.Run` (`paralleltest` linter checks missing t.Parallel() directives).

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
  by or tested in unit tests — unit tests assert real production code, so a
  unit test must never carry the `integration_test` tag (see §4.6.1).
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

**Scope — the ONLY reason to add this tag is `*_testexports.go` consumption:**

- A file gets `//go:build integration_test` **if and only if** it (or another
  file it shares a package with) references an accessor declared in a
  `*_testexports.go` implementation file. That is the whole rule. Today those
  accessors live in
  [app/gui/editor/window_testexports.go](app/gui/editor/window_testexports.go)
  and [app/gui/drivers/state_testexports.go](app/gui/drivers/state_testexports.go).
- **The tag is NOT a label for "this is an integration/performance test."** An
  integration or performance test that only touches production APIs must be
  written **without** any tag so it runs in a plain `go test ./test/...`. Do not
  blanket-apply the tag to every file in [test/integration/](test/integration/)
  or [test/performance/](test/performance/) — apply it file by file.
  (Example: [test/performance/template_generation_test.go](test/performance/template_generation_test.go)
  benchmarks the generator through production APIs only, so it carries no tag.)
- **Never tag a unit test with `integration_test`.** Unit tests assert real
  implementation code and must never see test-only exports; if a unit test
  appears to need one, the test is wrong, not the API.
- **Do NOT** run the whole suite with the tag, and **do NOT** set it as a global
  `go.testTags`/`go.buildTags`. A normal `go test ./...` must stay tag-free; the
  gated files then drop out of the build and are skipped rather than failing.
- Only files under [test/integration/](test/integration/) and
  [test/performance/](test/performance/) may reference the `integration_test`
  accessors. If another package needs an internal, do not widen this tag — add a
  test beside the code (`package X_test` in the same directory) instead.
- Because build constraints apply per file, a package can be **partly** gated:
  the untagged files compile in every run and the tagged ones only under the
  tag. Keep any shared `TestMain`/helpers in the tagged file only if the
  untagged files can run without them.

**Enforcement** — the rules of this section and of §4.6.2 are checked by
[cmd/testlayoutcheck](cmd/testlayoutcheck), a small Go program that walks the
repository and reports every misplaced build tag or misnamed unit-test file. Run
it (VS Code task *"Go: Check test build-tag layout"*) before handing work back:

```powershell
go run ./cmd/testlayoutcheck .
```

It exits `0` and prints `test-layout check passed` when clean, `1` with one line
per violation otherwise. A violation is a broken build — fix the test layout, do
not silence the checker.

### 4.6.2 The `gui` build tag (tests that need a GPU)

Some tests drive a real Gio window or rasterize frames through
`gioui.org/gpu/headless`, which requires a GPU/GL context. The CI pipeline has no
GPU, so those tests must never be picked up by a catch-all run.

**Scope — this tag marks GPU-dependent tests:**

- Add `gui` to any test that opens an `app.Window`, renders through
  `headless.Window` (snapshot tests), or otherwise cannot run without a GPU.
  In practice: everything under [test/integration/gui/](test/integration/gui/)
  and the window-driving benchmarks in [test/performance/](test/performance/)
  (e.g. [test/performance/window_tab_cycling_test.go](test/performance/window_tab_cycling_test.go)).
- The tag exists **to exclude** these tests from catch-all runs such as
  `go test ./test/...` or `go test -tags=integration_test ./test/...`. They are
  opt-in only.
- `gui` is orthogonal to `integration_test`: combine them
  (`//go:build integration_test && gui`) when a GPU test also needs test-only
  exports, and use `gui` alone when it does not.
- **Do NOT** set `gui` as a global/default test tag.

**Running them:**

```powershell
# Default run — everything EXCEPT the gated files (no tag):
go test ./test/... -count=1

# Integration:
go test -tags=integration_test ./test/integration/... -count=1

# UI Integration only (needs a GPU):
go test -tags='integration_test,gui' ./test/integration/gui/... -count=1

# Performance, GPU-free benchmarks:
go test -bench=. -run=xxx ./test/performance/... -benchtime=20x -timeout=120s

# Performance, window-driving benchmarks (needs a GPU):
go test -v -tags='integration_test,gui' -bench=BenchmarkEditorWindow_TabCycling -run=xxx ./test/performance/... -benchtime=20x -timeout=120s
```

In VS Code use the tasks in [.vscode/tasks.json](.vscode/tasks.json). gopls is
configured with `-tags=integration_test` for **analysis only** so the gated files
still get IntelliSense — that does not cause them to run.

### 4.6.3 The `wireinject` build tag (code generation only)

Dependency wiring is generated by [goforj/wire](https://github.com/goforj/wire), a
maintained fork of the archived `google/wire`. The injector **declaration** lives in
`internal/composition/wire.go` behind `//go:build wireinject`; the generated
**implementation** lives beside it in `wire_gen.go`, which carries the inverse
`//go:build !wireinject` constraint and **is committed to the repository**.

**Scope — this tag is for the `wire` generator ONLY:**

- Exactly one file in the repository carries it, and no test ever uses it.
- **Do NOT** pass `-tags=wireinject` to `go build` or `go test`, and **do NOT** add it
  to `go.testTags` / `go.buildTags` / `gopls.build.buildFlags` / `GOFLAGS`. The stub and
  the generated file declare the same function, so compiling both together is a
  duplicate-symbol error. `wire.go` showing as excluded in the editor is expected.
- `wire_gen.go` is generated: never hand-edit it, and never `gofmt`/lint-fix it in
  isolation. Regenerate instead.

**Regenerating** — after changing any provider set, constructor signature, or the
injector's return type, run the *"Go: Generate wire injectors"* task, or:

```powershell
wire gen ./internal/composition/...
```

The `wire` CLI is recorded in the `tool` directive of [tools/go.mod](tools/go.mod)
alongside `golangci-lint` and `gcov2lcov`; install it with
`go install github.com/goforj/wire/cmd/wire@latest`. A missing or ambiguous provider is
a **generation-time** failure — treat a broken `wire gen` as a broken build.

### 4.7 Comments

Don't overuse code comments. Comments describe how a thing is used, and move when the code moves.
To be used mostly to describe functions, not to annotate every line of behavior.

### 4.8 Writing the Plan

Save to `.agent/plans/<descriptive-name>.md` in the repository root. Use this self-documenting template:

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

- **Recommended length: <50 messages per session.**
- Around message **38**, warn the user that the session is approaching the
  recommended limit.
- At message **50** (or sooner if context feels saturated, tools start
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

### 5.4 Local memories

If you need to track any specific information between sessions, write it up in `.agents/memories`
instead of populating some arbitrary temporary user directory. Occasionally clear out stale and/or
irrelevant memories.

---

## 6. Communication Style

- Be brief. 1–3 sentences for simple answers; expand only for genuine
  complexity.
- Use Markdown. Wrap symbols in backticks; link files using
  `[path](path)` or `[path](path#L10-L20)` (workspace-relative, never
  inside backticks).
- Never name internal tools to the user ("I'll run the tests", not "I'll use `runTests`").
- No emojis unless the user asks.
- Don't use sloppy LLM speech, no em dashes, etc., put some soul into your responses.

---

## 7. Quick Reference

| Task                       | Command (Windows PowerShell / Linux bash)              |
| -------------------------- | ------------------------------------------------------ |
| Build                      | `go build ./...`                                       |
| Run GUI                    | `go run .`                                             |
| Run unit tests             | `go test ./test/unit... -count=1`                      |
| Run integration tests      | `go test -tags=integration_test ./test/integration/... -count=1` |
| Run ui integration tests   | `go test -tags=integration_test,gui ./test/integration/gui/... -count=1` |
| Run benchmarks (no GPU)    | `go test -bench=. -run=xxx ./test/performance/... -benchtime=20x -timeout=120s` |
| Run benchmarks (needs GPU) | `go test -tags=integration_test,gui -bench=BenchmarkEditorWindow_TabCycling -run=xxx ./test/performance/... -benchtime=20x -timeout=120s` |
| Run with race detector     | `go test -race ./test/...`                             |
| Check test build-tag layout | `go run ./cmd/testlayoutcheck .` (VS Code task *"Go: Check test build-tag layout"*; see §4.6.1) |
| Unit test coverage report  | `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...` then `go tool cover '-func=coverage.txt'` (see §2.3; VS Code task *"Go: Generate code coverage report"*) |
| Lint (report only)         | `golangci-lint-v2 run ./... --issues-exit-code=0` (VS Code task *"Go: Get Linter Results"*) |
| Lint (auto-fix)            | `golangci-lint-v2 run ./... --issues-exit-code=0 --fix` (VS Code task *"Go: Run Linter"*; clears gci/gofmt/golines formatting findings — re-run to verify) |
| Format Go code             | `gofmt -w .` (never run on `data/`)                    |
| Regenerate DI injectors    | `wire gen ./internal/composition/...` (VS Code task *"Go: Generate wire injectors"*; see §4.6.3) |
| Tidy modules               | `go mod tidy`                                          |

---

**TL;DR:** Don't touch [data/](data/),
[internal/entities/template/](internal/entities/template/) or [internal/registry/](internal/registry/). Never change
where `.rmg.json` is written and never persist the output directory — the game
only reads templates from its own folder. Stay cross-platform.
Cover everything you write with tests. Cap sessions at 38–50 messages and
hand off via `./.agent/session-carry-forward.md`.
