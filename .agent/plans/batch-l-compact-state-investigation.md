# Batch L: compact manual deltas and derivable editor state (review §2.4 / O11)

Investigate, measure and report whether manual zone/connection edits and derivable editor
state can be persisted and held live more compactly, then present alternatives for the
owner's decision gate. **No production code, schema, fixture or format change** is part of
this batch; all Go code is temporary and deleted before close.

Branch `AD/compact_state`, base `8cf9c80` (master, Batches H, I and K merged via PR #50).

## For Future Agents
As work proceeds: mark checkboxes `- [x]` as items complete; when a phase is done,
set its status to `Complete` and write its **Phase Summary** (what was done, key
decisions, anything needed to continue with zero context); run the phase's
**Verification Plan** and record the result before moving on. When all phases are
done, fill in **Final Recap** and **Deployment Plan**.

Read [AGENTS.md](../../AGENTS.md), [the handoff](../session-carry-forward.md) and review
[§2.4](../backlog/review-gpt-6-astra-09-07.md) first. Batches A–K are closed; never revisit
them. All raw artifacts live in gitignored `.agent/memories/batch-l/` (the owner's `air`
watcher can wipe `tmp/`).

## Owner decisions (2026-10-12)

| ID | Decision |
| --- | --- |
| L1 | Deliverable: the report [batch-l-compact-state.md](../investigations/batch-l-compact-state.md) (created in Phase 4); this plan links to it. |
| L2 | All measurement code is throwaway. It lives temporarily in `cmd/batchlprobe/` (a `main` package, so it may import `internal/`; `cmd/` is outside the coverage `-coverpkg`), source copies and raw output go to `.agent/memories/batch-l/`, and the folder is deleted with `Remove-Item` before close. Nothing in it is committed. |
| L3 | Delta encodings are prototyped (throwaway) and measured, not estimated. |
| L4 | Frozen fixture matrix (below), generated once because generation is unseeded, stored in `.agent/memories/batch-l/fixtures/`. |
| L5 | The owner's gitignored `output/research/Topology when loaded/Colosseum v3*.gen.json` may be read as a real-world sample; numbers only in the report, never copied into the repo. |
| L6 | The "snapshot wins over regenerated derived values" behaviour is an **unintended bug**. Verify it empirically, record it as a new review item §1.16 with a proposed separate fix batch; **no fix in L**. |
| L7 | Live state: measure per-frame `AutoRegenerate` cloning, `Clone`, `Equals`, `SetManualEdits`, save and load. Record only; no tuning. |
| L8 | Non-determinism is shown empirically (repeated generation per topology), throwaway. |
| L9 | Design depth: semantics plus illustrative JSON shapes, flagging every §2.5/§2.6 coordination point. No full schema. |
| L10 | The report ends with a reasoned recommendation; the owner decides. |
| L11 | Lighter gate: build, unit suite, coverage, testlayoutcheck, lint. Run before the first temporary file and again after deletion (proving the tree is back to baseline). |
| L12 | Reviews: plan by GPT-6.1 Sol; report by Claude Fable 5.1. |
| L13 | **Approved by the owner 2026-10-12.** While `cmd/batchlprobe/` exists it is a temporary, documented exception to AGENTS.md §4.4 (outside packages enter through handlers) and to the layering test, which scans `cmd/` and does not permit it to name DTOs or call services/mappers directly. The probe routes through handlers and `drivers.State` wherever possible; direct model/entity/mapper use is limited to byte measurement and the prototypes. No allow-list is edited. The unit suite is therefore run only before the probe exists (Phase 0) and after it is deleted (Phase 5); while it exists, only `go build ./...`, `go vet ./cmd/batchlprobe`, testlayoutcheck and lint are run, and an architecture-test failure attributed to the probe is expected and recorded. No `integration_test`, `gui` or `wireinject` tag and no test export is used. |

**Fixture matrix (L4).** Layouts: Square 2p small; GeometricHub 6p mapSize 368
(Colosseum-like); Random 8p large; Tournament 2p. Edits per layout: untouched; one zone
move; zone added; zone deleted; connection added; connection deleted; connection edited;
neutral quality change; explicit content/field override; nil/empty values (`Road` nil and
`false`, empty pools). Plus the committed v0/v1/v2 fixtures in
`test/test_helpers/testdata/` for old-snapshot compatibility.

## Facts established during scoping (code-verified at `8cf9c80`)

- The manual snapshot is the **whole** variant-0 layout: `drivers.State.ApplyEditedZones` →
  `EditorState.SetManualEdits` stores every zone and connection, and
  `reapplyManualEdits` → `templateHandler.UpdateTemplate` assigns them wholesale. There is
  no per-zone matching by name, index or ID; connections reference zones by `From`/`To`
  name strings.
- Persisted form: `editor_state.ManualZoneSave{Zone template_entity.Zone, GeneratorPosition,
  GeneratorRing, ManualPosition, Quality}` and `ManualConnectionSave{Connection, IsUserAdded}`.
  Schema v2; v0/v1 migrate through `EditorStateMigrator`.
- Sizes: committed fixtures (1 manual zone) ~5.6 KB; states without manual edits ~5.3 KB;
  `Colosseum v3 - after update.gen.json` (6p GeometricHub, 31 zones, 42 connections)
  **190,805 B**; one of its zones ~4.4 KB compact.
- Generation uses unseeded global `math/rand/v2` (`randomTopology.go`,
  `topologyConnectionService.go`); no seed is persisted.
- Non-layout-defining settings (`LayoutDefiningOptionsChanged` excludes, among others, zone
  sizes, guard randomization, resource/structure density, neutral stack and border guard
  strength) regenerate the template, but the stale snapshot then replaces the zones. Castle
  changes are reapplied to the snapshot (`ManualReapplyService.ApplyCastleSettingChanges`),
  and `UpdateTemplate` rebuilds template-level mandatory content and reconciles roads; no
  other derived zone/connection value is refreshed. This is L6; Phase 2 enumerates every
  setting.
- `editor.Window` calls `State.AutoRegenerate` every frame. It always clones the current
  state; previous and next are cloned only when present (`GetPreviousState`/`GetNextState`
  return nil otherwise). Idle after a generation is normally two full clones, a pending
  debounce three, each carrying the full snapshot.
- `drivers.NewUIState(handler, fileSystem, regeneration, false)` constructs the driver
  without a window; the exported composition injectors in `internal/composition` provide
  the production handlers.

## Phase 0: Baseline
Status: Complete

- [x] Confirm branch `AD/compact_state`; record `git status --short`. (Final baseline at owner HEAD `45876cf`, clean tree.)
- [x] `golangci-lint cache clean` (the owner's local binary is now named `golangci-lint`).
- [x] `go build ./...`.
- [x] `go test ./test/unit/... -count=1`.
- [x] Coverage: `go test '-coverpkg=./internal/...,./app/...' '-coverprofile=.agent/memories/batch-l/baseline/coverage.txt' ./test/unit/...` then `go tool cover '-func=…'`. Planned 73.2% (6441 / 8781); owner-adopted baseline **73.0% (5587 / 7646)** under go1.27.2.
- [x] `go run ./cmd/testlayoutcheck .` → `test-layout check passed`.
- [x] `golangci-lint run ./... --issues-exit-code=0`: **8 `exhaustive` issues** under v2.14.0 (see the summary).
- [x] Save all outputs under `.agent/memories/batch-l/baseline/`.

### Verification Plan
- Every command exits 0; coverage total equals 73.2% (6441 / 8781); lint 0. Any deviation
  stops the batch and is reported to the owner.

### Phase Summary
**Complete, with owner-adopted toolchain deviations (2026-10-12).**
- After the first run the owner committed `a9ae09c` (go 1.27.2, Gio 0.10.3, gofakeit 7.18.0,
  x/* bumps; `tools/go.mod` only its `go` line) and `45876cf` (this plan, the linter rename
  in `AGENTS.md`/`tasks.json`). The baseline was rerun on `45876cf`, clean tree.
- Build, unit, testlayoutcheck PASS.
- Coverage **73.0% (5587 / 7646)**. The same deduped count on Batch K's go1.27.0 profile
  reproduces 6915 / 9263, so the method is unchanged. There is no Go diff between
  `27b2933` and `8cf9c80`, so the denominator shift is go1.27.2's instrumentation (14
  functions change percentage, e.g. `createTopologyAdjacency` 70.0% → 57.1%). The owner
  adopted 73.0% as the Batch L baseline.
- Lint: the old `golangci-lint-v2` (built with go1.27.0) could not typecheck under 1.27.2
  (printed "0 issues", exited 7). The owner reinstalled it as `golangci-lint`, which reports
  **v2.14.0** built with go1.27.2, and exits 0 with **8 `exhaustive` findings** in
  pre-existing code (`bonusOptions.go`, `bonusPickerDialog.go`, `bonusesPanel.go`,
  `retiredTopologies.go`, `bonusEntryService.go`, `zoneLabelProvider.go`, and two unit
  tests). These 8 are the baseline; Phase 5 must reproduce exactly them. Not fixed in L.
- **Noted for the owner:** `tools/go.mod` still pins golangci-lint **v2.13.1** (CI's single
  source, Batch H), while the local binary is **v2.14.0**, so local and CI lint can differ.
- Outputs in `.agent/memories/batch-l/baseline/`.

## Phase 1: Probe, fixtures and current-cost measurements
Status: Complete

- [x] Create `cmd/batchlprobe/` (subcommands per phase) under L13. It uses the composition
  injectors, the handlers and a window-free `drivers.NewUIState(..., false)`; save goes
  through the state handler (the driver's save opens a file dialog). No production file is
  edited.
- [x] Build the L4 matrix: generate each layout once through `State.Generate`, apply each
  edit through `State.ApplyEditedZones`, and save the resulting state through the state
  handler into `.agent/memories/batch-l/fixtures/`.
- [x] Define and implement the comparators used by every later phase, each recording
  existing normalization separately so it is never attributed to a prototype:
  - **raw snapshot:** deep equality of manual zones/connections, nil vs empty preserved;
  - **persisted-semantic:** equality after the production save → load path (known
    normalizations: empty top-level snapshots, `omitempty` nested slices);
  - **validated state:** equality after load validation/migration (`fixIssues`);
  - **export:** byte equality of the final `.rmg.json` produced by the production mapper;
  - **preview:** equality of the production preview layout.
  Legacy checks use only the supported flat v0/v1/v2 fixtures; the retired-topology
  fixtures are rejection controls.
- [x] Persisted size per fixture: bytes as written; compact JSON of the same written bytes
  (`json.Compact`); gzip at a fixed level of the written bytes, decompression verified
  (reference only). Breakdown into settings / `manualZones` / `manualConnections`;
  per-zone field-group bytes (pools, guard reaction distribution, content limits,
  mandatory content, value fields, main objects, roads, biomes, position sidecars) and
  per-connection bytes.
- [x] Same breakdown for both Colosseum v3 samples (L5) and the supported v0/v1/v2 fixtures.
- [x] Live-state costs via `testing.Benchmark` (ns/op, B/op, allocs/op), with and without
  manual edits for each layout: `EditorState.Clone`; `ManualEditSettings.Equals` (equal and
  changed cases separately); `EqualsIgnoringManualEdits`; `SetManualEdits` (unchanged and
  changed inputs separately); encode/decode in memory, separate from filesystem save/load
  (the atomic writer syncs and may retry renames); and one `AutoRegenerate` call measured
  directly on the driver in each branch (idle after generation; pending debounce), with the
  clone count per branch recorded. Setup stays outside timing, results are retained to
  defeat elimination, and mutable inputs are reset per iteration.
- [x] Copy probe sources and raw outputs to `.agent/memories/batch-l/`.

### Verification Plan
- `go build ./...`, `go vet ./cmd/batchlprobe`, testlayoutcheck and lint run while the
  probe exists; any finding attributed to the probe is recorded under L13.
- Every matrix fixture loads back through the state handler with no rejection and passes
  the persisted-semantic comparator.
- Benchmark method: fixed toolchain (`go version` recorded), fixed `GOMAXPROCS`, fixed
  `-benchtime`-equivalent duration, 5 independent process runs; report median and
  min–max spread. Machine load is noted.

### Phase Summary
**Complete (2026-10-12).** Probe files: `main.go`, `fixtures.go`, `sizes.go`, `bench.go`,
`verify.go` (copies in `.agent/memories/batch-l/probe-src/`; `go build ./...` and
`go vet ./cmd/batchlprobe` exit 0). Subcommands: `fixtures`, `verify`, `sizes`, `bench <run>`,
`bench-summary`.
- **Matrix:** 44 fixtures (4 layouts × 11 edits: `untouched`, `applied-unchanged` and the 9
  L4 edits) in `fixtures/`, each with `.base.rmg.json` (its own fresh generation) and
  `.final.rmg.json`. Each edit runs on its own fresh generation (generation is unseeded).
  Every fixture: no apply error, persisted-semantic PASS, raw deep-equal after reload PASS,
  validated-state (`fixIssues` true vs false) equal with 0 warnings, export from the
  reloaded snapshot byte-equal (both in-process and via the stored base `.rmg.json`).
  `phase1-fixtures.txt`, `phase1-verify.txt`. Supported v0/v1/v2 fixtures load (1 zone, 1
  connection, migrate to v2); retired controls are rejected. **No normalization was
  observed** on any matrix fixture. The preview comparator is first exercised in Phase 3.
- **Approximation:** the probe freezes manual positions from `GeneratorPosition`; the
  dialog uses canvas geometry. The first fixture run had a probe bug (add-zone connected by
  label instead of zone name); it was discarded and regenerated.
- **Sizes** (`phase1-sizes.tsv`): no manual edits ≈ 5.3 KB (geohub6p 8.0 KB). One move:
  square2p 35.1 KB, tournament2p 33.1 KB, random8p 167.6 KB, geohub6p **202.4 KB**
  (Colosseum samples 187.5–190.8 KB). Compact JSON ≈ 56% of written; gzip-9 ≈ 3% (geohub6p
  move 6,480 B), i.e. the snapshot is extremely repetitive. Snapshot share of the file:
  > 95% for the large layouts. Largest zone groups (geohub6p move, compact value bytes):
  pools 28.1 KB, roads 15.9 KB, main objects 14.0 KB, content limits 12.8 KB, other 6.2 KB,
  values 5.3 KB, biomes 4.9 KB, position sidecars 3.6 KB, mandatory content 1.5 KB, guard
  reaction 1.4 KB; connections total 13.0 KB.
- **Live state** (`bench/phase1-bench-summary.tsv`; go1.27.2, GOMAXPROCS=4, 1 s
  benchtime, 5 process runs, medians): with geohub6p's snapshot, `Clone` 43 µs / 100 KB /
  980 allocs; `ManualEditSettings.Equals` 166–186 µs / 244–254 KB (it converts both sides to
  entities before `reflect.DeepEqual`); `SetManualEdits` 222–243 µs / 360–371 KB; encode
  442 µs; decode 552 µs. `AutoRegenerate` per frame: idle **80 µs / 201 KB / 1,960 allocs**
  (exactly 2 clones), pending **128 µs / 301 KB / 2,940 allocs** (3 clones); without manual
  edits 6.3 µs / 10.5 KB idle. `EqualsIgnoringManualEdits` is snapshot-independent
  (0.4–0.8 µs, 0 B). A first bench run seeded tournament drivers from a non-tournament state,
  so the mode transition dropped the snapshot; the driver now seeds settings first and
  asserts the snapshot survived, and all 5 runs were redone.
- **L13 while the probe exists:** testlayoutcheck PASS; lint = the 8 baseline `exhaustive`
  findings + 19 probe-only findings (forbidigo, funlen, gochecknoglobals, gocognit, goconst,
  golines, perfsprint, staticcheck, unused), `phase1-lint.txt`. Unit suite not run (L13).

## Phase 2: Determinism, identity and the snapshot-wins evidence
Status: Complete

- [x] Generate the same state 20 times for each of the 7 topologies and a 2p tournament;
  compare zone name sets, zone order, connection endpoint pairs, positions and per-zone
  value fields. Report which parts are stable per topology.
- [x] Identity: check zone-name and connection-name uniqueness and stability; document that
  `NextFreeZoneLabel` can reuse a deleted label and that `EnsureConnectionNames` names only
  nameless connections with order-dependent suffixes and does not repair duplicates. List
  every place that references a zone by name (connections, roads, main-object refs,
  mandatory content). Specify, for each delta option, identities for generated vs
  user-added zones, delete then re-add, rename (or rename explicitly unsupported), parallel
  edges, endpoint changes and reordered nameless edges; separate identity from display
  name; define reference rewriting and order preservation.
- [x] L6 settings inventory: classify **every** `EditorState` field as layout-defining,
  castle-diffed, zone/connection-affecting, template-only, inactive, or validation-only,
  including at least `MapSize`, `ExperimentalMapSizes`, zone sizes, guard randomization,
  densities, strengths, abandoned outposts, `MatchPlayerCastleFactions`, remote footholds,
  `TemplateName`, `GameMode`, hero counts, bans, value overrides, bonuses, the six content
  row slices, faction-law/astrology XP, lost-start rules, city hold, arena delay/count,
  tournament day/interval/points/save-army, and victory/tournament/arena changes that keep
  the effective modes. Inactive castle-count fields (advanced vs simple) are tested
  separately. `SchemaVersion` is persistence-only and excluded.
- [x] L6 evidence on the real driver: `NewUIState(..., false)` → `Generate` →
  `ApplyEditedZones` (one zone move) → `UpdateState` (change one setting) → `Generate`.
  For the paired comparison, capture and clone the fresh template generated in that same
  regeneration **before** reapplication (through a handler wrapper embedding
  `IGuiHandler`, as Batch G did), and diff it against the reapplied result: zones,
  connections and template-level content. Produce the exact list of stale fields per
  setting.
- [x] Record raw output in `.agent/memories/batch-l/`.

### Verification Plan
- The determinism run has two controls: a template vs its own `Clone` must report
  identical, and a template vs a copy with one field changed must report different.
- The L6 run has a castle control (player castle count) that must show the castle main
  objects actually changed in the reapplied result, not merely agree between outputs.

### Phase Summary
**Complete (2026-10-12).** Probe adds `phase2.go` (`determinism`, `tldiff`, `l6`).
- **Controls pass:** clone → all equal; one changed field → `zoneValues` 0/1 with exactly
  `guardCutoffValue`. L6: TemplateName no-op shows no difference; the castle control shows
  `mainObjects`/`roads` changed in R (R≠S0).
- **Determinism** (`phase2-determinism.tsv`, 20 runs each, compared with run 1): Circles,
  GeometricHub, Square, Geometric, Cross, Fractal at 4p, the 2p tournament (Random), and the
  square2p, geohub6p (Colosseum settings) and tournament2p fixture layouts are **byte-identical
  in every compared aspect (19/19)**. **Random** topology: zone set stable, but zone order,
  endpoint pairs, connection names, positions, roads and `variants[0].orientation` differ in
  almost every run (`phase2-template-level-diff.txt`). **RandomPortals** on Square: zones and
  positions stable, portal endpoint pairs/names and roads differ. The only `math/rand` call
  sites are `randomTopology.go` and `CreateRandomPortalConnections`'s shuffle/derangement.
  Determinism is empirical for these settings, not a contract (no seed, no generator version).
- **Identity:** in all 260 determinism generations zone names and connection names were unique and
  non-empty. References by name: connection `from`/`to`/`guardZone`; roads reference
  connections (`{"type":"Connection","args":[name]}`) and main objects; `mandatoryContent`
  uses label-derived names (`mandatory_content_neutral_G`); `guardMatchGroup` is
  label-derived; `contentCountLimits` lists player-pair names that depend only on the player
  count. `NextFreeZoneLabel` returns the first unused label, so a deleted zone's label is
  reused; `EnsureConnectionNames` names only nameless connections with order-dependent
  `-2`, `-3` suffixes and never repairs duplicates. Per-option identity rules are written in
  the report.
- **L6** (`phase2-l6.tsv`; Square/4p, GeometricHub/4p, Tournament/2p; one move each).
  Stale after a regeneration (fresh value differs, the result keeps the snapshot value):
  - `MapSize`: all six content value fields of every zone;
  - `ResourceDensityPercent`: `resourcesValue`, `resourcesValuePerArea`;
  - `StructureDensityPercent`: the four guarded/unguarded value fields;
  - `PlayerZoneSize`, `NeutralZoneSize`, `HubZoneSize`: `size`;
  - `NeutralStackStrengthPercent`: `guardMultiplier` and main objects;
  - `BorderGuardStrengthPercent`: main-object guards **and connection `guardValue`**;
  - `SpawnAbandonedOutposts`, `AbandonedOutpostCount`: main objects and roads (outposts never
    appear);
  - `CityHold` (GeometricHub): the hub zone's main objects, roads and biomes;
  - `HubZoneCastles`: castles/roads are reapplied, but the hub's three biome fields stay stale.
  Not stale: castle settings (reapplied), remote footholds (roads reconciled by
  `UpdateTemplate`), content rows (template-level mandatory content is rebuilt), XP,
  lost-start, city-hold days, game mode, hero counts, `ExperimentalMapSizes`,
  `MatchPlayerCastleFactions` and `GuardRandomization` (no zone effect observed), a
  mode-preserving victory change, inactive fields (simple-mode tier castles, arena count
  day, tournament interval outside tournament) and all four tournament fields in tournament
  mode. Bans, value overrides and bonuses feed only `gameRulesProvider` (template-level) by
  code reading. Layout-defining and effective-mode changes discard the snapshot (Batch E).
- **Side observation:** for Random topology the reapplied template takes `orientation` from
  the fresh generation while zones come from the snapshot.

## Phase 3: Delta prototypes and lossless-omission analysis
Status: Complete

- [x] Prototype, inside the probe, encoders/decoders for:
  - **A. Exact stored base + delta:** base persisted (or referenced) exactly; edits as
    added / changed-field / deleted records keyed by zone name and connection identity.
  - **B. Regenerated base + delta:** seed + generator version + delta. Not implementable
    today (unseeded); prototyped against the in-memory base to measure delta size only,
    with the determinism infrastructure it would need listed, not built.
  - **C. Structural graph + derived-field overrides:** the layout (zones, kinds/quality,
    positions, connections) stored compactly; every derivable field omitted when it equals
    its derivation from settings/quality/kind, stored as an explicit override otherwise.
  - **D. Representation-only:** today's semantics, dropping only fields proven lossless
    plus compact JSON. The cheapest comparison point.
- [x] Storage accounting is all-in: bytes include any stored base, metadata (generator
  version, seed) and references. For A, state where the base lives and whether the
  document stays portable without it. B is labelled **delta-only**, not a reconstructable
  persisted design today. For C, define the derivation inputs (settings, map size, zone
  count, quality, kind), the default-version policy, explicit null/empty/removal records,
  and the provenance limitation: today's snapshots cannot tell "explicitly set equal to the
  current default" from "never touched".
- [x] For each option and fixture: bytes, encode/decode ns and allocations (Phase 1
  method), and round trip through the Phase 1 comparators (raw snapshot, export bytes,
  preview layout).
- [x] Lossless-omission table for every candidate field (generated pools, guard reaction
  distribution, content limits, zone classification, default content rows, plus any field
  Phase 1 shows to be large): lossless today yes/no, conditions, what a manual edit can
  change, and how a future default/profile change would alter a saved document.
- [x] Record raw output in `.agent/memories/batch-l/`.

### Verification Plan
- Each prototype has a negative control: mutating one field of a reconstructed snapshot
  must make the equivalence check fail.
- Every fixture in the matrix round-trips for options A, C and D (B against the in-memory
  base); any failure is reported in the report as a limitation, not hidden.

### Phase Summary
**Complete (2026-10-12).** Probe adds `prototypes.go` (`prototypes`, `prototype-bench`);
output `phase3-prototypes.tsv`, `bench/phase3-prototype-bench-run{1..5}.tsv`.
- **Encodings:** zone and connection saves are flattened to top-level JSON keys
  (`zone.<key>`, sidecars, `connection.<key>`). A/B = per-record field delta keyed by name
  (`changed` set/unset, `added` full records, `deleted` names, `order` only when it differs).
  C = structure always stored (name, positions, ring, quality; connection name, endpoints,
  type, `isUserAdded`; a `castles` count for editor-added zones) plus overrides against a
  derivation (same-named zone of a **fresh generation with the current settings**, or an
  editor-made neutral zone for added or re-tiered zones; connections from
  `CreateZoneEditorConnection` + type change). C1 stores roads in the structure, C2 omits them
  (derived by `UpdateTemplate` reconciliation), C3 treats roads as an ordinary override.
  D = today's semantics; only pools, guard reaction, content limits and mandatory content
  may be omitted (all settings-independent); connections are stored whole. C/D are encoded
  against one fresh generation and decoded against an **independent second generation**.
- **Base for A/B:** regenerated and byte-compared with the stored base. Exact for every
  deterministic layout; random8p uses the stored `.rmg.json` with generator stamps and tier
  copied from the snapshot (a tier edit is then invisible to the delta; recorded).
- **Round trips (all 40 manual-edit fixtures):** A/B, C1, C3 and D pass raw, export and
  preview. **C2 fails raw and export everywhere**: `Reconcile` validates castle roads but
  never creates them, so roads are not derivable today.
- **Negative controls:** `GuardCutoffValue+1` → raw N, export N (preview Y, it does not use
  the field); a moved first zone → raw N, export Y (positions are not exported), preview N
  wherever all zones carry manual positions. Each comparator was shown able to fail.
- **Bytes (compact, manual section only; geohub6p move / random8p move / square2p move):**
  today 108,326 / 86,238 / 17,035; A/B delta 2,711 / 2,293 / 561 (29 B for applied-unchanged);
  A all-in (base + delta) 130,932 / 105,704 / 19,832; B all-in (delta + ~77 B metadata)
  2,788 / 2,370 / 638; C1 48,350 / 47,722 / 6,072; C3 23,297 / 48,058 / 2,457;
  D 99,638 / 88,096 / 13,834. Prototype keys carry a `zone.`/`connection.` prefix, so C/D are
  slightly pessimistic.
- **Derivation rates (C1 stats, 670 zone records):** every non-structure zone key matched
  its derivation in ≥ 99.4% of zones; overrides appear only where the edit changed the
  field (size, guard multiplier, pools on nil-empty, values on quality edits).
- **Encode+decode, geohub6p move, 5 runs:** today 1.11 ms median (1.00–1.24) / 785 KB /
  3,806 allocs; naive generic A/B prototype 2.81 ms (2.64–3.19) / 1.93 MB / 11,143 allocs.
- **Settings rows:** the default player content rows are 1,946 B of square2p's ~3.4 KB
  compact settings (14 rows).

## Phase 4: Report, review item and independent review
Status: Complete

- [x] Write [batch-l-compact-state.md](../investigations/batch-l-compact-state.md) covering
  everything §2.4 requires: live and persisted measurements; alternatives A–D with
  measured benefits, risks and migration limitations; exact-base vs regenerated-base
  semantics; seed/random requirements; stable identities and endpoint references;
  generator-version upgrades; default evolution; derived vs explicitly overridden values;
  compatibility with old full snapshots; dirty tracking, cloning, Undo/Revert-to-Base,
  setting changes; final preview/export equivalence; the lossless-omission table;
  reconstruction/round-trip and size/performance acceptance tests for the chosen design;
  the verification an approved implementation must ship (dedicated model/delta tests,
  migration tests, wire-format integration coverage, manual-edit GUI regression flows);
  illustrative JSON shapes; every §2.5/§2.6 coordination point, with one coordinated
  migration plan required; the L6 finding; and a recommendation (L10).
- [x] Review document (never renumber, §8 verbatim):
  - add **§1.16 🔴 High** for the L6 bug (evidence, why wrong, proposed fix scope, owner
    decision), placed after §1.15;
  - add a §9 row for a separate fix batch (next free letter, Q) for §1.16, and note its
    ordering relative to M;
  - add a Progress paragraph to §2.4 (investigation delivered, decision pending; §2.4 stays
    open, not FIXED);
  - update the progress line (24 fixed, **9 remaining: 1 High, 5 Medium, 3 Low**) and
    note §1.16 as an owner-added item;
  - keep the original audit total and add the current total beside it:
    **33 actionable, 9 High, 17 Medium, 7 Low** (owner-added §1.16 included).
- [x] Independent review of the report by Claude Fable 5.1 (L12); apply or explicitly
  decline every finding with a reason recorded here.

### Verification Plan
- Every number in the report traces to a raw file in `.agent/memories/batch-l/`.
- Every relative link in the report, plan and review edits resolves (throwaway check).
- `git diff` of the review shows §8 byte-identical and no renumbering.

### Phase Summary
**Complete (2026-10-12).** Report written; review edited (§1.16 🔴 added after §1.15, §2.4
Progress paragraph, §9 rows L updated and Q added, progress line 24 fixed / 9 remaining:
1 High, 5 Medium, 3 Low, current total 33: 9/17/7 beside the original 32). git diff shows
no hunk inside §8 and nothing renumbered. A throwaway link check passes for the report and
the plan. The review has two pre-existing broken links (webTopology.go and the deleted
Batch K plan), not touched here.
- **Independent review (Claude Fable 5.1): APPROVE WITH CHANGES, 13 findings, all applied.**
  (1) Five live-state numbers re-copied from phase1-bench-summary.tsv (they came from a
  discarded run). (2) GuardRandomization is inactive in simple mode: an advanced-mode L6
  base was added, which found guardRandomization stale; report §6 and review §1.16
  updated. (3) The live-state half of C stated (§7). (4) An R≠F column was added; §6 now
  separates "reapplied" from "no zone effect"; castle and foothold roads differ from a fresh
  generation. (5) The C decoder reading the snapshot is recorded in §14, and §13 requires
  deriving from rebuilt zones. (6) The export comparator's stored-base limitation is in §9
  and §14. (7) B's round trip is 30/30 regenerated; random8p uses the stored base only.
  (8) The settings-row size is now in raw phase3-settings-rows.txt. (9) Written-size
  labels fixed (202 → 113 KB file). (10) "C3 ≈ C1 (+0.7%)". (11) The dirty-tracking claim
  is marked "expected". (12) The manual-position sentence was added to §10. (13) §1.16's
  "Fix scope" heading aligned to "Fix".

## Phase 5: Cleanup and end-state gate
Status: Complete

- [x] Confirm the probe sources are copied to `.agent/memories/batch-l/probe-src/`, then
  `Remove-Item -Recurse cmd/batchlprobe`.
- [x] Rerun the Phase 0 gate (L11); coverage, lint and testlayoutcheck must equal the
  baseline exactly.
- [x] `git status --short` lists only docs: this plan, the report, the review, the handoff.
- [x] Update [the handoff](../session-carry-forward.md).

### Verification Plan
- Gate results identical to Phase 0 (coverage 73.0%, 5587 / 7646; lint the same 8
  `exhaustive` findings; testlayoutcheck passed); `cmd/batchlprobe` absent; no `.go` file in
  `git status`.

### Phase Summary
**Complete (2026-10-12).** Probe sources verified byte-identical in
`.agent/memories/batch-l/probe-src/` (7 files), then `cmd/batchlprobe` and the built binary
were removed with `Remove-Item`. End-state gate (`.agent/memories/batch-l/endstate/`):
- build PASS; unit PASS;
- coverage **73.0% (5587 / 7646)**, identical to Phase 0;
- testlayoutcheck PASS;
- lint exit 0, with the **same 8 `exhaustive` findings** (line-for-line compare).

`git status --short`: the review and this plan modified, the report untracked, and the
handoff updated. No `.go` file appears.

## Plan review (GPT-6.1 Sol, 2026-10-12)

APPROVE WITH CHANGES, 10 findings, all applied: L6 runs on the real driver with a paired
same-generation comparison (1); full settings inventory (2); L13 architecture exception
documented (3); explicit comparators and normalization (4); identity lifecycle cases (5);
all-in accounting, B delta-only, C provenance (6); benchmark method and a known-different
control (7); scoping facts corrected, `AutoRegenerate` measured directly (8); implementation
verification checklist carried into the report (9); total finding counts (10).

**Owner approved the plan on 2026-10-12, including L13.**

## Out of scope
Any representation, schema, migration, seed or determinism change; fixing §1.16; tuning
live-state cloning; §2.5/§2.6 implementation; touching `data/`, the template schema or the
registry; the output path. Batches A–K settled behaviour stays as is.

## Final Recap
Batch L delivered the §2.4 feasibility investigation without touching production code.
- **Today's format:** one manual edit persists the whole layout (6-player hub: 202 KB vs
  8 KB) and costs 80 µs / 201 KB on every idle frame.
- **Determinism:** generation is deterministic for 6 of 7 topologies and for tournaments;
  only Random and RandomPortals vary.
- **Alternatives:** four were prototyped and round-trip export and preview on all 40 edit
  fixtures. The exception is derived roads (C2), which fails because reconciliation never
  re-creates castle roads.
- **Recommendation:** C3 (structure + derived-field overrides), −78% / −44%, with no seed and
  §1.16 fixed for derived fields.
- **New bug:** the investigation found §1.16, recorded in the review with a batch Q row.
- **Reviews:** the plan was reviewed by GPT-6.1 Sol and the report by Claude Fable 5.1; all
  23 findings were applied.
- **Baseline:** owner-adopted 73.0% (5587 / 7646) under go1.27.2, plus 8 pre-existing
  `exhaustive` findings with the local `golangci-lint` v2.14.0. Unchanged at the end.

## Deployment Plan
Nothing to deploy: documentation only, and no code, schema, Wire or output-path change. The
owner reviews and commits the report, this plan, the review and the handoff, then answers the
report's §15 decisions. Any implementation starts as a new planned batch (M and/or Q).
