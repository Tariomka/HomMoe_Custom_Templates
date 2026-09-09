# Repository review: GPT-6 Astra, 2026-09-07

**Reviewed revision:** `f4f4cf63f22e84040754a7231b4d1dc793070af1`, branch `master`, initially clean working tree. **Platform:** Windows/amd64. **Toolchain:** Go `go1.27.0`; installed golangci-lint `2.13.1`, built with Go `1.27.0`. Both modules declare Go `1.27.0`; the tools module pins golangci-lint `2.12.2`. **Lint baseline: 0 issues. Unit statement coverage: 74.4%.**

**Scope and authority.** The original read-only review followed [AGENTS.md](../../AGENTS.md) and [the review contract](../promt_templates/review-prompt.md). No implementation, tests, configuration, game data, or existing coverage reports were modified. A temporary public-API verification program was created, run, and deleted; diagnostic reports were written outside the repository. Owner-directed curation on 2026-09-08 updated this backlog, its carry-forward, and completed entries in the two source backlogs. Baselines below remain measurements from the reviewed revision, not new test runs.

**Supersession.** This document supersedes the review/disposition status of [owner_findings.md](owner_findings.md) and [test_observations.md](test_observations.md), without deleting either and without cancelling their owner-requested features. It supersedes the *current-state claims* of [session-carry-forward.md](../session-carry-forward.md), not its historical decisions. No earlier review document or general backlog survives in the current tree. The latest commit deleted historical review/backlog files; their absent contents are not reconstructed from memory or silently described as audited. Every item in the two surviving backlog documents is accounted for below. The owner explicitly selected **“Include verified owner items”** during this review.

**Severity:** 🔴 High: bug/correctness/user-visible · 🟠 Medium: architecture, performance, CI gaps · 🟡 Low: readability/hygiene · ⚪ Informational.

**Finding count:** **32 actionable items: 8 High, 17 Medium, 7 Low.** This includes owner-requested architecture/product work scoped on 2026-09-08, not just proved defects. Informational observations and prior-item dispositions are not included in that count. Findings are source-verified unless a runtime reproduction is explicitly recorded. Performance claims are reasoned, not benchmark measurements. In-game behavior was not tested.

**Fix-session protocol:** ask → plan → owner approves → implement and verify → owner commits → mark the stable item `✅ FIXED`. Owner-confirmed scope below is not permission to implement before the per-item plan is approved. Protected-directory changes remain owner-approved and owner-applied only. Each future change must follow the test layout, AAA, `t.Parallel()`, testify, coverage, build, and tag requirements in [AGENTS.md](../../AGENTS.md). Proposed new test paths below are deliberately not hyperlinks until the files exist.

## §0 Disposition of prior reviews

Owner-item identifiers O01–O20 identify the original prose blocks in [owner_findings.md](owner_findings.md); T01–T06 identify its original six nested “After adding Template model” bullets. They remain stable after completed entries are removed from that source backlog. Observation descriptions identify entries in [test_observations.md](test_observations.md), including completed entries retained here as evidence. These identifiers are local accounting labels, not historical review section numbers.

### Fixed ✅

| Prior item | Current evidence and limits |
| --- | --- |
| O01: Geometric Hub save/load loses preview positions | **New saves fixed at the persistence seam:** [manualZoneSave.go](../../internal/models/editor_state_model/manualZoneSave.go#L10-L45) now carries generator position/ring both ways; [layout dispatch](../../internal/services/preview_service/previewLayoutService.go#L77-L89) consumes them. The owner's old schema-v1 artifact did not contain those positions, so migration cannot recover them. This is not a claim that an old artifact's missing coordinates were repaired or that every preview issue is solved. |
| T01: move Template conversion to mapper | [TemplateMapper](../../internal/mappers/templateMapper.go#L15-L58) owns `ToModel`/`ToEntity`; [Template](../../internal/models/template_model/template.go) no longer owns those methods. |
| T04: generator/builders should use models, not RMG entities | [generator return and assembly](../../internal/services/template_generator/templateGenerator.go#L57-L91) operate on models; [architecture enforcement](../../test/unit/architecture/dependency/layering_test.go#L45-L69) retains only the documented file-service entity exception. Current architecture tests pass. |
| T06: replace two square-map dimensions with one model size | [Template](../../internal/models/template_model/template.go#L10-L20) has `MapSize`; [mapper](../../internal/mappers/templateMapper.go#L44-L58) assigns it to both entity dimensions. Protected output schema unchanged. |

### Invalidated/accepted ✖

“Accepted” means a deliberate convention, explicit owner decision, or a proposal requiring product decisions rather than a newly proved defect. It does not mean an unreproduced bug is fixed.

| Prior item | Disposition and evidence |
| --- | --- |
| O05: OCR utility | Owner tooling proposal, not a missing application requirement. No OCR implementation found; retain as optional owner work, not a numbered defect. |
| O09: zone-content presets | Deferred by the owner on 2026-09-08: implementation has not been investigated or planned. Not an actionable item in this review. Existing [zone controls](../../app/gui/panels/layoutPanelZones.go) remain unchanged. |
| T02: stamp Quality during template mapping | [ToZoneModel](../../internal/mappers/templateMapper.go#L263-L301) leaves editor-only quality unstamped. [GetZoneQuality](../../internal/handlers/zoneEditorHandler.go#L58-L64) explicitly falls back to inference for raw RMG zones. No wrong current result established; owner must decide eager vs lazy inference before adding mapper service dependencies. |
| T03: remove all converter bridges | Two plural zone converters remain unused; their intentional 0% coverage is recorded below. Singular converters still bridge editor-state models without introducing a mapper→model→mapper cycle. [zone.go](../../internal/models/template_model/template_variant_model/zone.go#L85-L162), [converters.go](../../internal/models/template_model/converters.go). Delete-or-migrate is owner work; do not write coverage-padding tests for dead functions. |
| T05: move templateRevision into Template | [state.go](../../app/gui/drivers/state.go#L35-L40) defines revision as replacement/cache view state. Current behavior is coherent; moving it into the domain model is an owner design choice, not a missing invalidation found here. |

Test-observation dispositions, including the individual documented gaps:

| Observation | Disposition / present evidence |
| --- | --- |
| Bootstrap: StartApplication, eventLoop, getAndConfigureWindow; DestroyEvent failure and process exit | Accepted process/Gio-bound behavior in [program.go](../../app/gui/program.go). Returning errors to main remains owner-deferred; no new `os.Exit` finding. |
| Button constructors and private button semantics | Accepted Gio composition in [buttonWidget.go](../../app/gui/widgets/buttonWidget.go); dedicated [buttonPositionLogger tests](../../test/unit/app/gui/utils/buttonPositionLogger) cover public utility behavior. |
| Slider row composition | Accepted UI-only [sliderRowWidget.go](../../app/gui/widgets/sliderRowWidget.go); formatter tests remain separately unit-testable. |
| File explorer confirmation, overwrite, selection, confirm enablement, folder creation, listing, hidden toggle, selection, descent, scroll | Accepted GUI territory with [confirmation suite](../../test/integration/gui/fileExplorerDialog_integration_test.go) and [listing suite](../../test/integration/gui/fileExplorerDialogListing_integration_test.go). The current unit percentages are not evidence those GUI interactions are absent. |
| Windows hidden-attribute end-to-end gap | Accepted platform fixture limitation. [hiddenAttribute_windows.go](../../internal/services/file_system/hiddenAttribute_windows.go) has unit coverage, not full end-to-end attribute manipulation. |
| Whitespace save-target rejection unreachable from sanitized Save To | Accepted distinction between direct dialog construction and production [SaveTo](../../app/gui/drivers/stateFiles.go#L32-L42). Do not expose test seams to recover an unreachable UI path. |
| Zone editor geometry extraction and rendering/click remainder | Accepted split: [geometry service](../../internal/services/connection_editor/zoneEditorGeometryService.go) is unit-tested; [zoneEditorCanvas.go](../../app/gui/dialogs/zoneEditorCanvas.go) is GUI territory. New gaps are §1.13–§1.15, not a demand to unit-test Gio rendering. |
| Read-only zone name | Accepted: [zone property rows](../../app/gui/dialogs/zoneEditorZoneProps.go) offer a label, not a rename editor. No missing typing test. |
| Hub-specific property-row coordinates in the harness | Accepted known harness limitation in [zone-editor test helpers](../../test/test_helpers/integration_common); unchanged row semantics are covered by neutral/player cases. Recalibrate only when adding a Hub-specific interaction test. |
| Layout/Preview panels' LoadFromState/SaveToState | Accepted integration coverage. [layoutPanel.go](../../app/gui/panels/layoutPanel.go) and [previewPanel.go](../../app/gui/panels/previewPanel.go) remain Gio composition; numeric/domain validation is separate. |
| State.GetOutputPathWidget | Accepted UI-only method in [state.go](../../app/gui/drivers/state.go#L122-L124). |
| State.NewUIState host-dependent fallback | The **coverage limitation** is accepted; the fallback **behavior** is carried to §1.2 because it violates the current hard output rule. These are different claims. |
| handleSaveState/handleLoadState callbacks and suggestDirectory | Current [stateFiles.go](../../app/gui/drivers/stateFiles.go#L79-L131) has callbacks plus `getWorkingDirectory`; the old `suggestDirectory` name is stale. Accepted integration-only paths; see §7.2 for stale note handling. |
| SaveTo records currentPath only on successful save | Fixed behavior still present at [stateFiles.go](../../app/gui/drivers/stateFiles.go#L79-L93), covered by [stateSaveTo integration](../../test/integration/stateSaveTo_integration_test.go). No regression claimed. |
| PickOutputDir/RevealOutputDir | Accepted dialog wiring in [stateFiles.go](../../app/gui/drivers/stateFiles.go#L60-L76); the single-session override remains deliberately unpersisted. |
| getWorkingDirectory currentPath branch | Accepted integration coverage and private-path limitation; [stateFiles.go](../../app/gui/drivers/stateFiles.go#L125-L131). |
| reapplyManualEdits castle-change branch | Existing [manualCastleReapply integration](../../test/integration/manualCastleReapply_integration_test.go) remains relevant. §1.5/§1.10 identify different unasserted invariants. |
| integration_common GPU helpers vs pure snapshot helpers | Accepted tag separation. Current pure helpers live under [snapshot](../../test/test_helpers/integration_common/snapshot), not the older flat paths in the observation. |
| Private topology connection policy | Accepted coverage through public [TopologyBase](../../internal/services/template_generator/providers/topology/base/topologyBase.go); do not add test-only public services. New roads-policy gap is §1.6. |
| connectInteriorStables empty-input defense | Accepted valid-generation unreachable guard in [geometricHubLayout.go](../../internal/services/template_generator/providers/topology/geometricHubLayout.go). |
| preview provider asset-decode failure | Accepted embedded-asset failure seam in [previewGeneratorProvider.go](../../internal/composition/previewGeneratorProvider.go); no production seam solely for coverage. |
| Atomic writer private-file coverage, Close/Sync failures, truncation | Accepted through repository APIs; [atomicFileWriter.go](../../internal/repositories/atomicFileWriter.go) is 95.0%. Existing preservation-on-encoding-failure tests remain. Do not force disk-full or invent a close-failure test seam just for the number. |
| Steam/VDF/registry discovery chain | Accepted host-dependent gap in [io.go](../../internal/helpers/io.go) and [io_windows.go](../../internal/helpers/io_windows.go). Improving detector architecture can justify fixtures, but is not required merely to reach a coverage target. |
| buildShiftDerangement randomized fallback | Accepted defensive fallback in [topologyConnectionService.go](../../internal/services/template_generator/providers/topology/base/topologyConnectionService.go); no global random-seed seam proposed. |
| Test-local dto/stateDto names | Accepted cosmetic naming decision; no mechanical rename campaign. |
| No automated Gio allocation threshold | Accepted owner rejection of a flaky threshold. The statement that GUI tests never run in CI is now stale: [GUI job](../../.github/workflows/pr-validation.yml#L235-L278) uses Mesa/Xvfb. It still does not enforce the TabCycling allocation budget. |
| Dead ToZoneModels/ToZoneEntities | Accepted deliberate 0% coverage pending removal or a real caller; [zone.go](../../internal/models/template_model/template_variant_model/zone.go#L156-L162). Do not confuse this with the live singular converters. |

### Carried forward ❗

| Prior item | Re-verified outcome |
| --- | --- |
| O02: portals absent from PNG | **§1.3**, reproduced. Merely finding a dashed-line branch and a passing “different from solid” test did not establish that portal pixels were drawn. |
| O03: bonuses/bans ineffective in game | **§10 owner validation queue.** [gameRulesProvider.go](../../internal/services/template_generator/providers/gameRulesProvider.go#L21-L71) emits the data. Repository tests cannot prove the game's acceptance of the SID vocabulary. Not marked fixed or a newly proved application bug. |
| O04: added parallel/arched connection missing in PNG | **§2.1** establishes differing curve implementations; **§1.3** explains sufficiently short dashed chords. Current preview code retains multiple endpoint-pair edges. Exact old solid-edge artifact remains unreproduced; retain game/artifact follow-up in §10, not “fixed.” |
| O06: tournament player count and topology redesign | **§1.12** for the verified two-player enforcement gap. The specific retirement and fallback policy in **§2.3** is now scoped under O08; introducing new tournament designs remains separate owner-gated product work. |
| O07: hero-hire ban in Single Hero/FinalBattle/LostStartHero | **§10 owner validation queue.** Source emits `HeroHireBan` only for single-hero mode; broader desired game policy must be confirmed. Not silently discarded. |
| O08: remove Ring/Hub/Chain/Shared Web | **§2.3.** Remove ordinary implementations and corresponding tournament builders; reject retired saved selections and use balanced tournament generation for surviving fallback cases. |
| O10: quality change must recalculate connection guard value/preset | **§1.11**, source-verified missing propagation. |
| O11: save only manual deltas and remove derivable state | **§2.4.** Investigate both live state and persistence first; owner reviews feasibility before selecting reconstruction semantics or authorizing implementation. |
| O12: bans/overrides as lists | **§2.5.** Structured entries throughout, following `BonusEntry`, with existing parsing/export semantics retained and a coordinated migration. |
| O13: sectioned rather than flat JSON | **§2.6.** Nest current groups except `TemplateIdentity`, `MapSettings`, and `SchemaOptions`, whose fields remain flat. Coordinate the schema migration with §2.5 and §2.4's outcome. |
| O14: put LayoutPanel setters/getters on state | **§2.7.** Move the six zone-content accessors to `app/gui/drivers.State`; retain validated dirty-tracked updates and snapshot isolation. |
| O15: zone-content service returns DTO | **§2.2. Owner explicitly reopened the DTO removal request during this review**, after being shown the conflicting accepted-exception comment in the current architecture gate. This is newly authorized architecture backlog, not an undisclosed violation of the old gate. The `bonuses` exception remains accepted. |
| O16: rename zones/zone_interfaces packages | **§2.8.** Rename the implementation package and its nested interface package, including imports, mirrored tests, and generated wiring. |
| O17: intermittent preview artifacts | **§1.13** reproduces one real nondeterminism; **§2.1** records graph-renderer divergence. Neither is asserted to explain the historical artifact without a current reproducer. Keep that artifact queued in §10. |
| O18: move panel state into section subpackages | **§2.9.** Split General/Layout/Bonuses into private cohesive sections with top-level public aliases; preserve behavior and leave Preview unchanged. |
| O19: move neutralRowsForQuality | **§2.10.** Promote the lookup to a `GeneratorConfig` method without changing quality mapping or ownership semantics. |
| O20: use Vec2 methods instead of handwritten vector arithmetic | **§4.1.** Audit permitted production code, replace equivalent arithmetic, and add missing operations only when justified and tested. |

**Memory/history invalidations applied:** current branch is `master`, not the branch named in the handoff; schema-v2 work is present in the merged revision; missing backlog/plan paths are not active instructions to push/delete anything; old GUI-CI absence claims are stale. The current architecture gate accepts two DTO exceptions, but the owner explicitly reopened the zone-content one in this review; the bonuses exception is unchanged. The frozen v1 boundary, plain typed repositories, output-path non-persistence, nil/ring-zero semantics, and per-frame clone tradeoff were respected.

## §1 Bugs and correctness

### 1.1 🔴 Applying manual edits does not mark the document unsaved

**Progress (2026-09-09).** Implemented and verified on Windows Go 1.27.0; changed
manual commits now set dirty and re-arm Exit, while identical, rejected, and
accepted-warning no-op Apply cases preserve the appropriate state. Awaiting owner
commit; this finding is not yet marked fixed.

**Evidence.** [ApplyEditedZones](../../app/gui/drivers/stateManualEdits.go#L22-L35) ends with `this.innerState.SetManualEdits(request.Zones, request.Connections)` or `ClearManualEdits()`, without changing `unsaved` or `confirmExit`. The only dirty setter is [UpdateState](../../app/gui/drivers/state.go#L135-L142): `if this.innerState.WasStateChanged() { this.unsaved = true }`. That comparison deliberately [ignores manual edits](../../internal/models/editor_state_model/editorState.go#L129-L157). [Exit](../../app/gui/drivers/stateFiles.go#L48-L57) asks for confirmation only when `this.unsaved && !this.confirmExit`.

**Why wrong.** Load/save a clean state, move a zone or change a connection, Apply, then Exit. The persisted layout changed but `IsUnsaved()` remains false; the application can close without its normal warning. Revert-to-Base also changes persisted manual state and needs dirty tracking. Later per-frame scalar comparison does not repair this omission.

**Fix.** At the manual-edit commit seam, compare the persisted manual snapshot before/after; mark the document dirty and reset `confirmExit` when a successful apply changes it. Keep manual fields excluded from automatic regeneration comparison. Preserve clean state on rejected/no-template applies; decide whether semantically identical Apply should remain clean. Extend [applyEditedZones_test.go](../../test/unit/app/gui/drivers/stateManualEdits/applyEditedZones_test.go) with dirty-after-edit, dirty-after-revert, unchanged-apply, and rejection cases. Extend [exit_test.go](../../test/unit/app/gui/drivers/stateFiles/exit_test.go) to prove the first Exit after a manual edit does not invoke `onExit` and that a later edit invalidates a previous confirmation.

**Owner decision.** Confirm no-op Apply semantics; do not change close/reset/load policy incidentally.

### 1.2 🔴 Failed game-directory detection silently authorizes export to the working directory

**Progress (2026-09-09).** Implemented and verified on Windows Go 1.27.0:
`FindGameTemplateDirectory` now belongs to `PathResolutionService`/
`IPathResolutionService`, and failed detection leaves export unset until explicit
session-only picker confirmation. Awaiting owner commit; this finding is not yet
marked fixed.

**Evidence.** [NewUIState](../../app/gui/drivers/state.go#L77-L89): `templateDir = state.getWorkingDirectory()` after failed detection, then `state.outputPath.SetText(templateDir)`. One status even says `"using fallback directory."` with the error flag false. [SaveTemplate](../../internal/handlers/templateHandler.go#L107-L119) rejects only an empty output path before writing through the file service.

**Why wrong.** An undetected game installation produces a nonempty but unrelated export destination. The game never reads the resulting template. This violates the explicit machine-detected output-path requirement, independent of the accepted host-dependent coverage gap.

**Fix.** Leave the output path empty on discovery failure, report actionable failure, and require the existing session-only folder picker. Preserve automatic detection on each launch and never persist the selected directory. Extend [newUIState_test.go](../../test/unit/app/gui/drivers/state/newUIState_test.go) and [saveTemplate_test.go](../../test/unit/internal/handlers/templateHandler/saveTemplate_test.go) for success, failure, and export refusal. If deterministic constructor coverage needs a detector collaborator, introduce it only as a real composition responsibility, not a private test export; regenerate Wire. Add a GUI picker recovery case in the existing file-explorer suite.

**Owner decision.** This enforces the existing hard rule; no permission to choose a new default export folder or remember a path. Verify that opening the session folder picker from an empty output path still resolves a usable browsing start directory; a browsing start is not an authorized export destination.

### 1.3 🔴 Short portal connections rasterize to zero pixels in PNG previews

**Evidence.** [drawLine](../../internal/services/preview_service/previewGeneratorService.go#L131-L147): `steps := max(math.Abs(delta.X), math.Abs(delta.Y))`, followed by `for i := range int(steps)`. [drawDashedLine](../../internal/services/preview_service/previewGeneratorService.go#L113-L129) splits curves into 96 chords; solid curves use 24. Positive subpixel chords pass `steps <= 0` but truncate to zero loop iterations.

**Reproduction.** A public `IPreviewLayoutService` stub returned endpoints `(300,350)` and `(400,350)`, midpoint `(350,350)`, radius 21, and one portal. `CreatePreviewImage` returned pixels **identical to the no-connection background**. Changing only the type to Direct produced different pixels. No game installation or GPU was needed. The existing [portal test](../../test/unit/internal/services/preview_service/previewGeneratorService/createPreviewImage_test.go#L68-L84) only asserts “different from solid,” which a completely absent line satisfies.

**Why wrong.** The GUI uses vector paths while PNG export uses this rasterizer, explaining why an on-screen connection can disappear only in the export. Short solid segments are susceptible too.

**Fix.** Use an integer ceiling-based sample count and include a sample for every nonzero segment; ensure shared endpoints do not introduce visible dash artifacts. Do not merely lower `segmentsDashed`. Through `CreatePreviewImage`, test short/long horizontal, vertical, and diagonal portals and solid edges against a no-edge baseline, including visible connector-colored pixels outside zone artwork. Extend [createPreviewImage_test.go](../../test/unit/internal/services/preview_service/previewGeneratorService/createPreviewImage_test.go). Keep fixtures in tests, not the protected asset tree.

**Owner decision.** Inspect changed PNG expectations. Do not automatically regenerate unrelated Gio goldens: this rasterizer is not their vector drawing path.

### 1.4 🔴 Manual Apply and castle reapplication re-enable disabled roads

**Evidence.** [RebuildZoneConnectionRoads](../../internal/services/connection_editor/zoneEditorService.go#L117-L135) calls `CreateOuterZoneRoads(nil, mainObjectCount, 0, true)` and `CreateConnectorZoneRoads(names, true)` and directly appends main-object connection roads. [RebuildCastleRoads](../../internal/services/connection_editor/zoneEditorService.go#L279-L291) also passes `true`. [UpdateTemplate](../../internal/handlers/templateHandler.go#L79-L90) calls the rebuild unconditionally although its DTO carries the editor settings.

**Reproduction.** Two roadless spawn zones with one connection gained one road each when passed to the real `RebuildZoneConnectionRoads`. In production, uncheck Generate roads, generate, open the editor and Apply, even without substantive edits. Castle changes and neutral-quality reprofile also reach the unconditional road rebuild.

**Why wrong.** The exported template contradicts its road-generation option and can keep doing so on every subsequent manual-snapshot reapply.

**Fix.** Thread `GenerateRoads` through road rebuild, quality-edit, and castle-reapply service contracts. Guard both factory and hand-built road paths. Agree how to treat preserved user/imported roads when disabled; ensure no stale generated road survives. Source the setting from `TemplateUpdateDto.EditorState`, current generator configuration, and the dialog's existing options. Extend [rebuildZoneConnectionRoads_test.go](../../test/unit/internal/services/connection_editor/zoneEditorService/rebuildZoneConnectionRoads_test.go), [rebuildCastleRoads_test.go](../../test/unit/internal/services/connection_editor/zoneEditorService/rebuildCastleRoads_test.go), [applyNeutralZoneQuality_test.go](../../test/unit/internal/services/connection_editor/zoneEditorService/applyNeutralZoneQuality_test.go), and [updateTemplate_test.go](../../test/unit/internal/handlers/templateHandler/updateTemplate_test.go). Cover false, true, toggle-off, and castle-change reapply.

**Owner decision.** Define `UpdateTemplate` behavior when its optional `EditorState` is nil; do not guess a default that destroys imported roads.

### 1.5 🔴 Tournament/arena changes can be overwritten by stale manual snapshots

**Evidence.** [LayoutDefiningOptionsChanged](../../internal/models/editor_state_model/editorState.go#L98-L106) compares players, topology, roads, portals and zone counts, but not `Tournament`, `VictoryCondition`, or `GladiatorArena`. [TopologyProvider](../../internal/services/template_generator/providers/topologyProvider.go#L24-L30) switches the whole graph for tournament mode. [Generate](../../internal/services/template_generator/templateGenerator.go#L88-L91) places the arena before returning; [UpdateTemplate](../../internal/handlers/templateHandler.go#L83-L90) later replaces the variant zones/connections without reapplying that placement. [DecideManualEditReapplication](../../internal/services/editor/regenerationDecisionService.go#L72-L91) relies on the incomplete predicate.

**Why wrong.** With an existing manual snapshot, enabling a two-player tournament can discard the newly generated tournament graph and restore the old one while retaining tournament rules. Enabling/disabling the arena can similarly retain stale markers or omit the newly placed marker.

**Fix.** Compare the *effective layout modes*, including victory-condition aliases, not just checkbox bits. For a true tournament graph switch, invalidate the incompatible manual snapshot. For arena changes, choose explicitly between clearing the snapshot and recomputing arena placement on the final edited variant; the latter must also remove obsolete markers. Extend [layoutDefiningOptionsChanged_test.go](../../test/unit/internal/models/editor_state_model/editorState/layoutDefiningOptionsChanged_test.go), [decideManualEditReapplication_test.go](../../test/unit/internal/services/editor/regenerationDecisionService/decideManualEditReapplication_test.go), and [updateTemplate_test.go](../../test/unit/internal/handlers/templateHandler/updateTemplate_test.go). Cover enable/disable and checkbox/victory-selector equivalents with manual edits present.

**Owner decision.** Clearing edits or relocating an arena is user-visible. Confirm the policy before implementation; do not drop edits on an unrelated victory-condition change unnecessarily.

### 1.6 🔴 Connectivity repair both ignores road policy and uses roads as graph connectivity

**Evidence.** [spawnZoneHasConnection](../../internal/services/template_generator/providers/topology/base/topologyConnectionService.go#L255-L264) scans `zone.Roads`, not connection endpoints. [CreateMissingPlayerConnections](../../internal/services/template_generator/providers/topology/base/topologyConnectionService.go#L87-L125) creates fallback connections when that predicate is false. [appendSpawnFallbackRoads and appendBridgeRoads](../../internal/services/template_generator/providers/topology/base/topologyConnectionService.go#L266-L325) append roads with no `GenerateRoads` condition. [positioned builder](../../internal/services/template_generator/providers/topology/positionedTopologyBuilder.go#L55-L69) calls these repair paths after constructing the ordinary graph.

**Why wrong.** A valid roadless spawn is treated as disconnected even when actual edges exist. For example, a positioned/Random isolated-start map with two players and neutral zones can already connect both players through neutrals; disabling roads makes both spawn road lists empty, so the repair adds an unnecessary direct player link and defeats isolation. Independently, Random with two players, zero neutrals, isolated starts and roads disabled needs the accepted fallback edge but still must not receive roads. Existing [generator road test](../../test/unit/internal/services/template_generator/templateGenerator/generate_test.go#L414-L432) exercises a non-isolated Ring case, not repair.

**Fix.** Determine graph connectivity from connections (including newly accumulated repairs), independently of roads. Thread `GenerateRoads` through public TopologyBase repair calls and guard road construction only, not connectivity repair. Retain the documented safety fallback where isolation is impossible. Extend [createMissingPlayerConnections_test.go](../../test/unit/internal/services/template_generator/providers/topology/base/topologyBase/createMissingPlayerConnections_test.go), [createMissingConnections_test.go](../../test/unit/internal/services/template_generator/providers/topology/base/topologyBase/createMissingConnections_test.go), and [generate_test.go](../../test/unit/internal/services/template_generator/templateGenerator/generate_test.go). Assert both no unnecessary new edge for already-connected roadless players and zero roads after required repairs.

**Owner decision.** Preserve the existing “connect an otherwise impossible isolated map” policy; do not confuse that accepted fallback with the incorrect road-based predicate.

### 1.7 🔴 Tournament “Save army” is ignored

**Evidence.** [GeneralPanel.SaveToState](../../app/gui/panels/generalPanel.go#L179-L184) saves `settings.TournamentSaveArmy = this.checkTournamentSaveArmy.Value`. [setTournamentRules](../../internal/services/template_generator/providers/gameRulesProvider.go#L133-L143) writes `winConditions.TournamentSaveArmy = true` instead of consuming the configured value.

**Reproduction.** Calling the real `CreateGameRules` with tournament enabled and `TournamentRules.SaveArmy=false` produced `TournamentSaveArmy=true`.

**Why wrong.** The checkbox survives editor persistence but never controls the generated rules. Statement coverage is 99.3% for this provider, yet this value-level contract is untested.

**Fix.** Assign the configured value. Replace the existing false-input/true-output expectation with independent true/false cases in [createGameRules_test.go](../../test/unit/internal/services/template_generator/providers/gameRulesProvider/createGameRules_test.go). Add serialized-output verification at the file-service/generator integration seam. The protected entity uses `omitempty`; confirm that omitted false represents “off” in the game before claiming full in-game correctness.

**Owner decision.** If the game requires explicit false and the existing protected schema cannot express it, stop and request owner approval for a schema change; never edit it as part of a routine fix.

### 1.8 🔴 Untouched Revert-to-Base is compared after the handler mutates it

**Progress (2026-09-09).** Implemented and verified on Windows Go 1.27.0:
revert identity is captured before handler mutation, and the real-handler
roads-off regression confirms an untouched revert clears the snapshot. Awaiting
owner commit; this finding is not yet marked fixed.

**Evidence.** [ApplyEditedZones](../../app/gui/drivers/stateManualEdits.go#L22-L35) calls `this.handleUpdateTemplate(...)` before `matchesZoneSet(request, pendingBase)`. The comparison is [reflect.DeepEqual](../../app/gui/drivers/stateManualEdits.go#L66-L69). [UpdateTemplate](../../internal/handlers/templateHandler.go#L79-L85) rebuilds the request's roads in place. The dialog [copies the top-level zone slice](../../app/gui/dialogs/zoneEditorDialog.go#L179-L197), separating those field replacements from the stored pending base.

**Why wrong.** An untouched revert should clear the manual snapshot. Road reconstruction can make its payload differ from the pre-rebuild base, particularly for roadless generation, so the equality fails and the base is pinned as another manual edit. This is a source-traced interaction; no full GUI reproduction was run. The current [revert unit test](../../test/unit/app/gui/drivers/stateManualEdits/applyEditedZones_test.go#L135-L160) uses a nonmutating mock and passes identical base slices directly, missing the dialog/handler interaction.

**Fix.** Compute the untouched-revert predicate against the unmodified payload before calling the updating handler; commit the clear only after an accepted update. Coordinate with §1.1 dirty tracking. Extend [applyEditedZones_test.go](../../test/unit/app/gui/drivers/stateManualEdits/applyEditedZones_test.go) with separate top-level slices and a mock that mutates request roads; extend [zoneEditorRevertToBase_integration_test.go](../../test/integration/zoneEditorRevertToBase_integration_test.go) through the real handler. Test untouched, edited, and rejected applies.

**Owner decision.** None beyond preserving the existing commit-on-Apply/revert design.

### 1.9 🟠 Manual-zone snapshot setters and getters share nested storage

**Progress (2026-09-09).** Implemented and verified on Windows Go 1.27.0: manual
zones now deep-clone on ingress and egress, including zone-only snapshot
isolation checks. The pre-existing connection `Road` pointer/placement alias is
unchanged by this scoped work. Awaiting owner commit; this finding is not yet
marked fixed.

**Evidence.** [EditorState.SetManualEdits/GetManualZones](../../app/gui/models/editorState.go#L96-L113) both use `slices.Clone`, while the model's [Clone](../../internal/models/editor_state_model/editorState.go#L77-L91) correctly calls `template_model.Zone.Clone` for every zone.

**Reproduction.** Set a zone with `ManualPosition=(0.2,0.3)`. Changing the input pointer's X to 0.9 changed the stored snapshot. Changing Y on the result of `GetManualZones` to 0.8 also changed stored state. Both occurred through public APIs. Arrays such as roads and nested main objects have the same shallow-copy boundary.

**Why wrong.** Callers can change supposedly retained manual state without a commit, validation, dirty tracking, or revision update. `UpdateTemplate` also holds the input zones as its live variant, so its data can share storage with the snapshot. This is once-per-apply ownership, not the accepted per-frame deep-clone tradeoff.

**Fix.** Deep-clone zones on both ingress and egress using the existing `Zone.Clone`; preserve nil semantics. Extend [setManualEdits_test.go](../../test/unit/app/gui/models/editorState/setManualEdits_test.go) and [getManualZones_test.go](../../test/unit/app/gui/models/editorState/getManualZones_test.go) with input/output mutation isolation for pointers, slices, nested reference fields, and empty/nil values. Reuse existing clone drift guards rather than testing only a single scalar.

**Owner decision.** None; retain intentionally mutable local editing copies, not aliases into stored snapshots.

### 1.10 🟠 Foothold changes leave roads targeting removed mandatory content

**Evidence.** [road preservation](../../internal/services/connection_editor/zoneEditorService.go#L102-L121) keeps every road that is neither a connection road nor a castle road. [RoadFactory](../../internal/services/zones/roadFactory.go#L69-L77) creates foothold roads targeting `name_remote_foothold_N`. [CreateContentsForZones](../../internal/services/template_generator/providers/mandatoryContentProvider.go#L79-L85) recalculates foothold content from the current enabled/count settings, independently of those preserved roads. [UpdateTemplate](../../internal/handlers/templateHandler.go#L81-L90) combines those two results.

**Why wrong.** With manual edits present, disabling remote footholds or reducing their count removes content definitions but retains old road references. Export contains a reference to a mandatory object that no longer exists. Increasing counts can leave newly added footholds without matching roads.

**Fix.** Rebuild generated foothold roads from final zone main objects and effective foothold configuration, or validate preserved targets against final mandatory content and add missing required roads. Do not delete arbitrary user-defined mandatory-content roads by name prefix alone. Extend [rebuildZoneConnectionRoads_test.go](../../test/unit/internal/services/connection_editor/zoneEditorService/rebuildZoneConnectionRoads_test.go) and [updateTemplate_test.go](../../test/unit/internal/handlers/templateHandler/updateTemplate_test.go) for enabled→disabled, count decrease/increase, no-castle zones, and roads-off. Verify every generated road target resolves in final export.

**Owner decision.** Distinguish generated foothold policy from imported/custom roads. Coordinate with §1.4 rather than adding a second inconsistent filter.

### 1.11 🟠 Neutral-quality edits cannot update incident connection guard presets

**Evidence.** [ApplyZoneEditorQuality](../../internal/handlers/zoneEditorHandler.go#L76-L83) receives/returns one zone and calls `ApplyNeutralZoneQuality` without any connections. [service](../../internal/services/connection_editor/zoneEditorService.go#L193-L226) changes zone pools, values, castles, and roads, not incident connection guards. The owner explicitly requested preservation of the selected preset tier while recomputing its value.

**Why wrong.** Raising or lowering a neutral tier leaves adjacent guard values derived from the old quality. The existing default connection factory uses [quality-derived guard strength](../../internal/services/connection_editor/connectionEditorService.go#L25-L40), so old and newly created connections around the same zone can disagree.

**Fix.** Have the service/handler seam operate on the affected zone plus incident connections, using both endpoints' effective quality and the existing guard-quality service. Preserve preset identity before reprofile and recalculate its value afterwards. If a numeric custom value cannot be mapped unambiguously to a preset, require explicit policy instead of guessing. Extend [applyNeutralZoneQuality_test.go](../../test/unit/internal/services/connection_editor/zoneEditorService/applyNeutralZoneQuality_test.go), [applyZoneEditorQuality_test.go](../../test/unit/internal/handlers/zoneEditorHandler/applyZoneEditorQuality_test.go), and [zoneEditorProperties_integration_test.go](../../test/integration/gui/zoneEditorProperties_integration_test.go). Test each preset, stronger opposite endpoint, custom values, and unrelated edges unchanged.

**Owner decision.** Confirm custom guard-value and preset identity behavior; do not add fields to protected connections.

### 1.12 🟠 Tournament remains selectable with more than two players

**Evidence.** [Players control/writeback](../../app/gui/panels/generalPanel.go#L150-L195) always permits 2–8. [TopologyProvider](../../internal/services/template_generator/providers/topologyProvider.go#L24-L30) only uses tournament topology when `configuration.IsTournamentMode() && len(playerLabels) == 2`; [game rules](../../internal/services/template_generator/providers/gameRulesProvider.go#L114-L117) enable tournament independently of player count.

**Why wrong.** The UI permits a configuration the dedicated topology path refuses, yielding ordinary topology with tournament rules. This contradicts the owner's stated two-player requirement.

**Fix.** Enforce the effective tournament-mode constraint in domain validation and render the player control consistently disabled/fixed at two. Handle loaded states and the victory-condition selector as well as the checkbox. Extend [validateEditorState_test.go](../../test/unit/internal/validators/editorStateValidator/validateEditorState_test.go), [createTopologyVariant_test.go](../../test/unit/internal/services/template_generator/providers/topologyProvider/createTopologyVariant_test.go), and the General-tab GUI integration suite. Test loaded 3–8 player states, switching modes, warning text, and non-tournament preservation.

**Owner decision.** Confirm whether to remember the previous non-tournament count in session view state. Topology retirement is separately scoped in §2.3; new tournament designs from O06 are not part of this constraint fix.

### 1.13 🟠 Symmetric obstacles make connection curves nondeterministic

**Evidence.** [obstacleBulge](../../internal/services/connection_editor/zoneEditorGeometryService.go#L237-L260) ranges a position map and selects a winner only when `math.Abs(signed) > bestMagnitude`. Equal-magnitude obstacles on opposite sides therefore choose whichever map entry is visited first.

**Reproduction.** Through `BuildGeometry`, endpoints `(100,350)`/`(600,350)`, obstacles `(350,340)`/`(350,360)`, radius 21, and 200 identical calls produced **two control-point Y values, 300 and 400** (18 and 182 calls respectively). The counts are not a stable distribution; the two different geometries for identical input are the failure.

**Why wrong.** Dirty rebuilds of unchanged/symmetric geometry can flip the curve and its hit target. This is a reproducible candidate for preview/editor artifacts, not proof of the owner's historical artifact's exact cause.

**Fix.** Add a deterministic tie-break, preferably without sorting/allocation on every edge (e.g. consistent sign/zone identity). Extend [buildGeometry_test.go](../../test/unit/internal/services/connection_editor/zoneEditorGeometryService/buildGeometry_test.go) with symmetric opposing obstacles, permuted map construction, and a repeated result assertion. Keep existing clearance and reverse-endpoint tests.

**Owner decision.** Choose the preferred tie direction; do not claim this solves all overlapping-obstacle routing.

### 1.14 🟠 Batched pointer events can resolve stale edge indices after deletion

**Evidence.** [handlePointer](../../app/gui/dialogs/zoneEditorCanvas.go#L81-L109) drains all events before the [geometry rebuild](../../app/gui/dialogs/zoneEditorCanvas.go#L56-L67). [edgeConnection](../../app/gui/dialogs/zoneEditorCanvas.go#L185-L192) indexes the current working slice by a cached `ConnectionIndex`. [deleteConnection](../../app/gui/dialogs/zoneEditorDialog.go#L359-L371) compacts that slice and only marks geometry dirty.

**Why wrong.** Two queued presses in a frame can use the first deleted edge's old index against a shifted slice. The second hit can select/delete another connection; a bounds check only prevents out-of-range access, not identity mismatch. This is a concrete static event-batching path, not reproduced with a real mouse during this review.

**Fix.** Keep an immutable pointer/identity mapping corresponding to the cached geometry until events referencing that frame are consumed, or explicitly rebuild/reconcile after each structural mutation. Avoid simply breaking the input loop unless queued-event delivery is verified. Extend [zoneEditorPointer_integration_test.go](../../test/integration/gui/zoneEditorPointer_integration_test.go) with two presses queued before one frame, same-edge double press, different-edge press, and deletion followed by selection. Use real input routing; no unit-test-only exports.

**Owner decision.** None; geometry-to-object identity must survive input batching.

### 1.15 🟠 Editor and preview classify portal connections differently

**Evidence.** [drawEdges](../../app/gui/dialogs/zoneEditorCanvas.go#L209-L221) recognizes `strings.EqualFold(connection.ConnectionType, "Portal")`. [getPreviewConnectionType](../../internal/services/preview_service/previewLayoutService.go#L224-L243) also classifies connections with `PortalPlacementRulesFrom`/`To` as portals. The latter comment explicitly states the game treats those as portals.

**Why wrong.** A placement-rule portal is presented as a direct connection in the manual editor but as a portal in the preview. Editing a type while retaining placement rules or loading such a manual connection exposes the discrepancy. This is separate from PNG pixels disappearing in §1.3.

**Fix.** Centralize model-level connection classification and pass the result through the existing geometry/handler seam to rendering. Use registry values, but do not stop at replacing the literal string: placement-rule semantics must also match. Extend [buildGeometry_test.go](../../test/unit/internal/services/connection_editor/zoneEditorGeometryService/buildGeometry_test.go), [buildPreviewLayout_test.go](../../test/unit/internal/services/preview_service/previewLayoutService/buildPreviewLayout_test.go), and the zone-editor GUI cases for explicit portal, placement-rule-only portal, and ordinary direct connection.

**Owner decision.** Confirm whether changing the editor type should clear old placement rules or merely display the effective type. Never silently rewrite protected schema vocabulary.

## §2 Architecture, product scope, and duplication

### 2.1 🟡 Connection-curve construction is duplicated with divergent visual policies

**Evidence.** [editor buildEdges](../../internal/services/connection_editor/zoneEditorGeometryService.go#L151-L197) canonicalizes endpoint order, groups pair edges, computes normal/spread/control/midpoint, uses `parallelEdgeGapPx = 18.0`, and applies `obstacleBulge`. [preview buildPreviewConnections](../../internal/services/preview_service/previewLayoutService.go#L164-L234) repeats the grouping/canonicalization/control-point construction with `spacingBetweenEdges = 21.0` and no obstacle deflection. This is a manual structural comparison, not a dupl-reported issue.

**Why it matters.** Endpoint grouping fixes and curve policy must be maintained twice. A curve visible around a zone in the editor can cross that zone in preview/export. This is verified divergence, not a proven explanation for every absent parallel edge in the owner's artifact.

**Fix.** First decide which visual differences are intentional. Extract only the common pair grouping/canonical Bézier geometry into a pure model/helper or service seam; retain intentional renderer-specific spacing as explicit parameters. Do not turn the editor's pixel hit tolerance into a game-domain field. Keep dedicated tests in both [buildGeometry_test.go](../../test/unit/internal/services/connection_editor/zoneEditorGeometryService/buildGeometry_test.go) and [buildPreviewLayout_test.go](../../test/unit/internal/services/preview_service/previewLayoutService/buildPreviewLayout_test.go); assert reversed pairs, three parallel edges, obstructed edges, and missing endpoints. Test any new public helper in its own mirrored function file.

**Owner decision.** Decide whether exact editor/PNG curve agreement is a requirement before changing goldens.

### Architecture inventory and non-findings

- Layering gates pass. Model-bearing DTOs and GUI-held models are intentional. `file_service` entity conversion and migrator entity access remain accepted. The owner has reopened the zone-content DTO exception in §2.2; the current tests still permit it until that work is implemented.
- `TemplateMapper` is the largest inspected production implementation: 648 physical lines, 590 nonblank lines. `ZoneEditorDialog` is 530 physical lines, 498 nonblank lines. Large catalog literals are not god objects.
- No actionable configured `dupl`, `funlen`, or `gocognit` findings were emitted. Existing method-split dialog/panel files are not automatically debt. TUI/web folders contain README placeholders and are not claimed implemented front ends.
- `PreviewLayoutService.layout` is mutable per-call scratch state shared by consumers. Current GUI calls are serialized; no current concurrent caller was established. Treat concurrent PNG rendering as requiring a call-local layout refactor and a race test, **not** as an already observed production race. Likewise generator `SetConfiguration`/`Generate` is a serialized protocol.

Decomposition guidance for work already justified elsewhere, not additional numbered findings:

| Function/family | Suggested boundary | Why / prerequisite |
| --- | --- | --- |
| TemplateMapper conversion/list methods | Template root, zone/variant, content, rules in same-package method files | Owner-approved maintainability split only; retain one primary mapper and preserve all round trips. |
| RebuildZoneConnectionRoads | Effective road policy, connection naming, generated-road reconstruction, preserved-target validation | §1.4/§1.10 need one coherent policy, not more literal flags. |
| TopologyConnectionService repair methods | Graph connectivity analysis, repair edge creation, optional road materialization | §1.6 must stop inferring graph connectivity from decorative roads. |
| ZoneEditorDialog.layoutStatus | Cached graph summary plus rendering | §3.1; avoid redoing graph analysis in a label closure. |
| buildEdges/buildPreviewConnections | Common endpoint-pair geometry, explicit render policy | §2.1; do not abstract intentional differences blindly. |

### 2.2 🟠 Owner-reopened zone-content DTO removal remains unimplemented

**Evidence.** [ZoneContentEditorService](../../internal/services/zone_content/zoneContentEditorService.go#L22-L49) exposes `ComposeContentRule(request dtos.ContentRuleCompositionRequestDto) dtos.ContentRuleCompositionResultDto`; [validRule](../../internal/services/zone_content/zoneContentEditorService.go#L120-L122) constructs the result DTO, while [zoneContentHandler](../../internal/handlers/zoneContentHandler.go#L27-L30) simply forwards it. Other service methods also consume DTO option/description shapes. [The architecture gate](../../test/unit/architecture/dependency/layering_test.go#L59-L69) currently calls this an accepted exception. When shown that conflict, the owner explicitly selected **“Reopen the DTO removal request.”**

**Why it matters.** Under that renewed direction, consumer DTO construction remains in the service instead of the handler. This is approved architecture debt, not a claim that the currently passing gate was secretly violated or that the form's behavior is broken.

**Fix.** Confirm the precise service result: the current composition unit returns a `ContentRuleRow`, whereas the original owner note mentions `ZoneContentRow`; do not conflate the two. Move request/result translation into the handler, use existing model/scalar inputs where meaningful, and retain a clear validity result for rejected selections. Inventory the remaining DTO signatures before deciding whether this batch removes the whole zone-content exception or only the composition seam. Never widen either allow-list; remove the zone-content entry only after all such dependencies are gone. Extend [composeContentRule_test.go](../../test/unit/internal/services/zone_content/zoneContentEditorService/composeContentRule_test.go) and [composeContentRule_test.go](../../test/unit/internal/handlers/zoneContentHandler/composeContentRule_test.go), covering each rule kind, invalid indices, nil/value preservation, and handler DTO construction. Run the architecture gate, full tests, coverage, and relevant content-dialog integration flows.

**Owner decision.** Reopening is authorized; the API/result shape and whole-exception vs single-seam scope still require the normal ask/plan/approval protocol. The bonuses DTO exception is not reopened by this decision.

### 2.3 🟠 Retire Ring, Hub, Chain, and Shared Web topologies and their tournament builders

**Evidence / owner request O08.** [Topology descriptors](../../internal/common/common_topologies/topologyDescriptors.go) still expose all four ordinary choices. Their serialized IDs in [mapTopology.go](../../internal/entities/topology/mapTopology.go) are `Default`, `HubAndSpoke`, `Chain`, and `SharedWeb`. [TopologyServiceLookup](../../internal/services/template_generator/providers/topologyServiceLookup.go) uses Ring as a fallback. [TournamentTopology](../../internal/services/template_generator/providers/topology/tournamentTopology.go) selects dedicated Ring/Hub builders, a balanced builder for Circles, and a chain fallback for other IDs, including surviving choices such as Random.

**Approved scope.** Remove the four ordinary choices and their implementations, dedicated wiring, and now-unused helpers/tests. Remove the corresponding Ring/Hub/chain tournament builders too. All surviving tournament selections that used the chain fallback must use the balanced builder, including Random, Square, Geometric, Cross, Fractal, and Geometric Hub. Preserve their ordinary implementations, Circles, and shared hub-zone concepts; preserving ordinary Geometric Hub does not mean preserving its former tournament fallback output. This is owner-requested product simplification, not a claim that all old algorithms are broken.

**Compatibility.** Reject an existing saved state selecting a retired topology with a clear message; do not silently substitute, accept via a default fallback, or discard its manual edits. Rejection must leave the current document unchanged. Keep only the minimal legacy-ID recognition needed for this error. Preserve supported saved states and distinguish retired IDs from unknown/invalid ones. The implementation plan must replace the ordinary lookup's Ring fallback with an explicit, panic-free invalid-selection contract, including zero/unknown IDs and any documented legacy default handling; do not conflate this with the approved balanced fallback for valid surviving tournament selections. No new tournament designs belong in this item.

**Implementation / verification.** Inventory descriptors, validation/load paths, ordinary/tournament dispatch, preview helpers, zone labeling, factories, constructors, benchmarks, and all test fixtures before deletion; retain code shared with surviving topologies. Update the [README topology catalogue](../../README.md) and other active user documentation. Regenerate Wire. Replace retired topology tests with explicit rejection and surviving-dispatch coverage, including Random tournament fallback, Circles, Geometric Hub, and all remaining supported choices. Extend [topology-provider tests](../../test/unit/internal/services/template_generator/providers/topologyProvider/createTopologyVariant_test.go), [tournament tests](../../test/unit/internal/services/template_generator/providers/topology/tournamentTopology/createTopologyVariant_test.go), and [wire-format integration tests](../../test/integration/editorStateWireFormat_integration_test.go); cover old-version rejected loads and unchanged current state after failure. Compile/run the appropriate unit, integration, GUI, and performance suites explicitly so retired constants cannot remain hidden behind tags. Add GUI selection coverage and remeasure coverage after the deliberate code/test deletions. Keep §1.12's two-player enforcement as a separate coordinated item.

### 2.4 🟠 Investigate compact manual deltas and derivable editor state before redesign

**Evidence / owner request O11.** [ManualZoneSave](../../internal/models/editor_state_model/manualZoneSave.go) persists complete edited zones; [EditorState](../../internal/models/editor_state_model/editorState.go) retains manual zone/connection snapshots. A single manual edit can therefore persist substantially more than the changed data. Redundant candidates include generated pools, guard reaction distribution, content limits, inferable zone classification, and unchanged default content rows. These are candidates to verify, not fields approved for deletion.

**Approved first deliverable.** A focused feasibility investigation covering **both live editor state and persisted data**, with representative byte/allocation measurements and a proposed design for added/changed/deleted zones and connections. Compare exact generated-base reconstruction against regenerating untouched parts and applying edits. Explain random generation/seed requirements, stable identities and endpoint references, generator-version upgrades, default evolution, derived versus explicitly overridden values, and compatibility with old full snapshots. Include dirty tracking, cloning, Undo/Revert-to-Base, setting changes, and final preview/export equivalence. Do not assume regeneration produces an identical base or silently erase fields merely because current generators usually derive them.

**Decision gate.** The owner chose investigation first and has **not** selected exact-base versus regenerated-base semantics. Present alternatives, measured benefits, risks, and migration limitations for owner approval before changing representation, adding deterministic generation infrastructure, or implementing delta persistence. The investigation must also state whether each proposed derived/default field can be omitted losslessly and how future default changes affect saved documents.

**Verification plan.** Analyze fixtures representing untouched generation, single-zone moves, additions/deletions, connection edits, explicit content overrides, nil/empty values, and old saved states. Define reconstruction/round-trip and size/performance acceptance tests for the chosen design; do not claim improvement without a baseline. If implementation is approved, add dedicated model/delta tests, migration tests, and [wire-format integration coverage](../../test/integration/editorStateWireFormat_integration_test.go), plus manual-edit GUI regression flows. Coordinate its approved format outcome with §2.5/§2.6 in one migration plan; no speculative migration is authorized by this investigation item.

### 2.5 🟠 Represent bans and value overrides as structured lists throughout

**Evidence / owner request O12.** [ContentSettings model](../../internal/models/editor_state_model/contentSettings.go) and [entity](../../internal/entities/editor_state/contentSettings.go) still store `BannedItems`, `BannedMagics`, and `ValueOverridesText` as strings. [GameRulesProvider](../../internal/services/template_generator/providers/gameRulesProvider.go) parses the strings during generation. In contrast, [BonusEntry entity](../../internal/entities/editor_state/bonusEntry.go), [model wrapper](../../internal/models/editor_state_model/bonusEntry.go), and [config aliases](../../internal/models/config/types.go) already provide the requested pattern.

**Approved scope.** Use structured entries throughout editor entities/models, model-bearing DTOs, configuration, and handler/service crossings, following the `BonusEntry` approach. Ban entries carry `SID`; override entries carry `SID` and `GuardValue`. Preserve the current fixed output `Variant=-1`. Text parsing/formatting belongs at UI or legacy migration boundaries, not in generation over already-typed data. Retain existing order, duplicate handling, malformed-input behavior, empty/nil semantics, and export meaning. If typed migration cannot preserve a legacy case, surface the concrete conflict for owner approval rather than inventing a cleanup policy.

**Implementation / verification.** Inventory all readers, parsers, cloning/conversion seams, and UI controls before replacing the string fields; do not widen layer exceptions or edit the protected RMG schema. Preserve malformed-text warnings at the UI/legacy parsing boundaries: the implementation plan must trace their delivery to the user and specify the resulting typed generation API, rather than silently dropping warnings or testing malformed text at a seam that no longer parses it. Use one coordinated new-schema migration with §2.6, informed by §2.4's decision gate, while continuing to load supported older states. Extend [override generation tests](../../test/unit/internal/services/template_generator/providers/gameRulesProvider/createValueOverrides_test.go), content-state clone/mapping tests, and [wire-format integration tests](../../test/integration/editorStateWireFormat_integration_test.go). Test multiple entries, blanks, malformed legacy text, duplicates/order, round trips, mutation isolation, and unchanged exported rules; add Bonuses-panel editing/save/load GUI coverage. This item does not authorize a separate bonuses DTO-exception refactor.

### 2.6 🟠 Persist editor settings in sections while retaining flat identity/map/schema fields

**Evidence / owner request O13.** [EditorState entity](../../internal/entities/editor_state/editorState.go) embeds ten groups into a flat aggregate; [JSON tests](../../test/unit/internal/entities/editor_state/editorState/editorStateJson_test.go) and [wire-format integration tests](../../test/integration/editorStateWireFormat_integration_test.go) pin that current representation.

**Approved shape.** Use the existing groups as sections for `PlayerSettings`, `NeutralZoneSettings`, `CastleSettings`, `GenerationSettings`, `GameRuleSettings`, `ContentSettings`, and `ManualEditSettings`. Keep the fields of `TemplateIdentity`, `MapSettings`, and `SchemaOptions` at the root, not under three extra section objects. In particular, root-level `schemaVersion` remains available for version dispatch. Propose exact JSON section keys in the implementation plan using the current lower-camel convention; do not silently regroup settings by GUI tabs or change field meanings.

**Migration / verification.** Coordinate §2.5 and any owner-approved §2.4 format work into **one migration**, preserving loading of all supported older schema versions and leaving frozen v1 snapshots unchanged. Separate versioned wire shapes from the live model as needed within the existing layering rules; repositories remain typed decoders. Extend dedicated serialization, migrator, and wire-format tests for flat legacy input, the new nested shape, the three flat groups, typed content entries, manual geometry, nil/empty values, unsupported versions, and current-version round trips. Save/load/preview/export must retain the same user choices and data; test malformed sections and failed loads without overwriting the current document. No change to RMG output schema or output-path persistence belongs here.

### 2.7 🟠 Move zone-content accessors from LayoutPanel closures to the state driver

**Evidence / owner request O14.** [layoutPanelZones.go](../../app/gui/panels/layoutPanelZones.go) defines six inline getter/setter pairs for Player, Lowest, Low, Medium, High, and Hub rows, then passes them to `openZoneContentDialog`. Updates currently go through [drivers.State.UpdateState](../../app/gui/drivers/state.go), which owns validation/dirty tracking around the GUI state model.

**Approved scope.** Put the zone/tier content access API on **`app/gui/drivers.State`**, not the domain `EditorState` and not the GUI model. Replace the panel's repeated callbacks with a clear zone/tier selection passed to driver accessors. Keep the driver responsible only for view-state access/update orchestration; it must not take over content composition or generation policy. Preserve existing validation, dirty detection, and cloned-snapshot ownership; getters must not expose mutable aliases into retained state. Dialog cancellation remains non-mutating.

**Verification.** Add mirrored unit tests for every new public driver accessor, covering all six targets, unrelated rows unchanged, empty/nil values, input/output mutation isolation, no-op versus changed updates, and invalid selectors if the API permits them. Use existing production update paths rather than test-only exports. Extend content-dialog GUI flows for open, edit, cancel, and Apply. Coordinate with §2.9 so the panel move reuses the driver API rather than recreating closure-based setters inside the new sections.

### 2.8 🟡 Rename the zones service package and its nested interface package

**Evidence / owner request O16.** [ZoneFactory](../../internal/services/zones/zoneFactory.go) lives under the zones service package, with [IZoneFactory](../../internal/services/zones/zone_interfaces/zoneFactoryInterface.go) in its interface subpackage.

**Approved targets.** Rename the implementation package to internal/services/zone_services and its **nested** contract package to internal/services/zone_services/zone_service_interfaces. Update package declarations, imports, qualifiers, mirrored test directories, relevant documentation references, and Wire provider references together. Keep interfaces prefixed `I` in their separate `*Interface.go` files. This is a naming-only owner request: do not combine it with behavior changes, new abstractions, or unrelated renames.

**Verification.** Regenerate Wire rather than hand-editing generated code. Run build, test-layout checks, full unit tests, architecture gates, coverage, lint, and affected integration suites on supported platforms. Search active imports/declarations for the old paths; historical review evidence may still describe the reviewed revision. Update active documentation links after the move. Existing tests move with implementations and retain dedicated public-method coverage; add tests only for uncovered adjacent logic, not assertions about package spelling.

### 2.9 🟠 Split General, Layout, and Bonuses panels into private cohesive sections

**Evidence / owner request O18.** [GeneralPanel](../../app/gui/panels/generalPanel.go), [LayoutPanel](../../app/gui/panels/layoutPanel.go), and [BonusesPanel](../../app/gui/panels/bonusesPanel.go) each own many controls and their load/save/click handling. [Window](../../app/gui/editor/window.go) constructs the panels and coordinates per-frame saving and load refresh. The Layout method-file split does not give its sections their own state ownership.

**Approved structure.** Move those three panels into app/gui/panels/general_panel, app/gui/panels/layout_panel, and app/gui/panels/bonuses_panel respectively. Each cohesive section gets its own private struct owning its widgets, click handling, and state load/save operations; the enclosing panel composes sections and coordinates them. Expose the public panel structs through aliases in a new app/gui/panels/types.go and preserve existing constructors/public access with thin forwarding functions where needed. Leave Preview unchanged. Child packages must not import the parent alias package; keep shared interfaces in their existing cycle-free location.

**Behavior boundary.** Preserve appearance, ordering, keyboard/pointer behavior, frame-to-frame widget identity, save/load semantics, and public API behavior. No new content presets, domain logic in GUI, live state pointers, or business-layer refactors. Inventory section boundaries in the implementation plan before moving fields. Coordinate zone-content controls with §2.7 and the Bonuses controls with §2.5 to avoid overlapping rewrites.

**Verification.** Keep one primary struct per file and mirrored test layout. Move the existing tagged General/Layout panel test-export methods with their owning structs: Go cannot define methods on aliases of non-local types. Preserve per-file tags and run both gated integration and GUI suites that consume those exports. Run build and architecture checks after package moves, existing window/tab save/load and content-dialog GUI suites, and snapshot comparisons without blindly updating goldens. Exercise each section's state round trip and event routing through integration tests; any extracted non-Gio logic requires its own unit coverage. Compare TabCycling performance before/after because Gio widget identity and retained state are load-bearing; do not introduce an automated global allocation threshold.

### 2.10 🟡 Move neutral content-row selection to GeneratorConfig

**Evidence / owner request O19.** The package-level [neutralRowsForQuality](../../internal/services/template_generator/providers/mandatoryContentProvider.go) maps Lowest/Low/Medium/High to the corresponding configured neutral rows, Highest to Hub rows, and unknown quality to nil. [GeneratorConfig](../../internal/models/config/generatorConfig.go) already owns these configuration values.

**Approved scope.** Replace the private provider lookup with a public `GeneratorConfig` method and update both callers, retaining their respective direct-quality and resolved-quality inputs. Preserve every quality mapping, including Highest→Hub, nil versus empty rows, and the existing read-only borrowed-slice ownership contract. Do not introduce new cloning/allocation or expand lookup behavior incidentally; consumers must not mutate returned configuration rows. Do not add a competing public helper or alter standalone Hub generation policy.

**Verification.** Add a dedicated mirrored test file for the new public method, with separate quality cases including unknown, nil, empty, and populated rows. Retain [mandatory content tests](../../test/unit/internal/services/template_generator/providers/mandatoryContentProvider/createContents_test.go), especially Highest-tier content, and confirm generated content is unchanged. Coordinate with §2.5 if shared content types move.

## §3 Performance

### 3.1 🟠 Editor graph diagnostics are rebuilt on every frame despite geometry caching

**Evidence.** [layoutStatus](../../app/gui/dialogs/zoneEditorDialog.go#L231-L256) calls `derefConnections(this.working)` and `DescribeZoneEditorGraph(...)` on each layout. [derefConnections](../../app/gui/dialogs/zoneEditorDialog.go#L524-L530) allocates a connection slice. [DescribeZoneEditorGraph](../../internal/handlers/zoneEditorHandler.go#L85-L92) constructs both checks, while [FindIsolatedZones](../../internal/services/connection_editor/connectionEditorService.go#L45-L62) scans connections once per zone. Canvas geometry already has a [dirty gate](../../app/gui/dialogs/zoneEditorCanvas.go#L61-L67).

**Why it matters.** Reasoned cost: allocation plus O(zones × connections) diagnostic work on unchanged frames. No frame-time/allocations benchmark was taken, so neither a claimed percentage improvement nor a perceptible slowdown is asserted.

**Fix.** Cache a graph summary with an invalidation flag updated for structural changes; account for toolbar layout occurring before canvas input so the summary is correct on the next frame. Optionally implement isolation detection in O(zones + connections) using a referenced-name set. Extend [findIsolatedZones_test.go](../../test/unit/internal/services/connection_editor/connectionEditorService/findIsolatedZones_test.go) and [describeZoneEditorGraph_test.go](../../test/unit/internal/handlers/zoneEditorHandler/describeZoneEditorGraph_test.go). Add an untagged public-API benchmark in test/performance/zoneEditorGraph_test.go; use GUI integration to count handler calls across idle vs mutation frames. Measure before deciding whether broader caching is worthwhile.

**Owner decision.** Do not reintroduce the rejected global Gio allocation threshold or clone-free live state pointers.

## §4 Readability, comments, and maintainability inventory

`git grep` over tracked Go source found **zero TODO/FIXME/HACK/XXX markers**; configured `godox` also emitted zero issues. There is therefore no omitted source TODO list. Free-form suggestions remain in [zoneEditorGeometryService.go](../../internal/services/connection_editor/zoneEditorGeometryService.go#L133-L135) (`"Probably can just be a Vec3"`) and its [guide/snap helper comments](../../internal/services/connection_editor/zoneEditorGeometryService.go#L272-L295). Disposition: accepted owner follow-up/style suggestion, no demonstrated benefit from replacing three scalar offsets or allocating different guide arrays. A commented-out experimental banner is not a correctness finding.

Magic drawing constants mostly have descriptive local names. The portal literal is addressed in §1.15. `getCurrentMapSize` indexes the complete map-size list by the selector index; the normal list is currently its prefix and index clamping exists. This is a verified current non-issue, not a panic claim. A proposed “reused zone name leaves stale controls” bug was rejected: `addZoneAt` explicitly resets `syncedZoneFor`.

### 4.1 🟡 Replace equivalent handwritten vector arithmetic with Vec2 operations

**Evidence / owner request O20.** [Vec2](../../internal/helpers/data/vec2.go) provides value-returning arithmetic, distance, and related operations. [FindOpenPosition](../../internal/services/connection_editor/zoneEditorService.go) already uses some of them. Geometry files are audit candidates, not blanket rewrite targets; polar placement in [balanced-ring layout](../../internal/services/preview_service/layoutBalancedRings.go) is not itself evidence of a missing vector-method replacement. This is an owner-requested consistency audit, not a finding that all scalar math is incorrect. The earlier Vec3 comment disposition does not exclude equivalent Vec2 replacements from this audit.

**Approved scope.** Audit all non-protected production `data.Vec2` usage and replace calculations where an existing vector operation is equivalent and clearer. Record justified exclusions rather than forcing unrelated scalar/angle math into vectors. Add missing vector methods only for concrete callers where they improve the code, with their own tests. Exclude protected game-data/schema/registry trees and generated code; regenerate any generated output only through its normal generator. No bulk in-place rewrite.

**Verification.** Preserve operation order where floating-point rounding matters, zero-length/degenerate behavior, coordinate units, and value ownership. Any intentional numerical change requires separate approval, not concealment as cleanup. Add or extend each affected public API's mirrored tests, including representative axis/diagonal/negative/degenerate cases, and test new methods directly. Run geometry/layout/PNG tests and relevant GUI snapshots, comparing outputs rather than automatically accepting new goldens. Measure hot-path allocations/timing where touched and avoid new allocations or slower general-purpose abstractions. Coordinate overlapping geometry edits with §1.13/§2.1 and preserve dedicated regression assertions.

## §5 Testing and coverage interpretation

The default suite passed **183 packages**, including **181 unit-test packages**; tagged integration/performance passed two packages, with the performance package reporting **no tests to run** in the non-benchmark command. Geometry, PNG generation, and untagged integration were rerun with `-count=20`, all passing. These repetitions do not cover the newly identified batched-input or symmetric-obstacle scenarios until tests are added.

Coverage is statement coverage, not branch/path or value-contract coverage. In particular:

- Preview image tests permit “portal omitted” to satisfy “portal differs from solid” (§1.3).
- Manual-apply mocks omit the mutating behavior of the production handler (§1.8).
- Deep-clone tests on `EditorState.Clone` do not test the shallow `SetManualEdits`/`GetManualZones` boundaries (§1.9).
- Roads-off generation tests do not exercise repair and manual apply (§1.4/§1.6).
- A true-output assertion currently endorses a false-input tournament SaveArmy result (§1.7).
- `ComputeHasErrors` checks missing endpoints, not all road/content referential integrity (§1.10).

These are specific test-plan gaps attached to the production findings, not six additional duplicate findings. Current test-layout enforcement passes; no blanket build tags or unit imports of test-only exports are proposed. Genuinely Gio-bound files remain integration territory. The complete below-100% file inventory is in §11; protected registry/schema gaps are not instructions to modify them.

**Local execution limits:** no new local full race run, GUI suite, or benchmark run was made. The review contract's required default and tagged non-GUI commands were run. CI contains Linux race and Mesa/Xvfb GUI jobs, verified as configuration, not as a claim about this revision's remote job result. Windows results are measured here; Linux execution is not.

## §6 CI/CD and dependency tooling

### 6.2 🟠 The tools module installs a different linter than CI and this baseline

**Evidence.** [tools/go.mod](../../tools/go.mod#L90-L98) pins `github.com/golangci/golangci-lint/v2 v2.12.2`; [PR lint](../../.github/workflows/pr-validation.yml#L72-L77) specifies `version: v2.13.1`. The installed executable used here reports `2.13.1`.

**Why it matters.** Rebuilding the declared tool command does not reproduce the measured lint baseline or the CI analyzer set. Both versions may be clean today; that does not make them equivalent.

**Fix.** Select one version, update module/workflow/documentation together, and run tidy/build of the tools module plus report-only lint with the chosen binary. Keep warning-only CI exclusions explicit. Verify `go version -m`/tool version matches the pin and rerun the normal tests/coverage if module changes affect the application. No new Go unit test is needed for a dependency pin alone.

**Owner decision.** Choose the supported linter version; do not infer the tools module's broader dependencies should be upgraded wholesale.

### 6.3 🟡 Module/checksum files lack the LF policy their tidy checks require on Windows

**Evidence.** [.gitattributes](../../.gitattributes#L1-L5) says `* text=auto` and only `*.go text eol=lf`, despite the comment mentioning modules. `git ls-files --eol` reports `i/lf w/crlf attr/text=auto` for both module files and both checksum files. Both `go mod tidy -diff` runs exit **1** and produce only checksum-file line replacements: root **45 removed/45 added**, tools **988/988**. Comparing normalized removed/added lines yielded **zero content differences** in both.

**Why it matters.** A clean Windows checkout fails the documented dry-run consistency check even though dependency content is unchanged. This is an EOL reproducibility defect, not module drift or evidence of a failing Linux CI tidy job.

**Fix.** Add explicit LF attributes for module/checksum files, then have the owner approve narrowly scoped normalization of those four paths. Re-run tidy dry-runs from Windows and Linux and verify no dependency-content changes. Avoid a repository-wide normalization/rewrite. No new Go test is warranted.

**Owner decision.** Attribute/normalization change only; do not stage or commit during the fix session.

### 6.5 🟠 Release tag input is interpolated directly into shell source

**Evidence.** [release build](../../.github/workflows/release.yml#L60-L63) contains `go build -trimpath -ldflags "-s -w -X main.version=${{ github.event.inputs.tag || github.ref_name }}" ...`. [workflow_dispatch](../../.github/workflows/release.yml#L7-L11) accepts a free-form tag input.

**Why it matters.** Expression substitution happens before shell parsing; double quotes do not convert substituted source into inert data. A specially formed, checkout-resolvable ref can alter shell evaluation. This is a trusted-maintainer/tag-controlled hardening surface, **not** an arbitrary-PR exploit. The build job is read-only, while publication is a separate write-permission job.

**Fix.** Pass the selected release tag via an environment variable and expand it as data; validate an approved tag/version format before checkout/build/publish. Keep checkout ref, injected version, release name, and published tag consistent. Add workflow validation cases for ordinary tags, whitespace, quote/shell-metacharacter rejection, and invalid/missing refs. If strict semantic versioning would reject supported tags, agree the format first. No production Go test is required unless validation is implemented as a Go helper, which then needs its mirrored tests.

**Owner decision.** Confirm allowed release tag formats and manual-release behavior. Do not test this by publishing or dispatching an actual release.

### Verified CI/CD protections and limits

- Both root/tools modules and workflow Go versions are `1.27.0`, matching the local compiler. No Go-version drift finding.
- Linux build/vet, Windows build/unit tests, Linux unit race detection, coverage trend and minimum, integration/performance, and opt-in GUI/Mesa jobs exist in [pr-validation.yml](../../.github/workflows/pr-validation.yml). The configured coverage floor is **60.0%**, not the historical memory floor; the current measured 74.4% and no-decrease rule should anchor fixes.
- Release already uses `-trimpath`, version injection, SHA-256 checksum output, correct dispatch checkout ref, separate publish permissions, and concurrency with `cancel-in-progress: false`. They are not missing gates to add again.
- Scheduled security scanning is deliberately schedule-only. PR vulnerability scanning also exists. Do not add workflow_dispatch or remove the reduced post-merge job subset without owner direction.
- `git ls-files -ci --exclude-standard` returned no tracked ignored artifacts. 2,864 files are tracked. Coverage reports and output templates are ignored/untracked; committed generated Wire code and GUI goldens are intentional.
- The latest squash touches 1,281 files, +22,396/−14,060, including relocation and deletion of old review records. This justifies reviewing behavior across seams; volume alone is not misconduct or a defect.
- Current gopls analysis includes `integration_test,gui`, unlike older notes. No global Go build/test tag was identified. Analysis-only flags are not proof that GPU tests run by default; do not mislabel them as that defect. Ask the owner before changing editor visibility policy.

## §7 Docs and developer experience

### 7.1 🟡 README's dependency, pipeline, and license descriptions are stale

**Evidence.** [README](../../README.md#L88-L91) advertises `gioui.org v0.10.0`; [go.mod](../../go.mod#L5-L12) requires `v0.10.2`. The [pipeline diagram](../../README.md#L243-L257) ends in `entities.RmgTemplate`, whereas [Generate](../../internal/services/template_generator/templateGenerator.go#L57-L91) returns a model. [license text](../../README.md#L314-L316) says “See the main project repository,” although [LICENSE](../../LICENSE) is present.

**Why it matters.** Contributors following the architectural diagram can reintroduce the entity coupling the repository deliberately removed; dependency/license guidance is needlessly inaccurate.

**Fix.** Update these three descriptions together, naming the Model→Entity persistence seam and linking the local license. Prefer linking the module over repeating a patch version if no version-specific explanation is needed. Check README/QUICKSTART commands and internal links against the tree. Documentation-only change needs link/command verification, not new Go tests.

**Owner decision.** None; do not expand this into a documentation rewrite.

### 7.2 🟡 Remaining observation paths and CI claims need reconciliation

**Evidence / progress.** The obsolete branch/push instructions identified in the original review have been replaced by the owner-requested [current handoff](../session-carry-forward.md). That portion is addressed. Remaining [test observations](test_observations.md) still refer to older flat snapshot-helper paths, `suggestDirectory`, and GPU tests never running in CI. Current source uses [snapshot helpers](../../test/test_helpers/integration_common/snapshot), [getWorkingDirectory](../../app/gui/drivers/stateFiles.go), and a [Mesa/Xvfb GUI job](../../.github/workflows/pr-validation.yml#L235-L278). The old benchmark figures are historical, not current measurements.

**Why it matters.** Coverage limitations remain useful only if their paths and execution claims reflect the current test architecture. Updating the handoff does not automatically resolve these remaining observations or invalidate settled testability decisions.

**Remaining work.** Correct moved paths and stale symbols/CI claims; distinguish historical benchmark figures from current evidence and keep the intentional absence of an automated allocation threshold. Verify surviving local links and branch statements without fabricating remote status. Do not re-create retired operational instructions or remove useful coverage-limit explanations. No new Go unit tests are needed; mark the whole item fixed only after its remaining documentation scope and owner-commit protocol are satisfied.

## §8 Full linter disposition

Command: `golangci-lint-v2 run ./... --issues-exit-code=0`, installed version **2.13.1**, exit **0**, final summary **0 issues**. Three unused-exclusion warnings appeared (TODO/godot and comment-form exclusions); these are configuration warnings, not linter findings. No formatting auto-fix was run.

| Linter / class | Count | Disposition |
| --- | ---: | --- |
| dupl | 0 | No configured clone findings; manually verified divergence is §2.1. |
| funlen | 0 | No configured overlength findings; decomposition candidates are advisory in §2. |
| gocognit | 0 | No configured complexity findings. |
| godox | 0 | Zero tracked Go TODO/FIXME/HACK/XXX markers independently checked. |
| depguard | 0 | Current architecture exceptions are documented decisions, not hidden violations. |
| gosec | 0 | Not proof that workflow shell interpolation is safe; §6.5 is outside this Go lint result. |
| paralleltest / testpackage / testifylint / tparallel | 0 each | Current test-style baseline clean. |
| All other enabled analyzers/formatters | 0 each | No emitted finding needing a fix, exclusion, or acceptance entry. |
| **Total emitted issues** | **0** | **Every current finding accounted for. No new exclusions proposed.** |

The configured run includes existing exclusions (protected registry duplication, test-specific rules, accepted house naming). A zero result means zero *reported* issues under that configuration, not zero duplication or technical risk in all source. CI additionally disables `godox,dupl,unparam,gochecknoglobals`; do not equate its narrower result with the full local baseline.

## §9 Suggested execution order

| Batch | Items | Dependencies / verification |
| --- | --- | --- |
| A: prevent silent loss/misplacement | §1.1, §1.2 | Dirty/exit tests and failed-detection export refusal. Highest user impact; preserve output hard rule. |
| B: PNG and value contract | §1.3, §1.7 | Independent small fixes with public-API pixel and tournament-value tests. Confirm game's omitted-false semantics. |
| C: road and graph invariants | §1.4, §1.6, §1.10 | Agree generated/custom-road handling, then separate graph connectivity from road policy. End-to-end road-target integrity and roads-off matrices. |
| D: manual state lifecycle | §1.8, §1.9 | Coordinate with A's dirty tracking and C's road-rebuild changes; retain compare-before-mutation tests even after roads-off is fixed. Land after C or include §1.8 in C. |
| E: effective modes and guard propagation | §1.5, §1.11, §1.12 | Owner decisions on edit invalidation, custom guard values, player-count restoration. Do not include topology deletion/redesign automatically. |
| F: editor geometry | §1.13, §1.14, §1.15; optionally §2.1 | Determinism first; batched real-input tests; shared classification; owner approval before visual-policy consolidation. |
| G: measured performance | §3.1 | Establish public-API/GUI baseline before caching; no flaky global allocation threshold. |
| H: CI/tooling hardening | §6.2, §6.3, §6.5 | Independent configuration changes: linter alignment, Windows/Linux EOL checks, and safe release-tag validation. |
| I: docs | §7.1, §7.2 | Can run independently after owner approves retirement scope. |
| J: reopened service boundary | §2.2 | Confirm the composition result/API and exception-removal scope first; preserve bonuses exception. Independent from geometry and road fixes. |
| K: topology retirement | §2.3 | Inventory shared builders before removal; reject retired saved IDs, reroute surviving tournament fallbacks to balanced generation, and coordinate §1.12. |
| L: compact-state investigation | §2.4 | Can begin independently; owner reviews live/persisted feasibility and reconstruction semantics before implementation or schema design. |
| M: coordinated persistence format | §2.5, §2.6; approved outcome of §2.4 only | Follow L's decision gate. One migration plan, supported legacy loading, typed entries, and explicit root-versus-section wire assertions. |
| N: panel/state organization | §2.7, §2.9 | Driver accessors first; preserve public panel API, GUI snapshots, and state round trips. Coordinate Bonuses work with M. |
| O: focused naming and lookup | §2.8, §2.10 | Separate behavior-preserving changes; schedule rename outside overlapping service work, regenerate Wire, test the config lookup directly. |
| P: vector arithmetic audit | §4.1 | Coordinate with F; equivalent permitted production edits only, dedicated numerical tests and hot-path checks. |

For every code batch: capture coverage before/after, build, unit suite, test-layout check, full lint baseline, relevant tagged integration/GUI suites, and Wire regeneration when constructors/provider sets change. Preserve **74.4% or higher comparable Windows coverage** and **0 lint issues**, unless a measured denominator/platform change is separately explained and approved. Owner stages/commits. Mark items in place; never renumber.

## §10 Security, dependencies, owner validation, and verified non-issues

**Vulnerability status.** `govulncheck` is not installed locally. An attempted tool lookup via the tools module was not a successful scan; no vulnerability count is claimed. The contract allows confirmation of CI coverage: [PR vulnerability job](../../.github/workflows/pr-validation.yml#L79-L99) and [scheduled security scan](../../.github/workflows/security-scan.yml#L14-L33) run `golang/govulncheck-action@v1` against `./...`. Actual remote results were not fetched. No dependency is labelled vulnerable/unmaintained merely because it is old or indirect.

**File-system safety checked.** [SanitizeFilename](../../internal/helpers/string.go) removes path separators/unsafe characters; [PathResolutionService](../../internal/services/file_system/pathResolutionService.go) checks names and resolves destinations; [atomicFileWriter](../../internal/repositories/atomicFileWriter.go) uses temporary write/sync/close/rename and cleanup. No concrete template-name traversal or command execution in application template handling was found. The dangerous output fallback is §1.2, not a demand to persist a directory.

**Owner/in-game validation queue (not numbered application findings):**

1. Bonuses/bans: generate known valid SIDs, inspect exported fields, then compare in-game behavior with a shipped example. If the game's accepted vocabulary differs, obtain owner approval before any protected data/schema change. Source emission alone cannot close O03.
2. Hero hire: [CreateGameRules](../../internal/services/template_generator/providers/gameRulesProvider.go#L21-L34) uses `HeroHireBan: configuration.IsSingleHeroMode()`; [LostStartHero](../../internal/services/template_generator/providers/gameRulesProvider.go#L105-L110) additionally enables for arena/explicit loss. Confirm whether the latter modes must ban hiring; then add mode-matrix tests in the existing game-rules provider test file. This is an owner-requested policy check, not proven game malfunction.
3. Parallel-edge artifact: retain a current saved state plus PNG showing a missing **solid** curve after §1.3; distinguish node occlusion and §2.1 geometry differences from edge omission. Do not mark O04 fixed just because pair grouping exists.
4. Geometric Hub v1 artifact: new schema-v2 saves preserve stamps, but missing historical coordinates cannot be reconstructed by plain migration. Recreate/re-save under current code and compare before/after load. Do not discard manual edits silently to manufacture a new layout.
5. Historical preview artifacts: §1.13 proves a tie nondeterminism, not every saved artifact's cause. Retain an exact new reproducer if artifacts remain.

**Other verified non-issues:**

- `EqualsIgnoringManualEdits` covers current non-manual fields; excluding manual data from *regeneration* is intentional. Dirty tracking needs its own fix, §1.1.
- Nil vs empty mapping, nil quality fallback, ring zero, versioned migration, and frozen v1 inputs remain load-bearing. The migrator's double open and typed repositories are explicit owner decisions.
- Existing template clones, neutral/hub content-row cloning, error checks in template generation, and missing-endpoint guards avoid several obvious nil/value-copy false positives.
- SharedWeb ensures a neutral plan before modulo-based spoke assignment; Delaunay bounds callers guard small/empty inputs; connectivity repair reduces components and has a terminating fallback. No verified infinite loop found.
- The accepted isolated-map safety fallback is not itself a defect. Its road-based connectivity predicate is (§1.6).
- Per-frame deep-copy overhead was accepted by the owner. No speculative live-pointer access API or threshold is proposed.
- A spawn with zero extra castles still contains a Spawn main object. Source suggests different road anchoring in that mode, but whether that is incorrect game policy was not established; it is not promoted as a proved bug.
- Mutable shared preview scratch and content-slice sharing without a current mutating caller are limitations to consider when adding concurrency/mutation, not reported current races or corruption.

## §11 Measured baselines and complete per-file gaps

### Verification ledger

| Check | Measured result |
| --- | --- |
| go version | `go version go1.27.0 windows/amd64` |
| golangci-lint-v2 version | `2.13.1`, built with Go `1.27.0` |
| go build ./... | Exit 0 |
| go vet -tags=integration_test ./... | Exit 0 |
| go test ./test/... -count=1 | Exit 0, 183 passing packages |
| go test -tags=integration_test ./test/integration/... ./test/performance/... -count=1 | Exit 0; integration passed, performance package had no tests to run |
| Unit coverage command with coverpkg=./internal/...,./app/... | Exit 0, 181 passing unit packages; profile directed to temporary storage instead of overwriting the existing report |
| go tool cover -func on that profile | **74.4% total statements** |
| Go-generated coverage HTML file inventory | **282 instrumented files: 192 at 100%, 52 partial, 38 at 0%** |
| Report-only golangci-lint full configured run | Exit 0, **0 issues**, 3 unused-exclusion warnings |
| go run ./cmd/testlayoutcheck . | Exit 0, `test-layout check passed` |
| Root / tools go mod tidy -diff | Exit 1 / 1, **EOL-only checksum differences**, no normalized dependency-content changes (§6.3) |
| Twenty repetitions: geometry unit package, PNG unit package, untagged integration package | Exit 0 for all three |
| Public-API scratch reproductions | Portal absent; Set/GetManualZones alias; SaveArmy ignored; manual roads restored; two obstacle-curve results; program deleted |
| govulncheck | Local scan not run; PR and scheduled CI coverage confirmed, remote outcome unknown |
| git tracked-vs-ignore check | 0 tracked ignored artifacts; 2,864 tracked paths |
| Local Linux / full race / GUI / performance execution | Not run; do not infer pass from configuration or historical notes |

**Coverage method caution.** The profile contains repeated blocks from multiple test executables. A naïve sum of every profile line produced invalid per-file numbers and was discarded. Per-file percentages below come from **Go's own generated HTML**, and total **74.4%** comes from **Go's own `-func` output**. Do not replace that total with a homemade deduplication result or an average of file percentages. Files not instrumented in this Windows unit run (including some GUI entry/panel packages and other-OS files) are **not** implicitly 100% or 0%; they are outside this profile.

### All zero-coverage files (38)

All 36 app files below are Gio/rendering/view composition and should be considered with the accepted integration observations, not subjected to artificial unit seams. The two internal files have explicit host/protected-data dispositions.

| File | % | Disposition |
| --- | ---: | --- |
| [dropdownItem.go](../../app/gui/components/dropdownItem.go) | 0 | UI composition |
| [dropdownSelector.go](../../app/gui/components/dropdownSelector.go) | 0 | UI integration |
| [segmentButton.go](../../app/gui/components/segmentButton.go) | 0 | UI composition |
| [segmentButtonGroup.go](../../app/gui/components/segmentButtonGroup.go) | 0 | UI integration |
| [bonusPickerDialog.go](../../app/gui/dialogs/bonusPickerDialog.go) | 0 | GUI integration |
| [fileExplorerDialogConfirm.go](../../app/gui/dialogs/fileExplorerDialogConfirm.go) | 0 | GUI integration |
| [fileExplorerDialogToolbar.go](../../app/gui/dialogs/fileExplorerDialogToolbar.go) | 0 | GUI integration |
| [pickerDialog.go](../../app/gui/dialogs/pickerDialog.go) | 0 | GUI integration |
| [ruleDialog.go](../../app/gui/dialogs/ruleDialog.go) | 0 | GUI integration |
| [zoneContentDialog.go](../../app/gui/dialogs/zoneContentDialog.go) | 0 | GUI integration |
| [zoneContentRow.go](../../app/gui/dialogs/zoneContentRow.go) | 0 | UI composition |
| [zoneContentSection.go](../../app/gui/dialogs/zoneContentSection.go) | 0 | GUI integration |
| [zoneEditorCanvas.go](../../app/gui/dialogs/zoneEditorCanvas.go) | 0 | Add targeted GUI coverage for §1.14/§1.15 |
| [zoneEditorConnectionProps.go](../../app/gui/dialogs/zoneEditorConnectionProps.go) | 0 | GUI integration |
| [zoneEditorDialog.go](../../app/gui/dialogs/zoneEditorDialog.go) | 0 | GUI integration; §1.8/§3.1 seams |
| [zoneEditorInteractionState.go](../../app/gui/dialogs/zoneEditorInteractionState.go) | 0 | GUI interaction |
| [zoneEditorSnap.go](../../app/gui/dialogs/zoneEditorSnap.go) | 0 | GUI interaction |
| [zoneEditorZoneProps.go](../../app/gui/dialogs/zoneEditorZoneProps.go) | 0 | GUI integration; §1.11 |
| [tab.go](../../app/gui/drivers/tab.go) | 0 | GUI view state |
| [theme.go](../../app/gui/themes/theme.go) | 0 | GUI theme construction |
| [draw.go](../../app/gui/utils/draw.go) | 0 | Gio rendering |
| [buttonWidget.go](../../app/gui/widgets/buttonWidget.go) | 0 | GUI composition |
| [centeredMessageWidget.go](../../app/gui/widgets/centeredMessageWidget.go) | 0 | GUI composition |
| [emptyWidget.go](../../app/gui/widgets/emptyWidget.go) | 0 | GUI composition |
| [labelWidget.go](../../app/gui/widgets/labelWidget.go) | 0 | GUI composition |
| [labeledCheckboxRowWidget.go](../../app/gui/widgets/labeledCheckboxRowWidget.go) | 0 | GUI composition |
| [labeledRowWidget.go](../../app/gui/widgets/labeledRowWidget.go) | 0 | GUI composition |
| [labeledSliderWidget.go](../../app/gui/widgets/labeledSliderWidget.go) | 0 | GUI composition |
| [panelWidget.go](../../app/gui/widgets/panelWidget.go) | 0 | GUI composition |
| [sectionWidget.go](../../app/gui/widgets/sectionWidget.go) | 0 | GUI composition |
| [sliderRowWidget.go](../../app/gui/widgets/sliderRowWidget.go) | 0 | GUI composition |
| [spacerWidget.go](../../app/gui/widgets/spacerWidget.go) | 0 | GUI composition |
| [splitWidget.go](../../app/gui/widgets/splitWidget.go) | 0 | GUI composition |
| [textboxWidget.go](../../app/gui/widgets/textboxWidget.go) | 0 | GUI composition |
| [titleBarWidget.go](../../app/gui/widgets/titleBarWidget.go) | 0 | GUI composition |
| [warningBannerWidget.go](../../app/gui/widgets/warningBannerWidget.go) | 0 | GUI composition |
| [io_windows.go](../../internal/helpers/io_windows.go) | 0 | Real Steam registry lookup; accepted host-dependent gap |
| [orientationModeValues.go](../../internal/registry/orientationModeValues.go) | 0 | Protected game values; no coverage work requested |

### All partially covered files (52)

“Gap” here records unexecuted statements, not an unproved bug. Test plans must target reachable public behavior and respect the observation registry. All other instrumented files are 100% in this run; that does not prove their contracts correct, as §1 demonstrates.

| File | % |
| --- | ---: |
| [spells.go](../../app/gui/constants/spells.go) | 98.5 |
| [fileExplorerDialog.go](../../app/gui/dialogs/fileExplorerDialog.go) | 24.1 |
| [fileExplorerDialogEntries.go](../../app/gui/dialogs/fileExplorerDialogEntries.go) | 29.0 |
| [fileExplorerDialogModes.go](../../app/gui/dialogs/fileExplorerDialogModes.go) | 58.8 |
| [dialogHost.go](../../app/gui/drivers/dialogHost.go) | 2.1 |
| [state.go](../../app/gui/drivers/state.go) | 89.1 |
| [stateFiles.go](../../app/gui/drivers/stateFiles.go) | 36.5 |
| [stateGeneration.go](../../app/gui/drivers/stateGeneration.go) | 92.2 |
| [stateManualEdits.go](../../app/gui/drivers/stateManualEdits.go) | 82.1 |
| [buttonPositionLogger.go](../../app/gui/utils/buttonPositionLogger.go) | 78.9 |
| [math.go](../../app/gui/utils/math.go) | 90.0 |
| [previewGeneratorProvider.go](../../internal/composition/previewGeneratorProvider.go) | 60.0 |
| [wire_gen.go](../../internal/composition/wire_gen.go) | 92.3 |
| [gameRules.go](../../internal/entities/template_entity/template_rule_entity/gameRules.go) | 90.9 |
| [winConditions.go](../../internal/entities/template_entity/template_rule_entity/winConditions.go) | 80.8 |
| [contentRuleHandler.go](../../internal/handlers/contentRuleHandler.go) | 94.6 |
| [regenerationHandler.go](../../internal/handlers/regenerationHandler.go) | 92.9 |
| [bonusEntry.go](../../internal/helpers/config_helpers/bonusEntry.go) | 80.0 |
| [zoneContentRow.go](../../internal/helpers/editor_state_helpers/zoneContentRow.go) | 87.5 |
| [io.go](../../internal/helpers/io.go) | 17.9 |
| [string_windows.go](../../internal/helpers/string_windows.go) | 91.7 |
| [zone.go](../../internal/models/template_model/template_variant_model/zone.go) | 92.0 |
| [contentIncludeListValues.go](../../internal/registry/contentIncludeListValues.go) | 46.4 |
| [guardedContentPoolValues.go](../../internal/registry/guardedContentPoolValues.go) | 83.3 |
| [mapObjectArtifactValues.go](../../internal/registry/mapObjectArtifactValues.go) | 97.1 |
| [spellSidValues.go](../../internal/registry/spellSidValues.go) | 83.3 |
| [unguardedContentPoolValues.go](../../internal/registry/unguardedContentPoolValues.go) | 83.3 |
| [atomicFileWriter.go](../../internal/repositories/atomicFileWriter.go) | 95.0 |
| [assetProvider.go](../../internal/services/asset_provider/assetProvider.go) | 95.5 |
| [manualReapplyService.go](../../internal/services/connection_editor/manualReapplyService.go) | 88.9 |
| [zoneEditorGeometryService.go](../../internal/services/connection_editor/zoneEditorGeometryService.go) | 97.6 |
| [zoneEditorService.go](../../internal/services/connection_editor/zoneEditorService.go) | 98.4 |
| [directoryBrowserService.go](../../internal/services/file_system/directoryBrowserService.go) | 96.4 |
| [hiddenAttribute_windows.go](../../internal/services/file_system/hiddenAttribute_windows.go) | 80.0 |
| [pathResolutionService.go](../../internal/services/file_system/pathResolutionService.go) | 88.1 |
| [assetFitter.go](../../internal/services/preview_service/assetFitter.go) | 86.4 |
| [layoutBalancedRings.go](../../internal/services/preview_service/layoutBalancedRings.go) | 98.6 |
| [layoutGeometry.go](../../internal/services/preview_service/layoutGeometry.go) | 98.3 |
| [layoutRingHub.go](../../internal/services/preview_service/layoutRingHub.go) | 99.2 |
| [layoutScatter.go](../../internal/services/preview_service/layoutScatter.go) | 94.6 |
| [previewGeneratorService.go](../../internal/services/preview_service/previewGeneratorService.go) | 95.6 |
| [previewLayoutService.go](../../internal/services/preview_service/previewLayoutService.go) | 99.2 |
| [gameRulesProvider.go](../../internal/services/template_generator/providers/gameRulesProvider.go) | 99.3 |
| [mandatoryContentProvider.go](../../internal/services/template_generator/providers/mandatoryContentProvider.go) | 98.8 |
| [topologyConnectionService.go](../../internal/services/template_generator/providers/topology/base/topologyConnectionService.go) | 95.1 |
| [crossTopology.go](../../internal/services/template_generator/providers/topology/crossTopology.go) | 98.5 |
| [fractalTopology.go](../../internal/services/template_generator/providers/topology/fractalTopology.go) | 99.3 |
| [geometricHubLayout.go](../../internal/services/template_generator/providers/topology/geometricHubLayout.go) | 98.7 |
| [geometryHelpers.go](../../internal/services/template_generator/providers/topology/geometryHelpers.go) | 94.4 |
| [balancedClusterService.go](../../internal/services/template_generator/providers/topology/tournament_variant/balancedClusterService.go) | 99.5 |
| [webTopology.go](../../internal/services/template_generator/providers/topology/webTopology.go) | 98.8 |
| [zoneLabelProvider.go](../../internal/services/zones/zoneLabelProvider.go) | 90.4 |

Protected entity/registry rows above are informational only. Generated Wire code must be regenerated, never hand-edited. Dead plural converters are already dispositioned. Host/embedded-asset/disk-fault gaps retain the explicit observation limitations. Lower GUI percentages do not supersede the integration-only policy.

**Final review baseline:** Windows Go **1.27.0**; build/vet/default/tagged integration/layout checks **pass**; unit coverage **74.4%**; full configured lint **0 issues**; selected repeated suites **pass 20 times**; module tidy dry-runs **fail only on Windows checksum EOLs**; local vulnerability/race/GUI/Linux outcomes **not measured**. Subsequent backlog curation is documentation-only and does not constitute a fresh application verification run.
