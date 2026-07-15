# Geometric Hub Topology + Lowest/Highest Zone Tiers

Add a new "Geometric Hub" topology (player hexagons around a central Hub, per the
`output/*.png` proof-of-concepts), two new neutral zone quality profiles
(Lowest/"Plastic" and Highest/"Platinum"), platinum preview assets, Hub zone
switched to the Highest profile, Lowest tier exposed in the UI, and the
`preview.Zone` tier refactored onto the `preview.ZoneTier` enum.

## For Future Agents

As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read `AGENTS.md` first. Hard rules: never touch `data/`, `internal/entities/template/`,
`internal/registry/` (read-only game data); stay cross-platform (Windows+Linux,
`path/filepath` only); every change ships with unit tests under `test/unit/`
mirroring the impl path (§4.6) and must not drop coverage (§2.3); receivers are
named `this`; one struct per file, camelCase file names.

### User-confirmed decisions (do not re-ask)

1. **Growth order** (each step round-robin across players/pairs, deterministic):
   hub → stable zones (2 per player) → merged corner (1 per adjacent-player pair)
   → 1 interior per hexagon → corner splits into pairs (1 more per pair) →
   additional interiors (unbounded, chained per rule 8).
   Degradation below full stable count: stable zones shared between adjacent
   players (zone connects both players + hub) → hub-only star (players connect
   directly to hub).
2. **Portals**: ALL connections that touch the Hub are `Portal` connections.
3. **Interiors**: no cap; keep chaining per rule 8 (each new interior connects to
   hub, to the other interiors of its hexagon, and to exactly one non-central
   hexagon ring zone).
4. **Tier→position**: interiors pop the highest-quality available plans, corners
   pop the lowest, stable zones take the remainder (naturally medium-preferring).
   Hub is always Highest.
5. **`neutralZone.Quality` renumbering approved**: Lowest=0, Low=1, Medium=2,
   High=3, Highest=4. No `.gen.json` migration needed (quality is derived from
   pool SIDs, counts are plain ints).
6. **Lowest gets its own mandatory-content rows** ("Edit zone content..."); Hub
   keeps using the existing `HubZoneMandatoryContent`.
7. **Platinum sprite look**: icy blue-white platinum (harmonizes with the Hub GUI
   blue `#82B4C8`), clearly distinct from silver; create `neutral_highest.png`,
   `neutral_highest_castle.png`, `neutral_highest_arena.png`.
8. **Hub sliders** (size/castles/mandatory content) also apply to Geometric Hub.
   The hub zone is ALWAYS created, as an EXTRA zone on top of the user-selected
   neutral counts (0 neutrals selected → hub + star). 
9. **"Double T5" pools** = the T5 pool SID list included TWICE (double weight),
   no T4 entries.
10. Plan lives in `plans/` (per AGENTS.md §4.7 — user overrode their original
    `./todo/` suggestion).

### Topology algorithm spec (derived from output/*.rmg.json PoCs, user-confirmed)

Let `P` = player count, `M` = user-selected neutral zone count (hub NOT included).
Zones: `Spawn-<letter>` players, `Hub`, `Neutral-<letter>` neutrals (letters
continue the generator's label sequence from `neutralZone.Plans`).

Slot classes, in fill order (round-robin for fairness, deterministic iteration):

| # | Slot class            | Count | Connections when present |
|---|-----------------------|-------|--------------------------|
| 0 | Hub                   | 1 (always, extra) | see below |
| 1 | Stable `s`            | up to 2·P (2 per player) | `P_i--s`; `s--corner` when corner exists, else `s--Hub` |
| 2 | Merged corner `c`     | up to P (1 per adjacent pair) | `sR_i--c`, `c--sL_j`, `c--Hub` (portal) |
| 3 | Interior r1 `x_i`     | up to P (1 per hexagon) | `x--sL_i`, `x--sR_i`, `x--Hub` (portal); NEVER to player |
| 4 | Corner split          | up to P (merged corner becomes pair `c_i`,`c_j`) | `sR_i--c_i`, `c_i--c_j` (single common edge), `c_j--sL_j`, `c_i--Hub`, `c_j--Hub` (portals) |
| 5 | Interior r2+ `x_i^k`  | unbounded | `x^k--Hub` (portal), `x^k--x^m` (all interiors of same hexagon), `x^k--` exactly ONE ring zone (a stable not yet used by another interior; matches 2P-7H PoC: G--I, G--C, G--M) |

Degradation (`M < 2P`, before class 1 completes):
- `M == 0`: star — every `P_i--Hub` (direct? PoC 4P-1M uses direct; but rule 11
  says hub connections are portals → use Portal, they touch the hub).
- `1 ≤ M ≤ P`: shared stable per adjacent pair: `P_i--z`, `z--P_j`, `z--Hub`
  (portal). Unserved players connect `P_i--Hub` directly (portal).
- `P < M < 2P`: continue round-robin: pairs with budget get 2 zones
  (`P_i--z1--z2--P_j`, both `z--Hub`), rest keep 1 shared (matches
  3P-0L-6M pattern where every stable also connects hub while corners are absent).
- Once corners fill in (class 2), the `s--Hub` edges are REPLACED by the hexagon
  edge path `s--c--Hub` (rule 5: pure hexagon has only edge connections; stables
  do NOT touch hub in the full shape — verified in 2P base PoC).

Rule 5 check: full hexagon = P, sL, cL, Hub, cR, sR with only adjacent edges,
never P--Hub, never sL--sR, never diagonal cuts.

Geometry (normalized [0,1], stamp `zone.GeneratorPosition`, use `layoutCenter`
const and `geometryHelpers.circlePoint`): hub at center; player `i` at angle
`θ_i` radius ≈0.45; stables at `θ_i ± 0.5·sectorHalf`, radius ≈0.42; merged
corner at pair mid-angle radius ≈0.24; split corners at mid-angle ±0.35·sectorHalf
radius ≈0.30 (each nudged toward its own hexagon); interior r1 at `θ_i` radius
≈0.26; interiors r2+ fan around `θ_i ± k·offset` radius ≈0.28. Tune visually
against the PoC PNGs; exact constants are implementer's choice, but the rendered
shape must read as hexagons (compare with `output/2Players-2Low-4Medium-1High.png`,
`3Players-3Low-6Medium-4High.png`, `4Players-0Low-13Medium-0High.png`).

Tier assignment: build the plan multiset from `neutralZone.Plans`; pop for
interiors (highest first), then corners (lowest first), stables get the rest
sorted so the pair closest to each player is the higher of what remains
(rule 9 volatile/stable behavior). Castle counts ride along with each plan.

Guards: border guard between zones via existing `GetBorderGuardValue`
(quality-based 10k/15k/20k/25k/30k after renumbering). Portal connections follow
`CreateRandomPortalConnections` precedent (`WithConnectionTypePortal`, road true,
weekly increment 0.15) but with quality-scaled guard values, names `Portal-Hub-<X>`.

### New profile numeric values (extrapolated from Low/Medium/High deltas)

Existing (see [internal/models/neutralZone/neutralZoneProfile.go](../internal/models/neutralZone/neutralZoneProfile.go)):
guard multipliers 1.1/1.4/1.8 (additive steps 0.3, 0.4 → pattern 0.2,0.3,0.4,0.5),
guarded 120k/240k/480k, guarded-per-area 1000/2000/3000, unguarded 25k/38k/80k,
ung-per-area 200/300/620, resources 30k/55k/90k, res-per-area 240/420/580,
primary/extra city guard 4k/2k, 8k/4k, 16k/8k.

**Lowest ("Plastic")** — ratio Low²/Medium, rounded:
- Layout `Sides`, guard reactions `[0,10,10,10,10,0]` (same as Low)
- GuardMultiplier **0.9**
- Guarded/unguarded pools: **T1** lists (`GetGuardedContentPoolT1List` / `GetUnguardedContentPoolT1List`)
- Resources pool: `registry StartZoneVeryPoor` (`content_pool_general_resources_start_zone_very_poor`)
- Guarded **60_000**, per-area **500**; Unguarded **16_000**, per-area **130**
- Resources **16_000**, per-area **140**
- Primary city guard **2_000**, extra **1_000**
- Primary + extra construction: `ExtraPoor` (`extra_poor_buildings_construction`)
- `Quality.GetGuardValue()`: **10_000**

**Highest ("Platinum")** — ratio High²/Medium, rounded:
- Layout **Center** (Highest is hub-reserved; hub keeps center layout via profile)
- Guard reactions `[0,10,10,20,10,0]` (same as High)
- GuardMultiplier **2.3**
- Guarded/unguarded pools: **T5 list twice** (append `GetGuardedContentPoolT5List()` to itself; same for unguarded)
- Resources pool: `registry TreasureZoneRich` (`content_pool_general_resources_treasure_zone_rich`)
- Guarded **960_000**, per-area **4_500**; Unguarded **168_000**, per-area **1_280**
- Resources **147_000**, per-area **800**
- Primary city guard **32_000**, extra **16_000**
- Primary + extra construction: `UltraRich` (`ultra_rich_buildings_construction`)
- `Quality.GetGuardValue()`: **30_000**

All registry SIDs verified to exist already ([internal/registry/resourcesContentPoolValues.go](../internal/registry/resourcesContentPoolValues.go),
[internal/registry/buildingsConstructionSidValues.go](../internal/registry/buildingsConstructionSidValues.go),
T1/T5 pool getters in guarded/unguardedContentPoolValues.go) — registry is
read-only and needs NO changes.

### Subagent strategy (AGENTS.md §3.4)

- Mechanical/repetitive plumbing (Phase 4 DTO/mapper/validator field fan-out,
  test-file scaffolding): delegate to cheap models (gpt-5.5 / sonnet-5) with
  precise self-contained briefs listing exact files+fields.
- Topology algorithm (Phase 5) and the profile/hub semantics (Phase 1): fable-5
  or opus-4.8 (design-heavy, taste > 7 required).
- Asset recoloring (Phase 3): implement directly; verify sprites visually
  (vision) or via the asciiview technique from repo memory.
- Read-only exploration/verification runs: cheap models, parallelized.
- Final review of the topology implementation: fable-5/opus-4.8.

---

## Phase 1: Quality enum + Lowest/Highest profiles + Hub on Highest profile
Status: Not started

- [ ] Renumber `neutralZone.Quality`: `QualityLowest=0, QualityLow=1, QualityMedium=2, QualityHigh=3, QualityHighest=4` in [internal/models/neutralZone/neutralZoneQuality.go](../internal/models/neutralZone/neutralZoneQuality.go); extend `GetGuardValue` (10k/15k/20k/25k/30k).
- [ ] Extend `GetQualityFrom(zone)`: Highest when pools contain `_t5_` AND resources pool is `treasure_zone_rich` (or construction is `ultra_rich_*`); `_t4_/_t5_` → High; `_t3_` → Medium; `_t2_` → Low; `_t1_` → Lowest.
- [ ] Add Lowest/Highest cases to `NewNeutralZoneProfile` in [neutralZoneProfile.go](../internal/models/neutralZone/neutralZoneProfile.go) with the numeric values above (one private factory per tier, matching the existing style).
- [ ] Update every quality switch site (found by exploration, complete list):
  `Plan.GetBalanceScore` + `Plans.GetTier` (extend monotonically: Lowest=0 … Highest=4) + `Plans.sort` (ordering keeps working via renumber) in [neutralZonePlan.go](../internal/models/neutralZone/neutralZonePlan.go);
  `bucketNeutralsPerPlayer` in [fractalTopology.go](../internal/services/template_generator/providers/topology/fractalTopology.go) (bucket Lowest with low, Highest with high);
  `CreateContents` + `neutralRowsForQuality` in [mandatoryContentProvider.go](../internal/services/template_generator/providers/mandatoryContentProvider.go) (Lowest rows → new config field, Highest → hub content, no separate rows);
  `neutralCastleTarget` in [manualReapply.go](../internal/services/connection_editor/manualReapply.go);
  `QualityLabels` in [zoneEditor.go](../internal/services/connection_editor/zoneEditor.go) → `["Lowest","Low","Medium","High"]` (Highest deliberately NOT offered; indices align with renumbered enum).
- [ ] `CreateHubZone` in [topologyBase.go](../internal/services/template_generator/providers/topology/base/topologyBase.go#L184): replace ALL inlined profile values (pools, resource pool, guarded/unguarded/resource values+per-area, construction sids, guard multiplier, reactions, layout) with `neutralZone.NewNeutralZoneProfile(neutralZone.QualityHighest)` — same application pattern as `CreateNeutralZone`. Keep hub-specific structural fields (name, size param, `GuardCutoffValue 2000`, `DiplomacyModifier -0.5`, castles, roads, mandatory content, biome).
- [ ] Update/extend unit tests for every touched public function (folders under [test/unit/internal/models/neutralZone/](../test/unit/internal/models/neutralZone/), [test/unit/internal/services/connection_editor/](../test/unit/internal/services/connection_editor/), topologyBase, mandatoryContentProvider). Hub golden expectations in tests/fixtures WILL change — update deliberately, verifying the new values equal the Highest profile.

### Verification Plan
- `go build ./...` then `go test ./test/... -count=1` — all green.
- Coverage task before/after; total must not drop (§2.3).
- Spot-check: generate a Hub-and-Spoke template via existing unit/integration harness; assert hub zone has T5-doubled pools, `treasure_zone_rich` resources, `ultra_rich_buildings_construction`.

### Phase Summary
_(write when phase completes)_

## Phase 2: preview.ZoneTier enum refactor + classifier + legend
Status: Not started

- [ ] Change `preview.Zone.Tier` from `int` to `preview.ZoneTier` in [internal/models/preview/previewZone.go](../internal/models/preview/previewZone.go); delete the stale comment.
- [ ] `ClassifyZoneTier` in [zoneClassifier.go](../internal/services/preview_service/zoneClassifier.go) returns `preview.ZoneTier`: Platinum (t5 + rich-treasure/ultra-rich detection, mirror `GetQualityFrom`), Gold (t4/t5), Silver (t3), Bronze (t2), Plastic (t1 and default/unknown — per user, tier 0 was unused and Plastic becomes the default).
- [ ] Update all `Tier` readers: `buildPreviewZones` ([previewLayoutService.go](../internal/services/preview_service/previewLayoutService.go)), `getNeutralZoneAsset` switch in [assetProvider.go](../internal/services/asset_provider/assetProvider.go) (Platinum→`neutral_highest`, Gold→`neutral_high`, Silver→`neutral_medium`, Bronze→`neutral_low`, Plastic/default→`neutral_none`), `zoneColors`/`zoneLabel` in [app/gui/utils/draw.go](../app/gui/utils/draw.go) (add Plastic + Platinum colors/labels; hub override stays).
- [ ] Themes: add `ColorsPreview` entries for Plastic (muted grey-brown, distinct from bronze) and Platinum (icy blue-white) in [app/gui/themes/colors.go](../app/gui/themes/colors.go) — colors live ONLY in themes.
- [ ] Legend ([app/gui/constants/legend.go](../app/gui/constants/legend.go)): entries become `Player`, `Plastic (T0)`, `Bronze (T1)`, `Silver (T2)`, `Gold (T3)`, `Hub (T4)`, `Road`, `Portal`.
- [ ] `connection_editor.GetZoneTier` ([connectionEditor.go](../internal/services/connection_editor/connectionEditor.go)): map `_t1_`→Bronze presets (keep 3-tier guard preset enum unchanged — guard presets are a separate concept), hub/platinum zones → `ZoneTierGold` instead of Bronze (fixes hub connections getting bronze guard presets).
- [ ] Manual zone editor dialog quality dropdown keeps working with 4 labels; `preview` tier drawing in [zoneEditorDialog*.go](../app/gui/dialogs/zoneEditorDialog.go) picks up new tiers via `zoneColors`.
- [ ] Update all affected unit tests (zoneClassifier, asset_provider, connection_editor, preview_service).

### Verification Plan
- `go build ./...`, `go test ./test/... -count=1`, coverage no-drop.
- Manual-editor tier smoke: unit test asserting a t1-pool zone classifies Plastic and a hub-profile zone classifies Platinum.

### Phase Summary
_(write when phase completes)_

## Phase 3: Platinum preview assets
Status: Not started

- [ ] Write a small cross-platform Go generator under `tools/` (module `tools/go.mod` exists; do NOT use `tmp/` — the user's air watcher deletes it) that derives icy blue-white platinum sprites from the gold ones: `neutral_high.png`→`neutral_highest.png`, `neutral_high_castle.png`→`neutral_highest_castle.png`, `neutral_high_arena.png`→`neutral_highest_arena.png` in [internal/services/asset_provider/assets/](../internal/services/asset_provider/assets/). Recolor only saturated (gold) pixels: desaturate toward luminance then tint cool (approx `R'=lum·0.82, G'=lum·0.96, B'=min(255, lum·1.18)`), keep ring/shading/alpha. GOTCHA (repo memory): use `image.NRGBA`/`SetNRGBA`, NOT premultiplied `color.RGBA`, or edges corrupt.
- [ ] Visually verify sprites (view the PNGs) — must be clearly distinguishable from `neutral_medium*` silver at 96×96 and at typical preview scale; iterate multipliers if too close to silver.
- [ ] Register in `neutralAssetNames` in [assetProvider.go](../internal/services/asset_provider/assetProvider.go): add `neutral_highest`, `neutral_highest_castle`, AND the missing `neutral_none`, `neutral_none_castle` (needed by Plastic tier from Phase 2).
- [ ] Unit tests: asset loader returns the new sprites; `getNeutralZoneAsset` Platinum/Plastic paths covered.

### Verification Plan
- `go test ./test/... -count=1`; render a preview PNG through the existing preview generator test harness with a platinum-tier zone and confirm no fallback-to-bronze.
- Visual check of generated `.png` files (vision or asciiview technique from repo memory).

### Phase Summary
_(write when phase completes)_

## Phase 4: Lowest tier through config + UI
Status: Not started

Mechanical fan-out — good subagent candidate (precise field list below).

- [ ] `dtos.EditorStateDto` ([internal/dtos/editorStateDto.go](../internal/dtos/editorStateDto.go)): add `NeutralLowestNoCastleCount`, `NeutralLowestCastleCount`, `NeutralLowestZoneCastles` (+ JSON tags matching existing naming), `LowestNeutralContentRows`; include in `zoneCountOptionsChanged`, `zoneOptionScalarsEqual`, `EqualsIgnoringManualEdits` field-group comparators, defaults (castles-per-zone default 1 like other tiers), and the reflection tripwire test will force completeness.
- [ ] `CastleSettingChanges` ([internal/dtos/editor_state_dto/castleSettingChanges.go](../internal/dtos/editor_state_dto/castleSettingChanges.go)) + `DiffCastleSettings`: add `NeutralLowest`.
- [ ] `config_inner.AdvancedSettings` ([advancedSettings.go](../internal/models/config/config_inner/advancedSettings.go)) + `GeneratorConfig` mandatory-content field for Lowest rows; `generatorConfigMapper.FromEditorState`/`mapZoneConfig` ([generatorConfigMapper.go](../internal/mappers/generatorConfigMapper.go)).
- [ ] `zoneLabelProvider.CreateNeutralZonePlans` ([zoneLabelProvider.go](../internal/services/zones/zoneLabelProvider.go)): Lowest plans created FIRST (before low no-castle), so label order stays tier-ascending.
- [ ] `editorStateValidator` ranges; `EditorState.UpdateCurrentState` simple-mode clearing ([app/gui/models/editorState.go](../app/gui/models/editorState.go)); `ApplyCastleSettingChanges` (manualReapply) per-tier propagation for Lowest.
- [ ] GUI [layoutPanelZones.go](../app/gui/panels/layoutPanelZones.go) + [layoutPanel.go](../app/gui/panels/layoutPanel.go): "Lowest tier" section ABOVE "Low tier" (same `getNeutralTierSectionWidget` structure: No castle / With castle / Neutral castles per zone / Edit zone content...), new sliders + content-rows dialog wiring, LoadFromState/SaveToState.
- [ ] Zone content default rows for Lowest (mirror how Low defaults are built; check `DefaultPlayerZoneContentRows`-style helpers and the zone content dialog plumbing in [app/gui/drivers/zoneContent*.go](../app/gui/drivers/)).
- [ ] `mandatoryContentProvider`: Lowest rows emitted as `mandatory_content_lowest`-style group consistent with existing naming; `CreateContentsForZones` re-derivation handles Lowest via `GetQualityFrom` (already wired in Phase 1).
- [ ] Highest is NOT added to any UI surface (hub-reserved).
- [ ] Unit tests for every touched function per §4.6 layout; the DTO reflection tripwire and fuzz-parity tests must be extended with the new fields.

### Verification Plan
- `go build ./...`, `go test ./test/... -count=1`, coverage no-drop.
- Round-trip test: EditorStateDto with Lowest counts → mapper → config → `CreateNeutralZonePlans` yields Lowest plans first; save/load `.gen.json` preserves the new fields.
- Legacy `.gen.json` without the new keys loads with zero counts (silent default).

### Phase Summary
_(write when phase completes)_

## Phase 5: Geometric Hub topology
Status: Not started

- [ ] Enum: `TopologyGeometricHub MapTopology = "GeometricHub"` in [mapTopology.go](../internal/models/config/config_inner/mapTopology.go) + alias in [types.go](../internal/models/config/types.go).
- [ ] Descriptor in [internal/common/topologies.go](../internal/common/topologies.go): field `GeometricHub`, label "Geometric Hub", description; dropdown order: directly after `HubAndSpoke` (existing Geometric/Hub stay for now — future removal is a separate task).
- [ ] New service `internal/services/template_generator/providers/topology/geometricHubTopology.go` (`GeometricHubTopologyService`, embeds `base.TopologyBase`), implementing the algorithm spec above:
  - slot construction + round-robin fill from `neutralZone.Plans`,
  - tier→slot assignment (interiors high-first, corners low-first, stables medium),
  - hub via `CreateHubZone` (Highest profile from Phase 1; hub size/castles/mandatory content from config — reuse Hub-and-Spoke wiring incl. `hubIsHoldCity`),
  - neutral zones via `CreateNeutralZone`, spawns via `CreateSpawnZone`,
  - connection names/roads consistent with base helpers; all hub-touching connections Portal-type with quality-scaled guards,
  - `GeneratorPosition` stamped for every zone (hexagon geometry),
  - honor `NoDirectPlayerConnections` trivially (topology has no P–P edges except none), support `RandomPortals` via `CreateRandomPortalConnections` like Hub-and-Spoke.
  - Nest private helper structs/files in a sibling folder if they're single-use (§4.4).
- [ ] Dispatch: `case config.TopologyGeometricHub` in [topologyProvider.go](../internal/services/template_generator/providers/topologyProvider.go) (pass `holdCityNeutralLabel`/hub flag like Hub-and-Spoke).
- [ ] Preview: add to `isFixedGeometryTopology` in [layoutGeometry.go](../internal/services/preview_service/layoutGeometry.go) (exact positions, NO scatter relaxation — repo memory: relaxation destroys deterministic geometry). Hub detection by name `"Hub"` already works.
- [ ] Hub UI gating: hub size/castles sliders + hub mandatory content visible when topology is HubAndSpoke OR GeometricHub ([layoutPanelZones.go](../app/gui/panels/layoutPanelZones.go), search `TopologyHubAndSpoke` gates); `mandatoryContentProvider.hubContentGroup` gate likewise; check `config.CanHonorNeutralSeparation` and any other topology-switch sites (`grep TopologyHubAndSpoke` across internal/ + app/).
- [ ] Unit tests (`test/unit/internal/services/template_generator/providers/topology/geometricHubTopology/`): per public function; scenario tests for M=0 (star), M≤P (shared stables), M=2P, full shape, corner splits, multi-interior chains; assertions on portal types touching hub, no P--Hub edge when corners exist... (careful: star/degraded cases DO have P--Hub), rule-5 no-diagonals, tier placement, castle counts. Validate with `-count=20` (repo memory: topology tests can be flaky at 1 run if randomized; prefer deterministic construction so they aren't).
- [ ] Regenerate visual sanity previews: generate templates matching several PoC configs (2P-2L-4M-1H, 3P-3L-6M-4H, 4P-0L-13M-0H, 2P-2L-4M-7H) and compare rendered PNG + connection graph against `output/` PoCs (graph shape, not exact letters).

### Verification Plan
- `go build ./...`, `go test ./test/... -count=1` (+ `-count=20` on the new topology package), `go test -tags=integration_test ./test/integration/... ./test/performance/... -count=1`.
- Graph equivalence spot-checks vs PoC JSONs for the four configs above.
- Coverage no-drop; lint (`golangci-lint-v2 run ./... --issues-exit-code=0`) no new findings.

### Phase Summary
_(write when phase completes)_

## Phase 6: Final verification & review
Status: Not started

- [ ] Full suite: `go build ./...`; `go test ./test/... -count=1`; gated `go test -tags=integration_test ./test/integration/... ./test/performance/... -count=1`.
- [ ] Coverage report task; compare against pre-work baseline (record baseline % before Phase 1 starts: run the coverage task and note it here: baseline = ____).
- [ ] Lint report; fix new findings (funcorder: unexported methods below exported; golines).
- [ ] GUI smoke: `go run .` — select Geometric Hub, drag zone counts across the degradation ladder (0→1→P→2P→full→splits→interiors), watch live preview; verify legend shows the five tier entries; manual zone editor offers Lowest; PNG export renders hub with platinum sprite and lowest zones with the open-ring plastic sprite.
- [ ] Independent review of the topology service + profile changes by fable-5 or opus-4.8 (per AGENTS.md §3.4), optionally gpt-5.5 as extra perspective.
- [ ] Update repo memory (conventions.md: new tiers, new topology, enum renumber, asset additions).

### Verification Plan
- All commands above green; user visually accepts the generated previews.

### Phase Summary
_(write when phase completes)_

## Final Recap
_(write when all phases complete: summary of the entire piece of work)_

## Deployment Plan
_(write when all phases complete: step-by-step deployment instructions)_
