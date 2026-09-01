# Backlog — 2026-08-11 (Claude Opus 5)

Compiled backlog in the format of [review-prompt.md](../promt_templates/review-prompt.md)
(*Output format*). Every item below was re-verified against the source on
2026-08-11 at the state described in `.agent/session-carry-forward.md` (all 46
findings of [review-opus5-08-04.md](review-opus5-08-04.md) closed except §1.4
and §1.9, working tree clean apart from the owner's own staged files).

**Sources folded into this document**

| Source | Items taken |
| --- | --- |
| [review-opus5-08-04.md](review-opus5-08-04.md) | §1.4, §1.9 (the only two left open) |
| `backlog.md` (deleted in `cd7ad10`) | all 9 items |
| [test_observations.md](test_observations.md) | the 3 gaps that became writable once `AppRunner` gained `DragTo`/`MoveTo`/`InputText` |
| [runnerHandler.go](../../test/test_helpers/integration_common/runnerHandler.go#L61-L106) | the owner's 46-line design comment, split into §5.4 and §5.5 |

**Supersession.** This document **superseded `todo/backlog.md` in full**
— every one of its items is dispositioned in §0 and restated below with fresh
evidence. The owner deleted `backlog.md` in `cd7ad10` once this document was
accepted; §0.2 is now its only record.
It does **not** supersede [test_observations.md](test_observations.md), which
stays authoritative for intentional test gaps; only the three gaps promoted to
§5 leave it (and only once those items are done).
It does **not** supersede [review-opus5-08-04.md](review-opus5-08-04.md); §1.4
and §1.9 there stay numbered as they are and are restated here as §1.1 and §1.4.

**Severity legend:**
🔴 **High** — bug / correctness / user-visible data loss ·
🟠 **Medium** — architecture, performance, latent correctness ·
🟡 **Low** — readability, hygiene, docs, test gaps ·
⚪ **Informational / deferred decision** — no work until the owner rules.

⚠ marks an item that **must not be started without an explicit owner
go-ahead**, because it edits a protected directory (AGENTS.md §2.1) or reverses
a decision the owner already made.

**Item count: 0 🔴 · 7 🟠 · 10 🟡 · 3 ⚪ (21 total).**

**✅ Completed 2026-08-12:** §1.4 (batch A) · §1.2, §1.3 and §3.3 (batch B) ·
§3.1, §3.2 and §3.4 (batch C) · **2026-08-14:** §1.1 (batch D), §5.3 (batch F) ·
**2026-08-19:** §2.3 (batch G) · §5.1 and §5.2 (batch H) · **2026-09-01:**
§2.1 and §1.5 (batch I) — **14 done, 7 open.**
Batch D spun off §1.5 (render-path clone cost), which batch I then absorbed;
batch I spun off §2.6 (entities named outside the permitted layers).

**Baselines to hold (AGENTS.md §2.3):** unit coverage **73.9 %**, floor
**72.5 %** · `golangci-lint-v2 run ./...` **0 issues** · `gofmt -l` empty ·
`go run ./cmd/testlayoutcheck .` passes · build + vet clean under both
`integration_test` and `integration_test,gui`.

---

## 0. Disposition of the source documents

### 0.1 Carried forward from [review-opus5-08-04.md](review-opus5-08-04.md) ❗

| Prior item | Status re-verified 2026-08-11 | New section |
| --- | --- | --- |
| §1.4 Editor-state copies are shallow | Still open, verbatim — no `Clone` exists on `EditorStateDto` | §1.1 |
| §1.9 Fatal window error logged to a discard handler | ✅ Fixed 2026-08-12 | §1.4 |

All other review items are `✅ FIXED`, `❌ WILL NOT FIX`, or owner-excluded and
are **not** re-reported here.

### 0.2 Disposition of every `backlog.md` item

| Backlog item | Disposition | New section |
| --- | --- | --- |
| Preview sub-pixel precision (`Vec2`) | Carried, scope decided by owner 2026-08-11 (floats to the draw call, incl. `ZoneRadius` and editor hit-testing) | §2.3 |
| Remove `[2]float64` from template entities | Carried, ⚠ protected dir | §2.4 |
| `createTopologyAdjacency` dead Chain/Ring branches | Carried as a **deferred decision**, both options restated | §6.1 |
| Untracked zone tier property on entities | Carried, ⚠ protected dir, two branches (entity field vs. `Quality`/`Profile`) | §2.2 |
| "Use common values for template generation" | Carried, scoped by owner to four concrete literal families | §3.1 – §3.4 |
| Hub / random-portal guard values are wrong | Carried, **root cause confirmed**, split into two items | §1.2, §1.3 |
| Rework `EditorStateDto` | Carried, target package decided by owner: new **non-protected** `internal/entities/editor_state/` | §2.1 |
| Move `entities/types.go` into `template/`, rename `template` → `template_entity` | Carried, ⚠ protected dir | §2.5 |
| "Save As" is really "Save To" | ✅ Done (batch E, 2026-08-11) | §4.1 |

### 0.3 Disposition of [test_observations.md](test_observations.md)

**Promoted to backlog items** (they were recorded as "future work" only because
no synthetic-pointer seam existed; [appRunner.go](../../test/test_helpers/integration_common/appRunner.go#L152-L217)
now provides `ClickAt`, `MoveTo`, `DragTo` and `InputText`):

| Observation | New section |
| --- | --- |
| Zone editor: drag-to-connect, zone drag + snapping ("need synthetic pointer events — still future work") | §5.1 |
| Zone editor: property panels' `widget.Editor` / dropdown paths | §5.2 |
| File explorer: hidden-file toggle + pointer-driven row/scroll interactions (was "owner decision — excluded"; owner re-opened it 2026-08-11) | §5.3 ✅ done |

**Left in place as accepted, intentional gaps — do not re-report:** all the
Gio-widget/panel entries (`buttonWidget`, `sliderRowWidget`, `layoutPanel*`,
`previewPanel`), the `drivers.State` dialog-callback branches, the
`topologyConnectionService` private-policy note, and every entry under
*Unreachable defensive branches* (`connectInteriorStables` `len == 0`,
`providePreviewGenerator` `err != nil`, `atomicFileWriter` `Close`/`Sync`,
`helpers/io.go` Steam discovery, `buildShiftDerangement`).

### 0.4 Disposition of the [runnerHandler.go](../../test/test_helpers/integration_common/runnerHandler.go#L61-L106) design comment

The comment is a design note, not a to-do list, and it mixes two unrelated
problems. It is split as follows and **the comment itself should shrink to a
short pointer at this document** once §5.4 lands:

| Thought in the comment | Where it went |
| --- | --- |
| Per-tab / per-dialog handlers, fluent transitions, what embeds `baseHandler` | §5.4 (d) |
| "This mask should be updated so that only unstable parts are covered" | §5.4 (c) |
| Handlers must track layout-shifting state (checkboxes, inline dropdowns, sliders) | §5.4 (e) |
| "Calculate positions in place or hardcode them, or a combination" | §5.4 (b) |
| "There will also be an issue with scrolling, if a button is not visible" | §5.4 (f) |
| "Ideally, this should not use any code…from `*_testexports.go`" | §5.4 (g) |
| "For some reason there is a difference of some of the rendered text between local and CI" | §5.5 (separate item — a rendering/tolerance defect, not framework design) |

Two things the comment does not mention were found while reading it and are
folded into §5.4 (a): the file name violates AGENTS.md §4.1, and `NewHandler`
returns an unexported type.

---

## 1. Bugs & correctness

### 1.1 ✅ DONE 🟠 Editor-state copies are shallow, so snapshots alias live slices

*(= review-opus5-08-04 §1.4, restated with re-verified line numbers.)*

**Evidence.** `EditorStateDto` holds nine slice fields —
[editorStateDto.go](../../internal/dtos/editorStateDto.go#L82-L94):

```go
	Bonuses            []config.BonusEntry `json:"bonuses"`
	PlayerZoneContentRows    []models.ZoneContentRowSave `json:"playerZoneContentRows,omitempty"`
	LowestNeutralContentRows []models.ZoneContentRowSave `json:"lowestNeutralContentRows,omitempty"`
	LowNeutralContentRows    []models.ZoneContentRowSave `json:"lowNeutralContentRows,omitempty"`
	MediumNeutralContentRows []models.ZoneContentRowSave `json:"mediumNeutralContentRows,omitempty"`
	HighNeutralContentRows   []models.ZoneContentRowSave `json:"highNeutralContentRows,omitempty"`
	HubZoneContentRows       []models.ZoneContentRowSave `json:"hubZoneContentRows,omitempty"`
	ManualZones       []editor_state_dto.ManualZoneSave       `json:"manualZones,omitempty"`
	ManualConnections []editor_state_dto.ManualConnectionSave `json:"manualConnections,omitempty"`
```

Three places copy the struct by value and treat the result as an independent
snapshot — [editorState.go](../../app/gui/models/editorState.go#L31-L33):

```go
func (this *EditorState) GetCurrentState() dtos.EditorStateDto {
	return *this.current
}
```

[editorState.go](../../app/gui/models/editorState.go#L40-L44):

```go
func (this *EditorState) SnapshotCurrentState() {
	previousState := *this.current
	this.previous = &previousState
	this.next = nil
}
```

[stateHandler.go](../../internal/handlers/stateHandler.go#L55-L70) —
`ValidateEditorState(stateDto dtos.EditorStateDto, …)` takes the DTO **by
value**, mutates it through `issue.Fix(&stateDto)` /
`normalizeInactiveNeutralCounts(&stateDto)`, and returns it inside
`dtos.EditorStateValidationDto{State: stateDto, …}`.
[editorState.go](../../app/gui/models/editorState.go#L46-L54) does the same
shallow trick again in `GetPreviousState`.

**Why it is wrong.** A struct copy duplicates slice *headers*, not backing
arrays, so `this.previous` shares element storage with `this.current`. Change
detection compares element-wise
([editorStateDto.go](../../internal/dtos/editorStateDto.go#L187-L199) →
`contentRowSlicesEqual` / `slices.Equal`), so **any in-place element write makes
the change invisible to the unsaved/regenerate machinery**: the editor would not
mark the file dirty and `AutoRegenerate` would not fire. The same aliasing leaks
live editor state out of `GetCurrentState()` to every panel and to
`previewPanel`, and out of `ValidateEditorState` to every caller.

This is **latent, not live** today: every current writer replaces the whole
slice rather than an element
([layoutPanelZones.go](../../app/gui/panels/layoutPanelZones.go#L133-L148),
[bonusesPanel.go](../../app/gui/panels/bonusesPanel.go#L248),
[editorState.go](../../app/gui/models/editorState.go#L116-L119)). One in-place edit
anywhere reintroduces the bug with no compiler or lint signal.

**Fix.** Add an explicit deep copy beside the existing methods in
[editorStateDto.go](../../internal/dtos/editorStateDto.go):

```go
// Clone returns a copy that shares no backing array with the receiver.
func (this *EditorStateDto) Clone() EditorStateDto {
	clone := *this
	clone.Bonuses = slices.Clone(this.Bonuses)
	clone.PlayerZoneContentRows = cloneContentRows(this.PlayerZoneContentRows)
	// … the five remaining content-row slices …
	clone.ManualZones = slices.Clone(this.ManualZones)
	clone.ManualConnections = slices.Clone(this.ManualConnections)
	return clone
}
```

Call it at **all three sites** (owner decision, 2026-08-11): `GetCurrentState`,
`SnapshotCurrentState`, and `ValidateEditorState` — the last one should clone on
entry so the validated DTO it hands back never aliases the caller's slices.
`GetPreviousState` gets it for free once `SnapshotCurrentState` clones, but
clone there too rather than leaving a second shallow copy on the books.

The **elements must be checked for nested slices**:
`models.ZoneContentRowSave` carries `Rules []models.ContentRuleRowSave`
([editorStateDto.go](../../internal/dtos/editorStateDto.go#L338-L346) shows rows
built with `Rules: []models.ContentRuleRowSave{…}`), so a `slices.Clone` of the
row slice is **not** enough — write a `cloneContentRows` helper that also clones
each row's `Rules`. Re-check `editor_state_dto.ManualZoneSave` /
`ManualConnectionSave`: they wrap `entities.Zone` / `entities.Connection`, which
**do** contain slices (`MandatoryContent`, placement rules, roads). ⚠ Those
element types live in the protected
[internal/entities/template/](../../internal/entities/template/) tree — do **not**
add `Clone` methods there. Clone their slice fields from the `dtos` side, or, if
that proves unreadable, stop at the top level and document the boundary in the
test.

**Tests to add.**

- `test/unit/internal/dtos/editorStateDto/clone_test.go` — one test per slice
  field: mutate element 0 of the clone **in place**, assert the original is
  unchanged. Plus one for the nested `Rules` slice.
- `test/unit/app/gui/models/editorState/snapshotCurrentState_test.go` →
  `TestWhenSnapshotTakenAndContentRowMutatedInPlace_ReportsStateChanged` — the
  regression test that locks the behaviour in.
- `test/unit/internal/handlers/stateHandler/validateEditorState/` →
  `TestWhenValidationFixesAContentRow_TheCallersSliceIsUnchanged`.

**Watch for:** `Clone` on the hot path. `GetCurrentState` is called from panel
layout code; measure with the existing benchmarks
(`BenchmarkEditorWindow_TabCycling`) before and after. If the per-frame cost is
visible, the fallback is to clone in `SnapshotCurrentState` +
`ValidateEditorState` only and make `GetCurrentState` return a pointer to an
explicitly read-only view — but take that decision to the owner first.

**✅ Resolution (batch D, 2026-08-14).** Implemented as specified, and cloned at
`GetNextState` / `OverrideState` / `SetNextState` / `UpdateCurrentState` too —
the latter three take the DTO by value and stored its address, so the caller
kept aliasing the slices.

- New `Clone` methods: [editorStateDto.go](../../internal/dtos/editorStateDto.go)
  (+ `cloneContentRows`), [zoneContentRowSave.go](../../internal/models/zoneContentRowSave.go),
  [contentRuleRowSave.go](../../internal/models/contentRuleRowSave.go) (three
  pointer fields, so `slices.Clone` of `Rules` alone was not enough),
  [manualZoneSave.go](../../internal/dtos/editor_state_dto/manualZoneSave.go) and
  [manualConnectionSave.go](../../internal/dtos/editor_state_dto/manualConnectionSave.go).
- `entities.Zone` turned out to be **17** slice/pointer fields deep
  (`MainObject.Faction *TypedRef`, `TypedRef.Args` on the three biome fields and
  both `Road` endpoints). Per §2.1 no `Clone` was added under
  `internal/entities/` — the fields are copied from the `dtos` side by
  `cloneZone` / `cloneMainObject` / `cloneRoad` / `cloneTypedRef`.
- `PlacementRule.Args []any` is cloned as a slice only; the elements are boxed
  scalars from JSON decoding. Documented in code at the boundary.
- New shared helper [pointer.go](../../internal/helpers/pointer.go) `ClonePointer`.
- `clone_test.go` carries a **recursive reflection drift guard** that walks the
  whole tree and fails if any clone shares a backing array or pointer, so a
  field added to the protected `entities` types cannot silently go uncloned.
  Verified to fail by removing a clone line.
- The test-local `deepCloneEditorState` in
  [equalsIgnoringManualEdits_test.go](../../test/unit/internal/dtos/editorStateDto/equalsIgnoringManualEdits_test.go)
  was replaced by the production `Clone`; it had silently omitted
  `LowestNeutralContentRows`.
- **Benchmark (`TabCycling`, 6×50x, steady state).** 2.88 M → 3.01 M ns/op
  (+4.6 %), 4,676 → 6,640 allocs/op (+42 %). Owner accepted the cost. The
  read-only-view fallback above was **not** taken; the residual is tracked as
  §1.5.

---

### 1.2 ✅ DONE 🟠 Hub-touching connections are guarded as *player borders*, not as hub borders

**Resolved 2026-08-12 (batch B).** `constants.IsHubLabel`
([hubLabel.go](../../internal/common/constants/hubLabel.go)) plus a private
`labelQuality` helper collapsed `GetBorderGuardValue` to
`max(labelQuality(a), labelQuality(b)).GetGuardValue()`; `hubTopology`,
`geometricHubTopology` and `hubClusterService` now pass the hub label instead of
a player anchor, so every hub-touching connection is `35 000`. The hub was not
added to `neutral_zone.Plans`. Player/plan pairs are bit-identical to before.
The golden generator tests needed **no** expectation change.

**Evidence.** The hub zone is built with the top tier —
[zoneFactory.go](../../internal/services/zones/zoneFactory.go#L120-L141):

```go
		Profile: common_zones.GetNeutralZoneProfile(neutral_zone.QualityHighest)
```

but the hub is deliberately **not** in `neutral_zone.Plans`, and every hub
topology therefore substitutes a *player* label as the guard anchor.

[hubTopology.go](../../internal/services/template_generator/providers/topology/hubTopology.go#L99-L108):

```go
		hubAnchor := label
		if len(playerLabels) > 0 {
			hubAnchor = playerLabels[0]
		}
		hubGuard := this.GetBorderGuardValue(hubAnchor, label, playerLabels, neutralZones, tuning)
```

[geometricHubTopology.go](../../internal/services/template_generator/providers/topology/geometricHubTopology.go#L143-L158):

```go
	guardAnchor := playerLabels[0]
	for _, label := range layout.hubPortalLabels {
		guardLabel := label
		if slices.Contains(playerLabels, label) {
			guardLabel = guardAnchor
		}
		…
			WithGuardValue(this.GetBorderGuardValue(guardAnchor, guardLabel, playerLabels, neutralZones, tuning)).
```

[hubClusterService.go](../../internal/services/template_generator/providers/topology/tournament_variant/hubClusterService.go#L94):

```go
			WithGuardValue(this.GetBorderGuardValue(playerLabel, spokeLabel, []string{playerLabel}, allNeutralZonePlans, tuning)).
```

`GetBorderGuardValue`
([topologyConnectionService.go](../../internal/services/template_generator/providers/topology/base/topologyConnectionService.go#L180-L201))
knows only *player* and *plan* labels, so it can never see the hub.

**Why it is wrong.** With both endpoints resolving to player labels the function
returns `QualityUnknown.GetGuardValue()` = **30 000**
([neutralZoneQuality.go](../../internal/models/neutral_zone/neutralZoneQuality.go#L14-L30));
hub-to-neutral returns the *neutral's* tier (as low as 10 000). The hub is the
richest zone on the map (`QualityHighest` = **35 000**), so its borders are
systematically under-guarded — players reach the hub's Platinum content behind a
Bronze/player-grade guard. The bug is invisible in tests because
[getBorderGuardValue_test.go](../../test/unit/internal/services/template_generator/providers/topology/base/topologyBase/getBorderGuardValue_test.go#L12-L112)
only exercises player and plan labels.

**Fix** (rule confirmed by the owner 2026-08-11: **hub edges = `max(Highest,
other endpoint)`, i.e. always 35 000**).

The three existing branches of `GetBorderGuardValue` are already *exactly*
`max` over per-label qualities, because `QualityUnknown` is `-1` and every real
tier is `0…4`. Collapse them and add the hub as a third label kind:

```go
// in internal/common/constants (beside zoneNames.go)
func IsHubLabel(label string) bool {
	return label == HubZoneName || strings.HasPrefix(label, HubZonePrefix)
}
```

```go
// in topologyConnectionService.go
func labelQuality(label string, playerLabels []string, neutralZones neutral_zone.Plans) neutral_zone.Quality {
	if constants.IsHubLabel(label) {
		return neutral_zone.QualityHighest
	}
	if slices.Contains(playerLabels, label) {
		return neutral_zone.QualityUnknown
	}
	return neutralZones.GetQuality(label)
}

func (this *TopologyConnectionService) GetBorderGuardValue(
	labelA, labelB string, playerLabels []string,
	neutralZones neutral_zone.Plans, tuning models.GenerationTuning,
) int {
	higher := max(
		labelQuality(labelA, playerLabels, neutralZones),
		labelQuality(labelB, playerLabels, neutralZones))
	return tuning.ScaleByBorderGuardStrength(higher.GetGuardValue())
}
```

`IsHubLabel` must match both `constants.HubZoneName` (`"Hub"`, used by
`hubTopology` and `geometricHubTopology`) and the tournament form
`constants.HubZonePrefix + playerLabel` (`"Hub-A"`,
[hubClusterService.go](../../internal/services/template_generator/providers/topology/tournament_variant/hubClusterService.go#L39)).
Zone labels are `A…AF`
([zoneLabels.go](../../internal/common/constants/zoneLabels.go)), so the prefix
cannot collide with a real label.

Then fix the three callers to pass the hub, not a player:

- [hubTopology.go](../../internal/services/template_generator/providers/topology/hubTopology.go#L99-L108) — delete `hubAnchor`, call `GetBorderGuardValue(constants.HubZoneName, label, …)`.
- [geometricHubTopology.go](../../internal/services/template_generator/providers/topology/geometricHubTopology.go#L143-L160) — delete `guardAnchor`/`guardLabel`, call `GetBorderGuardValue(constants.HubZoneName, label, …)`.
- [hubClusterService.go](../../internal/services/template_generator/providers/topology/tournament_variant/hubClusterService.go#L94) — pass `hubName` as the first argument instead of `playerLabel`.

**Explicitly do NOT** add the hub as an ordinary entry in `neutral_zone.Plans` —
the plans list drives neutral-zone *generation* (counts, castles, profiles), and
injecting the hub there would change zone generation, not just guards.

**Tests to add.**

- `test/unit/internal/services/template_generator/providers/topology/base/topologyBase/getBorderGuardValue_test.go`
  — extend with `TestWhenOneLabelIsTheHub_ReturnsTheHighestQualityGuardValue`,
  `TestWhenOneLabelIsATournamentHub_ReturnsTheHighestQualityGuardValue`, and a
  regression pair asserting the existing player/player (30 000) and
  player/neutral outcomes are byte-identical after the `max` collapse.
- `…/topology/hubTopology/createTopologyVariant_test.go`,
  `…/topology/geometricHubTopology/createTopologyVariant_test.go`,
  `…/tournament_variant/hubClusterService/createClusterVariant_test.go` — one
  test each: every connection whose `From`/`To` is the hub has guard value
  `35 000` (unscaled tuning).

**Verify no regression** in the golden generator tests
(`test/unit/internal/services/template_generator/templateGenerator/generate_test.go`)
— guard values there will legitimately change for hub topologies; update the
expectations deliberately, do not loosen the assertions.

---

### 1.3 ✅ DONE 🟠 Random portal guards ignore endpoint tiers (flat 25 000)

**Resolved 2026-08-12 (batch B).** `CreateRandomPortalConnections` gained a
trailing `neutralZones neutral_zone.Plans` parameter (interface, `TopologyBase`
pass-through and all seven call sites) and now calls `GetBorderGuardValue`, so
portals are guarded by the higher endpoint tier and scaled exactly once. The
`25 000` literal is gone, which closes §3.3 as well.

**Evidence.**
[topologyConnectionService.go](../../internal/services/template_generator/providers/topology/base/topologyConnectionService.go#L30-L68):

```go
func (this *TopologyConnectionService) CreateRandomPortalConnections(
	playerLabels, orderedLabels []string,
	tuning models.GenerationTuning,
	maxCount int,
) []entities.Connection {
	…
			WithGuardValue(tuning.ScaleByBorderGuardStrength(25000)).
```

The signature does not even accept `neutral_zone.Plans`, so tier information is
not available at the call. The current unit test pins the literal —
[createRandomPortalConnections_test.go](../../test/unit/internal/services/template_generator/providers/topology/base/topologyBase/createRandomPortalConnections_test.go#L53-L75).

**Why it is wrong.** A portal that drops a player straight into a Platinum
neutral zone is guarded at 25 000 (the `QualityHigh` value), while the direct
land border into the same zone is guarded at 35 000 — the portal is the cheap
back door. Symmetrically, a portal into a Plastic zone is over-guarded
(25 000 vs. 10 000). Owner decision 2026-08-11: portals use **the same
`max(endpoint qualities)` rule as direct borders**.

**Fix.** Thread the plans through and reuse the (now hub-aware, §1.2)
`GetBorderGuardValue`:

1. Add `neutralZones neutral_zone.Plans` to `CreateRandomPortalConnections` in
   [topologyConnectionServiceInterface.go](../../internal/services/template_generator/providers/topology/base/topologyConnectionServiceInterface.go#L10),
   [topologyConnectionService.go](../../internal/services/template_generator/providers/topology/base/topologyConnectionService.go#L30-L34)
   and the pass-through in
   [topologyBase.go](../../internal/services/template_generator/providers/topology/base/topologyBase.go#L153-L158).
2. Replace the literal with
   `WithGuardValue(this.GetBorderGuardValue(fromLabel, toLabel, playerLabels, neutralZones, tuning))`.
   **Do not double-scale** — `GetBorderGuardValue` already applies
   `tuning.ScaleByBorderGuardStrength`.
3. Update the seven call sites, each of which already has the plans in scope:
   [chainTopology.go](../../internal/services/template_generator/providers/topology/chainTopology.go#L47),
   [geometricHubTopology.go](../../internal/services/template_generator/providers/topology/geometricHubTopology.go#L57),
   [hubTopology.go](../../internal/services/template_generator/providers/topology/hubTopology.go#L46),
   [positionedTopologyBuilder.go](../../internal/services/template_generator/providers/topology/positionedTopologyBuilder.go#L57),
   [ringTopology.go](../../internal/services/template_generator/providers/topology/ringTopology.go#L45),
   [tournamentTopology.go](../../internal/services/template_generator/providers/topology/tournamentTopology.go#L71),
   [webTopology.go](../../internal/services/template_generator/providers/topology/webTopology.go#L51).

This removes the `25000` literal, so §3.3 disappears with it.

**Verified non-issue (do not "fix"):** random portals can never touch the hub —
`hubTopology` passes only `outerLabels`
([hubTopology.go](../../internal/services/template_generator/providers/topology/hubTopology.go#L39-L48))
and `geometricHubTopology` passes players + plan labels
([geometricHubTopology.go](../../internal/services/template_generator/providers/topology/geometricHubTopology.go#L46-L59)).
Hub portals in the geometric layout are created separately and are covered by
§1.2.

**Tests to add / update.**

- Rewrite the guard expectations in
  [createRandomPortalConnections_test.go](../../test/unit/internal/services/template_generator/providers/topology/base/topologyBase/createRandomPortalConnections_test.go)
  →  `TestWhenPortalJoinsTwoNeutralZones_UsesTheHigherTierGuardValue`,
  `TestWhenPortalJoinsTwoPlayerZones_UsesThePlayerBorderGuardValue`.
- Golden generator expectations will move for every topology with
  `RandomPortals` enabled — update deliberately.

**Ordering:** land §1.2 first; §1.3 depends on the collapsed
`GetBorderGuardValue`.

---

### 1.4 ✅ DONE 🟡 A fatal window error is logged to a discard handler, then the process exits silently

*(= review-opus5-08-04 §1.9, restated with re-verified line numbers.)*

**Resolved 2026-08-12 (batch A).** The `app.DestroyEvent` error branch in
[program.go](../../app/gui/program.go) now writes to `os.Stderr` before the
discarded `slog.Error`; the optional `os.Exit`/`main.go` refactor was **not**
done. `app/gui/program.go` was added to the Gio-UI section of
[test_observations.md](test_observations.md).

**Evidence.** [program.go](../../app/gui/program.go#L34-L41):

```go
	for {
		switch event := window.Event().(type) {
		case app.DestroyEvent:
			if event.Err != nil {
				slog.Error("Window destroyed with error", slog.String("error", event.Err.Error()))
				os.Exit(1)
			}
			os.Exit(0)
```

but `getAndConfigureWindow` — called at
[program.go](../../app/gui/program.go#L26), i.e. *before* the loop — installs a
discard logger at [program.go](../../app/gui/program.go#L58):

```go
	slog.SetDefault(slog.New(slog.DiscardHandler))
```

Logging is only re-enabled by the opt-in `-with-logging` flag
([program.go](../../app/gui/program.go#L69-L72)).

**Why it is wrong.** When Gio fails to create or maintain the window (no GPU,
missing X/Wayland libraries on Linux, driver fault), the user sees the app
vanish with exit code 1 and **zero diagnostic output**. The failure is
unreportable — exactly the first-run experience a Linux user hits.

**Fix.** Write the fatal path to stderr unconditionally, independent of the
configured `slog` default:

```go
		case app.DestroyEvent:
			if event.Err != nil {
				fmt.Fprintln(os.Stderr, "Window destroyed with error:", event.Err)
				slog.Error("Window destroyed with error", slog.String("error", event.Err.Error()))
				os.Exit(1)
			}
```

`depguard` denies `log` in non-main files
([.golangci.yml](../../.golangci.yml)) but not `fmt`, so this stays lint-clean.
Secondary, **optional**: `os.Exit` inside the loop skips deferred cleanup —
harmless today (nothing is deferred), but returning an error to
[main.go](../../main.go) would make the bootstrap testable. Do not do that as part
of this item unless the owner asks.

**Tests.** No unit test is practical for the Gio event loop. Add an entry for
`app/gui/program.go` to the Gio-UI section of
[test_observations.md](test_observations.md) recording the gap.

---

### 1.5 ✅ FIXED 🟡 Panels deep-clone the whole editor state on every frame

*(Spun off from §1.1, batch D, 2026-08-14. Fixed in batch I phase 6, 2026-08-31.)*

**What was wrong.** `GetCurrentState` deep-clones since batch D, and the render
path called it several times a frame. Batch D's clone-free scalar readers
(`GetTemplateName` / `GetMapSize` / `GetTopology` / `GetExperimentalMapSizes` on
[editorState.go](../../app/gui/models/editorState.go) and
[state.go](../../app/gui/drivers/state.go)) converted eight single-field call
sites but recovered only ~13 % of the added allocations, because the whole-state
readers dominated and no scalar getter reaches them.

**What the profile actually said.** The original write-up blamed five `Layout`
paths in the panels. That was wrong: four of the five are `LoadFromState`, which
runs from the panel constructors and from `Window.load()` (toolbar new/open),
never per frame. An `alloc_objects` profile of
`BenchmarkEditorWindow_TabCycling` put **75.2 % of every allocation in the
benchmark** inside `editor_state_model.EditorState.Clone` — 97 % of that in
`CloneZoneContentRows` — reached from `UpdateCurrentState` (37.4 %),
`stateHandler.ValidateEditorState` (29.3 %), `GetCurrentState` (24.1 %) and
`GetPreviousState` (9.2 %).

The cost was the clone *mechanism*, not the clone count. Every row-slice clone
ran `linq.FromSlice(x).Select(f).ToSlice()`, a lazy chain that allocates its
three closures and boxes `ToSlice`'s accumulator **before it looks at the
source**, then regrows the result with `append`. `EditorState.Clone` runs eight
such chains and six are empty on a default state — including the clone of the
embedded entity's `Rules`, which the model deliberately keeps nil. Roughly a
quarter of all allocations in the benchmark were chains projecting nothing.

**What was done.**

1. Added `linq.SelectSlice` ([slice.go](../../internal/helpers/linq/slice.go)) —
   the eager equivalent of `FromSlice(...).Select(...).ToSlice()`. Returns `nil`
   for an empty source without allocating and sizes the result exactly once;
   semantics match the lazy chain including empty → `nil`.
2. Pointed the five clone helpers on the frame path at it:
   `CloneContentRuleRows` / `CloneZoneContentRows` in both
   [editor_state_helpers](../../internal/helpers/editor_state_helpers/) and
   [editor_state_model](../../internal/models/editor_state_model/), plus the
   `ManualZones` / `ManualConnections` chains in `EditorState.Clone`.
3. Stopped `handleZoneContentDialogClicks`
   ([layoutPanelZones.go](../../app/gui/panels/layoutPanelZones.go)) cloning the
   whole state per frame. It cloned to read six row slices that only matter when
   one of six buttons was clicked; `openZoneContentDialog` now takes a getter and
   reads the state itself. The `switch` and its short-circuit over `Clicked(gtx)`
   are untouched — a click left unread this frame must stay readable next frame.

**Result** — `BenchmarkEditorWindow_TabCycling`, headless, `-benchtime=50x
-count=6`:

| | allocs/op | B/op | ns/op (median of 6) |
| --- | --- | --- | --- |
| Before batch D | 4,676 | 1,254 KB | 2.88 M |
| Batch D + scalar accessors | 6,640 | 1,435 KB | 3.01 M |
| Before batch I phase 6 | 12,690 | 1,045 KB | 4.05 M |
| **After** | **4,773** | **720 KB** | **3.57 M** |

−62.4 % allocations, −31.1 % bytes, −11.8 % wall clock against the tree it was
measured on. (The batch-D rows come from a different tree *and* different
hardware — compare only the last two rows.) Re-profiled after: the state path
fell from ~1.72 M sampled objects to ~0.20 M (−88 %), and everything above it in
the profile is now Gio rendering.

**What was deliberately *not* done.**

- **No per-panel view structs.** Fix option 3 below assumed the panels read whole
  state every frame. They do not. The one genuine per-frame reader was the click
  handler, and the fix there is to not read at all rather than to read into a new
  type. Six view structs for zero measured gain is what AGENTS.md §3.1 forbids.
- **The clones that hand state *out* of the model stay** — `UpdateCurrentState`,
  `ValidateEditorState`, `GetCurrentState`, `GetPreviousState` and
  `AutoRegenerate`. Undoing those would undo §1.1.
- `UpdateCurrentState` → `ValidateEditorState` is a genuine double deep-clone and
  is the largest remaining item, but collapsing it means letting `updateFunc`
  mutate through to the live `current`. With each clone now ~9× cheaper the case
  for taking that risk is weak. Recorded here, not scheduled.

**The three options this item originally listed** — a borrowed read-only
accessor, a per-frame cache keyed on `GetTemplateRevision`, per-panel view
structs — were all superseded by the measurement. None was implemented.

---

## 2. Architecture & modelling

### 2.1 ✅ FIXED 🟠 `EditorStateDto` is a flat god-DTO with no entity/model layer

*(Batch I, 2026-08-21 → 2026-09-01. The doctrine that came out of it is now
**AGENTS.md §4.4.1**; read that first, this entry is the history.)*

**What was wrong.** `internal/dtos/editorStateDto.go` was a single struct
carrying identity, map, player, neutral-zone, castle, advanced tier, faction,
generation, connectivity, density, topology, victory, arena, tournament and XP
fields, plus banned content, overrides, bonuses, six mandatory-content row
slices and both manual-edit slices — then added defaults, layout comparison,
castle diffing, manual-edit-insensitive equality and their private helpers. 72
fields and ~350 lines of behaviour in one file, imported by **twelve**
production packages.

It was simultaneously the persistence schema (`.gen.json` tags), the GUI's
working state, the validator's target and the mapper's input, so adding a field
meant touching the on-disk format, the panels, the validator and the mapper at
once, and the equality/diff logic (§1.1's aliasing surface) lived in the same
file as the JSON contract. No layer owned the *meaning* of editor state
independently of how it is serialised.

**The shape that shipped.** Three layers, and the ordering matters because the
first attempt got it backwards:

- **Entity** — `internal/entities/editor_state/`, nine behaviour-free groups
  (`templateIdentity`, `mapSettings`, `playerSettings`, `neutralZoneSettings`,
  `castleSettings`, `generationSettings`, `gameRuleSettings`, `contentSettings`,
  `manualEditSettings`) plus the leaf types `BonusEntry`, `ZoneContentRow`,
  `ContentRuleRow`, `ManualZoneSave`, `ManualConnectionSave`. JSON tags only.
  `MapTopology` lives in `internal/entities/topology/`.
- **Model** — `internal/models/editor_state_model/`. **The model owns the
  structure.** Each group embeds its entity group and adds the behaviour;
  `ContentSettings` and `ManualEditSettings` cannot embed, because their slices
  must carry *model* element types and Go slices do not interconvert, so they are
  re-declared with explicit `ToXModel` / `ToXEntity` converters. `ZoneContentRow`
  embeds the entity but shadows `Rules` with `[]ContentRuleRow`, leaving the
  embedded slice nil so there is one source of truth.
- **DTO** — `internal/dtos/editor_state_dto/`. The owner's rule: *"redefinition
  is expected in Models, but it should never happen in DTOs."* So
  `EditorStateDto` is literally `struct { editor_state_model.EditorState }` and
  declares nothing else. **A DTO embedding or carrying a Model is intended** —
  `EditorStateValidationDto.State`, `CastleSettingsReapplyRequestDto.Changes` and
  `ManualEditDecisionDto.ReapplyWithCastleChanges` all do it deliberately. Do not
  "fix" them.

`app/` **may hold a Model** as its working state; only the *crossing* into
`internal/` is a DTO. Conversion happens at exactly two seams:
`internal/handlers` (DTO ⇄ Model) and `internal/repositories` (Model ⇄ Entity).

**The non-negotiable invariant held.** The on-disk `.gen.json` field names and
shape never changed. Two frozen fixtures under `test/test_helpers/testdata/` and
the untagged `editorStateWireFormat_integration_test.go` passed **unchanged**
through every phase, comparing parsed objects rather than bytes.

**Course corrections worth remembering.**

1. **Phases 1–5 were built on a false premise** — that `EditorStateDto` is the
   persisted `.gen.json` shape. It is not; that is an Entity's job. Phase 5 was
   abandoned and phases 8–12 added to correct it.
2. **The first Phase 10 inverted the layers** — it made the Model a thin wrapper
   around the entity and the DTO a full structural rewrite (nine `*Dto` group
   structs plus a DI'd DTO⇄Model mapper). The owner reversed it; that mapper and
   those structs are deleted. Do not propose them again.
3. **Two dereference bugs were found in the regeneration path.** `nil` is
   load-bearing there: a nil `Previous` means "first generation" and a nil `Next`
   means "unarmed debounce". `new(getPreviousStateDto())` always yields a
   non-nil pointer and destroyed the signal. Any refactor of those four call
   sites must preserve nil.
4. **`LoadState` stopped returning `(*EditorStateDto, []string, error)`** and now
   returns `(*editor_state_dto.EditorStateValidationDto, error)` — the envelope
   the handler was unpacking only for every caller to re-pair.

**The gate.** §0.4 was violated silently for four phases because nothing checked
it, so phase 12 added
[layering_test.go](../../test/unit/architecture/dependency/layering_test.go):
entities may not import upward (0 violations), entities may be named only by
repositories/models/entities/mappers/`*_helpers`, DTOs may be named only by
handlers/dtos/`app/`. Its two seeded allow-lists are the residual breach, tracked
as **§2.6**; they only ever shrink.

**Related:** §1.5 (the per-frame clone cost) was folded in as phase 6.

---

### 2.2 🟡 ⚠ Zone tier has no single source of truth

**Evidence.** `neutral_zone.Quality`
([neutralZoneQuality.go](../../internal/models/neutral_zone/neutralZoneQuality.go#L3-L12))
is the tier enum. During generation the tier is explicit on
`neutral_zone.Plan`
([neutralZonePlan.go](../../internal/models/neutral_zone/neutralZonePlan.go#L11-L15)),
but the finished `entities.Zone` carries **no tier**, so
[zoneClassifier.go](../../internal/services/zones/zoneClassifier.go#L23-L177)
re-derives it from layout + resource pools + guarded/unguarded pool IDs across
three content-inference branches.

Consumers that all depend on that re-derivation:
[previewLayoutService.go](../../internal/services/preview_service/previewLayoutService.go#L103-L116),
[zoneEditorZoneProps.go](../../app/gui/dialogs/zoneEditorZoneProps.go#L61-L64),
[zoneEditorService.go](../../internal/services/connection_editor/zoneEditorService.go#L187-L222),
[manualReapplyService.go](../../internal/services/connection_editor/manualReapplyService.go#L88-L124),
[mandatoryContentProvider.go](../../internal/services/template_generator/providers/mandatoryContentProvider.go#L88-L115),
[gladiatorArenaProvider.go](../../internal/services/template_generator/providers/gladiatorArenaProvider.go#L85-L111),
[zoneEditorHandler.go](../../internal/handlers/zoneEditorHandler.go#L58-L81),
[connectionEditorService.go](../../internal/services/connection_editor/connectionEditorService.go#L24-L39).

**Why it is wrong.** Tier is decided once at plan time, thrown away, then
reconstructed by pattern-matching on content. Every feature that edits a zone's
content (the manual editor's re-profile action, mandatory-content regeneration)
can silently flip the inferred tier, and the inference rules must stay in lockstep
with the profile catalogue in
[neutralZoneProfile.go](../../internal/common/common_zones/neutralZoneProfile.go#L10-L25)
with nothing enforcing that.

**Fix — two branches, owner picks:**

- **Branch A (owner's first idea, ⚠ protected dir).** Add a runtime-only
  `Quality` field to `entities.Zone`
  ([zone.go](../../internal/entities/template/template_variant/zone.go)) tagged
  `json:"-"` exactly like the existing `GeneratorPosition`, set it in
  `ZoneFactory` at creation, and reduce `ZoneClassifier` to a fallback for zones
  that arrive without it (loaded manual snapshots). **Requires explicit owner
  approval** — it edits the protected `.rmg.json` schema package, even though a
  `json:"-"` field cannot change the output file. Add a round-trip test proving
  the emitted `.rmg.json` is byte-identical before and after.
- **Branch B (owner's second idea, no protected edits).** Keep `entities.Zone`
  as is and make `neutral_zone.Profile`
  ([neutralZoneProfile.go](../../internal/models/neutral_zone/neutralZoneProfile.go#L3-L20))
  carry its own `Quality`, then have every consumer above resolve tier through a
  single `IZoneTierService` that prefers a recorded profile and falls back to
  `ZoneClassifier`. No schema change, but the tier still has to be carried
  alongside the zone through the manual-edit snapshot.

**Recommendation:** start with **B**. It is reversible, needs no protected edit,
and if it turns out the profile cannot be carried through the `.gen.json`
manual-zone snapshot, that is the concrete evidence needed to justify A.

**Tests.** Whichever branch: a test per consumer proving the tier it reads
matches the tier the generator planned, plus a manual-editor test that a
re-profiled zone reports the new tier and a saved/loaded manual zone keeps it.

---

### 2.3 ✅ Preview geometry is integer-only although a `Vec2` already exists

**Evidence.** [previewLayout.go](../../internal/models/preview/previewLayout.go#L5-L11):

```go
type Layout struct {
	Positions   map[string]image.Point
	Zones       []Zone
	Connections []Connection
	ZoneRadius  int
}
```

Positions are rounded at
[layoutGeometry.go](../../internal/services/preview_service/layoutGeometry.go#L82-L93)
(`image.Pt(int(math.Round(px[i])), int(math.Round(py[i])))`), and again
independently in
[layoutRingHub.go](../../internal/services/preview_service/layoutRingHub.go#L63-L69)
and
[layoutBalancedRings.go](../../internal/services/preview_service/layoutBalancedRings.go#L119-L129);
`canvasMetrics.center()` **truncates** rather than rounds
([layoutGeometry.go](../../internal/services/preview_service/layoutGeometry.go#L65-L67)).
The integers then propagate into the editor:
[zoneEditorGeometryService.go](../../internal/services/connection_editor/zoneEditorGeometryService.go#L73-L79)
copies them straight into `models.ZoneEditorGeometry`, and hit-testing works in
`image.Point`
([zoneEditorGeometryService.go](../../internal/services/connection_editor/zoneEditorGeometryService.go#L81-L96)).

Meanwhile a generic float vector already exists —
[vec2.go](../../internal/helpers/data/vec2.go#L8-L69) (`Vec2[T]`, `Vec2FromPoint`,
`ToPoint`, `ToPointRounded`) — and `models.Position = data.Vec2[float64]`
([position.go](../../internal/models/position.go#L5-L10)) is what the topology
builders already work in.

**Why it is wrong.** Every layout strategy quantises to whole pixels *before*
the geometry is consumed, so zone centres, connection endpoints and the Bezier
control points computed from them all inherit up to half a pixel of error per
stage — visible as jitter when the preview is resized and as off-centre edges in
the zone editor. The rounding also happens in four places with two different
rules (round vs. truncate).

**Fix** (scope confirmed by the owner 2026-08-11 — **all four**):

1. `Layout.Positions` → `map[string]models.Position` (float64).
2. `Layout.ZoneRadius` → `float64`.
3. Push floats through `models.ZoneEditorGeometry` and the hit-testing /
   edge-building in
   [zoneEditorGeometryService.go](../../internal/services/connection_editor/zoneEditorGeometryService.go)
   (`HitTestNode`, `HitTestEdge`, `buildEdges`) — they already compute in
   `float64` internally via `math.Hypot`, so this mostly deletes conversions.
4. **Round only at the final draw call**: the PNG renderer
   ([previewGeneratorService.go](../../internal/services/preview_service/previewGeneratorService.go#L51-L72))
   and the Gio canvases
   ([previewPanel.go](../../app/gui/panels/previewPanel.go#L170-L239),
   [zoneEditorCanvas.go](../../app/gui/dialogs/zoneEditorCanvas.go#L92-L207)).
   Use `Vec2.ToPointRounded()` / `f32.Pt` there — the canvas already draws
   Beziers in `f32.Point`, so step 4 in the editor is mostly *removing* the
   round-trip through `image.Point`.

**⚠ Protected boundary:** `Zone.GeneratorPosition *[2]float64`
([zone.go](../../internal/entities/template/template_variant/zone.go#L7-L23))
stays exactly as it is for this item — converting it is §2.4 and needs owner
approval. Convert at the read site (`generatorCoords`,
[layoutGeometry.go](../../internal/services/preview_service/layoutGeometry.go#L95-L106)).

**Tests.** The preview layout services already have mirrored unit folders under
`test/unit/internal/services/preview_service/` — update the expected values and
add `TestWhenTwoZonesAreLessThanAPixelApart_TheirCentresDiffer`. The GPU-gated
snapshot suite
([window_snapshot_integration_test.go](../../test/integration/gui/window_snapshot_integration_test.go))
will need `-update`; **the owner must eyeball the regenerated snapshots**, since
sub-pixel changes are exactly what that suite is meant to catch.
The numeric geometry pins in
[zoneEditorGeometry_integration_test.go](../../test/integration/gui/zoneEditorGeometry_integration_test.go#L138-L163)
will shift — update them deliberately, do not relax them to `InDelta`.

**✅ Resolved 2026-08-19 (batch G).** Implemented as specified, all four steps.
`Layout.Positions` is `map[string]data.Vec2[float64]` and `Layout.ZoneRadius` is
`float64`; `preview.Zone.Center`, `preview.Connection.Start/Ctrl/End`,
`ZoneEditorEdge.MidPoint` and `ZoneEditorSnapResult.Position` were floated with
them, since leaving any of those integer would re-quantise one step downstream
and defeat the item. `internal/models/preview` uses `data.Vec2[float64]` rather
than the `models.Position` alias because `internal/models` imports it — the two
are the same type.

Rounding now happens exactly once per output: `app/gui/utils/draw.go` for the
Gio canvases (the `Vec2[float64]` → `f32.Point` bridge is
`app/gui/utils.ToF32Point`, kept out of `internal/` so no `internal` package
imports Gio) and the pixel loop in `assetProvider` for the PNG. Both former
rules are gone — `commitPositions` no longer rounds and `canvasMetrics.center()`
no longer truncates. Pointer input is no longer truncated in
`zoneEditorCanvas.go`, so hit-testing and snapping see the true position.

`Zone.GeneratorPosition` was **not** touched; `generatorCoords` still converts at
the read site.

Two expectations changed value, both by design: a snap that returned `x = 201`
now returns `x = 200 + 6/7`, and the ring-layout zero-angle slot yields
`47.99999999999997` instead of `48`. Normalised `manualPosition` values in
`.gen.json` gain sub-pixel precision; the field is already `float64`, so old and
new files stay mutually readable.

`TestWhenTwoZonesAreLessThanAPixelApart_TheirCentresDiffer` was added, and
`helpers.CalculatePointTowards` / `GetVectorOnQuadraticBezierCurve` — which had
no tests at all — got their first ones.
`helpers.GetPointOnQuadraticBezierCurve` became dead in the process and was
removed.

**No snapshots were regenerated.** The prediction in the item above turned out
to be wrong: the GPU suite passed unchanged on the first run, because the
preview canvas is masked by the harness and the zone-editor handler takes no
snapshots. No `-update` was run and no owner eyeballing was needed. The numeric
pins in `zoneEditorGeometry_integration_test.go` were updated to exact new
literals as instructed.

Coverage flat at **72.9 %**. The batch G plan file has since been deleted, so
this entry is the record.

---

### 2.4 ⚪ ⚠ Replace `[2]float64` with `Vec2` in the template entities

**Evidence.** [zone.go](../../internal/entities/template/template_variant/zone.go#L7-L23)
declares `GeneratorPosition *[2]float64` (`json:"-"`), and producers stamp it
as an array literal —
[geometricHubTopology.go](../../internal/services/template_generator/providers/topology/geometricHubTopology.go#L104-L119):

```go
	zones[0].GeneratorPosition = &[2]float64{layoutCenter, layoutCenter}
	…
		zone.GeneratorPosition = &[2]float64{position.X, position.Y}
```

— i.e. a `models.Position` (`data.Vec2[float64]`) is unpacked into an array at
every write site and re-packed at every read site.

**Why it matters.** Purely ergonomic: `[2]float64` has no `.X`/`.Y`, no
arithmetic, and no `ToPointRounded`, so the conversion noise above is repeated
in `positionedTopologyBuilder`, `geometricHubTopology`, `balancedClusterService`
and the preview `generatorCoords` helper.

**⚠ OWNER DECISION REQUIRED — protected directory.**
[internal/entities/template/](../../internal/entities/template/) is read-only under
AGENTS.md §2.1. The field is `json:"-"`, so changing its Go type **cannot**
change the emitted `.rmg.json`, but the rule is absolute: **do not start this
without the owner explicitly approving the edit to
`template_variant/zone.go`.**

**Fix, if approved.** Change the field to `*models.Position`
(= `data.Vec2[float64]`) — note this makes `internal/entities/template` import
`internal/models`, so **first check for an import cycle**: if
`internal/models` already imports `internal/entities` (it does, via
`ManualZoneSave` in `internal/dtos/editor_state_dto`), the type must instead be
`*data.Vec2[float64]` from
[vec2.go](../../internal/helpers/data/vec2.go), which has no such dependency.
Then delete the pack/unpack at every producer and consumer.

**Blocked by:** §2.3 — do that first so the preview side is already float-native
and this becomes a mechanical type swap.

**Tests.** No behaviour change is expected. The proof obligation is a golden
test: generate a template before and after, assert the `.rmg.json` bytes are
identical.

---

### 2.5 ⚪ ⚠ Move `entities/types.go` under `template/` and rename `template` → `template_entity`

**Evidence.** [types.go](../../internal/entities/types.go) is a single
alias-only file re-exporting 40 types out of the seven `template*` subpackages:

```go
type (
	RmgTemplate = template.RmgTemplate
	…
	Zone = template_variant.Zone
	…
)
```

It sits **outside** the protected tree; everything it aliases sits inside it.

**Why the owner wants it.** `template` is an extremely generic package name for
something that is specifically the `.rmg.json` **entity** schema, and the alias
file living one directory above the thing it aliases makes the ownership
boundary of the protected tree ambiguous.

**⚠ OWNER DECISION REQUIRED — protected directory.** The rename touches the
`package` clause of **every file** under
[internal/entities/template/](../../internal/entities/template/) and moves a file
*into* that tree. AGENTS.md §2.1 forbids both without explicit approval. There
is no functional benefit, only naming; **do not start without a go-ahead.**

**Fix, if approved.**
1. `git mv internal/entities/types.go internal/entities/template/types.go` and
   change its package clause; the aliases then reference sibling subpackages.
2. Rename the directory `internal/entities/template` →
   `internal/entities/template_entity` and update the `package template` clause
   in [rmgTemplate.go](../../internal/entities/template/rmgTemplate.go).
3. Update every import path across the repo.
4. ⚠ **Do not** run a bulk in-place rewrite (AGENTS.md §2.6). Use
   `gopls rename` / the language server, then `gofmt -w` on the explicit list
   produced by `gofmt -l`.

**Open question for the owner:** if `types.go` moves *inside* the protected
tree, the alias list becomes protected too — every future type alias then needs
owner approval. Confirm that is intended, or keep `types.go` where it is and do
only the rename.

**Tests.** Pure move/rename; `go build ./...`,
`go vet -tags='integration_test,gui' ./...`,
`go run ./cmd/testlayoutcheck .` and the full suite are the verification.

---

### 2.6 🟠 113 production files name an Entity from outside the permitted layers

**Evidence.** Batch I §12 turned the Entity/Model/DTO doctrine into a gate
([test/unit/architecture/dependency/layering_test.go](../../test/unit/architecture/dependency/layering_test.go)).
Turning it on exposed the pre-existing breach it had to be seeded with:
**113 files in 23 packages** name a type from `internal/entities` while sitting
outside the permitted namers (`internal/repositories`, `internal/models`,
`internal/entities`, `internal/mappers`, `internal/helpers/*_helpers`).

| Area | Files |
| --- | --- |
| `internal/services/**` (16 packages) | 85 |
| `internal/dtos` | 11 |
| `app/gui/**` (4 packages: `dialogs`, `drivers`, `editor`, `models`) | 11 |
| `internal/handlers` + `handler_interfaces` | 6 |

A second, much smaller list rides along in the same test: **6 files in 3
packages** (`internal/services/bonuses`, `internal/services/pickers`,
`internal/services/zone_content`) consume DTOs below the handler boundary.

**Why it is not simply a defect.** Base `internal/entities` is the `.rmg.json`
schema vocabulary — `Zone`, `Connection`, `RmgTemplate` — and §0.5.3 of the
Batch I plan deliberately kept it as a repository-wide alias façade. The whole
generator exists to build those types, so a service naming `entities.Zone` is
not the same kind of mistake as a service naming `editor_state.EditorState`.
Scoped to the entity layer Batch I actually created
(`internal/entities/editor_state` + `/topology`) the breach is **one file**:
[fileService.go](../../internal/services/file_service/fileService.go), and only
as the generic argument of `repositories.IFileRepository[editor_state.EditorState]`.

**Fix (incremental, one package per batch — the allow-list only ever shrinks).**

1. `internal/services/{bonuses,pickers,zone_content}` — give each the treatment
   `internal/services/editor` got in Batch I Phase 10: a model-side
   request/result pair with the handler mapping onto it. Removes the DTO
   allow-list entirely.
2. `internal/dtos` and `internal/handlers` — decide whether the DTO layer
   naming `entities.Zone` / `entities.RmgTemplate` is a breach at all, or
   whether the schema vocabulary deserves a documented carve-out like
   `internal/helpers/data` already has.
3. `app/gui/**` — the 11 files are the zone editor (`entities.Zone`,
   `entities.Connection`) and the drivers. They need model wrappers before they
   can drop the import.
4. `internal/services/**` — the large tail. Only worth doing if step 2 rules
   that the schema vocabulary is genuinely off limits below the repositories.

**Do not** widen the allow-lists in `layering_test.go` to make a new package
compile; clean the package instead.

**Tests.** The gate is the test. Removing an allow-list entry must leave
`go test ./test/unit/architecture/...` green; it currently fails with the exact
file list when an entry is dropped, which is how it was verified.

---

## 3. Duplicated values that belong in `internal/common`

Owner decision 2026-08-11: **all four families move**. These are small,
independent, and safe to batch together.

Already centralised and correctly used — do **not** re-report: guard weekly
increments
([guardWeeklyIncrement.go](../../internal/common/common_connections/guardWeeklyIncrement.go)),
guard-strength presets
([guardStrength.go](../../internal/common/common_connections/guardStrength.go)),
content-distance presets
([distancePresets.go](../../internal/common/common_distances/distancePresets.go)),
zone names/prefixes and labels
([zoneNames.go](../../internal/common/constants/zoneNames.go),
[zoneLabels.go](../../internal/common/constants/zoneLabels.go)), map sizes
([mapSizes.go](../../internal/common/mapSizes.go)), permissions
([permissions.go](../../internal/common/constants/permissions.go)).

### 3.1 ✅ DONE 🟡 `"mandatory_content_hub"` is repeated in four production files

**Resolved 2026-08-12 (batch C)** together with §3.2 — see the note at the end
of §3.2.

**Evidence.**
[mandatoryContentProvider.go](../../internal/services/template_generator/providers/mandatoryContentProvider.go#L120)
and [#L146](../../internal/services/template_generator/providers/mandatoryContentProvider.go#L146),
[hubTopology.go](../../internal/services/template_generator/providers/topology/hubTopology.go#L77),
[geometricHubTopology.go](../../internal/services/template_generator/providers/topology/geometricHubTopology.go#L94),
[hubClusterService.go](../../internal/services/template_generator/providers/topology/tournament_variant/hubClusterService.go#L66).

**Why it is wrong.** The string is a cross-file contract: the topology writes it
into `Zone.MandatoryContent`, the content provider must emit a group with the
same name, and the parallel C# editor reads it
([mandatoryContentProvider.go](../../internal/services/template_generator/providers/mandatoryContentProvider.go#L132)
documents that). A typo in one of five places produces a template the game
silently generates without hub content.

**Fix.** New file `internal/common/constants/mandatoryContentNames.go`:

```go
package constants

const (
	MandatoryContentHubName      = "mandatory_content_hub"
	MandatoryContentNeutralPrefix = "mandatory_content_neutral_"
	MandatoryContentSidePrefix    = "mandatory_content_side_"
)
```

Replace all five literals. **Keep the exact strings** — they are the on-disk
contract with the game.

**Tests.** The existing tests already assert the literal
(e.g. [createContents_test.go](../../test/unit/internal/services/template_generator/providers/mandatoryContentProvider/createContents_test.go#L233)) —
leave the *test-side* literals as literals so they still catch an accidental
change to the constant's value. Add
`test/unit/internal/common/constants/mandatoryContentNames/` only if a helper
function (not a bare constant) ends up being introduced.

### 3.2 ✅ DONE 🟡 `"mandatory_content_neutral_"` / `"mandatory_content_side_"` prefixes repeated

**Evidence.** Neutral prefix:
[mandatoryContentProvider.go](../../internal/services/template_generator/providers/mandatoryContentProvider.go#L61),
[#L106](../../internal/services/template_generator/providers/mandatoryContentProvider.go#L106),
[topologyBase.go](../../internal/services/template_generator/providers/topology/base/topologyBase.go#L120).
Side prefix:
[mandatoryContentProvider.go](../../internal/services/template_generator/providers/mandatoryContentProvider.go#L48),
[#L92](../../internal/services/template_generator/providers/mandatoryContentProvider.go#L92),
[zoneFactory.go](../../internal/services/zones/zoneFactory.go#L68).

**Fix.** Same constants file as §3.1. Consider two tiny helpers
(`NeutralMandatoryContentName(label string) string`,
`SideMandatoryContentName(label string) string`) so the concatenation itself is
written once; if so, they get a mirrored unit-test folder per AGENTS.md §4.6.

**Resolved 2026-08-12 (batch C)**, §3.1 and §3.2 together, with the helper
option taken and a naming convention the owner set: a name **builder** is
`Get<X>For(label)` — `Get` separates it from the constant, `For` says it derives
a new name rather than returning one.

- [contentNames.go](../../internal/common/constants/contentNames.go) —
  `HubContentName`, `NeutralContentPrefix`, `SideContentPrefix`,
  `GetNeutralContentNameFor`, `GetSideContentNameFor`. No
  `"mandatory_content_*"` literal remains outside this file; test-side literals
  were deliberately left in place.
- [zoneNames.go](../../internal/common/constants/zoneNames.go) gained the matching
  `GetHubZoneNameFor` / `GetPlayerZoneNameFor` / `GetNeutralZoneNameFor`, and
  every production `XZonePrefix + label` now goes through them.
- The prefix **checks** were routed to the pre-existing
  `zone_helpers.IsZoneName*` predicates, which also let
  `internal/common/constants/hubLabel.go` (`IsHubLabel`, added in batch B) be
  deleted in favour of `zone_helpers.IsZoneNameHub`.
- Scope the owner added on review: connection names got the same treatment in
  [connectionNames.go](../../internal/common/constants/connectionNames.go) — 19
  unexported prefixes plus `Get*ConnectionNameFor` builders, wired through the
  chain/ring/web/geometric/tournament topologies.

Not converted, deliberately: the exact-`"Hub"`-vs-`"Hub-"` distinction in
[layoutRingHub.go](../../internal/services/preview_service/layoutRingHub.go#L32-L43),
where `IsZoneNameHub` would conflate the shared hub with a per-player hub and
change the preview layout.

### 3.3 ✅ DONE 🟡 Portal guard `25000` magic literal

**Evidence.**
[topologyConnectionService.go](../../internal/services/template_generator/providers/topology/base/topologyConnectionService.go#L64):
`WithGuardValue(tuning.ScaleByBorderGuardStrength(25000))`.

**Disposition: superseded by §1.3.** Landing §1.3 deletes the literal outright.
Only if §1.3 is deferred should this be done on its own, as a named constant in
[common_connections](../../internal/common/common_connections/) — note the value is
numerically identical to `neutral_zone.QualityHigh.GetGuardValue()`, so name it
for its meaning, not its number.

**Resolved 2026-08-12 (batch B)** by §1.3; no constant was needed.

### 3.4 ✅ DONE 🟡 Foothold placement distances are inline literals

**Evidence.**
[mandatoryContentProvider.go](../../internal/services/template_generator/providers/mandatoryContentProvider.go#L161-L189):

```go
					BuildCrossroadsRule(models.DistancePreset{Min: 0.2, Max: 0.3}, 0),
	…
						WithDistance(models.DistancePreset{Min: 0.2, Max: 0.4}).
	…
						WithDistance(models.DistancePreset{Min: 0.5, Max: 0.5}).
```

**Why it is wrong.** These are tuning constants for remote-foothold placement
expressed as anonymous `DistancePreset` literals in the middle of a builder
chain, while every other distance in the project comes from the named catalogue
in
[distancePresets.go](../../internal/common/common_distances/distancePresets.go).
Nothing tells a future maintainer what `0.5/0.5` means or that it is meant to
differ from `Far`.

**Fix.** Add a `GetFootholdDistancePresets()` (or three named presets) to
[common_distances](../../internal/common/common_distances/), following the shape of
`GetContentDistancePresets`, and reference them from the provider.
**Do not silently reuse the existing content presets** — `0.2–0.3` and
`0.2–0.4` do not match any of `Next To`/`Near`/`Medium`/`Far`/`Very Far`, so
folding them in would change generation.

**Tests.** `test/unit/internal/common/common_distances/distancePresets/` gains a
test per new accessor; the existing
`test/unit/internal/services/template_generator/providers/mandatoryContentProvider/`
tests must keep asserting the *numeric* values so a constant rename cannot
silently retune placement.

**Resolved 2026-08-12 (batch C).**
[footholdDistancePresets.go](../../internal/common/common_distances/footholdDistancePresets.go)
exposes `GetFootholdDistancePresets()` returning a named struct
(`Crossroads`, `NearMainCastle`, `NearSecondCastle`); the bounds are unchanged
and the presets are deliberately **not** in the user-facing
`GetContentDistancePresets` catalogue. The new `Name` fields cannot reach the
output because `WithDistance` copies only `Min`/`Max`.

---

## 4. UI / UX

### 4.1 ✅ DONE — "Save As" is really "Save To" — the UI offered a filename it then discarded

**Resolved 2026-08-11 (batch E).** The save dialog no longer lets the user type
a filename it cannot honour. Outcome, so this entry stands on its own:

- `getSaveRowWidget` shows a **read-only** textbox labelled `"Will save as:"`,
  populated with the resolved name rather than a hint.
- `State.SaveTo` (was `SaveAs`) sanitizes and trims the template name and
  appends `.gen.json` **only when the result is non-empty** — an unnamed
  template resolves to no filename at all, instead of the old `".gen.json"`
  that the repository would have silently rewritten to
  `Generated_Template.gen.json`.
- A `hasResolvedSaveName` guard disables the confirm button and renders
  `missingSaveNameMessage` ("Template name is required.") inline under the row,
  mirroring the new-folder row's error pattern.
- `onEntryClicked`'s `modeSaveFile` branch — which retargeted the save to a
  clicked file row — was removed: with a read-only preview it would have
  changed the target silently. Row clicks now only mean something in open mode.
- Toolbar button `"Save As"` → `"Save To"`, dialog title `"Save File"` →
  `"Save To"`. `NewSaveFileDialog`, `modeSaveFile` and `onSave` kept their
  names (owner decision): they name the *explorer mode and callback*, not the
  toolbar action.
- `SetFilename` was dropped from the testexports (it existed only to drive an
  editable field) in favour of the read accessors `ResolvedSaveName` and
  `SaveNameReadOnly`.
- The whitespace-only disabled-confirm test was deleted: that state is no
  longer reachable through any production path. See
  [test_observations.md](test_observations.md) for the recorded branch.
- Docs updated in [QUICKSTART.md](../../QUICKSTART.md) and
  [README.md](../../README.md); the 10 window snapshot goldens were regenerated
  locally for the new button label.

<details>
<summary>Original finding</summary>

**Evidence.** Writing editor state as `{TemplateName}.gen.json` is **intended**
(review §1.1, owner-approved):
[fileService.go](../../internal/services/file_service/fileService.go#L41-L48)
passes `filepath.Dir(filePath)` and `editorState.TemplateName` to the
repository, and
[saveSettings_test.go](../../test/unit/internal/services/file_service/fileService/saveSettings_test.go#L14-L60)
pins that behaviour deliberately.

The defect is the UI. The dialog still offers an editable **name** field whose
value is silently dropped —
[fileExplorerDialogToolbar.go](../../app/gui/dialogs/fileExplorerDialogToolbar.go#L35-L51):

```go
	hint := fmt.Sprintf("filename%s", saveFileSuffix)
	…
				layout.Rigid(widgets.NewLabelBigWidget(theme, "Save as:", themes.ColorsBase.TextDim)),
				widgets.NewDefaultComponentSpacer(),
				layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.filenameEd, hint, false)),
```

**Why it is wrong.** The user types `my-template`, presses Save, and gets
`Whatever The Template Name Is.gen.json` in that folder. The dialog picks a
*directory*, nothing more, and the whole vocabulary around it ("Save As", "Save
File", an editable filename box) promises otherwise.

</details>

---

## 5. Testing

### 5.1 ✅ FIXED 🟠 Zone editor pointer flows are untested (drag-to-connect, zone drag + snapping)

**Fixed in batch H, 2026-08-11.** The pointer flows landed as eight tests in
[zoneEditorPointer_integration_test.go](../../test/integration/gui/zoneEditorPointer_integration_test.go),
driven through the real window over the Geometric Hub layout: a zone dragged to
a new position and committed through Apply's normalized manual position, a drag
that snaps onto another zone's centre line and the grid, a drag-to-connect in
Add connection mode, a drag that ends on empty canvas and creates nothing, a
right-click that deletes a curve, a zone placed from Add zone mode, the placed
zone landing where it was clicked, and a drag inside the 6 px dead zone that
moves nothing.

Each test asserts a state change as well as its golden, since a golden alone
cannot tell "the drag did the right thing" from "nothing happened". The canvas
coordinates come from the post-§2.3 float geometry and are pinned exactly - the
snap test pins `(290, 251.88571428571436)` - because a pin that rounds would
hide the very regression §2.3 was about. Snap **guides** are deliberately not
asserted: every zone in this layout shares x = 290, so several guides propose
the same correction and which one is reported depends on map iteration order.

⚠ **The exact float pin is an unverified portability risk.**
`251.88571428571436` is the full-precision result of the amd64 computation, and
it has only ever been run there. If the geometry ever evaluates differently on
arm64 the assertion breaks, and the failure will look like a snapping bug rather
than a platform difference. Rounding the pin is **not** the fix — that is exactly
what §2.3 was about. If it does break, pin with a tolerance far tighter than one
pixel (an `InDelta` of ~1e-9) so a real regression still fails.

[test_observations.md](test_observations.md) was updated to point at the new
tests instead of calling them future work.

### 5.2 ✅ FIXED 🟡 Zone editor property panels (`widget.Editor` / dropdown paths) are untested

**Fixed in batch H, 2026-08-11.** The panels landed as eighteen tests in
[zoneEditorProperties_integration_test.go](../../test/integration/gui/zoneEditorProperties_integration_test.go):
the zone `Size`, `Guard x` and `Weekly +` editors including `Size`'s clamp to
0.1–2.0 and its rounding to two decimals; the neutral-zone `Quality` and
`Castles` dropdowns, which exercise the `ApplyZoneEditorQuality` reprofile path;
the connection guard value typed and a non-numeric value rejected; the `Type`,
`Guard zone`, `Guard preset` and `Weekly +` connection dropdowns; and the
Advanced options checkbox together with the `Match group`, `Guard escape` and
`Sim turn squad` rows it reveals.

Zone assertions read the editor's own zones, connection assertions go through
`ClickApply()` and `CurrentState().ManualConnections` - committed state either
way, never widget text.

Two items from the original sketch changed:

- `TestWhenAZoneNameIsTyped_…` **was not written**: the zone name is a read-only
  `material.Body1` label and the dialog offers no rename, so there is no path to
  drive. Recorded in the test file and in [test_observations.md](test_observations.md).
- Typing **inserts at the caret, which sits at the start of a focused field**,
  so the expectations are written for insertion rather than replacement (typing
  `1` into `35000` yields `135000`).

One harness gap is filed rather than fixed: `integration_common`'s
`zoneEditorZone*Y` rows were measured on a zone whose note wraps to one line,
which a player spawn and a neutral zone do but the shared `Hub` does not, so the
Hub's property rows cannot be clicked through the handler. The rows are the same
code for every zone, so no behaviour is left uncovered.

⚠ **Open: is a golden per handler action the right granularity?** Batches F–H
adopted "one snapshot per action", which left roughly 145 new goldens beside the
~173 already committed — on the order of 21 MB of PNGs in the repository. The
question was never settled, and every future driving handler compounds it. The
trade-off: a golden per action is what localises a visual regression to the
action that caused it, but most of these snapshots differ from their neighbour in
a few pixels of one widget. If the size becomes a problem, the lever is to keep a
golden only for actions that change layout and assert committed **state** for the
rest — not to loosen the comparison tolerance, which §5.5 already closed off.

### 5.3 ✅ FIXED 🟡 File explorer: hidden-file toggle and pointer-driven row/scroll interactions

**Fixed in batch F, 2026-08-14.** The listing
behaviours landed as five new tests in
[fileExplorerDialogListing_integration_test.go](../../test/integration/gui/fileExplorerDialogListing_integration_test.go)
(toggle on, toggle off, row selection in open mode, directory descent, wheel
scroll), and the twelve existing tests in
[fileExplorerDialog_integration_test.go](../../test/integration/gui/fileExplorerDialog_integration_test.go)
were migrated onto the same handler, so the whole dialog is now driven by real
pointer events through the real toolbar with a golden per action. Every `Click*`
test-export on the dialog was deleted as a result. The listing rows had no
accessibility label to address, so `getEntryRowWidget` now emits button
semantics; that is the one production line the tests required.

**Original report follows.**

**Evidence.** [test_observations.md](test_observations.md) for the file
explorer: *"Still uncovered: the hidden-file toggle and the pointer-driven
row/scroll interactions (owner decision - excluded from the scenario set)."*
The owner re-opened this on 2026-08-11.

The toggle exists at
[fileExplorerDialogToolbar.go](../../app/gui/dialogs/fileExplorerDialogToolbar.go#L31)
(`widgets.NewToggleButtonWidget(theme, "Show hidden", &this.hiddenToggle, this.showHidden)`),
and the underlying filtering policy is already unit-tested in
`internal/services/file_system`.

**Why it matters.** The policy is tested; the *wiring* is not. A toggle that
flips a field but never re-lists, or a row click that resolves the wrong entry,
would pass every existing test.

**Fix.** Extend
[fileExplorerDialog_integration_test.go](../../test/integration/gui/fileExplorerDialog_integration_test.go):

- `TestWhenShowHiddenIsToggledOn_HiddenEntriesAppearInTheListing`
- `TestWhenShowHiddenIsToggledOff_HiddenEntriesDisappearAgain`
- `TestWhenARowIsClicked_ThatEntryBecomesTheSelection`
- `TestWhenADirectoryRowIsClicked_TheListingDescendsIntoIt`

Create the hidden fixture entries cross-platform: on Linux a leading `.`
suffices, on Windows the attribute must be set — put that behind a
build-tagged helper in `test/test_helpers/` rather than skipping the test on one
OS (AGENTS.md §2.2).

**If investigation shows** the scroll interaction cannot be driven without a
scrollbar hit-target that `AppRunner` cannot reach, cover the row click only and
leave the scroll note in [test_observations.md](test_observations.md) with the
reason.

**Coordinate with §4.1** — it rewrites the save-filename scenarios in the same
file. Do §4.1 first or expect a merge conflict.

### 5.4 ✅ DONE 🟡 The GUI handler framework is a design comment, not a design

> **Landed in batches L ((a)–(c)) and M ((d)–(f)).** Everything below is the
> original diagnosis, annotated with what was actually built. The plan files that
> tracked the work have been deleted, so **this section is now the only record** —
> keep it self-contained. (g) was kept as a standing guideline rather than a task.
>
> **Re-measuring the coordinates** — the one procedure that is unrecoverable
> otherwise. Do this whenever a UI change moves a widget; never nudge a constant
> by trial and error:
>
> 1. Add a temporary accessor on `AppRunner` returning the last frame's `*op.Ops`,
>    and a temporary `//go:build integration_test` test that drives the UI to the
>    state you want to measure.
> 2. In that test build a **fresh** `input.Router`, then
>    `replay.Frame(ops); nodes := replay.AppendSemantics(nil)`. The flat slice
>    holds every node once, with `Desc.Bounds` in absolute window coordinates at
>    the harness geometry (1600×900, `PxPerDp = 1`).
> 3. **Only `Desc.Class == semantic.Button` nodes have trustworthy bounds.**
>    `material.CheckBox` and `dropdownSelector.getTriggerWidget` never call
>    `utils.AddButtonSemantics`, so they emit no button node at all. For those,
>    read the **neighbouring label** node's `bounds.Min`, which is correct — its
>    `Max` and `center` are garbage (values like `1000072`) because the replay
>    router never pops the clip stack.
> 4. **Confirm every non-button coordinate by driving it** and asserting the state
>    change it causes. A position read from a label is a guess until it clicks the
>    thing you meant.
> 5. Delete the probe accessor and the probe test.
>
> Constants are points **inside** a widget, not exact centres: a tab button's
> bounds shift 1–3 px depending on which tab is selected (`General` starts at
> x = 612 when selected, x = 610 otherwise).
>
> **Regenerating goldens.** Locally, on a real GPU, **never in CI** (§5.5), with
> `go test -tags='integration_test,gui' -run '<TestName>' ./test/integration/gui/...
> -count=1 -update`. Always name the tests — a bare `-update` rewrites goldens you
> did not intend to touch. `.golden` files are PNGs; copy to `.png` to view.

**Evidence.** [runnerHandler.go](../../test/test_helpers/integration_common/runnerHandler.go#L61-L106)
carries a 46-line first-person design note above a struct that currently
implements three tab clicks and one mask. The whole framework today is:

```go
func (this *baseHandler) ClickGeneralTab() *baseHandler {
	this.runner.ClickAt(f32.Pt(672, 60))
	this.runner.VerifySnapshot()
	return this
}
```

Consumers: [window_snapshot_integration_test.go](../../test/integration/gui/window_snapshot_integration_test.go#L20-L25)
and [window_tab_cycling_test.go](../../test/performance/window_tab_cycling_test.go#L25).

**Why it matters.** Gio has no Playwright/Selenium, so this handler *is* the
GUI test API. §5.1, §5.2 and §5.3 will each be written against it; if they land
first, three more test files acquire their own ad-hoc `f32.Pt` literals and the
framework never happens. The comment is also the only place this design lives —
a doc comment is not a backlog and will rot.

**Fix — seven separable pieces, in this order.**

**(a) ✅ Hygiene (mechanical, do first).** The file is `runnerHandler.go` but the
struct is `baseHandler`; AGENTS.md §4.1 requires `baseHandler.go`. `NewHandler`
also returns the unexported `*baseHandler` from an exported function — either
export the type or (once (d) exists) return the handler contract. Interfaces
follow AGENTS.md §4.2.1/§4.2.2: with fewer than five implementations they stay
in this package as `*Interface.go` files.

> *Done.* Renamed in L; the type was **exported** to `*BaseHandler` in M by owner
> decision — no `IBaseHandler`, because there is one implementation and an
> interface would only obstruct the embedding in (d).

**(b) ✅ Settle the coordinate strategy.** Recommendation: **keep the coordinates
literal but stop scattering them.** Gio exposes no widget-rect lookup without a
new `*_testexports.go` seam, and computing positions from the layout code would
re-implement the thing under test. Put them in one `handlerCoordinates.go` as
named constants derived where possible from
[app/gui/constants/ui.go](../../app/gui/constants/ui.go) (`DefaultPadding`,
`DefaultLabelWidth`, `DefaultPreviewWidthMaximum`), so a padding change is one
edit rather than a hunt. Only compute at runtime where a value genuinely varies
(slider value → x, list row index → y).

> *Done.* The literal-but-centralised approach held, with one revision: rather
> than deriving from `ui.go`, every coordinate was **measured** and then confirmed
> by driving the widget — see the procedure at the top of this section. The
> dropdown row pitch is the one runtime-computed value.

**(c) ✅ Narrow the mask.** `setRandomTopology` blanks
`image.Rect(WindowWidth-470, 0, WindowWidth, WindowHeight)` — a 470 px column
over the **full window height**, i.e. roughly a third of every snapshot,
including stable chrome. Replace the single catch-all with named helpers that
mask only what is genuinely nondeterministic: the preview canvas (random
topology), the output-directory field (differs per machine/OS/Steam layout —
see AGENTS.md §2.7), and the timestamp/path inside the status message. Derive
the preview rect from `DefaultPreviewWidthMaximum` (440) rather than the
unexplained 470.

> *Done.* **Masked area 423 000 → 208 000 px**, putting the canvas border, the
> legend, the template name, Browse/Reveal and the Generate and Save Template
> buttons back under test. Exactly three regions are nondeterministic:
> preview canvas `(1163,203)-(1577,627)`, status block `(1157,726)-(1583,775)`,
> output directory `(1157,809)-(1583,838)`.
>
> **How to re-measure a mask** (reusable): capture two *unmasked* runs of the
> same action and diff them — the differing pixels **are** the nondeterministic
> regions. Everything outside those three rectangles differed only by
> anti-aliasing noise of delta 1–8.

**(d) ✅ Per-tab and per-dialog handlers.** `baseHandler` gains
`ClickGeneralTab() *generalTabHandler` etc.; tab handlers embed `baseHandler`
(tab bar and toolbar stay live), dialog handlers **do not** (a dialog disables
the background, so inheriting its clicks would be a lie). One struct per file
(AGENTS.md §4.1).

> *Done.* Three tab handlers embed `*BaseHandler` (a pointer, so the shift state
> of (e) is shared rather than copied). Two dialog handlers — `FileExplorer` and
> `ZoneEditor` — hold `base *BaseHandler` without embedding and expose only
> `IsOpen()` / `Close()`. They were **reachability-only**: they owned no canvas
> coordinates and took no snapshot, because driving a dialog through the window
> is nondeterministic by construction (the explorer lists the machine-detected
> templates directory, AGENTS.md §2.7; the zone editor draws a freshly randomised
> layout). The dialog *behaviour* tests of §5.1–§5.3 constructed the dialog
> directly through the package-local `newDialogContext`.
>
> **Superseded for both handlers.** Batch F gave `FileExplorerHandler` a seeded
> fixture directory and a golden per action, and batch H did the same for
> `ZoneEditorHandler`: `WithFixtureDirectory()` masks the per-machine path and
> `LayoutAndZonesTabHandler.SelectTopology` takes the layout off Random, which is
> what made snapshotting deterministic. Both are now **driving** handlers with
> canvas and side-panel coordinates and a snapshot per action, and §5.1–§5.3 are
> written against them rather than against `newDialogContext`. Only the five
> zone-editor cases the window cannot reach are still dialog-direct; each one
> names its reason in the test file itself.
>
> Toolbar methods are `ClickNew` / `ClickLoad` / `ClickSaveAs`; **`Exit` is never
> reachable** (`State.Exit()` calls `os.Exit(0)` and would kill the test process)
> and `Save` is excluded because it writes into the real output directory.

**(e) ✅ Layout-shifting state.** A handler must remember what it toggled, because
some widgets move the ones below them: *Allow non-official larger map sizes*
expands the Map Size dropdown, and dropdowns render inline on the canvas rather
than floating. Store those flags on the handler and offset subsequent
coordinates from them; a handler that clicks blind will click the wrong widget.

> *Done, with two corrections to the framing above.* Only **one** flag is stored,
> `isExperimentalMapSizes`, because it decides whether the dropdown has 11 rows
> or 28 and a row index is meaningless without that count. The `isSingleHero`
> shift, which this section elsewhere calls the larger one, needs **no**
> bookkeeping: it adds or removes three slider rows in the *right* column only,
> and no handler coordinate sits below them.
>
> A third shifter was added that this section did not anticipate: the panel's own
> **scroll offset**, tracked by reading the real `widget.List` position through a
> `*_testexports.go` accessor rather than by accumulating injected deltas —
> `layout.List` clamps at both ends, so an accumulator diverges from the truth at
> the first clamp.
>
> Separately: a click is processed on the frame it is queued against, but a panel
> only writes its widget values back to the editor state in `SaveToState` on the
> **following** layout. Every state-mutating handler method therefore ends with a
> `commit()` frame, before any snapshot.

**(f) ✅ Scrolling.** `AppRunner` has `ClickAt`, `MoveTo`, `DragTo` and
`InputText` — but no scroll. Anything below the fold is unreachable, which is a
hard blocker for §5.2. Add `Scroll(point f32.Point, delta f32.Point)` injecting
a `pointer.Scroll` event through the same router, mirroring `ClickAt`.

> *Done*, plus `BaseHandler.ScrollPanel(delta)`. Note that a golden alone cannot
> prove a scroll happened — a `Scroll` that silently did nothing would still match
> a golden captured from an equally unscrolled frame — so the snapshot is backed by
> headless assertions on the real list position. Note also that no panel has a
> useful scroll range out of the box: Bonuses & Bans never overflows and Layout &
> Zones overflows by ~18 px, which is why the handler set gained one method beyond
> the agreed scope, `ToggleAdvancedZoneControl()`, taking the range to ~386 px.

**(g) ⬜ Keep `*_testexports.go` out of the handler.** `SelectedTabIndex`,
`TabCount`, `DialogsOpen`, `CurrentState` and `Status` are test-only exports.
Handlers should drive pixels and assert through snapshots; reach for an
accessor only where a pixel genuinely cannot express the assertion, and say why
in a one-line comment.

> *Kept as a standing guideline, not built as a task.* It was applied while
> writing (d)–(f); each of the three new accessors carries the required one-line
> justification.

**Scope guard.** Do **not** build (d)–(g) speculatively. (a)–(c) are worth
doing on their own; the rest grows one method at a time as §5.1–§5.3 need it.
The comment's own worry — *"this looks like an entire framework"* — is the
correct instinct: an unused handler method is untested test code.

> *Superseded by owner decision:* (d)–(f) were built standalone and ahead of
> §5.3, so that §5.1–§5.3 are written against a finished API instead of growing
> one. The guard's intent was still honoured — every method added has at least one
> test that exercises it.

**When it lands,** replace the design comment with a two-line pointer to this
section so the design lives in one place.

> *Done in batch L.*

### 5.5 ✅ DONE 🟡 GUI snapshots differ between local and CI, and the tolerance hiding it also hides regressions

> **Landed in batch L.** Steps 1 and 3 are done; step 2 was **rejected by the
> owner**. The plan file that tracked the work has been deleted, so this section
> is the only record — keep it self-contained.
>
> **Root cause was a real production bug, not a rasterizer difference.**
> `themes.NewTheme` built its shaper without `text.NoSystemFonts()`, so any glyph
> missing from `gofont.Collection()` was resolved through whatever face the OS
> happened to offer. `◆` (U+25C6), used in **every** section header, is one such
> glyph; CI's substitute face is 1 px shorter and 4 px narrower, and the shortfall
> **accumulates** down the panel. Evidence, from a per-band best-offset search:
> the checkbox rows at y 430–580 scored **1.2957 %** at zero offset and
> **0.0029 % at dy = −3**. Fixed by adding `text.NoSystemFonts()` and substituting
> seven glyphs.
>
> **Glyph inventory** (verify before using a new symbol in `app/`):
> missing from `gofont.Collection()` — `◆ ✕ ↺ ⚠ ✘ ✗ ↻ ⟲ ↩ ⚡ △ ◇ ✦ ⚔`;
> present — `♦ ◊ × ← ‼ ▼ ▲ → · “ ” ⌂ ± … –`. Substitutions applied:
> `◆`→`♦`, `✕`→`×`, `↺`→`←`, `⚠`→`‼`. This is enforced, not remembered:
> `newTheme_test.go` AST-scans every string and char literal under `app/` and
> fails on any rune the bundled fonts cannot render.
>
> **CI is verified green** against goldens generated locally on a real GPU, with
> no pinning. llvmpipe's residual difference — exactly 0.75× the AA coverage of a
> real GPU inside rounded-rect clips, backgrounds bit-identical — peaks at a
> per-pixel delta of ~40–45, below the 64 tolerance, so it barely registers.
> **Treat a future GUI-suite failure as a real rendering change, not as noise.**

**Evidence.** The last paragraph of the same comment: *"for some reason there
is a difference of some of the rendered text between local … and CI (looks like
some of the text is grayed out like not finishing a rerender)"*. The tolerance
that absorbs it is [comparer.go](../../test/test_helpers/integration_common/snapshot/comparer.go#L8-L11):

```go
// DefaultSnapshotThreshold is the maximum allowed normalized mean color
// distance between a golden snapshot and an actual screenshot (2%).
// Pipeline has discrepancies, I don't want to investigate them right now.
const DefaultSnapshotThreshold = 0.02
```

The two environments are genuinely different: goldens are generated locally on
Windows against a real GPU, while
[pr-validation.yml](../../.github/workflows/pr-validation.yml#L244-L257) runs the
suite under `xvfb-run` with `LIBGL_ALWAYS_SOFTWARE=1` (Mesa llvmpipe).

**Why it matters.** `Compare` returns a **mean** distance over the whole
1600×900 frame, and `Matches` passes anything under 2 %. A 40×40 button that
turns entirely black moves that mean by ~0.1 % — it passes. The threshold was
raised to swallow a font-rasterization difference and now swallows real
regressions with it, which quietly weakens every test §5.1–§5.4 will add.

**Fix — diagnose before tuning.**

1. **Rule out a half-rendered frame first.** "Grayed out like not finishing a
   rerender" describes a capture-timing bug, not a rasterizer difference.
   Capture two consecutive frames for the same action and compare them: if
   frame 2 differs from frame 1, `captureScreenshot` is racing the frame and
   the fix belongs in
   [appRunnerSnapshots.go](../../test/test_helpers/integration_common/appRunnerSnapshots.go#L120),
   not in the threshold. Download the `gui-snapshot-failures` artifact the
   workflow already uploads to see the actual CI pixels.

   > *Done, and the hypothesis was **disproved**.* Consecutive frames for the
   > same action are identical; `captureScreenshot` is not racing the frame. The
   > cause was the font fallback described above, so no capture-timing fix was
   > made. Do not resurrect this theory.

2. ~~**If it is genuinely llvmpipe text anti-aliasing**, make CI the reference:
   regenerate goldens in the CI environment (run the update job, download the
   artifact, commit the images) so the *comparison* is exact and only local
   runs are lenient.~~

   > **Rejected by the owner.** CI is a software renderer and must not become the
   > reference. Goldens are generated **locally on a real GPU only**, never in CI,
   > and an agent never runs CI — ask the owner. Step 3 is what makes that safe.

3. **Then replace the metric.** A mean is the wrong measure for UI diffs. Fail
   on the fraction of pixels whose per-pixel distance exceeds a small
   per-pixel tolerance (e.g. "> 0.5 % of pixels differ by > 10 %") so wide
   faint AA noise passes while a small solid change fails. Keep
   `Comparer.Threshold` configurable and update
   [test/unit/test/test_helpers/integration_common/snapshot/comparer/](../../test/unit/test/test_helpers/integration_common/snapshot/comparer/)
   alongside.

   > *Done.* The comparer is now two-gate: `Compare` returns
   > `Difference{MeanDistance, ChangedPixelFraction}` judged against
   > `DefaultMeanThreshold` 0.0025 (down from 0.02), `DefaultPixelTolerance` 64
   > and `DefaultChangedPixelThreshold` 0.0005. A pixel counts as changed when the
   > **largest** of its three channel deltas exceeds the tolerance. The broken CI
   > frames scored mean 0.66 / 1.22 / 0.66 % — all three passed the old 2 % gate —
   > against fraction **1.38 / 2.83 / 1.38 %**, tripping the new one by 27–57×.
   > Both measurements and both limits are named in the failure message.
   >
   > **If CI ever needs loosening, raise the fraction floor, never the mean**, and
   > only to a measured value.

**If investigation shows** the difference is unavoidable and CI-generated
goldens are impractical, say so in the constant's comment with the measured
worst-case difference — a documented 2 % is fine, an unexplained one is not.

---

## 6. Deferred decisions

### 6.1 ⚪ `createTopologyAdjacency` dead Chain/Ring branches

**Evidence.**
[zoneLabelProvider.go](../../internal/services/zones/zoneLabelProvider.go#L212)
declares `createTopologyAdjacency`; its only call site is
[#L82](../../internal/services/zones/zoneLabelProvider.go#L82). The
`case TopologyChain` and `case TopologyRing, TopologyCircles` branches (and the
`isIsolated` guard they use) are unreachable: the only production caller,
`GetHoldCityLabel`, gates on `IsHubCityToHold()`
(= `Topology == HubAndSpoke && IsCityHoldMode()`), so the switch always takes
`default`. Verified 2026-07 (review §5.5): single private call site, single
production caller (`templateGenerator.Generate`), all three symbols born
together in `bb50aab`, 0 % coverage on the dead branches.

**History.** A removal was implemented and verified green (build/tests/lint,
coverage 64.1 → 64.2) and then **rolled back by the owner** to keep the
topology-aware adjacency as a starting point.

**Still undecided — pick one:**

- **(a) Delete the branches.** Pure coverage-ratio win, zero behaviour change,
  removes the misleading impression that hold-city works on Chain/Ring.
- **(b) Start using them.** Extend hold-city (or another adjacency-based
  feature) to Chain/Ring/Circles topologies. This would additionally fix the
  `default` branch, which currently models Hub & Spoke as a *sequential ring*
  rather than its real star graph — i.e. the one branch that *does* run is
  arguably the wrong model.

**Note for whoever picks (b):** the `default` branch's ring model is the more
interesting defect. If hold-city ever needs true adjacency (e.g. "the city to
hold must border every player"), the current model gives wrong answers today,
silently.

**Owner-decision flag ⚠:** this reverses a rollback the owner performed
deliberately. **Do not act without an explicit instruction.**

---

## 7. Owner-decision summary

| Item | Why it needs the owner | Blocking |
| --- | --- | --- |
| §2.2 Branch A (tier on the entity) | Edits protected `internal/entities/template/` | Branch B is unblocked |
| §2.4 `[2]float64` → `Vec2` | Edits protected `internal/entities/template/` | Yes |
| §2.5 `template` → `template_entity` | Renames/moves inside protected tree | Yes |
| §6.1 dead adjacency branches | Reverses an owner rollback | Yes |
| §2.3 snapshot regeneration | Regenerated GPU snapshots need a human eyeball | Review, not approval |
| §1.2 / §1.3 guard values | Changes generated `.rmg.json` guard numbers for every hub topology and every portal | Approved 2026-08-11, landed 2026-08-12 |

---

## 8. Suggested execution order

Bugs first, protected-directory work last, nothing batched with something it
blocks. Each batch is one PR-sized unit; the owner reviews and commits.

| Batch | Items | Notes |
| --- | --- | --- |
| ✅ **A** | §1.4 | **Done 2026-08-12.** Two-line stderr fix + a `test_observations.md` entry. |
| ✅ **B** | §1.2 → §1.3 | **Done 2026-08-12.** No golden-template churn was required — the generator tests do not pin hub or portal guard values. |
| ✅ **C** | §3.1, §3.2, §3.4 | **Done 2026-08-12.** Extended on review with `constants/connectionNames.go` (connection-name builders). No behaviour change. |
| ✅ **D** | §1.1 | **Done 2026-08-14.** Deep `Clone` + regression tests. Cost +4.6 % frame time / +42 % allocs on `TabCycling`; spun the residual off as §1.5. |
| ✅ **E** | §4.1 | **Done 2026-08-11.** Save To rename + read-only resolved-name preview + blank-name guard. Regenerated the 10 window goldens for the new button label. |
| ✅ **F** | §5.3 | **Done 2026-08-14.** File-explorer pointer/hidden-file tests plus a full migration of the existing dialog tests onto `FileExplorerHandler`; the save-mode row-click test applies to **open mode only**, per §4.1. Coverage flat at 72.9 %. |
| ✅ **G** | §2.3 | **Done 2026-08-19.** Float preview geometry end to end, rounded once at the draw boundary. **No goldens moved** — the preview canvas is masked and the zone-editor handler takes no snapshots, so no `-update` was needed. Coverage flat at 72.9 %. Record: §2.3. |
| ✅ **H** | §5.1, §5.2 | **Done 2026-08-11.** Zone-editor pointer + property-panel tests against the post-§2.3 float coordinates: eight pointer tests and eighteen property tests, all driven through the real window with a golden per action. Turned `ZoneEditorHandler` from a reachability-only handler into a driving one (canvas, side-panel and Apply actions). `TestWhenAZoneNameIsTyped_…` dropped — the zone name is a read-only label. Coverage flat. Record: §5.1, §5.2. |
| ✅ **I** | §2.1, §1.5 | **Done 2026-09-01.** `EditorStateDto` rework across twelve phases (5 and 11 superseded mid-flight), folding in §1.5 as phase 6. Entity/Model/DTO split with the **Model owning the structure**; `.gen.json` shape unchanged throughout. Phase 6 cut render-path allocations by 62 %; phase 12 added the layering gate and spun off §2.6. Doctrine now lives in **AGENTS.md §4.4.1**. Records: §2.1, §1.5, §2.6. |
| **J** | §2.2 Branch B | Zone tier single source of truth without a protected edit. Benefits from §2.1's model layer. |
| **⚠ K** | §2.2 Branch A, §2.4, §2.5, §6.1 | Owner-gated. Do not schedule until each is explicitly approved. §2.4 depends on §2.3. |
| ✅ **L** | §5.4 (a–c), §5.5 | **Done 2026-08-14.** GUI test-harness groundwork: handler hygiene, named mask helpers (423 k → 208 k masked px), coordinate constants, two-gate snapshot comparer, and a real font-fallback bug in `themes.NewTheme`. §5.5 step 2 rejected — CI never becomes the golden reference. Full record in §5.4/§5.5 above. |
| ✅ **M** | §5.4 (d–g) | **Done 2026-08-14.** Built **standalone and ahead of F** by owner decision, not grown from it. Three tab handlers, two reachability-only dialog handlers, three toolbar methods, the `Scroll` seam, and layout-shift tracking. (g) kept as a standing guideline. Full record in §5.4 above. |
| ✅ **N** | §1.5 | **Folded into batch I phase 6, 2026-08-31.** Never ran standalone — the measurement showed the cost was the clone *mechanism* (lazy `linq` chains allocating for empty slices), not the panel read sites this item named. Record: §1.5. |
| **O** | §2.6 | Drain the two layering allow-lists seeded by batch I phase 12. Four independent steps, one package at a time; step 1 (the three DTO-consuming services) is the smallest and clears a list entirely. |

**Note on L/M.** Both are done; they sit last in the table only because it is
otherwise ordered by dependency.

**Coverage note.** Run the coverage task before and after **every** batch
(AGENTS.md §2.3) — the floor is **72.5 %** and the current figure is **73.9 %**
(72.5 % through batch B; batch C added the helper tests, batch D the clone and
accessor tests, batch I the entity/model/converter tests).

---

## 9. Baselines to hold

| Gate | Command | Expected |
| --- | --- | --- |
| Build | `go build ./...` | clean |
| Vet (default) | `go vet ./...` | clean |
| Vet (gated) | `go vet -tags='integration_test,gui' ./...` | clean |
| Test layout | `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| Layering | `go test ./test/unit/architecture/... -count=1` | pass (also runs inside the unit gate; see AGENTS.md §4.4.1) |
| Unit | `go test ./test/unit/... -count=1` | pass |
| Integration | `go test -tags=integration_test ./test/integration/... -count=1` | pass |
| GUI integration | `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | pass (needs GPU) |
| Coverage | `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...` then `go tool cover '-func=coverage.txt'` | **≥ 72.5 %**, currently **73.9 %** |
| Lint | `golangci-lint-v2 run ./... --issues-exit-code=0` | **0 issues** |
| Format | `gofmt -l ./app ./internal ./test ./cmd` | empty |
| Wire | `wire diff ./internal/composition/...` | no diff |
