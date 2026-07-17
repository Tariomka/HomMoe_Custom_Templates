# Code Quality Review

Reviewed: 2026-07-16

## Scope And Method

This is a read-only, adversarial review of the production Go implementation,
its direct callers, and the focused unit/integration tests. `data/`,
`internal/entities/template/`, and `internal/registry/` were inspected only to
understand contracts; this review proposes no edits to those protected paths.

Each finding below has a concrete correctness, maintenance, or measurable
runtime benefit. The review deliberately does not recommend replacing the
project's `linq` helpers with ad-hoc loops or replacing the entity builders with
struct literals: both patterns improve local clarity and establish useful
contracts. Existing open items in `todo/review-gpt-5.6-sol-07-13.md` (such as
atomic output writes and preview-layout caching) are not duplicated here.

Priority meaning:

- P0: generated output or a user workflow can be incorrect now.
- P1: duplicated domain logic can produce inconsistent behaviour as features
  evolve, or a common user workflow is broken.
- P2: bounded simplification, robustness, or test coverage improvement.

## Findings

### CQ-01 - Lowest-tier mandatory content does not lift content limits - DONE

**Priority:** P0 - correctness

**Evidence:** `ContentLimitProvider.CreateContentCountLimits` in
`internal/services/template_generator/providers/contentLimitProvider.go`
tallies Player, Low, Medium, High, and Hub mandatory-content rows, but omits
`settings.LowestNeutralMandatoryContent`. The Lowest tier is otherwise fully
supported: `GeneratorConfigMapper` writes it, and `MandatoryContentProvider`
emits both planned and manually rebuilt Lowest-tier content.

Consequently, a Lowest-tier configuration that requests more copies of a SID
than its default content limit can generate a template whose mandatory-content
group requests the copies but whose `content_limits_side_*` definitions still
cap the SID below that amount. The game may silently discard the excess.

**Exact fix:** Add the following tally directly after the player tally, keeping
the existing tier order:

```go
tally(settings.LowestNeutralMandatoryContent)
```

**Why this is net-positive:** this is a one-line correction that restores the
same invariant for all six configured content lists: the generated limit must
be at least the configured mandatory count.

**Tests:** add `TestWhenLowestTierRequestsMoreThanDefaultCap_LiftsLimit` in
`test/unit/internal/services/template_generator/providers/contentLimitProvider/createContentCountLimits_test.go`.
Use an existing default-capped SID, configure two or more Lowest rows, and
assert the matching output limit is raised. Keep the existing Player/Low/Medium
tests as regression coverage for all tally sources.

### CQ-02 - Three incompatible implementations infer a zone's quality from pool SIDs

**Priority:** P1 - domain correctness and deduplication

**Evidence:** the same domain fact is parsed in three places:

- `neutralZone.GetQualityFrom` in
  `internal/models/neutralZone/neutralZoneQuality.go` reads only the first
  guarded and resource pools, is case-sensitive, and defaults to Medium.
- `preview_service.ClassifyZoneTier` in
  `internal/services/preview_service/zoneClassifier.go` scans every guarded
  and unguarded pool case-insensitively, detects the rich-treasure Highest
  profile, then uses imported-template layout/name fallbacks and defaults to
  Bronze.
- `connection_editor.GetZoneTier` in
  `internal/services/connection_editor/connectionEditor.go` reads only the
  first guarded pool, is case-sensitive, and maps the result straight into its
  three neutral guard brackets.

The implementations have already drifted with the five-tier model: the
connection editor has no direct awareness of Lowest or Highest quality, while
the preview and domain model do. A zone with a marker in a later pool, in an
unguarded pool, or with different casing can be rendered as one tier while the
zone editor and manual castle reapplication treat it as another. That produces
an internally inconsistent UI and can apply the wrong advanced castle/content
settings to imported or manually changed templates.

**Exact fix:** make `internal/models/neutralZone` the owner of pool-marker
recognition, without merging the intentionally different presentation and
guard-bracket policies.

1. Add a pure helper next to `GetQualityFrom`, for example
   `ClassifyPoolQuality(guardedPools, unguardedPools, resourcePools []string)
   (Quality, bool)`. It must lower-case inputs, scan all supplied pools, return
   the strongest detected tier, and return Highest only for a tier-five marker
   combined with the rich-treasure resource marker.
2. Retain `GetQualityFrom(entities.Zone) Quality` as the compatibility wrapper:
   call the helper for a Zone and preserve its existing Medium fallback when no
   marker is present.
3. In `ClassifyZoneTier`, call the helper for the pool-based branch and map
   `Quality` to `preview.ZoneTier`. Retain the Spawn check and its layout/name
   fallbacks: those are preview-specific compatibility behaviour for incomplete
   external templates, not domain-quality inference.
4. In `GetZoneTier`, retain the player and Hub name rules, then map the shared
   `Quality` to the existing guard brackets:
   `Lowest/Low -> Bronze`, `Medium -> Silver`, `High/Highest -> Gold`.
   Do not merge `ZoneTier` with `neutralZone.Quality`; the former intentionally
   represents guard-strength bands, including Player-to-Player.

**Why this is net-positive:** one parser owns the tier grammar, so future tier
additions or SID changes cannot leave generation, preview, manual editing, and
castle propagation disagreeing. The mapping layers remain explicit where their
meaning intentionally differs.

**Tests:**

- Extend `test/unit/internal/models/neutralZone/neutralZoneQuality/getQualityFrom_test.go`
  with later-pool, unguarded-pool, uppercase, mixed-tier, rich-treasure, and
  no-marker cases.
- Keep the current fallback contracts explicit: `GetQualityFrom` without a
  marker remains Medium; preview fallback without a pool hint remains Bronze;
  an unknown neutral zone in the connection editor remains Silver.
- Extend
  `test/unit/internal/services/preview_service/zoneClassifier/classifyZoneTier_test.go`
  with parity cases for pool-tagged zones.
- Keep `getZoneTier_test.go` and `higherTierOf_test.go` under
  `test/unit/internal/services/connection_editor/connectionEditor/` green to
  prove the guard-bracket mapping did not change accidentally.

### CQ-03 - Mandatory-content tier selection is duplicated in the same provider

**Priority:** P2 - maintenance and future correctness

**Evidence:** `MandatoryContentProvider.CreateContents` contains an inline
`switch zone.Quality` selecting the configured Lowest/Low/Medium/High lists.
The same mapping appears again in `neutralRowsForQuality` in
`internal/services/template_generator/providers/mandatoryContentProvider.go`,
which is used by `CreateContentsForZones`. Both intentionally collapse Highest
onto the High configuration. CQ-01 demonstrates the risk of maintaining tier
lists in several places: the newest tier was missed in a nearby tally.

**Exact fix:** replace the inline switch with:

```go
content := cloneContentItems(neutralRowsForQuality(configuration, zone.Quality))
```

Keep `neutralRowsForQuality` as the single explicit policy for mapping domain
quality to configurable mandatory-content rows.

**Why this is net-positive:** it removes one tier-to-content decision table
without adding an abstraction or changing output. Future tier additions have
one mapping to update.

**Tests:** existing `createContents_test.go` and
`createContentsForZones_test.go` under the MandatoryContentProvider test
directory cover both paths. Add no test-only seam; run those focused tests to
prove byte-identical output.

### CQ-04 - Manually added neutral zones create orphan mandatory-content groups

**Priority:** P1 - correctness decision required

**Evidence:** `connection_editor.NewDefaultNeutralZone` in
`internal/services/connection_editor/zoneEditor.go` creates a zone using the
same generator builder, then clears `zone.MandatoryContent`. Later,
`MandatoryContentProvider.CreateContentsForZones` generates a
`mandatory_content_neutral_<label>` group for every `Neutral-*` zone,
regardless of whether that zone references the group. The generated group is
therefore not reachable from a manually added zone, and configured tier
content is not placed there. At the same time, the output carries an unused
group.

**Exact fix:** choose and document one coherent policy:

1. **Recommended:** manually added neutral zones should receive the configured
   tier content. Set `zone.MandatoryContent` to
   `entities.StringList{"mandatory_content_neutral_" + label}` in
   `NewDefaultNeutralZone`, matching generator-owned neutral zones.
2. **Alternative:** manual zones intentionally have no mandatory content. In
   that case, make `CreateContentsForZones` skip neutral zones with no matching
   mandatory-content reference, avoiding unreachable JSON.

The first option matches the expectation that an added zone behaves like a
generated zone; the second is valid only if the UI deliberately exposes a
content-free manual zone.

**Why this is net-positive:** either choice eliminates a contradictory model
where data is generated but cannot be consumed. The recommended option also
keeps manual zones inside the configuration's content policy.

**Tests:** update
`test/unit/internal/services/connection_editor/zoneEditor/newDefaultNeutralZone_test.go`
and
`test/unit/internal/services/template_generator/providers/mandatoryContentProvider/createContentsForZones_test.go`.
For the recommended policy, assert both the zone reference and matching group
are present after `GUIHandler.UpdateTemplate`; for the alternative, assert no
orphan group is emitted.

### CQ-05 - Generator/editor duplicate zone labels, castle counting, and castle-road construction

**Priority:** P2 - targeted deduplication

**Evidence:**

- `connection_editor.zoneLabels` duplicates the complete `A` through `AF`
  label pool from `services/zones/zoneLabelProvider.go`. Its own comment says
  it mirrors the generator pool.
- `connection_editor.CountZoneCastles` and the provider-local
  `countCityMainObjects` both count `MainObject.Type == "City"` using the same
  `strings.EqualFold` loop.
- `connection_editor.buildCastleRoads` duplicates the generator's stone-road
  loop from `topology/base.TopologyBase.createOuterZoneRoads`. This algorithm
  previously had an off-by-one repair, making a second copy an avoidable drift
  risk.

**Exact fix:**

1. Add a clone-returning `zones.AllZoneLabels() []string` around the existing
   generator pool and use it in `NextFreeZoneLabel`. Return a copy so callers
   cannot mutate generator state.
2. Move City-object counting to a small pure helper in the neutral-zone domain
   package, for example `neutralZone.CountCastles(entities.Zone) int`. Replace
   both current implementations. Do not make generation depend on
   `connection_editor` just to reuse its helper.
3. Extract a narrowly named exported generator helper such as
   `base.BuildCastleRoads(mainObjectCount int) []entities.Road`, and call it
   from both the generator's outer-road builder and the connection editor.
   At the same time, replace the editor's duplicate literal TypedRef type names
   with the existing canonical road-connection values already used by the
   generator.

**Why this is net-positive:** these are not generic utilities. Each represents
one shared game invariant that must remain identical across generation and
manual editing. Centralization reduces the number of locations that can
reintroduce duplicate zone names, wrong castle totals, or dangling roads.

**Tests:** retain and run `nextFreeZoneLabel_test.go`,
`countZoneCastles_test.go`, `applyNeutralZoneQuality_test.go`, and
`rebuildZoneConnectionRoads_test.go` in the connection-editor test tree, plus
the topology-base road tests. Add a zones-package test that proves
`AllZoneLabels` is equivalent to the player-label prefix and mutation-safe.

### CQ-06 - Hub-zone predicate is duplicated across GUI and generation

**Priority:** P2 - domain policy deduplication

**Evidence:** `topologyUsesHubZone` in
`app/gui/panels/layoutPanelZones.go` and `usesHubZone` in
`internal/services/template_generator/providers/mandatoryContentProvider.go`
both implement exactly:

```go
topology == config.TopologyHubAndSpoke || topology == config.TopologyGeometricHub
```

The predicate controls different user-visible behaviour: showing hub controls
and producing hub mandatory content. A future topology that creates a hub can
easily be added to only one list.

**Exact fix:** define `UsesHubZone() bool` on `config.MapTopology` (or an
equivalent pure function in that model package), then replace both private
predicates. Keep the callers' presentation and content decisions local.

**Why this is net-positive:** one two-case domain rule becomes a named,
discoverable capability of a topology. This is small now and prevents a
cross-layer feature omission later.

**Tests:** add a table-driven `usesHubZone_test.go` beside the MapTopology
model tests covering all topology constants. Existing GUI/generator tests then
continue to cover their consumers indirectly.

### CQ-07 - Bonus picker applies defaults using the previous selected type

**Priority:** P1 - user-visible correctness

**Evidence:** `BonusPickerDialog.Body` in
`app/gui/dialogs/bonusPickerDialog.go` captures `presetType` before it lays out
the dropdown. `DropdownSelector.GetWidget` updates selection during layout, but
when `WasUpdated` is true the code calls `applyTypeDefaults(presetType)` with
the old value. For example, changing Gold to Wood seeds the Gold amount
(`10000`) rather than Wood's amount (`20`); changing the initial Town Portal
type directly to a resource can leave the amount empty.

**Exact fix:** in the `WasUpdated` branch, call:

```go
this.applyTypeDefaults(this.getSelectedType())
```

The pre-layout `presetType` can remain for rendering the current frame; the
next invalidated frame will render the controls for the new type.

**Why this is net-positive:** one local change corrects data entry without
altering the established immediate-mode event pattern.

**Tests:** add a focused `integration_test` dialog scenario using the existing
synthetic input infrastructure: switch Gold to Wood and assert the amount text
is `20`; switch from the initial type to Wood and assert it is not empty. If
the existing test seams cannot inspect the textbox value, add only the smallest
test-only accessor or record the Gio-bound gap in `todo/test_observations.md`.

### CQ-08 - Load and Save As bypass their own directory suggestion policy

**Priority:** P1 - common workflow and dead-code recovery

**Evidence:** `State.Load` and `State.SaveAs` in
`app/gui/drivers/stateFiles.go` use the process working directory whenever
`os.Getwd()` succeeds. `suggestDirectory`, which correctly prefers the loaded
state's directory and then the configured output directory, runs only when
`Getwd()` fails. Its final fallback invokes `Getwd()` again. In normal use,
Open and Save As therefore ignore the currently loaded project location.

**Exact fix:** replace the two `Getwd` preambles with:

```go
dir := this.suggestDirectory()
```

Leave the working-directory fallback solely inside `suggestDirectory`. In the
same small change, trim before checking the path returned by `PickOutputDir`,
so a whitespace-only value cannot be stored as an output directory.

**Why this is net-positive:** removes duplicated fallback code, makes the
existing heuristic live, and keeps save/load navigation where the user is
already working.

**Tests:** extend the existing State-files tests with current-path, output-path,
and working-directory cases. Assert the fake dialog receives the directory of
the loaded settings file first, then the configured output directory, then CWD.

### CQ-09 - File-name sanitization is not sufficient for Windows paths

**Priority:** P1 - cross-platform correctness

**Evidence:** `helpers.SanitizeFilename` replaces visible separators and
wildcards but leaves Windows reserved device basenames (`CON`, `PRN`, `AUX`,
`NUL`, `COM1` through `COM9`, `LPT1` through `LPT9`), trailing periods/spaces,
and control runes. `file_service.SaveTemplate` and `SavePreviewImage` use the
result as a user-controlled filename. A template named `con`, `My map.`, or
containing a control character consequently fails to save on Windows despite
being legal GUI text.

**Exact fix:** after replacing unsafe runes:

1. replace every rune below U+0020 with `_`;
2. apply `strings.TrimRight(result, ". ")`;
3. compare the basename before the first `.` case-insensitively against the
   reserved device names and prefix `_` on a match;
4. retain the FileService empty-name fallback.

Extract FileService's repeated `SanitizeFilename` plus `Generated_Template`
fallback into one private `safeFileName` helper while making this change.

**Why this is net-positive:** it fulfills the repository's Windows/Linux
requirement and gives template and preview saving one filename policy rather
than two copied call-site blocks.

**Tests:** extend
`test/unit/internal/helpers/string/sanitizeFilename_test.go` with reserved
names, extensions on device names, trailing dots/spaces, control runes, and an
all-trimmed input. Existing SaveTemplate and SavePreviewImage fallback tests
should continue to prove the shared fallback behaviour.

### CQ-10 - Steam path discovery can report unreadable paths as successful

**Priority:** P2 - error propagation

**Evidence:** `FindOldenEraTemplatesDir` and `getVDFFilePath` in
`internal/helpers/io.go` return an error only when `os.Stat` satisfies
`os.IsNotExist(err)`. Permission, invalid-path, and I/O errors fall through as
successful paths. `getVDFFilePath` also uses unused named return values while
shadowing its `err` variable.

**Exact fix:** replace both checks with direct error propagation:

```go
if statErr := os.Stat(path); statErr != nil {
    return "", statErr
}
```

Make `getVDFFilePath` use ordinary unnamed return values. Preserve callers'
existing fallback logic; only false success changes.

**Why this is net-positive:** a broken install location fails at the actual
filesystem operation rather than being reported later as an unrelated VDF or
glob failure. It also simplifies the function signature.

**Tests:** extend `findOldenEraTemplatesDir_test.go` with a missing resolved
install-path case that asserts `fs.ErrNotExist`. Permission-error injection is
not portable, so the production code should stay simple and platform-neutral
rather than adding a fake filesystem solely for this branch.

### CQ-11 - Picker group badges do repeated full scans on every frame

**Priority:** P2 - measurable interactive-path simplification

**Evidence:** `MultiSelectPicker.getRowWidgets` in
`app/gui/dialogs/pickerDialog.go` iterates matching entries once to emit rows,
then `appendGroup` iterates every entry again for each distinct group to count
its matching badge. For a picker with $E$ entries and $G$ groups, filtering is
$O(E + G \cdot E)$ substring checks per Gio frame instead of $O(E)$. Item and
spell pickers contain enough grouped entries for this to occur while the user
types.

**Exact fix:** make one pre-pass over `this.entries` that filters each entry and
increments `counts[entry.group]`. Then emit headers using `counts[group]` and
emit the already filtered entries. Keep ordering and `this.grouped` behaviour
unchanged.

**Why this is net-positive:** it is a small local refactor that removes
repeated work in a typing/rendering path without changing the widget API or
introducing caching invalidation complexity.

**Tests:** no new production seam is justified for a private Gio layout helper.
When picker interaction coverage is next expanded, assert header counts along
with the existing visible-row behaviour. Run the existing GUI/integration suite
after the refactor.

### CQ-12 - A few public behaviours lack dedicated unit coverage

**Priority:** P2 - regression prevention

**Evidence:** `internal/common` has no mirrored unit-test tree even though
`mapSizes.go` contains public fallback and tie-breaking logic
(`GetMapSize`, `GetNearestMapSize`, `GetMapSizes`) and `topologies.go` contains
descriptor fallback logic. In `internal/helpers/math.go`, public
`CalculatePointTowards` has a distinct coincident-point `ok == false` branch
but no dedicated test file. These are pure, cross-platform functions and do
not require Gio seams.

**Exact fix:** add the required per-file mirrored test folders and test files:

- `test/unit/internal/common/mapSizes/` for exact lookup, unknown fallback,
  experimental-size gating, and an equal-distance tie selecting the smaller
  map size.
- `test/unit/internal/common/topologies/` for valid descriptors, unknown type,
  and out-of-range index fallback.
- `test/unit/internal/helpers/math/calculatePointTowards_test.go` for
  coincident points, axis-aligned movement, diagonal movement, and intentional
  overshoot.

**Why this is net-positive:** these fallback policies are subtle user-facing
configuration rules. Tests make them explicit and bring the implementation in
line with the repository's per-public-function test convention.

## Low-Risk Cleanup Batch

These are individually small and should be bundled only when touching the
surrounding file; they do not justify standalone churn.

- Delete `app/gui/constants/roadDistances.go` after confirming its zero
  production references. The rule dialog already uses
  `content_rules.GetDistanceDisplayNames()`.
- Remove the unused `_ = i` in `app/gui/dialogs/zoneContent.go` and the stale
  commented assignment in `app/gui/models/editorState.go`.
- Correct stale comments in `internal/models/zoneContentRowSave.go` and
  `internal/services/content_rules/contentRuleManager.go`: they describe
  deprecated flat rule fields and RoadDistance normalization that no longer
  exist. Accurate comments prevent a future compatibility assumption from
  becoming a silent data-loss bug.
- Rename FileService parameters named `filepath` and `image`, which shadow
  packages imported by the same file, to `path` and `previewImage`.
- Replace `slices.EqualFunc(..., func(left, right config.BonusEntry) bool {
  return left == right })` in `EditorStateDto.EqualsIgnoringManualEdits` with
  `slices.Equal`; `BonusEntry` is comparable and the result is identical.

## Explicit Non-Findings

The following were examined and are deliberately not recommendations:

- `internal/helpers/linq`: its clarity benefit is appropriate for the current
  call sites; replacing it with inline loops would make the code less uniform
  for negligible runtime gain.
- Entity and content builders: they encode required/default state and are
  preferable to repeated large struct literals.
- Preview layout caching, non-atomic writes, and several prior GUI/editor
  findings: these remain tracked in the existing review/backlog and should be
  handled there rather than copied into another work item.
- The preview's layout/name tier heuristics and the connection editor's
  Player-to-Player bracket: they model presentation and guard policies that are
  intentionally different from neutral-zone quality. CQ-02 centralizes only
  the duplicated pool-marker grammar.

## Recommended Implementation Order

1. CQ-01 and CQ-03 together: they are small mandatory-content changes with a
   direct generated-output fix.
2. CQ-02 with focused classifier parity tests before any additional tier work.
3. Resolve the product decision in CQ-04, then implement and test it as one
   end-to-end manual-zone change.
4. CQ-07, CQ-08, and CQ-09: independent user-workflow fixes.
5. CQ-05, CQ-06, CQ-10, CQ-11, and CQ-12 as isolated maintenance PRs.

For every implementation change, follow `AGENTS.md`: run the focused tests,
the unit coverage command, `go build ./...`, `go test ./test/... -count=1`, and
the tagged integration/performance suites when editor internals are touched.
