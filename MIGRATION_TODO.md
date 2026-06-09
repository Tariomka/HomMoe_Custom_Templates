# Migration TODO — Parity with `Olden Era - Template Editor` (C# WPF)

This document tracks the work required to bring the Go port
(`HomMoe_Custom_Templates`) up to functional parity with the upstream C# WPF
project at `..\Olden-Era---Template-Generator\Olden Era - Template Editor\`.

The upstream project recently received a batch of features that have **not** yet
been ported. Source commits of interest (newest first):

| Commit | Feature |
|--------|---------|
| `9dad383`, `ecd8f88`, `d9e3891`, `efcbc93`, `34a16b2` | Visual node/zone-connection editor (curved parallel edges, drag-to-create, **Proximity type removed**, per-tier guard strength presets incl. player-to-player) |
| `8070062` → `9e5141d` | Polymorphic zone content rule system (`IContentRule`), Add/Manage rule dialogs, rule markers, **Solo Encounter** rule |
| `2ab7099`, `f9feaed`, `f5b8700` | Item **variants** (`RuleVariant` + `VariantMapping`), expanded Pandora Box options |
| `05b6f91` | Multiple starting castles |
| `42533b7` | Single hero mode |
| `ccf97c4`, `fc4cb46` | Neutral mandatory content |
| `4565e6c` | Wood & ore start bonuses |

> Legend in the checklists below:
> ✅ already present in Go · ⚠️ partially present · ❌ missing

Work is split into the three requested sections: **Model**, **Service/Backend**,
and **UI**. Implement them in that order — UI work depends on the model and
service layers being in place first.

---

## 1. Model Portion

Data-struct changes required so the persisted/serialized data matches upstream.

### 1.1 Polymorphic content-rule serialization ❌ (highest priority)

Upstream replaced the flat per-row content fields with a polymorphic rule list.

- **C# reference:** `Models/Generator/ZoneContentRowSave.cs`
  (`ZoneContentRowSave`, `ContentRuleRowSave`), `Models/Generator/ZoneContentItemUI.cs`.
- **Go target:** [internal/models/zoneContentRowSave.go](internal/models/zoneContentRowSave.go)

Current Go `ZoneContentRowSave` is the *legacy* flat shape:
`Sid, Count, IsGroup, IsGuarded, NearCastle, RoadDistance, IsMine`.

- [ ] Add a `ContentRuleRowSave` struct (new file
  `internal/models/contentRuleRowSave.go`) mirroring the C# fields:
  `Name string`, `DistanceName string`, `IsGuarded *bool`,
  `IsSoloEncounter *bool`, `VariantId *int` (all `omitempty`).
- [ ] Add `Rules []ContentRuleRowSave \`json:"rules,omitempty"\`` to
  `ZoneContentRowSave`.
- [ ] Keep `IsGuarded`, `NearCastle`, `RoadDistance` as **deprecated**
  fields for backward-compatible load of existing `.gen.json` files, but stop
  writing them. Match the C# comment markers (`deprecated`).
- [ ] Update `Normalised()` to migrate legacy flat fields into the new
  `Rules` list on load (so old saves keep working). This is the load-time
  upgrade path equivalent to C#'s `CreateRuleFromSavedRule`.

### 1.2 EditorState row collections ⚠️

- **Go target:** [internal/models/editorStateModel.go](internal/models/editorStateModel.go)
- The five row collections (`PlayerZoneContentRows`, `LowNeutralContentRows`,
  `MediumNeutralContentRows`, `HighNeutralContentRows`, `HubZoneContentRows`)
  already exist and will carry the new `Rules[]` automatically once 1.1 lands.
- [ ] Verify round-trip (save → load) keeps `Rules` intact; add/adjust tests in
  `test/models`.

### 1.3 `.rmg.json` content item — already mostly done ✅ / ⚠️

- **Go target:** [internal/models/template/template_inner/content/mandatoryContentItem.go](internal/models/template/template_inner/content/mandatoryContentItem.go)
- `MandatoryContentItem` already exposes `SoloEncounter bool` and
  `Variant *int`, matching C# `ContentItem.SoloEncounter` / `Variant`. ✅
- [ ] Confirm `Variant` (pointer int) is emitted with the correct
  `variant` JSON key and omitted when nil (matches C# `Variant?`). ⚠️ verify
  against example templates.

### 1.4 Connection model — remove `Proximity`, add user-edit tracking ⚠️

- **C# reference:** `Models/Unfrozen/Connection.cs` (note new `[JsonIgnore]
  IsUserAdded`), `Models/Unfrozen/Miscellaneous.cs`.
- **Go target:**
  [internal/models/template/template_inner/variant/connection.go](internal/models/template/template_inner/variant/connection.go)
- [ ] Add an unserialized `IsUserAdded bool` (json:"-") field to support the
  visual editor's "user added this edge" tracking (needed by §3.1).
- [ ] The `Length` (Proximity) field and `Proximity` connection type are being
  retired upstream — see §2.5 and §3.1. Keep the struct field for reading
  legacy templates, but stop generating it.

### 1.5 Variant-mapping data ❌

Needed for item variants (Pandora Box, Dragon Utopia, Monty Hall, etc.).

- **C# reference:** `Services/ContentManagement/VariantMapping.cs`
  (`VariantMapping`, `VariantMappingManager`).
- **Go target:** new file, e.g.
  `internal/models/variantMapping.go` (or under `internal/registry/`).
- [ ] Add a `VariantMapping` struct: `{ Sid string; Variants map[int]string }`
  (variant id → human description), keyed by content SID.
- [ ] Port the variant tables: `utopiaVariants` (4), `pandoraBoxVariants`
  (~28: Gold T1–T4, Exp T1–T4, Units T1–T7, All Stats T1–T4, Magic Schools ×4,
  Spells T1–T5), `montyHallVariants` (4).

---

## 2. Service / Backend Portion

Generation and content-management logic changes.

### 2.1 Polymorphic content-rule engine ❌ (highest priority)

Upstream replaced the inline "if NearCastle / if RoadDistance" logic with a
polymorphic `IContentRule` system.

- **C# reference:** `Services/ContentManagement/Rules/`
  (`IContentRule`, `ContentRuleManager`, `RuleSoloEncounter`, `RuleGuarded`,
  `RuleVariant`, `RuleDistanceToRoad`, `RuleDistanceToTown`).
- **Go target:** new package, e.g.
  `internal/services/content_rules/` plus existing
  [internal/services/zoneContentManager.go](internal/services/zoneContentManager.go).

- [ ] Define a `ContentRule` interface mirroring `IContentRule`:
  `Name()`, `Description()`, `Marker()` (1–2 char UI badge),
  `DisplayText()`, `SerializeToRowSave() ContentRuleRowSave`.
- [ ] Implement the five rule types:
  - `RuleSoloEncounter` → sets `MandatoryContentItem.SoloEncounter`; marker `S` / `!S`.
  - `RuleGuarded` → sets `IsGuarded`; marker `G` / `!G`.
  - `RuleVariant` → sets `Variant`; needs the SID's `VariantMapping`; marker `""`.
  - `RuleDistanceToRoad` → emits `PlacementRule{Type:"Road", TargetMin/Max}`; marker `R`.
  - `RuleDistanceToTown` → emits `PlacementRule{Type:"MainObject", Args:["0"], TargetMin/Max}`; marker `T`.
- [ ] Add a `ContentRuleManager` equivalent:
  - `GetRules()` — list of available rule prototypes.
  - `ApplyRulesToFinalContentItem(item, rules)` — convert UI rules → final
    `MandatoryContentItem` fields/`Rules[]` (replaces today's flat conversion).
  - `CreateRuleFromSavedRule(ContentRuleRowSave, sid)` — reconstruct a rule
    from a saved row (used by the model migration in §1.1).
- [ ] Enforce "at most one rule per type" semantics (matches the C# dialogs).

### 2.2 Rewire `zoneContentManager` / `MandatoryContentProvider` ⚠️

- **Go target:**
  [internal/services/zoneContentManager.go](internal/services/zoneContentManager.go),
  `internal/services/template_generator/providers/mandatoryContentProvider.go`.
- [ ] Replace `rowToMandatoryItem` / `createContentItemFrom` flat logic
  (`row.NearCastle`, `row.RoadDistance`, `row.IsGuarded`) with calls to the new
  rule engine (`ApplyRulesToFinalContentItem`).
- [ ] Keep `StripNearCastleRules` behaviour (drop `MainObject[0]` distance-to-town
  rules when the zone has 0 castles) — now expressed in terms of the new rules.
- [ ] Solo Encounter: ensure `RuleSoloEncounter` flows through to
  `MandatoryContentItem.SoloEncounter` in generated `.rmg.json`.

### 2.3 Variant support in generation ❌

- **C# reference:** `RuleVariant` + `ContentItemBuilder`/`ContentItem.Variant`.
- **Go target:** generation providers + `zoneContentManager`.
- [ ] When a row has a variant rule, set `MandatoryContentItem.Variant`.
- [ ] Confirm content count limits and pools still work with variants
  (`internal/services/template_generator/providers/contentLimitProvider.go`).

### 2.4 Single hero mode & multiple starting castles ⚠️ (verify)

These appear largely ported already; confirm parity.

- **C# reference:** `Services/ContentManagement/.../GameRulesProvider`,
  `GeneratorSettings.SingleHeroMode`, `ZoneConfiguration.*Castles`.
- **Go target:** `internal/services/template_generator/providers/gameRulesProvider.go`,
  [internal/models/config/config_inner/zoneConfig.go](internal/models/config/config_inner/zoneConfig.go).
- [ ] Verify SingleHero mode clamps hero count to 1 and sets `HeroHireBan=true`.
- [ ] Verify `PlayerZoneCastles`, `NeutralZoneCastles`, `HubZoneCastles` are all
  honoured per tier and that >1 starting castle generates the extra
  `MainObject` city entries (upstream `05b6f91`).

### 2.5 Connection generation — retire `Proximity` ⚠️

- **Go target:** [internal/services/previewLayout.go](internal/services/previewLayout.go) (line ~214
  still treats `"Proximity"` as a portal-like edge), topology providers under
  `internal/services/template_generator/providers/topology/`,
  [internal/registry/connectionTypeValues.go](internal/registry/connectionTypeValues.go).
- [ ] Stop emitting `Proximity` connections / `Length` from the generator.
- [ ] Remove `Proximity` from the user-selectable connection types
  (keep parsing it for legacy template reads only).

### 2.6 Guard-strength presets (per-tier, incl. player-to-player) ❌

Backend data tables consumed by the visual editor (§3.1).

- **C# reference:** `ZoneConnectionEditorWindow.xaml.cs` strength tables.
- **Go target:** new constants, e.g. `internal/constants/guardStrengths.go`.
- [ ] Define per-tier guard value presets with 5 strength levels
  (Weak/Moderate/Medium/High/Very High) for Bronze (~3k–16k), Silver
  (~18k–30k), Gold (~36k–60k) and **PlayerToPlayer** (~10k–58k).
- [ ] Define a "Generator Default" value per tier (commit `34a16b2`/`efcbc93`).
- [ ] Define weekly-increment presets (Slow 5% … Very Fast 25%).

---

## 3. UI Portion

Gio (Go) GUI changes so the interface is functionally equivalent. Existing GUI:
[internal/gui/window.go](internal/gui/window.go) + tab panels under
`internal/gui/components/`.

### 3.1 Visual zone-connection (node) editor ❌ (largest item)

- **C# reference:** `ZoneConnectionEditorWindow.xaml(.cs)`.
- **Go target:** new component, e.g.
  `internal/gui/components/connectionEditor/` + a launch button on the Layout
  tab and/or preview panel.
- [ ] Canvas rendering of zones as draggable nodes, colour-coded by tier
  (Player=green, Bronze/Silver/Gold, Hub=blue) — reuse colours from
  `internal/gui/components/themes/previewTheme.go` / `internal/constants/legend.go`.
- [ ] Render connections as **curved** lines, with curved offsets for
  **parallel edges** (commit `ecd8f88`/`9dad383`); Direct=gold, Portal=blue.
- [ ] **Drag-to-create**: click node A → drag to node B → create connection,
  with rubber-band preview.
- [ ] Edge selection → property panel: Connection Type (Direct/Portal — **no
  Proximity**), Guard Value with the per-tier strength preset buttons and
  "Generator Default" (§2.6), Weekly Increment dropdown, Guard Escape and
  Sim Turn Squad toggles.
- [ ] Right-click (or delete button) removes a connection.
- [ ] Track edits via `Connection.IsUserAdded` (§1.4); surface
  `ConnectionsWereModified` and `HasUnresolvedErrors` (block export when a
  connection references a non-existent zone).
- [ ] Persist user-edited connections across regeneration (decide storage:
  extend `EditorStateModel` with a saved-connections list).
- [ ] Warning when recreating a template after manual edits (commit `8cdf375`).

### 3.2 Zone-content rule dialogs ❌

- **C# reference:** `Windows/AddZoneContentRuleWindow.xaml(.cs)`,
  `Windows/ManageZoneContentRulesWindow.xaml(.cs)`.
- **Go target:** new dialogs under `internal/gui/components/` + hook into
  the Zone Content tab rows
  (`internal/gui/components/content/zoneContent.go`).
- [ ] Replace the per-row flat controls (Guarded checkbox, Near-Castle checkbox,
  RoadDistance dropdown) with a single **"Manage Rules"** button per row plus a
  compact **rule-marker** display (`S`, `G`, `R`, `T`, `!G` …) — commit `f37b78a`.
- [ ] **Manage Rules** dialog: list applied rules, Add / Edit (double-click) /
  Remove, enforcing one rule per type.
- [ ] **Add/Edit Rule** dialog: polymorphic editor switching on rule type:
  - Distance to Road / Distance to Town → distance dropdown
    (`Any, Next To, Near, Medium, Far, Very Far`).
  - Guarded → checkbox.
  - Solo Encounter → checkbox.
  - Variant → dropdown of `VariantMapping` options for the row's SID (§1.5).
- [ ] Per-row rule descriptions/tooltips for the user (commit `df1e378`).

### 3.3 Item / Spell / Value-override / Bonus pickers ❌

Currently all of these are text-area driven in Go (Bonuses & Bans tab,
`internal/gui/components/.../bonusesPanel.go`).

- **C# reference:** `Windows/ItemPickerWindow`, `SpellPickerWindow`,
  `ValueOverridePickerWindow`, `BonusPickerWindow`.
- **Go target:** new picker dialogs.
- [ ] **Item picker** — searchable, category-grouped, multi-select; returns SIDs.
- [ ] **Spell picker** — searchable, grouped by magic school, tier-sorted,
  multi-select, "free spell" toggle.
- [ ] **Value-override picker** — searchable SID list, multi-select; emits
  `sid=guardValue` lines into the existing override text.
- [ ] **Bonus picker** — section-based UI (spell / unit multiplier / movement /
  starting item / resource amounts incl. **Wood & Ore**), receiver filter
  dropdown (`all_heroes` / `start_hero`). Wood/Ore presets already exist in the
  Go bonus model (`BonusPresetType.StartingWood/StartingOre`) — this only adds
  the picker UI (commit `4565e6c`).

### 3.4 Neutral mandatory content UI ⚠️ (verify)

- The Go Zone Content tab already has the five tiers
  (Player, Low/Medium/High Neutral, Hub) via the tier `SegmentButtonGroup`. ✅
- [ ] Confirm each neutral tier persists its own mandatory-content rows and that
  per-tier defaults are seeded equivalently to C#
  `ZoneContentManager.Build*NeutralMandatoryContent` (commits `ccf97c4`, `fc4cb46`).

### 3.5 Single-hero / multiple-castle UI ⚠️ (verify)

- Game mode `SingleHero` toggle and castle-count sliders already exist on the
  General/Layout tabs. ✅
- [ ] Confirm UI exposes `HubZoneCastles` and per-tier castle counts and that the
  "default generator strength" selection is reachable (commits `42533b7`,
  `05b6f91`, `d9e3891`).

### 3.6 Misc UI polish from upstream ⚠️

- [ ] Zone-visualization clarity tweaks (commits `bc630c9`, `99844fd`).
- [ ] Connection display when many connections exist (commit `ecd8f88`).
- [ ] Fix/verify MainWindow icons equivalent (commit `6edc025`) — N/A if Go uses
  its own toolbar glyphs, but confirm parity of actions.

---

## Suggested Implementation Order

1. **Model:** §1.1 `ContentRuleRowSave` + `Rules[]` and load-time migration,
   §1.5 variant-mapping data, §1.4 connection field.
2. **Service:** §2.1 rule engine, §2.2 rewire managers, §2.3 variants,
   §2.6 guard presets, §2.5 Proximity retirement; verify §2.4.
3. **UI:** §3.2 rule dialogs (depends on 1+2), §3.3 pickers, then §3.1 visual
   connection editor (largest), finishing with §3.4–§3.6 verification/polish.

Each step should be accompanied by tests under `test/models` and
`test/services` and validated with `go test ./test/... -count=1`.
