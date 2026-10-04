# Batch K: Topology retirement (review §2.3 / O08)

Retire the Ring, Hub, Chain and Shared Web topologies and their Ring/Hub/Chain tournament
builders. Every tournament uses the balanced builder. Saved states with a retired or unknown
topology are rejected at load with a clear message, and the ordinary lookup gets an explicit,
panic-free invalid-selection contract. Binding sources: review
[§2.3](../backlog/review-gpt-6-astra-09-07.md) and handoff §8 (§2.3/O08) and §9.

## For Future Agents
As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Context: branch `AD/pbi_resolution`, clean tree at HEAD `61ecb3b` when this plan was written.
Read [AGENTS.md](../../AGENTS.md) and [the handoff](../session-carry-forward.md) first. Never
stage, commit, push, stash or switch branches. Never touch `data/`,
`internal/entities/template_entity/` or `internal/registry/`. Never hand-edit `wire_gen.go`.
Never keep a blanket `-update` snapshot run. Searches use the editor's workspace search or
PowerShell `Get-ChildItem -Recurse -File -Filter *.go | Select-String -Pattern '<regex>'`
(no Unix `grep`/`rg` on this machine).

**Every phase must end with build, vet and the untagged and tagged suites green.** Test
rewrites therefore land in the same phase as the production change that would break them.

### Owner decisions (2026-10-04, binding)

| ID | Decision |
| --- | --- |
| K1 | Delete the four constants from `internal/entities/topology/mapTopology.go` and their aliases in `internal/models/config/types.go`. |
| K2 | Legacy recognition is a private table in `internal/common/common_topologies` (raw ID → former label: `Default`→Ring, `HubAndSpoke`→Hub, `Chain`→Chain, `SharedWeb`→Shared Web), exposed only through `GetRetiredTopologyLabel(id config.MapTopology) (string, bool)`. |
| K3/F1 | `ValidationIssue` gains a generic blocking flag; blocking issues have no fix. `stateHandler.LoadState` returns an error whenever any issue is blocking, regardless of `fixIssues`, so the GUI load fails before anything is applied. `GenerateTemplate` returns the same error. `ValidateEditorState` lists blocking messages in `Warnings` but never fixes them. Only the topology check uses the flag. |
| K4/F2 | One sentinel `common_errors.ErrUnsupportedTopology` (`"unsupported topology"`), wrapped (so `errors.Is` works) but not repeated in the text. GUI output: retired → `Load failed: topology "Ring" (saved as "Default") has been retired and is no longer supported; re-create the template with a supported topology.`; unknown → `Load failed: topology "Foo" is not a known topology; re-create the template with a supported topology.` |
| K5 | An explicit saved `"Default"` is the retired Ring and is rejected. A missing `topology` key keeps loading as Random (seeded entity default, both current and legacy paths). |
| K6 | Unknown IDs are hard-rejected (blocking) with the distinct message above. |
| F3/F3a | Explicit `"topology": ""` is **not** rejected: a non-blocking fixable issue `topology is empty; using Random`, fix → Random. The GUI shows `Loaded … (adjusted: …)`. |
| K7/F4 | `ITopologyServiceLookup.Resolve(id) (TopologyVariantCreator, bool)`; `ITopologyProvider.CreateTopologyVariant(...) (template_model.Variant, error)`, wrapping `ErrUnsupportedTopology` with the ID; `TemplateGenerator.Generate() (*template_model.Template, []string, error)`. `templateHandler.GenerateTemplate` returns that error; its existing nil-template check stays. The Ring fallback is deleted. |
| K8 | `TopologyProvider` checks validity (via `Resolve`) **before** the tournament branch, so an unknown ID fails with the K7 contract in both modes. Tournament dispatch has a single balanced builder; `selectClusterService` and the Ring/Hub/Chain cluster services are deleted; `TournamentTopologyService` takes one `IClusterService` (interface kept). |
| K9 | `GetTopologyDescriptorFromType` falls back to the **Random** descriptor (not Ring). `GetTopologyCapabilities` follows. |
| K10 | Rename `models.TopologyLayoutRingHub` → `TopologyLayoutGeneric` (still iota zero, still the generic preview fallback `layoutRingOrHub`); fix stale comments naming the retired topologies. |
| K11 | `zoneLabelProvider`: remove the Shared Web plan rule and the Chain/Ring adjacency branches; keep the Circles branch and its comment. Remove `CreateBalancedChainZoneLabels` and the `isRing` parameter of `CreateOrderedZoneLabels` (no production caller remains), plus the helpers only they used. |
| K12 | v0/v1/v2 migration fixtures and `test_helpers/allFieldsEditorState.go` switch `Chain` → `Square`. New committed rejection fixtures: one v2 per retired ID plus one legacy v0 `"Chain"`. |
| K13/F5 | Benchmarks: `template_generation` Ring → Circles (8p), HubAndSpoke → Square (8p), Tournament(Chain) → Tournament(Random). `preview_layout` RingLarge → SquareLarge (8p/16n). |
| K14 | GUI coverage: the dropdown is exactly Random, Circles, Geometric Hub, Square, Geometric, Cross, Fractal; a GUI load of a retired file shows `Load failed` and the document is unchanged; `zoneEditorGeometry` tests that select Ring/Hub move to survivors. |
| K15 | Goldens: accept only the `.failure` files a plain run produces; list each with its reason in the phase summary for owner approval. |
| K16 | Coverage: no surviving file loses coverage; all new/changed code fully covered; a total below 74.4% must be explained and needs owner approval. |
| F6 | Unknown-ID coverage uses temp files in unit tests plus one integration test reading an unknown-ID temp file; no extra committed fixture. |
| — | Out of K: §1.12 (Batch E tournament lock untouched), new tournament designs, direct `GeneratorConfig` rejection, DTO/schema/package/output-path work beyond the K3 rejection field below, allocation tuning. Reviews by GPT-6.1 Sol. |

### Design details fixed by this plan (within the owner decisions)

- **D1, the error route.** `EditorStateValidationDto` gains one field, `Rejection error` (nil
  when no issue is blocking). `stateHandler.ValidateEditorState` sets it from the first
  blocking issue; `LoadState` returns `nil, validation.Rejection` when it is set;
  `templateHandler.GenerateTemplate` returns it before the mapper and generator are called.
  No warning text is parsed. `app/gui/models.EditorState.UpdateCurrentState` keeps ignoring
  it: the dropdown only ever writes catalogued types (`GetTopologyDescriptorFromIndex` falls
  back to Random), so a GUI edit cannot produce a blocking topology; a unit test proves it.
- **D2, the error type.** `validators/blockingIssueError.go`: an unexported struct whose
  `Error()` returns the issue message verbatim (no trailing period) and whose `Unwrap()`
  returns the cause (`ErrUnsupportedTopology`). `ValidationIssue` holds it in an unexported
  field and exposes `IsBlocking() bool` and `Rejection() error`; `Fix` is a no-op when
  blocking. The GUI's `"Load failed: %v."` supplies the single terminal period.
- **D3, the generator error.** `TopologyProvider` returns `fmt.Errorf("%w: %q", ErrUnsupportedTopology, id)`.
  It is unreachable from the GUI (D1 rejects first) and is not user-facing copy.
- **D4, legacy loads.** Rejection runs after migration, so v0, v1 and v2 share it. Both
  `NewDefaultEntity` and `newV1Seed` seed Random, so a missing key stays Random on every
  version.
- **Recorded, not assigned:** `SaveState` and `UpdateTemplate` do not validate, so a
  programmatic caller bypassing the GUI is not covered by D1. No reachable GUI path exists;
  out of K.

### Plan reviews (GPT-6.1 Sol)

- Round 1: REJECT, 9 findings (3 High: the generation error route, phase greenness, the
  missing `ITemplateGenerator`/mock/tournament helper; 4 Medium; 2 Low). All applied.
- Round 2: APPROVE WITH CHANGES. 8 of 9 resolved, 1 partial (dedicated `blockingIssueError`
  tests were conditional): applied. D1–D4 verified against the code.

### Inventory (taken at `61ecb3b`; line numbers indicative)

- **IDs:** [mapTopology.go](../../internal/entities/topology/mapTopology.go) L7–L10; aliases in
  [types.go](../../internal/models/config/types.go) L36–L39.
- **Catalogue:** [topologyDescriptors.go](../../internal/common/common_topologies/topologyDescriptors.go)
  (four entries, display order, Ring fallback); [models/topologyDescriptors.go](../../internal/models/topologyDescriptors.go)
  fields; [topologyLayoutKind.go](../../internal/models/topologyLayoutKind.go).
- **Validation/load:** [editorStateValidator.go](../../internal/validators/editorStateValidator.go)
  `validateTopology`; [validationIssue.go](../../internal/validators/validationIssue.go);
  [stateHandler.go](../../internal/handlers/stateHandler.go);
  [editorStateValidationDto.go](../../internal/dtos/editor_state_dto/editorStateValidationDto.go);
  [stateFiles.go](../../app/gui/drivers/stateFiles.go) `handleLoadState` (returns before any
  mutation on error); [editorState.go](../../app/gui/models/editorState.go) `UpdateCurrentState`;
  [fileService.go](../../internal/services/file_service/fileService.go) and the
  `editor_state_migrator` (decode/migrate only).
- **Generator:** [templateGenerator.go](../../internal/services/template_generator/templateGenerator.go),
  [templateGeneratorInterface.go](../../internal/services/template_generator/templateGeneratorInterface.go),
  [templateHandler.go](../../internal/handlers/templateHandler.go) `GenerateTemplate`;
  test double [templateGeneratorMock.go](../../test/test_helpers/templateGeneratorMock.go)
  and every `Generate` return setup in tests.
- **Ordinary dispatch:** [topologyServiceLookup.go](../../internal/services/template_generator/providers/topologyServiceLookup.go),
  [topologyProvider.go](../../internal/services/template_generator/providers/topologyProvider.go),
  `providers/provider_interfaces/{topologyServiceLookupInterface,topologyProviderInterface}.go`,
  hand-written wiring [topologyServiceProvider.go](../../internal/composition/topologyServiceProvider.go)
  (`wire_gen.go` only calls `provideTopologyServices`), `providers/topology/topologyServiceAssertions.go`,
  test wiring [test_helpers/topologyServiceLookup.go](../../test/test_helpers/topologyServiceLookup.go).
- **Retired services:** `providers/topology/{ringTopology,hubTopology,chainTopology,webTopology}.go`.
  No survivor embeds or calls them.
- **Tournament:** [tournamentTopology.go](../../internal/services/template_generator/providers/topology/tournamentTopology.go)
  (`selectClusterService`: Hub→hub, Circles→balanced, Ring→ring, default→chain);
  `tournament_variant/{ringClusterService,hubClusterService,chainClusterService}.go` retire;
  `balancedClusterService.go`, `clusterServiceInterface.go`, `misc/` stay. Test helper
  [tournamentTopologyDependencies.go](../../test/test_helpers/tournamentTopologyDependencies.go)
  builds all four cluster services and returns the 8 constructor arguments.
- **Retirement-only helpers** (each used only by deleted code; delete with their dedicated
  tests under `test/unit/internal/common/constants/connectionNames/` and
  `test/unit/internal/services/zones/utils/misc/`):
  - [connectionNames.go](../../internal/common/constants/connectionNames.go):
    `GetNeutralRingConnectionNameFor`, `GetTournamentRingConnectionNameFor`,
    `GetTournamentHubAndSpokeConnectionNameFor`, `GetTournamentHubRingConnectionNameFor`,
    `GetTournamentChainConnectionNameFor`, `GetWebConnectionNameFor`,
    `GetChainConnectionNameFor`, `GetRingConnectionNameFor`, `GetHubSpokeConnectionNameFor`,
    and their prefixes. Re-check usages after the deletions before removing each.
    `GetHubZoneNameFor`/`HubZonePrefix` stay (Geometric Hub).
  - [zones/utils/misc.go](../../internal/services/zones/utils/misc.go): `OrderEdgeGap` and
    `AssignNeutralZonesToGaps`'s `preferInterior` parameter (its only `true` caller is
    `CreateBalancedChainZoneLabels`); drop the parameter and its branch, keep the function.
- **Shared code that stays:** `base.TopologyBase`, `positionedTopology*`, `geometryHelpers.go`
  (`circlePoint`, `squarePerimeterPoint`, `nearestIndexInRange`, `pairBuilder`),
  `position_layout/`, `geoHubGeometry.go`, `geometricHubLayout.go`,
  `CreateBalancedRingZoneLabels`, `CreateBalancedNeutralRingZoneLabels`, Delaunay helpers.
- **Labels/config:** [zoneLabelProvider.go](../../internal/services/zones/zoneLabelProvider.go)
  and `zones/zone_interfaces/zoneLabelProviderInterface.go`;
  [generatorConfig.go](../../internal/models/config/generatorConfig.go) `IsHubCityToHold`.
- **Preview:** [previewLayoutService.go](../../internal/services/preview_service/previewLayoutService.go)
  `dispatchClusterLayout`; [layoutRingHub.go](../../internal/services/preview_service/layoutRingHub.go)
  comment; `layoutBalancedRings.go` fallback.
- **GUI:** [layoutPanel.go](../../app/gui/panels/layoutPanel.go) (dropdown from the descriptor
  sequence); `UsesHub` consumers stay (Geometric Hub).
- **Tests referencing the retired constants or labels:** run the Phase 1 search. Known heavy
  users: `previewLayoutService/buildPreviewLayout_test.go` (~50× Ring, mostly as "a topology
  without metadata"), `previewGeneratorService/createPreviewImage_test.go`, `templateGenerator/*`
  (`generate`, `generateAllTopologies`, `generateCastles`, `generateStructure` with retired
  labels, `generateTournament`, `generateZoneTiers`), `mandatoryContentProvider/createContents_test.go`
  (Hub content), `guiHandler` preview/save, `common_topologies/topologies/*`, `isHubCityToHold`,
  zone-label tests, `editorState`/`drivers.State` `getTopology`, `equalsIgnoringManualEdits`,
  zone-editor options, generator-config mapper, layout-defining options, regeneration decision,
  `v1ToV2/migrateToV2_test.go`, `previewLayoutCache/get_test.go`, `nullPreviewGeneratorService`,
  `roadPolicyApply_integration_test.go`, GUI `zoneEditorGeometry_integration_test.go`
  (labels "Ring"/"Hub", edge `Ring-A-B`), performance `template_generation_test.go`,
  `preview_layout_test.go`. Raw `"Ring-…"` names in `topologyBase` tests are plain strings,
  not topology references, and may stay.
- **Fixtures:** `test/test_helpers/testdata/editorState_v{0,1,2}_flat.gen.json` (L35 `"Chain"`),
  `test/test_helpers/allFieldsEditorState.go` L69.
- **Docs:** [README.md](../../README.md) topology table; `.agent/memories/generator-domain.md`
  (gitignored). No GUI golden is named after a retired topology; `output/` has no retired IDs.

## Phase 0: Baseline (no Go edits)
Status: Complete

- [x] Record HEAD, `git status --short`, `go version`, `go env GOFLAGS` (must be empty).
- [x] `go build ./...`; `go vet ./...`; `go vet -tags=integration_test ./...`; `go vet -tags=integration_test,gui ./...`.
- [x] Coverage: `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...`, then `go tool cover '-func=coverage.txt'` saved to `tmp/batch-k-coverage-before.txt` (gitignored). Expect 74.5% (6915 / 9263).
- [x] `go test ./test/... -count=1`; `go test -tags=integration_test ./test/integration/... -count=1`; `go test -tags=integration_test,gui ./test/integration/gui/... -count=1` (record times, no `.failure` files).
- [x] Benchmarks (record only) to `tmp/batch-k-bench-before.txt`: `go test -bench=. -run=xxx ./test/performance/... -benchtime=20x -timeout=120s` and `go test -tags='integration_test,gui' -bench=BenchmarkEditorWindow_TabCycling -run=xxx ./test/performance/... -benchtime=20x -timeout=120s`.
- [x] `go run ./cmd/testlayoutcheck .`; `gofmt -l` over Go dirs (never `data/`); lint report-only (0 issues expected).
- [x] Probe: the balanced cluster builder with 0 neutral plans for a player and with uneven per-player counts (odd neutral total) must not panic (`createSortedPairs` slices `tierKeys[:len-1]`). Throwaway test in gitignored `tmp/`, deleted afterwards. If it panics, **stop and ask the owner** before Phase 2.

### Verification Plan
- All commands exit 0; coverage equals the recorded baseline; probe result recorded.

### Phase Summary
Complete, 2026-10-04, Windows/amd64, Go 1.27.0, empty `GOFLAGS`. HEAD `d50f34d` ("plan",
the owner's commit of this file), clean tree, branch `AD/pbi_resolution`.
- Build PASS; `go vet` PASS under no tag, `integration_test` and `integration_test,gui`.
- Coverage PASS: **74.5% (6915 / 9263)**, identical to the record. The func report is in
  `tmp/batch-k-coverage-before.txt`, the raw profile in `tmp/batch-k-coverage-before.out`.
- `go test ./test/...` PASS; tagged integration PASS (root 3.116s); tagged GUI PASS (28.946s);
  0 `.failure` files.
- Benchmarks recorded: `tmp/batch-k-bench-before-nogpu.txt` (GPU-free) and
  `tmp/batch-k-bench-before.txt` (TabCycling, 3239085 ns/op). Key rows: Generate Ring 45160,
  HubAndSpoke 49670, GeometricHub 42385, Fractal 45425, Tournament(Chain) 13890 ns/op;
  BuildPreviewLayout RingLarge 14805, CirclesLarge 49855 ns/op.
- testlayoutcheck PASS; `gofmt -l` clean; lint **0 issues** after `golangci-lint-v2 cache
  clean`. The first run reported 3 `nolintlint` hits in two untouched `editorState` test files,
  which were stale cache artifacts (same v2.13.1 binary).
- Probe (Circles, the balanced builder today): 0, 1, 3 and 5 neutral zones, with portals on and
  off, all ran without a panic (0 → 2 zones, 0 connections). The probe was deleted.
- **PowerShell 5.1 gotcha:** an unquoted `-bench=.` is split into `-bench=` and `.`, which
  silently runs nothing. Quote every flag containing `.` or `,` (`'-bench=.'`,
  `'-tags=integration_test,gui'`).
- **Artifacts:** `tmp/` can be wiped by the owner's `air` watcher, so all four baseline files
  are also copied to gitignored `.agent/memories/batch-k-baseline/`. Compare against that copy.

## Phase 1: Catalogue, validation and load rejection
Status: Complete

Production (the retired services and constants still exist after this phase):
- [x] `common_topologies`: delete the four descriptors and their display-order entries (order: Random, Circles, Geometric Hub, Square, Geometric, Cross, Fractal); `GetTopologyDescriptorFromType` falls back to `descriptorValues.Random` (K9). Add the private retired table and `GetRetiredTopologyLabel` (K2) in a new file `retiredTopologies.go`, keyed by raw `config.MapTopology("Default")` etc. so it survives K1.
- [x] `models.TopologyDescriptors`: drop the four fields. Rename `TopologyLayoutRingHub` → `TopologyLayoutGeneric` (K10) everywhere; fix the stale comments in `layoutRingHub.go` and `previewLayoutService.go`.
- [x] `common_errors`: add `ErrUnsupportedTopology` to `handlerErrors.go`.
- [x] `validators`: `blockingIssueError.go` (D2); `ValidationIssue` gains the unexported rejection, `IsBlocking()`, `Rejection()`, and the no-op `Fix` for blocking issues.
- [x] `validateTopology`: `""` → non-blocking, `topology is empty; using Random`, fix → Random. Retired → blocking with the K4 retired message. Otherwise not catalogued → blocking with the K4 unknown message. Catalogued → nil.
- [x] `EditorStateValidationDto.Rejection` (D1). `stateHandler.ValidateEditorState` sets it and skips blocking fixes; `LoadState` returns `nil, Rejection` when set, regardless of `fixIssues`.
- [x] `templateHandler.GenerateTemplate`: return `validation.Rejection` when set, before the mapper or generator.
- [x] Fixtures: v0/v1/v2 `"Chain"` → `"Square"`; `allFieldsEditorState` → `config.TopologySquare`. New in `test/test_helpers/testdata/`: `editorState_v2_retired_{Default,HubAndSpoke,Chain,SharedWeb}.gen.json` (copies of the v2 fixture, topology only changed) and `editorState_v0_retired_Chain.gen.json` (copy of the v0 fixture).

Tests moved in this phase because the catalogue change breaks them:
- [x] Every test that **validates, describes or selects** a retired topology through the catalogue or the GUI: `generateStructure_test.go` (retired labels), `mandatoryContentProvider/createContents_test.go` (Hub content via `UsesHub`), `common_topologies/topologies/*`, the preview tests whose behaviour depends on `GetTopologyCapabilities` (retired IDs now resolve to Random's capabilities, so Ring-as-"no metadata" cases move to zones without positions or to a survivor with the same path, checked per test), `guiHandler` preview/save, `editorState`/`drivers.State` tests that run the validator, regeneration decision, layout-defining options, zone-editor options, and GUI `zoneEditorGeometry_integration_test.go` (Ring → a survivor per test intent, Hub → Geometric Hub; the `Ring-A-B` edge lookup changes with it). Tests that call a retired **service** directly stay until Phase 2.
- [x] Run the search `TopologyRing\b|TopologyHubAndSpoke|TopologyChain|TopologySharedWeb|"Ring"|"Hub"|"Chain"|"Shared Web"` over `test/` and `app/` and record each remaining hit as "Phase 2" (service-level) or fixed.

New tests (mirrored layout, triple-A, `t.Parallel()`, one assertion each, `gofakeit` where meaningful):
- [x] `test/unit/internal/common/common_topologies/retiredTopologies/getRetiredTopologyLabel_test.go`: each retired ID → its label and true (named subtests); a survivor, an unknown ID and `""` → false.
- [x] `common_topologies/topologies/*`: 7 entries in exact order; `getTopologyDescriptorFromType` unknown → Random; `getTopologyCapabilities` unknown → Random's capabilities.
- [x] `test/unit/internal/validators/validationIssue/{fix,isBlocking,rejection}_test.go`: the blocking `Fix` leaves the state unchanged; the non-blocking `Fix` applies the correction; `IsBlocking`; `Rejection` is nil for non-blocking and `errors.Is(…, ErrUnsupportedTopology)` with the exact text for blocking.
- [x] `test/unit/internal/validators/blockingIssueError/{error,unwrap}_test.go` (unconditional, AGENTS.md §4.6): obtain the error through production validation and `ValidationIssue.Rejection()`; assert the exact text and the `ErrUnsupportedTopology` cause without exporting the private type.
- [x] `editorStateValidator`: retired (table over the 4 IDs) → blocking + exact message; unknown → blocking + exact message; `""` → non-blocking, exact message, fix sets Random; each survivor → no topology issue.
- [x] `stateHandler/validateEditorState_test.go`: blocking listed in `Warnings`, `Rejection` set, state not fixed; non-blocking `Rejection` nil.
- [x] `stateHandler/loadState_test.go`: blocking → nil result and the wrapped error with exact text, for both `fixIssues` values; `""` → loads with the warning.
- [x] `templateHandler/generateTemplate_test.go`: `Rejection` set → returned error; mapper and generator not called.
- [x] `guiHandler/loadState_test.go` (real file service): each committed retired fixture rejected; an unknown-ID temp file rejected (F6).
- [x] `app/gui/models/editorState/updateCurrentState_test.go`: a dropdown-style update to every catalogued topology yields no rejection (D1 unreachability). **Replaced, see summary:** proven by existing and new tests instead of a new test file.
- [x] Integration, through the production load handler (not only the migrator): a table over v0/v1/v2 for missing key → Random, explicit `""` → Random with the warning, and the converted fixtures loading as Square; the four v2 retired fixtures and the v0 Chain fixture → `ErrUnsupportedTopology` with exact text; a v1 retired case from a temp file; an unknown-ID temp file → the unknown message. Existing migrator-level assertions stay. Place it in `editorStateWireFormat_integration_test.go` or a sibling, untagged unless it needs `*_testexports.go` (§4.6.1).
- [x] Unchanged-document-after-failure, legacy and current: covered in Phase 3's GUI test (both a v2 and the v0 fixture). If `drivers.State` load can be unit-tested without test exports, add it here too. **Done here instead, see summary:** `integration_test`-tagged, GPU-free.

### Verification Plan
- `go build ./...`; `go vet` under no tag, `integration_test`, `integration_test,gui`.
- `go test ./test/... -count=1`; `go test -tags=integration_test ./test/integration/... -count=1`; `go test -tags=integration_test,gui ./test/integration/gui/... -count=1` (handle any `.failure` per K15).
- Search `descriptorValues\.(Default|HubAndSpoke|Chain|SharedWeb)|TopologyLayoutRingHub` → no hits.

### Phase Summary
Complete, 2026-10-04. All verification passed: build; `go vet` × 3 tag sets; `go test ./test/...`;
tagged integration; tagged GUI (30.081s, **0 `.failure` files, no golden moved**); forbidden-symbol
search empty; testlayoutcheck PASS; lint 0 issues (fixed 1 `godoclint`, 1 `golines`); new files
`gofmt -w`'d by explicit list (CRLF from file creation).

Production (as planned): catalogue down to 7 in K14 order with the Random fallback;
`retiredTopologies.go` (`GetRetiredTopologyLabel`); `TopologyLayoutGeneric` (one-line doc) and the
`layoutRingOrHub` comment; `ErrUnsupportedTopology`; `blockingIssueError.go`; `ValidationIssue`
`rejection`/`IsBlocking`/`Rejection` with a no-op blocking `Fix`; `validateTopology` (`""`
→ fixable, retired/unknown → blocking, shared `recreateTemplateAdvice` const);
`EditorStateValidationDto.Rejection`; `stateHandler` (first blocking issue wins, `LoadState` fails);
`templateHandler.GenerateTemplate` returns the rejection before mapping.

Deviations:
- **No `updateCurrentState_test.go`.** D1 unreachability already follows from two tested facts:
  `GetTopologyDescriptorFromIndex` falls back to Random (existing tests), and every catalogued
  type validates without a topology issue (new `TestWhenTopologyIsSupported_ReturnsNoTopologyIssue`).
  A GUI-model test would need the real validator behind the handler, which `app/` unit tests
  mock.
- **The unchanged-document check moved here** from Phase 3, as GPU-free
  `test/integration/topologyLoadRejection_integration_test.go` (`integration_test` tag, because it
  uses `LoadStateFromFile`/`SetCurrentPath`). For the v2 `Default` and v0 `Chain` fixtures it pins
  the exact `Load failed: ….` status, the error flag, and an unchanged
  state/path/unsaved/template snapshot. **Negative control:** with the `LoadState` rejection
  disabled, all 13 rejection subtests failed; reverted, and the diff was re-checked.
- **Untagged load matrix** in `test/integration/editorStateTopologyLoad_integration_test.go`,
  through `composition.InitializeGuiHandler().LoadState`. Temp files are written by
  re-marshalling a fixture's keys, so no line-ending assumptions are made.
- **GUI geometry tests.** No survivor's default 2-player template has a doubled pair or a zone on
  a chord (a throwaway `cmd/tmpprobe` printed all 7 layouts, then was deleted). The four Ring/Hub
  tests now draw connections with the editor:
  - grouping uses Geometric Hub plus drawn A→B, B→A and Hub→A;
  - spread and label use Square plus a drawn A→B: control points `(275.15076, 304.84924)` and
    `(304.84924, 275.15076)`, ±21px off the diagonal chord as in Batch F; the label midpoint is
    `(P0 + 2C + P2) / 4`;
  - bulge uses Geometric Hub plus a drawn A→B: control point `(380.97144, 290)`, clearing the
    hub by the same 90.97px the old Hub case did, on the positive side.
- **Preview tests did not fail.** Ring now resolves to Random's capabilities, and those zones carry
  no positions, so they still take the generic layout. They still name the constant, so they move
  in Phase 2.
- **Also moved early, because their file was touched anyway:** `roadPolicyApply` (castleless hub
  arena: Hub → Geometric Hub, passes), the `mandatoryContentProvider` hub tests (Hub → Geometric
  Hub, Ring → Circles), the description tests (Chain → Square; Ring/Hub/Shared Web →
  Circles/Geometric Hub/Cross).
- **Phase 2 search inventory** (files still naming a retired constant): the production
  `mapTopology.go`, `types.go`, `generatorConfig.go`, `topologyServiceLookup.go`,
  `tournamentTopology.go` and `zoneLabelProvider.go`, plus 37 test files. The heaviest are
  `previewLayoutService/buildPreviewLayout_test.go` (56), `templateGenerator/generate_test.go` (28)
  and `previewGeneratorService/createPreviewImage_test.go` (16).

## Phase 2: Generator contract, tournament rewiring and code deletion
Status: Not started

- [ ] `ITopologyServiceLookup.Resolve` → `(TopologyVariantCreator, bool)`; delete the `ring` field and fallback. `NewTopologyServiceLookup` drops the ring/hub/chain/sharedWeb parameters.
- [ ] `ITopologyProvider.CreateTopologyVariant` → `(template_model.Variant, error)`. `TopologyProvider` resolves first; on `!ok` it returns D3's error. Only then is the tournament branch taken.
- [ ] `ITemplateGenerator.Generate` and `TemplateGenerator.Generate` → `(*template_model.Template, []string, error)`; return `nil, nil, err` on a provider error, before the description and content work. `templateHandler.GenerateTemplate` returns that error; the nil-template check stays. Update `TemplateGeneratorMock` and every `Generate` return setup and two-result caller.
- [ ] `TournamentTopologyService`: one `clusterService tournament_variant.IClusterService`; delete `selectClusterService`; the constructor takes one cluster service. Update `internal/composition/topologyServiceProvider.go` and `test/test_helpers/tournamentTopologyDependencies.go` (it builds real services; no mocks assumed) and every consumer of that helper.
- [ ] Delete `ringTopology.go`, `hubTopology.go`, `chainTopology.go`, `webTopology.go`, `tournament_variant/{ring,hub,chain}ClusterService.go`, their assertions and their wiring (composition and `test/test_helpers/topologyServiceLookup.go`).
- [ ] `zoneLabelProvider` (K11): delete the Shared Web plan rule; `createTopologyAdjacency` keeps the Circles branch (and its comment) and the default; delete `CreateBalancedChainZoneLabels` and the `isRing` parameter (interface, implementation, any test double). Then `OrderEdgeGap` and `preferInterior` (inventory).
- [ ] Delete the retirement-only connection-name helpers and prefixes (inventory), each after confirming no remaining caller.
- [ ] `GeneratorConfig.IsHubCityToHold`: drop `TopologyHubAndSpoke`.
- [ ] Delete the four constants (K1) and their aliases.
- [ ] Rewrite benchmark sources per K13/F5 in `test/performance/template_generation_test.go` and `preview_layout_test.go` (they are untagged and stop compiling here).
- [ ] Run `wire gen ./internal/composition/...`; `wire_gen.go` is expected byte-identical (the provider set is unchanged). If it changes, record the diff and keep the regenerated file.

Tests:
- [ ] Delete the dedicated retired test folders (four ordinary services, three cluster services), `createBalancedChainZoneLabels_test.go`, `orderEdgeGap_test.go`, the deleted connection-name tests, and the `preferInterior` cases in `assignNeutralZonesToGaps_test.go`.
- [ ] `topologyServiceLookup/resolve_test.go`: each of the 7 survivors → its creator and true; each retired raw ID, an unknown and `""` → nil and false (named subtests). `tournament_test.go` uses a survivor.
- [ ] `topologyProvider/createTopologyVariant_test.go`: an unknown ID returns the wrapped error in ordinary **and** tournament mode, with neither creator called; a valid tournament calls only the tournament creator; each survivor dispatches to its creator (Random, Circles, Geometric Hub explicitly).
- [ ] `tournamentTopology/createTopologyVariant_test.go`: Random, Square, Geometric, Cross, Fractal, Geometric Hub and Circles all produce balanced output (`TBal-` prefix, left/right halves); 0 neutral zones and an odd neutral count do not panic; random portals unchanged.
- [ ] `templateGenerator`: remaining tests move to survivors; an unknown topology returns the wrapped error and a nil template; `generateTournament_test.go` covers a Random tournament → balanced.
- [ ] `templateHandler`/`guiHandler` generate: a generator error propagates through the mock.
- [ ] Zone labels: `getHoldCityLabel` with Geometric Hub; `createOrderedZoneLabels` without `isRing`; zone plans without the Shared Web rule; `isHubCityToHold` (Geometric Hub only); `assignNeutralZonesToGaps` without `preferInterior`.
- [ ] Every remaining service-level reference recorded in Phase 1 moves to a survivor that keeps the test's intent.

### Verification Plan
- `go build ./...`; `go vet` under no tag, `integration_test`, `integration_test,gui` (proves no retired symbol hides behind a tag).
- `go test ./test/... -count=1`; tagged integration; tagged GUI (K15 for any `.failure`).
- `go test -bench=. -run=xxx ./test/performance/... -benchtime=1x -timeout=120s` compiles and runs.
- Search `TopologyRing\b|TopologyHubAndSpoke|TopologyChain|TopologySharedWeb|RingTopologyService|HubTopologyService|ChainTopologyService|SharedWebTopologyService|RingClusterService|HubClusterService|ChainClusterService|CreateBalancedChainZoneLabels|OrderEdgeGap` over `*.go` → no hits. The raw legacy IDs appear only in `retiredTopologies.go` and test fixtures/tests.
- `git diff --stat -- internal/composition/wire_gen.go` recorded.

### Phase Summary
_(write when phase completes)_

## Phase 3: GUI coverage, benchmarks and documentation
Status: Not started

- [ ] GUI test (`integration_test && gui`, `test/integration/gui/`): the topology dropdown lists exactly the 7 K14 labels in order.
- [x] GUI test: with a modified current document, loading the v2 retired fixture and then the v0 retired fixture each show the exact K4 status `Load failed: topology "…" (saved as "…") has been retired and is no longer supported; re-create the template with a supported topology.`; the state, current path, unsaved flag and generated template are unchanged. **Delivered in Phase 1** as the GPU-free `topologyLoadRejection_integration_test.go`.
- [ ] Plain GUI run; for each `.failure`, confirm the batch explains it (tournament output for non-Circles topologies, moved fixtures or selections), accept only those files, and list each with its reason in the Phase Summary for owner approval (K15). Never a blanket `-update`.
- [ ] Benchmarks after: both commands from Phase 0 to `tmp/batch-k-bench-after.txt` (recorded, no gate).
- [ ] README topology table → 7 rows (Random, Circles, Square, Geometric, Geometric Hub, Cross, Fractal). Search README/QUICKSTART/docs for the retired names and any Ring fallback.
- [ ] `.agent/memories/generator-domain.md`: replace the retirement notes with the new contract.

### Verification Plan
- `go test -tags=integration_test,gui ./test/integration/gui/... -count=1` passes; no leftover `.failure` files.
- Both benchmark commands pass.
- The touched README section's links resolve.

### Phase Summary
_(write when phase completes)_

## Phase 4: Final verification, review and close-out
Status: Not started

- [ ] Full gate: build; vet × 3 tag sets; `go test ./test/... -count=1`; tagged integration; tagged GUI; both benchmark commands; testlayoutcheck; `gofmt -l` on an explicit list; lint report-only (0 issues; `--fix` only for formatter findings in this batch's files).
- [ ] Coverage after; compare per file with `tmp/batch-k-coverage-before.txt`: no surviving file drops; every new or changed function at 100%; total ≥ 74.4%, or the denominator change is explained for owner approval (K16).
- [ ] Independent implementation review (GPT-6.1 Sol); apply or justify every finding.
- [ ] Review document: §2.3 marked FIXED **pending owner commit**, the §9 row K and the progress line. Do not touch §8 or renumber.
- [ ] `.agent/backlog/test_observations.md`: record any new GUI-only statement uncovered by unit tests.
- [ ] Rewrite the handoff per AGENTS.md §5.2.

### Verification Plan
- Every command exits 0; the coverage comparison is recorded in the Phase Summary.

### Phase Summary
_(write when phase completes)_

## Final Recap
_(write when all phases complete: summary of the entire piece of work)_

## Deployment Plan
_(write when all phases complete: step-by-step deployment instructions)_
