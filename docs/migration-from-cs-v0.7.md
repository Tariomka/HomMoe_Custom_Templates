# Migration spec: C# editor v0.7 → Go port

This document maps the changes in
`Olden Era - Template Editor` (C# WPF source, currently at tag **v0.7**,
commit `6c97732`) onto the Go port in this repository. The last Go sync
roughly tracked the C# tree as of **2026-05-12** (just before commit
`6dd070e` "Added balanced map layout"). Everything in this document
describes work needed to bring the Go port back to parity with v0.7.

User decisions (already taken):

- **No backward compatibility.** Drop old Go-only fields. `.gen.json`
  files written by previous Go builds may load with some content missing.
- **Mirror the C# test suite.** Each Go phase should land with the
  equivalent of the matching changes in
  `Olden Era - Template Editor.Tests/TemplateGeneratorTests.cs`.

Throughout the doc, file references use this convention:

- `cs:<path>` — file under `D:\Git\some_personal_shit\Olden-Era---Template-Generator\Olden Era - Template Editor\`
- `go:<path>` — file under this repo (`d:\Git\some_personal_shit\HomMoe_Custom_Templates\`)
- "commit X" without a repo prefix means the C# repo.

---

## 0. Commit roll-up (since last sync)

C# commits to port, grouped by area. Each section below references these
hashes; you can `git show <hash>` in the C# repo for the exact diff.

| Area | Commits |
| --- | --- |
| **Data model / persistence** | `5f7cca5`, `ccf97c4`, `fc4cb46` |
| **MapTopology.Balanced** (first-class) | `6dd070e`, `0e31576`, `2d4cace`, `e268446`, `02bb65d`, `047370a`, `d2166eb` |
| **Zone content manager refactor** | `5f7cca5`, `ccf97c4`, `9151ffd`, `48deacc`, `7029355` |
| **Tournament layout** (chain/ring/balanced) | `6d28ef0`, `a040c98`, `df8a115`, `cd3a46e`, `ab4c1d3`, `6606755`, `e91e79f`, `c20b40d` |
| **Random map fixes** | `3f46209`, `cc5d7ee`, `c63ea67` |
| **Wood/Ore bonus + bonus picker** | `4565e6c`, `fb3c203`, `16580b3` |
| **Guard scaling / fairness** | `89cba4b`, `2c62fab`, `de73435` |
| **Preview** | `432eeb9`, plus preview hunks of all balanced/tournament commits |
| **UI refactor & naming** | `ca7556b`, `fc4cb46`, `30a235f`, `fb3c203` |
| **Tests** | `7303848`, `cc63ecc`, `c400677` |

---

## 1. Recommended order

Each phase is self-contained, ships with its own tests, and leaves the
program building and the existing test suite green. Do not skip ahead;
later phases depend on data model / enum changes added in early phases.

1. **Phase 1 – Persistence parity** (`SettingsFile`, `ZoneContentRowSave`,
   bonus storage as a JSON string). Mechanical; no behavior change.
2. **Phase 2 – `MapTopology.Balanced` enum + UI label.** Dispatch
   temporarily falls through to existing logic. Required prerequisite
   for the balanced algorithm.
3. **Phase 3 – `ZoneContentManager` refactor.** Drop the hardcoded
   "preset" content; builders now return foothold + user-configured rows.
   Requires phase 1 fields.
4. **Phase 4 – Neutral / Hub mandatory-content UI.** Surfaces the four
   new content sets in the Zone Content tab. Requires phases 1 and 3.
5. **Phase 5 – Bonus picker UI + wood/ore bonus types.** Bonuses must
   round-trip through `BonusesJson` from phase 1.
6. **Phase 6 – Algorithm port (TemplateGenerator).** The big one.
   Split into 6a balanced layout, 6b tournament balanced cluster,
   6c connection fairness / zone sorting fixes, 6d guard scaling.
7. **Phase 7 – Preview improvements.** Hub-castle count badge,
   balanced/tournament preview layouts.
8. **Phase 8 – Test parity.** Port any test cases not already mirrored.

---

## 2. Phase 1 — Persistence parity

### 2.1 New file: `internal/models/zoneContentRowSave.go`

Mirror `cs:Models/Generator/ZoneContentRowSave.cs`. Lightweight record
preserving one mandatory-content row exactly as configured in the UI.

```go
package models

// ZoneContentRowSave is the serialisation record for a single mandatory-
// content UI row. Preserves the row exactly as the user configured it,
// including Count — so two separate sawmill rows stay as two rows after
// a round-trip.
type ZoneContentRowSave struct {
    SID          string `json:"sid"`                    // SID of content item or include-list
    Count        int    `json:"count"`                  // spinner / count value (default 1)
    IsGroup      bool   `json:"isGroup"`                // true when SID is an include-list group
    IsGuarded    bool   `json:"isGuarded"`              //
    NearCastle   bool   `json:"nearCastle"`             // MainObject placement rule active
    RoadDistance string `json:"roadDistance"`           // "Any" | "Next To" | "Near" | "Medium" | "Far" | "Very Far"
    IsMine       bool   `json:"isMine"`                 // affects IsMine on the generated MandatoryContentItem
}
```

Default values when unmarshalled: a row with `Count == 0` and
`RoadDistance == ""` should be treated as `Count = 1`,
`RoadDistance = "Any"`. Add a normalising helper if needed.

### 2.2 Update: `go:internal/models/settingsFile.go`

Add the C# fields verbatim, **remove** the legacy
`PlayerZoneMandatoryContent []ZoneContentItem` field, and update the
defaults to match C#.

```go
type SettingsFile struct {
    // …existing scalar/topology fields…

    // ── Banned / overrides / bonuses ─────────────────────────────────
    BannedItems         string `json:"bannedItems"`
    BannedMagics        string `json:"bannedMagics"`
    ValueOverridesText  string `json:"valueOverrides"`
    BonusesJson         string `json:"bonuses"`         // pipe-separated, newline-separated; see §5.3

    // ── Mandatory content per zone type ──────────────────────────────
    PlayerZoneContentRows    []ZoneContentRowSave `json:"playerZoneContentRows,omitempty"`
    LowNeutralContentRows    []ZoneContentRowSave `json:"lowNeutralContentRows,omitempty"`
    MediumNeutralContentRows []ZoneContentRowSave `json:"mediumNeutralContentRows,omitempty"`
    HighNeutralContentRows   []ZoneContentRowSave `json:"highNeutralContentRows,omitempty"`
    HubZoneContentRows       []ZoneContentRowSave `json:"hubZoneContentRows,omitempty"`

    // Legacy v0.2 setting — when present seeds both density sliders.
    ContentDensityPercent *int `json:"contentDensity,omitempty"`
}
```

Defaults to change in `NewSettingsFile()`:

- `Topology: generator.TopologyBalanced` (was `TopologyRandom`).
- `HubZoneCastles: 0` (already 0; leave).
- `NeutralZoneCount: 0`, `NeutralZoneCastles: 1` — already correct.

### 2.3 Update: `go:internal/services/settingsFileLoader.go`

- Stop reading the old `playerZoneMandatoryContent` JSON key.
- When `ContentDensityPercent != nil && ResourceDensityPercent == nil`,
  seed `ResourceDensityPercent`/`StructureDensityPercent` from it (already
  done; verify).
- No other migration step is required: missing `*ContentRows` keys mean
  the user has no per-tier mandatory content, which is the new default.

### 2.4 Update: `go:internal/models/generator/generatorSettings.go`

Add list fields to mirror `cs:Models/Generator/GeneratorSettings.cs`:

```go
type GeneratorSettings struct {
    // …existing…
    BannedItems                  string
    BannedMagics                 string
    ValueOverridesText           string
    Bonuses                      []models.BonusEntry            // see §5.2

    PlayerZoneMandatoryContent   []template.MandatoryContentItem
    LowNeutralMandatoryContent   []template.MandatoryContentItem
    MediumNeutralMandatoryContent []template.MandatoryContentItem
    HighNeutralMandatoryContent  []template.MandatoryContentItem
    HubZoneMandatoryContent      []template.MandatoryContentItem
}
```

The translator from `SettingsFile` → `GeneratorSettings` (currently in
`go:internal/services/settingsFileLoader.go`, search for
`PlayerZoneMandatoryContent`) must:

1. Map each `ZoneContentRowSave` to a `template.MandatoryContentItem`
   (a fresh `internal/services/zoneContentRowMapper.go` is a good home
   for this — see §3.1).
2. Repeat the conversion for the four neutral/hub lists.
3. Parse `BonusesJson` into `[]models.BonusEntry` (see §5.2/§5.3).

### 2.5 Tests

Add `go:test/models/settingsFile_test.go` (new file). Verify:

- A `SettingsFile` with all new fields round-trips through
  `encoding/json` losslessly.
- Loading a fixture containing legacy `contentDensity` produces
  resource/structure densities equal to the legacy value.
- Loading a fixture with no `*ContentRows` keys yields empty slices.

---

## 3. Phase 2 — `MapTopology.Balanced`

### 3.1 Update: `go:internal/models/generator/mapTopology.go`

Append `TopologyBalanced`:

```go
const (
    TopologyDefault     MapTopology = "Default"
    TopologyHubAndSpoke MapTopology = "HubAndSpoke"
    TopologyChain       MapTopology = "Chain"
    TopologySharedWeb   MapTopology = "SharedWeb"
    TopologyRandom      MapTopology = "Random"
    TopologyBalanced    MapTopology = "Balanced"
)
```

### 3.2 Update: `go:internal/constants/ui.go`

Add `"Balanced"` to both `TopologyLabels` and `TopologyValues` (insert it
last so existing dropdown indices keep working).

C# label map (`TopologyDisplayName` in `TemplateGenerator.cs`):

| Enum | Label |
| --- | --- |
| Default | "Ring" |
| HubAndSpoke | "Hub and Spoke" |
| Chain | "Chain" |
| SharedWeb | "Shared Web" |
| Random | "Random" |
| Balanced | "Balanced" |

### 3.3 Update: `go:internal/services/templateGenerator.go`

In `describeTopology(...)`: add `TopologyBalanced => "Balanced"`.
In the variant builder switch: add `case TopologyBalanced: return
buildVariantBalanced(...)`. Until phase 6a lands, alias
`buildVariantBalanced` to `buildVariantRandom` so the program still
builds.

In the existing description string builder, **drop** the
"balanced zone placement" option from the option list (it's no longer an
experimental flag).

Mark `SettingsFile.ExperimentalBalancedZonePlacement` deprecated. When
loading an old file that has it set, upgrade it to `Topology = Balanced`
and clear the flag — this is the only piece of legacy migration we keep.

### 3.4 Tests

Update `go:test/services/services_test.go` to cover:

- A settings with `Topology = TopologyBalanced` generates without panic.
- An older settings file with `experimentalBalancedZonePlacement: true`
  is upgraded to `Topology = Balanced` after loading.

---

## 4. Phase 3 — `ZoneContentManager` refactor

C# commit `5f7cca5` plus `ccf97c4` reduce each `Build…MandatoryContent`
method to: optional foothold + verbatim copy of the user's configured
list. **All the hardcoded watchtower/market/include-list rows currently
baked into Go must be removed.**

### 4.1 Replace the body of each builder in `go:internal/services/zoneContentManager.go`

```go
func BuildPlayerZoneMandatoryContent(s *models.GeneratorSettings) []template.MandatoryContentItem {
    var out []template.MandatoryContentItem
    if s.SpawnRemoteFootholds {
        out = append(out, presetRemoteFoothold(s.ZoneCfg.PlayerZoneCastles))
    }
    return append(out, s.PlayerZoneMandatoryContent...)
}

func BuildLowNeutralMandatoryContent(s *models.GeneratorSettings) []template.MandatoryContentItem {
    var out []template.MandatoryContentItem
    if s.SpawnRemoteFootholds {
        out = append(out, presetRemoteFoothold(s.ZoneCfg.PlayerZoneCastles))
    }
    return append(out, s.LowNeutralMandatoryContent...)
}
// MediumNeutral and HighNeutral follow the same shape.
```

Add a new builder for hub zones:

```go
func BuildHubZoneMandatoryContent(s *models.GeneratorSettings) []template.MandatoryContentItem {
    var out []template.MandatoryContentItem
    if s.SpawnRemoteFootholds {
        out = append(out, presetRemoteFoothold(s.ZoneCfg.HubZoneCastles))
    }
    return append(out, s.HubZoneMandatoryContent...)
}
```

### 4.2 Port `StripNearCastleRules`

```go
func StripNearCastleRules(items []template.MandatoryContentItem) []template.MandatoryContentItem {
    out := make([]template.MandatoryContentItem, 0, len(items))
    for _, item := range items {
        hasNearCastle := false
        for _, r := range item.Rules {
            if r.Type == "MainObject" { hasNearCastle = true; break }
        }
        if !hasNearCastle { out = append(out, item); continue }
        stripped := item
        stripped.Rules = nil
        for _, r := range item.Rules {
            if r.Type != "MainObject" { stripped.Rules = append(stripped.Rules, r) }
        }
        out = append(out, stripped)
    }
    return out
}
```

Wire it in the algorithm wherever a zone has no castle (look for
`PlayerZoneCastles == 0` / `NeutralZoneCastles == 0` branches in
`templateGenerator.go`).

### 4.3 Row → MandatoryContentItem mapper

Add `go:internal/services/zoneContentRowMapper.go` (or a small section
in `zoneContentManager.go`):

```go
// rowsToMandatoryItems converts UI rows to engine content items.
func rowsToMandatoryItems(rows []models.ZoneContentRowSave) []template.MandatoryContentItem {
    out := make([]template.MandatoryContentItem, 0, len(rows))
    for _, r := range rows {
        item := template.MandatoryContentItem{IsGuarded: r.IsGuarded, IsMine: r.IsMine}
        if r.IsGroup {
            item.IncludeLists = []string{r.SID}
        } else {
            item.SID = r.SID
        }
        if r.NearCastle {
            item.Rules = append(item.Rules, ruleNearCastle(1))
        }
        if dp, ok := lookupDistancePreset(r.RoadDistance); ok {
            item.Rules = append(item.Rules, ruleRoadDistance(dp, 1))
        }
        // Count > 1 expands into Count copies.
        for i := 0; i < max(r.Count, 1); i++ { out = append(out, item) }
    }
    return out
}
```

Distance preset lookup table to add (currently only `distNextTo` is
declared in Go; C# defines all six):

| Label | Min | Max |
| --- | --- | --- |
| `"Any"` | — (no rule) |
| `"Next To"` | 0.05 | 0.10 |
| `"Near"` | 0.10 | 0.25 |
| `"Medium"` | 0.25 | 0.50 |
| `"Far"` | 0.50 | 0.75 |
| `"Very Far"` | 0.75 | 0.90 |

### 4.4 Content limits

`BuildAllContentCountLimits` should now lift caps from **all five**
mandatory-content lists, not just `PlayerZoneMandatoryContent`. Iterate
all five, accumulate SID counts, then bump default `MaxCount`.

### 4.5 Update template generator

Anywhere the algorithm currently calls `BuildLow/Medium/High…` with
parameters `(castleCount, spawnFootholds)`, switch to passing
`*GeneratorSettings` so the new per-tier user content is included.

### 4.6 Tests

Move the asserts that count *specific* hardcoded items in
`go:test/services/services_test.go` to instead assert the
shape returned with explicit user-configured rows. C# tests do the same.

---

## 5. Phase 4/5 — UI (Zone Content tab + bonuses)

### 5.1 Zone Content tab — new sub-tabs

`cs:MainWindow.xaml` after `ccf97c4` exposes five mandatory-content
editors: **Player**, **Low**, **Medium**, **High**, **Hub** (the latter
is only visible when `HubZoneSize > 0`). The Go equivalent is
`go:internal/gui/components/zoneContentPanel.go`.

Implementation outline:

- Wrap the existing single editor in a five-way tab strip rendered with
  `content/segmentButtonGroup.go`.
- Each tab edits a different slice on `state.SettingsFile`
  (`PlayerZoneContentRows`, `LowNeutralContentRows`, …).
- Reuse `content/zoneContent.go` for the row renderer; bind the active
  slice via a closure.
- Hide the Hub tab when `state.SettingsFile.HubZoneCastles == 0 &&
  state.SettingsFile.HubZoneSize == 0`.

### 5.2 Bonus model

New file `go:internal/models/bonus.go`:

```go
package models

type BonusPresetType int

const (
    BonusTownPortalFree BonusPresetType = 0
    BonusSpell          BonusPresetType = 1
    BonusUnitMultiplier BonusPresetType = 2
    BonusMovement       BonusPresetType = 3
    BonusStartingItem   BonusPresetType = 4
    BonusStartingGold   BonusPresetType = 5
    BonusStartingGems   BonusPresetType = 6
    BonusStartingCrystals BonusPresetType = 7
    BonusStartingMercury  BonusPresetType = 8
    BonusStartingWood     BonusPresetType = 9   // added in 4565e6c
    BonusStartingOre      BonusPresetType = 10  // added in 4565e6c
)

type BonusEntry struct {
    PresetType     BonusPresetType
    ReceiverFilter string // "start_hero" | "all_heroes"
    Param          string
    Param2         string // for Spell: "1" = free, "0" = normal
}
```

### 5.3 `BonusesJson` format

Pipe-separated, newline-separated. Each non-empty line is one
`BonusEntry`:

```
<PresetType>|<ReceiverFilter>|<Param>|<Param2>
```

`<PresetType>` accepts both the enum name (preferred) and the integer
ordinal (legacy). See `BonusEntry.FromString` in
`cs:Models/Unfrozen/BonusEntry.cs`. The C# editor stores the full list
inside a single JSON string field rather than as a JSON array. Mirror
that for round-trip compatibility.

### 5.4 `BonusEntry → []template.Bonus` expansion

Each `BonusEntry` expands to one or two raw `Bonus` rows. Port
`ToBonuses()` exactly:

| Preset | Sid | Parameters |
| --- | --- | --- |
| TownPortalFree | `add_bonus_hero_spell` then `add_bonus_hero_stat` | `["neutral_magic_town_portal"]` / `["magicCostSidSet", "neutral_magic_town_portal", "-999", "0"]` |
| Spell | `add_bonus_hero_spell` (+ `add_bonus_hero_stat` if free) | `[Param]` / `["magicCostSidSet", Param, "-999", "0"]` |
| UnitMultiplier | `add_bonus_hero_unit_multipler` *(note typo, must be preserved)* | `[Param]` |
| MovementBonus | `add_bonus_hero_stat` | `["movementBonus", Param]` |
| StartingItem | `add_bonus_hero_item` | `[Param]` |
| StartingGold | `add_bonus_res` | `["gold", Param]` |
| StartingGems | `add_bonus_res` | `["gemstones", Param]` |
| StartingCrystals | `add_bonus_res` | `["crystals", Param]` |
| StartingMercury | `add_bonus_res` | `["mercury", Param]` |
| StartingWood | `add_bonus_res` | `["wood", Param]` |
| StartingOre | `add_bonus_res` | `["ore", Param]` |

All entries use `ReceiverSide = -1` and the row's `ReceiverFilter`.

Wire the conversion into the existing
`go:internal/services/templateGenerator.go` build phase that emits
`template.Bonuses`. Existing Go bonus injection should be removed in
favor of this list.

### 5.5 Bonus picker UI

New file `go:internal/gui/components/bonusPickerPanel.go` (or modal
window via Gio). Visible fields, mirroring `cs:BonusPickerWindow.xaml`:

- Type dropdown (all 11 `BonusPresetType` values).
- Receiver dropdown (`start_hero`, `all_heroes`) — **hidden** when
  type is a starting-resource preset.
- Per-type input panel that swaps based on selected type:
  - Spell: text + "Pick spell…" button + "Make free" checkbox.
  - UnitMultiplier: numeric, default `"2"`.
  - MovementBonus: numeric, default `"0"`.
  - StartingItem: text + "Pick item…" button.
  - Resource presets: amount field with type-appropriate label/default
    (`Gold 10000`, `Gems 15`, `Crystals 15`, `Mercury 15`, `Wood 20`,
    `Ore 20`).
- An "Add" button that appends to the bonus list.

Spell / item pickers (`SpellPickerWindow`, `ItemPickerWindow`) are full
WPF child windows in C#. The Go port can ship MVP-level
combobox+freeform-text equivalents in this phase and upgrade later.

### 5.6 Tests

- Bonus round-trip: `BonusEntry.String()`/`Parse()` are inverses for
  every preset type, including legacy numeric form.
- `Generate(...)` injects the right number of `template.Bonuses` for a
  list containing one of each preset.

---

## 6. Phase 6 — `TemplateGenerator` algorithm port

This phase is the largest; the C# generator grew by ~840 lines since the
last Go sync. **Do not attempt all sub-phases at once.** Each subphase
should land independently with green tests.

### 6a. Balanced layout (`6dd070e`, `0e31576`, `2d4cace`, `e268446`, `02bb65d`, `047370a`, `d2166eb`)

The key extraction in `6dd070e` is splitting `BuildVariantRandom` into:

- `BuildVariantRandom` — Delaunay over uniformly-random positions.
- `BuildVariantBalanced` — Delaunay over **concentric-ring** positions
  (players on outer ring, low-quality neutrals next, …, high-quality
  innermost).
- `BuildVariantFromDelaunay(filterByTier bool, …)` — shared body. When
  `filterByTier` is true, prune edges that skip a quality tier (allowed:
  same-tier, player↔low, low↔medium, medium↔high; tier-skip allowed only
  when an intermediate tier is absent).

Helpers to port (study `TemplateGenerator.cs` around lines 1100–1700 in
the C# repo, post-merge):

| C# method | Go equivalent (rename to lowerCamelCase) |
| --- | --- |
| `BuildBalancedRingLetters(...)` | already exists; verify tier ordering matches |
| `BuildBalancedRandomPositions(...)` | port: assigns `(x,y)` on rings whose radius scales with tier |
| `BuildBalancedNeutralRing(...)` | port: distributes neutrals around a single ring while balancing distance between players |
| `DelaunayEdges(...)` | already exists; verify identical behaviour for collinear input |
| `ZoneQualityGroup(letter, players, neutralByLetter)` | already exists |
| `EnsureFullConnectivity(...)` | already exists; verify the new `tierSkipPenalty` weighting |
| `EnsurePlayerZonesConnected(...)` | already exists |

Specific fixes layered on top of the initial port:

- `02bb65d` + `e268446` + `d2166eb`: inner rings (high-quality and hub
  zones) must always form a connected sub-graph. Add a fallback pass
  that adds intra-ring edges between adjacent angular positions when the
  Delaunay edge set leaves them disconnected.
- `047370a`: in `EnsureFullConnectivity`, prefer adding an edge that
  bridges *all* of: a disconnected sub-component, the nearest already-
  connected zone, and the smallest tier skip. Use `(tierGap*1000 +
  distance)` as the edge cost.
- `2d4cace`: scale per-tier ring radii by player count so that 2-player
  maps don't end up with overlapping rings. Formula from C#:
  `radius(tier) = 0.45 - 0.10 * tier` for player count ≥ 4, and shrink
  outer ring linearly for smaller player counts.

### 6b. Tournament balanced cluster (`6d28ef0`, `a040c98`, `df8a115`, `cd3a46e`, `ab4c1d3`, `6606755`, `e91e79f`, `c20b40d`)

Add three cluster builders mirrored from C#:

| C# | Go target |
| --- | --- |
| `BuildTournamentRingCluster` | `buildTournamentRingCluster` |
| `BuildTournamentHubCluster` | `buildTournamentHubCluster` (currently chain — rename + adjust) |
| `BuildTournamentBalancedCluster` | new `buildTournamentBalancedCluster` |

`buildVariantTournament` switches on the cluster topology (`Default`,
`HubAndSpoke`, `Chain`, `Balanced`) per cluster — not per template.
Guard match-group naming must follow:

- Ring cluster: `tourney_ring_guard_{cluster}`
- Hub cluster: `tourney_hub_guard_{cluster}`
- Chain cluster: `tourney_chain_guard_{cluster}`
- Balanced cluster: `tourney_bal_guard_{cluster}`

`c20b40d` permits `RandomPortals` in tournament mode (previously
suppressed); just drop the early `if (Tournament) return;` guard.

`e91e79f` changes neutral-zone sorting inside a chain/ring cluster from
"input order" to "highest quality first, then by castle count". Port the
new `OrderTournamentClusterNeutrals` comparator.

### 6c. Random / shared-web fixes (`3f46209`, `cc5d7ee`, `c63ea67`)

- `3f46209`: when previewing `Random` topology with a tiny zone count
  (≤ 3 players, no neutrals), short-circuit to a hand-laid layout so the
  preview doesn't degenerate.
- `cc5d7ee`: in `BuildVariantChain` ensure player zones occupy chain
  endpoints (positions 0 and N-1) when `PlayerCount == 2`; otherwise the
  spawn zones may sit interior and break the "no direct connection"
  guarantee.

### 6d. Guard scaling (`89cba4b`)

Replace the current single `borderGuardMultiplier` value with a tiered
function:

```go
func scaleNeutralGuardValue(base int, quality constants.NeutralZoneQuality) int {
    switch quality {
    case constants.QualityLow:    return base
    case constants.QualityMedium: return int(float64(base) * 1.5)
    case constants.QualityHigh:   return int(float64(base) * 2.5)
    }
    return base
}
```

and equivalent `scaleBorderGuardValue` based on the zone-tier pair
(player↔neutral uses the neutral's tier; neutral↔neutral uses the
**higher** of the two). Wire it through every place currently calling
the old multiplier.

`de73435` only reformats error strings; trivial port.

`2c62fab` adds defensive panic guards around division-by-zero paths in
the tournament balanced cluster. Mirror as Go `if … { return nil, err }`
returns rather than panics.

### 6e. Tests

`TemplateGeneratorTests.cs` gained the following relevant test methods
between v0.6 and v0.7 (port each into
`go:test/services/services_test.go`):

- `Generate_BalancedTopologyProducesRingedLayout`
- `Generate_BalancedTopologyKeepsHighQualityNeutralsCentral`
- `Generate_TournamentBalancedClusterIsolatesPlayers`
- `Generate_TournamentRingClusterEmitsCorrectGuardGroupNames`
- `Generate_TournamentAllowsRandomPortals`
- `Generate_RandomTopologyShortCircuitsForTinyMaps`
- `Generate_GuardValuesScaleWithNeutralQuality`
- `Generate_StripsNearCastleRulesOnCastlelessZones`

(Names paraphrased; the exact names are in `cs:Olden Era - Template
Editor.Tests/TemplateGeneratorTests.cs`. Grep there for `[Fact]`.)

---

## 7. Phase 7 — Preview rendering

**Status: ✅ Complete.** `internal/services/previewLayout.go` now has a
dedicated Balanced branch (concentric tier rings, per-cluster canvas
split via floodfill so tournament-balanced renders as two side-by-side
sub-canvases). Hub-castle count was already drawn whenever
`HasCastle && Castles > 0` (set for both `Spawn` and `City` main
objects), so 7.1 needs no extra renderer work. A panic-free smoke test
`TestRenderPreviewImage_DoesNotPanic_AllTopologies` covers every
topology × `PlayerCount ∈ {2,3,4,6,8}`, both vanilla and tournament.

### 7.1 Hub castle count badge (`432eeb9`)

In `go:internal/services/previewRenderer.go` the hub zone is drawn as a
single chip. Add a small castle-count overlay (a tower glyph followed by
`HubZoneCastles`) similar to how player/neutral chips already show their
castle count.

`cs:Services/TemplatePreviewPngWriter.cs` does this with a tiny pill
roughly 20×14 in the top-right corner of the hub chip; use the same
position in the Go renderer for visual parity.

### 7.2 Balanced layout preview

`previewLayout.go` currently chooses positions only for the existing
topologies. Add a `Balanced` branch that mirrors the algorithm:

1. Concentric rings: player ring at radius `R_player`, then low/med/high
   neutral rings, hub at centre if present.
2. Within each ring, place letters at equal angular intervals starting
   at a phase offset that aligns players with neutral gaps.
3. Compute connection polylines from the same Delaunay edge set returned
   by the generator so preview matches generated template.

### 7.3 Tournament balanced preview

Group player+neutrals into clusters as in `buildVariantTournament`,
position each cluster on its own concentric area, and reuse 7.2 within
each cluster.

### 7.4 Tests

Preview is rendered to PNG. Add a smoke test that the renderer doesn't
panic for every topology under `PlayerCount ∈ {2,3,4,6,8}`.

---

## 8. Phase 8 — Test parity & misc

**Status: ✅ Complete.** All actionable items ported:

- `30a235f` — "The Gorge" SID display name → "Carrion Pile"
  (`internal/constants/contentIds.go`).
- `48deacc` / `7029355` / `9151ffd` — full content-group catch-up:
  - Added ~50 missing SIDs to `internal/constants/contentIds.go`
    (Altar/Magic Amplifier 1–4, Scroll Box variants, Twilight Bloom,
    Storage * resources, all new resource-bank / utility / hero-improvement
    structures).
  - Added 17 new entries to `internal/constants/includeListIds.go`
    (Guarded Banks T1–T3, Basic Storage, Rare Mines, Magic Buildings T1–T2,
    Hero Stats/Skills T1–T3, etc.) plus the C# naming fixes (`Any Tier`,
    weighted variant).
  - Restructured `internal/constants/contentItemGroups.go` to mirror
    the v0.7 six-group layout (`Mines`, `UtilityStructures`, `Treasures`,
    `UnitRecruitment`, `ResourceBanks`, `HeroImprovementStructures`).
    `HireBuildings` is kept as an alias of `UnitRecruitment` for
    backwards compatibility with existing UI callers.
- `fc4cb46` — neutral-zone breakdown no longer gated by Advanced
  toggle: `buildNeutralZonePlan` in `internal/services/templateGenerator.go`
  now always reads the breakdown fields when any of them are non-zero,
  falling back to the legacy single-bucket form when the breakdown is
  empty (preserves back-compat for older settings files).
- `7303848` / `cc63ecc` — pure C# test renames/refactors; Go test suite
  already uses the new `Advanced.NeutralMediumCastleCount` shape and
  the v0.7 guard values are exercised by the Phase 6d tier-aware guard
  scaling tests.

### 8.1 Original notes (kept for reference)

- `7303848` / `cc63ecc` / `c400677` are pure test fixes; port the test
  changes (renamed asserts, removed obsolete tests) directly.
- `30a235f` renames "The Gorge" UI label; mirror in
  `go:internal/constants/contentIds.go` or wherever the display name
  comes from.
- `9151ffd` / `48deacc` / `7029355` add new content groups. Inventory
  diff:

```
git show 9151ffd -- 'Olden Era - Template Editor/Services/ContentManagement/ContentItemGroup.cs'
git show 48deacc -- 'Olden Era - Template Editor/Services/ContentManagement/ContentItemGroup.cs'
git show 7029355 -- 'Olden Era - Template Editor/Services/ContentManagement/ContentItemGroup.cs'
```

then port any missing groups into
`go:internal/constants/contentItemGroups.go`.

- `ca7556b` UI refactor and `fb3c203` bonus-picker polish are largely
  cosmetic; cross-reference once phases 4/5 are in.
- `fc4cb46` moves the neutral-zone breakdown UI out of "advanced
  settings". In Go that means moving the Low/Medium/High count fields
  in `basicSetupPanel.go` (or `generationPanel.go`) out of the
  "Advanced Mode" gate so they're always visible.

---

## 9. Suggested PR shape

| PR | Phases | Notes |
| --- | --- | --- |
| #1 | 1 | persistence-only, no behavior change |
| #2 | 2 | balanced enum + UI label, dispatch falls through to random |
| #3 | 3 | content manager refactor; updates existing tests |
| #4 | 4 | Zone Content tab gets 5 sub-tabs |
| #5 | 5 | bonus model + picker UI + generator wiring |
| #6a–6d | 6 | algorithm port — split into four PRs for review sanity |
| #7 | 7 | preview improvements |
| #8 | 8 | test parity cleanup + content-group catch-up |

---

## 10. Useful commands

```powershell
# See files changed by any given commit
cd 'D:\Git\some_personal_shit\Olden-Era---Template-Generator'
git show --stat <hash>
git show <hash> -- 'Olden Era - Template Editor/Services/TemplateGenerator.cs'

# Range diff for one file across the whole window
git log --since=2026-05-12 --reverse --format='%h %s' -- '<file>'
git diff 6d28ef0~1 6c97732 -- '<file>'

# Side-by-side method extraction
git show 6dd070e -- 'Olden Era - Template Editor/Services/TemplateGenerator.cs' | code -
```

Keep this document up to date as each phase lands; tick off the
corresponding entries in §0 / §9.
