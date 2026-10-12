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
Status: In progress

- [x] Confirm branch `AD/compact_state`, HEAD `8cf9c80` plus only this plan; record `git status --short`.
- [x] `golangci-lint-v2 cache clean`.
- [x] `go build ./...`.
- [x] `go test ./test/unit/... -count=1`.
- [x] Coverage: `go test '-coverpkg=./internal/...,./app/...' '-coverprofile=.agent/memories/batch-l/baseline/coverage.txt' ./test/unit/...` then `go tool cover '-func=…'`; expect **73.2% (6441 / 8781)**. (Got 73.0%, 5587 / 7646; see the summary.)
- [x] `go run ./cmd/testlayoutcheck .` → `test-layout check passed`.
- [ ] `golangci-lint-v2 run ./... --issues-exit-code=0` → 0 issues. (Typecheck failure; see the summary.)
- [x] Save all outputs under `.agent/memories/batch-l/baseline/`.

### Verification Plan
- Every command exits 0; coverage total equals 73.2% (6441 / 8781); lint 0. Any deviation
  stops the batch and is reported to the owner.

### Phase Summary
**Blocked on the owner (2026-10-12): toolchain changed since Batch K.** The machine now runs
**go1.27.2** (Batch K used go1.27.0); there is no Go diff between `27b2933` and `8cf9c80`.
- Build PASS; unit PASS; testlayoutcheck PASS.
- Coverage **73.0% (5587 / 7646)** by `go tool cover` and the deduped block count. The same
  dedupe on Batch K's go1.27.0 profile reproduces 6915 / 9263, so the method is unchanged; the
  denominator moved because go1.27.2 instruments blocks differently (14 functions change
  percentage with no code change, e.g. `createTopologyAdjacency` 70.0% → 57.1%).
- Lint prints "0 issues" but exits 7 with a **typechecking error**: the installed
  golangci-lint-v2 cannot read go1.27.2 export data ("export data version 5 is greater than
  maximum supported version 4"). The lint baseline is therefore **invalid**, not clean.
- Outputs in `.agent/memories/batch-l/baseline/`.

## Phase 1: Probe, fixtures and current-cost measurements
Status: Not started

- [ ] Create `cmd/batchlprobe/` (subcommands per phase) under L13. It uses the composition
  injectors, the handlers and a window-free `drivers.NewUIState(..., false)`; save goes
  through the state handler (the driver's save opens a file dialog). No production file is
  edited.
- [ ] Build the L4 matrix: generate each layout once through `State.Generate`, apply each
  edit through `State.ApplyEditedZones`, and save the resulting state through the state
  handler into `.agent/memories/batch-l/fixtures/`.
- [ ] Define and implement the comparators used by every later phase, each recording
  existing normalization separately so it is never attributed to a prototype:
  - **raw snapshot:** deep equality of manual zones/connections, nil vs empty preserved;
  - **persisted-semantic:** equality after the production save → load path (known
    normalizations: empty top-level snapshots, `omitempty` nested slices);
  - **validated state:** equality after load validation/migration (`fixIssues`);
  - **export:** byte equality of the final `.rmg.json` produced by the production mapper;
  - **preview:** equality of the production preview layout.
  Legacy checks use only the supported flat v0/v1/v2 fixtures; the retired-topology
  fixtures are rejection controls.
- [ ] Persisted size per fixture: bytes as written; compact JSON of the same written bytes
  (`json.Compact`); gzip at a fixed level of the written bytes, decompression verified
  (reference only). Breakdown into settings / `manualZones` / `manualConnections`;
  per-zone field-group bytes (pools, guard reaction distribution, content limits,
  mandatory content, value fields, main objects, roads, biomes, position sidecars) and
  per-connection bytes.
- [ ] Same breakdown for both Colosseum v3 samples (L5) and the supported v0/v1/v2 fixtures.
- [ ] Live-state costs via `testing.Benchmark` (ns/op, B/op, allocs/op), with and without
  manual edits for each layout: `EditorState.Clone`; `ManualEditSettings.Equals` (equal and
  changed cases separately); `EqualsIgnoringManualEdits`; `SetManualEdits` (unchanged and
  changed inputs separately); encode/decode in memory, separate from filesystem save/load
  (the atomic writer syncs and may retry renames); and one `AutoRegenerate` call measured
  directly on the driver in each branch (idle after generation; pending debounce), with the
  clone count per branch recorded. Setup stays outside timing, results are retained to
  defeat elimination, and mutable inputs are reset per iteration.
- [ ] Copy probe sources and raw outputs to `.agent/memories/batch-l/`.

### Verification Plan
- `go build ./...`, `go vet ./cmd/batchlprobe`, testlayoutcheck and lint run while the
  probe exists; any finding attributed to the probe is recorded under L13.
- Every matrix fixture loads back through the state handler with no rejection and passes
  the persisted-semantic comparator.
- Benchmark method: fixed toolchain (`go version` recorded), fixed `GOMAXPROCS`, fixed
  `-benchtime`-equivalent duration, 5 independent process runs; report median and
  min–max spread. Machine load is noted.

### Phase Summary
_(write when phase completes)_

## Phase 2: Determinism, identity and the snapshot-wins evidence
Status: Not started

- [ ] Generate the same state 20 times for each of the 7 topologies and a 2p tournament;
  compare zone name sets, zone order, connection endpoint pairs, positions and per-zone
  value fields. Report which parts are stable per topology.
- [ ] Identity: check zone-name and connection-name uniqueness and stability; document that
  `NextFreeZoneLabel` can reuse a deleted label and that `EnsureConnectionNames` names only
  nameless connections with order-dependent suffixes and does not repair duplicates. List
  every place that references a zone by name (connections, roads, main-object refs,
  mandatory content). Specify, for each delta option, identities for generated vs
  user-added zones, delete then re-add, rename (or rename explicitly unsupported), parallel
  edges, endpoint changes and reordered nameless edges; separate identity from display
  name; define reference rewriting and order preservation.
- [ ] L6 settings inventory: classify **every** `EditorState` field as layout-defining,
  castle-diffed, zone/connection-affecting, template-only, inactive, or validation-only,
  including at least `MapSize`, `ExperimentalMapSizes`, zone sizes, guard randomization,
  densities, strengths, abandoned outposts, `MatchPlayerCastleFactions`, remote footholds,
  `TemplateName`, `GameMode`, hero counts, bans, value overrides, bonuses, the six content
  row slices, faction-law/astrology XP, lost-start rules, city hold, arena delay/count,
  tournament day/interval/points/save-army, and victory/tournament/arena changes that keep
  the effective modes. Inactive castle-count fields (advanced vs simple) are tested
  separately. `SchemaVersion` is persistence-only and excluded.
- [ ] L6 evidence on the real driver: `NewUIState(..., false)` → `Generate` →
  `ApplyEditedZones` (one zone move) → `UpdateState` (change one setting) → `Generate`.
  For the paired comparison, capture and clone the fresh template generated in that same
  regeneration **before** reapplication (through a handler wrapper embedding
  `IGuiHandler`, as Batch G did), and diff it against the reapplied result: zones,
  connections and template-level content. Produce the exact list of stale fields per
  setting.
- [ ] Record raw output in `.agent/memories/batch-l/`.

### Verification Plan
- The determinism run has two controls: a template vs its own `Clone` must report
  identical, and a template vs a copy with one field changed must report different.
- The L6 run has a castle control (player castle count) that must show the castle main
  objects actually changed in the reapplied result, not merely agree between outputs.

### Phase Summary
_(write when phase completes)_

## Phase 3: Delta prototypes and lossless-omission analysis
Status: Not started

- [ ] Prototype, inside the probe, encoders/decoders for:
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
- [ ] Storage accounting is all-in: bytes include any stored base, metadata (generator
  version, seed) and references. For A, state where the base lives and whether the
  document stays portable without it. B is labelled **delta-only**, not a reconstructable
  persisted design today. For C, define the derivation inputs (settings, map size, zone
  count, quality, kind), the default-version policy, explicit null/empty/removal records,
  and the provenance limitation: today's snapshots cannot tell "explicitly set equal to the
  current default" from "never touched".
- [ ] For each option and fixture: bytes, encode/decode ns and allocations (Phase 1
  method), and round trip through the Phase 1 comparators (raw snapshot, export bytes,
  preview layout).
- [ ] Lossless-omission table for every candidate field (generated pools, guard reaction
  distribution, content limits, zone classification, default content rows, plus any field
  Phase 1 shows to be large): lossless today yes/no, conditions, what a manual edit can
  change, and how a future default/profile change would alter a saved document.
- [ ] Record raw output in `.agent/memories/batch-l/`.

### Verification Plan
- Each prototype has a negative control: mutating one field of a reconstructed snapshot
  must make the equivalence check fail.
- Every fixture in the matrix round-trips for options A, C and D (B against the in-memory
  base); any failure is reported in the report as a limitation, not hidden.

### Phase Summary
_(write when phase completes)_

## Phase 4: Report, review item and independent review
Status: Not started

- [ ] Write [batch-l-compact-state.md](../investigations/batch-l-compact-state.md) covering
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
- [ ] Review document (never renumber, §8 verbatim):
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
- [ ] Independent review of the report by Claude Fable 5.1 (L12); apply or explicitly
  decline every finding with a reason recorded here.

### Verification Plan
- Every number in the report traces to a raw file in `.agent/memories/batch-l/`.
- Every relative link in the report, plan and review edits resolves (throwaway check).
- `git diff` of the review shows §8 byte-identical and no renumbering.

### Phase Summary
_(write when phase completes)_

## Phase 5: Cleanup and end-state gate
Status: Not started

- [ ] Confirm the probe sources are copied to `.agent/memories/batch-l/probe-src/`, then
  `Remove-Item -Recurse cmd/batchlprobe`.
- [ ] Rerun the Phase 0 gate (L11); coverage, lint and testlayoutcheck must equal the
  baseline exactly.
- [ ] `git status --short` lists only docs: this plan, the report, the review, the handoff.
- [ ] Update [the handoff](../session-carry-forward.md).

### Verification Plan
- Gate results identical to Phase 0 (coverage 73.2%, 6441 / 8781; lint 0; testlayoutcheck
  passed); `cmd/batchlprobe` absent; no `.go` file in `git status`.

### Phase Summary
_(write when phase completes)_

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
_(write when all phases complete)_

## Deployment Plan
_(write when all phases complete)_
