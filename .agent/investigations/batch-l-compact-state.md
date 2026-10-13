# Batch L: compact manual deltas and derivable editor state, feasibility report

Review item [§2.4 / O11](../backlog/review-gpt-6-astra-09-07.md). Plan:
[batch-l-compact-state-investigation.md](../plans/batch-l-compact-state-investigation.md).
Base `45876cf` on `AD/compact_state`, Windows/amd64, go1.27.2. Every number below traces to a
raw file in gitignored `.agent/memories/batch-l/` (named in brackets). No production code,
schema or fixture changed; all measurement code was a throwaway `cmd/batchlprobe`.

**This is a decision document.** It presents alternatives with measured benefits, risks and
migration limits. It authorizes nothing: no representation change, no seed or determinism
infrastructure, no delta persistence and no migration.

## 1. Summary

- **What is saved today.** One manual edit saves the **entire** variant-0 layout: every zone
  as a full `.rmg.json` zone plus position/tier sidecars, and every connection.
  `reapplyManualEdits` then replaces the generated zones and connections wholesale. One drag
  in a 6-player GeometricHub layout takes the file from 8 KB to **202 KB** (the owner's
  Colosseum samples: 188–191 KB). More than 95% of that is the snapshot, and it is highly
  redundant: gzip-9 shrinks it to about 3%.
- **What it costs live.** With that snapshot, every idle frame deep-clones the state twice
  (**80 µs, 201 KB, 1,960 allocations**), and three times while a regeneration is pending
  (128 µs, 301 KB). Without manual edits an idle frame costs 6 µs and 10 KB.
- **The snapshot overrides newer settings.** After any manual edit, changing map size,
  either density, zone sizes, neutral-stack or border-guard strength, abandoned outposts,
  City Hold or (in advanced mode) guard randomization regenerates the template, but the old snapshot values are written back over
  it. Those settings then silently stop applying (§6). Recorded as review **§1.16**.
- **Generation is deterministic except where it is random by design.** In 20 repeats, six of
  the seven topologies and the tournament builder produced byte-identical templates; only
  the Random topology and the RandomPortals option vary (§5).
- **Every prototype round-trips losslessly** on all 40 manual-edit fixtures (snapshot,
  export bytes and preview layout) except one variant that tried to derive roads (§7).
- **Recommendation** (§12): option **C3**, which stores the structure plus per-field
  overrides, with every other field derived from the current settings. It cuts
  deterministic layouts by about 78% (108 → 23 KB for the 6-player move), roughly halves
  Random layouts, needs no seed, stays portable, and fixes §1.16 for every derived field. Its
  costs are a behaviour change for derived fields, a provenance gap when migrating old
  snapshots, and a determinism contract the test suite must enforce.

## 2. Current design (code-verified)

- **Persisted:** [`ManualZoneSave`](../../internal/entities/editor_state/manualZoneSave.go)
  `{Zone template_entity.Zone, GeneratorPosition, GeneratorRing, ManualPosition, Quality}`
  and `ManualConnectionSave{Connection, IsUserAdded}`, inside the flat schema-v2 root.
- **Live:** `editor_state_model.ManualEditSettings{ManualZones []template_model.Zone,
  ManualConnections []ManualConnectionSave}` lives inside every `EditorState` copy:
  current, previous and next in [`app/gui/models.EditorState`](../../app/gui/models/editorState.go).
- **Apply:** `drivers.State.ApplyEditedZones` → `UpdateTemplate` → `SetManualEdits`. That
  stores **every** zone and connection, and the equality check converts both sides to
  entities before `reflect.DeepEqual`.
- **Reapply:** [`handleGenerateTemplate`](../../app/gui/drivers/stateGeneration.go) generates
  fresh, decides with `DecideManualEditReapplication`, and calls `reapplyManualEdits`. That
  runs `ReapplyCastleSettings` when castle counts changed, then `UpdateTemplate`, which
  assigns the saved lists wholesale and rebuilds template-level mandatory content and roads.
  Nothing is matched per zone.
- **Discard:** layout-defining options (`LayoutDefiningOptionsChanged`) and effective-mode
  transitions (Batch E) drop the snapshot.
- **Identity:** zones and connections are referenced by name only. There is no ID, seed or
  generator version anywhere in the document.

## 3. Persisted-size measurements [phase1-sizes.tsv, phase1-fixtures.txt]

The fixture matrix covers 4 layouts × 11 cases (untouched, applied-unchanged and 9 edits).
Each fixture was built through the real driver (`Generate` → `ApplyEditedZones`) and saved
by the production handler. All 44 reload with no rejection and no normalization, and
reproduce the export byte for byte.

| Fixture (move) | Zones / conns | Written | Compact | gzip-9 | No edits (written) |
| --- | --- | --- | --- | --- | --- |
| square2p (Square, 96) | 6 / 6 | 35,138 | 20,482 | 2,427 | 5,345 |
| tournament2p (Random, 160) | 6 / 4 | 33,093 | 19,729 | 2,454 | 5,349 |
| random8p (Random, 240) | 24 / 60 | 167,585 | 89,687 | 5,718 | 5,347 |
| geohub6p (Colosseum settings, 368) | 31 / 48 | 202,434 | 113,211 | 6,480 | 8,015 |
| Colosseum v3 (owner sample) | 31 / 42 | 190,805 | 107,545 | 6,074 | n/a |

Largest groups in geohub6p-move (compact value bytes, all zones):

| Group | Bytes |
| --- | --- |
| content pools | 28,130 |
| roads | 15,900 |
| main objects | 14,003 |
| content limits | 12,772 |
| other scalars | 6,178 |
| content values | 5,348 |
| biomes | 4,927 |
| position sidecars | 3,603 |
| mandatory content | 1,530 |
| guard reaction | 1,395 |
| all connections | 13,013 |

The supported legacy fixtures (v0/v1/v2, 1 zone each) are 5.6–5.7 KB and load through the
migrator. In settings, the 14 default player content rows take 1,946 of square2p's ~3.4 KB.

## 4. Live-state measurements [bench/phase1-bench-summary.tsv]

The setup was go1.27.2, `GOMAXPROCS=4`, 1 s benchtime, 5 independent processes. The table
gives the median with the min–max range.

| geohub6p | No edits | With move snapshot |
| --- | --- | --- |
| `EditorState.Clone` | 2.7 µs / 5.2 KB / 55 | 43 µs (41–46) / 100 KB / 980 |
| `ManualEditSettings.Equals` (equal) | 3 ns | 186 µs (182–199) / 254 KB / 436 |
| `SetManualEdits` (unchanged) | 7 ns | 243 µs (226–280) / 371 KB / 1,322 |
| encode (in memory) | 27 µs | 442 µs (440–583) / 342 KB |
| decode (in memory) | 30 µs | 552 µs (507–628) / 368 KB / 3,443 |
| `AutoRegenerate`, idle frame | 6.3 µs / 10.5 KB / 110 | **80 µs (77–97) / 201 KB / 1,960** |
| `AutoRegenerate`, debounce pending | 9.0 µs / 15.7 KB | **128 µs (117–146) / 301 KB / 2,940** |

`AutoRegenerate`'s bytes are exactly 2× and 3× `Clone`, confirming two clones when idle
(current and previous) and three when pending. `EqualsIgnoringManualEdits` costs 0.4–0.8 µs
and allocates nothing in every case. random8p is about 85% of geohub6p; the 2-player layouts
are about 18%. Tuning is out of scope (L7). A smaller in-memory snapshot or a non-cloning
regeneration check would remove most of this cost, and this report records it without
assigning it.

## 5. Determinism and identity [phase2-determinism.tsv, phase2-template-level-diff.txt]

Each case was generated 20 times and compared with run 1. Controls: a clone reports equal,
and one changed field is detected.

| Case | Result |
| --- | --- |
| Circles, GeometricHub, Square, Geometric, Cross, Fractal (4p); tournament 2p; square2p, geohub6p, tournament2p layouts | 19/19 identical: zone set and order, endpoints, names, positions, every zone and connection value, template level |
| Random (4p and random8p) | zone **set** stable; order, endpoint pairs, connection names, positions, roads and `orientation` differ in ~every run; all other zone values identical |
| Square + RandomPortals | zones and positions stable; portal pairs, names and roads differ |

- **Random sources:** the only `math/rand` call sites are
  [`randomTopology.go`](../../internal/services/template_generator/providers/topology/randomTopology.go)
  and the portal derangement in
  [`topologyConnectionService.go`](../../internal/services/template_generator/providers/topology/base/topologyConnectionService.go).
  Determinism elsewhere is empirical, not a contract: there is no seed, no generator version,
  and nothing stops a future change (or map iteration) from adding variance.
- **Side observation:** for the Random topology a reapplied template takes `orientation` from
  the fresh generation while its zones come from the snapshot.
- **Identity facts:**
  - In all 260 determinism generations zone and connection names were unique and non-empty.
  - Zones are referenced by connection `from`/`to`/`guardZone`.
  - Roads reference connections (`{"type":"Connection","args":[name]}`) and main objects.
  - `mandatoryContent`, `guardMatchGroup` and connection names derive from zone labels.
  - `contentCountLimits` depends only on the player count.
  - [`NextFreeZoneLabel`](../../internal/services/connection_editor/zoneEditorService.go)
    reuses a deleted label.
  - `EnsureConnectionNames` names only nameless connections (order-dependent `-2`, `-3`
    suffixes) and never repairs duplicates. It runs inside `UpdateTemplate`, so a stored
    snapshot never holds a nameless connection.
  - The editor has no rename and no endpoint change: an endpoint change is a delete plus an
    add.

## 6. The snapshot-wins bug (review §1.16) [phase2-l6.tsv]

The test ran on the real driver: `Generate` → one move → `ApplyEditedZones` → change one
setting → `Generate`. The fresh template was captured before reapplication, so B (base), S0
(snapshot), F (fresh) and R (result) are paired. A field counts as stale when B≠F and R=S0≠F.
Controls: renaming the template changes nothing; a player-castle change does update the
castles in R.

| Setting | Stale in the result |
| --- | --- |
| `MapSize` | all six content values of every zone |
| `ResourceDensityPercent` | `resourcesValue`, `resourcesValuePerArea` |
| `StructureDensityPercent` | four guarded/unguarded values |
| Player / neutral / hub zone size | `size` |
| `NeutralStackStrengthPercent` | `guardMultiplier`, main-object guards |
| `BorderGuardStrengthPercent` | main-object guards and **connection `guardValue`** |
| `SpawnAbandonedOutposts`, `AbandonedOutpostCount` | main objects and roads (outposts never appear) |
| `CityHold` (hub topology) | hub main objects, roads, biomes |
| `HubZoneCastles` | hub biomes (castles themselves are reapplied) |
| `GuardRandomization` (advanced mode only) | `guardRandomization` of every zone |

An R≠F column separates settings that are **reapplied** (R≠S0) from settings with **no
zone effect**. Castle settings and remote-foothold counts are reapplied, but their result
roads still differ from a fresh generation (R≠F on `roads`). So the reapply path is
consistent, not identical to regeneration.

**No zone effect at the tested settings:** content rows, XP, lost-start rules, City Hold
days, game mode, hero counts, `ExperimentalMapSizes`, `MatchPlayerCastleFactions`,
mode-preserving victory changes and tournament rules. **Inactive** in the mode tested:
`GuardRandomization` in simple mode (fixed at 0.05), simple-mode tier castles,
advanced-mode `NeutralZoneCastles`, arena count day, tournament fields outside tournament.
Bans, value overrides and bonuses are template-level only. Layout-defining and mode changes
discard the snapshot.

**Coupling with this item.** A correct fix needs to know which fields the user changed. The
current snapshot cannot say: an untouched derived value and a deliberate edit look the same.
The two stop-gaps available without provenance are both lossy. One treats zone-affecting
settings as layout-defining, so edits are discarded with a warning. The other always
re-derives those fields, which also overwrites deliberate per-zone edits of the same fields.
Option C provides the provenance.

## 7. Alternatives and prototypes [phase3-prototypes.tsv, bench/phase3-prototype-bench-run*.tsv]

All the prototypes flatten each zone or connection save to its top-level keys and
reconstruct `template_model` values. Each was checked with three comparators: raw snapshot
deep-equal, `.rmg.json` export bytes through the production mapper, and preview layout.
Negative controls proved each comparator able to fail. C and D were encoded against one
fresh generation and decoded against an **independent** second generation.

| Option | Idea | geohub6p move | random8p move | square2p move | Round trip |
| --- | --- | --- | --- | --- | --- |
| Today | full snapshot, compact | 108,326 | 86,238 | 17,035 | n/a |
| Today, manual section indented | | 194,426 | 162,245 | 29,800 | n/a |
| **A** | exact stored base + delta | 130,932 all-in | 105,704 | 19,832 | 40/40 |
| **B** | seed + generator version + delta vs regenerated base | 2,788 all-in | not reconstructable | 638 | 30/30 against a regenerated base; random8p only against its stored base |
| **C1** | structure + overrides, roads always stored | 48,350 | 47,722 | 6,072 | 40/40 |
| C2 | as C1, roads derived by `UpdateTemplate` | 23,297 | 17,996 | 2,457 | **0/40** (castle roads are never re-created) |
| **C3** | structure + overrides, roads as ordinary overrides | **23,297** | 48,058 | **2,457** | 40/40 |
| **D** | today's semantics, omit only settings-independent fields | 99,638 | 88,096 | 13,834 | 40/40 |

The bytes cover the manual section only. Settings add 3.4–4.8 KB to every option, and the
prototype keys carry a `zone.`/`connection.` prefix, so C and D are slightly pessimistic.
Across all 40 fixtures, the delta of the unchanged-apply case is 29 B, and the largest
deltas (up to 7.2 KB) come from adding or deleting a zone.

**Encode + decode, geohub6p move:** today 1.11 ms (1.00–1.24) / 785 KB / 3,806 allocs. The
generic map-based A/B prototype takes 2.81 ms (2.64–3.19) / 1.93 MB / 11,143 allocs. It is
naive; a typed implementation would not keep the per-key maps.

### A. Exact stored base + delta

- **Benefit:** an explicit, minimal edit record. Revert-to-base and "what did I change" are
  trivial.
- **Cost:** the base has to live somewhere. Inside the document, all-in grows 16–23% over
  today. Outside it, the document is no longer portable and breaks when the base is lost.
- **Verdict:** no size benefit. Only useful as an in-memory structure.

### B. Regenerated base + delta (seed + generator version)

- **Benefit:** the smallest by far (≈2.7 KB all-in for 6 players).
- **Requirements:**
  - a seeded, injectable RNG through `randomTopology` and the portal derangement (not
    authorized here);
  - a persisted seed;
  - a pinned generator version.
- **Risk:** the delta is only meaningful against byte-identical regeneration. Any generator
  change silently re-bases every saved document, including profile tweaks, ordering and new
  fields. Old documents would need the old generator, or a fallback to a full snapshot plus
  a refusal or warning.
- **Old full snapshots:** cannot be converted, because their seed is unknown.
- **Verdict:** the highest-risk option. Today it works only by accident for six of the seven
  topologies.

### C. Structure + derived-field overrides

- **Stored per zone:** name, tier, generator position and ring, manual position, plus a
  castle count and `origin: added` for editor-made zones.
- **Stored per connection:** name, endpoints, type, `isUserAdded`.
- **Everything else is derived and stored only when different:**
  - generated zones: the same-named zone of a fresh generation with the **current**
    settings;
  - editor-made or re-tiered neutral zones: `CreateZoneEditorNeutralZone`;
  - connections: `CreateZoneEditorConnection` plus `ChangeZoneEditorConnectionType`.
- **Benefits:**
  - −78% for deterministic layouts (C3) and −44% for Random (C3 ≈ C1, +0.7%, because
    Random roads differ per generation);
  - portable, and no seed needed;
  - derived fields follow settings, which fixes §1.16 for every non-overridden field;
  - revert-to-base becomes "drop all overrides".
- **Risks:**
  - **Behaviour change:** untouched fields now track settings and generator improvements.
  - **Provenance gap at migration:** old snapshots cannot tell "set equal to the derived
    value" from "never touched".
  - **Determinism contract:** non-random zone fields must stay deterministic for given
    settings. This holds for every topology in §5, including Random's non-road fields, and
    the decode-against-an-independent-generation round trip exercised it.
  - **Roads:** C2 shows they are not derivable by reconciliation. They are derived from the
    fresh generation, so for Random they are stored.
- **Verdict:** recommended (§12).

**Live state under C.** The persisted form alone does not change the per-frame cost: if
`ManualEditSettings` keeps expanded zones live, `Clone` and `AutoRegenerate` copy exactly
what they copy today. To shrink them, the live model would hold the compact structure and
overrides and expand on demand: for the editor dialog, regeneration and preview. That makes
the expansion a per-use cost instead of a per-frame copy. The alternative, which is
independent of any format, is to stop deep-cloning the snapshot in the regeneration check
(§12.3). Neither is measured here.

### D. Representation-only (today's semantics)

- **Saving:** small: −8% (geohub6p) to −25% (tournament2p), and +2% for random8p, because
  only pools, guard reaction, content limits and mandatory content are
  settings-independent.
- **Zero-semantics alternatives:** compact (unindented) JSON alone cuts the whole geohub6p
  move file by 44% (202 → 113 KB), and gzip by 97%; both cost human readability.
- **Verdict:** a cheap complement to any option, not a solution. It also leaves §1.16
  unfixed.

## 8. Lossless omission of candidate fields (today's generator, matrix-wide)

The derivation rates come from C1 over 670 zone records. "Lossless" means equal to a
derivation in every zone the edit did not touch.

| Field | Derived from | Lossless today | Under today's semantics | If defaults or profiles change |
| --- | --- | --- | --- | --- |
| guarded/unguarded/resources pools | zone kind + tier ([`neutralZoneProfile.go`](../../internal/common/common_zones/neutralZoneProfile.go)) | yes (666–670/670) | yes | the saved document follows the new profile |
| `guardReactionDistribution` | kind + tier | yes (670/670) | yes | follows |
| `contentCountLimits` | player count (layout-defining) | yes (670/670) | yes | follows |
| `mandatoryContent` (group names) | zone label | yes (overrides only on edited zones) | yes | follows |
| zone classification (kind) | name prefix ([`zoneNameType.go`](../../internal/helpers/zone_helpers/zoneNameType.go)) | yes, never stored | n/a | n/a |
| zone tier (`quality` sidecar) | inferable from pools, but inference can disagree after edits | **no**, keep it | no | n/a |
| content values ×6, `size`, `guardMultiplier`, main-object guards, connection `guardValue` | settings via `GenerationTuning` | yes (≥ 99.4%) | **no**: omitting changes semantics (fixes §1.16) | follows settings and generator |
| main objects | kind, tier, castle count, settings | yes (670/670) | no | follows |
| biomes | kind, tier, hub castles | yes (670/670) | no | follows |
| roads | connection names + castles; not reproducible by `Reconcile` | only for deterministic topologies, via a fresh generation | no | follows |
| generator position and ring | topology | deterministic topologies only | n/a | must stay stored for Random and old saves |
| manual position | user | **no** | n/a | n/a |
| default player content rows (settings, 1.9 KB) | `GetDefaultPlayerZoneContentRows` | yes when unchanged | yes | **silently changes a saved document's rules**; needs a defaults version |

**Nil and empty values:** with the repository options, nil and `[]` stay distinct for the
non-`omitempty` pools. `omitempty` slices already collapse empty into absent today. A
derived representation needs an explicit `unset` list to store "absent where the derivation
has a value" (e.g. `road` nil or a removed `mandatoryContent`). The prototypes carry one,
and the nil-empty fixtures round-trip.

## 9. Semantics the owner must choose

- **Exact base vs regenerated base.** The review asked for these two to be compared. Exact
  base (A) is safe but has no size benefit. Regenerated base (B) is tiny but brittle: it
  needs seed infrastructure and breaks on every generator change. C sidesteps both: the
  *structure* is stored exactly, and only the *content* is derived.
- **Random generation and seeds.** C needs none. B needs a seeded RNG threaded through
  `randomTopology` and `CreateRandomPortalConnections`, a stored seed and a generator version.
- **Stable identity and endpoints.**
  - Identity is the zone name and the connection name (both verified unique). Display name
    equals identity; rename is explicitly unsupported (the UI has none).
  - If renames are added later, they must rewrite connection `from`/`to`/`guardZone`/
    `guardMatchGroup`, road `Connection` refs and `mandatoryContent` names.
  - A zone deleted and re-added under the same label: A/B treat it as a changed record; C
    stores `origin: added`, so it derives from the editor factory rather than from the
    generated zone that used to carry the name.
  - Parallel edges are distinct by their suffixed names. An endpoint change is delete plus
    add. Nameless connections are never stored.
  - Order: C stores the zone and connection order; A/B store it only when it differs from
    "base minus deleted, plus added".
- **Generator-version upgrades.** Under C, derived fields follow the new generator by design.
  Recording a generator version lets the editor warn "N derived fields changed since this
  file was saved". Under B, an upgrade invalidates every delta.
- **Default evolution.** Omitting default content rows or profile-derived fields means a
  default change rewrites the meaning of old documents. This is acceptable for C's derived
  fields (that is the point). Content rows are user rules, so they need a defaults version
  or must stay explicit.
- **Derived vs explicitly overridden.** C records an override only when the value differs
  from the derivation at save time. A user who types the default value back is "derived".
  Owner choice: accept this, or record explicit-override flags from the editor (more
  plumbing, exact intent).
- **Compatibility with old full snapshots.** v0/v1/v2 load as today. Converting to C at the
  first save is lossless at that instant (all 40 fixtures round-trip), but every field equal
  to its derivation becomes derived and starts following settings. The alternative is to mark
  every converted field as overridden: no semantic change, no size win until the user edits
  again.
- **Dirty tracking, cloning, Undo, Revert-to-Base.**
  - Dirty tracking: `SetManualEdits`' change detection would compare the compact form,
    which is expected (not measured) to be cheaper than today's entity conversion plus
    `DeepEqual`.
  - Cloning: the per-frame cost shrinks with the live representation; if the live model
    stays expanded, nothing changes until cloning is addressed separately.
  - Session Undo stays dialog-local and unchanged.
  - Revert-to-Base: under C, drop all overrides and edited structure; under A/B, drop the
    delta.
- **Setting changes.** Under C, derived fields follow and overrides persist; for example, a
  zone whose size the user set to 1.5 keeps it when `NeutralZoneSize` changes. The owner
  must confirm that precedence.
- **Final preview and export equivalence.** This was proven for every prototype and fixture
  by export-byte and preview-layout equality, rebuilt on the fixture's **stored** base
  template. That demonstrates zone and connection reconstruction. A real reload of a Random
  layout also regenerates template-level parts such as `orientation` (§5), so whole-template
  equality after reload is not claimed. Each implemented design must keep the two
  comparators as tests.

## 10. Illustrative JSON (option C under §2.6 sections; not a schema)

`manualPosition` is stored only once the user froze or moved the layout. The editor stamps
it on every zone at the first move, which is why `Spawn-A` carries one equal to its
generator position.

```json
{
  "schemaVersion": 3,
  "templateName": "Colosseum v4",
  "mapSize": 368,
  "manualEditSettings": {
    "generatorVersion": "v0.4.0",
    "zones": [
      { "name": "Spawn-A", "generatorPosition": { "X": 0.12, "Y": 0.5 }, "generatorRing": 0,
        "manualPosition": { "X": 0.12, "Y": 0.5 } },
      { "name": "Neutral-G", "quality": 2, "generatorPosition": { "X": 0.41, "Y": 0.63 },
        "generatorRing": 1, "manualPosition": { "X": 0.46, "Y": 0.63 },
        "fieldOverrides": { "size": 1.5, "guardMultiplier": 1.3 } },
      { "name": "Neutral-Z", "origin": "added", "quality": 2, "castles": 1,
        "manualPosition": { "X": 0.8, "Y": 0.2 }, "unset": ["mandatoryContent"] }
    ],
    "connections": [
      { "name": "GeoHub-A-S", "from": "Spawn-A", "to": "Neutral-S", "type": "Direct" },
      { "name": "Manual-G-Z", "from": "Neutral-G", "to": "Neutral-Z", "type": "Direct",
        "isUserAdded": true, "fieldOverrides": { "guardValue": 80000 }, "unset": ["road"] }
    ]
  }
}
```

## 11. §2.5 / §2.6 coordination points

1. **One migration.** §2.5's typed ban and override lists, §2.6's sections and any approved
   §2.4 representation ship as a single schema version (v3). The migrator chain becomes
   v0/v1 → v2 → v3, and v2 needs a **frozen** entity package (as `editor_state_v1` is today),
   because v3 changes the shape.
2. **Placement.** The manual section becomes §2.6's `manualEditSettings` object.
   `schemaVersion`, `templateName`, `gameMode`, `mapSize` and `experimentalMapSizes` stay at
   the root.
3. **Naming collision.** §2.5 introduces "value overrides" (SID/GuardValue). C's per-field
   data must use a distinct key (`fieldOverrides` above) and a distinct Go type name.
4. **Content rows and defaults.** Whether default content rows may be omitted is a
   ContentSettings question that sits beside §2.5's typed lists. Decide both in the same plan
   and with the same defaults-version policy.
5. **Wire-format pins.** [`editorStateWireFormat_integration_test.go`](../../test/integration/editorStateWireFormat_integration_test.go)
   and the entity JSON tests are rewritten once, for v3, with the legacy fixtures kept as
   inputs.
6. **Layering.** The derivation runs in services (generation, zone editor factories). The
   repository stays a typed decoder, and Model ⇄ Entity conversion stays at the existing
   seams (`file_service` today). The handler crosses DTO ⇄ Model only.

## 12. Recommendation

1. **Adopt C3 semantics.**
   - Store the structure; derive content from the current settings and the generator.
   - Store explicit overrides and `unset` lists.
   - Store roads as overrides (needed for Random).
   - Record a `generatorVersion`.
   - Decide migration as "derived when equal at first save", with a one-time notice.
2. **Fix §1.16 together with C3, or after it.** A stop-gap without provenance is lossy
   either way. If a quick fix is wanted first, the least surprising one is to warn and
   re-derive only the stale fields listed in §6 and accept that same-field edits are lost.
   That needs owner sign-off.
3. **Treat live cloning as a separate allocation item.** It needs no format change: the
   regeneration check does not need deep clones of the snapshot.
4. **Optional, independent quick wins:** compact JSON (−44% written bytes, at a readability
   cost) and default-row omission (−1.9 KB per file, needs a defaults version).
5. **Reject B** unless deterministic, seeded generation is wanted for its own sake. **Reject
   A** for persistence.

## 13. Acceptance tests for the chosen design

- **Round trip:** for every fixture class (the 11 cases × deterministic and Random layouts,
  plus a 30-zone hub layout), reconstruct a raw-equal snapshot, byte-equal export and equal
  preview. Port the probe's comparators, but derive connections from the **rebuilt** zones
  (the probe derived them from the original snapshot).
- **Determinism contract:** for every topology, two generations with the same settings have
  identical non-road zone content. For non-Random topologies without RandomPortals, roads,
  endpoints, names and positions must match as well. This test guards C against a future
  random field.
- **Settings precedence (§1.16 regression):** each stale setting in §6 changes the derived
  fields after an edit; explicit overrides survive.
- **Migration:** v0/v1/v2 fixtures, a large real-world-shaped v2 fixture and retired-topology
  rejections. Converted documents reproduce export and preview; unsupported newer versions
  are refused; a failed load leaves the document unchanged.
- **Identity:** delete then re-add the same label, parallel edges, user-added connection
  naming, preserved order.
- **Size:** the 6-player GeometricHub move fixture's manual section stays ≤ 30% of today's
  compact size. Absolute numbers are recorded, not gated.
- **Performance:** encode + decode of the large fixture is not slower than today's 1.1 ms
  median, and per-frame `AutoRegenerate` allocations do not grow (G2-style counting).
- **The verification an approved implementation must ship:**
  - dedicated model and delta tests;
  - migration tests;
  - wire-format integration coverage for v3;
  - manual-edit GUI regression flows: move, add, delete, quality, connection edit, apply,
    revert to base, save, load, settings change.

## 14. Limitations of this investigation

- Windows only, go1.27.2. Native Linux, Steam Deck and the race detector were unavailable,
  and no in-game result is claimed.
- The probe froze manual positions from `GeneratorPosition`; the dialog uses canvas
  geometry. The values differ, but the shape and field counts are identical.
- random8p's A/B base copies generator stamps and tier from the snapshot, so a tier edit is
  invisible to that delta. Every other base was regenerated and byte-verified.
- The owner's Colosseum files were measured for size only (no base template exists for them).
- The probe's C derivation of connections, and the zone count it passes to the tuning, read
  the original snapshot on both the encode and decode sides. Raw zone equality is asserted
  separately, so the round trip is not vacuous, but a real decoder only has the rebuilt zones.
- C/D exports were rebuilt on the stored base template, not on a fresh reload (§9).
- The default-content-rows size comes from a separate PowerShell measurement
  (`phase3-settings-rows.txt`).
- Determinism is shown for the tested settings (4–8 players, simple and advanced neutral
  counts); it is not proven for every combination.

## 15. Decisions requested

1. Representation: C3 (recommended), D only, A, or B.
2. Migration of old snapshots: derived when equal (size win, behaviour change) or
   all-overridden (no change, no win).
3. Override semantics: value-difference (as prototyped) or explicit editor flags.
4. §1.16: fix with C3, or a lossy stop-gap first.
5. `generatorVersion` and a "derived values changed" notice: yes or no.
6. Independent quick wins: compact JSON, default-row omission.
7. Live cloning: open a separate allocation item, or leave it.
