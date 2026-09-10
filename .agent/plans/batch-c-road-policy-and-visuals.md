# Batch C: road policy, connectivity and visual road status

Implement review §1.4/§1.6/§1.10 with the owner's clarified between-zone road policy, plus road-state styling in the Preview panel, manual editor and exported PNG. **Scope and written-plan implementation approved on 2026-09-10**, owner: “looks good, please proceed”. Phase 1 is committed at `07ca8b6`; Phase 2 is committed at `abf0d36`. Phase 3 is next, explicitly authorized by the owner: “Changes reviewed, you can proceed”. Resume in a fresh session without another approval gate.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../../AGENTS.md), [current handoff](../session-carry-forward.md) and
[review](../backlog/review-gpt-6-astra-09-07.md). A/D and B are complete: no validation,
closure or retired-document lookup. Current safety checks and fresh baselines are
not reopening those batches. Do not stage, unstage, commit, push or switch branches.
Protected data/schema/registry paths remain read-only; generated Wire is regenerated,
never edited. No global test tags, fake test seams, output-path changes or persistence.

## Approved behavior contract

### Road policy and graph

- The GUI label is **Generate roads between zones**, not disable every road.
- With settings present, `GenerateRoads` authoritatively stamps `Connection.Road`
  true/on or false/off on **every non-explicit-Portal connection**, including
  Direct, Default, empty, GladiatorArena, Proximity and custom/imported connections.
  On overwrites imported false; off overwrites imported true. Off→on restores roads.
- Off removes the whole zone-road segment if either endpoint references one of
  these disabled connections. A mixed direct/portal segment is removed whole;
  rebuild eligible portal approaches without the disabled endpoint.
- Keep graph connections, guards, placement rules and arena markers. Repair uses
  actual graph endpoints, including accumulated repair edges, never road presence.
  Preserve the accepted impossible-isolation fallback edge.
- Only an explicit Portal type is exempt. Use the existing case-insensitive editor
  convention for that type; portal-placement rules alone do not grant exemption.
- Preserve explicit Portal road flags and existing valid approaches, even when a
  portal's own flag is false. The checkbox does not control portal approaches.
  Generated portals retain true. Do not reinterpret `Road=false` as authority to
  erase all portal approach segments; that was not approved.
- Internal castle/object/foothold roads remain generated and rebuilt regardless of
  the checkbox, in fresh generation, new zones, quality edits and castle reapply.
- Preserve valid custom internal roads; remove only confirmed invalid MainObject,
  incident Connection or named MandatoryContent references. Preserve unknown/opaque
  reference types. No arbitrary prefix-only deletion or provenance schema.
- Resolve mandatory item targets through the final groups referenced by each zone,
  not a global item-name lookup. Derive generated foothold additions from effective
  settings and final named content; verify that final main-object anchors exist.
  No MainObject:0 road when the zone has no main objects; retain existing connector
  behavior, do not invent a new foothold anchor policy for castleless zones.
- `TemplateUpdateDto.EditorState == nil`: preserve supplied zone roads and connection
  flags, and the template's existing mandatory content. Skip all road reconciliation,
  including dangling-target cleanup. Call `EnsureConnectionNames` directly (it is
  currently nested inside the rebuild being skipped) and retain normal graph checks;
  do not guess generation settings or add a nil-state rejection.
- Pending manual connection creation/type changes must reflect the same settings
  policy before Apply via a model/service/handler seam. Cancel remains non-mutating.
  No new per-connection road checkbox is requested.

### Visual contract

- The shared road classifier is: explicit false → roadless; explicit true → roaded;
  nil → roaded only when the declared type is Portal. Do not infer from zone roads.
  The owner will test these game semantics after implementation; source inspection
  did not establish the engine's omitted-value default.
- Preview panel AND manual editor: roadless non-portals `#B0B0B0`, roadless explicit
  portals `#90EE90`, opaque. Existing roaded colors and normal widths stay unchanged.
- Selected editor edges keep their road-state color, **including roaded edges**;
  use the existing thicker selected width without the orange color override.
- GUI legend entries: `Road`, `No road`, `Portal`, `Portal without road`. Expose the
  key where each canvas needs it without redesigning the panels/dialog.
- PNG: roadless strokes use the existing dark stroke color composited once at 50%
  per edge. Individual brush stamps of one edge must not compound opacity; distinct
  intersecting edges may compound. Keep PNG portal dashes, curve geometry, spacing,
  clipping, line widths and fully opaque arena/zone markers unchanged.
- Keep existing effective preview connection classification for shape/dash/markers.
  Road color classification uses explicit type separately. A Direct connection with
  portal placement rules therefore keeps its existing portal-shaped PNG rendering,
  but its nil/off road state is non-portal/roadless. This is intentional scope
  separation, not authorization to implement review §1.15 or §2.1.

## Verified starting point

- 2026-09-10: clean `master`, HEAD `5b1feb0eb315a4df54ca1ea9e1defd070f6061d8`.
  No staged/unstaged/untracked owner changes reported; do not assume that persists.
- Windows/amd64 Go 1.27.0; `go env GOFLAGS` empty.
- Fresh full unit coverage task (`-count=1`): PASS. Go's `cover -func` total **74.5%**.
  Existing ignored coverage.txt/coverage.html/lcov.info refreshed; Git still clean.
- Fresh report-only lint: **0 issues**, three existing unused-exclusion warnings.
- Relevant statement coverage: RebuildZoneConnectionRoads 97.1%; RebuildCastleRoads,
  ApplyNeutralZoneQuality and UpdateTemplate 100%; CreateContentsForZones 95.2%;
  CreateMissingPlayerConnections 92.6%; CreateMissingConnections 100%. These percentages
  do not cover the missing policy/value scenarios.
- No build, GUI, Linux, race or performance baseline execution claimed this session.
- Production manual Apply always supplies EditorState. Nil is an API-surface contract,
  not a missing GUI scenario to manufacture.
- Source verification: [road rebuild](../../internal/services/connection_editor/zoneEditorService.go),
  [UpdateTemplate](../../internal/handlers/templateHandler.go),
  [mandatory content](../../internal/services/template_generator/providers/mandatoryContentProvider.go),
  [road factory](../../internal/services/zones/roadFactory.go),
  [graph repair](../../internal/services/template_generator/providers/topology/base/topologyConnectionService.go).
- Both current previews retain edges with roads off: [preview projection](../../internal/services/preview_service/previewLayoutService.go)
  drops Road state, [GUI drawing](../../app/gui/utils/draw.go) and
  [PNG drawing](../../internal/services/preview_service/previewGeneratorService.go)
  use type only. [Manual editor](../../app/gui/dialogs/zoneEditorCanvas.go) also ignores Road.

## Design and implementation boundaries

1. Put shared explicit-portal/road-status logic on the connection model, with a
   small model enum in its own file if needed. No model→GUI/DTO dependency. Preview
   render data must carry road display state separately from existing effective Type.
  Existing preview shape classification is case-sensitive, editor shape classification
  is case-insensitive; leave both unchanged. The new road classifier is case-insensitive.
2. Introduce one cohesive road-policy/reconciliation service under internal/services/zones,
   with its interface in the existing zone_interfaces subpackage and any request
   model under internal/models. Reuse RoadFactory for construction. Keep generation,
   manual Apply and pending connection changes on this one policy rather than forks.
   Prefer existing files/seams except where a new cohesive responsibility is needed.
3. Compute final mandatory content in UpdateTemplate before reconciliation. Pass
   content data into the service; never inject MandatoryContentProvider into the
   zone editor (the provider already depends on it, creating a DI cycle).
4. Generation finalization must see the complete connections, content and arena type
   changes. Stamp final flags and reconcile approach/foothold targets after assembly.
   Do not broadly rebuild all main-object roads after arena placement: that would
   newly connect arena markers as castles. Keep ordinary internal-road construction
   in existing factory/castle-edit seams and preserve existing spawn anchoring.
5. Split RoadFactory's internal loops from its connection-road gate. Repair APIs
   receive `generateRoads` explicitly and guard only road materialization, not edges.
   Update TopologyBase/interface and positioned, Ring, Chain, SharedWeb and balanced
   tournament callers. These topologies remain supported until their separate removal
   batch; do not delete them here. Final reconciliation is not a substitute for correct
   standalone public repair behavior.
6. Do not indiscriminately replace valid custom MainObject→MainObject roads. Preserve
   valid records and add missing generated internal routes without duplicating them;
   remove genuinely stale anchors when castle counts change. Preserve custom road
   attributes, endpoint order and nil/reference ownership where not changed by policy.
7. PNG roadless drawing uses one reusable per-image stroke mask/scratch buffer, bounded
   clear between edges and a single masked Over composite per edge. Opaque drawing
   retains its current stamp positions and Src semantics. Use an opaque source with
   a half-opacity mask (8-bit 128/255 approximation), not unpremultiplied color.RGBA.
8. GUI decides how to paint, not what the road setting means. Pending edits invoke
   the existing handler boundary. Avoid per-frame deep clones/allocations solely for
   styling; resolve/classify on current model data or existing invalidation seams.
   No unrelated event-batching or graph-diagnostic cache refactor.

## Phase 1: Public policy and graph regressions

Status: Complete

- [x] Obtain owner approval of this written plan; inspect current Git safely. Approval received 2026-09-10. Clean `AD/road_and_graph_invariants` at `5240389aa0a8eb6a0a04d0552369bc2437540e60`; owner committed only the plan/handoff since the measured baseline. Recheck Git before code edits next session.
- [x] Add dedicated connection-model classifier tests, true/false/nil × explicit
  Portal/non-Portal (including placement-rule-only portal and type-case handling).
- [x] Add red public-API graph tests for roadless-but-connected players, newly
  accumulated repairs, required impossible-isolation fallback and component bridges.
- [x] Implement endpoint-based player connectivity and road-gated repair calls across
  TopologyBase, its service/interface and every current caller.
- [x] Test no unnecessary fallback for fixed A-neutral-B layout; required fallback
  for two-player/no-neutral isolated Random; bridge edges remain when roads are off.

### Verification Plan: Phase 1

- Run affected mirrored unit packages, including topology and generator suites.
- Run build, architecture tests and test-layout checker; verify protected/generated
  files unchanged except regenerated Wire if genuinely needed.

### Phase Summary: Phase 1

Started 2026-09-10. Current branch `AD/road_and_graph_invariants`, HEAD
`87ae8789d84178854a395e8dd34e172095a9ebc6`; working tree and index were clean at start.
Only plan/handoff documentation differs from the recorded source baseline.
Fresh pre-change full unit coverage task passed; Go total remains **74.5%**.

Completed 2026-09-10:

- Connection model now exposes `IsExplicitPortal()` / `HasRoad()` with the approved
  case-insensitive explicit-type classifier. No preview/GUI consumer changed yet.
- Player repair uses valid incident endpoints (not transitive reachability or roads),
  including newly created repairs; missing/self endpoints cannot mark a spawn connected.
- Both repair APIs take final `generateRoads bool`; TopologyBase and every current
  caller forward the setting. Edges/guards remain when approach-road creation is off.
- Occupied fallback/bridge names use the first free `-2`, `-3`, etc. suffix. This
  closes the existing bridge collision path that linked adjacency without creating
  an edge. Only real repair edges update connectivity; existing records are untouched.
- Added dedicated classifier and adjacent clone/mapping tests, direct service tests,
  base regressions, fixed positioned A-neutral-B and real Random isolation cases.
  Initial focused graph run reproduced eight failures; repaired suites pass.
- Independent Claude Opus 5 review approved after vacuous/combined test assertions
  were corrected and collision cases added. No approved policy was reopened.

Verification: production build PASS; affected generator/model suites and architecture
tests PASS; fresh post-change full unit coverage task (`-count=1`) PASS; final cached-
eligible full coverage rerun after test-only additions PASS, **74.6%** (baseline 74.5%).
Classifier methods, both service/base repair methods, fallback construction and unique-
name helper are **100%** statement-covered; value/branch matrices cover roads on/off,
endpoint validity, accumulated repairs and collisions. `go test ./test/...` PASS,
including untagged integration; test-layout checker PASS; report-only lint **0 issues**
(three existing exclusion warnings); `git diff --check` PASS.
No protected data/schema/registry or generated Wire changes. Index remains empty;
all Phase 1 work is unstaged/untracked. No constructor changes or Wire regeneration.
Windows only: no Linux, gated integration/GUI, race, or in-game validation claimed.
Phase 2 remains unstarted; full road-flag policy and internal/foothold reconciliation
are not implemented by this phase. Existing preview shape/type classification remains.

## Phase 2: Road policy and final content reconciliation

Status: Complete

- [x] Implement flag policy, internal-vs-approach factory separation and coherent
  final-zone reconciliation using scoped mandatory-content references.
- [x] Integrate generation and UpdateTemplate with correct finalization ordering.
- [x] Preserve nil EditorState roads/flags/content exactly; cover graph warnings and
  existing naming behavior separately from the skipped road work.
- [x] Cover roads on/off/toggle both directions for all non-Portal types; explicit
  Portal true/false/nil flags preserved, including mixed direct/portal approaches.
- [x] Cover foothold disable, count decrease/increase, no-main-object zones, custom
  valid targets with matching prefixes, unknown refs, confirmed invalid indices,
  missing/nonincident connections, and per-zone mandatory group resolution.
- [x] Cover valid custom internal road preservation/attributes, no duplicate generated
  roads, actual castle changes, unchanged arena routing and existing spawn anchoring.
- [x] Regenerate Wire for changed providers/constructors; update real test mocks.
- [x] Update [shared generated-template golden](../../test/test_helpers/defaultTemplate.json)
  and derived fixtures for explicit road flags, preserving full golden equality.
  Never change protected example templates or weaken assertions to hide the delta.
- [x] Re-specify the existing nil-state UpdateTemplate mock test that expects a road
  rebuild: require direct naming, no reconciliation, unchanged roads/flags/content.

### Verification Plan: Phase 2

- Extend existing mirrored RoadFactory, ZoneEditorService, ManualReapplyService,
  MandatoryContentProvider, TemplateHandler and TemplateGenerator tests, plus dedicated
  tests per new public method. Existing no-roads-anywhere assertions must be replaced
  by the approved internal/portal/border distinction, not merely deleted.
- Run a real generation→Apply→serialization integration matrix with in-memory JSON
  inspection (or existing authorized test repositories), asserting flags and every
  generated road target in final output. Do not export diagnostics to arbitrary folders.
- Run affected manual castle/revert integration regressions without reopening their
  completed batch lifecycle. Preserve compare-before-mutation and snapshot isolation.

### Phase Summary: Phase 2

Started 2026-09-10 at `07ca8b6` (`AD/road_and_graph_invariants`): the owner committed
Phase 1. The only inherited edit is the staged handoff addition, preserved untouched.
Fresh pre-change full unit coverage task passed, total **74.6%**. Implementation
approval remains in force; no policy questions reopened.

Completed 2026-09-10:

- Added `RoadPolicyService`, its interface, `RoadReconciliationRequest`, and private
  resolved scope. Non-explicit-Portal flags follow the setting; explicit Portal flags
  and eligible approaches survive unchanged, including false/nil flags.
- Factory internal roads are independent of the checkbox. Reconciliation removes
  whole segments with confirmed invalid/disabled endpoints, validates mandatory
  names against each zone's final groups (including empty final content), preserves
  opaque/custom references and attributes, and adds missing approaches/foothold roads.
- `Generate` reconciles after content and arena placement; `UpdateTemplate` names
  connections directly and rebuilds final content before policy only with settings.
  Nil settings preserve roads/flags/content and still compute graph diagnostics.
- Castle reconciliation preserves custom routes and clones before pruning/appending
  to prevent alias mutation. Arena markers are excluded from generated anchors;
  marker-only zones retain connector routing. Imported shifted anchors are rebased
  to actual non-arena indices; original spawn anchoring is retained.
- Legacy settings-free `RebuildZoneConnectionRoads` remains public and delegates to
  policy with explicitly unknown content and roads on (documented flag stamping).
  No production caller remains; removal is not part of this phase.
- Added dedicated public-policy tests and real generation→Apply→JSON regression
  matrices, plus nil-handler, roads-disabled, arena-only and slice-alias tests.
  Existing mandatory-content and manual-reapply suites pass unchanged. Wire was
  regenerated. Full golden equality remains; its only delta is `road: true` on
  `Rnd-A-B` in the shared test fixture.
- Off→on restores eligible approach targets, not byte-identical redundant connector
  segments from initial topology assembly. Existing valid custom records are not
  collapsed or rewritten to make the two forms identical.

Verification: Windows production build PASS; full fresh unit run PASS; fresh full
coverage task and final cache-eligible reruns PASS, **74.6% → 74.9%**. Policy service,
resolved scope, road factory and changed handler/editor entry points are 100%
statement-covered after the final injected-factory reference matrix. Default
`go test ./test/...`, tagged integration, and
`go test -tags='integration_test,gui' ./test/integration/...` PASS (GUI included).
Test-layout and diff-whitespace checks PASS; final report-only lint **0 issues**
(three existing exclusion warnings). Independent Claude Opus 5 review approved
after arena-only anchoring, shared-array mutation and a one-endpoint assertion
were corrected. Regression tests were observed failing with those fixes removed.

Protected trees unchanged; generated Wire changed only through regeneration. All
Phase 2 source/tests/plan edits remain unstaged or untracked, and the owner's staged
five-line handoff addition is untouched. No Linux, race, benchmark or in-game
validation claimed. GUI/PNG road-state visuals are still Phase 3/4, not implemented.
Resume from this plan rather than the older staged handoff; later scope recorded
in that handoff remains unchanged. No extra approval question is needed for Phase 3.

## Phase 3: Pending editor policy and GUI road styles

Status: Not started

- [ ] Thread the current checkbox through existing CreateZoneEditorConnection. Add
  a handler/service operation for type-change policy (currently direct GUI writeback),
  with interface, real mocks and mirrored public-method tests. Show pending output
  before Apply; cancel must not mutate retained data.
- [ ] Carry road display state in preview data; keep effective Type for shape unchanged.
- [ ] Add approved theme colors and four legend entries; apply to both GUI canvases.
- [ ] Selected editor edge uses thicker width and retains its road color for every state.
- [ ] Cover changes to road flags/types without stale preview/editor styles, including
  off→on regeneration and selection/unselection. Do not widen geometry/cache refactors.

### Verification Plan: Phase 3

- Model/service/handler and preview projection unit matrices, including explicit false
  portal, nil Direct, nil Portal and Direct with portal rules.
- Tagged GUI tests driving real pending create/type-change, Apply/cancel and both canvases.
- Focused image assertions/snapshots for four styles, selected edges and legends.
  Inspect only intentional changed/new goldens; do not bulk-update unrelated snapshots.

### Phase Summary: Phase 3

Not started. Owner reviewed and committed Phase 2 at `abf0d36` and authorized the
next phase on 2026-09-10: “Changes reviewed, you can proceed”. Git was clean before
the documentation-only handoff update. The preceding session reached its budget;
begin Phase 3 directly in a fresh session after inspecting current edits. Read
[the refreshed handoff](../session-carry-forward.md), preserve all existing work,
and use 74.9% as the last verified unit coverage reference before a fresh baseline.

## Phase 4: PNG per-edge opacity

Status: Not started

- [ ] Implement reusable per-edge mask compositing at 50%, keeping opaque raster path,
  ceiling sample counts, floating-point spacing, clipping and dash policy unchanged.
- [ ] Cover opaque/roadless direct and portal short/long/curved/clipped paths, dash gaps,
  repeated overlapping stamps, two intersecting roadless edges and opaque markers.
- [ ] Assert exact expected composite pixels against the actual background; absent
  strokes must fail. No environment-gated capture-only test or generated asset changes.
- [ ] Provide focused before/after visual evidence for owner review without claiming
  in-game validation. Prefer in-memory assertions and existing GUI/snapshot facilities.

### Verification Plan: Phase 4

- Run preview-layout and raster unit suites; opaque fixtures remain pixel-equivalent.
- Check memory/allocation effect of reusable scratch with focused measurements rather
  than an automated global GUI allocation threshold. No per-edge full-image allocations.

### Phase Summary: Phase 4

Pending phase completion.

## Phase 5: Final verification and owner handoff

Status: Not started

- [ ] Independent Claude Opus 5 implementation review; fix current-batch issues here.
- [ ] `go build ./...`; full unit suite and before/after coverage using the prescribed
  coverpkg set and Go's authoritative total. Fresh implementation run uses `-count=1`;
  avoid needless forced repetitions for subsequent unchanged runs.
- [ ] `go test ./test/...`; tagged integration; tagged GUI integration across affected
  editor flows; test-layout checker; report-only lint (target zero issues).
- [ ] Run tagged integration+GUI under test/integration as required by editor changes.
  No global tags; do not mix wireinject into builds/tests. Regenerate Wire normally.
- [ ] Inspect current Git/diff, preserve owner staging, verify no protected changes.
- [ ] Record platform limits honestly: Windows measurement is not Linux validation.
- [ ] Owner reviews the visual output and tests engine true/false/nil behavior. Record
  in-game outcome separately; do not claim engine correctness from JSON examples.
- [ ] Update surviving review/decision/handoff records; owner stages/commits. Mark
  review findings FIXED only following the owner-commit protocol. Leave plan deletion
  to the owner; no completed-plan links in the next-batch handoff.

### Verification Plan: Phase 5

- Coverage must not fall below comparable **74.5%**; inspect touched function branches
  and value contracts rather than relying solely on the percentage.
- Record test/build/lint/layout/GUI results, reviewed visual changes and engine caveat.
- Final `git status --short` and diff checks; no Git index/branch mutation by the agent.

### Phase Summary: Phase 5

Pending phase completion.

## Plan review

Pre-draft independent Claude Opus 5 design review completed. Useful findings included
the mandatory-provider DI cycle, generation-after-arena ordering, opaque-vs-effective
portal separation and alpha-mask compositing. Explicitly rejected suggestions that
contradicted owner decisions: cleanup during nil-state updates; orange selected roaded
edges; removal of portal approaches because their flag is false; relying solely on
end-of-generation cleanup for public repair APIs; and new stone roads to arena markers.
Written-plan Claude Opus 5 review completed: approve with corrections. Both blocking
gaps were incorporated: preserve whole-template golden equality while updating explicit
flags, and retain direct naming on the nil-state path with the old mock expectation
re-specified. Clarifications for the new type-change seam and unchanged case-sensitive
preview shape classifier were also incorporated. Owner implementation approval received
2026-09-10: “looks good, please proceed”. No further approval gate for the existing scope.

## Final Recap

Pending implementation and verification.

## Deployment Plan

Pending completion: owner review/commit/build and session-only authorized export;
no database/schema migration or persisted output path.
