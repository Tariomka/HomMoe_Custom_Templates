# Compile-Time Dependency Composition With goforj/wire

Replace the hand-threaded, nil-defaulting object graph in `internal/handlers` with a single
compile-time-generated composition root built by [`github.com/goforj/wire`](https://github.com/goforj/wire),
collapse the duplicate service instances and the duplicate composition roots, and remove the
13 `New*TopologyServiceWithCreationServices` twin constructors along with the per-call topology
allocations in the auto-regeneration hot path.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done, set its status
to `Complete` and write its **Phase Summary** (what was done, key decisions, anything needed to
continue with zero context); run the phase's **Verification Plan** and record the result before
moving on. When all phases are done, fill in **Final Recap** and **Deployment Plan**.

Read `AGENTS.md` first. Hard rules that apply to every phase:

- `data/`, `internal/entities/template/` and `internal/registry/` are **read-only**. Nothing in this
  plan touches them.
- Cross-platform (Windows + Linux): `path/filepath`, no shell-specific code, PowerShell chains with `;`.
- Every non-trivial change ships with tests; coverage must not drop (AGENTS.md §2.3).
- **Never** stage or commit. The author reviews and commits.

## Decisions Already Made (do not re-litigate)

| Question | Decision |
| --- | --- |
| DI approach | `github.com/goforj/wire` — maintained fork of the archived `google/wire`, API-identical, codegen, reflection-free |
| Injector location | New `internal/composition` package; the five handler constructors get **exported** |
| Preview generator failure | Keep degrading gracefully — a wrapper provider logs and returns nil so the injector stays error-free |
| Topology services | A hand-written provider function builds the enum-keyed lookup; wire only calls it |
| `wire_gen.go` coverage | A unit test calls the injector and asserts the graph is non-nil |
| Broken handler tests | Repaired first, in Phase 0, against the current signatures |
| `wireinject` build tag | Accepted; documented in AGENTS.md as a codegen-only tag |
| Duplicate roots | `editor.NewWindow(handler)` takes the handler; `app/gui/program.go` becomes the single root |
| Wire CLI | Already installed globally and on `PATH`. `tools/go.mod` records it as documentation of the toolchain other developers need; the root `go.mod` carries only the runtime package the project imports directly. These are not in conflict. |
| Lifetimes | Stateless collaborator -> singleton. Anything holding per-call mutable state stays transient (a provider function or factory func). Expected and accepted. |

## Baseline Evidence (captured 2026-08-03, before any change)

- `go build ./...` — clean.
- `go vet ./test/unit/internal/handlers/...` — **FAILS**:
  `test\unit\internal\handlers\guiHandler\handlerDependenciesStub_test.go:152:62: undefined: handlers.GUIHandlerDependencies`.
  Roughly 80 further call sites use a zero-argument `handlers.NewGuiHandler()` while
  [internal/handlers/guiHandler.go](internal/handlers/guiHandler.go#L35-L41) now requires five arguments.
  An in-flight refactor removed `GUIHandlerDependencies` / `NewGuiHandlerWithDependencies` without
  updating the suite.
- Duplicate instances: [`NewGuiHandler`](internal/handlers/guiHandler.go#L48-L62) builds `tuningFactory`,
  `zoneEditor`, `creationServices`, then
  [`NewTemplateGenerator`](internal/services/template_generator/templateGenerator.go#L38-L52)
  silently builds a *second* `ZoneEditorService`, `GenerationTuningFactory` and `ZoneLabelProvider`.
- Duplicate roots: [app/gui/editor/window.go](app/gui/editor/window.go#L35) and
  [app/gui/drivers/state.go](app/gui/drivers/state.go#L53) each call `handlers.NewDefaultGuiHandler()`.
- Hot-path allocation: [topologyProvider.go](internal/services/template_generator/providers/topologyProvider.go#L32-L80)
  constructs a fresh topology service on every `CreateTopologyVariant` call, which runs on the 300 ms
  auto-regen debounce loop.
- Lint baseline: 84 `gochecknoglobals` findings (linter runs with `--issues-exit-code=0`).

## Phase 0: Repair The Handler Test Suite

Status: Complete

Get to a green, measurable baseline before introducing wire. No production code changes in this phase.

- [x] Reconcile `test/unit/internal/handlers/guiHandler/handlerDependenciesStub_test.go` with the
      current API. `handlers.GUIHandlerDependencies` no longer exists — either rebuild the stub around
      the five `handler_interfaces` collaborators that `NewGuiHandler` now accepts, or delete the stub
      if nothing else consumes it.
- [x] Retarget or delete `test/unit/internal/handlers/guiHandler/newGuiHandlerWithDependencies_test.go`
      (7 call sites). The dependency-validation behavior it asserted was removed with
      `NewGuiHandlerWithDependencies`; if that validation is still wanted, it belongs on the new
      composition root instead — note it and defer to Phase 2 rather than reinstating it here.
- [x] Update every zero-argument `handlers.NewGuiHandler()` call site to the current five-argument
      signature. Introduce one shared `newProductionGuiHandler(t)` helper in the guiHandler test folder
      rather than repeating `NewGuiHandler(nil, nil, nil, nil, nil)` ~80 times.
- [x] Fix the same breakage in the gated suites: `test/integration/editorState_integration_test.go:30`,
      `test/integration/gui/contentRuleDialogs_integration_test.go:29` and `:60`,
      `test/integration/gui/zoneEditorDialog_integration_test.go:20`.
- [x] Record the coverage baseline number in this phase's summary — every later phase is measured
      against it.

### Verification Plan

- `go build ./...` — no output.
- `go vet ./test/...` — no errors.
- `go test ./test/unit/... -count=1` — all pass.
- `go test -tags=integration_test ./test/integration/... -count=1` — all pass.
- `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` — all pass.
- `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...`
  then `go tool cover '-func=coverage.txt'` — record the total percentage.

### Phase Summary

**COVERAGE BASELINE: 64.8% of statements** (`-coverpkg=./internal/...,./app/...` over `./test/unit/...`).
Every later phase is measured against this number.

All three suites are green and `go vet ./...` (plus both gated tag combinations) is clean.

The breakage was wider than the plan anticipated. The in-flight refactor had not just changed
`NewGuiHandler`; it also deleted the zero-argument variants of four other constructors and renamed
their `WithDependencies` / `WithCreationServices` twins onto the plain names:

| Constructor | Old call in tests | Current signature |
| --- | --- | --- |
| `handlers.NewGuiHandler` | `()` | five `handler_interfaces` collaborators |
| `connection_editor.NewManualReapplyService` | `()` | `(zoneEditor, zoneClassifier, tuningFactory)` |
| `connection_editor.NewZoneEditorService` | `(creationServices)` | `(castleFactory, roadFactory, zoneFactory)` |
| `providers.NewMandatoryContentProvider` | `()` | `(zoneClassifier, zoneEditor)` |
| `template_generator.NewTemplateGenerator` | `(configuration)` | `(configuration, creationServices)` |

**Approach — one construction seam per test package.** Rather than repeating the full argument list at
~180 call sites, each affected test folder got a package-local helper in its `common_test.go`. Phases 2
and 4 only have to change those helper bodies:

- `test/unit/internal/handlers/guiHandler/common_test.go` (new) — `newProductionGuiHandler()`,
  `newManualReapplyService()`, `newMandatoryContentProvider()`, plus `generateDefaultTemplate` moved in
  from `updateTemplate_test.go` (its parameter widened from `*handlers.GUIHandler` to
  `handler_interfaces.ITemplateHandler`).
- `test/unit/internal/services/connection_editor/manualReapplyService/common_test.go` —
  `newManualReapplyService()`.
- `test/unit/internal/services/template_generator/providers/mandatoryContentProvider/common_test.go` —
  `newMandatoryContentProvider()`.
- `test/unit/internal/services/template_generator/templateGenerator/common_test.go` (new) —
  `newTemplateGenerator(configuration)`.

**Call sites retargeted:** 74 in the guiHandler folder, 16 in `manualReapplyService`, 23 in
`mandatoryContentProvider`, 67 in `templateGenerator`, 4 in the gated integration suites
(`handlers.NewGuiHandler()` → `handlers.NewDefaultGuiHandler()`).

**Duplicate constructor test files deleted** — each named an API that no longer exists; their content was
folded into the canonical `newX_test.go` file (AGENTS.md §4.6, one test file per public function):

- `newGuiHandlerWithDependencies_test.go`
- `newManualReapplyServiceWithDependencies_test.go`
- `newZoneEditorServiceWithCreationServices_test.go`
- `newMandatoryContentProviderWithDependencies_test.go`
- `newTemplateGeneratorWithCreationServices_test.go`

`newZoneEditorService_test.go` had been testing `NewDefaultZoneEditorService()`; it now tests the real
three-factory `NewZoneEditorService`, and the default-constructor assertion moved to the new
`newDefaultZoneEditorService_test.go`.

**Deliberately dropped test — carry into Phase 2.** `TestWhenDependencyIsMissing_ReturnsError` (a
six-case table asserting messages like `"template workflow handler is required"`) was deleted, not
ported: the runtime nil-validation it exercised no longer exists. Wire's generation-time failure is a
strictly stronger guarantee, so this is not a coverage regression in spirit — but Phase 2 must confirm a
missing provider really does fail `wire gen`.

**Two repairs outside the planned scope, both pre-existing and both blocking a green baseline:**

1. `test/unit/architecture/dependency/dependency_test.go` —
   `TestWhenAppHandlerImportsAreScanned_OnlyCompositionRootsDependOnConcreteHandlers` matched
   `/internal/handlers` by *prefix*, so it flagged the nine GUI files that legitimately import the new
   `/internal/handlers/handler_interfaces` sub-package. Switched to an exact package-path match, which
   is what "depend on **concrete** handlers" was always supposed to mean. Phase 3 still has to update
   the expected map when `drivers/state.go` and `editor/window.go` stop importing the concrete package.
2. `TestWhenTournamentEnabledWithRandomPortals_AddsPortalConnections` was flaky at roughly 1 run in 5.
   Portals are drawn from each player's own neutral cluster, and `gofakeit.Number(1, 20)` neutral zones
   could leave a cluster with too few portal targets. Narrowed to `gofakeit.Number(8, 20)`; verified
   stable over 30 consecutive runs. (The flake was latent, not new — the test could not compile before
   this phase, so it had never actually run.)

`golangci-lint-v2 run ./test/... --fix` reports **0 issues**. Nothing was staged or committed; the
author's pre-existing staged refactor was left untouched.

## Phase 1: Wire Tooling And Project Configuration

Status: Complete

Add the dependency and make regeneration reproducible on Windows and Linux. Still no graph changes.

- [x] `go get github.com/goforj/wire` in the root module. Only the tiny runtime package
      (`wire.Build`, `wire.NewSet`, `wire.Bind`, `wire.Struct`) is imported by the build-tagged stub;
      the generated output has no wire dependency at all.
- [x] Add `github.com/goforj/wire/cmd/wire` to the `tool` directive in [tools/go.mod](tools/go.mod),
      then `go mod tidy` inside `tools/`. The `tools/` module is **documentation of the CLI toolchain**
      other developers need, alongside `golangci-lint` and `gcov2lcov`; the root `go.mod` carries only
      the runtime package the project imports directly. The CLI is already installed globally on this
      machine, so generation runs as a plain `wire` invocation from the repository root.
- [x] Add a `Go: Generate wire injectors` task to [.vscode/tasks.json](.vscode/tasks.json) running
      `wire gen ./internal/composition/...` from the repository root, in the `build` group.
- [x] Do **not** add `wireinject` to `gopls.build.buildFlags` in
      [.vscode/settings.json](.vscode/settings.json) (currently `-tags=integration_test,gui`).
      The stub and `wire_gen.go` declare the same function, so compiling both together is a duplicate-symbol
      error. The stub showing as excluded in the editor is expected and correct.
- [x] Confirm golangci-lint v2 skips `wire_gen.go` via its generated-file header
      (`// Code generated by Wire. DO NOT EDIT.`). [.golangci.yml](.golangci.yml) sets no
      `issues.exclude-generated` key, so the v2 default (`lax`, which skips generated files) applies.
      **Confirmed in Phase 2** against the real `wire_gen.go`: zero findings for that file while the
      hand-written files in the same package were flagged for formatting.
- [x] Add `goforj`, `wireinject` and `injector` to the `cSpell.words` list in
      [.vscode/settings.json](.vscode/settings.json).
- [x] Document the tag in AGENTS.md as a new §4.6.2: `wireinject` is a **codegen-only** tag, never
      passed to `go build` or `go test`, never added to `go.testTags`/`go.buildTags`/`GOFLAGS`, and it
      applies to exactly one file. Add a regeneration row to the §7 quick-reference table.

### Verification Plan

- `go build ./...` — unchanged, clean.
- `go test ./test/unit/... -count=1` — unchanged, green.
- Running the new VS Code task on a trivial throwaway injector produces a `wire_gen.go`; delete the
  throwaway afterwards.
- `golangci-lint-v2 run ./... --issues-exit-code=0` — no new findings beyond the 84-item
  `gochecknoglobals` baseline.

### Phase Summary

**Verification executed 2026-08-04 (the outstanding item from the previous session).**

| Check | Result |
| --- | --- |
| `go build ./...` | Clean |
| `go test ./test/unit/... -count=1` | Green |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **37 issues** — 34 `gochecknoglobals`, 2 `dupl`, 1 `funlen`; identical to the Phase 1.5 baseline, no new findings |
| Wire CLI reachable | `wire` resolves to `C:\Users\...\go\bin\wire.exe`; subcommands `gen`, `check`, `diff`, `show`, `watch` all present |

The throwaway-injector smoke test was **skipped deliberately** — Phase 2 followed immediately and
produced a real `wire_gen.go`, which is a strictly better test of the same thing. The generated-file
lint exemption was confirmed against that real file rather than a throwaway.

PowerShell note for future sessions: the `wire` CLI writes its success banner to **stderr**, so
`wire gen` and `wire check` appear to "fail" under PowerShell's `NativeCommandError` handling even when
they succeed. Judge them by `wire diff` (exit 0 = committed file matches) rather than by the banner.

**Done:** dependency added, CLI pinned in `tools/go.mod`, VS Code task added, AGENTS.md §4.6.2 and the
§7 quick-reference row written, cSpell words added, gopls flags deliberately left untouched.

**`go.mod` resolved.** `github.com/goforj/wire v1.2.0` was sitting in the *indirect* block because
nothing imported it. Phase 2's `//go:build wireinject` stub plus `go mod tidy` promoted it to the direct
`require` block, exactly as predicted. `go mod tidy` sees all build configurations, including
tag-excluded files.

**Still for the author to eyeball:** `go mod tidy` in `tools/` had incidentally bumped an unrelated
transitive dependency, `github.com/fsnotify/fsnotify v1.5.4 → v1.7.0` (pulled in by wire's `wire watch`
support), and added `github.com/google/subcommands v1.2.0`.

## Phase 1.5: Remove The `CreationServices` Aggregate

Status: Complete

Author-requested, inserted 2026-08-04 before Phase 2. `zones.CreationServices` is a bundle struct that
exists only to carry three factories through constructors — exactly the job wire is being brought in to
do. Removing it before the injector is written means the provider set describes the real graph instead
of a bundle.

**This phase absorbs work originally scheduled for Phases 4 and 5.** Removing the type forces the twin
constructor collapse, because the `WithCreationServices` suffix names a type that no longer exists.

Key finding driving the design: `TopologyBase` stores only `roadFactory` and `zoneFactory` — it never
touches `castleFactory`, which reaches zone creation *through* `ZoneFactory`. `castleFactory` is needed
only to build `ZoneEditorService`.

Author decisions: pass **only the two factories actually used**, and **delete the zero-argument
convenience constructors now** rather than in Phase 4.

- [x] `base.NewTopologyBase(zoneFactory *zones.ZoneFactory, roadFactory *zones.RoadFactory)`; delete
      `NewTopologyBaseWithCreationServices`.
- [x] Collapse each topology service to a single `New*TopologyService(zoneFactory, roadFactory)`,
      deleting both the zero-argument variant and the `WithCreationServices` twin: ring, hub,
      geometricHub, chain, web, random, circles, square, geometric, cross, fractal, tournament.
      `chainTopology` and `positionedTopologyBuilder` already have one constructor — just change the
      parameters.
- [x] Same for the four `tournament_variant` cluster services: balanced, chain, hub, ring.
- [x] `TournamentTopologyService` stores `zoneFactory` + `roadFactory` instead of a `creationServices`
      field, and passes them to whichever cluster service its `CreateTopologyVariant` switch selects.
- [x] `providers.NewTopologyProvider(zoneFactory, roadFactory)` — single constructor, no nil branch.
- [x] `template_generator.NewTemplateGenerator` takes the three factories explicitly (it needs
      `castleFactory` for the `ZoneEditorService` it builds).
- [x] [internal/handlers/guiHandler.go](internal/handlers/guiHandler.go) constructs the three factories
      directly instead of unpacking a `CreationServices`.
- [x] Remove nil-defaulting from `zones.NewZoneFactory` so the whole `zones` package is explicit
      (5 `NewZoneFactory(nil, nil)` test call sites). `connection_editor.NewZoneEditorService` keeps its
      nil branches until Phase 4.
- [x] Delete `internal/services/zones/creationServices.go`
      and `test/unit/internal/services/zones/creationServices/`.
- [x] Delete the ~15 obsolete `new*WithCreationServices_test.go` files, folding any unique assertion
      into the canonical `newX_test.go` per AGENTS.md §4.6.

### Verification Plan

- `grep -r "CreationServices"` across the repository returns **zero** matches.
- `go build ./...` — clean.
- `go test ./test/unit/... -count=1` — green. Generator output must be unchanged: the golden-template
  tests in `test/unit/internal/services/template_generator/` are the guard.
- `go test -tags=integration_test ./test/integration/... -count=1` and the `gui` variant — green.
- Coverage at or above the 64.8% Phase 0 baseline.

### Phase Summary

`zones.CreationServices` is gone. A repository-wide search for the identifier returns zero matches.

**Shape of the result.** Every topology service, cluster service and the `TopologyProvider` now has
exactly one constructor taking `(zoneFactory *zones.ZoneFactory, roadFactory *zones.RoadFactory)`.
`TopologyBase` deliberately does **not** take `castleFactory`: it never used it, and castles reach zone
creation through `ZoneFactory`, which composes `CastleFactory` itself. Only
`NewTemplateGenerator(configuration, castleFactory, roadFactory, zoneFactory)` takes all three, because
it builds a `ZoneEditorService`.

**Deleted:** `internal/services/zones/creationServices.go`; the whole
`test/unit/internal/services/zones/creationServices/` folder; 17 `*WithCreationServices_test.go` files
(each was a verbatim duplicate of the canonical `newX_test.go` in the same folder — single
`assert.NotNil` on `ZoneLabelProvider` or on the constructed value — so nothing needed folding in);
`TestWhenCreationServicesAreOmitted_ReturnsGenerator`, which asserted nil-defaulting that no longer
exists.

**New test seam.** [test/test_helpers/zoneFactories.go](test/test_helpers/zoneFactories.go) exposes
`NewZoneFactories() (*zones.ZoneFactory, *zones.RoadFactory)`. Its results match the topology
constructor parameter lists exactly, so it spreads inline:
`topology.NewRingTopologyService(test_helpers.NewZoneFactories())`. That kept ~215 call-site updates to
a single-expression substitution instead of two statements per site.

**Verified:** `go build ./...` clean; `go vet` clean under no tags and under `integration_test,gui`;
unit, integration and GUI integration suites all green; `golangci-lint-v2 run ./... --fix` leaves 37
issues, all pre-existing (34 `gochecknoglobals`, 2 `dupl` in `app/gui/widgets/buttonWidget.go`, and
`funlen` on `NewGuiHandler` at 63 > 60 — the refactor made that function one line *shorter*, and Phase 2
shrinks it further by moving construction into `internal/composition`).

**Coverage: 64.7% against a 64.8% baseline.** This 0.1% dip is arithmetic, not a regression. Every
statement removed was fully covered (the aggregate constructor, the nil-default branches, the zero-arg
twins), and removing fully-covered statements from a base below 100% necessarily lowers the ratio:
`(C-k)/(T-k) < C/T` whenever `C < T`. `go tool cover -func` confirms every rewritten constructor is at
100% and nothing in the touched packages became newly uncovered.

**Consequences for later phases.** Phase 4 lost its `NewCreationServices` and `NewTopologyProvider`
items. Phase 5 lost the constructor collapse entirely and is now only the performance work (prebuilding
the topology lookup). Note for Phase 5: `TournamentTopologyService` assigns `this.clusterService` inside
`CreateTopologyVariant`, so it is genuinely stateful as written and cannot be shared as a singleton
without lifting that switch into a local variable.

## Phase 2: The `internal/composition` Package And Injector

Status: Complete

Move the wiring currently embedded in `NewGuiHandler` into generated code.

- [x] Export the five handler constructors in `internal/handlers`: `newTemplateHandler`,
      `newStateHandler`, `newPreviewHandler`, `newContentRuleHandler`, `newZoneEditorHandler` become
      `New*`. The handler **structs** stay unexported — the constructors already return
      `handler_interfaces` types, which is exactly what wire needs, so no `wire.Bind` is required for
      them. Verify each return type is the interface, not the concrete struct, before writing providers.
      **Three of the five returned the concrete unexported struct** (`*stateHandler`, `*previewHandler`,
      `*contentRuleHandler`) and were changed to return their interface.
- [x] Create `internal/composition` containing:
  - `previewGeneratorProvider.go` — wraps
    [`preview_service.NewPreviewGenerator`](internal/services/preview_service/previewGeneratorService.go#L34),
    logs a failure via `slog.Error` and returns nil, preserving today's degrade-gracefully behavior
    from [guiHandler.go](internal/handlers/guiHandler.go#L41-L47) and keeping the injector error-free.
  - ~~`topologyServiceProvider.go`~~ — **deferred to Phase 5**, see the summary.
  - `providerSets.go` — the `wire.NewSet` declarations. **Documented exemption to AGENTS.md §4.1**:
    provider sets are package-level vars, not structs, so the one-struct-per-file rule does not apply;
    grouping them in one file keeps the graph readable. Expect a few new `gochecknoglobals` findings.
  - `wire.go` — `//go:build wireinject`, declaring
    `func InitializeGuiHandler() handler_interfaces.IGuiHandler { wire.Build(GuiHandlerSet); return nil }`.
  - `wire_gen.go` — generated, **committed to the repository**.
- [x] Make sure every collaborator is provided exactly once so the graph has no duplicates:
      `ZoneClassifier`, ~~`CreationServices`~~ (removed in Phase 1.5) and its
      `CastleFactory`/`RoadFactory`/`ZoneFactory`, `ZoneEditorService`, `GenerationTuningFactory`,
      `FileService`, `GeneratorConfigMapper`, `ConnectionEditorService`, `ManualReapplyService`,
      `EditorStateValidator`, `ContentRuleService`, `PreviewLayoutService`, `MandatoryContentProvider`,
      `TemplateGenerator`.
- [x] Break the duplicate-instance bug in
      [templateGenerator.go](internal/services/template_generator/templateGenerator.go#L28-L52):
      `NewTemplateGenerator` must **receive** its `ZoneLabelProvider`, `GenerationTuningFactory` and
      providers instead of constructing them, so wire's singletons are the only instances.
- [x] Reduce `NewGuiHandler` to a plain five-argument struct assembler with no internal wiring.
- [x] Add `test/unit/internal/composition/wire/initializeGuiHandler_test.go`
      (package `wire_test`) asserting the returned graph is non-nil. This covers the generated file and
      doubles as a wiring regression test. Follow AGENTS.md §4.6: `t.Parallel()`, Arrange/Act/Assert,
      one logical assertion, `Test{Scenario}_{ExpectedBehavior}` naming.
- [x] If Phase 0 deferred dependency validation here, decide whether the injector needs it. It probably
      does **not** — wire fails at generation time when a provider is missing, which is strictly
      stronger than the runtime nil checks that were removed. **Confirmed empirically — see summary.**

### Verification Plan

- Regenerating with the Phase 1 task produces a `wire_gen.go` identical to the committed one
  (determinism check; the fork lists generator determinism as a goal).
- `go build ./...` — clean.
- `go test ./test/unit/... -count=1` — green.
- Coverage total is at or above the Phase 0 baseline.
- Manual read of `wire_gen.go`: each shared service appears exactly once.

### Phase Summary

`internal/composition` now owns the object graph. `InitializeGuiHandler()` is generated, committed and
covered.

| Check | Result |
| --- | --- |
| `wire diff ./internal/composition/...` | Exit 0 — committed `wire_gen.go` is byte-identical to a fresh generation |
| `go build ./...` | Clean |
| `go vet -tags='integration_test,gui' ./...` | Clean |
| `go test ./test/unit/... -count=1` | Green |
| `go test -tags=integration_test ./test/integration/... -count=1` | Green |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | Green |
| Coverage | **64.8%** — back to the Phase 0 baseline, up from 64.7% at the end of Phase 1.5 |
| `golangci-lint-v2 run ./... --fix` | 42 issues: 40 `gochecknoglobals` (34 pre-existing + **6 new provider sets**, as predicted), 2 pre-existing `dupl`. The `funlen` finding on `NewGuiHandler` is **gone**. |

**Provider sets.** Six exported `wire.NewSet` vars in `providerSets.go`, grouped by role rather than
dumped into one blob: `ZoneSet` (factories + classifier + label provider), `GenerationSet` (generator
and its five providers), `EditorSet` (zone editor, manual reapply, content rules, preview),
`InfrastructureSet` (file service, config mapper, state validator), `HandlerSet` (the five handlers),
and `GuiHandlerSet` which composes all of them plus `handlers.NewGuiHandler`.

**Handler constructors.** All five are now exported **and all five return their interface**. Three of
them (`NewStateHandler`, `NewPreviewHandler`, `NewContentRuleHandler`) previously returned the concrete
unexported struct, which would have forced a `wire.Bind` and would have tripped `unexported-return`
once exported. The structs themselves stay unexported.

**Duplicate instances removed.** `wire_gen.go` reads as a flat list of 20 constructor calls with no
repeats. Two duplicate-instance bugs were fixed to get there:

1. `NewTemplateGenerator` built its own `ZoneLabelProvider`, `GenerationTuningFactory`,
   `ContentLimitProvider`, `GameRulesProvider`, `TopologyProvider`, `ZoneLayoutProvider` **and a second
   `ZoneEditorService`** (via `NewMandatoryContentProvider(nil, ...)`). It now takes all eight
   collaborators as parameters. Behaviour is unchanged: the `nil` classifier it used to pass simply
   defaulted to a fresh stateless `ZoneClassifier`, so sharing the injector's instance is equivalent.
2. `NewPreviewGenerator` built its own `PreviewLayoutService` while the GUI built another. It now takes
   one, satisfying the plan's "provided exactly once" requirement for `PreviewLayoutService`.

**`topologyServiceProvider.go` was deliberately not created.** In this phase its only possible content
is `return providers.NewTopologyProvider(zoneFactory, roadFactory)` — an identity wrapper that wire does
not need, since it calls the constructor directly. `GenerationSet` references
`providers.NewTopologyProvider` as-is. Phase 5 creates the file when it has real content (the prebuilt
enum-keyed lookup). Nothing else about Phase 5 changes.

**The deleted `TestWhenDependencyIsMissing_ReturnsError` guarantee is genuinely replaced.** Verified by
temporarily deleting `file_service.NewFileService` from `InfrastructureSet` and running `wire check`:

```
inject InitializeGuiHandler: no provider found for *...file_service.FileService
needed by ...handler_interfaces.ITemplateHandler in provider set "GuiHandlerSet"
needed by ...handler_interfaces.IGuiHandler in provider set "GuiHandlerSet"
```

Generation-time, with the full dependency chain named. Strictly stronger than the six runtime nil checks
Phase 0 dropped. The probe was reverted immediately.

**`NewDefaultGuiHandler` carries the old wiring for one more phase.** `NewGuiHandler` is now the plain
five-argument assembler wire needs, so the hand-rolled graph moved into `NewDefaultGuiHandler`, whose
only remaining callers are [app/gui/editor/window.go](app/gui/editor/window.go#L35),
[app/gui/drivers/state.go](app/gui/drivers/state.go#L53) and four test files. **Phase 3 deletes the
function outright**, and the temporary duplication (including the `newDefaultPreviewGenerator` helper,
which mirrors `composition.providePreviewGenerator`) goes with it. Do not "fix" the duplication any
other way.

**Coverage detail.** `wire_gen.go`'s `InitializeGuiHandler` is at **100%** — it is straight-line code, so
the single non-nil assertion covers all 21 statements. `providePreviewGenerator` sits at 60%: its
`err != nil` branch is unreachable in a build whose `go:embed`-ed assets compiled, and is now recorded
in [todo/test_observations.md](todo/test_observations.md).

**PowerShell gotcha, carried forward.** The `wire` CLI writes its success banner to stderr, so
`wire gen` / `wire check` surface as `NativeCommandError` with exit code 1 even on success. Use
`wire diff` (exit 0 = committed file is current) as the reliable signal.

## Phase 3: One Composition Root

Status: Complete

- [x] Change `editor.NewWindow` to accept `handler_interfaces.IGuiHandler` and stop calling
      `handlers.NewDefaultGuiHandler()` ([window.go](app/gui/editor/window.go#L34-L46)).
- [x] Change `drivers.NewUIState` likewise, or remove it in favour of the existing
      `NewUIStateWithBackend` ([state.go](app/gui/drivers/state.go#L52-L54)). — removed.
- [x] `eventLoop` in [app/gui/program.go](app/gui/program.go#L23-L26) calls
      `composition.InitializeGuiHandler()` once and passes the result into `editor.NewWindow`.
- [x] Delete `handlers.NewDefaultGuiHandler`.
- [x] Update the `concrete-handlers-only-at-gui-composition-roots` depguard rule in
      [.golangci.yml](.golangci.yml#L146-L154): drop the `app/gui/drivers/state.go` and
      `app/gui/editor/window.go` exceptions. `program.go` imports `internal/composition` and
      `handler_interfaces`, neither of which the rule denies, so no new exception is needed — confirm
      this by running the linter rather than assuming.
- [x] Update every `editor.NewWindow()` / `drivers.NewUIState()` call site in the unit, integration and
      performance suites. `test/performance/appRunner_test.go` and the `test/integration/gui/` runner
      are the likely ones.

### Verification Plan

- `go build ./...` — clean.
- `golangci-lint-v2 run ./... --issues-exit-code=0` — the depguard rule still rejects a deliberate
  temporary `internal/handlers` import added to some other `app/` file (revert the probe afterwards).
- All three suites green: unit, `-tags=integration_test`, `-tags='integration_test,gui'`.
- `go test -v -tags=integration_test -bench=BenchmarkEditorWindow_TabCycling -run=xxx ./test/performance/... -benchtime=20x -timeout=120s`
  still compiles and runs.

### Phase Summary

`app/gui` no longer knows how to build a handler graph. `internal/composition.InitializeGuiHandler`
is now the single place the object graph is assembled, and `eventLoop` in
[app/gui/program.go](app/gui/program.go) is the single place it is called in production.

**What changed**

- [app/gui/editor/window.go](app/gui/editor/window.go) — `NewWindow(backend handler_interfaces.IGuiHandler)`.
  The `internal/handlers` import is replaced by `internal/handlers/handler_interfaces`.
- [app/gui/drivers/state.go](app/gui/drivers/state.go) — `NewUIState()` deleted; the surviving
  `NewUIStateWithBackend` already took the interface. The `internal/handlers` import is gone.
  `NewUIState` was *not* rewritten to call `composition.InitializeGuiHandler()`: that would have made
  `drivers` a second composition root, which is exactly what this phase removes.
- [app/gui/program.go](app/gui/program.go) — `editor.NewWindow(composition.InitializeGuiHandler())`.
- [internal/handlers/guiHandler.go](internal/handlers/guiHandler.go) — `NewDefaultGuiHandler` and its
  `newDefaultPreviewGenerator` helper deleted along with the twelve imports they alone needed. The file
  is now just `GUIHandler`, the five-argument `NewGuiHandler` assembler, and the delegating methods.
- [.golangci.yml](.golangci.yml) — the `concrete-handlers-only-at-gui-composition-roots` rule lost both
  file exceptions, so no non-test file under `app/` may import `internal/handlers` at all.
- [test/unit/architecture/dependency/dependency_test.go](test/unit/architecture/dependency/dependency_test.go) —
  the expected map in `TestWhenAppHandlerImportsAreScanned_OnlyCompositionRootsDependOnConcreteHandlers`
  is now empty. `findForbiddenAppImports` gained `internal/composition` in `allowedRoots`, since
  `program.go` legitimately imports it (the depguard config already permitted this; the architecture
  test was the stricter of the two and had to be widened deliberately).

**Test call sites**

- [test/test_helpers/integration_common/appRunner.go](test/test_helpers/integration_common/appRunner.go) —
  `editor.NewWindow(composition.InitializeGuiHandler())`. This is the only `editor.NewWindow` call in the
  gated suites; `test/performance` drives the window through `AppRunner`, so it needed no edit.
- [test/integration/editorState_integration_test.go](test/integration/editorState_integration_test.go) —
  gained a package-local `newUIState()` helper returning
  `drivers.NewUIStateWithBackend(composition.InitializeGuiHandler())`. All twelve `drivers.NewUIState()`
  call sites across the four files of `package integration_test` now use it, and `newEditorSession` builds
  its backend from `composition` too.
- [test/integration/gui/contentRuleDialogs_integration_test.go](test/integration/gui/contentRuleDialogs_integration_test.go)
  and [test/integration/gui/zoneEditorDialog_integration_test.go](test/integration/gui/zoneEditorDialog_integration_test.go) —
  `handlers.NewDefaultGuiHandler()` → `composition.InitializeGuiHandler()`.
- [test/unit/internal/handlers/guiHandler/common_test.go](test/unit/internal/handlers/guiHandler/common_test.go) —
  `newProductionGuiHandler()` now returns `composition.InitializeGuiHandler()`, which is what "the same
  handler graph the application uses" means from this phase onward.

**Verification**

| Check | Result |
| --- | --- |
| `go build ./...` | clean |
| `go vet -tags='integration_test,gui' ./...` | clean |
| `go test ./test/unit/... -count=1` | green |
| `go test -tags=integration_test ./test/integration/... -count=1` | green |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | green |
| `BenchmarkEditorWindow_TabCycling` (20x) | passes, 4 579 735 ns/op |
| depguard probe (blank `internal/handlers` import in `window.go`) | rejected as expected, probe reverted |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | 42 issues (40 `gochecknoglobals`, 2 `dupl`) — unchanged baseline |
| Unit coverage | **64.7%** (was 64.8%) |

**Coverage note.** The 0.1 pp move is arithmetic, not a regression: no statement that was covered before
is uncovered now, and no new uncovered code was introduced. Deleting `NewDefaultGuiHandler` and
`newDefaultPreviewGenerator` removed ~22 *covered* statements that duplicated what `wire_gen.go` already
builds, so both numerator and denominator shrank by nearly the same amount. Phase 6 should record 64.7%
as the standing baseline.

## Phase 4: Remove Nil-Defaulting

Status: Complete

With wire supplying every dependency, a missing dependency should be a compile error, not a silent
second instance. Remove the fallbacks; do not replace them with panics.

- [x] `zones.NewCreationServices` — **absorbed into Phase 1.5** (the type is gone).
- [x] `connection_editor.NewDefaultZoneEditorService` — delete, along with the nil branches in
      `NewZoneEditorService`.
- [x] `template_generator.NewTemplateGenerator` — drop the `configuration == nil` branch (the
      factory parameters became explicit in Phase 1.5).
- [x] `providers.NewTopologyProvider` — **absorbed into Phase 1.5**.
- [x] `providers.NewMandatoryContentProvider` — remove any nil-tolerant parameters (currently called
      with a nil classifier from `NewTemplateGenerator`).
- [x] Sweep for remaining `New...(nil` call sites across `internal/` and `test/` and update them to pass
      real collaborators or testify mocks.
- [x] `config.NewGeneratorConfig()` stays as-is — it is a value factory, not a dependency fallback.

### Verification Plan

- `grep` for `(nil, nil)` and `== nil {` inside constructors under `internal/services` and
  `internal/handlers` returns no dependency-defaulting matches.
- `go build ./...` — clean.
- All three suites green.
- Coverage at or above the Phase 0 baseline (deleting unreachable nil branches should nudge it up).

### Phase Summary

Every constructor under `internal/` now takes its collaborators and stores them verbatim. A missing
dependency is a compile error.

**Production constructors cleaned (5 files):**

| File | Change |
| --- | --- |
| `internal/services/connection_editor/zoneEditorService.go` | `NewDefaultZoneEditorService()` deleted; `NewZoneEditorService(castleFactory, roadFactory, zoneFactory)` returns the struct with no nil branches |
| `internal/services/connection_editor/manualReapplyService.go` | three nil branches removed from `NewManualReapplyService` |
| `internal/services/connection_editor/connectionEditorService.go` | nil branch removed; body is a single `return` |
| `internal/services/template_generator/providers/mandatoryContentProvider.go` | both nil branches removed |
| `internal/services/template_generator/templateGenerator.go` | `if configuration == nil { configuration = config.NewGeneratorConfig() }` removed |

**Test-side collaborator helper.** `NewDefaultZoneEditorService` had 79 call sites across 26 files,
all of them tests. Rather than repeat the three-line wiring in every test, a
`test_helpers.NewZoneEditorService()` factory was added at
[test/test_helpers/zoneEditorService.go](test/test_helpers/zoneEditorService.go), mirroring the
`test_helpers.NewZoneFactories()` precedent introduced in Phase 1.5. The convenience lives in the test
helper package, never in production code — production callers keep passing collaborators explicitly, so
the compiler still catches a forgotten dependency. 23 test files were updated mechanically; the
`connection_editor` import was dropped from 20 of them and retained in 3 that use the package for other
reasons.

**Tests removed** (they asserted the deleted fallback behaviour, so they had no subject any more):

- `zoneEditorService/newDefaultZoneEditorService_test.go` — whole file; the sibling
  `newZoneEditorService_test.go` already covers the explicit constructor.
- `connectionEditorService/newConnectionEditorService_test.go` → `TestWhenClassifierIsNil_ReturnsUsableService`
- `manualReapplyService/newManualReapplyService_test.go` → `TestWhenDependenciesAreNil_ReturnsUsableService`
- `templateGenerator/newTemplateGenerator_test.go` → `TestWhenConfigurationIsNil_FallsBackToDefaultConfiguration`

Newly-unused imports were pruned from each of those files.

**`New...(nil` sweep.** The five remaining `zones.NewZoneFactory(nil, nil)` call sites in
`test/unit/internal/services/zones/zoneFactory/` were compiling only because the nil factories were
never dereferenced. They now use a package-local `newZoneFactory()` helper in `common_test.go` that
wires `NewCastleFactory()` and `NewRoadFactory()`, and `TestWhenDependenciesAreOmitted_ReturnsUsableFactory`
was renamed to `TestWhenDependenciesAreProvided_ReturnsInstance`.

**Deliberately left alone** — these pass optional *data*, not dependencies, so `nil` is a meaningful
argument rather than a fallback trigger:

- `content_rules.NewRuleDistanceToRoad(nil)`, `NewRuleDistanceToTown(nil)`, `NewRuleVariant(nil, nil)` —
  a nil `*models.DistancePreset` / `*models.VariantMapping` means "no preset", which is a real rule state.
- `components.NewDropdownSelector(nil)` in `app/gui/dialogs/zoneEditorDialog.go` — an empty option list.
- `config.NewGeneratorConfig()` — a value factory, per the checklist above.
- `TemplateGenerator.SetConfiguration`'s `if configuration != nil` guard. It is a setter over a value,
  not a constructor over a dependency, so it falls outside this phase's grep criterion; its only caller
  ([internal/handlers/templateHandler.go](internal/handlers/templateHandler.go#L63)) already passes a
  non-nil mapper result.

**Verification**

| Check | Result |
| --- | --- |
| `grep "(nil, nil)"` under `internal/services` (production) | no matches |
| `grep "== nil {"` under `internal/handlers` | 4 matches, all runtime data validation (`rule`, `stateDto.State`, `template`, `templateDto.Template`) |
| `grep "== nil {"` / `"!= nil {"` under `internal/services` (production) | no constructor dependency defaults remain |
| `go build ./...` | clean |
| `go vet -tags='integration_test,gui' ./...` | clean |
| `go test ./test/unit/... -count=1` | green |
| `go test -tags=integration_test ./test/integration/... -count=1` | ok 0.574s |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | ok 2.849s |
| `wire diff ./internal/composition/...` | exit 0 — no constructor signature changed, `wire_gen.go` still current |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | 42 issues (40 `gochecknoglobals`, 2 `dupl`) — back to the Phase 3 baseline after `--fix` cleared the 18 `gci`/`gofmt`/`golines` findings the mechanical edits introduced |
| Unit coverage | 64.7% — unchanged |

Coverage held rather than rose: the deleted nil branches were exercised by the very tests that were
deleted with them, so numerator and denominator moved together again. 64.7% remains the standing
baseline for Phase 6.

## Phase 5: Prebuild The Topology Lookup

Status: Complete

The constructor collapse itself moved to Phase 1.5. What remains here is the performance half: the
provider still allocates a fresh topology service on every call inside the auto-regen loop.

- [x] **First, classify each service.** Stateless -> singleton, built once and shared. Holds per-call
      mutable state -> stays transient behind a factory function. They currently receive
      `configuration`, `playerLabels`, `neutralZones` and `tuning` per call, which suggests stateless,
      but `RandomTopologyService` and the `TopologyBase` embedders must be read before deciding.
      `TournamentTopologyService` assigns `this.clusterService` inside `CreateTopologyVariant`, so it
      **is** stateful as written — either keep it transient or lift that switch into a local variable.
      Record the classification for each service in this phase's summary.
- [x] `TopologyProvider` builds the lookup once in its constructor and
      [`CreateTopologyVariant`](internal/services/template_generator/providers/topologyProvider.go#L22-L36)
      resolves through it instead of constructing. Singletons are looked up directly; transients are
      looked up as a factory func and invoked. Preserve the two-player tournament short-circuit and
      the `default:` → ring fallback exactly. **Author decision: the lookup is built by the composition
      provider and injected**, see the summary.
- [x] Wire the lookup through `composition/topologyServiceProvider.go`.

### Verification Plan

- `grep -r "WithCreationServices"` returns nothing.
- Golden-template tests still pass — topology output must be byte-identical:
  `go test ./test/unit/internal/services/template_generator/... -count=1`.
- `go test -tags=integration_test ./test/integration/... -count=1` — green.
- Benchmark comparison against the Phase 0 numbers:
  `go test -v -tags=integration_test -bench=BenchmarkEditorWindow_TabCycling -run=xxx ./test/performance/... -benchmem -benchtime=20x -timeout=120s`
  — allocations per op should drop, and must not rise.

### Phase Summary

`CreateTopologyVariant` no longer constructs anything. Every topology service is built once, at
composition time, and the hot path is a map lookup plus a call through a func value.

**Classification (all read before deciding).**

| Service | Verdict | Evidence |
| --- | --- | --- |
| `Ring`, `Chain`, `SharedWeb`, `Hub`, `GeometricHub` | **Stateless -> singleton** | embed only `base.TopologyBase`; no field assignment outside the constructor |
| `Random`, `Circles`, `Square`, `Geometric`, `Cross`, `Fractal` | **Stateless -> singleton** | embed `PositionedTopologyBuilder` (itself only `TopologyBase`); `Circles` additionally holds a `PositionLayoutService`, which is `struct{}` |
| `TopologyBase` | Stateless | `ZoneLabelProvider`, both factories and `connectionService` are set once in `NewTopologyBase` |
| `tournament_variant` `Hub`/`Balanced`/`Ring`/`Chain` cluster services | **Stateless -> singleton** | same shape; no per-call field writes |
| `TournamentTopologyService` | **Was stateful, now stateless -> singleton** | `this.clusterService` was assigned inside `CreateTopologyVariant`; see below |

The only `this.X = ...` hits outside constructors in the whole topology tree are in
`geometricHubLayout.go` and `geometryHelpers.go` — both are per-call helper objects
(`newGeometricHubLayout`, `newPairBuilder`) created inside a single `CreateTopologyVariant` call, not
services. Nothing needed a transient factory func, so the "transients are looked up as a factory func"
half of the checklist did not apply.

**`TournamentTopologyService` was made stateless.** The four cluster services are now built in the
constructor and held in four dedicated fields; the `switch` moved into a private
`selectClusterService(config.MapTopology)` that returns one of them into a local variable. The
`zoneFactory`/`roadFactory` fields it kept only to feed that switch are gone. The case-to-service
mapping — including the `default:` chain-per-cluster fallback for Chain/SharedWeb/Random — is byte-for-
byte the same.

**Where the lookup is assembled — author decision.** The plan's decision table
("a hand-written provider function builds the enum-keyed lookup; wire only calls it") and the Phase 5
checklist ("`TopologyProvider` builds the lookup in its constructor") could not both hold without one
artifact becoming the identity wrapper Phase 2 rejected. The author chose the decision-table reading:

- [internal/composition/topologyServiceProvider.go](internal/composition/topologyServiceProvider.go) —
  `provideTopologyServices(zoneFactory, roadFactory)` constructs all twelve services once and returns
  the lookup. Added to `GenerationSet`; `wire_gen.go` now calls it exactly once.
- [providers/topologyServiceLookup.go](internal/services/template_generator/providers/topologyServiceLookup.go) —
  `TopologyServiceLookup` owns the enum-to-service mapping. `Resolve` returns the mapped creator or
  falls back to ring, which reproduces the old `default:` arm for both `config.TopologyRing`
  (whose value is the string `"Default"`) and any unknown topology. `Tournament()` exposes the
  short-circuit creator separately, because tournament mode is selected by generation mode, not by
  topology.
- [providers/topologyVariantCreator.go](internal/services/template_generator/providers/topologyVariantCreator.go) —
  the one func type every `CreateTopologyVariant` implements directly.
- `NewTopologyProvider(services *TopologyServiceLookup)`; `TopologyProvider` keeps
  `shufflePlayerZones`, `copyLabels` and the two-player tournament short-circuit unchanged.

The mapping stays in `internal/services` rather than in `composition` on purpose: which enum picks
which service is generation behaviour, not wiring. Composition only constructs.

**Signatures were unified instead of adapted.** The three `CreateTopologyVariant` shapes were first
reconciled by two private adapters (`newTournamentVariantCreator`, `newHubCityVariantCreator`). The
author then moved the `configuration.IsHubCityToHold()` lookup inside `hubTopology.createZones` and
`geometricHubTopology.createZones`, making their `hubIsHoldCity bool` parameter redundant. All twelve
services now declare the identical `TopologyVariantCreator` signature — hub/geometricHub and tournament
take `_ string` for the hold-city label they do not use — so both adapters were deleted and the map
stores bare method values. Behaviour is unchanged; the affected unit tests pass `""` instead of `false`.

**Test seams.** [test/test_helpers/topologyServiceLookup.go](test/test_helpers/topologyServiceLookup.go)
mirrors the composition provider, and
[test/test_helpers/topologyProvider.go](test/test_helpers/topologyProvider.go) wraps it as
`NewTopologyProvider()`. This follows the `NewZoneFactories()` / `NewZoneEditorService()` precedent from
Phases 1.5 and 4 — the twelve-line construction list is duplicated in the helper package rather than
adding a convenience constructor to production code, so production callers still have to supply the
lookup explicitly. All 11 test call sites were retargeted.
[test/test_helpers/templateGenerator.go](test/test_helpers/templateGenerator.go) was added later so the
new benchmark and the `templateGenerator` unit suite share one generator arrangement.

**New tests** in `test/unit/internal/services/template_generator/providers/topologyServiceLookup/`:
`newTopologyServiceLookup_test.go`, `resolve_test.go` (mapped-topology table, hub resolution,
unmapped-and-`TopologyRing`-fall-back-to-ring, and the `IsHubCityToHold()` pass-through) and
`tournament_test.go`. `topologyVariantCreator.go` is a bare type declaration and needs none
(AGENTS.md §4.6). `provideTopologyServices` and `selectClusterService` are private and covered
indirectly — both report 100%.
**Verification**

| Check | Result |
| --- | --- |
| `grep "WithCreationServices"` | no matches outside this plan and the carry-forward doc |
| `grep "topology.New"` in `topologyProvider.go` | no matches — the hot path constructs nothing |
| `go build ./...` | clean |
| `go vet -tags='integration_test,gui' ./...` | clean |
| `go test ./test/unit/... -count=1` | green, including the golden-template suites |
| `go test -tags=integration_test ./test/integration/... -count=1` | ok 2.667s |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | ok 2.936s |
| `wire diff ./internal/composition/...` | exit 0 |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | 42 issues (40 `gochecknoglobals`, 2 `dupl`) — the standing baseline, after `--fix` cleared the `gci`/`gofmt`/`golines` findings the new files introduced |
| Unit coverage | **64.6%** — 0.1 pp below the 64.7% baseline, purely the arithmetic effect of deleting the two fully-covered adapters; no function lost coverage and every new or rewritten function reports 100% |

**Benchmark.** `BenchmarkEditorWindow_TabCycling` clicks through editor tabs and never triggers
generation, so it does not exercise the topology path at all — it is not a valid before/after measure
for this phase. Four headless runs gave 7.58–9.84 ms/op (heavy run-to-run noise) at a flat
**8360–8364 allocs/op**, confirming only that nothing regressed. Phase 3's recorded 4.58 ms/op is not
comparable: that run used the task's `-args headed`.

[test/performance/template_generation_test.go](test/performance/template_generation_test.go) was added
at the author's request to measure the real path — `BenchmarkTemplateGenerator_Generate` runs a full
`TemplateGenerator.Generate()` per topology and needs no window. Post-change figures
(`-benchtime=50x -benchmem`):

| Case | ns/op | B/op | allocs/op |
| --- | --- | --- | --- |
| Ring (8p) | 127150 | 60339 | 608 |
| HubAndSpoke (8p) | 88476 | 80251 | 637 |
| GeometricHub (8p) | 107398 | 67605 | 619 |
| Fractal (8p) | 122438 | 68731 | 707 |
| Tournament (2p) | 23200 | 18939 | 163 |

These are the new baseline for Phase 6. The allocation win itself is structural and visible in the
diff — each `CreateTopologyVariant` call previously allocated a topology service plus its
`TopologyBase` (which itself allocates a `ZoneLabelProvider` and a `topologyConnectionService`), and in
tournament mode a cluster service on top; it now allocates none.

## Phase 6: Final Verification And Documentation

Status: Complete

- [x] `go build ./...`.
- [x] `go test ./test/unit/... -count=1`.
- [x] `go test -tags=integration_test ./test/integration/... -count=1`.
- [x] `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1`.
- [x] Coverage report; confirm the total is at or above the Phase 0 baseline and that
      `internal/composition` is covered by the injector test.
- [x] `golangci-lint-v2 run ./... --issues-exit-code=0`; confirm the only new findings are the expected
      `gochecknoglobals` entries for the wire provider sets, and record the new baseline count.
- [x] Confirm `wire_gen.go` regenerates byte-identically from a clean checkout.
- [x] Confirm `wireinject` appears in **no** build or test command, no `GOFLAGS`, and no VS Code
      `go.buildTags` / `go.testTags`.
- [x] Update `README.md` / `QUICKSTART.md` only if they document the build or contribution workflow.
- [x] Leave everything unstaged for author review.

### Verification Plan

- Every command above exits successfully with the recorded expectations met.

### Phase Summary

**Verification**

| Check | Result |
| --- | --- |
| `go build ./...` | clean |
| `go vet -tags='integration_test,gui' ./...` | clean |
| `go test ./test/unit/... -count=1` | green |
| `go test -tags=integration_test ./test/integration/... -count=1` | ok 0.841s |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | ok 1.382s |
| `wire gen` then `Get-FileHash internal/composition/wire_gen.go` | byte-identical before and after |
| `wire diff ./internal/composition/...` | exit 0 |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | 42 issues (40 `gochecknoglobals`, 2 `dupl`) |
| Unit coverage | **64.6%** |

**Coverage — 0.2 pp below the Phase 0 baseline, and why that is the correct outcome.** The trajectory
was 64.8% (Phase 0 and 2) → 64.7% (Phase 3) → 64.7% (Phase 4) → 64.6% (Phase 5). Both 0.1 pp steps are
denominator arithmetic from deleting *covered* code, documented in place: Phase 3 removed
`NewDefaultGuiHandler` and `newDefaultPreviewGenerator` (~22 covered statements that duplicated what
`wire_gen.go` builds), Phase 5 removed the two fully-covered signature adapters. No statement that was
covered at Phase 0 is uncovered now, and no uncovered code was added. `internal/composition` is covered
by the injector test: `wire_gen.go InitializeGuiHandler` 100%, `provideTopologyServices` 100%,
`providePreviewGenerator` 60% — the missing 40% is the asset-load failure branch, which is only
reachable with a broken installation and would need a production seam to force (AGENTS.md §4.6
forbids that).

**Lint.** 42 issues against the Phase 0 baseline of 84 `gochecknoglobals` — a net improvement of 42.
The six new findings are exactly the expected ones, all in
[internal/composition/providerSets.go](internal/composition/providerSets.go): `ZoneSet`,
`GenerationSet`, `EditorSet`, `InfrastructureSet`, `HandlerSet`, `GuiHandlerSet`. `wire.NewSet` values
must be package-level, so these are inherent to the tool. The 2 `dupl` findings are pre-existing in
[app/gui/widgets/buttonWidget.go](app/gui/widgets/buttonWidget.go) and unrelated to this work.

**`wireinject` containment.** The tag appears only in the two build constraints
([wire.go](internal/composition/wire.go) `//go:build wireinject`,
[wire_gen.go](internal/composition/wire_gen.go) `//go:build !wireinject`), in the AGENTS.md §4.6.2
prose, in a comment in [.vscode/tasks.json](.vscode/tasks.json), and in the `cSpell.words` list in
[.vscode/settings.json](.vscode/settings.json). `go env GOFLAGS` is empty; `gopls.build.buildFlags` is
`-tags=integration_test,gui` only; no `go.buildTags` / `go.testTags` entry exists.

**Documentation.** [README.md](README.md) gained two things: a *Building & Running* paragraph stating
that `wire_gen.go` is committed (so a plain `go build` needs no extra step), the
`wire gen ./internal/composition/...` regeneration command, and the warning never to pass
`-tags=wireinject`; and an eighth *Architecture → Layers* entry for `internal/composition`.
[QUICKSTART.md](QUICKSTART.md) is end-user documentation (`go run .` / `go build .`) and documents no
contribution workflow, so it was left unchanged. AGENTS.md §4.6.2 was already written in Phase 1 and
remains accurate.

**Working tree.** Only `README.md` is modified; nothing staged, nothing committed by the agent.

## Risks And Watch-Outs

- **`wire_gen.go` in the coverage denominator.** The injector test should cover nearly all of it, but
  if wire emits cleanup/error paths that the happy path never reaches, coverage can still dip. Measure
  at the end of Phase 2, not at the end of the effort.
- **The two `//go:build` worlds.** `wire.go` and `wire_gen.go` declare the same symbol. Any tool
  configured with `wireinject` will report a duplicate declaration. This is expected, not a bug.
- **Behavior parity is non-negotiable.** Generated `.rmg.json` output must be unchanged. The golden
  fixture tests in `test/unit/internal/services/template_generator/` are the guard.
- **Phase 5 sharing assumption.** Stateless topology services become singletons; any that hold
  per-call mutable state stay transient behind a factory function. Determine which is which by reading
  the service, before the collapse - a stateful service shared as a singleton would leak state across
  generations in a way the golden tests may not catch.

## Final Recap

The application now has a single compile-time composition root. Before this work, dependencies were
constructed ad hoc: `NewGuiHandler` built one set of collaborators and `NewTemplateGenerator` silently
built a second copy of three of them, two call sites each created an independent handler graph, and
`TopologyProvider` allocated a fresh topology service on every call inside a 300 ms auto-regeneration
debounce loop.

| Phase | Outcome |
| --- | --- |
| 0 | Repaired the handler test suite left broken by an in-flight refactor; established the 64.8% coverage baseline every later phase was measured against. |
| 1 | Adopted `github.com/goforj/wire` (maintained fork of the archived `google/wire`), recorded the CLI in `tools/go.mod`, added the *Go: Generate wire injectors* task, and documented `wireinject` in AGENTS.md §4.6.2 as a codegen-only tag. |
| 1.5 | Deleted the `CreationServices` aggregate; constructors now take the collaborators they actually use. |
| 2 | Created `internal/composition` — six provider sets, the `//go:build wireinject` injector stub, the committed `wire_gen.go`, and an injector test. |
| 3 | Collapsed to one composition root: `window.go` and `drivers/state.go` both go through `composition.InitializeGuiHandler()`; `NewDefaultGuiHandler` and `newDefaultPreviewGenerator` deleted. |
| 4 | Removed nil-defaulting fallbacks — a missing dependency is now a compile error, not a silent second instance. |
| 5 | Prebuilt the topology lookup: twelve singleton services behind `TopologyServiceLookup`, `TournamentTopologyService` made stateless, all twelve `CreateTopologyVariant` signatures unified on `TopologyVariantCreator`, and a `TemplateGenerator.Generate` benchmark added. |
| 6 | Final verification and documentation. |

**Behaviour parity held throughout.** Generated `.rmg.json` output is unchanged — the golden-template
suites in `test/unit/internal/services/template_generator/` guarded every phase and stayed green.

**Net measurements.** Lint 84 → 42 issues. Coverage 64.8% → 64.6%, entirely from deleting covered
duplicate code (see the Phase 6 summary). `CreateTopologyVariant` went from allocating a topology
service plus its `TopologyBase` — and in tournament mode a cluster service — on every call, to
allocating none; `BenchmarkTemplateGenerator_Generate` in
[test/performance/template_generation_test.go](test/performance/template_generation_test.go) is the
standing baseline for that path.

## Deployment Plan

This is a desktop application built from source; there is no server or release pipeline to coordinate.

1. **Review the working tree.** `git status --short` — the agent staged and committed nothing. Review
   the diff, then stage and commit on `AD/refactoring-07-21`.
2. **Confirm `wire_gen.go` is included in the commit.** It is generated but committed on purpose; a
   checkout without it does not build.
3. **Run the gate locally before pushing:**

   ```powershell
   go build ./...
   go vet -tags='integration_test,gui' ./...
   go test ./test/unit/... -count=1
   go test -tags=integration_test ./test/integration/... -count=1
   go test -tags='integration_test,gui' ./test/integration/gui/... -count=1
   wire diff ./internal/composition/...
   ```

   `wire diff` exiting 0 proves the committed `wire_gen.go` matches its providers. Note that the `wire`
   CLI writes its banner to stderr, so PowerShell may surface a `NativeCommandError` on success — judge
   by the exit code, not the banner.
4. **CI.** The existing PR validation (`go build ./...`, `go vet -tags=integration_test ./...`,
   `go test -race ./test/unit/...`) needs no change. `wire` is **not** required on CI agents because
   `wire_gen.go` is committed; optionally add `wire diff` as a drift check, which would require
   `go install github.com/goforj/wire/cmd/wire@latest` in the job.
5. **Contributor onboarding.** Anyone changing a provider set or a constructor signature must install
   the CLI (`go install github.com/goforj/wire/cmd/wire@latest`) and run the *Go: Generate wire
   injectors* task. README *Building & Running* and AGENTS.md §4.6.2 both state this.
6. **Rollback.** Purely a `git revert` of the merge; there is no migration, no persisted format change
   and no state to unwind — `.gen.json` and `.rmg.json` are byte-compatible with the previous build.
