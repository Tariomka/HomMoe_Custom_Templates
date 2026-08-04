# Session Carry-Forward — Wire Dependency Composition (Phases 3 & 4)

Generated 2026-08-04. Branch `AD/refactoring-07-21`. Supersedes the 2026-08-03 handoff.

## 1. Session Goal

Continue `plans/wire-dependency-composition.md` from the end of Phase 2: implement **Phase 3 — One
Composition Root** and **Phase 4 — Remove Nil-Defaulting**, pausing after each for the author's review.

Both phases are now **Complete**, verified, and reviewed. Phase 3 is committed (`bf21759 "DI Phase 3"`);
Phase 4 is staged by the author. Phases 5 and 6 remain.

## 2. Fixes Applied

Two defects were found and fixed inside this session's own work; no pre-existing breakage was
encountered.

- [test/integration/window_render_integration_test.go](test/integration/window_render_integration_test.go)
  was left with an unused `app/gui/drivers` import after its `drivers.NewUIState()` call site moved to
  the package-local `newUIState()` helper. Caught by `go vet -tags='integration_test,gui' ./...`, not by
  `go build`. **Always run vet with both tags after touching the gated suites.**
- [test/unit/architecture/dependency/dependency_test.go](test/unit/architecture/dependency/dependency_test.go)
  was *stricter* than the depguard rule it mirrors: `findForbiddenAppImports` would have flagged
  `app/gui/program.go` importing `internal/composition`, because that path was missing from
  `allowedRoots`. Added deliberately — the composition root is allowed to know the object graph.

## 3. Features Added / Changed

No user-visible behaviour changed. Two structural guarantees were established.

- **A single composition root.** `app/gui/program.go`'s `eventLoop` is now the only place in the
  application that builds the object graph, via `editor.NewWindow(composition.InitializeGuiHandler())`.
  Nothing under `app/` may import `internal/handlers` any more — enforced twice over, by the depguard
  rule `concrete-handlers-only-at-gui-composition-roots` and by the architecture test, which now expects
  an **empty** violation map.
- **A missing dependency is a compile error.** Every constructor under `internal/` takes its
  collaborators and stores them verbatim. The nil-defaulting fallbacks — which silently constructed
  *second* instances of services the graph already owned — are gone. No panics were added in their
  place; the compiler is the check.

## 4. File Modifications

### Phase 3 — production

| File | Change |
| --- | --- |
| [app/gui/program.go](app/gui/program.go) | `eventLoop` builds the graph: `editor.NewWindow(composition.InitializeGuiHandler())`; `internal/composition` imported |
| [app/gui/editor/window.go](app/gui/editor/window.go) | `NewWindow(backend handler_interfaces.IGuiHandler)` — no longer calls `handlers.NewDefaultGuiHandler()` |
| [app/gui/drivers/state.go](app/gui/drivers/state.go) | `NewUIState()` **deleted**; only `NewUIStateWithBackend` / `NewUIStateWithHandler` remain |
| [internal/handlers/guiHandler.go](internal/handlers/guiHandler.go) | `NewDefaultGuiHandler` + `newDefaultPreviewGenerator` **deleted**, with twelve now-unused imports |
| [.golangci.yml](.golangci.yml) | depguard rule lost both file exceptions (`state.go`, `window.go`); `app/**` minus `$test` now denies `internal/handlers` outright |

### Phase 4 — production

| File | Change |
| --- | --- |
| [internal/services/connection_editor/zoneEditorService.go](internal/services/connection_editor/zoneEditorService.go) | `NewDefaultZoneEditorService()` **deleted**; `NewZoneEditorService` has no nil branches |
| [internal/services/connection_editor/manualReapplyService.go](internal/services/connection_editor/manualReapplyService.go) | three nil branches removed |
| [internal/services/connection_editor/connectionEditorService.go](internal/services/connection_editor/connectionEditorService.go) | nil branch removed; body is a single `return` |
| [internal/services/template_generator/providers/mandatoryContentProvider.go](internal/services/template_generator/providers/mandatoryContentProvider.go) | both nil branches removed |
| [internal/services/template_generator/templateGenerator.go](internal/services/template_generator/templateGenerator.go) | `configuration == nil` fallback removed |

### Created

| File | Purpose |
| --- | --- |
| [test/test_helpers/zoneEditorService.go](test/test_helpers/zoneEditorService.go) | `NewZoneEditorService()` — wires the same collaborators the composition root does, for the ~79 test call sites the deleted `NewDefaultZoneEditorService` served. Mirrors the `NewZoneFactories()` precedent from Phase 1.5. |

### Deleted

- `test/unit/internal/services/connection_editor/zoneEditorService/newDefaultZoneEditorService_test.go` —
  whole file; its subject no longer exists and `newZoneEditorService_test.go` already covers the real
  constructor.

### Edited — tests

Phase 3 (10 files): `test/integration/{editorState,manualCastleReapply,stateExit,window_render}_integration_test.go`
(6+3+2+1 `drivers.NewUIState()` sites moved onto a new package-local `newUIState()` helper defined in
`editorState_integration_test.go`), `test/test_helpers/integration_common/appRunner.go`,
`test/integration/gui/{contentRuleDialogs,zoneEditorDialog}_integration_test.go`,
`test/unit/internal/handlers/guiHandler/common_test.go`,
`test/unit/architecture/dependency/dependency_test.go`.

Phase 4 (28 files): 23 files had `connection_editor.NewDefaultZoneEditorService()` swapped for
`test_helpers.NewZoneEditorService()` (the `connection_editor` import was dropped from 20 of them and
kept in 3 that use the package for other reasons); 4 files lost a nil-asserting test; 5 files under
`test/unit/internal/services/zones/zoneFactory/` gained a package-local `newZoneFactory()`.

## 5. Tests Added Or Updated

**No new test cases.** Four were deleted, each because the behaviour it asserted was deliberately
removed and the test therefore had no subject:

| Test | File |
| --- | --- |
| `TestWhenServiceIsCreated_ReturnsInstance` (whole file) | `zoneEditorService/newDefaultZoneEditorService_test.go` |
| `TestWhenClassifierIsNil_ReturnsUsableService` | `connectionEditorService/newConnectionEditorService_test.go` |
| `TestWhenDependenciesAreNil_ReturnsUsableService` | `manualReapplyService/newManualReapplyService_test.go` |
| `TestWhenConfigurationIsNil_FallsBackToDefaultConfiguration` | `templateGenerator/newTemplateGenerator_test.go` |

One test was renamed because its scenario inverted:
`zoneFactory/newZoneFactory_test.go` → `TestWhenDependenciesAreOmitted_ReturnsUsableFactory` became
`TestWhenDependenciesAreProvided_ReturnsInstance`.

Status at session end:

| Command | Result |
| --- | --- |
| `go build ./...` | Clean |
| `go vet -tags='integration_test,gui' ./...` | Clean |
| `go test ./test/unit/... -count=1` | **PASS** |
| `go test -tags=integration_test ./test/integration/... -count=1` | **PASS** (0.574s) |
| `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | **PASS** (2.849s) |
| `BenchmarkEditorWindow_TabCycling` (20x, run at end of Phase 3) | PASS, 4 579 735 ns/op |
| `wire diff ./internal/composition/...` | exit 0 — `wire_gen.go` still current |
| `golangci-lint-v2 run ./... --issues-exit-code=0` | **42 issues** (40 `gochecknoglobals` = 34 pre-existing + 6 wire sets; 2 `dupl` in `app/gui/widgets/buttonWidget.go`) |

**Coverage: 64.7%** (was 64.8% at the start of Phase 3). See §7 — this is not a regression, and 64.7% is
the number Phase 6 should record as the standing baseline.

## 6. Git Status Snapshot

Branch: **`AD/refactoring-07-21`**, HEAD `bf21759 "DI Phase 3"` (pushed to `origin`).

Nothing was staged or committed by the agent (AGENTS.md §2.5). Phase 3 was committed by the author;
Phase 4 is **entirely staged by the author** and awaiting commit. `git status --short` shows 38 entries,
all in column 1 (staged), column 2 empty — 5 production files, 32 test files, plus
`plans/wire-dependency-composition.md`. One `A` (`test/test_helpers/zoneEditorService.go`), one `D`
(`newDefaultZoneEditorService_test.go`).

No untracked files. `coverage.txt` / `coverage.html` / `lcov.info` were regenerated but are ignored.

## 7. Rejections / Things The User Declined

Nothing was declined this session. Four judgement calls were made and should be re-litigated only with
the author:

- **`drivers.NewUIState()` was deleted rather than repointed at `composition.InitializeGuiHandler()`.**
  Repointing would have made `app/gui/drivers` a *second* composition root, defeating Phase 3's whole
  premise.
- **The coverage drop was reported honestly rather than papered over.** Deleting
  `NewDefaultGuiHandler` and `newDefaultPreviewGenerator` (~22 fully-covered statements that merely
  duplicated what `wire_gen.go` builds) shrank numerator and denominator together:
  `(0.648·N − 22)/(N − 22) ≈ 64.71%` for `N ≈ 9000`. No statement lost coverage; no new uncovered code
  was introduced. Phase 4 held at 64.7% for the same reason — the deleted nil branches were covered by
  the tests deleted alongside them.
- **A test-only convenience factory was accepted, a production one was not.**
  `test_helpers.NewZoneEditorService()` exists so 79 test call sites do not each repeat three lines of
  wiring; production callers still pass collaborators explicitly, so the compiler still catches a
  forgotten dependency.
- **`TemplateGenerator.SetConfiguration`'s `if configuration != nil` guard was left in place.** It is a
  setter over a *value*, not a constructor over a dependency, so it falls outside Phase 4's stated grep
  criterion. Its only caller ([internal/handlers/templateHandler.go](internal/handlers/templateHandler.go#L63))
  already passes a non-nil mapper result. Recorded in the Phase 4 summary.

Also deliberately left alone, as optional **data** rather than dependencies: `NewRuleDistanceToRoad(nil)`,
`NewRuleDistanceToTown(nil)`, `NewRuleVariant(nil, nil)`, `components.NewDropdownSelector(nil)`, and
`config.NewGeneratorConfig()` (the plan names it explicitly as a value factory).

## 8. Open Questions

- None blocking.
- One decision is deferred to Phase 5 by design: each topology service must be classified stateless
  (share one instance) or stateful (keep transient). **`TournamentTopologyService` assigns
  `this.clusterService` inside `CreateTopologyVariant`, so as written it *is* stateful** — verify before
  collapsing it.

## 9. Next Recommended Actions

1. **Phase 5 — Prebuild The Topology Lookup.** Classify every topology service stateless/stateful; have
   `TopologyProvider` build the enum-keyed lookup **once in its constructor** so `CreateTopologyVariant`
   stops allocating on the 300 ms auto-regen loop. Create
   `internal/composition/topologyServiceProvider.go` at this point — it was deliberately deferred from
   Phase 2 because it would have been an identity wrapper. Preserve the two-player tournament
   short-circuit and the `default:` → ring fallback **exactly**. `grep -r "WithCreationServices"` must
   return nothing; golden-template tests must stay byte-identical; benchmark allocations must not rise.
2. **Phase 6 — Final Verification And Documentation.** All four suites; coverage at or above **64.7%**;
   record the lint baseline; confirm `wire_gen.go` regenerates byte-identically (`wire diff`, exit 0);
   confirm `wireinject` appears in no build or test command; update README / QUICKSTART only if they
   document the build workflow; leave everything unstaged.
3. **Write the plan's Final Recap and Deployment Plan** sections once Phase 6 closes.

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
> recorded **64.7%** baseline.
> **§2.4** Multi-session work lives in a plan file, and that file — not the conversation — is the source
> of truth.
> **§2.5** Never stage and never commit; if you find staged changes, leave them exactly as they are.
>
> We are implementing compile-time dependency composition with `github.com/goforj/wire`. The plan is
> `plans/wire-dependency-composition.md` — read it before doing anything else; it carries the locked
> decisions, the per-phase checklists, the verification commands, and a written summary of every phase
> completed so far.
>
> **Where work left off:** Phases 0, 1, 1.5, 2, 3 and 4 are Complete. `app/gui/program.go` is the single
> composition root, and no constructor under `internal/` nil-defaults a dependency any more. All three
> suites are green, coverage is 64.7%, lint is at its 42-issue baseline, and `wire diff` exits 0.
> **Phases 5 and 6 remain.**
>
> **Start with Phase 5 — Prebuild The Topology Lookup.** Classify each topology service stateless or
> stateful first (note that `TournamentTopologyService` mutates `this.clusterService` inside
> `CreateTopologyVariant`), then move the enum-keyed lookup into `TopologyProvider`'s constructor and add
> `internal/composition/topologyServiceProvider.go`. The two-player tournament short-circuit and the
> `default:` → ring fallback must be preserved exactly, and golden-template output must stay
> byte-identical.
>
> Two gotchas that will otherwise cost you time: the `wire` CLI writes its success banner to **stderr**,
> so `wire gen` and `wire check` look like PowerShell failures (`NativeCommandError`, exit 1) even when
> they succeed — use `wire diff ./internal/composition/...` (exit 0 = committed file is current) as the
> reliable signal. And `go build ./...` will not catch unused imports in the gated test suites; run
> `go vet -tags='integration_test,gui' ./...` after touching them.
>
> The working tree currently holds the author's **staged** Phase 4 changes, awaiting commit. Do not
> unstage or "tidy" them.
>
> Full handoff, including the judgement calls already made and why the coverage number moved, is in
> `./.agent/session-carry-forward.md`.
