# Session Carry-Forward — Wire Dependency Composition

Generated 2026-08-03. Branch `AD/refactoring-07-21`.

## 1. Session Goal

Evaluate whether a dependency-injection mechanism is worth adopting now that handlers and services take
interfaces, then plan and begin implementing one — compile-time, reflection-free — across the object
graph.

## 2. Fixes Applied

All of these were **pre-existing breakage** discovered while establishing a green baseline; none were
caused by this session's work.

- The entire `test/unit/internal/handlers/guiHandler` package failed to compile
  (`undefined: handlers.GUIHandlerDependencies`). An in-flight refactor had removed
  `GUIHandlerDependencies` / `NewGuiHandlerWithDependencies` and changed `NewGuiHandler` to take five
  collaborators, without updating ~74 call sites.
- Four further constructors had lost their zero-argument variants and had their
  `WithDependencies` / `WithCreationServices` twins renamed onto the plain names, breaking three more
  test packages: [manualReapplyService.go](internal/services/connection_editor/manualReapplyService.go),
  [zoneEditorService.go](internal/services/connection_editor/zoneEditorService.go),
  [mandatoryContentProvider.go](internal/services/template_generator/providers/mandatoryContentProvider.go),
  [templateGenerator.go](internal/services/template_generator/templateGenerator.go).
- [test/unit/architecture/dependency/dependency_test.go](test/unit/architecture/dependency/dependency_test.go)
  matched `/internal/handlers` by **prefix**, so it flagged the nine GUI files that legitimately import
  the new `/internal/handlers/handler_interfaces` sub-package. Switched to an exact package-path match.
- `TestWhenTournamentEnabledWithRandomPortals_AddsPortalConnections` was flaky at roughly 1 run in 5.
  Portals are drawn from each player's own neutral cluster, and `gofakeit.Number(1, 20)` neutral zones
  could leave a cluster with too few portal targets. Narrowed to `gofakeit.Number(8, 20)`; verified
  stable over 30 consecutive runs. The flake was latent — the test had never compiled, so it had never
  actually run.

## 3. Features Added / Changed

No production behaviour changed this session. Two infrastructure additions:

- **`github.com/goforj/wire` adopted** as the project's DI mechanism — a maintained, API-identical fork
  of the archived `google/wire`. Codegen, no runtime container, no reflection, failures surface at
  generation time rather than startup. Rationale: the codebase already *had* a hand-written container
  inside `NewGuiHandler`; the goal is to make it explicit and machine-checked, not to add indirection.
- **`wireinject` build tag documented** in AGENTS.md as a new §4.6.2, scoped as codegen-only and
  explicitly banned from `go build` / `go test` / `go.testTags` / `go.buildTags` / `GOFLAGS`, because
  the stub and the generated file declare the same symbol.

## 4. File Modifications

### Created

| File | Purpose |
| --- | --- |
| [plans/wire-dependency-composition.md](plans/wire-dependency-composition.md) | The 7-phase plan. **Source of truth** — read it before resuming. |
| `test/unit/internal/handlers/guiHandler/common_test.go` | `newProductionGuiHandler()`, `newManualReapplyService()`, `newMandatoryContentProvider()`, `generateDefaultTemplate` |
| `test/unit/internal/services/template_generator/templateGenerator/common_test.go` | `newTemplateGenerator(configuration)` |
| `test/unit/internal/services/connection_editor/zoneEditorService/newDefaultZoneEditorService_test.go` | Split out of `newZoneEditorService_test.go` |

### Edited — configuration (Phase 1)

| File | Change |
| --- | --- |
| [AGENTS.md](AGENTS.md) | New §4.6.2 (`wireinject` tag); new §7 quick-reference row for regeneration |
| [.vscode/tasks.json](.vscode/tasks.json) | New `Go: Generate wire injectors` task (`wire gen ./internal/composition/...`) |
| [.vscode/settings.json](.vscode/settings.json) | cSpell: `goforj`, `injector`, `wireinject`. gopls `buildFlags` deliberately unchanged. |
| [go.mod](go.mod) | `+ github.com/goforj/wire v1.2.0 // indirect` |
| [tools/go.mod](tools/go.mod) | `+ github.com/goforj/wire/cmd/wire` in the `tool` directive; incidental `fsnotify v1.5.4 → v1.7.0`, `+ google/subcommands v1.2.0` |

### Edited — tests (Phase 0)

Call sites retargeted onto the new per-package construction helpers: **74** in `guiHandler`, **16** in
`manualReapplyService`, **23** in `mandatoryContentProvider`, **67** in `templateGenerator`, **4** in the
gated integration suites (`handlers.NewGuiHandler()` → `handlers.NewDefaultGuiHandler()`).

Also edited: `test/unit/architecture/dependency/dependency_test.go`,
`test/test_helpers/templateHandlerMock.go`, `test/integration/editorState_integration_test.go`,
`test/integration/gui/contentRuleDialogs_integration_test.go`,
`test/integration/gui/zoneEditorDialog_integration_test.go`, and the individual `*_test.go` files across
the four affected unit folders.

### Deleted

Five test files, each naming an API that no longer exists; their content was folded into the canonical
`newX_test.go` per AGENTS.md §4.6 (one test file per public function):

- `newGuiHandlerWithDependencies_test.go`
- `newManualReapplyServiceWithDependencies_test.go`
- `newZoneEditorServiceWithCreationServices_test.go`
- `newMandatoryContentProviderWithDependencies_test.go`
- `newTemplateGeneratorWithCreationServices_test.go`

> `plans/clean-architecture-refactoring.md` also shows as deleted in `git status`, but that was the
> author's own staged change, not this session's.

## 5. Tests Added Or Updated

No new test *cases* were written — Phase 0 was a repair, not an expansion. One test case was
**deliberately dropped**: `TestWhenDependencyIsMissing_ReturnsError`, a six-case table asserting runtime
messages like `"template workflow handler is required"`. The nil-validation it exercised no longer
exists, and wire's generation-time failure is a strictly stronger guarantee — but **Phase 2 must confirm
that a missing provider really does fail `wire gen`**, otherwise this is a genuine regression.

Status at session end:

| Command | Result |
| --- | --- |
| `go build ./...` | Clean |
| `go vet ./test/...` | Clean |
| `go test ./test/unit/... -count=1` | **PASS** |
| `go test -tags=integration_test ./test/integration/... -count=1` | PASS (as of end of Phase 0) |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | PASS (as of end of Phase 0) |
| `golangci-lint-v2 run ./test/... --fix` | 0 issues (as of end of Phase 0) |

**Coverage baseline: 64.8% of statements** (`-coverpkg=./internal/...,./app/...` over `./test/unit/...`).
Every later phase is measured against this number (AGENTS.md §2.3).

## 6. Git Status Snapshot

Branch: **`AD/refactoring-07-21`**. Nothing was staged or committed by this session (AGENTS.md §2.5).

The working tree mixes two sets of changes — read the two columns of `git status --short` carefully:

- **Column 1 (staged)** — the author's own pre-existing refactor: the `internal/handlers` split into
  `templateHandler` / `stateHandler` / `zoneEditorHandler`, the move of the GUI port interfaces from
  `app/gui/interfaces/` into `internal/handlers/handler_interfaces/`, the new `internal/repositories/`
  package, the deletion of `plans/clean-architecture-refactoring.md`. **Left untouched — do not
  unstage.**
- **Column 2 (unstaged) + untracked** — this session's work: everything listed in §4 above.

Untracked files the next session inherits:

```
test/unit/internal/handlers/guiHandler/common_test.go
test/unit/internal/services/connection_editor/zoneEditorService/newDefaultZoneEditorService_test.go
test/unit/internal/services/template_generator/templateGenerator/common_test.go
```

## 7. Rejections / Things The User Declined

- **A DI container was rejected** in favour of compile-time codegen. Ruled out: `uber/dig`, `uber/fx`
  (heavy reflection, runtime failures, built for server lifecycles the app does not have) and
  `samber/do` (runtime registry). A plain hand-written composition root was recommended first but the
  author chose wire for long-term maintainability.
- **`google/wire` rejected** by the author in favour of the maintained `goforj/wire` fork.
- **Moving the wire CLI into the root `go.mod` was rejected.** The author clarified the intended
  convention: `tools/go.mod` documents the CLI toolchain other developers should install, while the root
  `go.mod` carries only packages the project imports directly. Both entries coexisting is correct, not a
  contradiction. The CLI is already installed globally on this machine (`wire` is on `PATH`).
- **`plans/clean-architecture-refactoring.md` is not authoritative** — the author confirmed it is an
  uncleaned artifact of a previous refactor. Its recorded non-goal "do not introduce a
  dependency-injection framework" was knowingly overridden. Do not cite it.
- **Phase 1 was halted mid-execution** at the author's request; the session was wrapped up here.

## 8. Open Questions

- None blocking. Two items need the author's eye during review, both flagged in the Phase 1 summary:
  1. `github.com/goforj/wire` sits in the **indirect** block of `go.mod` because nothing imports it yet.
     Phase 2's stub plus `go mod tidy` promotes it to `require`. Expected.
  2. `go mod tidy` in `tools/` incidentally bumped `fsnotify v1.5.4 → v1.7.0` and added
     `google/subcommands v1.2.0`.
- One decision is deferred to Phase 5 by design: each topology service must be classified as stateless
  (share one instance) or stateful (keep transient). The author confirmed both outcomes are acceptable;
  the classification just has to happen **before** the collapse, not after.

## 9. Next Recommended Actions

1. **Close out Phase 1's verification** — it was never run. Execute `golangci-lint-v2 run ./... --issues-exit-code=0`
   and confirm no new findings beyond the 84-item `gochecknoglobals` baseline. The throwaway-injector
   smoke test can be skipped if Phase 2 starts immediately, since Phase 2 produces a real `wire_gen.go`.
2. **Phase 2** — build `internal/composition`:
   export the five handler constructors, write the providers (including the log-and-return-nil wrapper
   for `preview_service.NewPreviewGenerator`), write `wire.go` behind `//go:build wireinject`, generate
   and commit `wire_gen.go`, and add the injector unit test that keeps generated code out of the
   coverage hole. Then confirm the deleted `TestWhenDependencyIsMissing_ReturnsError` guarantee is
   genuinely replaced by a `wire gen` failure.
3. **Phase 3** — single composition root: `editor.NewWindow(handler)`, `app/gui/program.go` as the only
   root, delete `handlers.NewDefaultGuiHandler`, update the `concrete-handlers-only-at-gui-composition-roots`
   depguard rule **and** the expected map in `test/unit/architecture/dependency/dependency_test.go`.
4. **Phase 4** — remove nil-defaulting from `NewCreationServices`, `NewDefaultZoneEditorService`,
   `NewTemplateGenerator`, `NewTopologyProvider`, `NewMandatoryContentProvider`.
5. **Phase 5** — classify then collapse the 13 topology twin constructors; prebuild the enum-keyed
   lookup so `CreateTopologyVariant` stops allocating on the 300 ms auto-regen loop.
6. **Phase 6** — full verification sweep, coverage at or above 64.8%, leave everything unstaged.

## 10. Carry-Forward Prompt

> Read `AGENTS.md` first, in full, before touching anything.
>
> Hard rules, one line each:
> **§2.1** `data/`, `internal/entities/template/` and `internal/registry/` are read-only — read them
> freely, never edit them.
> **§2.2** Everything must build and run on Windows *and* Linux — `path/filepath`, no OS-specific
> syscalls without build tags, and this workspace's shell is PowerShell so chain commands with `;`,
> never `&&`.
> **§2.3** Every non-trivial change ships with tests, and total unit coverage must not drop below the
> recorded 64.8% baseline.
> **§2.4** Multi-session work lives in a plan file, and that file — not the conversation — is the source
> of truth.
> **§2.5** Never stage and never commit; if you find staged changes, leave them exactly as they are.
>
> We are implementing compile-time dependency composition with `github.com/goforj/wire`. The plan is
> `plans/wire-dependency-composition.md` — read it before doing anything else; it carries the locked
> decisions, the per-phase checklists, and the verification commands.
>
> **Where work left off:** Phase 0 (repairing a test suite that did not compile) is Complete and all
> three suites are green at 64.7% coverage. Phase 1 (wire tooling and project configuration) is
> In progress — the dependency, the `tools/go.mod` CLI pin, the `Go: Generate wire injectors` VS Code
> task, the AGENTS.md §4.6.2 documentation and the cSpell entries are all done, but Phase 1's own
> verification plan was never executed. Phases 2 through 6 are untouched.
>
> **Start by** re-running Phase 1's verification, then begin Phase 2: create the `internal/composition`
> package with its providers, the `//go:build wireinject` stub, the generated-and-committed
> `wire_gen.go`, and the injector unit test.
>
> Note that the working tree contains the author's own **staged** in-flight refactor of
> `internal/handlers` alongside this session's unstaged work. Do not unstage or "tidy" it.
>
> Full handoff, including the pre-existing bugs that were repaired and the decisions the author already
> ruled on, is in `./.agent/session-carry-forward.md`.
