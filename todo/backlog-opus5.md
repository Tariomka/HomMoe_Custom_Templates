# Backlog — 2026-08-11 (Claude Opus 5)

Compiled backlog in the format of [review-prompt.md](review-prompt.md)
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
| [runnerHandler.go](../test/test_helpers/integration_common/runnerHandler.go#L61-L106) | the owner's 46-line design comment, split into §5.4 and §5.5 |

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

**Item count: 0 🔴 · 6 🟠 · 10 🟡 · 3 ⚪ (19 total).**

**Baselines to hold (AGENTS.md §2.3):** unit coverage **72.5 %**, floor
**69.3 %** · `golangci-lint-v2 run ./...` **0 issues** · `gofmt -l` empty ·
`go run ./cmd/testlayoutcheck .` passes · build + vet clean under both
`integration_test` and `integration_test,gui`.

---

## 0. Disposition of the source documents

### 0.1 Carried forward from [review-opus5-08-04.md](review-opus5-08-04.md) ❗

| Prior item | Status re-verified 2026-08-11 | New section |
| --- | --- | --- |
| §1.4 Editor-state copies are shallow | Still open, verbatim — no `Clone` exists on `EditorStateDto` | §1.1 |
| §1.9 Fatal window error logged to a discard handler | Still open, verbatim | §1.4 |

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
| "Save As" is really "Save To" | Carried, UI treatment confirmed (read-only textbox) | §4.1 |

### 0.3 Disposition of [test_observations.md](test_observations.md)

**Promoted to backlog items** (they were recorded as "future work" only because
no synthetic-pointer seam existed; [appRunner.go](../test/test_helpers/integration_common/appRunner.go#L152-L217)
now provides `ClickAt`, `MoveTo`, `DragTo` and `InputText`):

| Observation | New section |
| --- | --- |
| Zone editor: drag-to-connect, zone drag + snapping ("need synthetic pointer events — still future work") | §5.1 |
| Zone editor: property panels' `widget.Editor` / dropdown paths | §5.2 |
| File explorer: hidden-file toggle + pointer-driven row/scroll interactions (was "owner decision — excluded"; owner re-opened it 2026-08-11) | §5.3 |

**Left in place as accepted, intentional gaps — do not re-report:** all the
Gio-widget/panel entries (`buttonWidget`, `sliderRowWidget`, `layoutPanel*`,
`previewPanel`), the `drivers.State` dialog-callback branches, the
`topologyConnectionService` private-policy note, and every entry under
*Unreachable defensive branches* (`connectInteriorStables` `len == 0`,
`providePreviewGenerator` `err != nil`, `atomicFileWriter` `Close`/`Sync`,
`helpers/io.go` Steam discovery, `buildShiftDerangement`).

### 0.4 Disposition of the [runnerHandler.go](../test/test_helpers/integration_common/runnerHandler.go#L61-L106) design comment

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

### 1.1 🟠 Editor-state copies are shallow, so snapshots alias live slices

*(= review-opus5-08-04 §1.4, restated with re-verified line numbers.)*

**Evidence.** `EditorStateDto` holds nine slice fields —
[editorStateDto.go](../internal/dtos/editorStateDto.go#L82-L94):

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
snapshot — [editorState.go](../app/gui/models/editorState.go#L31-L33):

```go
func (this *EditorState) GetCurrentState() dtos.EditorStateDto {
	return *this.current
}
```

[editorState.go](../app/gui/models/editorState.go#L40-L44):

```go
func (this *EditorState) SnapshotCurrentState() {
	previousState := *this.current
	this.previous = &previousState
	this.next = nil
}
```

[stateHandler.go](../internal/handlers/stateHandler.go#L55-L70) —
`ValidateEditorState(stateDto dtos.EditorStateDto, …)` takes the DTO **by
value**, mutates it through `issue.Fix(&stateDto)` /
`normalizeInactiveNeutralCounts(&stateDto)`, and returns it inside
`dtos.EditorStateValidationDto{State: stateDto, …}`.
[editorState.go](../app/gui/models/editorState.go#L46-L54) does the same
shallow trick again in `GetPreviousState`.

**Why it is wrong.** A struct copy duplicates slice *headers*, not backing
arrays, so `this.previous` shares element storage with `this.current`. Change
detection compares element-wise
([editorStateDto.go](../internal/dtos/editorStateDto.go#L187-L199) →
`contentRowSlicesEqual` / `slices.Equal`), so **any in-place element write makes
the change invisible to the unsaved/regenerate machinery**: the editor would not
mark the file dirty and `AutoRegenerate` would not fire. The same aliasing leaks
live editor state out of `GetCurrentState()` to every panel and to
`previewPanel`, and out of `ValidateEditorState` to every caller.

This is **latent, not live** today: every current writer replaces the whole
slice rather than an element
([layoutPanelZones.go](../app/gui/panels/layoutPanelZones.go#L133-L148),
[bonusesPanel.go](../app/gui/panels/bonusesPanel.go#L248),
[editorState.go](../app/gui/models/editorState.go#L116-L119)). One in-place edit
anywhere reintroduces the bug with no compiler or lint signal.

**Fix.** Add an explicit deep copy beside the existing methods in
[editorStateDto.go](../internal/dtos/editorStateDto.go):

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
([editorStateDto.go](../internal/dtos/editorStateDto.go#L338-L346) shows rows
built with `Rules: []models.ContentRuleRowSave{…}`), so a `slices.Clone` of the
row slice is **not** enough — write a `cloneContentRows` helper that also clones
each row's `Rules`. Re-check `editor_state_dto.ManualZoneSave` /
`ManualConnectionSave`: they wrap `entities.Zone` / `entities.Connection`, which
**do** contain slices (`MandatoryContent`, placement rules, roads). ⚠ Those
element types live in the protected
[internal/entities/template/](../internal/entities/template/) tree — do **not**
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

---

### 1.2 🟠 Hub-touching connections are guarded as *player borders*, not as hub borders

**Evidence.** The hub zone is built with the top tier —
[zoneFactory.go](../internal/services/zones/zoneFactory.go#L120-L141):

```go
		Profile: common_zones.GetNeutralZoneProfile(neutral_zone.QualityHighest)
```

but the hub is deliberately **not** in `neutral_zone.Plans`, and every hub
topology therefore substitutes a *player* label as the guard anchor.

[hubTopology.go](../internal/services/template_generator/providers/topology/hubTopology.go#L99-L108):

```go
		hubAnchor := label
		if len(playerLabels) > 0 {
			hubAnchor = playerLabels[0]
		}
		hubGuard := this.GetBorderGuardValue(hubAnchor, label, playerLabels, neutralZones, tuning)
```

[geometricHubTopology.go](../internal/services/template_generator/providers/topology/geometricHubTopology.go#L143-L158):

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

[hubClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/hubClusterService.go#L94):

```go
			WithGuardValue(this.GetBorderGuardValue(playerLabel, spokeLabel, []string{playerLabel}, allNeutralZonePlans, tuning)).
```

`GetBorderGuardValue`
([topologyConnectionService.go](../internal/services/template_generator/providers/topology/base/topologyConnectionService.go#L180-L201))
knows only *player* and *plan* labels, so it can never see the hub.

**Why it is wrong.** With both endpoints resolving to player labels the function
returns `QualityUnknown.GetGuardValue()` = **30 000**
([neutralZoneQuality.go](../internal/models/neutral_zone/neutralZoneQuality.go#L14-L30));
hub-to-neutral returns the *neutral's* tier (as low as 10 000). The hub is the
richest zone on the map (`QualityHighest` = **35 000**), so its borders are
systematically under-guarded — players reach the hub's Platinum content behind a
Bronze/player-grade guard. The bug is invisible in tests because
[getBorderGuardValue_test.go](../test/unit/internal/services/template_generator/providers/topology/base/topologyBase/getBorderGuardValue_test.go#L12-L112)
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
[hubClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/hubClusterService.go#L39)).
Zone labels are `A…AF`
([zoneLabels.go](../internal/common/constants/zoneLabels.go)), so the prefix
cannot collide with a real label.

Then fix the three callers to pass the hub, not a player:

- [hubTopology.go](../internal/services/template_generator/providers/topology/hubTopology.go#L99-L108) — delete `hubAnchor`, call `GetBorderGuardValue(constants.HubZoneName, label, …)`.
- [geometricHubTopology.go](../internal/services/template_generator/providers/topology/geometricHubTopology.go#L143-L160) — delete `guardAnchor`/`guardLabel`, call `GetBorderGuardValue(constants.HubZoneName, label, …)`.
- [hubClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/hubClusterService.go#L94) — pass `hubName` as the first argument instead of `playerLabel`.

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

### 1.3 🟠 Random portal guards ignore endpoint tiers (flat 25 000)

**Evidence.**
[topologyConnectionService.go](../internal/services/template_generator/providers/topology/base/topologyConnectionService.go#L30-L68):

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
[createRandomPortalConnections_test.go](../test/unit/internal/services/template_generator/providers/topology/base/topologyBase/createRandomPortalConnections_test.go#L53-L75).

**Why it is wrong.** A portal that drops a player straight into a Platinum
neutral zone is guarded at 25 000 (the `QualityHigh` value), while the direct
land border into the same zone is guarded at 35 000 — the portal is the cheap
back door. Symmetrically, a portal into a Plastic zone is over-guarded
(25 000 vs. 10 000). Owner decision 2026-08-11: portals use **the same
`max(endpoint qualities)` rule as direct borders**.

**Fix.** Thread the plans through and reuse the (now hub-aware, §1.2)
`GetBorderGuardValue`:

1. Add `neutralZones neutral_zone.Plans` to `CreateRandomPortalConnections` in
   [topologyConnectionServiceInterface.go](../internal/services/template_generator/providers/topology/base/topologyConnectionServiceInterface.go#L10),
   [topologyConnectionService.go](../internal/services/template_generator/providers/topology/base/topologyConnectionService.go#L30-L34)
   and the pass-through in
   [topologyBase.go](../internal/services/template_generator/providers/topology/base/topologyBase.go#L153-L158).
2. Replace the literal with
   `WithGuardValue(this.GetBorderGuardValue(fromLabel, toLabel, playerLabels, neutralZones, tuning))`.
   **Do not double-scale** — `GetBorderGuardValue` already applies
   `tuning.ScaleByBorderGuardStrength`.
3. Update the seven call sites, each of which already has the plans in scope:
   [chainTopology.go](../internal/services/template_generator/providers/topology/chainTopology.go#L47),
   [geometricHubTopology.go](../internal/services/template_generator/providers/topology/geometricHubTopology.go#L57),
   [hubTopology.go](../internal/services/template_generator/providers/topology/hubTopology.go#L46),
   [positionedTopologyBuilder.go](../internal/services/template_generator/providers/topology/positionedTopologyBuilder.go#L57),
   [ringTopology.go](../internal/services/template_generator/providers/topology/ringTopology.go#L45),
   [tournamentTopology.go](../internal/services/template_generator/providers/topology/tournamentTopology.go#L71),
   [webTopology.go](../internal/services/template_generator/providers/topology/webTopology.go#L51).

This removes the `25000` literal, so §3.3 disappears with it.

**Verified non-issue (do not "fix"):** random portals can never touch the hub —
`hubTopology` passes only `outerLabels`
([hubTopology.go](../internal/services/template_generator/providers/topology/hubTopology.go#L39-L48))
and `geometricHubTopology` passes players + plan labels
([geometricHubTopology.go](../internal/services/template_generator/providers/topology/geometricHubTopology.go#L46-L59)).
Hub portals in the geometric layout are created separately and are covered by
§1.2.

**Tests to add / update.**

- Rewrite the guard expectations in
  [createRandomPortalConnections_test.go](../test/unit/internal/services/template_generator/providers/topology/base/topologyBase/createRandomPortalConnections_test.go)
  →  `TestWhenPortalJoinsTwoNeutralZones_UsesTheHigherTierGuardValue`,
  `TestWhenPortalJoinsTwoPlayerZones_UsesThePlayerBorderGuardValue`.
- Golden generator expectations will move for every topology with
  `RandomPortals` enabled — update deliberately.

**Ordering:** land §1.2 first; §1.3 depends on the collapsed
`GetBorderGuardValue`.

---

### 1.4 🟡 A fatal window error is logged to a discard handler, then the process exits silently

*(= review-opus5-08-04 §1.9, restated with re-verified line numbers.)*

**Evidence.** [program.go](../app/gui/program.go#L34-L41):

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
[program.go](../app/gui/program.go#L26), i.e. *before* the loop — installs a
discard logger at [program.go](../app/gui/program.go#L58):

```go
	slog.SetDefault(slog.New(slog.DiscardHandler))
```

Logging is only re-enabled by the opt-in `-with-logging` flag
([program.go](../app/gui/program.go#L69-L72)).

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
([.golangci.yml](../.golangci.yml)) but not `fmt`, so this stays lint-clean.
Secondary, **optional**: `os.Exit` inside the loop skips deferred cleanup —
harmless today (nothing is deferred), but returning an error to
[main.go](../main.go) would make the bootstrap testable. Do not do that as part
of this item unless the owner asks.

**Tests.** No unit test is practical for the Gio event loop. Add an entry for
`app/gui/program.go` to the Gio-UI section of
[test_observations.md](test_observations.md) recording the gap.

---

## 2. Architecture & modelling

### 2.1 🟠 `EditorStateDto` is a flat god-DTO with no entity/model layer

**Evidence.** [editorStateDto.go](../internal/dtos/editorStateDto.go) is a
single struct carrying identity, map, player, neutral-zone, castle, advanced
tier, faction, generation, connectivity, density, topology, victory, arena,
tournament and XP fields, plus banned content, overrides, bonuses, six
mandatory-content row slices, and both manual-edit slices
([editorStateDto.go](../internal/dtos/editorStateDto.go#L15-L94)) — then adds
defaults, layout comparison, castle diffing, manual-edit-insensitive equality
and their private helpers
([editorStateDto.go](../internal/dtos/editorStateDto.go#L97-L447)).

It is imported by **twelve** production packages: `app/gui/drivers`,
`app/gui/editor`, `app/gui/models`, `app/gui/panels`, `internal/dtos`,
`internal/handlers`, `internal/handlers/handler_interfaces`, `internal/mappers`,
`internal/repositories`, `internal/services/editor`,
`internal/services/file_service`, `internal/validators`.

**Why it is wrong.** The DTO is simultaneously the persistence schema
(`.gen.json` tags), the GUI's working state, the validator's target and the
mapper's input, so every one of those concerns is coupled to the others: adding
a field means touching the on-disk format, the panels, the validator and the
mapper at once, and the equality/diff logic (§1.1's aliasing surface) lives in
the same file as the JSON contract. There is no layer that owns the *meaning* of
editor state independently of how it is serialised.

**Fix** (target decided by owner 2026-08-11 — a **new, non-protected** package,
so `internal/entities/template/` is untouched).

1. Create `internal/entities/editor_state/` holding the behaviour-free value
   groups the DTO currently flattens: e.g. `mapSettings.go`, `playerSettings.go`,
   `neutralZoneSettings.go`, `castleSettings.go`, `generationSettings.go`,
   `gameRuleSettings.go`, `contentSettings.go`, `manualEditSettings.go`
   (file-per-struct, camelCase names, AGENTS.md §4.1).
2. Add `internal/models/editorStateModel.go` — a model that **wraps** those
   entities and owns the behaviour currently on the DTO
   (`LayoutDefiningOptionsChanged`, `DiffCastleSettings`,
   `EqualsIgnoringManualEdits`, `HasManualEdits`, the `Clone` from §1.1, the
   defaults factory).
3. Reduce `dtos.EditorStateDto` to the model **embedded** plus the JSON tags,
   so the serialisation contract is the only thing left in `dtos`.

**⚠ Do this as a plan file under `plans/`, not in one pass** (AGENTS.md §2.4) —
twelve packages and the `.gen.json` round trip are at stake. Suggested phasing:
extract entities → introduce the model with the behaviour moved (DTO delegates)
→ embed the model in the DTO → migrate consumers package by package → delete the
delegation shims.

**Non-negotiable invariant:** the on-disk `.gen.json` field names and shape must
not change. Every phase must keep
`test/integration/` load/save round trips green, and the phase that touches the
DTO must add a golden-file test: write a `.gen.json` with the pre-change code,
load it with the post-change code, assert every field survives.

**Tests.** Each new entity/model file gets its mirrored folder under
`test/unit/`; the behaviour tests currently under
`test/unit/internal/dtos/editorStateDto/` move with the methods.

**Blocked by:** §1.1 should land first — `Clone` belongs on the model, and
writing it twice is wasted work.

---

### 2.2 🟡 ⚠ Zone tier has no single source of truth

**Evidence.** `neutral_zone.Quality`
([neutralZoneQuality.go](../internal/models/neutral_zone/neutralZoneQuality.go#L3-L12))
is the tier enum. During generation the tier is explicit on
`neutral_zone.Plan`
([neutralZonePlan.go](../internal/models/neutral_zone/neutralZonePlan.go#L11-L15)),
but the finished `entities.Zone` carries **no tier**, so
[zoneClassifier.go](../internal/services/zones/zoneClassifier.go#L23-L177)
re-derives it from layout + resource pools + guarded/unguarded pool IDs across
three content-inference branches.

Consumers that all depend on that re-derivation:
[previewLayoutService.go](../internal/services/preview_service/previewLayoutService.go#L103-L116),
[zoneEditorZoneProps.go](../app/gui/dialogs/zoneEditorZoneProps.go#L61-L64),
[zoneEditorService.go](../internal/services/connection_editor/zoneEditorService.go#L187-L222),
[manualReapplyService.go](../internal/services/connection_editor/manualReapplyService.go#L88-L124),
[mandatoryContentProvider.go](../internal/services/template_generator/providers/mandatoryContentProvider.go#L88-L115),
[gladiatorArenaProvider.go](../internal/services/template_generator/providers/gladiatorArenaProvider.go#L85-L111),
[zoneEditorHandler.go](../internal/handlers/zoneEditorHandler.go#L58-L81),
[connectionEditorService.go](../internal/services/connection_editor/connectionEditorService.go#L24-L39).

**Why it is wrong.** Tier is decided once at plan time, thrown away, then
reconstructed by pattern-matching on content. Every feature that edits a zone's
content (the manual editor's re-profile action, mandatory-content regeneration)
can silently flip the inferred tier, and the inference rules must stay in lockstep
with the profile catalogue in
[neutralZoneProfile.go](../internal/common/common_zones/neutralZoneProfile.go#L10-L25)
with nothing enforcing that.

**Fix — two branches, owner picks:**

- **Branch A (owner's first idea, ⚠ protected dir).** Add a runtime-only
  `Quality` field to `entities.Zone`
  ([zone.go](../internal/entities/template/template_variant/zone.go)) tagged
  `json:"-"` exactly like the existing `GeneratorPosition`, set it in
  `ZoneFactory` at creation, and reduce `ZoneClassifier` to a fallback for zones
  that arrive without it (loaded manual snapshots). **Requires explicit owner
  approval** — it edits the protected `.rmg.json` schema package, even though a
  `json:"-"` field cannot change the output file. Add a round-trip test proving
  the emitted `.rmg.json` is byte-identical before and after.
- **Branch B (owner's second idea, no protected edits).** Keep `entities.Zone`
  as is and make `neutral_zone.Profile`
  ([neutralZoneProfile.go](../internal/models/neutral_zone/neutralZoneProfile.go#L3-L20))
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

### 2.3 🟡 Preview geometry is integer-only although a `Vec2` already exists

**Evidence.** [previewLayout.go](../internal/models/preview/previewLayout.go#L5-L11):

```go
type Layout struct {
	Positions   map[string]image.Point
	Zones       []Zone
	Connections []Connection
	ZoneRadius  int
}
```

Positions are rounded at
[layoutGeometry.go](../internal/services/preview_service/layoutGeometry.go#L82-L93)
(`image.Pt(int(math.Round(px[i])), int(math.Round(py[i])))`), and again
independently in
[layoutRingHub.go](../internal/services/preview_service/layoutRingHub.go#L63-L69)
and
[layoutBalancedRings.go](../internal/services/preview_service/layoutBalancedRings.go#L119-L129);
`canvasMetrics.center()` **truncates** rather than rounds
([layoutGeometry.go](../internal/services/preview_service/layoutGeometry.go#L65-L67)).
The integers then propagate into the editor:
[zoneEditorGeometryService.go](../internal/services/connection_editor/zoneEditorGeometryService.go#L73-L79)
copies them straight into `models.ZoneEditorGeometry`, and hit-testing works in
`image.Point`
([zoneEditorGeometryService.go](../internal/services/connection_editor/zoneEditorGeometryService.go#L81-L96)).

Meanwhile a generic float vector already exists —
[vec2.go](../internal/helpers/data/vec2.go#L8-L69) (`Vec2[T]`, `Vec2FromPoint`,
`ToPoint`, `ToPointRounded`) — and `models.Position = data.Vec2[float64]`
([position.go](../internal/models/position.go#L5-L10)) is what the topology
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
   [zoneEditorGeometryService.go](../internal/services/connection_editor/zoneEditorGeometryService.go)
   (`HitTestNode`, `HitTestEdge`, `buildEdges`) — they already compute in
   `float64` internally via `math.Hypot`, so this mostly deletes conversions.
4. **Round only at the final draw call**: the PNG renderer
   ([previewGeneratorService.go](../internal/services/preview_service/previewGeneratorService.go#L51-L72))
   and the Gio canvases
   ([previewPanel.go](../app/gui/panels/previewPanel.go#L170-L239),
   [zoneEditorCanvas.go](../app/gui/dialogs/zoneEditorCanvas.go#L92-L207)).
   Use `Vec2.ToPointRounded()` / `f32.Pt` there — the canvas already draws
   Beziers in `f32.Point`, so step 4 in the editor is mostly *removing* the
   round-trip through `image.Point`.

**⚠ Protected boundary:** `Zone.GeneratorPosition *[2]float64`
([zone.go](../internal/entities/template/template_variant/zone.go#L7-L23))
stays exactly as it is for this item — converting it is §2.4 and needs owner
approval. Convert at the read site (`generatorCoords`,
[layoutGeometry.go](../internal/services/preview_service/layoutGeometry.go#L95-L106)).

**Tests.** The preview layout services already have mirrored unit folders under
`test/unit/internal/services/preview_service/` — update the expected values and
add `TestWhenTwoZonesAreLessThanAPixelApart_TheirCentresDiffer`. The GPU-gated
snapshot suite
([window_snapshot_integration_test.go](../test/integration/gui/window_snapshot_integration_test.go))
will need `-update`; **the owner must eyeball the regenerated snapshots**, since
sub-pixel changes are exactly what that suite is meant to catch.
The numeric geometry pins in
[zoneEditorGeometry_integration_test.go](../test/integration/gui/zoneEditorGeometry_integration_test.go#L138-L163)
will shift — update them deliberately, do not relax them to `InDelta`.

---

### 2.4 ⚪ ⚠ Replace `[2]float64` with `Vec2` in the template entities

**Evidence.** [zone.go](../internal/entities/template/template_variant/zone.go#L7-L23)
declares `GeneratorPosition *[2]float64` (`json:"-"`), and producers stamp it
as an array literal —
[geometricHubTopology.go](../internal/services/template_generator/providers/topology/geometricHubTopology.go#L104-L119):

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
[internal/entities/template/](../internal/entities/template/) is read-only under
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
[vec2.go](../internal/helpers/data/vec2.go), which has no such dependency.
Then delete the pack/unpack at every producer and consumer.

**Blocked by:** §2.3 — do that first so the preview side is already float-native
and this becomes a mechanical type swap.

**Tests.** No behaviour change is expected. The proof obligation is a golden
test: generate a template before and after, assert the `.rmg.json` bytes are
identical.

---

### 2.5 ⚪ ⚠ Move `entities/types.go` under `template/` and rename `template` → `template_entity`

**Evidence.** [types.go](../internal/entities/types.go) is a single
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
[internal/entities/template/](../internal/entities/template/) and moves a file
*into* that tree. AGENTS.md §2.1 forbids both without explicit approval. There
is no functional benefit, only naming; **do not start without a go-ahead.**

**Fix, if approved.**
1. `git mv internal/entities/types.go internal/entities/template/types.go` and
   change its package clause; the aliases then reference sibling subpackages.
2. Rename the directory `internal/entities/template` →
   `internal/entities/template_entity` and update the `package template` clause
   in [rmgTemplate.go](../internal/entities/template/rmgTemplate.go).
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

## 3. Duplicated values that belong in `internal/common`

Owner decision 2026-08-11: **all four families move**. These are small,
independent, and safe to batch together.

Already centralised and correctly used — do **not** re-report: guard weekly
increments
([guardWeeklyIncrement.go](../internal/common/common_connections/guardWeeklyIncrement.go)),
guard-strength presets
([guardStrength.go](../internal/common/common_connections/guardStrength.go)),
content-distance presets
([distancePresets.go](../internal/common/common_distances/distancePresets.go)),
zone names/prefixes and labels
([zoneNames.go](../internal/common/constants/zoneNames.go),
[zoneLabels.go](../internal/common/constants/zoneLabels.go)), map sizes
([mapSizes.go](../internal/common/mapSizes.go)), permissions
([permissions.go](../internal/common/constants/permissions.go)).

### 3.1 🟡 `"mandatory_content_hub"` is repeated in four production files

**Evidence.**
[mandatoryContentProvider.go](../internal/services/template_generator/providers/mandatoryContentProvider.go#L120)
and [#L146](../internal/services/template_generator/providers/mandatoryContentProvider.go#L146),
[hubTopology.go](../internal/services/template_generator/providers/topology/hubTopology.go#L77),
[geometricHubTopology.go](../internal/services/template_generator/providers/topology/geometricHubTopology.go#L94),
[hubClusterService.go](../internal/services/template_generator/providers/topology/tournament_variant/hubClusterService.go#L66).

**Why it is wrong.** The string is a cross-file contract: the topology writes it
into `Zone.MandatoryContent`, the content provider must emit a group with the
same name, and the parallel C# editor reads it
([mandatoryContentProvider.go](../internal/services/template_generator/providers/mandatoryContentProvider.go#L132)
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
(e.g. [createContents_test.go](../test/unit/internal/services/template_generator/providers/mandatoryContentProvider/createContents_test.go#L233)) —
leave the *test-side* literals as literals so they still catch an accidental
change to the constant's value. Add
`test/unit/internal/common/constants/mandatoryContentNames/` only if a helper
function (not a bare constant) ends up being introduced.

### 3.2 🟡 `"mandatory_content_neutral_"` / `"mandatory_content_side_"` prefixes repeated

**Evidence.** Neutral prefix:
[mandatoryContentProvider.go](../internal/services/template_generator/providers/mandatoryContentProvider.go#L61),
[#L106](../internal/services/template_generator/providers/mandatoryContentProvider.go#L106),
[topologyBase.go](../internal/services/template_generator/providers/topology/base/topologyBase.go#L120).
Side prefix:
[mandatoryContentProvider.go](../internal/services/template_generator/providers/mandatoryContentProvider.go#L48),
[#L92](../internal/services/template_generator/providers/mandatoryContentProvider.go#L92),
[zoneFactory.go](../internal/services/zones/zoneFactory.go#L68).

**Fix.** Same constants file as §3.1. Consider two tiny helpers
(`NeutralMandatoryContentName(label string) string`,
`SideMandatoryContentName(label string) string`) so the concatenation itself is
written once; if so, they get a mirrored unit-test folder per AGENTS.md §4.6.

### 3.3 🟡 Portal guard `25000` magic literal

**Evidence.**
[topologyConnectionService.go](../internal/services/template_generator/providers/topology/base/topologyConnectionService.go#L64):
`WithGuardValue(tuning.ScaleByBorderGuardStrength(25000))`.

**Disposition: superseded by §1.3.** Landing §1.3 deletes the literal outright.
Only if §1.3 is deferred should this be done on its own, as a named constant in
[common_connections](../internal/common/common_connections/) — note the value is
numerically identical to `neutral_zone.QualityHigh.GetGuardValue()`, so name it
for its meaning, not its number.

### 3.4 🟡 Foothold placement distances are inline literals

**Evidence.**
[mandatoryContentProvider.go](../internal/services/template_generator/providers/mandatoryContentProvider.go#L161-L189):

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
[distancePresets.go](../internal/common/common_distances/distancePresets.go).
Nothing tells a future maintainer what `0.5/0.5` means or that it is meant to
differ from `Far`.

**Fix.** Add a `GetFootholdDistancePresets()` (or three named presets) to
[common_distances](../internal/common/common_distances/), following the shape of
`GetContentDistancePresets`, and reference them from the provider.
**Do not silently reuse the existing content presets** — `0.2–0.3` and
`0.2–0.4` do not match any of `Next To`/`Near`/`Medium`/`Far`/`Very Far`, so
folding them in would change generation.

**Tests.** `test/unit/internal/common/common_distances/distancePresets/` gains a
test per new accessor; the existing
`test/unit/internal/services/template_generator/providers/mandatoryContentProvider/`
tests must keep asserting the *numeric* values so a constant rename cannot
silently retune placement.

---

## 4. UI / UX

### 4.1 🟡 "Save As" is really "Save To" — the UI offers a filename it then discards

**Evidence.** Writing editor state as `{TemplateName}.gen.json` is **intended**
(review §1.1, owner-approved):
[fileService.go](../internal/services/file_service/fileService.go#L41-L48)
passes `filepath.Dir(filePath)` and `editorState.TemplateName` to the
repository, and
[saveSettings_test.go](../test/unit/internal/services/file_service/fileService/saveSettings_test.go#L14-L60)
pins that behaviour deliberately.

The defect is the UI. The dialog still offers an editable **name** field whose
value is silently dropped —
[fileExplorerDialogToolbar.go](../app/gui/dialogs/fileExplorerDialogToolbar.go#L35-L51):

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

**Fix** (UI treatment confirmed by the owner 2026-08-11: **read-only textbox**,
not a label — the value is filename-shaped and belongs in a field).

- [fileExplorerDialogToolbar.go](../app/gui/dialogs/fileExplorerDialogToolbar.go#L35-L51) —
  `getSaveRowWidget`: pass `readonly = true` to
  [`NewTextboxWidget`](../app/gui/widgets/textboxWidget.go#L17-L38) (the
  parameter already exists and already sets `textEditor.Editor.ReadOnly`), and
  relabel `"Save as:"` → `"Will save as:"` so the row reads as a preview of the
  resolved name. The field must be **populated with the resolved name**, not a
  hint, since the user can no longer type it.
- [toolbar.go](../app/gui/editor/toolbar.go#L13-L73) — button `"Save As"` →
  `"Save To"`, field `buttonSaveAs` → `buttonSaveTo`.
- [stateFiles.go](../app/gui/drivers/stateFiles.go#L26-L42) — `State.SaveAs` →
  `State.SaveTo`, including the `Save` fallback at
  [#L27-L30](../app/gui/drivers/stateFiles.go#L27-L30).
- [fileExplorerDialogModes.go](../app/gui/dialogs/fileExplorerDialogModes.go#L29-L42) —
  dialog title `"Save File"` → `"Save To"`. **`NewSaveFileDialog`,
  `modeSaveFile` and `onSave` stay** (owner decision): they name the *explorer
  mode and callback*, not the toolbar action.
- Docs: [QUICKSTART.md](../QUICKSTART.md#L48) and
  [QUICKSTART.md](../QUICKSTART.md#L109), [README.md](../README.md#L154).

**Tests.**

- Rename `test/integration/stateSaveAs_integration_test.go` →
  `stateSaveTo_integration_test.go` with its helper `newSaveAsProbe` and both
  test names (`TestWhenSaveAsFails_CurrentPathIsNotRecorded`,
  `TestWhenSaveAsSucceeds_CurrentPathIsRecorded`). It uses the `ConfirmSave`
  testexport, not `SetFilename`, so it needs no behavioural change.
- Rename `test/unit/app/gui/drivers/stateFiles/saveAs_test.go` →
  `saveTo_test.go` (`TestWhenSaveAsIsCalled_DialogIsOpened`), and fix the name
  reference in
  [save_test.go](../test/unit/app/gui/drivers/stateFiles/save_test.go#L11).
- [fileExplorerDialog_integration_test.go](../test/integration/gui/fileExplorerDialog_integration_test.go#L73-L119):
  `TestWhenSaveDialogIsConfirmed_TheTypedNameBecomesTheSaveTarget` and
  `TestWhenSaveDialogIsConfirmedThroughTheDriver_AFileLandsInTheChosenDirectory`
  both drive `SetFilename`; rewrite them as "the field shows the resolved name
  and cannot be edited" / "the file lands in the chosen directory under the
  resolved name".
  `TestWhenTheFilenameIsWhitespaceOnly_TheConfirmButtonIsDisabled` and
  `TestWhenTheFilenameIsValid_TheConfirmButtonIsEnabled` need a **new trigger**
  — the user can no longer produce a whitespace filename. **If investigation
  shows** the whitespace state is now unreachable through any public path,
  delete the disabled-case test and record the removed branch in
  [test_observations.md](test_observations.md) rather than adding a test-only
  seam.
- [fileExplorerDialog_testexports.go](../app/gui/dialogs/fileExplorerDialog_testexports.go#L5-L45):
  `SetFilename` exists only to drive an editable field — remove it, or replace
  it with a read accessor. `ConfirmSave` stays.

---

## 5. Testing

### 5.1 🟠 Zone editor pointer flows are untested (drag-to-connect, zone drag + snapping)

**Evidence.** [test_observations.md](test_observations.md) records for the zone
editor: *"Still uncovered: … the pointer flows (drag-to-connect, zone drag +
snapping), which need synthetic pointer events - the test/performance AppRunner
pattern is the way in, and it is still future work."*

That seam now exists:
[appRunner.go](../test/test_helpers/integration_common/appRunner.go#L152-L217)
provides `ClickAt`, `MoveTo`, `DragTo(from, to image.Point)` (press +
8 interpolated moves + release) and `InputText`. `DragTo` is already used to
drive `gesture.Drag` widgets.

**Why it matters.** Drag-to-connect and zone-drag-with-snapping are the two
interactions where the manual zone editor actually changes user data. The
supporting geometry is unit-tested at ≥ 92.9 % in
[internal/services/connection_editor](../internal/services/connection_editor/),
but nothing proves the *dialog wires the pointer to it correctly* — a swapped
argument or a missing mode check would ship silently.

**Fix.** Add scenarios to
[zoneEditorDialog_integration_test.go](../test/integration/gui/zoneEditorDialog_integration_test.go)
(or a sibling `zoneEditorPointer_integration_test.go`), tagged
`//go:build integration_test && gui` since they drive real frames:

- `TestWhenAZoneIsDraggedToANewPosition_TheAppliedLayoutRecordsIt`
- `TestWhenAZoneIsDraggedNearAGuide_ItSnapsToTheGuide`
- `TestWhenADragStartsOnAZoneInAddConnectionMode_AConnectionIsCreated`
- `TestWhenADragEndsOnEmptyCanvas_NoConnectionIsCreated`

Read results back through the existing
`zoneEditorDialog_testexports.go` accessors and `State.ApplyEditedZones`; use
the canvas coordinates already pinned by
[zoneEditorGeometry_integration_test.go](../test/integration/gui/zoneEditorGeometry_integration_test.go#L138-L163)
so the tests do not hard-code fresh magic points.

**⚠ Ordering:** if §2.3 (float geometry) is scheduled, land it **first** — these
tests pin canvas coordinates and would need rewriting otherwise.

**Then:** delete the "still future work" sentence from
[test_observations.md](test_observations.md) and replace it with a pointer to
the new tests.

### 5.2 🟡 Zone editor property panels (`widget.Editor` / dropdown paths) are untested

**Evidence.** Same entry in [test_observations.md](test_observations.md):
*"Still uncovered: the property panels' `widget.Editor`/dropdown paths…"*.
`AppRunner.InputText` plus `ClickAt` now cover both input kinds, and
[pickerDialog_integration_test.go](../test/integration/gui/pickerDialog_integration_test.go)
already demonstrates the dropdown-click pattern.

**Why it matters.** The zone/connection property panels are where a user renames
a zone, changes its castle count and re-profiles its tier — the exact paths
§2.2 will touch.

**Fix.** New `test/integration/gui/zoneEditorProperties_integration_test.go`
(`//go:build integration_test && gui`):

- `TestWhenAZoneNameIsTyped_TheAppliedZoneCarriesTheNewName`
- `TestWhenAZoneQualityIsPickedFromTheDropdown_TheZoneIsReprofiled`
- `TestWhenAConnectionGuardValueIsTyped_TheAppliedConnectionCarriesIt`
- `TestWhenANonNumericGuardValueIsTyped_TheValueIsRejected`

Focus the field with `ClickAt` before `InputText` (the helper's doc comment
requires it). Assert through `State.ApplyEditedZones` output, not through
dialog internals.

### 5.3 🟡 File explorer: hidden-file toggle and pointer-driven row/scroll interactions

**Evidence.** [test_observations.md](test_observations.md) for the file
explorer: *"Still uncovered: the hidden-file toggle and the pointer-driven
row/scroll interactions (owner decision - excluded from the scenario set)."*
The owner re-opened this on 2026-08-11.

The toggle exists at
[fileExplorerDialogToolbar.go](../app/gui/dialogs/fileExplorerDialogToolbar.go#L31)
(`widgets.NewToggleButtonWidget(theme, "Show hidden", &this.hiddenToggle, this.showHidden)`),
and the underlying filtering policy is already unit-tested in
`internal/services/file_system`.

**Why it matters.** The policy is tested; the *wiring* is not. A toggle that
flips a field but never re-lists, or a row click that resolves the wrong entry,
would pass every existing test.

**Fix.** Extend
[fileExplorerDialog_integration_test.go](../test/integration/gui/fileExplorerDialog_integration_test.go):

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

### 5.4 🟡 The GUI handler framework is a design comment, not a design

**Evidence.** [runnerHandler.go](../test/test_helpers/integration_common/runnerHandler.go#L61-L106)
carries a 46-line first-person design note above a struct that currently
implements three tab clicks and one mask. The whole framework today is:

```go
func (this *baseHandler) ClickGeneralTab() *baseHandler {
	this.runner.ClickAt(f32.Pt(672, 60))
	this.runner.VerifySnapshot()
	return this
}
```

Consumers: [window_snapshot_integration_test.go](../test/integration/gui/window_snapshot_integration_test.go#L20-L25)
and [window_tab_cycling_test.go](../test/performance/window_tab_cycling_test.go#L25).

**Why it matters.** Gio has no Playwright/Selenium, so this handler *is* the
GUI test API. §5.1, §5.2 and §5.3 will each be written against it; if they land
first, three more test files acquire their own ad-hoc `f32.Pt` literals and the
framework never happens. The comment is also the only place this design lives —
a doc comment is not a backlog and will rot.

**Fix — seven separable pieces, in this order.**

**(a) Hygiene (mechanical, do first).** The file is `runnerHandler.go` but the
struct is `baseHandler`; AGENTS.md §4.1 requires `baseHandler.go`. `NewHandler`
also returns the unexported `*baseHandler` from an exported function — either
export the type or (once (d) exists) return the handler contract. Interfaces
follow AGENTS.md §4.2.1/§4.2.2: with fewer than five implementations they stay
in this package as `*Interface.go` files.

**(b) Settle the coordinate strategy.** Recommendation: **keep the coordinates
literal but stop scattering them.** Gio exposes no widget-rect lookup without a
new `*_testexports.go` seam, and computing positions from the layout code would
re-implement the thing under test. Put them in one `handlerCoordinates.go` as
named constants derived where possible from
[app/gui/constants/ui.go](../app/gui/constants/ui.go) (`DefaultPadding`,
`DefaultLabelWidth`, `DefaultPreviewWidthMaximum`), so a padding change is one
edit rather than a hunt. Only compute at runtime where a value genuinely varies
(slider value → x, list row index → y).

**(c) Narrow the mask.** `setRandomTopology` blanks
`image.Rect(WindowWidth-470, 0, WindowWidth, WindowHeight)` — a 470 px column
over the **full window height**, i.e. roughly a third of every snapshot,
including stable chrome. Replace the single catch-all with named helpers that
mask only what is genuinely nondeterministic: the preview canvas (random
topology), the output-directory field (differs per machine/OS/Steam layout —
see AGENTS.md §2.7), and the timestamp/path inside the status message. Derive
the preview rect from `DefaultPreviewWidthMaximum` (440) rather than the
unexplained 470.

**(d) Per-tab and per-dialog handlers.** `baseHandler` gains
`ClickGeneralTab() *generalTabHandler` etc.; tab handlers embed `baseHandler`
(tab bar and toolbar stay live), dialog handlers **do not** (a dialog disables
the background, so inheriting its clicks would be a lie). One struct per file
(AGENTS.md §4.1).

**(e) Layout-shifting state.** A handler must remember what it toggled, because
some widgets move the ones below them: *Allow non-official larger map sizes*
expands the Map Size dropdown, and dropdowns render inline on the canvas rather
than floating. Store those flags on the handler and offset subsequent
coordinates from them; a handler that clicks blind will click the wrong widget.

**(f) Scrolling.** `AppRunner` has `ClickAt`, `MoveTo`, `DragTo` and
`InputText` — but no scroll. Anything below the fold is unreachable, which is a
hard blocker for §5.2. Add `Scroll(point f32.Point, delta f32.Point)` injecting
a `pointer.Scroll` event through the same router, mirroring `ClickAt`.

**(g) Keep `*_testexports.go` out of the handler.** `SelectedTabIndex`,
`TabCount`, `DialogsOpen`, `CurrentState` and `Status` are test-only exports.
Handlers should drive pixels and assert through snapshots; reach for an
accessor only where a pixel genuinely cannot express the assertion, and say why
in a one-line comment.

**Scope guard.** Do **not** build (d)–(g) speculatively. (a)–(c) are worth
doing on their own; the rest grows one method at a time as §5.1–§5.3 need it.
The comment's own worry — *"this looks like an entire framework"* — is the
correct instinct: an unused handler method is untested test code.

**When it lands,** replace the design comment with a two-line pointer to this
section so the design lives in one place.

### 5.5 🟡 GUI snapshots differ between local and CI, and the tolerance hiding it also hides regressions

**Evidence.** The last paragraph of the same comment: *"for some reason there
is a difference of some of the rendered text between local … and CI (looks like
some of the text is grayed out like not finishing a rerender)"*. The tolerance
that absorbs it is [comparer.go](../test/test_helpers/integration_common/snapshot/comparer.go#L8-L11):

```go
// DefaultSnapshotThreshold is the maximum allowed normalized mean color
// distance between a golden snapshot and an actual screenshot (2%).
// Pipeline has discrepancies, I don't want to investigate them right now.
const DefaultSnapshotThreshold = 0.02
```

The two environments are genuinely different: goldens are generated locally on
Windows against a real GPU, while
[pr-validation.yml](../.github/workflows/pr-validation.yml#L244-L257) runs the
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
   [appRunnerSnapshots.go](../test/test_helpers/integration_common/appRunnerSnapshots.go#L120),
   not in the threshold. Download the `gui-snapshot-failures` artifact the
   workflow already uploads to see the actual CI pixels.
2. **If it is genuinely llvmpipe text anti-aliasing**, make CI the reference:
   regenerate goldens in the CI environment (run the update job, download the
   artifact, commit the images) so the *comparison* is exact and only local
   runs are lenient.
3. **Then replace the metric.** A mean is the wrong measure for UI diffs. Fail
   on the fraction of pixels whose per-pixel distance exceeds a small
   per-pixel tolerance (e.g. "> 0.5 % of pixels differ by > 10 %") so wide
   faint AA noise passes while a small solid change fails. Keep
   `Comparer.Threshold` configurable and update
   [test/unit/test/test_helpers/integration_common/snapshot/comparer/](../test/unit/test/test_helpers/integration_common/snapshot/comparer/)
   alongside.

**If investigation shows** the difference is unavoidable and CI-generated
goldens are impractical, say so in the constant's comment with the measured
worst-case difference — a documented 2 % is fine, an unexplained one is not.

---

## 6. Deferred decisions

### 6.1 ⚪ `createTopologyAdjacency` dead Chain/Ring branches

**Evidence.**
[zoneLabelProvider.go](../internal/services/zones/zoneLabelProvider.go#L212)
declares `createTopologyAdjacency`; its only call site is
[#L82](../internal/services/zones/zoneLabelProvider.go#L82). The
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
| §1.2 / §1.3 guard values | Changes generated `.rmg.json` guard numbers for every hub topology and every portal | Approved 2026-08-11 |

---

## 8. Suggested execution order

Bugs first, protected-directory work last, nothing batched with something it
blocks. Each batch is one PR-sized unit; the owner reviews and commits.

| Batch | Items | Notes |
| --- | --- | --- |
| **A** | §1.4 | Two-line stderr fix + a `test_observations.md` entry. Smallest possible warm-up, no dependencies. |
| **B** | §1.2 → §1.3 | Strictly ordered: §1.3 depends on the collapsed `GetBorderGuardValue`. Expect golden-template churn; regenerate expectations deliberately. |
| **C** | §3.1, §3.2, §3.4 | Constant extraction, mechanical, no behaviour change. §3.3 is already gone if B landed. |
| **D** | §1.1 | Deep copy + regression tests. Benchmark before/after. |
| **E** | §4.1 | Save To rename. Touches the file-explorer integration tests — do **before** §5.3. |
| **F** | §5.3 | File-explorer pointer/hidden-file tests, in the file §4.1 just rewrote. |
| **G** | §2.3 | Float preview geometry. Regenerates GPU snapshots — owner review required. Do **before** §5.1. |
| **H** | §5.1, §5.2 | Zone-editor pointer + property-panel tests, against the post-§2.3 coordinates. |
| **I** | §2.1 | `EditorStateDto` rework. **Needs a `plans/` file** (AGENTS.md §2.4) — multi-phase, twelve packages. Depends on §1.1 for `Clone`. |
| **J** | §2.2 Branch B | Zone tier single source of truth without a protected edit. Benefits from §2.1's model layer. |
| **⚠ K** | §2.2 Branch A, §2.4, §2.5, §6.1 | Owner-gated. Do not schedule until each is explicitly approved. §2.4 depends on §2.3. |
| **L** | §5.4 (a–c), §5.5 | GUI test-harness groundwork: handler hygiene, named mask helpers, coordinate constants, snapshot-comparison fix. Slot **before F**, so F/H write against the settled shape and a comparison that can actually fail. |
| **M** | §5.4 (d–g) | Grows with F and H — the scroll seam (f) is a hard prerequisite for §5.2. Never build ahead of the test that needs it. |

**Note on L/M.** They are listed last only because the table is otherwise
ordered by dependency; L should actually run before batch F. §5.5 step 1 can be
done at any time — it needs a CI artifact, not a code freeze.

**Coverage note.** Batches C, D, F, H and J add tests; B, E and G mostly move
existing behaviour. Run the coverage task before and after **every** batch
(AGENTS.md §2.3) — the floor is **69.3 %** and the current figure is **72.5 %**.

---

## 9. Baselines to hold

| Gate | Command | Expected |
| --- | --- | --- |
| Build | `go build ./...` | clean |
| Vet (default) | `go vet ./...` | clean |
| Vet (gated) | `go vet -tags='integration_test,gui' ./...` | clean |
| Test layout | `go run ./cmd/testlayoutcheck .` | `test-layout check passed` |
| Unit | `go test ./test/unit/... -count=1` | pass |
| Integration | `go test -tags=integration_test ./test/integration/... -count=1` | pass |
| GUI integration | `go test -tags='integration_test,gui' ./test/integration/gui/... -count=1` | pass (needs GPU) |
| Coverage | `go test -count=1 '-coverpkg=./internal/...,./app/...' '-coverprofile=coverage.txt' ./test/unit/...` then `go tool cover '-func=coverage.txt'` | **≥ 69.3 %**, currently **72.5 %** |
| Lint | `golangci-lint-v2 run ./... --issues-exit-code=0` | **0 issues** |
| Format | `gofmt -l ./app ./internal ./test ./cmd` | empty |
| Wire | `wire diff ./internal/composition/...` | no diff |
